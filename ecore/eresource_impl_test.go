package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestResourceURI(t *testing.T) {
	u, err := ParseURI("https://example.com/foo%2fbar")
	assert.Nil(t, err)
	r := NewEResourceImpl()
	r.SetURI(u)
	assert.Equal(t, u, r.GetURI())
}

func TestResourceURINotifications(t *testing.T) {
	r := NewEResourceImpl()
	mockEAdapter := NewMockEAdapter(t)
	mockEAdapter.EXPECT().SetTarget(r).Once()
	r.EAdapters().Add(mockEAdapter)
	mock.AssertExpectationsForObjects(t, mockEAdapter)

	u, err := ParseURI("https://example.com/foo%2fbar")
	assert.Nil(t, err)

	mockEAdapter.EXPECT().NotifyChanged(mock.Anything).Once()
	r.SetURI(u)
	assert.Equal(t, u, r.GetURI())
	mock.AssertExpectationsForObjects(t, mockEAdapter)
}

func TestResourceContents(t *testing.T) {
	r := NewEResourceImpl()

	mockEObjectInternal := NewMockEObjectInternal(t)
	mockEObjectInternal.EXPECT().ESetResource(r, mock.Anything).Return(nil)
	r.GetContents().Add(mockEObjectInternal)

	mockEObjectInternal.EXPECT().ESetResource(nil, mock.Anything).Return(nil)
	r.GetContents().Remove(mockEObjectInternal)
}

func TestResourceLoadInvalid(t *testing.T) {
	r := NewEResourceImpl()
	r.SetURI(NewURI("testdata/invalid.xml"))
	r.Load()
	assert.False(t, r.IsLoaded())
	assert.False(t, r.GetErrors().Empty())
}

func TestResourceGetURIFragment(t *testing.T) {

	// id attribute
	{
		r := NewEResourceImpl()
		mockObject := NewMockEObject(t)
		mockClass := NewMockEClass(t)
		mockAttribute := NewMockEAttribute(t)
		mockDataType := NewMockEDataType(t)
		mockPackage := NewMockEPackage(t)
		mockFactory := NewMockEFactory(t)
		mockIDValue := 0

		mockObject.EXPECT().EClass().Return(mockClass).Once()
		mockClass.EXPECT().GetEIDAttribute().Return(mockAttribute).Once()
		mockObject.EXPECT().EIsSet(mockAttribute).Return(true).Once()
		mockObject.EXPECT().EGet(mockAttribute).Return(mockIDValue).Once()
		mockAttribute.EXPECT().GetEAttributeType().Return(mockDataType).Once()
		mockDataType.EXPECT().GetEPackage().Return(mockPackage).Once()
		mockPackage.EXPECT().GetEFactoryInstance().Return(mockFactory).Once()
		mockFactory.EXPECT().ConvertToString(mockDataType, mockIDValue).Return("id1").Once()
		assert.Equal(t, "id1", r.GetURIFragment(mockObject))
		mock.AssertExpectationsForObjects(t, mockObject, mockClass, mockAttribute, mockDataType, mockPackage, mockFactory)
	}
	// root
	{
		r := NewEResourceImpl()
		mockObject := NewMockEObjectInternal(t)
		mockObject.EXPECT().ESetResource(r, nil).Return(nil).Once()
		r.GetContents().Add(mockObject)
		mock.AssertExpectationsForObjects(t, mockObject)

		mockClass := NewMockEClass(t)
		mockObject.EXPECT().EClass().Return(mockClass).Twice()
		mockClass.EXPECT().GetEIDAttribute().Return(nil).Twice()
		mockObject.EXPECT().EInternalResource().Return(r).Once()

		assert.Equal(t, "/", r.GetURIFragment(mockObject))
		mock.AssertExpectationsForObjects(t, mockObject, mockClass)
	}
	// two roots
	{
		r := NewEResourceImpl()
		mockObject1 := NewMockEObjectInternal(t)
		mockObject2 := NewMockEObjectInternal(t)
		mockObject1.EXPECT().ESetResource(r, mock.Anything).Return(nil).Once()
		mockObject2.EXPECT().ESetResource(r, mock.Anything).Return(nil).Once()
		r.GetContents().AddAll(NewImmutableEList([]any{mockObject1, mockObject2}))
		mock.AssertExpectationsForObjects(t, mockObject1, mockObject2)

		mockClass := NewMockEClass(t)
		mockObject1.EXPECT().EClass().Return(mockClass).Twice()
		mockClass.EXPECT().GetEIDAttribute().Return(nil).Twice()
		mockObject1.EXPECT().EInternalResource().Return(r).Once()

		assert.Equal(t, "/0", r.GetURIFragment(mockObject1))
		mock.AssertExpectationsForObjects(t, mockObject1, mockObject2, mockClass)
	}
	// element - no id
	{
		r := NewEResourceImpl()
		mockRoot := NewMockEObjectInternal(t)
		mockObject := NewMockEObjectInternal(t)

		mockClass := NewMockEClass(t)
		mockFeature := NewMockEStructuralFeature(t)
		mockObject.EXPECT().EClass().Return(mockClass).Twice()
		mockClass.EXPECT().GetEIDAttribute().Return(nil).Twice()
		mockObject.EXPECT().EInternalResource().Return(nil).Once()
		mockObject.EXPECT().EInternalContainer().Return(mockRoot).Once()
		mockObject.EXPECT().EContainingFeature().Return(mockFeature).Once()
		mockRoot.EXPECT().EURIFragmentSegment(mockFeature, mockObject).Return("@fragment").Once()
		mockRoot.EXPECT().EInternalResource().Return(r).Once()
		assert.Equal(t, "//@fragment", r.GetURIFragment(mockObject))
		mock.AssertExpectationsForObjects(t, mockObject, mockRoot, mockClass, mockFeature)
	}
	// element - no id attribute - id manager
	{
		r := NewEResourceImpl()
		mockIDManager := &MockEObjectIDManager{}
		r.SetObjectIDManager(mockIDManager)

		mockRoot := NewMockEObjectInternal(t)
		mockObject := NewMockEObjectInternal(t)
		mockClass := NewMockEClass(t)
		mockObject.EXPECT().EClass().Return(mockClass).Once()
		mockClass.EXPECT().GetEIDAttribute().Return(nil).Once()
		mockObject.EXPECT().EInternalResource().Return(nil).Once()
		mockObject.EXPECT().EInternalContainer().Return(mockRoot).Once()
		mockIDManager.EXPECT().GetID(mockObject).Return("objectID").Once()
		mockRoot.EXPECT().EInternalResource().Return(r).Once()
		assert.Equal(t, "objectID", r.GetURIFragment(mockObject))
	}
}

