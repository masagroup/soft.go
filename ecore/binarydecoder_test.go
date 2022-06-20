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

func TestBinaryDecoder_Complex(t *testing.T) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	//
	uri := &URI{Path: "testdata/library.complex.bin"}
	eResource := NewEResourceImpl()
	eResource.SetURI(uri)
	eResourceSet := NewEResourceSetImpl()
	eResourceSet.GetResources().Add(eResource)
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)

	// file
	f, err := os.Open(uri.Path)
	require.Nil(t, err)

	binaryDecoder := NewBinaryDecoder(eResource, f, nil)
	binaryDecoder.Decode()
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))

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
	eDocumentRoot, _ := eResource.GetContents().Get(0).(EObject)
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
	eBookCategoryAttribute, _ := eBookClass.GetEStructuralFeatureFromName("category").(EAttribute)
	require.NotNil(t, eBookCategoryAttribute)
	eBookAuthorReference, _ := eBookClass.GetEStructuralFeatureFromName("author").(EReference)
	require.NotNil(t, eBookAuthorReference)

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
	assert.Equal(t, time.Date(2015, time.September, 6, 4, 24, 46, 0, time.UTC), *date)

	// check book category
	category := eBook.EGet(eBookCategoryAttribute)
	assert.Equal(t, 2, category)

	// check author
	author := eBook.EGet(eBookAuthorReference).(EObject)
	require.NotNil(t, author)

	eWriterClass, _ := ePackage.GetEClassifier("Writer").(EClass)
	require.NotNil(t, eWriterClass)
	eWriterNameAttribute := eWriterClass.GetEStructuralFeatureFromName("firstName").(EAttribute)
	require.NotNil(t, eWriterNameAttribute)
	authorName := author.EGet(eWriterNameAttribute).(string)
	assert.Equal(t, "First Name 0", authorName)
}

func TestBinaryDecoder_ComplexWithID(t *testing.T) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	//
	uri := &URI{Path: "testdata/library.complex.id.bin"}
	idManager := NewUniqueIDManager(20)
	eResource := NewEResourceImpl()
	eResource.SetURI(uri)
	eResource.SetObjectIDManager(idManager)
	eResourceSet := NewEResourceSetImpl()
	eResourceSet.GetResources().Add(eResource)
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)

	// file
	f, err := os.Open(uri.Path)
	require.Nil(t, err)

	binaryDecoder := NewBinaryDecoder(eResource, f, nil)
	binaryDecoder.Decode()
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))

	// retrieve document root class , library class & library name attribute
	eDocumentRootClass, _ := ePackage.GetEClassifier("DocumentRoot").(EClass)
	require.NotNil(t, eDocumentRootClass)
	eDocumentRootLibraryFeature, _ := eDocumentRootClass.GetEStructuralFeatureFromName("library").(EReference)
	require.NotNil(t, eDocumentRootLibraryFeature)

	// check ids for document root and library
	eDocumentRoot, _ := eResource.GetContents().Get(0).(EObject)
	require.NotNil(t, eDocumentRoot)
	eLibrary, _ := eDocumentRoot.EGet(eDocumentRootLibraryFeature).(EObject)
	require.NotNil(t, eLibrary)
	assert.Equal(t, "h0Rz1FjVeBXUgaW3OzT2frUce90=", idManager.GetID(eDocumentRoot))
	assert.Equal(t, "d13pf-ypXLeIySkWAX03JcP-TbA=", idManager.GetID(eLibrary))

}

func TestBinaryDecoder_SimpleWithEDataTypeList(t *testing.T) {
	// load package
	ePackage := loadPackage("library.datalist.ecore")
	require.NotNil(t, ePackage)

	//
	uri := &URI{Path: "testdata/library.datalist.bin"}
	idManager := NewUniqueIDManager(20)
	eResource := NewEResourceImpl()
	eResource.SetURI(uri)
	eResource.SetObjectIDManager(idManager)
	eResourceSet := NewEResourceSetImpl()
	eResourceSet.GetResources().Add(eResource)
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)

	// file
	f, err := os.Open(uri.Path)
	require.Nil(t, err)

	binaryDecoder := NewBinaryDecoder(eResource, f, nil)
	binaryDecoder.Decode()
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))

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

func TestBinaryDecoder_ComplexBig(t *testing.T) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	//
	uri := &URI{Path: "testdata/library.complex.big.bin"}
	eResource := NewEResourceImpl()
	eResource.SetURI(uri)
	eResourceSet := NewEResourceSetImpl()
	eResourceSet.GetResources().Add(eResource)
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)

	// file
	f, err := os.Open(uri.Path)
	require.Nil(t, err)

	binaryDecoder := NewBinaryDecoder(eResource, f, nil)
	binaryDecoder.Decode()
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
}

func TestBinaryDecoder_Maps(t *testing.T) {
	// load package
	ePackage := loadPackage("emap.ecore")
	require.NotNil(t, ePackage)

	//
	uri := &URI{Path: "testdata/emap.bin"}
	eResource := NewEResourceImpl()
	eResource.SetURI(uri)
	eResourceSet := NewEResourceSetImpl()
	eResourceSet.GetResources().Add(eResource)
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)

	// file
	f, err := os.Open(uri.Path)
	require.Nil(t, err)

	binaryDecoder := NewBinaryDecoder(eResource, f, nil)
	binaryDecoder.Decode()
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))

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

	mapTest := eResource.GetContents().Get(0).(EObject)
	require.Equal(t, eMapTestClass, mapTest.EClass())
	mapKeyToValue, _ := mapTest.EGet(eMapTestKeyToValueReference).(EMap)
	require.NotNil(t, mapKeyToValue)
	assert.Equal(t, 5, mapKeyToValue.Size())
	mapKeyToInt, _ := mapTest.EGet(eMapTestKeyToIntReference).(EMap)
	require.NotNil(t, mapKeyToInt)
	assert.Equal(t, 5, mapKeyToInt.Size())
}

func BenchmarkBinaryDecoderLibraryComplexBig(b *testing.B) {
	// load package
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(b, ePackage)

	// create resource
	uri := &URI{Path: "testdata/library.complex.big.bin"}
	eResource := NewEResourceImpl()
	eResource.SetURI(uri)
	eResourceSet := NewEResourceSetImpl()
	eResourceSet.GetResources().Add(eResource)
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)

	// get file content
	content, err := ioutil.ReadFile(uri.Path)
	require.Nil(b, err)
	require.Nil(b, err)
	r := bytes.NewReader(content)

	for i := 0; i < b.N; i++ {
		r.Seek(0, io.SeekStart)
		binaryDecoder := NewBinaryDecoder(eResource, r, nil)
		binaryDecoder.Decode()
		require.True(b, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	}
}
