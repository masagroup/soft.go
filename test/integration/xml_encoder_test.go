// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package integration

import (
	"os"
	"strings"
	"testing"

	"github.com/masagroup/soft.go/codec/xml"
	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestXMLEncoderLibraryNoRootWithOptions(t *testing.T) {
	// load package
	ePackage := loadPackage("library.noroot.ecore")
	assert.NotNil(t, ePackage)

	// load resource
	options := map[string]any{xml.OPTION_SUPPRESS_DOCUMENT_ROOT: true, xml.OPTION_EXTENDED_META_DATA: xml.NewExtendedMetaData()}
	xmlProcessor := xml.NewProcessor(xml.ProcessorPackages([]ecore.EPackage{ePackage}))
	eResource := xmlProcessor.LoadWithOptions(ecore.NewURI("testdata/xml/library.noroot.xml"), options)
	require.NotNil(t, eResource)
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	assert.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	// save
	eResource.SetURI(ecore.NewURI("testdata/xml/library.noroot.result.xml"))
	xmlProcessor.SaveWithOptions(eResource, options)

	// result
	src, err := os.ReadFile("testdata/xml/library.noroot.xml")
	assert.Nil(t, err)

	result, err := os.ReadFile("testdata/xml/library.noroot.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(src), "\r\n", "\n"), strings.ReplaceAll(string(result), "\r\n", "\n"))
}

func TestXMLEncoderLibraryNoRootWithReaderWriter(t *testing.T) {
	// load package
	ePackage := loadPackage("library.noroot.ecore")
	assert.NotNil(t, ePackage)

	// xml processor
	xmlProcessor := xml.NewProcessor(xml.ProcessorPackages([]ecore.EPackage{ePackage}))
	options := map[string]any{xml.OPTION_SUPPRESS_DOCUMENT_ROOT: true, xml.OPTION_EXTENDED_META_DATA: xml.NewExtendedMetaData()}

	// load resource
	reader, error := os.Open("testdata/xml/library.noroot.xml")
	require.Nil(t, error)
	eResource := xmlProcessor.LoadWithReader(reader, options)
	require.NotNil(t, eResource)
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	assert.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	// save
	writer, error := os.Create("testdata/xml/library.noroot.result.xml")
	require.Nil(t, error)
	xmlProcessor.SaveWithWriter(writer, eResource, options)

	// result
	src, err := os.ReadFile("testdata/xml/library.noroot.xml")
	assert.Nil(t, err)

	result, err := os.ReadFile("testdata/xml/library.noroot.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(src), "\r\n", "\n"), strings.ReplaceAll(string(result), "\r\n", "\n"))
}

func TestXMLEncoderLibraryComplex(t *testing.T) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	assert.NotNil(t, ePackage)

	// load resource
	xmlProcessor := xml.NewProcessor(xml.ProcessorPackages([]ecore.EPackage{ePackage}))
	eResource := xmlProcessor.Load(ecore.NewURI("testdata/xml/library.complex.xml"))
	require.NotNil(t, eResource)
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	assert.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))
	// save resource
	result := xmlProcessor.SaveToString(eResource, nil)

	bytes, err := os.ReadFile("testdata/xml/library.complex.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(result, "\r\n", "\n"))
}

func TestXMLEncoderEMaps(t *testing.T) {
	// load package
	ePackage := loadPackage("emap.ecore")
	assert.NotNil(t, ePackage)

	// load resource
	xmlProcessor := xml.NewProcessor(xml.ProcessorPackages([]ecore.EPackage{ePackage}))
	eResource := xmlProcessor.Load(ecore.NewURI("testdata/xml/emap.xml"))
	require.NotNil(t, eResource)
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	assert.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))
	// save resource
	result := xmlProcessor.SaveToString(eResource, nil)

	bytes, err := os.ReadFile("testdata/xml/emap.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(result, "\r\n", "\n"))
}

