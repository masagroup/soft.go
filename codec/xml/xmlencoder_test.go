package xml

import (
	"bytes"
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func encodeResource(t *testing.T, eResource ecore.EResource, options map[string]any) string {
	t.Helper()
	var buffer bytes.Buffer
	NewEncoder(eResource, &buffer, options).EncodeResource()
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	return buffer.String()
}

func decodeResource(t *testing.T, eResource ecore.EResource, document string, options map[string]any) ecore.EObject {
	t.Helper()
	NewDecoder(eResource, strings.NewReader(document), options).DecodeResource()
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	require.Equal(t, 1, eResource.GetContents().Size())
	return eResource.GetContents().Get(0).(ecore.EObject)
}

// a library whose references are in its resource: containments are elements, references are
// IDREF attributes, an object of a sub class of the containment type has its xsi:type
func newSameResourceLibrary(m *libraryModel) (library, b1, b2, w1 ecore.EObject) {
	library = m.create(m.library, "l")
	b1 = m.create(m.book, "b1")
	b1.ESet(m.pages, 10)
	b2 = m.create(m.novel, "b2")
	w1 = m.create(m.writer, "w1")
	library.EGet(m.books).(ecore.EList).AddAll(ecore.NewImmutableEList([]any{b1, b2}))
	library.EGet(m.writers).(ecore.EList).Add(w1)
	library.ESet(m.main, m.create(m.writer, "m"))
	library.EGet(m.tags).(ecore.EList).AddAll(ecore.NewImmutableEList([]any{"t1", "t2"}))
	library.ESet(m.best, b2)
	library.EGet(m.favorites).(ecore.EList).AddAll(ecore.NewImmutableEList([]any{b1, b2}))
	library.ESet(m.any, w1)
	w1.EGet(m.written).(ecore.EList).Add(b1)
	return
}

func TestXMLEncoder_SameResourceReferences(t *testing.T) {
	m := newLibraryModel()
	library, b1, _, w1 := newSameResourceLibrary(m)
	rs := m.newResourceSet()
	eResource := rs.CreateResource(ecore.NewURI("file:///library.xml"))
	eResource.GetContents().Add(library)

	document := encodeResource(t, eResource, nil)
	assert.True(t, strings.HasPrefix(document, `<?xml version="1.0" encoding="UTF-8"?>`))
	assert.Contains(t, document, `<lib:Library xmlns:lib="http://lib"`)
	assert.Contains(t, document, `xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"`)
	assert.Contains(t, document, `best="#//@books.1"`)
	assert.Contains(t, document, `favorites="#//@books.0 #//@books.1"`)
	assert.Contains(t, document, `any="#//@writers.0"`)
	assert.Contains(t, document, `<books name="b1" pages="10" author="#//@writers.0"/>`)
	assert.Contains(t, document, `<books xsi:type="lib:Novel" name="b2"/>`)
	assert.Contains(t, document, `<writers name="w1" books="#//@books.0"/>`)
	assert.Contains(t, document, `<main name="m"/>`)
	assert.Contains(t, document, `<tags>t1</tags>`)
	assert.Contains(t, document, `<tags>t2</tags>`)

	eDecoded := decodeResource(t, m.newResourceSet().CreateResource(ecore.NewURI("file:///library.xml")), document, nil)
	assert.True(t, ecore.Equals(library, eDecoded))
	books := eDecoded.EGet(m.books).(ecore.EList)
	writers := eDecoded.EGet(m.writers).(ecore.EList)
	assert.Equal(t, m.novel, books.Get(1).(ecore.EObject).EClass())
	assert.Same(t, books.Get(1), eDecoded.EGet(m.best))
	assert.Same(t, writers.Get(0), books.Get(0).(ecore.EObject).EGet(m.author))
	assert.Equal(t, b1.EGet(m.pages), books.Get(0).(ecore.EObject).EGet(m.pages))
	assert.Equal(t, w1.EGet(m.name), writers.Get(0).(ecore.EObject).EGet(m.name))
}

// the references to objects of another resource are href elements, with the xsi:type of an
// object of a type other than the abstract reference type
func TestXMLEncoder_CrossResourceReferences(t *testing.T) {
	m := newLibraryModel()
	rs := m.newResourceSet()
	other := m.create(m.library, "other")
	x := m.create(m.novel, "x")
	other.EGet(m.books).(ecore.EList).Add(x)
	eOther := rs.CreateResource(ecore.NewURI("file:///other.xml"))
	eOther.GetContents().Add(other)

	library := m.create(m.library, "l")
	library.ESet(m.best, x)
	library.ESet(m.any, x)
	library.EGet(m.favorites).(ecore.EList).Add(x)
	proxy := m.ePackage.GetEFactoryInstance().Create(m.book).(ecore.EObjectInternal)
	proxy.ESetProxyURI(ecore.NewURI("file:///missing.xml#//@books.0"))
	library.EGet(m.books).(ecore.EList).Add(proxy)
	eResource := rs.CreateResource(ecore.NewURI("file:///library.xml"))
	eResource.GetContents().Add(library)

	document := encodeResource(t, eResource, nil)
	assert.Contains(t, document, `<books href="missing.xml#//@books.0"/>`)
	assert.Contains(t, document, `<best href="other.xml#//@books.0"/>`)
	assert.Contains(t, document, `<favorites href="other.xml#//@books.0"/>`)
	assert.Contains(t, document, `<any xsi:type="lib:Novel" href="other.xml#//@books.0"/>`)

	// decoded as proxies, resolved in a resource set holding the other resource
	drs := m.newResourceSet()
	decodeResource(t, drs.CreateResource(ecore.NewURI("file:///other.xml")), encodeResource(t, eOther, nil), nil)
	eDecoded := decodeResource(t, drs.CreateResource(ecore.NewURI("file:///library.xml")), document, nil)
	assert.Equal(t, m.novel, eDecoded.EGetResolve(m.any, false).(ecore.EObject).EClass())
	assert.Equal(t, m.book, eDecoded.EGetResolve(m.best, false).(ecore.EObject).EClass())
	eBookProxy := eDecoded.EGetResolve(m.books, false).(ecore.EList).Get(0).(ecore.EObjectInternal)
	assert.Equal(t, "file:///missing.xml#//@books.0", eBookProxy.EProxyURI().String())
	for _, feature := range []ecore.EStructuralFeature{m.best, m.any} {
		eProxy := eDecoded.EGetResolve(feature, false).(ecore.EObjectInternal)
		require.True(t, eProxy.EIsProxy(), feature.GetName())
		assert.Equal(t, "file:///other.xml#//@books.0", eProxy.EProxyURI().String())
		resolved := eDecoded.EGet(feature).(ecore.EObject)
		assert.False(t, resolved.EIsProxy())
		assert.Equal(t, "x", resolved.EGet(m.name))
	}
	favorite := eDecoded.EGet(m.favorites).(ecore.EList).Get(0).(ecore.EObject)
	assert.Equal(t, "x", favorite.EGet(m.name))
}

// an object referenced in a list of cross references and contained by no resource is reported
func TestXMLEncoder_DanglingReference(t *testing.T) {
	m := newLibraryModel()
	rs := m.newResourceSet()
	other := m.create(m.library, "other")
	x := m.create(m.book, "x")
	other.EGet(m.books).(ecore.EList).Add(x)
	rs.CreateResource(ecore.NewURI("other.xml")).GetContents().Add(other)

	library := m.create(m.library, "l")
	library.EGet(m.favorites).(ecore.EList).AddAll(ecore.NewImmutableEList([]any{x, m.create(m.book, "dangling")}))
	eResource := rs.CreateResource(ecore.NewURI("file:///library.xml"))
	eResource.GetContents().Add(library)
	var buffer bytes.Buffer
	NewEncoder(eResource, &buffer, nil).EncodeResource()
	require.Equal(t, 1, eResource.GetErrors().Size())
	assert.Contains(t, diagnosticError(eResource.GetErrors()), "is not contained in a resource")
	assert.Contains(t, buffer.String(), `<favorites href="other.xml#//@books.0"/>`)
}

// the unsettable features set to nil are xsi:nil elements, an unsettable list set empty an empty
// attribute
func TestXMLEncoder_Nil(t *testing.T) {
	m := newLibraryModel()
	library := &nilSetEObject{EObject: m.create(m.library, "l"), nils: []ecore.EStructuralFeature{m.nick, m.optional, m.maybe}}
	others := library.EGet(m.others).(ecore.EList)
	book := m.create(m.book, "b")
	others.Add(book)
	others.Remove(book)
	library.EGet(m.tags).(ecore.EList).Add(nil)
	eResource := m.newResourceSet().CreateResource(ecore.NewURI("file:///library.xml"))

	var buffer bytes.Buffer
	require.NoError(t, NewEncoder(eResource, &buffer, nil).EncodeObject(library))
	document := buffer.String()
	assert.Contains(t, document, `<nick xsi:nil="true"/>`)
	assert.Contains(t, document, `<optional xsi:nil="true"/>`)
	assert.Contains(t, document, `<maybe xsi:nil="true"/>`)
	assert.Contains(t, document, `<tags xsi:nil="true"/>`)
	assert.Contains(t, document, `others=""`)
	assert.Contains(t, document, `xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"`)
}

// nilSetEObject is an object whose nils features are set to nil: a dynamic object does not keep an
// unsettable feature set to nil set, and reads the default value of an attribute
type nilSetEObject struct {
	ecore.EObject
	nils []ecore.EStructuralFeature
}

func (o *nilSetEObject) EIsSet(feature ecore.EStructuralFeature) bool {
	return slices.Contains(o.nils, feature) || o.EObject.EIsSet(feature)
}

func (o *nilSetEObject) EGetResolve(feature ecore.EStructuralFeature, resolve bool) any {
	if slices.Contains(o.nils, feature) {
		return nil
	}
	return o.EObject.EGetResolve(feature, resolve)
}

// the values of the unsettable features: an attribute, a containment and references
func TestXMLEncoder_UnsettableValues(t *testing.T) {
	m := newLibraryModel()
	library := m.create(m.library, "l")
	book := m.create(m.book, "b")
	library.EGet(m.books).(ecore.EList).Add(book)
	library.ESet(m.nick, "n")
	library.ESet(m.optional, m.create(m.writer, "o"))
	library.ESet(m.maybe, book)
	library.EGet(m.others).(ecore.EList).Add(book)
	eResource := m.newResourceSet().CreateResource(ecore.NewURI("file:///library.xml"))
	eResource.GetContents().Add(library)

	document := encodeResource(t, eResource, nil)
	assert.Contains(t, document, `nick="n"`)
	assert.Contains(t, document, `maybe="#//@books.0"`)
	assert.Contains(t, document, `others="#//@books.0"`)
	assert.Contains(t, document, `<optional name="o"/>`)

	eDecoded := decodeResource(t, m.newResourceSet().CreateResource(ecore.NewURI("file:///library.xml")), document, nil)
	assert.True(t, ecore.Equals(library, eDecoded))
}

// the references to objects of the resource of another resource
func TestXMLEncoder_UnsettableCrossReferences(t *testing.T) {
	m := newLibraryModel()
	rs := m.newResourceSet()
	other := m.create(m.library, "other")
	x := m.create(m.book, "x")
	other.EGet(m.books).(ecore.EList).Add(x)
	rs.CreateResource(ecore.NewURI("other.xml")).GetContents().Add(other)

	library := m.create(m.library, "l")
	library.ESet(m.maybe, x)
	library.EGet(m.others).(ecore.EList).Add(x)
	eResource := rs.CreateResource(ecore.NewURI("file:///library.xml"))
	eResource.GetContents().Add(library)

	document := encodeResource(t, eResource, nil)
	assert.Contains(t, document, `<maybe href="other.xml#//@books.0"/>`)
	assert.Contains(t, document, `<others href="other.xml#//@books.0"/>`)
}

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) {
	return 0, errors.New("write failed")
}

