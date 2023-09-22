package ecore

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBinaryCodec_NewEncoder(t *testing.T) {
	c := &BinaryCodec{}
	mockResource := NewMockEResource(t)
	mockResource.EXPECT().GetURI().Return(nil).Once()
	e := c.NewEncoder(mockResource, nil, nil)
	require.NotNil(t, e)
	mock.AssertExpectationsForObjects(t, mockResource)
}

func TestBinaryCodec_NewDecoder(t *testing.T) {
	c := &BinaryCodec{}
	mockResource := NewMockEResource(t)
	mockResource.EXPECT().GetURI().Return(nil).Once()
	e := c.NewDecoder(mockResource, nil, nil)
	require.NotNil(t, e)
	mock.AssertExpectationsForObjects(t, mockResource)
}

func TestBinaryCodec_GetFeatureKind_Unknown(t *testing.T) {
	mockFeature := NewMockEStructuralFeature(t)
	assert.Equal(t, binaryFeatureKind(-1), getBinaryCodecFeatureKind(mockFeature))
}

func TestBinaryCodec_GetFeatureKind_Reference(t *testing.T) {
	mockReference := NewMockEReference(t)
	mockReference.EXPECT().IsContainment().Return(true).Once()
	mockReference.EXPECT().IsResolveProxies().Return(true).Once()
	mockReference.EXPECT().IsMany().Return(true).Once()
	assert.Equal(t, bfkObjectContainmentListProxy, getBinaryCodecFeatureKind(mockReference))
	mockReference.AssertExpectations(t)

	mockReference.EXPECT().IsContainment().Return(true).Once()
	mockReference.EXPECT().IsResolveProxies().Return(true).Once()
	mockReference.EXPECT().IsMany().Return(false).Once()
	assert.Equal(t, bfkObjectContainmentProxy, getBinaryCodecFeatureKind(mockReference))
	mockReference.AssertExpectations(t)

	mockReference.EXPECT().IsContainment().Return(true).Once()
	mockReference.EXPECT().IsResolveProxies().Return(false).Once()
	mockReference.EXPECT().IsMany().Return(true).Once()
	assert.Equal(t, bfkObjectContainmentList, getBinaryCodecFeatureKind(mockReference))
	mockReference.AssertExpectations(t)

	mockReference.EXPECT().IsContainment().Return(true).Once()
	mockReference.EXPECT().IsResolveProxies().Return(false).Once()
	mockReference.EXPECT().IsMany().Return(false).Once()
	assert.Equal(t, bfkObjectContainment, getBinaryCodecFeatureKind(mockReference))
	mockReference.AssertExpectations(t)

	mockReference.EXPECT().IsContainment().Return(false).Once()
	mockReference.EXPECT().IsContainer().Return(true).Once()
	mockReference.EXPECT().IsResolveProxies().Return(true).Once()
	assert.Equal(t, bfkObjectContainerProxy, getBinaryCodecFeatureKind(mockReference))
	mockReference.AssertExpectations(t)

	mockReference.EXPECT().IsContainment().Return(false).Once()
	mockReference.EXPECT().IsContainer().Return(true).Once()
	mockReference.EXPECT().IsResolveProxies().Return(false).Once()
	assert.Equal(t, bfkObjectContainer, getBinaryCodecFeatureKind(mockReference))
	mockReference.AssertExpectations(t)

	mockReference.EXPECT().IsContainment().Return(false).Once()
	mockReference.EXPECT().IsContainer().Return(false).Once()
	mockReference.EXPECT().IsResolveProxies().Return(true).Once()
	mockReference.EXPECT().IsMany().Return(true).Once()
	assert.Equal(t, bfkObjectListProxy, getBinaryCodecFeatureKind(mockReference))
	mockReference.AssertExpectations(t)

	mockReference.EXPECT().IsContainment().Return(false).Once()
	mockReference.EXPECT().IsContainer().Return(false).Once()
	mockReference.EXPECT().IsResolveProxies().Return(true).Once()
	mockReference.EXPECT().IsMany().Return(false).Once()
	assert.Equal(t, bfkObjectProxy, getBinaryCodecFeatureKind(mockReference))
	mockReference.AssertExpectations(t)

	mockReference.EXPECT().IsContainment().Return(false).Once()
	mockReference.EXPECT().IsContainer().Return(false).Once()
	mockReference.EXPECT().IsResolveProxies().Return(false).Once()
	mockReference.EXPECT().IsMany().Return(true).Once()
	assert.Equal(t, bfkObjectList, getBinaryCodecFeatureKind(mockReference))
	mockReference.AssertExpectations(t)

	mockReference.EXPECT().IsContainment().Return(false).Once()
	mockReference.EXPECT().IsContainer().Return(false).Once()
	mockReference.EXPECT().IsResolveProxies().Return(false).Once()
	mockReference.EXPECT().IsMany().Return(false).Once()
	assert.Equal(t, bfkObject, getBinaryCodecFeatureKind(mockReference))
	mockReference.AssertExpectations(t)
}

