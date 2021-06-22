package ecore

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBinaryEncoder_ComplexBig(t *testing.T) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	// load resource
	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	eResource := xmlProcessor.LoadWithOptions(&URI{Path: "testdata/library.complex.big.xml"}, nil)
	require.NotNil(t, eResource)
	require.True(t, eResource.IsLoaded())
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	require.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	// file
	f, err := os.Create("testdata/library.complex.bin")
	require.Nil(t, err)

	binaryEncoder := NewBinaryEncoder(eResource, f, nil)
	binaryEncoder.Encode()
}
