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
	"time"
)

type sqlDecoderClassData struct {
	eClass   EClass
	eFactory EFactory
}

type SQLDecoder struct {
	resource        EResource
	reader          io.Reader
	driver          string
	db              *sql.DB
	schema          *sqlSchema
	packages        map[int64]EPackage
	objects         map[int64]EObject
	classes         map[int64]*sqlDecoderClassData
	enums           map[int64]any
	idAttributeName string
	baseURI         *URI
}

func NewSQLDecoder(resource EResource, r io.Reader, options map[string]any) *SQLDecoder {
	// options
	schemaOptions := []sqlSchemaOption{}
	driver := "sqlite"
	idAttributeName := ""
	if options != nil {
		if d, isDriver := options[SQL_OPTION_DRIVER]; isDriver {
			driver = d.(string)
		}

		idAttributeName, _ = options[JSON_OPTION_ID_ATTRIBUTE_NAME].(string)
		if resource.GetObjectIDManager() != nil && len(idAttributeName) > 0 {
			schemaOptions = append(schemaOptions, withIDAttributeName(idAttributeName))
		}
	}
	var baseURI *URI
	if uri := resource.GetURI(); uri != nil {
		baseURI = uri
	}

	return &SQLDecoder{
		resource:        resource,
		reader:          r,
		driver:          driver,
		schema:          newSqlSchema(schemaOptions...),
		packages:        map[int64]EPackage{},
		objects:         map[int64]EObject{},
		classes:         map[int64]*sqlDecoderClassData{},
		enums:           map[int64]any{},
		idAttributeName: idAttributeName,
		baseURI:         baseURI,
	}
}

func (d *SQLDecoder) createDB() (*sql.DB, error) {
	fileName := filepath.Base(d.resource.GetURI().Path())
	dbPath, err := sqlTmpDB(fileName)
	if err != nil {
		return nil, err
	}

	dbFile, err := os.Create(dbPath)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(dbFile, d.reader)
	if err != nil {
		dbFile.Close()
		return nil, err
	}
	dbFile.Close()

	return sql.Open(d.driver, dbPath)
}

func (d *SQLDecoder) DecodeResource() {
	var err error
	d.db, err = d.createDB()
	if err != nil {
		d.addError(err)
		return
	}

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
		eObject := d.objects[objectID]
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

		packageRegistry := GetPackageRegistry()
		resourceSet := d.resource.GetResourceSet()
		if resourceSet != nil {
			packageRegistry = resourceSet.GetPackageRegistry()
		}
		ePackage := packageRegistry.GetPackage(packageURI)
		if ePackage == nil {
			return fmt.Errorf("unable to find package '%s'", packageURI)
		}
		d.packages[packageID] = ePackage
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
		ePackage := d.packages[packageID]
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

		// create map enum
		d.enums[enumID] = ePackage.GetEFactoryInstance().CreateFromString(eEnum, literalValue)
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
		ePackage := d.packages[packageID]
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

		// compute class hierarchy
		classHierarchy := []EClass{eClass}
		for itClass := eClass.GetEAllSuperTypes().Iterator(); itClass.HasNext(); {
			classHierarchy = append(classHierarchy, itClass.Next().(EClass))
		}

		d.classes[classID] = &sqlDecoderClassData{
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

		// class data
		classData := d.classes[classID]
		if classData == nil {
			return fmt.Errorf("unable to find class with id '%v'", classID)
		}

		// create & register object
		eObject := classData.eFactory.Create(classData.eClass)
		d.objects[objectID] = eObject

		// set its id
		if len(values) > 2 {
			uniqueID, isString := values[2].(string)
			if !isString {
				return fmt.Errorf("%v is not a string value", values[2])
			}
			objectIDManager := d.resource.GetObjectIDManager()
			if err := objectIDManager.SetID(eObject, uniqueID); err != nil {
				return err
			}
		}

		return nil
	})
}

