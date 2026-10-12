package ecore

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

// tests of the writer of the store : one goroutine, one transaction per batch

// newWriterTestStore creates a store on a copy of testdata/library.store.sqlite,
// with the Lendable class and its copies feature
func newWriterTestStore(t *testing.T, options map[string]any) (s *SQLStore, dbPath string, eClass EClass, eFeature EStructuralFeature) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)
	eClass, _ = ePackage.GetEClassifier("Lendable").(EClass)
	require.NotNil(t, eClass)
	eFeature = eClass.GetEStructuralFeatureFromName("copies")
	require.NotNil(t, eFeature)

	dbPath = filepath.Join(t.TempDir(), "library.store.sqlite")
	require.NoError(t, copyFile("testdata/library.store.sqlite", dbPath))

	s, err := NewSQLStore(dbPath, NewURI(""), nil, nil, options)
	require.NoError(t, err)
	require.NotNil(t, s)
	return
}

// newLendable creates a mock of the lendable whose id is objectID in the test database
func newLendable(t *testing.T, eClass EClass, objectID int64) *MockSQLObject {
	mockObject := NewMockSQLObject(t)
	mockObject.EXPECT().GetSQLID().Return(objectID).Maybe()
	mockObject.EXPECT().SetSQLID(objectID).Return().Maybe()
	mockObject.EXPECT().EClass().Return(eClass).Maybe()
	return mockObject
}

// blockWriter gives the writer an operation that blocks until the returned function is called :
// the operations enqueued meanwhile stay pending
func blockWriter(t *testing.T, s *SQLStore) (release func()) {
	blocking := newLendable(t, nil, -1)
	started := make(chan struct{})
	releaseCh := make(chan struct{})
	op := newOperation("block", operationWrite, blocking, nil, -1, nil, func() (any, error) {
		close(started)
		<-releaseCh
		return nil, nil
	})
	s.enqueue(op)
	<-started
	return func() { close(releaseCh) }
}

func countCopies(t *testing.T, dbPath string, lendableID int64) int {
	conn, err := sqlite.OpenConn(dbPath, sqlite.OpenReadOnly)
	require.NoError(t, err)
	defer conn.Close()
	copies := -1
	require.NoError(t, sqlitex.ExecuteTransient(conn, "SELECT copies FROM Lendable WHERE lendableID=?", &sqlitex.ExecOptions{
		Args: []any{lendableID},
		ResultFunc: func(stmt *sqlite.Stmt) error {
			copies = stmt.ColumnInt(0)
			return nil
		},
	}))
	return copies
}

func TestSQLStore_Writer_GroupCommit(t *testing.T) {
	s, dbPath, eClass, eFeature := newWriterTestStore(t, nil)
	defer s.Close()
	lendable := newLendable(t, eClass, 6)

	// the writer is blocked in a transaction : the writes pile up in the queue
	release := blockWriter(t, s)
	require.True(t, s.isTransactionOpen())
	for i := range 100 {
		s.Set(lendable, eFeature, -1, 10+i, false)
	}
	s.mutexWriter.Lock()
	queued := len(s.queue)
	s.mutexWriter.Unlock()
	require.Equal(t, 100, queued)

	// nothing is committed yet
	require.Equal(t, 4, countCopies(t, dbPath, 6))

	// released : the writes are executed and committed in one transaction
	release()
	require.NoError(t, s.WaitOperations(context.Background(), nil))
	require.False(t, s.isTransactionOpen())
	require.Equal(t, 109, countCopies(t, dbPath, 6))
	require.Equal(t, 109, s.Get(lendable, eFeature, -1))
}

func TestSQLStore_Writer_ReadPendingWrite(t *testing.T) {
	s, _, eClass, eFeature := newWriterTestStore(t, nil)
	defer s.Close()
	lendable := newLendable(t, eClass, 6)

	release := blockWriter(t, s)
	defer release()

	// a pending Set is read without waiting for the writer
	s.Set(lendable, eFeature, -1, 9, false)
	require.Equal(t, 9, s.Get(lendable, eFeature, -1))
	require.True(t, s.IsSet(lendable, eFeature))

	// an UnSet too
	go s.UnSet(lendable, eFeature)
	require.Eventually(t, func() bool {
		return s.pendingOperation(lendable, eFeature, -1) != nil && s.pendingOperation(lendable, eFeature, -1).cmd == "UnSet"
	}, time.Second, time.Millisecond)
	require.Nil(t, s.Get(lendable, eFeature, -1))
	require.False(t, s.IsSet(lendable, eFeature))
}

