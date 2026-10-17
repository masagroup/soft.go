package ecore

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/panjf2000/ants/v2"
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
	}, withWait())
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
	require.Equal(t, 0, s.Get(lendable, eFeature, -1)) // the default
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
	require.NoError(t, os.WriteFile(serializedPath, bytes, 0644))
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
	}, withFlush(), withWait()))
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

	dateSchema, err := s.getFeatureSchema(book, m.publicationDate)
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
	book := newLendable(t, m.bookClass, 6)
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

	for _, table := range []string{"lendable", "library_books", `".contents"`} {
		require.NoError(t, s.ExecuteQuery(context.Background(), "DROP TABLE "+table, nil))
	}
	require.NoError(t, s.WaitOperations(context.Background(), nil))

	// the reads return zero values
	require.Nil(t, s.Get(book6, m.copies, -1))
	require.False(t, s.IsSet(book6, m.copies))
	require.Equal(t, 0, s.Size(library, m.libraryBooks))
	require.True(t, s.IsEmpty(library, m.libraryBooks))
	require.False(t, s.Contains(library, m.libraryBooks, book6))
	require.Equal(t, -1, s.IndexOf(library, m.libraryBooks, book6))
	require.Empty(t, s.ToArray(library, m.libraryBooks))
	require.Empty(t, s.GetRoots())

	// the writes reading a missing table fail, the errors are reported by the next wait
	s.Set(book6, m.copies, -1, 9, false)
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
	require.Nil(t, s.Get(book, m.copies, -1))
	require.False(t, s.IsSet(book, m.copies))
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

func TestSQLStore_Writer_PromisesAfterClose(t *testing.T) {
	m := newLibraryModel(t)
	s, _, _, _ := newWriterTestStore(t, nil)
	book := newLendable(t, m.bookClass, 6)
	require.NoError(t, s.Close())

	// the promises of a closed store are rejected, they don't block
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := s.UnRegister(book, false).Await(ctx)
	require.ErrorContains(t, err, "closed")
	_, err = s.Serialize(context.Background()).Await(ctx)
	require.ErrorContains(t, err, "closed")
}

func TestSQLBase_PromisePoolReleased(t *testing.T) {
	pool, err := ants.NewPool(1)
	require.NoError(t, err)
	promisePool := newAntsPromisePool(pool)

	ran := make(chan struct{})
	promisePool.Go(func() { close(ran) })
	<-ran

	// a released pool refuses the tasks : they run in a goroutine of their own
	pool.Release()
	ran = make(chan struct{})
	promisePool.Go(func() { close(ran) })
	select {
	case <-ran:
	case <-time.After(5 * time.Second):
		require.Fail(t, "the task was not run")
	}
}

// inheritanceModel is a model with multiple inheritance : Derived extends Named and Parent,
// so the id of the inherited children feature is 0 in Parent and 1 in Derived
type inheritanceModel struct {
	ePackage     EPackage
	namedClass   EClass
	name         EAttribute
	parentClass  EClass
	children     EReference
	derivedClass EClass
	childClass   EClass
}

func newInheritanceModel() *inheritanceModel {
	m := &inheritanceModel{}
	f := GetFactory()
	m.ePackage = f.CreateEPackage()
	m.ePackage.SetName("inheritance")
	m.ePackage.SetNsPrefix("inheritance")
	m.ePackage.SetNsURI("http:///inheritance.ecore")
	m.ePackage.SetEFactoryInstance(f.CreateEFactory())

	m.childClass = f.CreateEClass()
	m.childClass.SetName("Child")

	m.namedClass = f.CreateEClass()
	m.namedClass.SetName("Named")
	m.name = f.CreateEAttribute()
	m.name.SetName("name")
	m.name.SetEType(GetPackage().GetEString())
	m.namedClass.GetEStructuralFeatures().Add(m.name)

	m.parentClass = f.CreateEClass()
	m.parentClass.SetName("Parent")
	m.children = f.CreateEReference()
	m.children.SetName("children")
	m.children.SetEType(m.childClass)
	m.children.SetUpperBound(UNBOUNDED_MULTIPLICITY)
	m.children.SetContainment(true)
	m.parentClass.GetEStructuralFeatures().Add(m.children)

	m.derivedClass = f.CreateEClass()
	m.derivedClass.SetName("Derived")
	m.derivedClass.GetESuperTypes().AddAll(NewImmutableEList([]any{m.namedClass, m.parentClass}))

	m.ePackage.GetEClassifiers().AddAll(NewImmutableEList([]any{m.childClass, m.namedClass, m.parentClass, m.derivedClass}))
	return m
}

