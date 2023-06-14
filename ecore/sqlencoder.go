package ecore

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

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

type sqlColumnOption interface {
	apply(col *sqlColumn)
}

type funcSqlColumnOption struct {
	f func(col *sqlColumn)
}

func (fdo *funcSqlColumnOption) apply(col *sqlColumn) {
	fdo.f(col)
}

func newFuncSqlColumnOption(f func(col *sqlColumn)) *funcSqlColumnOption {
	return &funcSqlColumnOption{f: f}
}

func withSqlColumnName(columnName string) sqlColumnOption {
	return newFuncSqlColumnOption(func(col *sqlColumn) {
		col.columnName = columnName
	})
}

func withSqlColumnPrimary(primary bool) sqlColumnOption {
	return newFuncSqlColumnOption(func(col *sqlColumn) {
		col.primary = primary
	})
}

func withSqlColumnAuto(auto bool) sqlColumnOption {
	return newFuncSqlColumnOption(func(col *sqlColumn) {
		col.auto = auto
	})
}

func newSqlAttributeColumn(columnName string, columnType string, options ...sqlColumnOption) *sqlColumn {
	col := &sqlColumn{
		columnName: columnName,
		columnType: columnType,
	}
	for _, opt := range options {
		opt.apply(col)
	}
	return col
}

func newSqlReferenceColumn(reference *sqlTable, options ...sqlColumnOption) *sqlColumn {
	col := &sqlColumn{
		columnName: reference.key.columnName,
		columnType: reference.key.columnType,
		reference:  reference,
	}
	for _, opt := range options {
		opt.apply(col)
	}
	return col
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
	for i, column := range columns {
		t.initColumn(column, i)
	}
	return t
}

func (t *sqlTable) addColumn(column *sqlColumn) {
	t.initColumn(column, len(t.columns))
	t.columns = append(t.columns, column)
}

func (t *sqlTable) initColumn(column *sqlColumn, index int) {
	column.index = index
	if column.primary {
		t.key = column
	}
}

