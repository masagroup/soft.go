package ecore

import (
	"database/sql"
	"fmt"
	"strings"
)

type stmtOrError struct {
	stmt *sql.Stmt
	err  error
}

type sqlSingleStmts struct {
	db         *sql.DB
	column     *sqlColumn
	updateStmt *stmtOrError
	selectStmt *stmtOrError
}

func (ss *sqlSingleStmts) getUpdateStmt() (*sql.Stmt, error) {
	if ss.updateStmt == nil {
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
		// stmt
		ss.updateStmt = &stmtOrError{}
		ss.updateStmt.stmt, ss.updateStmt.err = ss.db.Prepare(query.String())
	}
	return ss.updateStmt.stmt, ss.updateStmt.err
}

func (ss *sqlSingleStmts) getSelectStmt() (*sql.Stmt, error) {
	if ss.selectStmt == nil {
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
		// stmt
		ss.selectStmt = &stmtOrError{}
		ss.selectStmt.stmt, ss.selectStmt.err = ss.db.Prepare(query.String())
	}
	return ss.selectStmt.stmt, ss.selectStmt.err
}

type sqlManyStmts struct {
	db           *sql.DB
	table        *sqlTable
	updateStmt   *stmtOrError
	selectStmt   *stmtOrError
	existsStmt   *stmtOrError
	clearStmt    *stmtOrError
	countStmt    *stmtOrError
	containsStmt *stmtOrError
	indexOfStmt  *stmtOrError
}

func (ss *sqlManyStmts) getUpdateStmt() (*sql.Stmt, error) {
	if ss.updateStmt == nil {
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
		// stmt
		ss.updateStmt = &stmtOrError{}
		ss.updateStmt.stmt, ss.updateStmt.err = ss.db.Prepare(query.String())
	}
	return ss.updateStmt.stmt, ss.updateStmt.err
}

func (ss *sqlManyStmts) getSelectStmt() (*sql.Stmt, error) {
	if ss.selectStmt == nil {
		// query
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
		// stmt
		ss.selectStmt = &stmtOrError{}
		ss.selectStmt.stmt, ss.selectStmt.err = ss.db.Prepare(query.String())
	}
	return ss.selectStmt.stmt, ss.selectStmt.err
}

func (ss *sqlManyStmts) getExistsStmt() (*sql.Stmt, error) {
	if ss.existsStmt == nil {
		// query
		var query strings.Builder
		query.WriteString("SELECT 1 FROM")
		query.WriteString(sqlEscapeIdentifier(ss.table.name))
		query.WriteString(" WHERE ")
		query.WriteString(ss.table.keyName())
		query.WriteString("=?")
		// stmt
		ss.existsStmt = &stmtOrError{}
		ss.existsStmt.stmt, ss.existsStmt.err = ss.db.Prepare(query.String())
	}
	return ss.existsStmt.stmt, ss.existsStmt.err
}

func (ss *sqlManyStmts) getClearStmt() (*sql.Stmt, error) {
	if ss.clearStmt == nil {
		// query
		var query strings.Builder
		query.WriteString("DELETE FROM")
		query.WriteString(sqlEscapeIdentifier(ss.table.name))
		query.WriteString(" WHERE ")
		query.WriteString(ss.table.keyName())
		query.WriteString("=?")
		// stmt
		ss.clearStmt = &stmtOrError{}
		ss.clearStmt.stmt, ss.clearStmt.err = ss.db.Prepare(query.String())
	}
	return ss.clearStmt.stmt, ss.clearStmt.err
}

func (ss *sqlManyStmts) getCountStmt() (*sql.Stmt, error) {
	if ss.countStmt == nil {
		// query
		var query strings.Builder
		query.WriteString("SELECT COUNT(*) FROM")
		query.WriteString(sqlEscapeIdentifier(ss.table.name))
		query.WriteString(" WHERE ")
		query.WriteString(ss.table.keyName())
		query.WriteString("=?")
		// stmt
		ss.countStmt = &stmtOrError{}
		ss.countStmt.stmt, ss.countStmt.err = ss.db.Prepare(query.String())
	}
	return ss.countStmt.stmt, ss.countStmt.err
}