func TestSQLStore_InheritedFeatures(t *testing.T) {
	m := newLibraryModel(t)
	s, dbPath, _, _ := newWriterTestStore(t, nil)
	defer s.Close()

	// copies and borrowers are declared by Lendable, publicationDate by Item
	book := newLendable(t, m.bookClass, 6)
	borrowers := m.lendableClass.GetEStructuralFeatureFromName("borrowers")
	require.Equal(t, 4, s.Get(book, m.copies, -1))
	require.True(t, s.IsSet(book, m.copies))
	require.Equal(t, 0, s.Size(book, borrowers))
	require.True(t, s.IsEmpty(book, borrowers))
	date, isDate := s.Get(book, m.publicationDate, -1).(*time.Time)
	require.True(t, isDate)
	require.Equal(t, time.Date(2015, time.September, 6, 4, 24, 46, 0, time.UTC), date.UTC())

	// written in the table of the declaring class
	s.Set(book, m.copies, -1, 9, false)
	require.NoError(t, s.WaitOperations(context.Background(), nil))
	require.Equal(t, 9, countCopies(t, dbPath, 6))
	require.Equal(t, 9, s.Get(book, m.copies, -1))

	// a feature neither declared nor inherited is unknown
	require.Nil(t, s.Get(book, m.libraryBooks, -1))
}

func TestSQLStore_InheritedContainmentFeatureID(t *testing.T) {
	m := newInheritanceModel()
	require.Equal(t, 0, m.parentClass.GetFeatureID(m.children))
	require.Equal(t, 1, m.derivedClass.GetFeatureID(m.children))

	dbPath := filepath.Join(t.TempDir(), "inheritance.sqlite")
	packageRegistry := NewMockEPackageRegistry(t)
	packageRegistry.EXPECT().GetPackage(m.ePackage.GetNsURI()).Return(m.ePackage).Maybe()
	s, err := NewSQLStore(dbPath, NewURI(""), nil, packageRegistry, nil)
	require.NoError(t, err)
	defer s.Close()

	eFactory := m.ePackage.GetEFactoryInstance()
	derived := eFactory.Create(m.derivedClass)
	derived.ESet(m.name, "derived")
	child := eFactory.Create(m.childClass)
	derived.EGet(m.children).(EList).Add(child)
	s.AddRoot(derived)
	require.NoError(t, s.WaitOperations(context.Background(), nil))

	// the child is contained by the feature 1 of Derived, not by its feature 0 (name)
	derivedID, isDerivedID := s.sqlIDManager.GetObjectID(derived)
	require.True(t, isDerivedID)
	childID, isChildID := s.sqlIDManager.GetObjectID(child)
	require.True(t, isChildID)
	containerID, featureID, err := containerOf(s, childID)
	require.NoError(t, err)
	require.Equal(t, derivedID, containerID)
	require.Equal(t, int64(1), featureID)
	container, feature := s.GetContainer(child)
	require.Equal(t, derived, container)
	require.Equal(t, m.children, feature)

	// the inherited features are read and written through the derived object
	require.Equal(t, "derived", s.Get(derived, m.name, -1))
	require.Equal(t, 1, s.Size(derived, m.children))
	s.Set(derived, m.name, -1, "renamed", false)
	require.NoError(t, s.WaitOperations(context.Background(), nil))
	require.Equal(t, "renamed", s.Get(derived, m.name, -1))
}

// newRegistryStore creates a store on a copy of testdata/library.store.sqlite with a registry
// of the library package : the objects of the database can be decoded
func newRegistryStore(t *testing.T, m *libraryModel) (*SQLStore, string) {
	dbPath := filepath.Join(t.TempDir(), "library.store.sqlite")
	require.NoError(t, copyFile("testdata/library.store.sqlite", dbPath))
	registry := NewEPackageRegistryImpl()
	registry.RegisterPackage(m.ePackage)
	s, err := NewSQLStore(dbPath, NewURI(""), nil, registry, nil)
	require.NoError(t, err)
	return s, dbPath
}

