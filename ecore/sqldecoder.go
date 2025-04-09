package ecore

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/chebyrash/promise"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

func decodeAny(stmt *sqlite.Stmt, i int) any {
	switch stmt.ColumnType(i) {
	case sqlite.TypeNull:
		return nil
	case sqlite.TypeInteger:
		return stmt.ColumnInt64(i)
	case sqlite.TypeText:
		return stmt.ColumnText(i)
	case sqlite.TypeFloat:
		return stmt.ColumnFloat(i)
	case sqlite.TypeBlob:
		bytes := make([]byte, stmt.ColumnLen(i))
		stmt.ColumnBytes(i, bytes)
		return bytes
	default:
		panic("sqlite type not supported")
	}
}

type sqlObjectManager interface {
	registerObject(EObject)
}

type SQLDecoderIDManager interface {
	SQLCodecIDManager
	// Returns package with a specific id
	// Returns (package with request id otherwise nil, decoded as true if package is decoded otherwise false)
	GetPackageFromID(int64) (EPackage, bool)
	// Returns object with a specific id
	// Returns (object with request id otherwise nil, decoded as true if object is decoded otherwise false)
	GetObjectFromID(int64) (EObject, bool)
	// Returns class with a specific id
	// Returns (class with request id otherwise nil, decoded as true if class is decoded otherwise false)
	GetClassFromID(int64) (EClass, bool)
	// Returns enum literal with a specific id
	// Returns (enum literal with request id otherwise nil, decoded as true if enum literal is decoded otherwise false)
	GetEnumLiteralFromID(int64) (EEnumLiteral, bool)
}

type sqlDecoderClassData struct {
	eClass   EClass
	eFactory EFactory
}

type sqlDecoderIDManagerImpl struct {
	packages     map[int64]EPackage
	objects      map[int64]EObject
	classes      map[int64]EClass
	enumLiterals map[int64]EEnumLiteral
}

func newSQLDecoderIDManager() SQLDecoderIDManager {
	return &sqlDecoderIDManagerImpl{
		packages:     map[int64]EPackage{},
		objects:      map[int64]EObject{},
		classes:      map[int64]EClass{},
		enumLiterals: map[int64]EEnumLiteral{},
	}
}

func (r *sqlDecoderIDManagerImpl) SetPackageID(p EPackage, id int64) {
	r.packages[id] = p
}

func (r *sqlDecoderIDManagerImpl) GetPackageFromID(id int64) (p EPackage, b bool) {
	p, b = r.packages[id]
	return
}

func (r *sqlDecoderIDManagerImpl) SetObjectID(o EObject, id int64) {
	r.objects[id] = o

	// set sql id if created object is an sql object
	if sqlObject, _ := o.(SQLObject); sqlObject != nil {
		sqlObject.SetSQLID(id)
	}
}

func (r *sqlDecoderIDManagerImpl) GetObjectFromID(id int64) (o EObject, b bool) {
	o, b = r.objects[id]
	return
}

func (r *sqlDecoderIDManagerImpl) SetClassID(c EClass, id int64) {
	r.classes[id] = c
}

func (r *sqlDecoderIDManagerImpl) GetClassFromID(id int64) (c EClass, b bool) {
	c, b = r.classes[id]
	return
}

func (r *sqlDecoderIDManagerImpl) SetEnumLiteralID(e EEnumLiteral, id int64) {
	r.enumLiterals[id] = e
}

func (r *sqlDecoderIDManagerImpl) GetEnumLiteralFromID(id int64) (e EEnumLiteral, b bool) {
	e, b = r.enumLiterals[id]
	return
}

type sqlDecoderObjectManager struct {
}

func newSqlDecoderObjectManager() *sqlDecoderObjectManager {
	return &sqlDecoderObjectManager{}
}

func (r *sqlDecoderObjectManager) registerObject(EObject) {
}

type sqlDecoder struct {
	*sqlBase
	classDataMap     map[EClass]*sqlDecoderClassData
	packageRegistry  EPackageRegistry
	sqlIDManager     SQLDecoderIDManager
	sqlObjectManager sqlObjectManager
}

func (d *sqlDecoder) resolveURI(uri *URI) *URI {
	if d.uri != nil {
		return d.uri.Resolve(uri)
	}
	return uri
}

