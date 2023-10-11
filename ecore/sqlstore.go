package ecore

import (
	"database/sql"
	"fmt"
)

type SQLStore struct {
	sqlEncoder
	errorHandler func(error)
}

func NewSQLStore(dbPath string, uri *URI, idManager EObjectIDManager, options map[string]any) (*SQLStore, error) {
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

	return &SQLStore{
		sqlEncoder: sqlEncoder{
			db:              db,
			uri:             uri,
			idAttributeName: idAttributeName,
			idManager:       idManager,
			schema:          newSqlSchema(schemaOptions...),
			insertStmts:     map[*sqlTable]*sql.Stmt{},
			classDataMap:    map[EClass]*sqlEncoderClassData{},
			packageIDs:      map[EPackage]int64{},
			objectIDs:       map[EObject]int64{},
			enumLiteralIDs:  map[string]int64{},
		},
		errorHandler: errorHandler,
	}, nil
}

func (s *SQLStore) Get(object EObject, feature EStructuralFeature, index int) any {
	return nil
}

func (s *SQLStore) Set(object EObject, feature EStructuralFeature, index int, value any) any {
	classData, err := s.getClassData(feature.GetEContainingClass())
	if err != nil {
		s.errorHandler(err)
		return nil
	}

	featureData, isFeatureData := classData.features[feature]
	if !isFeatureData {
		s.errorHandler(fmt.Errorf("feature %s is unknown", feature.GetName()))
		return nil
	}

	if featureColumn := featureData.schema.column; featureColumn != nil {
		// feature is encoded as a column

	} else if featureTable := featureData.schema.table; featureTable != nil {
		// feature is encoded in a external table
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
