package ecore

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

type sqlEncoderFeatureData struct {
	schema   *sqlFeatureSchema
	dataType EDataType
	factory  EFactory
}

type sqlEncoderClassData struct {
	id        int64
	schema    *sqlClassSchema
	hierarchy []EClass
	features  *linkedHashMap[EStructuralFeature, *sqlEncoderFeatureData]
}

type sqlEncoderIDManager interface {
	sqlCodecIDManager
	getPackageID(EPackage) (int64, bool)
	getObjectID(EObject) (int64, bool)
	getClassID(EClass) (int64, bool)
	getEnumLiteralID(EEnumLiteral) (int64, bool)
}

type sqlEncoderIDManagerImpl struct {
	packages     map[EPackage]int64
	classes      map[EClass]int64
	objects      map[EObject]int64
	enumLiterals map[EEnumLiteral]int64
}

func newSqlEncoderIDManager() *sqlEncoderIDManagerImpl {
	return &sqlEncoderIDManagerImpl{
		packages:     map[EPackage]int64{},
		classes:      map[EClass]int64{},
		objects:      map[EObject]int64{},
		enumLiterals: map[EEnumLiteral]int64{},
	}
}

func (r *sqlEncoderIDManagerImpl) setPackageID(p EPackage, id int64) {
	r.packages[p] = id
}

func (r *sqlEncoderIDManagerImpl) getPackageID(p EPackage) (id int64, b bool) {
	id, b = r.packages[p]
	return
}

func (r *sqlEncoderIDManagerImpl) setObjectID(o EObject, id int64) {
	// set sql id if created object is an sql object
	if sqlObject, _ := o.(SQLObject); sqlObject != nil {
		sqlObject.SetSqlID(id)
	} else {
		// otherwise store it in map
		r.objects[o] = id
	}
}

func (r *sqlEncoderIDManagerImpl) getObjectID(o EObject) (id int64, b bool) {
	if sqlObject, _ := o.(SQLObject); sqlObject != nil {
		id = sqlObject.GetSqlID()
		b = id != 0
	} else {
		id, b = r.objects[o]
	}
	return
}

func (r *sqlEncoderIDManagerImpl) setClassID(c EClass, id int64) {
	r.classes[c] = id
}

func (r *sqlEncoderIDManagerImpl) getClassID(c EClass) (id int64, b bool) {
	id, b = r.classes[c]
	return
}

func (r *sqlEncoderIDManagerImpl) setEnumLiteralID(e EEnumLiteral, id int64) {
	r.enumLiterals[e] = id
}

func (r *sqlEncoderIDManagerImpl) getEnumLiteralID(e EEnumLiteral) (id int64, b bool) {
	id, b = r.enumLiterals[e]
	return
}

type sqlEncoderObjectManager struct {
}

func newSqlEncoderObjectManager() *sqlEncoderObjectManager {
	return &sqlEncoderObjectManager{}
}

func (r *sqlEncoderObjectManager) registerObject(EObject) {
}

type sqlEncoder struct {
	*sqlBase
	keepDefaults     bool
	classDataMap     map[EClass]*sqlEncoderClassData
	sqlIDManager     sqlEncoderIDManager
	sqlObjectManager sqlObjectManager
}

func (e *sqlEncoder) encodeVersion(conn *sqlite.Conn) error {
	return sqlitex.ExecuteTransient(conn, fmt.Sprintf(`PRAGMA user_version = %v;`, e.codecVersion), nil)
}

func (e *sqlEncoder) encodeProperties(conn *sqlite.Conn) error {
	// synchronous mode
	if err := sqlitex.ExecuteTransient(conn, `PRAGMA synchronous = NORMAL;`, nil); err != nil {
		return err
	}

	// journal mode
	if err := sqlitex.ExecuteTransient(conn, `PRAGMA journal_mode = WAL;`, nil); err != nil {
		return err
	}

	return nil
}

func (e *sqlEncoder) encodeSchema(conn *sqlite.Conn) error {
	// tables
	for _, table := range []*sqlTable{
		e.schema.packagesTable,
		e.schema.classesTable,
		e.schema.objectsTable,
		e.schema.contentsTable,
		e.schema.enumsTable,
	} {
		if err := sqlitex.ExecuteScript(conn, table.createQuery(), nil); err != nil {
			return err
		}
	}
	return nil
}

func (e *sqlEncoder) encodeContent(conn *sqlite.Conn, eObject EObject) error {
	objectID, err := e.encodeObject(conn, eObject)
	if err != nil {
		return err
	}
	return sqlitex.Execute(conn, e.schema.contentsTable.insertQuery(), &sqlitex.ExecOptions{Args: []any{objectID}})
}