func (d *sqlDecoder) decodeVersion() error {
	var version int64
	if err := d.executeQueryTransient("PRAGMA user_version;", &sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			version = stmt.ColumnInt64(0)
			return nil
		},
	}); err != nil {
		return err
	}

	if version > 0 && version != d.codecVersion {
		return fmt.Errorf("codec version %v is not supported", version)
	}
	return nil
}

func (d *sqlDecoder) decodeSchema(schemaOptions []sqlSchemaOption) error {
	// properties
	properties, err := d.decodeProperties()
	if err != nil {
		return nil
	}

	// object id column name
	if propertyObjectID := properties["objectID"]; d.objectIDManager != nil && len(propertyObjectID) > 0 {
		d.objectIDName = propertyObjectID
		d.isObjectID = d.objectIDName != "objectID" && d.objectIDManager != nil
	}

	// container id
	if propertyContainerID, isPropertyContainerID := properties["containerID"]; isPropertyContainerID {
		d.isContainerID, err = strconv.ParseBool(propertyContainerID)
		if err != nil {
			return err
		}
	}

	// create schema
	if d.isObjectID {
		schemaOptions = append(schemaOptions, withObjectIDName(d.objectIDName))
	}
	if d.isContainerID {
		schemaOptions = append(schemaOptions, withContainerID(d.isContainerID))
	}

	d.schema = newSqlSchema(schemaOptions...)
	return nil
}

func (d *sqlDecoder) decodePackage(id int64) (EPackage, error) {
	ePackage, isPackage := d.sqlIDManager.GetPackageFromID(id)
	if !isPackage {
		table := d.schema.packagesTable
		var packageID int64
		var packageURI string
		if err := d.executeQuery(
			table.selectQuery(nil, table.keyName()+"=?", ""),
			&sqlitex.ExecOptions{
				Args: []any{id},
				ResultFunc: func(stmt *sqlite.Stmt) error {
					packageID = stmt.ColumnInt64(0)
					packageURI = stmt.ColumnText(1)
					return nil
				},
			}); err != nil {
			return nil, err
		}

		// retrieve package
		if d.packageRegistry == nil {
			return nil, fmt.Errorf("package registry not defined in sql decoder")
		}
		ePackage = d.packageRegistry.GetPackage(packageURI)
		if ePackage == nil {
			return nil, fmt.Errorf("unable to find package '%s'", packageURI)
		}

		// register package id
		d.sqlIDManager.SetPackageID(ePackage, packageID)
	}
	return ePackage, nil
}

func (d *sqlDecoder) decodeClass(id int64) (*sqlDecoderClassData, error) {
	eClass, isClass := d.sqlIDManager.GetClassFromID(id)
	if !isClass {
		table := d.schema.classesTable
		var className string
		var packageID int64
		if err := d.executeQuery(
			table.selectQuery(nil, table.keyName()+"=?", ""),
			&sqlitex.ExecOptions{
				Args: []any{id},
				ResultFunc: func(stmt *sqlite.Stmt) error {
					packageID = stmt.ColumnInt64(1)
					className = stmt.ColumnText(2)
					return nil
				},
			}); err != nil {
			return nil, err
		}

		// retrieve package
		ePackage, err := d.decodePackage(packageID)
		if err != nil {
			return nil, err
		}
		// retrieve class
		eClass, _ = ePackage.GetEClassifier(className).(EClass)
		if eClass == nil {
			return nil, fmt.Errorf("unable to find class '%s' in package '%s'", className, ePackage.GetNsURI())
		}
		// set class id
		d.sqlIDManager.SetClassID(eClass, id)
	}
	return d.getDecoderClassData(eClass), nil

}

func (d *sqlDecoder) getDecoderClassData(eClass EClass) *sqlDecoderClassData {
	classData, isClassData := d.classDataMap[eClass]
	if !isClassData {
		classData = &sqlDecoderClassData{
			eClass:   eClass,
			eFactory: eClass.GetEPackage().GetEFactoryInstance(),
		}
		d.classDataMap[eClass] = classData
	}
	return classData
}