func (ss *sqlManyStmts) getContainsStmt() (*sql.Stmt, error) {
	if ss.containsStmt == nil {
		column := ss.table.columns[len(ss.table.columns)-1]
		// query
		var query strings.Builder
		query.WriteString("SELECT rowid FROM")
		query.WriteString(sqlEscapeIdentifier(ss.table.name))
		query.WriteString(" WHERE ")
		query.WriteString(ss.table.keyName())
		query.WriteString("=? AND ")
		query.WriteString(sqlEscapeIdentifier(column.columnName))
		query.WriteString("=?")
		// stmt
		ss.containsStmt = &stmtOrError{}
		ss.containsStmt.stmt, ss.containsStmt.err = ss.db.Prepare(query.String())
	}
	return ss.containsStmt.stmt, ss.containsStmt.err
}

func (ss *sqlManyStmts) getIndexOfStmt() (*sql.Stmt, error) {
	if ss.indexOfStmt == nil {
		column := ss.table.columns[len(ss.table.columns)-1]
		// query
		var query strings.Builder
		query.WriteString("SELECT COUNT(*) WHERE ")
		query.WriteString(ss.table.keyName())
		query.WriteString("=? AND idx < (SELECT idx FROM ")
		query.WriteString(sqlEscapeIdentifier(ss.table.name))
		query.WriteString(" WHERE ")
		query.WriteString(ss.table.keyName())
		query.WriteString("=? AND ")
		query.WriteString(sqlEscapeIdentifier(column.columnName))
		query.WriteString("=?)")
		// stmt
		ss.indexOfStmt = &stmtOrError{}
		ss.indexOfStmt.stmt, ss.indexOfStmt.err = ss.db.Prepare(query.String())
	}
	return ss.indexOfStmt.stmt, ss.indexOfStmt.err
}

type SQLStore struct {
	*sqlBase
	sqlDecoder
	sqlEncoder
	errorHandler func(error)
	singleStmts  map[*sqlColumn]*sqlSingleStmts
	manyStmts    map[*sqlTable]*sqlManyStmts
}

func NewSQLStore(dbPath string, uri *URI, idManager EObjectIDManager, packageRegistry EPackageRegistry, options map[string]any) (*SQLStore, error) {
	// options
	schemaOptions := []sqlSchemaOption{withCreateIfNotExists(true)}
	driver := "sqlite"
	idAttributeName := ""
	errorHandler := func(error) {}
	if options != nil {
		if d, isDriver := options[SQL_OPTION_DRIVER]; isDriver {
			driver = d.(string)
		}

		idAttributeName, _ = options[SQL_OPTION_ID_ATTRIBUTE_NAME].(string)
		if idManager != nil && len(idAttributeName) > 0 {
			schemaOptions = append(schemaOptions, withIDAttributeName(idAttributeName))
		}

		if eh, isErrorHandler := options[SQL_OPTION_ERROR_HANDLER]; isErrorHandler {
			errorHandler = eh.(func(error))
		}
	}

	// Open db
	db, err := sql.Open(driver, dbPath)
	if err != nil {
		return nil, err
	}

	// properties
	propertiesQuery := `
	PRAGMA synchronous = NORMAL;
	PRAGMA journal_mode = WAL;
	`
	_, err = db.Exec(propertiesQuery)
	if err != nil {
		return nil, err
	}

	// version
	if row := db.QueryRow("PRAGMA user_version;"); row == nil {
		// create version
		versionQuery := fmt.Sprintf(`PRAGMA user_version = %v`, sqlCodecVersion)
		_, err = db.Exec(versionQuery)
		if err != nil {
			return nil, err
		}
	} else {
		// retrieve version
		var v int
		if err := row.Scan(&v); err != nil {
			return nil, err
		}
		if v != sqlCodecVersion {
			return nil, fmt.Errorf("history version %v is not supported", v)
		}
	}

	// create sql base
	base := &sqlBase{
		db:              db,
		uri:             uri,
		idAttributeName: idAttributeName,
		idManager:       idManager,
		schema:          newSqlSchema(schemaOptions...),
	}

	// create sql store
	return &SQLStore{
		sqlBase: base,
		sqlDecoder: sqlDecoder{
			sqlBase:         base,
			packageRegistry: packageRegistry,
			packages:        map[int64]EPackage{},
			objects:         map[int64]EObject{},
			classes:         map[int64]*sqlDecoderClassData{},
			enums:           map[int64]any{},
			selectStmts:     map[*sqlTable]*sql.Stmt{},
		},
		sqlEncoder: sqlEncoder{
			sqlBase:        base,
			insertStmts:    map[*sqlTable]*sql.Stmt{},
			classDataMap:   map[EClass]*sqlEncoderClassData{},
			packageIDs:     map[EPackage]int64{},
			objectIDs:      map[EObject]int64{},
			enumLiteralIDs: map[string]int64{},
		},
		errorHandler: errorHandler,
		singleStmts:  map[*sqlColumn]*sqlSingleStmts{},
		manyStmts:    map[*sqlTable]*sqlManyStmts{},
	}, nil
}

