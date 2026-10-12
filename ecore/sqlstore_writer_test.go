package ecore

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
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
	copies, err := committedCopies(dbPath, lendableID)
	require.NoError(t, err)
	return copies
}

// committedCopies reads the copies of a lendable committed in the database file, from any goroutine
func committedCopies(dbPath string, lendableID int64) (int, error) {
	conn, err := sqlite.OpenConn(dbPath, sqlite.OpenReadOnly)
	if err != nil {
		return -1, err
	}
	defer conn.Close()
	copies := -1
	err = sqlitex.ExecuteTransient(conn, "SELECT copies FROM Lendable WHERE lendableID=?", &sqlitex.ExecOptions{
		Args: []any{lendableID},
		ResultFunc: func(stmt *sqlite.Stmt) error {
			copies = stmt.ColumnInt(0)
			return nil
		},
	})
	return copies, err
}

// onWriter gives the writer an operation executing fn, on an object of its own; the returned
// channel is closed once it is executed
func onWriter(t *testing.T, s *SQLStore, fn func()) <-chan struct{} {
	op := newOperation("test", operationWrite, newLendable(t, nil, -1), nil, -1, nil, func() (any, error) {
		fn()
		return nil, nil
	})
	op.done = make(chan struct{})
	s.enqueue(op)
	return op.done
}