func (e *sqlEncoder) encodeObject(conn *sqlite.Conn, eObject EObject) (id int64, err error) {
	defer sqlitex.Save(conn)(&err)
	objectID, isObjectID := e.sqlIDManager.getObjectID(eObject)
	if !isObjectID {

		// encode object class
		eClass := eObject.EClass()
		classData, err := e.encodeClass(conn, eClass)
		if err != nil {
			return -1, fmt.Errorf("getData('%s') error : %w", eClass.GetName(), err)
		}

		// insert object in table
		args := []any{classData.id}
		if idManager := e.idManager; idManager != nil && len(e.idAttributeName) > 0 {
			args = append(args, fmt.Sprintf("%v", idManager.GetID(eObject)))
		}

		if err := sqlitex.Execute(
			conn,
			e.schema.objectsTable.insertQuery(),
			&sqlitex.ExecOptions{Args: args}); err != nil {
			return -1, err
		}

		// retrieve object id
		objectID = conn.LastInsertRowID()

		// register object in registry
		e.sqlIDManager.setObjectID(eObject, objectID)

		// collection of statements
		// used to avoid nested transactions
		for _, eClass := range classData.hierarchy {
			classData, err := e.getEncoderClassData(conn, eClass)
			if err != nil {
				return -1, err
			}
			classTable := classData.schema.table

			// encode features columnValues in table columns
			columnValues := classTable.defaultValues()
			columnValues[classTable.key.index] = objectID
			for itFeature := classData.features.Iterator(); itFeature.Next(); {
				eFeature := itFeature.Key()
				if eObject.EIsSet(eFeature) || (e.keepDefaults && len(eFeature.GetDefaultValueLiteral()) > 0) {
					featureData := itFeature.Value()
					if featureColumn := featureData.schema.column; featureColumn != nil {
						featureValue := eObject.EGetResolve(eFeature, false)
						columnValue, err := e.encodeFeatureValue(conn, featureData, featureValue)
						if err != nil {
							return -1, err
						}
						columnValues[featureColumn.index] = columnValue
					} else if featureTable := featureData.schema.table; featureTable != nil {
						// feature is encoded in a external table
						featureValue := eObject.EGetResolve(eFeature, false)
						featureList, _ := featureValue.(EList)
						if featureList == nil {
							return -1, errors.New("feature value is not a list")
						}
						// for each list element, insert its value
						index := 1.0
						for itList := featureList.Iterator(); itList.HasNext(); {
							value := itList.Next()
							converted, err := e.encodeFeatureValue(conn, featureData, value)
							if err != nil {
								return -1, err
							}
							if err := sqlitex.Execute(
								conn,
								featureTable.insertQuery(),
								&sqlitex.ExecOptions{Args: []any{objectID, index, converted}},
							); err != nil {
								return -1, err
							}
							index++
						}
					}
				}
			}

			// insert new row in class column
			if err := sqlitex.Execute(
				conn,
				classTable.insertQuery(),
				&sqlitex.ExecOptions{
					Args: columnValues,
				}); err != nil {
				return -1, err
			}
		}

		// register in sql object manager
		// (must be done at the end otherwise internal data of eObject may disappear if its a EStoreEObject)
		e.sqlObjectManager.registerObject(eObject)

	}
	return objectID, nil
}