func (d *sqlDecoder) decodeObject(id int64) (EObject, error) {
	eObject, isObject := d.sqlIDManager.GetObjectFromID(id)
	if !isObject {
		table := d.schema.objectsTable
		var classID int64
		var objectID string
		columns := []string{"classID"}
		if d.isObjectID {
			columns = append(columns, d.objectIDName)
		}
		if err := d.executeQuery(
			table.selectQuery(columns, table.keyName()+"=?", ""),
			&sqlitex.ExecOptions{
				Args: []any{id},
				ResultFunc: func(stmt *sqlite.Stmt) error {
					isObject = true
					classID = stmt.ColumnInt64(0)
					if d.isObjectID {
						objectID = stmt.ColumnText(1)
					}
					return nil
				},
			}); err != nil {
			return nil, err
		}

		if !isObject {
			return nil, fmt.Errorf("object with id '%v' doesn't exists", id)
		}

		// retrieve class data
		classData, err := d.decodeClass(classID)
		if err != nil {
			return nil, err
		}

		// create object
		eObject = classData.eFactory.Create(classData.eClass)

		// register its id
		if d.objectIDManager != nil {
			if d.isObjectID {
				// object id is column id
				if err := d.objectIDManager.SetID(eObject, objectID); err != nil {
					return nil, err
				}
			} else if d.objectIDName == table.key.columnName {
				// object id is sql id
				if err := d.objectIDManager.SetID(eObject, id); err != nil {
					return nil, err
				}
			}

		}

		// register in sql id manager
		d.sqlIDManager.SetObjectID(eObject, id)

		// register in sql object maneger
		d.sqlObjectManager.registerObject(eObject)
	}
	return eObject, nil
}

func (d *sqlDecoder) decodeEnumLiteral(id int64) (EEnumLiteral, error) {
	eEnumLiteral, isEnumLiteral := d.sqlIDManager.GetEnumLiteralFromID(id)
	if !isEnumLiteral {
		table := d.schema.enumsTable
		var enumID int64
		var packageID int64
		var enumName string
		var literalValue string
		if err := d.executeQuery(
			table.selectQuery(nil, table.keyName()+"=?", ""),
			&sqlitex.ExecOptions{
				Args: []any{id},
				ResultFunc: func(stmt *sqlite.Stmt) error {
					enumID = stmt.ColumnInt64(0)
					packageID = stmt.ColumnInt64(1)
					enumName = stmt.ColumnText(2)
					literalValue = stmt.ColumnText(3)
					return nil
				},
			}); err != nil {
			return nil, err
		}

		// package
		ePackage, err := d.decodePackage(packageID)
		if err != nil {
			return nil, err
		}

		// enum
		eEnum, _ := ePackage.GetEClassifier(enumName).(EEnum)
		if eEnum == nil {
			return nil, fmt.Errorf("unable to find enum '%s' in package '%s'", enumName, ePackage.GetName())
		}

		eEnumLiteral = eEnum.GetEEnumLiteralByLiteral(literalValue)
		if eEnumLiteral == nil {
			return nil, fmt.Errorf("unable to find enum literal '%s' in enum '%s'", literalValue, eEnum.GetName())
		}

		// register enum value
		d.sqlIDManager.SetEnumLiteralID(eEnumLiteral, enumID)
	}
	return eEnumLiteral, nil
}

