package xmi

import (
	"bytes"
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The Ecore package round-trips in XMI: some of its operations and parameters have no type, an
// unsettable eType set to nil, written as an xsi:nil element.
func TestXMICodec_EcorePackageRoundTrip(t *testing.T) {
	source := ecore.GetPackage()
	var buffer bytes.Buffer
	NewEncoder(source.EResource(), &buffer, nil).EncodeResource()
	require.True(t, source.EResource().GetErrors().Empty())
	require.Contains(t, buffer.String(), `<eType xsi:nil="true"/>`)

	decoded := ecore.NewEResourceImpl()
	decoded.SetURI(ecore.NewURI(source.GetNsURI()))
	ecore.NewEResourceSetImpl().GetResources().Add(decoded)
	NewDecoder(decoded, &buffer, nil).DecodeResource()
	require.True(t, decoded.GetErrors().Empty(), "errors: %v", decoded.GetErrors().ToArray())
	require.Equal(t, 1, decoded.GetContents().Size())
	pkg := decoded.GetContents().Get(0).(ecore.EPackage)
	eObject := pkg.GetEClassifier("EObject").(ecore.EClass)
	assert.Equal(t, ecore.GetPackage().GetEObject().GetEOperations().Size(), eObject.GetEOperations().Size())
	for i, operation := range ecore.GetPackage().GetEObject().GetEOperations().ToArray() {
		want := operation.(ecore.EOperation)
		got := eObject.GetEOperations().Get(i).(ecore.EOperation)
		assert.Equal(t, want.GetName(), got.GetName())
		assert.Equal(t, want.GetEParameters().Size(), got.GetEParameters().Size(), want.GetName())
	}
}