func (d *SQLDecoder) decodeFeatures() error {
	decoded := map[EClass]struct{}{}
	// for each leaf class
	for _, classData := range d.classes {
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
	classSchema, err := d.schema.getClassSchema(eClass)
	if err != nil {
		return err
	}

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
		eObject := d.objects[objectID]
		if eObject == nil {
			return fmt.Errorf("unable to find object with id '%v'", objectID)
		}

		// for each column
		for i, columnData := range columnFeatures {
			columnValue, err := d.decodeFeatureValue(columnData, values[i+1])
			if err != nil {
				return err
			}
			eObject.ESet(columnData.feature, columnValue)
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
	eObject := d.objects[objectID]
	if eObject == nil {
		return fmt.Errorf("unable to find object with id '%v'", objectID)
	}
	eList := eObject.EGetResolve(feature, false).(EList)
	eList.AddAll(NewImmutableEList(values))
	return nil
}

func (d *SQLDecoder) decodeFeatureValue(featureData *sqlFeatureSchema, value driver.Value) (any, error) {
	switch featureData.featureKind {
	case sfkObject, sfkObjectList:
		switch v := value.(type) {
		case nil:
			return nil, nil
		case int64:
			return d.objects[v], nil
		default:
			return nil, fmt.Errorf("%v is not supported as a object id", v)
		}
	case sfkObjectReference, sfkObjectReferenceList:
		switch v := value.(type) {
		case string:
			// no reference
			if len(v) == 0 {
				return nil, nil
			}
			// resolve uri reference
			uri := d.resolveURI(NewURI(v))

			// create proxy
			eFeature := featureData.feature
			eClass := eFeature.GetEType().(EClass)
			eFactory := eClass.GetEPackage().GetEFactoryInstance()
			eObject := eFactory.Create(eClass)
			eObjectInternal := eObject.(EObjectInternal)
			eObjectInternal.ESetProxyURI(uri)
			return eObject, nil
		default:
			return nil, fmt.Errorf("%v is not supported as a object reference uri", v)
		}
	case sfkBool:
		switch v := value.(type) {
		case bool:
			return v, nil
		default:
			var defaultBool bool
			return defaultBool, fmt.Errorf("%v is not a bool value", v)
		}
	case sfkByte:
		// TODO
		var defaultByte byte
		return defaultByte, nil
	case sfkInt:
		switch v := value.(type) {
		case int64:
			return int(v), nil
		default:
			var defaultInt int
			return defaultInt, fmt.Errorf("%v is not a int value", v)
		}
	case sfkInt64:
		switch v := value.(type) {
		case int64:
			return v, nil
		default:
			var defaultInt int64
			return defaultInt, fmt.Errorf("%v is not a int64 value", v)
		}
	case sfkInt32:
		switch v := value.(type) {
		case int64:
			return int32(v), nil
		default:
			var defaultInt int32
			return defaultInt, fmt.Errorf("%v is not a int32 value", v)
		}
	case sfkInt16:
		switch v := value.(type) {
		case int64:
			return int16(v), nil
		default:
			var defaultInt int16
			return defaultInt, fmt.Errorf("%v is not a int16 value", v)
		}
	case sfkEnum:
		enumID, isInt := value.(int64)
		if !isInt {
			return nil, fmt.Errorf("%v is not a int64 value", value)
		}
		return d.enums[enumID], nil
	case sfkString:
		switch v := value.(type) {
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
			return "", fmt.Errorf("%v is not a time value", v)
		}
	case sfkFloat64:
		switch v := value.(type) {
		case float64:
			return v, nil
		default:
			var defaultFloat float64
			return defaultFloat, fmt.Errorf("%v is not a float64 value", value)
		}
	case sfkFloat32:
		switch v := value.(type) {
		case float64:
			return float32(v), nil
		default:
			var defaultFloat float32
			return defaultFloat, fmt.Errorf("%v is not a float64 value", value)
		}
	case sfkData, sfkDataList:
		switch v := value.(type) {
		case string:
			eFeature := featureData.feature
			eDataType := eFeature.GetEType().(EDataType)
			eFactory := eDataType.GetEPackage().GetEFactoryInstance()
			return eFactory.CreateFromString(eDataType, v), nil
		default:
			return "", fmt.Errorf("%v is not a data value", value)
		}
	}

	return nil, nil
}

func (d *SQLDecoder) resolveURI(uri *URI) *URI {
	if d.baseURI != nil {
		return d.baseURI.Resolve(uri)
	}
	return uri
}

func (d *SQLDecoder) query(q string, cb func(values []driver.Value) error, args ...driver.Value) error {
	con, err := d.db.Conn(context.Background())
	if err != nil {
		return err
	}
	defer con.Close()

	return con.Raw(func(driverConn any) error {

		driverQuery, _ := driverConn.(driver.Queryer)
		if driverQuery == nil {
			return errors.New("driver is not a driver.Queryer")
		}

		rows, err := driverQuery.Query(q, args)
		if err != nil {
			return err
		}

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