func (d *sqlDecoder) decodeFeatureValue(featureData *sqlFeatureSchema, value any) (any, error) {
	switch featureData.featureKind {
	case sfkObject, sfkObjectList:
		switch v := value.(type) {
		case nil:
			return nil, nil
		case int64:
			return d.decodeObject(v)
		default:
			return nil, fmt.Errorf("%v is not supported as a object id", v)
		}
	case sfkObjectReference, sfkObjectReferenceList:
		switch v := value.(type) {
		case nil:
			return nil, nil
		case string:
			// no reference
			if len(v) == 0 {
				return nil, nil
			} else {
				sqlID, err := strconv.ParseInt(v, 10, 64)
				if err == nil {
					// string is an int => object is in the resource
					// decode it
					return d.decodeObject(sqlID)
				} else {
					// string is an uri
					proxyURI := NewURI(v)
					resolvedURI := d.resolveURI(proxyURI)
					// create proxy
					eFeature := featureData.feature
					eClass := eFeature.GetEType().(EClass)
					eFactory := eClass.GetEPackage().GetEFactoryInstance()
					eObject := eFactory.Create(eClass)
					eObjectInternal := eObject.(EObjectInternal)
					eObjectInternal.ESetProxyURI(resolvedURI)
					return eObject, nil
				}
			}
		default:
			return nil, fmt.Errorf("%v is not supported as a object reference uri", v)
		}
	case sfkBool:
		switch v := value.(type) {
		case nil:
			return false, nil
		case int64:
			return v == 1, nil
		default:
			return nil, fmt.Errorf("%v is not a bool value", v)
		}
	case sfkByte:
		switch v := value.(type) {
		case nil:
			return byte(0), nil
		case int64:
			return byte(v), nil
		default:
			return nil, fmt.Errorf("%v is not a byte value", v)
		}
	case sfkInt:
		switch v := value.(type) {
		case nil:
			return 0, nil
		case int64:
			return int(v), nil
		default:
			return nil, fmt.Errorf("%v is not a int value", v)
		}
	case sfkInt64:
		switch v := value.(type) {
		case nil:
			return int64(0), nil
		case int64:
			return v, nil
		default:
			return nil, fmt.Errorf("%v is not a int64 value", v)
		}
	case sfkInt32:
		switch v := value.(type) {
		case nil:
			return int32(0), nil
		case int64:
			return int32(v), nil
		default:
			return nil, fmt.Errorf("%v is not a int32 value", v)
		}
	case sfkInt16:
		switch v := value.(type) {
		case nil:
			return int16(0), nil
		case int64:
			return int16(v), nil
		default:
			return nil, fmt.Errorf("%v is not a int16 value", v)
		}
	case sfkEnum:
		switch v := value.(type) {
		case nil:
			return featureData.feature.GetDefaultValue(), nil
		case int64:
			enumLiteral, err := d.decodeEnumLiteral(v)
			if err != nil {
				return nil, err
			}
			instance := enumLiteral.GetInstance()
			if instance == nil {
				instance = d.decodeFeatureData(featureData, enumLiteral.GetLiteral())
			}
			return instance, nil
		default:
			return nil, fmt.Errorf("%v is not a enum value", value)
		}
	case sfkString:
		switch v := value.(type) {
		case nil:
			return "", nil
		case string:
			return v, nil
		default:
			return "", fmt.Errorf("%v is not a string value", v)
		}
	case sfkByteArray:
		switch v := value.(type) {
		case nil:
			return nil, nil
		case []byte:
			return v, nil
		default:
			return "", fmt.Errorf("%v is not a byte array value", v)
		}
	case sfkDate:
		switch v := value.(type) {
		case nil:
			return nil, nil
		case time.Time:
			return &v, nil
		case string:
			t, err := time.Parse(time.RFC3339, v)
			if err != nil {
				return nil, err
			}
			return &t, nil
		default:
			return nil, fmt.Errorf("%v is not a time value", v)
		}
	case sfkFloat64:
		switch v := value.(type) {
		case nil:
			return float64(0), nil
		case float64:
			return v, nil
		default:
			return nil, fmt.Errorf("%v is not a float64 value", value)
		}
	case sfkFloat32:
		switch v := value.(type) {
		case nil:
			return float32(0), nil
		case float64:
			return float32(v), nil
		default:
			return nil, fmt.Errorf("%v is not a float64 value", value)
		}
	case sfkData, sfkDataList:
		switch v := value.(type) {
		case nil:
			return nil, nil
		case []byte:
			return d.decodeFeatureData(featureData, string(v)), nil
		case string:
			return d.decodeFeatureData(featureData, v), nil
		default:
			return nil, fmt.Errorf("%v is not a data value", value)
		}
	}

	return nil, nil
}

func (d *sqlDecoder) decodeFeatureData(featureSchema *sqlFeatureSchema, v string) any {
	eFeature := featureSchema.feature
	eDataType := eFeature.GetEType().(EDataType)
	eFactory := eDataType.GetEPackage().GetEFactoryInstance()
	return eFactory.CreateFromString(eDataType, v)
}

type SQLDecoder struct {
	sqlDecoder
	resource     EResource
	withFeatures bool
	withObjects  bool
}

const SQLITE_MAX_ALLOCATION_SIZE = 2147483391

