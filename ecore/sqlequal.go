package ecore

import (
	"fmt"
	"maps"

	"github.com/davecgh/go-spew/spew"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/stretchr/testify/require"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

var spewConfig = spew.ConfigState{
	Indent:                  " ",
	DisablePointerAddresses: true,
	DisableCapacities:       true,
	SortKeys:                true,
	DisableMethods:          true,
	MaxDepth:                10,
}

func getDBTables(conn *sqlite.Conn) (map[string]struct{}, error) {
	tables := map[string]struct{}{}
	if err := sqlitex.ExecuteTransient(
		conn,
		"SELECT name FROM sqlite_schema WHERE type ='table' AND name NOT LIKE 'sqlite_%'",
		&sqlitex.ExecOptions{
			ResultFunc: func(stmt *sqlite.Stmt) error {
				tables[stmt.ColumnText(0)] = struct{}{}
				return nil
			},
		}); err != nil {
		return nil, err
	}
	return tables, nil
}

func getDBRows(conn *sqlite.Conn, table string) ([][]any, error) {
	resultList := [][]any{}
	if err := sqlitex.ExecuteTransient(conn, "SELECT * FROM '"+table+"'", &sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			results := make([]any, stmt.ColumnCount())
			for i := 0; i < stmt.ColumnCount(); i++ {
				results[i] = decodeAny(stmt, i)
			}
			resultList = append(resultList, results)
			return nil
		},
	}); err != nil {
		return nil, err
	}
	return resultList, nil
}

func RequireEqualDB(t require.TestingT, expected, actual *sqlite.Conn) {
	te, err := getDBTables(expected)
	if err != nil {
		require.Fail(t, err.Error())
	}
	ta, err := getDBTables(actual)
	if err != nil {
		require.Fail(t, err.Error())
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
		require.Fail(t, diff)
	}
	for table := range te {
		e, err := getDBRows(expected, table)
		if err != nil {
			require.Fail(t, err.Error())
		}
		a, err := getDBRows(actual, table)
		if err != nil {
			require.Fail(t, err.Error())
		}
		require.Equal(t, e, a, fmt.Sprintf("rows for tables '%s' are different", table))
	}
}