func TestBinaryCodec_GetFeatureKind_Attribute(t *testing.T) {
	mockAttribute := NewMockEAttribute(t)
	mockAttribute.EXPECT().IsMany().Return(true).Once()
	assert.Equal(t, bfkDataList, getBinaryCodecFeatureKind(mockAttribute))
	mockAttribute.AssertExpectations(t)

	mockEnum := NewMockEEnum(t)
	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().GetEAttributeType().Return(mockEnum).Once()
	assert.Equal(t, bfkEnum, getBinaryCodecFeatureKind(mockAttribute))
	mock.AssertExpectationsForObjects(t, mockAttribute, mockEnum)

	mockDataType := NewMockEDataType(t)
	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().GetEAttributeType().Return(mockDataType).Once()
	mockDataType.EXPECT().GetInstanceTypeName().Return("float64").Once()
	assert.Equal(t, bfkFloat64, getBinaryCodecFeatureKind(mockAttribute))
	mock.AssertExpectationsForObjects(t, mockAttribute, mockDataType)

	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().GetEAttributeType().Return(mockDataType).Once()
	mockDataType.EXPECT().GetInstanceTypeName().Return("float32").Once()
	assert.Equal(t, bfkFloat32, getBinaryCodecFeatureKind(mockAttribute))
	mock.AssertExpectationsForObjects(t, mockAttribute, mockDataType)

	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().GetEAttributeType().Return(mockDataType).Once()
	mockDataType.EXPECT().GetInstanceTypeName().Return("int").Once()
	assert.Equal(t, bfkInt, getBinaryCodecFeatureKind(mockAttribute))
	mock.AssertExpectationsForObjects(t, mockAttribute, mockDataType)

	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().GetEAttributeType().Return(mockDataType).Once()
	mockDataType.EXPECT().GetInstanceTypeName().Return("int64").Once()
	assert.Equal(t, bfkInt64, getBinaryCodecFeatureKind(mockAttribute))
	mock.AssertExpectationsForObjects(t, mockAttribute, mockDataType)

	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().GetEAttributeType().Return(mockDataType).Once()
	mockDataType.EXPECT().GetInstanceTypeName().Return("int32").Once()
	assert.Equal(t, bfkInt32, getBinaryCodecFeatureKind(mockAttribute))
	mock.AssertExpectationsForObjects(t, mockAttribute, mockDataType)

	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().GetEAttributeType().Return(mockDataType).Once()
	mockDataType.EXPECT().GetInstanceTypeName().Return("int16").Once()
	assert.Equal(t, bfkInt16, getBinaryCodecFeatureKind(mockAttribute))
	mock.AssertExpectationsForObjects(t, mockAttribute, mockDataType)

	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().GetEAttributeType().Return(mockDataType).Once()
	mockDataType.EXPECT().GetInstanceTypeName().Return("byte").Once()
	assert.Equal(t, bfkByte, getBinaryCodecFeatureKind(mockAttribute))
	mock.AssertExpectationsForObjects(t, mockAttribute, mockDataType)

	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().GetEAttributeType().Return(mockDataType).Once()
	mockDataType.EXPECT().GetInstanceTypeName().Return("bool").Once()
	assert.Equal(t, bfkBool, getBinaryCodecFeatureKind(mockAttribute))
	mock.AssertExpectationsForObjects(t, mockAttribute, mockDataType)

	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().GetEAttributeType().Return(mockDataType).Once()
	mockDataType.EXPECT().GetInstanceTypeName().Return("string").Once()
	assert.Equal(t, bfkString, getBinaryCodecFeatureKind(mockAttribute))
	mock.AssertExpectationsForObjects(t, mockAttribute, mockDataType)

	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().GetEAttributeType().Return(mockDataType).Once()
	mockDataType.EXPECT().GetInstanceTypeName().Return("[]byte").Once()
	assert.Equal(t, bfkByteArray, getBinaryCodecFeatureKind(mockAttribute))
	mock.AssertExpectationsForObjects(t, mockAttribute, mockDataType)

	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().GetEAttributeType().Return(mockDataType).Once()
	mockDataType.EXPECT().GetInstanceTypeName().Return("*time/time.Time").Once()
	assert.Equal(t, bfkDate, getBinaryCodecFeatureKind(mockAttribute))
	mock.AssertExpectationsForObjects(t, mockAttribute, mockDataType)

	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().GetEAttributeType().Return(mockDataType).Once()
	mockDataType.EXPECT().GetInstanceTypeName().Return("data").Once()
	assert.Equal(t, bfkData, getBinaryCodecFeatureKind(mockAttribute))
	mock.AssertExpectationsForObjects(t, mockAttribute, mockDataType)
}

