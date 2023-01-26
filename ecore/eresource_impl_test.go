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
	mockEAdapter.On("SetTarget", r).Once()
	r.EAdapters().Add(mockEAdapter)
	mock.AssertExpectationsForObjects(t, mockEAdapter)

	u, err := ParseURI("https://example.com/foo%2fbar")
	assert.Nil(t, err)

	mockEAdapter.On("NotifyChanged", mock.Anything).Once()
	r.SetURI(u)
	assert.Equal(t, u, r.GetURI())
	mock.AssertExpectationsForObjects(t, mockEAdapter)
}

func TestResourceContents(t *testing.T) {
	r := NewEResourceImpl()

	mockEObjectInternal := NewMockEObjectInternal(t)
	mockEObjectInternal.On("ESetResource", r, mock.Anything).Return(nil)
	r.GetContents().Add(mockEObjectInternal)

	mockEObjectInternal.On("ESetResource", nil, mock.Anything).Return(nil)
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

		mockObject.On("EClass").Return(mockClass).Once()
		mockClass.On("GetEIDAttribute").Return(mockAttribute).Once()
		mockObject.On("EIsSet", mockAttribute).Return(true).Once()
		mockObject.On("EGet", mockAttribute).Return(mockIDValue).Once()
		mockAttribute.On("GetEAttributeType").Return(mockDataType).Once()
		mockDataType.On("GetEPackage").Return(mockPackage).Once()
		mockPackage.On("GetEFactoryInstance").Return(mockFactory).Once()
		mockFactory.On("ConvertToString", mockDataType, mockIDValue).Return("id1").Once()
		assert.Equal(t, "id1", r.GetURIFragment(mockObject))
		mock.AssertExpectationsForObjects(t, mockObject, mockClass, mockAttribute, mockDataType, mockPackage, mockFactory)
	}
	// root
	{
		r := NewEResourceImpl()
		mockObject := NewMockEObjectInternal(t)
		mockObject.On("ESetResource", r, nil).Return(nil).Once()
		r.GetContents().Add(mockObject)
		mock.AssertExpectationsForObjects(t, mockObject)

		mockClass := NewMockEClass(t)
		mockObject.On("EClass").Return(mockClass).Twice()
		mockClass.On("GetEIDAttribute").Return(nil).Twice()
		mockObject.On("EInternalResource").Return(r).Once()

		assert.Equal(t, "/", r.GetURIFragment(mockObject))
		mock.AssertExpectationsForObjects(t, mockObject, mockClass)
	}
	// two roots
	{
		r := NewEResourceImpl()
		mockObject1 := NewMockEObjectInternal(t)
		mockObject2 := NewMockEObjectInternal(t)
		mockObject1.On("ESetResource", r, mock.Anything).Return(nil).Once()
		mockObject2.On("ESetResource", r, mock.Anything).Return(nil).Once()
		r.GetContents().AddAll(NewImmutableEList([]any{mockObject1, mockObject2}))
		mock.AssertExpectationsForObjects(t, mockObject1, mockObject2)

		mockClass := NewMockEClass(t)
		mockObject1.On("EClass").Return(mockClass).Twice()
		mockClass.On("GetEIDAttribute").Return(nil).Twice()
		mockObject1.On("EInternalResource").Return(r).Once()

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
		mockObject.On("EClass").Return(mockClass).Twice()
		mockClass.On("GetEIDAttribute").Return(nil).Twice()
		mockObject.On("EInternalResource").Return(nil).Once()
		mockObject.On("EInternalContainer").Return(mockRoot).Once()
		mockObject.On("EContainingFeature").Return(mockFeature).Once()
		mockRoot.On("EURIFragmentSegment", mockFeature, mockObject).Return("@fragment").Once()
		mockRoot.On("EInternalResource").Return(r).Once()
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
		mockObject.On("EClass").Return(mockClass).Once()
		mockClass.On("GetEIDAttribute").Return(nil).Once()
		mockObject.On("EInternalResource").Return(nil).Once()
		mockObject.On("EInternalContainer").Return(mockRoot).Once()
		mockIDManager.On("GetID", mockObject).Return("objectID").Once()
		mockRoot.On("EInternalResource").Return(r).Once()
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
	mockIDManager.On("Register", eLibrary).Once()
	eResource.GetContents().Add(eLibrary)
	mock.AssertExpectationsForObjects(t, mockIDManager)

	// create 2 books and add them to library
	eBookClass := ePackage.GetEClassifier("Book").(EClass)
	eLibraryBooksReference := eLibraryClass.GetEStructuralFeatureFromName("books").(EReference)
	eBookList := eLibrary.EGet(eLibraryBooksReference).(EList)
	eBook1 := eFactory.Create(eBookClass)
	eBook2 := eFactory.Create(eBookClass)
	mockIDManager.On("Register", eBook1).Once()
	mockIDManager.On("Register", eBook2).Once()
	eBookList.AddAll(NewImmutableEList([]any{eBook1, eBook2}))
	mock.AssertExpectationsForObjects(t, mockIDManager)
}

func TestResourceListeners(t *testing.T) {
	mockListener := NewMockEResourceListener(t)
	mockObject := NewMockEObject(t)
	eResource := NewEResourceImpl()
	eResource.GetResourceListeners().Add(mockListener)

	mockListener.On("Attached", mockObject).Once()
	mockObject.On("EContents").Return(NewEmptyImmutableEList())
	eResource.Attached(mockObject)

	mockListener.On("Detached", mockObject).Once()
	mockObject.On("EContents").Return(NewEmptyImmutableEList())
	eResource.Detached(mockObject)
}
