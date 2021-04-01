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
	mockEAdapter := new(MockEAdapter)
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

	mockEObjectInternal := new(MockEObjectInternal)
	mockEObjectInternal.On("ESetResource", r, mock.Anything).Return(nil)
	r.GetContents().Add(mockEObjectInternal)

	mockEObjectInternal.On("ESetResource", nil, mock.Anything).Return(nil)
	r.GetContents().Remove(mockEObjectInternal)
}

func TestResourceLoadInvalid(t *testing.T) {
	r := NewEResourceImpl()
	r.SetURI(&URI{Path: "testdata/invalid.xml"})
	r.Load()
	assert.False(t, r.IsLoaded())
	assert.False(t, r.GetErrors().Empty())
}

func TestResourceGetURIFragment(t *testing.T) {

	// id attribute
	{
		r := NewEResourceImpl()
		mockObject := &MockEObject{}
		mockClass := &MockEClass{}
		mockAttribute := &MockEAttribute{}
		mockDataType := &MockEDataType{}
		mockPackage := &MockEPackage{}
		mockFactory := &MockEFactory{}
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
		mockObject := &MockEObjectInternal{}
		mockObject.On("ESetResource", r, nil).Return(nil).Once()
		r.GetContents().Add(mockObject)
		mock.AssertExpectationsForObjects(t, mockObject)

		mockClass := &MockEClass{}
		mockObject.On("EClass").Return(mockClass).Twice()
		mockClass.On("GetEIDAttribute").Return(nil).Twice()
		mockObject.On("EInternalResource").Return(r).Once()

		assert.Equal(t, "/", r.GetURIFragment(mockObject))
		mock.AssertExpectationsForObjects(t, mockObject, mockClass)
	}
	// two roots
	{
		r := NewEResourceImpl()
		mockObject1 := &MockEObjectInternal{}
		mockObject2 := &MockEObjectInternal{}
		mockObject1.On("ESetResource", r, mock.Anything).Return(nil).Once()
		mockObject2.On("ESetResource", r, mock.Anything).Return(nil).Once()
		r.GetContents().AddAll(NewImmutableEList([]interface{}{mockObject1, mockObject2}))
		mock.AssertExpectationsForObjects(t, mockObject1, mockObject2)

		mockClass := &MockEClass{}
		mockObject1.On("EClass").Return(mockClass).Twice()
		mockClass.On("GetEIDAttribute").Return(nil).Twice()
		mockObject1.On("EInternalResource").Return(r).Once()

		assert.Equal(t, "/0", r.GetURIFragment(mockObject1))
		mock.AssertExpectationsForObjects(t, mockObject1, mockObject2, mockClass)
	}
	// element - no id
	{
		r := NewEResourceImpl()
		mockRoot := &MockEObjectInternal{}
		mockObject := &MockEObjectInternal{}

		mockClass := &MockEClass{}
		mockFeature := &MockEStructuralFeature{}
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

		mockRoot := &MockEObjectInternal{}
		mockObject := &MockEObjectInternal{}
		mockClass := &MockEClass{}
		mockFeature := &MockEStructuralFeature{}
		mockObject.On("EClass").Return(mockClass).Once()
		mockClass.On("GetEIDAttribute").Return(nil).Once()
		mockObject.On("EInternalResource").Return(nil).Once()
		mockObject.On("EInternalContainer").Return(mockRoot).Once()
		mockObject.On("EContainingFeature").Return(mockFeature).Once()
		mockRoot.On("EURIFragmentSegment", mockFeature, mockObject).Return("@fragment").Once()
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
