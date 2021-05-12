package ecore

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func m(a, b interface{}) []interface{} {
	return []interface{}{a, b}
}

func TestXmlNamespacesNoContext(t *testing.T) {
	n := newXmlNamespaces()
	assert.Equal(t, m("", false), m(n.getURI("prefix")))
	assert.Equal(t, m("", false), m(n.getPrefix("uri")))
}

func TestXmlNamespacesEmpty(t *testing.T) {
	n := newXmlNamespaces()
	n.pushContext()
	assert.Equal(t, m("", false), m(n.getURI("prefix")))
	assert.Equal(t, m("", false), m(n.getPrefix("uri")))
	c := n.popContext()
	assert.Equal(t, 0, len(c))
}

func TestXmlNamespacesContext(t *testing.T) {
	n := newXmlNamespaces()
	n.pushContext()
	assert.False(t, n.declarePrefix("prefix", "uri"))
	assert.Equal(t, m("uri", true), m(n.getURI("prefix")))
	assert.Equal(t, m("prefix", true), m(n.getPrefix("uri")))

	n.pushContext()
	assert.False(t, n.declarePrefix("prefix", "uri2"))
	assert.Equal(t, m("uri2", true), m(n.getURI("prefix")))
	assert.Equal(t, m("prefix", true), m(n.getPrefix("uri2")))

	n.popContext()
	assert.Equal(t, m("uri", true), m(n.getURI("prefix")))
	assert.Equal(t, m("prefix", true), m(n.getPrefix("uri")))

	n.popContext()
	assert.Equal(t, m("", false), m(n.getURI("prefix")))
	assert.Equal(t, m("", false), m(n.getPrefix("uri")))
}

func TestXmlNamespacesContextRemap(t *testing.T) {
	n := newXmlNamespaces()
	n.pushContext()
	assert.False(t, n.declarePrefix("prefix", "uri"))
	assert.Equal(t, m("uri", true), m(n.getURI("prefix")))
	assert.Equal(t, m("prefix", true), m(n.getPrefix("uri")))

	assert.True(t, n.declarePrefix("prefix", "uri2"))
	assert.Equal(t, m("uri2", true), m(n.getURI("prefix")))
	assert.Equal(t, m("prefix", true), m(n.getPrefix("uri2")))
}

func TestXmlNamespacesContextNoRemap(t *testing.T) {
	n := newXmlNamespaces()
	n.pushContext()
	assert.False(t, n.declarePrefix("prefix", "uri"))
	assert.Equal(t, m("uri", true), m(n.getURI("prefix")))
	assert.Equal(t, m("prefix", true), m(n.getPrefix("uri")))

	n.pushContext()
	assert.False(t, n.declarePrefix("prefix", "uri2"))
	assert.Equal(t, m("uri2", true), m(n.getURI("prefix")))
	assert.Equal(t, m("prefix", true), m(n.getPrefix("uri2")))
}

func loadPackage(packageFileName string) EPackage {
	xmiProcessor := NewXMIProcessor()
	eResource := xmiProcessor.Load(&URI{Path: "testdata/" + packageFileName})
	if eResource.IsLoaded() && eResource.GetContents().Size() > 0 {
		ePackage, _ := eResource.GetContents().Get(0).(EPackage)
		ePackage.SetEFactoryInstance(NewEFactoryExt())
		return ePackage
	} else {
		return nil
	}
}