// queueLen is the number of operations waiting for the writer
func queueLen(s *SQLStore) int {
	s.mutexWriter.Lock()
	defer s.mutexWriter.Unlock()
	return len(s.queue)
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

func TestSQLStore_Writer_CommitMaxOperations(t *testing.T) {
	s, dbPath, eClass, eFeature := newWriterTestStore(t, map[string]any{SQL_OPTION_COMMIT_MAX_OPERATIONS: 10})
	defer s.Close()
	lendable6 := newLendable(t, eClass, 6)
	lendable7 := newLendable(t, eClass, 7)

	// the blocking operation and 9 Set are 10 operations : the writer commits after them,
	// although its queue is not empty
	release := blockWriter(t, s)
	for i := range 9 {
		s.Set(lendable6, eFeature, -1, 10+i, false)
	}
	type state struct {
		copies     int
		err        error
		clean      bool
		stillDirty bool
	}
	got := make(chan state, 1)
	onWriter(t, s, func() {
		copies, err := committedCopies(dbPath, 6)
		got <- state{
			copies: copies,
			err:    err,
			// the object written before the commit is clean, the one of a still queued operation is not
			clean:      s.canReadConcurrently(newOperation("Get", operationRead, lendable6, eFeature, -1, nil, nil)),
			stillDirty: !s.canReadConcurrently(newOperation("Get", operationRead, lendable7, eFeature, -1, nil, nil)),
		}
	})
	s.Set(lendable7, eFeature, -1, 20, false)
	release()

	st := <-got
	require.NoError(t, st.err)
	require.Equal(t, 18, st.copies)
	require.True(t, st.clean)
	require.True(t, st.stillDirty)
	require.NoError(t, s.WaitOperations(context.Background(), nil))
	require.Equal(t, 20, countCopies(t, dbPath, 7))
}

func TestSQLStore_Writer_CommitMaxDuration(t *testing.T) {
	s, dbPath, eClass, eFeature := newWriterTestStore(t, map[string]any{SQL_OPTION_COMMIT_MAX_DURATION: 50 * time.Millisecond})
	defer s.Close()
	lendable := newLendable(t, eClass, 6)

	// a transaction open longer than the duration is committed, although the queue is not empty
	release := blockWriter(t, s)
	s.Set(lendable, eFeature, -1, 9, false)
	onWriter(t, s, func() { time.Sleep(100 * time.Millisecond) })
	type result struct {
		copies int
		err    error
	}
	got := make(chan result, 1)
	onWriter(t, s, func() {
		copies, err := committedCopies(dbPath, 6)
		got <- result{copies, err}
	})
	release()

	r := <-got
	require.NoError(t, r.err)
	require.Equal(t, 9, r.copies)
}

func TestSQLStore_Writer_ErrorRollsTransactionBack(t *testing.T) {
	s, dbPath, eClass, eFeature := newWriterTestStore(t, nil)
	defer s.Close()
	lendable6 := newLendable(t, eClass, 6)
	lendable7 := newLendable(t, eClass, 7)

	// a write, then an error rolling the transaction back, then a write
	release := blockWriter(t, s)
	s.Set(lendable6, eFeature, -1, 9, false)
	errCh := make(chan error, 1)
	go func() {
		errCh <- s.ExecuteQuery(context.Background(), `INSERT OR ROLLBACK INTO ".objects" (objectID, classID) VALUES (1, 1)`, nil)
	}()
	require.Eventually(t, func() bool { return queueLen(s) == 2 }, time.Second, time.Millisecond)
	s.Set(lendable7, eFeature, -1, 8, false)
	release()

	require.ErrorContains(t, <-errCh, "constraint")
	require.Error(t, s.WaitOperations(context.Background(), nil))
	require.False(t, s.isTransactionOpen())

	// the write before the error is lost with the transaction
	require.Equal(t, 4, countCopies(t, dbPath, 6))
	require.Equal(t, 4, s.Get(lendable6, eFeature, -1))
	// the writer began a new transaction for the next write
	require.Equal(t, 8, countCopies(t, dbPath, 7))
}

func TestSQLStore_Writer_BeginError(t *testing.T) {
	s, dbPath, eClass, eFeature := newWriterTestStore(t, nil)
	defer s.Close()
	lendable := newLendable(t, eClass, 6)

	// the writer fails at once when the database is locked
	<-onWriter(t, s, func() { s.writerConn.SetBusyTimeout(0) })
	require.NoError(t, s.WaitOperations(context.Background(), nil))

	// another connection holds the write lock : the writer can't begin its transaction
	conn, err := sqlite.OpenConn(dbPath, sqlite.OpenReadWrite|sqlite.OpenWAL)
	require.NoError(t, err)
	defer conn.Close()
	require.NoError(t, sqlitex.ExecuteTransient(conn, "BEGIN IMMEDIATE;", nil))
	s.Set(lendable, eFeature, -1, 9, false)
	require.ErrorContains(t, s.WaitOperations(context.Background(), nil), "locked")
	require.False(t, s.isTransactionOpen())

	// the failed write is not pending anymore
	require.NoError(t, sqlitex.ExecuteTransient(conn, "ROLLBACK;", nil))
	require.Equal(t, 4, s.Get(lendable, eFeature, -1))

	// the next write begins a transaction
	s.Set(lendable, eFeature, -1, 10, false)
	require.NoError(t, s.WaitOperations(context.Background(), nil))
	require.Equal(t, 10, countCopies(t, dbPath, 6))
}

func TestSQLStore_Writer_CommitError(t *testing.T) {
	s, dbPath, eClass, eFeature := newWriterTestStore(t, nil)
	defer s.Close()
	lendable := newLendable(t, eClass, 6)

	// foreign keys are enabled outside of a transaction : a flush operation has none
	_, err := s.execute(newOperation("foreign keys", operationWrite, nil, nil, -1, nil, func() (any, error) {
		return nil, sqlitex.ExecuteTransient(s.writerConn, "PRAGMA foreign_keys=ON;", nil)
	}, withFlush()), true)
	require.NoError(t, err)
	require.NoError(t, s.ExecuteQuery(context.Background(), "CREATE TABLE parent (id INTEGER PRIMARY KEY)", nil))
	require.NoError(t, s.ExecuteQuery(context.Background(), "CREATE TABLE child (parentID INTEGER REFERENCES parent(id) DEFERRABLE INITIALLY DEFERRED)", nil))
	require.NoError(t, s.WaitOperations(context.Background(), nil))

	// a write, and a deferred violation detected by the commit
	release := blockWriter(t, s)
	s.Set(lendable, eFeature, -1, 9, false)
	errCh := make(chan error, 1)
	go func() {
		errCh <- s.ExecuteQuery(context.Background(), "INSERT INTO child VALUES (42)", nil)
	}()
	require.Eventually(t, func() bool { return queueLen(s) == 2 }, time.Second, time.Millisecond)
	release()
	require.NoError(t, <-errCh)

	// the commit error is reported by the next wait, and the transaction is rolled back
	require.ErrorContains(t, s.WaitOperations(context.Background(), nil), "FOREIGN KEY")
	require.False(t, s.isTransactionOpen())
	require.Equal(t, 4, countCopies(t, dbPath, 6))

	// the next write begins a transaction
	s.Set(lendable, eFeature, -1, 10, false)
	require.NoError(t, s.WaitOperations(context.Background(), nil))
	require.Equal(t, 10, countCopies(t, dbPath, 6))
}

func TestSQLStore_Writer_OperationTimeout(t *testing.T) {
	core, logs := observer.New(zap.WarnLevel)
	s, dbPath, _, _ := newWriterTestStore(t, map[string]any{
		SQL_OPTION_LOGGER:            zap.New(core),
		SQL_OPTION_OPERATION_TIMEOUT: 20 * time.Millisecond,
	})
	defer s.Close()

	// the writer is blocked longer than the timeout : the wait is logged and goes on
	release := blockWriter(t, s)
	time.AfterFunc(100*time.Millisecond, release)
	require.NoError(t, s.ExecuteQuery(context.Background(), "UPDATE lendable SET copies=5 WHERE lendableID=6", nil))
	require.NoError(t, s.WaitOperations(context.Background(), nil))
	require.Equal(t, 5, countCopies(t, dbPath, 6))
	require.Equal(t, 1, logs.FilterMessage("waiting").Len())

	// a context done after the warning stops the wait
	release = blockWriter(t, s)
	defer release()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Millisecond)
	defer cancel()
	require.ErrorIs(t, s.ExecuteQuery(ctx, "UPDATE lendable SET copies=6 WHERE lendableID=6", nil), context.DeadlineExceeded)
	require.Equal(t, 2, logs.FilterMessage("waiting").Len())
}

