package ecore

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

func TestSqlDecoder_DecodeResource(t *testing.T) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	// create resource & resourceset
	uri := NewURI("testdata/library.complex.sqlite")
	eResource := NewEResourceImpl()
	eResource.SetURI(uri)
	eResourceSet := NewEResourceSetImpl()
	eResourceSet.GetResources().Add(eResource)
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)

	r, err := os.Open(uri.String())
	require.NoError(t, err)
	defer r.Close()

	sqlDecoder := NewSQLReaderDecoder(r, eResource, nil)
	sqlDecoder.DecodeResource()
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
}

func TestSqlDecoder_DecodeResource_Memory(t *testing.T) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	// create resource & resourceset
	uri := NewURI("testdata/library.complex.sqlite")
	eResource := NewEResourceImpl()
	eResource.SetURI(uri)
	eResourceSet := NewEResourceSetImpl()
	eResourceSet.GetResources().Add(eResource)
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)

	r, err := os.Open(uri.String())
	require.NoError(t, err)
	defer r.Close()

	sqlDecoder := NewSQLReaderDecoder(r, eResource, map[string]any{SQL_OPTION_IN_MEMORY_DATABASE: true})
	sqlDecoder.DecodeResource()
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
}

func TestSqlDecoder_EMaps(t *testing.T) {
	// load package
	ePackage := loadPackage("emap.ecore")
	require.NotNil(t, ePackage)

	// create resource & resourceset
	sqlURI := NewURI("testdata/emap.sqlite")
	sqlResource := NewEResourceImpl()
	sqlResource.SetURI(sqlURI)

	eResourceSet := NewEResourceSetImpl()
	eResourceSet.GetResources().Add(sqlResource)
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)

	sqlReader, err := os.Open(sqlURI.String())
	require.NoError(t, err)
	defer sqlReader.Close()

	sqlDecoder := NewSQLReaderDecoder(sqlReader, sqlResource, nil)
	sqlDecoder.DecodeResource()
	require.True(t, sqlResource.GetErrors().Empty(), diagnosticError(sqlResource.GetErrors()))

}

func TestSqlDecoder_SimpleNoIDs_NoObjectIDManager(t *testing.T) {
	// load package
	ePackage := loadPackage("library.simple.ecore")
	require.NotNil(t, ePackage)

	// create resource & resourceset
	sqlURI := NewURI("testdata/library.simple.sqlite")
	sqlResource := NewEResourceImpl()
	sqlResource.SetURI(sqlURI)

	eResourceSet := NewEResourceSetImpl()
	eResourceSet.GetResources().Add(sqlResource)
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)

	sqlReader, err := os.Open(sqlURI.String())
	require.NoError(t, err)
	defer sqlReader.Close()

	sqlDecoder := NewSQLReaderDecoder(sqlReader, sqlResource, nil)
	sqlDecoder.DecodeResource()
	require.True(t, sqlResource.GetErrors().Empty(), diagnosticError(sqlResource.GetErrors()))
}

func TestSqlDecoder_SimpleNoIDs(t *testing.T) {
	// load package
	ePackage := loadPackage("library.simple.ecore")
	require.NotNil(t, ePackage)

	objectIDManager := NewIncrementalIDManager()

	// create resource & resourceset
	sqlURI := NewURI("testdata/library.simple.sqlite")
	sqlResource := NewEResourceImpl()
	sqlResource.SetURI(sqlURI)
	sqlResource.SetObjectIDManager(objectIDManager)

	eResourceSet := NewEResourceSetImpl()
	eResourceSet.GetResources().Add(sqlResource)
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)

	sqlReader, err := os.Open(sqlURI.String())
	require.NoError(t, err)
	defer sqlReader.Close()

	sqlDecoder := NewSQLReaderDecoder(sqlReader, sqlResource, nil)
	sqlDecoder.DecodeResource()
	require.True(t, sqlResource.GetErrors().Empty(), diagnosticError(sqlResource.GetErrors()))

	require.False(t, sqlResource.GetContents().Empty())
	eRoot, _ := sqlResource.GetContents().Get(0).(EObject)
	require.NotNil(t, eRoot)
	require.Equal(t, int64(0), objectIDManager.GetID(eRoot))
}

func TestSqlDecoder_SimpleWithIDs(t *testing.T) {

	// load package
	ePackage := loadPackage("library.simple.ecore")
	require.NotNil(t, ePackage)

	// create resource & resourceset
	objectIDManager := NewIncrementalIDManager()
	sqlURI := NewURI("testdata/library.simple.ids.sqlite")
	sqlResource := NewEResourceImpl()
	sqlResource.SetURI(sqlURI)
	sqlResource.SetObjectIDManager(objectIDManager)

	eResourceSet := NewEResourceSetImpl()
	eResourceSet.GetResources().Add(sqlResource)
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)

	sqlReader, err := os.Open(sqlURI.String())
	require.NoError(t, err)
	defer sqlReader.Close()

	sqlDecoder := NewSQLReaderDecoder(sqlReader, sqlResource, nil)
	sqlDecoder.DecodeResource()
	require.True(t, sqlResource.GetErrors().Empty(), diagnosticError(sqlResource.GetErrors()))

	// check id of the root
	require.False(t, sqlResource.GetContents().Empty())
	eRoot, _ := sqlResource.GetContents().Get(0).(EObject)
	require.NotNil(t, eRoot)
	require.Equal(t, int64(1), objectIDManager.GetID(eRoot))

}