func TestXMLEncoderLibraryComplexSubElement(t *testing.T) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	assert.NotNil(t, ePackage)

	// load resource
	xmlProcessor := xml.NewProcessor(xml.ProcessorPackages([]ecore.EPackage{ePackage}))
	eResource := xmlProcessor.Load(ecore.NewURI("testdata/xml/library.complex.xml"))
	require.NotNil(t, eResource)
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	assert.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	eObject := eResource.GetEObject("//@library/@employees.0")
	require.NotNil(t, eObject)
	eContainer := eObject.EContainer()
	require.NotNil(t, eContainer)

	// create a new resource
	eNewResource := eResource.GetResourceSet().CreateResource(ecore.NewURI("testdata/xml/library.complex.sub.xml"))
	// add object to new resource
	eNewResource.GetContents().Add(eObject)
	// save it
	result := xmlProcessor.SaveToString(eNewResource, nil)

	// check result
	bytes, err := os.ReadFile("testdata/xml/library.complex.sub.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(result, "\r\n", "\n"))

	// attach to original resource
	eLibraryClass, _ := ePackage.GetEClassifier("Library").(ecore.EClass)
	require.NotNil(t, eLibraryClass)
	eLibraryEmployeesFeature := eLibraryClass.GetEStructuralFeatureFromName("employees")
	require.NotNil(t, eLibraryEmployeesFeature)
	eList := eContainer.EGet(eLibraryEmployeesFeature).(ecore.EList)
	eList.Add(eObject)
	assert.Equal(t, eResource, eObject.EResource())
}

func TestXMLEncoderLibraryComplexWithOptions(t *testing.T) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	assert.NotNil(t, ePackage)

	options := map[string]any{xml.OPTION_SUPPRESS_DOCUMENT_ROOT: true, xml.OPTION_EXTENDED_META_DATA: xml.NewExtendedMetaData()}

	// load resource
	xmlProcessor := xml.NewProcessor(xml.ProcessorPackages([]ecore.EPackage{ePackage}))
	eResource := xmlProcessor.LoadWithOptions(ecore.NewURI("testdata/xml/library.complex.noroot.xml"), options)
	require.NotNil(t, eResource)
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	assert.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	// save resource
	result := xmlProcessor.SaveToString(eResource, options)

	bytes, err := os.ReadFile("testdata/xml/library.complex.noroot.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(result, "\r\n", "\n"))
}

func TestXMLEncoderSimpleEscapeXML(t *testing.T) {
	// load libray simple ecore	package
	ePackage := loadPackage("library.simple.ecore")
	assert.NotNil(t, ePackage)

	// retrive library class & library name attribute
	eLibraryClass, _ := ePackage.GetEClassifier("Library").(ecore.EClass)
	assert.NotNil(t, eLibraryClass)
	eLibraryLocationAttribute, _ := eLibraryClass.GetEStructuralFeatureFromName("location").(ecore.EAttribute)
	assert.NotNil(t, eLibraryLocationAttribute)

	eFactory := ePackage.GetEFactoryInstance()
	eLibrary := eFactory.Create(eLibraryClass)
	eLibrary.ESet(eLibraryLocationAttribute, "a<b")

	xmlProcessor := xml.NewProcessor(xml.ProcessorPackages([]ecore.EPackage{ePackage}))
	eResource := xmlProcessor.GetResourceSet().CreateResource(ecore.NewURI("testdata/xml/library.simple.escape.output.xml"))
	eResource.GetContents().Add(eLibrary)
	result := xmlProcessor.SaveToString(eResource, nil)

	bytes, err := os.ReadFile("testdata/xml/library.simple.escape.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(result, "\r\n", "\n"))
}

func TestXMLEncoderSimpleXMLWithIDs(t *testing.T) {

	// load libray simple ecore	package
	ePackage := loadPackage("library.simple.ecore")
	assert.NotNil(t, ePackage)

	eResourceSet := ecore.NewEResourceSetImpl()
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)
	eResource := eResourceSet.CreateResource(ecore.NewURI("testdata/xml/library.simple.xml"))
	require.NotNil(t, eResource)
	eResource.SetObjectIDManager(ecore.NewIncrementalIDManager())
	eResource.Load()
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))

	var strbuff strings.Builder
	eResource.SaveWithWriter(&strbuff, map[string]any{xml.OPTION_ID_ATTRIBUTE_NAME: "id"})

	bytes, err := os.ReadFile("testdata/xml/library.simple.ids.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(strbuff.String(), "\r\n", "\n"))
}