func TestSQLStore_Writer_ReadDirtyObjectOrdered(t *testing.T) {
	s, _, eClass, eFeature := newWriterTestStore(t, nil)
	defer s.Close()
	lendable := newLendable(t, eClass, 6)

	release := blockWriter(t, s)
	s.Add(lendable, eFeature, 0, nil) // an invalid write on the object : it is dirty, the read waits for the writer
	got := make(chan any, 1)
	go func() { got <- s.Get(lendable, eFeature, -1) }()
	select {
	case <-got:
		require.Fail(t, "the read of a dirty object was executed before the pending write")
	case <-time.After(50 * time.Millisecond):
	}
	release()
	select {
	case v := <-got:
		require.Equal(t, 4, v)
	case <-time.After(5 * time.Second):
		require.Fail(t, "the read was not executed")
	}
}

func TestSQLStore_Writer_ReadCleanObjectNotBlocked(t *testing.T) {
	s, _, eClass, eFeature := newWriterTestStore(t, nil)
	defer s.Close()
	dirty := newLendable(t, eClass, 6)
	clean := newLendable(t, eClass, 7)

	release := blockWriter(t, s)
	defer release()
	s.Set(dirty, eFeature, -1, 9, false)

	// the read of an object untouched by the pending writes runs on a pool connection
	start := time.Now()
	require.Equal(t, 3, s.Get(clean, eFeature, -1))
	require.Less(t, time.Since(start), time.Second)
}

func TestSQLStore_Writer_MemoryDatabase_ReadAfterWrite(t *testing.T) {
	s, dbPath, eClass, eFeature := newWriterTestStore(t, map[string]any{SQL_OPTION_IN_MEMORY_DATABASE: true})
	lendable := newLendable(t, eClass, 6)

	s.Set(lendable, eFeature, -1, 9, false)
	require.Equal(t, 9, s.Get(lendable, eFeature, -1))
	require.NoError(t, s.WaitOperations(context.Background(), nil))
	require.Equal(t, 9, s.Get(lendable, eFeature, -1))

	// the memory database is written to the file at close
	require.NoError(t, s.Close())
	require.Equal(t, 9, countCopies(t, dbPath, 6))
}

func TestSQLStore_Writer_ClosePendingOperations(t *testing.T) {
	s, dbPath, eClass, eFeature := newWriterTestStore(t, nil)
	lendable := newLendable(t, eClass, 6)

	release := blockWriter(t, s)
	for i := range 50 {
		s.Set(lendable, eFeature, -1, 100+i, false)
	}
	go func() {
		time.Sleep(20 * time.Millisecond)
		release()
	}()
	require.NoError(t, s.Close())
	require.Equal(t, 149, countCopies(t, dbPath, 6))

	// operations after close are rejected
	require.Error(t, s.ExecuteQuery(context.Background(), "SELECT 1", nil))
}

func TestSQLStore_Writer_ErrorNoPoison(t *testing.T) {
	s, _, eClass, eFeature := newWriterTestStore(t, nil)
	defer s.Close()
	lendable := newLendable(t, eClass, 6)

	// a failing write
	require.Error(t, s.ExecuteQuery(context.Background(), "INSERT INTO unknown VALUES (1)", nil))

	// the next operations work, the error is reported once
	s.Set(lendable, eFeature, -1, 9, false)
	require.NoError(t, s.ExecuteQuery(context.Background(), "SELECT 1", nil))
	require.Equal(t, 9, s.Get(lendable, eFeature, -1))
	require.Error(t, s.WaitOperations(context.Background(), nil))
	require.NoError(t, s.WaitOperations(context.Background(), nil))
	require.Equal(t, 9, s.Get(lendable, eFeature, -1))
}

func TestSQLStore_Writer_FlushReentrant(t *testing.T) {
	ePackage := loadPackage("library.datalist.ecore")
	require.NotNil(t, ePackage)
	eClass, _ := ePackage.GetEClassifier("Book").(EClass)
	require.NotNil(t, eClass)
	eFeature := eClass.GetEStructuralFeatureFromName("contents")
	require.NotNil(t, eFeature)

	dbPath := filepath.Join(t.TempDir(), "library.datalist.sqlite")
	require.NoError(t, copyFile("testdata/library.datalist.sqlite", dbPath))
	s, err := NewSQLStore(dbPath, NewURI(""), nil, NewMockEPackageRegistry(t), nil)
	require.NoError(t, err)
	defer s.Close()

	book := NewMockSQLObject(t)
	book.EXPECT().GetSQLID().Return(int64(5)).Maybe()
	book.EXPECT().SetSQLID(int64(5)).Return().Maybe()
	book.EXPECT().EClass().Return(eClass).Maybe()

	// a write makes the book dirty : its reads are executed by the writer
	size := s.Size(book, eFeature)
	s.Add(book, eFeature, size, "last")

	// the writer iterates the contents; the iteration calls the store back :
	// a nested iteration of the same feature, and a flush
	outer := 0
	for range s.All(book, eFeature) {
		outer++
		inner := 0
		for range s.All(book, eFeature) {
			inner++
		}
		require.Equal(t, size+1, inner)
		require.NoError(t, s.WaitOperations(context.Background(), nil))
	}
	require.Equal(t, size+1, outer)
}

