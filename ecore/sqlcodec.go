package ecore

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	SQL_OPTION_DRIVER            = "DRIVER_NAME"       // value of the sql driver
	SQL_OPTION_ID_ATTRIBUTE_NAME = "ID_ATTRIBUTE_NAME" // value of the id attribute
)

type SQLCodec struct {
}

const sqlCodecVersion int = 1

func sqlTmpDB(prefix string) (string, error) {
	try := 0
	for {
		randBytes := make([]byte, 16)
		rand.Read(randBytes)
		f := filepath.Join(os.TempDir(), prefix+"."+hex.EncodeToString(randBytes)+".sqlite")
		_, err := os.Stat(f)
		if os.IsExist(err) {
			if try++; try < 10000 {
				continue
			}
			return "", &fs.PathError{Op: "sqlTmpDB", Path: prefix, Err: fs.ErrExist}
		}
		return f, nil
	}
}

func (d *SQLCodec) NewEncoder(resource EResource, w io.Writer, options map[string]any) EEncoder {
	return NewSQLEncoder(resource, w, options)
}
func (d *SQLCodec) NewDecoder(resource EResource, r io.Reader, options map[string]any) EDecoder {
	return NewSQLDecoder(resource, r, options)
}

type sqlFeatureKind int

const (
	sfkTransient sqlFeatureKind = iota
	sfkFloat64
	sfkFloat32
	sfkInt
	sfkInt64
	sfkInt32
	sfkInt16
	sfkByte
	sfkBool
	sfkString
	sfkByteArray
	sfkEnum
	sfkDate
	sfkData
	sfkDataList
	sfkObject
	sfkObjectList
	sfkObjectReference
	sfkObjectReferenceList
)

func getSQLCodecFeatureKind(eFeature EStructuralFeature) sqlFeatureKind {
	if eFeature.IsTransient() {
		return sfkTransient
	} else if eReference, _ := eFeature.(EReference); eReference != nil {
		if eReference.IsContainment() {
			if eReference.IsMany() {
				return sfkObjectList
			} else {
				return sfkObject
			}
		}
		opposite := eReference.GetEOpposite()
		if opposite != nil && opposite.IsContainment() {
			return sfkTransient
		}
		if eReference.IsResolveProxies() {
			if eReference.IsMany() {
				return sfkObjectReferenceList
			} else {
				return sfkObjectReference
			}
		}
		if eReference.IsContainer() {
			return sfkTransient
		}
		if eReference.IsMany() {
			return sfkObjectList
		} else {
			return sfkObject
		}
	} else if eAttribute, _ := eFeature.(EAttribute); eAttribute != nil {
		if eAttribute.IsMany() {
			return sfkDataList
		} else {
			eDataType := eAttribute.GetEAttributeType()
			if eEnum, _ := eDataType.(EEnum); eEnum != nil {
				return sfkEnum
			}

			switch eDataType.GetInstanceTypeName() {
			case "float64", "java.lang.Double", "double":
				return sfkFloat64
			case "float32", "java.lang.Float", "float":
				return sfkFloat32
			case "int", "java.lang.Integer":
				return sfkInt
			case "int64", "java.lang.Long", "java.math.BigInteger", "long":
				return sfkInt64
			case "int32":
				return sfkInt32
			case "int16", "java.lang.Short", "short":
				return sfkInt16
			case "byte":
				return sfkByte
			case "bool", "java.lang.Boolean", "boolean":
				return sfkBool
			case "string", "java.lang.String":
				return sfkString
			case "[]byte", "java.util.ByteArray":
				return sfkByteArray
			case "*time.Time", "java.util.Date":
				return sfkDate
			}

			return sfkData
		}
	}
	return -1
}

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

type sqlClassData[FD sqlFeatureData] interface {
	setFeatureData(eFeature EStructuralFeature, featureData FD)
	getTable() *sqlTable
}

