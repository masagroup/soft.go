package ecore

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type sqlObjectManager interface {
	registerObject(EObject)
}

type sqlDecoderIDManager interface {
	sqlCodecIDManager
	getPackageFromID(int64) (EPackage, bool)
	getObjectFromID(int64) (EObject, bool)
	getClassFromID(int64) (EClass, bool)
	getEnumLiteralFromID(int64) (EEnumLiteral, bool)
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

func newSqlDecoderIDManager() *sqlDecoderIDManagerImpl {
	return &sqlDecoderIDManagerImpl{
		packages:     map[int64]EPackage{},
		objects:      map[int64]EObject{},
		classes:      map[int64]EClass{},
		enumLiterals: map[int64]EEnumLiteral{},
	}
}

func (r *sqlDecoderIDManagerImpl) setPackageID(p EPackage, id int64) {
	r.packages[id] = p
}

func (r *sqlDecoderIDManagerImpl) getPackageFromID(id int64) (p EPackage, b bool) {
	p, b = r.packages[id]
	return
}

func (r *sqlDecoderIDManagerImpl) setObjectID(o EObject, id int64) {
	r.objects[id] = o

	// set sql id if created object is an sql object
	if sqlObject, _ := o.(SQLObject); sqlObject != nil {
		sqlObject.SetSqlID(id)
	}
}

func (r *sqlDecoderIDManagerImpl) getObjectFromID(id int64) (o EObject, b bool) {
	o, b = r.objects[id]
	return
}

func (r *sqlDecoderIDManagerImpl) setClassID(c EClass, id int64) {
	r.classes[id] = c
}

func (r *sqlDecoderIDManagerImpl) getClassFromID(id int64) (c EClass, b bool) {
	c, b = r.classes[id]
	return
}

func (r *sqlDecoderIDManagerImpl) setEnumLiteralID(e EEnumLiteral, id int64) {
	r.enumLiterals[id] = e
}

func (r *sqlDecoderIDManagerImpl) getEnumLiteralFromID(id int64) (e EEnumLiteral, b bool) {
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
	selectStmts      map[*sqlTable]*sqlSafeStmt
	classDataMap     map[EClass]*sqlDecoderClassData
	packageRegistry  EPackageRegistry
	sqlIDManager     sqlDecoderIDManager
	sqlObjectManager sqlObjectManager
}

func (d *sqlDecoder) resolveURI(uri *URI) *URI {
	if d.uri != nil {
		return d.uri.Resolve(uri)
	}
	return uri
}

func (s *sqlDecoder) getSelectStmt(table *sqlTable, query func() string) (stmt *sqlSafeStmt, err error) {
	stmt = s.selectStmts[table]
	if stmt == nil {
		stmt, err = s.db.Prepare(query())
		s.selectStmts[table] = stmt
	}
	return
}

func (d *sqlDecoder) decodePackage(id int64) (EPackage, error) {
	ePackage, isPackage := d.sqlIDManager.getPackageFromID(id)
	if !isPackage {
		// get select stmt
		table := d.schema.packagesTable
		stmt, err := d.getSelectStmt(table, func() string {
			return table.selectQuery(nil, table.keyName()+"=?", "")
		})
		if err != nil {
			return nil, err
		}

		// query package infos
		row := stmt.QueryRow(id)
		var packageID int64
		var packageURI string
		if err := row.Scan(&packageID, &packageURI); err != nil {
			return nil, err
		}

		// retrieve package
		if d.packageRegistry == nil {
			panic(fmt.Errorf("package registry not defined in sql decoder"))
		}
		ePackage = d.packageRegistry.GetPackage(packageURI)
		if ePackage == nil {
			return nil, fmt.Errorf("unable to find package '%s'", packageURI)
		}

		// register package id
		d.sqlIDManager.setPackageID(ePackage, packageID)
	}
	return ePackage, nil
}

func (d *sqlDecoder) decodeClass(id int64) (*sqlDecoderClassData, error) {
	eClass, isClass := d.sqlIDManager.getClassFromID(id)
	if !isClass {
		// get select stmt
		table := d.schema.classesTable
		stmt, err := d.getSelectStmt(table, func() string {
			return table.selectQuery(nil, table.keyName()+"=?", "")
		})
		if err != nil {
			return nil, err
		}

		// query package infos
		row := stmt.QueryRow(id)
		var classID int64
		var className string
		var packageID int64
		if err := row.Scan(&classID, &packageID, &className); err != nil {
			return nil, err
		}

		// retrieve package
		ePackage, err := d.decodePackage(packageID)
		if err != nil {
			return nil, err
		}
		// retrieve class
		eClass, _ := ePackage.GetEClassifier(className).(EClass)
		if eClass == nil {
			return nil, fmt.Errorf("unable to find class '%s' in package '%s'", className, ePackage.GetNsURI())
		}

		eClassData := &sqlDecoderClassData{
			eClass:   eClass,
			eFactory: ePackage.GetEFactoryInstance(),
		}

		d.classDataMap[eClass] = eClassData
		d.sqlIDManager.setClassID(eClass, id)
		return eClassData, nil
	}
	return d.classDataMap[eClass], nil

}

func (d *sqlDecoder) decodeContents() ([]EObject, error) {
	table := d.schema.contentsTable
	stmt, err := d.getSelectStmt(table, func() string {
		return table.selectQuery(nil, "", "")
	})
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, nil
	}
	defer rows.Close()

	contents := []EObject{}
	for rows.Next() {
		var objectID int64
		if err := rows.Scan(&objectID); err != nil {
			return nil, err
		}

		object, err := d.decodeObject(objectID)
		if err != nil {
			return nil, err
		}

		contents = append(contents, object)
	}
	return contents, nil
}