func TestResourceGetEObject(t *testing.T) {
	{
		r := NewEResourceImpl()
		assert.Nil(t, r.GetEObject("Test"))
	}
}

func TestResourceIDManager(t *testing.T) {
	// load libray simple ecore	package
	ePackage := loadPackage("library.simple.ecore")
	assert.NotNil(t, ePackage)

	// create a resource with an id manager
	mockIDManager := &MockEObjectIDManager{}
	eResource := NewEResourceImpl()
	eResource.SetObjectIDManager(mockIDManager)

	// create a library and add it to resource
	eFactory := ePackage.GetEFactoryInstance()
	eLibraryClass := ePackage.GetEClassifier("Library").(EClass)
	eLibrary := eFactory.Create(eLibraryClass)
	mockIDManager.EXPECT().Register(eLibrary).Once()
	eResource.GetContents().Add(eLibrary)
	mock.AssertExpectationsForObjects(t, mockIDManager)

	// create 2 books and add them to library
	eBookClass := ePackage.GetEClassifier("Book").(EClass)
	eLibraryBooksReference := eLibraryClass.GetEStructuralFeatureFromName("books").(EReference)
	eBookList := eLibrary.EGet(eLibraryBooksReference).(EList)
	eBook1 := eFactory.Create(eBookClass)
	eBook2 := eFactory.Create(eBookClass)
	mockIDManager.EXPECT().Register(eBook1).Once()
	mockIDManager.EXPECT().Register(eBook2).Once()
	eBookList.AddAll(NewImmutableEList([]any{eBook1, eBook2}))
	mock.AssertExpectationsForObjects(t, mockIDManager)
}

func TestResourceListeners(t *testing.T) {
	mockListener := NewMockEResourceListener(t)
	mockObject := NewMockEObject(t)
	eResource := NewEResourceImpl()
	eResource.GetResourceListeners().Add(mockListener)

	mockListener.EXPECT().Attached(mockObject).Once()
	mockObject.EXPECT().EContents().Return(NewEmptyImmutableEList())
	eResource.Attached(mockObject)

	mockListener.EXPECT().Detached(mockObject).Once()
	mockObject.EXPECT().EContents().Return(NewEmptyImmutableEList())
	eResource.Detached(mockObject)
}
