package ecore

import (
	"database/sql"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type stringBuilder strings.Builder

func (b *stringBuilder) String() string {
	return (*strings.Builder)(b).String()
}

func (b *stringBuilder) WriteString(s string) {
	(*strings.Builder)(b).WriteString(s)
}

func (b *stringBuilder) Grow(n int) {
	(*strings.Builder)(b).Grow(n)
}

func (b *stringBuilder) WriteArray(a []string, sep string) {
	switch len(a) {
	case 0:
	case 1:
		b.WriteString(a[0])
	}
	n := len(sep) * (len(a) - 1)
	for i := 0; i < len(a); i++ {
		n += len(a[i])
	}

	b.Grow(n)
	b.WriteString(a[0])
	for _, s := range a[1:] {
		b.WriteString(sep)
		b.WriteString(s)
	}
}

type sqlEncoderFeatureTableData struct {
	tableName string
	keyName   string
}

type sqlEncoderFeatureColumnData struct {
	ndx    int
	encode func(any) (any, error)
}

type sqlEncoderObjectData struct {
	id        int64
	classData *sqlEncoderClassData
}

type sqlEncoderClassData struct {
	id          int64
	tableName   string
	keyName     string
	tableQuery  string
	insertQuery string
	insertStmt  *sql.Stmt
	columnData  map[EStructuralFeature]*sqlEncoderFeatureColumnData
	tableData   map[EStructuralFeature]*sqlEncoderFeatureTableData
}

type sqlEncoderPackageData struct {
	id int64
}

type SQLEncoder struct {
	resource          EResource
	writer            io.Writer
	driver            string
	db                *sql.DB
	packageDataMap    map[EPackage]*sqlEncoderPackageData
	classDataMap      map[EClass]*sqlEncoderClassData
	insertContentStmt *sql.Stmt
	insertPackageStmt *sql.Stmt
	insertClassStmt   *sql.Stmt
	insertObjectStmt  *sql.Stmt
	idAttributeName   string
}

func NewSQLEncoder(resource EResource, w io.Writer, options map[string]any) *SQLEncoder {
	e := &SQLEncoder{
		resource:       resource,
		writer:         w,
		driver:         "sqlite",
		packageDataMap: map[EPackage]*sqlEncoderPackageData{},
		classDataMap:   map[EClass]*sqlEncoderClassData{},
	}
	if options != nil {
		if driver, isDriver := options[SQL_OPTION_DRIVER]; isDriver {
			e.driver = driver.(string)
		}
		e.idAttributeName, _ = options[JSON_OPTION_ID_ATTRIBUTE_NAME].(string)
	}
	return e
}

func (e *SQLEncoder) createDB() (*sql.DB, error) {
	fileName := filepath.Base(e.resource.GetURI().Path())
	dbPath, err := sqlTmpDB(fileName)
	if err != nil {
		return nil, err
	}

	// open db
	db, err := sql.Open(e.driver, dbPath)
	if err != nil {
		return nil, err
	}

	// versionQuery info
	versionQuery := fmt.Sprintf(`PRAGMA user_version = %v`, sqlCodecVersion)
	_, err = db.Exec(versionQuery)
	if err != nil {
		return nil, err
	}

	// propertiesQuery infos
	propertiesQuery := `
	PRAGMA synchronous = NORMAL;
	PRAGMA journal_mode = WAL;
	`
	_, err = db.Exec(propertiesQuery)
	if err != nil {
		return nil, err
	}

	// tables
	var tablesQuery strings.Builder
	// packages table
	tablesQuery.WriteString(`
	CREATE TABLE packages ( 
		packageID INTEGER PRIMARY KEY AUTOINCREMENT,
		uri TEXT
	);`)
	// classes table
	tablesQuery.WriteString(`
	CREATE TABLE classes (
		classID INTEGER PRIMARY KEY AUTOINCREMENT,
		packageID INTEGER,
		name TEXT,
		FOREIGN KEY(packageID) REFERENCES packages(packageID)
	);`)
	// objects table
	tablesQuery.WriteString(`
	CREATE TABLE objects (
		objectID INTEGER PRIMARY KEY AUTOINCREMENT,
		classID INTEGER,
	`)
	if e.isObjectWithID() {
		tablesQuery.WriteString(e.idAttributeName)
		tablesQuery.WriteString(` TEXT,`)
	}
	tablesQuery.WriteString(`FOREIGN KEY(classID) REFERENCES classes(classID)
	);`)
	// contents
	tablesQuery.WriteString(`
	CREATE TABLE contents (
		objectID INTEGER,
		FOREIGN KEY(objectID) REFERENCES objects(objectID)
	);
	`)
	// exec query
	_, err = db.Exec(tablesQuery.String())
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (e *SQLEncoder) EncodeResource() {
	var err error
	e.db, err = e.createDB()
	if err != nil {
		e.resource.GetErrors().Add(NewEDiagnosticImpl(err.Error(), e.resource.GetURI().String(), 0, 0))
		return
	}
	defer func() {
		_ = e.db.Close()
	}()

	if contents := e.resource.GetContents(); !contents.Empty() {
		object := contents.Get(0).(EObject)
		if err := e.encodeContent(object); err != nil {
			e.resource.GetErrors().Add(NewEDiagnosticImpl(err.Error(), e.resource.GetURI().String(), 0, 0))
			return
		}
	}
}

func (e *SQLEncoder) EncodeObject(object EObject) error {
	return nil
}

func (e *SQLEncoder) encodeContent(eObject EObject) error {
	objectData, err := e.encodeObject(eObject)
	if err != nil {
		return err
	}

	insertContentStmt, err := e.getInsertContentStmt()
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
	insertObjectStmt, err := e.getInsertObjectStmt()
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

	// encode features values
	// for featureID, featureData := range classData.featureData {
	// 	e.encodeFeatureValue(eObject, featureID, objectID, featureData)
	// }

	return &sqlEncoderObjectData{
		id:        objectID,
		classData: classData,
	}, nil
}

// func (e *SQLEncoder) encodeFeatureValue(eObject EObject, featureID int, objectID int64, featureData *sqlEncoderFeatureData) {

// }

func (e *SQLEncoder) getClassData(eClass EClass) (*sqlEncoderClassData, error) {
	classData := e.classDataMap[eClass]
	if classData == nil {
		// encode package
		ePackage := eClass.GetEPackage()
		packageData, err := e.getPackageData(ePackage)
		if err != nil {
			return nil, err
		}

		// insert class in sql
		insertClassStmt, err := e.getInsertClassStmt()
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

		// create class data
		classData = e.newClassData(eClass, id)
		e.classDataMap[eClass] = classData
	}
	return classData, nil
}

func (e *SQLEncoder) getPackageData(ePackage EPackage) (*sqlEncoderPackageData, error) {
	ePackageData := e.packageDataMap[ePackage]
	if ePackageData == nil {
		// insert new package
		insertPackageStmt, err := e.getInsertPackageStmt()
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
		ePackageData = e.newPackageData(id)
		e.packageDataMap[ePackage] = ePackageData
	}
	return ePackageData, nil
}

func (e *SQLEncoder) getInsertContentStmt() (*sql.Stmt, error) {
	if e.insertContentStmt == nil {
		insertContentStmt, err := e.db.Prepare(`INSERT INTO contents (objectID) VALUES (?)`)
		if err != nil {
			return nil, err
		}
		e.insertContentStmt = insertContentStmt
	}
	return e.insertContentStmt, nil
}

func (e *SQLEncoder) getInsertObjectStmt() (*sql.Stmt, error) {
	if e.insertObjectStmt == nil {
		columns := []string{"classID"}
		bindings := []string{"?"}
		if e.isObjectWithID() {
			columns = append(columns, e.idAttributeName)
			bindings = append(bindings, "?")
		}

		var insertObjectQuery stringBuilder
		insertObjectQuery.WriteString(`INSERT INTO objects (`)
		insertObjectQuery.WriteArray(columns, ",")
		insertObjectQuery.WriteString(`) `)
		insertObjectQuery.WriteString(`VALUES (`)
		insertObjectQuery.WriteArray(bindings, ",")
		insertObjectQuery.WriteString(`)`)

		insertObjectStmt, err := e.db.Prepare(insertObjectQuery.String())
		if err != nil {
			return nil, err
		}
		e.insertObjectStmt = insertObjectStmt
	}
	return e.insertObjectStmt, nil
}

func (e *SQLEncoder) getInsertClassStmt() (*sql.Stmt, error) {
	if e.insertClassStmt == nil {
		stmt, err := e.db.Prepare(`INSERT INTO classes (packageID,name) VALUES (?,?)`)
		if err != nil {
			return nil, err
		}
		e.insertClassStmt = stmt
	}
	return e.insertClassStmt, nil
}

func (e *SQLEncoder) getInsertPackageStmt() (*sql.Stmt, error) {
	// create statement if needed
	if e.insertPackageStmt == nil {
		stmt, err := e.db.Prepare(`INSERT INTO packages (uri) VALUES (?)`)
		if err != nil {
			return nil, err
		}
		e.insertPackageStmt = stmt
	}
	return e.insertPackageStmt, nil
}

func (e *SQLEncoder) newPackageData(id int64) *sqlEncoderPackageData {
	return &sqlEncoderPackageData{
		id: id,
	}
}

func (e *SQLEncoder) newClassData(eClass EClass, id int64) (*sqlEncoderClassData, error) {
	// create data
	ePackage := eClass.GetEPackage()
	eFeatures := eClass.GetEAllStructuralFeatures()
	classData := &sqlEncoderClassData{
		id:         id,
		tableName:  ePackage.GetNsPrefix() + "_" + strings.ToLower(eClass.GetName()),
		keyName:    strings.ToLower(eClass.GetName()) + "ID",
		columnData: map[EStructuralFeature]*sqlEncoderFeatureColumnData{},
		tableData:  map[EStructuralFeature]*sqlEncoderFeatureTableData{},
	}

	// table query
	var tableQuery strings.Builder
	tableQuery.WriteString("CREATE TABLE ")
	tableQuery.WriteString(classData.tableName)
	tableQuery.WriteString(" (")
	tableQuery.WriteString(classData.keyName)
	tableQuery.WriteString(" INTEGER PRIMARY KEY")

	// insert query
	var insertQuery strings.Builder
	insertQuery.WriteString("INSERT INTO ")
	insertQuery.WriteString(classData.tableName)
	insertQuery.WriteString(" (")
	insertQuery.WriteString(classData.keyName)

	// this function registers feature as a column in the table class
	newColumnData := func(eFeature EStructuralFeature, columnType string, encode func(any) (any, error)) {
		// table
		tableQuery.WriteString(",\n")
		tableQuery.WriteString(eFeature.GetName())
		tableQuery.WriteString(columnType)
		// insert
		insertQuery.WriteString(",?")
		// data
		classData.columnData[eFeature] = &sqlEncoderFeatureColumnData{
			ndx:    len(classData.columnData),
			encode: encode,
		}
	}

	newForeignKey := func(eFeature EStructuralFeature, tableName, keyName string) {
		tableQuery.WriteString(",\n")
		tableQuery.WriteString("FOREIGN KEY(")
		tableQuery.WriteString(eFeature.GetName())
		tableQuery.WriteString(") REFERENCES ")
		tableQuery.WriteString(tableName)
		tableQuery.WriteString("(")
		tableQuery.WriteString(keyName)
		tableQuery.WriteString(")")
	}

	for itFeature := eFeatures.Iterator(); itFeature.HasNext(); {
		eFeature := itFeature.Next().(EStructuralFeature)
		featureKind := getSQLCodecFeatureKind(eFeature)
		switch featureKind {
		case sfkObject:
			eReference := eFeature.(EReference)
			referenceData, err := e.getClassData(eReference.GetEReferenceType())
			if err != nil {
				return nil, err
			}
			newColumnData(eFeature, "INTEGER", func(value any) (any, error) {
				object := value.(EObject)
				objectData, err := e.encodeObject(object)
				if err != nil {
					return nil, err
				}
				return objectData.id, nil
			})
			newForeignKey(eFeature, referenceData.tableName, referenceData.keyName)
		case sfkObjectReference:
			newColumnData(eFeature, "TEXT", func(value any) (any, error) {
				object := value.(EObject)
				uri := GetURI(object)
				return uri.String(), nil
			})
		case sfkObjectList:

		case sfkObjectReferenceList:

		case sfkBool, sfkByte, sfkInt, sfkInt16, sfkInt32, sfkInt64, sfkEnum:
			newColumnData(eFeature, "INTEGER", identity)
		case sfkDate:
			newColumnData(eFeature, "TEXT", func(value any) (any, error) {
				t := value.(*time.Time)
				return t.Format(time.RFC3339), nil
			})
		case sfkString:
			newColumnData(eFeature, "TEXT", identity)
		case sfkByteArray:
			newColumnData(eFeature, "BLOB", identity)
		case sfkData:
			eAttribute := eFeature.(EAttribute)
			eDataType := eAttribute.GetEAttributeType()
			eFactory := eDataType.GetEPackage().GetEFactoryInstance()
			newColumnData(eFeature, "TEXT", func(value any) (any, error) {
				return eFactory.ConvertToString(eDataType, value), nil
			})
		case sfkDataList:
			newTableData(eFeature)
		}
	}

	// end
	tableQuery.WriteString("\n);")
	insertQuery.WriteString(");")

	// build queries
	classData.insertQuery = insertQuery.String()
	classData.tableQuery = tableQuery.String()
	return classData
}

func identity(v any) (any, error) { return v, nil }

func (e *SQLEncoder) isObjectWithID() bool {
	idManager := e.resource.GetObjectIDManager()
	return idManager != nil && len(e.idAttributeName) > 0
}

// func (e *SQLEncoder) newFeatureData(eFeature EStructuralFeature) *sqlEncoderFeatureData {
// 	featureData := &sqlEncoderFeatureData{
// 		featureKind: getSQLCodecFeatureKind(eFeature),
// 	}
// 	if eReference, _ := eFeature.(EReference); eReference != nil {
// 		featureData.isTransient = eReference.IsTransient() || (eReference.IsContainer() && !eReference.IsResolveProxies())
// 	} else if eAttribute, _ := eFeature.(EAttribute); eAttribute != nil {
// 		eDataType := eAttribute.GetEAttributeType()
// 		featureData.isTransient = eAttribute.IsTransient()
// 		featureData.dataType = eDataType
// 		featureData.factory = eDataType.GetEPackage().GetEFactoryInstance()
// 	}
// 	return featureData
// }
