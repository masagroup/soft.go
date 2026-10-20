package integration

import (
	"bytes"
	"io"
	"os"

	"strings"
	"testing"

	_ "github.com/masagroup/soft.go/codec/defaults"
	"github.com/masagroup/soft.go/codec/binarycodec"
	"github.com/masagroup/soft.go/codec/xmlcodec"
	"github.com/masagroup/soft.go/ecore"
	"github.com/masagroup/soft.go/test/library"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func diagnosticError(errors ecore.EList) string {
	if errors.Empty() {
		return ""
	} else {
		return errors.Get(0).(ecore.EDiagnostic).GetMessage()
	}
}

func TestSerializationLoadSimpleDefaultXML(t *testing.T) {
	xmlProcessor := xmlcodec.NewXMLProcessor(xmlcodec.XMLProcessorPackages([]ecore.EPackage{library.GetPackage()}))
	resource := xmlProcessor.Load(ecore.CreateFileURI("testdata/library/library.simple.default.xml"))
	assert.True(t, resource.IsLoaded())
	assert.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	assert.True(t, resource.GetWarnings().Empty(), diagnosticError(resource.GetWarnings()))
}

func TestSerializationLoadSimplePrefixXML(t *testing.T) {
	xmlProcessor := xmlcodec.NewXMLProcessor(xmlcodec.XMLProcessorPackages([]ecore.EPackage{library.GetPackage()}))
	resource := xmlProcessor.Load(ecore.CreateFileURI("testdata/library/library.simple.prefix.xml"))
	assert.True(t, resource.IsLoaded())
	assert.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	assert.True(t, resource.GetWarnings().Empty(), diagnosticError(resource.GetWarnings()))
}

func TestSerializationLoadProprietaryXML(t *testing.T) {
	xmlProcessor := xmlcodec.NewXMLProcessor(xmlcodec.XMLProcessorPackages([]ecore.EPackage{library.GetPackage()}))
	resource := xmlProcessor.Load(ecore.CreateFileURI("testdata/library/library.proprietary.xml"))
	assert.True(t, resource.IsLoaded())
	assert.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	assert.True(t, resource.GetWarnings().Empty(), diagnosticError(resource.GetWarnings()))
}

func TestSerializationLoadComplexXML(t *testing.T) {
	xmlProcessor := xmlcodec.NewXMLProcessor(xmlcodec.XMLProcessorPackages([]ecore.EPackage{library.GetPackage()}))
	resource := xmlProcessor.Load(ecore.CreateFileURI("testdata/library/library.complex.xml"))
	assert.True(t, resource.IsLoaded())
	assert.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	assert.True(t, resource.GetWarnings().Empty(), diagnosticError(resource.GetWarnings()))
}

func TestSerializationSaveSimpleXml(t *testing.T) {
	// create a library model with a single employee
	root := library.GetFactory().CreateDocumentRoot()

	lib := library.GetFactory().CreateLibrary()
	lib.SetName("My Library")
	lib.SetAddress("My Library Adress")
	root.SetLibrary(lib)

	employee := library.GetFactory().CreateEmployee()
	employee.SetFirstName("First Name")
	employee.SetLastName("Last Name")
	employee.SetAddress("adress")
	lib.GetEmployees().Add(employee)

	// save library model with a resource
	xmlProcessor := xmlcodec.NewXMLProcessor(xmlcodec.XMLProcessorPackages([]ecore.EPackage{library.GetPackage()}))
	xmlProcessor.SaveObject(ecore.CreateFileURI("testdata/library/dynamic.simple.output.xml"), root)

	bytesInput, errInput := os.ReadFile("testdata/library/dynamic.simple.result.xml")
	assert.Nil(t, errInput)
	bytesOutput, errOutput := os.ReadFile("testdata/library/dynamic.simple.output.xml")
	assert.Nil(t, errOutput)
	assert.Equal(t, strings.ReplaceAll(string(bytesInput), "\r\n", "\n"), strings.ReplaceAll(string(bytesOutput), "\r\n", "\n"))

}

func TestSerializationLoadSaveSimpleXML(t *testing.T) {
	xmlProcessor := xmlcodec.NewXMLProcessor(xmlcodec.XMLProcessorPackages([]ecore.EPackage{library.GetPackage()}))
	resource := xmlProcessor.Load(ecore.CreateFileURI("testdata/library/library.simple.default.xml"))

	var strbuff strings.Builder
	resource.SaveWithWriter(&strbuff, nil)

	bytes, err := os.ReadFile("testdata/library/library.simple.default.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(strbuff.String(), "\r\n", "\n"))
}

func TestSerializationLoadSavePrefixXML(t *testing.T) {
	xmlProcessor := xmlcodec.NewXMLProcessor(xmlcodec.XMLProcessorPackages([]ecore.EPackage{library.GetPackage()}))
	resource := xmlProcessor.Load(ecore.CreateFileURI("testdata/library/library.simple.prefix.xml"))

	var strbuff strings.Builder
	resource.SaveWithWriter(&strbuff, nil)

	bytes, err := os.ReadFile("testdata/library/library.simple.prefix.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(strbuff.String(), "\r\n", "\n"))
}

func TestSerializationLoadSaveComplexXML(t *testing.T) {
	xmlProcessor := xmlcodec.NewXMLProcessor(xmlcodec.XMLProcessorPackages([]ecore.EPackage{library.GetPackage()}))
	resource := xmlProcessor.Load(ecore.CreateFileURI("testdata/library/library.complex.xml"))

	var strbuff strings.Builder
	resource.SaveWithWriter(&strbuff, nil)

	bytes, err := os.ReadFile("testdata/library/library.complex.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(strbuff.String(), "\r\n", "\n"))
}

func TestSerializationSaveComplexSqlite(t *testing.T) {
	srcURI := ecore.CreateFileURI("testdata/library/library.complex.xml")
	destURI := ecore.CreateFileURI("testdata/library/library.output.sqlite")
	xmlProcessor := xmlcodec.NewXMLProcessor(xmlcodec.XMLProcessorPackages([]ecore.EPackage{library.GetPackage()}))
	resource := xmlProcessor.Load(srcURI)
	resource.SetURI(destURI)
	resource.Save()
	assert.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	assert.True(t, resource.GetWarnings().Empty(), diagnosticError(resource.GetWarnings()))
}

func TestSerializationLoadComplexSqlite(t *testing.T) {
	srcURI := ecore.CreateFileURI("testdata/library/library.complex.sqlite")
	xmlProcessor := xmlcodec.NewXMLProcessor(xmlcodec.XMLProcessorPackages([]ecore.EPackage{library.GetPackage()}))
	resource := xmlProcessor.Load(srcURI)
	assert.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	assert.True(t, resource.GetWarnings().Empty(), diagnosticError(resource.GetWarnings()))
}

func TestDeepOperations(t *testing.T) {
	xmlProcessor := xmlcodec.NewXMLProcessor(xmlcodec.XMLProcessorPackages([]ecore.EPackage{library.GetPackage()}))
	resource := xmlProcessor.Load(ecore.CreateFileURI("testdata/library/library.complex.xml"))

	eObject := resource.GetContents().Get(0).(ecore.EObject)
	eCopyObject := ecore.Copy(eObject)
	assert.True(t, ecore.Equals(eObject, eCopyObject))
}

func BenchmarkXMLDecoderLibraryComplexBig(b *testing.B) {
	// create resource
	uri := ecore.CreateFileURI("testdata/library/library.complex.xml")
	eResource := ecore.NewEResourceImpl()
	eResource.SetURI(uri)
	eResourceSet := ecore.NewEResourceSetImpl()
	eResourceSet.GetResources().Add(eResource)
	eResourceSet.GetPackageRegistry().RegisterPackage(library.GetPackage())

	// get file content
	content, err := os.ReadFile(uri.String())
	require.Nil(b, err)
	r := bytes.NewReader(content)

	for b.Loop() {
		_, err := r.Seek(0, io.SeekStart)
		require.Nil(b, err)
		xmlDecoder := xmlcodec.NewXMLDecoder(eResource, r, nil)
		xmlDecoder.DecodeResource()
		require.True(b, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	}
}

func BenchmarkXMLEncoderLibraryComplexBig(b *testing.B) {
	// load resource
	xmlProcessor := xmlcodec.NewXMLProcessor(xmlcodec.XMLProcessorPackages([]ecore.EPackage{library.GetPackage()}))
	eResource := xmlProcessor.LoadWithOptions(ecore.CreateFileURI("testdata/library/library.complex.xml"), nil)
	require.NotNil(b, eResource)
	require.True(b, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	for b.Loop() {
		var strbuff strings.Builder
		binaryEncoder := xmlcodec.NewXMLEncoder(eResource, &strbuff, nil)
		binaryEncoder.EncodeResource()
		require.True(b, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	}
}

func BenchmarkBinaryDecoderLibraryComplexBig(b *testing.B) {
	// create resource
	uri := ecore.CreateFileURI("testdata/library/library.complex.bin")
	eResource := ecore.NewEResourceImpl()
	eResource.SetURI(uri)
	eResourceSet := ecore.NewEResourceSetImpl()
	eResourceSet.GetResources().Add(eResource)
	eResourceSet.GetPackageRegistry().RegisterPackage(library.GetPackage())

	// get file content
	content, err := os.ReadFile(uri.String())
	require.Nil(b, err)
	r := bytes.NewReader(content)

	for b.Loop() {
		_, err := r.Seek(0, io.SeekStart)
		require.Nil(b, err)
		xmlDecoder := binarycodec.NewBinaryDecoder(eResource, r, nil)
		xmlDecoder.DecodeResource()
		require.True(b, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	}
}

func BenchmarkBinaryEncoderLibraryComplexBig(b *testing.B) {
	// load resource
	xmlProcessor := xmlcodec.NewXMLProcessor(xmlcodec.XMLProcessorPackages([]ecore.EPackage{library.GetPackage()}))
	eResource := xmlProcessor.LoadWithOptions(ecore.CreateFileURI("testdata/library/library.complex.xml"), nil)
	require.NotNil(b, eResource)
	require.True(b, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	for b.Loop() {
		var strbuff strings.Builder
		binaryEncoder := binarycodec.NewBinaryEncoder(eResource, &strbuff, nil)
		binaryEncoder.EncodeResource()
		require.True(b, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	}
}
