package xml

import (
	"bytes"
	"testing"
	"unicode/utf16"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// encodingModel is a package of one class C with a string attribute "single", registered in a
// resource set.
type encodingModel struct {
	ePackage     ecore.EPackage
	eClass       ecore.EClass
	single       ecore.EAttribute
	eResourceSet ecore.EResourceSet
}

func newEncodingModel() *encodingModel {
	f := ecore.GetFactory()
	m := &encodingModel{}
	m.ePackage = f.CreateEPackage()
	m.ePackage.SetName("p")
	m.ePackage.SetNsURI("http://p")
	m.ePackage.SetNsPrefix("p")
	m.ePackage.SetEFactoryInstance(f.CreateEFactory())
	m.eClass = f.CreateEClass()
	m.eClass.SetName("C")
	m.ePackage.GetEClassifiers().Add(m.eClass)
	m.single = f.CreateEAttribute()
	m.single.SetName("single")
	m.single.SetEType(ecore.GetPackage().GetEString())
	m.eClass.GetEStructuralFeatures().Add(m.single)
	m.eResourceSet = ecore.NewEResourceSetImpl()
	m.eResourceSet.GetPackageRegistry().RegisterPackage(m.ePackage)
	return m
}

// decode reads document into a new resource, without error, and returns its decoder, the resource
// and the value of single.
func (m *encodingModel) decode(t *testing.T, document []byte) (*Decoder, ecore.EResource, any) {
	t.Helper()
	eResource := m.eResourceSet.CreateResource(ecore.NewURI("file:///encoding.xml"))
	t.Cleanup(func() { m.eResourceSet.GetResources().Remove(eResource) })
	decoder := NewDecoder(eResource, bytes.NewReader(document), nil)
	decoder.DecodeResource()
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	require.Equal(t, 1, eResource.GetContents().Size())
	return decoder, eResource, eResource.GetContents().Get(0).(ecore.EObject).EGet(m.single)
}

// decodeError reads document into a new resource and returns its first error, "" for none.
func (m *encodingModel) decodeError(document []byte) string {
	eResource := m.eResourceSet.CreateResource(ecore.NewURI("file:///encoding-error.xml"))
	defer m.eResourceSet.GetResources().Remove(eResource)
	NewDecoder(eResource, bytes.NewReader(document), nil).DecodeResource()
	if eResource.GetErrors().Empty() {
		return ""
	}
	return eResource.GetErrors().Get(0).(ecore.EDiagnostic).GetMessage()
}

// encode writes eResource in encoding with the document type documentType, and returns its bytes,
// without error.
func encode(t *testing.T, eResource ecore.EResource, encoding, documentType string) []byte {
	t.Helper()
	var buffer bytes.Buffer
	encoder := NewEncoder(eResource, &buffer, map[string]any{OPTION_ENCODING: encoding})
	encoder.SetDocumentType(documentType)
	encoder.EncodeResource()
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	return buffer.Bytes()
}

func toUTF16(s string, bigEndian bool, bom bool) []byte {
	var b []byte
	units := utf16.Encode([]rune(s))
	if bom {
		units = append([]uint16{0xFEFF}, units...)
	}
	for _, u := range units {
		if bigEndian {
			b = append(b, byte(u>>8), byte(u))
		} else {
			b = append(b, byte(u), byte(u>>8))
		}
	}
	return b
}

// A document in ISO-8859-1 is read as ISO-8859-1, not windows-1252, and written in the encoding its
// decoder gives: a character the encoding cannot represent is written as a character reference.
func TestXMLEncoding_ISO88591(t *testing.T) {
	m := newEncodingModel()
	document := []byte("<?xml version=\"1.0\" encoding=\"ISO-8859-1\"?>\n<p:C xmlns:p=\"http://p\" single=\"\xe9\x80\"/>\n")
	decoder, eResource, value := m.decode(t, document)
	assert.Equal(t, "é\u0080", value)
	assert.Equal(t, "ISO-8859-1", decoder.GetEncoding())
	assert.Equal(t, document, encode(t, eResource, decoder.GetEncoding(), ""))

	eResource.GetContents().Get(0).(ecore.EObject).ESet(m.single, "é€")
	assert.Equal(t, "<?xml version=\"1.0\" encoding=\"ISO-8859-1\"?>\n<p:C xmlns:p=\"http://p\" single=\"\xe9&#8364;\"/>\n", string(encode(t, eResource, "ISO-8859-1", "")))
}

// OPTION_ENCODING writes a document in another encoding than UTF-8, the default.
func TestXMLEncoding_OptionEncoding(t *testing.T) {
	m := newEncodingModel()
	_, eResource, _ := m.decode(t, []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<p:C xmlns:p=\"http://p\" single=\"é€\"/>\n"))
	assert.Equal(t, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<p:C xmlns:p=\"http://p\" single=\"é€\"/>\n", string(encode(t, eResource, "", "")))
	ascii := encode(t, eResource, "US-ASCII", "")
	assert.Equal(t, "<?xml version=\"1.0\" encoding=\"US-ASCII\"?>\n<p:C xmlns:p=\"http://p\" single=\"&#233;&#8364;\"/>\n", string(ascii))
	_, _, value := m.decode(t, ascii)
	assert.Equal(t, "é€", value)
}