func newMemoryConnectionPool(dbName string, dbPath string, r io.Reader) (*sqlitex.Pool, error) {
	// create a memory connection
	var connSrc *sqlite.Conn

	// initialize db with reader bytes
	if r != nil {
		bytes, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}

		// set journal mode as rolling back ( WAL is not supported )
		bytes[18] = 0x01
		bytes[19] = 0x01

		// sqlite.Conn.Deserialize method has a max allocation size
		// so we must use an intermediate file

		// create tmp file
		if dbPath == "" {
			dbPath, err = sqlTmpDB(dbName)
			if err != nil {
				return nil, err
			}
		}

		tmpDBFile, err := os.Create(dbPath)
		if err != nil {
			return nil, err
		}

		// write bytes inside it
		_, err = tmpDBFile.Write(bytes)
		if err != nil {
			tmpDBFile.Close()
			return nil, err
		}

		// close tmp file
		if err = tmpDBFile.Close(); err != nil {
			return nil, err
		}

		// open connection to tmp file
		connSrc, err = sqlite.OpenConn(dbPath, sqlite.OpenReadOnly)
		if err != nil {
			return nil, err
		}
		defer connSrc.Close()
	}

	// create connection pool
	dbPath = fmt.Sprintf("file:%s?mode=memory&cache=shared", dbName)
	connPool, err := sqlitex.NewPool(dbPath, sqlitex.PoolOptions{Flags: sqlite.OpenCreate | sqlite.OpenReadWrite | sqlite.OpenURI})
	if err != nil {
		return nil, err
	}

	// create connection pool
	connDst, err := connPool.Take(context.Background())
	if err != nil {
		return nil, err
	}
	defer connPool.Put(connDst)

	if connSrc != nil {
		// backup src db to dst db
		backup, err := sqlite.NewBackup(connDst, "main", connSrc, "main")
		if err != nil {
			return nil, err
		}
		if more, err := backup.Step(-1); err != nil {
			return nil, err
		} else if more {
			return nil, errors.New("full backup step with remaining pages")
		}
		if err := backup.Close(); err != nil {
			return nil, err
		}

	}

	return connPool, nil
}

func newFileConnectionPool(dbName string, dbPath string, r io.Reader) (p *sqlitex.Pool, err error) {
	if dbPath == "" {
		dbPath, err = sqlTmpDB(dbName)
		if err != nil {
			return nil, err
		}
	}

	dbFile, err := os.Create(dbPath)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(dbFile, r)
	if err != nil {
		dbFile.Close()
		return nil, err
	}
	dbFile.Close()

	return sqlitex.NewPool(dbPath, sqlitex.PoolOptions{Flags: sqlite.OpenReadOnly})
}

func NewSQLReaderDecoder(r io.Reader, resource EResource, options map[string]any) *SQLDecoder {
	return newSQLDecoder(
		func() (*sqlitex.Pool, error) {
			dbName := filepath.Base(resource.GetURI().Path())
			dbPath := ""
			inMemoryDatabase := false
			if options != nil {
				inMemoryDatabase, _ = options[SQL_OPTION_IN_MEMORY_DATABASE].(bool)
				dbPath, _ = options[SQL_OPTION_DECODER_DB_PATH].(string)
			}
			if inMemoryDatabase {
				return newMemoryConnectionPool(dbName, dbPath, r)
			} else {
				return newFileConnectionPool(dbName, dbPath, r)
			}
		},
		func(connPool *sqlitex.Pool) error {
			return connPool.Close()
		},
		resource,
		options)
}

func NewSQLDBDecoder(connPool *sqlitex.Pool, resource EResource, options map[string]any) *SQLDecoder {
	return newSQLDecoder(
		func() (*sqlitex.Pool, error) {
			return connPool, nil
		},
		func(pool *sqlitex.Pool) error {
			return nil
		},
		resource,
		options,
	)
}

type zapLogger struct {
	*zap.Logger
}

func (l *zapLogger) Printf(format string, args ...any) {
	l.Info(fmt.Sprintf(format, args...))
}

