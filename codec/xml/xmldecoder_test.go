package xml

import (
	"strings"
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const libraryDocument = `<?xml version="1.0" encoding="ISO-8859-1"?>
<!-- a library -->
<!DOCTYPE library>
<lib:Library xmlns:lib="http://lib" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
    xsi:schemaLocation="http://lib lib.xsd" xsi:noNamespaceSchemaLocation="lib.xsd"
    name="l" best="#//@books.1" favorites="#//@books.1 #//@books.0">
  <books name="b1" author="#//@writers.0"/>
  <books xsi:type="lib:Novel" name="b2"/>
  <writers name="w1" books="#//@books.0"/>
</lib:Library>
`

func TestXMLDecoder_Document(t *testing.T) {
	m := newLibraryModel()
	eResource := m.newResourceSet().CreateResource(ecore.NewURI("file:///library.xml"))
	decoder := NewDecoder(eResource, strings.NewReader(libraryDocument), nil)
	decoder.DecodeResource()
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	assert.Equal(t, "1.0", decoder.GetXMLVersion())
	assert.Equal(t, "ISO-8859-1", decoder.GetEncoding())

	require.Equal(t, 1, eResource.GetContents().Size())
	library := eResource.GetContents().Get(0).(ecore.EObject)
	books := library.EGet(m.books).(ecore.EList)
	writers := library.EGet(m.writers).(ecore.EList)
	require.Equal(t, 2, books.Size())
	require.Equal(t, 1, writers.Size())
	b1 := books.Get(0).(ecore.EObject)
	b2 := books.Get(1).(ecore.EObject)
	w1 := writers.Get(0).(ecore.EObject)
	assert.Equal(t, m.novel, b2.EClass())
	// forward references
	assert.Same(t, b2, library.EGet(m.best))
	assert.Equal(t, []any{b2, b1}, library.EGet(m.favorites).(ecore.EList).ToArray())
	assert.Same(t, w1, b1.EGet(m.author))
	assert.Equal(t, []any{b1}, w1.EGet(m.written).(ecore.EList).ToArray())
}

// the references resolved at the end of the document, the roots attached at the end of the
// document: the same objects
func TestXMLDecoder_DeferredOptions(t *testing.T) {
	m := newLibraryModel()
	expected := decodeResource(t, m.newResourceSet().CreateResource(ecore.NewURI("file:///library.xml")), libraryDocument, nil)
	for _, option := range []string{OPTION_DEFERRED_REFERENCE_RESOLUTION, OPTION_DEFERRED_ROOT_ATTACHMENT} {
		eResource := m.newResourceSet().CreateResource(ecore.NewURI("file:///library.xml"))
		library := decodeResource(t, eResource, libraryDocument, map[string]any{option: true})
		assert.True(t, ecore.Equals(expected, library), option)
	}
}

func TestXMLDecoder_DecodeObject(t *testing.T) {
	m := newLibraryModel()
	eResource := m.newResourceSet().CreateResource(ecore.NewURI("file:///library.xml"))
	eObject, err := NewDecoder(eResource, strings.NewReader(`<lib:Library xmlns:lib="http://lib" name="l"/>`), nil).DecodeObject()
	require.NoError(t, err)
	require.NotNil(t, eObject)
	assert.Equal(t, "l", eObject.EGet(m.name))
	assert.True(t, eResource.GetContents().Empty())

	// the first error is returned
	_, err = NewDecoder(eResource, strings.NewReader(`<lib:Library xmlns:lib="http://lib" unknown="u"><lib:Library`), nil).DecodeObject()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Feature unknown not found")
	assert.True(t, eResource.GetErrors().Empty())
}

func TestXMLDecoder_Errors(t *testing.T) {
	m := newLibraryModel()
	for _, test := range []struct {
		name     string
		document string
		message  string
	}{
		{"malformed", `<lib:Library xmlns:lib="http://lib">`, "XML syntax error on line 1: unexpected EOF"},
		{"unknown package", `<x:Library xmlns:x="http://unknown"/>`, "Package {'x'='http://unknown'} not found"},
		{"unknown class", `<lib:Nope xmlns:lib="http://lib"/>`, "Class {'http://lib':'Nope'} not found"},
		{"abstract class", `<lib:Named xmlns:lib="http://lib"/>`, "Class {'http://lib':'Named'} not found"},
		{"unknown attribute", `<lib:Library xmlns:lib="http://lib" nope="x"/>`, "Feature nope not found"},
		{"unknown element", `<lib:Library xmlns:lib="http://lib"><nope/></lib:Library>`, "Feature nope not found"},
		{"unknown type", `<lib:Library xmlns:lib="http://lib" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"><books xsi:type="lib:Nope"/></lib:Library>`, "Class {'http://lib':'Nope'} not found"},
		{"unknown type package", `<lib:Library xmlns:lib="http://lib" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"><books xsi:type="x:Book"/></lib:Library>`, "Package {'x'=''} not found"},
		{"unresolved reference", `<lib:Library xmlns:lib="http://lib" best="#//@books.3"/>`, "Unresolved reference '//@books.3'"},
	} {
		eResource := m.newResourceSet().CreateResource(ecore.NewURI("file:///library.xml"))
		NewDecoder(eResource, strings.NewReader(test.document), nil).DecodeResource()
		assert.Equal(t, test.message, diagnosticError(eResource.GetErrors()), test.name)
	}
}

// the references to objects of another document are proxies: href attributes of elements and
// uris of IDREF attributes, with or without the type of the proxy
func TestXMLDecoder_Proxies(t *testing.T) {
	m := newLibraryModel()
	eResource := m.newResourceSet().CreateResource(ecore.NewURI("file:///library.xml"))
	library := decodeResource(t, eResource, `<lib:Library xmlns:lib="http://lib" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
    favorites="other.xml#//@books.0 lib:Novel other.xml#//@books.1 #//@books.0" others="">
  <books name="b"/>
  <best href="other.xml#//@books.2"/>
  <any xsi:type="lib:Writer" href="other.xml#//@writers.0"/>
</lib:Library>`, nil)

	proxyURI := func(value any) string {
		eObject := value.(ecore.EObjectInternal)
		require.True(t, eObject.EIsProxy())
		return eObject.EProxyURI().String()
	}
	favorites := library.EGetResolve(m.favorites, false).(ecore.EList)
	require.Equal(t, 3, favorites.Size())
	assert.Equal(t, "file:///other.xml#//@books.0", proxyURI(favorites.Get(0)))
	assert.Equal(t, m.book, favorites.Get(0).(ecore.EObject).EClass())
	assert.Equal(t, "file:///other.xml#//@books.1", proxyURI(favorites.Get(1)))
	assert.Equal(t, m.novel, favorites.Get(1).(ecore.EObject).EClass())
	assert.Same(t, library.EGet(m.books).(ecore.EList).Get(0), favorites.Get(2))
	assert.Equal(t, "file:///other.xml#//@books.2", proxyURI(library.EGetResolve(m.best, false)))
	assert.Equal(t, "file:///other.xml#//@writers.0", proxyURI(library.EGetResolve(m.any, false)))
	assert.Equal(t, m.writer, library.EGetResolve(m.any, false).(ecore.EObject).EClass())
	assert.True(t, library.EGet(m.others).(ecore.EList).Empty())
}

// a proxy to an object of the same document is replaced by the object in the list holding it
func TestXMLDecoder_SameDocumentProxy(t *testing.T) {
	m := newLibraryModel()
	for name, document := range map[string]string{
		// the opposite of the book is unset when the proxy is resolved
		"proxy first": `<lib:Library xmlns:lib="http://lib">
  <writers name="w"><books href="#//@books.0"/></writers>
  <books name="b"/>
</lib:Library>`,
		// the book is in the list of the writer, after the proxy
		"book first": `<lib:Library xmlns:lib="http://lib">
  <writers name="w"><books href="#//@books.0"/></writers>
  <books name="b" author="#//@writers.0"/>
</lib:Library>`,
	} {
		eResource := m.newResourceSet().CreateResource(ecore.NewURI("file:///library.xml"))
		library := decodeResource(t, eResource, document, nil)
		book := library.EGet(m.books).(ecore.EList).Get(0).(ecore.EObject)
		writer := library.EGet(m.writers).(ecore.EList).Get(0).(ecore.EObject)
		assert.Equal(t, []any{book}, writer.EGetResolve(m.written, false).(ecore.EList).ToArray(), name)
		assert.Same(t, writer, book.EGet(m.author), name)
	}
}

// a prefix declared again in a nested element maps to the namespace declared there
func TestXMLDecoder_PrefixRedeclared(t *testing.T) {
	m := newLibraryModel()
	eResource := m.newResourceSet().CreateResource(ecore.NewURI("file:///library.xml"))
	library := decodeResource(t, eResource, `<lib:Library xmlns:lib="http://lib" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <books xmlns:lib="http://lib" xsi:type="lib:Novel" name="b"/>
  <books xmlns:lib="http://lib" xsi:type="lib:Novel" name="c"/>
</lib:Library>`, nil)
	books := library.EGet(m.books).(ecore.EList)
	require.Equal(t, 2, books.Size())
	assert.Equal(t, m.novel, books.Get(1).(ecore.EObject).EClass())
}

// the hooks of a decoder embedding Decoder are called instead of its own
type hookedDecoder struct {
	*Decoder
	attributes int
}

func (d *hookedDecoder) GetXSIType() string {
	return d.Decoder.GetXSIType()
}

func (d *hookedDecoder) HandleAttributes(eObject ecore.EObject) {
	d.attributes++
	d.Decoder.HandleAttributes(eObject)
}

func TestXMLDecoder_Interfaces(t *testing.T) {
	m := newLibraryModel()
	eResource := m.newResourceSet().CreateResource(ecore.NewURI("file:///library.xml"))
	decoder := &hookedDecoder{Decoder: NewDecoder(eResource, strings.NewReader(libraryDocument), nil)}
	decoder.SetInterfaces(decoder)
	decoder.DecodeResource()
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	assert.Equal(t, 4, decoder.attributes)
}
