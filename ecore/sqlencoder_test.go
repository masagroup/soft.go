package ecore

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSqlEncoder_Complex(t *testing.T) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	// load resource
	xmlProcessor := NewXMLProcessor(XMLProcessorPackages([]EPackage{ePackage}))
	eResource := xmlProcessor.LoadWithOptions(NewURI("testdata/library.complex.xml"), nil)
	require.NotNil(t, eResource)
	require.True(t, eResource.IsLoaded())
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	require.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	w := &bytes.Buffer{}
	sqliteEncoder := NewSQLEncoder(eResource, w, nil)
	sqliteEncoder.EncodeResource()
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
}