func TestXMLEncoderSimpleXMLWithEDataTypeList(t *testing.T) {

	// load libray simple ecore	package
	ePackage := loadPackage("library.datalist.ecore")
	assert.NotNil(t, ePackage)

	xmlProcessor := xml.NewProcessor(xml.ProcessorPackages([]ecore.EPackage{ePackage}))
	eResource := xmlProcessor.Load(ecore.NewURI("testdata/xml/library.datalist.xml"))
	require.NotNil(t, eResource)
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	assert.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))
	// save resource
	result := xmlProcessor.SaveToString(eResource, nil)

	bytes, err := os.ReadFile("testdata/xml/library.datalist.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(result, "\r\n", "\n"))
}

func TestXMLEncoderSimpleXMLRootObjects(t *testing.T) {
	// load libray simple ecore	package
	ePackage := loadPackage("library.simple.ecore")
	assert.NotNil(t, ePackage)

	// load model file
	xmlProcessor := xml.NewProcessor(xml.ProcessorPackages([]ecore.EPackage{ePackage}))
	eResource := xmlProcessor.Load(ecore.NewURI("testdata/xml/library.simple.xml"))
	require.NotNil(t, eResource)
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))

	// retrieve second book
	libraryClass, _ := ePackage.GetEClassifier("Library").(ecore.EClass)
	require.NotNil(t, libraryClass)
	libraryBooksFeature := libraryClass.GetEStructuralFeatureFromName("books")
	require.NotNil(t, libraryBooksFeature)

	require.Equal(t, 1, eResource.GetContents().Size())
	eLibrary, _ := eResource.GetContents().Get(0).(ecore.EObject)
	require.NotNil(t, eLibrary)

	eBooks, _ := eLibrary.EGet(libraryBooksFeature).(ecore.EList)
	require.NotNil(t, eBooks)
	require.Equal(t, 4, eBooks.Size())
	eBook := eBooks.Get(1).(ecore.EObject)

	// save it now
	var strbuff strings.Builder
	eResource.SaveWithWriter(&strbuff, map[string]any{xml.OPTION_ROOT_OBJECTS: ecore.NewImmutableEList([]any{eBook})})

	bytes, err := os.ReadFile("testdata/xml/book.simple.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(strbuff.String(), "\r\n", "\n"))
}

func TestXMLEncoderSimpleObject(t *testing.T) {
	// load libray simple ecore	package
	ePackage := loadPackage("library.simple.ecore")
	assert.NotNil(t, ePackage)

	// load model file
	xmlProcessor := xml.NewProcessor(xml.ProcessorPackages([]ecore.EPackage{ePackage}))
	eResource := xmlProcessor.Load(ecore.NewURI("testdata/xml/library.simple.xml"))
	require.NotNil(t, eResource)
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))

	// retrieve second book
	libraryClass, _ := ePackage.GetEClassifier("Library").(ecore.EClass)
	require.NotNil(t, libraryClass)
	libraryBooksFeature := libraryClass.GetEStructuralFeatureFromName("books")
	require.NotNil(t, libraryBooksFeature)

	require.Equal(t, 1, eResource.GetContents().Size())
	eLibrary, _ := eResource.GetContents().Get(0).(ecore.EObject)
	require.NotNil(t, eLibrary)

	eBooks, _ := eLibrary.EGet(libraryBooksFeature).(ecore.EList)
	require.NotNil(t, eBooks)
	require.Equal(t, 4, eBooks.Size())
	eBook := eBooks.Get(1).(ecore.EObject)

	var strbuff strings.Builder
	encoder := xml.NewEncoder(eResource, &strbuff, nil)
	err := encoder.EncodeObject(eBook)
	assert.Nil(t, err)

	bytes, err := os.ReadFile("testdata/xml/book.simple.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(strbuff.String(), "\r\n", "\n"))
}

