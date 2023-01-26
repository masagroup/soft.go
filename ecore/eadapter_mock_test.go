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

type mockMockEAdapterRun struct {
	mock.Mock
}

func (m *mockMockEAdapterRun) Run(args ...any) {
	m.Called(args...)
}

type mockConstructorTestingTMockMockEAdapterRun interface {
	mock.TestingT
	Cleanup(func())
}

// newMockEListRun creates a new instance of MockEList. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockMockEAdapterRun(t mockConstructorTestingTMockMockEAdapterRun, args ...any) *mockMockEAdapterRun {
	mock := &mockMockEAdapterRun{}
	mock.Test(t)
	mock.On("Run", args...).Once()
	t.Cleanup(func() { mock.AssertExpectations(t) })
	return mock
}

func TestMockEAdapter_GetTarget(t *testing.T) {
	mockNotifier := NewMockENotifier(t)
	mockAdapter := NewMockEAdapter(t)
	m := newMockMockEAdapterRun(t)
	mockAdapter.EXPECT().GetTarget().Return(mockNotifier).Run(func() { m.Run() }).Once()
	mockAdapter.EXPECT().GetTarget().Once().Return(func() ENotifier { return mockNotifier })
	assert.Equal(t, mockNotifier, mockAdapter.GetTarget())
	assert.Equal(t, mockNotifier, mockAdapter.GetTarget())
}

func TestMockEAdapter_NotifyChanged(t *testing.T) {
	mockNotification := NewMockENotification(t)
	mockAdapter := NewMockEAdapter(t)
	m := newMockMockEAdapterRun(t, mockNotification)
	mockAdapter.EXPECT().NotifyChanged(mockNotification).Return().Run(func(notification ENotification) { m.Run(notification) }).Once()
	mockAdapter.NotifyChanged(mockNotification)
}

func TestMockEAdapter_SetTarget(t *testing.T) {
	mockNotifier := NewMockENotifier(t)
	mockAdapter := NewMockEAdapter(t)
	m := newMockMockEAdapterRun(t, mockNotifier)
	mockAdapter.EXPECT().SetTarget(mockNotifier).Return().Run(func(_a0 ENotifier) { m.Run(_a0) }).Once()
	mockAdapter.SetTarget(mockNotifier)
}

func TestMockEAdapter_UnSetTarget(t *testing.T) {
	mockNotifier := NewMockENotifier(t)
	mockAdapter := NewMockEAdapter(t)
	m := newMockMockEAdapterRun(t, mockNotifier)
	mockAdapter.EXPECT().UnSetTarget(mockNotifier).Return().Run(func(_a0 ENotifier) { m.Run(_a0) }).Once()
	mockAdapter.UnSetTarget(mockNotifier)
}
