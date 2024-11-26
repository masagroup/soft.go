package ecore

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

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

func newSQLStoreIDManager() *sqlStoreIDManagerImpl {
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

type SQLStore struct {
	*sqlBase
	sqlDecoder
	sqlEncoder
	mutex         sync.Mutex
	pool          *sqlitex.Pool
	errorHandler  func(error)
	sqlIDManager  *sqlStoreIDManagerImpl
	singleQueries map[*sqlColumn]*sqlSingleQueries
	manyQueries   map[*sqlTable]*sqlManyQueries
}

func NewSQLStore(databasePath string, resourceURI *URI, idManager EObjectIDManager, packageRegistry EPackageRegistry, options map[string]any) (*SQLStore, error) {
	// options
	schemaOptions := []sqlSchemaOption{withCreateIfNotExists(true)}
	idAttributeName := ""
	storeVersion := sqlCodecVersion
	errorHandler := func(error) {}
	sqlIDManager := newSQLStoreIDManager()
	sqlObjectManager := newSQLStoreObjectManager()
	if options != nil {
		idAttributeName, _ = options[SQL_OPTION_OBJECT_ID_NAME].(string)
		if idManager != nil && len(idAttributeName) > 0 {
			schemaOptions = append(schemaOptions, withIDAttributeName(idAttributeName))
		}
		if eh, isErrorHandler := options[SQL_OPTION_ERROR_HANDLER]; isErrorHandler {
			errorHandler = eh.(func(error))
		}
		if v, isVersion := options[SQL_OPTION_CODEC_VERSION].(int64); isVersion {
			storeVersion = v
		}
	}

	pool, err := sqlitex.NewPool(databasePath, sqlitex.PoolOptions{
		Flags: sqlite.OpenReadWrite | sqlite.OpenCreate | sqlite.OpenURI,
	})
	if err != nil {
		return nil, err
	}

	conn, err := pool.Take(context.Background())
	if err != nil {
		return nil, err
	}
	defer pool.Put(conn)

	// retrieve version
	var version int64
	if err := sqlitex.ExecuteTransient(conn, `PRAGMA user_version;`, &sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			version = stmt.ColumnInt64(0)
			return nil
		},
	}); err != nil {
		return nil, err
	}

	// encode version
	if version > 0 {
		if version != storeVersion {
			return nil, fmt.Errorf("history version %v is not supported", version)
		}
	} else {
		if err := sqlitex.ExecuteTransient(conn, fmt.Sprintf(`PRAGMA user_version = %v`, storeVersion), nil); err != nil {
			return nil, err
		}
	}
	// create sql base
	base := &sqlBase{
		uri:             resourceURI,
		idAttributeName: idAttributeName,
		idManager:       idManager,
		schema:          newSqlSchema(schemaOptions...),
	}

	// initialize
	// create sql store
	store := &SQLStore{
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
			classDataMap:     map[EClass]*sqlEncoderClassData{},
			sqlIDManager:     sqlIDManager,
			sqlObjectManager: sqlObjectManager,
		},
		pool:          pool,
		sqlIDManager:  sqlIDManager,
		errorHandler:  errorHandler,
		singleQueries: map[*sqlColumn]*sqlSingleQueries{},
		manyQueries:   map[*sqlTable]*sqlManyQueries{},
	}

	// set store in sql object manager
	sqlObjectManager.store = store

	// encode properties
	if err := store.encodePragmas(conn); err != nil {
		return nil, err
	}

	// encode schema
	if err := store.encodeSchema(conn); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *SQLStore) Close() error {
	return s.pool.Close()
}

func (s *SQLStore) getSingleQueries(column *sqlColumn) *sqlSingleQueries {
	s.mutex.Lock()
	stmts := s.singleQueries[column]
	if stmts == nil {
		stmts = &sqlSingleQueries{
			column: column,
		}
		s.singleQueries[column] = stmts
	}
	s.mutex.Unlock()
	return stmts
}

