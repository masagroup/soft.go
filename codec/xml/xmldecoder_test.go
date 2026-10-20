// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************
package xml_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/masagroup/soft.go/codec/xmi"
	"github.com/masagroup/soft.go/codec/xml"
	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadPackage(packageFileName string) ecore.EPackage {
	xmiProcessor := xmi.NewProcessor()
	eResource := xmiProcessor.Load(ecore.NewURI("testdata/" + packageFileName))
	if eResource.IsLoaded() && eResource.GetContents().Size() > 0 {
		ePackage, _ := eResource.GetContents().Get(0).(ecore.EPackage)
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
	xmlProcessor := xml.NewProcessor(xml.ProcessorPackages([]ecore.EPackage{ePackage}))
	eResource := xmlProcessor.Load(ecore.NewURI("testdata/library.noroot.xml"))
	require.NotNil(t, eResource)
	require.True(t, eResource.IsLoaded())
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	require.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	// retrive library class & library name attribute
	eLibraryClass, _ := ePackage.GetEClassifier("Library").(ecore.EClass)
	assert.NotNil(t, eLibraryClass)
	eLibraryNameAttribute, _ := eLibraryClass.GetEStructuralFeatureFromName("name").(ecore.EAttribute)
	assert.NotNil(t, eLibraryNameAttribute)

	// check library name
	eLibrary, _ := eResource.GetContents().Get(0).(ecore.EObject)
	assert.Equal(t, "My Library", eLibrary.EGet(eLibraryNameAttribute))
}

func TestXMLDecoderLibraryComplex(t *testing.T) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	// load resource
	xmlProcessor := xml.NewProcessor(xml.ProcessorPackages([]ecore.EPackage{ePackage}))
	eResource := xmlProcessor.Load(ecore.NewURI("testdata/library.complex.xml"))
	require.NotNil(t, eResource)
	require.True(t, eResource.IsLoaded())
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	require.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	// retrieve document root class , library class & library name attribute
	eDocumentRootClass, _ := ePackage.GetEClassifier("DocumentRoot").(ecore.EClass)
	assert.NotNil(t, eDocumentRootClass)
	eDocumentRootLibraryFeature, _ := eDocumentRootClass.GetEStructuralFeatureFromName("library").(ecore.EReference)
	assert.NotNil(t, eDocumentRootLibraryFeature)
	eLibraryClass, _ := ePackage.GetEClassifier("Library").(ecore.EClass)
	assert.NotNil(t, eLibraryClass)
	eLibraryNameAttribute, _ := eLibraryClass.GetEStructuralFeatureFromName("name").(ecore.EAttribute)
	assert.NotNil(t, eLibraryNameAttribute)

	// check library name
	eDocumentRoot := eResource.GetContents().Get(0).(ecore.EObject)
	assert.NotNil(t, eDocumentRoot)
	eLibrary, _ := eDocumentRoot.EGet(eDocumentRootLibraryFeature).(ecore.EObject)
	assert.NotNil(t, eLibrary)
	assert.Equal(t, "My Library", eLibrary.EGet(eLibraryNameAttribute))

	// book class and attributes
	eLibraryBooksRefeference, _ := eLibraryClass.GetEStructuralFeatureFromName("books").(ecore.EReference)
	assert.NotNil(t, eLibraryBooksRefeference)
	eBookClass, _ := ePackage.GetEClassifier("Book").(ecore.EClass)
	require.NotNil(t, eBookClass)
	eBookTitleAttribute, _ := eBookClass.GetEStructuralFeatureFromName("title").(ecore.EAttribute)
	require.NotNil(t, eBookTitleAttribute)
	eBookDateAttribute, _ := eBookClass.GetEStructuralFeatureFromName("publicationDate").(ecore.EAttribute)
	require.NotNil(t, eBookDateAttribute)

	// retrive book
	eBooks, _ := eLibrary.EGet(eLibraryBooksRefeference).(ecore.EList)
	assert.NotNil(t, eBooks)
	eBook := eBooks.Get(0).(ecore.EObject)
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
	xmlProcessor := xml.NewProcessor(xml.ProcessorPackages([]ecore.EPackage{ePackage}))
	eResource := xmlProcessor.LoadWithOptions(
		ecore.NewURI("testdata/library.complex.noroot.xml"),
		map[string]any{
			xml.OPTION_SUPPRESS_DOCUMENT_ROOT: true,
			xml.OPTION_EXTENDED_META_DATA:     xml.NewExtendedMetaData()})
	require.NotNil(t, eResource)
	require.True(t, eResource.IsLoaded())
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	require.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	// retrive library class & library name attribute
	eLibraryClass, _ := ePackage.GetEClassifier("Library").(ecore.EClass)
	assert.NotNil(t, eLibraryClass)
	eLibraryNameAttribute, _ := eLibraryClass.GetEStructuralFeatureFromName("name").(ecore.EAttribute)
	assert.NotNil(t, eLibraryNameAttribute)

	// check library name
	eLibrary, _ := eResource.GetContents().Get(0).(ecore.EObject)
	assert.Equal(t, "My Library", eLibrary.EGet(eLibraryNameAttribute))
}

