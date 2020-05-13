// *****************************************************************************
//
// This file is part of a MASA library or program.
// Refer to the included end-user license agreement for restrictions.
//
// Copyright (c) 2020 MASA Group
//
// *****************************************************************************

package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEcoreUtilsConvertToString(t *testing.T) {
	mockObject := &MockEObject{}
	mockDataType := &MockEDataType{}
	mockPackage := &MockEPackage{}
	mockFactory := &MockEFactory{}
	mockDataType.On("GetEPackage").Once().Return(mockPackage)
	mockPackage.On("GetEFactoryInstance").Once().Return(mockFactory)
	mockFactory.On("ConvertToString", mockDataType, mockObject).Once().Return("test")
	assert.Equal(t, "test", ConvertToString(mockDataType, mockObject))
}

func TestEcoreUtilsCreateFromString(t *testing.T) {
	mockObject := &MockEObject{}
	mockDataType := &MockEDataType{}
	mockPackage := &MockEPackage{}
	mockFactory := &MockEFactory{}
	mockDataType.On("GetEPackage").Once().Return(mockPackage)
	mockPackage.On("GetEFactoryInstance").Once().Return(mockFactory)
	mockFactory.On("CreateFromString", mockDataType, "test").Once().Return(mockObject)
	assert.Equal(t, mockObject, CreateFromString(mockDataType, "test"))
}
