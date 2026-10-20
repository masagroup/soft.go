package json

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const librarySimpleNsURI = "http:///org/eclipse/emf/examples/library/library.simple.ecore/1.0.0"

// loadXMLResource loads an XML model of testdata with its packages, all in one resource set so
// that the references between packages resolve to the same classes. With ids, the resource has an
// incremental ID manager, as the model of testdata/orders.json.
func loadXMLResource(t *testing.T, model string, ids bool, packages ...string) (ecore.EResource, []ecore.EPackage) {
	eResourceSet := ecore.NewEResourceSetImpl()
	ePackages := []ecore.EPackage{}
	for _, p := range packages {
		_, ePackage := loadTestPackage(t, eResourceSet, ecore.NewURI("testdata/"+p))
		ePackages = append(ePackages, ePackage)
	}
	eResource := eResourceSet.CreateResource(ecore.NewURI("testdata/" + model))
	if ids {
		eResource.SetObjectIDManager(ecore.NewIncrementalIDManager())
	}
	eResource.Load()
	require.True(t, eResource.IsLoaded())
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	return eResource, ePackages
}

// newJSONResource creates a resource with uri in a resource set holding packages.
func newJSONResource(uri *ecore.URI, packages []ecore.EPackage) ecore.EResource {
	eResourceSet := ecore.NewEResourceSetImpl()
	for _, ePackage := range packages {
		eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)
	}
	return eResourceSet.CreateResource(uri)
}

func encode(t *testing.T, eResource ecore.EResource, options map[string]any) []byte {
	buffer := &bytes.Buffer{}
	NewEncoder(eResource, buffer, options).EncodeResource()
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	return buffer.Bytes()
}

func decode(eResource ecore.EResource, data string, options map[string]any) {
	NewDecoder(eResource, strings.NewReader(data), options).DecodeResource()
}

func TestJSONDecoder_RoundTrip(t *testing.T) {
	for _, test := range []struct {
		model    string
		json     string
		ids      bool
		packages []string
	}{
		{"library.simple.xml", "library.simple.json", false, []string{"library.simple.ecore"}},
		{"library.complex.xml", "library.complex.json", false, []string{"library.complex.ecore"}},
		{"orders.xml", "orders.json", true, []string{"shop.ecore", "orders.ecore"}},
	} {
		t.Run(test.model, func(t *testing.T) {
			var options map[string]any
			if test.ids {
				options = map[string]any{OPTION_ID_ATTRIBUTE_NAME: "id"}
			}
			xmlResource, ePackages := loadXMLResource(t, test.model, test.ids, test.packages...)
			data := encode(t, xmlResource, options)
			expected, err := os.ReadFile("testdata/" + test.json)
			require.NoError(t, err)
			require.Equal(t, strings.ReplaceAll(string(expected), "\r\n", "\n"), string(data))

			// decode in a resource of the same URI: the relative references resolve the same way
			jsonResource := newJSONResource(xmlResource.GetURI(), ePackages)
			if test.ids {
				jsonResource.SetObjectIDManager(ecore.NewIncrementalIDManager())
			}
			decode(jsonResource, string(data), options)
			require.True(t, jsonResource.GetErrors().Empty(), diagnosticError(jsonResource.GetErrors()))
			// the transient features, as the namespaces of an XML document root, are not written
			unsetTransientFeatures(xmlResource)
			assert.True(t, ecore.EqualsAll(xmlResource.GetContents(), jsonResource.GetContents()))

			// and encode it again
			assert.Equal(t, string(data), string(encode(t, jsonResource, options)))
		})
	}
}

func TestJSONDecoder_LoadFromRegistry(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)
	eResource := newJSONResource(ecore.CreateFileURI("testdata/library.complex.json"), []ecore.EPackage{ePackage})
	eResource.Load()
	require.True(t, eResource.IsLoaded())
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	require.Equal(t, 1, eResource.GetContents().Size())

	// the books of the writer are references to the books of the library, resolved
	eDocumentRoot := eResource.GetContents().Get(0).(ecore.EObject)
	eLibrary := eDocumentRoot.EGet(eDocumentRoot.EClass().GetEStructuralFeatureFromName("library")).(ecore.EObject)
	eLibraryClass := eLibrary.EClass()
	eWriter := eLibrary.EGet(eLibraryClass.GetEStructuralFeatureFromName("writers")).(ecore.EList).Get(0).(ecore.EObject)
	eWriterBooks := eWriter.EGetResolve(eWriter.EClass().GetEStructuralFeatureFromName("books"), false).(ecore.EList)
	eLibraryBooks := eLibrary.EGet(eLibraryClass.GetEStructuralFeatureFromName("books")).(ecore.EList)
	require.Equal(t, 2, eWriterBooks.Size())
	for i := range 2 {
		book := eWriterBooks.Get(i).(ecore.EObject)
		assert.False(t, book.EIsProxy())
		assert.Same(t, eLibraryBooks.Get(i), book)
	}
}

