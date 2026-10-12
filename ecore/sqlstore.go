package ecore

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/chebyrash/promise"
	"github.com/panjf2000/ants/v2"
	"github.com/petermattis/goid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

// log queries

type sqlSingleQueries struct {
	column      *sqlColumn
	updateQuery string
	selectQuery string
	removeQuery string
}

func (ss *sqlSingleQueries) getUpdateQuery() string {
	if len(ss.updateQuery) == 0 {
		// query
		table := ss.column.table
		var query strings.Builder
		query.WriteString("UPDATE ")
		query.WriteString(sqlEscapeIdentifier(table.name))
		query.WriteString(" SET ")
		query.WriteString(sqlEscapeIdentifier(ss.column.columnName))
		query.WriteString("=? WHERE ")
		query.WriteString(table.keyName())
		query.WriteString("=?")
		ss.updateQuery = query.String()
	}
	return ss.updateQuery
}

func (ss *sqlSingleQueries) getSelectQuery() string {
	if len(ss.selectQuery) == 0 {
		table := ss.column.table
		// query
		var query strings.Builder
		query.WriteString("SELECT ")
		query.WriteString(sqlEscapeIdentifier(ss.column.columnName))
		query.WriteString(" FROM ")
		query.WriteString(sqlEscapeIdentifier(table.name))
		query.WriteString(" WHERE ")
		query.WriteString(table.keyName())
		query.WriteString("=?")
		ss.selectQuery = query.String()
	}
	return ss.selectQuery
}

