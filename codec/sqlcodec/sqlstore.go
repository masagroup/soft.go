package sqlcodec

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/chebyrash/promise"
	"github.com/masagroup/soft.go/ecore"
	"github.com/panjf2000/ants/v2"
	"github.com/petermattis/goid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

// queries of the store, built once : they are shared by the writer and the pool connections

type sqlSingleQueries struct {
	column      *sqlColumn
	updateQuery string
	selectQuery string
	removeQuery string
}

func newSqlSingleQueries(column *sqlColumn) *sqlSingleQueries {
	table := column.table
	tableName := sqlEscapeIdentifier(table.name)
	columnName := sqlEscapeIdentifier(column.columnName)
	// a table without key (contents) is only removed from
	keyName := ""
	if table.key != nil {
		keyName = table.keyName()
	}
	return &sqlSingleQueries{
		column:      column,
		updateQuery: "UPDATE " + tableName + " SET " + columnName + "=? WHERE " + keyName + "=?",
		selectQuery: "SELECT " + columnName + " FROM " + tableName + " WHERE " + keyName + "=?",
		removeQuery: "DELETE FROM " + tableName + " WHERE " + columnName + "=?",
	}
}

func (ss *sqlSingleQueries) getUpdateQuery() string {
	return ss.updateQuery
}

func (ss *sqlSingleQueries) getSelectQuery() string {
	return ss.selectQuery
}

func (ss *sqlSingleQueries) getRemoveQuery() string {
	return ss.removeQuery
}

type sqlManyQueries struct {
	table            *sqlTable
	updateValueQuery string
	updateIdxQuery   string
	selectOneQuery   string
	selectAllQuery   string
	existsQuery      string
	clearQuery       string
	countQuery       string
	containsQuery    string
	indexOfQuery     string
	lastIndexOfQuery string
	idxToListIndex   string
	listIndexToIdx   string
	removeQuery      string
	insertQuery      string
}

func newSqlManyQueries(table *sqlTable) *sqlManyQueries {
	tableName := sqlEscapeIdentifier(table.name)
	keyName := table.keyName()
	valueName := sqlEscapeIdentifier(table.columns[len(table.columns)-1].columnName)
	// the rowid of the element at the list index given as parameter
	rowAtIndex := "(SELECT rowid FROM " + tableName + " WHERE " + keyName + "=? ORDER BY " + keyName + " ASC, idx ASC LIMIT 1 OFFSET ?)"
	return &sqlManyQueries{
		table:            table,
		insertQuery:      table.insertQuery(),
		updateValueQuery: "UPDATE " + tableName + " SET " + valueName + "=? WHERE rowid IN " + rowAtIndex,
		updateIdxQuery:   "UPDATE " + tableName + " SET idx=? WHERE rowid IN " + rowAtIndex + " RETURNING " + valueName,
		removeQuery:      "DELETE FROM " + tableName + " WHERE rowid IN " + rowAtIndex + " RETURNING " + valueName,
		selectOneQuery:   "SELECT " + valueName + " FROM " + tableName + " WHERE " + keyName + "=? ORDER BY " + keyName + " ASC, idx ASC LIMIT 1 OFFSET ?",
		selectAllQuery:   "SELECT " + valueName + " FROM " + tableName + " WHERE " + keyName + "=? ORDER BY " + keyName + " ASC, idx ASC",
		existsQuery:      "SELECT EXISTS(SELECT 1 FROM " + tableName + " WHERE " + keyName + "=?)",
		clearQuery:       "DELETE FROM " + tableName + " WHERE " + keyName + "=?",
		countQuery:       "SELECT COUNT(*) FROM " + tableName + " WHERE " + keyName + "=?",
		containsQuery:    "SELECT rowid FROM " + tableName + " WHERE " + keyName + "=? AND " + valueName + "=?",
		indexOfQuery:     "SELECT idx FROM " + tableName + " WHERE " + keyName + "=? AND " + valueName + "=? ORDER BY idx ASC LIMIT 1",
		lastIndexOfQuery: "SELECT idx FROM " + tableName + " WHERE " + keyName + "=? AND " + valueName + "=? ORDER BY idx DESC LIMIT 1",
		idxToListIndex:   "SELECT COUNT(*) FROM " + tableName + " WHERE " + keyName + "=? AND idx<?",
		listIndexToIdx:   "SELECT idx FROM " + tableName + " WHERE " + keyName + "=? ORDER BY idx ASC LIMIT ? OFFSET ?",
	}
}

func (ss *sqlManyQueries) getInsertQuery() string {
	return ss.insertQuery
}

func (ss *sqlManyQueries) getUpdateValueQuery() string {
	return ss.updateValueQuery
}

func (ss *sqlManyQueries) getUpdateIdxQuery() string {
	return ss.updateIdxQuery
}

func (ss *sqlManyQueries) getSelectOneQuery() string {
	return ss.selectOneQuery
}

func (ss *sqlManyQueries) getSelectAllQuery() string {
	return ss.selectAllQuery
}

func (ss *sqlManyQueries) getExistsQuery() string {
	return ss.existsQuery
}

func (ss *sqlManyQueries) getClearQuery() string {
	return ss.clearQuery
}

func (ss *sqlManyQueries) getCountQuery() string {
	return ss.countQuery
}

func (ss *sqlManyQueries) getContainsQuery() string {
	return ss.containsQuery
}

func (ss *sqlManyQueries) getIndexOfQuery() string {
	return ss.indexOfQuery
}

func (ss *sqlManyQueries) getLastIndexOfQuery() string {
	return ss.lastIndexOfQuery
}

func (ss *sqlManyQueries) getIdxToListIndexQuery() string {
	return ss.idxToListIndex
}

func (ss *sqlManyQueries) getListIndexToIdxQuery() string {
	return ss.listIndexToIdx
}

func (ss *sqlManyQueries) getRemoveQuery() string {
	return ss.removeQuery
}

type sqlStoreObjectManager struct {
	store *SQLStore
}

func newSQLStoreObjectManager() *sqlStoreObjectManager {
	return &sqlStoreObjectManager{}
}

func (r *sqlStoreObjectManager) registerObject(o ecore.EObject, locked bool) {
	// an object encoded by the open transaction is read by the writer until the commit
	r.store.registerDirty(o)
	// set store object
	if lockedObject, _ := o.(ecore.ELockedStoreProvider); locked && lockedObject != nil {
		lockedObject.SetEStoreLocked(r.store)
	} else if storeObject, _ := o.(ecore.EStoreProvider); storeObject != nil {
		storeObject.SetEStore(r.store)
	}
}

type SQLStoreIDManager interface {
	SQLDecoderIDManager
	SQLEncoderIDManager
	ClearObjectID(eObject ecore.EObject)
}

type sqlStoreIDManagerImpl struct {
	sqlDecoderIDManagerImpl
	sqlEncoderIDManagerImpl
	mutex sync.Mutex
}

func newSQLStoreIDManager() SQLStoreIDManager {
	return &sqlStoreIDManagerImpl{
		sqlDecoderIDManagerImpl: sqlDecoderIDManagerImpl{
			packages:     map[int64]ecore.EPackage{},
			objects:      map[int64]ecore.EObject{},
			classes:      map[int64]ecore.EClass{},
			enumLiterals: map[int64]ecore.EEnumLiteral{},
		},
		sqlEncoderIDManagerImpl: sqlEncoderIDManagerImpl{
			packages:     map[ecore.EPackage]int64{},
			objects:      map[ecore.EObject]int64{},
			classes:      map[ecore.EClass]int64{},
			enumLiterals: map[ecore.EEnumLiteral]int64{},
		},
	}
}

func (r *sqlStoreIDManagerImpl) SetPackageID(p ecore.EPackage, id int64) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.sqlDecoderIDManagerImpl.packages[id] = p
	r.sqlEncoderIDManagerImpl.packages[p] = id
}

func (r *sqlStoreIDManagerImpl) GetPackageID(p ecore.EPackage) (int64, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.sqlEncoderIDManagerImpl.GetPackageID(p)
}

func (r *sqlStoreIDManagerImpl) GetPackageFromID(id int64) (ecore.EPackage, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.sqlDecoderIDManagerImpl.GetPackageFromID(id)
}

func (r *sqlStoreIDManagerImpl) SetClassID(c ecore.EClass, id int64) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.sqlDecoderIDManagerImpl.classes[id] = c
	r.sqlEncoderIDManagerImpl.classes[c] = id
}

func (r *sqlStoreIDManagerImpl) GetClassID(c ecore.EClass) (int64, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.sqlEncoderIDManagerImpl.GetClassID(c)
}

func (r *sqlStoreIDManagerImpl) GetClassFromID(id int64) (ecore.EClass, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.sqlDecoderIDManagerImpl.GetClassFromID(id)
}

func (r *sqlStoreIDManagerImpl) GetObjectFromID(id int64) (o ecore.EObject, b bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.sqlDecoderIDManagerImpl.GetObjectFromID(id)
}

func (r *sqlStoreIDManagerImpl) GetObjectID(o ecore.EObject) (id int64, b bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if sqlObject, _ := o.(SQLObject); sqlObject != nil {
		// sql object with an id
		id = sqlObject.GetSQLID()
		// check if registered
		_, b = r.sqlDecoderIDManagerImpl.objects[id]
	} else {
		id, b = r.sqlEncoderIDManagerImpl.GetObjectID(o)
	}
	return
}

func (r *sqlStoreIDManagerImpl) SetObjectID(o ecore.EObject, id int64) {
	r.mutex.Lock()
	r.sqlDecoderIDManagerImpl.objects[id] = o
	if sqlObject, _ := o.(SQLObject); sqlObject != nil {
		// set sql id if created object is an sql object
		sqlObject.SetSQLID(id)
	} else {
		// otherwse initialize map
		r.sqlEncoderIDManagerImpl.objects[o] = id
	}
	r.mutex.Unlock()
}

func (r *sqlStoreIDManagerImpl) ClearObjectID(o ecore.EObject) {
	r.mutex.Lock()
	var id int64
	if sqlObject, _ := o.(SQLObject); sqlObject != nil {
		id = sqlObject.GetSQLID()
	} else {
		id = r.sqlEncoderIDManagerImpl.objects[o]
		delete(r.sqlEncoderIDManagerImpl.objects, o)
	}
	delete(r.sqlDecoderIDManagerImpl.objects, id)
	r.mutex.Unlock()
}

