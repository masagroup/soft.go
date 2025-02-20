package ecore

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/chebyrash/promise"
	"github.com/panjf2000/ants/v2"
	"github.com/petermattis/goid"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
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
		query.WriteString("SELECT 1 FROM ")
		query.WriteString(sqlEscapeIdentifier(ss.table.name))
		query.WriteString(" WHERE ")
		query.WriteString(ss.table.keyName())
		query.WriteString("=?")
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
	store EStore
}

func newSQLStoreObjectManager() *sqlStoreObjectManager {
	return &sqlStoreObjectManager{}
}

func (r *sqlStoreObjectManager) registerObject(o EObject) {
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
	r.sqlDecoderIDManagerImpl.packages[id] = p
	r.sqlEncoderIDManagerImpl.packages[p] = id
}

func (r *sqlStoreIDManagerImpl) SetClassID(c EClass, id int64) {
	r.sqlDecoderIDManagerImpl.classes[id] = c
	r.sqlEncoderIDManagerImpl.classes[c] = id
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
	r.sqlDecoderIDManagerImpl.enumLiterals[id] = e
	r.sqlEncoderIDManagerImpl.enumLiterals[e] = id
}

var operationID atomic.Int64

type operationType uint8

const (
	operationRead  = 1 << 0
	operationWrite = 1 << 1
)

type operation struct {
	id      int64
	type_   operationType
	cmd     string
	object  EObject
	feature EStructuralFeature
	index   int
	value   any
	fn      func() (any, error)
	promise *promise.Promise[any]
}

func newOperation(
	cmd string,
	type_ operationType,
	object EObject,
	feature EStructuralFeature,
	index int,
	value any,
	fn func() (any, error)) *operation {
	return &operation{
		id:      operationID.Add(1),
		type_:   type_,
		cmd:     cmd,
		object:  object,
		feature: feature,
		index:   index,
		value:   value,
		fn:      fn,
	}
}