func (ss *sqlSingleQueries) getRemoveQuery() string {
	if len(ss.removeQuery) == 0 {
		// query
		table := ss.column.table
		var query strings.Builder
		query.WriteString("DELETE FROM ")
		query.WriteString(sqlEscapeIdentifier(table.name))
		query.WriteString(" WHERE ")
		query.WriteString(sqlEscapeIdentifier(ss.column.columnName))
		query.WriteString("=?")
		ss.removeQuery = query.String()
	}
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

func (ss *sqlManyQueries) getInsertQuery() string {
	if len(ss.insertQuery) == 0 {
		ss.insertQuery = ss.table.insertQuery()
	}
	return ss.insertQuery
}

func (ss *sqlManyQueries) getUpdateValueQuery() string {
	if len(ss.updateValueQuery) == 0 {
		column := ss.table.columns[len(ss.table.columns)-1]
		var query strings.Builder
		query.WriteString("UPDATE ")
		query.WriteString(sqlEscapeIdentifier(ss.table.name))
		query.WriteString(" SET ")
		query.WriteString(sqlEscapeIdentifier(column.columnName))
		query.WriteString("=? WHERE rowid IN (SELECT rowid FROM ")
		query.WriteString(sqlEscapeIdentifier(ss.table.name))
		query.WriteString(" WHERE ")
		query.WriteString(ss.table.keyName())
		query.WriteString("=?")
		query.WriteString(" ORDER BY ")
		query.WriteString(ss.table.keyName())
		query.WriteString(" ASC, idx ASC LIMIT 1 OFFSET ?)")
		ss.updateValueQuery = query.String()
	}
	return ss.updateValueQuery
}

func (ss *sqlManyQueries) getUpdateIdxQuery() string {
	if len(ss.updateIdxQuery) == 0 {
		column := ss.table.columns[len(ss.table.columns)-1]
		var query strings.Builder
		query.WriteString("UPDATE ")
		query.WriteString(sqlEscapeIdentifier(ss.table.name))
		query.WriteString(" SET idx=? WHERE rowid IN (SELECT rowid FROM ")
		query.WriteString(sqlEscapeIdentifier(ss.table.name))
		query.WriteString(" WHERE ")
		query.WriteString(ss.table.keyName())
		query.WriteString("=?")
		query.WriteString(" ORDER BY ")
		query.WriteString(ss.table.keyName())
		query.WriteString(" ASC, idx ASC LIMIT 1 OFFSET ?) RETURNING ")
		query.WriteString(sqlEscapeIdentifier(column.columnName))
		ss.updateIdxQuery = query.String()
	}
	return ss.updateIdxQuery
}

func (ss *sqlManyQueries) getSelectOneQuery() string {
	if len(ss.selectOneQuery) == 0 {
		column := ss.table.columns[len(ss.table.columns)-1]
		var query strings.Builder
		query.WriteString("SELECT ")
		query.WriteString(sqlEscapeIdentifier(column.columnName))
		query.WriteString(" FROM ")
		query.WriteString(sqlEscapeIdentifier(ss.table.name))
		query.WriteString(" WHERE ")
		query.WriteString(ss.table.keyName())
		query.WriteString("=? ORDER BY ")
		query.WriteString(ss.table.keyName())
		query.WriteString(" ASC, idx ASC LIMIT 1 OFFSET ?")
		ss.selectOneQuery = query.String()
	}
	return ss.selectOneQuery
}

func (ss *sqlManyQueries) getSelectAllQuery() string {
	if len(ss.selectAllQuery) == 0 {
		column := ss.table.columns[len(ss.table.columns)-1]
		var query strings.Builder
		query.WriteString("SELECT ")
		query.WriteString(sqlEscapeIdentifier(column.columnName))
		query.WriteString(" FROM ")
		query.WriteString(sqlEscapeIdentifier(ss.table.name))
		query.WriteString(" WHERE ")
		query.WriteString(ss.table.keyName())
		query.WriteString("=? ORDER BY ")
		query.WriteString(ss.table.keyName())
		query.WriteString(" ASC, idx ASC")
		ss.selectAllQuery = query.String()
	}
	return ss.selectAllQuery
}

func (ss *sqlManyQueries) getExistsQuery() string {
	if len(ss.existsQuery) == 0 {
		var query strings.Builder
		query.WriteString("SELECT EXISTS(SELECT 1 FROM ")
		query.WriteString(sqlEscapeIdentifier(ss.table.name))
		query.WriteString(" WHERE ")
		query.WriteString(ss.table.keyName())
		query.WriteString("=?)")
		ss.existsQuery = query.String()
	}
	return ss.existsQuery
}

func (ss *sqlManyQueries) getClearQuery() string {
	if len(ss.clearQuery) == 0 {
		var query strings.Builder
		query.WriteString("DELETE FROM ")
		query.WriteString(sqlEscapeIdentifier(ss.table.name))
		query.WriteString(" WHERE ")
		query.WriteString(ss.table.keyName())
		query.WriteString("=?")
		ss.clearQuery = query.String()
	}
	return ss.clearQuery
}

func (ss *sqlManyQueries) getCountQuery() string {
	if len(ss.countQuery) == 0 {
		var query strings.Builder
		query.WriteString("SELECT COUNT(*) FROM ")
		query.WriteString(sqlEscapeIdentifier(ss.table.name))
		query.WriteString(" WHERE ")
		query.WriteString(ss.table.keyName())
		query.WriteString("=?")
		ss.countQuery = query.String()
	}
	return ss.countQuery
}

func (ss *sqlManyQueries) getContainsQuery() string {
	if len(ss.containsQuery) == 0 {
		column := ss.table.columns[len(ss.table.columns)-1]
		var query strings.Builder
		query.WriteString("SELECT rowid FROM ")
		query.WriteString(sqlEscapeIdentifier(ss.table.name))
		query.WriteString(" WHERE ")
		query.WriteString(ss.table.keyName())
		query.WriteString("=? AND ")
		query.WriteString(sqlEscapeIdentifier(column.columnName))
		query.WriteString("=?")
		ss.containsQuery = query.String()
	}
	return ss.containsQuery
}

func (ss *sqlManyQueries) getIndexOfQuery() string {
	if len(ss.indexOfQuery) == 0 {
		column := ss.table.columns[len(ss.table.columns)-1]
		var query strings.Builder
		query.WriteString("SELECT idx FROM ")
		query.WriteString(sqlEscapeIdentifier(ss.table.name))
		query.WriteString(" WHERE ")
		query.WriteString(ss.table.keyName())
		query.WriteString("=? AND ")
		query.WriteString(sqlEscapeIdentifier(column.columnName))
		query.WriteString("=? ORDER BY idx ASC LIMIT 1")
		ss.indexOfQuery = query.String()
	}
	return ss.indexOfQuery
}

func (ss *sqlManyQueries) getLastIndexOfQuery() string {
	if len(ss.lastIndexOfQuery) == 0 {
		column := ss.table.columns[len(ss.table.columns)-1]
		var query strings.Builder
		query.WriteString("SELECT idx FROM ")
		query.WriteString(sqlEscapeIdentifier(ss.table.name))
		query.WriteString(" WHERE ")
		query.WriteString(ss.table.keyName())
		query.WriteString("=? AND ")
		query.WriteString(sqlEscapeIdentifier(column.columnName))
		query.WriteString("=? ORDER BY idx DESC LIMIT 1")
		ss.lastIndexOfQuery = query.String()
	}
	return ss.lastIndexOfQuery
}

func (ss *sqlManyQueries) getIdxToListIndexQuery() string {
	if len(ss.idxToListIndex) == 0 {
		var query strings.Builder
		query.WriteString("SELECT COUNT(*) FROM ")
		query.WriteString(sqlEscapeIdentifier(ss.table.name))
		query.WriteString(" WHERE ")
		query.WriteString(ss.table.keyName())
		query.WriteString("=? AND idx<?")
		ss.idxToListIndex = query.String()
	}
	return ss.idxToListIndex
}

func (ss *sqlManyQueries) getListIndexToIdxQuery() string {
	if len(ss.listIndexToIdx) == 0 {
		var query strings.Builder
		query.WriteString("SELECT idx FROM ")
		query.WriteString(sqlEscapeIdentifier(ss.table.name))
		query.WriteString(" WHERE ")
		query.WriteString(ss.table.keyName())
		query.WriteString("=? ORDER BY idx ASC LIMIT ? OFFSET ?")
		ss.listIndexToIdx = query.String()
	}
	return ss.listIndexToIdx
}

func (ss *sqlManyQueries) getRemoveQuery() string {
	if len(ss.removeQuery) == 0 {
		column := ss.table.columns[len(ss.table.columns)-1]
		var query strings.Builder
		query.WriteString("DELETE FROM ")
		query.WriteString(sqlEscapeIdentifier(ss.table.name))
		query.WriteString(" WHERE rowid IN (SELECT rowid FROM ")
		query.WriteString(sqlEscapeIdentifier(ss.table.name))
		query.WriteString(" WHERE ")
		query.WriteString(ss.table.keyName())
		query.WriteString("=?")
		query.WriteString(" ORDER BY ")
		query.WriteString(ss.table.keyName())
		query.WriteString(" ASC, idx ASC LIMIT 1 OFFSET ?) RETURNING ")
		query.WriteString(sqlEscapeIdentifier(column.columnName))
		ss.removeQuery = query.String()
	}
	return ss.removeQuery
}

type sqlStoreObjectManager struct {
	store *SQLStore
}

func newSQLStoreObjectManager() *sqlStoreObjectManager {
	return &sqlStoreObjectManager{}
}

func (r *sqlStoreObjectManager) registerObject(o EObject) {
	// an object encoded by the open transaction is read by the writer until the commit
	r.store.registerDirty(o)
	// set store object
	if storeObject, _ := o.(EStoreProvider); storeObject != nil {
		storeObject.SetEStore(r.store)
	}
}

type SQLStoreIDManager interface {
	SQLDecoderIDManager
	SQLEncoderIDManager
	ClearObjectID(eObject EObject)
}

type sqlStoreIDManagerImpl struct {
	sqlDecoderIDManagerImpl
	sqlEncoderIDManagerImpl
	mutex sync.Mutex
}

func newSQLStoreIDManager() SQLStoreIDManager {
	return &sqlStoreIDManagerImpl{
		sqlDecoderIDManagerImpl: sqlDecoderIDManagerImpl{
			packages:     map[int64]EPackage{},
			objects:      map[int64]EObject{},
			classes:      map[int64]EClass{},
			enumLiterals: map[int64]EEnumLiteral{},
		},
		sqlEncoderIDManagerImpl: sqlEncoderIDManagerImpl{
			packages:     map[EPackage]int64{},
			objects:      map[EObject]int64{},
			classes:      map[EClass]int64{},
			enumLiterals: map[EEnumLiteral]int64{},
		},
	}
}

func (r *sqlStoreIDManagerImpl) SetPackageID(p EPackage, id int64) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.sqlDecoderIDManagerImpl.packages[id] = p
	r.sqlEncoderIDManagerImpl.packages[p] = id
}