func TestBinaryCodec_EncodeDecodeEcore(t *testing.T) {
	// load package
	ePackage := loadPackage("library.simple.ecore")
	require.NotNil(t, ePackage)

	// encode package resource in binary format
	buffer := bytes.Buffer{}
	c := &BinaryCodec{}
	encoder := c.NewEncoder(ePackage.EResource(), &buffer, nil)
	encoder.EncodeResource()

	// decode buffer into another resource
	eNewResource := NewEResourceImpl()
	decoder := c.NewDecoder(eNewResource, &buffer, nil)
	decoder.DecodeResource()
	require.True(t, eNewResource.GetErrors().Empty(), diagnosticError(eNewResource.GetErrors()))

	eNewPackage, _ := eNewResource.GetContents().Get(0).(EPackage)
	require.NotNil(t, eNewPackage)

	// retrieve document root class , library class & library name attribute
	eLibraryClass, _ := eNewPackage.GetEClassifier("Library").(EClass)
	require.NotNil(t, eLibraryClass)
	eLibraryOwnerAttribute, _ := eLibraryClass.GetEStructuralFeatureFromName("owner").(EAttribute)
	require.NotNil(t, eLibraryOwnerAttribute)
	eDataType := eLibraryOwnerAttribute.GetEAttributeType()
	require.NotNil(t, eDataType)
	assert.Equal(t, "EString", eDataType.GetName())

}

func loadTestPackage(t *testing.T, resourceSet EResourceSet, packageURI *URI) (EResource, EPackage) {
	// load package
	r := resourceSet.CreateResource(packageURI)
	r.SetObjectIDManager(NewIncrementalIDManager())
	r.Load()
	assert.True(t, r.IsLoaded())
	assert.True(t, r.GetErrors().Empty(), diagnosticError(r.GetErrors()))
	assert.True(t, r.GetWarnings().Empty(), diagnosticError(r.GetWarnings()))

	// retrieve package
	ePackage, _ := r.GetContents().Get(0).(EPackage)
	require.NotNil(t, ePackage)
	resourceSet.GetPackageRegistry().RegisterPackage(ePackage)
	return r, ePackage
}

func loadTestModel(t *testing.T, resourceSet EResourceSet, modelURI *URI) (EResource, EObject) {
	// load package
	r := resourceSet.CreateResource(modelURI)
	r.SetObjectIDManager(NewIncrementalIDManager())
	r.Load()
	require.True(t, r.IsLoaded())
	require.True(t, r.GetErrors().Empty(), diagnosticError(r.GetErrors()))
	require.True(t, r.GetWarnings().Empty(), diagnosticError(r.GetWarnings()))
	require.Equal(t, 1, r.GetContents().Size())

	// retrieve root object
	return r, r.GetContents().Get(0).(EObject)
}