func TestSqlDecoder_SimpleWithULIDs(t *testing.T) {
	// load package
	ePackage := loadPackage("library.simple.ecore")
	require.NotNil(t, ePackage)

	// create resource & resourceset
	sqlURI := NewURI("testdata/library.simple.ulids.sqlite")
	sqlResource := NewEResourceImpl()
	sqlResource.SetURI(sqlURI)

	eResourceSet := NewEResourceSetImpl()
	eResourceSet.GetResources().Add(sqlResource)
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)

	sqlReader, err := os.Open(sqlURI.String())
	require.NoError(t, err)
	defer sqlReader.Close()

	sqlDecoder := NewSQLReaderDecoder(sqlReader, sqlResource, nil)
	sqlDecoder.DecodeResource()
	require.True(t, sqlResource.GetErrors().Empty(), diagnosticError(sqlResource.GetErrors()))
}

func TestSqlDecoder_SimpleWithContainerIDs(t *testing.T) {
	// load package
	ePackage := loadPackage("library.simple.ecore")
	require.NotNil(t, ePackage)

	// create resource & resourceset
	sqlURI := NewURI("testdata/library.container.sqlite")
	sqlResource := NewEResourceImpl()
	sqlResource.SetURI(sqlURI)

	eResourceSet := NewEResourceSetImpl()
	eResourceSet.GetResources().Add(sqlResource)
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)

	sqlReader, err := os.Open(sqlURI.String())
	require.NoError(t, err)
	defer sqlReader.Close()

	sqlDecoder := NewSQLReaderDecoder(sqlReader, sqlResource, nil)
	sqlDecoder.DecodeResource()
	require.True(t, sqlResource.GetErrors().Empty(), diagnosticError(sqlResource.GetErrors()))
}

func TestSqlDecodr_SharedMemoryPool_CreateDB(t *testing.T) {
	fileName := "test.sqlite"
	dbPath := fmt.Sprintf("file:%s?mode=memory&cache=shared", fileName)
	connPool, err := sqlitex.NewPool(dbPath, sqlitex.PoolOptions{Flags: sqlite.OpenCreate | sqlite.OpenReadWrite | sqlite.OpenURI})
	require.NoError(t, err)
	defer connPool.Close()

	// create connection pool
	conn, err := connPool.Take(context.Background())
	require.NoError(t, err)
	defer connPool.Put(conn)

	// create a table
	// Create a table using the first connection
	err = sqlitex.ExecuteTransient(conn, "CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT)", nil)
	require.NoError(t, err)

	// Insert data using the first connection
	err = sqlitex.ExecuteTransient(conn, "INSERT INTO users (name) VALUES ('Alice')", nil)
	require.NoError(t, err)

	conn2, err := connPool.Take(context.Background())
	require.NoError(t, err)
	defer connPool.Put(conn2)

	var id int64
	var name string
	err = sqlitex.ExecuteTransient(
		conn,
		"SELECT id, name FROM users",
		&sqlitex.ExecOptions{
			ResultFunc: func(stmt *sqlite.Stmt) error {
				id = stmt.ColumnInt64(0)
				name = stmt.ColumnText(1)
				return nil
			},
		})
	require.NoError(t, err)
	require.True(t, id > 0)
	require.True(t, name != "")

}

func TestSqlDecodr_SharedMemoryPool_DeserializeDB(t *testing.T) {
	fileName := "test.sqlite"
	dbPath := fmt.Sprintf("file:%s?mode=memory&cache=shared", fileName)
	connPool, err := sqlitex.NewPool(dbPath, sqlitex.PoolOptions{Flags: sqlite.OpenCreate | sqlite.OpenReadWrite | sqlite.OpenURI})
	require.NoError(t, err)
	defer connPool.Close()

	// create connection pool
	conn, err := connPool.Take(context.Background())
	require.NoError(t, err)
	defer connPool.Put(conn)

	f, err := os.Open("testdata/library.complex.sqlite")
	require.NoError(t, err)
	defer f.Close()

	err = sqlitex.ExecuteScript(conn, "ATTACH DATABASE 'file::memory:' AS input;", nil)
	require.NoError(t, err)

	// initialize db with reader bytes
	bytes, err := io.ReadAll(f)
	require.NoError(t, err)

	// deserialize db in input
	err = conn.Deserialize("input", bytes)
	require.NoError(t, err)

	sqls := []string{}
	tables := []string{}
	err = sqlitex.ExecuteTransient(conn, "SELECT sql,name FROM input.sqlite_master WHERE type='table' and name NOT LIKE 'sqlite_%'", &sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			sqls = append(sqls, stmt.ColumnText(0))
			tables = append(tables, stmt.ColumnText(1))
			return nil
		},
	})
	require.NoError(t, err)

	err = sqlitex.ExecuteScript(conn, strings.Join(sqls, ";"), nil)
	require.NoError(t, err)

	for _, name := range tables {
		var query strings.Builder
		query.WriteString(`INSERT INTO "`)
		query.WriteString(name)
		query.WriteString(`" SELECT * FROM "input"."`)
		query.WriteString(name)
		query.WriteString(`"`)
		err = sqlitex.ExecuteTransient(conn, query.String(), nil)
		require.NoError(t, err)
	}

	err = sqlitex.ExecuteScript(conn, "DETACH DATABASE input;", nil)
	require.NoError(t, err)

	conn2, err := connPool.Take(context.Background())
	require.NoError(t, err)
	defer connPool.Put(conn2)

	count := 0
	sqlitex.ExecuteTransient(conn2, "SELECT COUNT(*) FROM '.packages'", &sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			count = stmt.ColumnInt(0)
			return nil
		},
	})
	require.True(t, count > 0)

}