func TestXMLDecoderSimpleInvalidXML(t *testing.T) {
	// load libray simple ecore	package
	ePackage := loadPackage("library.simple.ecore")
	require.NotNil(t, ePackage)

	xmlProcessor := xml.NewProcessor(xml.ProcessorPackages([]ecore.EPackage{ePackage}))
	eResource := xmlProcessor.Load(ecore.NewURI("testdata/library.simple.invalid.xml"))
	require.NotNil(t, eResource)
	require.True(t, eResource.IsLoaded())
	require.False(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
}

func TestXMLDecoderSimpleEscapeXML(t *testing.T) {
	// load libray simple ecore	package
	ePackage := loadPackage("library.simple.ecore")
	require.NotNil(t, ePackage)

	xmlProcessor := xml.NewProcessor(xml.ProcessorPackages([]ecore.EPackage{ePackage}))
	eResource := xmlProcessor.Load(ecore.NewURI("testdata/library.simple.escape.xml"))
	require.NotNil(t, eResource)
	require.True(t, eResource.IsLoaded())
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))

	// retrive library class & library name attribute
	eLibraryClass, _ := ePackage.GetEClassifier("Library").(ecore.EClass)
	assert.NotNil(t, eLibraryClass)
	eLibraryLocationAttribute, _ := eLibraryClass.GetEStructuralFeatureFromName("location").(ecore.EAttribute)
	assert.NotNil(t, eLibraryLocationAttribute)

	// check library name
	eLibrary, _ := eResource.GetContents().Get(0).(ecore.EObject)
	assert.Equal(t, "a<b", eLibrary.EGet(eLibraryLocationAttribute))
}

func TestXMLDecoderSimpleXMLWithIDs(t *testing.T) {
	idManager := ecore.NewIncrementalIDManager()

	// load libray simple ecore	package
	ePackage := loadPackage("library.simple.ecore")
	require.NotNil(t, ePackage)

	eResourceSet := ecore.NewEResourceSetImpl()
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)
	eResource := eResourceSet.CreateResource(ecore.NewURI("testdata/library.simple.ids.xml"))
	require.NotNil(t, eResource)
	eResource.SetObjectIDManager(idManager)
	eResource.LoadWithOptions(map[string]any{xml.OPTION_ID_ATTRIBUTE_NAME: "id"})
	require.True(t, eResource.IsLoaded())
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))

	// retrive library class & library name attribute
	libraryClass, _ := ePackage.GetEClassifier("Library").(ecore.EClass)
	require.NotNil(t, libraryClass)
	libraryBooksFeature := libraryClass.GetEStructuralFeatureFromName("books")
	require.NotNil(t, libraryBooksFeature)

	require.Equal(t, 1, eResource.GetContents().Size())
	eLibrary, _ := eResource.GetContents().Get(0).(ecore.EObject)
	require.NotNil(t, eLibrary)
	assert.Equal(t, int64(0), idManager.GetID(eLibrary))

	eBooks, _ := eLibrary.EGet(libraryBooksFeature).(ecore.EList)
	require.NotNil(t, eBooks)
	require.Equal(t, 4, eBooks.Size())
	assert.Equal(t, int64(1), idManager.GetID(eBooks.Get(0).(ecore.EObject)))
	assert.Equal(t, int64(2), idManager.GetID(eBooks.Get(1).(ecore.EObject)))
	assert.Equal(t, int64(3), idManager.GetID(eBooks.Get(2).(ecore.EObject)))
	assert.Equal(t, int64(4), idManager.GetID(eBooks.Get(3).(ecore.EObject)))
}

