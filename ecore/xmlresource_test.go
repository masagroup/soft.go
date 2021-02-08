package ecore

import (
	"io/ioutil"
	"net/url"
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
	eResource := xmiProcessor.Load(&url.URL{Path: "testdata/" + packageFileName})
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
	eResource := xmlProcessor.Load(&url.URL{Path: "testdata/library.noroot.xml"})
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

func TestXmlLoadSaveLibraryNoRootWithOptions(t *testing.T) {
	// load package
	ePackage := loadPackage("library.noroot.ecore")
	assert.NotNil(t, ePackage)

	// load resource
	options := map[string]interface{}{OPTION_SUPPRESS_DOCUMENT_ROOT: true, OPTION_EXTENDED_META_DATA: NewExtendedMetaData()}
	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	eResource := xmlProcessor.LoadWithOptions(&url.URL{Path: "testdata/library.noroot.xml"}, options)
	require.NotNil(t, eResource)
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	assert.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	// save
	eResource.SetURI(&url.URL{Path: "testdata/library.noroot.result.xml"})
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
	eResource := xmlProcessor.Load(&url.URL{Path: "testdata/library.complex.xml"})
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
	eResource := xmlProcessor.LoadWithOptions(&url.URL{Path: "testdata/library.complex.noroot.xml"}, map[string]interface{}{OPTION_SUPPRESS_DOCUMENT_ROOT: true, OPTION_EXTENDED_META_DATA: NewExtendedMetaData()})
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
	eResource := xmlProcessor.Load(&url.URL{Path: "testdata/library.complex.xml"})
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

func TestXmlLoadSaveLibraryComplexWithOptions(t *testing.T) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	assert.NotNil(t, ePackage)

	options := map[string]interface{}{OPTION_SUPPRESS_DOCUMENT_ROOT: true, OPTION_EXTENDED_META_DATA: NewExtendedMetaData()}

	// load resource
	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	eResource := xmlProcessor.LoadWithOptions(&url.URL{Path: "testdata/library.complex.noroot.xml"}, options)
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
	mockIDManager := &MockEResourceIDManager{}
	eResource := newXMLResourceImpl()
	eResource.SetIDManager(mockIDManager)

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