type SQLStore struct {
	*sqlBase
	*taskManager
	sqlDecoder
	sqlEncoder
	errorHandler     func(error)
	sqlIDManager     SQLStoreIDManager
	singleQueries    map[*sqlColumn]*sqlSingleQueries
	manyQueries      map[*sqlTable]*sqlManyQueries
	mutexQueries     sync.Mutex
	objectOperations map[EObject]map[EStructuralFeature][]*operation
	mutexOperations  sync.Mutex
	logger           *zap.Logger
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
		return newSQLStore(
			func() (*sqlitex.Pool, error) {
				connSrc, err := sqlite.OpenConn(databasePath)
				if err != nil {
					return nil, err
				}
				defer connSrc.Close()

				// create connection pool
				name := filepath.Base(databasePath)
				dbPath := fmt.Sprintf("file:%s?mode=memory&cache=shared", name)
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
			func(connPool *sqlitex.Pool) (err error) {
				defer func() {
					// close pool
					err = connPool.Close()
				}()

				// destination connection is the store db
				connDst, err := sqlite.OpenConn(databasePath)
				if err != nil {
					return err
				}
				defer connDst.Close()

				// source connection is the memory db
				connSrc, err := connPool.Take(context.Background())
				if err != nil {
					return err
				}
				defer connPool.Put(connSrc)

				// backup src db to dst db
				if err := backupDB(connDst, connSrc); err != nil {
					return err
				}

				return nil
			},
			resourceURI, idManager, packageRegistry, options,
		)
	} else {
		return newSQLStore(
			func() (*sqlitex.Pool, error) {
				return sqlitex.NewPool(databasePath, sqlitex.PoolOptions{})
			},
			func(pool *sqlitex.Pool) error {
				return pool.Close()
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
	resourceURI *URI,
	idManager EObjectIDManager,
	packageRegistry EPackageRegistry,
	options map[string]any) (store *SQLStore, err error) {

	objectIDName := ""
	codecVersion := sqlCodecVersion
	errorHandler := func(error) {}
	sqlIDManager := newSQLStoreIDManager()
	sqlObjectManager := newSQLStoreObjectManager()
	logger := zap.NewNop()
	if options != nil {
		objectIDName, _ = options[SQL_OPTION_OBJECT_ID].(string)
		if eh, isErrorHandler := options[SQL_OPTION_ERROR_HANDLER]; isErrorHandler {
			errorHandler = eh.(func(error))
		}
		if v, isVersion := options[SQL_OPTION_CODEC_VERSION].(int64); isVersion {
			codecVersion = v
		}
		if m, isSQLIDManager := options[SQL_OPTION_SQL_ID_MANAGER].(SQLStoreIDManager); isSQLIDManager {
			sqlIDManager = m
		}
		if l, isLogger := options[SQL_OPTION_LOGGER]; isLogger {
			logger = l.(*zap.Logger)
		}
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

	// ants promise pool
	antsPool, _ := ants.NewPool(-1, ants.WithLogger(&zapLogger{logger.Named("ants")}))
	promisePool := promise.FromAntsPool(antsPool)

	// create sql base
	base := &sqlBase{
		codecVersion:     codecVersion,
		uri:              resourceURI,
		objectIDName:     objectIDName,
		objectIDManager:  idManager,
		isContainerID:    true,
		isObjectID:       len(objectIDName) > 0 && objectIDName != "objectID" && idManager != nil,
		sqliteManager:    newTaskManager(promisePool, logger.Named("sqlite")),
		logger:           logger,
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
			isObjectExists:   true,
			classDataMap:     map[EClass]*sqlEncoderClassData{},
			sqlIDManager:     sqlIDManager,
			sqlObjectManager: sqlObjectManager,
			sqlLockManager:   newSqlEncoderLockManager(),
		},
		taskManager:      newTaskManager(promisePool, logger.Named("tasks")),
		sqlIDManager:     sqlIDManager,
		errorHandler:     errorHandler,
		singleQueries:    map[*sqlColumn]*sqlSingleQueries{},
		manyQueries:      map[*sqlTable]*sqlManyQueries{},
		objectOperations: map[EObject]map[EStructuralFeature][]*operation{},
		logger:           logger,
	}

	// set store in sql object manager
	sqlObjectManager.store = store

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

	return store, nil
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
		s.errorHandler(err)
		return nil, fmt.Errorf("class %s is unknown", class.GetName())
	}

	// retrieve feature data
	featureData, isFeatureData := classData.features.Get(feature)
	if !isFeatureData {
		err := fmt.Errorf("feature %s is unknown", feature.GetName())
		s.errorHandler(err)
		return nil, err
	}

	return featureData, nil
}

func (e *SQLStore) encodeFeatureValue(sqlObjectID int64, featureData *sqlEncoderFeatureData, value any) (encoded any, err error) {
	if value != nil {
		switch featureData.schema.featureKind {
		case sfkObject, sfkObjectList:
			eObject := value.(EObject)
			return e.getSQLID(eObject)
		default:
			return e.sqlEncoder.encodeFeatureValue(sqlObjectID, featureData, value)
		}
	}
	return nil, nil
}

// get object sql id
func (e *SQLStore) getSQLID(eObject EObject) (int64, error) {
	// retrieve sql id for eObject
	sqlObjectID, isSqlObjectID := e.sqlIDManager.GetObjectID(eObject)
	if !isSqlObjectID {
		// object is not in store - check if it exists in db
		objectExists := false
		if sqlObjectID != 0 && e.isObjectExists {
			if err := e.executeQuery(
				`SELECT objectID FROM ".objects" WHERE objectID=?`,
				&sqlitex.ExecOptions{
					Args: []any{sqlObjectID},
					ResultFunc: func(stmt *sqlite.Stmt) error {
						objectExists = true
						return nil
					},
				}); err != nil {
				return 0, err
			}
		}
		if objectExists {
			// object exists - register it as decoded
			e.sqlIDManager.SetObjectID(eObject, sqlObjectID)
		} else {
			// object doesn't exists in db - encode it with encoder
			return e.encodeObject(eObject, -1, -1)
		}
	}
	return sqlObjectID, nil
}

func (s *SQLStore) waitOperations(context context.Context, object any) error {
	// compute operations to wait for
	s.mutexOperations.Lock()
	allOperations := []*operation{}
	if object == nil {
		allOperations = make([]*operation, 0)
		for _, objectOperations := range s.objectOperations {
			for _, operations := range objectOperations {
				allOperations = append(allOperations, operations...)
			}
		}
	} else {
		for _, operations := range s.objectOperations[object.(EObject)] {
			allOperations = append(allOperations, operations...)
		}
	}
	s.mutexOperations.Unlock()

	// wait for operations to be finished
	if len(allOperations) > 0 {
		logger := s.logger.Named("ops")
		// debug
		if e := logger.Check(zap.DebugLevel, "waiting operations"); e != nil {
			e.Write(zap.Int64s("operations", mapSlice(allOperations, func(index int, op *operation) int64 { return op.id })))
		}
		// compute promises
		allPromises := mapSlice(allOperations, func(index int, op *operation) *promise.Promise[any] { return op.promise })
		// wait for promises
		_, err := promise.AllWithPool(context, s.pool, allPromises...).Await(context)
		if err != nil {
			return err
		}
		logger.Debug("waiting operations finished")
	}
	return nil
}

func (s *SQLStore) scheduleOperation(context context.Context, op *operation) *promise.Promise[any] {
	logger := s.logger.Named("ops")
	if e := logger.Check(zap.DebugLevel, "schedule"); e != nil {
		var objectStringer fmt.Stringer
		if op.object != nil {
			objectStringer = op.object.(fmt.Stringer)
		}
		var featureString string
		if op.feature != nil {
			featureString = op.feature.GetName()
		}
		e.Write(zap.Int64("goid", goid.Get()),
			zap.String("op", op.cmd),
			zap.Stringer("object", objectStringer),
			zap.String("feature", featureString),
			zap.Int("index", op.index),
			zap.Any("value", op.value))
	}

	type objectFeature struct {
		object  EObject
		feature EStructuralFeature
	}

	// objects features = { object feature operations [, value object lock operations] }
	objectFeatures := []objectFeature{{op.object, op.feature}}
	if object, isObject := op.value.(EObject); isObject {
		objectFeatures = append(objectFeatures, objectFeature{object, nil})
	}

	// previous features = { general lock operations, object lock operations , object feature operations }
	previousFeatures := map[objectFeature]struct{}{
		{nil, nil}:              {},
		{op.object, nil}:        {},
		{op.object, op.feature}: {},
	}

	// lock operations
	s.mutexOperations.Lock()

	// compute previous operations to wait for
	previous := []*operation{}
	for objectFeature := range previousFeatures {
		if objectOperations := s.objectOperations[objectFeature.object]; objectOperations != nil {
			if operations := objectOperations[objectFeature.feature]; len(operations) > 0 {
				switch op.type_ {
				case TaskRead:
					for i := len(operations) - 1; i >= 0; i-- {
						operation := operations[i]
						if operation.type_ == TaskWrite {
							previous = append(previous, operation)
							break
						}
					}
				case TaskWrite:
					for i := len(operations) - 1; i >= 0; i-- {
						operation := operations[i]
						previous = append(previous, operation)
						if operation.type_ == TaskWrite {
							break
						}
					}
				}
			}
		}
	}

	// create promise
	op.promise = promise.NewWithPool(func(resolve func(any), reject func(error)) {
		logger := logger.With(zap.Int64("goid", goid.Get()), zap.Int64("id", op.id))
		// wait for previous operations
		if len(previous) > 0 {
			if e := logger.Check(zap.DebugLevel, "waiting previous operations"); e != nil {
				e.Write(zap.Int64s("previous", mapSlice(previous, func(index int, operation *operation) int64 { return op.id })))
			}
			promises := mapSlice(previous, func(index int, operation *operation) *promise.Promise[any] { return operation.promise })
			if _, err := promise.All(context, promises...).Await(context); err != nil {
				logger.Debug("error in previous operation", zap.Error(err))
				reject(err)
				return
			}
		}

		// execute operation
		logger.Debug("execute")
		result, err := op.fn()

		// clean operations
		logger.Debug("cleaning")
		s.mutexOperations.Lock()
		defer s.mutexOperations.Unlock()
		for _, objectFeature := range objectFeatures {
			if err := s.unregisterOperation(objectFeature.object, objectFeature.feature, op); err != nil {
				reject(err)
				return
			}
		}
		logger.Debug("cleaned")

		if len(s.objectOperations) == 0 {
			logger.Debug("no pending")
		}

		// result or error
		if err != nil {
			s.errorHandler(err)
			reject(err)
		} else {
			resolve(result)
		}
	}, s.promisePool)

	// register operation for object and value if its an object
	// to lock objects for future updates
	for _, of := range objectFeatures {
		s.registerOperation(of.object, of.feature, op)
	}
	s.mutexOperations.Unlock()
	return op.promise
}

func (s *SQLStore) scheduleExclusive(context context.Context, op *operation) *promise.Promise[any] {
	logger := s.logger.Named("ops")
	logger.Debug("schedule", zap.Int64("goid", goid.Get()), zap.String("op", op.cmd))
	return promise.NewWithPool(func(resolve func(any), reject func(error)) {
		s.mutexOperations.Lock()
		objects := slices.Collect(maps.Keys(s.objectOperations))
		size := int64(len(objects))
		s.mutexOperations.Unlock()

		// create a lock semaphore
		run := make(chan struct{})
		locked := semaphore.NewWeighted(size)
		if err := locked.Acquire(context, size); err != nil {
			reject(err)
			return
		}

		// schedule lock operation
		for _, object := range objects {
			s.scheduleOperation(context, newOperation("child-exclusive", op.type_, object, nil, -1, nil, func() (any, error) {
				// the object is locked
				locked.Release(1)

				// wait for the op to be run
				<-run

				return nil, nil
			}))
		}

		// register exclusive operation
		s.mutexOperations.Lock()
		s.registerOperation(nil, nil, op)
		s.mutexOperations.Unlock()

		// wait for all objects to be locked
		logger.Debug("waiting for exclusive access")
		if err := locked.Acquire(context, size); err != nil {
			reject(err)
			return
		}

		// indicate all child exclusive that operation is ready to run
		defer close(run)

		// execute exclusive operation
		logger.Debug("executing with exclusive access")
		result, err := op.fn()

		// unregister exclusive operation
		s.mutexOperations.Lock()
		s.unregisterOperation(nil, nil, op)
		s.mutexOperations.Unlock()

		if err != nil {
			reject(err)
		} else {
			resolve(result)
		}
	}, s.promisePool)
}

func (s *SQLStore) getLastOperation(type_ operationType, object EObject, feature EStructuralFeature, index int) *operation {
	s.mutexOperations.Lock()
	defer s.mutexOperations.Unlock()
	if objectOperations := s.objectOperations[object]; objectOperations != nil {
		if operations := objectOperations[feature]; operations != nil {
			for i := len(operations) - 1; i >= 0; i-- {
				if operation := operations[i]; operation.type_ == type_ && operation.index == index {
					return operation
				}
			}
		}
	}
	return nil
}

func (s *SQLStore) registerOperation(object EObject, feature EStructuralFeature, op *operation) {
	objectOperations := s.objectOperations[object]
	if objectOperations == nil {
		objectOperations = map[EStructuralFeature][]*operation{}
		s.objectOperations[object] = objectOperations
	}
	objectOperations[feature] = append(objectOperations[feature], op)
}

func (s *SQLStore) unregisterOperation(object EObject, feature EStructuralFeature, op *operation) error {
	objectOperations := s.objectOperations[object]
	operations := objectOperations[feature]
	// retrieve operation index
	index := slices.Index(operations, op)
	if index == -1 {
		return errors.New("unable to find task index")
	}
	// remove operation from collection
	copy(operations[index:], operations[index+1:])
	operations[len(operations)-1] = nil
	operations = operations[:len(operations)-1]
	// set object operations for feature
	if len(operations) == 0 {
		delete(objectOperations, feature)
		// cleanup object operations
		if len(objectOperations) == 0 {
			delete(s.objectOperations, object)
		}
	} else {
		objectOperations[feature] = operations
	}
	return nil
}

func awaitPromise[T any](p *promise.Promise[any]) T {
	var def T
	if r, err := p.Await(context.Background()); err != nil {
		return def
	} else if result, isResult := (*r).(T); isResult {
		return result
	} else {
		return def
	}
}

func (s *SQLStore) Get(object EObject, feature EStructuralFeature, index int) any {
	// optimization : if we have a non many feature, check if we have a pending write operation
	if !feature.IsMany() {
		if operation := s.getLastOperation(operationWrite, object, feature, index); operation != nil {
			return operation.value
		}
	}
	return awaitPromise[any](s.scheduleOperation(context.Background(), newOperation("Get", operationRead, object, feature, index, nil, func() (any, error) {
		return s.doGet(object, feature, index)
	})))
}

func (s *SQLStore) doGet(object EObject, feature EStructuralFeature, index int) (any, error) {
	// retrieve value from sqlite db
	sqlID, err := s.getSQLID(object)
	if err != nil {
		return err, nil
	}
	featureSchema, err := s.getFeatureSchema(object, feature)
	if err != nil {
		return err, nil
	}
	return s.getValue(sqlID, featureSchema, index)
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

func (s *SQLStore) Set(object EObject, feature EStructuralFeature, index int, value any, isOldValue bool) (oldValue any) {
	p := s.scheduleOperation(context.Background(), newOperation("Set", operationWrite, object, feature, index, value, func() (any, error) {
		return s.doSet(object, feature, index, value, isOldValue)
	}))
	if isOldValue {
		oldValue = awaitPromise[any](p)
	}
	return
}

func (s *SQLStore) doSet(object EObject, feature EStructuralFeature, index int, value any, isOldValue bool) (oldValue any, err error) {
	// get object sql id
	sqlObjectID, err := s.getSQLID(object)
	if err != nil {
		return nil, err
	}

	// get encoder feature data
	featureData, err := s.getEncoderFeatureData(object, feature)
	if err != nil {
		return nil, err
	}

	if isOldValue {
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
		if operation := s.getLastOperation(operationWrite, object, feature, -1); operation != nil {
			return operation.value != nil
		}
	}
	return awaitPromise[bool](s.scheduleOperation(context.Background(), newOperation("IsSet", operationRead, object, feature, -1, nil, func() (any, error) {
		return s.doIsSet(object, feature)
	})))
}

func (s *SQLStore) doIsSet(object EObject, feature EStructuralFeature) (bool, error) {
	sqlObjectID, err := s.getSQLID(object)
	if err != nil {
		return false, err
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
		var value any
		if err := s.executeQuery(
			s.getManyQueries(featureTable).getExistsQuery(),
			&sqlitex.ExecOptions{
				Args: []any{sqlObjectID},
				ResultFunc: func(stmt *sqlite.Stmt) error {
					value = decodeAny(stmt, 0)
					return nil
				}}); err != nil {
			return false, err
		}
		return value != nil, nil
	}
	return false, nil
}

func (s *SQLStore) UnSet(object EObject, feature EStructuralFeature) {
	awaitPromise[any](s.scheduleOperation(context.Background(), newOperation("UnSet", operationWrite, object, feature, -1, nil, func() (any, error) {
		return s.doUnSet(object, feature)
	})))
}

func (s *SQLStore) doUnSet(object EObject, feature EStructuralFeature) (any, error) {
	sqlObjectID, err := s.getSQLID(object)
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
	return awaitPromise[bool](s.scheduleOperation(context.Background(), newOperation("IsEmpty", operationRead, object, feature, -1, nil, func() (any, error) {
		return s.doIsEmpty(object, feature)
	})))
}

func (s *SQLStore) doIsEmpty(object EObject, feature EStructuralFeature) (bool, error) {
	sqlObjectID, err := s.getSQLID(object)
	if err != nil {
		return false, err
	}

	// retrieve table
	featureTable, err := s.getFeatureTable(object, feature)
	if err != nil {
		return false, err
	}

	// retrieve statement
	var value any
	if err := s.executeQuery(
		s.getManyQueries(featureTable).getExistsQuery(),
		&sqlitex.ExecOptions{
			Args: []any{sqlObjectID},
			ResultFunc: func(stmt *sqlite.Stmt) error {
				value = decodeAny(stmt, 0)
				return nil
			}}); err != nil {
		return true, nil
	}
	return value == nil, nil

}

func (s *SQLStore) Size(object EObject, feature EStructuralFeature) int {
	return awaitPromise[int](s.scheduleOperation(context.Background(), newOperation("Size", operationRead, object, feature, -1, nil, func() (any, error) {
		return s.doSize(object, feature)
	})))
}

func (s *SQLStore) doSize(object EObject, feature EStructuralFeature) (int, error) {
	sqlObjectID, err := s.getSQLID(object)
	if err != nil {
		return 0, err
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
	return awaitPromise[bool](s.scheduleOperation(context.Background(), newOperation("Contains", operationRead, object, feature, -1, value, func() (any, error) {
		return s.doContains(object, feature, value)
	})))
}

func (s *SQLStore) doContains(object EObject, feature EStructuralFeature, value any) (bool, error) {
	sqlObjectID, err := s.getSQLID(object)
	if err != nil {
		return false, err
	}

	// retrieve table
	featureData, err := s.getEncoderFeatureData(object, feature)
	if err != nil {
		return false, err
	}

	// query statement
	encoded, err := s.encodeFeatureValue(sqlObjectID, featureData, value)
	if err != nil {
		return false, err
	}

	var rowid int64
	featureTable := featureData.schema.table
	if err := s.executeQuery(
		s.getManyQueries(featureTable).getContainsQuery(),
		&sqlitex.ExecOptions{
			Args: []any{sqlObjectID, encoded},
			ResultFunc: func(stmt *sqlite.Stmt) error {
				rowid = stmt.ColumnInt64(0)
				return nil
			}}); err != nil {
		return false, err
	}
	return rowid != 0, nil
}

func (s *SQLStore) doIndexOf(object EObject, feature EStructuralFeature, value any, getIndexOfQuery func(*sqlManyQueries) string) (int, error) {
	sqlObjectID, err := s.getSQLID(object)
	if err != nil {
		return -1, err
	}

	// retrieve table
	featureData, err := s.getEncoderFeatureData(object, feature)
	if err != nil {
		return -1, err
	}
	featureTable := featureData.schema.table

	// compute parameters
	encoded, err := s.encodeFeatureValue(sqlObjectID, featureData, value)
	if err != nil {
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
	return awaitPromise[int](s.scheduleOperation(context.Background(), newOperation("IndexOf", operationRead, object, feature, -1, value, func() (any, error) {
		return s.doIndexOf(object, feature, value, func(sms *sqlManyQueries) string {
			return sms.getIndexOfQuery()
		})
	})))
}

func (s *SQLStore) LastIndexOf(object EObject, feature EStructuralFeature, value any) int {
	return awaitPromise[int](s.scheduleOperation(context.Background(), newOperation("LastIndexOf", operationRead, object, feature, -1, value, func() (any, error) {
		return s.doIndexOf(object, feature, value, func(sms *sqlManyQueries) string {
			return sms.getLastIndexOfQuery()
		})
	})))
}

// AddRoot add object as store root
func (s *SQLStore) AddRoot(object EObject) {
	if err := s.encodeContent(object); err != nil {
		s.errorHandler(err)
	}
}

// RemoveRoot implements EStore.
func (s *SQLStore) RemoveRoot(object EObject) {
	sqlObjectID, err := s.getSQLID(object)
	if err != nil {
		s.errorHandler(err)
		return
	}

	if err := s.executeQuery(
		s.getSingleQueries(s.schema.contentsTable.key).getRemoveQuery(),
		&sqlitex.ExecOptions{Args: []any{sqlObjectID}},
	); err != nil {
		s.errorHandler(err)
	}
}

// GetRoot return root objects
func (s *SQLStore) GetRoots() []EObject {
	contents, err := s.decodeContents()
	if err != nil {
		s.errorHandler(err)
	}
	return contents
}

func (s *SQLStore) UnRegister(object EObject) {
	s.sqlIDManager.ClearObjectID(object)
}

func (s *SQLStore) Add(object EObject, feature EStructuralFeature, index int, value any) {
	awaitPromise[any](s.scheduleOperation(context.Background(), newOperation("Add", operationWrite, object, feature, index, value, func() (any, error) {
		return s.doAdd(object, feature, index, value)
	})))
}

func (s *SQLStore) doAdd(object EObject, feature EStructuralFeature, index int, value any) (any, error) {
	sqlObjectID, err := s.getSQLID(object)
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
	awaitPromise[any](s.scheduleOperation(context.Background(), newOperation("AddAll", operationWrite, object, feature, index, c, func() (any, error) {
		return s.doAddAll(object, feature, index, c)
	})))
}

func (s *SQLStore) doAddAll(object EObject, feature EStructuralFeature, index int, c Collection) (any, error) {
	sqlObjectID, err := s.getSQLID(object)
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
			s.errorHandler(err)
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
			s.errorHandler(err)
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

func (s *SQLStore) Remove(object EObject, feature EStructuralFeature, index int) any {
	return awaitPromise[any](s.scheduleOperation(context.Background(), newOperation("Remove", operationWrite, object, feature, index, nil, func() (any, error) {
		return s.doRemove(object, feature, index)
	})))
}

func (s *SQLStore) doRemove(object EObject, feature EStructuralFeature, index int) (any, error) {
	sqlObjectID, err := s.getSQLID(object)
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
	decoded, err := s.decodeFeatureValue(featureData.schema, value)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

type moveIndexes struct {
	sourceIndex int
	targetIndex int
}

func (mi moveIndexes) String() string {
	return fmt.Sprintf("{ sourceIndex : %v, targetIndex %v}", mi.sourceIndex, mi.targetIndex)
}

func (s *SQLStore) Move(object EObject, feature EStructuralFeature, sourceIndex int, targetIndex int) any {
	return awaitPromise[any](s.scheduleOperation(context.Background(), newOperation("Move", operationWrite, object, feature, -1, moveIndexes{sourceIndex, targetIndex}, func() (any, error) {
		return s.doMove(object, feature, sourceIndex, targetIndex)
	})))
}

func (s *SQLStore) doMove(object EObject, feature EStructuralFeature, sourceIndex int, targetIndex int) (any, error) {
	sqlObjectID, err := s.getSQLID(object)
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
	decoded, err := s.decodeFeatureValue(featureSchema, value)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}

func (s *SQLStore) Clear(object EObject, feature EStructuralFeature) {
	awaitPromise[any](s.scheduleOperation(context.Background(), newOperation("Clear", operationWrite, object, feature, -1, nil, func() (any, error) {
		return s.doClear(object, feature)
	})))
}

func (s *SQLStore) doClear(object EObject, feature EStructuralFeature) (any, error) {
	sqlObjectID, err := s.getSQLID(object)
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
	if result := awaitPromise[*containerAndFeature](s.scheduleOperation(context.Background(), newOperation("GetContainer", operationRead, object, nil, -1, nil, func() (any, error) {
		return s.doGetContainer(object)
	}))); result != nil {
		return result.container, result.feature
	}
	return nil, nil
}

func (s *SQLStore) doGetContainer(object EObject) (*containerAndFeature, error) {
	sqlID, err := s.getSQLID(object)
	if err != nil {
		return nil, err
	}

	containerID := int64(-1)
	containerFeatureID := int64(-1)
	if err := s.executeQuery(`SELECT containerID,containerFeatureID FROM ".objects" WHERE objectID=?`, &sqlitex.ExecOptions{
		Args: []any{sqlID},
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
		return nil, fmt.Errorf("unable to find container for object '%v'", sqlID)
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

func (s *SQLStore) All(object EObject, feature EStructuralFeature) iter.Seq[any] {
	return func(yield func(any) bool) {
		awaitPromise[any](s.scheduleOperation(context.Background(), newOperation("All", operationRead, object, feature, -1, nil, func() (any, error) {
			interrupted := errors.New("interrupted")
			sqlID, err := s.getSQLID(object)
			if err != nil {
				return nil, err
			}
			featureSchema, err := s.getFeatureSchema(object, feature)
			if err != nil {
				return nil, err
			}
			featureTable := featureSchema.table
			if err := s.executeQuery(
				s.getManyQueries(featureTable).getSelectAllQuery(),
				&sqlitex.ExecOptions{
					Args: []any{sqlID},
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
		})))
	}
}

func (s *SQLStore) ToArray(object EObject, feature EStructuralFeature) []any {
	return slices.Collect(s.All(object, feature))
}

func (s *SQLStore) Serialize(ctx context.Context) ([]byte, error) {
	p := s.scheduleExclusive(ctx, newOperation("Serialize", operationRead, nil, nil, -1, nil, func() (any, error) {
		return s.doSerialize(ctx)
	}))
	r, err := p.Await(ctx)
	if err != nil {
		return nil, err
	}
	return (*r).([]byte), nil
}

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

	// open connection to serialize
	conn, err := s.connPool.Take(ctx)
	if err != nil {
		return nil, err
	}
	defer s.connPool.Put(conn)

	// supports big databases
	if dbSize > SQLITE_MAX_ALLOCATION_SIZE {
		dbPath, err := sqlTmpDB("store-serialize")
		if err != nil {
			return nil, err
		}

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

func (s *SQLStore) Close() error {
	s.logger.Named("ops").Debug("Close",
		zap.Int64("goid", goid.Get()),
	)

	if err := s.waitOperations(context.Background(), nil); err != nil {
		return err
	}

	if err := s.taskManager.Close(); err != nil {
		return err
	}

	if err := s.connPoolClose(s.connPool); err != nil {
		return err
	}

	s.antsPool.Release()

	return nil
}
