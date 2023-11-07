package ecore

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

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
	sqliteEncoder := NewSQLEncoder(eResource, w, nil)
	sqliteEncoder.EncodeResource()
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
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
	sqliteEncoder := NewSQLEncoder(eResource, w, nil)
	sqliteEncoder.EncodeResource()
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
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
	sqliteEncoder := NewSQLEncoder(eResource, w, nil)
	sqliteEncoder.EncodeResource()
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
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

	//w := &bytes.Buffer{}
	w, err := os.Create("testdata/emap.sqlite")
	require.NoError(t, err)
	defer w.Close()

	binaryEncoder := NewSQLEncoder(eResource, w, nil)
	binaryEncoder.EncodeResource()
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
}
