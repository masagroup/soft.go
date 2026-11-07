package xml

import (
	"strings"
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// An xsi:nil element of an object feature, as the encoder writes an unsettable feature set to nil,
// creates no object, and the elements after it are read in the object holding it: it took that
// object off the stack of the decoder, so the next element was read in its parent, and an element
// of a concrete type created an object and read "nil" as one of its features.
func TestXMLDecoder_NilElement(t *testing.T) {
	m := newLibraryModel()
	document := `<?xml version="1.0" encoding="UTF-8"?>
<lib:Library xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:lib="http://lib" name="l" favorites="//@books.0">
  <books name="b"/>
  <optional xsi:nil="true"/>
  <maybe xsi:nil="true"/>
  <books name="c"/>
</lib:Library>
`
	decoded := m.newResourceSet().CreateResource(ecore.NewURI("file:///decoded.xml"))
	NewDecoder(decoded, strings.NewReader(document), nil).DecodeResource()
	require.True(t, decoded.GetErrors().Empty(), "errors: %v", decoded.GetErrors().ToArray())
	require.Equal(t, 1, decoded.GetContents().Size())
	library := decoded.GetContents().Get(0).(ecore.EObject)
	assert.Nil(t, library.EGet(m.optional), "no object created for an xsi:nil containment")
	assert.Nil(t, library.EGet(m.maybe), "no object created for an xsi:nil reference")
	books := library.EGet(m.books).(ecore.EList)
	require.Equal(t, 2, books.Size(), "the element after the xsi:nil ones is read in the library")
	assert.Equal(t, "c", books.Get(1).(ecore.EObject).EGet(m.name))
	favorites := library.EGet(m.favorites).(ecore.EList)
	require.Equal(t, 1, favorites.Size())
	assert.Same(t, books.Get(0), favorites.Get(0))
}
