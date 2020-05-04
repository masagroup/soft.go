package library

import (
	"io/ioutil"
	"net/url"
	"strings"
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
)

func TestSerializationLoadXML(t *testing.T) {
	ecore.GetPackageRegistry().RegisterPackage(GetPackage())

	fileURI := &url.URL{Path: "testdata/library.xml"}
	resourceFactory := ecore.GetResourceFactoryRegistry().GetFactory(fileURI)
	resource := resourceFactory.CreateResource(fileURI)
	resource.Load()
	assert.True(t, resource.IsLoaded())
	assert.True(t, resource.GetErrors().Empty())
	assert.True(t, resource.GetWarnings().Empty())
}

func TestSerializationSaveXmlFile(t *testing.T) {
	// create a library model with a single employee
	library := GetFactory().CreateLibrary()
	employee := GetFactory().CreateEmployee()
	employee.SetFirstName("First Name")
	employee.SetLastName("Last Name")
	employee.SetAddress("adress")
	library.GetEmployees().Add(employee)

	// save library model with a resource
	ecore.GetPackageRegistry().RegisterPackage(GetPackage())
	fileURI := &url.URL{Path: "testdata/simple.output.xml"}
	resourceFactory := ecore.GetResourceFactoryRegistry().GetFactory(fileURI)
	resource := resourceFactory.CreateResource(fileURI)
	resource.GetContents().Add(library)
	resource.Save()

	bytesInput, errInput := ioutil.ReadFile("testdata/simple.input.xml")
	assert.Nil(t, errInput)
	bytesOutput, errOutput := ioutil.ReadFile("testdata/simple.output.xml")
	assert.Nil(t, errOutput)
	assert.Equal(t, strings.ReplaceAll(string(bytesInput), "\r\n", "\n"), strings.ReplaceAll(string(bytesOutput), "\r\n", "\n"))

}