func TestXmlLoadLibraryNoRoot(t *testing.T) {
	// load package
	ePackage := loadPackage("library.noroot.ecore")
	assert.NotNil(t, ePackage)

	// load resource
	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	eResource, _ := xmlProcessor.Load(&URI{Path: "testdata/library.noroot.xml"}).(XMLResource)
	require.NotNil(t, eResource)
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	assert.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))
	assert.Equal(t, "1.0", eResource.GetXMLVersion())
	assert.Equal(t, "UTF-8", eResource.GetEncoding())

	// retrive library class & library name attribute
	eLibraryClass, _ := ePackage.GetEClassifier("Library").(EClass)
	assert.NotNil(t, eLibraryClass)
	eLibraryNameAttribute, _ := eLibraryClass.GetEStructuralFeatureFromName("name").(EAttribute)
	assert.NotNil(t, eLibraryNameAttribute)

	// check library name
	eLibrary, _ := eResource.GetContents().Get(0).(EObject)
	assert.Equal(t, "My Library", eLibrary.EGet(eLibraryNameAttribute))
}

func TestXmlLoadSaveLibraryNoRootWithOptions(t *testing.T) {
	// load package
	ePackage := loadPackage("library.noroot.ecore")
	assert.NotNil(t, ePackage)

	// load resource
	options := map[string]interface{}{OPTION_SUPPRESS_DOCUMENT_ROOT: true, OPTION_EXTENDED_META_DATA: NewExtendedMetaData()}
	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	eResource := xmlProcessor.LoadWithOptions(&URI{Path: "testdata/library.noroot.xml"}, options)
	require.NotNil(t, eResource)
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	assert.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	// save
	eResource.SetURI(&URI{Path: "testdata/library.noroot.result.xml"})
	xmlProcessor.SaveWithOptions(eResource, options)

	// result
	src, err := ioutil.ReadFile("testdata/library.noroot.xml")
	assert.Nil(t, err)

	result, err := ioutil.ReadFile("testdata/library.noroot.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(src), "\r\n", "\n"), strings.ReplaceAll(string(result), "\r\n", "\n"))
}

func TestXmlLoadSaveLibraryNoRootWithReaderWriter(t *testing.T) {
	// load package
	ePackage := loadPackage("library.noroot.ecore")
	assert.NotNil(t, ePackage)

	// xml processor
	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	options := map[string]interface{}{OPTION_SUPPRESS_DOCUMENT_ROOT: true, OPTION_EXTENDED_META_DATA: NewExtendedMetaData()}

	// load resource
	reader, error := os.Open("testdata/library.noroot.xml")
	require.Nil(t, error)
	eResource := xmlProcessor.LoadWithReader(reader, options)
	require.NotNil(t, eResource)
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	assert.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	// save
	writer, error := os.Create("testdata/library.noroot.result.xml")
	require.Nil(t, error)
	xmlProcessor.SaveWithWriter(writer, eResource, options)

	// result
	src, err := ioutil.ReadFile("testdata/library.noroot.xml")
	assert.Nil(t, err)

	result, err := ioutil.ReadFile("testdata/library.noroot.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(src), "\r\n", "\n"), strings.ReplaceAll(string(result), "\r\n", "\n"))
}

func TestXmlLoadLibraryComplex(t *testing.T) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	assert.NotNil(t, ePackage)

	// load resource
	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	eResource := xmlProcessor.Load(&URI{Path: "testdata/library.complex.xml"})
	require.NotNil(t, eResource)
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	assert.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	// retrieve document root class , library class & library name attribute
	eDocumentRootClass, _ := ePackage.GetEClassifier("DocumentRoot").(EClass)
	assert.NotNil(t, eDocumentRootClass)
	eDocumentRootLibraryFeature, _ := eDocumentRootClass.GetEStructuralFeatureFromName("library").(EReference)
	assert.NotNil(t, eDocumentRootLibraryFeature)
	eLibraryClass, _ := ePackage.GetEClassifier("Library").(EClass)
	assert.NotNil(t, eLibraryClass)
	eLibraryNameAttribute, _ := eLibraryClass.GetEStructuralFeatureFromName("name").(EAttribute)
	assert.NotNil(t, eLibraryNameAttribute)

	// check library name
	eDocumentRoot := eResource.GetContents().Get(0).(EObject)
	assert.NotNil(t, eDocumentRoot)
	eLibrary, _ := eDocumentRoot.EGet(eDocumentRootLibraryFeature).(EObject)
	assert.NotNil(t, eLibrary)
	assert.Equal(t, "My Library", eLibrary.EGet(eLibraryNameAttribute))
}

func TestXmlLoadLibraryComplexWithOptions(t *testing.T) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	assert.NotNil(t, ePackage)

	// load resource
	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	eResource := xmlProcessor.LoadWithOptions(&URI{Path: "testdata/library.complex.noroot.xml"}, map[string]interface{}{OPTION_SUPPRESS_DOCUMENT_ROOT: true, OPTION_EXTENDED_META_DATA: NewExtendedMetaData()})
	require.NotNil(t, eResource)
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	assert.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	// retrive library class & library name attribute
	eLibraryClass, _ := ePackage.GetEClassifier("Library").(EClass)
	assert.NotNil(t, eLibraryClass)
	eLibraryNameAttribute, _ := eLibraryClass.GetEStructuralFeatureFromName("name").(EAttribute)
	assert.NotNil(t, eLibraryNameAttribute)

	// check library name
	eLibrary, _ := eResource.GetContents().Get(0).(EObject)
	assert.Equal(t, "My Library", eLibrary.EGet(eLibraryNameAttribute))
}