func (t *sqlTable) createQuery() string {
	var tableQuery strings.Builder
	tableQuery.WriteString("CREATE TABLE ")
	tableQuery.WriteString(t.name)
	tableQuery.WriteString(" (")
	// columns
	for i, c := range t.columns {
		if i != 0 {
			tableQuery.WriteString(",")
		}
		tableQuery.WriteString(c.columnName)
		tableQuery.WriteString(" ")
		tableQuery.WriteString(c.columnType)
		if c.primary {
			tableQuery.WriteString(" PRIMARY KEY")
			if c.auto {
				tableQuery.WriteString(" AUTOINCREMENT")
			}
		}
	}
	// constraints
	for _, c := range t.columns {
		if c.reference != nil {
			tableQuery.WriteString(",FOREIGN KEY(")
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
	}
	tableQuery.WriteString(")")
	return tableQuery.String()
}

func (t *sqlTable) defaultValues() []any {
	values := make([]any, len(t.columns))
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

type sqlEncoderFeatureData struct {
	featureKind sqlFeatureKind
	dataType    EDataType
	factory     EFactory
	column      *sqlColumn
	table       *sqlTable
}

type sqlEncoderClassData struct {
	id       int64
	table    *sqlTable
	features []*sqlEncoderFeatureData
}

type sqlEncoderPackageData struct {
	id int64
}

type sqlStmt struct {
	stmt *sql.Stmt
	args []any
}

type sqlStmts struct {
	db    *sql.DB
	stmts []*sqlStmt
}

func newSqlStmts(db *sql.DB) *sqlStmts {
	return &sqlStmts{db: db}
}

func (s *sqlStmts) add(stmt *sql.Stmt, args ...any) {
	s.stmts = append(s.stmts, &sqlStmt{stmt: stmt, args: args})
}

func (s *sqlStmts) exec() error {
	tx, err := s.db.Begin()
	if err != nil {
		return nil
	}
	txStmts := map[*sql.Stmt]*sql.Stmt{}
	for _, t := range s.stmts {
		stmt := t.stmt
		txStmt := txStmts[stmt]
		if txStmt == nil {
			txStmt = tx.Stmt(stmt)
			txStmts[stmt] = txStmt
		}
		_, err := txStmt.Exec(t.args...)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit()
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
		newSqlAttributeColumn("packageID", "INTEGER", withSqlColumnPrimary(true), withSqlColumnAuto(true)),
		newSqlAttributeColumn("uri", "TEXT"),
	)
	classesTable := newSqlTable(
		"classes",
		newSqlAttributeColumn("classID", "INTEGER", withSqlColumnPrimary(true), withSqlColumnAuto(true)),
		newSqlReferenceColumn(packagesTable),
		newSqlAttributeColumn("name", "TEXT"),
	)
	objectsTable := newSqlTable(
		"objects",
		newSqlAttributeColumn("objectID", "INTEGER", withSqlColumnPrimary(true), withSqlColumnAuto(true)),
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
		if _, err := db.Exec(table.createQuery()); err != nil {
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

	// collection of statements
	// used to avoid nested transactions
	insertStmts := newSqlStmts(e.db)

	// encode features columnValues in table columns
	columnValues := classData.table.defaultValues()
	columnValues[classData.table.key.index] = objectID
	for featureID, featureData := range classData.features {
		if featureData.column != nil {
			// feature is encoded as a column
			featureValue := eObject.(EObjectInternal).EGetFromID(featureID, false)
			columnValue, err := e.convertFeatureValue(featureData, featureValue)
			if err != nil {
				return nil, err
			}
			columnValues[featureData.column.index] = columnValue
		} else if featureData.table != nil {
			// feature is encoded in a external table
			featureValue := eObject.(EObjectInternal).EGetFromID(featureID, false)
			featureList, _ := featureValue.(EList)
			if featureList == nil {
				return nil, errors.New("feature value is not a list")
			}
			// retrieve insert statement
			insertStmt, err := e.getInsertStmt(featureData.table)
			if err != nil {
				return nil, err
			}
			// for each list element, insert its value
			index := 0.0
			for itList := featureList.Iterator(); itList.HasNext(); {
				value := itList.Next()
				converted, err := e.convertFeatureValue(featureData, value)
				if err != nil {
					return nil, err
				}
				insertStmts.add(insertStmt, objectID, index, converted)
				index++
			}
		}
	}

	// insert new row in class column
	insertStmt, err := e.getInsertStmt(classData.table)
	if err != nil {
		return nil, err
	}
	insertStmts.add(insertStmt, columnValues...)

	// execute all statements
	if err := insertStmts.exec(); err != nil {
		return nil, err
	}

	return &sqlEncoderObjectData{
		id:        objectID,
		classData: classData,
	}, nil
}

func (e *SQLEncoder) convertFeatureValue(featureData *sqlEncoderFeatureData, value any) (any, error) {
	if value != nil {
		switch featureData.featureKind {
		case sfkObject, sfkObjectList:
			objectData, err := e.encodeObject(value.(EObject))
			if err != nil {
				return nil, err
			}
			return objectData.id, nil
		case sfkObjectReference, sfkObjectReferenceList:
			return GetURI(value.(EObject)).String(), nil
		case sfkBool, sfkByte, sfkInt, sfkInt16, sfkInt32, sfkInt64, sfkEnum, sfkString, sfkByteArray:
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
		return e.newClassData(eClass, id)
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
	classTable.addColumn(newSqlAttributeColumn(strings.ToLower(eClass.GetName())+"ID", "INTEGER", withSqlColumnPrimary(true)))

	// compute table columns and external tables
	classData := &sqlEncoderClassData{
		id:       id,
		table:    classTable,
		features: make([]*sqlEncoderFeatureData, 0, eFeatures.Size()),
	}

	// register class data now to handle correctly cycles references
	e.classDataMap[eClass] = classData

	newFeatureReferenceColumn := func(featureData *sqlEncoderFeatureData, eFeature EStructuralFeature, table *sqlTable) {
		column := newSqlReferenceColumn(table, withSqlColumnName(eFeature.GetName()))
		classTable.addColumn(column)
		featureData.column = column
	}

	newFeatureAttributeColumn := func(featureData *sqlEncoderFeatureData, eFeature EStructuralFeature, columnType string) {
		column := newSqlAttributeColumn(eFeature.GetName(), columnType)
		classTable.addColumn(column)
		featureData.column = column
	}

	newFeatureTable := func(featureData *sqlEncoderFeatureData, eFeature EStructuralFeature, columns ...*sqlColumn) {
		featureData.table = newSqlTable(
			classTable.name+"_"+eFeature.GetName(),
			columns...,
		)
	}

	for itFeature := eFeatures.Iterator(); itFeature.HasNext(); {
		eFeature := itFeature.Next().(EStructuralFeature)
		// new feature data
		featureData := e.newFeatureData(eFeature)
		classData.features = append(classData.features, featureData)

		// compute class table columns or children tables
		switch featureData.featureKind {
		case sfkObject:
			// retrieve object reference type
			eReference := eFeature.(EReference)
			referenceData, err := e.getClassData(eReference.GetEReferenceType())
			if err != nil {
				return nil, err
			}
			newFeatureReferenceColumn(featureData, eFeature, referenceData.table)
		case sfkObjectReference:
			newFeatureAttributeColumn(featureData, eFeature, "TEXT")
		case sfkObjectList:
			// internal reference
			eReference := eFeature.(EReference)
			referenceData, err := e.getClassData(eReference.GetEReferenceType())
			if err != nil {
				return nil, err
			}
			newFeatureTable(featureData, eFeature,
				newSqlReferenceColumn(classData.table),
				newSqlAttributeColumn("idx", "REAL"),
				newSqlReferenceColumn(referenceData.table, withSqlColumnName(eFeature.GetName())),
			)
		case sfkObjectReferenceList:
			newFeatureTable(featureData, eFeature,
				newSqlReferenceColumn(classData.table),
				newSqlAttributeColumn("idx", "REAL"),
				newSqlAttributeColumn("uri", "TEXT"),
			)
		case sfkBool, sfkByte, sfkInt, sfkInt16, sfkInt32, sfkInt64, sfkEnum:
			newFeatureAttributeColumn(featureData, eFeature, "INTEGER")
		case sfkDate, sfkString, sfkData:
			newFeatureAttributeColumn(featureData, eFeature, "TEXT")
		case sfkByteArray:
			newFeatureAttributeColumn(featureData, eFeature, "BLOB")
		case sfkDataList:
			newFeatureTable(featureData, eFeature,
				newSqlReferenceColumn(classData.table),
				newSqlAttributeColumn("idx", "REAL"),
				newSqlAttributeColumn(eFeature.GetName(), "TEXT"),
			)
		}
	}

	// create class table
	if _, err := e.db.Exec(classData.table.createQuery()); err != nil {
		return nil, err
	}

	// create children tables
	for _, featureData := range classData.features {
		if table := featureData.table; table != nil {
			if _, err := e.db.Exec(table.createQuery()); err != nil {
				return nil, err
			}
		}
	}

	return classData, nil
}

func (e *SQLEncoder) newFeatureData(eFeature EStructuralFeature) *sqlEncoderFeatureData {
	featureData := &sqlEncoderFeatureData{
		featureKind: getSQLCodecFeatureKind(eFeature),
	}
	if eAttribute, _ := eFeature.(EAttribute); eAttribute != nil {
		eDataType := eAttribute.GetEAttributeType()
		featureData.dataType = eDataType
		featureData.factory = eDataType.GetEPackage().GetEFactoryInstance()
	}
	return featureData
}

func (e *SQLEncoder) isObjectWithID() bool {
	idManager := e.resource.GetObjectIDManager()
	return idManager != nil && len(e.idAttributeName) > 0
}
