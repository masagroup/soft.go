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
	value := NewMockEDataType(t)
	d.defaultValue = value
	assert.Equal(t, value, d.GetDefaultValue())
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
