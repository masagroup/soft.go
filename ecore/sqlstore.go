package ecore

import (
	"database/sql"
	"fmt"
	"strings"
)

type SQLStore struct {
	*sqlBase
	sqlDecoder
	sqlEncoder
	errorHandler func(error)
	updateStmts  map[*sqlColumn]*sql.Stmt
	selectStmts  map[*sqlColumn]*sql.Stmt
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
		updateStmts:  map[*sqlColumn]*sql.Stmt{},
		selectStmts:  map[*sqlColumn]*sql.Stmt{},
	}, nil
}

func (s *SQLStore) Close() error {
	return s.db.Close()
}

func (s *SQLStore) getUpdateStmt(column *sqlColumn, query func() string) (stmt *sql.Stmt, err error) {
	stmt = s.updateStmts[column]
	if stmt == nil {
		stmt, err = s.db.Prepare(query())
	}
	return
}

func (s *SQLStore) getSelectStmt(column *sqlColumn, query func() string) (stmt *sql.Stmt, err error) {
	stmt = s.selectStmts[column]
	if stmt == nil {
		stmt, err = s.db.Prepare(query())
	}
	return
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

	var row *sql.Row
	if featureColumn := featureSchema.column; featureColumn != nil {
		stmt, err := s.getSelectStmt(featureColumn, func() string {
			var query strings.Builder
			query.WriteString("SELECT ")
			query.WriteString(sqlEscapeIdentifier(featureColumn.columnName))
			query.WriteString(" FROM ")
			query.WriteString(sqlEscapeIdentifier(featureColumn.table.name))
			query.WriteString(" WHERE ")
			query.WriteString(featureColumn.table.keyName())
			query.WriteString("=?")
			return query.String()
		})
		if err != nil {
			s.errorHandler(err)
			return nil
		}
		row = stmt.QueryRow(sqlID)

	} else if featureTable := featureSchema.table; featureTable != nil {
		featureColumn := featureTable.columns[len(featureTable.columns)-1]
		stmt, err := s.getSelectStmt(featureColumn, func() string {
			var query strings.Builder
			query.WriteString("SELECT ")
			query.WriteString(sqlEscapeIdentifier(featureColumn.columnName))
			query.WriteString(" FROM ")
			query.WriteString(sqlEscapeIdentifier(featureColumn.table.name))
			query.WriteString(" WHERE ")
			query.WriteString(featureColumn.table.keyName())
			query.WriteString("=? ORDER BY ")
			query.WriteString(featureTable.keyName())
			query.WriteString(" ASC, idx ASC LIMIT 1 OFFSET ?")
			return query.String()
		})
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

	// feature is encoded as a column
	v, err := s.encodeFeatureValue(featureData, value)
	if err != nil {
		s.errorHandler(err)
		return nil
	}

	if featureColumn := featureData.schema.column; featureColumn != nil {

		stmt, err := s.getUpdateStmt(featureColumn, func() string {
			var query strings.Builder
			query.WriteString("UPDATE ")
			query.WriteString(sqlEscapeIdentifier(featureColumn.table.name))
			query.WriteString(" SET ")
			query.WriteString(sqlEscapeIdentifier(featureColumn.columnName))
			query.WriteString("=? WHERE ")
			query.WriteString(featureColumn.table.keyName())
			query.WriteString("=?")
			return query.String()
		})
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
		featureColumn := featureTable.columns[len(featureTable.columns)-1]
		stmt, err := s.getUpdateStmt(featureColumn, func() string {
			var query strings.Builder
			query.WriteString("UPDATE ")
			query.WriteString(sqlEscapeIdentifier(featureTable.name))
			query.WriteString(" SET ")
			query.WriteString(sqlEscapeIdentifier(featureTable.columns[len(featureTable.columns)-1].columnName))
			query.WriteString("=? WHERE rowid IN (SELECT rowid FROM ")
			query.WriteString(sqlEscapeIdentifier(featureTable.name))
			query.WriteString(" WHERE ")
			query.WriteString(featureTable.keyName())
			query.WriteString("=?")
			query.WriteString(" ORDER BY ")
			query.WriteString(featureTable.keyName())
			query.WriteString(" ASC, idx ASC LIMIT 1 OFFSET ?)")
			return query.String()
		})
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
	return nil
}

func (s *SQLStore) IsSet(object EObject, feature EStructuralFeature) bool {
	return false
}

func (s *SQLStore) UnSet(object EObject, feature EStructuralFeature) {

}

func (s *SQLStore) IsEmpty(object EObject, feature EStructuralFeature) bool {
	return false
}

func (s *SQLStore) Size(object EObject, feature EStructuralFeature) int {
	return 0
}

func (s *SQLStore) Contains(object EObject, feature EStructuralFeature, value any) bool {
	return false
}

func (s *SQLStore) IndexOf(object EObject, feature EStructuralFeature, value any) int {
	return 0
}

func (s *SQLStore) LastIndexOf(object EObject, feature EStructuralFeature, value any) int {
	return 0
}

func (s *SQLStore) Add(object EObject, feature EStructuralFeature, index int, value any) {
}

func (s *SQLStore) Remove(object EObject, feature EStructuralFeature, index int) any {
	return nil
}

func (s *SQLStore) Move(object EObject, feature EStructuralFeature, targetIndex int, sourceIndex int) any {
	return nil
}

func (s *SQLStore) Clear(object EObject, feature EStructuralFeature) {

}