func (r *sqlStoreIDManagerImpl) SetEnumLiteralID(e ecore.EEnumLiteral, id int64) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.sqlDecoderIDManagerImpl.enumLiterals[id] = e
	r.sqlEncoderIDManagerImpl.enumLiterals[e] = id
}

func (r *sqlStoreIDManagerImpl) GetEnumLiteralID(e ecore.EEnumLiteral) (int64, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.sqlEncoderIDManagerImpl.GetEnumLiteralID(e)
}

func (r *sqlStoreIDManagerImpl) GetEnumLiteralFromID(id int64) (ecore.EEnumLiteral, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.sqlDecoderIDManagerImpl.GetEnumLiteralFromID(id)
}

var operationID atomic.Int64

type operationType uint8

func (ot operationType) String() string {
	switch ot {
	case operationRead:
		return "read"
	case operationWrite:
		return "write"
	}
	return ""
}

const (
	operationRead  operationType = 1 << 0
	operationWrite operationType = 1 << 1
)

// operation is an operation of the store : a write, executed by the writer goroutine
// inside the current transaction, or a read, executed by the writer when it depends
// on a pending write and on a pool connection otherwise. An operation without object
// is on the whole store (roots, arbitrary queries, serialization)
type operation struct {
	id           int64
	type_        operationType
	cmd          string
	object       ecore.EObject
	feature      ecore.EStructuralFeature
	index        int
	value        any
	fn           func() (any, error)
	abandoned    atomic.Bool     // the caller stopped waiting : a read not executed yet is skipped
	ctx          context.Context // context of the caller, nil for background
	commitBefore bool            // commit the transaction before the operation
	flush        bool            // commit the transaction after the operation
	done         chan struct{}   // closed when the operation is executed, nil when nobody waits
	result       any
	err          error
	enqueued     time.Time
}

type operationOption func(*operation)

// withContext gives the operation the context of its caller : it is not executed once the
// context is done
func withContext(ctx context.Context) operationOption {
	return func(op *operation) { op.ctx = ctx }
}

// withCommitBefore commits the transaction before the operation
func withCommitBefore() operationOption {
	return func(op *operation) { op.commitBefore = true }
}

// withFlush commits the transaction after the operation
func withFlush() operationOption {
	return func(op *operation) { op.flush = true }
}

// withWait makes the caller wait for the result of the operation : without it, a write is
// executed in the background and its error reported by the next Flush
func withWait() operationOption {
	return func(op *operation) { op.done = make(chan struct{}) }
}

func newOperation(
	cmd string,
	type_ operationType,
	object ecore.EObject,
	feature ecore.EStructuralFeature,
	index int,
	value any,
	fn func() (any, error),
	options ...operationOption) *operation {
	op := &operation{
		id:      operationID.Add(1),
		type_:   type_,
		cmd:     cmd,
		object:  object,
		feature: feature,
		index:   index,
		value:   value,
		fn:      fn,
	}
	for _, option := range options {
		option(op)
	}
	return op
}

type operationMarshaler struct {
	op *operation
}

func (m *operationMarshaler) MarshalLogObject(e zapcore.ObjectEncoder) error {
	op := m.op
	var objectString string
	if op.object != nil {
		objectString = fmt.Sprintf("%s(%p)", op.object.EClass().GetName(), op.object)
	}
	var featureString string
	if op.feature != nil {
		featureString = op.feature.GetName()
	}
	e.AddInt64("id", op.id)
	e.AddString("op", op.cmd)
	e.AddString("type", op.type_.String())
	e.AddString("object", objectString)
	e.AddString("feature", featureString)
	e.AddInt("index", op.index)
	switch v := op.value.(type) {
	case ecore.EObject:
		e.AddString("value", fmt.Sprintf("%s(%p)", v.EClass().GetName(), v))
	default:
		if err := e.AddReflected("value", op.value); err != nil {
			return err
		}
	}
	return nil
}

type objectFeature struct {
	object  ecore.EObject
	feature ecore.EStructuralFeature
}

type SQLStore struct {
	*sqlBase
	sqlDecoder
	sqlEncoder
	isClosed         atomic.Bool
	sqlIDManager     SQLStoreIDManager
	singleQueries    map[*sqlColumn]*sqlSingleQueries
	manyQueries      map[*sqlTable]*sqlManyQueries
	mutexQueries     sync.Mutex
	loggerOperations *zap.Logger
	timeoutOperation time.Duration
	maxAllocSize     int64

	// writer : one goroutine executes the write operations in order on its own connection,
	// inside one transaction committed when the queue is empty or past commitMaxOperations
	// or commitMaxDuration (group commit)
	writerConn          *sqlite.Conn
	writerConnClose     func(conn *sqlite.Conn) error
	writerGoid          atomic.Int64
	writerSignal        chan struct{}
	writerDone          chan struct{}
	commitMaxOperations int
	commitMaxDuration   time.Duration
	concurrentReads     bool // reads of objects untouched by the transaction run on pool connections
	txOpen              bool
	txStart             time.Time
	txOperations        int
	mutexWriter         sync.Mutex // guards the fields below
	queue               []*operation
	busy                bool                         // the writer is executing an operation
	closed              bool                         // no operation is accepted anymore
	dirty               map[ecore.EObject]struct{}   // objects of the pending operations and of the open transaction
	dirtyGlobal         bool                         // the pending operations or the open transaction touch the roots or arbitrary tables
	pending             map[objectFeature]*operation // last pending Set or UnSet by object and feature
	errs                []error                      // errors of the operations since the last flush
	failure             error                        // a transaction was lost : no operation is accepted anymore
}

// ErrSQLStoreFailed is the error of the operations of a store that lost a transaction (the
// database aborted it or its commit failed) : the writes of the transaction, already accepted,
// are not in the database and the ids the store registered for their objects are wrong, so the
// store accepts no operation anymore. The database holds the transactions committed before.
var ErrSQLStoreFailed = errors.New("sql store failed")

// ErrSQLStoreOperationTimeout is the error of an operation not executed by the writer within
// the operation timeout. A read is then not executed, a write is still executed
var ErrSQLStoreOperationTimeout = errors.New("sql store operation timeout")

func backupDB(dstConn, srcConn *sqlite.Conn) error {
	backup, err := sqlite.NewBackup(dstConn, "main", srcConn, "main")
	if err != nil {
		return err
	}
	if more, err := backup.Step(-1); err != nil {
		return err
	} else if more {
		return errors.New("full backup step with remaining pages")
	}
	if err := backup.Close(); err != nil {
		return err
	}
	return nil
}

func NewSQLStore(
	databasePath string,
	resourceURI *ecore.URI,
	idManager ecore.EObjectIDManager,
	packageRegistry ecore.EPackageRegistry,
	options map[string]any) (store *SQLStore, err error) {

	inMemoryDatabase := false
	if options != nil {
		inMemoryDatabase, _ = options[SQL_OPTION_IN_MEMORY_DATABASE].(bool)
	}
	if inMemoryDatabase {
		// memory database shared by the connections of the pool and the writer
		name := filepath.Base(databasePath)
		dbPath := fmt.Sprintf("file:%s?mode=memory&cache=shared", name)
		return newSQLStore(
			func() (*sqlitex.Pool, error) {
				connSrc, err := sqlite.OpenConn(databasePath)
				if err != nil {
					return nil, err
				}
				defer connSrc.Close()

				// create connection pool
				connPool, err := sqlitex.NewPool(dbPath, sqlitex.PoolOptions{Flags: sqlite.OpenCreate | sqlite.OpenReadWrite | sqlite.OpenURI})
				if err != nil {
					return nil, err
				}

				// create connection pool
				connDst, err := connPool.Take(context.Background())
				if err != nil {
					return nil, err
				}
				defer connPool.Put(connDst)

				// backup src db to dst db
				if err := backupDB(connDst, connSrc); err != nil {
					return nil, err
				}

				return connPool, nil
			},
			func(connPool *sqlitex.Pool) error {
				return connPool.Close()
			},
			func() (*sqlite.Conn, error) {
				return sqlite.OpenConn(dbPath, sqlite.OpenCreate|sqlite.OpenReadWrite|sqlite.OpenURI)
			},
			func(conn *sqlite.Conn) (err error) {
				// the writer is closed last and committed : backup the memory db to the store db
				defer func() {
					if err2 := conn.Close(); err == nil {
						err = err2
					}
				}()
				connDst, err := sqlite.OpenConn(databasePath)
				if err != nil {
					return err
				}
				defer connDst.Close()
				return backupDB(connDst, conn)
			},
			resourceURI, idManager, packageRegistry, options,
		)
	} else {
		return newSQLStore(
			func() (*sqlitex.Pool, error) {
				return sqlitex.NewPool(databasePath, sqlitex.PoolOptions{
					Flags: sqlite.OpenReadWrite | sqlite.OpenCreate | sqlite.OpenWAL,
					PrepareConn: func(conn *sqlite.Conn) error {
						// connection is in synchronous mode
						return sqlitex.ExecuteTransient(conn, "PRAGMA synchronous=normal", nil)
					},
				})
			},
			func(pool *sqlitex.Pool) error {
				return pool.Close()
			},
			func() (*sqlite.Conn, error) {
				conn, err := sqlite.OpenConn(databasePath, sqlite.OpenReadWrite|sqlite.OpenCreate|sqlite.OpenWAL)
				if err != nil {
					return nil, err
				}
				if err := sqlitex.ExecuteTransient(conn, "PRAGMA synchronous=normal", nil); err != nil {
					conn.Close()
					return nil, err
				}
				return conn, nil
			},
			func(conn *sqlite.Conn) error {
				return conn.Close()
			},
			resourceURI,
			idManager,
			packageRegistry,
			options,
		)
	}
}

