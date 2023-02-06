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
	mock "github.com/stretchr/testify/mock"
)

type MockEContentAdapter struct {
	mock.Mock
	EContentAdapter
}

type MockEContentAdapter_Expecter struct {
	*mock.Mock
}

func (eAdapter *MockEContentAdapter) EXPECT() *MockEContentAdapter_Expecter {
	return &MockEContentAdapter_Expecter{Mock: &eAdapter.Mock}
}

func (adapter *MockEContentAdapter) NotifyChanged(notification ENotification) {
	adapter.Called(notification)
	adapter.EContentAdapter.NotifyChanged(notification)
}

type MockEClass_NotifyChanged_Call struct {
	*mock.Call
}

func (e *MockEContentAdapter_Expecter) NotifyChanged(notification any) *MockEClass_NotifyChanged_Call {
	return &MockEClass_NotifyChanged_Call{Call: e.Mock.On("NotifyChanged", notification)}
}

type mockConstructorTestingTNewMockEContentAdapter interface {
	mock.TestingT
	Cleanup(func())
}

func NewMockEContentAdapter(t mockConstructorTestingTNewMockEContentAdapter) *MockEContentAdapter {
	mock := &MockEContentAdapter{}
	mock.SetInterfaces(mock)
	mock.Mock.Test(t)
	t.Cleanup(func() { mock.AssertExpectations(t) })
	return mock
}