func TestSQLStore_Writer_SerializeCancelledWhileQueued(t *testing.T) {
	s, _, _, _ := newWriterTestStore(t, nil)
	defer s.Close()

	release := blockWriter(t, s)
	s.AddRoot(nil) // a global write : the serialization waits for the writer
	ctx, cancel := context.WithCancel(context.Background())
	p := s.Serialize(ctx)
	cancel()
	_, err := p.Await(context.Background())
	require.ErrorIs(t, err, context.Canceled)
	release()
	require.Error(t, s.WaitOperations(context.Background(), nil)) // the AddRoot of nil
}

func TestSQLStore_Writer_SerializeCommitsPending(t *testing.T) {
	s, _, eClass, eFeature := newWriterTestStore(t, nil)
	defer s.Close()
	lendable := newLendable(t, eClass, 6)

	s.Set(lendable, eFeature, -1, 9, false)
	bytes, err := s.Serialize(context.Background()).Await(context.Background())
	require.NoError(t, err)

	// the serialized database holds the write
	serializedPath := filepath.Join(t.TempDir(), "serialized.sqlite")
	require.NoError(t, os.WriteFile(serializedPath, *bytes, 0644))
	require.Equal(t, 9, countCopies(t, serializedPath, 6))
}

func TestSQLStore_Writer_ContainsUnknownObject(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)
	libraryClass, _ := ePackage.GetEClassifier("Library").(EClass)
	require.NotNil(t, libraryClass)
	booksFeature := libraryClass.GetEStructuralFeatureFromName("books")
	require.NotNil(t, booksFeature)
	bookClass, _ := ePackage.GetEClassifier("Book").(EClass)
	require.NotNil(t, bookClass)

	dbPath := filepath.Join(t.TempDir(), "library.complex.sqlite")
	require.NoError(t, copyFile("testdata/library.complex.sqlite", dbPath))
	s, err := NewSQLStore(dbPath, NewURI(""), nil, nil, nil)
	require.NoError(t, err)
	defer s.Close()

	countObjects := func() (count int) {
		require.NoError(t, s.ExecuteQuery(context.Background(), `SELECT COUNT(*) FROM ".objects"`, &sqlitex.ExecOptions{
			ResultFunc: func(stmt *sqlite.Stmt) error {
				count = stmt.ColumnInt(0)
				return nil
			}}))
		return
	}
	objects := countObjects()

	// an object unknown to the store is not in the list, and looking for it does not write it
	library := newLendable(t, libraryClass, 2)
	unknown := NewMockSQLObject(t)
	unknown.EXPECT().GetSQLID().Return(int64(0)).Maybe()
	unknown.EXPECT().EClass().Return(bookClass).Maybe()
	require.False(t, s.Contains(library, booksFeature, unknown))
	require.Equal(t, -1, s.IndexOf(library, booksFeature, unknown))
	require.Equal(t, -1, s.LastIndexOf(library, booksFeature, unknown))
	require.NoError(t, s.WaitOperations(context.Background(), nil))
	require.Equal(t, objects, countObjects())

	// a known one is found
	book := newLendable(t, bookClass, 6)
	require.True(t, s.Contains(library, booksFeature, book))
	require.GreaterOrEqual(t, s.IndexOf(library, booksFeature, book), 0)
}

func TestSQLStore_Writer_ParallelWritersReaders(t *testing.T) {
	s, dbPath, eClass, eFeature := newWriterTestStore(t, nil)
	defer s.Close()
	lendable := newLendable(t, eClass, 6)

	wg := sync.WaitGroup{}
	for i := range 50 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Set(lendable, eFeature, -1, i, i%2 == 0)
			v := s.Get(lendable, eFeature, -1)
			assert.IsType(t, 0, v)
		}()
	}
	wg.Wait()
	require.NoError(t, s.WaitOperations(context.Background(), nil))
	require.Equal(t, s.Get(lendable, eFeature, -1), countCopies(t, dbPath, 6))
}