func TestXMLDecoderSimpleXMLWithEDataTypeList(t *testing.T) {
	// load libray simple ecore	package
	ePackage := loadPackage("library.datalist.ecore")
	require.NotNil(t, ePackage)

	eResourceSet := ecore.NewEResourceSetImpl()
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)
	eResource := eResourceSet.CreateResource(ecore.NewURI("testdata/library.datalist.xml"))
	require.NotNil(t, eResource)
	eResource.Load()
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	require.True(t, eResource.IsLoaded())

	// retrieve library class & library name attribute
	libraryClass, _ := ePackage.GetEClassifier("Library").(ecore.EClass)
	require.NotNil(t, libraryClass)
	libraryBooksFeature := libraryClass.GetEStructuralFeatureFromName("books")
	require.NotNil(t, libraryBooksFeature)
	bookClass, _ := ePackage.GetEClassifier("Book").(ecore.EClass)
	require.NotNil(t, bookClass)
	bookContentsFeature := bookClass.GetEStructuralFeatureFromName("contents")
	require.NotNil(t, bookContentsFeature)

	require.Equal(t, 1, eResource.GetContents().Size())
	eLibrary, _ := eResource.GetContents().Get(0).(ecore.EObject)
	require.NotNil(t, eLibrary)

	eBooks, _ := eLibrary.EGet(libraryBooksFeature).(ecore.EList)
	require.NotNil(t, eBooks)
	require.Equal(t, 4, eBooks.Size())

	eLastBook, _ := eBooks.Get(3).(ecore.EObject)
	require.NotNil(t, eLastBook)
	eContents, _ := eLastBook.EGet(bookContentsFeature).(ecore.EList)
	require.NotNil(t, eContents)
	assert.Equal(t, 3, eContents.Size())
	assert.Equal(t, "c31", eContents.Get(0))
}

func TestXMLDecoderLibraryComplexBig(t *testing.T) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	// load resource
	xmlProcessor := xml.NewProcessor(xml.ProcessorPackages([]ecore.EPackage{ePackage}))
	eResource := xmlProcessor.LoadWithOptions(ecore.NewURI("testdata/library.complex.big.xml"), nil)
	require.NotNil(t, eResource)
	require.True(t, eResource.IsLoaded())
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	require.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))
}