func newDataListStore(t *testing.T) (s *SQLStore, book *MockSQLObject, contents EStructuralFeature) {
	ePackage := loadPackage("library.datalist.ecore")
	require.NotNil(t, ePackage)
	eClass, _ := ePackage.GetEClassifier("Book").(EClass)
	require.NotNil(t, eClass)
	contents = eClass.GetEStructuralFeatureFromName("contents")
	require.NotNil(t, contents)
	dbPath := filepath.Join(t.TempDir(), "library.datalist.sqlite")
	require.NoError(t, copyFile("testdata/library.datalist.sqlite", dbPath))
	s, err := NewSQLStore(dbPath, NewURI(""), nil, NewMockEPackageRegistry(t), nil)
	require.NoError(t, err)
	book = newLendable(t, eClass, 4)
	return
}

func TestSQLStore_Writer_FlushReportsCommitError(t *testing.T) {
	s, dbPath, eClass, eFeature := newWriterTestStore(t, nil)
	defer s.Close()
	lendable := newLendable(t, eClass, 6)
	_, err := s.execute(newOperation("foreign keys", operationWrite, nil, nil, -1, nil, func() (any, error) {
		return nil, sqlitex.ExecuteTransient(s.writerConn, "PRAGMA foreign_keys=ON;", nil)
	}, withFlush(), withWait()))
	require.NoError(t, err)
	require.NoError(t, s.ExecuteQuery(context.Background(), "CREATE TABLE parent (id INTEGER PRIMARY KEY)", nil))
	require.NoError(t, s.ExecuteQuery(context.Background(), "CREATE TABLE child (parentID INTEGER REFERENCES parent(id) DEFERRABLE INITIALLY DEFERRED)", nil))
	require.NoError(t, s.WaitOperations(context.Background(), nil))

	// a write, a deferred violation, and a flush queued before the commit : the flush
	// commits first and reports the error of the commit
	release := blockWriter(t, s)
	s.Set(lendable, eFeature, -1, 9, false)
	go func() { _ = s.ExecuteQuery(context.Background(), "INSERT INTO child VALUES (42)", nil) }()
	require.Eventually(t, func() bool { return queueLen(s) == 2 }, time.Second, time.Millisecond)
	flushErr := make(chan error, 1)
	go func() { flushErr <- s.WaitOperations(context.Background(), nil) }()
	require.Eventually(t, func() bool { return queueLen(s) == 3 }, time.Second, time.Millisecond)
	release()
	require.ErrorContains(t, <-flushErr, "FOREIGN KEY")
	require.NoError(t, s.WaitOperations(context.Background(), nil))
	require.Equal(t, 4, countCopies(t, dbPath, 6))
}

func TestSQLStore_AddAll_Order(t *testing.T) {
	s, book, contents := newDataListStore(t)
	defer s.Close()
	wait := func() { require.NoError(t, s.WaitOperations(context.Background(), nil)) }

	// between two elements, before the first, after the last
	s.AddAll(book, contents, 1, NewImmutableEList([]any{"x", "y"}))
	wait()
	require.Equal(t, []any{"c21", "x", "y", "c22", "c23"}, s.ToArray(book, contents))
	s.AddAll(book, contents, 1, NewImmutableEList([]any{"p", "q", "r"}))
	wait()
	require.Equal(t, []any{"c21", "p", "q", "r", "x", "y", "c22", "c23"}, s.ToArray(book, contents))
	s.AddAll(book, contents, 0, NewImmutableEList([]any{"a", "b", "c"}))
	wait()
	require.Equal(t, []any{"a", "b", "c", "c21", "p", "q", "r", "x", "y", "c22", "c23"}, s.ToArray(book, contents))
	s.AddAll(book, contents, 11, NewImmutableEList([]any{"z1", "z2"}))
	wait()
	require.Equal(t, []any{"a", "b", "c", "c21", "p", "q", "r", "x", "y", "c22", "c23", "z1", "z2"}, s.ToArray(book, contents))
	s.AddAll(book, contents, 0, NewImmutableEList([]any{"first"}))
	wait()
	require.Equal(t, "first", s.Get(book, contents, 0))
	require.Equal(t, "a", s.Get(book, contents, 1))

	// out of bounds
	s.AddAll(book, contents, 20, NewImmutableEList([]any{"out"}))
	require.ErrorContains(t, s.WaitOperations(context.Background(), nil), "not in list bounds")

	// repeated inserts at the same place keep the order
	for i := range 20 {
		s.Add(book, contents, 2, fmt.Sprintf("i%d", i))
	}
	wait()
	values := s.ToArray(book, contents)
	require.Equal(t, "first", values[0])
	require.Equal(t, "a", values[1])
	for i := range 20 {
		require.Equal(t, fmt.Sprintf("i%d", 19-i), values[2+i])
	}
	require.Equal(t, "b", values[22])
}