func newSQLDecoder(connectionPoolProvider func() (*sqlitex.Pool, error), connectionPoolClose func(conn *sqlitex.Pool) error, resource EResource, options map[string]any) *SQLDecoder {
	codecVersion := sqlCodecVersion
	sqlIDManager := newSQLDecoderIDManager()
	logger := zap.NewNop()
	withFeatures := true
	withObjects := true
	if options != nil {
		if v, isVersion := options[SQL_OPTION_CODEC_VERSION].(int64); isVersion {
			codecVersion = v
		} else if v, isVersion := options[SQL_OPTION_CODEC_VERSION].(int); isVersion {
			codecVersion = int64(v)
		}
		if m, isSQLIDManager := options[SQL_OPTION_SQL_ID_MANAGER].(SQLDecoderIDManager); isSQLIDManager {
			sqlIDManager = m
		}
		if l, isLogger := options[SQL_OPTION_LOGGER]; isLogger {
			logger = l.(*zap.Logger)
		}
		if b, isWithFeatures := options[SQL_OPTION_DECODER_WITH_FEATURES].(bool); isWithFeatures {
			withFeatures = b
		}
		if b, isWithObjects := options[SQL_OPTION_DECODER_WITH_OBJECTS].(bool); isWithObjects {
			withObjects = b
		}
	}

	// package registry
	packageRegistry := GetPackageRegistry()
	resourceSet := resource.GetResourceSet()
	if resourceSet != nil {
		packageRegistry = resourceSet.GetPackageRegistry()
	}

	// ants pool
	antsPool, _ := ants.NewPool(-1, ants.WithLogger(&zapLogger{logger.Named("ants")}))
	promisePool := promise.FromAntsPool(antsPool)

	return &SQLDecoder{
		sqlDecoder: sqlDecoder{
			sqlBase: &sqlBase{
				codecVersion:     codecVersion,
				uri:              resource.GetURI(),
				objectIDManager:  resource.GetObjectIDManager(),
				sqliteQueries:    map[string][]*query{},
				logger:           logger,
				antsPool:         antsPool,
				promisePool:      promisePool,
				connPoolProvider: connectionPoolProvider,
				connPoolClose:    connectionPoolClose,
			},
			packageRegistry:  packageRegistry,
			sqlIDManager:     sqlIDManager,
			sqlObjectManager: newSqlDecoderObjectManager(),
			classDataMap:     map[EClass]*sqlDecoderClassData{},
		},
		resource:     resource,
		withFeatures: withFeatures,
		withObjects:  withObjects,
	}
}

func (d *SQLDecoder) DecodeResource() {
	var err error
	if d.connPool, err = d.connPoolProvider(); err != nil {
		d.addError(err)
		return
	}
	defer func() {
		if err := d.connPoolClose(d.connPool); err != nil {
			d.addError(err)
		}

		d.antsPool.Release()
	}()

	if err := d.decodeVersion(); err != nil {
		d.addError(err)
		return
	}

	if err := d.decodeSchema(nil); err != nil {
		d.addError(err)
		return
	}

	if err := d.decodePackages(); err != nil {
		d.addError(err)
		return
	}

	if err := d.decodeEnums(); err != nil {
		d.addError(err)
		return
	}

	if err := d.decodeClasses(); err != nil {
		d.addError(err)
		return
	}

	if d.withObjects {
		if err := d.decodeObjects(); err != nil {
			d.addError(err)
			return
		}
	}

	if d.withFeatures {
		if err := d.decodeFeatures(); err != nil {
			d.addError(err)
			return
		}
	}

	if err := d.decodeContents(); err != nil {
		d.addError(err)
		return
	}
}

func (d *SQLDecoder) DecodeObject() (EObject, error) {
	panic("SQLDecoder doesn't support object decoding")
}

func (d *SQLDecoder) decodeContents() error {
	return d.executeQueryTransient(
		d.schema.contentsTable.selectQuery(nil, "", ""),
		&sqlitex.ExecOptions{
			ResultFunc: func(stmt *sqlite.Stmt) error {
				objectID := stmt.ColumnInt64(0)
				object, err := d.decodeObject(objectID)
				if err != nil {
					return err
				}
				d.resource.GetContents().Add(object)
				return nil
			},
		})
}

func (d *SQLDecoder) decodePackages() error {
	return d.executeQueryTransient(
		d.schema.packagesTable.selectQuery(nil, "", ""),
		&sqlitex.ExecOptions{
			ResultFunc: func(stmt *sqlite.Stmt) error {
				packageID := stmt.ColumnInt64(0)
				packageURI := stmt.ColumnText(1)
				ePackage := d.packageRegistry.GetPackage(packageURI)
				if ePackage == nil {
					return fmt.Errorf("unable to find package '%s'", packageURI)
				}
				d.sqlIDManager.SetPackageID(ePackage, packageID)
				return nil
			},
		})
}