func TestSQLStore_Writer_ReadDirtyValueOrdered(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)
	libraryClass, _ := ePackage.GetEClassifier("Library").(EClass)
	require.NotNil(t, libraryClass)
	libraryBooks := libraryClass.GetEStructuralFeatureFromName("books")
	writerClass, _ := ePackage.GetEClassifier("Writer").(EClass)
	require.NotNil(t, writerClass)
	writerBooks := writerClass.GetEStructuralFeatureFromName("books")
	bookClass, _ := ePackage.GetEClassifier("Book").(EClass)
	require.NotNil(t, bookClass)

	s, _, _, _ := newWriterTestStore(t, nil)
	defer s.Close()
	library := newLendable(t, libraryClass, 2)
	writer := newLendable(t, writerClass, 3)
	book6 := newLendable(t, bookClass, 6)
	book7 := newLendable(t, bookClass, 7)

	// a pending write with book6 as value : book6 is dirty
	release := blockWriter(t, s)
	s.Add(writer, writerBooks, 2, book6)
	require.False(t, s.canReadConcurrently(newOperation("Contains", operationRead, library, libraryBooks, -1, book6, nil)))
	require.True(t, s.canReadConcurrently(newOperation("Contains", operationRead, library, libraryBooks, -1, book7, nil)))

	// a read of a clean object with the dirty value waits for the writer
	got := make(chan bool, 1)
	go func() { got <- s.Contains(library, libraryBooks, book6) }()
	select {
	case <-got:
		require.Fail(t, "the read of a dirty value was executed before the pending write")
	case <-time.After(50 * time.Millisecond):
	}
	release()
	select {
	case v := <-got:
		require.True(t, v)
	case <-time.After(5 * time.Second):
		require.Fail(t, "the read was not executed")
	}
	require.NoError(t, s.WaitOperations(context.Background(), nil))
	require.Equal(t, 2, s.IndexOf(writer, writerBooks, book6))
}

func TestSQLStore_Writer_MemoryDatabase_ReadsByWriter(t *testing.T) {
	s, _, eClass, eFeature := newWriterTestStore(t, map[string]any{SQL_OPTION_IN_MEMORY_DATABASE: true})
	defer s.Close()
	clean := newLendable(t, eClass, 7)

	// tables are locked in a memory database : even the read of a clean object waits for the writer
	release := blockWriter(t, s)
	got := make(chan any, 1)
	go func() { got <- s.Get(clean, eFeature, -1) }()
	select {
	case <-got:
		require.Fail(t, "the read was executed on a pool connection")
	case <-time.After(50 * time.Millisecond):
	}
	release()
	select {
	case v := <-got:
		require.Equal(t, 3, v)
	case <-time.After(5 * time.Second):
		require.Fail(t, "the read was not executed")
	}
}

