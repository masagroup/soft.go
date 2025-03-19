package ecore

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

func TestSQLDecoder_DecodeResource(t *testing.T) {
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

func TestSQLDecoder_DecodeResource_Memory(t *testing.T) {
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

func BenchmarkSQLDecoder_Complex(b *testing.B) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(b, ePackage)

	// resource
	uri := NewURI("testdata/library.complex.sqlite")
	eResource := NewEResourceImpl()
	eResource.SetURI(uri)
	eResourceSet := NewEResourceSetImpl()
	eResourceSet.GetResources().Add(eResource)
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)

	r, err := os.Open(uri.String())
	require.NoError(b, err)
	defer r.Close()

	// initialize db with reader bytes
	array, err := io.ReadAll(r)
	require.NoError(b, err)

	for n := 0; n < b.N; n++ {
		r := bytes.NewReader(array)
		sqlDecoder := NewSQLReaderDecoder(r, eResource, nil)
		sqlDecoder.DecodeResource()
		require.True(b, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	}
}

func BenchmarkSQLDecoder_Complex_Memory(b *testing.B) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(b, ePackage)

	// resource
	uri := NewURI("testdata/library.complex.sqlite")
	eResource := NewEResourceImpl()
	eResource.SetURI(uri)
	eResourceSet := NewEResourceSetImpl()
	eResourceSet.GetResources().Add(eResource)
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)

	r, err := os.Open(uri.String())
	require.NoError(b, err)
	defer r.Close()

	// initialize db with reader bytes
	array, err := io.ReadAll(r)
	require.NoError(b, err)

	for n := 0; n < b.N; n++ {
		r := bytes.NewReader(array)
		sqlDecoder := NewSQLReaderDecoder(r, eResource, map[string]any{SQL_OPTION_IN_MEMORY_DATABASE: true})
		sqlDecoder.DecodeResource()
		require.True(b, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	}
}

func TestSQLDecoder_EMaps(t *testing.T) {
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

func TestSQLDecoder_SimpleNoIDs_NoObjectIDManager(t *testing.T) {
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

func TestSQLDecoder_SimpleNoIDs(t *testing.T) {
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

func TestSQLDecoder_SimpleWithIDs(t *testing.T) {

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

func TestSQLDecoder_SimpleWithULIDs(t *testing.T) {
	// load package
	ePackage := loadPackage("library.simple.ecore")
	require.NotNil(t, ePackage)

	// create resource & resourceset
	sqlURI := NewURI("testdata/library.simple.ulids.sqlite")
	sqlResource := NewEResourceImpl()
	sqlResource.SetURI(sqlURI)
	sqlResource.SetObjectIDManager(NewULIDManager())

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

func TestSQLDecoder_SimpleWithContainerIDs(t *testing.T) {
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

func TestSQLDecoder_SharedMemoryPool_CreateDB(t *testing.T) {
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

func TestSQLDecoder_SharedMemoryPool_DeserializeDB_NotWorking(t *testing.T) {
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

	// initialize db with reader bytes
	bytes, err := io.ReadAll(f)
	require.NoError(t, err)

	// set journal mode as rolling back
	bytes[18] = 0x01
	bytes[19] = 0x01

	// deserialize db in input
	err = conn.Deserialize("main", bytes)
	require.NoError(t, err)

	conn2, err := connPool.Take(context.Background())
	require.NoError(t, err)
	defer connPool.Put(conn2)

	err = sqlitex.ExecuteTransient(conn2, "SELECT COUNT(*) FROM '.packages'", &sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			return nil
		},
	})
	require.Error(t, err)
}

func TestSQLDecoder_SharedMemoryPool_DeserializeDB(t *testing.T) {
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

	// initialize db with reader bytes
	bytes, err := io.ReadAll(f)
	require.NoError(t, err)

	// set journal mode as rolling back
	bytes[18] = 0x01
	bytes[19] = 0x01

	input, err := sqlite.OpenConn("file::memory:", sqlite.OpenCreate|sqlite.OpenReadWrite|sqlite.OpenURI)
	require.NoError(t, err)
	defer input.Close()

	// deserialize db in input
	err = input.Deserialize("main", bytes)
	require.NoError(t, err)

	// back to conn
	back, err := sqlite.NewBackup(conn, "main", input, "main")
	require.NoError(t, err)
	done, err := back.Step(-1)
	require.NoError(t, err)
	require.False(t, done)
	err = back.Close()
	require.NoError(t, err)

	conn2, err := connPool.Take(context.Background())
	require.NoError(t, err)
	defer connPool.Put(conn2)

	count := 0
	err = sqlitex.ExecuteTransient(conn2, "SELECT COUNT(*) FROM '.packages'", &sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			count = stmt.ColumnInt(0)
			return nil
		},
	})
	require.NoError(t, err)
	require.True(t, count > 0)
}

func TestSQLDecoder_AllTypes(t *testing.T) {
	// load package
	ePackage := loadPackage("alltypes.ecore")
	require.NotNil(t, ePackage)

	// create resource & resourceset
	sqlURI := NewURI("testdata/alltypes.sqlite")
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

func TestSQLDecoder_InvalidVersion(t *testing.T) {
	// load package
	ePackage := loadPackage("alltypes.ecore")
	require.NotNil(t, ePackage)

	// create resource & resourceset
	sqlURI := NewURI("testdata/alltypes.sqlite")
	sqlResource := NewEResourceImpl()
	sqlResource.SetURI(sqlURI)

	eResourceSet := NewEResourceSetImpl()
	eResourceSet.GetResources().Add(sqlResource)
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)

	sqlReader, err := os.Open(sqlURI.String())
	require.NoError(t, err)
	defer sqlReader.Close()

	sqlDecoder := NewSQLReaderDecoder(sqlReader, sqlResource, map[string]any{SQL_OPTION_CODEC_VERSION: 2})
	sqlDecoder.DecodeResource()
	require.False(t, sqlResource.GetErrors().Empty(), diagnosticError(sqlResource.GetErrors()))
}
