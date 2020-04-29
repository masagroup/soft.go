// *****************************************************************************
//
// This file is part of a MASA library or program.
// Refer to the included end-user license agreement for restrictions.
//
// Copyright (c) 2020 MASA Group
//
// *****************************************************************************

// *****************************************************************************
//
// Warning: This file was generated by soft.generator.go Generator
//
// *****************************************************************************

package ecore

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func discardEGenericType() {
	_ = assert.Equal
	_ = mock.Anything
	_ = testing.Coverage

}

func TestEGenericTypeEUpperBoundSet(t *testing.T) {
	var newValue *MockEGenericType = &MockEGenericType{}
	obj := newEGenericTypeImpl()
	mockAdapter := &MockEAdapter{}
	mockAdapter.On("SetTarget", obj).Once()
	mockAdapter.On("NotifyChanged", mock.Anything).Once()
	obj.EAdapters().Add(mockAdapter)
	obj.SetEUpperBound(newValue)
	mockAdapter.AssertExpectations(t)
}

func TestEGenericTypeETypeArgumentsGetList(t *testing.T) {
	obj := newEGenericTypeImpl()
	assert.NotNil(t, obj.GetETypeArguments())
}

func TestEGenericTypeELowerBoundSet(t *testing.T) {
	var newValue *MockEGenericType = &MockEGenericType{}
	obj := newEGenericTypeImpl()
	mockAdapter := &MockEAdapter{}
	mockAdapter.On("SetTarget", obj).Once()
	mockAdapter.On("NotifyChanged", mock.Anything).Once()
	obj.EAdapters().Add(mockAdapter)
	obj.SetELowerBound(newValue)
	mockAdapter.AssertExpectations(t)
}

func TestEGenericTypeETypeParameterGet(t *testing.T) {
	var newValue *MockETypeParameter = &MockETypeParameter{}
	obj := newEGenericTypeImpl()
	obj.SetETypeParameter(newValue)
	assert.Equal(t, newValue, obj.GetETypeParameter())
}

func TestEGenericTypeETypeParameterSet(t *testing.T) {
	var newValue *MockETypeParameter = &MockETypeParameter{}
	obj := newEGenericTypeImpl()
	mockAdapter := &MockEAdapter{}
	mockAdapter.On("SetTarget", obj).Once()
	mockAdapter.On("NotifyChanged", mock.Anything).Once()
	obj.EAdapters().Add(mockAdapter)
	obj.SetETypeParameter(newValue)
	mockAdapter.AssertExpectations(t)
}

func TestEGenericTypeEClassifierGet(t *testing.T) {
	var newValue *MockEClassifier = &MockEClassifier{}
	newValue.On("EIsProxy").Return(false)
	obj := newEGenericTypeImpl()
	obj.SetEClassifier(newValue)
	assert.Equal(t, newValue, obj.GetEClassifier())
}

func TestEGenericTypeEClassifierSet(t *testing.T) {
	var newValue *MockEClassifier = &MockEClassifier{}
	obj := newEGenericTypeImpl()
	mockAdapter := &MockEAdapter{}
	mockAdapter.On("SetTarget", obj).Once()
	mockAdapter.On("NotifyChanged", mock.Anything).Once()
	obj.EAdapters().Add(mockAdapter)
	obj.SetEClassifier(newValue)
	mockAdapter.AssertExpectations(t)
}

func TestEGenericTypeIsInstanceOperation(t *testing.T) {
	obj := newEGenericTypeImpl()
	assert.Panics(t, func() { obj.IsInstance(nil) })
}