func newSQLStore(
	connectionPoolProvider func() (*sqlitex.Pool, error),
	connectionPoolClose func(pool *sqlitex.Pool) error,
	writerConnProvider func() (*sqlite.Conn, error),
	writerConnClose func(conn *sqlite.Conn) error,
	resourceURI *ecore.URI,
	idManager ecore.EObjectIDManager,
	packageRegistry ecore.EPackageRegistry,
	options map[string]any) (store *SQLStore, err error) {

	objectIDName := ""
	codecVersion := sqlCodecVersion
	sqlIDManager := newSQLStoreIDManager()
	sqlObjectManager := newSQLStoreObjectManager()
	logger := zap.NewNop()
	isKeepDefaults := false
	timeoutOperation := unlockOperationTimeout
	maxAllocSize := SQLITE_MAX_ALLOCATION_SIZE
	commitMaxOperations := defaultCommitMaxOperations
	commitMaxDuration := defaultCommitMaxDuration
	inMemoryDatabase := false
	if options != nil {
		objectIDName, _ = options[SQL_OPTION_OBJECT_ID].(string)
		if v, isVersion := options[SQL_OPTION_CODEC_VERSION].(int64); isVersion {
			codecVersion = v
		}
		if m, isSQLIDManager := options[SQL_OPTION_SQL_ID_MANAGER].(SQLStoreIDManager); isSQLIDManager {
			sqlIDManager = m
		}
		if l, isLogger := options[SQL_OPTION_LOGGER]; isLogger {
			logger = l.(*zap.Logger)
		}
		if b, isBool := options[SQL_OPTION_KEEP_DEFAULTS].(bool); isBool {
			isKeepDefaults = b
		}
		if t, isTimeout := options[SQL_OPTION_OPERATION_TIMEOUT].(time.Duration); isTimeout {
			timeoutOperation = t
		}
		if v, isMaxAllocSize := options[SQL_OPTION_MAX_ALLOC_SIZE].(int); isMaxAllocSize {
			maxAllocSize = v
		}
		if v, isInt := options[SQL_OPTION_COMMIT_MAX_OPERATIONS].(int); isInt && v > 0 {
			commitMaxOperations = v
		}
		if v, isDuration := options[SQL_OPTION_COMMIT_MAX_DURATION].(time.Duration); isDuration && v > 0 {
			commitMaxDuration = v
		}
		inMemoryDatabase, _ = options[SQL_OPTION_IN_MEMORY_DATABASE].(bool)
	}

	// retrieve connection pool
	connPool, err := connectionPoolProvider()
	if err != nil {
		return nil, err
	}

	// close pool if there is an error
	defer func() {
		if err != nil {
			_ = connectionPoolClose(connPool)
		}
	}()

	// writer connection
	writerConn, err := writerConnProvider()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = writerConn.Close()
		}
	}()

	// ants promise pool
	antsPool, _ := ants.NewPool(math.MaxInt32, ants.WithExpiryDuration(5*time.Second), ants.WithLogger(&zapLogger{logger.Named("ants")}))
	promisePool := newAntsPromisePool(antsPool)

	// create sql base
	base := &sqlBase{
		codecVersion:     codecVersion,
		uri:              resourceURI,
		objectIDName:     objectIDName,
		objectIDManager:  idManager,
		isContainerID:    true,
		isObjectID:       len(objectIDName) > 0 && objectIDName != "objectID" && idManager != nil,
		sqliteQueries:    map[string][]*query{},
		antsPool:         antsPool,
		promisePool:      promisePool,
		connPool:         connPool,
		connPoolProvider: connectionPoolProvider,
		connPoolClose:    connectionPoolClose,
	}

	// create sql store
	store = &SQLStore{
		sqlBase: base,
		sqlDecoder: sqlDecoder{
			sqlBase:          base,
			packageRegistry:  packageRegistry,
			sqlIDManager:     sqlIDManager,
			sqlObjectManager: sqlObjectManager,
			classDataMap:     map[ecore.EClass]*sqlDecoderClassData{},
			objectLocks:      newSqlEncoderLockManager(),
		},
		sqlEncoder: sqlEncoder{
			sqlBase:          base,
			isForced:         false,
			isKeepDefaults:   isKeepDefaults,
			classDataMap:     map[ecore.EClass]*sqlEncoderClassData{},
			sqlIDManager:     sqlIDManager,
			sqlObjectManager: sqlObjectManager,
			sqlLockManager:   newSqlEncoderLockManager(),
		},
		sqlIDManager:        sqlIDManager,
		singleQueries:       map[*sqlColumn]*sqlSingleQueries{},
		manyQueries:         map[*sqlTable]*sqlManyQueries{},
		timeoutOperation:    timeoutOperation,
		maxAllocSize:        int64(maxAllocSize),
		writerConn:          writerConn,
		writerConnClose:     writerConnClose,
		writerSignal:        make(chan struct{}, 1),
		writerDone:          make(chan struct{}),
		commitMaxOperations: commitMaxOperations,
		commitMaxDuration:   commitMaxDuration,
		// in a memory database, the tables are locked, not versioned : a pool connection
		// reading a table written by the open transaction would wait for its commit
		concurrentReads: !inMemoryDatabase,
		dirty:           map[ecore.EObject]struct{}{},
		pending:         map[objectFeature]*operation{},
	}

	// set store logger
	store.setLogger(logger)

	// set store in sql object manager
	sqlObjectManager.store = store

	// initialize the db with the writer connection
	unbind := store.bindConn(writerConn)
	defer unbind()

	// decode version
	if err = store.decodeVersion(); err != nil {
		return nil, err
	}

	// decode schema
	if err = store.decodeSchema([]sqlSchemaOption{withCreateIfNotExists(true)}); err != nil {
		return nil, err
	}

	// encode version
	if err = store.encodeVersion(); err != nil {
		return nil, err
	}

	// encode schema
	if err = store.encodeSchema(); err != nil {
		return nil, err
	}

	// encode schema
	if err = store.encodeProperties(); err != nil {
		return nil, err
	}

	// launch writer
	go store.runWriter()

	return store, nil
}

func (s *SQLStore) setLogger(logger *zap.Logger) {
	s.sqlBase.setLogger(logger)
	s.sqlEncoder.setLogger(logger)
	s.loggerOperations = logger.Named("ops")
}

func (s *SQLStore) getSingleQueries(column *sqlColumn) *sqlSingleQueries {
	s.mutexQueries.Lock()
	defer s.mutexQueries.Unlock()
	stmts := s.singleQueries[column]
	if stmts == nil {
		stmts = newSqlSingleQueries(column)
		s.singleQueries[column] = stmts
	}
	return stmts
}

func (s *SQLStore) getManyQueries(table *sqlTable) *sqlManyQueries {
	s.mutexQueries.Lock()
	defer s.mutexQueries.Unlock()
	stmts := s.manyQueries[table]
	if stmts == nil {
		stmts = newSqlManyQueries(table)
		s.manyQueries[table] = stmts
	}
	return stmts
}

// getManyFeatureTable returns the table of a many feature
func (s *SQLStore) getManyFeatureTable(featureSchema *sqlFeatureSchema) (*sqlTable, error) {
	if featureSchema.table == nil {
		return nil, fmt.Errorf("feature %s is not a list", featureSchema.feature.GetName())
	}
	return featureSchema.table, nil
}

func (s *SQLStore) getFeatureSchema(object ecore.EObject, feature ecore.EStructuralFeature) (*sqlFeatureSchema, error) {
	// retrieve the schema of the feature, declared or inherited by the object class
	featureSchema := s.sqlDecoder.schema.getFeatureSchema(object.EClass(), feature)
	if featureSchema == nil {
		return nil, fmt.Errorf("feature %s is unknown", feature.GetName())
	}
	return featureSchema, nil
}

func (s *SQLStore) getFeatureTable(object ecore.EObject, feature ecore.EStructuralFeature) (*sqlTable, error) {
	featureSchema, err := s.getFeatureSchema(object, feature)
	if err != nil {
		return nil, err
	}
	return s.getManyFeatureTable(featureSchema)
}

func (s *SQLStore) getEncoderFeatureData(object ecore.EObject, feature ecore.EStructuralFeature) (*sqlEncoderFeatureData, error) {
	// retrieve class schema
	class := object.EClass()

	// retrieve class data
	classData, err := s.getEncoderClassData(class)
	if err != nil {
		return nil, fmt.Errorf("class %s is unknown %w", class.GetName(), err)
	}

	// retrieve feature data, declared or inherited by the object class
	featureData, isFeatureData := classData.getFeatureData(feature)
	if !isFeatureData {
		return nil, fmt.Errorf("feature %s is unknown", feature.GetName())
	}

	return featureData, nil
}

func (e *SQLStore) encodeFeatureValue(sqlObjectID int64, featureData *sqlEncoderFeatureData, value any) (encoded any, err error) {
	if value != nil {
		switch featureData.schema.featureKind {
		case sfkObject, sfkObjectList:
			eObject := value.(ecore.EObject)
			return e.encodeSQLID(eObject, -1, -1)
		default:
			return e.sqlEncoder.encodeFeatureValue(sqlObjectID, featureData, value)
		}
	}
	return nil, nil
}

func (s *SQLStore) getSQLID(eObject ecore.EObject) (int64, bool, error) {
	// retrieve sql id for eObject
	sqlObjectID, isSqlObjectID := s.sqlIDManager.GetObjectID(eObject)
	if !isSqlObjectID {
		// object is not in store - check if it exists in db
		if sqlObjectID != 0 {
			if err := s.executeQuery(
				`SELECT objectID FROM ".objects" WHERE objectID=?`,
				&sqlitex.ExecOptions{
					Args: []any{sqlObjectID},
					ResultFunc: func(stmt *sqlite.Stmt) error {
						isSqlObjectID = true
						return nil
					},
				}); err != nil {
				return 0, false, err
			}
		}
		if isSqlObjectID {
			// object exists - register it as decoded
			s.sqlIDManager.SetObjectID(eObject, sqlObjectID)
		}
	}
	return sqlObjectID, isSqlObjectID, nil
}

// encode object sql id
func (s *SQLStore) encodeSQLID(eObject ecore.EObject, sqlContainerID int64, containerFeatureID int64) (int64, error) {
	if sqlObjectID, isSqlObjectID, err := s.getSQLID(eObject); err != nil {
		return 0, err
	} else if isSqlObjectID {
		return sqlObjectID, nil
	} else {
		return s.encodeObject(eObject, sqlContainerID, containerFeatureID)
	}
}

func mapSet[S ~map[E]struct{}, E comparable, R any](m S, mapper func(E) R) []R {
	i := 0
	mappedSlice := make([]R, len(m))
	for e := range m {
		mappedSlice[i] = mapper(e)
		i++
	}
	return mappedSlice
}

