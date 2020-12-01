package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEReferenceExt_GetEReferenceType(t *testing.T) {
	r := newEReferenceExt()
	mockType := &MockEClass{}
	mockType.On("EIsProxy").Return(false).Once()
	r.SetEType(mockType)
	assert.Equal(t, mockType, r.GetEReferenceType())
	mock.AssertExpectationsForObjects(t, mockType)
}

func TestEReferenceExt_basicGetEReferenceType(t *testing.T) {
	r := newEReferenceExt()
	mockType := &MockEClass{}
	r.SetEType(mockType)
	assert.Equal(t, mockType, r.basicGetEReferenceType())
	mock.AssertExpectationsForObjects(t, mockType)
}
