package ecore

import (
	"testing"

	"github.com/stretchr/testify/require"
	"zombiezen.com/go/sqlite"
)

type MockT struct {
	Failed bool
}

func (t *MockT) FailNow() {
	t.Failed = true
}

func (t *MockT) Errorf(format string, args ...interface{}) {
	_, _ = format, args
}

func TestRequireEqualDB_DifferentTables(t *testing.T) {
	mockT := &MockT{}
	conn1, err := sqlite.OpenConn("testdata/library.simple.sqlite")
	require.NoError(t, err)
	defer conn1.Close()
	conn2, err := sqlite.OpenConn("testdata/library.simple.diff.table.sqlite")
	require.NoError(t, err)
	defer conn2.Close()
	RequireEqualDB(mockT, conn1, conn2)
	require.True(t, mockT.Failed)
}

func TestRequireEqualDB_DifferentContent(t *testing.T) {
	mockT := &MockT{}
	conn1, err := sqlite.OpenConn("testdata/library.simple.sqlite")
	require.NoError(t, err)
	defer conn1.Close()
	conn2, err := sqlite.OpenConn("testdata/library.simple.diff.content.sqlite")
	require.NoError(t, err)
	defer conn2.Close()
	RequireEqualDB(mockT, conn1, conn2)
	require.True(t, mockT.Failed)
}

func TestRequireEqualDB_Same(t *testing.T) {
	mockT := &MockT{}
	conn1, err := sqlite.OpenConn("testdata/library.simple.sqlite")
	require.NoError(t, err)
	defer conn1.Close()
	conn2, err := sqlite.OpenConn("testdata/library.simple.sqlite")
	require.NoError(t, err)
	defer conn2.Close()
	RequireEqualDB(mockT, conn1, conn2)
	require.False(t, mockT.Failed)
}