const (
	// unlockOperationTimeout is the time an operation waits for the writer before failing with
	// ErrSQLStoreOperationTimeout
	unlockOperationTimeout = 5 * time.Second
	// defaultCommitMaxOperations is the number of operations after which the writer commits
	// when the queue is never empty
	defaultCommitMaxOperations = 10000
	// defaultCommitMaxDuration is the duration after which the writer commits when the queue
	// is never empty
	defaultCommitMaxDuration = time.Second
	// maxOperationErrors is the number of errors kept for the next Flush
	maxOperationErrors = 100
)

// WaitOperations waits for all the pending operations to be executed and committed.
// object is ignored : the operations are executed in order, so waiting for the
// operations of one object is waiting for all the operations before them.
// It returns the errors of the operations since the last wait.
func (s *SQLStore) WaitOperations(ctx context.Context, object any) error {
	return s.Flush(ctx)
}

// Flush waits for all the pending operations to be executed and committed, and
// returns the errors of the operations since the last flush. Once the store has failed,
// it returns ErrSQLStoreFailed
func (s *SQLStore) Flush(ctx context.Context) error {
	if goid.Get() == s.writerGoid.Load() {
		// from an operation being executed : the operations before it are executed
		return nil
	}
	op := newOperation("Flush", operationWrite, nil, nil, -1, nil, func() (any, error) {
		s.mutexWriter.Lock()
		errs := s.errs
		s.errs = nil
		s.mutexWriter.Unlock()
		return nil, errors.Join(errs...)
	}, withContext(ctx), withFlush(), withWait())
	_, err := s.execute(op)
	return err
}

// Err returns ErrSQLStoreFailed, wrapping its cause, once the store has lost a transaction
func (s *SQLStore) Err() error {
	s.mutexWriter.Lock()
	defer s.mutexWriter.Unlock()
	return s.failure
}

// fail makes the store refuse every operation : the transaction is lost
func (s *SQLStore) fail(cause error) {
	s.mutexWriter.Lock()
	if s.failure == nil {
		s.failure = fmt.Errorf("%w: transaction lost: %w", ErrSQLStoreFailed, cause)
	}
	s.mutexWriter.Unlock()
	s.loggerOperations.Error("store failed", zap.Error(cause))
}

// isTransactionOpen is true while the writer has a transaction open
func (s *SQLStore) isTransactionOpen() bool {
	s.mutexWriter.Lock()
	defer s.mutexWriter.Unlock()
	return s.txOpen
}

// markDirty registers the objects an operation writes : until the transaction
// commits, their reads are executed by the writer. must be called with mutexWriter
func (s *SQLStore) markDirty(op *operation) {
	if op.type_ != operationWrite {
		return
	}
	if op.object == nil {
		s.dirtyGlobal = true
	} else {
		s.dirty[op.object] = struct{}{}
	}
	if value, isObject := op.value.(ecore.EObject); isObject {
		s.dirty[value] = struct{}{}
	}
}

// registerDirty registers an object encoded or decoded by the open transaction
func (s *SQLStore) registerDirty(object ecore.EObject) {
	if goid.Get() == s.writerGoid.Load() {
		s.mutexWriter.Lock()
		s.dirty[object] = struct{}{}
		s.mutexWriter.Unlock()
	}
}

// enqueue gives an operation to the writer
func (s *SQLStore) enqueue(op *operation) {
	s.mutexWriter.Lock()
	if s.closed || s.failure != nil {
		failure := s.failure
		s.mutexWriter.Unlock()
		if failure != nil {
			op.err = failure
		} else {
			op.err = errors.New("store is closed")
		}
		s.loggerOperations.Error("rejected", zap.Object("operation", &operationMarshaler{op: op}), zap.Error(op.err))
		if op.done != nil {
			close(op.done)
		}
		return
	}
	op.enqueued = time.Now()
	s.queue = append(s.queue, op)
	s.markDirty(op)
	if op.type_ == operationWrite && op.feature != nil && (op.cmd == "Set" || op.cmd == "UnSet") {
		s.pending[objectFeature{op.object, op.feature}] = op
	}
	s.mutexWriter.Unlock()
	if e := s.loggerOperations.Check(zap.DebugLevel, "enqueue"); e != nil {
		e.Write(zap.Int64("goid", goid.Get()), zap.Object("operation", &operationMarshaler{op: op}))
	}
	select {
	case s.writerSignal <- struct{}{}:
	default:
	}
}

// pendingOperation returns the pending Set or UnSet of an object feature at index, if any
func (s *SQLStore) pendingOperation(object ecore.EObject, feature ecore.EStructuralFeature, index int) *operation {
	s.mutexWriter.Lock()
	defer s.mutexWriter.Unlock()
	if op := s.pending[objectFeature{object, feature}]; op != nil && op.index == index {
		return op
	}
	return nil
}

// await waits for an enqueued operation to be executed
func (s *SQLStore) await(op *operation) (any, error) {
	if op.done == nil {
		return nil, fmt.Errorf("operation %s was not created withWait", op.cmd)
	}
	ctx := op.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	timer := time.NewTimer(s.timeoutOperation)
	defer timer.Stop()
	select {
	case <-op.done:
		return op.result, op.err
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-timer.C:
		// the caller may hold a lock the writer needs (a store list read while the writer
		// encodes its owner) : it stops waiting and releases it. A read not executed yet is
		// skipped, a write is executed anyway
		op.abandoned.Store(true)
		s.loggerOperations.Error("timeout", zap.Int64("goid", goid.Get()), zap.Object("operation", &operationMarshaler{op: op}), zap.Duration("timeout", s.timeoutOperation))
		return nil, ErrSQLStoreOperationTimeout
	}
}

// canReadConcurrently is true when a read doesn't depend on a pending or uncommitted
// write : it can be executed on a pool connection, against the committed state
func (s *SQLStore) canReadConcurrently(op *operation) bool {
	if !s.concurrentReads {
		return false
	}
	s.mutexWriter.Lock()
	defer s.mutexWriter.Unlock()
	if s.closed || s.failure != nil || s.dirtyGlobal {
		return false
	}
	if op.object == nil {
		// the whole store : only when idle
		return !s.busy && !s.txOpen && len(s.queue) == 0
	}
	if _, isDirty := s.dirty[op.object]; isDirty {
		return false
	}
	if value, isObject := op.value.(ecore.EObject); isObject {
		if _, isDirty := s.dirty[value]; isDirty {
			return false
		}
	}
	return true
}

// execute executes an operation :
// - from the writer goroutine (an operation calling the store), now on the writer connection
// - a read independent of the pending writes, now on a pool connection
// - otherwise by the writer, waiting for its result if the operation was created withWait
func (s *SQLStore) execute(op *operation) (any, error) {
	if goid.Get() == s.writerGoid.Load() {
		return op.fn()
	}
	if op.type_ == operationRead && s.canReadConcurrently(op) {
		return s.executeConcurrently(op)
	}
	s.enqueue(op)
	if op.done == nil {
		return nil, nil
	}
	return s.await(op)
}

// executeConcurrently executes a read on a pool connection, in the current goroutine
func (s *SQLStore) executeConcurrently(op *operation) (result any, err error) {
	ctx := op.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	conn, err := s.connPool.Take(ctx)
	if err != nil {
		return nil, err
	}
	unbind := s.bindConn(conn)
	defer func() {
		if r := recover(); r != nil {
			result, err = nil, s.panicError(op, r)
			// reset the statements left stepped by the panic
			conn.SetInterrupt(nil)
		}
		unbind()
		s.connPool.Put(conn)
	}()
	if e := s.loggerOperations.Check(zap.DebugLevel, "execute"); e != nil {
		e.Write(zap.Int64("goid", goid.Get()), zap.Object("operation", &operationMarshaler{op: op}))
	}
	result, err = op.fn()
	if err != nil {
		s.loggerOperations.Error("error", zap.Object("operation", &operationMarshaler{op: op}), zap.Error(err))
	}
	return
}

// resultAs converts the result of an operation, or returns the zero value on error
func resultAs[T any](result any, err error) T {
	var def T
	if err != nil {
		return def
	} else if r, isResult := result.(T); isResult {
		return r
	}
	return def
}

// indexResult converts the result of an index operation, or returns -1 on error
func indexResult(result any, err error) int {
	if index, isIndex := result.(int); err == nil && isIndex {
		return index
	}
	return -1
}

func (s *SQLStore) panicError(op *operation, recovered any) error {
	var err error
	switch v := recovered.(type) {
	case error:
		err = v
	default:
		err = fmt.Errorf("%+v", v)
	}
	s.loggerOperations.Error("panic",
		zap.Object("operation", &operationMarshaler{op: op}),
		zap.Error(err),
		zap.Stack("stack"))
	return err
}

func (s *SQLStore) addError(err error) {
	s.mutexWriter.Lock()
	if len(s.errs) < maxOperationErrors {
		s.errs = append(s.errs, err)
	}
	s.mutexWriter.Unlock()
}

// runWriter executes the operations of the queue in order on the writer connection
func (s *SQLStore) runWriter() {
	defer close(s.writerDone)
	s.writerGoid.Store(goid.Get())
	unbind := s.bindConn(s.writerConn)
	defer unbind()
	for {
		op := s.nextOperation()
		if op == nil {
			// closed and drained
			s.commit()
			return
		}
		s.executeOperation(op)
		s.mutexWriter.Lock()
		empty := len(s.queue) == 0
		s.mutexWriter.Unlock()
		if op.flush || empty || s.txOperations >= s.commitMaxOperations || time.Since(s.txStart) >= s.commitMaxDuration {
			s.commit()
		}
		s.mutexWriter.Lock()
		s.busy = false
		s.mutexWriter.Unlock()
		if op.done != nil {
			close(op.done)
		}
	}
}

// nextOperation waits for the next operation of the queue; nil when the store is closed and the queue empty
func (s *SQLStore) nextOperation() *operation {
	for {
		s.mutexWriter.Lock()
		if len(s.queue) > 0 {
			op := s.queue[0]
			s.queue[0] = nil
			s.queue = s.queue[1:]
			s.busy = true
			s.mutexWriter.Unlock()
			return op
		}
		closed := s.closed
		s.mutexWriter.Unlock()
		if closed {
			return nil
		}
		<-s.writerSignal
	}
}

