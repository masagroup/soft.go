package ecore

import (
	"bytes"
	"database/sql"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func requireCopyBytes(t require.TestingT, filePath string, buff []byte) {
	actualFile, err := os.Create(filePath)
	require.NoError(t, err)
	defer actualFile.Close()

	_, err = io.Copy(actualFile, bytes.NewBuffer(buff))
	require.NoError(t, err)
}

func requireSameDB(t require.TestingT, expectedPath string, actualBytes []byte) {
	// check buffers
	expectedBytes, err := os.ReadFile(expectedPath)
	require.NoError(t, err)
	if bytes.Equal(expectedBytes, actualBytes) {
		return
	}

	// open expected db file
	expectedDB, err := sql.Open("sqlite", expectedPath)
	require.NoError(t, err)
	defer func() {
		_ = expectedDB.Close()
	}()

	// open actual db
	actualPath, err := sqlTmpDB("")
	require.NoError(t, err)
	requireCopyBytes(t, actualPath, actualBytes)

	actualDB, err := sql.Open("sqlite", actualPath)
	require.NoError(t, err)
	defer func() {
		_ = actualDB.Close()
	}()

	// check that db are equal
	RequireEqualDB(t, expectedDB, actualDB)
}

func TestSqlEncoder_Complex(t *testing.T) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	// load resource
	xmlProcessor := NewXMLProcessor(XMLProcessorPackages([]EPackage{ePackage}))
	eResource := xmlProcessor.LoadWithOptions(NewURI("testdata/library.complex.xml"), nil)
	require.NotNil(t, eResource)
	require.True(t, eResource.IsLoaded())
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	require.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	// w, err := os.Create("testdata/library.complex.sqlite")
	// require.NoError(t, err)
	// defer w.Close()
	w := &bytes.Buffer{}
	sqliteEncoder := NewSQLWriterEncoder(w, eResource, nil)
	sqliteEncoder.EncodeResource()
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))

	// compare expected and actual bytes
	requireSameDB(t, "testdata/library.complex.sqlite", w.Bytes())
}

func TestSqlEncoder_DataList(t *testing.T) {
	// load package
	ePackage := loadPackage("library.datalist.ecore")
	require.NotNil(t, ePackage)

	// load resource
	xmlProcessor := NewXMLProcessor(XMLProcessorPackages([]EPackage{ePackage}))
	eResource := xmlProcessor.LoadWithOptions(NewURI("testdata/library.datalist.xml"), nil)
	require.NotNil(t, eResource)
	require.True(t, eResource.IsLoaded())
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	require.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	// w, err := os.Create("testdata/library.datalist.sqlite")
	// require.NoError(t, err)
	// defer w.Close()
	w := &bytes.Buffer{}
	sqliteEncoder := NewSQLWriterEncoder(w, eResource, nil)
	sqliteEncoder.EncodeResource()
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))

	// compare expected and actual bytes
	requireSameDB(t, "testdata/library.datalist.sqlite", w.Bytes())
}

func TestSqlEncoder_ComplexWithOwner(t *testing.T) {
	// load package and retrieve library / person features
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eLibraryClass, _ := ePackage.GetEClassifier("Library").(EClass)
	require.NotNil(t, eLibraryClass)

	eLibraryOwnerFeature, _ := eLibraryClass.GetEStructuralFeatureFromName("ownerPdg").(EReference)
	require.NotNil(t, eLibraryOwnerFeature)

	ePersonClass, _ := ePackage.GetEClassifier("Person").(EClass)
	require.NotNil(t, ePersonClass)

	ePersonAdressAttribute, _ := ePersonClass.GetEStructuralFeatureFromName("address").(EAttribute)
	require.NotNil(t, ePersonAdressAttribute)

	ePersonFirstNameAttribute, _ := ePersonClass.GetEStructuralFeatureFromName("firstName").(EAttribute)
	require.NotNil(t, ePersonFirstNameAttribute)

	ePersonLastNameAttribute, _ := ePersonClass.GetEStructuralFeatureFromName("lastName").(EAttribute)
	require.NotNil(t, ePersonLastNameAttribute)

	// create a library with a owner
	eFactory := ePackage.GetEFactoryInstance()

	aPerson := eFactory.Create(ePersonClass)
	aPerson.ESet(ePersonAdressAttribute, "owner adress")
	aPerson.ESet(ePersonFirstNameAttribute, "owner first name")
	aPerson.ESet(ePersonLastNameAttribute, "owner last name")

	aLibrary := eFactory.Create(eLibraryClass)
	aLibrary.ESet(eLibraryOwnerFeature, aPerson)

	eResourceSet := CreateEResourceSet([]EPackage{ePackage})
	eResource := eResourceSet.CreateResource(NewURI("testdata/library.owner.sqlite"))
	eResource.GetContents().Add(aLibrary)

	// w, err := os.Create("testdata/library.owner.sqlite")
	// require.NoError(t, err)
	// defer w.Close()
	w := &bytes.Buffer{}
	sqliteEncoder := NewSQLWriterEncoder(w, eResource, nil)
	sqliteEncoder.EncodeResource()
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))

	// compare expected and actual bytes
	requireSameDB(t, "testdata/library.owner.sqlite", w.Bytes())
}

func TestSQLEncoder_Maps(t *testing.T) {
	// load package
	ePackage := loadPackage("emap.ecore")
	require.NotNil(t, ePackage)

	// load resource
	xmlProcessor := NewXMLProcessor(XMLProcessorPackages([]EPackage{ePackage}))
	eResource := xmlProcessor.LoadWithOptions(NewURI("testdata/emap.xml"), nil)
	require.NotNil(t, eResource)
	require.True(t, eResource.IsLoaded())
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	require.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	// w, err := os.Create("testdata/emap.sqlite")
	// require.NoError(t, err)
	// defer w.Close()
	w := &bytes.Buffer{}
	sqliteEncoder := NewSQLWriterEncoder(w, eResource, nil)
	sqliteEncoder.EncodeResource()
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))

	// compare expected and actual bytes
	requireSameDB(t, "testdata/emap.sqlite", w.Bytes())
}