func (s *SQLStore) Close() error {
	return s.db.Close()
}

func (s *SQLStore) getSingleStmts(column *sqlColumn) *sqlSingleStmts {
	stmts := s.singleStmts[column]
	if stmts == nil {
		stmts = &sqlSingleStmts{
			db:     s.db,
			column: column,
		}
		s.singleStmts[column] = stmts
	}
	return stmts
}

func (s *SQLStore) getManyStmts(table *sqlTable) *sqlManyStmts {
	stmts := s.manyStmts[table]
	if stmts == nil {
		stmts = &sqlManyStmts{
			db:    s.db,
			table: table,
		}
		s.manyStmts[table] = stmts
	}
	return stmts
}

func (s *SQLStore) Get(object EObject, feature EStructuralFeature, index int) any {
	sqlObject := object.(SQLObject)
	sqlID := sqlObject.GetSqlID()

	// retrieve class schema
	classSchema := s.sqlDecoder.schema.getClassSchema(object.EClass())

	// retrieve feature schema
	featureSchema := classSchema.getFeatureSchema(feature)
	if featureSchema == nil {
		s.errorHandler(fmt.Errorf("feature %s is unknown", feature.GetName()))
		return nil
	}

	return s.getValue(sqlID, featureSchema, index)
}

func (s *SQLStore) getValue(sqlID int64, featureSchema *sqlFeatureSchema, index int) any {
	var row *sql.Row
	if featureColumn := featureSchema.column; featureColumn != nil {
		stmt, err := s.getSingleStmts(featureColumn).getSelectStmt()
		if err != nil {
			s.errorHandler(err)
			return nil
		}
		row = stmt.QueryRow(sqlID)

	} else if featureTable := featureSchema.table; featureTable != nil {
		stmt, err := s.getManyStmts(featureTable).getSelectStmt()
		if err != nil {
			s.errorHandler(err)
			return nil
		}
		row = stmt.QueryRow(sqlID, index)

	}

	var v any
	if err := row.Scan(&v); err != nil {
		if err != sql.ErrNoRows {
			s.errorHandler(err)
		}
		return nil
	}

	value, err := s.decodeFeatureValue(featureSchema, v)
	if err != nil {
		s.errorHandler(err)
	}

	return value
}

func (s *SQLStore) Set(object EObject, feature EStructuralFeature, index int, value any) any {
	sqlObject := object.(SQLObject)
	sqlID := sqlObject.GetSqlID()

	// retrieve class data
	classData, err := s.getClassData(object.EClass())
	if err != nil {
		s.errorHandler(err)
		return nil
	}

	// retrieve feature data
	featureData, isFeatureData := classData.features[feature]
	if !isFeatureData {
		s.errorHandler(fmt.Errorf("feature %s is unknown", feature.GetName()))
		return nil
	}

	// retrieve previous value
	oldValue := s.getValue(sqlID, featureData.schema, index)

	// encode value
	v, err := s.encodeFeatureValue(featureData, value)
	if err != nil {
		s.errorHandler(err)
		return nil
	}

	if featureColumn := featureData.schema.column; featureColumn != nil {
		stmt, err := s.getSingleStmts(featureColumn).getUpdateStmt()
		if err != nil {
			s.errorHandler(err)
			return nil
		}
		_, err = stmt.Exec(v, sqlID)
		if err != nil {
			s.errorHandler(err)
			return nil
		}

	} else if featureTable := featureData.schema.table; featureTable != nil {
		stmt, err := s.getManyStmts(featureTable).getUpdateStmt()
		if err != nil {
			s.errorHandler(err)
			return nil
		}
		_, err = stmt.Exec(v, sqlID, index)
		if err != nil {
			s.errorHandler(err)
			return nil
		}
	}

	return oldValue
}

