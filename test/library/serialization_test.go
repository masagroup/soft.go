package library

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
)

func diagnosticError(errors ecore.EList) string {
	if errors.Empty() {
		return ""
	} else {
		return errors.Get(0).(ecore.EDiagnostic).GetMessage()
	}
}

func TestSerializationLoadSimpleDefaultXML(t *testing.T) {
	xmlProcessor := ecore.NewXMLProcessor([]ecore.EPackage{GetPackage()})
	resource := xmlProcessor.Load(ecore.CreateFileURI("testdata/library.simple.default.xml"))
	assert.True(t, resource.IsLoaded())
	assert.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	assert.True(t, resource.GetWarnings().Empty(), diagnosticError(resource.GetWarnings()))
}

func TestSerializationLoadSimplePrefixXML(t *testing.T) {
	xmlProcessor := ecore.NewXMLProcessor([]ecore.EPackage{GetPackage()})
	resource := xmlProcessor.Load(ecore.CreateFileURI("testdata/library.simple.prefix.xml"))
	assert.True(t, resource.IsLoaded())
	assert.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	assert.True(t, resource.GetWarnings().Empty(), diagnosticError(resource.GetWarnings()))
}

func TestSerializationLoadOwnerXML(t *testing.T) {
	xmlProcessor := ecore.NewXMLProcessor([]ecore.EPackage{GetPackage()})
	resource := xmlProcessor.Load(ecore.CreateFileURI("testdata/library.owner.xml"))
	assert.True(t, resource.IsLoaded())
	assert.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	assert.True(t, resource.GetWarnings().Empty(), diagnosticError(resource.GetWarnings()))
}

func TestSerializationLoadComplexXML(t *testing.T) {
	xmlProcessor := ecore.NewXMLProcessor([]ecore.EPackage{GetPackage()})
	resource := xmlProcessor.Load(ecore.CreateFileURI("testdata/library.complex.xml"))
	assert.True(t, resource.IsLoaded())
	assert.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	assert.True(t, resource.GetWarnings().Empty(), diagnosticError(resource.GetWarnings()))
}

func TestSerializationSaveSimpleXml(t *testing.T) {
	// create a library model with a single employee
	root := GetFactory().CreateDocumentRoot()

	library := GetFactory().CreateLibrary()
	library.SetName("My Library")
	library.SetAddress("My Library Adress")
	root.SetLibrary(library)

	employee := GetFactory().CreateEmployee()
	employee.SetFirstName("First Name")
	employee.SetLastName("Last Name")
	employee.SetAddress("adress")
	library.GetEmployees().Add(employee)

	// save library model with a resource
	xmlProcessor := ecore.NewXMLProcessor([]ecore.EPackage{GetPackage()})
	xmlProcessor.SaveObject(ecore.CreateFileURI("testdata/dynamic.simple.output.xml"), root)

	bytesInput, errInput := ioutil.ReadFile("testdata/dynamic.simple.result.xml")
	assert.Nil(t, errInput)
	bytesOutput, errOutput := ioutil.ReadFile("testdata/dynamic.simple.output.xml")
	assert.Nil(t, errOutput)
	assert.Equal(t, strings.ReplaceAll(string(bytesInput), "\r\n", "\n"), strings.ReplaceAll(string(bytesOutput), "\r\n", "\n"))

}

func TestSerializationLoadSaveSimpleXML(t *testing.T) {
	xmlProcessor := ecore.NewXMLProcessor([]ecore.EPackage{GetPackage()})
	resource := xmlProcessor.Load(ecore.CreateFileURI("testdata/library.simple.default.xml"))

	var strbuff strings.Builder
	resource.SaveWithWriter(&strbuff, nil)

	bytes, err := ioutil.ReadFile("testdata/library.simple.default.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(strbuff.String(), "\r\n", "\n"))
}

func TestSerializationLoadSavePrefixXML(t *testing.T) {
	xmlProcessor := ecore.NewXMLProcessor([]ecore.EPackage{GetPackage()})
	resource := xmlProcessor.Load(ecore.CreateFileURI("testdata/library.simple.prefix.xml"))

	var strbuff strings.Builder
	resource.SaveWithWriter(&strbuff, nil)

	bytes, err := ioutil.ReadFile("testdata/library.simple.prefix.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(strbuff.String(), "\r\n", "\n"))
}

func TestSerializationLoadSaveComplexXML(t *testing.T) {
	xmlProcessor := ecore.NewXMLProcessor([]ecore.EPackage{GetPackage()})
	resource := xmlProcessor.Load(ecore.CreateFileURI("testdata/library.complex.xml"))

	var strbuff strings.Builder
	resource.SaveWithWriter(&strbuff, nil)

	bytes, err := ioutil.ReadFile("testdata/library.complex.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(strbuff.String(), "\r\n", "\n"))
}

func TestDeepOperations(t *testing.T) {
	xmlProcessor := ecore.NewXMLProcessor([]ecore.EPackage{GetPackage()})
	resource := xmlProcessor.Load(ecore.CreateFileURI("testdata/library.complex.xml"))

	eObject := resource.GetContents().Get(0).(ecore.EObject)
	eCopyObject := ecore.Copy(eObject)
	assert.True(t, ecore.Equals(eObject, eCopyObject))
}
