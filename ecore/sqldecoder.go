package ecore

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
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
	objects         map[int]EObject
	classes         map[int]EClass
	classToData     map[EClass]*sqlDecoderClassData
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
		objects:         map[int]EObject{},
		classes:         map[int]EClass{},
		classToData:     map[EClass]*sqlDecoderClassData{},
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

	if err := d.decodeClasses(); err != nil {
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

func (d *SQLDecoder) addError(err error) {
	d.resource.GetErrors().Add(NewEDiagnosticImpl(err.Error(), d.resource.GetURI().String(), 0, 0))
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
	rows, err := d.db.Query("SELECT * FROM contents")
	if err != nil {
		return err
	}
	defer rows.Close()

	rawBuffer := make([]sql.RawBytes, 1)
	scanCallArgs := make([]any, len(rawBuffer))
	for i := range rawBuffer {
		scanCallArgs[i] = &rawBuffer[i]
	}

	for rows.Next() {
		// retrieve object id
		if err := rows.Scan(scanCallArgs...); err != nil {
			return err
		}
		objectID, _ := strconv.Atoi(string(rawBuffer[0]))

		// decode object
		eObject, err := d.decodeObject(objectID)
		if err != nil {
			return err
		}

		// add object to resource contents
		d.resource.GetContents().Add(eObject)
	}
	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}

func (d *SQLDecoder) decodeObject(objectID int) (EObject, error) {
	eObject := d.objects[objectID]
	if eObject == nil {
		eClass, uniqueID, err := d.decodeObjectClassAndID(objectID)
		if err != nil {
			return nil, err
		}

		classData := d.classToData[eClass]
		if classData == nil {
			return nil, fmt.Errorf("unable to find class for object '%v'", objectID)
		}

		// create object & set its unique id if any
		eObject := classData.eFactory.Create(classData.eClass)
		if uniqueID.Valid {
			objectIDManager := d.resource.GetObjectIDManager()
			objectIDManager.SetID(eObject, uniqueID.String)
		}

		// register object
		d.objects[objectID] = eObject

		// decode object feature values
		for _, eClass := range classData.hierarchy {
			classData := d.classToData[eClass]
			if err := d.decodeObjectColumnFeatures(objectID, eObject, classData); err != nil {
				return nil, err
			}
			if err := d.decodeObjectTableFeatures(objectID, eObject, classData); err != nil {
				return nil, err
			}
		}

	}
	return eObject, nil
}

func (d *SQLDecoder) decodeObjectColumnFeatures(objectID int, eObject EObject, classData *sqlDecoderClassData) error {
	// stmt
	selectStmt, err := d.getSelectStmt(classData.schema.table)
	if err != nil {
		return err
	}

	// query
	rows, err := selectStmt.Query(objectID)
	if err != nil {
		return err
	}
	defer rows.Close()

	// one row
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return err
		}
		return sql.ErrNoRows
	}

	// retrieve column values
	rawBuffer := make([]sql.RawBytes, len(classData.columnFeatures)+1)
	scanCallArgs := make([]any, len(rawBuffer))
	for i := range rawBuffer {
		scanCallArgs[i] = &rawBuffer[i]
	}
	if err := rows.Scan(scanCallArgs...); err != nil {
		return err
	}

	// decode feature values in this table
	// first value is objectID in rawBuffer so we skip it
	for i, featureData := range classData.columnFeatures {
		columnValue, err := d.decodeFeatureValue(featureData, rawBuffer[i+1])
		if err != nil {
			return err
		}
		eObject.ESet(featureData.eFeature, columnValue)
	}

	return nil
}