func (s *SQLStore) executeOperation(op *operation) {
	// executed, failed or skipped : the write is not pending anymore
	defer s.clearPending(op)
	if op.ctx != nil && op.ctx.Err() != nil {
		// nobody waits for the result anymore
		op.err = op.ctx.Err()
		return
	}
	if op.type_ == operationRead && op.abandoned.Load() {
		op.err = ErrSQLStoreOperationTimeout
		return
	}
	if err := s.Err(); err != nil {
		// enqueued before the failure
		op.err = err
		return
	}
	if op.commitBefore {
		s.commit()
	}
	if op.flush {
		// committed first : the errors collected include the ones of the commit
		s.commit()
		op.result, op.err = op.fn()
		return
	}
	if !s.txOpen {
		if err := s.begin(); err != nil {
			op.err = err
			s.addError(err)
			return
		}
	}
	s.txOperations++
	if e := s.loggerOperations.Check(zap.DebugLevel, "execute"); e != nil {
		e.Write(zap.Int64("goid", goid.Get()), zap.Object("operation", &operationMarshaler{op: op}))
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				op.result, op.err = nil, s.panicError(op, r)
			}
		}()
		op.result, op.err = op.fn()
	}()
	if op.err != nil {
		s.loggerOperations.Error("error", zap.Object("operation", &operationMarshaler{op: op}), zap.Error(op.err))
		s.addError(op.err)
		// reset the statements left stepped by a panic
		s.writerConn.SetInterrupt(nil)
		// some errors roll the transaction back : its writes are lost
		if s.txOpen && s.writerConn.AutocommitEnabled() {
			s.loggerOperations.Error("transaction rolled back", zap.Object("operation", &operationMarshaler{op: op}), zap.Error(op.err))
			s.mutexWriter.Lock()
			s.txOpen = false
			s.txOperations = 0
			s.mutexWriter.Unlock()
			s.fail(op.err)
		}
	}
}

// clearPending forgets op as the pending Set or UnSet of its object feature
func (s *SQLStore) clearPending(op *operation) {
	if op.type_ == operationWrite && op.feature != nil {
		s.mutexWriter.Lock()
		if key := (objectFeature{op.object, op.feature}); s.pending[key] == op {
			delete(s.pending, key)
		}
		s.mutexWriter.Unlock()
	}
}

func (s *SQLStore) begin() error {
	if err := sqlitex.ExecuteTransient(s.writerConn, "BEGIN IMMEDIATE;", nil); err != nil {
		s.loggerOperations.Error("begin", zap.Error(err))
		return err
	}
	s.mutexWriter.Lock()
	s.txOpen = true
	s.txStart = time.Now()
	s.txOperations = 0
	s.mutexWriter.Unlock()
	s.loggerOperations.Debug("begin")
	return nil
}

// commit commits the open transaction. the objects of the operations still queued stay dirty
func (s *SQLStore) commit() {
	if !s.txOpen {
		return
	}
	if err := sqlitex.ExecuteTransient(s.writerConn, "COMMIT;", nil); err != nil {
		s.loggerOperations.Error("commit", zap.Error(err))
		s.addError(err)
		if !s.writerConn.AutocommitEnabled() {
			_ = sqlitex.ExecuteTransient(s.writerConn, "ROLLBACK;", nil)
		}
		s.fail(err)
	} else if e := s.loggerOperations.Check(zap.DebugLevel, "commit"); e != nil {
		e.Write(zap.Int("operations", s.txOperations), zap.Duration("duration", time.Since(s.txStart)))
	}
	s.mutexWriter.Lock()
	s.txOpen = false
	s.txOperations = 0
	s.dirty = map[ecore.EObject]struct{}{}
	s.dirtyGlobal = false
	for _, op := range s.queue {
		s.markDirty(op)
	}
	s.mutexWriter.Unlock()
}

func (s *SQLStore) Get(object ecore.EObject, feature ecore.EStructuralFeature, index int) any {
	// optimization : if we have a non many feature, check if we have a pending write operation
	if !feature.IsMany() {
		if op := s.pendingOperation(object, feature, index); op != nil {
			if op.cmd == "UnSet" {
				return feature.GetDefaultValue()
			}
			return op.value
		}
	}
	return resultAs[any](s.execute(
		newOperation("Get", operationRead, object, feature, index, nil, func() (any, error) {
			return s.doGet(object, feature, index)
		}, withWait())))
}

func (s *SQLStore) doGet(object ecore.EObject, feature ecore.EStructuralFeature, index int) (any, error) {
	// retrieve value from sqlite db
	sqlObjectID, isSQLObjectID, err := s.getSQLID(object)
	if err != nil {
		return nil, err
	} else if !isSQLObjectID {
		return nil, nil
	}
	featureSchema, err := s.getFeatureSchema(object, feature)
	if err != nil {
		return nil, err
	}
	return s.getValue(sqlObjectID, featureSchema, index)
}

func (s *SQLStore) getValue(sqlID int64, featureSchema *sqlFeatureSchema, index int) (any, error) {
	var query string
	var args []any
	if featureColumn := featureSchema.column; featureColumn != nil {
		query = s.getSingleQueries(featureColumn).getSelectQuery()
		args = []any{sqlID}
	} else if featureTable := featureSchema.table; featureTable != nil {
		query = s.getManyQueries(featureTable).getSelectOneQuery()
		args = []any{sqlID, index}
	}

	var value any
	if err := s.executeQuery(
		query,
		&sqlitex.ExecOptions{
			Args: args,
			ResultFunc: func(stmt *sqlite.Stmt) error {
				value = decodeAny(stmt, 0)
				return nil
			}}); err != nil {
		return nil, err
	}

	return s.decodeFeatureValue(featureSchema, value)
}

func (s *SQLStore) Set(object ecore.EObject, feature ecore.EStructuralFeature, index int, value any, needResult bool) (oldValue any) {
	var options []operationOption
	if needResult {
		options = append(options, withWait())
	}
	result, err := s.execute(
		newOperation("Set", operationWrite, object, feature, index, value, func() (any, error) {
			return s.doSet(object, feature, index, value, needResult)
		}, options...))
	if needResult {
		oldValue = resultAs[any](result, err)
	}
	return
}

func (s *SQLStore) doSet(object ecore.EObject, feature ecore.EStructuralFeature, index int, value any, needResult bool) (oldValue any, err error) {
	// get object sql id
	sqlObjectID, err := s.encodeSQLID(object, -1, -1)
	if err != nil {
		return nil, err
	}

	// get encoder feature data
	featureData, err := s.getEncoderFeatureData(object, feature)
	if err != nil {
		return nil, err
	}

	if needResult {
		// retrieve previous value
		if oldValue, err = s.getValue(sqlObjectID, featureData.schema, index); err != nil {
			return nil, err
		}
	}

	// encode value
	encoded, err := s.encodeFeatureValue(sqlObjectID, featureData, value)
	if err != nil {
		return nil, err
	}

	var query string
	var args []any
	if featureColumn := featureData.schema.column; featureColumn != nil {
		query = s.getSingleQueries(featureColumn).getUpdateQuery()
		args = []any{encoded, sqlObjectID}
	} else if featureTable := featureData.schema.table; featureTable != nil {
		query = s.getManyQueries(featureTable).getUpdateValueQuery()
		args = []any{encoded, sqlObjectID, index}
	}

	if err := s.executeQuery(query, &sqlitex.ExecOptions{Args: args}); err != nil {
		return nil, err
	}

	return oldValue, nil
}

func (s *SQLStore) IsSet(object ecore.EObject, feature ecore.EStructuralFeature) bool {
	// optimization : if we have a non many feature, check if we have a pending write operation
	if !feature.IsMany() {
		if op := s.pendingOperation(object, feature, -1); op != nil {
			return op.value != nil
		}
	}
	return resultAs[bool](s.execute(
		newOperation("IsSet", operationRead, object, feature, -1, nil, func() (any, error) {
			return s.doIsSet(object, feature)
		}, withWait())))
}

func (s *SQLStore) doIsSet(object ecore.EObject, feature ecore.EStructuralFeature) (bool, error) {
	sqlObjectID, isSQLObjectID, err := s.getSQLID(object)
	if err != nil {
		return false, err
	} else if !isSQLObjectID {
		return false, nil
	}
	featureSchema, err := s.getFeatureSchema(object, feature)
	if err != nil {
		return false, err
	}
	if featureColumn := featureSchema.column; featureColumn != nil {
		var value any
		if err := s.executeQuery(
			s.getSingleQueries(featureColumn).getSelectQuery(),
			&sqlitex.ExecOptions{
				Args: []any{sqlObjectID},
				ResultFunc: func(stmt *sqlite.Stmt) error {
					value = decodeAny(stmt, 0)
					return nil
				}}); err != nil {
			return false, err
		}
		// a null column is not set; an object is set; otherwise the decoded value is compared
		// with the default
		if value == nil {
			return false, nil
		}
		switch featureSchema.featureKind {
		case sfkObject, sfkObjectReference:
			return true, nil
		}
		decoded, err := s.decodeFeatureValue(featureSchema, value)
		if err != nil {
			return false, err
		}
		return !reflect.DeepEqual(decoded, featureSchema.feature.GetDefaultValue()), nil
	} else if featureTable := featureSchema.table; featureTable != nil {
		isSet := false
		if err := s.executeQuery(
			s.getManyQueries(featureTable).getExistsQuery(),
			&sqlitex.ExecOptions{
				Args: []any{sqlObjectID},
				ResultFunc: func(stmt *sqlite.Stmt) error {
					isSet = stmt.ColumnBool(0)
					return nil
				}}); err != nil {
			return false, err
		}
		return isSet, nil
	}
	return false, nil
}

func (s *SQLStore) UnSet(object ecore.EObject, feature ecore.EStructuralFeature) {
	_, _ = s.execute(
		newOperation("UnSet", operationWrite, object, feature, -1, nil, func() (any, error) {
			return s.doUnSet(object, feature)
		}, withWait()))
}

