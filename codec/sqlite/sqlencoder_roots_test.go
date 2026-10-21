package sqlite

import (
	"bytes"
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/require"
)

// all the roots of a resource are encoded, and decoded in their order
func TestSQLEncoder_Roots(t *testing.T) {
	f := ecore.GetFactory()
	eResource := ecore.NewEResourceImpl()
	eResource.SetURI(ecore.NewURI("file:///roots.sqlite"))
	for _, name := range []string{"a", "b", "c"} {
		ePackage := f.CreateEPackage()
		ePackage.SetName(name)
		eResource.GetContents().Add(ePackage)
	}
	var buffer bytes.Buffer
	NewWriterEncoder(&buffer, eResource, map[string]any{OPTION_IN_MEMORY_DATABASE: true}).EncodeResource()
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))

	eDecoded := ecore.NewEResourceImpl()
	eDecoded.SetURI(ecore.NewURI("file:///roots.sqlite"))
	NewReaderDecoder(bytes.NewReader(buffer.Bytes()), eDecoded, map[string]any{OPTION_IN_MEMORY_DATABASE: true}).DecodeResource()
	require.True(t, eDecoded.GetErrors().Empty(), diagnosticError(eDecoded.GetErrors()))
	names := []string{}
	for root := range eDecoded.GetContents().All() {
		names = append(names, root.(ecore.EPackage).GetName())
	}
	require.Equal(t, []string{"a", "b", "c"}, names)
}
