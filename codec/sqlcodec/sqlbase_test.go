package sqlcodec

import (
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

func TestSQLBase_NewQuery(t *testing.T) {
	tests := []struct {
		cmd    string
		type_  queryType
		tables []string
	}{
		{"SELECT * FROM a", queryRead, []string{"a"}},
		{"SELECT * FROM a JOIN b ON a.id=b.id", queryRead, []string{"a", "b"}},
		{"SELECT * FROM (a)", queryRead, []string{"a"}},
		{"INSERT INTO a VALUES (1)", queryWrite, []string{"a", ""}},
		{"UPDATE a SET x=1", queryWrite, []string{"a", ""}},
		{"DELETE FROM a", queryWrite, []string{"a", ""}},
		{"CREATE TABLE a (x INTEGER)", queryWrite, []string{"a", ""}},
		// an unknown statement or a parser error : a write of all the tables
		{"DROP TABLE a", queryWrite, []string{""}},
		{"NOT SQL", queryWrite, []string{""}},
	}
	for _, test := range tests {
		t.Run(test.cmd, func(t *testing.T) {
			q := newQuery(test.cmd)
			assert.Equal(t, test.type_, q.type_)
			assert.Equal(t, test.tables, q.tables)
		})
	}
}

// blockingRead executes a read of the lendable table through the query scheduler, from a goroutine
// without connection : it blocks in its first row until release is closed
func blockingRead(s *SQLStore, release <-chan struct{}) (started <-chan struct{}, done <-chan error) {
	startedCh := make(chan struct{})
	doneCh := make(chan error, 1)
	once := sync.Once{}
	go func() {
		doneCh <- s.executeQuery("SELECT lendableID FROM lendable", &sqlitex.ExecOptions{
			ResultFunc: func(stmt *sqlite.Stmt) error {
				once.Do(func() { close(startedCh) })
				<-release
				return nil
			}})
	}()
	return startedCh, doneCh
}

func TestSQLBase_Scheduler_WriteWaitsForRead(t *testing.T) {
	// a debug logger : the scheduler logs the queries a query waits for
	core, logs := observer.New(zap.DebugLevel)
	s, dbPath, _, _ := newWriterTestStore(t, map[string]any{SQL_OPTION_LOGGER: zap.New(core)})
	defer s.Close()

	release := make(chan struct{})
	started, readDone := blockingRead(s, release)
	<-started

	// the write of the table waits for the read
	writeDone := make(chan error, 1)
	go func() {
		writeDone <- s.executeQuery("UPDATE lendable SET copies=7 WHERE lendableID=6", nil)
	}()
	select {
	case <-writeDone:
		require.Fail(t, "the write was executed before the read of its table")
	case <-time.After(50 * time.Millisecond):
	}
	close(release)
	require.NoError(t, <-readDone)
	require.NoError(t, <-writeDone)
	require.Equal(t, 7, countCopies(t, dbPath, 6))
	require.NotZero(t, logs.FilterMessage("waiting previous queries").Len())
}

func TestSQLBase_Scheduler_PreviousError(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	s, _, _, _ := newWriterTestStore(t, map[string]any{SQL_OPTION_LOGGER: zap.New(core)})
	defer s.Close()

	release := make(chan struct{})
	started, readDone := blockingRead(s, release)
	<-started

	// a failing write waiting for the read, and a read waiting for the write
	writeDone := make(chan error, 1)
	go func() {
		writeDone <- s.executeQuery("INSERT INTO lendable (lendableID, copies) VALUES (6, 1)", nil)
	}()
	require.Eventually(t, func() bool {
		s.sqliteMutex.Lock()
		defer s.sqliteMutex.Unlock()
		return len(s.sqliteQueries["lendable"]) == 2
	}, time.Second, time.Millisecond)
	nextDone := make(chan error, 1)
	go func() {
		nextDone <- s.executeQuery("SELECT copies FROM lendable", nil)
	}()
	require.Eventually(t, func() bool {
		s.sqliteMutex.Lock()
		defer s.sqliteMutex.Unlock()
		return len(s.sqliteQueries["lendable"]) == 3
	}, time.Second, time.Millisecond)

	// the error of the write is the error of the read waiting for it
	close(release)
	require.NoError(t, <-readDone)
	writeErr := <-writeDone
	require.ErrorContains(t, writeErr, "constraint")
	require.ErrorContains(t, <-nextDone, "constraint")
	require.NotZero(t, logs.FilterMessage("error in previous query").Len())

	// the next queries work
	require.NoError(t, s.executeQuery("SELECT copies FROM lendable", nil))
}