func TestSQLStore_Add_NotAList(t *testing.T) {
	s, _, eClass, eFeature := newWriterTestStore(t, nil)
	defer s.Close()
	lendable := newLendable(t, eClass, 6)

	// a list operation on a single valued feature is an error, not a panic holding a lock
	s.Add(lendable, eFeature, 0, 1)
	require.ErrorContains(t, s.WaitOperations(context.Background(), nil), "is not a list")
	s.Clear(lendable, eFeature)
	require.ErrorContains(t, s.WaitOperations(context.Background(), nil), "is not a list")
	require.Equal(t, 0, s.Size(lendable, eFeature))
	require.True(t, s.IsEmpty(lendable, eFeature))
	require.Equal(t, 4, s.Get(lendable, eFeature, -1))
}

func TestSQLStore_IsSet_Default(t *testing.T) {
	m := newLibraryModel(t)
	s, _ := newRegistryStore(t, m)
	defer s.Close()
	book := newLendable(t, m.bookClass, 6)
	pages := m.bookClass.GetEStructuralFeatureFromName("pages")
	wait := func() { require.NoError(t, s.WaitOperations(context.Background(), nil)) }

	// a numeric feature holding its default is not set
	require.True(t, s.IsSet(book, m.copies))
	s.Set(book, m.copies, -1, 0, false)
	wait()
	require.False(t, s.IsSet(book, m.copies))
	require.True(t, s.IsSet(book, pages))
	s.Set(book, pages, -1, 100, false)
	wait()
	require.False(t, s.IsSet(book, pages))
	s.Set(book, pages, -1, 101, false)
	wait()
	require.True(t, s.IsSet(book, pages))

	// an enum holding its default is not set
	require.True(t, s.IsSet(book, m.bookCategory))
	s.UnSet(book, m.bookCategory)
	wait()
	require.False(t, s.IsSet(book, m.bookCategory))

	// an object
	library := newLendable(t, m.libraryClass, 2)
	ownerPdg := m.libraryClass.GetEStructuralFeatureFromName("ownerPdg")
	require.False(t, s.IsSet(library, ownerPdg))
}

func TestSQLStore_UnSet_Enum(t *testing.T) {
	m := newLibraryModel(t)
	s, _ := newRegistryStore(t, m)
	defer s.Close()
	book := newLendable(t, m.bookClass, 6)
	require.Equal(t, 2, s.Get(book, m.bookCategory, -1)) // Biography

	// the default literal is encoded : it is read back
	s.UnSet(book, m.bookCategory)
	require.NoError(t, s.WaitOperations(context.Background(), nil))
	require.Equal(t, m.bookCategory.GetDefaultValue(), s.Get(book, m.bookCategory, -1))
	require.Equal(t, 0, s.Get(book, m.bookCategory, -1)) // Mystery

	// the date too
	s.UnSet(book, m.publicationDate)
	require.NoError(t, s.WaitOperations(context.Background(), nil))
	require.Nil(t, s.Get(book, m.publicationDate, -1))
}

func TestSQLStore_Writer_NestedReadKeepsBinding(t *testing.T) {
	m := newLibraryModel(t)
	s, _ := newRegistryStore(t, m)
	defer s.Close()
	library := newLendable(t, m.libraryClass, 2)
	book7 := newLendable(t, m.bookClass, 7)

	// a read on a pool connection from a read on a pool connection : the connection of
	// the outer read is bound again once the inner read is done
	books := 0
	for range s.All(library, m.libraryBooks) {
		require.NotNil(t, s.boundConn())
		require.Equal(t, 3, s.Get(book7, m.copies, -1))
		require.NotNil(t, s.boundConn())
		books++
	}
	require.Equal(t, 2, books)
	require.Nil(t, s.boundConn())
}

