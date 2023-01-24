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

type mockENotifierRun struct {
	mock.Mock
}

func (m *mockENotifierRun) Run(args ...any) {
	m.Called(args...)
}

type mockConstructorTestingTMockENotifierRun interface {
	mock.TestingT
	Cleanup(func())
}

// newMockENotifierRun creates a new instance of MockEList. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockENotifierRun(t mockConstructorTestingTMockENotifierRun, args ...any) *mockENotifierRun {
	mock := &mockENotifierRun{}
	mock.Test(t)
	mock.On("Run", args...).Once()
	t.Cleanup(func() { mock.AssertExpectations(t) })
	return mock
}

func TestMockENotifierEAdapters(t *testing.T) {
	n := NewMockENotifier(t)
	a := NewMockEList(t)
	m := newMockENotifierRun(t)
	n.EXPECT().EAdapters().Run(func() { m.Run() }).Return(a).Once()
	n.EXPECT().EAdapters().Once().Return(func() EList { return a })
	assert.Equal(t, a, n.EAdapters())
	assert.Equal(t, a, n.EAdapters())
}

func TestMockENotifierEDeliver(t *testing.T) {
	n := NewMockENotifier(t)
	m := newMockENotifierRun(t)
	n.EXPECT().EDeliver().Run(func() { m.Run() }).Return(false).Once()
	n.EXPECT().EDeliver().Once().Return(func() bool { return true }).Once()
	assert.False(t, n.EDeliver())
	assert.True(t, n.EDeliver())
}

func TestMockENotifierESetDeliver(t *testing.T) {
	n := NewMockENotifier(t)
	m := newMockENotifierRun(t, true)
	n.EXPECT().ESetDeliver(true).Run(func(_a0 bool) { m.Run(_a0) }).Once()
	n.ESetDeliver(true)
}

func TestMockENotifierENotify(t *testing.T) {
	n := NewMockENotifier(t)
	notif := NewMockENotification(t)
	m := newMockENotifierRun(t, notif)
	n.EXPECT().ENotify(notif).Run(func(_a0 ENotification) { m.Run(_a0) }).Once()
	n.ENotify(notif)
}