func TestBinaryCodec_EncodeDecodeResource_WithReferences(t *testing.T) {
	eResourceSet := NewEResourceSetImpl()
	binaryCodecOptions := map[string]any{BINARY_OPTION_ID_ATTRIBUTE: true}
	// load packages & models
	eShopPackageResource, eShopPackage := loadTestPackage(t, eResourceSet, NewURI("testdata/shop.ecore"))
	require.NotNil(t, eShopPackage)
	require.NotNil(t, eShopPackageResource)
	eShopModelResource, eShopModel := loadTestModel(t, eResourceSet, NewURI("testdata/shop.xml"))
	require.NotNil(t, eShopModel)
	require.NotNil(t, eShopModelResource)

	eOrdersPackageResource, eOrdersPackage := loadTestPackage(t, eResourceSet, NewURI("testdata/orders.ecore"))
	require.NotNil(t, eOrdersPackageResource)
	require.NotNil(t, eOrdersPackage)
	eOrdersModelResource, eOrdersModel := loadTestModel(t, eResourceSet, NewURI("testdata/orders.xml"))
	require.NotNil(t, eOrdersModelResource)
	require.NotNil(t, eOrdersModel)
	ResolveAll(eShopModel)
	ResolveAll(eOrdersModel)

	// encode orders resource
	var buffer bytes.Buffer
	binaryEncoder := NewBinaryEncoder(eOrdersModelResource, &buffer, binaryCodecOptions)
	binaryEncoder.EncodeResource()
	require.True(t, eOrdersModelResource.GetErrors().Empty(), diagnosticError(eOrdersModelResource.GetErrors()))
	eResourceSet.GetResources().Remove(eOrdersModelResource)

	// decode orders resource
	eOrdersModelResource = NewEResourceImpl()
	eOrdersModelResource.SetObjectIDManager(NewIncrementalIDManager())
	eOrdersModelResource.SetURI(NewURI("testdata/orders.xml"))
	eResourceSet.GetResources().Add(eOrdersModelResource)
	binaryDecoder := NewBinaryDecoder(eOrdersModelResource, &buffer, binaryCodecOptions)
	binaryDecoder.DecodeResource()
	require.True(t, eOrdersModelResource.GetErrors().Empty(), diagnosticError(eOrdersModelResource.GetErrors()))

	eProducClass, _ := eShopPackage.GetEClassifier("Product").(EClass)
	require.NotNil(t, eProducClass)
	eProductNameAttribute, _ := eProducClass.GetEStructuralFeatureFromName("name").(EAttribute)
	require.NotNil(t, eProductNameAttribute)

	eOrdersClass, _ := eOrdersPackage.GetEClassifier("Orders").(EClass)
	require.NotNil(t, eOrdersClass)
	eOrderReference, _ := eOrdersClass.GetEStructuralFeatureFromName("order").(EReference)
	require.NotNil(t, eOrderReference)
	eOrderClass, _ := eOrdersPackage.GetEClassifier("Order").(EClass)
	require.NotNil(t, eOrderClass)
	eProductReference, _ := eOrderClass.GetEStructuralFeatureFromName("product").(EReference)
	require.NotNil(t, eProductReference)

	eOrders, _ := eOrdersModelResource.GetContents().Get(0).(EObject)
	require.NotNil(t, eOrders)
	assert.Equal(t, eOrdersClass, eOrders.EClass())

	eOrderList, _ := eOrders.EGet(eOrderReference).(EList)
	require.NotNil(t, eOrderList)
	eOrder, _ := eOrderList.Get(0).(EObject)
	require.NotNil(t, eOrder)
	assert.Equal(t, eOrderClass, eOrder.EClass())

	eOrderProduct, _ := eOrder.EGet(eProductReference).(EObject)
	require.NotNil(t, eOrderProduct)
	assert.False(t, eOrderProduct.EIsProxy())
	eOrderProductName, _ := eOrderProduct.EGet(eProductNameAttribute).(string)
	assert.Equal(t, "Product 0", eOrderProductName)
}