func (d *sqlDecoder) decodeObject(id int64) (EObject, error) {
	eObject, isObject := d.sqlIDManager.getObjectFromID(id)
	if !isObject {
		table := d.schema.objectsTable
		stmt, err := d.getSelectStmt(table, func() string {
			return table.selectQuery(nil, table.keyName()+"=?", "")
		})
		if err != nil {
			return nil, err
		}

		// query object infos
		rows, err := stmt.Query(id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		// check if we have object id column
		columns, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		isObjectID := len(columns) > 2

		// check if we have a row
		if !rows.Next() {
			return nil, sql.ErrNoRows
		}

		// retrieve all ids
		var classID int64
		var sqlID int64
		var objectID string
		scanArgs := []any{&sqlID, &classID}
		if isObjectID {
			scanArgs = append(scanArgs, &objectID)
		}
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}

		// retrieve class data
		classData, err := d.decodeClass(classID)
		if err != nil {
			return nil, err
		}

		// create object
		eObject = classData.eFactory.Create(classData.eClass)

		// register its id
		if isObjectID && d.idManager != nil {
			if err := d.idManager.SetID(eObject, objectID); err != nil {
				return nil, err
			}
		}

		// register in sql id manager
		d.sqlIDManager.setObjectID(eObject, id)

		// register in sql object maneger
		d.sqlObjectManager.registerObject(eObject)
	}
	return eObject, nil
}

