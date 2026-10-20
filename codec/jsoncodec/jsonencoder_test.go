package jsoncodec

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/masagroup/soft.go/codec/xmlcodec"
	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONEncoder_EncodeResourceSimple(t *testing.T) {
	// load libray simple ecore	package
	ePackage := loadPackage("library.simple.ecore")
	assert.NotNil(t, ePackage)

	// load model file
	xmlProcessor := xmlcodec.NewXMLProcessor(xmlcodec.XMLProcessorPackages([]ecore.EPackage{ePackage}))
	eResource := xmlProcessor.Load(ecore.CreateFileURI("testdata/library.simple.xml"))
	require.NotNil(t, eResource)
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))

	buffer := &bytes.Buffer{}
	encoder := NewJSONEncoder(eResource, buffer, nil)
	encoder.EncodeResource()

	bytes, err := os.ReadFile("testdata/library.simple.json")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(buffer.String(), "\r\n", "\n"))

}

func TestJSONEncoder_EncodeResourceComplex(t *testing.T) {
	// load libray simple ecore	package
	ePackage := loadPackage("library.complex.ecore")
	assert.NotNil(t, ePackage)

	// load model file
	xmlProcessor := xmlcodec.NewXMLProcessor(xmlcodec.XMLProcessorPackages([]ecore.EPackage{ePackage}))
	eResource := xmlProcessor.Load(ecore.CreateFileURI("testdata/library.complex.xml"))
	require.NotNil(t, eResource)
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))

	buffer := &bytes.Buffer{}
	encoder := NewJSONEncoder(eResource, buffer, nil)
	encoder.EncodeResource()

	bytes, err := os.ReadFile("testdata/library.complex.json")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(buffer.String(), "\r\n", "\n"))
}

func TestJSONEncoder_EncodeObject(t *testing.T) {
	eResourceSet := ecore.NewEResourceSetImpl()
	// load packages & models
	eShopPackageResource, eShopPackage := loadTestPackage(t, eResourceSet, ecore.NewURI("testdata/shop.ecore"))
	require.NotNil(t, eShopPackage)
	require.NotNil(t, eShopPackageResource)
	eShopModelResource, eShopModel := loadTestModel(t, eResourceSet, ecore.NewURI("testdata/shop.xml"))
	require.NotNil(t, eShopModel)
	require.NotNil(t, eShopModelResource)

	eOrdersPackageResource, eOrdersPackage := loadTestPackage(t, eResourceSet, ecore.NewURI("testdata/orders.ecore"))
	require.NotNil(t, eOrdersPackageResource)
	require.NotNil(t, eOrdersPackage)
	eOrdersModelResource, eOrdersModel := loadTestModel(t, eResourceSet, ecore.NewURI("testdata/orders.xml"))
	require.NotNil(t, eOrdersModelResource)
	require.NotNil(t, eOrdersModel)

	codecOptions := map[string]any{JSON_OPTION_ID_ATTRIBUTE_NAME: "id"}
	var buffer bytes.Buffer
	jsonEncoder := NewJSONEncoder(eOrdersModelResource, &buffer, codecOptions)
	require.NoError(t, jsonEncoder.EncodeObject(eOrdersModel))

	//os.WriteFile("testdata/orders.json", buffer.Bytes(), 0644)

	bytes, err := os.ReadFile("testdata/orders.json")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(buffer.String(), "\r\n", "\n"))
}
