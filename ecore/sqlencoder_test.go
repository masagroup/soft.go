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
	actualConn, err := sqlite.OpenConn(":memory:")
	require.NoError(t, err)

	// set journal mode as rolling back
	// wal journal mode cannot be loaded in memory db
	actualBytes[18] = 0x01
	actualBytes[19] = 0x01

	// deserialize bytes into db
	err = actualConn.Deserialize("", actualBytes)
	require.NoError(t, err)
	defer func() {
		_ = actualConn.Close()
	}()

	// check that db are equal
	RequireEqualDB(t, expectedConn, actualConn)
}

var fileWriter bool = false

func testSQLEncoder(t require.TestingT, eResource EResource, dbPath string, options map[string]any) {
	if fileWriter {
		w, err := os.Create(dbPath)
		require.NoError(t, err)
		defer w.Close()
		sqliteEncoder := NewSQLWriterEncoder(w, eResource, options)
		sqliteEncoder.EncodeResource()
		require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	} else {
		w := &bytes.Buffer{}
		sqliteEncoder := NewSQLWriterEncoder(w, eResource, options)
		sqliteEncoder.EncodeResource()
		require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
		requireSameDB(t, dbPath, w.Bytes())
	}
}

func SQLEncoder_Complex(t *testing.T) {
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

	// test encoder
	testSQLEncoder(t, eResource, "testdata/library.complex.sqlite", nil)
}

func SQLEncoder_Complex_Memory(t *testing.T) {
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

	// test encoder
	testSQLEncoder(t, eResource, "testdata/library.complex.sqlite", map[string]any{SQL_OPTION_IN_MEMORY_DATABASE: true})
}

func BenchmarkSQLEncoder_Complex(b *testing.B) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(b, ePackage)

	// load resource
	xmlProcessor := NewXMLProcessor(XMLProcessorPackages([]EPackage{ePackage}))
	eResource := xmlProcessor.LoadWithOptions(NewURI("testdata/library.complex.xml"), nil)
	require.NotNil(b, eResource)
	require.True(b, eResource.IsLoaded())
	require.True(b, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	require.True(b, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	for n := 0; n < b.N; n++ {
		w := &bytes.Buffer{}
		sqliteEncoder := NewSQLWriterEncoder(w, eResource, map[string]any{SQL_OPTION_IN_MEMORY_DATABASE: false})
		sqliteEncoder.EncodeResource()
		require.True(b, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	}
}

func BenchmarkSQLEncoder_Complex_Memory(b *testing.B) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(b, ePackage)

	// load resource
	xmlProcessor := NewXMLProcessor(XMLProcessorPackages([]EPackage{ePackage}))
	eResource := xmlProcessor.LoadWithOptions(NewURI("testdata/library.complex.xml"), nil)
	require.NotNil(b, eResource)
	require.True(b, eResource.IsLoaded())
	require.True(b, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	require.True(b, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	for n := 0; n < b.N; n++ {
		w := &bytes.Buffer{}
		sqliteEncoder := NewSQLWriterEncoder(w, eResource, map[string]any{SQL_OPTION_IN_MEMORY_DATABASE: true})
		sqliteEncoder.EncodeResource()
		require.True(b, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	}
}

func SQLEncoder_DataList(t *testing.T) {
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

	// test encoder
	testSQLEncoder(t, eResource, "testdata/library.datalist.sqlite", nil)
}

func SQLEncoder_ComplexWithOwner(t *testing.T) {
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

	// test encoder
	testSQLEncoder(t, eResource, "testdata/library.owner.sqlite", nil)
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

	// test encoder
	testSQLEncoder(t, eResource, "testdata/emap.sqlite", nil)
}

func TestSQLEncoder_Simple(t *testing.T) {
	// load package
	ePackage := loadPackage("library.simple.ecore")
	require.NotNil(t, ePackage)

	// load resource
	xmlProcessor := NewXMLProcessor(XMLProcessorPackages([]EPackage{ePackage}))
	eResource := xmlProcessor.LoadWithOptions(NewURI("testdata/library.simple.xml"), nil)
	require.NotNil(t, eResource)
	require.True(t, eResource.IsLoaded())
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	require.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	// test encoder
	testSQLEncoder(t, eResource, "testdata/library.simple.sqlite", nil)
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
	eResourceSet := NewEResourceSetImpl()
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)
	eResource := eResourceSet.CreateResource(NewURI("testdata/library.simple.xml"))
	eResource.SetObjectIDManager(idManager)
	eResource.Load()
	require.NotNil(t, eResource)
	require.True(t, eResource.IsLoaded())
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	require.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	// test encoder
	testSQLEncoder(t, eResource, "testdata/library.simple.ulids.sqlite", map[string]any{SQL_OPTION_OBJECT_ID: "esyncID"})
}

func TestSQLEncoder_SimpleWithObjectID(t *testing.T) {
	// load package
	ePackage := loadPackage("library.simple.ecore")
	require.NotNil(t, ePackage)
	// id manager
	idManager := NewIncrementalIDManager()
	// load resource
	eResourceSet := NewEResourceSetImpl()
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)
	eResource := eResourceSet.CreateResource(NewURI("testdata/library.simple.xml"))
	eResource.SetObjectIDManager(idManager)
	eResource.Load()
	require.NotNil(t, eResource)
	require.True(t, eResource.IsLoaded())
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	require.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	// test encoder
	testSQLEncoder(t, eResource, "testdata/library.simple.ids.sqlite", map[string]any{SQL_OPTION_OBJECT_ID: "objectID"})
}

func TestSQLEncoder_SimpleWithContainerID(t *testing.T) {
	// load package
	ePackage := loadPackage("library.simple.ecore")
	require.NotNil(t, ePackage)

	// load resource
	xmlProcessor := NewXMLProcessor(XMLProcessorPackages([]EPackage{ePackage}))
	eResource := xmlProcessor.LoadWithOptions(NewURI("testdata/library.simple.xml"), nil)
	require.NotNil(t, eResource)
	require.True(t, eResource.IsLoaded())
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	require.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	// test encoder
	testSQLEncoder(t, eResource, "testdata/library.simple.container.sqlite", map[string]any{SQL_OPTION_CONTAINER_ID: true})
}

func TestSQLEncoder_ComplexWithContainerID(t *testing.T) {
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

	// test encoder
	testSQLEncoder(t, eResource, "testdata/library.complex.container.sqlite", map[string]any{SQL_OPTION_CONTAINER_ID: true})
}