type sqlFeatureData interface {
	getFeatureKind() sqlFeatureKind
	getColumn() *sqlColumn
	setColumn(column *sqlColumn)
	getTable() *sqlTable
	setTable(table *sqlTable)
}

// newSqlClassData creates class data struct with its tables/columns structure
func newSqlClassData[CD sqlClassData[FD], FD sqlFeatureData](
	eClass EClass,
	classID int64,
	classDataMap map[EClass]CD,
	getClassData func(eClass EClass) (CD, error),
	newClassData func(eClass EClass, classID int64, classTable *sqlTable, hierarchy []EClass) CD,
	newFeatureData func(eFeature EStructuralFeature) FD) (CD, error) {
	// create data
	ePackage := eClass.GetEPackage()
	eFeatures := eClass.GetEStructuralFeatures()

	// create table descriptor
	classTable := newSqlTable(ePackage.GetNsPrefix() + "_" + strings.ToLower(eClass.GetName()))
	classTable.addColumn(newSqlAttributeColumn(strings.ToLower(eClass.GetName())+"ID", "INTEGER", withSqlColumnPrimary(true)))

	// compute eclass super types
	hierarchy := []EClass{eClass}
	for itClass := eClass.GetEAllSuperTypes().Iterator(); itClass.HasNext(); {
		hierarchy = append(hierarchy, itClass.Next().(EClass))
	}

	// compute table columns and external tables
	classData := newClassData(eClass, classID, classTable, hierarchy)

	// register class data now to handle correctly cycles references
	classDataMap[eClass] = classData

	newFeatureReferenceColumn := func(featureData FD, eFeature EStructuralFeature, table *sqlTable) {
		column := newSqlReferenceColumn(table, withSqlColumnName(eFeature.GetName()))
		classTable.addColumn(column)
		featureData.setColumn(column)
	}

	newFeatureAttributeColumn := func(featureData FD, eFeature EStructuralFeature, columnType string) {
		column := newSqlAttributeColumn(eFeature.GetName(), columnType)
		classTable.addColumn(column)
		featureData.setColumn(column)
	}

	newFeatureTable := func(featureData FD, eFeature EStructuralFeature, columns ...*sqlColumn) {
		featureData.setTable(newSqlTable(
			classTable.name+"_"+eFeature.GetName(),
			columns...,
		))
	}

	for itFeature := eFeatures.Iterator(); itFeature.HasNext(); {
		eFeature := itFeature.Next().(EStructuralFeature)
		// new feature data
		featureData := newFeatureData(eFeature)
		classData.setFeatureData(eFeature, featureData)

		// compute class table columns or children tables
		switch featureData.getFeatureKind() {
		case sfkObject:
			// retrieve object reference type
			eReference := eFeature.(EReference)
			referenceData, err := getClassData(eReference.GetEReferenceType())
			if err != nil {
				var cd CD
				return cd, err
			}
			newFeatureReferenceColumn(featureData, eFeature, referenceData.getTable())
		case sfkObjectReference:
			newFeatureAttributeColumn(featureData, eFeature, "TEXT")
		case sfkObjectList:
			// internal reference
			eReference := eFeature.(EReference)
			referenceData, err := getClassData(eReference.GetEReferenceType())
			if err != nil {
				var cd CD
				return cd, err
			}
			newFeatureTable(featureData, eFeature,
				newSqlReferenceColumn(classTable),
				newSqlAttributeColumn("idx", "REAL"),
				newSqlReferenceColumn(referenceData.getTable(), withSqlColumnName(eFeature.GetName())),
			)
		case sfkObjectReferenceList:
			newFeatureTable(featureData, eFeature,
				newSqlReferenceColumn(classTable),
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
				newSqlReferenceColumn(classTable),
				newSqlAttributeColumn("idx", "REAL"),
				newSqlAttributeColumn(eFeature.GetName(), "TEXT"),
			)
		}
	}

	return classData, nil
}
