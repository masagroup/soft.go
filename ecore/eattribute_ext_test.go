package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEAttributeEClass(t *testing.T) {
	assert.Equal(t, GetPackage().GetEAttribute(), GetFactory().CreateEAttribute().EClass())
}

func TestEAttribute_GetEAttributeType(t *testing.T) {
	mockType := new(MockEDataType)
	a := newEAttributeExt()
	a.SetEType(mockType)

	mockType.On("EIsProxy").Return(false).Once()
	assert.Equal(t, mockType, a.GetEAttributeType())
	mockType.AssertExpectations(t)
}

func TestEAttribute_BasicGetEAttributeType(t *testing.T) {
	mockType := new(MockEDataType)
	a := newEAttributeExt()
	a.SetEType(mockType)
	assert.Equal(t, mockType, a.basicGetEAttributeType())
}
