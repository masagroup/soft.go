package ecore

import (
	"strings"
	"sync"
)

type sqlColumn struct {
	table      *sqlTable
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

func (o *funcSqlColumnOption) apply(col *sqlColumn) {
	o.f(col)
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

type sqlTableOption interface {
	apply(t *sqlTable)
}

type funcSqlTableOption struct {
	f func(t *sqlTable)
}

func (o *funcSqlTableOption) apply(t *sqlTable) {
	o.f(t)
}

func newFuncSqlTableOption(f func(t *sqlTable)) *funcSqlTableOption {
	return &funcSqlTableOption{f: f}
}

func withSqlTableColumns(columns ...*sqlColumn) sqlTableOption {
	return newFuncSqlTableOption(func(t *sqlTable) {
		t.columns = columns
		for i, column := range columns {
			t.initColumn(column, i)
		}
	})
}

func withSqlTableCreateIfNotExists(createIfNotExists bool) sqlTableOption {
	return newFuncSqlTableOption(func(t *sqlTable) {
		t.createIfNotExists = createIfNotExists
	})
}

type sqlTable struct {
	name              string
	key               *sqlColumn
	columns           []*sqlColumn
	indexes           [][]*sqlColumn
	createIfNotExists bool
}

func newSqlTable(name string, options ...sqlTableOption) *sqlTable {
	t := &sqlTable{
		name: name,
	}
	for _, opt := range options {
		opt.apply(t)
	}
	return t
}

func (t *sqlTable) addColumn(column *sqlColumn) {
	t.initColumn(column, len(t.columns))
	t.columns = append(t.columns, column)
}

func (t *sqlTable) initColumn(column *sqlColumn, index int) {
	column.table = t
	column.index = index
	if column.primary {
		t.key = column
	}
}

func (t *sqlTable) createQuery() string {
	var tableQuery strings.Builder
	tableQuery.WriteString("CREATE TABLE ")
	if t.createIfNotExists {
		tableQuery.WriteString("IF NOT EXISTS ")
	}
	tableQuery.WriteString(sqlEscapeIdentifier(t.name))
	tableQuery.WriteString(" (")
	// columns
	for i, c := range t.columns {
		if i != 0 {
			tableQuery.WriteString(",")
		}
		tableQuery.WriteString(sqlEscapeIdentifier(c.columnName))
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
			tableQuery.WriteString(sqlEscapeIdentifier(c.columnName))
			tableQuery.WriteString(") REFERENCES ")
			tableQuery.WriteString(sqlEscapeIdentifier(c.reference.name))
			tableQuery.WriteString("(")
			tableQuery.WriteString(sqlEscapeIdentifier(c.reference.key.columnName))
			tableQuery.WriteString(")")
		}
	}
	tableQuery.WriteString(");")
	for _, index := range t.indexes {
		tableQuery.WriteString("\n")
		tableQuery.WriteString("CREATE INDEX ")
		if t.createIfNotExists {
			tableQuery.WriteString("IF NOT EXISTS ")
		}
		tableQuery.WriteString("\"idx_")
		tableQuery.WriteString(t.name)
		for _, c := range index {
			tableQuery.WriteString("_")
			tableQuery.WriteString(c.columnName)
		}
		tableQuery.WriteString("\" ON ")
		tableQuery.WriteString(sqlEscapeIdentifier(t.name))
		tableQuery.WriteString("(")
		for i, c := range index {
			if i != 0 {
				tableQuery.WriteString(",")
			}
			tableQuery.WriteString(sqlEscapeIdentifier(c.columnName))
		}
		tableQuery.WriteString(");")
	}
	return tableQuery.String()
}

func (t *sqlTable) insertQuery() string {
	var tableQuery strings.Builder
	tableQuery.WriteString("INSERT INTO ")
	tableQuery.WriteString(sqlEscapeIdentifier(t.name))
	tableQuery.WriteString(" (")
	for i, c := range t.columns {
		if i != 0 {
			tableQuery.WriteString(",")
		}
		tableQuery.WriteString(sqlEscapeIdentifier(c.columnName))
	}
	tableQuery.WriteString(") VALUES (")
	for i := range t.columns {
		if i != 0 {
			tableQuery.WriteString(",")
		}
		tableQuery.WriteString("?")
	}
	tableQuery.WriteString(") RETURNING ")
	if t.key != nil {
		tableQuery.WriteString(sqlEscapeIdentifier(t.key.columnName))
	} else {
		tableQuery.WriteString("rowid")
	}
	return tableQuery.String()
}

func (t *sqlTable) insertOrReplaceQuery() string {
	var tableQuery strings.Builder
	tableQuery.WriteString("INSERT OR REPLACE INTO ")
	tableQuery.WriteString(sqlEscapeIdentifier(t.name))
	tableQuery.WriteString(" (")
	for i, c := range t.columns {
		if i != 0 {
			tableQuery.WriteString(",")
		}
		tableQuery.WriteString(sqlEscapeIdentifier(c.columnName))
	}
	tableQuery.WriteString(") VALUES (")
	for i := range t.columns {
		if i != 0 {
			tableQuery.WriteString(",")
		}
		tableQuery.WriteString("?")
	}
	tableQuery.WriteString(")")
	return tableQuery.String()
}

func (t *sqlTable) defaultValues() []any {
	return make([]any, len(t.columns))
}

func (t *sqlTable) keyName() string {
	return sqlEscapeIdentifier(t.key.columnName)
}

func (t *sqlTable) selectQuery(columns []string, selection string, orderBy string) string {
	var selectQuery strings.Builder
	selectQuery.WriteString("SELECT ")
	if len(columns) == 0 {
		selectQuery.WriteString("*")
	} else {
		for i, column := range columns {
			if i != 0 {
				selectQuery.WriteString(",")
			}
			selectQuery.WriteString(sqlEscapeIdentifier(column))
		}
	}
	selectQuery.WriteString(" FROM ")
	selectQuery.WriteString(sqlEscapeIdentifier(t.name))
	if len(selection) > 0 {
		selectQuery.WriteString(" WHERE ")
		selectQuery.WriteString(selection)
	}
	if len(orderBy) > 0 {
		selectQuery.WriteString(" ORDER BY ")
		selectQuery.WriteString(orderBy)
	}
	return selectQuery.String()
}

type sqlClassSchema struct {
	table    *sqlTable
	features []*sqlFeatureSchema
}

func (s *sqlClassSchema) getFeatureSchema(feature EStructuralFeature) *sqlFeatureSchema {
	for _, f := range s.features {
		if f.feature == feature {
			return f
		}
	}
	return nil
}

type sqlFeatureSchema struct {
	featureKind sqlFeatureKind
	feature     EStructuralFeature
	column      *sqlColumn
	table       *sqlTable
}

type sqlSchema struct {
	mutex             sync.Mutex
	propertiesTable   *sqlTable
	packagesTable     *sqlTable
	classesTable      *sqlTable
	objectsTable      *sqlTable
	contentsTable     *sqlTable
	enumsTable        *sqlTable
	classSchemaMap    map[EClass]*sqlClassSchema
	createIfNotExists bool
	isContainerID     bool
	objectIDName      string
}

type sqlSchemaOption interface {
	apply(s *sqlSchema)
}

type funcSqlSchemaOption struct {
	f func(s *sqlSchema)
}

func (o *funcSqlSchemaOption) apply(s *sqlSchema) {
	o.f(s)
}

func newFuncSqlSchemaOption(f func(s *sqlSchema)) *funcSqlSchemaOption {
	return &funcSqlSchemaOption{f: f}
}

func withObjectIDName(objectIDName string) sqlSchemaOption {
	return newFuncSqlSchemaOption(func(s *sqlSchema) {
		s.objectIDName = objectIDName
	})
}

func withContainerID(isContainerID bool) sqlSchemaOption {
	return newFuncSqlSchemaOption(func(s *sqlSchema) {
		s.isContainerID = isContainerID
	})
}

func withCreateIfNotExists(createIfNotExists bool) sqlSchemaOption {
	return newFuncSqlSchemaOption(func(s *sqlSchema) {
		s.createIfNotExists = createIfNotExists
	})
}

func newSqlSchema(options ...sqlSchemaOption) *sqlSchema {
	// create scheam and apply options
	s := &sqlSchema{
		classSchemaMap: map[EClass]*sqlClassSchema{},
	}
	for _, opt := range options {
		opt.apply(s)
	}

	// common tables definitions
	s.propertiesTable = newSqlTable(
		".properties",
		withSqlTableColumns(
			newSqlAttributeColumn("key", "TEXT", withSqlColumnPrimary(true)),
			newSqlAttributeColumn("value", "TEXT"),
		),
		withSqlTableCreateIfNotExists(s.createIfNotExists),
	)
	s.packagesTable = newSqlTable(
		".packages",
		withSqlTableColumns(
			newSqlAttributeColumn("packageID", "INTEGER", withSqlColumnPrimary(true), withSqlColumnAuto(true)),
			newSqlAttributeColumn("uri", "TEXT"),
		),
		withSqlTableCreateIfNotExists(s.createIfNotExists),
	)
	s.classesTable = newSqlTable(
		".classes",
		withSqlTableColumns(
			newSqlAttributeColumn("classID", "INTEGER", withSqlColumnPrimary(true), withSqlColumnAuto(true)),
			newSqlReferenceColumn(s.packagesTable),
			newSqlAttributeColumn("name", "TEXT"),
		),
		withSqlTableCreateIfNotExists(s.createIfNotExists),
	)
	s.objectsTable = newSqlTable(
		".objects",
		withSqlTableColumns(
			newSqlAttributeColumn("objectID", "INTEGER", withSqlColumnPrimary(true), withSqlColumnAuto(true)),
			newSqlReferenceColumn(s.classesTable),
		),
		withSqlTableCreateIfNotExists(s.createIfNotExists),
	)
	// container and feayure id in objects table
	if s.isContainerID {
		s.objectsTable.addColumn(newSqlReferenceColumn(s.objectsTable, withSqlColumnName("containerID")))
		s.objectsTable.addColumn(newSqlAttributeColumn("containerFeatureID", "INTEGER"))
	}
	// add id attribute column if name is not object table primary key
	if len(s.objectIDName) > 0 && s.objectIDName != s.objectsTable.key.columnName {
		s.objectsTable.addColumn(newSqlAttributeColumn(s.objectIDName, "TEXT"))
	}
	s.contentsTable = newSqlTable(
		".contents",
		withSqlTableColumns(
			newSqlReferenceColumn(s.objectsTable),
		),
		withSqlTableCreateIfNotExists(s.createIfNotExists),
	)
	s.enumsTable = newSqlTable(
		".enums",
		withSqlTableColumns(
			newSqlAttributeColumn("enumID", "INTEGER", withSqlColumnPrimary(true), withSqlColumnAuto(true)),
			newSqlReferenceColumn(s.packagesTable),
			newSqlAttributeColumn("name", "TEXT"),
			newSqlAttributeColumn("literal", "TEXT"),
		),
		withSqlTableCreateIfNotExists(s.createIfNotExists),
	)
	return s
}

func (s *sqlSchema) getClassSchema(eClass EClass) *sqlClassSchema {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.getOrComputeClassSchema(eClass)
}

func (s *sqlSchema) getOrComputeClassSchema(eClass EClass) *sqlClassSchema {
	classSchema := s.classSchemaMap[eClass]
	if classSchema == nil {
		// create table descriptor
		classTable := newSqlTable(strings.ToLower(eClass.GetName()), withSqlTableCreateIfNotExists(s.createIfNotExists))
		classTable.addColumn(newSqlAttributeColumn(strings.ToLower(eClass.GetName())+"ID", "INTEGER", withSqlColumnPrimary(true)))

		// compute table columns and external tables
		eFeatures := eClass.GetEStructuralFeatures()
		classSchema = &sqlClassSchema{
			table:    classTable,
			features: make([]*sqlFeatureSchema, 0, eFeatures.Size()),
		}

		// register class data now to handle correctly cycles references
		s.classSchemaMap[eClass] = classSchema

		newFeatureReferenceColumn := func(featureSchema *sqlFeatureSchema, eFeature EStructuralFeature, table *sqlTable) {
			column := newSqlReferenceColumn(table, withSqlColumnName(eFeature.GetName()))
			classTable.addColumn(column)
			featureSchema.column = column
		}

		newFeatureAttributeColumn := func(featureSchema *sqlFeatureSchema, eFeature EStructuralFeature, columnType string) {
			column := newSqlAttributeColumn(eFeature.GetName(), columnType)
			classTable.addColumn(column)
			featureSchema.column = column
		}

		newFeatureTable := func(featureSchema *sqlFeatureSchema, eFeature EStructuralFeature, columns ...*sqlColumn) {
			table := newSqlTable(
				classTable.name+"_"+eFeature.GetName(),
				withSqlTableColumns(columns...),
				withSqlTableCreateIfNotExists(s.createIfNotExists),
			)
			table.key = columns[0]
			table.indexes = [][]*sqlColumn{{columns[0], columns[1]}}
			featureSchema.table = table
		}

		for itFeature := eFeatures.Iterator(); itFeature.HasNext(); {
			eFeature := itFeature.Next().(EStructuralFeature)
			// new feature data
			featureSchema := &sqlFeatureSchema{
				feature:     eFeature,
				featureKind: getSQLCodecFeatureKind(eFeature),
			}
			classSchema.features = append(classSchema.features, featureSchema)

			// compute class table columns or children tables
			switch featureSchema.featureKind {
			case sfkObject:
				// retrieve object reference type
				eReference := eFeature.(EReference)
				referenceSchema := s.getOrComputeClassSchema(eReference.GetEReferenceType())
				newFeatureReferenceColumn(featureSchema, eFeature, referenceSchema.table)
			case sfkObjectReference:
				newFeatureAttributeColumn(featureSchema, eFeature, "TEXT")
			case sfkObjectList:
				// internal reference
				eReference := eFeature.(EReference)
				referenceSchema := s.getOrComputeClassSchema(eReference.GetEReferenceType())
				newFeatureTable(featureSchema, eFeature,
					newSqlReferenceColumn(classTable),
					newSqlAttributeColumn("idx", "REAL"),
					newSqlReferenceColumn(referenceSchema.table, withSqlColumnName(eFeature.GetName())),
				)
			case sfkObjectReferenceList:
				newFeatureTable(featureSchema, eFeature,
					newSqlReferenceColumn(classTable),
					newSqlAttributeColumn("idx", "REAL"),
					newSqlAttributeColumn("uri", "TEXT"),
				)
			case sfkEnum:
				newFeatureReferenceColumn(featureSchema, eFeature, s.enumsTable)
			case sfkBool, sfkByte, sfkInt, sfkInt16, sfkInt32, sfkInt64:
				newFeatureAttributeColumn(featureSchema, eFeature, "INTEGER")
			case sfkDate, sfkString, sfkData:
				newFeatureAttributeColumn(featureSchema, eFeature, "TEXT")
			case sfkByteArray:
				newFeatureAttributeColumn(featureSchema, eFeature, "BLOB")
			case sfkDataList:
				newFeatureTable(featureSchema, eFeature,
					newSqlReferenceColumn(classTable),
					newSqlAttributeColumn("idx", "REAL"),
					newSqlAttributeColumn(eFeature.GetName(), "TEXT"),
				)
			}
		}
	}
	return classSchema
}

var sqliteKeywWords = []string{
	"ABORT", "ACTION", "ADD", "AFTER", "ALL", "ALTER", "ALWAYS", "ANALYZE", "AND", "AS", "ASC", "ATTACH",
	"AUTOINCREMENT", "BEFORE", "BEGIN", "BETWEEN", "BY", "CASCADE", "CASE", "CAST", "CHECK", "COLLATE",
	"COLUMN", "COMMIT", "CONFLICT", "CONSTRAINT", "CREATE", "CROSS", "CURRENT", "CURRENT_DATE",
	"CURRENT_TIME", "CURRENT_TIMESTAMP", "DATABASE", "DEFAULT", "DEFERRABLE", "DEFERRED", "DELETE",
	"DESC", "DETACH", "DISTINCT", "DO", "DROP", "EACH", "ELSE", "END", "ESCAPE", "EXCEPT", "EXCLUDE",
	"EXCLUSIVE", "EXISTS", "EXPLAIN", "FAIL", "FILTER", "FIRST", "FOLLOWING", "FOR", "FOREIGN", "FROM",
	"FULL", "GENERATED", "GLOB", "GROUP", "GROUPS", "HAVING", "IF", "IGNORE", "IMMEDIATE", "IN", "INDEX",
	"INDEXED", "INITIALLY", "INNER", "INSERT", "INSTEAD", "INTERSECT", "INTO", "IS", "ISNULL", "JOIN",
	"KEY", "LAST", "LEFT", "LIKE", "LIMIT", "MATCH", "MATERIALIZED", "NATURAL", "NO", "NOT", "NOTHING",
	"NOTNULL", "NULL", "NULLS", "OF", "OFFSET", "ON", "OR", "ORDER", "OTHERS", "OUTER", "OVER", "PARTITION",
	"PLAN", "PRAGMA", "PRECEDING", "PRIMARY", "QUERY", "RAISE", "RANGE", "RECURSIVE", "REFERENCES", "REGEXP",
	"REINDEX", "RELEASE", "RENAME", "REPLACE", "RESTRICT", "RETURNING", "RIGHT", "ROLLBACK", "ROW", "ROWS", "SAVEPOINT",
	"SELECT", "SET", "TABLE", "TEMP", "TEMPORARY", "THEN", "TIES", "TO", "TRANSACTION", "TRIGGER", "UNBOUNDED", "UNION",
	"UNIQUE", "UPDATE", "USING", "VACUUM", "VALUES", "VIEW", "VIRTUAL", "WHEN", "WHERE", "WINDOW", "WITH", "WITHOUT",
}

var sqliteKeyWordsAsSet = func() (result map[string]struct{}) {
	result = map[string]struct{}{}
	for _, keyword := range sqliteKeywWords {
		result[keyword] = struct{}{}
	}
	return
}()

func isSqliteKeyWord(k string) bool {
	if _, isKeyWord := sqliteKeyWordsAsSet[k]; isKeyWord {
		return true
	}
	if _, isKeyWord := sqliteKeyWordsAsSet[strings.ToUpper(k)]; isKeyWord {
		return true
	}
	return false
}

func sqlEscapeIdentifier(id string) string {
	if id[0] == '.' || isSqliteKeyWord(id) {
		return "\"" + id + "\""
	}
	return id
}
