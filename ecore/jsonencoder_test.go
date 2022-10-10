package ecore

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONEncoder_EncodeResourceSimple(t *testing.T) {
	// load libray simple ecore	package
	ePackage := loadPackage("library.simple.ecore")
	assert.NotNil(t, ePackage)

	// load model file
	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	eResource := xmlProcessor.Load(CreateFileURI("testdata/library.simple.xml"))
	require.NotNil(t, eResource)
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))

	f, error := os.Create("testdata/library.simple.result.json")
	require.Nil(t, error)
	defer f.Close()

	encoder := NewJSONEncoder(eResource, f, nil)
	encoder.Encode()
}

func TestJSONEncoder_EncodeResourceComplex(t *testing.T) {
	// load libray simple ecore	package
	ePackage := loadPackage("library.complex.ecore")
	assert.NotNil(t, ePackage)

	// load model file
	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	eResource := xmlProcessor.Load(CreateFileURI("testdata/library.complex.xml"))
	require.NotNil(t, eResource)
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))

	f, error := os.Create("testdata/library.complex.result.json")
	require.Nil(t, error)
	defer f.Close()

	encoder := NewJSONEncoder(eResource, f, nil)
	encoder.Encode()
}