func (d *SQLDecoder) decodeEnums() error {
	return d.executeQueryTransient(
		d.schema.enumsTable.selectQuery(nil, "", ""),
		&sqlitex.ExecOptions{
			ResultFunc: func(stmt *sqlite.Stmt) error {
				enumID := stmt.ColumnInt64(0)
				packageID := stmt.ColumnInt64(1)
				enumName := stmt.ColumnText(2)
				literalValue := stmt.ColumnText(3)

				ePackage, _ := d.sqlIDManager.GetPackageFromID(packageID)
				if ePackage == nil {
					return fmt.Errorf("unable to find package with id '%d'", packageID)
				}

				eEnum, _ := ePackage.GetEClassifier(enumName).(EEnum)
				if eEnum == nil {
					return fmt.Errorf("unable to find enum '%s' in package '%s'", enumName, ePackage.GetNsURI())
				}

				eEnumLiteral := eEnum.GetEEnumLiteralByLiteral(literalValue)
				if eEnumLiteral == nil {
					return fmt.Errorf("unable to find enum literal '%s' in enum '%s'", literalValue, eEnum.GetName())
				}

				// create map enum
				d.sqlIDManager.SetEnumLiteralID(eEnumLiteral, enumID)
				return nil
			},
		})
}

func (d *SQLDecoder) decodeClasses() error {
	return d.executeQueryTransient(
		d.schema.classesTable.selectQuery(nil, "", ""),
		&sqlitex.ExecOptions{
			ResultFunc: func(stmt *sqlite.Stmt) error {
				classID := stmt.ColumnInt64(0)
				packageID := stmt.ColumnInt64(1)
				className := stmt.ColumnText(2)
				ePackage, _ := d.sqlIDManager.GetPackageFromID(packageID)
				if ePackage == nil {
					return fmt.Errorf("unable to find package with id '%d'", packageID)
				}
				eClass, _ := ePackage.GetEClassifier(className).(EClass)
				if eClass == nil {
					return fmt.Errorf("unable to find class '%s' in package '%s'", className, ePackage.GetNsURI())
				}

				d.sqlIDManager.SetClassID(eClass, classID)

				d.classDataMap[eClass] = &sqlDecoderClassData{
					eClass:   eClass,
					eFactory: ePackage.GetEFactoryInstance(),
				}
				return nil
			},
		})
}

func (d *SQLDecoder) decodeObjects() error {
	return d.executeQueryTransient(
		d.schema.objectsTable.selectQuery(nil, "", ""),
		&sqlitex.ExecOptions{
			ResultFunc: func(stmt *sqlite.Stmt) error {
				sqlObjectID := stmt.ColumnInt64(0)
				classID := stmt.ColumnInt64(1)
				eClass, _ := d.sqlIDManager.GetClassFromID(classID)
				if eClass == nil {
					return fmt.Errorf("unable to find class with id '%v'", classID)
				}

				// class data
				classData := d.classDataMap[eClass]
				if classData == nil {
					return fmt.Errorf("unable to find class data with id '%v'", classID)
				}

				// create & register object
				eObject := classData.eFactory.Create(classData.eClass)
				d.sqlIDManager.SetObjectID(eObject, sqlObjectID)

				// set its id
				if d.objectIDManager != nil {
					if d.isObjectID {
						switch stmt.ColumnType(2) {
						case sqlite.TypeNull:
						case sqlite.TypeText:
							objectID := stmt.ColumnText(2)
							if err := d.objectIDManager.SetID(eObject, objectID); err != nil {
								return err
							}
						}
					} else if d.objectIDName == d.schema.objectsTable.key.columnName {
						// object id is sql id
						if err := d.objectIDManager.SetID(eObject, sqlObjectID); err != nil {
							return err
						}
					}
				}
				return nil
			},
		})
}

