package ecore

import (
	"bytes"
	"os"
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/require"
	"zombiezen.com/go/sqlite"
)

func requireSameDB(t require.TestingT, expectedPath string, actualBytes []byte) {
	// check buffers
	expectedBytes, err := os.ReadFile(expectedPath)
	require.NoError(t, err)
	if bytes.Equal(expectedBytes, actualBytes) {
		return
	}

	// open expected db file
	expectedConn, err := sqlite.OpenConn(expectedPath)
	require.NoError(t, err)
	defer func() {
		_ = expectedConn.Close()
	}()

	// open actual db
	actualPath, err := sqlTmpDB("")
	require.NoError(t, err)
	err = os.WriteFile(actualPath, actualBytes, 0644)
	require.NoError(t, err)

	actualConn, err := sqlite.OpenConn(actualPath)
	require.NoError(t, err)
	defer func() {
		_ = actualConn.Close()
	}()

	// check that db are equal
	RequireEqualDB(t, expectedConn, actualConn)
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

	w, err := os.Create("testdata/library.complex.sqlite")
	require.NoError(t, err)
	defer w.Close()
	// w := &bytes.Buffer{}
	sqliteEncoder := NewSQLWriterEncoder(w, eResource, nil)
	sqliteEncoder.EncodeResource()
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))

	// compare expected and actual bytes
	// requireSameDB(t, "testdata/library.complex.sqlite", w.Bytes())
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

func TestSQLEncoder_Simple(t *testing.T) {
	// load package
	ePackage := loadPackage("library.simple.ecore")
	require.NotNil(t, ePackage)

	// load resource
	xmlProcessor := NewXMLProcessor(XMLProcessorPackages([]EPackage{ePackage}))
	resource := xmlProcessor.LoadWithOptions(NewURI("testdata/library.simple.xml"), nil)
	require.NotNil(t, resource)
	require.True(t, resource.IsLoaded())
	require.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	require.True(t, resource.GetWarnings().Empty(), diagnosticError(resource.GetWarnings()))

	// w, err := os.Create("testdata/library.simple.sqlite")
	// require.NoError(t, err)
	// defer w.Close()
	w := &bytes.Buffer{}
	sqliteEncoder := NewSQLWriterEncoder(w, resource, nil)
	sqliteEncoder.EncodeResource()
	require.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))

	// compare expected and actual bytes
	requireSameDB(t, "testdata/library.simple.sqlite", w.Bytes())

}

func TestSQLEncoder_SimpleWithULIDs(t *testing.T) {
	// object id manager with predefined ulids
	ids := []ulid.ULID{}
	for _, u := range []string{
		"01HVKKK6XFK7E245WTVEMXV1T9",
		"01HVKKK6XFK7E245WTVG8BGNXB",
		"01HVKKK6XFK7E245WTVJ5G90X9",
		"01HVKKK6XFK7E245WTVKH0QADF",
		"01HVKKK6XFK7E245WTVPZ40C1T"} {
		ids = append(ids, ulid.MustParse(u))
	}
	idManager := NewULIDManager()
	idManager.newID = func() ulid.ULID {
		require.True(t, len(ids) > 0)
		id := ids[0]
		ids = ids[1:]
		return id
	}

	// load package
	ePackage := loadPackage("library.simple.ecore")
	require.NotNil(t, ePackage)

	// load resource
	resourceSet := NewEResourceSetImpl()
	resourceSet.GetPackageRegistry().RegisterPackage(ePackage)
	resource := resourceSet.CreateResource(NewURI("testdata/library.simple.xml"))
	resource.SetObjectIDManager(idManager)
	resource.Load()
	require.NotNil(t, resource)
	require.True(t, resource.IsLoaded())
	require.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	require.True(t, resource.GetWarnings().Empty(), diagnosticError(resource.GetWarnings()))

	// w, err := os.Create("testdata/library.simple.ulids.sqlite")
	// require.NoError(t, err)
	// defer w.Close()
	w := &bytes.Buffer{}
	sqliteEncoder := NewSQLWriterEncoder(w, resource, map[string]any{SQL_OPTION_OBJECT_ID: "esyncID"})
	sqliteEncoder.EncodeResource()
	require.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))

	// compare expected and actual bytes
	requireSameDB(t, "testdata/library.simple.ulids.sqlite", w.Bytes())

}

func TestSQLEncoder_SimpleWithObjectID(t *testing.T) {
	// load package
	ePackage := loadPackage("library.simple.ecore")
	require.NotNil(t, ePackage)
	// id manager
	idManager := NewIncrementalIDManager()
	// load resource
	resourceSet := NewEResourceSetImpl()
	resourceSet.GetPackageRegistry().RegisterPackage(ePackage)
	resource := resourceSet.CreateResource(NewURI("testdata/library.simple.xml"))
	resource.SetObjectIDManager(idManager)
	resource.Load()
	require.NotNil(t, resource)
	require.True(t, resource.IsLoaded())
	require.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	require.True(t, resource.GetWarnings().Empty(), diagnosticError(resource.GetWarnings()))

	// w, err := os.Create("testdata/library.simple.ids.sqlite")
	// require.NoError(t, err)
	// defer w.Close()
	w := &bytes.Buffer{}
	sqliteEncoder := NewSQLWriterEncoder(w, resource, map[string]any{SQL_OPTION_OBJECT_ID: "objectID"})
	sqliteEncoder.EncodeResource()
	require.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))

	// compare expected and actual bytes
	requireSameDB(t, "testdata/library.simple.ids.sqlite", w.Bytes())
}

func TestSQLEncoder_SimpleWithContainerID(t *testing.T) {
	// load package
	ePackage := loadPackage("library.simple.ecore")
	require.NotNil(t, ePackage)

	// load resource
	xmlProcessor := NewXMLProcessor(XMLProcessorPackages([]EPackage{ePackage}))
	resource := xmlProcessor.LoadWithOptions(NewURI("testdata/library.simple.xml"), nil)
	require.NotNil(t, resource)
	require.True(t, resource.IsLoaded())
	require.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	require.True(t, resource.GetWarnings().Empty(), diagnosticError(resource.GetWarnings()))

	// w, err := os.Create("testdata/library.simple.container.sqlite")
	// require.NoError(t, err)
	// defer w.Close()
	w := &bytes.Buffer{}
	sqliteEncoder := NewSQLWriterEncoder(w, resource, map[string]any{SQL_OPTION_CONTAINER_ID: true})
	sqliteEncoder.EncodeResource()
	require.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))

	// compare expected and actual bytes
	requireSameDB(t, "testdata/library.simple.container.sqlite", w.Bytes())
}

func TestSQLEncoder_ComplexWithContainerID(t *testing.T) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	// load resource
	xmlProcessor := NewXMLProcessor(XMLProcessorPackages([]EPackage{ePackage}))
	resource := xmlProcessor.LoadWithOptions(NewURI("testdata/library.complex.xml"), nil)
	require.NotNil(t, resource)
	require.True(t, resource.IsLoaded())
	require.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	require.True(t, resource.GetWarnings().Empty(), diagnosticError(resource.GetWarnings()))

	// w, err := os.Create("testdata/library.complex.container.sqlite")
	// require.NoError(t, err)
	// defer w.Close()
	w := &bytes.Buffer{}
	sqliteEncoder := NewSQLWriterEncoder(w, resource, map[string]any{SQL_OPTION_CONTAINER_ID: true})
	sqliteEncoder.EncodeResource()
	require.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))

	// compare expected and actual bytes
	requireSameDB(t, "testdata/library.complex.container.sqlite", w.Bytes())
}