func (d *SQLDecoder) decodeObjectTableFeatures(objectID int, eObject EObject, classData *sqlDecoderClassData) error {
	for _, featureData := range classData.tableFeatures {
		values := []any{}
		indexes := map[int]float64{}

		// stmt
		selectStmt, err := d.getSelectStmt(classData.schema.table)
		if err != nil {
			return err
		}

		// query
		rows, err := selectStmt.Query(objectID)
		if err != nil {
			return err
		}
		defer rows.Close()

		// retrieve column values
		rawBuffer := make([]sql.RawBytes, 3)
		scanCallArgs := make([]any, len(rawBuffer))
		for i := range rawBuffer {
			scanCallArgs[i] = &rawBuffer[i]
		}

		index := 0
		for rows.Next() {
			// retrieve object id
			if err := rows.Scan(scanCallArgs...); err != nil {
				return err
			}

			// index
			indexes[index], err = strconv.ParseFloat(string(rawBuffer[2]), 64)
			if err != nil {
				return err
			}

			// value
			value, err := d.decodeFeatureValue(featureData, rawBuffer[2])
			if err != nil {
				return err
			}

			values = append(values, value)

		}
		if err = rows.Err(); err != nil {
			return err
		}

		// sort values according to their indexes
		sort.Slice(values, func(i, j int) bool { return indexes[i] < indexes[2] })

		// add values to the list
		list := eObject.EGetResolve(featureData.eFeature, false).(EList)
		list.AddAll(NewImmutableEList(values))
	}

	return nil

}

