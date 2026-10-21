package xml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// packages of the same prefix: the prefixes of the next ones are numbered
func TestXMLEncoder_PrefixCollisions(t *testing.T) {
	f := ecore.GetFactory()
	m := newLibraryModel()
	rs := m.newResourceSet()
	library := m.create(m.library, "l")
	for i := 1; i <= 3; i++ {
		ePackage := f.CreateEPackage()
		ePackage.SetName("lib")
		ePackage.SetNsURI(fmt.Sprintf("http://lib%d", i))
		ePackage.SetNsPrefix("lib")
		ePackage.SetEFactoryInstance(f.CreateEFactory())
		eClass := f.CreateEClass()
		eClass.SetName("Book")
		eClass.GetESuperTypes().Add(m.book)
		ePackage.GetEClassifiers().Add(eClass)
		rs.GetPackageRegistry().RegisterPackage(ePackage)
		library.EGet(m.books).(ecore.EList).Add(ePackage.GetEFactoryInstance().Create(eClass))
	}
	eResource := rs.CreateResource(ecore.NewURI("library.xml"))
	eResource.GetContents().Add(library)

	var document string
	done := make(chan struct{})
	go func() {
		defer close(done)
		document = encodeResource(t, eResource, nil)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		require.FailNow(t, "the encoder does not end")
	}
	assert.Contains(t, document, `xmlns:lib_1="http://lib1"`)
	assert.Contains(t, document, `xmlns:lib_2="http://lib2"`)
	assert.Contains(t, document, `xmlns:lib_3="http://lib3"`)
	assert.Contains(t, document, `<books xsi:type="lib_3:Book"/>`)

	eDecoded := decodeResource(t, rs.CreateResource(ecore.NewURI("decoded.xml")), document, nil)
	assert.True(t, ecore.Equals(library, eDecoded))
}

// the hooks of an encoder embedding Encoder are called instead of its own
type hookedEncoder struct {
	*Encoder
}

func (e *hookedEncoder) SaveNamespaces() {
	e.Encoder.SaveNamespaces()
	e.AddAttribute("version", "2")
}

func TestXMLEncoder_Interfaces(t *testing.T) {
	m := newLibraryModel()
	eResource := m.newResourceSet().CreateResource(ecore.NewURI("library.xml"))
	eResource.GetContents().Add(m.create(m.library, "l"))
	var buffer bytes.Buffer
	encoder := &hookedEncoder{Encoder: NewEncoder(eResource, &buffer, nil)}
	encoder.SetInterfaces(encoder)
	encoder.SetExtendedMetaData(nil)
	encoder.EncodeResource()
	require.True(t, eResource.GetErrors().Empty())
	assert.Contains(t, buffer.String(), `<lib:Library xmlns:lib="http://lib" version="2" name="l"/>`)
}

// the decoder without extended metadata, an attribute declared as not a feature
func TestXMLDecoder_NotFeatures(t *testing.T) {
	m := newLibraryModel()
	eResource := m.newResourceSet().CreateResource(ecore.NewURI("library.xml"))
	decoder := NewDecoder(eResource, strings.NewReader(`<lib:Library xmlns:lib="http://lib" version="2" name="l"/>`), nil)
	decoder.SetExtendedMetaData(nil)
	decoder.AddNotFeatures(xml.Name{Local: "version"})
	decoder.DecodeResource()
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	assert.Equal(t, "l", eResource.GetContents().Get(0).(ecore.EObject).EGet(m.name))
}