func (d *sqlDecoder) decodeEnumLiteral(id int64) (EEnumLiteral, error) {
	eEnumLiteral, isEnumLiteral := d.sqlIDManager.getEnumLiteralFromID(id)
	if !isEnumLiteral {
		table := d.schema.enumsTable
		stmt, err := d.getSelectStmt(table, func() string {
			return table.selectQuery(nil, table.keyName()+"=?", "")
		})
		if err != nil {
			return nil, err
		}

		// query enum infos
		row := stmt.QueryRow(id)
		var enumID int64
		var packageID int64
		var enumName string
		var literalValue string
		if err := row.Scan(&enumID, &packageID, &enumName, &literalValue); err != nil {
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
		d.sqlIDManager.setEnumLiteralID(eEnumLiteral, enumID)
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
			return nil, nil
		case bool:
			return v, nil
		default:
			return nil, fmt.Errorf("%v is not a bool value", v)
		}
	case sfkByte:
		switch v := value.(type) {
		case nil:
			return nil, nil
		case byte:
			return v, nil
		default:
			return nil, fmt.Errorf("%v is not a bool value", v)
		}
	case sfkInt:
		switch v := value.(type) {
		case nil:
			return nil, nil
		case int64:
			return int(v), nil
		default:
			return nil, fmt.Errorf("%v is not a int value", v)
		}
	case sfkInt64:
		switch v := value.(type) {
		case nil:
			return nil, nil
		case int64:
			return v, nil
		default:
			return nil, fmt.Errorf("%v is not a int64 value", v)
		}
	case sfkInt32:
		switch v := value.(type) {
		case nil:
			return nil, nil
		case int64:
			return int32(v), nil
		default:
			return nil, fmt.Errorf("%v is not a int32 value", v)
		}
	case sfkInt16:
		switch v := value.(type) {
		case nil:
			return nil, nil
		case int64:
			return int16(v), nil
		default:
			return nil, fmt.Errorf("%v is not a int16 value", v)
		}
	case sfkEnum:
		switch v := value.(type) {
		case nil:
			return nil, nil
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
			return nil, nil
		case float64:
			return v, nil
		default:
			return nil, fmt.Errorf("%v is not a float64 value", value)
		}
	case sfkFloat32:
		switch v := value.(type) {
		case nil:
			return nil, nil
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
	resource   EResource
	driver     string
	dbProvider func(driver string) (*sql.DB, error)
	dbClose    func(db *sql.DB) error
}

func NewSQLReaderDecoder(r io.Reader, resource EResource, options map[string]any) *SQLDecoder {
	return newSQLDecoder(
		func(driver string) (*sql.DB, error) {
			fileName := filepath.Base(resource.GetURI().Path())
			dbPath, err := sqlTmpDB(fileName)
			if err != nil {
				return nil, err
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

			return sql.Open(driver, dbPath)

		},
		func(db *sql.DB) error {
			return db.Close()
		},
		resource,
		options)
}

func NewSQLDBDecoder(db *sql.DB, resource EResource, options map[string]any) *SQLDecoder {
	return newSQLDecoder(
		func(driver string) (*sql.DB, error) {
			return db, nil
		},
		func(db *sql.DB) error {
			return nil
		},
		resource,
		options,
	)
}

func newSQLDecoder(dbProvider func(driver string) (*sql.DB, error), dbClose func(db *sql.DB) error, resource EResource, options map[string]any) *SQLDecoder {
	// options
	schemaOptions := []sqlSchemaOption{}
	driver := "sqlite"
	idAttributeName := ""
	if options != nil {
		if d, isDriver := options[SQL_OPTION_DRIVER]; isDriver {
			driver = d.(string)
		}

		idAttributeName, _ = options[SQL_OPTION_ID_ATTRIBUTE_NAME].(string)
		if resource.GetObjectIDManager() != nil && len(idAttributeName) > 0 {
			schemaOptions = append(schemaOptions, withIDAttributeName(idAttributeName))
		}
	}

	// package registry
	packageRegistry := GetPackageRegistry()
	resourceSet := resource.GetResourceSet()
	if resourceSet != nil {
		packageRegistry = resourceSet.GetPackageRegistry()
	}

	return &SQLDecoder{
		sqlDecoder: sqlDecoder{
			sqlBase: &sqlBase{
				uri:             resource.GetURI(),
				idManager:       resource.GetObjectIDManager(),
				idAttributeName: idAttributeName,
				schema:          newSqlSchema(schemaOptions...),
			},
			packageRegistry:  packageRegistry,
			sqlIDManager:     newSqlDecoderIDManager(),
			sqlObjectManager: newSqlDecoderObjectManager(),
			selectStmts:      map[*sqlTable]*sqlSafeStmt{},
			classDataMap:     map[EClass]*sqlDecoderClassData{},
		},
		resource:   resource,
		dbProvider: dbProvider,
		dbClose:    dbClose,
		driver:     driver,
	}
}

func (d *SQLDecoder) DecodeResource() {
	// retrieve db
	db, err := d.dbProvider(d.driver)
	if err != nil {
		d.addError(err)
		return
	}
	defer func() {
		if err := d.dbClose(db); err != nil {
			d.addError(err)
		}
	}()

	// create safe db
	d.db = newSQLSafeDB(db)

	if err := d.decodeVersion(); err != nil {
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

	if err := d.decodeObjects(); err != nil {
		d.addError(err)
		return
	}

	if err := d.decodeFeatures(); err != nil {
		d.addError(err)
		return
	}

	if err := d.decodeContents(); err != nil {
		d.addError(err)
		return
	}
}

func (d *SQLDecoder) DecodeObject() (EObject, error) {
	panic("SQLDecoder doesn't support object decoding")
}

func (d *SQLDecoder) decodeVersion() error {
	if row := d.db.QueryRow("PRAGMA user_version;"); row == nil {
		return fmt.Errorf("unable to retrieve user version")
	} else {
		var v int
		if err := row.Scan(&v); err != nil {
			return err
		}
		if v != sqlCodecVersion {
			return fmt.Errorf("history version %v is not supported", v)
		}
		return nil
	}
}

func (d *SQLDecoder) decodeContents() error {
	// read all object contents
	return d.query(d.schema.contentsTable.selectQuery(nil, "", ""), func(values []driver.Value) error {
		// retrieve id
		objectID, isInt := values[0].(int64)
		if !isInt {
			return fmt.Errorf("%v is not an int64 value", values[0])
		}

		// decode object
		eObject, _ := d.sqlIDManager.getObjectFromID(objectID)
		if eObject == nil {
			return fmt.Errorf("unable to find object with id '%v'", objectID)
		}

		// add object to resource contents
		d.resource.GetContents().Add(eObject)
		return nil
	})
}

func (d *SQLDecoder) decodePackages() error {
	return d.query(d.schema.packagesTable.selectQuery(nil, "", ""), func(values []driver.Value) error {
		packageID, isInt := values[0].(int64)
		if !isInt {
			return fmt.Errorf("%v is not a int64 value", values[0])
		}

		packageURI, isString := values[1].(string)
		if !isString {
			return fmt.Errorf("%v is not a string value", values[1])
		}

		ePackage := d.packageRegistry.GetPackage(packageURI)
		if ePackage == nil {
			return fmt.Errorf("unable to find package '%s'", packageURI)
		}
		d.sqlIDManager.setPackageID(ePackage, packageID)
		return nil
	})
}

func (d *SQLDecoder) decodeEnums() error {
	return d.query(d.schema.enumsTable.selectQuery(nil, "", ""), func(values []driver.Value) error {
		enumID, isInt := values[0].(int64)
		if !isInt {
			return fmt.Errorf("%v is not a int64 value", values[0])
		}

		// package
		packageID, isInt := values[1].(int64)
		if !isInt {
			return fmt.Errorf("%v is not a string value", values[1])
		}
		ePackage, _ := d.sqlIDManager.getPackageFromID(packageID)
		if ePackage == nil {
			return fmt.Errorf("unable to find package with id '%d'", packageID)
		}

		// enum
		enumName, isString := values[2].(string)
		if !isString {
			return fmt.Errorf("%v is not a string value", values[2])
		}
		eEnum, _ := ePackage.GetEClassifier(enumName).(EEnum)
		if eEnum == nil {
			return fmt.Errorf("unable to find enum '%s' in package '%s'", enumName, ePackage.GetNsURI())
		}

		// enum literal value
		literalValue, isString := values[3].(string)
		if !isString {
			return fmt.Errorf("%v is not a string value", values[3])
		}

		eEnumLiteral := eEnum.GetEEnumLiteralByLiteral(literalValue)
		if eEnumLiteral == nil {
			return fmt.Errorf("unable to find enum literal '%s' in enum '%s'", literalValue, eEnum.GetName())
		}

		// create map enum
		d.sqlIDManager.setEnumLiteralID(eEnumLiteral, enumID)
		return nil
	})
}

func (d *SQLDecoder) decodeClasses() error {
	return d.query(d.schema.classesTable.selectQuery(nil, "", ""), func(values []driver.Value) error {
		classID, isInt := values[0].(int64)
		if !isInt {
			return fmt.Errorf("%v is not a int64 value", values[0])
		}

		packageID, isInt := values[1].(int64)
		if !isInt {
			return fmt.Errorf("%v is not a int64 value", values[1])
		}
		ePackage, _ := d.sqlIDManager.getPackageFromID(packageID)
		if ePackage == nil {
			return fmt.Errorf("unable to find package with id '%d'", packageID)
		}
		className, isString := values[2].(string)
		if !isString {
			return fmt.Errorf("%v is not a string value", values[2])
		}
		eClass, _ := ePackage.GetEClassifier(className).(EClass)
		if eClass == nil {
			return fmt.Errorf("unable to find class '%s' in package '%s'", className, ePackage.GetNsURI())
		}

		d.sqlIDManager.setClassID(eClass, classID)

		d.classDataMap[eClass] = &sqlDecoderClassData{
			eClass:   eClass,
			eFactory: ePackage.GetEFactoryInstance(),
		}

		return nil
	})

}

func (d *SQLDecoder) decodeObjects() error {
	return d.query(d.schema.objectsTable.selectQuery(nil, "", ""), func(values []driver.Value) error {
		// object id
		objectID, isInt := values[0].(int64)
		if !isInt {
			return fmt.Errorf("%v is not a int64 value", values[0])
		}

		// class id
		classID, isInt := values[1].(int64)
		if !isInt {
			return fmt.Errorf("%v is not a int64 value", values[1])
		}

		eClass, _ := d.sqlIDManager.getClassFromID(classID)
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
		d.sqlIDManager.setObjectID(eObject, objectID)

		// set its id
		if len(values) > 2 {
			switch v := values[2].(type) {
			case nil:
			case string:
				if d.idManager != nil {
					if err := d.idManager.SetID(eObject, v); err != nil {
						return err
					}
				}
			default:
				return fmt.Errorf("%v is not a string value", values[2])
			}
		}

		return nil
	})
}

func (d *SQLDecoder) decodeFeatures() error {
	decoded := map[EClass]struct{}{}
	// for each leaf class
	for _, classData := range d.classDataMap {
		eClass := classData.eClass
		itSuper := eClass.GetEAllSuperTypes().Iterator()
		for eClass != nil {
			// decode class features
			if _, idDecoded := decoded[eClass]; !idDecoded {
				decoded[eClass] = struct{}{}
				if err := d.decodeClassFeatures(eClass); err != nil {
					return err
				}
			}

			// next super class
			if itSuper.HasNext() {
				eClass = itSuper.Next().(EClass)
			} else {
				eClass = nil
			}
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
	return d.query(table.selectQuery(nil, "", ""), func(values []driver.Value) error {
		// object id
		objectID, isInt := values[0].(int64)
		if !isInt {
			return fmt.Errorf("%v is not a int64 value", values[0])
		}

		// retrieve EObject
		eObject, _ := d.sqlIDManager.getObjectFromID(objectID)
		if eObject == nil {
			return fmt.Errorf("unable to find object with id '%v'", objectID)
		}

		// for each column
		for i, columnData := range columnFeatures {
			columnValue, err := d.decodeFeatureValue(columnData, values[i+1])
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
	})

}

func (d *SQLDecoder) decodeTableFeature(table *sqlTable, tableFeature *sqlFeatureSchema) error {
	column := sqlEscapeIdentifier(table.columns[len(table.columns)-1].columnName)
	query := table.selectQuery([]string{table.keyName(), column}, "", table.keyName()+" ASC, idx ASC")

	feature := tableFeature.feature
	currentValues := []any{}
	var currentID int64 = -1

	if err := d.query(query, func(values []driver.Value) error {
		// object id
		objectID, isInt := values[0].(int64)
		if !isInt {
			return fmt.Errorf("%v is not a int64 value", values[0])
		}

		value, err := d.decodeFeatureValue(tableFeature, values[1])
		if err != nil {
			return err
		}

		if currentID == -1 {
			currentID = objectID
		} else if currentID != objectID {
			if err := d.decodeFeatureList(currentID, feature, currentValues); err != nil {
				return err
			}
			currentValues = nil
		}

		currentID = objectID
		currentValues = append(currentValues, value)

		return nil
	}); err != nil {
		return err
	}

	if currentID != -1 {
		if err := d.decodeFeatureList(currentID, feature, currentValues); err != nil {
			return err
		}
	}

	return nil
}

func (d *SQLDecoder) decodeFeatureList(objectID int64, feature EStructuralFeature, values []any) error {
	if len(values) == 0 {
		return nil
	}
	eObject, _ := d.sqlIDManager.getObjectFromID(objectID)
	if eObject == nil {
		return fmt.Errorf("unable to find object with id '%v'", objectID)
	}
	eList := eObject.EGetResolve(feature, false).(EList)
	eList.AddAll(NewImmutableEList(values))
	return nil
}

func (d *SQLDecoder) query(q string, cb func(values []driver.Value) error, args ...driver.Value) error {
	con, err := d.db.Conn(context.Background())
	if err != nil {
		return err
	}
	defer con.Close()

	return con.Raw(func(driverConn any) error {
		//lint:ignore SA1019 driver.Queryer has been deprecated since Go 1.8: Drivers should implement QueryerContext instead
		driverQuery, _ := driverConn.(driver.Queryer)
		if driverQuery == nil {
			return errors.New("driver is not a driver.Queryer")
		}

		rows, err := driverQuery.Query(q, args)
		if err != nil {
			return err
		}
		defer rows.Close()

		results := make([]driver.Value, len(rows.Columns()))
		for {
			// retrieve results
			if err := rows.Next(results); err != nil {
				if err == io.EOF {
					return nil
				} else {
					return err
				}
			}
			// and call cb function
			if err := cb(results); err != nil {
				return err
			}
		}
	})
}

func (d *SQLDecoder) addError(err error) {
	d.resource.GetErrors().Add(NewEDiagnosticImpl(err.Error(), d.resource.GetURI().String(), 0, 0))
}
