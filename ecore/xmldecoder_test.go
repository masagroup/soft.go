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
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadPackage(packageFileName string) EPackage {
	xmiProcessor := NewXMIProcessor()
	eResource := xmiProcessor.Load(&URI{Path: "testdata/" + packageFileName})
	if eResource.IsLoaded() && eResource.GetContents().Size() > 0 {
		ePackage, _ := eResource.GetContents().Get(0).(EPackage)
		ePackage.SetEFactoryInstance(NewEFactoryExt())
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
	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	eResource := xmlProcessor.Load(&URI{Path: "testdata/library.noroot.xml"})
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
	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	eResource := xmlProcessor.Load(&URI{Path: "testdata/library.complex.xml"})
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
	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	eResource := xmlProcessor.LoadWithOptions(
		&URI{Path: "testdata/library.complex.noroot.xml"},
		map[string]interface{}{
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

	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	eResource := xmlProcessor.Load(&URI{Path: "testdata/library.simple.invalid.xml"})
	require.NotNil(t, eResource)
	require.True(t, eResource.IsLoaded())
	require.False(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
}

func TestXMLDecoderSimpleEscapeXML(t *testing.T) {
	// load libray simple ecore	package
	ePackage := loadPackage("library.simple.ecore")
	require.NotNil(t, ePackage)

	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	eResource := xmlProcessor.Load(&URI{Path: "testdata/library.simple.escape.xml"})
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
	eResource := eResourceSet.CreateResource(&URI{Path: "testdata/library.simple.ids.xml"})
	require.NotNil(t, eResource)
	eResource.SetObjectIDManager(idManager)
	eResource.LoadWithOptions(map[string]interface{}{XML_OPTION_ID_ATTRIBUTE_NAME: "id"})
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

func TestXMLDecoderLibraryComplexBig(t *testing.T) {
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
	eResource := eResourceSet.CreateResource(&URI{Path: "$tmp.xml"})

	f, err := os.Open("testdata/book.simple.xml")
	require.NotNil(t, f)
	require.Nil(t, err)

	decoder := NewXMLDecoder(eResource, f, nil)
	eObject, err := decoder.DecodeObject()
	require.Nil(t, err)
	require.NotNil(t, eObject)
	assert.Equal(t, "Book 1", eObject.EGet(eBookNameAttribute))
}

func BenchmarkXMLDecoderLibraryComplexBig(b *testing.B) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(b, ePackage)

	// create resource
	uri := &URI{Path: "testdata/library.complex.big.xml"}
	eResource := NewEResourceImpl()
	eResource.SetURI(uri)
	eResourceSet := NewEResourceSetImpl()
	eResourceSet.GetResources().Add(eResource)
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)

	// get file content
	content, err := ioutil.ReadFile(uri.Path)
	require.Nil(b, err)
	r := bytes.NewReader(content)

	for i := 0; i < b.N; i++ {
		r.Seek(0, io.SeekStart)
		xmlDecoder := NewXMLDecoder(eResource, r, nil)
		xmlDecoder.Decode()
		require.True(b, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	}
}
