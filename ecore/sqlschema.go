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

func (t *sqlTable) selectAllQuery() string {
	var selectQuery strings.Builder
	selectQuery.WriteString("SELECT * from ")
	selectQuery.WriteString(t.name)
	return selectQuery.String()
}

type sqlClassSchema struct {
	table    *sqlTable
	features map[EStructuralFeature]*sqlFeatureSchema
}

type sqlFeatureSchema struct {
	featureKind sqlFeatureKind
	column      *sqlColumn
	table       *sqlTable
}

type sqlSchema struct {
	packagesTable  *sqlTable
	classesTable   *sqlTable
	objectsTable   *sqlTable
	contentsTable  *sqlTable
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
	s := &sqlSchema{
		packagesTable:  packagesTable,
		classesTable:   classesTable,
		objectsTable:   objectsTable,
		contentsTable:  contentsTable,
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
		// create data
		ePackage := eClass.GetEPackage()
		eFeatures := eClass.GetEStructuralFeatures()

		// create table descriptor
		classTable := newSqlTable(ePackage.GetNsPrefix() + "_" + strings.ToLower(eClass.GetName()))
		classTable.addColumn(newSqlAttributeColumn(strings.ToLower(eClass.GetName())+"ID", "INTEGER", withSqlColumnPrimary(true)))

		// compute table columns and external tables
		classSchema = &sqlClassSchema{
			table:    classTable,
			features: map[EStructuralFeature]*sqlFeatureSchema{},
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
			featureSchema.table = newSqlTable(
				classTable.name+"_"+eFeature.GetName(),
				columns...,
			)
		}

		for itFeature := eFeatures.Iterator(); itFeature.HasNext(); {
			eFeature := itFeature.Next().(EStructuralFeature)
			// new feature data
			featureSchema := &sqlFeatureSchema{
				featureKind: getSQLCodecFeatureKind(eFeature),
			}
			classSchema.features[eFeature] = featureSchema

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
			case sfkBool, sfkByte, sfkInt, sfkInt16, sfkInt32, sfkInt64, sfkEnum:
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