func TestXMLEncoder_EncodeObject(t *testing.T) {
	m := newLibraryModel()
	library := m.create(m.library, "l")
	eResource := m.newResourceSet().CreateResource(ecore.NewURI("file:///library.xml"))

	var buffer bytes.Buffer
	encoder := NewEncoder(eResource, &buffer, nil)
	encoder.SetXMLVersion("1.1")
	encoder.SetEncoding("ISO-8859-1")
	require.NoError(t, encoder.EncodeObject(library))
	assert.Equal(t, "<?xml version=\"1.1\" encoding=\"ISO-8859-1\"?>\n<lib:Library xmlns:lib=\"http://lib\" name=\"l\"/>\n", buffer.String())

	// the error of the writer is returned
	err := NewEncoder(eResource, failingWriter{}, nil).EncodeObject(library)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "write failed")
	assert.True(t, eResource.GetErrors().Empty())

	// and reported in the resource by EncodeResource
	eResource.GetContents().Add(library)
	NewEncoder(eResource, failingWriter{}, nil).EncodeResource()
	assert.Equal(t, "write failed", diagnosticError(eResource.GetErrors()))
}

func TestXMLEncoder_EmptyResource(t *testing.T) {
	eResource := ecore.NewEResourceImpl()
	var buffer bytes.Buffer
	NewEncoder(eResource, &buffer, nil).EncodeResource()
	assert.Empty(t, buffer.String())
}

