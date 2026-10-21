package xml

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func fileURI(path string) *ecore.URI {
	return ecore.NewURI("file:///" + strings.TrimPrefix(filepath.ToSlash(path), "/"))
}

func TestXMLProcessor_SaveLoad(t *testing.T) {
	m := newLibraryModel()
	library, _, _, _ := newSameResourceLibrary(m)
	p := NewProcessor(ProcessorPackages([]ecore.EPackage{m.ePackage}))
	assert.NotNil(t, p.GetExtendedMetaData())
	path := filepath.Join(t.TempDir(), "library.xml")

	// SaveObject saves a copy of the object
	eResource := p.SaveObject(fileURI(path), library)
	require.NotNil(t, eResource)
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	assert.Nil(t, library.EResource())
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Contains(t, string(content), `<lib:Library xmlns:lib="http://lib"`)

	eLoaded := p.Load(fileURI(path))
	require.True(t, eLoaded.IsLoaded())
	require.True(t, eLoaded.GetErrors().Empty(), diagnosticError(eLoaded.GetErrors()))
	assert.True(t, ecore.Equals(library, eLoaded.GetContents().Get(0).(ecore.EObject)))

	// Save of a loaded resource writes the same content
	require.NoError(t, os.Remove(path))
	p.Save(eLoaded)
	require.True(t, eLoaded.GetErrors().Empty(), diagnosticError(eLoaded.GetErrors()))
	saved, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, string(content), string(saved))

	// to a string or a writer, and from a reader
	document := p.SaveToString(eLoaded, nil)
	assert.Equal(t, string(content), document)
	var buffer bytes.Buffer
	p.SaveWithWriter(&buffer, eLoaded, nil)
	assert.Equal(t, document, buffer.String())
	eRead := p.LoadWithReader(strings.NewReader(document), nil)
	require.True(t, eRead.GetErrors().Empty(), diagnosticError(eRead.GetErrors()))
	assert.True(t, ecore.Equals(library, eRead.GetContents().Get(0).(ecore.EObject)))
}

func TestXMLProcessor_ResourceSet(t *testing.T) {
	m := newLibraryModel()
	rs := m.newResourceSet()
	p := NewProcessor(ProcessorResourceSet(rs))
	assert.Same(t, rs, p.GetResourceSet())
	eResource := p.LoadWithReader(strings.NewReader(`<lib:Library xmlns:lib="http://lib" name="l"/>`), nil)
	assert.Same(t, rs, eResource.GetResourceSet())
	assert.Equal(t, 1, eResource.GetContents().Size())

	// a missing file is reported in the resource
	eMissing := p.LoadWithOptions(fileURI(filepath.Join(t.TempDir(), "missing.xml")), nil)
	assert.False(t, eMissing.GetErrors().Empty())
}

// documentRootModel is the library model with a document root, as generated from a schema: its
// features are the root elements of the documents
type documentRootModel struct {
	*libraryModel
	documentRoot   ecore.EClass
	rootLibrary    ecore.EReference
	prefixMap      ecore.EReference
	schemaLocation ecore.EReference
}

func newDocumentRootModel() *documentRootModel {
	f := ecore.GetFactory()
	m := &documentRootModel{libraryModel: newLibraryModel()}
	setMetaData := func(eElement ecore.EModelElement, details map[string]string) {
		eAnnotation := f.CreateEAnnotation()
		eAnnotation.SetSource(annotationURI)
		for key, value := range details {
			eAnnotation.GetDetails().Put(key, value)
		}
		eElement.GetEAnnotations().Add(eAnnotation)
	}
	newReference := func(name string, eType ecore.EClass, details map[string]string) ecore.EReference {
		eReference := f.CreateEReference()
		eReference.SetName(name)
		eReference.SetEType(eType)
		eReference.SetContainment(true)
		setMetaData(eReference, details)
		m.documentRoot.GetEStructuralFeatures().Add(eReference)
		return eReference
	}
	m.documentRoot = f.CreateEClass()
	m.documentRoot.SetName("DocumentRoot")
	setMetaData(m.documentRoot, map[string]string{"name": ""})
	m.ePackage.GetEClassifiers().Add(m.documentRoot)
	m.prefixMap = newReference("xMLNSPrefixMap", ecore.GetPackage().GetEStringToStringMapEntry(), map[string]string{"name": "xmlns:prefix"})
	m.prefixMap.SetUpperBound(ecore.UNBOUNDED_MULTIPLICITY)
	m.prefixMap.SetTransient(true)
	m.schemaLocation = newReference("xSISchemaLocation", ecore.GetPackage().GetEStringToStringMapEntry(), map[string]string{"name": "xsi:schemaLocation"})
	m.schemaLocation.SetUpperBound(ecore.UNBOUNDED_MULTIPLICITY)
	m.schemaLocation.SetTransient(true)
	m.rootLibrary = newReference("library", m.library, map[string]string{"name": "library", "namespace": "##targetNamespace"})
	return m
}

func TestXMLDocumentRoot(t *testing.T) {
	m := newDocumentRootModel()
	emd := NewExtendedMetaData()
	assert.Same(t, m.documentRoot, emd.GetDocumentRoot(m.ePackage))
	assert.Same(t, m.prefixMap, emd.GetXMLNSPrefixMapFeature(m.documentRoot))
	assert.Same(t, m.schemaLocation, emd.GetXSISchemaLocationMapFeature(m.documentRoot))
	assert.Nil(t, emd.GetXSISchemaLocationMapFeature(m.library))
	assert.Equal(t, "http://lib", emd.GetNamespace(m.rootLibrary))

	document := `<?xml version="1.0" encoding="UTF-8"?>
<lib:library xmlns:lib="http://lib" name="l">
  <books name="b"/>
</lib:library>
`
	options := map[string]any{OPTION_EXTENDED_META_DATA: emd}

	// decoded in a document root holding the prefixes of the document
	eResource := m.newResourceSet().CreateResource(ecore.NewURI("library.xml"))
	root := decodeResource(t, eResource, document, options)
	require.Same(t, m.documentRoot, root.EClass())
	library := root.EGet(m.rootLibrary).(ecore.EObject)
	assert.Equal(t, "l", library.EGet(m.name))
	assert.Equal(t, "http://lib", root.EGet(m.prefixMap).(ecore.EMap).GetValue("lib"))

	// encoded with the prefixes of the document root
	root.EGet(m.prefixMap).(ecore.EMap).Put("lib", "http://lib")
	assert.Equal(t, document, encodeResource(t, eResource, options))

	// the document root is suppressed: the library is the root
	eSuppressed := m.newResourceSet().CreateResource(ecore.NewURI("library.xml"))
	library = decodeResource(t, eSuppressed, document, map[string]any{OPTION_EXTENDED_META_DATA: emd, OPTION_SUPPRESS_DOCUMENT_ROOT: true})
	require.Same(t, m.library, library.EClass())
	assert.Nil(t, library.EContainer())
	assert.Equal(t, 1, library.EGet(m.books).(ecore.EList).Size())

	// a library is encoded with the name of its document root feature
	assert.Equal(t, document, encodeResource(t, eSuppressed, options))
}
