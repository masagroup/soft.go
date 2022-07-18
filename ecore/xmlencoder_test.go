// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestXMLEncoderLibraryNoRootWithOptions(t *testing.T) {
	// load package
	ePackage := loadPackage("library.noroot.ecore")
	assert.NotNil(t, ePackage)

	// load resource
	options := map[string]interface{}{XML_OPTION_SUPPRESS_DOCUMENT_ROOT: true, XML_OPTION_EXTENDED_META_DATA: NewExtendedMetaData()}
	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	eResource := xmlProcessor.LoadWithOptions(NewURI("testdata/library.noroot.xml"), options)
	require.NotNil(t, eResource)
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	assert.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	// save
	eResource.SetURI(NewURI("testdata/library.noroot.result.xml"))
	xmlProcessor.SaveWithOptions(eResource, options)

	// result
	src, err := ioutil.ReadFile("testdata/library.noroot.xml")
	assert.Nil(t, err)

	result, err := ioutil.ReadFile("testdata/library.noroot.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(src), "\r\n", "\n"), strings.ReplaceAll(string(result), "\r\n", "\n"))
}

func TestXMLEncoderLibraryNoRootWithReaderWriter(t *testing.T) {
	// load package
	ePackage := loadPackage("library.noroot.ecore")
	assert.NotNil(t, ePackage)

	// xml processor
	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	options := map[string]interface{}{XML_OPTION_SUPPRESS_DOCUMENT_ROOT: true, XML_OPTION_EXTENDED_META_DATA: NewExtendedMetaData()}

	// load resource
	reader, error := os.Open("testdata/library.noroot.xml")
	require.Nil(t, error)
	eResource := xmlProcessor.LoadWithReader(reader, options)
	require.NotNil(t, eResource)
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	assert.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	// save
	writer, error := os.Create("testdata/library.noroot.result.xml")
	require.Nil(t, error)
	xmlProcessor.SaveWithWriter(writer, eResource, options)

	// result
	src, err := ioutil.ReadFile("testdata/library.noroot.xml")
	assert.Nil(t, err)

	result, err := ioutil.ReadFile("testdata/library.noroot.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(src), "\r\n", "\n"), strings.ReplaceAll(string(result), "\r\n", "\n"))
}

func TestXMLEncoderLibraryComplex(t *testing.T) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	assert.NotNil(t, ePackage)

	// load resource
	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	eResource := xmlProcessor.Load(NewURI("testdata/library.complex.xml"))
	require.NotNil(t, eResource)
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	assert.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))
	// save resource
	result := xmlProcessor.SaveToString(eResource, nil)

	bytes, err := ioutil.ReadFile("testdata/library.complex.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(result, "\r\n", "\n"))
}

func TestXMLEncoderEMaps(t *testing.T) {
	// load package
	ePackage := loadPackage("emap.ecore")
	assert.NotNil(t, ePackage)

	// load resource
	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	eResource := xmlProcessor.Load(NewURI("testdata/emap.xml"))
	require.NotNil(t, eResource)
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	assert.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))
	// save resource
	result := xmlProcessor.SaveToString(eResource, nil)

	bytes, err := ioutil.ReadFile("testdata/emap.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(result, "\r\n", "\n"))
}

func TestXMLEncoderLibraryComplexSubElement(t *testing.T) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	assert.NotNil(t, ePackage)

	// load resource
	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	eResource := xmlProcessor.Load(NewURI("testdata/library.complex.xml"))
	require.NotNil(t, eResource)
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	assert.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	eObject := eResource.GetEObject("//@library/@employees.0")
	require.NotNil(t, eObject)
	eContainer := eObject.EContainer()
	require.NotNil(t, eContainer)

	// create a new resource
	eNewResource := eResource.GetResourceSet().CreateResource(NewURI("testdata/library.complex.sub.xml"))
	// add object to new resource
	eNewResource.GetContents().Add(eObject)
	// save it
	result := xmlProcessor.SaveToString(eNewResource, nil)

	// check result
	bytes, err := ioutil.ReadFile("testdata/library.complex.sub.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(result, "\r\n", "\n"))

	// attach to original resource
	eLibraryClass, _ := ePackage.GetEClassifier("Library").(EClass)
	require.NotNil(t, eLibraryClass)
	eLibraryEmployeesFeature := eLibraryClass.GetEStructuralFeatureFromName("employees")
	require.NotNil(t, eLibraryEmployeesFeature)
	eList := eContainer.EGet(eLibraryEmployeesFeature).(EList)
	eList.Add(eObject)
	assert.Equal(t, eResource, eObject.EResource())
}

func TestXMLEncoderLibraryComplexWithOptions(t *testing.T) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	assert.NotNil(t, ePackage)

	options := map[string]interface{}{XML_OPTION_SUPPRESS_DOCUMENT_ROOT: true, XML_OPTION_EXTENDED_META_DATA: NewExtendedMetaData()}

	// load resource
	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	eResource := xmlProcessor.LoadWithOptions(NewURI("testdata/library.complex.noroot.xml"), options)
	require.NotNil(t, eResource)
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	assert.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	// save resource
	result := xmlProcessor.SaveToString(eResource, options)

	bytes, err := ioutil.ReadFile("testdata/library.complex.noroot.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(result, "\r\n", "\n"))
}

