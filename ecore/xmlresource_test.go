package ecore

import (
	"io/ioutil"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	eResource := newXMIResourceImpl()
	eResource.SetURI(&url.URL{Path: "testdata/" + packageFileName})
	eResource.Load()
	if eResource.IsLoaded() && eResource.GetContents().Size() > 0 {
		ePackage, _ := eResource.GetContents().Get(0).(EPackage)
		ePackage.SetEFactoryInstance(NewEFactoryExt())
		return ePackage
	} else {
		return nil
	}
}

func TestXmlLoadLibraryNoRoot(t *testing.T) {
	ePackage := loadPackage("library.noroot.ecore")
	assert.NotNil(t, ePackage)

	// create a resource set
	eResourceSet := NewEResourceSetImpl()
	eResourceSet.GetPackageRegistry().RegisterPackage(GetPackage())
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)

	// create a resource
	eResource := newXMLResourceImpl()
	eResourceSet.GetResources().Add(eResource)

	// load it
	eResource.SetURI(&url.URL{Path: "testdata/library.noroot.xml"})
	eResource.Load()
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

func TestXmlLoadLibraryComplex(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	assert.NotNil(t, ePackage)

	// create a resource set
	eResourceSet := NewEResourceSetImpl()
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)

	// create a resource
	eResource := newXMLResourceImpl()
	eResourceSet.GetResources().Add(eResource)

	// load it
	eResource.SetURI(&url.URL{Path: "testdata/library.complex.xml"})
	eResource.Load()
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
	ePackage := loadPackage("library.complex.ecore")
	assert.NotNil(t, ePackage)

	// create a resource set
	eResourceSet := NewEResourceSetImpl()
	eResourceSet.GetPackageRegistry().RegisterPackage(GetPackage())
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)

	// create a resource
	eResource := newXMLResourceImpl()
	eResourceSet.GetResources().Add(eResource)

	extendedMetaData := NewExtendedMetaData()

	// load it
	eResource.SetURI(&url.URL{Path: "testdata/library.complex.noroot.xml"})
	eResource.LoadWithOptions(map[string]interface{}{OPTION_SUPPRESS_DOCUMENT_ROOT: true, OPTION_EXTENDED_META_DATA: extendedMetaData})
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
	ePackage := loadPackage("library.complex.ecore")
	assert.NotNil(t, ePackage)

	// create a resource set
	eResourceSet := NewEResourceSetImpl()
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)

	// create a resource
	eResource := newXMLResourceImpl()
	eResourceSet.GetResources().Add(eResource)

	// load it
	eResource.SetURI(&url.URL{Path: "testdata/library.complex.xml"})
	eResource.Load()
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	assert.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	var strbuff strings.Builder
	eResource.SaveWithWriter(&strbuff, nil)

	bytes, err := ioutil.ReadFile("testdata/library.complex.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(strbuff.String(), "\r\n", "\n"))
}

func TestXmlLoadSaveLibraryComplexWithOptions(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	assert.NotNil(t, ePackage)

	// create a resource set
	eResourceSet := NewEResourceSetImpl()
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)

	// create a resource
	eResource := newXMLResourceImpl()
	eResourceSet.GetResources().Add(eResource)

	extendedMetaData := NewExtendedMetaData()
	options := map[string]interface{}{OPTION_SUPPRESS_DOCUMENT_ROOT: true, OPTION_EXTENDED_META_DATA: extendedMetaData}
	// load it
	eResource.SetURI(&url.URL{Path: "testdata/library.complex.noroot.xml"})
	eResource.LoadWithOptions(options)
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	assert.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	// save it
	var strbuff strings.Builder
	eResource.SaveWithWriter(&strbuff, options)

	bytes, err := ioutil.ReadFile("testdata/library.complex.noroot.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(strbuff.String(), "\r\n", "\n"))
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