func TestSQLStore_DecodeObject_Concurrent(t *testing.T) {
	m := newLibraryModel(t)
	s, _ := newRegistryStore(t, m)
	defer s.Close()
	library := newLendable(t, m.libraryClass, 2)

	// the first book is decoded by many goroutines at once : one object
	decoded := make(chan any, 50)
	wg := sync.WaitGroup{}
	for range cap(decoded) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			decoded <- s.Get(library, m.libraryBooks, 0)
		}()
	}
	wg.Wait()
	close(decoded)
	first := <-decoded
	require.NotNil(t, first)
	for book := range decoded {
		require.Same(t, first, book)
	}
}

// a reader holds an object while it waits for the writer : the writer must never need
// the object, in particular to walk the contents it unregisters
func TestSQLStore_Writer_UnRegisterContents_ReaderHoldingObject(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)
	libraryClass, _ := ePackage.GetEClassifier("Library").(EClass)
	require.NotNil(t, libraryClass)
	booksFeature := libraryClass.GetEStructuralFeatureFromName("books")
	require.NotNil(t, booksFeature)

	dbPath := filepath.Join(t.TempDir(), "library.complex.sqlite")
	require.NoError(t, copyFile("testdata/library.complex.sqlite", dbPath))
	s, err := NewSQLStore(dbPath, NewURI(""), nil, nil, nil)
	require.NoError(t, err)
	defer s.Close()

	// a store object without cache : its features are read from the store
	library := NewEStoreEObjectImpl(false)
	library.SetEClass(libraryClass)
	library.SetEStore(s)

	// the writer is blocked, an unregistration of the library and its contents is pending
	release := blockWriter(t, s)
	unregistered := s.UnRegister(library, true)

	// a reader holds the library while it waits for the writer, the library being dirty
	read := make(chan struct{})
	go func() {
		defer close(read)
		library.EGet(booksFeature)
	}()
	time.Sleep(20 * time.Millisecond)
	release()

	// both finish
	_, err = unregistered.Await(context.Background())
	require.NoError(t, err)
	select {
	case <-read:
	case <-time.After(5 * time.Second):
		require.Fail(t, "the reader holding the object and the writer unregistering it wait for each other")
	}
}

// a reader holds an object while it waits for the writer : the writer must never need
// the object, in particular to encode it when it is the value of a pending write
func TestSQLStore_Writer_EncodeObject_ReaderHoldingObject(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)
	libraryClass, _ := ePackage.GetEClassifier("Library").(EClass)
	require.NotNil(t, libraryClass)
	booksFeature := libraryClass.GetEStructuralFeatureFromName("books")
	require.NotNil(t, booksFeature)
	bookClass, _ := ePackage.GetEClassifier("Book").(EClass)
	require.NotNil(t, bookClass)
	titleFeature := bookClass.GetEStructuralFeatureFromName("title")
	require.NotNil(t, titleFeature)

	dbPath := filepath.Join(t.TempDir(), "library.sqlite")
	s, err := NewSQLStore(dbPath, NewURI(""), nil, nil, nil)
	require.NoError(t, err)
	defer s.Close()

	library := NewEStoreEObjectImpl(false)
	library.SetEClass(libraryClass)
	library.SetEStore(s)

	// a new store object without cache : its features are read from the store
	book := NewEStoreEObjectImpl(false)
	book.SetEClass(bookClass)
	book.SetEStore(s)

	// the writer is blocked, an addition of the book is pending : the writer encodes the book
	release := blockWriter(t, s)
	s.Add(library, booksFeature, 0, book)

	// a reader holds the book while it waits for the writer, the book being dirty
	read := make(chan struct{})
	go func() {
		defer close(read)
		book.EGet(titleFeature)
	}()
	time.Sleep(20 * time.Millisecond)
	release()

	select {
	case <-read:
	case <-time.After(5 * time.Second):
		require.Fail(t, "the reader holding the object and the writer encoding it wait for each other")
	}
	require.NoError(t, s.Flush(context.Background()))
}
