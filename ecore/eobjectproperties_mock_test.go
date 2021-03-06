package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMockEObjectInternal_EDynamicGet(t *testing.T) {
	o := &MockEObjectProperties{}
	obj := &MockEObject{}
	o.On("EDynamicGet", 1).Once().Return(obj)
	o.On("EDynamicGet", 1).Once().Return(func(dynamicFeatureID int) interface{} {
		return obj
	})
	assert.Equal(t, obj, o.EDynamicGet(1))
	assert.Equal(t, obj, o.EDynamicGet(1))
	mock.AssertExpectationsForObjects(t, o, obj)
}

func TestMockEObjectInternal_EDynamicSet(t *testing.T) {
	o := &MockEObjectProperties{}
	obj := &MockEObject{}
	o.On("EDynamicSet", 1, obj).Once()
	o.EDynamicSet(1, obj)
	mock.AssertExpectationsForObjects(t, o, obj)
}

func TestMockEObjectInternal_EDynamicUnset(t *testing.T) {
	o := &MockEObjectProperties{}
	o.On("EDynamicUnset", 1).Once()
	o.EDynamicUnset(1)
	mock.AssertExpectationsForObjects(t, o)
}
