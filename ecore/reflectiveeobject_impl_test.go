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

	o.ESetFromID(0, mockObject)

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

	o.ESetFromID(0, mockObject)
	o.EUnsetFromID(0)

	assert.Nil(t, o.EGetFromID(0, true))

	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute)
}
