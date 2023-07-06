package ecore

import (
	"database/sql"
	"strings"
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

func withSqlTableIndexes(tableIndexes [][]string) sqlTableOption {
	return newFuncSqlTableOption(func(t *sqlTable) {
		nameToColumn := map[string]*sqlColumn{}
		for _, column := range t.columns {
			nameToColumn[column.columnName] = column
		}
		for _, indexes := range tableIndexes {
			indexColumns := []*sqlColumn{}
			for _, index := range indexes {
				indexColumns = append(indexColumns, nameToColumn[index])
			}
			t.indexes = append(t.indexes, indexColumns)
		}
	})
}

type sqlTable struct {
	name    string
	key     *sqlColumn
	columns []*sqlColumn
	indexes [][]*sqlColumn
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
	column.index = index
	if column.primary {
		t.key = column
	}
}

func (t *sqlTable) createQuery() string {
	var tableQuery strings.Builder
	tableQuery.WriteString("CREATE TABLE ")
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
		tableQuery.WriteString("CREATE INDEX \"idx_")
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

func (t *sqlTable) columnNames(first, last int) []string {
	if last == -1 {
		last = len(t.columns)
	}
	columns := t.columns[first:last]
	names := make([]string, 0, len(columns))
	for _, column := range columns {
		names = append(names, sqlEscapeIdentifier(column.columnName))
	}
	return names
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
			selectQuery.WriteString(column)
		}
	}
	selectQuery.WriteString(" from ")
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

type sqlFeatureSchema struct {
	featureKind sqlFeatureKind
	feature     EStructuralFeature
	column      *sqlColumn
	table       *sqlTable
}

type sqlSchema struct {
	packagesTable  *sqlTable
	classesTable   *sqlTable
	objectsTable   *sqlTable
	contentsTable  *sqlTable
	enumsTable     *sqlTable
	classSchemaMap map[EClass]*sqlClassSchema
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

func withIDAttributeName(idAttributeName string) sqlSchemaOption {
	return newFuncSqlSchemaOption(func(s *sqlSchema) {
		if len(idAttributeName) > 0 {
			s.objectsTable.addColumn(newSqlAttributeColumn(idAttributeName, "TEXT"))
		}
	})
}

func newSqlSchema(options ...sqlSchemaOption) *sqlSchema {

	// common tables definitions
	packagesTable := newSqlTable(
		".packages",
		withSqlTableColumns(
			newSqlAttributeColumn("packageID", "INTEGER", withSqlColumnPrimary(true), withSqlColumnAuto(true)),
			newSqlAttributeColumn("uri", "TEXT"),
		),
	)
	classesTable := newSqlTable(
		".classes",
		withSqlTableColumns(
			newSqlAttributeColumn("classID", "INTEGER", withSqlColumnPrimary(true), withSqlColumnAuto(true)),
			newSqlReferenceColumn(packagesTable),
			newSqlAttributeColumn("name", "TEXT"),
		),
	)
	objectsTable := newSqlTable(
		".objects",
		withSqlTableColumns(
			newSqlAttributeColumn("objectID", "INTEGER", withSqlColumnPrimary(true), withSqlColumnAuto(true)),
			newSqlReferenceColumn(classesTable),
		),
	)
	contentsTable := newSqlTable(
		".contents",
		withSqlTableColumns(
			newSqlReferenceColumn(objectsTable),
		),
	)
	enumsTable := newSqlTable(
		".enums",
		withSqlTableColumns(
			newSqlAttributeColumn("enumID", "INTEGER", withSqlColumnPrimary(true), withSqlColumnAuto(true)),
			newSqlAttributeColumn("value", "TEXT"),
		),
	)
	s := &sqlSchema{
		packagesTable:  packagesTable,
		classesTable:   classesTable,
		objectsTable:   objectsTable,
		contentsTable:  contentsTable,
		enumsTable:     enumsTable,
		classSchemaMap: map[EClass]*sqlClassSchema{},
	}
	for _, opt := range options {
		opt.apply(s)
	}
	return s
}

func (s *sqlSchema) getClassSchema(eClass EClass) (*sqlClassSchema, error) {
	classSchema := s.classSchemaMap[eClass]
	if classSchema == nil {
		// create table descriptor
		classTable := newSqlTable(strings.ToLower(eClass.GetName()))
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
				referenceSchema, err := s.getClassSchema(eReference.GetEReferenceType())
				if err != nil {
					return nil, err
				}
				newFeatureReferenceColumn(featureSchema, eFeature, referenceSchema.table)
			case sfkObjectReference:
				newFeatureAttributeColumn(featureSchema, eFeature, "TEXT")
			case sfkObjectList:
				// internal reference
				eReference := eFeature.(EReference)
				referenceSchema, err := s.getClassSchema(eReference.GetEReferenceType())
				if err != nil {
					return nil, err
				}
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
	return classSchema, nil
}

func sqlEscapeIdentifier(id string) string {
	return "\"" + id + "\""
}