func TestXmlLoadSaveLibraryComplex(t *testing.T) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	assert.NotNil(t, ePackage)

	// load resource
	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	eResource := xmlProcessor.Load(&URI{Path: "testdata/library.complex.xml"})
	require.NotNil(t, eResource)
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	assert.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))
	// save resource
	result := xmlProcessor.SaveToString(eResource, nil)

	bytes, err := ioutil.ReadFile("testdata/library.complex.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(result, "\r\n", "\n"))
}

func TestXmlLoadSaveLibraryComplexSubElement(t *testing.T) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	assert.NotNil(t, ePackage)

	// load resource
	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	eResource := xmlProcessor.Load(&URI{Path: "testdata/library.complex.xml"})
	require.NotNil(t, eResource)
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	assert.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	eObject := eResource.GetEObject("//@library/@employees.0")
	require.NotNil(t, eObject)
	eContainer := eObject.EContainer()
	require.NotNil(t, eContainer)

	// create a new resource
	eNewResource := eResource.GetResourceSet().CreateResource(&URI{Path: "testdata/library.complex.sub.xml"})
	// add object to new resource
	eNewResource.GetContents().Add(eObject)
	// save it
	result := xmlProcessor.SaveToString(eNewResource, nil)

	// check result
	bytes, err := ioutil.ReadFile("testdata/library.complex.sub.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(result, "\r\n", "\n"))

	// attach to original resource
	eLibraryClass, _ := ePackage.GetEClassifier("Library").(EClass)
	require.NotNil(t, eLibraryClass)
	eLibraryEmployeesFeature := eLibraryClass.GetEStructuralFeatureFromName("employees")
	require.NotNil(t, eLibraryEmployeesFeature)
	eList := eContainer.EGet(eLibraryEmployeesFeature).(EList)
	eList.Add(eObject)
	assert.Equal(t, eResource, eObject.EResource())
}

func TestXmlLoadSaveLibraryComplexWithOptions(t *testing.T) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	assert.NotNil(t, ePackage)

	options := map[string]interface{}{OPTION_SUPPRESS_DOCUMENT_ROOT: true, OPTION_EXTENDED_META_DATA: NewExtendedMetaData()}

	// load resource
	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	eResource := xmlProcessor.LoadWithOptions(&URI{Path: "testdata/library.complex.noroot.xml"}, options)
	require.NotNil(t, eResource)
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	assert.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	// save resource
	result := xmlProcessor.SaveToString(eResource, options)

	bytes, err := ioutil.ReadFile("testdata/library.complex.noroot.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(result, "\r\n", "\n"))
}

