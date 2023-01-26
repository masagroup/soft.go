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

type mockENotificationChainRun struct {
	mock.Mock
}

func (m *mockENotificationChainRun) Run(args ...any) {
	m.Called(args...)
}

type mockConstructorTestingTmockENotificationChainRun interface {
	mock.TestingT
	Cleanup(func())
}

// newMockENotificationRun creates a new instance of MockEList. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockENotificationChainRun(t mockConstructorTestingTmockENotificationChainRun, args ...any) *mockENotificationChainRun {
	mock := &mockENotificationChainRun{}
	mock.Test(t)
	mock.On("Run", args...).Once()
	t.Cleanup(func() { mock.AssertExpectations(t) })
	return mock
}

func TestMockENotificationChainAdd(t *testing.T) {
	nc := NewMockENotificationChain(t)
	n := NewMockENotification(t)
	m := newMockENotificationChainRun(t, n)
	nc.EXPECT().Add(n).Return(true).Run(func(n ENotification) { m.Run(n) }).Once()
	nc.EXPECT().Add(n).Call.Return(func(ENotification) bool { return false }).Once()
	assert.True(t, nc.Add(n))
	assert.False(t, nc.Add(n))
}

func TestMockENotificationChainDispatch(t *testing.T) {
	nc := NewMockENotificationChain(t)
	m := newMockENotificationChainRun(t)
	nc.EXPECT().Dispatch().Return().Run(func() { m.Run() }).Once()
	nc.Dispatch()
}
