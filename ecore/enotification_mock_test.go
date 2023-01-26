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
	n := NewMockENotification(t)
	m := NewMockRun(t)
	n.EXPECT().GetEventType().Return(SET).Run(func() { m.Run() }).Once()
	n.EXPECT().GetEventType().Call.Return(func() EventType { return SET }).Once()
	assert.Equal(t, SET, n.GetEventType())
	assert.Equal(t, SET, n.GetEventType())
}

func TestMockENotificationGetNotifier(t *testing.T) {
	n := NewMockENotification(t)
	no := NewMockENotifier(t)
	m := NewMockRun(t)
	n.EXPECT().GetNotifier().Return(no).Run(func() { m.Run() }).Once()
	n.EXPECT().GetNotifier().Call.Return(func() ENotifier { return no }).Once()
	assert.Equal(t, no, n.GetNotifier())
	assert.Equal(t, no, n.GetNotifier())
}

func TestMockENotificationGetFeature(t *testing.T) {
	n := NewMockENotification(t)
	f := &MockEStructuralFeature{}
	m := NewMockRun(t)
	n.EXPECT().GetFeature().Run(func() { m.Run() }).Return(f).Once()
	n.EXPECT().GetFeature().Once().Return(func() EStructuralFeature { return f })
	assert.Equal(t, f, n.GetFeature())
	assert.Equal(t, f, n.GetFeature())
}

func TestMockENotificationGetFeatureID(t *testing.T) {
	n := NewMockENotification(t)
	m := NewMockRun(t)
	n.EXPECT().GetFeatureID().Run(func() { m.Run() }).Return(0).Once()
	n.EXPECT().GetFeatureID().Once().Return(func() int { return 1 })
	assert.Equal(t, 0, n.GetFeatureID())
	assert.Equal(t, 1, n.GetFeatureID())
}

func TestMockENotificationGetOldValue(t *testing.T) {
	n := NewMockENotification(t)
	v := &MockEObject{}
	m := NewMockRun(t)
	n.EXPECT().GetOldValue().Run(func() { m.Run() }).Return(v).Once()
	n.EXPECT().GetOldValue().Once().Return(func() any { return v })
	assert.Equal(t, v, n.GetOldValue())
	assert.Equal(t, v, n.GetOldValue())
	mock.AssertExpectationsForObjects(t, n, v)
}

func TestMockENotificationGetNewValue(t *testing.T) {
	n := NewMockENotification(t)
	v := &MockEObject{}
	m := NewMockRun(t)
	n.EXPECT().GetNewValue().Run(func() { m.Run() }).Return(v).Once()
	n.EXPECT().GetNewValue().Once().Return(func() any { return v })
	assert.Equal(t, v, n.GetNewValue())
	assert.Equal(t, v, n.GetNewValue())
	mock.AssertExpectationsForObjects(t, n, v)
}

func TestMockENotificationGetPosition(t *testing.T) {
	n := NewMockENotification(t)
	m := NewMockRun(t)
	n.EXPECT().GetPosition().Run(func() { m.Run() }).Return(0).Once()
	n.EXPECT().GetPosition().Once().Return(func() int { return 1 })
	assert.Equal(t, 0, n.GetPosition())
	assert.Equal(t, 1, n.GetPosition())
}

func TestMockENotificationMerge(t *testing.T) {
	n := NewMockENotification(t)
	no := NewMockENotification(t)
	m := NewMockRun(t, no)
	n.EXPECT().Merge(no).Run(func(_a0 ENotification) { m.Run(_a0) }).Return(false).Once()
	n.EXPECT().Merge(no).Once().Return(func(ENotification) bool { return true })
	assert.False(t, n.Merge(no))
	assert.True(t, n.Merge(no))
}
