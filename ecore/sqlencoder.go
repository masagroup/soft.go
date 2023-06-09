package ecore

import (
	"database/sql"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

type sqlColumn struct {
	index      int
	columnName string
	columnType string
	primary    bool
	auto       bool
	reference  *sqlTable
}

func newSqlAttributeColumn(columnName string, columnType string) *sqlColumn {
	return &sqlColumn{
		columnName: columnName,
		columnType: columnType,
	}
}

func newSqlReferenceColumn(reference *sqlTable) *sqlColumn {
	return &sqlColumn{
		columnName: reference.key.columnName,
		columnType: reference.key.columnType,
		reference:  reference,
	}
}

func (c *sqlColumn) setPrimary(primary bool) *sqlColumn {
	c.primary = primary
	return c
}

func (c *sqlColumn) setAuto(auto bool) *sqlColumn {
	c.auto = auto
	return c
}

type sqlTable struct {
	name    string
	key     *sqlColumn
	columns []*sqlColumn
}

func newSqlTable(name string, columns ...*sqlColumn) *sqlTable {
	t := &sqlTable{
		name:    name,
		columns: columns,
	}
	for _, column := range columns {
		t.addColumn(column)
	}
	return t
}

func (t *sqlTable) addColumn(column *sqlColumn) {
	column.index = len(t.columns)
	t.columns = append(t.columns, column)
	if column.primary {
		t.key = column
	}
}

func (t *sqlTable) createQuery() string {
	var tableQuery strings.Builder
	tableQuery.WriteString("CREATE TABLE ")
	tableQuery.WriteString(t.name)
	tableQuery.WriteString(" (")
	for i, c := range t.columns {
		if i != 0 {
			tableQuery.WriteString(",")
		}
		tableQuery.WriteString("\n")
		tableQuery.WriteString(c.columnName)
		tableQuery.WriteString(" ")
		tableQuery.WriteString(c.columnType)
		if c.primary {
			tableQuery.WriteString(" PRIMARY KEY")
			if c.auto {
				tableQuery.WriteString(" AUTOINCREMENT")
			}
		}
		if c.reference != nil {
			tableQuery.WriteString(",\n")
			tableQuery.WriteString("FOREIGN KEY(")
			tableQuery.WriteString(c.columnName)
			tableQuery.WriteString(") REFERENCES ")
			tableQuery.WriteString(c.reference.name)
			tableQuery.WriteString("(")
			tableQuery.WriteString(c.reference.key.columnName)
			tableQuery.WriteString(")")
		}
	}
	tableQuery.WriteString(")")
	return tableQuery.String()
}

func (t *sqlTable) insertQuery() string {
	var tableQuery strings.Builder
	tableQuery.WriteString("INSERT INTO ")
	tableQuery.WriteString(t.name)
	tableQuery.WriteString(" (")
	for i, c := range t.columns {
		if i != 0 {
			tableQuery.WriteString(",")
		}
		tableQuery.WriteString(c.columnName)
	}
	tableQuery.WriteString(") VALUES (")
	for i, c := range t.columns {
		if i != 0 {
			tableQuery.WriteString(",")
		}
		if c.auto {
			tableQuery.WriteString("NULL")
		} else {
			tableQuery.WriteString("?")
		}
		tableQuery.WriteString(c.columnName)
	}
	tableQuery.WriteString(")")
	return tableQuery.String()
}

func (t *sqlTable) insertValues() []any {
	values := make([]any, 0, len(t.columns))
	for i, c := range t.columns {
		if c.auto {
			switch c.columnType {
			case "TEXT":
				values[i] = sql.NullString{}
			case "INTEGER":
				values[i] = sql.NullInt64{}
			}
		}
	}
	return values
}

type sqlEncoderObjectData struct {
	id        int64
	classData *sqlEncoderClassData
}

type sqlEncoderClassData struct {
	id             int64
	table          *sqlTable
	featureColumns map[EStructuralFeature]*sqlColumn
	featureTables  map[EStructuralFeature]*sqlTable
}

func newClassData(id int64, table *sqlTable) *sqlEncoderClassData {
	return &sqlEncoderClassData{
		id:             id,
		table:          table,
		featureColumns: map[EStructuralFeature]*sqlColumn{},
		featureTables:  map[EStructuralFeature]*sqlTable{},
	}
}

func (classData *sqlEncoderClassData) addReferenceColumn(eFeature EStructuralFeature, reference *sqlTable) *sqlColumn {
	column := newSqlReferenceColumn(reference)
	column.columnName = eFeature.GetName()
	classData.table.addColumn(column)
	classData.featureColumns[eFeature] = column
	return column
}

func (classData *sqlEncoderClassData) addAttributeColumn(eFeature EStructuralFeature, columnType string) *sqlColumn {
	column := newSqlAttributeColumn(eFeature.GetName(), columnType)
	classData.table.addColumn(column)
	classData.featureColumns[eFeature] = column
	return column
}

func (classData *sqlEncoderClassData) addFeatureTable(eFeature EStructuralFeature, table *sqlTable) {
	classData.featureTables[eFeature] = table
}

type sqlEncoderPackageData struct {
	id int64
}

type SQLEncoder struct {
	resource        EResource
	writer          io.Writer
	driver          string
	db              *sql.DB
	insertStmts     map[*sqlTable]*sql.Stmt
	packageDataMap  map[EPackage]*sqlEncoderPackageData
	classDataMap    map[EClass]*sqlEncoderClassData
	packagesTable   *sqlTable
	classesTable    *sqlTable
	objectsTable    *sqlTable
	contentsTable   *sqlTable
	idAttributeName string
}

func NewSQLEncoder(resource EResource, w io.Writer, options map[string]any) *SQLEncoder {
	// common tables definitions
	packagesTable := newSqlTable(
		"packages",
		newSqlAttributeColumn("packageID", "INTEGER").setPrimary(true).setAuto(true),
		newSqlAttributeColumn("uri", "TEXT"),
	)
	classesTable := newSqlTable(
		"classes",
		newSqlAttributeColumn("classID", "INTEGER").setPrimary(true).setAuto(true),
		newSqlAttributeColumn("name", "TEXT"),
		newSqlReferenceColumn(packagesTable),
	)
	objectsTable := newSqlTable(
		"objects",
		newSqlAttributeColumn("objectID", "INTEGER").setPrimary(true).setAuto(true),
		newSqlReferenceColumn(classesTable),
	)
	contentsTable := newSqlTable(
		"contents",
		newSqlReferenceColumn(objectsTable),
	)

	// encoder structure
	e := &SQLEncoder{
		resource:       resource,
		writer:         w,
		driver:         "sqlite",
		packagesTable:  packagesTable,
		classesTable:   classesTable,
		objectsTable:   objectsTable,
		contentsTable:  contentsTable,
		packageDataMap: map[EPackage]*sqlEncoderPackageData{},
		classDataMap:   map[EClass]*sqlEncoderClassData{},
		insertStmts:    map[*sqlTable]*sql.Stmt{},
	}

	// options
	if options != nil {
		if driver, isDriver := options[SQL_OPTION_DRIVER]; isDriver {
			e.driver = driver.(string)
		}

		e.idAttributeName, _ = options[JSON_OPTION_ID_ATTRIBUTE_NAME].(string)
		if e.isObjectWithID() {
			e.objectsTable.addColumn(newSqlAttributeColumn(e.idAttributeName, "TEXT"))
		}
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
		e.packagesTable,
		e.classesTable,
		e.objectsTable,
		e.contentsTable,
	} {
		if _, err := e.db.Exec(table.createQuery()); err != nil {
			return nil, err
		}
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

	insertContentStmt, err := e.getInsertStmt(e.contentsTable)
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
	insertObjectStmt, err := e.getInsertStmt(e.objectsTable)
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
		insertClassStmt, err := e.getInsertStmt(e.classesTable)
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
		classData, err = e.newClassData(eClass, id)
		if err != nil {
			return nil, err
		}

		// register class data
		e.classDataMap[eClass] = classData
	}
	return classData, nil
}

func (e *SQLEncoder) getPackageData(ePackage EPackage) (*sqlEncoderPackageData, error) {
	ePackageData := e.packageDataMap[ePackage]
	if ePackageData == nil {
		// insert new package
		insertPackageStmt, err := e.getInsertStmt(e.packagesTable)
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

func (e *SQLEncoder) getInsertStmt(table *sqlTable) (stmt *sql.Stmt, err error) {
	stmt = e.insertStmts[table]
	if stmt == nil {
		stmt, err = e.db.Prepare(table.insertQuery())
	}
	return
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
	// create table descriptor
	classTable := newSqlTable(ePackage.GetNsPrefix() + "_" + strings.ToLower(eClass.GetName()))
	classTable.addColumn(newSqlAttributeColumn(strings.ToLower(eClass.GetName())+"ID", "INTEGER").setPrimary(true))

	// compute table columns and external tables
	classData := newClassData(id, classTable)
	for itFeature := eFeatures.Iterator(); itFeature.HasNext(); {
		eFeature := itFeature.Next().(EStructuralFeature)
		featureKind := getSQLCodecFeatureKind(eFeature)
		switch featureKind {
		case sfkObject:
			// retrieve object reference type
			eReference := eFeature.(EReference)
			referenceData, err := e.getClassData(eReference.GetEReferenceType())
			if err != nil {
				return nil, err
			}
			classData.addReferenceColumn(eFeature, referenceData.table)
		case sfkObjectReference:
			classData.addAttributeColumn(eFeature, "TEXT")
		case sfkObjectList:
			// internal reference
			eReference := eFeature.(EReference)
			referenceData, err := e.getClassData(eReference.GetEReferenceType())
			if err != nil {
				return nil, err
			}
			// table descriptor
			referenceTable := newSqlTable(
				classTable.name+"_"+eFeature.GetName(),
				newSqlReferenceColumn(classData.table),
				newSqlReferenceColumn(referenceData.table),
			)
			// table create
			tableQuery := referenceTable.createQuery()
			if _, err := e.db.Exec(tableQuery); err != nil {
				return nil, err
			}
			// table registering
			classData.addFeatureTable(eFeature, referenceTable)
		case sfkObjectReferenceList:
			// table descriptor
			referenceTable := newSqlTable(
				classTable.name+"_"+eFeature.GetName(),
				newSqlReferenceColumn(classData.table),
				newSqlAttributeColumn("uri", "TEXT"),
			)
			// table create
			tableQuery := referenceTable.createQuery()
			if _, err := e.db.Exec(tableQuery); err != nil {
				return nil, err
			}
			// table registering
			classData.addFeatureTable(eFeature, referenceTable)
		case sfkBool, sfkByte, sfkInt, sfkInt16, sfkInt32, sfkInt64, sfkEnum:
			classData.addAttributeColumn(eFeature, "INTEGER")
		case sfkDate:
			classData.addAttributeColumn(eFeature, "TEXT")
		case sfkString:
			classData.addAttributeColumn(eFeature, "TEXT")
		case sfkByteArray:
			classData.addAttributeColumn(eFeature, "BLOB")
		case sfkData:
			classData.addAttributeColumn(eFeature, "TEXT")
		case sfkDataList:
			dataTable := newSqlTable(
				classTable.name+"_"+eFeature.GetName(),
				newSqlReferenceColumn(classData.table),
				newSqlAttributeColumn(eFeature.GetName(), "TEXT"),
			)
			// table create
			tableQuery := dataTable.createQuery()
			if _, err := e.db.Exec(tableQuery); err != nil {
				return nil, err
			}
			// table registering
			classData.addFeatureTable(eFeature, dataTable)
		}
	}

	// create class table
	tableQuery := classData.table.createQuery()
	if _, err := e.db.Exec(tableQuery); err != nil {
		return nil, err
	}

	return classData, nil
}

func (e *SQLEncoder) isObjectWithID() bool {
	idManager := e.resource.GetObjectIDManager()
	return idManager != nil && len(e.idAttributeName) > 0
}
