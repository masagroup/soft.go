package ecore

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBinaryEncoder_Complex(t *testing.T) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	// load resource
	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	eResource := xmlProcessor.LoadWithOptions(&URI{Path: "testdata/library.complex.xml"}, nil)
	require.NotNil(t, eResource)
	require.True(t, eResource.IsLoaded())
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	require.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	// w, err := os.Create("testdata/library.complex.bin")
	w := &bytes.Buffer{}
	binaryEncoder := NewBinaryEncoder(eResource, w, nil)
	binaryEncoder.Encode()
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))

	bytes, err := ioutil.ReadFile("testdata/library.complex.bin")
	assert.Nil(t, err)
	assert.Equal(t, bytes, w.Bytes())
}

func TestBinaryEncoder_ComplexWithID(t *testing.T) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	// load resource
	uri := &URI{Path: "testdata/library.complex.id.xml"}
	eResource := NewEResourceImpl()
	eResource.SetURI(uri)
	eResource.SetObjectIDManager(NewUniqueIDManager(20))
	eResourceSet := NewEResourceSetImpl()
	eResourceSet.GetResources().Add(eResource)
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)
	eResource.LoadWithOptions(map[string]interface{}{XML_OPTION_ID_ATTRIBUTE_NAME: "id"})
	require.NotNil(t, eResource)
	require.True(t, eResource.IsLoaded())
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	require.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	// set DocumentRoot uuid, because it is not defined in xml and is always regenerated
	eDocumentRoot, _ := eResource.GetContents().Get(0).(EObject)
	require.NotNil(t, eDocumentRoot)
	require.Nil(t, eResource.GetObjectIDManager().SetID(eDocumentRoot, "h0Rz1FjVeBXUgaW3OzT2frUce90="))

	w := &bytes.Buffer{}
	binaryEncoder := NewBinaryEncoder(eResource, w, map[string]interface{}{BINARY_OPTION_ID_ATTRIBUTE: true})
	binaryEncoder.Encode()
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))

	bytes, err := ioutil.ReadFile("testdata/library.complex.id.bin")
	assert.Nil(t, err)
	assert.Equal(t, bytes, w.Bytes())
}

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

	// w, err := os.Create("testdata/library.complex.big.bin")
	w := &bytes.Buffer{}
	binaryEncoder := NewBinaryEncoder(eResource, w, map[string]interface{}{})
	binaryEncoder.Encode()
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))

	bytes, err := ioutil.ReadFile("testdata/library.complex.big.bin")
	assert.Nil(t, err)
	assert.Equal(t, bytes, w.Bytes())
}

func BenchmarkBinaryEncoderLibraryComplexBig(b *testing.B) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(b, ePackage)

	// load resource
	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	eResource := xmlProcessor.LoadWithOptions(&URI{Path: "testdata/library.complex.big.xml"}, nil)
	require.NotNil(b, eResource)
	require.True(b, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	for i := 0; i < b.N; i++ {
		var buffer bytes.Buffer
		binaryEncoder := NewBinaryEncoder(eResource, &buffer, nil)
		binaryEncoder.Encode()
		require.True(b, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	}
}