func TestJSONDecoder_Proxy(t *testing.T) {
	ePackages := []ecore.EPackage{loadPackage("shop.ecore"), loadPackage("orders.ecore")}
	uri := ecore.CreateFileURI("testdata/orders.json")
	eResource := newJSONResource(uri, ePackages)
	data, err := os.ReadFile("testdata/orders.json")
	require.NoError(t, err)
	eResource.SetObjectIDManager(ecore.NewIncrementalIDManager())
	decode(eResource, string(data), map[string]any{OPTION_ID_ATTRIBUTE_NAME: "id"})
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))

	eOrders := eResource.GetContents().Get(0).(ecore.EObject)
	eOrder := eOrders.EGet(eOrders.EClass().GetEStructuralFeatureFromName("order")).(ecore.EList).Get(0).(ecore.EObject)
	eProduct := eOrder.EGetResolve(eOrder.EClass().GetEStructuralFeatureFromName("product"), false).(ecore.EObject)
	require.True(t, eProduct.EIsProxy())
	assert.Equal(t, "Product", eProduct.EClass().GetName())
	assert.Equal(t, uri.Resolve(ecore.NewURI("shop.xml#//@products.0")).String(), eProduct.(ecore.EObjectInternal).EProxyURI().String())
}

func TestJSONDecoder_IDs(t *testing.T) {
	ePackage := loadPackage("library.simple.ecore")
	require.NotNil(t, ePackage)
	xmlResource := newJSONResource(ecore.CreateFileURI("testdata/library.simple.xml"), []ecore.EPackage{ePackage})
	xmlResource.SetObjectIDManager(ecore.NewULIDManager())
	xmlResource.Load()
	require.True(t, xmlResource.GetErrors().Empty(), diagnosticError(xmlResource.GetErrors()))

	options := map[string]any{OPTION_ID_ATTRIBUTE_NAME: "id"}
	data := encode(t, xmlResource, options)
	assert.Contains(t, string(data), `"id":"`)

	jsonResource := newJSONResource(ecore.CreateFileURI("testdata/library.simple.json"), []ecore.EPackage{ePackage})
	jsonResource.SetObjectIDManager(ecore.NewULIDManager())
	decode(jsonResource, string(data), options)
	require.True(t, jsonResource.GetErrors().Empty(), diagnosticError(jsonResource.GetErrors()))

	xmlObjects := xmlResource.GetAllContents()
	jsonObjects := jsonResource.GetAllContents()
	count := 0
	for xmlObjects.HasNext() {
		require.True(t, jsonObjects.HasNext())
		xmlObject := xmlObjects.Next().(ecore.EObject)
		jsonObject := jsonObjects.Next().(ecore.EObject)
		id := xmlResource.GetObjectIDManager().GetID(xmlObject)
		require.NotNil(t, id)
		assert.Equal(t, id, jsonResource.GetObjectIDManager().GetID(jsonObject))
		count++
	}
	assert.False(t, jsonObjects.HasNext())
	assert.Greater(t, count, 1)
}

func TestJSONDecoder_Errors(t *testing.T) {
	ePackage := loadPackage("library.simple.ecore")
	require.NotNil(t, ePackage)
	library := `{"eClass":"` + librarySimpleNsURI + `#//Library"`
	for _, test := range []struct {
		name    string
		data    string
		message string
		root    bool // the root object is decoded
	}{
		{"syntax", `{"eClass":`, "unexpected EOF", false},
		{"invalid", `{"eClass":"x"}`, "Invalid class 'x'", false},
		{"unknown package", `{"eClass":"http://unknown#//Library"}`, "Unknown package 'http://unknown'", false},
		{"unknown class", `{"eClass":"` + librarySimpleNsURI + `#//Unknown"}`, "Unknown class 'Unknown'", false},
		{"missing class", `{"owner":"Owner"}`, "Missing eClass", false},
		{"unknown feature", library + `,"unknown":{"a":[1,{"b":2}]},"owner":"Owner"}`, "Unknown feature 'unknown' in class 'Library'", true},
		{"wrong type", library + `,"books":"x"}`, "Expected '[', found string 'x'", false},
		{"trailing content", library + `}{}`, "Unexpected content after the root object", true},
	} {
		t.Run(test.name, func(t *testing.T) {
			eResource := newJSONResource(ecore.NewURI("library.json"), []ecore.EPackage{ePackage})
			decode(eResource, test.data, nil)
			require.False(t, eResource.GetErrors().Empty())
			assert.Contains(t, diagnosticError(eResource.GetErrors()), test.message)
			assert.Equal(t, test.root, eResource.GetContents().Size() == 1)
			if test.name == "unknown feature" {
				eLibrary := eResource.GetContents().Get(0).(ecore.EObject)
				assert.Equal(t, "Owner", eLibrary.EGet(eLibrary.EClass().GetEStructuralFeatureFromName("owner")))
			}
		})
	}
}