func (r *sqlStoreIDManagerImpl) GetPackageID(p EPackage) (int64, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.sqlEncoderIDManagerImpl.GetPackageID(p)
}

func (r *sqlStoreIDManagerImpl) GetPackageFromID(id int64) (EPackage, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.sqlDecoderIDManagerImpl.GetPackageFromID(id)
}

func (r *sqlStoreIDManagerImpl) SetClassID(c EClass, id int64) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.sqlDecoderIDManagerImpl.classes[id] = c
	r.sqlEncoderIDManagerImpl.classes[c] = id
}

func (r *sqlStoreIDManagerImpl) GetClassID(c EClass) (int64, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.sqlEncoderIDManagerImpl.GetClassID(c)
}

func (r *sqlStoreIDManagerImpl) GetClassFromID(id int64) (EClass, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.sqlDecoderIDManagerImpl.GetClassFromID(id)
}

func (r *sqlStoreIDManagerImpl) GetObjectFromID(id int64) (o EObject, b bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.sqlDecoderIDManagerImpl.GetObjectFromID(id)
}

func (r *sqlStoreIDManagerImpl) GetObjectID(o EObject) (id int64, b bool) {
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

func (r *sqlStoreIDManagerImpl) SetObjectID(o EObject, id int64) {
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

func (r *sqlStoreIDManagerImpl) ClearObjectID(o EObject) {
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

func (r *sqlStoreIDManagerImpl) SetEnumLiteralID(e EEnumLiteral, id int64) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.sqlDecoderIDManagerImpl.enumLiterals[id] = e
	r.sqlEncoderIDManagerImpl.enumLiterals[e] = id
}

func (r *sqlStoreIDManagerImpl) GetEnumLiteralID(e EEnumLiteral) (int64, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.sqlEncoderIDManagerImpl.GetEnumLiteralID(e)
}

func (r *sqlStoreIDManagerImpl) GetEnumLiteralFromID(id int64) (EEnumLiteral, bool) {
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
	object       EObject
	feature      EStructuralFeature
	index        int
	value        any
	fn           func() (any, error)
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

func newOperation(
	cmd string,
	type_ operationType,
	object EObject,
	feature EStructuralFeature,
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
	case EObject:
		e.AddString("value", fmt.Sprintf("%s(%p)", v.EClass().GetName(), v))
	default:
		if err := e.AddReflected("value", op.value); err != nil {
			return err
		}
	}
	return nil
}

type objectFeature struct {
	object  EObject
	feature EStructuralFeature
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
	dirty               map[EObject]struct{}         // objects of the pending operations and of the open transaction
	dirtyGlobal         bool                         // the pending operations or the open transaction touch the roots or arbitrary tables
	pending             map[objectFeature]*operation // last pending Set or UnSet by object and feature
	errs                []error                      // errors of the operations since the last flush
}

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
	resourceURI *URI,
	idManager EObjectIDManager,
	packageRegistry EPackageRegistry,
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
	resourceURI *URI,
	idManager EObjectIDManager,
	packageRegistry EPackageRegistry,
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
	promisePool := promise.FromAntsPool(antsPool)

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
			classDataMap:     map[EClass]*sqlDecoderClassData{},
		},
		sqlEncoder: sqlEncoder{
			sqlBase:          base,
			isForced:         false,
			isKeepDefaults:   isKeepDefaults,
			classDataMap:     map[EClass]*sqlEncoderClassData{},
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
		dirty:           map[EObject]struct{}{},
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
	stmts := s.singleQueries[column]
	if stmts == nil {
		stmts = &sqlSingleQueries{
			column: column,
		}
		s.singleQueries[column] = stmts
	}
	s.mutexQueries.Unlock()
	return stmts
}

func (s *SQLStore) getManyQueries(table *sqlTable) *sqlManyQueries {
	s.mutexQueries.Lock()
	stmts := s.manyQueries[table]
	if stmts == nil {
		stmts = &sqlManyQueries{
			table: table,
		}
		s.manyQueries[table] = stmts
	}
	s.mutexQueries.Unlock()
	return stmts
}

func (s *SQLStore) getFeatureSchema(object EObject, feature EStructuralFeature) (*sqlFeatureSchema, error) {
	// retrieve class schema
	class := object.EClass()
	classSchema := s.sqlDecoder.schema.getClassSchema(class)
	if classSchema == nil {
		return nil, fmt.Errorf("class %s is unknown", class.GetName())
	}
	// retrieve feature schema
	featureSchema := classSchema.getFeatureSchema(feature)
	if featureSchema == nil {
		return nil, fmt.Errorf("feature %s is unknown", feature.GetName())
	}
	return featureSchema, nil
}

func (s *SQLStore) getFeatureTable(object EObject, feature EStructuralFeature) (*sqlTable, error) {
	featureSchema, err := s.getFeatureSchema(object, feature)
	if err != nil {
		return nil, err
	}
	return featureSchema.table, nil
}

func (s *SQLStore) getEncoderFeatureData(object EObject, feature EStructuralFeature) (*sqlEncoderFeatureData, error) {
	// retrieve class schema
	class := object.EClass()

	// retrieve class data
	classData, err := s.getEncoderClassData(class)
	if err != nil {
		return nil, fmt.Errorf("class %s is unknown %w", class.GetName(), err)
	}

	// retrieve feature data
	featureData, isFeatureData := classData.features.Get(feature)
	if !isFeatureData {
		return nil, fmt.Errorf("feature %s is unknown", feature.GetName())
	}

	return featureData, nil
}

func (e *SQLStore) encodeFeatureValue(sqlObjectID int64, featureData *sqlEncoderFeatureData, value any) (encoded any, err error) {
	if value != nil {
		switch featureData.schema.featureKind {
		case sfkObject, sfkObjectList:
			eObject := value.(EObject)
			return e.encodeSQLID(eObject, -1, -1)
		default:
			return e.sqlEncoder.encodeFeatureValue(sqlObjectID, featureData, value)
		}
	}
	return nil, nil
}

func (s *SQLStore) getSQLID(eObject EObject) (int64, bool, error) {
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
func (s *SQLStore) encodeSQLID(eObject EObject, sqlContainerID int64, containerFeatureID int64) (int64, error) {
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

func mapSlice[S ~[]E, E, R any](slice S, mapper func(int, E) R) []R {
	mappedSlice := make([]R, len(slice))
	for i, v := range slice {
		mappedSlice[i] = mapper(i, v)
	}
	return mappedSlice
}

const (
	// unlockOperationTimeout is the time an operation waits for the writer before a warning is logged
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
// returns the errors of the operations since the last flush
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
	}, withContext(ctx), withFlush())
	op.done = make(chan struct{})
	s.enqueue(op)
	_, err := s.await(op)
	return err
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
	if value, isObject := op.value.(EObject); isObject {
		s.dirty[value] = struct{}{}
	}
}

// registerDirty registers an object encoded or decoded by the open transaction
func (s *SQLStore) registerDirty(object EObject) {
	if goid.Get() == s.writerGoid.Load() {
		s.mutexWriter.Lock()
		s.dirty[object] = struct{}{}
		s.mutexWriter.Unlock()
	}
}

// enqueue gives an operation to the writer
func (s *SQLStore) enqueue(op *operation) {
	s.mutexWriter.Lock()
	if s.closed {
		s.mutexWriter.Unlock()
		op.err = errors.New("store is closed")
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
func (s *SQLStore) pendingOperation(object EObject, feature EStructuralFeature, index int) *operation {
	s.mutexWriter.Lock()
	defer s.mutexWriter.Unlock()
	if op := s.pending[objectFeature{object, feature}]; op != nil && op.index == index {
		return op
	}
	return nil
}

// await waits for an enqueued operation to be executed
func (s *SQLStore) await(op *operation) (any, error) {
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
		s.loggerOperations.Warn("waiting", zap.Int64("goid", goid.Get()), zap.Object("operation", &operationMarshaler{op: op}), zap.Duration("timeout", s.timeoutOperation))
	}
	select {
	case <-op.done:
		return op.result, op.err
	case <-ctx.Done():
		return nil, ctx.Err()
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
	if s.closed || s.dirtyGlobal {
		return false
	}
	if op.object == nil {
		// the whole store : only when idle
		return !s.busy && !s.txOpen && len(s.queue) == 0
	}
	if _, isDirty := s.dirty[op.object]; isDirty {
		return false
	}
	if value, isObject := op.value.(EObject); isObject {
		if _, isDirty := s.dirty[value]; isDirty {
			return false
		}
	}
	return true
}

// execute executes an operation :
// - from the writer goroutine (an operation calling the store), now on the writer connection
// - a read independent of the pending writes, now on a pool connection
// - otherwise by the writer, waiting for its result if wait
func (s *SQLStore) execute(op *operation, wait bool) (any, error) {
	if goid.Get() == s.writerGoid.Load() {
		return op.fn()
	}
	if op.type_ == operationRead && s.canReadConcurrently(op) {
		return s.executeConcurrently(op)
	}
	if wait {
		op.done = make(chan struct{})
	}
	s.enqueue(op)
	if !wait {
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
	if op.ctx != nil && op.ctx.Err() != nil {
		// nobody waits for the result anymore
		op.err = op.ctx.Err()
		return
	}
	if op.commitBefore {
		s.commit()
	}
	if op.flush {
		// collects the errors : no transaction needed
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
	if op.type_ == operationWrite && op.feature != nil {
		s.mutexWriter.Lock()
		if key := (objectFeature{op.object, op.feature}); s.pending[key] == op {
			delete(s.pending, key)
		}
		s.mutexWriter.Unlock()
	}
	if op.err != nil {
		s.loggerOperations.Error("error", zap.Object("operation", &operationMarshaler{op: op}), zap.Error(op.err))
		s.addError(op.err)
		// reset the statements left stepped by a panic
		s.writerConn.SetInterrupt(nil)
		// some errors roll the transaction back
		if s.txOpen && s.writerConn.AutocommitEnabled() {
			s.loggerOperations.Error("transaction rolled back", zap.Object("operation", &operationMarshaler{op: op}), zap.Error(op.err))
			s.mutexWriter.Lock()
			s.txOpen = false
			s.txOperations = 0
			s.mutexWriter.Unlock()
		}
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
	} else if e := s.loggerOperations.Check(zap.DebugLevel, "commit"); e != nil {
		e.Write(zap.Int("operations", s.txOperations), zap.Duration("duration", time.Since(s.txStart)))
	}
	s.mutexWriter.Lock()
	s.txOpen = false
	s.txOperations = 0
	s.dirty = map[EObject]struct{}{}
	s.dirtyGlobal = false
	for _, op := range s.queue {
		s.markDirty(op)
	}
	s.mutexWriter.Unlock()
}

func (s *SQLStore) Get(object EObject, feature EStructuralFeature, index int) any {
	// optimization : if we have a non many feature, check if we have a pending write operation
	if !feature.IsMany() {
		if op := s.pendingOperation(object, feature, index); op != nil {
			return op.value
		}
	}
	return resultAs[any](s.execute(
		newOperation("Get", operationRead, object, feature, index, nil, func() (any, error) {
			return s.doGet(object, feature, index)
		}), true))
}

func (s *SQLStore) doGet(object EObject, feature EStructuralFeature, index int) (any, error) {
	// retrieve value from sqlite db
	sqlObjectID, isSQLObjectID, err := s.getSQLID(object)
	if err != nil {
		return nil, err
	} else if !isSQLObjectID {
		return nil, nil
	}
	featureSchema, err := s.getFeatureSchema(object, feature)
	if err != nil {
		return err, nil
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
		return err, nil
	}

	return s.decodeFeatureValue(featureSchema, value)
}

func (s *SQLStore) Set(object EObject, feature EStructuralFeature, index int, value any, needResult bool) (oldValue any) {
	result, err := s.execute(
		newOperation("Set", operationWrite, object, feature, index, value, func() (any, error) {
			return s.doSet(object, feature, index, value, needResult)
		}), needResult)
	if needResult {
		oldValue = resultAs[any](result, err)
	}
	return
}

func (s *SQLStore) doSet(object EObject, feature EStructuralFeature, index int, value any, needResult bool) (oldValue any, err error) {
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

func (s *SQLStore) IsSet(object EObject, feature EStructuralFeature) bool {
	// optimization : if we have a non many feature, check if we have a pending write operation
	if !feature.IsMany() {
		if op := s.pendingOperation(object, feature, -1); op != nil {
			return op.value != nil
		}
	}
	return resultAs[bool](s.execute(
		newOperation("IsSet", operationRead, object, feature, -1, nil, func() (any, error) {
			return s.doIsSet(object, feature)
		}), true))
}

func (s *SQLStore) doIsSet(object EObject, feature EStructuralFeature) (bool, error) {
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
		return value != featureSchema.feature.GetDefaultValue(), nil
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

func (s *SQLStore) UnSet(object EObject, feature EStructuralFeature) {
	_, _ = s.execute(
		newOperation("UnSet", operationWrite, object, feature, -1, nil, func() (any, error) {
			return s.doUnSet(object, feature)
		}), true)
}

func (s *SQLStore) doUnSet(object EObject, feature EStructuralFeature) (any, error) {
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
		query = s.getSingleQueries(featureColumn).getUpdateQuery()
		args = []any{feature.GetDefaultValue(), sqlObjectID}
	} else if featureTable := featureData.schema.table; featureTable != nil {
		query = s.getManyQueries(featureTable).getClearQuery()
		args = []any{sqlObjectID}
	}
	if err := s.executeQuery(query, &sqlitex.ExecOptions{Args: args}); err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *SQLStore) IsEmpty(object EObject, feature EStructuralFeature) bool {
	return resultAs[bool](s.execute(
		newOperation("IsEmpty", operationRead, object, feature, -1, nil, func() (any, error) {
			return s.doIsEmpty(object, feature)
		}), true))
}

func (s *SQLStore) doIsEmpty(object EObject, feature EStructuralFeature) (bool, error) {
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
		return true, nil
	}
	return isEmpty, nil

}

func (s *SQLStore) Size(object EObject, feature EStructuralFeature) int {
	return resultAs[int](s.execute(
		newOperation("Size", operationRead, object, feature, -1, nil, func() (any, error) {
			return s.doSize(object, feature)
		}), true))
}

func (s *SQLStore) doSize(object EObject, feature EStructuralFeature) (int, error) {
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

func (s *SQLStore) Contains(object EObject, feature EStructuralFeature, value any) bool {
	return resultAs[bool](s.execute(
		newOperation("Contains", operationRead, object, feature, -1, value, func() (any, error) {
			return s.doContains(object, feature, value)
		}), true))
}

// encodeLookupValue encodes a value to look it up in a feature table, without writing :
// an object or an enum literal unknown to the store is not in the table, and isEncoded is false
func (s *SQLStore) encodeLookupValue(featureSchema *sqlFeatureSchema, value any) (encoded any, isEncoded bool, err error) {
	if value == nil {
		return nil, true, nil
	}
	switch featureSchema.featureKind {
	case sfkObject, sfkObjectList:
		sqlID, isSQLID, err := s.getSQLID(value.(EObject))
		if err != nil || !isSQLID {
			return nil, false, err
		}
		return sqlID, true, nil
	case sfkObjectReference, sfkObjectReferenceList:
		eObject := value.(EObject)
		if sqlID, _ := s.sqlIDManager.GetObjectID(eObject); sqlID != 0 {
			return strconv.FormatInt(sqlID, 10), true, nil
		}
		return s.uri.Relativize(GetURI(eObject)).String(), true, nil
	case sfkEnum:
		eAttribute := featureSchema.feature.(EAttribute)
		eEnum := eAttribute.GetEAttributeType().(EEnum)
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
		eDataType := featureSchema.feature.(EAttribute).GetEAttributeType()
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

func (s *SQLStore) doContains(object EObject, feature EStructuralFeature, value any) (bool, error) {
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

func (s *SQLStore) doIndexOf(object EObject, feature EStructuralFeature, value any, getIndexOfQuery func(*sqlManyQueries) string) (int, error) {
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

func (s *SQLStore) IndexOf(object EObject, feature EStructuralFeature, value any) int {
	return resultAs[int](s.execute(
		newOperation("IndexOf", operationRead, object, feature, -1, value, func() (any, error) {
			return s.doIndexOf(object, feature, value, func(sms *sqlManyQueries) string {
				return sms.getIndexOfQuery()
			})
		}), true))
}

func (s *SQLStore) LastIndexOf(object EObject, feature EStructuralFeature, value any) int {
	return resultAs[int](s.execute(
		newOperation("LastIndexOf", operationRead, object, feature, -1, value, func() (any, error) {
			return s.doIndexOf(object, feature, value, func(sms *sqlManyQueries) string {
				return sms.getLastIndexOfQuery()
			})
		}), true))
}

// AddRoot add object as store root
func (s *SQLStore) AddRoot(object EObject) {
	_, _ = s.execute(newOperation("AddRoot", operationWrite, nil, nil, -1, nil, func() (any, error) {
		return s.doAddRoot(object)
	}), false)
}

func (s *SQLStore) doAddRoot(object EObject) (any, error) {
	return nil, s.encodeContent(object)
}

// RemoveRoot implements EStore.
func (s *SQLStore) RemoveRoot(object EObject) {
	_, _ = s.execute(newOperation("RemoveRoot", operationWrite, nil, nil, -1, nil, func() (any, error) {
		return s.doRemoveRoot(object)
	}), false)
}

func (s *SQLStore) doRemoveRoot(object EObject) (any, error) {
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
func (s *SQLStore) GetRoots() []EObject {
	return resultAs[[]EObject](s.execute(
		newOperation("GetRoots", operationRead, nil, nil, -1, nil, func() (any, error) {
			return s.doGetRoots()
		}), true))
}

func (s *SQLStore) doGetRoots() ([]EObject, error) {
	table := s.schema.contentsTable
	contents := []EObject{}
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

func postOrderTraverse(object EObject, yield func(EObject) bool) bool {
	for child := range object.EContents().All() {
		if !postOrderTraverse(child.(EObject), yield) {
			return false
		}
	}
	return yield(object)
}

func postOrderIterator(object EObject) iter.Seq[EObject] {
	return func(yield func(EObject) bool) {
		postOrderTraverse(object, yield)
	}
}

// UnRegister forgets the ids of object, and of its contents if withContents, once the
// operations before it are executed : the objects are encoded again when written again
func (s *SQLStore) UnRegister(object EObject, withContents bool) *promise.Promise[any] {
	op := newOperation("UnRegister", operationWrite, object, nil, -1, nil, func() (any, error) {
		if withContents {
			// we use a post order iterator to avoid registering again
			// parent
			for object := range postOrderIterator(object) {
				s.sqlIDManager.ClearObjectID(object)
			}
		} else {
			s.sqlIDManager.ClearObjectID(object)
		}
		return nil, nil
	})
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
	op.done = make(chan struct{})
	s.enqueue(op)
	return promise.NewWithPool(func(resolve func(any), reject func(error)) {
		<-op.done
		if op.err != nil {
			reject(op.err)
		} else {
			resolve(op.result)
		}
	}, s.promisePool)
}

func (s *SQLStore) Add(object EObject, feature EStructuralFeature, index int, value any) {
	_, _ = s.execute(
		newOperation("Add", operationWrite, object, feature, index, value, func() (any, error) {
			return s.doAdd(object, feature, index, value)
		}), false)
}

func (s *SQLStore) doAdd(object EObject, feature EStructuralFeature, index int, value any) (any, error) {
	sqlObjectID, err := s.encodeSQLID(object, -1, -1)
	if err != nil {
		return nil, err
	}
	featureData, err := s.getEncoderFeatureData(object, feature)
	if err != nil {
		return nil, err
	}
	featureTable := featureData.schema.table
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

func (s *SQLStore) AddAll(object EObject, feature EStructuralFeature, index int, c Collection) {
	_, _ = s.execute(
		newOperation("AddAll", operationWrite, object, feature, index, c, func() (any, error) {
			return s.doAddAll(object, feature, index, c)
		}), false)
}

func (s *SQLStore) doAddAll(object EObject, feature EStructuralFeature, index int, c Collection) (any, error) {
	sqlObjectID, err := s.encodeSQLID(object, -1, -1)
	if err != nil {
		return nil, err
	}
	featureData, err := s.getEncoderFeatureData(object, feature)
	if err != nil {
		return nil, err
	}
	featureTable := featureData.schema.table
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
	if index == 0 {
		// first row in the list
		idx := 1.0
		withElements := false
		if err := s.executeQuery(
			s.getManyQueries(table).getListIndexToIdxQuery(),
			&sqlitex.ExecOptions{
				Args: []any{sqlID, 1, 0},
				ResultFunc: func(stmt *sqlite.Stmt) error {
					withElements = true
					idx = stmt.ColumnFloat(0)
					return nil
				}}); err != nil {
			return 0.0, 0.0, err
		}
		if withElements {
			increment := 1.0 / (float64(nb) + 1)
			return idx * increment, increment, nil
		} else {
			return 1.0, 1.0, nil
		}
	} else {
		count := 0
		idx := 0.0
		if err := s.executeQuery(
			s.getManyQueries(table).getListIndexToIdxQuery(),
			&sqlitex.ExecOptions{
				Args: []any{sqlID, 2, index - 1},
				ResultFunc: func(stmt *sqlite.Stmt) error {
					idx += stmt.ColumnFloat(0)
					count++
					return nil
				}}); err != nil {
			return 0.0, 0.0, err
		}
		switch count {
		case 0:
			return 0, 0, fmt.Errorf("invalid index in table %v for object %v : %v not in list bounds", table.name, sqlID, index)
		case 1:
			// at the end
			return idx + 1, 1, nil
		default:
			// before
			increment := 1.0 / (float64(nb) + 1)
			return idx * increment, increment, nil
		}
	}
}

func (s *SQLStore) Remove(object EObject, feature EStructuralFeature, index int, needResult bool) any {
	result, err := s.execute(
		newOperation("Remove", operationWrite, object, feature, index, nil, func() (any, error) {
			return s.doRemove(object, feature, index, needResult)
		}), needResult)
	if needResult {
		return resultAs[any](result, err)
	}
	return nil
}

func (s *SQLStore) doRemove(object EObject, feature EStructuralFeature, index int, needResult bool) (decoded any, err error) {
	sqlObjectID, err := s.encodeSQLID(object, -1, -1)
	if err != nil {
		return nil, err
	}
	featureData, err := s.getEncoderFeatureData(object, feature)
	if err != nil {
		return nil, err
	}
	featureTable := featureData.schema.table

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

func (s *SQLStore) Move(object EObject, feature EStructuralFeature, sourceIndex int, targetIndex int, needResult bool) any {
	return resultAs[any](s.execute(
		newOperation("Move", operationWrite, object, feature, -1, moveIndexes{sourceIndex, targetIndex}, func() (any, error) {
			return s.doMove(object, feature, sourceIndex, targetIndex, needResult)
		}), true))
}

func (s *SQLStore) doMove(object EObject, feature EStructuralFeature, sourceIndex int, targetIndex int, needResult bool) (decoded any, err error) {
	sqlObjectID, err := s.encodeSQLID(object, -1, -1)
	if err != nil {
		return nil, err
	}
	featureSchema, err := s.getFeatureSchema(object, feature)
	if err != nil {
		return nil, err
	}
	featureTable := featureSchema.table

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

func (s *SQLStore) Clear(object EObject, feature EStructuralFeature) {
	_, _ = s.execute(newOperation("Clear", operationWrite, object, feature, -1, nil, func() (any, error) {
		return s.doClear(object, feature)
	}), false)
}

func (s *SQLStore) doClear(object EObject, feature EStructuralFeature) (any, error) {
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
	container EObject
	feature   EStructuralFeature
}

func (s *SQLStore) GetContainer(object EObject) (EObject, EStructuralFeature) {
	if result := resultAs[*containerAndFeature](s.execute(
		newOperation("GetContainer", operationRead, object, nil, -1, nil, func() (any, error) {
			return s.doGetContainer(object)
		}), true)); result != nil {
		return result.container, result.feature
	}
	return nil, nil
}

func (s *SQLStore) doGetContainer(object EObject) (*containerAndFeature, error) {
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

		containerInternal := container.(EObjectInternal)
		containerClass := containerInternal.EClass()
		feature := containerClass.GetEStructuralFeature(containerInternal.EStaticFeatureCount() + int(containerFeatureID))
		return &containerAndFeature{container, feature}, nil
	}
}

func (s *SQLStore) SetContainer(object EObject, container EObject, feature EStructuralFeature) {
	// called from an operation being executed (an object encoded by the writer is bound to the
	// store), execute runs doSetContainer in the writer goroutine
	_, _ = s.execute(newOperation("SetContainer", operationWrite, object, feature, -1, container, func() (any, error) {
		return s.doSetContainer(object, container, feature)
	}), false)
}

func (s *SQLStore) doSetContainer(object EObject, container EObject, feature EStructuralFeature) (any, error) {
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

func (s *SQLStore) All(object EObject, feature EStructuralFeature) iter.Seq[any] {
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
				featureTable := featureSchema.table
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
			}), true)
	}
}

func (s *SQLStore) ToArray(object EObject, feature EStructuralFeature) []any {
	return slices.Collect(s.All(object, feature))
}

// Serialize returns the content of the database, with the pending operations committed
func (s *SQLStore) Serialize(ctx context.Context) *promise.Promise[[]byte] {
	op := newOperation("Serialize", operationRead, nil, nil, -1, nil, func() (any, error) {
		return s.doSerialize(ctx)
	}, withContext(ctx), withCommitBefore())
	return promise.NewWithPool(func(resolve func([]byte), reject func(error)) {
		result, err := s.execute(op, true)
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
	}, withContext(ctx)), true)
	return err
}

// Close executes and commits the pending operations, and closes the connections
func (s *SQLStore) Close() error {
	if !s.isClosed.Swap(true) {
		s.loggerOperations.Debug("Close", zap.Int64("goid", goid.Get()))

		// error is an an operation error and is handled with the logger
		// ignore it
		_ = s.Flush(context.Background())

		// stop the writer once the operations enqueued meanwhile are executed
		s.mutexWriter.Lock()
		s.closed = true
		s.mutexWriter.Unlock()
		select {
		case s.writerSignal <- struct{}{}:
		default:
		}
		<-s.writerDone

		// close the writer connection first : in memory mode, it backs the database up
		if err := s.writerConnClose(s.writerConn); err != nil {
			return err
		}

		if err := s.connPoolClose(s.connPool); err != nil {
			return err
		}

		s.antsPool.Release()
	}
	return nil
}
