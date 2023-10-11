package ecore

import (
	"database/sql"
	"fmt"
	"time"
)

type sqlStorePackageData struct {
	id int64
}

type sqlStoreClassData struct {
	id int64
}

type sqlStoreFeatureData struct {
	schema   *sqlFeatureSchema
	dataType EDataType
	factory  EFactory
}

type SQLStore struct {
	schema             *sqlSchema
	db                 *sql.DB
	uri                *URI
	errorHandler       func(error)
	insertStmts        map[*sqlTable]*sql.Stmt
	packageDataMap     map[EPackage]*sqlStorePackageData
	classDataMap       map[EClass]*sqlStoreClassData
	enumLiteralToIDMap map[string]int64
	featureDataMap     map[EStructuralFeature]*sqlStoreFeatureData
}

func NewSQLStore(dbPath string, uri *URI, options map[string]any) (*SQLStore, error) {
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
		if len(idAttributeName) > 0 {
			schemaOptions = append(schemaOptions, withIDAttributeName(idAttributeName))
		}

		if eh, isErrorHandler := options[SQL_OPTION_DRIVER]; isErrorHandler {
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
		schema:       newSqlSchema(schemaOptions...),
		db:           db,
		uri:          uri,
		errorHandler: errorHandler,
	}, nil
}

func (s *SQLStore) Get(object EObject, feature EStructuralFeature, index int) any {
	return nil
}

func (s *SQLStore) Set(object EObject, feature EStructuralFeature, index int, value any) any {
	featureData, err := s.getFeatureData(feature)
	if err != nil {
		s.errorHandler(err)
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

func (s *SQLStore) getFeatureData(feature EStructuralFeature) (*sqlStoreFeatureData, error) {
	featureData := s.featureDataMap[feature]
	if featureData == nil {
		class := feature.GetEContainingClass()
		classSchema := s.schema.getClassSchema(class)
		featureSchema := classSchema.getFeatureSchema(feature)

		// create feature table if any
		if table := featureSchema.table; table != nil {
			if _, err := s.db.Exec(table.createQuery()); err != nil {
				return nil, err
			}
		}

		// create & initialize feature data
		featureData = &sqlStoreFeatureData{
			schema: featureSchema,
		}
		if eAttribute, _ := feature.(EAttribute); eAttribute != nil {
			eDataType := eAttribute.GetEAttributeType()
			featureData.dataType = eDataType
			featureData.factory = eDataType.GetEPackage().GetEFactoryInstance()
		}

		// register feature data
		s.featureDataMap[feature] = featureData
	}
	return featureData, nil
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

func (s *SQLStore) encodeFeatureValue(featureData *sqlStoreFeatureData, value any) (any, error) {
	if value != nil {
		switch featureData.schema.featureKind {
		case sfkObject, sfkObjectList:
			objectData, err := s.encodeObject(value.(EObject))
			if err != nil {
				return nil, err
			}
			return objectData.id, nil
		case sfkObjectReference, sfkObjectReferenceList:
			ref := GetURI(value.(EObject))
			uri := s.uri.Relativize(ref)
			return uri.String(), nil
		case sfkEnum:
			eEnum := featureData.dataType.(EEnum)
			literal := featureData.factory.ConvertToString(featureData.dataType, value)
			return s.encodeEnumLiteral(eEnum, literal)
		case sfkBool, sfkByte, sfkInt, sfkInt16, sfkInt32, sfkInt64, sfkString, sfkByteArray, sfkFloat32, sfkFloat64:
			return value, nil
		case sfkDate:
			t := value.(*time.Time)
			return t.Format(time.RFC3339), nil
		case sfkData, sfkDataList:
			return featureData.factory.ConvertToString(featureData.dataType, value), nil
		}
	}
	return nil, nil
}

func (s *SQLStore) encodeEnumLiteral(eEnum EEnum, literal string) (any, error) {
	if enumID, isEnumID := s.enumLiteralToIDMap[literal]; isEnumID {
		return enumID, nil
	} else {
		ePackage := eEnum.GetEPackage()
		packageData, err := s.encodePackage(ePackage)
		if err != nil {
			return nil, err
		}
		// insert enum in sql
		insertEnumStmt, err := s.getInsertStmt(s.schema.enumsTable)
		if err != nil {
			return nil, err
		}
		sqlResult, err := insertEnumStmt.Exec(packageData.id, eEnum.GetName(), literal)
		if err != nil {
			return nil, err
		}

		// retrieve enum index
		enumID, err := sqlResult.LastInsertId()
		if err != nil {
			return nil, err
		}
		s.enumLiteralToIDMap[literal] = enumID
		return enumID, nil
	}
}

func (s *SQLStore) getInsertStmt(table *sqlTable) (stmt *sql.Stmt, err error) {
	stmt = s.insertStmts[table]
	if stmt == nil {
		stmt, err = s.db.Prepare(table.insertQuery())
	}
	return
}