func TestXMLEncoder_RootObjects(t *testing.T) {
	m := newLibraryModel()
	eResource := m.newResourceSet().CreateResource(ecore.NewURI("file:///library.xml"))
	eResource.GetContents().Add(m.create(m.library, "first"))
	second := m.create(m.library, "second")
	document := encodeResource(t, eResource, map[string]any{OPTION_ROOT_OBJECTS: ecore.NewImmutableEList([]any{second})})
	assert.Contains(t, document, `name="second"`)
	assert.NotContains(t, document, `name="first"`)
}

// the ids of the objects are written in the id attribute, and the references by id
func TestXMLEncoder_IDs(t *testing.T) {
	m := newLibraryModel()
	library := m.create(m.library, "l")
	book := m.create(m.book, "b")
	library.EGet(m.books).(ecore.EList).Add(book)
	library.ESet(m.best, book)
	eResource := m.newResourceSet().CreateResource(ecore.NewURI("file:///library.xml"))
	eResource.SetObjectIDManager(newStringIDManager())
	eResource.GetContents().Add(library)
	require.NoError(t, eResource.GetObjectIDManager().SetID(library, "L"))
	require.NoError(t, eResource.GetObjectIDManager().SetID(book, "B<1>"))
	options := map[string]any{OPTION_ID_ATTRIBUTE_NAME: "id"}

	document := encodeResource(t, eResource, options)
	assert.Contains(t, document, `id="L"`)
	assert.Contains(t, document, `<books id="B&lt;1&gt;" name="b"/>`)
	assert.Contains(t, document, `best="#B&lt;1&gt;"`)

	eDecodedResource := m.newResourceSet().CreateResource(ecore.NewURI("file:///library.xml"))
	eDecodedResource.SetObjectIDManager(newStringIDManager())
	eDecoded := decodeResource(t, eDecodedResource, document, options)
	eBook := eDecoded.EGet(m.books).(ecore.EList).Get(0).(ecore.EObject)
	assert.Equal(t, "L", eDecodedResource.GetObjectIDManager().GetID(eDecoded))
	assert.Equal(t, "B<1>", eDecodedResource.GetObjectIDManager().GetID(eBook))
	assert.Same(t, eBook, eDecoded.EGet(m.best))
}