func TestJSONDecoder_UnresolvedReference(t *testing.T) {
	ePackages := []ecore.EPackage{loadPackage("shop.ecore"), loadPackage("orders.ecore")}
	data, err := os.ReadFile("testdata/orders.json")
	require.NoError(t, err)
	// a reference to a product of the document, which has none
	broken := strings.Replace(string(data), `"eRef":"shop.xml#//@products.0"`, `"eRef":"#//@products.0"`, 1)
	eResource := newJSONResource(ecore.NewURI("orders.json"), ePackages)
	eResource.SetObjectIDManager(ecore.NewIncrementalIDManager())
	decode(eResource, broken, map[string]any{OPTION_ID_ATTRIBUTE_NAME: "id"})
	require.Equal(t, 1, eResource.GetErrors().Size())
	assert.Equal(t, "Unresolved reference '//@products.0'", diagnosticError(eResource.GetErrors()))
}

func TestJSONDecoder_DecodeObject(t *testing.T) {
	ePackage := loadPackage("library.simple.ecore")
	require.NotNil(t, ePackage)
	eResource := newJSONResource(ecore.NewURI("library.json"), []ecore.EPackage{ePackage})

	data, err := os.ReadFile("testdata/library.simple.json")
	require.NoError(t, err)
	eObject, err := NewDecoder(eResource, bytes.NewReader(data), nil).DecodeObject()
	require.NoError(t, err)
	require.NotNil(t, eObject)
	assert.Equal(t, "Library", eObject.EClass().GetName())
	assert.Nil(t, eObject.EResource())

	_, err = NewDecoder(eResource, strings.NewReader(`{"eClass":"x"}`), nil).DecodeObject()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid class 'x'")
}

func TestJSONWriter_Escaping(t *testing.T) {
	for _, s := range []string{
		"", "plain", `quote " backslash \ slash /`, "<html> & 'quotes'",
		"control \b\f\n\r\t \x00 \x1f \x7f", "accents éàü", "emoji 😀", "separators    ",
		"invalid \xff utf8 \xe2\x82", "mixed \x01é \"",
	} {
		buffer := &bytes.Buffer{}
		encoder := json.NewEncoder(buffer)
		encoder.SetEscapeHTML(false)
		require.NoError(t, encoder.Encode(s))
		expected := strings.TrimSuffix(buffer.String(), "\n")
		assert.Equal(t, expected, string(appendString(nil, s)), s)
	}
}

func TestJSONWriter_Structure(t *testing.T) {
	buffer := &bytes.Buffer{}
	w := newJSONWriter(buffer)
	w.beginObject()
	w.key("a")
	w.stringValue("1")
	w.key("b")
	w.beginArray()
	w.stringValue("x")
	w.beginObject()
	w.endObject()
	w.beginArray()
	w.endArray()
	w.endArray()
	w.key("c")
	w.beginObject()
	w.key("d")
	w.stringValue("")
	w.endObject()
	w.endObject()
	require.NoError(t, w.flush())
	assert.Equal(t, `{"a":"1","b":["x",{},[]],"c":{"d":""}}`, buffer.String())
	assert.True(t, json.Valid(buffer.Bytes()))
}

func unsetTransientFeatures(eResource ecore.EResource) {
	for it := eResource.GetAllContents(); it.HasNext(); {
		eObject := it.Next().(ecore.EObject)
		for itFeature := eObject.EClass().GetEAllStructuralFeatures().Iterator(); itFeature.HasNext(); {
			if eFeature := itFeature.Next().(ecore.EStructuralFeature); eFeature.IsTransient() && eFeature.IsChangeable() {
				eObject.EUnset(eFeature)
			}
		}
	}
}
