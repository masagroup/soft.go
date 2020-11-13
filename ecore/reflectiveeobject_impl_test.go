package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestReflectiveEObjectImpl_GetAttribute(t *testing.T) {
	mockAttribute := new(MockEAttribute)
	mockAttribute.On("IsMany").Return(false).Once()

	mockClass := new(MockEClass)
	mockClass.On("GetFeatureCount").Return(2).Once()
	mockClass.On("GetEStructuralFeature", 0).Return(mockAttribute).Once()

	o := NewReflectiveEObjectImpl()
	o.setEClass(mockClass)

	// get unitialized value
	assert.Nil(t, o.EGetFromID(0, true))

	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute)
}

func TestReflectiveEObjectImpl_SetAttribute(t *testing.T) {
	mockAttribute := new(MockEAttribute)

	mockClass := new(MockEClass)
	mockClass.On("GetFeatureCount").Return(2).Once()
	mockClass.On("GetEStructuralFeature", 0).Return(mockAttribute).Twice()

	mockObject := new(MockEObject)

	o := NewReflectiveEObjectImpl()
	o.setEClass(mockClass)

	// set
	o.ESetFromID(0, mockObject)

	// check that value is well set
	assert.Equal(t, mockObject, o.EGetFromID(0, true))

	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute)
}

func TestReflectiveEObjectImpl_UnsetAttribute(t *testing.T) {
	mockAttribute := new(MockEAttribute)
	mockAttribute.On("IsMany").Return(false).Once()

	mockClass := new(MockEClass)
	mockClass.On("GetFeatureCount").Return(2).Once()
	mockClass.On("GetEStructuralFeature", 0).Return(mockAttribute).Times(3)

	mockObject := new(MockEObject)

	o := NewReflectiveEObjectImpl()
	o.setEClass(mockClass)

	// set - unset
	o.ESetFromID(0, mockObject)
	o.EUnsetFromID(0)

	// check that value is unset
	assert.Nil(t, o.EGetFromID(0, true))

	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute)
}

func TestReflectiveEObjectImpl_GetContainer(t *testing.T) {
	mockOpposite := new(MockEReference)
	mockOpposite.On("IsContainment").Return(true).Once()

	mockReference := new(MockEReference)
	mockReference.On("GetEOpposite").Return(mockOpposite).Once()

	mockClass := new(MockEClass)
	mockClass.On("GetEStructuralFeature", 0).Return(mockReference).Once()
	mockClass.On("GetFeatureID", mockReference).Return(0).Once()

	o := NewReflectiveEObjectImpl()
	o.setEClass(mockClass)

	// get non initialized container
	assert.Nil(t, o.EGetFromID(0, true))

	mock.AssertExpectationsForObjects(t, mockClass, mockReference, mockOpposite)
}

func TestReflectiveEObjectImpl_SetContainer(t *testing.T) {
	mockObject := new(MockEObjectInternal)
	mockObject.On("EResource").Return(nil).Once()
	mockObject.On("EIsProxy").Return(false).Once()

	mockOpposite := new(MockEReference)
	mockOpposite.On("IsContainment").Return(true).Times(3)

	mockReference := new(MockEReference)
	mockReference.On("GetEOpposite").Return(mockOpposite).Times(3)

	mockClass := new(MockEClass)
	mockClass.On("GetEStructuralFeature", 0).Return(mockReference).Times(3)
	mockClass.On("GetFeatureID", mockReference).Return(0).Times(3)

	o := NewReflectiveEObjectImpl()
	o.setEClass(mockClass)

	mockObject.On("EInverseAdd", o.GetInterfaces(), 0, nil).Return(nil)

	// set object reference as container
	o.ESetFromID(0, mockObject)

	// get unresolved
	assert.Equal(t, mockObject, o.EGetFromID(0, false))

	// get resolved
	assert.Equal(t, mockObject, o.EGetFromID(0, true))

	mock.AssertExpectationsForObjects(t, mockClass, mockReference, mockOpposite, mockObject)
}
