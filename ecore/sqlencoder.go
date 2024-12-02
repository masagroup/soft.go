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

type SQLEncoderIDManager interface {
	SQLCodecIDManager
	// Returns package id , decoded as true if package is already encoded otherwise false
	GetPackageID(EPackage) (int64, bool)
	// Returns object id , decoded as true if package is already encoded otherwise false
	GetObjectID(EObject) (int64, bool)
	// Returns class id , decoded as true if package is already encoded otherwise false
	GetClassID(EClass) (int64, bool)
	// Returns enum literal id , decoded as true if package is already encoded otherwise false
	GetEnumLiteralID(EEnumLiteral) (int64, bool)
}

type sqlEncoderIDManagerImpl struct {
	packages     map[EPackage]int64
	classes      map[EClass]int64
	objects      map[EObject]int64
	enumLiterals map[EEnumLiteral]int64
}

func newSQLEncoderIDManager() SQLEncoderIDManager {
	return &sqlEncoderIDManagerImpl{
		packages:     map[EPackage]int64{},
		classes:      map[EClass]int64{},
		objects:      map[EObject]int64{},
		enumLiterals: map[EEnumLiteral]int64{},
	}
}

func (r *sqlEncoderIDManagerImpl) SetPackageID(p EPackage, id int64) {
	r.packages[p] = id
}

func (r *sqlEncoderIDManagerImpl) GetPackageID(p EPackage) (id int64, b bool) {
	id, b = r.packages[p]
	return
}

func (r *sqlEncoderIDManagerImpl) SetObjectID(o EObject, id int64) {
	// store it in map
	r.objects[o] = id
	// set sql id if created object is an sql object
	if sqlObject, _ := o.(SQLObject); sqlObject != nil {
		sqlObject.SetSQLID(id)
	}
}

func (r *sqlEncoderIDManagerImpl) GetObjectID(o EObject) (id int64, b bool) {
	if id, b = r.objects[o]; !b {
		if sqlObject, _ := o.(SQLObject); sqlObject != nil {
			id = sqlObject.GetSQLID()
		}
	}
	return
}

func (r *sqlEncoderIDManagerImpl) SetClassID(c EClass, id int64) {
	r.classes[c] = id
}

func (r *sqlEncoderIDManagerImpl) GetClassID(c EClass) (id int64, b bool) {
	id, b = r.classes[c]
	return
}

func (r *sqlEncoderIDManagerImpl) SetEnumLiteralID(e EEnumLiteral, id int64) {
	r.enumLiterals[e] = id
}

func (r *sqlEncoderIDManagerImpl) GetEnumLiteralID(e EEnumLiteral) (id int64, b bool) {
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
	isKeepDefaults     bool
	isContainerEncoded bool
	classDataMap       map[EClass]*sqlEncoderClassData
	sqlIDManager       SQLEncoderIDManager
	sqlObjectManager   sqlObjectManager
}

func (e *sqlEncoder) encodeVersion(conn *sqlite.Conn) error {
	return sqlitex.ExecuteTransient(conn, fmt.Sprintf(`PRAGMA user_version = %v;`, e.codecVersion), nil)
}

