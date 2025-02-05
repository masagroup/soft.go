package ecore

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/chebyrash/promise"
	"go.uber.org/zap"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

// log queries
var logQuery bool = false

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

type operation struct {
	type_    OperationType
	promise_ *promise.Promise[any]
}

type SQLStore struct {
	*sqlBase
	sqlDecoder
	sqlEncoder
	pool                *sqlitex.Pool
	errorHandler        func(error)
	sqlIDManager        SQLStoreIDManager
	singleQueries       map[*sqlColumn]*sqlSingleQueries
	manyQueries         map[*sqlTable]*sqlManyQueries
	mutexQueries        sync.Mutex
	connectionPoolClose func(conn *sqlitex.Pool) error
	executeQuery        func(conn *sqlite.Conn, query string, opts *sqlitex.ExecOptions) error
	operations          map[any][]*operation
	mutexOperations     sync.Mutex
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

	// log sqlite queries
	executeQuery := sqlitex.Execute
	if logQuery {
		executeQuery = func(conn *sqlite.Conn, query string, opts *sqlitex.ExecOptions) error {
			start := time.Now()
			if err := sqlitex.Execute(conn, query, opts); err != nil {
				return err
			}
			logger.Debug("execute",
				zap.String("query", query),
				zap.Any("args", opts.Args),
				zap.Duration("duration", time.Since(start)))
			return nil
		}
	}

	// retrieve connection pool
	pool, err := connectionPoolProvider()
	if err != nil {
		return nil, err
	}

	// close pool if there is an error
	defer func() {
		if err != nil {
			_ = connectionPoolClose(pool)
		}
	}()

	// create sql base
	base := &sqlBase{
		codecVersion:    codecVersion,
		uri:             resourceURI,
		objectIDName:    objectIDName,
		objectIDManager: idManager,
		isContainerID:   true,
		isObjectID:      len(objectIDName) > 0 && objectIDName != "objectID" && idManager != nil,
		logger:          logger,
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
			classDataMap:     map[EClass]*sqlEncoderClassData{},
			sqlIDManager:     sqlIDManager,
			sqlObjectManager: sqlObjectManager,
		},
		pool:                pool,
		sqlIDManager:        sqlIDManager,
		errorHandler:        errorHandler,
		singleQueries:       map[*sqlColumn]*sqlSingleQueries{},
		manyQueries:         map[*sqlTable]*sqlManyQueries{},
		connectionPoolClose: connectionPoolClose,
		executeQuery:        executeQuery,
		operations:          map[any][]*operation{},
	}

	// set store in sql object manager
	sqlObjectManager.store = store

	// decode version
	if err = store.decodeVersion(pool); err != nil {
		return nil, err
	}

	// decode schema
	if err = store.decodeSchema(pool, []sqlSchemaOption{withCreateIfNotExists(true)}); err != nil {
		return nil, err
	}

	// write connection
	conn, err := pool.Take(context.Background())
	if err != nil {
		return nil, err
	}
	defer pool.Put(conn)

	// encode version
	if err = store.encodeVersion(conn); err != nil {
		return nil, err
	}

	// encode schema
	if err = store.encodeSchema(conn); err != nil {
		return nil, err
	}

	// encode schema
	if err = store.encodeProperties(conn); err != nil {
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

func (s *SQLStore) getEncoderFeatureData(conn *sqlite.Conn, object EObject, feature EStructuralFeature) (*sqlEncoderFeatureData, error) {
	// retrieve class schema
	class := object.EClass()

	// retrieve class data
	classData, err := s.getEncoderClassData(conn, class)
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

func (e *SQLStore) encodeFeatureValue(conn *sqlite.Conn, featureData *sqlEncoderFeatureData, value any) (encoded any, err error) {
	if value != nil {
		switch featureData.schema.featureKind {
		case sfkObject, sfkObjectList:
			eObject := value.(EObject)
			return e.getSQLID(conn, eObject)
		default:
			return e.sqlEncoder.encodeFeatureValue(conn, featureData, value)
		}
	}
	return nil, nil
}

// get object sql id
func (s *SQLStore) getSQLID(conn *sqlite.Conn, eObject EObject) (int64, error) {
	// retrieve sql id for eObject
	sqlID, isSQLID := s.sqlIDManager.GetObjectID(eObject)
	if !isSQLID {
		// object is not in store - check if it exists in db
		objectExists := false
		objectsTable := s.schema.objectsTable
		if err := s.executeQuery(
			conn,
			s.getSingleQueries(objectsTable.key).getSelectQuery(),
			&sqlitex.ExecOptions{
				Args: []any{sqlID},
				ResultFunc: func(stmt *sqlite.Stmt) error {
					objectExists = true
					return nil
				},
			}); err != nil {
			return 0, err
		}
		if objectExists {
			// object exists - register it as decoded
			s.sqlIDManager.SetObjectID(eObject, sqlID)
		} else {
			// object doesn't exists in db - encode it with encoder
			return s.sqlEncoder.encodeObject(conn, eObject)
		}
	}
	return sqlID, nil
}

func (s *SQLStore) Get(object EObject, feature EStructuralFeature, index int) any {
	conn, err := s.pool.Take(context.Background())
	if err != nil {
		s.errorHandler(err)
		return nil
	}
	defer s.pool.Put(conn)

	sqlID, err := s.getSQLID(conn, object)
	if err != nil {
		s.errorHandler(err)
		return nil
	}
	featureSchema, err := s.getFeatureSchema(object, feature)
	if err != nil {
		s.errorHandler(err)
		return nil
	}
	return s.getValue(conn, sqlID, featureSchema, index)
}

func (s *SQLStore) getValue(conn *sqlite.Conn, sqlID int64, featureSchema *sqlFeatureSchema, index int) any {
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
		conn,
		query,
		&sqlitex.ExecOptions{
			Args: args,
			ResultFunc: func(stmt *sqlite.Stmt) error {
				value = decodeAny(stmt, 0)
				return nil
			}}); err != nil {
		s.errorHandler(err)
		return nil
	}

	decoded, err := s.decodeFeatureValue(conn, featureSchema, value)
	if err != nil {
		s.errorHandler(err)
	}
	return decoded
}

func (s *SQLStore) Set(object EObject, feature EStructuralFeature, index int, value any, isOldValue bool) (oldValue any) {
	// connection
	conn, err := s.pool.Take(context.Background())
	if err != nil {
		s.errorHandler(err)
		return
	}
	defer s.pool.Put(conn)

	// get object sql id
	sqlID, err := s.getSQLID(conn, object)
	if err != nil {
		s.errorHandler(err)
		return
	}

	// get encoder feature data
	featureData, err := s.getEncoderFeatureData(conn, object, feature)
	if err != nil {
		s.errorHandler(err)
		return
	}

	if isOldValue {
		// retrieve previous value
		oldValue = s.getValue(conn, sqlID, featureData.schema, index)
	}

	// encode value
	encoded, err := s.encodeFeatureValue(conn, featureData, value)
	if err != nil {
		s.errorHandler(err)
		return
	}

	var query string
	var args []any
	if featureColumn := featureData.schema.column; featureColumn != nil {
		query = s.getSingleQueries(featureColumn).getUpdateQuery()
		args = []any{encoded, sqlID}
	} else if featureTable := featureData.schema.table; featureTable != nil {
		query = s.getManyQueries(featureTable).getUpdateValueQuery()
		args = []any{encoded, sqlID, index}
	}

	if err := s.executeQuery(conn, query, &sqlitex.ExecOptions{Args: args}); err != nil {
		s.errorHandler(err)
		return
	}

	return
}

func (s *SQLStore) IsSet(object EObject, feature EStructuralFeature) bool {
	conn, err := s.pool.Take(context.Background())
	if err != nil {
		s.errorHandler(err)
		return false
	}
	defer s.pool.Put(conn)

	sqlID, err := s.getSQLID(conn, object)
	if err != nil {
		s.errorHandler(err)
		return false
	}

	featureSchema, err := s.getFeatureSchema(object, feature)
	if err != nil {
		s.errorHandler(err)
		return false
	}
	if featureColumn := featureSchema.column; featureColumn != nil {
		var value any
		if err := s.executeQuery(
			conn,
			s.getSingleQueries(featureColumn).getSelectQuery(),
			&sqlitex.ExecOptions{
				Args: []any{sqlID},
				ResultFunc: func(stmt *sqlite.Stmt) error {
					value = decodeAny(stmt, 0)
					return nil
				}}); err != nil {
			s.errorHandler(err)
		}
		return value != featureSchema.feature.GetDefaultValue()
	} else if featureTable := featureSchema.table; featureTable != nil {
		var value any
		if err := s.executeQuery(
			conn,
			s.getManyQueries(featureTable).getExistsQuery(),
			&sqlitex.ExecOptions{
				Args: []any{sqlID},
				ResultFunc: func(stmt *sqlite.Stmt) error {
					value = decodeAny(stmt, 0)
					return nil
				}}); err != nil {
			s.errorHandler(err)
		}
		return value != nil
	}
	return false
}

func (s *SQLStore) UnSet(object EObject, feature EStructuralFeature) {
	conn, err := s.pool.Take(context.Background())
	if err != nil {
		s.errorHandler(err)
		return
	}
	defer s.pool.Put(conn)

	sqlID, err := s.getSQLID(conn, object)
	if err != nil {
		s.errorHandler(err)
		return
	}

	featureData, err := s.getEncoderFeatureData(conn, object, feature)
	if err != nil {
		s.errorHandler(err)
		return
	}

	var query string
	var args []any
	if featureColumn := featureData.schema.column; featureColumn != nil {
		query = s.getSingleQueries(featureColumn).getUpdateQuery()
		args = []any{feature.GetDefaultValue(), sqlID}
	} else if featureTable := featureData.schema.table; featureTable != nil {
		query = s.getManyQueries(featureTable).getClearQuery()
		args = []any{sqlID}
	}
	if err := s.executeQuery(conn, query, &sqlitex.ExecOptions{Args: args}); err != nil {
		s.errorHandler(err)
	}
}

func (s *SQLStore) IsEmpty(object EObject, feature EStructuralFeature) bool {
	conn, err := s.pool.Take(context.Background())
	if err != nil {
		s.errorHandler(err)
		return false
	}
	defer s.pool.Put(conn)

	sqlID, err := s.getSQLID(conn, object)
	if err != nil {
		s.errorHandler(err)
		return false
	}

	// retrieve table
	featureTable, err := s.getFeatureTable(object, feature)
	if err != nil {
		s.errorHandler(err)
		return false
	}

	// retrieve statement
	var value any
	if err := s.executeQuery(
		conn,
		s.getManyQueries(featureTable).getExistsQuery(),
		&sqlitex.ExecOptions{
			Args: []any{sqlID},
			ResultFunc: func(stmt *sqlite.Stmt) error {
				value = decodeAny(stmt, 0)
				return nil
			}}); err != nil {
		s.errorHandler(err)
	}
	return value == nil
}

func (s *SQLStore) Size(object EObject, feature EStructuralFeature) int {
	conn, err := s.pool.Take(context.Background())
	if err != nil {
		s.errorHandler(err)
		return 0
	}
	defer s.pool.Put(conn)

	sqlID, err := s.getSQLID(conn, object)
	if err != nil {
		s.errorHandler(err)
		return 0
	}

	// retrieve table
	featureTable, err := s.getFeatureTable(object, feature)
	if err != nil {
		s.errorHandler(err)
		return 0
	}

	var size int
	if err := s.executeQuery(
		conn,
		s.getManyQueries(featureTable).getCountQuery(),
		&sqlitex.ExecOptions{
			Args: []any{sqlID},
			ResultFunc: func(stmt *sqlite.Stmt) error {
				size = stmt.ColumnInt(0)
				return nil
			}}); err != nil {
		s.errorHandler(err)
	}
	return size
}

func (s *SQLStore) Contains(object EObject, feature EStructuralFeature, value any) bool {
	conn, err := s.pool.Take(context.Background())
	if err != nil {
		s.errorHandler(err)
		return false
	}
	defer s.pool.Put(conn)

	sqlID, err := s.getSQLID(conn, object)
	if err != nil {
		s.errorHandler(err)
		return false
	}

	// retrieve table
	featureData, err := s.getEncoderFeatureData(conn, object, feature)
	if err != nil {
		s.errorHandler(err)
		return false
	}

	// query statement
	encoded, err := s.encodeFeatureValue(conn, featureData, value)
	if err != nil {
		s.errorHandler(err)
		return false
	}

	var rowid int64
	featureTable := featureData.schema.table
	if err := s.executeQuery(
		conn,
		s.getManyQueries(featureTable).getContainsQuery(),
		&sqlitex.ExecOptions{
			Args: []any{sqlID, encoded},
			ResultFunc: func(stmt *sqlite.Stmt) error {
				rowid = stmt.ColumnInt64(0)
				return nil
			}}); err != nil {
		s.errorHandler(err)
	}
	return rowid != 0
}

func (s *SQLStore) indexOf(object EObject, feature EStructuralFeature, value any, getIndexOfQuery func(*sqlManyQueries) string) int {
	conn, err := s.pool.Take(context.Background())
	if err != nil {
		s.errorHandler(err)
		return -1
	}
	defer s.pool.Put(conn)

	sqlID, err := s.getSQLID(conn, object)
	if err != nil {
		s.errorHandler(err)
		return -1
	}

	// retrieve table
	featureData, err := s.getEncoderFeatureData(conn, object, feature)
	if err != nil {
		s.errorHandler(err)
		return -1
	}
	featureTable := featureData.schema.table

	// compute parameters
	encoded, err := s.encodeFeatureValue(conn, featureData, value)
	if err != nil {
		s.errorHandler(err)
		return -1
	}
	// retrieve row idx in table
	idx := -1.0
	if err := s.executeQuery(
		conn,
		getIndexOfQuery(s.getManyQueries(featureTable)),
		&sqlitex.ExecOptions{
			Args: []any{sqlID, encoded},
			ResultFunc: func(stmt *sqlite.Stmt) error {
				idx = stmt.ColumnFloat(0)
				return nil
			}}); err != nil {
		s.errorHandler(err)
		return -1
	}
	if idx == -1.0 {
		return -1
	}

	// convert idx to list index - index is the count of rows where idx < expected idx
	index := -1
	if err := s.executeQuery(
		conn,
		s.getManyQueries(featureTable).getIdxToListIndexQuery(),
		&sqlitex.ExecOptions{
			Args: []any{sqlID, idx},
			ResultFunc: func(stmt *sqlite.Stmt) error {
				index = stmt.ColumnInt(0)
				return nil
			}}); err != nil {
		s.errorHandler(err)
		return -1
	}
	return index
}

func (s *SQLStore) IndexOf(object EObject, feature EStructuralFeature, value any) int {
	return s.indexOf(object, feature, value, func(sms *sqlManyQueries) string {
		return sms.getIndexOfQuery()
	})
}

func (s *SQLStore) LastIndexOf(object EObject, feature EStructuralFeature, value any) int {
	return s.indexOf(object, feature, value, func(sms *sqlManyQueries) string {
		return sms.getLastIndexOfQuery()
	})
}

// AddRoot add object as store root
func (s *SQLStore) AddRoot(object EObject) {
	conn, err := s.pool.Take(context.Background())
	if err != nil {
		s.errorHandler(err)
		return
	}
	defer s.pool.Put(conn)

	if err := s.encodeContent(conn, object); err != nil {
		s.errorHandler(err)
	}
}

// RemoveRoot implements EStore.
func (s *SQLStore) RemoveRoot(object EObject) {
	conn, err := s.pool.Take(context.Background())
	if err != nil {
		s.errorHandler(err)
		return
	}
	defer s.pool.Put(conn)

	sqlID, err := s.getSQLID(conn, object)
	if err != nil {
		s.errorHandler(err)
		return
	}

	if err := s.executeQuery(
		conn,
		s.getSingleQueries(s.schema.contentsTable.key).getRemoveQuery(),
		&sqlitex.ExecOptions{Args: []any{sqlID}},
	); err != nil {
		s.errorHandler(err)
	}
}

// GetRoot return root objects
func (s *SQLStore) GetRoots() []EObject {
	conn, err := s.pool.Take(context.Background())
	if err != nil {
		s.errorHandler(err)
		return nil
	}
	defer s.pool.Put(conn)

	contents, err := s.decodeContents(conn)
	if err != nil {
		s.errorHandler(err)
	}
	return contents
}

func (s *SQLStore) UnRegister(object EObject) {
	s.sqlIDManager.ClearObjectID(object)
}

func (s *SQLStore) Add(object EObject, feature EStructuralFeature, index int, value any) {
	conn, err := s.pool.Take(context.Background())
	if err != nil {
		s.errorHandler(err)
		return
	}
	defer s.pool.Put(conn)

	sqlID, err := s.getSQLID(conn, object)
	if err != nil {
		s.errorHandler(err)
		return
	}
	featureData, err := s.getEncoderFeatureData(conn, object, feature)
	if err != nil {
		s.errorHandler(err)
		return
	}
	featureTable := featureData.schema.table
	idx, _, err := s.getInsertIdx(conn, featureTable, sqlID, index, 1)
	if err != nil {
		s.errorHandler(err)
		return
	}
	encoded, err := s.encodeFeatureValue(conn, featureData, value)
	if err != nil {
		s.errorHandler(err)
		return
	}
	if err := s.executeQuery(
		conn,
		s.getManyQueries(featureTable).getInsertQuery(),
		&sqlitex.ExecOptions{Args: []any{sqlID, idx, encoded}},
	); err != nil {
		s.errorHandler(err)
	}
}

func (s *SQLStore) AddAll(object EObject, feature EStructuralFeature, index int, c Collection) {
	conn, err := s.pool.Take(context.Background())
	if err != nil {
		s.errorHandler(err)
		return
	}
	defer s.pool.Put(conn)

	sqlID, err := s.getSQLID(conn, object)
	if err != nil {
		s.errorHandler(err)
		return
	}
	featureData, err := s.getEncoderFeatureData(conn, object, feature)
	if err != nil {
		s.errorHandler(err)
		return
	}
	featureTable := featureData.schema.table
	idx, inc, err := s.getInsertIdx(conn, featureTable, sqlID, index, c.Size())
	if err != nil {
		s.errorHandler(err)
		return
	}
	// encode each value
	query := s.getManyQueries(featureTable).getInsertQuery()
	for it := c.Iterator(); it.HasNext(); {
		value := it.Next()
		v, err := s.encodeFeatureValue(conn, featureData, value)
		if err != nil {
			s.errorHandler(err)
			return
		}
		if err := s.executeQuery(
			conn,
			query,
			&sqlitex.ExecOptions{Args: []any{sqlID, idx, v}},
		); err != nil {
			s.errorHandler(err)
			return
		}
		idx += inc
	}
}

// compute insert idx of element in the list whose index is index. nb is the number of elements to be inserted
// return idx, increment (for each inserted element) and error if any
func (s *SQLStore) getInsertIdx(conn *sqlite.Conn, table *sqlTable, sqlID int64, index int, nb int) (float64, float64, error) {
	if index == 0 {
		// first row in the list
		idx := 1.0
		withElements := false
		if err := s.executeQuery(
			conn,
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
			conn,
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
			panic(fmt.Sprintf("invalid index in table %v for object %v : %v not in list bounds", index, table.name, sqlID))
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
	conn, err := s.pool.Take(context.Background())
	if err != nil {
		s.errorHandler(err)
		return nil
	}
	defer s.pool.Put(conn)

	sqlID, err := s.getSQLID(conn, object)
	if err != nil {
		s.errorHandler(err)
		return nil
	}
	featureData, err := s.getEncoderFeatureData(conn, object, feature)
	if err != nil {
		s.errorHandler(err)
		return nil
	}
	featureTable := featureData.schema.table

	var value any
	if err := s.executeQuery(
		conn,
		s.getManyQueries(featureTable).getRemoveQuery(),
		&sqlitex.ExecOptions{
			Args: []any{sqlID, index},
			ResultFunc: func(stmt *sqlite.Stmt) error {
				value = decodeAny(stmt, 0)
				return nil
			}}); err != nil {
		s.errorHandler(err)
		return nil
	}

	// decode previous value
	decoded, err := s.decodeFeatureValue(conn, featureData.schema, value)
	if err != nil {
		s.errorHandler(err)
		return nil
	}
	return decoded
}

func (s *SQLStore) Move(object EObject, feature EStructuralFeature, sourceIndex int, targetIndex int) any {
	conn, err := s.pool.Take(context.Background())
	if err != nil {
		s.errorHandler(err)
		return nil
	}
	defer s.pool.Put(conn)

	sqlID, err := s.getSQLID(conn, object)
	if err != nil {
		s.errorHandler(err)
		return nil
	}
	featureSchema, err := s.getFeatureSchema(object, feature)
	if err != nil {
		s.errorHandler(err)
		return nil
	}
	featureTable := featureSchema.table

	// compute target index
	if targetIndex > sourceIndex {
		targetIndex++
	}
	idx, _, err := s.getInsertIdx(conn, featureTable, sqlID, targetIndex, 1)
	if err != nil {
		s.errorHandler(err)
		return nil
	}

	// update idx of source index row with target idx
	var value any
	if err := s.executeQuery(
		conn,
		s.getManyQueries(featureTable).getUpdateIdxQuery(),
		&sqlitex.ExecOptions{
			Args: []any{idx, sqlID, sourceIndex},
			ResultFunc: func(stmt *sqlite.Stmt) error {
				value = decodeAny(stmt, 0)
				return nil
			}}); err != nil {
		s.errorHandler(err)
		return nil
	}

	// decode value
	decoded, err := s.decodeFeatureValue(conn, featureSchema, value)
	if err != nil {
		s.errorHandler(err)
		return nil
	}

	return decoded
}

func (s *SQLStore) Clear(object EObject, feature EStructuralFeature) {
	conn, err := s.pool.Take(context.Background())
	if err != nil {
		s.errorHandler(err)
		return
	}
	defer s.pool.Put(conn)

	sqlID, err := s.getSQLID(conn, object)
	if err != nil {
		s.errorHandler(err)
		return
	}
	featureTable, err := s.getFeatureTable(object, feature)
	if err != nil {
		s.errorHandler(err)
		return
	}

	if err := s.executeQuery(
		conn,
		s.getManyQueries(featureTable).getClearQuery(),
		&sqlitex.ExecOptions{Args: []any{sqlID}},
	); err != nil {
		s.errorHandler(err)
	}
}

func (s *SQLStore) GetContainer(object EObject) (container EObject, feature EStructuralFeature) {
	conn, err := s.pool.Take(context.Background())
	if err != nil {
		s.errorHandler(err)
		return
	}
	defer s.pool.Put(conn)

	sqlID, err := s.getSQLID(conn, object)
	if err != nil {
		s.errorHandler(err)
		return
	}

	containerID := int64(-1)
	containerFeatureID := int64(-1)
	if err := s.executeQuery(conn, `SELECT containerID,containerFeatureID FROM ".objects" WHERE objectID=?`, &sqlitex.ExecOptions{
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
		s.errorHandler(err)
		return
	}

	switch containerID {
	case -1:
		s.errorHandler(fmt.Errorf("unable to find container for object '%v'", sqlID))
	case 0:
	default:
		container, err = s.decodeObject(conn, containerID)
		if err != nil {
			s.errorHandler(err)
			return
		}

		containerInternal := container.(EObjectInternal)
		containerClass := containerInternal.EClass()
		feature = containerClass.GetEStructuralFeature(containerInternal.EStaticFeatureCount() + int(containerFeatureID))
	}
	return
}

func (s *SQLStore) SetContainer(object EObject, container EObject, feature EStructuralFeature) {
	conn, err := s.pool.Take(context.Background())
	if err != nil {
		s.errorHandler(err)
		return
	}
	defer s.pool.Put(conn)

	sqlObjectID, err := s.getSQLID(conn, object)
	if err != nil {
		s.errorHandler(err)
		return
	}

	var sqlContainerID any
	if container != nil {
		sqlContainerID, err = s.getSQLID(conn, container)
		if err != nil {
			s.errorHandler(err)
			return
		}
	}

	var featureID any
	if container != nil && feature != nil {
		featureID = container.EClass().GetFeatureID(feature)
	}

	if err := s.executeQuery(conn, `UPDATE ".objects" SET containerID=?,containerFeatureID=? WHERE objectID=?`, &sqlitex.ExecOptions{
		Args: []any{sqlContainerID, featureID, sqlObjectID},
	}); err != nil {
		s.errorHandler(err)
		return
	}

}

func (s *SQLStore) All(object EObject, feature EStructuralFeature) iter.Seq[any] {
	return func(yield func(any) bool) {
		interrupted := errors.New("interrupted")
		conn, err := s.pool.Take(context.Background())
		if err != nil {
			s.errorHandler(err)
			return
		}
		defer s.pool.Put(conn)

		sqlID, err := s.getSQLID(conn, object)
		if err != nil {
			s.errorHandler(err)
			return
		}
		featureSchema, err := s.getFeatureSchema(object, feature)
		if err != nil {
			s.errorHandler(err)
			return
		}
		featureTable := featureSchema.table
		if err := s.executeQuery(
			conn,
			s.getManyQueries(featureTable).getSelectAllQuery(),
			&sqlitex.ExecOptions{
				Args: []any{sqlID},
				ResultFunc: func(stmt *sqlite.Stmt) error {
					value := decodeAny(stmt, 0)
					decoded, err := s.decodeFeatureValue(conn, featureSchema, value)
					if err != nil {
						return err
					}
					if !yield(decoded) {
						return interrupted
					}
					return nil
				}}); err != nil {
			if err != interrupted {
				s.errorHandler(err)
			}
		}
	}
}

func (s *SQLStore) ToArray(object EObject, feature EStructuralFeature) []any {
	return slices.Collect(s.All(object, feature))
}

func (s *SQLStore) Serialize(ctx context.Context) ([]byte, error) {
	// connection
	conn, err := s.pool.Take(ctx)
	if err != nil {
		return nil, err
	}
	defer s.pool.Put(conn)

	// retrive database size
	var dbSize int64
	if err := sqlitex.ExecuteTransient(conn, "SELECT page_count * page_size as size FROM pragma_page_count(), pragma_page_size();", &sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			dbSize = stmt.ColumnInt64(0)
			return nil
		},
	}); err != nil {
		return nil, err
	}

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
	if err := s.WaitOperations(context.Background(), nil); err != nil {
		return err
	}

	return s.connectionPoolClose(s.pool)
}

func (s *SQLStore) WaitOperations(context context.Context, object any) error {
	var promises []*promise.Promise[any]
	s.mutexOperations.Lock()
	if object == nil {
		promises = make([]*promise.Promise[any], 0)
		for _, operations := range s.operations {
			for _, operation := range operations {
				promises = append(promises, operation.promise_)
			}
		}
	} else {
		operations := s.operations[object]
		promises = make([]*promise.Promise[any], 0, len(operations))
		for _, operation := range operations {
			promises = append(promises, operation.promise_)
		}
	}
	s.mutexOperations.Unlock()

	if len(promises) > 0 {
		allOperations := promise.All(context, promises...)
		_, err := allOperations.Await(context)
		return err
	}
	return nil
}

func (s *SQLStore) AsyncOperation(object any, operationType OperationType, op func() any) *promise.Promise[any] {
	s.mutexOperations.Lock()
	// compute previous operations
	previous := []*promise.Promise[any]{}
	operations := s.operations[object]
	switch operationType {
	case ReadOperation:
		for i := len(operations) - 1; i >= 0; i-- {
			operation := operations[i]
			if operation.type_ == WriteOperation {
				previous = append(previous, operation.promise_)
				break
			}
		}
	case WriteOperation:
		for i := len(operations) - 1; i >= 0; i-- {
			operation := operations[i]
			previous = append(previous, operation.promise_)
			if operation.type_ == WriteOperation {
				break
			}
		}
	}
	// create operation
	operation := &operation{
		type_: operationType,
		promise_: promise.New(func(resolve func(any), reject func(error)) {
			// wait for all previous promises
			if len(previous) > 0 {
				allPrevious := promise.All(context.Background(), previous...)
				if _, err := allPrevious.Await(context.Background()); err != nil {
					reject(err)
					return
				}
			}
			resolve(op())
		}),
	}
	// insertion index
	index := len(operations)
	// add operation
	s.operations[object] = append(operations, operation)
	s.mutexOperations.Unlock()
	// remove operation when finished with result or error
	return promise.New(func(resolve func(any), reject func(error)) {
		r, err := operation.promise_.Await(context.Background())
		s.mutexOperations.Lock()
		operations = s.operations[object]
		// remove operation from collection
		copy(operations[index:], operations[index+1:])
		operations[len(operations)-1] = nil
		s.operations[object] = operations[:len(operations)-1]
		s.mutexOperations.Unlock()
		if err != nil {
			reject(err)
		} else {
			resolve(r)
		}
	})
}