func TestXMLEncoderSimpleEscapeXML(t *testing.T) {
	// load libray simple ecore	package
	ePackage := loadPackage("library.simple.ecore")
	assert.NotNil(t, ePackage)

	// retrive library class & library name attribute
	eLibraryClass, _ := ePackage.GetEClassifier("Library").(EClass)
	assert.NotNil(t, eLibraryClass)
	eLibraryLocationAttribute, _ := eLibraryClass.GetEStructuralFeatureFromName("location").(EAttribute)
	assert.NotNil(t, eLibraryLocationAttribute)

	eFactory := ePackage.GetEFactoryInstance()
	eLibrary := eFactory.Create(eLibraryClass)
	eLibrary.ESet(eLibraryLocationAttribute, "a<b")

	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	eResource := xmlProcessor.GetResourceSet().CreateResource(NewURI("testdata/library.simple.escape.output.xml"))
	eResource.GetContents().Add(eLibrary)
	result := xmlProcessor.SaveToString(eResource, nil)

	bytes, err := ioutil.ReadFile("testdata/library.simple.escape.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(result, "\r\n", "\n"))
}

func TestXMLEncoderSimpleXMLWithIDs(t *testing.T) {

	// load libray simple ecore	package
	ePackage := loadPackage("library.simple.ecore")
	assert.NotNil(t, ePackage)

	eResourceSet := NewEResourceSetImpl()
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)
	eResource := eResourceSet.CreateResource(NewURI("testdata/library.simple.xml"))
	require.NotNil(t, eResource)
	eResource.SetObjectIDManager(NewIncrementalIDManager())
	eResource.Load()
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))

	var strbuff strings.Builder
	eResource.SaveWithWriter(&strbuff, map[string]interface{}{XML_OPTION_ID_ATTRIBUTE_NAME: "id"})

	bytes, err := ioutil.ReadFile("testdata/library.simple.ids.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(strbuff.String(), "\r\n", "\n"))
}

func TestXMLEncoderSimpleXMLWithEDataTypeList(t *testing.T) {

	// load libray simple ecore	package
	ePackage := loadPackage("library.datalist.ecore")
	assert.NotNil(t, ePackage)

	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	eResource := xmlProcessor.Load(NewURI("testdata/library.datalist.xml"))
	require.NotNil(t, eResource)
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	assert.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))
	// save resource
	result := xmlProcessor.SaveToString(eResource, nil)

	bytes, err := ioutil.ReadFile("testdata/library.datalist.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(result, "\r\n", "\n"))
}

func TestXMLEncoderSimpleXMLRootObjects(t *testing.T) {
	// load libray simple ecore	package
	ePackage := loadPackage("library.simple.ecore")
	assert.NotNil(t, ePackage)

	// load model file
	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	eResource := xmlProcessor.Load(NewURI("testdata/library.simple.xml"))
	require.NotNil(t, eResource)
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))

	// retrieve second book
	libraryClass, _ := ePackage.GetEClassifier("Library").(EClass)
	require.NotNil(t, libraryClass)
	libraryBooksFeature := libraryClass.GetEStructuralFeatureFromName("books")
	require.NotNil(t, libraryBooksFeature)

	require.Equal(t, 1, eResource.GetContents().Size())
	eLibrary, _ := eResource.GetContents().Get(0).(EObject)
	require.NotNil(t, eLibrary)

	eBooks, _ := eLibrary.EGet(libraryBooksFeature).(EList)
	require.NotNil(t, eBooks)
	require.Equal(t, 4, eBooks.Size())
	eBook := eBooks.Get(1).(EObject)

	// save it now
	var strbuff strings.Builder
	eResource.SaveWithWriter(&strbuff, map[string]interface{}{XML_OPTION_ROOT_OBJECTS: NewImmutableEList([]interface{}{eBook})})

	bytes, err := ioutil.ReadFile("testdata/book.simple.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(strbuff.String(), "\r\n", "\n"))
}

func TestXMLEncoderSimpleObject(t *testing.T) {
	// load libray simple ecore	package
	ePackage := loadPackage("library.simple.ecore")
	assert.NotNil(t, ePackage)

	// load model file
	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	eResource := xmlProcessor.Load(NewURI("testdata/library.simple.xml"))
	require.NotNil(t, eResource)
	assert.True(t, eResource.IsLoaded())
	assert.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))

	// retrieve second book
	libraryClass, _ := ePackage.GetEClassifier("Library").(EClass)
	require.NotNil(t, libraryClass)
	libraryBooksFeature := libraryClass.GetEStructuralFeatureFromName("books")
	require.NotNil(t, libraryBooksFeature)

	require.Equal(t, 1, eResource.GetContents().Size())
	eLibrary, _ := eResource.GetContents().Get(0).(EObject)
	require.NotNil(t, eLibrary)

	eBooks, _ := eLibrary.EGet(libraryBooksFeature).(EList)
	require.NotNil(t, eBooks)
	require.Equal(t, 4, eBooks.Size())
	eBook := eBooks.Get(1).(EObject)

	var strbuff strings.Builder
	encoder := NewXMLEncoder(eResource, &strbuff, nil)
	encoder.EncodeObject(eBook)

	bytes, err := ioutil.ReadFile("testdata/book.simple.xml")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(strbuff.String(), "\r\n", "\n"))
}

func BenchmarkXMLEncoderLibraryComplexBig(b *testing.B) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(b, ePackage)
	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	eResource := xmlProcessor.Load(NewURI("testdata/library.complex.big.xml"))
	require.NotNil(b, eResource)
	for i := 0; i < b.N; i++ {
		var strbuff strings.Builder
		eResource.SaveWithWriter(&strbuff, nil)
		assert.True(b, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	}
}
