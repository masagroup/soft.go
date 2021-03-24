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

func TestMockENotificationGetEventType(t *testing.T) {
	n := &MockENotification{}
	n.On("GetEventType").Return(SET).Once()
	n.On("GetEventType").Return(func() EventType {
		return SET
	}).Once()
	assert.Equal(t, SET, n.GetEventType())
	assert.Equal(t, SET, n.GetEventType())
	mock.AssertExpectationsForObjects(t, n)
}

func TestMockENotificationGetNotifier(t *testing.T) {
	n := &MockENotification{}
	no := &MockENotifier{}
	n.On("GetNotifier").Return(no).Once()
	n.On("GetNotifier").Return(func() ENotifier {
		return no
	}).Once()
	assert.Equal(t, no, n.GetNotifier())
	assert.Equal(t, no, n.GetNotifier())
	mock.AssertExpectationsForObjects(t, n)
}

func TestMockENotificationGetFeature(t *testing.T) {
	n := &MockENotification{}
	f := &MockEStructuralFeature{}
	n.On("GetFeature").Return(f).Once()
	n.On("GetFeature").Return(func() EStructuralFeature {
		return f
	}).Once()
	assert.Equal(t, f, n.GetFeature())
	assert.Equal(t, f, n.GetFeature())
	mock.AssertExpectationsForObjects(t, n)
}

func TestMockENotificationGetFeatureID(t *testing.T) {
	n := &MockENotification{}
	n.On("GetFeatureID").Return(0).Once()
	n.On("GetFeatureID").Return(func() int {
		return 1
	}).Once()
	assert.Equal(t, 0, n.GetFeatureID())
	assert.Equal(t, 1, n.GetFeatureID())
	mock.AssertExpectationsForObjects(t, n)
}

func TestMockENotificationGetOldValue(t *testing.T) {
	n := &MockENotification{}
	v := &MockEObject{}
	n.On("GetOldValue").Return(v).Once()
	n.On("GetOldValue").Return(func() interface{} {
		return v
	}).Once()
	assert.Equal(t, v, n.GetOldValue())
	assert.Equal(t, v, n.GetOldValue())
	mock.AssertExpectationsForObjects(t, n, v)
}

func TestMockENotificationGetNewValue(t *testing.T) {
	n := &MockENotification{}
	v := &MockEObject{}
	n.On("GetNewValue").Return(v).Once()
	n.On("GetNewValue").Return(func() interface{} {
		return v
	}).Once()
	assert.Equal(t, v, n.GetNewValue())
	assert.Equal(t, v, n.GetNewValue())
	mock.AssertExpectationsForObjects(t, n, v)
}

func TestMockENotificationGetPosition(t *testing.T) {
	n := &MockENotification{}
	n.On("GetPosition").Return(0).Once()
	n.On("GetPosition").Return(func() int {
		return 1
	}).Once()
	assert.Equal(t, 0, n.GetPosition())
	assert.Equal(t, 1, n.GetPosition())
	mock.AssertExpectationsForObjects(t, n)
}

func TestMockENotificationMerge(t *testing.T) {
	n := &MockENotification{}
	no := &MockENotification{}
	n.On("Merge", no).Return(false).Once()
	n.On("Merge", no).Return(func(ENotification) bool {
		return true
	}).Once()
	assert.False(t, n.Merge(no))
	assert.True(t, n.Merge(no))
	mock.AssertExpectationsForObjects(t, n)
}
