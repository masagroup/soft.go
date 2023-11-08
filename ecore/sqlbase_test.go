package ecore

import (
	"database/sql"
	"fmt"
	"maps"
	"reflect"

	"github.com/davecgh/go-spew/spew"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/stretchr/testify/assert"
)

var spewConfig = spew.ConfigState{
	Indent:                  " ",
	DisablePointerAddresses: true,
	DisableCapacities:       true,
	SortKeys:                true,
	DisableMethods:          true,
	MaxDepth:                10,
}

func getDBTables(db *sql.DB) (map[string]struct{}, error) {
	rows, err := db.Query("SELECT name FROM sqlite_schema WHERE type ='table' AND name NOT LIKE 'sqlite_%'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tables := map[string]struct{}{}
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, err
		}
		tables[table] = struct{}{}
	}
	return tables, nil
}

func getDBRows(db *sql.DB, table string) ([][]any, error) {
	rows, err := db.Query("SELECT * FROM " + sqlEscapeIdentifier(table))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// get column type info
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	// used for allocation & dereferencing
	rowValues := make([]reflect.Value, len(columnTypes))
	for i := 0; i < len(columnTypes); i++ {
		// allocate reflect.Value representing a **T value
		rowValues[i] = reflect.New(reflect.PtrTo(columnTypes[i].ScanType()))
	}

	resultList := [][]any{}
	for rows.Next() {
		// initially will hold pointers for Scan, after scanning the
		// pointers will be dereferenced so that the slice holds actual values
		rowResult := make([]any, len(columnTypes))
		for i := 0; i < len(columnTypes); i++ {
			// get the **T value from the reflect.Value
			rowResult[i] = rowValues[i].Interface()
		}

		// scan each column value into the corresponding **T value
		if err := rows.Scan(rowResult...); err != nil {
			return nil, err
		}

		// dereference pointers
		for i := 0; i < len(rowValues); i++ {
			// first pointer deref to get reflect.Value representing a *T value,
			// if rv.IsNil it means column value was NULL
			if rv := rowValues[i].Elem(); rv.IsNil() {
				rowResult[i] = nil
			} else {
				// second deref to get reflect.Value representing the T value
				// and call Interface to get T value from the reflect.Value
				rowResult[i] = rv.Elem().Interface()
			}
		}

		resultList = append(resultList, rowResult)

	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return resultList, nil
}

func assertEqualDB(t assert.TestingT, expected, actual *sql.DB) bool {
	te, err := getDBTables(expected)
	if err != nil {
		return assert.Fail(t, err.Error())
	}
	ta, err := getDBTables(actual)
	if err != nil {
		return assert.Fail(t, err.Error())
	}
	if !maps.Equal(te, ta) {
		e := spewConfig.Sdump(te)
		a := spewConfig.Sdump(ta)
		diff, _ := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
			A:        difflib.SplitLines(e),
			B:        difflib.SplitLines(a),
			FromFile: "Expected",
			FromDate: "",
			ToFile:   "Actual",
			ToDate:   "",
			Context:  1,
		})
		return assert.Fail(t, diff)
	}
	for table := range te {
		e, err := getDBRows(expected, table)
		if err != nil {
			return assert.Fail(t, err.Error())
		}
		a, err := getDBRows(actual, table)
		if err != nil {
			return assert.Fail(t, err.Error())
		}
		if !assert.Equal(t, e, a, fmt.Sprintf("rows for tables '%s' are different", table)) {
			return false
		}
	}
	return true
}