func TestBinaryCodec_EncodeDecodeObject_WithExternalReferences(t *testing.T) {
	eResourceSet := NewEResourceSetImpl()
	binaryCodecOptions := map[string]any{BINARY_OPTION_ID_ATTRIBUTE: true}
	eLibraryPackageResource, eLibraryPackage := loadTestPackage(t, eResourceSet, NewURI("testdata/library.complex.ecore"))
	require.NotNil(t, eLibraryPackageResource)
	require.NotNil(t, eLibraryPackage)
	eLibraryModelResource, eLibraryModel := loadTestModel(t, eResourceSet, NewURI("testdata/library.complex.xml"))
	require.NotNil(t, eLibraryModelResource)
	require.NotNil(t, eLibraryModel)

	// retrieve document root class , library class & library name attribute
	eDocumentRootClass, _ := eLibraryPackage.GetEClassifier("DocumentRoot").(EClass)
	assert.NotNil(t, eDocumentRootClass)
	eDocumentRootLibraryFeature, _ := eDocumentRootClass.GetEStructuralFeatureFromName("library").(EReference)
	assert.NotNil(t, eDocumentRootLibraryFeature)
	eLibraryClass, _ := eLibraryPackage.GetEClassifier("Library").(EClass)
	assert.NotNil(t, eLibraryClass)

	// book class and attributes
	eLibraryBooksRefeference, _ := eLibraryClass.GetEStructuralFeatureFromName("books").(EReference)
	assert.NotNil(t, eLibraryBooksRefeference)
	eBookClass, _ := eLibraryPackage.GetEClassifier("Book").(EClass)
	require.NotNil(t, eBookClass)
	eBookTitleAttribute, _ := eBookClass.GetEStructuralFeatureFromName("title").(EAttribute)
	require.NotNil(t, eBookTitleAttribute)
	eBookDateAttribute, _ := eBookClass.GetEStructuralFeatureFromName("publicationDate").(EAttribute)
	require.NotNil(t, eBookDateAttribute)
	eBookCategoryAttribute, _ := eBookClass.GetEStructuralFeatureFromName("category").(EAttribute)
	require.NotNil(t, eBookCategoryAttribute)
	eBookAuthorReference, _ := eBookClass.GetEStructuralFeatureFromName("author").(EReference)
	require.NotNil(t, eBookAuthorReference)

	// retrieve library
	eLibrary, _ := eLibraryModel.EGet(eDocumentRootLibraryFeature).(EObject)
	assert.NotNil(t, eLibrary)

	// retrieve book
	eBooks, _ := eLibrary.EGet(eLibraryBooksRefeference).(EList)
	assert.NotNil(t, eBooks)
	eBook := eBooks.Get(0).(EObject)
	require.NotNil(t, eBook)

	// check book name
	assert.Equal(t, "Title 0", eBook.EGet(eBookTitleAttribute))

	// retrieve author
	author := eBook.EGet(eBookAuthorReference).(EObject)
	require.NotNil(t, author)

	// encode book
	var buffer bytes.Buffer
	binaryEncoder := NewBinaryEncoder(eLibraryModelResource, &buffer, binaryCodecOptions)
	err := binaryEncoder.EncodeObject(eBook)
	require.Nil(t, err)

	// decode new book
	binaryDecoder := NewBinaryDecoder(eLibraryModelResource, &buffer, binaryCodecOptions)
	eNewBook, err := binaryDecoder.DecodeObject()
	require.Nil(t, err)
	require.NotEqual(t, eBook, eNewBook)

	// check book name
	assert.Equal(t, "Title 0", eNewBook.EGet(eBookTitleAttribute))

	// retrieve new author
	// new book is not in the resource and then in the hierarchy - so author must remain a proxy
	eProxyAuthor := eNewBook.EGet(eBookAuthorReference).(EObject)
	require.NotNil(t, eProxyAuthor)
	assert.True(t, eProxyAuthor.EIsProxy())

	// retrieve new author
	// add new book in the resource and then resolve author
	eBooks.Add(eNewBook)
	eNewAuthor := eNewBook.EGet(eBookAuthorReference).(EObject)
	require.NotNil(t, eNewAuthor)
	assert.False(t, eNewAuthor.EIsProxy())
	assert.Equal(t, author, eNewAuthor)

}

