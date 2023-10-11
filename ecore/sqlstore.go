package ecore

import (
	"database/sql"
	"fmt"
)

type sqlStorePackageData struct {
	id int64
}

type sqlStoreClassData struct {
	id int64
}

type SQLStore struct {
	schema         *sqlSchema
	db             *sql.DB
	insertStmts    map[*sqlTable]*sql.Stmt
	packageDataMap map[EPackage]*sqlStorePackageData
	classDataMap   map[EClass]*sqlStoreClassData
}

func NewSQLStore(driver string, dbPath string, options map[string]any) (*SQLStore, error) {
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

	// schema options
	schemaOptions := []sqlSchemaOption{withCreateIfNotExists(true)}
	idAttributeName, _ := options[SQL_OPTION_ID_ATTRIBUTE_NAME].(string)
	if len(idAttributeName) > 0 {
		schemaOptions = append(schemaOptions, withIDAttributeName(idAttributeName))
	}

	return &SQLStore{
		schema: newSqlSchema(schemaOptions...),
		db:     db,
	}, nil
}

func (s *SQLStore) Get(object EObject, feature EStructuralFeature, index int) any {
	return nil
}

func (s *SQLStore) Set(object EObject, feature EStructuralFeature, index int, value any) any {
	class := feature.GetEContainingClass()
	classSchema := s.schema.getClassSchema(class)
	featureSchema := classSchema.getFeatureSchema(feature)

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

func (s *SQLStore) encodePackage(ePackage EPackage) (*sqlStorePackageData, error) {
	ePackageData := s.packageDataMap[ePackage]
	if ePackageData == nil {
		// insert new package
		insertPackageStmt, err := s.getInsertStmt(s.schema.packagesTable)
		if err != nil {
			return nil, err
		}
		sqlResult, err := insertPackageStmt.Exec(ePackage.GetNsURI())
		if err != nil {
			return nil, err
		}
		// retrieve package index
		id, err := sqlResult.LastInsertId()
		if err != nil {
			return nil, err
		}
		// create data
		ePackageData = &sqlStorePackageData{id: id}
		s.packageDataMap[ePackage] = ePackageData
	}
	return ePackageData, nil
}

func (s *SQLStore) encodeClass(eClass EClass) (*sqlStoreClassData, error) {
	eClassData := s.classDataMap[eClass]
	if eClassData == nil {
		// encode package
		ePackage := eClass.GetEPackage()
		packageData, err := s.encodePackage(ePackage)
		if err != nil {
			return nil, err
		}

		// insert class in sql
		insertClassStmt, err := s.getInsertStmt(s.schema.classesTable)
		if err != nil {
			return nil, err
		}
		sqlResult, err := insertClassStmt.Exec(packageData.id, eClass.GetName())
		if err != nil {
			return nil, err
		}

		// retrieve class index
		id, err := sqlResult.LastInsertId()
		if err != nil {
			return nil, err
		}

		eClassData = &sqlStoreClassData{id: id}
		s.classDataMap[eClass] = eClassData
	}
	return eClassData, nil
}

func (s *SQLStore) getInsertStmt(table *sqlTable) (stmt *sql.Stmt, err error) {
	stmt = s.insertStmts[table]
	if stmt == nil {
		stmt, err = s.db.Prepare(table.insertQuery())
	}
	return
}
