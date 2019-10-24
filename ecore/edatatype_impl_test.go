// *****************************************************************************
//
// This file is part of a MASA library or program.
// Refer to the included end-user license agreement for restrictions.
//
// Copyright (c) 2019 MASA Group
//
// *****************************************************************************

// *****************************************************************************
//
// Warning: This file was generated by soft.generator.go Generator
//
// *****************************************************************************

package ecore

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func discardEDataType() {
	_ = assert.Equal
	_ = mock.Anything
	_ = testing.Coverage

	_ = time.Now()
}

func TestEDataTypeSerializableGet(t *testing.T) {
	obj := newEDataTypeImpl()
	obj.SetSerializable(true)
	assert.Equal(t, true, obj.IsSerializable())
}

func TestEDataTypeSerializableSet(t *testing.T) {
	obj := newEDataTypeImpl()
	mockAdapter := &MockEAdapter{}
	mockAdapter.On("SetTarget", obj).Once()
	mockAdapter.On("NotifyChanged", mock.Anything).Once()
	obj.EAdapters().Add(mockAdapter)
	obj.SetSerializable(true)
	mockAdapter.AssertExpectations(t)
}

func TestEDataTypeSerializableEGet(t *testing.T) {
	obj := newEDataTypeImpl()
	{
		assert.Equal(t, obj.IsSerializable(), obj.EGetFromID(EDATA_TYPE__SERIALIZABLE, false, false))
	}
}

func TestEDataTypeSerializableEIsSet(t *testing.T) {
	obj := newEDataTypeImpl()
	{
		_ = obj
	}
}

func TestEDataTypeSerializableEUnset(t *testing.T) {
	obj := newEDataTypeImpl()
	{
		_ = obj
	}
}

func TestEDataTypeSerializableESet(t *testing.T) {
	obj := newEDataTypeImpl()
	{
		_ = obj
	}
}