func TestBinaryCodec_EncodeDecodeObject_WithInternalReferences(t *testing.T) {
	eResourceSet := NewEResourceSetImpl()
	binaryCodecOptions := map[string]any{BINARY_OPTION_ID_ATTRIBUTE: true}
	eLibraryPackageResource, eLibraryPackage := loadTestPackage(t, eResourceSet, NewURI("testdata/library.complex.ecore"))
	require.NotNil(t, eLibraryPackageResource)
	require.NotNil(t, eLibraryPackage)
	eLibraryModelResource, eLibraryModel := loadTestModel(t, eResourceSet, NewURI("testdata/library.complex.xml"))
	require.NotNil(t, eLibraryModelResource)
	require.NotNil(t, eLibraryModel)

	// retrieve document root class , library class & library name attribute
	eDocumentRootClass, _ := eLibraryPackage.GetEClassifier("DocumentRoot").(EClass)
	assert.NotNil(t, eDocumentRootClass)
	eDocumentRootLibraryFeature, _ := eDocumentRootClass.GetEStructuralFeatureFromName("library").(EReference)
	assert.NotNil(t, eDocumentRootLibraryFeature)
	eLibraryClass, _ := eLibraryPackage.GetEClassifier("Library").(EClass)
	assert.NotNil(t, eLibraryClass)

	// writers
	eLibraryWritersRefeference, _ := eLibraryClass.GetEStructuralFeatureFromName("writers").(EReference)
	assert.NotNil(t, eLibraryWritersRefeference)
	eWriterClass, _ := eLibraryPackage.GetEClassifier("Writer").(EClass)
	require.NotNil(t, eWriterClass)
	eWriterAddressAttribute := eWriterClass.GetEStructuralFeatureFromName("address")
	require.NotNil(t, eWriterAddressAttribute)
	eWriterBooksReference, _ := eWriterClass.GetEStructuralFeatureFromName("books").(EReference)
	require.NotNil(t, eWriterBooksReference)
	// book
	eBookClass, _ := eLibraryPackage.GetEClassifier("Book").(EClass)
	require.NotNil(t, eBookClass)
	eBookTitleAttribute, _ := eBookClass.GetEStructuralFeatureFromName("title").(EAttribute)
	require.NotNil(t, eBookTitleAttribute)
	eBookAuthorReference, _ := eBookClass.GetEStructuralFeatureFromName("author").(EReference)
	require.NotNil(t, eBookAuthorReference)

	// retrieve library
	eLibrary, _ := eLibraryModel.EGet(eDocumentRootLibraryFeature).(EObject)
	assert.NotNil(t, eLibrary)

	// encode book
	var buffer bytes.Buffer
	binaryEncoder := NewBinaryEncoder(eLibraryModelResource, &buffer, binaryCodecOptions)
	err := binaryEncoder.EncodeObject(eLibrary)
	require.Nil(t, err)

	// decode new book
	binaryDecoder := NewBinaryDecoder(eLibraryModelResource, &buffer, binaryCodecOptions)
	eNewLibrary, err := binaryDecoder.DecodeObject()
	require.Nil(t, err)
	require.NotEqual(t, eLibrary, eNewLibrary)

	// first author
	authors, _ := eNewLibrary.EGet(eLibraryWritersRefeference).(EList)
	require.NotNil(t, authors)
	assert.Equal(t, 1, authors.Size())
	author := authors.Get(0).(EObject)
	assert.Equal(t, "Adress 0", author.EGet(eWriterAddressAttribute))

	// check author books
	books, _ := author.EGet(eWriterBooksReference).(EList)
	require.NotNil(t, books)
	assert.Equal(t, 2, books.Size())
	book := books.Get(0).(EObject)
	require.NotNil(t, book)
	assert.Equal(t, "Title 0", book.EGet(eBookTitleAttribute))

	// check that author books is exactly the same object
	bookAuthor, _ := book.EGet(eBookAuthorReference).(EObject)
	require.NotNil(t, bookAuthor)
	assert.Equal(t, author, bookAuthor)

}
