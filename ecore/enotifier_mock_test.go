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

func TestMockENotifierEAdapters(t *testing.T) {
	n := &MockENotifier{}
	a := &MockEList{}
	n.On("EAdapters").Return(a).Once()
	n.On("EAdapters").Return(func() EList {
		return a
	}).Once()
	assert.Equal(t, a, n.EAdapters())
	assert.Equal(t, a, n.EAdapters())
	mock.AssertExpectationsForObjects(t, n, a)
}

func TestMockENotifierEDeliver(t *testing.T) {
	n := &MockENotifier{}
	n.On("EDeliver").Return(false).Once()
	n.On("EDeliver").Return(func() bool {
		return true
	}).Once()
	assert.False(t, n.EDeliver())
	assert.True(t, n.EDeliver())
	mock.AssertExpectationsForObjects(t, n)
}

func TestMockENotifierESetDeliver(t *testing.T) {
	n := &MockENotifier{}
	n.On("ESetDeliver", true).Once()
	n.ESetDeliver(true)
	mock.AssertExpectationsForObjects(t, n)
}

func TestMockENotifierENotify(t *testing.T) {
	n := &MockENotifier{}
	notif := &MockENotification{}
	n.On("ENotify", notif).Once()
	n.ENotify(notif)
	mock.AssertExpectationsForObjects(t, n, notif)
}

func TestMockENotifierENotificationRequired(t *testing.T) {
	n := &MockENotifier{}
	n.On("ENotificationRequired").Return(false).Once()
	n.On("ENotificationRequired").Return(func() bool {
		return true
	}).Once()
	assert.False(t, n.ENotificationRequired())
	assert.True(t, n.ENotificationRequired())
	mock.AssertExpectationsForObjects(t, n)
}