// A document in UTF-16 is read by its byte order mark, or by its declaration without one, and
// written in UTF-16 with a byte order mark.
func TestXMLEncoding_UTF16(t *testing.T) {
	m := newEncodingModel()
	text := "<?xml version=\"1.0\" encoding=\"UTF-16\"?>\n<p:C xmlns:p=\"http://p\" single=\"é€\"/>\n"
	for name, document := range map[string][]byte{
		"little endian with bom":    toUTF16(text, false, true),
		"big endian with bom":       toUTF16(text, true, true),
		"little endian without bom": toUTF16(text, false, false),
		"big endian without bom":    toUTF16(text, true, false),
	} {
		t.Run(name, func(t *testing.T) {
			decoder, eResource, value := m.decode(t, document)
			assert.Equal(t, "é€", value)
			assert.Equal(t, "UTF-16", decoder.GetEncoding())
			assert.Equal(t, toUTF16(text, true, true), encode(t, eResource, decoder.GetEncoding(), ""))
		})
	}
}

// A document detected as UTF-16 that declares another encoding, or that declares an unknown one, is
// an error; so is an encoding the encoder cannot write, which writes nothing.
func TestXMLEncoding_Errors(t *testing.T) {
	m := newEncodingModel()
	assert.Contains(t, m.decodeError(toUTF16("<?xml version=\"1.0\" encoding=\"ISO-8859-1\"?><p:C xmlns:p=\"http://p\"/>", false, true)), `encoding "ISO-8859-1" of a document in UTF-16`)
	assert.Contains(t, m.decodeError([]byte("<?xml version=\"1.0\" encoding=\"X-NOPE\"?><p:C xmlns:p=\"http://p\"/>")), `unsupported encoding "X-NOPE"`)

	_, eResource, _ := m.decode(t, []byte("<p:C xmlns:p=\"http://p\" single=\"v\"/>"))
	var buffer bytes.Buffer
	NewEncoder(eResource, &buffer, map[string]any{OPTION_ENCODING: "X-NOPE"}).EncodeResource()
	require.Equal(t, 1, eResource.GetErrors().Size())
	assert.Contains(t, eResource.GetErrors().Get(0).(ecore.EDiagnostic).GetMessage(), `unsupported encoding "X-NOPE"`)
	assert.Empty(t, buffer.Bytes())
}

// A label IANA does not support is read with its WHATWG encoding, as before.
func TestXMLEncoding_WHATWGLabel(t *testing.T) {
	m := newEncodingModel()
	decoder, _, value := m.decode(t, []byte("<?xml version=\"1.0\" encoding=\"GB2312\"?><p:C xmlns:p=\"http://p\" single=\"\xc4\xe3\"/>"))
	assert.Equal(t, "你", value)
	assert.Equal(t, "GB2312", decoder.GetEncoding())
}

// The document type declaration of a document, <!DOCTYPE TS> of a Qt Linguist file, is given by its
// decoder, and written after the XML declaration.
func TestXMLEncoding_DocumentType(t *testing.T) {
	m := newEncodingModel()
	document := "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n<!DOCTYPE TS>\n<p:C xmlns:p=\"http://p\" single=\"v\"/>\n"
	decoder, eResource, _ := m.decode(t, []byte(document))
	assert.Equal(t, "<!DOCTYPE TS>", decoder.GetDocumentType())
	assert.Equal(t, document, string(encode(t, eResource, decoder.GetEncoding(), decoder.GetDocumentType())))
}