func TestXmlResourceIDManager(t *testing.T) {
	// load libray simple ecore	package
	ePackage := loadPackage("library.simple.ecore")
	assert.NotNil(t, ePackage)

	// create a resource with an id manager
	mockIDManager := &MockEObjectIDManager{}
	eResource := newXMLResourceImpl()
	eResource.SetObjectIDManager(mockIDManager)

	// create a library and add it to resource
	eFactory := ePackage.GetEFactoryInstance()
	eLibraryClass := ePackage.GetEClassifier("Library").(EClass)
	eLibrary := eFactory.Create(eLibraryClass)
	mockIDManager.On("Register", eLibrary).Once()
	eResource.GetContents().Add(eLibrary)
	mock.AssertExpectationsForObjects(t, mockIDManager)

	// create 2 books and add them to library
	eBookClass := ePackage.GetEClassifier("Book").(EClass)
	eLibraryBooksReference := eLibraryClass.GetEStructuralFeatureFromName("books").(EReference)
	eBookList := eLibrary.EGet(eLibraryBooksReference).(EList)
	eBook1 := eFactory.Create(eBookClass)
	eBook2 := eFactory.Create(eBookClass)
	mockIDManager.On("Register", eBook1).Once()
	mockIDManager.On("Register", eBook2).Once()
	eBookList.AddAll(NewImmutableEList([]interface{}{eBook1, eBook2}))
	mock.AssertExpectationsForObjects(t, mockIDManager)
}

func TestSerializationLoadSimpleInvalidXML(t *testing.T) {
	// load libray simple ecore	package
	ePackage := loadPackage("library.simple.ecore")
	assert.NotNil(t, ePackage)

	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	eResource := xmlProcessor.Load(&URI{Path: "testdata/library.simple.invalid.xml"})
	require.NotNil(t, eResource)
	assert.True(t, eResource.IsLoaded())
	assert.False(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
}

func TestSerializationLoadSimpleEscapeXML(t *testing.T) {
	// load libray simple ecore	package
	ePackage := loadPackage("library.simple.ecore")
	assert.NotNil(t, ePackage)

	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	eResource := xmlProcessor.Load(&URI{Path: "testdata/library.simple.escape.xml"})
	require.NotNil(t, eResource)
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))

	// retrive library class & library name attribute
	eLibraryClass, _ := ePackage.GetEClassifier("Library").(EClass)
	assert.NotNil(t, eLibraryClass)
	eLibraryLocationAttribute, _ := eLibraryClass.GetEStructuralFeatureFromName("location").(EAttribute)
	assert.NotNil(t, eLibraryLocationAttribute)

	// check library name
	eLibrary, _ := eResource.GetContents().Get(0).(EObject)
	assert.Equal(t, "a<b", eLibrary.EGet(eLibraryLocationAttribute))
}