// libraryModel holds the classes and features of library.complex.ecore used by the tests
type libraryModel struct {
	ePackage        EPackage
	libraryClass    EClass
	libraryBooks    EStructuralFeature
	writerClass     EClass
	writerBooks     EStructuralFeature
	bookClass       EClass
	bookCategory    EStructuralFeature
	lendableClass   EClass
	copies          EStructuralFeature
	itemClass       EClass
	publicationDate EStructuralFeature
}

func newLibraryModel(t *testing.T) *libraryModel {
	m := &libraryModel{}
	m.ePackage = loadPackage("library.complex.ecore")
	require.NotNil(t, m.ePackage)
	m.libraryClass, _ = m.ePackage.GetEClassifier("Library").(EClass)
	require.NotNil(t, m.libraryClass)
	m.libraryBooks = m.libraryClass.GetEStructuralFeatureFromName("books")
	m.writerClass, _ = m.ePackage.GetEClassifier("Writer").(EClass)
	require.NotNil(t, m.writerClass)
	m.writerBooks = m.writerClass.GetEStructuralFeatureFromName("books")
	m.bookClass, _ = m.ePackage.GetEClassifier("Book").(EClass)
	require.NotNil(t, m.bookClass)
	m.bookCategory = m.bookClass.GetEStructuralFeatureFromName("category")
	// the store reads a feature with the class declaring it
	m.lendableClass, _ = m.ePackage.GetEClassifier("Lendable").(EClass)
	require.NotNil(t, m.lendableClass)
	m.copies = m.lendableClass.GetEStructuralFeatureFromName("copies")
	m.itemClass, _ = m.ePackage.GetEClassifier("Item").(EClass)
	require.NotNil(t, m.itemClass)
	m.publicationDate = m.itemClass.GetEStructuralFeatureFromName("publicationDate")
	require.NotNil(t, m.libraryBooks)
	require.NotNil(t, m.writerBooks)
	require.NotNil(t, m.bookCategory)
	require.NotNil(t, m.copies)
	require.NotNil(t, m.publicationDate)
	return m
}

// containerOf reads the container columns of an object, from any goroutine
func containerOf(s *SQLStore, objectID int64) (container any, feature any, err error) {
	err = s.ExecuteQuery(context.Background(), `SELECT containerID,containerFeatureID FROM ".objects" WHERE objectID=?`, &sqlitex.ExecOptions{
		Args: []any{objectID},
		ResultFunc: func(stmt *sqlite.Stmt) error {
			container = decodeAny(stmt, 0)
			feature = decodeAny(stmt, 1)
			return nil
		}})
	return
}

func TestSQLStore_Writer_SetContainer(t *testing.T) {
	m := newLibraryModel(t)
	s, _, _, _ := newWriterTestStore(t, nil)
	defer s.Close()
	library := newLendable(t, m.libraryClass, 2)
	book := newLendable(t, m.bookClass, 6)
	branches := m.libraryClass.GetEStructuralFeatureFromName("branches")
	branchesID := int64(m.libraryClass.GetFeatureID(branches))

	// a container and its feature
	s.SetContainer(book, library, branches)
	require.NoError(t, s.WaitOperations(context.Background(), nil))
	container, feature, err := containerOf(s, 6)
	require.NoError(t, err)
	require.Equal(t, int64(2), container)
	require.Equal(t, branchesID, feature)

	// no container
	s.SetContainer(book, nil, nil)
	require.NoError(t, s.WaitOperations(context.Background(), nil))
	container, feature, err = containerOf(s, 6)
	require.NoError(t, err)
	require.Nil(t, container)
	require.Nil(t, feature)
	c, f := s.GetContainer(book)
	require.Nil(t, c)
	require.Nil(t, f)

	// from an operation being executed : executed at once, on the writer connection
	type result struct {
		container, feature any
		err                error
	}
	got := make(chan result, 1)
	<-onWriter(t, s, func() {
		s.SetContainer(book, library, m.libraryBooks)
		container, feature, err := containerOf(s, 6)
		got <- result{container, feature, err}
	})
	r := <-got
	require.NoError(t, r.err)
	require.Equal(t, int64(2), r.container)
	require.Equal(t, int64(m.libraryClass.GetFeatureID(m.libraryBooks)), r.feature)
	require.NoError(t, s.WaitOperations(context.Background(), nil))
}

