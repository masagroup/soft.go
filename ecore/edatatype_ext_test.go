// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEDataTypeExtSetDefaultValue(t *testing.T) {
	d := newEDataTypeExt()
	mockValue := NewMockEObject(t)
	mockAdapter := NewMockEAdapter(t)
	mockAdapter.EXPECT().SetTarget(d).Once()
	d.EAdapters().Add(mockAdapter)
	mock.AssertExpectationsForObjects(t, mockAdapter)

	mockAdapter.EXPECT().NotifyChanged(mock.MatchedBy(func(notification ENotification) bool {
		return notification.GetEventType() == SET && notification.GetNewValue() == mockValue
	})).Once()
	d.SetDefaultValue(mockValue)
	mock.AssertExpectationsForObjects(t, mockAdapter)
}

func TestEDataTypeExtGetDefaultValue(t *testing.T) {
	d := newEDataTypeExt()
	mockValue := NewMockEDataType(t)
	d.defaultValue = mockValue
	assert.Equal(t, mockValue, d.GetDefaultValue())
}

func TestEDataTypeExtGetDefaultValueFloat64(t *testing.T) {
	d := newEDataTypeExt()
	d.SetInstanceTypeName("float64")
	assert.Equal(t, float64(0), d.GetDefaultValue())
}

func TestEDataTypeExtGetDefaultValueFloat32(t *testing.T) {
	d := newEDataTypeExt()
	d.SetInstanceTypeName("float32")
	assert.Equal(t, float32(0), d.GetDefaultValue())
}

func TestEDataTypeExtGetDefaultValueInt(t *testing.T) {
	d := newEDataTypeExt()
	d.SetInstanceTypeName("int")
	assert.Equal(t, 0, d.GetDefaultValue())
}

func TestEDataTypeExtGetDefaultValueInt64(t *testing.T) {
	d := newEDataTypeExt()
	d.SetInstanceTypeName("int64")
	assert.Equal(t, int64(0), d.GetDefaultValue())
}

func TestEDataTypeExtGetDefaultValueInt32(t *testing.T) {
	d := newEDataTypeExt()
	d.SetInstanceTypeName("int32")
	assert.Equal(t, int32(0), d.GetDefaultValue())
}

func TestEDataTypeExtGetDefaultValueInt16(t *testing.T) {
	d := newEDataTypeExt()
	d.SetInstanceTypeName("int16")
	assert.Equal(t, int16(0), d.GetDefaultValue())
}

func TestEDataTypeExtGetDefaultValueInt8(t *testing.T) {
	d := newEDataTypeExt()
	d.SetInstanceTypeName("int8")
	assert.Equal(t, int8(0), d.GetDefaultValue())
}

func TestEDataTypeExtGetDefaultValueBool(t *testing.T) {
	d := newEDataTypeExt()
	d.SetInstanceTypeName("bool")
	assert.Equal(t, false, d.GetDefaultValue())
}

func TestEDataTypeExtGetDefaultValueString(t *testing.T) {
	d := newEDataTypeExt()
	d.SetInstanceTypeName("string")
	assert.Equal(t, "", d.GetDefaultValue())
}

func TestEDataTypeExtGetInstanceTypeName(t *testing.T) {
	mockEDataType := NewMockEDataType(t)
	mockEDataType.EXPECT().GetEAnnotation("http://net.masagroup/soft/2019/GenGo").Return(nil).Once()
	mockEDataType.EXPECT().GetInstanceTypeName().Return("type").Once()
	assert.Equal(t, "type", getInstanceTypeName(mockEDataType))

	mockEAnnotation := NewMockEAnnotation(t)
	mockEMap := NewMockEMap(t)
	mockEDataType.EXPECT().GetEAnnotation("http://net.masagroup/soft/2019/GenGo").Return(mockEAnnotation).Once()
	mockEAnnotation.EXPECT().GetDetails().Return(mockEMap).Once()
	mockEMap.EXPECT().GetValue("instanceTypeName").Return("type")
	assert.Equal(t, "type", getInstanceTypeName(mockEDataType))
}