func (s *SQLStore) getManyQueries(table *sqlTable) *sqlManyQueries {
	s.mutex.Lock()
	stmts := s.manyQueries[table]
	if stmts == nil {
		stmts = &sqlManyQueries{
			table: table,
		}
		s.manyQueries[table] = stmts
	}
	s.mutex.Unlock()
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

func (s *SQLStore) getSQLID(conn *sqlite.Conn, eObject EObject) (int64, error) {
	// retrieve sql id for eObject
	sqlID, isSQLID := s.sqlIDManager.GetObjectID(eObject)
	if !isSQLID {
		// object is not in store - check if it exists in db
		objectExists := false
		objectsTable := s.schema.objectsTable
		if err := sqlitex.Execute(
			conn,
			objectsTable.selectQuery(nil, objectsTable.keyName()+"=?", ""),
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
			s.sqlIDManager.sqlDecoderIDManagerImpl.objects[sqlID] = eObject
		} else {
			// object doesn't exists in db - encode it
			return s.encodeObject(conn, eObject)
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
	if err := sqlitex.Execute(
		conn, query, &sqlitex.ExecOptions{
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

func (s *SQLStore) Set(object EObject, feature EStructuralFeature, index int, value any) any {
	// connection
	conn, err := s.pool.Take(context.Background())
	if err != nil {
		s.errorHandler(err)
		return nil
	}
	defer s.pool.Put(conn)

	// get object sql id
	sqlID, err := s.getSQLID(conn, object)
	if err != nil {
		s.errorHandler(err)
		return nil
	}

	// get encoder feature data
	featureData, err := s.getEncoderFeatureData(conn, object, feature)
	if err != nil {
		s.errorHandler(err)
		return nil
	}

	// retrieve previous value
	oldValue := s.getValue(conn, sqlID, featureData.schema, index)

	// encode value
	encoded, err := s.encodeFeatureValue(conn, featureData, value)
	if err != nil {
		s.errorHandler(err)
		return nil
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

	if err := sqlitex.Execute(conn, query, &sqlitex.ExecOptions{Args: args}); err != nil {
		s.errorHandler(err)
		return nil
	}

	return oldValue
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
		if err := sqlitex.Execute(
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
		if err := sqlitex.Execute(
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
	if err := sqlitex.Execute(conn, query, &sqlitex.ExecOptions{Args: args}); err != nil {
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
	if err := sqlitex.Execute(
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
	if err := sqlitex.Execute(
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
	if err := sqlitex.Execute(
		conn,
		s.getManyQueries(featureData.schema.table).getContainsQuery(),
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
	// compute parameters
	encoded, err := s.encodeFeatureValue(conn, featureData, value)
	if err != nil {
		s.errorHandler(err)
		return -1
	}
	// retrieve row idx in table
	idx := -1.0
	if err := sqlitex.Execute(
		conn,
		getIndexOfQuery(s.getManyQueries(featureData.schema.table)),
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
	if err := sqlitex.Execute(
		conn,
		s.getManyQueries(featureData.schema.table).getIdxToListIndexQuery(),
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

	if err := sqlitex.Execute(
		conn,
		s.getSingleQueries(s.schema.contentsTable.key).getRemoveQuery(),
		&sqlitex.ExecOptions{Args: []any{sqlID}},
	); err != nil {
		s.errorHandler(err)
	}
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
	idx, _, err := s.getInsertIdx(conn, featureData.schema.table, sqlID, index, 1)
	if err != nil {
		s.errorHandler(err)
		return
	}
	encoded, err := s.encodeFeatureValue(conn, featureData, value)
	if err != nil {
		s.errorHandler(err)
		return
	}
	if err := sqlitex.Execute(
		conn,
		s.getManyQueries(featureData.schema.table).getInsertQuery(),
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
	idx, inc, err := s.getInsertIdx(conn, featureData.schema.table, sqlID, index, c.Size())
	if err != nil {
		s.errorHandler(err)
		return
	}
	// encode each value
	query := s.getManyQueries(featureData.schema.table).getInsertQuery()
	for it := c.Iterator(); it.HasNext(); {
		value := it.Next()
		v, err := s.encodeFeatureValue(conn, featureData, value)
		if err != nil {
			s.errorHandler(err)
			return
		}
		if err := sqlitex.Execute(
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
		if err := sqlitex.Execute(
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
		if err := sqlitex.Execute(
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

	var value any
	if err := sqlitex.Execute(
		conn,
		s.getManyQueries(featureData.schema.table).getRemoveQuery(),
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
	// compute target index
	if targetIndex > sourceIndex {
		targetIndex++
	}
	idx, _, err := s.getInsertIdx(conn, featureSchema.table, sqlID, targetIndex, 1)
	if err != nil {
		s.errorHandler(err)
		return nil
	}

	// update idx of source index row with target idx
	var value any
	if err := sqlitex.Execute(
		conn,
		s.getManyQueries(featureSchema.table).getUpdateIdxQuery(),
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

	if err := sqlitex.Execute(
		conn,
		s.getManyQueries(featureTable).getClearQuery(),
		&sqlitex.ExecOptions{Args: []any{sqlID}},
	); err != nil {
		s.errorHandler(err)
	}
}

func (s *SQLStore) ToArray(object EObject, feature EStructuralFeature) []any {
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

	values := []any{}
	if err := sqlitex.Execute(
		conn,
		s.getManyQueries(featureSchema.table).getSelectAllQuery(),
		&sqlitex.ExecOptions{
			Args: []any{sqlID},
			ResultFunc: func(stmt *sqlite.Stmt) error {
				value := decodeAny(stmt, 0)
				decoded, err := s.decodeFeatureValue(conn, featureSchema, value)
				if err != nil {
					return err
				}
				values = append(values, decoded)
				return nil
			}}); err != nil {
		s.errorHandler(err)
		return nil
	}
	return values
}