func TestSQLStore_Writer_LookupObjectReferences(t *testing.T) {
	m := newLibraryModel(t)
	s, _, _, _ := newWriterTestStore(t, nil)
	defer s.Close()
	writer := newLendable(t, m.writerClass, 3)

	// an object outside of the store is looked up by its uri
	proxy := func(uri string) EObject {
		o := NewDynamicEObjectImpl()
		o.SetEClass(m.bookClass)
		o.ESetProxyURI(NewURI(uri))
		return o
	}
	book1 := proxy("#//@library/@books.1")
	require.True(t, s.Contains(writer, m.writerBooks, book1))
	require.Equal(t, 1, s.IndexOf(writer, m.writerBooks, book1))
	require.Equal(t, 1, s.LastIndexOf(writer, m.writerBooks, book1))
	unknown := proxy("#//@library/@books.5")
	require.False(t, s.Contains(writer, m.writerBooks, unknown))
	require.Equal(t, -1, s.IndexOf(writer, m.writerBooks, unknown))

	// an object of the store by its id
	book6 := newLendable(t, m.bookClass, 6)
	require.False(t, s.Contains(writer, m.writerBooks, book6))
	s.Add(writer, m.writerBooks, 2, book6)
	require.NoError(t, s.WaitOperations(context.Background(), nil))
	require.True(t, s.Contains(writer, m.writerBooks, book6))
	require.Equal(t, 2, s.IndexOf(writer, m.writerBooks, book6))
}

func TestSQLStore_Writer_LookupData(t *testing.T) {
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
	book := newLendable(t, eClass, 4)

	require.True(t, s.Contains(book, eFeature, "c22"))
	require.Equal(t, 1, s.IndexOf(book, eFeature, "c22"))
	require.Equal(t, 1, s.LastIndexOf(book, eFeature, "c22"))
	require.False(t, s.Contains(book, eFeature, "c99"))
	require.Equal(t, -1, s.IndexOf(book, eFeature, "c99"))
	require.False(t, s.Contains(book, eFeature, nil))
}

func TestSQLStore_Writer_LookupEnumAndDate(t *testing.T) {
	// Contains and IndexOf are for many features : the enum and date encodings of a
	// lookup are tested directly
	m := newLibraryModel(t)
	s, _, _, _ := newWriterTestStore(t, nil)
	defer s.Close()
	book := newLendable(t, m.bookClass, 6)

	categorySchema, err := s.getFeatureSchema(book, m.bookCategory)
	require.NoError(t, err)
	// a literal of the database : looked up in the database, then known
	for range 2 {
		encoded, isEncoded, err := s.encodeLookupValue(categorySchema, 2) // Biography
		require.NoError(t, err)
		require.True(t, isEncoded)
		require.Equal(t, int64(1), encoded)
	}
	// a literal never written is not in the database
	encoded, isEncoded, err := s.encodeLookupValue(categorySchema, 0) // Mystery
	require.NoError(t, err)
	require.False(t, isEncoded)
	require.Nil(t, encoded)

	item := newLendable(t, m.itemClass, 6)
	dateSchema, err := s.getFeatureSchema(item, m.publicationDate)
	require.NoError(t, err)
	date := time.Date(2015, time.September, 6, 4, 24, 46, 0, time.UTC)
	encoded, isEncoded, err = s.encodeLookupValue(dateSchema, &date)
	require.NoError(t, err)
	require.True(t, isEncoded)
	require.Equal(t, "2015-09-06T04:24:46Z", encoded)
}

func TestSQLStore_Writer_UnRegister(t *testing.T) {
	m := newLibraryModel(t)
	s, _, _, _ := newWriterTestStore(t, nil)
	defer s.Close()
	book := newLendable(t, m.lendableClass, 6)
	isRegistered := func() bool {
		_, registered := s.sqlIDManager.GetObjectID(book)
		return registered
	}

	// without its contents
	require.Equal(t, 4, s.Get(book, m.copies, -1))
	require.True(t, isRegistered())
	_, err := s.UnRegister(book, false).Await(context.Background())
	require.NoError(t, err)
	require.False(t, isRegistered())

	// from an operation being executed : at once
	require.Equal(t, 4, s.Get(book, m.copies, -1))
	require.True(t, isRegistered())
	type result struct {
		err        error
		registered bool
	}
	got := make(chan result, 1)
	<-onWriter(t, s, func() {
		_, err := s.UnRegister(book, false).Await(context.Background())
		got <- result{err, isRegistered()}
	})
	r := <-got
	require.NoError(t, r.err)
	require.False(t, r.registered)
}