func (e *sqlEncoder) encodePragmas(conn *sqlite.Conn) error {
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
		e.schema.propertiesTable,
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

func (e *sqlEncoder) encodeProperties(conn *sqlite.Conn) error {
	if len(e.objectIDName) > 0 {
		if err := sqlitex.ExecuteTransient(conn, e.schema.propertiesTable.insertQuery(), &sqlitex.ExecOptions{
			Args: []any{"objectID", e.objectIDName},
		}); err != nil {
			return err
		}
	}
	if e.isContainerEncoded {
		if err := sqlitex.ExecuteTransient(conn, e.schema.propertiesTable.insertQuery(), &sqlitex.ExecOptions{
			Args: []any{"containerID", true},
		}); err != nil {
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
	sqlObjectID, isSqlObjectID := e.sqlIDManager.GetObjectID(eObject)
	if !isSqlObjectID {
		defer sqlitex.Save(conn)(&err)

		// encode object class
		eClass := eObject.EClass()
		classData, err := e.encodeClass(conn, eClass)
		if err != nil {
			return -1, fmt.Errorf("getData('%s') error : %w", eClass.GetName(), err)
		}

		objectTable := e.schema.objectsTable
		// args
		args := []any{}
		// object id
		if sqlObjectID != 0 {
			args = append(args, sqlObjectID)
		} else {
			args = append(args, nil)
		}
		// class id
		args = append(args, classData.id)
		// container and container feature id
		if e.isContainerEncoded {
			if container := eObject.EContainer(); container != nil {
				containerID, err := e.encodeObject(conn, eObject.EContainer())
				if err != nil {
					return -1, err
				}
				containerFeatureID := eObject.EContainingFeature().GetFeatureID()
				args = append(args, containerID, containerFeatureID)
			} else {
				args = append(args, nil, nil)
			}
		}
		// object id is serialized only if we have an id manager and
		// a corresponding column in objects table
		if e.objectIDManager != nil && len(objectTable.columns) > 2 {
			args = append(args, fmt.Sprintf("%v", e.objectIDManager.GetID(eObject)))
		}

		// query
		if err := sqlitex.Execute(
			conn,
			e.schema.objectsTable.insertQuery(),
			&sqlitex.ExecOptions{Args: args}); err != nil {
			return -1, err
		}

		if sqlObjectID == 0 {
			// if object id not defined, retrieve it
			// as the last insert row id
			sqlObjectID = conn.LastInsertRowID()
		}

		// register object in registry
		e.sqlIDManager.SetObjectID(eObject, sqlObjectID)

		// for all object hierarchy classes
		for _, eClass := range classData.hierarchy {
			classData, err := e.getEncoderClassData(conn, eClass)
			if err != nil {
				return -1, err
			}
			classTable := classData.schema.table

			// encode features columnValues in table columns
			columnValues := classTable.defaultValues()
			columnValues[classTable.key.index] = sqlObjectID
			for itFeature := classData.features.Iterator(); itFeature.HasNext(); {
				eFeature := itFeature.Key()
				if eObject.EIsSet(eFeature) || (e.isKeepDefaults && len(eFeature.GetDefaultValueLiteral()) > 0) {
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
								&sqlitex.ExecOptions{Args: []any{sqlObjectID, index, converted}},
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
	return sqlObjectID, nil
}

func (e *sqlEncoder) encodeFeatureValue(conn *sqlite.Conn, featureData *sqlEncoderFeatureData, value any) (encoded any, err error) {
	if value != nil {
		switch featureData.schema.featureKind {
		case sfkObject, sfkObjectList:
			eObject := value.(EObject)
			return e.encodeObject(conn, eObject)
		case sfkObjectReference, sfkObjectReferenceList:
			eObject := value.(EObject)
			sqlID, _ := e.sqlIDManager.GetObjectID(eObject)
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
	if enumLiteralID, isEnumLiteralID := e.sqlIDManager.GetEnumLiteralID(eEnumLiteral); isEnumLiteralID {
		return enumLiteralID, nil
	} else {
		eEnum := eEnumLiteral.GetEEnum()
		ePackage := eEnum.GetEPackage()
		packageID, err := e.encodePackage(conn, ePackage)
		if err != nil {
			return nil, err
		}

		// query args
		args := []any{}
		if enumLiteralID == 0 {
			args = append(args, nil)
		} else {
			args = append(args, enumLiteralID)
		}
		args = append(args, packageID, eEnum.GetName(), eEnumLiteral.GetLiteral())

		// insert enum
		if err := sqlitex.Execute(conn, e.schema.enumsTable.insertQuery(), &sqlitex.ExecOptions{
			Args: args,
		}); err != nil {
			return nil, err
		}

		if enumLiteralID == 0 {
			// retrieve enum index if not defined
			enumLiteralID = conn.LastInsertRowID()
		}

		// register enum literal
		e.sqlIDManager.SetEnumLiteralID(eEnumLiteral, enumLiteralID)

		return enumLiteralID, nil
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
		classID, isClassID := e.sqlIDManager.GetClassID(eClass)
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

			// args
			args := []any{}
			if classID == 0 {
				args = append(args, nil)
			} else {
				args = append(args, classID)
			}
			args = append(args, packageID, eClass.GetName())

			// insert new class
			if err := sqlitex.Execute(conn, e.schema.classesTable.insertQuery(), &sqlitex.ExecOptions{
				Args: args,
			}); err != nil {
				return nil, err
			}

			if classID == 0 {
				// retrieve class id
				classID = conn.LastInsertRowID()
			}

			// set class data id
			classData.id = classID

			// register eClass with its id in registry
			e.sqlIDManager.SetClassID(eClass, classID)
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
			classFeatures.Put(eFeature, featureData)
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
	packageID, isPackageID := e.sqlIDManager.GetPackageID(ePackage)
	if !isPackageID {
		// args
		args := []any{}
		if packageID == 0 {
			args = append(args, nil)
		} else {
			args = append(args, packageID)
		}
		args = append(args, ePackage.GetNsURI())

		// query
		if err := sqlitex.Execute(conn, e.schema.packagesTable.insertQuery(), &sqlitex.ExecOptions{
			Args: args,
		}); err != nil {
			return -1, err
		}

		if packageID == 0 {
			// retrieve package index
			packageID = conn.LastInsertRowID()
		}

		// set package id
		e.sqlIDManager.SetPackageID(ePackage, packageID)
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
	objectIDName := ""
	isContainerEncoded := false
	isKeepDefaults := false
	codecVersion := sqlCodecVersion
	sqlIDManager := newSQLEncoderIDManager()
	if options != nil {
		if id, isID := options[SQL_OPTION_OBJECT_ID].(string); isID && len(id) > 0 && resource.GetObjectIDManager() != nil {
			schemaOptions = append(schemaOptions, withObjectIDName(id))
			objectIDName = id
		}
		if b, isBool := options[SQL_OPTION_CONTAINER_ID].(bool); isBool {
			schemaOptions = append(schemaOptions, withContainerID(b))
			isContainerEncoded = b
		}
		if b, isBool := options[SQL_OPTION_KEEP_DEFAULTS].(bool); isBool {
			isKeepDefaults = b
		}
		if v, isVersion := options[SQL_OPTION_CODEC_VERSION].(int64); isVersion {
			codecVersion = v
		}
		if m, isSQLIDManager := options[SQL_OPTION_SQL_ID_MANAGER].(SQLEncoderIDManager); isSQLIDManager {
			sqlIDManager = m
		}
	}

	// encoder structure
	return &SQLEncoder{
		sqlEncoder: sqlEncoder{
			sqlBase: &sqlBase{
				codecVersion:    codecVersion,
				uri:             resource.GetURI(),
				objectIDManager: resource.GetObjectIDManager(),
				objectIDName:    objectIDName,
				schema:          newSqlSchema(schemaOptions...),
			},
			isKeepDefaults:     isKeepDefaults,
			isContainerEncoded: isContainerEncoded,
			classDataMap:       map[EClass]*sqlEncoderClassData{},
			sqlIDManager:       sqlIDManager,
			sqlObjectManager:   newSqlEncoderObjectManager(),
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

	if err := e.encodePragmas(conn); err != nil {
		e.addError(err)
		return
	}

	if err := e.encodeSchema(conn); err != nil {
		e.addError(err)
		return
	}

	if err := e.encodeProperties(conn); err != nil {
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