func (s *SQLStore) doUnSet(object ecore.EObject, feature ecore.EStructuralFeature) (any, error) {
	sqlObjectID, err := s.encodeSQLID(object, -1, -1)
	if err != nil {
		return nil, err
	}

	featureData, err := s.getEncoderFeatureData(object, feature)
	if err != nil {
		return nil, err
	}

	var query string
	var args []any
	if featureColumn := featureData.schema.column; featureColumn != nil {
		// the default is encoded as any value : an enum literal is its id, a data type its string
		encoded, err := s.encodeFeatureValue(sqlObjectID, featureData, feature.GetDefaultValue())
		if err != nil {
			return nil, err
		}
		query = s.getSingleQueries(featureColumn).getUpdateQuery()
		args = []any{encoded, sqlObjectID}
	} else if featureTable := featureData.schema.table; featureTable != nil {
		query = s.getManyQueries(featureTable).getClearQuery()
		args = []any{sqlObjectID}
	}
	if err := s.executeQuery(query, &sqlitex.ExecOptions{Args: args}); err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *SQLStore) IsEmpty(object ecore.EObject, feature ecore.EStructuralFeature) bool {
	// empty on error
	result, err := s.execute(
		newOperation("IsEmpty", operationRead, object, feature, -1, nil, func() (any, error) {
			return s.doIsEmpty(object, feature)
		}, withWait()))
	if err != nil {
		return true
	}
	return resultAs[bool](result, nil)
}

func (s *SQLStore) doIsEmpty(object ecore.EObject, feature ecore.EStructuralFeature) (bool, error) {
	sqlObjectID, isSQLObjectID, err := s.getSQLID(object)
	if err != nil {
		return true, err
	} else if !isSQLObjectID {
		return true, nil
	}
	// retrieve table
	featureTable, err := s.getFeatureTable(object, feature)
	if err != nil {
		return true, err
	}

	// retrieve statement
	isEmpty := true
	if err := s.executeQuery(
		s.getManyQueries(featureTable).getExistsQuery(),
		&sqlitex.ExecOptions{
			Args: []any{sqlObjectID},
			ResultFunc: func(stmt *sqlite.Stmt) error {
				isEmpty = !stmt.ColumnBool(0)
				return nil
			}}); err != nil {
		if isNoSuchTable(err) {
			return true, nil
		}
		return true, err
	}
	return isEmpty, nil

}

func (s *SQLStore) Size(object ecore.EObject, feature ecore.EStructuralFeature) int {
	return resultAs[int](s.execute(
		newOperation("Size", operationRead, object, feature, -1, nil, func() (any, error) {
			return s.doSize(object, feature)
		}, withWait())))
}

func (s *SQLStore) doSize(object ecore.EObject, feature ecore.EStructuralFeature) (int, error) {
	sqlObjectID, isSQLObjectID, err := s.getSQLID(object)
	if err != nil {
		return 0, err
	} else if !isSQLObjectID {
		return 0, nil
	}
	// retrieve table
	featureTable, err := s.getFeatureTable(object, feature)
	if err != nil {
		return 0, err
	}

	var size int
	if err := s.executeQuery(
		s.getManyQueries(featureTable).getCountQuery(),
		&sqlitex.ExecOptions{
			Args: []any{sqlObjectID},
			ResultFunc: func(stmt *sqlite.Stmt) error {
				size = stmt.ColumnInt(0)
				return nil
			}}); err != nil {
		return 0, err
	}
	return size, nil
}

func (s *SQLStore) Contains(object ecore.EObject, feature ecore.EStructuralFeature, value any) bool {
	return resultAs[bool](s.execute(
		newOperation("Contains", operationRead, object, feature, -1, value, func() (any, error) {
			return s.doContains(object, feature, value)
		}, withWait())))
}

// encodeLookupValue encodes a value to look it up in a feature table, without writing :
// an object or an enum literal unknown to the store is not in the table, and isEncoded is false
func (s *SQLStore) encodeLookupValue(featureSchema *sqlFeatureSchema, value any) (encoded any, isEncoded bool, err error) {
	if value == nil {
		return nil, true, nil
	}
	switch featureSchema.featureKind {
	case sfkObject, sfkObjectList:
		sqlID, isSQLID, err := s.getSQLID(value.(ecore.EObject))
		if err != nil || !isSQLID {
			return nil, false, err
		}
		return sqlID, true, nil
	case sfkObjectReference, sfkObjectReferenceList:
		eObject := value.(ecore.EObject)
		if sqlID, _ := s.sqlIDManager.GetObjectID(eObject); sqlID != 0 {
			return strconv.FormatInt(sqlID, 10), true, nil
		}
		return s.uri.Relativize(ecore.GetURI(eObject)).String(), true, nil
	case sfkEnum:
		eAttribute := featureSchema.feature.(ecore.EAttribute)
		eEnum := eAttribute.GetEAttributeType().(ecore.EEnum)
		literal := eEnum.GetEPackage().GetEFactoryInstance().ConvertToString(eEnum, value)
		eEnumLiteral := eEnum.GetEEnumLiteralByLiteral(literal)
		if eEnumLiteral == nil {
			return nil, false, fmt.Errorf("unable to find enum literal in enum '%s' with value '%v'", eEnum.GetName(), literal)
		}
		if id, isID := s.sqlIDManager.GetEnumLiteralID(eEnumLiteral); isID {
			return id, true, nil
		}
		// look for the literal in the db
		table := s.schema.enumsTable
		var id int64
		isID := false
		if err := s.executeQuery(
			table.selectQuery(nil, sqlEscapeIdentifier(table.columns[2].columnName)+"=? AND "+sqlEscapeIdentifier(table.columns[3].columnName)+"=?", ""),
			&sqlitex.ExecOptions{
				Args: []any{eEnum.GetName(), eEnumLiteral.GetLiteral()},
				ResultFunc: func(stmt *sqlite.Stmt) error {
					id = stmt.ColumnInt64(0)
					isID = true
					return nil
				}}); err != nil {
			return nil, false, err
		}
		if !isID {
			return nil, false, nil
		}
		s.sqlIDManager.SetEnumLiteralID(eEnumLiteral, id)
		return id, true, nil
	case sfkDate:
		return value.(*time.Time).Format(time.RFC3339), true, nil
	case sfkData, sfkDataList:
		eDataType := featureSchema.feature.(ecore.EAttribute).GetEAttributeType()
		return eDataType.GetEPackage().GetEFactoryInstance().ConvertToString(eDataType, value), true, nil
	default:
		return value, true, nil
	}
}

// isNoSuchTable is true for the error of a query on a table not created yet : the feature
// tables are created by the first write of their class
func isNoSuchTable(err error) bool {
	return err != nil && strings.Contains(err.Error(), "no such table")
}

func (s *SQLStore) doContains(object ecore.EObject, feature ecore.EStructuralFeature, value any) (bool, error) {
	sqlObjectID, isSQLObjectID, err := s.getSQLID(object)
	if err != nil {
		return false, err
	} else if !isSQLObjectID {
		return false, nil
	}
	// retrieve table
	featureSchema, err := s.getFeatureSchema(object, feature)
	if err != nil {
		return false, err
	}

	// query statement
	encoded, isEncoded, err := s.encodeLookupValue(featureSchema, value)
	if err != nil || !isEncoded {
		return false, err
	}

	var rowid int64
	featureTable := featureSchema.table
	if err := s.executeQuery(
		s.getManyQueries(featureTable).getContainsQuery(),
		&sqlitex.ExecOptions{
			Args: []any{sqlObjectID, encoded},
			ResultFunc: func(stmt *sqlite.Stmt) error {
				rowid = stmt.ColumnInt64(0)
				return nil
			}}); err != nil {
		if isNoSuchTable(err) {
			return false, nil
		}
		return false, err
	}
	return rowid != 0, nil
}

func (s *SQLStore) doIndexOf(object ecore.EObject, feature ecore.EStructuralFeature, value any, getIndexOfQuery func(*sqlManyQueries) string) (int, error) {
	sqlObjectID, isSQLObjectID, err := s.getSQLID(object)
	if err != nil {
		return -1, err
	} else if !isSQLObjectID {
		return -1, nil
	}
	// retrieve table
	featureSchema, err := s.getFeatureSchema(object, feature)
	if err != nil {
		return -1, err
	}
	featureTable := featureSchema.table

	// compute parameters
	encoded, isEncoded, err := s.encodeLookupValue(featureSchema, value)
	if err != nil || !isEncoded {
		return -1, err
	}
	// retrieve row idx in table
	idx := -1.0
	if err := s.executeQuery(
		getIndexOfQuery(s.getManyQueries(featureTable)),
		&sqlitex.ExecOptions{
			Args: []any{sqlObjectID, encoded},
			ResultFunc: func(stmt *sqlite.Stmt) error {
				idx = stmt.ColumnFloat(0)
				return nil
			}}); err != nil {
		if isNoSuchTable(err) {
			return -1, nil
		}
		return -1, err
	}
	if idx == -1.0 {
		return -1, err
	}

	// convert idx to list index - index is the count of rows where idx < expected idx
	index := -1
	if err := s.executeQuery(
		s.getManyQueries(featureTable).getIdxToListIndexQuery(),
		&sqlitex.ExecOptions{
			Args: []any{sqlObjectID, idx},
			ResultFunc: func(stmt *sqlite.Stmt) error {
				index = stmt.ColumnInt(0)
				return nil
			}}); err != nil {
		return -1, err
	}
	return index, err
}

func (s *SQLStore) IndexOf(object ecore.EObject, feature ecore.EStructuralFeature, value any) int {
	return indexResult(s.execute(
		newOperation("IndexOf", operationRead, object, feature, -1, value, func() (any, error) {
			return s.doIndexOf(object, feature, value, func(sms *sqlManyQueries) string {
				return sms.getIndexOfQuery()
			})
		}, withWait())))
}

func (s *SQLStore) LastIndexOf(object ecore.EObject, feature ecore.EStructuralFeature, value any) int {
	return indexResult(s.execute(
		newOperation("LastIndexOf", operationRead, object, feature, -1, value, func() (any, error) {
			return s.doIndexOf(object, feature, value, func(sms *sqlManyQueries) string {
				return sms.getLastIndexOfQuery()
			})
		}, withWait())))
}

// AddRoot add object as store root
func (s *SQLStore) AddRoot(object ecore.EObject) {
	_, _ = s.execute(newOperation("AddRoot", operationWrite, nil, nil, -1, nil, func() (any, error) {
		return s.doAddRoot(object)
	}))
}