func TestSQLStore_Writer_AllInterrupted(t *testing.T) {
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
	book := newLendable(t, eClass, 4)

	first := func() any {
		for v := range s.All(book, eFeature) {
			return v
		}
		return nil
	}

	// on a pool connection
	require.Equal(t, "c21", first())
	require.Equal(t, []any{"c21", "c22", "c23"}, s.ToArray(book, eFeature))

	// by the writer : the object is dirty
	release := blockWriter(t, s)
	s.Add(book, eFeature, 3, "c24")
	got := make(chan any, 1)
	go func() { got <- first() }()
	release()
	require.Equal(t, "c21", <-got)
	// the writer connection is usable : the next writes are committed
	s.Add(book, eFeature, 4, "c25")
	require.NoError(t, s.WaitOperations(context.Background(), nil))
	require.Equal(t, []any{"c21", "c22", "c23", "c24", "c25"}, s.ToArray(book, eFeature))
}

func TestSQLStore_Writer_MissingTables(t *testing.T) {
	m := newLibraryModel(t)
	s, _, _, _ := newWriterTestStore(t, nil)
	defer s.Close()
	library := newLendable(t, m.libraryClass, 2)
	book6 := newLendable(t, m.bookClass, 6)
	lendable6 := newLendable(t, m.lendableClass, 6)

	for _, table := range []string{"lendable", "library_books", `".contents"`} {
		require.NoError(t, s.ExecuteQuery(context.Background(), "DROP TABLE "+table, nil))
	}
	require.NoError(t, s.WaitOperations(context.Background(), nil))

	// the reads return zero values
	require.Nil(t, s.Get(lendable6, m.copies, -1))
	require.False(t, s.IsSet(lendable6, m.copies))
	require.Equal(t, 0, s.Size(library, m.libraryBooks))
	require.False(t, s.Contains(library, m.libraryBooks, book6))
	require.Equal(t, -1, s.IndexOf(library, m.libraryBooks, book6))
	require.Empty(t, s.ToArray(library, m.libraryBooks))
	require.Empty(t, s.GetRoots())

	// the writes reading a missing table fail, the errors are reported by the next wait
	s.Set(lendable6, m.copies, -1, 9, false)
	require.ErrorContains(t, s.WaitOperations(context.Background(), nil), "no such table: lendable")
	s.RemoveRoot(library)
	require.ErrorContains(t, s.WaitOperations(context.Background(), nil), "no such table: .contents")
}

func TestSQLStore_Writer_MissingObjectsTable(t *testing.T) {
	m := newLibraryModel(t)
	s, _, _, _ := newWriterTestStore(t, nil)
	defer s.Close()

	require.NoError(t, s.ExecuteQuery(context.Background(), `DROP TABLE ".objects"`, nil))
	require.NoError(t, s.WaitOperations(context.Background(), nil))

	// objects unknown to the store can not be looked up : the reads return zero values
	library := newLendable(t, m.libraryClass, 2)
	book := newLendable(t, m.bookClass, 6)
	lendable := newLendable(t, m.lendableClass, 6)
	require.Nil(t, s.Get(lendable, m.copies, -1))
	require.False(t, s.IsSet(lendable, m.copies))
	require.Equal(t, 0, s.Size(library, m.libraryBooks))
	require.False(t, s.Contains(library, m.libraryBooks, book))
	require.Equal(t, -1, s.IndexOf(library, m.libraryBooks, book))
	require.Equal(t, -1, s.LastIndexOf(library, m.libraryBooks, book))
	require.Empty(t, s.ToArray(library, m.libraryBooks))
	c, f := s.GetContainer(book)
	require.Nil(t, c)
	require.Nil(t, f)

	// and the writes fail
	s.RemoveRoot(library)
	require.Error(t, s.WaitOperations(context.Background(), nil))
}

func TestSQLStore_IndexResult(t *testing.T) {
	require.Equal(t, 2, indexResult(2, nil))
	require.Equal(t, -1, indexResult(-1, nil))
	// an error, or no index : not found
	require.Equal(t, -1, indexResult(nil, errors.New("error")))
	require.Equal(t, -1, indexResult(0, errors.New("error")))
	require.Equal(t, -1, indexResult(nil, nil))
}
