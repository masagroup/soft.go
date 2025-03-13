package ecore

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/chebyrash/promise"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

type sqlEncoderFeatureData struct {
	schema      *sqlFeatureSchema
	featureID   int64
	dataType    EDataType
	factory     EFactory
	isTransient bool
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

type sqlEncoderLock struct {
	mutex    sync.Mutex
	refCount atomic.Int64
}

type sqlEncoderLockManager struct {
	mutex sync.Mutex
	locks map[any]*sqlEncoderLock
}

func newSqlEncoderLockManager() *sqlEncoderLockManager {
	return &sqlEncoderLockManager{
		locks: map[any]*sqlEncoderLock{},
	}
}

func (l *sqlEncoderLockManager) lock(object any) {
	l.mutex.Lock()
	lock := l.locks[object]
	if lock == nil {
		lock = &sqlEncoderLock{}
		l.locks[object] = lock
	}
	l.mutex.Unlock()
	lock.refCount.Add(1)
	lock.mutex.Lock()

}

func (l *sqlEncoderLockManager) unlock(object any) {
	l.mutex.Lock()
	lock := l.locks[object]
	if lock != nil {
		if lock.refCount.Add(-1) == 0 {
			delete(l.locks, object)
		}
		lock.mutex.Unlock()
	}
	l.mutex.Unlock()
}

type sqlEncoder struct {
	*sqlBase
	isForced         bool
	isKeepDefaults   bool
	classDataMap     map[EClass]*sqlEncoderClassData
	sqlIDManager     SQLEncoderIDManager
	sqlObjectManager sqlObjectManager
	sqlLockManager   *sqlEncoderLockManager
}

func (e *sqlEncoder) encodeVersion() error {
	if !e.isForced {
		var version int64
		if err := e.executeQueryTransient("PRAGMA user_version;", &sqlitex.ExecOptions{
			ResultFunc: func(stmt *sqlite.Stmt) error {
				version = stmt.ColumnInt64(0)
				return nil
			},
		}); err != nil {
			return err
		}
		// already encoded
		if version == e.codecVersion {
			return nil
		}
	}
	// encode
	return e.executeQueryTransient(fmt.Sprintf(`PRAGMA user_version = %v;`, e.codecVersion), nil)
}

func (e *sqlEncoder) encodeSchema() error {
	// tables
	for _, table := range []*sqlTable{
		e.schema.propertiesTable,
		e.schema.packagesTable,
		e.schema.classesTable,
		e.schema.objectsTable,
		e.schema.contentsTable,
		e.schema.enumsTable,
	} {
		if err := e.executeQueryScript(table.createQuery(), nil); err != nil {
			return err
		}
	}
	return nil
}

func (e *sqlEncoder) encodeProperties() (err error) {
	previous := map[string]string{}
	if !e.isForced {
		previous, err = e.decodeProperties()
		if err != nil {
			return err
		}
	}
	properties := map[string]string{}
	if len(e.objectIDName) > 0 {
		properties["objectID"] = e.objectIDName
	}
	if e.isContainerID {
		properties["containerID"] = "true"
	}

	query := e.schema.propertiesTable.insertOrReplaceQuery()
	for k, v := range properties {
		if previous[k] != v {
			if err := e.executeQuery(query, &sqlitex.ExecOptions{
				Args: []any{k, v},
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *sqlEncoder) encodeContent(eObject EObject) error {
	objectID, err := e.encodeObject(eObject, -1, -1)
	if err != nil {
		return err
	}
	return e.executeQuery(e.schema.contentsTable.insertQuery(), &sqlitex.ExecOptions{Args: []any{objectID}})
}

type ptrField struct{ ptr any }

func (f ptrField) String() string {
	return fmt.Sprintf("%p", f.ptr)
}

func (e *sqlEncoder) shouldEncodeFeature(eObject EObject, eFeature EStructuralFeature) bool {
	return eObject.EIsSet(eFeature) || (e.isKeepDefaults && len(eFeature.GetDefaultValueLiteral()) > 0)
}

func (e *sqlEncoder) encodeObject(eObject EObject, sqlContainerID int64, containerFeatureID int64) (id int64, err error) {
	logger := e.logger.Named("encoder").With(zap.Stringer("object", ptrField{eObject}), zap.String("class", eObject.EClass().GetName()))
	logger.Debug("object encode")
	sqlObjectID, isSqlObjectID := e.sqlIDManager.GetObjectID(eObject)
	if !isSqlObjectID {
		e.sqlLockManager.lock(eObject)
		defer e.sqlLockManager.unlock(eObject)

		// object may be encoded in another goroutine
		sqlObjectID, isSqlObjectID = e.sqlIDManager.GetObjectID(eObject)
		if isSqlObjectID {
			logger.Debug("object encoded in another goroutine", zap.Int64("id", sqlObjectID))
			return sqlObjectID, nil
		}

		logger.Debug("object encoding")

		// encode object class
		eClass := eObject.EClass()
		classData, err := e.encodeClass(eClass)
		if err != nil {
			return -1, fmt.Errorf("getData('%s') error : %w", eClass.GetName(), err)
		}

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
		if e.isContainerID {
			if sqlContainerID > 0 {
				args = append(args, sqlContainerID, containerFeatureID)
			} else {
				args = append(args, nil, nil)
			}
		}
		// object id is serialized only if we have an id manager and
		// a corresponding column in objects table
		if e.isObjectID {
			args = append(args, fmt.Sprintf("%v", e.objectIDManager.GetID(eObject)))
		}

		// query
		if err := e.executeQuery(
			e.schema.objectsTable.insertQuery(),
			&sqlitex.ExecOptions{
				Args: args,
				ResultFunc: func(stmt *sqlite.Stmt) error {
					if sqlObjectID == 0 {
						sqlObjectID = stmt.ColumnInt64(0)
					}
					return nil
				},
			}); err != nil {
			return -1, err
		}

		// register object in registry
		e.sqlIDManager.SetObjectID(eObject, sqlObjectID)

		if lockProvider, isLockProvider := eObject.(ELockProvider); isLockProvider {
			lockProvider.Lock()
			defer lockProvider.Unlock()
		}

		// for all object hierarchy classes
		for _, eClass := range classData.hierarchy {
			classData, err := e.getEncoderClassData(eClass)
			if err != nil {
				return -1, err
			}
			classTable := classData.schema.table

			// encode features columnValues in table columns
			columnValues := classTable.defaultValues()
			columnValues[classTable.key.index] = sqlObjectID
			for itFeature := classData.features.Iterator(); itFeature.HasNext(); {
				eFeature := itFeature.Key()
				featureData := itFeature.Value()
				if !featureData.isTransient && e.shouldEncodeFeature(eObject, eFeature) {
					if featureColumn := featureData.schema.column; featureColumn != nil {
						featureValue := eObject.EGetResolve(eFeature, false)
						columnValue, err := e.encodeFeatureValue(sqlObjectID, featureData, featureValue)
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
						for value := range featureList.All() {
							converted, err := e.encodeFeatureValue(sqlObjectID, featureData, value)
							if err != nil {
								return -1, err
							}
							if err := e.executeQuery(
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
			if err := e.executeQuery(
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
	logger.Debug("object encoded", zap.Int64("id", sqlObjectID))
	return sqlObjectID, nil
}

func (e *sqlEncoder) encodeFeatureValue(containerID int64, featureData *sqlEncoderFeatureData, value any) (encoded any, err error) {
	if value != nil {
		switch featureData.schema.featureKind {
		case sfkObject, sfkObjectList:
			eObject := value.(EObject)
			return e.encodeObject(eObject, containerID, featureData.featureID)
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
			return e.encodeEnumLiteral(eEnumLiteral)
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

func (e *sqlEncoder) encodeEnumLiteral(eEnumLiteral EEnumLiteral) (any, error) {
	// lock class data
	e.sqlLockManager.lock(eEnumLiteral)
	defer e.sqlLockManager.unlock(eEnumLiteral)

	if enumLiteralID, isEnumLiteralID := e.sqlIDManager.GetEnumLiteralID(eEnumLiteral); isEnumLiteralID {
		return enumLiteralID, nil
	} else {
		eEnum := eEnumLiteral.GetEEnum()
		ePackage := eEnum.GetEPackage()
		packageID, err := e.encodePackage(ePackage)
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
		if err := e.executeQuery(e.schema.enumsTable.insertQuery(),
			&sqlitex.ExecOptions{
				Args: args,
				ResultFunc: func(stmt *sqlite.Stmt) error {
					if enumLiteralID == 0 {
						enumLiteralID = stmt.ColumnInt64(0)
					}
					return nil
				},
			}); err != nil {
			return nil, err
		}

		// register enum literal
		e.sqlIDManager.SetEnumLiteralID(eEnumLiteral, enumLiteralID)

		return enumLiteralID, nil
	}
}

func (e *sqlEncoder) encodeClass(eClass EClass) (*sqlEncoderClassData, error) {
	logger := e.logger.Named("encoder").With(zap.String("class", eClass.GetName()))
	logger.Debug("class encode")

	// retrieve class data
	classData, err := e.getEncoderClassData(eClass)
	if err != nil {
		return nil, err
	}

	// lock class data
	e.sqlLockManager.lock(classData)
	defer e.sqlLockManager.unlock(classData)

	// encode class
	if classData.id == -1 {
		logger.Debug("class register")

		// class is not encoded
		// check if class is registered in registry
		classID, isClassID := e.sqlIDManager.GetClassID(eClass)
		if isClassID {
			logger.Debug("class already registered", zap.Int64("id", classID))
			// already registered
			classData.id = classID
		} else {
			logger.Debug("class not registered", zap.Int64("id", classID))
			// not registered
			// got to insert in classes table and retirve its id

			// encode package
			ePackage := eClass.GetEPackage()
			packageID, err := e.encodePackage(ePackage)
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
			if err := e.executeQuery(e.schema.classesTable.insertQuery(), &sqlitex.ExecOptions{
				Args: args,
				ResultFunc: func(stmt *sqlite.Stmt) error {
					if classID == 0 {
						classID = stmt.ColumnInt64(0)
					}
					return nil
				},
			}); err != nil {
				return nil, err
			}

			// set class data id
			classData.id = classID

			// register eClass with its id in registry
			e.sqlIDManager.SetClassID(eClass, classID)
		}
	}
	logger.Debug("class encoded", zap.Int64("id", classData.id))
	return classData, nil
}

func (e *sqlEncoder) getEncoderClassData(eClass EClass) (*sqlEncoderClassData, error) {
	logger := e.logger.Named("encoder").With(zap.String("class", eClass.GetName()))

	// lock class
	e.sqlLockManager.lock(eClass)
	defer e.sqlLockManager.unlock(eClass)

	classData := e.classDataMap[eClass]
	if classData == nil {
		logger.Debug("class data create")

		// compute class data for super types
		for itClass := eClass.GetESuperTypes().Iterator(); itClass.HasNext(); {
			eClass := itClass.Next().(EClass)
			_, err := e.getEncoderClassData(eClass)
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
		if err := e.executeQueryScript(classSchema.table.createQuery(), nil); err != nil {
			return nil, err
		}

		logger.Debug("class data created", zap.String("table", classSchema.table.name))

		// computes features data
		classFeatures := newLinkedHashMap[EStructuralFeature, *sqlEncoderFeatureData]()
		for _, featureSchema := range classSchema.features {
			eFeature := featureSchema.feature

			// create feature table if any
			if table := featureSchema.table; table != nil {
				if err := e.executeQueryScript(table.createQuery(), nil); err != nil {
					return nil, err
				}
			}

			// create feature data
			featureData := &sqlEncoderFeatureData{
				schema:    featureSchema,
				featureID: int64(eClass.GetFeatureID(eFeature)),
			}
			if eReference, _ := eFeature.(EReference); eReference != nil {
				featureData.isTransient = eReference.IsTransient() || (eReference.IsContainer() && !eReference.IsResolveProxies())
			} else if eAttribute, _ := eFeature.(EAttribute); eAttribute != nil {
				eDataType := eAttribute.GetEAttributeType()
				featureData.isTransient = eAttribute.IsTransient()
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

func (e *sqlEncoder) encodePackage(ePackage EPackage) (int64, error) {
	e.sqlLockManager.lock(ePackage)
	defer e.sqlLockManager.unlock(ePackage)

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
		if err := e.executeQuery(e.schema.packagesTable.insertQuery(), &sqlitex.ExecOptions{
			Args: args,
			ResultFunc: func(stmt *sqlite.Stmt) error {
				if packageID == 0 {
					packageID = stmt.ColumnInt64(0)
				}
				return nil
			},
		}); err != nil {
			return -1, err
		}

		// set package id
		e.sqlIDManager.SetPackageID(ePackage, packageID)
	}
	return packageID, nil
}

type SQLEncoder struct {
	sqlEncoder
	resource EResource
}

func NewSQLWriterEncoder(w io.Writer, resource EResource, options map[string]any) *SQLEncoder {
	inMemoryDatabase := false
	if options != nil {
		inMemoryDatabase, _ = options[SQL_OPTION_IN_MEMORY_DATABASE].(bool)
	}
	if inMemoryDatabase {
		return newSQLEncoder(
			func() (*sqlitex.Pool, error) {
				return sqlitex.NewPool("file::memory:?mode=memory&cache=shared", sqlitex.PoolOptions{Flags: sqlite.OpenReadWrite | sqlite.OpenCreate | sqlite.OpenURI})
			},
			func(connPool *sqlitex.Pool) (err error) {
				// close pool
				defer func() {
					if err2 := connPool.Close(); err2 != nil && err == nil {
						err = err2
					}
				}()

				// open connection
				conn, err := connPool.Take(context.Background())
				if err != nil {
					return
				}
				defer connPool.Put(conn)

				// save bytes
				bytes, err := conn.Serialize("")
				if err != nil {
					return err
				}

				// set journal mode as WAL
				bytes[18] = 0x02
				bytes[19] = 0x02

				// write bytes to writer
				if _, err := w.Write(bytes); err != nil {
					return err
				}

				return nil
			},
			resource,
			options,
		)
	} else {
		// create a temp file for the database file
		fileName := filepath.Base(resource.GetURI().Path())
		dbPath, err := sqlTmpDB(fileName)
		if err != nil {
			return nil
		}
		return newSQLEncoder(
			func() (*sqlitex.Pool, error) {
				return sqlitex.NewPool(
					dbPath,
					sqlitex.PoolOptions{
						Flags: sqlite.OpenReadWrite | sqlite.OpenCreate | sqlite.OpenWAL,
						PrepareConn: func(conn *sqlite.Conn) error {
							// connection is in synchronous mode
							return sqlitex.ExecuteTransient(conn, "PRAGMA synchronous=normal", nil)
						}},
				)
			},
			func(connPool *sqlitex.Pool) error {
				// close db
				if err := connPool.Close(); err != nil {
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
}

func NewSQLDBEncoder(connPool *sqlitex.Pool, resource EResource, options map[string]any) *SQLEncoder {
	return newSQLEncoder(
		func() (*sqlitex.Pool, error) { return connPool, nil },
		func(conn *sqlitex.Pool) error { return nil },
		resource,
		options)
}

func newSQLEncoder(connectionPoolProvider func() (*sqlitex.Pool, error), connectionPoolClose func(conn *sqlitex.Pool) error, resource EResource, options map[string]any) *SQLEncoder {
	// options
	schemaOptions := []sqlSchemaOption{}
	objectIDName := ""
	isContainerID := false
	isObjectID := false
	isKeepDefaults := false
	codecVersion := sqlCodecVersion
	sqlIDManager := newSQLEncoderIDManager()
	objectIDManager := resource.GetObjectIDManager()
	logger := zap.NewNop()
	if options != nil {
		if id, isID := options[SQL_OPTION_OBJECT_ID].(string); isID && len(id) > 0 && objectIDManager != nil {
			schemaOptions = append(schemaOptions, withObjectIDName(id))
			objectIDName = id
			isObjectID = objectIDName != "objectID"
		}
		if b, isBool := options[SQL_OPTION_CONTAINER_ID].(bool); isBool {
			schemaOptions = append(schemaOptions, withContainerID(b))
			isContainerID = b
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
		if l, isLogger := options[SQL_OPTION_LOGGER]; isLogger {
			logger = l.(*zap.Logger)
		}
	}

	// ants pool
	antsPool, _ := ants.NewPool(-1, ants.WithLogger(&zapLogger{logger.Named("ants")}))
	promisePool := promise.FromAntsPool(antsPool)

	// encoder structure
	return &SQLEncoder{
		sqlEncoder: sqlEncoder{
			sqlBase: &sqlBase{
				codecVersion:     codecVersion,
				uri:              resource.GetURI(),
				objectIDManager:  objectIDManager,
				objectIDName:     objectIDName,
				isContainerID:    isContainerID,
				isObjectID:       isObjectID,
				schema:           newSqlSchema(schemaOptions...),
				sqliteQueries:    map[string][]*query{},
				logger:           logger,
				antsPool:         antsPool,
				promisePool:      promisePool,
				connPoolProvider: connectionPoolProvider,
				connPoolClose:    connectionPoolClose,
			},
			isForced:         true,
			isKeepDefaults:   isKeepDefaults,
			classDataMap:     map[EClass]*sqlEncoderClassData{},
			sqlIDManager:     sqlIDManager,
			sqlObjectManager: newSqlEncoderObjectManager(),
			sqlLockManager:   newSqlEncoderLockManager(),
		},
		resource: resource,
	}
}

func (e *SQLEncoder) EncodeResource() {
	var err error
	if e.connPool, err = e.connPoolProvider(); err != nil {
		e.addError(err)
		return
	}
	defer func() {
		if err := e.connPoolClose(e.connPool); err != nil {
			e.addError(err)
		}

		e.antsPool.Release()
	}()

	if err := e.encodeVersion(); err != nil {
		e.addError(err)
		return
	}

	if err := e.encodeSchema(); err != nil {
		e.addError(err)
		return
	}

	if err := e.encodeProperties(); err != nil {
		e.addError(err)
		return
	}

	// encode contents into db
	if contents := e.resource.GetContents(); !contents.Empty() {
		object := contents.Get(0).(EObject)
		if err := e.encodeContent(object); err != nil {
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
