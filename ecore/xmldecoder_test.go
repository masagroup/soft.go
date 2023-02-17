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
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadPackage(packageFileName string) EPackage {
	xmiProcessor := NewXMIProcessor()
	eResource := xmiProcessor.Load(NewURI("testdata/" + packageFileName))
	if eResource.IsLoaded() && eResource.GetContents().Size() > 0 {
		ePackage, _ := eResource.GetContents().Get(0).(EPackage)
		return ePackage
	} else {
		return nil
	}
}

func TestXMLDecoderLibraryNoRoot(t *testing.T) {
	// load package
	ePackage := loadPackage("library.noroot.ecore")
	require.NotNil(t, ePackage)

	// load resource
	xmlProcessor := NewXMLProcessor(XMLProcessorPackages([]EPackage{ePackage}))
	eResource := xmlProcessor.Load(NewURI("testdata/library.noroot.xml"))
	require.NotNil(t, eResource)
	require.True(t, eResource.IsLoaded())
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	require.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	// retrive library class & library name attribute
	eLibraryClass, _ := ePackage.GetEClassifier("Library").(EClass)
	assert.NotNil(t, eLibraryClass)
	eLibraryNameAttribute, _ := eLibraryClass.GetEStructuralFeatureFromName("name").(EAttribute)
	assert.NotNil(t, eLibraryNameAttribute)

	// check library name
	eLibrary, _ := eResource.GetContents().Get(0).(EObject)
	assert.Equal(t, "My Library", eLibrary.EGet(eLibraryNameAttribute))
}

func TestXMLDecoderLibraryComplex(t *testing.T) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	// load resource
	xmlProcessor := NewXMLProcessor(XMLProcessorPackages([]EPackage{ePackage}))
	eResource := xmlProcessor.Load(NewURI("testdata/library.complex.xml"))
	require.NotNil(t, eResource)
	require.True(t, eResource.IsLoaded())
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	require.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	// retrieve document root class , library class & library name attribute
	eDocumentRootClass, _ := ePackage.GetEClassifier("DocumentRoot").(EClass)
	assert.NotNil(t, eDocumentRootClass)
	eDocumentRootLibraryFeature, _ := eDocumentRootClass.GetEStructuralFeatureFromName("library").(EReference)
	assert.NotNil(t, eDocumentRootLibraryFeature)
	eLibraryClass, _ := ePackage.GetEClassifier("Library").(EClass)
	assert.NotNil(t, eLibraryClass)
	eLibraryNameAttribute, _ := eLibraryClass.GetEStructuralFeatureFromName("name").(EAttribute)
	assert.NotNil(t, eLibraryNameAttribute)

	// check library name
	eDocumentRoot := eResource.GetContents().Get(0).(EObject)
	assert.NotNil(t, eDocumentRoot)
	eLibrary, _ := eDocumentRoot.EGet(eDocumentRootLibraryFeature).(EObject)
	assert.NotNil(t, eLibrary)
	assert.Equal(t, "My Library", eLibrary.EGet(eLibraryNameAttribute))

	// book class and attributes
	eLibraryBooksRefeference, _ := eLibraryClass.GetEStructuralFeatureFromName("books").(EReference)
	assert.NotNil(t, eLibraryBooksRefeference)
	eBookClass, _ := ePackage.GetEClassifier("Book").(EClass)
	require.NotNil(t, eBookClass)
	eBookTitleAttribute, _ := eBookClass.GetEStructuralFeatureFromName("title").(EAttribute)
	require.NotNil(t, eBookTitleAttribute)
	eBookDateAttribute, _ := eBookClass.GetEStructuralFeatureFromName("publicationDate").(EAttribute)
	require.NotNil(t, eBookDateAttribute)

	// retrive book
	eBooks, _ := eLibrary.EGet(eLibraryBooksRefeference).(EList)
	assert.NotNil(t, eBooks)
	eBook := eBooks.Get(0).(EObject)
	require.NotNil(t, eBook)

	// check book name
	assert.Equal(t, "Title 0", eBook.EGet(eBookTitleAttribute))

	// check book date
	date, _ := eBook.EGet(eBookDateAttribute).(*time.Time)
	require.NotNil(t, date)
	expected := time.Date(2015, time.September, 6, 4, 24, 46, 0, time.UTC)
	assert.Equal(t, expected, *date)
}