func (d *SQLDecoder) decodeFeatures() error {
	decoded := map[EClass]struct{}{}
	decoding := []*promise.Promise[any]{}
	// for each leaf class
	for _, classData := range d.classDataMap {
		eClass := classData.eClass
		itSuper := eClass.GetEAllSuperTypes().Iterator()
		for eClass != nil {
			// decode class features
			decodingClass := eClass
			if _, isDecoded := decoded[decodingClass]; !isDecoded {
				decoded[decodingClass] = struct{}{}
				decoding = append(decoding, promise.NewWithPool(func(resolve func(any), reject func(error)) {
					if err := d.decodeClassFeatures(decodingClass); err != nil {
						reject(err)
					} else {
						resolve(nil)
					}
				}, d.promisePool))
			}

			// next super class
			if itSuper.HasNext() {
				eClass = itSuper.Next().(EClass)
			} else {
				eClass = nil
			}
		}
	}
	if len(decoding) > 0 {
		_, err := promise.All(context.Background(), decoding...).Await(context.Background())
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *SQLDecoder) decodeClassFeatures(eClass EClass) error {
	classSchema := d.schema.getClassSchema(eClass)
	columnFeatures := []*sqlFeatureSchema{}
	for _, featureData := range classSchema.features {
		if featureData.column != nil {
			columnFeatures = append(columnFeatures, featureData)
		} else if featureData.table != nil {
			if err := d.decodeTableFeature(featureData.table, featureData); err != nil {
				return err
			}
		}
	}

	return d.decodeColumnFeatures(classSchema.table, columnFeatures)
}

func (d *SQLDecoder) decodeColumnFeatures(table *sqlTable, columnFeatures []*sqlFeatureSchema) error {
	if len(columnFeatures) == 0 {
		return nil
	}
	return d.executeQuery(
		table.selectQuery(nil, "", ""),
		&sqlitex.ExecOptions{
			ResultFunc: func(stmt *sqlite.Stmt) error {
				objectID := stmt.ColumnInt64(0)
				eObject, _ := d.sqlIDManager.GetObjectFromID(objectID)
				if eObject == nil {
					return fmt.Errorf("unable to find object with id '%v'", objectID)
				}

				// for each column
				for i, columnData := range columnFeatures {
					value := decodeAny(stmt, i+1)
					columnValue, err := d.decodeFeatureValue(columnData, value)
					if err != nil {
						return err
					}
					// column value is nil, if column is not set
					// so use default value
					if columnValue != nil {
						eObject.ESet(columnData.feature, columnValue)
					}
				}
				return nil
			},
		})
}

func (d *SQLDecoder) decodeTableFeature(table *sqlTable, tableFeature *sqlFeatureSchema) error {
	column := sqlEscapeIdentifier(table.columns[len(table.columns)-1].columnName)
	query := table.selectQuery([]string{table.keyName(), column}, "", table.keyName()+" ASC, idx ASC")
	feature := tableFeature.feature
	values := []any{}
	var id int64 = -1
	if err := d.executeQuery(
		query,
		&sqlitex.ExecOptions{
			ResultFunc: func(stmt *sqlite.Stmt) error {
				// object id
				objectID := stmt.ColumnInt64(0)
				value := decodeAny(stmt, 1)
				decoded, err := d.decodeFeatureValue(tableFeature, value)
				if err != nil {
					return err
				}

				if id == -1 {
					id = objectID
				} else if id != objectID {
					if err := d.decodeFeatureList(id, feature, values); err != nil {
						return err
					}
					values = nil
				}
				id = objectID
				values = append(values, decoded)
				return nil
			},
		}); err != nil {
		return err
	}
	if id != -1 {
		if err := d.decodeFeatureList(id, feature, values); err != nil {
			return err
		}
	}
	return nil
}

func (d *SQLDecoder) decodeFeatureList(objectID int64, feature EStructuralFeature, values []any) error {
	if len(values) == 0 {
		return nil
	}
	eObject, _ := d.sqlIDManager.GetObjectFromID(objectID)
	if eObject == nil {
		return fmt.Errorf("unable to find object with id '%v'", objectID)
	}
	eList := eObject.EGetResolve(feature, false).(EList)
	eList.AddAll(NewImmutableEList(values))
	return nil
}

func (d *SQLDecoder) addError(err error) {
	d.resource.GetErrors().Add(NewEDiagnosticImpl(err.Error(), d.resource.GetURI().String(), 0, 0))
}
