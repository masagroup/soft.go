package ecore

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type sqlDecoderClassData struct {
	schema         *sqlClassSchema
	eClass         EClass
	eFactory       EFactory
	hierarchy      []EClass
	columnFeatures []*sqlDecoderFeatureData
	tableFeatures  []*sqlDecoderFeatureData
}

type sqlDecoderFeatureData struct {
	schema   *sqlFeatureSchema
	eFeature EStructuralFeature
	eFactory EFactory
	eType    EClassifier
}

type SQLDecoder struct {
	resource        EResource
	reader          io.Reader
	driver          string
	db              *sql.DB
	schema          *sqlSchema
	packages        map[int]EPackage
	objects         map[int]EObject
	classes         map[int]*sqlDecoderClassData
	selectStmts     map[*sqlTable]*sql.Stmt
	idAttributeName string
	baseURI         *URI
}

func NewSQLDecoder(resource EResource, r io.Reader, options map[string]any) *SQLDecoder {
	// options
	schemaOptions := []sqlSchemaOption{}
	driver := "sqlite"
	idAttributeName := ""
	if options != nil {
		if driver, isDriver := options[SQL_OPTION_DRIVER]; isDriver {
			driver = driver.(string)
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
		packages:        map[int]EPackage{},
		objects:         map[int]EObject{},
		classes:         map[int]*sqlDecoderClassData{},
		selectStmts:     map[*sqlTable]*sql.Stmt{},
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

	// if err := d.decodeClasses(); err != nil {
	// 	d.addError(err)
	// 	return
	// }

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
	rows, err := d.query(d.schema.contentsTable, d.selectAllQuery)
	if err != nil {
		return err
	}
	defer rows.Close()

	return d.forEachRow(rows, func(rb []sql.RawBytes) error {
		// objectID, err := strconv.Atoi(string(rb[0]))
		// if err != nil {
		// 	return err
		// }

		// // decode object
		// eObject, err := d.decodeObject(objectID)
		// if err != nil {
		// 	return err
		// }

		// // add object to resource contents
		// d.resource.GetContents().Add(eObject)
		return nil
	})
}

// func (d *SQLDecoder) decodeObject(objectID int) (EObject, error) {
// 	eObject := d.objects[objectID]
// 	if eObject == nil {
// 		eClass, uniqueID, err := d.decodeObjectClassAndID(objectID)
// 		if err != nil {
// 			return nil, err
// 		}

// 		classData := d.classToData[eClass]
// 		if classData == nil {
// 			return nil, fmt.Errorf("unable to find class for object '%v'", objectID)
// 		}

// 		// create object & set its unique id if any
// 		eObject = classData.eFactory.Create(classData.eClass)
// 		if uniqueID.Valid {
// 			objectIDManager := d.resource.GetObjectIDManager()
// 			objectIDManager.SetID(eObject, uniqueID.String)
// 		}

// 		// register object
// 		d.objects[objectID] = eObject

// 		// decode object feature values
// 		for _, eClass := range classData.hierarchy {
// 			classData := d.classToData[eClass]
// 			if err := d.decodeObjectColumnFeatures(objectID, eObject, classData); err != nil {
// 				return nil, err
// 			}
// 			if err := d.decodeObjectTableFeatures(objectID, eObject, classData); err != nil {
// 				return nil, err
// 			}
// 		}

// 	}
// 	return eObject, nil
// }

func (d *SQLDecoder) decodeObjectColumnFeatures(objectID int, eObject EObject, classData *sqlDecoderClassData) error {
	if len(classData.columnFeatures) > 0 {
		// stmt
		rows, err := d.query(classData.schema.table, func(table *sqlTable) string {
			return table.selectQuery(table.columnNames(1, -1), table.keyName()+"= ? ", "")
		}, objectID)
		if err != nil {
			return err
		}
		defer rows.Close()

		var rawBuffer []sql.RawBytes
		d.forOneRow(rows, func(rb []sql.RawBytes) error {
			rawBuffer = rb
			return nil
		})

		// decode feature values in this table
		// first value is objectID in rawBuffer so we skip it
		for i, featureData := range classData.columnFeatures {
			columnValue, err := d.decodeFeatureValue(featureData, rawBuffer[i])
			if err != nil {
				return err
			}
			eObject.ESet(featureData.eFeature, columnValue)
		}
	}
	return nil
}

func (d *SQLDecoder) decodeObjectTableFeatures(objectID int, eObject EObject, classData *sqlDecoderClassData) error {
	for _, featureData := range classData.tableFeatures {
		// decode each list values
		values := []any{}

		// query
		rows, err := d.query(featureData.schema.table, func(table *sqlTable) string {
			// column value is the last one
			valueColumn := table.columns[len(table.columns)-1]
			return table.selectQuery([]string{valueColumn.columnName}, table.keyName()+"= ?", "idx ASC")
		}, objectID)
		if err != nil {
			return err
		}
		defer rows.Close()

		// for each row, decode feature value and add value to the list
		d.forEachRow(rows, func(rb []sql.RawBytes) error {
			value, err := d.decodeFeatureValue(featureData, rb[0])
			if err != nil {
				return err
			}
			values = append(values, value)
			return nil
		})
		// set list element values
		list := eObject.EGetResolve(featureData.eFeature, false).(EList)
		list.AddAll(NewImmutableEList(values))
	}

	return nil

}

func (d *SQLDecoder) decodeFeatureValue(featureData *sqlDecoderFeatureData, bytes []byte) (any, error) {
	switch featureData.schema.featureKind {
	case sfkObject, sfkObjectList:
		if len(bytes) == 0 {
			return nil, nil
		}
		objectID, err := strconv.Atoi(string(bytes))
		if err != nil {
			return nil, err
		}
		//return d.decodeObject(objectID)
		return d.objects[objectID], nil
	case sfkObjectReference, sfkObjectReferenceList:
		if len(bytes) == 0 {
			return nil, nil
		}
		// uri
		uriStr := string(bytes)
		uri := d.baseURI
		if len(uriStr) > 0 {
			uri = d.resolveURI(NewURI(uriStr))
		}
		// create proxy
		eObject := featureData.eFactory.Create(featureData.eType.(EClass))
		eObjectInternal := eObject.(EObjectInternal)
		eObjectInternal.ESetProxyURI(uri)
		return eObject, nil
	case sfkBool:
		return strconv.ParseBool(string(bytes))
	case sfkByte:
		if len(bytes) == 0 {
			var defaultByte byte
			return defaultByte, errors.New("invalid bytes length")
		}
		return bytes[0], nil
	case sfkInt:
		i, err := strconv.ParseInt(string(bytes), 10, 0)
		if err != nil {
			var defaultInt int
			return defaultInt, err
		}
		return int(i), nil
	case sfkInt64:
		return strconv.ParseInt(string(bytes), 10, 64)
	case sfkInt32:
		i, err := strconv.ParseInt(string(bytes), 10, 32)
		if err != nil {
			var defaultInt32 int32
			return defaultInt32, err
		}
		return int32(i), nil
	case sfkInt16:
		i, err := strconv.ParseInt(string(bytes), 10, 16)
		if err != nil {
			var defaultInt16 int16
			return defaultInt16, err
		}
		return int16(i), nil
	case sfkEnum:
		return strconv.ParseInt(string(bytes), 10, 64)
	case sfkString:
		return string(bytes), nil
	case sfkByteArray:
		return bytes, nil
	case sfkDate:
		t, err := time.Parse(time.RFC3339, string(bytes))
		if err != nil {
			return nil, err
		}
		return &t, nil
	case sfkFloat64:
		return strconv.ParseFloat(string(bytes), 64)
	case sfkFloat32:
		f, err := strconv.ParseFloat(string(bytes), 32)
		if err != nil {
			var defaultFloat32 float32
			return defaultFloat32, err
		}
		return f, nil
	case sfkData, sfkDataList:
		return featureData.eFactory.CreateFromString(featureData.eType.(EDataType), string(bytes)), nil
	}

	return nil, nil
}

func (d *SQLDecoder) resolveURI(uri *URI) *URI {
	if d.baseURI != nil {
		return d.baseURI.Resolve(uri)
	}
	return uri
}

// func (d *SQLDecoder) decodeObjectClassAndID(objectID int) (EClass, sql.NullString, error) {
// 	// retrieve class id and unique id for this object
// 	var uniqueID sql.NullString
// 	var classID int

// 	// query
// 	objectsTable := d.schema.objectsTable
// 	columns := objectsTable.columnNames(1, -1)
// 	rows, err := d.query(objectsTable, func(table *sqlTable) string {
// 		return table.selectQuery(columns, table.keyName()+" = ?", "")
// 	}, objectID)
// 	if err != nil {
// 		return nil, uniqueID, err
// 	}
// 	defer rows.Close()

// 	// only one row
// 	d.forOneRow(rows, func(rb []sql.RawBytes) error {
// 		// classID
// 		id, err := strconv.Atoi(string(rb[0]))
// 		if err != nil {
// 			return err
// 		}
// 		classID = id
// 		// uniqueID
// 		if len(rb) > 1 {
// 			if err := uniqueID.Scan(string(rb[1])); err != nil {
// 				return err
// 			}
// 		}
// 		return nil
// 	})

// 	// retrieve class data
// 	eClass := d.classes[classID]
// 	if eClass == nil {
// 		return nil, uniqueID, fmt.Errorf("unable to find class with id '%v'", classID)
// 	}
// 	return eClass, uniqueID, nil
// }

// func (d *SQLDecoder) decodeClasses() error {
// 	// read packages
// 	packagesData, err := d.decodePackages()
// 	if err != nil {
// 		return err
// 	}

// 	// read classes
// 	classes := map[int]EClass{}
// 	classToData := map[EClass]*sqlDecoderClassData{}

// 	// query
// 	rows, err := d.query(d.schema.classesTable, d.selectAllQuery)
// 	if err != nil {
// 		return err
// 	}
// 	defer rows.Close()

// 	// for each row
// 	if err := d.forEachRow(rows, func(rb []sql.RawBytes) error {
// 		classID, err := strconv.Atoi(string(rb[0]))
// 		if err != nil {
// 			return err
// 		}
// 		packageID, err := strconv.Atoi(string(rb[1]))
// 		if err != nil {
// 			return err
// 		}
// 		ePackage := packagesData[packageID]
// 		if ePackage == nil {
// 			return fmt.Errorf("unable to find package with id '%d'", packageID)
// 		}
// 		className := string(rb[2])
// 		eClass, _ := ePackage.GetEClassifier(className).(EClass)
// 		if eClass == nil {
// 			return fmt.Errorf("unable to find class '%s' in package '%s'", className, ePackage.GetNsURI())
// 		}

// 		// init map
// 		classes[classID] = eClass

// 		// get class schema
// 		classSchema, err := d.schema.getClassSchema(eClass)
// 		if err != nil {
// 			return err
// 		}

// 		// compute class hierarchy
// 		classHierarchy := []EClass{eClass}
// 		for itClass := eClass.GetEAllSuperTypes().Iterator(); itClass.HasNext(); {
// 			classHierarchy = append(classHierarchy, itClass.Next().(EClass))
// 		}

// 		// compute class features
// 		classColumnFeatures := make([]*sqlDecoderFeatureData, 0, len(classSchema.features))
// 		classTableFeatures := make([]*sqlDecoderFeatureData, 0, len(classSchema.features))
// 		for _, featureSchema := range classSchema.features {
// 			eFeature := featureSchema.feature
// 			eFeatureData := &sqlDecoderFeatureData{
// 				eFeature: eFeature,
// 				schema:   featureSchema,
// 			}
// 			if eAttribute, _ := eFeature.(EAttribute); eAttribute != nil {
// 				eFeatureData.eType = eAttribute.GetEAttributeType()
// 				eFeatureData.eFactory = eFeatureData.eType.GetEPackage().GetEFactoryInstance()
// 			} else if eReference := eFeature.(EReference); eReference != nil {
// 				eFeatureData.eType = eReference.GetEReferenceType()
// 				eFeatureData.eFactory = eFeatureData.eType.GetEPackage().GetEFactoryInstance()
// 			}
// 			if featureSchema.column != nil {
// 				classColumnFeatures = append(classColumnFeatures, eFeatureData)
// 			} else if featureSchema.table != nil {
// 				classTableFeatures = append(classTableFeatures, eFeatureData)
// 			}

// 		}

// 		// register class data
// 		classToData[eClass] = &sqlDecoderClassData{
// 			eClass:         eClass,
// 			eFactory:       ePackage.GetEFactoryInstance(),
// 			schema:         classSchema,
// 			hierarchy:      classHierarchy,
// 			columnFeatures: classColumnFeatures,
// 			tableFeatures:  classTableFeatures,
// 		}
// 		return nil
// 	}); err != nil {
// 		return err
// 	}

// 	d.classes = classes
// 	d.classToData = classToData
// 	return nil
// }

func (d *SQLDecoder) decodeObjects() error {
	rows, err := d.query(d.schema.objectsTable, d.selectAllQuery)
	if err != nil {
		return err
	}
	defer rows.Close()

	// for each row, create object and retrieve used classes
	return d.forEachRow(rows, func(rb []sql.RawBytes) error {
		// object id
		objectID, err := strconv.Atoi(string(rb[0]))
		if err != nil {
			return err
		}

		// class id
		classID, err := strconv.Atoi(string(rb[1]))
		if err != nil {
			return err
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
		if len(rb) > 2 {
			uniqueID := string(rb[1])
			objectIDManager := d.resource.GetObjectIDManager()
			objectIDManager.SetID(eObject, uniqueID)
		}

		return nil
	})

}

func (d *SQLDecoder) decodePackages() error {
	// query
	rows, err := d.query(d.schema.packagesTable, d.selectAllQuery)
	if err != nil {
		return err
	}
	defer rows.Close()

	// for each row, retrieve package from registry
	return d.forEachRow(rows, func(rb []sql.RawBytes) error {
		packageID, _ := strconv.Atoi(string(rb[0]))
		packageURI := string(rb[1])
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

func (d *SQLDecoder) getSelectStmt(table *sqlTable, queryProvider func(table *sqlTable) string) (stmt *sql.Stmt, err error) {
	stmt = d.selectStmts[table]
	if stmt == nil {
		query := queryProvider(table)
		stmt, err = d.db.Prepare(query)
	}
	return
}

func (d *SQLDecoder) query(table *sqlTable, queryProvider func(table *sqlTable) string, args ...any) (*sql.Rows, error) {
	selectStmt, err := d.getSelectStmt(table, queryProvider)
	if err != nil {
		return nil, err
	}
	return selectStmt.Query(args...)
}

func (d *SQLDecoder) selectAllQuery(table *sqlTable) string {
	return table.selectQuery(nil, "", "")
}

func (d *SQLDecoder) addError(err error) {
	d.resource.GetErrors().Add(NewEDiagnosticImpl(err.Error(), d.resource.GetURI().String(), 0, 0))
}

func (d *SQLDecoder) forOneRow(rows *sql.Rows, cb func([]sql.RawBytes) error) error {
	// one row
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return err
		}
		return sql.ErrNoRows
	}

	//scan raw buffer
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	rawBuffer := make([]sql.RawBytes, len(columns))
	scanCallArgs := make([]any, len(rawBuffer))
	for i := range rawBuffer {
		scanCallArgs[i] = &rawBuffer[i]
	}
	if err := rows.Scan(scanCallArgs...); err != nil {
		return err
	}
	// callback
	return cb(rawBuffer)
}

func (d *SQLDecoder) forEachRow(rows *sql.Rows, cb func([]sql.RawBytes) error) error {
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	rawBuffer := make([]sql.RawBytes, len(columns))
	scanCallArgs := make([]any, len(columns))
	for i := range rawBuffer {
		scanCallArgs[i] = &rawBuffer[i]
	}

	for rows.Next() {
		if err := rows.Scan(scanCallArgs...); err != nil {
			return err
		}

		if err := cb(rawBuffer); err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}