func TestXMLDecoderLibraryComplexWithOptions(t *testing.T) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	// load resource
	xmlProcessor := NewXMLProcessor(XMLProcessorPackages([]EPackage{ePackage}))
	eResource := xmlProcessor.LoadWithOptions(
		NewURI("testdata/library.complex.noroot.xml"),
		map[string]any{
			XML_OPTION_SUPPRESS_DOCUMENT_ROOT: true,
			XML_OPTION_EXTENDED_META_DATA:     NewExtendedMetaData()})
	require.NotNil(t, eResource)
	require.True(t, eResource.IsLoaded())
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	require.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	// retrive library class & library name attribute
	eLibraryClass, _ := ePackage.GetEClassifier("Library").(EClass)
	assert.NotNil(t, eLibraryClass)
	eLibraryNameAttribute, _ := eLibraryClass.GetEStructuralFeatureFromName("name").(EAttribute)
	assert.NotNil(t, eLibraryNameAttribute)

	// check library name
	eLibrary, _ := eResource.GetContents().Get(0).(EObject)
	assert.Equal(t, "My Library", eLibrary.EGet(eLibraryNameAttribute))
}

func TestXMLDecoderSimpleInvalidXML(t *testing.T) {
	// load libray simple ecore	package
	ePackage := loadPackage("library.simple.ecore")
	require.NotNil(t, ePackage)

	xmlProcessor := NewXMLProcessor(XMLProcessorPackages([]EPackage{ePackage}))
	eResource := xmlProcessor.Load(NewURI("testdata/library.simple.invalid.xml"))
	require.NotNil(t, eResource)
	require.True(t, eResource.IsLoaded())
	require.False(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
}

func TestXMLDecoderSimpleEscapeXML(t *testing.T) {
	// load libray simple ecore	package
	ePackage := loadPackage("library.simple.ecore")
	require.NotNil(t, ePackage)

	xmlProcessor := NewXMLProcessor(XMLProcessorPackages([]EPackage{ePackage}))
	eResource := xmlProcessor.Load(NewURI("testdata/library.simple.escape.xml"))
	require.NotNil(t, eResource)
	require.True(t, eResource.IsLoaded())
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))

	// retrive library class & library name attribute
	eLibraryClass, _ := ePackage.GetEClassifier("Library").(EClass)
	assert.NotNil(t, eLibraryClass)
	eLibraryLocationAttribute, _ := eLibraryClass.GetEStructuralFeatureFromName("location").(EAttribute)
	assert.NotNil(t, eLibraryLocationAttribute)

	// check library name
	eLibrary, _ := eResource.GetContents().Get(0).(EObject)
	assert.Equal(t, "a<b", eLibrary.EGet(eLibraryLocationAttribute))
}

func TestXMLDecoderSimpleXMLWithIDs(t *testing.T) {
	idManager := NewIncrementalIDManager()

	// load libray simple ecore	package
	ePackage := loadPackage("library.simple.ecore")
	require.NotNil(t, ePackage)

	eResourceSet := NewEResourceSetImpl()
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)
	eResource := eResourceSet.CreateResource(NewURI("testdata/library.simple.ids.xml"))
	require.NotNil(t, eResource)
	eResource.SetObjectIDManager(idManager)
	eResource.LoadWithOptions(map[string]any{XML_OPTION_ID_ATTRIBUTE_NAME: "id"})
	require.True(t, eResource.IsLoaded())
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))

	// retrive library class & library name attribute
	libraryClass, _ := ePackage.GetEClassifier("Library").(EClass)
	require.NotNil(t, libraryClass)
	libraryBooksFeature := libraryClass.GetEStructuralFeatureFromName("books")
	require.NotNil(t, libraryBooksFeature)

	require.Equal(t, 1, eResource.GetContents().Size())
	eLibrary, _ := eResource.GetContents().Get(0).(EObject)
	require.NotNil(t, eLibrary)
	assert.Equal(t, int64(0), idManager.GetID(eLibrary))

	eBooks, _ := eLibrary.EGet(libraryBooksFeature).(EList)
	require.NotNil(t, eBooks)
	require.Equal(t, 4, eBooks.Size())
	assert.Equal(t, int64(1), idManager.GetID(eBooks.Get(0).(EObject)))
	assert.Equal(t, int64(2), idManager.GetID(eBooks.Get(1).(EObject)))
	assert.Equal(t, int64(3), idManager.GetID(eBooks.Get(2).(EObject)))
	assert.Equal(t, int64(4), idManager.GetID(eBooks.Get(3).(EObject)))
}