func (d *SQLDecoder) decodeFeatureValue(featureData *sqlDecoderFeatureData, bytes []byte) (any, error) {
	switch featureData.schema.featureKind {
	case sfkObject, sfkObjectList:
		objectID, err := strconv.Atoi(string(bytes))
		if err != nil {
			return nil, err
		}
		return d.decodeObject(objectID)
	case sfkObjectReference, sfkObjectReferenceList:
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
		return bytes[0], nil
	case sfkInt:
		i, err := strconv.ParseInt(string(bytes), 10, 0)
		if err != nil {
			return nil, err
		}
		return int(i), nil
	case sfkInt64:
		return strconv.ParseInt(string(bytes), 10, 64)
	case sfkInt32:
		i, err := strconv.ParseInt(string(bytes), 10, 32)
		if err != nil {
			return nil, err
		}
		return int32(i), nil
	case sfkInt16:
		i, err := strconv.ParseInt(string(bytes), 10, 16)
		if err != nil {
			return nil, err
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

func (d *SQLDecoder) decodeObjectClassAndID(objectID int) (EClass, sql.NullString, error) {
	// retrieve class id for this object
	var uniqueID sql.NullString
	selectObjectStmt, err := d.getSelectStmt(d.schema.objectsTable)
	if err != nil {
		return nil, uniqueID, err
	}

	rows, err := selectObjectStmt.Query(objectID)
	if err != nil {
		return nil, uniqueID, err
	}
	defer rows.Close()

	// one row
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, uniqueID, err
		}
		return nil, uniqueID, sql.ErrNoRows
	}

	// scan first row
	rawBufferSize := 2
	if d.isObjectWithUniqueID() {
		rawBufferSize++
	}
	rawBuffer := make([]sql.RawBytes, rawBufferSize)
	scanCallArgs := make([]any, len(rawBuffer))
	for i := range rawBuffer {
		scanCallArgs[i] = &rawBuffer[i]
	}
	if err := rows.Scan(scanCallArgs...); err != nil {
		return nil, uniqueID, err
	}

	// extract row args
	classID, _ := strconv.Atoi(string(rawBuffer[1]))
	if d.isObjectWithUniqueID() {
		uniqueID.Scan(string(rawBuffer[2]))
	}

	// retrieve class data
	eClass := d.classes[classID]
	if eClass == nil {
		return nil, uniqueID, fmt.Errorf("unable to find class with id '%v'", classID)
	}
	return eClass, uniqueID, nil
}

func (d *SQLDecoder) decodeClasses() error {
	// read packages
	packagesData, err := d.decodePackages()
	if err != nil {
		return err
	}

	// read classes
	rows, err := d.db.Query(d.schema.classesTable.selectAllQuery())
	if err != nil {
		return err
	}
	defer rows.Close()

	classes := map[int]EClass{}
	classToData := map[EClass]*sqlDecoderClassData{}
	rawBuffer := make([]sql.RawBytes, 3)
	scanCallArgs := make([]any, len(rawBuffer))
	for i := range rawBuffer {
		scanCallArgs[i] = &rawBuffer[i]
	}

	for rows.Next() {
		// retrieve EClass
		if err := rows.Scan(scanCallArgs...); err != nil {
			return err
		}
		classID, _ := strconv.Atoi(string(rawBuffer[0]))
		packageID, _ := strconv.Atoi(string(rawBuffer[1]))
		className := string(rawBuffer[2])
		ePackage := packagesData[packageID]
		if ePackage == nil {
			return fmt.Errorf("unable to find package with id '%d'", packageID)
		}
		eClass, _ := ePackage.GetEClassifier(className).(EClass)
		if eClass == nil {
			return fmt.Errorf("unable to find class '%s' in package '%s'", className, ePackage.GetNsURI())
		}

		// init map
		classes[classID] = eClass

		// get class schema
		classSchema, err := d.schema.getClassSchema(eClass)
		if err != nil {
			return err
		}

		// compute class hierarchy
		classHierarchy := []EClass{eClass}
		for itClass := eClass.GetEAllSuperTypes().Iterator(); itClass.HasNext(); {
			classHierarchy = append(classHierarchy, itClass.Next().(EClass))
		}

		// compute class features
		classColumnFeatures := make([]*sqlDecoderFeatureData, 0, len(classSchema.features))
		classTableFeatures := make([]*sqlDecoderFeatureData, 0, len(classSchema.features))
		for eFeature, featureSchema := range classSchema.features {
			eFeatureData := &sqlDecoderFeatureData{
				eFeature: eFeature,
				schema:   featureSchema,
			}
			if eAttribute, _ := eFeature.(EAttribute); eAttribute != nil {
				eFeatureData.eType = eAttribute.GetEAttributeType()
				eFeatureData.eFactory = eFeatureData.eType.GetEPackage().GetEFactoryInstance()
			} else if eReference := eFeature.(EReference); eReference != nil {
				eFeatureData.eType = eReference.GetEReferenceType()
				eFeatureData.eFactory = eFeatureData.eType.GetEPackage().GetEFactoryInstance()
			}
			if featureSchema.column != nil {
				classColumnFeatures = append(classColumnFeatures, eFeatureData)
			} else if featureSchema.table != nil {
				classTableFeatures = append(classTableFeatures, eFeatureData)
			}

		}

		// register class data
		classToData[eClass] = &sqlDecoderClassData{
			eClass:         eClass,
			eFactory:       ePackage.GetEFactoryInstance(),
			schema:         classSchema,
			hierarchy:      classHierarchy,
			columnFeatures: classColumnFeatures,
			tableFeatures:  classTableFeatures,
		}

	}
	if err = rows.Err(); err != nil {
		return err
	}

	d.classes = classes
	d.classToData = classToData
	return nil
}

func (d *SQLDecoder) decodePackages() (map[int]EPackage, error) {
	rows, err := d.db.Query(d.schema.packagesTable.selectAllQuery())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	packagesData := map[int]EPackage{}
	rawBuffer := make([]sql.RawBytes, 2)
	scanCallArgs := make([]any, len(rawBuffer))
	for i := range rawBuffer {
		scanCallArgs[i] = &rawBuffer[i]
	}

	for rows.Next() {
		if err := rows.Scan(scanCallArgs...); err != nil {
			return nil, err
		}
		packageID, _ := strconv.Atoi(string(rawBuffer[0]))
		packageURI := string(rawBuffer[1])
		packageRegistry := GetPackageRegistry()
		resourceSet := d.resource.GetResourceSet()
		if resourceSet != nil {
			packageRegistry = resourceSet.GetPackageRegistry()
		}
		ePackage := packageRegistry.GetPackage(packageURI)
		if ePackage == nil {
			return nil, fmt.Errorf("unable to find package '%s'", packageURI)
		}
		packagesData[packageID] = ePackage
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return packagesData, nil
}

func (d *SQLDecoder) getSelectStmt(table *sqlTable) (stmt *sql.Stmt, err error) {
	stmt = d.selectStmts[table]
	if stmt == nil {
		stmt, err = d.db.Prepare(table.selectWhereQuery())
	}
	return
}

func (d *SQLDecoder) isObjectWithUniqueID() bool {
	return d.resource.GetObjectIDManager() != nil && len(d.idAttributeName) > 0
}