func BenchmarkXMLEncoderDynamicLibraryComplexBig(b *testing.B) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(b, ePackage)
	xmlProcessor := xml.NewProcessor(xml.ProcessorPackages([]ecore.EPackage{ePackage}))
	eResource := xmlProcessor.Load(ecore.NewURI("testdata/xml/library.complex.big.xml"))
	require.NotNil(b, eResource)
	for b.Loop() {
		var strbuff strings.Builder
		eResource.SaveWithWriter(&strbuff, nil)
		assert.True(b, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	}
}

func BenchmarkXMLEncoderAllTypes(b *testing.B) {
	ePackage := loadPackage("alltypes.ecore")
	require.NotNil(b, ePackage)

	eResourceSet := ecore.NewEResourceSetImpl()
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)
	eResource := eResourceSet.CreateResource(ecore.NewURI("testdata/xml/alltypes.xml"))
	require.NotNil(b, eResource)

	// retrive library class & library name attribute
	objectClass, _ := ePackage.GetEClassifier("Object").(ecore.EClass)
	require.NotNil(b, objectClass)

	enumType := ePackage.GetEClassifier("EnumCategory").(ecore.EEnum)
	require.NotNil(b, enumType)

	objectF32Attribute := objectClass.GetEStructuralFeatureFromName("f32")
	require.NotNil(b, objectF32Attribute)

	objectF64Attribute := objectClass.GetEStructuralFeatureFromName("f64")
	require.NotNil(b, objectF64Attribute)

	objectStringAttribute := objectClass.GetEStructuralFeatureFromName("str")
	require.NotNil(b, objectStringAttribute)

	objectI8Attribute := objectClass.GetEStructuralFeatureFromName("i8")
	require.NotNil(b, objectI8Attribute)

	objectI16Attribute := objectClass.GetEStructuralFeatureFromName("i16")
	require.NotNil(b, objectI16Attribute)

	objectI32Attribute := objectClass.GetEStructuralFeatureFromName("i32")
	require.NotNil(b, objectI32Attribute)

	objectI64Attribute := objectClass.GetEStructuralFeatureFromName("i64")
	require.NotNil(b, objectI64Attribute)

	objectIntAttribute := objectClass.GetEStructuralFeatureFromName("i")
	require.NotNil(b, objectIntAttribute)

	objectBytesAttribute := objectClass.GetEStructuralFeatureFromName("bytes")
	require.NotNil(b, objectBytesAttribute)

	objectBoolAttribute := objectClass.GetEStructuralFeatureFromName("b")
	require.NotNil(b, objectBoolAttribute)

	objectEnumAttribute := objectClass.GetEStructuralFeatureFromName("e")
	require.NotNil(b, objectEnumAttribute)

	eFactory := ePackage.GetEFactoryInstance()
	eObject := eFactory.Create(objectClass)
	eObject.ESet(objectF64Attribute, float64(4.0))
	eObject.ESet(objectF32Attribute, float32(3.0))
	eObject.ESet(objectStringAttribute, "str")
	eObject.ESet(objectI8Attribute, byte('b'))
	eObject.ESet(objectI16Attribute, int16(2))
	eObject.ESet(objectI32Attribute, int32(1))
	eObject.ESet(objectI64Attribute, int64(0))
	eObject.ESet(objectIntAttribute, int(-1))
	eObject.ESet(objectBytesAttribute, []byte("bytes"))
	eObject.ESet(objectBoolAttribute, true)
	eObject.ESet(objectEnumAttribute, enumType.GetDefaultValue())

	eResource.GetContents().Add(eObject)
	eResource.Save()
}
