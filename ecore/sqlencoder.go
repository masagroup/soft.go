package ecore

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type sqlEncoderObjectData struct {
	id        int64
	classData *sqlEncoderClassData
}

type sqlEncoderFeatureData struct {
	schema   *sqlFeatureSchema
	dataType EDataType
	factory  EFactory
}

type sqlEncoderClassData struct {
	id        int64
	schema    *sqlClassSchema
	hierarchy []EClass
	features  map[EStructuralFeature]*sqlEncoderFeatureData
}

type sqlEncoderPackageData struct {
	id int64
}

func newSqlEncoderPackageData(id int64) *sqlEncoderPackageData {
	return &sqlEncoderPackageData{
		id: id,
	}
}

type sqlStmt struct {
	stmt *sql.Stmt
	args []any
}

type sqlStmts struct {
	db    *sql.DB
	stmts []*sqlStmt
}

func newSqlStmts(db *sql.DB) *sqlStmts {
	return &sqlStmts{db: db}
}

func (s *sqlStmts) add(stmt *sql.Stmt, args ...any) {
	s.stmts = append(s.stmts, &sqlStmt{stmt: stmt, args: args})
}

func (s *sqlStmts) exec() error {
	tx, err := s.db.Begin()
	if err != nil {
		return nil
	}
	txStmts := map[*sql.Stmt]*sql.Stmt{}
	for _, t := range s.stmts {
		stmt := t.stmt
		txStmt := txStmts[stmt]
		if txStmt == nil {
			txStmt = tx.Stmt(stmt)
			txStmts[stmt] = txStmt
		}
		_, err := txStmt.Exec(t.args...)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

type SQLEncoder struct {
	resource        EResource
	writer          io.Writer
	driver          string
	db              *sql.DB
	schema          *sqlSchema
	insertStmts     map[*sqlTable]*sql.Stmt
	packageDataMap  map[EPackage]*sqlEncoderPackageData
	classDataMap    map[EClass]*sqlEncoderClassData
	idAttributeName string
}

func NewSQLEncoder(resource EResource, w io.Writer, options map[string]any) *SQLEncoder {

	// options
	schemaOptions := []sqlSchemaOption{}
	driver := "sqlite"
	if options != nil {
		if driver, isDriver := options[SQL_OPTION_DRIVER]; isDriver {
			driver = driver.(string)
		}

		idAttributeName, _ := options[JSON_OPTION_ID_ATTRIBUTE_NAME].(string)
		if resource.GetObjectIDManager() != nil && len(idAttributeName) > 0 {
			schemaOptions = append(schemaOptions, withIDAttributeName(idAttributeName))
		}
	}

	// encoder structure
	return &SQLEncoder{
		resource:       resource,
		writer:         w,
		driver:         driver,
		schema:         newSqlSchema(schemaOptions...),
		packageDataMap: map[EPackage]*sqlEncoderPackageData{},
		classDataMap:   map[EClass]*sqlEncoderClassData{},
		insertStmts:    map[*sqlTable]*sql.Stmt{},
	}
}

func (e *SQLEncoder) createDB(dbPath string) (*sql.DB, error) {

	// open db
	db, err := sql.Open(e.driver, dbPath)
	if err != nil {
		return nil, err
	}

	// version
	versionQuery := fmt.Sprintf(`PRAGMA user_version = %v`, sqlCodecVersion)
	_, err = db.Exec(versionQuery)
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

	// tables
	for _, table := range []*sqlTable{
		e.schema.packagesTable,
		e.schema.classesTable,
		e.schema.objectsTable,
		e.schema.contentsTable,
	} {
		if _, err := db.Exec(table.createQuery()); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func (e *SQLEncoder) EncodeResource() {
	// create a temp file for the database file
	fileName := filepath.Base(e.resource.GetURI().Path())
	dbPath, err := sqlTmpDB(fileName)
	if err != nil {
		e.addError(err)
		return
	}

	// create db
	e.db, err = e.createDB(dbPath)
	if err != nil {
		e.addError(err)
		return
	}
	defer func() {
		_ = e.db.Close()
	}()

	// encode contents into db
	if contents := e.resource.GetContents(); !contents.Empty() {
		object := contents.Get(0).(EObject)
		if err := e.encodeContent(object); err != nil {
			e.addError(err)
			return
		}
	}

	// close db
	if err := e.db.Close(); err != nil {
		e.addError(err)
		return
	}

	// open db file
	dbFile, err := os.Open(dbPath)
	if err != nil {
		e.addError(err)
		return
	}
	defer func() {
		_ = dbFile.Close()
	}()

	// copy db file content to writer
	if _, err := io.Copy(e.writer, dbFile); err != nil {
		e.addError(err)
		return
	}

	// close db file
	if err := dbFile.Close(); err != nil {
		e.addError(err)
		return
	}

	// remove it from fs
	if err := os.Remove(dbPath); err != nil {
		e.addError(err)
		return
	}
}

func (e *SQLEncoder) addError(err error) {
	e.resource.GetErrors().Add(NewEDiagnosticImpl(err.Error(), e.resource.GetURI().String(), 0, 0))
}

func (e *SQLEncoder) EncodeObject(object EObject) error {
	return nil
}

func (e *SQLEncoder) encodeContent(eObject EObject) error {
	objectData, err := e.encodeObject(eObject)
	if err != nil {
		return err
	}

	insertContentStmt, err := e.getInsertStmt(e.schema.contentsTable)
	if err != nil {
		return err
	}

	if _, err := insertContentStmt.Exec(objectData.id); err != nil {
		return err
	}

	return nil
}

func (e *SQLEncoder) encodeObject(eObject EObject) (*sqlEncoderObjectData, error) {
	// encode object class
	eClass := eObject.EClass()
	classData, err := e.getClassData(eClass)
	if err != nil {
		return nil, err
	}

	// create table
	insertObjectStmt, err := e.getInsertStmt(e.schema.objectsTable)
	if err != nil {
		return nil, err
	}

	// insert object in table
	args := []any{classData.id}
	if idManager := e.resource.GetObjectIDManager(); idManager != nil && len(e.idAttributeName) > 0 {
		args = append(args, idManager.GetID(eObject))
	}

	sqlResult, err := insertObjectStmt.Exec(args...)
	if err != nil {
		return nil, err
	}

	// retrieve object id
	objectID, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	// collection of statements
	// used to avoid nested transactions
	insertStmts := newSqlStmts(e.db)
	for _, eClass := range classData.hierarchy {
		classData, err := e.getClassData(eClass)
		if err != nil {
			return nil, err
		}
		classTable := classData.schema.table

		// encode features columnValues in table columns
		columnValues := classTable.defaultValues()
		columnValues[classTable.key.index] = objectID
		for eFeature, featureData := range classData.features {
			if featureColumn := featureData.schema.column; featureColumn != nil {
				// feature is encoded as a column
				featureValue := eObject.(EObjectInternal).EGetResolve(eFeature, false)
				columnValue, err := e.encodeFeatureValue(featureData, featureValue)
				if err != nil {
					return nil, err
				}
				columnValues[featureColumn.index] = columnValue
			} else if featureTable := featureData.schema.table; featureTable != nil {
				// feature is encoded in a external table
				featureValue := eObject.(EObjectInternal).EGetResolve(eFeature, false)
				featureList, _ := featureValue.(EList)
				if featureList == nil {
					return nil, errors.New("feature value is not a list")
				}
				// retrieve insert statement
				insertStmt, err := e.getInsertStmt(featureTable)
				if err != nil {
					return nil, err
				}
				// for each list element, insert its value
				index := 0.0
				for itList := featureList.Iterator(); itList.HasNext(); {
					value := itList.Next()
					converted, err := e.encodeFeatureValue(featureData, value)
					if err != nil {
						return nil, err
					}
					insertStmts.add(insertStmt, objectID, index, converted)
					index++
				}
			}
		}

		// insert new row in class column
		insertStmt, err := e.getInsertStmt(classTable)
		if err != nil {
			return nil, err
		}
		insertStmts.add(insertStmt, columnValues...)
	}

	// execute all statements
	if err := insertStmts.exec(); err != nil {
		return nil, err
	}

	return &sqlEncoderObjectData{
		id:        objectID,
		classData: classData,
	}, nil
}

func (e *SQLEncoder) encodeFeatureValue(featureData *sqlEncoderFeatureData, value any) (any, error) {
	if value != nil {
		switch featureData.schema.featureKind {
		case sfkObject, sfkObjectList:
			objectData, err := e.encodeObject(value.(EObject))
			if err != nil {
				return nil, err
			}
			return objectData.id, nil
		case sfkObjectReference, sfkObjectReferenceList:
			ref := GetURI(value.(EObject))
			uri := e.resource.GetURI().Relativize(ref)
			return uri.String(), nil
		case sfkBool, sfkByte, sfkInt, sfkInt16, sfkInt32, sfkInt64, sfkEnum, sfkString, sfkByteArray:
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

func (e *SQLEncoder) getClassData(eClass EClass) (*sqlEncoderClassData, error) {
	classData := e.classDataMap[eClass]
	if classData == nil {
		// compute class data for super types
		for itClass := eClass.GetESuperTypes().Iterator(); itClass.HasNext(); {
			eClass := itClass.Next().(EClass)
			_, err := e.getClassData(eClass)
			if err != nil {
				return nil, err
			}
		}

		// encode package
		ePackage := eClass.GetEPackage()
		packageData, err := e.getPackageData(ePackage)
		if err != nil {
			return nil, err
		}

		// insert class in sql
		insertClassStmt, err := e.getInsertStmt(e.schema.classesTable)
		if err != nil {
			return nil, err
		}
		sqlResult, err := insertClassStmt.Exec(packageData.id, eClass.GetName())
		if err != nil {
			return nil, err
		}

		// retrieve class index
		classID, err := sqlResult.LastInsertId()
		if err != nil {
			return nil, err
		}

		// create class data
		classSchema, err := e.schema.getClassSchema(eClass)
		if err != nil {
			return nil, err
		}

		// compute class hierarchy
		classHierarchy := []EClass{eClass}
		for itClass := eClass.GetEAllSuperTypes().Iterator(); itClass.HasNext(); {
			classHierarchy = append(classHierarchy, itClass.Next().(EClass))
		}

		// create class tables
		if _, err := e.db.Exec(classSchema.table.createQuery()); err != nil {
			return nil, err
		}

		classFeatures := map[EStructuralFeature]*sqlEncoderFeatureData{}
		for _, featureSchema := range classSchema.features {

			// create feature table if any
			if table := featureSchema.table; table != nil {
				if _, err := e.db.Exec(table.createQuery()); err != nil {
					return nil, err
				}
			}

			// create feature data
			featureData := &sqlEncoderFeatureData{
				schema: featureSchema,
			}
			eFeature := featureSchema.feature
			if eAttribute, _ := eFeature.(EAttribute); eAttribute != nil {
				eDataType := eAttribute.GetEAttributeType()
				featureData.dataType = eDataType
				featureData.factory = eDataType.GetEPackage().GetEFactoryInstance()
			}

			classFeatures[eFeature] = featureData
		}

		// create & register class data
		classData = &sqlEncoderClassData{
			id:        classID,
			schema:    classSchema,
			features:  classFeatures,
			hierarchy: classHierarchy,
		}
		e.classDataMap[eClass] = classData

	}
	return classData, nil
}

func (e *SQLEncoder) getPackageData(ePackage EPackage) (*sqlEncoderPackageData, error) {
	ePackageData := e.packageDataMap[ePackage]
	if ePackageData == nil {
		// insert new package
		insertPackageStmt, err := e.getInsertStmt(e.schema.packagesTable)
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
		ePackageData = newSqlEncoderPackageData(id)
		e.packageDataMap[ePackage] = ePackageData
	}
	return ePackageData, nil
}

func (e *SQLEncoder) getInsertStmt(table *sqlTable) (stmt *sql.Stmt, err error) {
	stmt = e.insertStmts[table]
	if stmt == nil {
		stmt, err = e.db.Prepare(table.insertQuery())
	}
	return
}