func (e *sqlEncoder) encodeFeatureValue(conn *sqlite.Conn, featureData *sqlEncoderFeatureData, value any) (encoded any, err error) {
	if value != nil {
		switch featureData.schema.featureKind {
		case sfkObject, sfkObjectList:
			if sqlObject, isSqlObject := value.(SQLObject); isSqlObject {
				objectID := sqlObject.GetSqlID()
				if objectID == 0 {
					objectID, err = e.encodeObject(conn, sqlObject)
					if err != nil {
						return nil, err
					}
					sqlObject.SetSqlID(objectID)
				}
				return objectID, nil
			} else if eObject, isEObject := value.(EObject); isEObject {
				return e.encodeObject(conn, eObject)
			}
		case sfkObjectReference, sfkObjectReferenceList:
			sqlID := int64(0)
			eObject := value.(EObject)
			if sqlObject, isSqlObject := value.(SQLObject); isSqlObject {
				sqlID = sqlObject.GetSqlID()
			} else {
				sqlID, _ = e.sqlIDManager.getObjectID(eObject)
			}
			if sqlID != 0 {
				return strconv.FormatInt(sqlID, 10), nil
			} else {
				ref := GetURI(eObject)
				uri := e.uri.Relativize(ref)
				return uri.String(), nil
			}
		case sfkEnum:
			eEnum := featureData.dataType.(EEnum)
			literal := featureData.factory.ConvertToString(eEnum, value)
			eEnumLiteral := eEnum.GetEEnumLiteralByLiteral(literal)
			if eEnumLiteral == nil {
				return nil, fmt.Errorf("unable to find enum literal in enmu '%s' with value '%v'", eEnum.GetName(), literal)
			}
			return e.encodeEnumLiteral(conn, eEnumLiteral)
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

func (e *sqlEncoder) encodeEnumLiteral(conn *sqlite.Conn, eEnumLiteral EEnumLiteral) (any, error) {
	if enumLiteralID, isEnumLiteralID := e.sqlIDManager.getEnumLiteralID(eEnumLiteral); isEnumLiteralID {
		return enumLiteralID, nil
	} else {
		eEnum := eEnumLiteral.GetEEnum()
		ePackage := eEnum.GetEPackage()
		packageID, err := e.encodePackage(conn, ePackage)
		if err != nil {
			return nil, err
		}

		// insert enum
		if err := sqlitex.Execute(conn, e.schema.enumsTable.insertQuery(), &sqlitex.ExecOptions{
			Args: []any{packageID, eEnum.GetName(), eEnumLiteral.GetLiteral()},
		}); err != nil {
			return nil, err
		}

		// retrieve enum index
		enumID := conn.LastInsertRowID()

		// register enum literal
		e.sqlIDManager.setEnumLiteralID(eEnumLiteral, enumID)
		return enumID, nil
	}
}

func (e *sqlEncoder) encodeClass(conn *sqlite.Conn, eClass EClass) (*sqlEncoderClassData, error) {
	classData, err := e.getEncoderClassData(conn, eClass)
	if err != nil {
		return nil, err
	}
	if classData.id == -1 {
		// class is not encoded
		// check if class is registered in registry
		classID, isClassID := e.sqlIDManager.getClassID(eClass)
		if isClassID {
			// already registered
			classData.id = classID
		} else {
			// not registered
			// got to insert in classes table and retirve its id

			// encode package
			ePackage := eClass.GetEPackage()
			packageID, err := e.encodePackage(conn, ePackage)
			if err != nil {
				return nil, err
			}

			// insert new class
			if err := sqlitex.Execute(conn, e.schema.classesTable.insertQuery(), &sqlitex.ExecOptions{
				Args: []any{packageID, eClass.GetName()},
			}); err != nil {
				return nil, err
			}

			// retrieve class index
			classID = conn.LastInsertRowID()

			classData.id = classID

			// register eClass with its id in registry
			e.sqlIDManager.setClassID(eClass, classID)
		}
	}
	return classData, nil
}

func (e *sqlEncoder) getEncoderClassData(conn *sqlite.Conn, eClass EClass) (*sqlEncoderClassData, error) {
	classData := e.classDataMap[eClass]
	if classData == nil {
		// compute class data for super types
		for itClass := eClass.GetESuperTypes().Iterator(); itClass.HasNext(); {
			eClass := itClass.Next().(EClass)
			_, err := e.getEncoderClassData(conn, eClass)
			if err != nil {
				return nil, err
			}
		}

		// create class schema
		classSchema := e.schema.getClassSchema(eClass)

		// compute class hierarchy
		classHierarchy := []EClass{eClass}
		for itClass := eClass.GetEAllSuperTypes().Iterator(); itClass.HasNext(); {
			classHierarchy = append(classHierarchy, itClass.Next().(EClass))
		}

		// create class tables
		if err := sqlitex.ExecuteScript(conn, classSchema.table.createQuery(), nil); err != nil {
			return nil, err
		}

		// computes features data
		classFeatures := newLinkedHashMap[EStructuralFeature, *sqlEncoderFeatureData]()
		for _, featureSchema := range classSchema.features {

			// create feature table if any
			if table := featureSchema.table; table != nil {
				if err := sqlitex.ExecuteScript(conn, table.createQuery(), nil); err != nil {
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
			classFeatures.put(eFeature, featureData)
		}

		// create & register class data
		classData = &sqlEncoderClassData{
			id:        -1,
			schema:    classSchema,
			features:  classFeatures,
			hierarchy: classHierarchy,
		}
		e.classDataMap[eClass] = classData

	}
	return classData, nil
}

func (e *sqlEncoder) encodePackage(conn *sqlite.Conn, ePackage EPackage) (int64, error) {
	packageID, isPackageID := e.sqlIDManager.getPackageID(ePackage)
	if !isPackageID {
		// insert new package
		if err := sqlitex.Execute(conn, e.schema.packagesTable.insertQuery(), &sqlitex.ExecOptions{
			Args: []any{ePackage.GetNsURI()},
		}); err != nil {
			return -1, err
		}
		// retrieve package index
		packageID = conn.LastInsertRowID()
		// set package id
		e.sqlIDManager.setPackageID(ePackage, packageID)
	}
	return packageID, nil
}

type SQLEncoder struct {
	sqlEncoder
	resource     EResource
	connProvider func() (*sqlite.Conn, error)
	connClose    func(conn *sqlite.Conn) error
}

func NewSQLWriterEncoder(w io.Writer, resource EResource, options map[string]any) *SQLEncoder {
	// create a temp file for the database file
	fileName := filepath.Base(resource.GetURI().Path())
	dbPath, err := sqlTmpDB(fileName)
	if err != nil {
		return nil
	}
	return newSQLEncoder(
		func() (*sqlite.Conn, error) {
			return sqlite.OpenConn(dbPath, sqlite.OpenReadWrite|sqlite.OpenCreate)
		},
		func(conn *sqlite.Conn) error {
			// close db
			if err := conn.Close(); err != nil {
				return err
			}

			// open db file
			dbFile, err := os.Open(dbPath)
			if err != nil {
				return err
			}
			defer func() {
				_ = dbFile.Close()
			}()

			// copy db file content to writer
			if _, err := io.Copy(w, dbFile); err != nil {
				return err
			}

			// close db file
			if err := dbFile.Close(); err != nil {
				return err
			}

			// remove it from fs
			if err := os.Remove(dbPath); err != nil {
				return err
			}

			return nil
		},
		resource,
		options,
	)
}

func NewSQLDBEncoder(conn *sqlite.Conn, resource EResource, options map[string]any) *SQLEncoder {
	return newSQLEncoder(
		func() (*sqlite.Conn, error) { return conn, nil },
		func(conn *sqlite.Conn) error { return nil },
		resource,
		options)
}

func newSQLEncoder(connProvider func() (*sqlite.Conn, error), connClose func(conn *sqlite.Conn) error, resource EResource, options map[string]any) *SQLEncoder {
	// options
	schemaOptions := []sqlSchemaOption{}
	idAttributeName := ""
	keepDefaults := false
	codecVersion := sqlCodecVersion
	if options != nil {
		if id, isID := options[SQL_OPTION_ID_ATTRIBUTE_NAME].(string); isID && len(id) > 0 && resource.GetObjectIDManager() != nil {
			schemaOptions = append(schemaOptions, withIDAttributeName(id))
			idAttributeName = id
		}

		if b, isBool := options[SQL_OPTION_KEEP_DEFAULTS].(bool); isBool {
			keepDefaults = b
		}

		if v, isVersion := options[SQL_OPTION_CODEC_VERSION].(int64); isVersion {
			codecVersion = v
		}
	}

	// encoder structure
	return &SQLEncoder{
		sqlEncoder: sqlEncoder{
			sqlBase: &sqlBase{
				codecVersion:    codecVersion,
				uri:             resource.GetURI(),
				idManager:       resource.GetObjectIDManager(),
				idAttributeName: idAttributeName,
				schema:          newSqlSchema(schemaOptions...),
			},
			keepDefaults:     keepDefaults,
			classDataMap:     map[EClass]*sqlEncoderClassData{},
			sqlIDManager:     newSqlEncoderIDManager(),
			sqlObjectManager: newSqlEncoderObjectManager(),
		},
		resource:     resource,
		connProvider: connProvider,
		connClose:    connClose,
	}
}

func (e *SQLEncoder) EncodeResource() {
	// create conn
	// open conn
	conn, err := e.connProvider()
	if err != nil {
		e.addError(err)
		return
	}
	defer func() {
		err = e.connClose(conn)
		if err != nil {
			e.addError(err)
		}
	}()

	if err := e.encodeVersion(conn); err != nil {
		e.addError(err)
		return
	}

	if err := e.encodeProperties(conn); err != nil {
		e.addError(err)
		return
	}

	if err := e.encodeSchema(conn); err != nil {
		e.addError(err)
		return
	}

	// encode contents into db
	if contents := e.resource.GetContents(); !contents.Empty() {
		object := contents.Get(0).(EObject)
		if err := e.encodeContent(conn, object); err != nil {
			e.addError(err)
			return
		}
	}
}

func (e *SQLEncoder) addError(err error) {
	e.resource.GetErrors().Add(NewEDiagnosticImpl(err.Error(), e.resource.GetURI().String(), 0, 0))
}

func (e *SQLEncoder) EncodeObject(object EObject) error {
	return nil
}
