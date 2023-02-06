// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import mock "github.com/stretchr/testify/mock"

type MockRun struct {
	mock.Mock
}

func (m *MockRun) Run(args ...any) {
	m.Called(args...)
}

type mockConstructorTestingTMockRun interface {
	mock.TestingT
	Cleanup(func())
}

// newMockEIteratorRun creates a new instance of NewMockRun. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockRun(t mockConstructorTestingTMockRun, args ...any) *MockRun {
	mock := &MockRun{}
	mock.Test(t)
	mock.On("Run", args...).Once()
	t.Cleanup(func() { mock.AssertExpectations(t) })
	return mock
}
