package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBasicEStoreListAccessors(t *testing.T) {
	mockOwner := &MockEObject{}
	mockFeature := &MockEStructuralFeature{}
	mockStore := &MockEStore{}
	list := NewBasicEStoreList(mockOwner, mockFeature, mockStore)

	assert.Equal(t, mockOwner, list.GetOwner())
	assert.Equal(t, mockFeature, list.GetFeature())

	mockClass := &MockEClass{}
	mockClass.On("GetFeatureID", mockFeature).Return(0).Once()
	mockOwner.On("EClass").Return(mockClass).Once()
	assert.Equal(t, 0, list.GetFeatureID())
	mock.AssertExpectationsForObjects(t, mockOwner, mockClass, mockFeature, mockStore)
}