func TestXMLDecoderSimpleObject(t *testing.T) {
	// load package
	ePackage := loadPackage("library.simple.ecore")
	require.NotNil(t, ePackage)

	eBookClass, _ := ePackage.GetEClassifier("Book").(ecore.EClass)
	require.NotNil(t, eBookClass)
	eBookNameAttribute, _ := eBookClass.GetEStructuralFeatureFromName("name").(ecore.EAttribute)
	require.NotNil(t, eBookNameAttribute)

	eResourceSet := ecore.NewEResourceSetImpl()
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)
	eResource := eResourceSet.CreateResource(ecore.NewURI("$tmp.xml"))

	f, err := os.Open("testdata/book.simple.xml")
	require.NotNil(t, f)
	require.Nil(t, err)

	decoder := xml.NewDecoder(eResource, f, nil)
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
	xmlProcessor := xml.NewProcessor(xml.ProcessorPackages([]ecore.EPackage{ePackage}))
	eResource := xmlProcessor.LoadWithOptions(ecore.NewURI("testdata/emap.xml"), nil)
	require.NotNil(t, eResource)
	require.True(t, eResource.IsLoaded())
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	require.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	eMapTestClass, _ := ePackage.GetEClassifier("EMapTest").(ecore.EClass)
	require.NotNil(t, eMapTestClass)
	eMapTestKeyToValueReference, _ := eMapTestClass.GetEStructuralFeatureFromName("keyToValue").(ecore.EReference)
	require.NotNil(t, eMapTestKeyToValueReference)
	eMapTestKeyToIntReference, _ := eMapTestClass.GetEStructuralFeatureFromName("keyToInt").(ecore.EReference)
	require.NotNil(t, eMapTestKeyToIntReference)
	eKeyTypeClass, _ := ePackage.GetEClassifier("KeyType").(ecore.EClass)
	require.NotNil(t, eKeyTypeClass)
	eKeyTypeNameAttribute, _ := eKeyTypeClass.GetEStructuralFeatureFromName("name").(ecore.EAttribute)
	require.NotNil(t, eKeyTypeNameAttribute)
	eValueTypeClass, _ := ePackage.GetEClassifier("ValueType").(ecore.EClass)
	require.NotNil(t, eValueTypeClass)
	eValueTypeNameAttribute, _ := eValueTypeClass.GetEStructuralFeatureFromName("name").(ecore.EAttribute)
	require.NotNil(t, eValueTypeNameAttribute)
	eRefTypeClass, _ := ePackage.GetEClassifier("RefType").(ecore.EClass)
	require.NotNil(t, eRefTypeClass)
	eRefTypeNameAttribute, _ := eRefTypeClass.GetEStructuralFeatureFromName("name").(ecore.EAttribute)
	require.NotNil(t, eRefTypeNameAttribute)
	eMapTestRefTypeReference, _ := eMapTestClass.GetEStructuralFeatureFromName("refs").(ecore.EReference)
	require.NotNil(t, eMapTestRefTypeReference)
	eMapTestRefToIntsReference, _ := eMapTestClass.GetEStructuralFeatureFromName("refToInts").(ecore.EReference)
	require.NotNil(t, eMapTestRefToIntsReference)
	eRefToIntsMapEntryClass, _ := ePackage.GetEClassifier("RefToIntsMapEntry").(ecore.EClass)
	require.NotNil(t, eRefToIntsMapEntryClass)
	eRefToIntsMapEntryKeyReference, _ := eRefToIntsMapEntryClass.GetEStructuralFeatureFromName("key").(ecore.EReference)
	require.NotNil(t, eRefToIntsMapEntryKeyReference)
	eRefToIntsMapEntryValueAttribute, _ := eRefToIntsMapEntryClass.GetEStructuralFeatureFromName("value").(ecore.EAttribute)
	require.NotNil(t, eRefToIntsMapEntryValueAttribute)

	mapTest := eResource.GetContents().Get(0).(ecore.EObject)
	require.Equal(t, eMapTestClass, mapTest.EClass())

	// map key value
	keyToValueMap, _ := mapTest.EGet(eMapTestKeyToValueReference).(ecore.EMap)
	require.NotNil(t, keyToValueMap)
	assert.Equal(t, 5, keyToValueMap.Size())
	check := 0
	for k, v := range keyToValueMap.ToMap() {
		key, _ := k.(ecore.EObject)
		require.NotNil(t, key)
		assert.Equal(t, eKeyTypeClass, key.EClass())
		keyName := key.EGet(eKeyTypeNameAttribute).(string)
		var keyIndex int
		_, err := fmt.Sscanf(keyName, "key %d", &keyIndex)
		assert.NoError(t, err)

		value, _ := v.(ecore.EObject)
		require.NotNil(t, value)
		assert.Equal(t, eValueTypeClass, value.EClass())
		valueName := value.EGet(eValueTypeNameAttribute).(string)
		var valueIndex int
		_, err = fmt.Sscanf(valueName, "value %d", &valueIndex)
		assert.NoError(t, err)
		check += keyIndex + valueIndex + 2
	}
	assert.Equal(t, 30, check)

	// map key reference with a int list value
	refList, _ := mapTest.EGet(eMapTestRefTypeReference).(ecore.EList)
	require.NotNil(t, refList)
	refToIntsMap, _ := mapTest.EGet(eMapTestKeyToValueReference).(ecore.EMap)
	require.NotNil(t, refToIntsMap)
	assert.Equal(t, 5, refToIntsMap.Size())
	for i := range 0 {
		ref := refList.Get(i)
		l, _ := refToIntsMap.GetValue(ref).(ecore.EList)
		require.NotNil(t, l)
		assert.Equal(t, i, l.Size())
	}

}

func BenchmarkXMLDecoderLibraryComplexBig(b *testing.B) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(b, ePackage)

	// create resource
	uri := ecore.NewURI("testdata/library.complex.big.xml")
	eResource := ecore.NewEResourceImpl()
	eResource.SetURI(uri)
	eResourceSet := ecore.NewEResourceSetImpl()
	eResourceSet.GetResources().Add(eResource)
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)

	// get file content
	content, err := os.ReadFile(uri.String())
	require.Nil(b, err)
	r := bytes.NewReader(content)

	for b.Loop() {
		_, err = r.Seek(0, io.SeekStart)
		require.Nil(b, err)
		xmlDecoder := xml.NewDecoder(eResource, r, nil)
		xmlDecoder.DecodeResource()
		require.True(b, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	}
}