func TestXMLDecoderSimpleXMLWithEDataTypeList(t *testing.T) {
	// load libray simple ecore	package
	ePackage := loadPackage("library.datalist.ecore")
	require.NotNil(t, ePackage)

	eResourceSet := NewEResourceSetImpl()
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)
	eResource := eResourceSet.CreateResource(NewURI("testdata/library.datalist.xml"))
	require.NotNil(t, eResource)
	eResource.Load()
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	require.True(t, eResource.IsLoaded())

	// retrieve library class & library name attribute
	libraryClass, _ := ePackage.GetEClassifier("Library").(EClass)
	require.NotNil(t, libraryClass)
	libraryBooksFeature := libraryClass.GetEStructuralFeatureFromName("books")
	require.NotNil(t, libraryBooksFeature)
	bookClass, _ := ePackage.GetEClassifier("Book").(EClass)
	require.NotNil(t, bookClass)
	bookContentsFeature := bookClass.GetEStructuralFeatureFromName("contents")
	require.NotNil(t, bookContentsFeature)

	require.Equal(t, 1, eResource.GetContents().Size())
	eLibrary, _ := eResource.GetContents().Get(0).(EObject)
	require.NotNil(t, eLibrary)

	eBooks, _ := eLibrary.EGet(libraryBooksFeature).(EList)
	require.NotNil(t, eBooks)
	require.Equal(t, 4, eBooks.Size())

	eLastBook, _ := eBooks.Get(3).(EObject)
	require.NotNil(t, eLastBook)
	eContents, _ := eLastBook.EGet(bookContentsFeature).(EList)
	require.NotNil(t, eContents)
	assert.Equal(t, 3, eContents.Size())
	assert.Equal(t, "c1", eContents.Get(0))
}

func TestXMLDecoderLibraryComplexBig(t *testing.T) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	// load resource
	xmlProcessor := NewXMLProcessor(XMLProcessorPackages([]EPackage{ePackage}))
	eResource := xmlProcessor.LoadWithOptions(NewURI("testdata/library.complex.big.xml"), nil)
	require.NotNil(t, eResource)
	require.True(t, eResource.IsLoaded())
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	require.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))
}

func TestXMLDecoderSimpleObject(t *testing.T) {
	// load package
	ePackage := loadPackage("library.simple.ecore")
	require.NotNil(t, ePackage)

	eBookClass, _ := ePackage.GetEClassifier("Book").(EClass)
	require.NotNil(t, eBookClass)
	eBookNameAttribute, _ := eBookClass.GetEStructuralFeatureFromName("name").(EAttribute)
	require.NotNil(t, eBookNameAttribute)

	eResourceSet := NewEResourceSetImpl()
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)
	eResource := eResourceSet.CreateResource(NewURI("$tmp.xml"))

	f, err := os.Open("testdata/book.simple.xml")
	require.NotNil(t, f)
	require.Nil(t, err)

	decoder := NewXMLDecoder(eResource, f, nil)
	eObject, err := decoder.DecodeObject()
	require.Nil(t, err)
	require.NotNil(t, eObject)
	assert.Equal(t, "Book 1", eObject.EGet(eBookNameAttribute))
}

