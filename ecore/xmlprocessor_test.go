package ecore

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewXmlProcessor(t *testing.T) {
	p := NewXMLProcessor(nil)
	require.NotNil(t, p)
	assert.NotNil(t, p.GetResourceSet())
}

func TestNewSharedXmlProcessor(t *testing.T) {
	mockResourceSet := &MockEResourceSet{}
	p := NewSharedXMLProcessor(mockResourceSet)
	require.NotNil(t, p)
	assert.Equal(t, mockResourceSet, p.GetResourceSet())
}

func xmlProcessorLoad(t *testing.T, xmlProcessor *XMLProcessor, path string) EObject {
	resource := xmlProcessor.Load(&URI{Path: path})
	require.NotNil(t, resource)
	require.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	eObject, _ := resource.GetContents().Get(0).(EObject)
	require.NotNil(t, eObject)
	return eObject
}

func TestSaveObject(t *testing.T) {

	testsCases := []struct {
		name  string
		meta  string
		model string
	}{
		{"shop", "shop.ecore", "shop.xml"},
		{"orders", "orders.ecore", "orders.xml"},
	}

	resourceSet := NewEResourceSetImpl()
	packageRegistry := resourceSet.GetPackageRegistry()
	xmlProcessor := NewSharedXMLProcessor(resourceSet)
	for _, testCase := range testsCases {
		t.Run(testCase.name, func(t *testing.T) {
			ePackage, _ := xmlProcessorLoad(t, xmlProcessor, "testdata/"+testCase.meta).(EPackage)
			require.NotNil(t, ePackage)
			packageRegistry.RegisterPackage(ePackage)
			eModel, _ := xmlProcessorLoad(t, xmlProcessor, "testdata/"+testCase.model).(EObject)
			require.NotNil(t, eModel)

			resultName := "testdata/" + testCase.name + ".result.xml"
			xmlProcessor.SaveObject(&URI{Path: resultName}, eModel)

			// src
			src, err := ioutil.ReadFile("testdata/" + testCase.model)
			assert.Nil(t, err)

			// result
			result, err := ioutil.ReadFile(resultName)
			assert.Nil(t, err)
			assert.Equal(t, strings.ReplaceAll(string(src), "\r\n", "\n"), strings.ReplaceAll(string(result), "\r\n", "\n"))
		})
	}

}
