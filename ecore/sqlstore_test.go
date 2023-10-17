package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSQLStore_Constructor(t *testing.T) {
	s, err := NewSQLStore("testdata/library.store.sqlite", NewURI(""), nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
}

func TestSQLStore_SetSingleValue(t *testing.T) {
	xmiProcessor := NewXMIProcessor()
	resource := xmiProcessor.Load(NewURI("testdata/library.complex.ecore"))
	require.NotNil(t, resource)
	assert.True(t, resource.IsLoaded())
	assert.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	assert.True(t, resource.GetWarnings().Empty(), diagnosticError(resource.GetWarnings()))

	contents := resource.GetContents()
	require.Equal(t, 1, contents.Size())

	ePackage, _ := contents.Get(0).(EPackage)
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Lendable").(EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("copies")
	require.NotNil(t, eFeature)

	s, err := NewSQLStore("testdata/library.store.sqlite", NewURI(""), nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()
	mockObject := NewMockSQLObject(t)
	mockObject.EXPECT().GetSqlID().Return(int64(3)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	s.Set(mockObject, eFeature, -1, 5)
}

func TestSQLStore_SetListValue(t *testing.T) {
	xmiProcessor := NewXMIProcessor()
	resource := xmiProcessor.Load(NewURI("testdata/library.datalist.ecore"))
	require.NotNil(t, resource)
	assert.True(t, resource.IsLoaded())
	assert.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	assert.True(t, resource.GetWarnings().Empty(), diagnosticError(resource.GetWarnings()))

	contents := resource.GetContents()
	require.Equal(t, 1, contents.Size())

	ePackage, _ := contents.Get(0).(EPackage)
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Book").(EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("contents")
	require.NotNil(t, eFeature)

	s, err := NewSQLStore("testdata/library.datalist.sqlite", NewURI(""), nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()
	mockObject := NewMockSQLObject(t)
	mockObject.EXPECT().GetSqlID().Return(int64(5)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	s.Set(mockObject, eFeature, 1, "c4")
}