func (s *SQLStore) doAddRoot(object ecore.EObject) (any, error) {
	return nil, s.encodeContent(object)
}

// RemoveRoot implements EStore.
func (s *SQLStore) RemoveRoot(object ecore.EObject) {
	_, _ = s.execute(newOperation("RemoveRoot", operationWrite, nil, nil, -1, nil, func() (any, error) {
		return s.doRemoveRoot(object)
	}))
}

func (s *SQLStore) doRemoveRoot(object ecore.EObject) (any, error) {
	sqlObjectID, isSQLObjectID, err := s.getSQLID(object)
	if err != nil {
		return nil, err
	} else if !isSQLObjectID {
		return nil, err
	}
	contentColumn := s.schema.contentsTable.columns[0]
	if err := s.executeQuery(
		s.getSingleQueries(contentColumn).getRemoveQuery(),
		&sqlitex.ExecOptions{Args: []any{sqlObjectID}},
	); err != nil {
		return nil, err
	}
	return nil, nil
}

// GetRoot return root objects
func (s *SQLStore) GetRoots() []ecore.EObject {
	return resultAs[[]ecore.EObject](s.execute(
		newOperation("GetRoots", operationRead, nil, nil, -1, nil, func() (any, error) {
			return s.doGetRoots()
		}, withWait())))
}

func (s *SQLStore) doGetRoots() ([]ecore.EObject, error) {
	table := s.schema.contentsTable
	contents := []ecore.EObject{}
	if err := s.executeQuery(
		table.selectQuery(nil, "", ""),
		&sqlitex.ExecOptions{
			ResultFunc: func(stmt *sqlite.Stmt) error {
				// retrieve object id
				objectID := stmt.ColumnInt64(0)
				// decode object
				object, err := s.decodeObject(objectID)
				if err != nil {
					return err
				}
				// add object to contents
				contents = append(contents, object)
				return nil
			},
		}); err != nil {
		return nil, err
	}
	return contents, nil
}

func postOrderTraverse(object ecore.EObject, yield func(ecore.EObject) bool) bool {
	for child := range object.EContents().All() {
		if !postOrderTraverse(child.(ecore.EObject), yield) {
			return false
		}
	}
	return yield(object)
}

func postOrderIterator(object ecore.EObject) iter.Seq[ecore.EObject] {
	return func(yield func(ecore.EObject) bool) {
		postOrderTraverse(object, yield)
	}
}

// UnRegister forgets the ids of object, and of its contents if withContents, once the
// operations before it are executed : the objects are encoded again when written again
func (s *SQLStore) UnRegister(object ecore.EObject, withContents bool) *promise.Promise[any] {
	// the objects are collected now, in the caller's goroutine : walking the contents locks
	// the objects, and a reader may hold one while it waits for the writer
	objects := []ecore.EObject{object}
	if withContents {
		// post order : a parent is cleared after its children, so that clearing a child
		// never registers its parent again
		objects = slices.Collect(postOrderIterator(object))
	}
	op := newOperation("UnRegister", operationWrite, object, nil, -1, nil, func() (any, error) {
		for _, object := range objects {
			s.sqlIDManager.ClearObjectID(object)
		}
		return nil, nil
	}, withWait())
	if goid.Get() == s.writerGoid.Load() {
		// from an operation being executed
		result, err := op.fn()
		return promise.NewWithPool(func(resolve func(any), reject func(error)) {
			if err != nil {
				reject(err)
			} else {
				resolve(result)
			}
		}, s.promisePool)
	}
	s.enqueue(op)
	return promise.NewWithPool(func(resolve func(any), reject func(error)) {
		result, err := s.await(op)
		if err != nil {
			reject(err)
		} else {
			resolve(result)
		}
	}, s.promisePool)
}

func (s *SQLStore) Add(object ecore.EObject, feature ecore.EStructuralFeature, index int, value any) {
	_, _ = s.execute(
		newOperation("Add", operationWrite, object, feature, index, value, func() (any, error) {
			return s.doAdd(object, feature, index, value)
		}))
}

func (s *SQLStore) doAdd(object ecore.EObject, feature ecore.EStructuralFeature, index int, value any) (any, error) {
	sqlObjectID, err := s.encodeSQLID(object, -1, -1)
	if err != nil {
		return nil, err
	}
	featureData, err := s.getEncoderFeatureData(object, feature)
	if err != nil {
		return nil, err
	}
	featureTable, err := s.getManyFeatureTable(featureData.schema)
	if err != nil {
		return nil, err
	}
	idx, _, err := s.getInsertIdx(featureTable, sqlObjectID, index, 1)
	if err != nil {
		return nil, err
	}
	encoded, err := s.encodeFeatureValue(sqlObjectID, featureData, value)
	if err != nil {
		return nil, err
	}
	if err := s.executeQuery(
		s.getManyQueries(featureTable).getInsertQuery(),
		&sqlitex.ExecOptions{Args: []any{sqlObjectID, idx, encoded}},
	); err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *SQLStore) AddAll(object ecore.EObject, feature ecore.EStructuralFeature, index int, c ecore.Collection) {
	_, _ = s.execute(
		newOperation("AddAll", operationWrite, object, feature, index, c, func() (any, error) {
			return s.doAddAll(object, feature, index, c)
		}))
}

func (s *SQLStore) doAddAll(object ecore.EObject, feature ecore.EStructuralFeature, index int, c ecore.Collection) (any, error) {
	sqlObjectID, err := s.encodeSQLID(object, -1, -1)
	if err != nil {
		return nil, err
	}
	featureData, err := s.getEncoderFeatureData(object, feature)
	if err != nil {
		return nil, err
	}
	featureTable, err := s.getManyFeatureTable(featureData.schema)
	if err != nil {
		return nil, err
	}
	idx, inc, err := s.getInsertIdx(featureTable, sqlObjectID, index, c.Size())
	if err != nil {
		return nil, err
	}
	// encode each value
	query := s.getManyQueries(featureTable).getInsertQuery()
	for it := c.Iterator(); it.HasNext(); {
		value := it.Next()
		v, err := s.encodeFeatureValue(sqlObjectID, featureData, value)
		if err != nil {
			return nil, err
		}
		if err := s.executeQuery(
			query,
			&sqlitex.ExecOptions{Args: []any{sqlObjectID, idx, v}},
		); err != nil {
			return nil, err
		}
		idx += inc
	}
	return nil, nil
}

// compute insert idx of element in the list whose index is index. nb is the number of elements to be inserted
// return idx, increment (for each inserted element) and error if any
func (s *SQLStore) getInsertIdx(table *sqlTable, sqlID int64, index int, nb int) (float64, float64, error) {
	// idx of the elements at index-1 and index, if any
	var neighbors []float64
	limit, offset := 2, index-1
	if index == 0 {
		limit, offset = 1, 0
	}
	if err := s.executeQuery(
		s.getManyQueries(table).getListIndexToIdxQuery(),
		&sqlitex.ExecOptions{
			Args: []any{sqlID, limit, offset},
			ResultFunc: func(stmt *sqlite.Stmt) error {
				neighbors = append(neighbors, stmt.ColumnFloat(0))
				return nil
			}}); err != nil {
		return 0.0, 0.0, err
	}
	// the nb elements are spread evenly between the idx of their neighbors
	// (0 before the first element, the previous idx + 1 after the last)
	switch {
	case index == 0 && len(neighbors) == 0:
		// empty list
		return 1.0, 1.0, nil
	case index == 0:
		// before the first element
		increment := neighbors[0] / (float64(nb) + 1)
		return increment, increment, nil
	case len(neighbors) == 0:
		return 0, 0, fmt.Errorf("invalid index in table %v for object %v : %v not in list bounds", table.name, sqlID, index)
	case len(neighbors) == 1:
		// after the last element
		return neighbors[0] + 1, 1, nil
	default:
		// between two elements
		increment := (neighbors[1] - neighbors[0]) / (float64(nb) + 1)
		return neighbors[0] + increment, increment, nil
	}
}

func (s *SQLStore) Remove(object ecore.EObject, feature ecore.EStructuralFeature, index int, needResult bool) any {
	var options []operationOption
	if needResult {
		options = append(options, withWait())
	}
	result, err := s.execute(
		newOperation("Remove", operationWrite, object, feature, index, nil, func() (any, error) {
			return s.doRemove(object, feature, index, needResult)
		}, options...))
	if needResult {
		return resultAs[any](result, err)
	}
	return nil
}

func (s *SQLStore) doRemove(object ecore.EObject, feature ecore.EStructuralFeature, index int, needResult bool) (decoded any, err error) {
	sqlObjectID, err := s.encodeSQLID(object, -1, -1)
	if err != nil {
		return nil, err
	}
	featureData, err := s.getEncoderFeatureData(object, feature)
	if err != nil {
		return nil, err
	}
	featureTable, err := s.getManyFeatureTable(featureData.schema)
	if err != nil {
		return nil, err
	}

	var value any
	if err := s.executeQuery(
		s.getManyQueries(featureTable).getRemoveQuery(),
		&sqlitex.ExecOptions{
			Args: []any{sqlObjectID, index},
			ResultFunc: func(stmt *sqlite.Stmt) error {
				value = decodeAny(stmt, 0)
				return nil
			}}); err != nil {
		return nil, err
	}

	// decode previous value
	if needResult {
		decoded, err = s.decodeFeatureValue(featureData.schema, value)
		if err != nil {
			return nil, err
		}
	}
	return
}

type moveIndexes struct {
	sourceIndex int
	targetIndex int
}

func (mi moveIndexes) String() string {
	return fmt.Sprintf("{ sourceIndex : %v, targetIndex %v}", mi.sourceIndex, mi.targetIndex)
}

func (s *SQLStore) Move(object ecore.EObject, feature ecore.EStructuralFeature, sourceIndex int, targetIndex int, needResult bool) any {
	return resultAs[any](s.execute(
		newOperation("Move", operationWrite, object, feature, -1, moveIndexes{sourceIndex, targetIndex}, func() (any, error) {
			return s.doMove(object, feature, sourceIndex, targetIndex, needResult)
		}, withWait())))
}

