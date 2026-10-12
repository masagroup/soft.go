package ecore

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// BenchmarkSQLStore_Add measures fire-and-forget adds to a list of the store, waiting
// for them to be written every 1,000 adds : the cost of an operation of the store
func BenchmarkSQLStore_Add(b *testing.B) {
	ePackage := loadPackage("library.datalist.ecore")
	require.NotNil(b, ePackage)
	eClass, _ := ePackage.GetEClassifier("Book").(EClass)
	require.NotNil(b, eClass)
	eFeature := eClass.GetEStructuralFeatureFromName("contents")
	require.NotNil(b, eFeature)

	dbPath := filepath.Join(b.TempDir(), "library.datalist.sqlite")
	require.NoError(b, copyFile("testdata/library.datalist.sqlite", dbPath))
	s, err := NewSQLStore(dbPath, NewURI(""), nil, NewMockEPackageRegistry(b), nil)
	require.NoError(b, err)
	defer s.Close()

	book := NewMockSQLObject(b)
	book.EXPECT().GetSQLID().Return(int64(5)).Maybe()
	book.EXPECT().SetSQLID(int64(5)).Return().Maybe()
	book.EXPECT().EClass().Return(eClass).Maybe()
	index := s.Size(book, eFeature)

	b.ReportAllocs()
	for b.Loop() {
		s.Add(book, eFeature, index, "content")
		index++
		if index%1000 == 0 {
			require.NoError(b, s.WaitOperations(context.Background(), nil))
		}
	}
	require.NoError(b, s.WaitOperations(context.Background(), nil))
}

// BenchmarkSQLStore_Set measures fire-and-forget sets of an attribute of the store, waiting
// for them to be written every 1,000 sets
func BenchmarkSQLStore_Set(b *testing.B) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(b, ePackage)
	eClass, _ := ePackage.GetEClassifier("Lendable").(EClass)
	require.NotNil(b, eClass)
	eFeature := eClass.GetEStructuralFeatureFromName("copies")
	require.NotNil(b, eFeature)

	dbPath := filepath.Join(b.TempDir(), "library.store.sqlite")
	require.NoError(b, copyFile("testdata/library.store.sqlite", dbPath))
	s, err := NewSQLStore(dbPath, NewURI(""), nil, nil, nil)
	require.NoError(b, err)
	defer s.Close()

	lendable := NewMockSQLObject(b)
	lendable.EXPECT().GetSQLID().Return(int64(6)).Maybe()
	lendable.EXPECT().SetSQLID(int64(6)).Return().Maybe()
	lendable.EXPECT().EClass().Return(eClass).Maybe()

	b.ReportAllocs()
	i := 0
	for b.Loop() {
		s.Set(lendable, eFeature, -1, i, false)
		i++
		if i%1000 == 0 {
			require.NoError(b, s.WaitOperations(context.Background(), nil))
		}
	}
	require.NoError(b, s.WaitOperations(context.Background(), nil))
}