func (s *SQLStore) IsSet(object EObject, feature EStructuralFeature) bool {
	sqlObject := object.(SQLObject)
	sqlID := sqlObject.GetSqlID()

	// retrieve class schema
	classSchema := s.sqlDecoder.schema.getClassSchema(object.EClass())

	// retrieve feature schema
	featureSchema := classSchema.getFeatureSchema(feature)
	if featureSchema == nil {
		s.errorHandler(fmt.Errorf("feature %s is unknown", feature.GetName()))
		return false
	}

	if featureColumn := featureSchema.column; featureColumn != nil {
		stmt, err := s.getSingleStmts(featureColumn).getSelectStmt()
		if err != nil {
			s.errorHandler(err)
			return false
		}
		row := stmt.QueryRow(sqlID)
		var v any
		if err := row.Scan(&v); err != nil {
			return false
		}
		return v != featureSchema.feature.GetDefaultValue()

	} else if featureTable := featureSchema.table; featureTable != nil {
		stmt, err := s.getManyStmts(featureTable).getExistsStmt()
		if err != nil {
			s.errorHandler(err)
			return false
		}
		var v any
		row := stmt.QueryRow(sqlID)
		_ = row.Scan(&v)
		return v != nil
	}

	return false
}

func (s *SQLStore) UnSet(object EObject, feature EStructuralFeature) {
	sqlObject := object.(SQLObject)
	sqlID := sqlObject.GetSqlID()

	// retrieve class data
	classData, err := s.getClassData(object.EClass())
	if err != nil {
		s.errorHandler(err)
		return
	}

	// retrieve feature data
	featureData, isFeatureData := classData.features[feature]
	if !isFeatureData {
		s.errorHandler(fmt.Errorf("feature %s is unknown", feature.GetName()))
		return
	}

	if featureColumn := featureData.schema.column; featureColumn != nil {
		stmt, err := s.getSingleStmts(featureColumn).getUpdateStmt()
		if err != nil {
			s.errorHandler(err)
			return
		}
		v := feature.GetDefaultValue()
		_, err = stmt.Exec(v, sqlID)
		if err != nil {
			s.errorHandler(err)
			return
		}
	} else if featureTable := featureData.schema.table; featureTable != nil {
		stmt, err := s.getManyStmts(featureTable).getClearStmt()
		if err != nil {
			s.errorHandler(err)
			return
		}
		_, err = stmt.Exec(sqlID)
		if err != nil {
			s.errorHandler(err)
			return
		}
	}
}

func (s *SQLStore) getFeatureTable(object EObject, feature EStructuralFeature) (*sqlTable, error) {
	if !feature.IsMany() {
		panic(fmt.Sprintf("%s is not a many feature", feature.GetName()))
	}
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
	return featureSchema.table, nil
}

func (s *SQLStore) getFeatureData(object EObject, feature EStructuralFeature) (*sqlEncoderFeatureData, error) {
	if !feature.IsMany() {
		panic(fmt.Sprintf("%s is not a many feature", feature.GetName()))
	}
	// retrieve class schema
	class := object.EClass()

	// retrieve class data
	classData, err := s.getClassData(class)
	if err != nil {
		s.errorHandler(err)
		return nil, fmt.Errorf("class %s is unknown", class.GetName())
	}

	// retrieve feature data
	featureData, isFeatureData := classData.features[feature]
	if !isFeatureData {
		err := fmt.Errorf("feature %s is unknown", feature.GetName())
		s.errorHandler(err)
		return nil, err
	}

	return featureData, nil
}