func (s *SQLStore) doMove(object ecore.EObject, feature ecore.EStructuralFeature, sourceIndex int, targetIndex int, needResult bool) (decoded any, err error) {
	sqlObjectID, err := s.encodeSQLID(object, -1, -1)
	if err != nil {
		return nil, err
	}
	featureSchema, err := s.getFeatureSchema(object, feature)
	if err != nil {
		return nil, err
	}
	featureTable, err := s.getManyFeatureTable(featureSchema)
	if err != nil {
		return nil, err
	}

	// compute target index
	if targetIndex > sourceIndex {
		targetIndex++
	}
	idx, _, err := s.getInsertIdx(featureTable, sqlObjectID, targetIndex, 1)
	if err != nil {
		return nil, err
	}

	// update idx of source index row with target idx
	var value any
	if err := s.executeQuery(
		s.getManyQueries(featureTable).getUpdateIdxQuery(),
		&sqlitex.ExecOptions{
			Args: []any{idx, sqlObjectID, sourceIndex},
			ResultFunc: func(stmt *sqlite.Stmt) error {
				value = decodeAny(stmt, 0)
				return nil
			}}); err != nil {
		return nil, err
	}

	// decode value
	if needResult {
		decoded, err = s.decodeFeatureValue(featureSchema, value)
		if err != nil {
			return nil, err
		}
	}

	return
}

func (s *SQLStore) Clear(object ecore.EObject, feature ecore.EStructuralFeature) {
	_, _ = s.execute(newOperation("Clear", operationWrite, object, feature, -1, nil, func() (any, error) {
		return s.doClear(object, feature)
	}))
}

func (s *SQLStore) doClear(object ecore.EObject, feature ecore.EStructuralFeature) (any, error) {
	sqlObjectID, err := s.encodeSQLID(object, -1, -1)
	if err != nil {
		return nil, err
	}
	featureTable, err := s.getFeatureTable(object, feature)
	if err != nil {
		return nil, err
	}

	if err := s.executeQuery(
		s.getManyQueries(featureTable).getClearQuery(),
		&sqlitex.ExecOptions{Args: []any{sqlObjectID}},
	); err != nil {
		return nil, err
	}
	return nil, nil
}

type containerAndFeature struct {
	container ecore.EObject
	feature   ecore.EStructuralFeature
}

func (s *SQLStore) GetContainer(object ecore.EObject) (ecore.EObject, ecore.EStructuralFeature) {
	if result := resultAs[*containerAndFeature](s.execute(
		newOperation("GetContainer", operationRead, object, nil, -1, nil, func() (any, error) {
			return s.doGetContainer(object)
		}, withWait()))); result != nil {
		return result.container, result.feature
	}
	return nil, nil
}

func (s *SQLStore) doGetContainer(object ecore.EObject) (*containerAndFeature, error) {
	sqlObjectID, isSQLObjectID, err := s.getSQLID(object)
	if err != nil {
		return nil, err
	} else if !isSQLObjectID {
		return nil, nil
	}

	containerID := int64(-1)
	containerFeatureID := int64(-1)
	if err := s.executeQuery(`SELECT containerID,containerFeatureID FROM ".objects" WHERE objectID=?`, &sqlitex.ExecOptions{
		Args: []any{sqlObjectID},
		ResultFunc: func(stmt *sqlite.Stmt) error {
			switch stmt.ColumnType(0) {
			case sqlite.TypeNull:
				containerID = 0
			case sqlite.TypeInteger:
				containerID = stmt.ColumnInt64(0)
			}
			switch stmt.ColumnType(1) {
			case sqlite.TypeNull:
				containerFeatureID = 0
			case sqlite.TypeInteger:
				containerFeatureID = stmt.ColumnInt64(1)
			}
			return nil
		},
	}); err != nil {
		return nil, err
	}

	switch containerID {
	case -1:
		return nil, fmt.Errorf("unable to find container for object '%v'", sqlObjectID)
	case 0:
		return nil, nil
	default:
		container, err := s.decodeObject(containerID)
		if err != nil {
			return nil, err
		}

		containerInternal := container.(ecore.EObjectInternal)
		containerClass := containerInternal.EClass()
		feature := containerClass.GetEStructuralFeature(containerInternal.EStaticFeatureCount() + int(containerFeatureID))
		return &containerAndFeature{container, feature}, nil
	}
}

func (s *SQLStore) SetContainer(object ecore.EObject, container ecore.EObject, feature ecore.EStructuralFeature) {
	// called from an operation being executed (an object encoded by the writer is bound to the
	// store), execute runs doSetContainer in the writer goroutine
	_, _ = s.execute(newOperation("SetContainer", operationWrite, object, feature, -1, container, func() (any, error) {
		return s.doSetContainer(object, container, feature)
	}))
}

func (s *SQLStore) doSetContainer(object ecore.EObject, container ecore.EObject, feature ecore.EStructuralFeature) (any, error) {
	sqlObjectID, err := s.encodeSQLID(object, -1, -1)
	if err != nil {
		return nil, err
	}

	var sqlContainerID any
	if container != nil {
		sqlContainerID, err = s.encodeSQLID(container, -1, -1)
		if err != nil {
			return nil, err
		}
	}

	var featureID any
	if container != nil && feature != nil {
		featureID = container.EClass().GetFeatureID(feature)
	}

	if err := s.executeQuery(`UPDATE ".objects" SET containerID=?,containerFeatureID=? WHERE objectID=?`, &sqlitex.ExecOptions{
		Args: []any{sqlContainerID, featureID, sqlObjectID},
	}); err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *SQLStore) All(object ecore.EObject, feature ecore.EStructuralFeature) iter.Seq[any] {
	return func(yield func(any) bool) {
		_, _ = s.execute(
			newOperation("All", operationRead, object, feature, -1, nil, func() (any, error) {
				interrupted := errors.New("interrupted")
				sqlObjectID, isSQLObjectID, err := s.getSQLID(object)
				if err != nil {
					return nil, err
				} else if !isSQLObjectID {
					return nil, nil
				}
				featureSchema, err := s.getFeatureSchema(object, feature)
				if err != nil {
					return nil, err
				}
				featureTable, err := s.getManyFeatureTable(featureSchema)
				if err != nil {
					return nil, err
				}
				if err := s.executeQuery(
					s.getManyQueries(featureTable).getSelectAllQuery(),
					&sqlitex.ExecOptions{
						Args: []any{sqlObjectID},
						ResultFunc: func(stmt *sqlite.Stmt) error {
							value := decodeAny(stmt, 0)
							decoded, err := s.decodeFeatureValue(featureSchema, value)
							if err != nil {
								return err
							}
							if !yield(decoded) {
								return interrupted
							}
							return nil
						}}); err != nil && err != interrupted {
					return nil, err
				}
				return nil, nil
			}, withWait()))
	}
}

func (s *SQLStore) ToArray(object ecore.EObject, feature ecore.EStructuralFeature) []any {
	return slices.Collect(s.All(object, feature))
}

// Serialize returns the content of the database, with the pending operations committed
func (s *SQLStore) Serialize(ctx context.Context) *promise.Promise[[]byte] {
	op := newOperation("Serialize", operationRead, nil, nil, -1, nil, func() (any, error) {
		return s.doSerialize(ctx)
	}, withContext(ctx), withCommitBefore(), withWait())
	return promise.NewWithPool(func(resolve func([]byte), reject func(error)) {
		result, err := s.execute(op)
		if err != nil {
			reject(err)
		} else {
			resolve(resultAs[[]byte](result, nil))
		}
	}, s.promisePool)
}

// doSerialize serializes the database with the connection bound to the goroutine
func (s *SQLStore) doSerialize(ctx context.Context) ([]byte, error) {
	// retrieve database size
	var dbSize int64
	if err := s.executeQueryTransient("SELECT page_count * page_size as size FROM pragma_page_count(), pragma_page_size();", &sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			dbSize = stmt.ColumnInt64(0)
			return nil
		},
	}); err != nil {
		return nil, err
	}

	// connection to serialize
	bound := s.boundConn()
	if bound == nil {
		return nil, errors.New("serialize: no connection bound to the goroutine")
	}
	conn := bound.conn
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// supports big databases
	if dbSize > s.maxAllocSize {
		dbPath, err := sqlTmpDB("store-serialize")
		if err != nil {
			return nil, err
		}
		defer func() {
			// remove tmp file and ignore error
			_ = os.Remove(dbPath)
		}()

		dstConn, err := sqlite.OpenConn(dbPath)
		if err != nil {
			return nil, err
		}

		if err = backupDB(dstConn, conn); err != nil {
			dstConn.Close()
			return nil, err
		}

		if err := dstConn.Close(); err != nil {
			return nil, err
		}

		return os.ReadFile(dbPath)

	} else {
		return conn.Serialize("main")
	}
}

func (s *SQLStore) ExecuteQuery(ctx context.Context, query string, opts *sqlitex.ExecOptions) error {
	var operationType operationType
	switch operation, _, _ := strings.Cut(query, " "); operation {
	case "SELECT":
		operationType = operationRead
	default:
		operationType = operationWrite
	}
	_, err := s.execute(newOperation("ExecuteQuery", operationType, nil, nil, -1, nil, func() (any, error) {
		return nil, s.executeQuery(query, opts)
	}, withContext(ctx), withWait()))
	return err
}

// Close executes and commits the pending operations, and closes the connections.
// It returns the errors of the operations since the last flush, ErrSQLStoreFailed if the
// store has failed, and the errors of closing the connections
func (s *SQLStore) Close() error {
	if !s.isClosed.Swap(true) {
		s.loggerOperations.Debug("Close", zap.Int64("goid", goid.Get()))

		flushErr := s.Flush(context.Background())

		// stop the writer once the operations enqueued meanwhile are executed
		s.mutexWriter.Lock()
		s.closed = true
		s.mutexWriter.Unlock()
		select {
		case s.writerSignal <- struct{}{}:
		default:
		}
		<-s.writerDone

		// the last commit of the writer may have failed the store
		if failure := s.Err(); failure != nil && !errors.Is(flushErr, ErrSQLStoreFailed) {
			flushErr = errors.Join(flushErr, failure)
		}

		// close the writer connection first : in memory mode, it backs the database up
		writerErr := s.writerConnClose(s.writerConn)
		poolErr := s.connPoolClose(s.connPool)
		s.antsPool.Release()
		return errors.Join(flushErr, writerErr, poolErr)
	}
	return nil
}