func TestXMLDecoderMaps(t *testing.T) {
	// load package
	ePackage := loadPackage("emap.ecore")
	require.NotNil(t, ePackage)

	// load resource
	xmlProcessor := NewXMLProcessor(XMLProcessorPackages([]EPackage{ePackage}))
	eResource := xmlProcessor.LoadWithOptions(NewURI("testdata/emap.xml"), nil)
	require.NotNil(t, eResource)
	require.True(t, eResource.IsLoaded())
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	require.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	eMapTestClass, _ := ePackage.GetEClassifier("EMapTest").(EClass)
	require.NotNil(t, eMapTestClass)
	eMapTestKeyToValueReference, _ := eMapTestClass.GetEStructuralFeatureFromName("keyToValue").(EReference)
	require.NotNil(t, eMapTestKeyToValueReference)
	eMapTestKeyToIntReference, _ := eMapTestClass.GetEStructuralFeatureFromName("keyToInt").(EReference)
	require.NotNil(t, eMapTestKeyToIntReference)
	eKeyTypeClass, _ := ePackage.GetEClassifier("KeyType").(EClass)
	require.NotNil(t, eKeyTypeClass)
	eKeyTypeNameAttribute, _ := eKeyTypeClass.GetEStructuralFeatureFromName("name").(EAttribute)
	require.NotNil(t, eKeyTypeNameAttribute)
	eValueTypeClass, _ := ePackage.GetEClassifier("ValueType").(EClass)
	require.NotNil(t, eValueTypeClass)
	eValueTypeNameAttribute, _ := eValueTypeClass.GetEStructuralFeatureFromName("name").(EAttribute)
	require.NotNil(t, eValueTypeNameAttribute)
	eRefTypeClass, _ := ePackage.GetEClassifier("RefType").(EClass)
	require.NotNil(t, eRefTypeClass)
	eRefTypeNameAttribute, _ := eRefTypeClass.GetEStructuralFeatureFromName("name").(EAttribute)
	require.NotNil(t, eRefTypeNameAttribute)
	eMapTestRefTypeReference, _ := eMapTestClass.GetEStructuralFeatureFromName("refs").(EReference)
	require.NotNil(t, eMapTestRefTypeReference)
	eMapTestRefToIntsReference, _ := eMapTestClass.GetEStructuralFeatureFromName("refToInts").(EReference)
	require.NotNil(t, eMapTestRefToIntsReference)
	eRefToIntsMapEntryClass, _ := ePackage.GetEClassifier("RefToIntsMapEntry").(EClass)
	require.NotNil(t, eRefToIntsMapEntryClass)
	eRefToIntsMapEntryKeyReference, _ := eRefToIntsMapEntryClass.GetEStructuralFeatureFromName("key").(EReference)
	require.NotNil(t, eRefToIntsMapEntryKeyReference)
	eRefToIntsMapEntryValueAttribute, _ := eRefToIntsMapEntryClass.GetEStructuralFeatureFromName("value").(EAttribute)
	require.NotNil(t, eRefToIntsMapEntryValueAttribute)

	mapTest := eResource.GetContents().Get(0).(EObject)
	require.Equal(t, eMapTestClass, mapTest.EClass())

	// map key value
	keyToValueMap, _ := mapTest.EGet(eMapTestKeyToValueReference).(EMap)
	require.NotNil(t, keyToValueMap)
	assert.Equal(t, 5, keyToValueMap.Size())
	check := 0
	for k, v := range keyToValueMap.ToMap() {
		key, _ := k.(EObject)
		require.NotNil(t, key)
		assert.Equal(t, eKeyTypeClass, key.EClass())
		keyName := key.EGet(eKeyTypeNameAttribute).(string)
		var keyIndex int
		fmt.Sscanf(keyName, "key %d", &keyIndex)

		value, _ := v.(EObject)
		require.NotNil(t, value)
		assert.Equal(t, eValueTypeClass, value.EClass())
		valueName := value.EGet(eValueTypeNameAttribute).(string)
		var valueIndex int
		fmt.Sscanf(valueName, "value %d", &valueIndex)
		check += keyIndex + valueIndex + 2
	}
	assert.Equal(t, 30, check)

	// map key reference with a int list value
	refList, _ := mapTest.EGet(eMapTestRefTypeReference).(EList)
	require.NotNil(t, refList)
	refToIntsMap, _ := mapTest.EGet(eMapTestKeyToValueReference).(EMap)
	require.NotNil(t, refToIntsMap)
	assert.Equal(t, 5, refToIntsMap.Size())
	for i := 0; i < 0; i++ {
		ref := refList.Get(i)
		l, _ := refToIntsMap.GetValue(ref).(EList)
		require.NotNil(t, l)
		assert.Equal(t, i, l.Size())
	}

}

func BenchmarkXMLDecoderLibraryComplexBig(b *testing.B) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(b, ePackage)

	// create resource
	uri := NewURI("testdata/library.complex.big.xml")
	eResource := NewEResourceImpl()
	eResource.SetURI(uri)
	eResourceSet := NewEResourceSetImpl()
	eResourceSet.GetResources().Add(eResource)
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)

	// get file content
	content, err := os.ReadFile(uri.String())
	require.Nil(b, err)
	r := bytes.NewReader(content)

	for i := 0; i < b.N; i++ {
		_, err = r.Seek(0, io.SeekStart)
		require.Nil(b, err)
		xmlDecoder := NewXMLDecoder(eResource, r, nil)
		xmlDecoder.Decode()
		require.True(b, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	}
}