func TestSerializationSaveSimpleEscapeXML(t *testing.T) {
	// load libray simple ecore	package
	ePackage := loadPackage("library.simple.ecore")
	assert.NotNil(t, ePackage)

	// retrive library class & library name attribute
	eLibraryClass, _ := ePackage.GetEClassifier("Library").(EClass)
	assert.NotNil(t, eLibraryClass)
	eLibraryLocationAttribute, _ := eLibraryClass.GetEStructuralFeatureFromName("location").(EAttribute)
	assert.NotNil(t, eLibraryLocationAttribute)

	eFactory := ePackage.GetEFactoryInstance()
	eLibrary := eFactory.Create(eLibraryClass)
	eLibrary.ESet(eLibraryLocationAttribute, "a<b")

	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	eResource := xmlProcessor.GetResourceSet().CreateResource(&URI{Path: "testdata/library.simple.escape.output.xml"})
	eResource.GetContents().Add(eLibrary)
	result := xmlProcessor.SaveToString(eResource, nil)

	bytes, err := ioutil.ReadFile("testdata/library.simple.escape.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(result, "\r\n", "\n"))
}

func TestSerializationLoadSimpleXMLWithIDs(t *testing.T) {
	idManager := NewIncrementalIDManager()

	// load libray simple ecore	package
	ePackage := loadPackage("library.simple.ecore")
	assert.NotNil(t, ePackage)

	eResourceSet := NewEResourceSetImpl()
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)
	eResource := eResourceSet.CreateResource(&URI{Path: "testdata/library.simple.ids.xml"})
	require.NotNil(t, eResource)
	eResource.SetObjectIDManager(idManager)
	eResource.LoadWithOptions(map[string]interface{}{OPTION_ID_ATTRIBUTE_NAME: "id"})
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))

	// retrive library class & library name attribute
	libraryClass, _ := ePackage.GetEClassifier("Library").(EClass)
	require.NotNil(t, libraryClass)
	libraryBooksFeature := libraryClass.GetEStructuralFeatureFromName("books")
	require.NotNil(t, libraryBooksFeature)

	require.Equal(t, 1, eResource.GetContents().Size())
	eLibrary, _ := eResource.GetContents().Get(0).(EObject)
	require.NotNil(t, eLibrary)
	assert.Equal(t, 0, idManager.GetID(eLibrary))

	eBooks, _ := eLibrary.EGet(libraryBooksFeature).(EList)
	require.NotNil(t, eBooks)
	require.Equal(t, 4, eBooks.Size())
	assert.Equal(t, 1, idManager.GetID(eBooks.Get(0).(EObject)))
	assert.Equal(t, 2, idManager.GetID(eBooks.Get(1).(EObject)))
	assert.Equal(t, 3, idManager.GetID(eBooks.Get(2).(EObject)))
	assert.Equal(t, 4, idManager.GetID(eBooks.Get(3).(EObject)))
}

func TestSerializationSaveSimpleXMLWithIDs(t *testing.T) {

	// load libray simple ecore	package
	ePackage := loadPackage("library.simple.ecore")
	assert.NotNil(t, ePackage)

	eResourceSet := NewEResourceSetImpl()
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)
	eResource := eResourceSet.CreateResource(&URI{Path: "testdata/library.simple.xml"})
	require.NotNil(t, eResource)
	eResource.SetObjectIDManager(NewIncrementalIDManager())
	eResource.Load()
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))

	var strbuff strings.Builder
	eResource.SaveWithWriter(&strbuff, map[string]interface{}{OPTION_ID_ATTRIBUTE_NAME: "id"})

	bytes, err := ioutil.ReadFile("testdata/library.simple.ids.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(strbuff.String(), "\r\n", "\n"))
}

func TestSerializationSaveSimpleXMLRootObjects(t *testing.T) {
	// load libray simple ecore	package
	ePackage := loadPackage("library.simple.ecore")
	assert.NotNil(t, ePackage)

	// load model file
	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	eResource := xmlProcessor.Load(&URI{Path: "testdata/library.simple.xml"})
	require.NotNil(t, eResource)
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))

	// retrieve second book
	libraryClass, _ := ePackage.GetEClassifier("Library").(EClass)
	require.NotNil(t, libraryClass)
	libraryBooksFeature := libraryClass.GetEStructuralFeatureFromName("books")
	require.NotNil(t, libraryBooksFeature)

	require.Equal(t, 1, eResource.GetContents().Size())
	eLibrary, _ := eResource.GetContents().Get(0).(EObject)
	require.NotNil(t, eLibrary)

	eBooks, _ := eLibrary.EGet(libraryBooksFeature).(EList)
	require.NotNil(t, eBooks)
	require.Equal(t, 4, eBooks.Size())
	eBook := eBooks.Get(1).(EObject)

	// save it now
	var strbuff strings.Builder
	eResource.SaveWithWriter(&strbuff, map[string]interface{}{OPTION_ROOT_OBJECTS: NewImmutableEList([]interface{}{eBook})})

	bytes, err := ioutil.ReadFile("testdata/book.simple.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(strbuff.String(), "\r\n", "\n"))
}