func (s *SQLStore) IsEmpty(object EObject, feature EStructuralFeature) bool {
	// retrieve table
	featureTable, err := s.getFeatureTable(object, feature)
	if err != nil {
		s.errorHandler(err)
		return false
	}

	// retrieve statement
	stmt, err := s.getManyStmts(featureTable).getExistsStmt()
	if err != nil {
		s.errorHandler(err)
		return false
	}

	// query statement
	sqlObject := object.(SQLObject)
	sqlID := sqlObject.GetSqlID()
	var v any
	row := stmt.QueryRow(sqlID)
	_ = row.Scan(&v)
	return v == nil
}

func (s *SQLStore) Size(object EObject, feature EStructuralFeature) int {
	// retrieve table
	featureTable, err := s.getFeatureTable(object, feature)
	if err != nil {
		s.errorHandler(err)
		return 0
	}

	// retrieve statement
	stmt, err := s.getManyStmts(featureTable).getCountStmt()
	if err != nil {
		s.errorHandler(err)
		return 0
	}

	// query statement
	sqlObject := object.(SQLObject)
	sqlID := sqlObject.GetSqlID()
	row := stmt.QueryRow(sqlID)

	// retrieve count
	var count int
	_ = row.Scan(&count)
	return count
}

func (s *SQLStore) Contains(object EObject, feature EStructuralFeature, value any) bool {
	// retrieve table
	featureData, err := s.getFeatureData(object, feature)
	if err != nil {
		s.errorHandler(err)
		return false
	}
	// retrieve statement
	stmt, err := s.getManyStmts(featureData.schema.table).getContainsStmt()
	if err != nil {
		s.errorHandler(err)
		return false
	}

	// query statement
	sqlObject := object.(SQLObject)
	sqlID := sqlObject.GetSqlID()
	v, err := s.encodeFeatureValue(featureData, value)
	if err != nil {
		s.errorHandler(err)
		return false
	}
	row := stmt.QueryRow(sqlID, v)
	var rowid int
	_ = row.Scan(&rowid)
	return rowid != 0
}

func (s *SQLStore) IndexOf(object EObject, feature EStructuralFeature, value any) int {
	// retrieve table
	featureData, err := s.getFeatureData(object, feature)
	if err != nil {
		s.errorHandler(err)
		return -1
	}
	// retrieve statement
	stmt, err := s.getManyStmts(featureData.schema.table).getIndexOfStmt()
	if err != nil {
		s.errorHandler(err)
		return -1
	}
	sqlObject := object.(SQLObject)
	sqlID := sqlObject.GetSqlID()
	v, err := s.encodeFeatureValue(featureData, value)
	if err != nil {
		s.errorHandler(err)
		return -1
	}
	row := stmt.QueryRow(sqlID, v)
	var index int
	_ = row.Scan(&index)
	return index
}

func (s *SQLStore) LastIndexOf(object EObject, feature EStructuralFeature, value any) int {
	if !feature.IsMany() {
		panic(fmt.Sprintf("%s is not a many feature", feature.GetName()))
	}
	return 0
}

func (s *SQLStore) Add(object EObject, feature EStructuralFeature, index int, value any) {
	if !feature.IsMany() {
		panic(fmt.Sprintf("%s is not a many feature", feature.GetName()))
	}
}

func (s *SQLStore) Remove(object EObject, feature EStructuralFeature, index int) any {
	if !feature.IsMany() {
		panic(fmt.Sprintf("%s is not a many feature", feature.GetName()))
	}
	return nil
}

func (s *SQLStore) Move(object EObject, feature EStructuralFeature, targetIndex int, sourceIndex int) any {
	if !feature.IsMany() {
		panic(fmt.Sprintf("%s is not a many feature", feature.GetName()))
	}
	return nil
}

func (s *SQLStore) Clear(object EObject, feature EStructuralFeature) {
	if !feature.IsMany() {
		panic(fmt.Sprintf("%s is not a many feature", feature.GetName()))
	}
}
