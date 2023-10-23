package ecore

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSQLStore_Constructor(t *testing.T) {
	s, err := NewSQLStore("testdata/library.store.sqlite", NewURI(""), nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
}

func closeFile(f *os.File, reported *error) {
	if err := f.Close(); *reported == nil {
		*reported = err
	}
}

func copyFile(src, dest string) (err error) {
	if err = os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return
	}

	d, err := os.Create(dest)
	if err != nil {
		return
	}
	defer closeFile(d, &err)

	s, err := os.Open(src)
	if err != nil {
		return
	}
	defer closeFile(s, &err)

	_, err = io.Copy(d, s)
	return
}

func TestSQLStore_SetSingleValue(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Lendable").(EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("copies")
	require.NotNil(t, eFeature)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.store.sqlite")
	err := copyFile("testdata/library.store.sqlite", dbPath)
	require.Nil(t, err)

	// store
	s, err := NewSQLStore(dbPath, NewURI(""), nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockSQLObject(t)
	mockObject.EXPECT().GetSqlID().Return(int64(3)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	s.Set(mockObject, eFeature, -1, 5)
}

func TestSQLStore_GetSingleValue(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Lendable").(EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("copies")
	require.NotNil(t, eFeature)

	// store
	s, err := NewSQLStore("testdata/library.store.sqlite", NewURI(""), nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockSQLObject(t)
	mockObject.EXPECT().GetSqlID().Return(int64(3)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	v := s.Get(mockObject, eFeature, -1)
	require.NotNil(t, v)
}

func TestSQLStore_SetListValue(t *testing.T) {
	ePackage := loadPackage("library.datalist.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Book").(EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("contents")
	require.NotNil(t, eFeature)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.datalist.sqlite")
	err := copyFile("testdata/library.datalist.sqlite", dbPath)
	require.Nil(t, err)

	// create store
	s, err := NewSQLStore(dbPath, NewURI(""), nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockSQLObject(t)
	mockObject.EXPECT().GetSqlID().Return(int64(5)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	s.Set(mockObject, eFeature, 1, "c4")

	// load db and retrieve new value
}

func TestSQLStore_GetListValue(t *testing.T) {
	ePackage := loadPackage("library.datalist.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Book").(EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("contents")
	require.NotNil(t, eFeature)

	// create store
	s, err := NewSQLStore("testdata/library.datalist.sqlite", NewURI(""), nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockSQLObject(t)
	mockObject.EXPECT().GetSqlID().Return(int64(5)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	s.Get(mockObject, eFeature, 1)

	// load db and retrieve new value
}
