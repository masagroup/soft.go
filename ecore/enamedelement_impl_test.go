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

func discardENamedElement() {
	_ = assert.Equal
	_ = mock.Anything
	_ = testing.Coverage

}

func TestENamedElementNameGet(t *testing.T) {
	var newValue string = "Test String"
	obj := newENamedElementImpl()
	obj.SetName(newValue)
	assert.Equal(t, newValue, obj.GetName())
}

func TestENamedElementNameSet(t *testing.T) {
	var newValue string = "Test String"
	obj := newENamedElementImpl()
	mockAdapter := &MockEAdapter{}
	mockAdapter.On("SetTarget", obj).Once()
	mockAdapter.On("NotifyChanged", mock.Anything).Once()
	obj.EAdapters().Add(mockAdapter)
	obj.SetName(newValue)
	mockAdapter.AssertExpectations(t)
}
