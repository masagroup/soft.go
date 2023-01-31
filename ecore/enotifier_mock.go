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

// MockENotifier is an mock type for the ENotifier type
type MockENotifier struct {
	mock.Mock
	MockENotifier_Prototype
}

type MockENotifier_Prototype struct {
	MockENotifier_Declared_Prototype
}

func (_m *MockENotifier_Prototype) EXPECT() *MockENotifier_Expecter {
	expecter := &MockENotifier_Expecter{}
	expecter.SetMock(_m.mock)
	return expecter
}

type MockENotifier_Expecter struct {
	MockENotifier_Declared_Expecter
}

type MockENotifier_Declared_Prototype struct {
	mock *mock.Mock
}

func (_mdp *MockENotifier_Declared_Prototype) SetMock(mock *mock.Mock) {
	_mdp.mock = mock
}

type MockENotifier_Declared_Expecter struct {
	mock *mock.Mock
}

func (_mde *MockENotifier_Declared_Expecter) SetMock(mock *mock.Mock) {
	_mde.mock = mock
}

// EAdapters provides a mock function with given fields:
func (_m *MockENotifier_Prototype) EAdapters() EList {
	ret := _m.mock.Called()

	var r0 EList
	if rf, ok := ret.Get(0).(func() EList); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(EList)
		}
	}

	return r0
}

// MockENotifier_EAdapters_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EAdapters'
type MockENotifier_EAdapters_Call struct {
	*mock.Call
}

// EAdapters is a helper method to define mock.On call
func (_e *MockENotifier_Expecter) EAdapters() *MockENotifier_EAdapters_Call {
	return &MockENotifier_EAdapters_Call{Call: _e.mock.On("EAdapters")}
}

func (_c *MockENotifier_EAdapters_Call) Run(run func()) *MockENotifier_EAdapters_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockENotifier_EAdapters_Call) Return(_a0 EList) *MockENotifier_EAdapters_Call {
	_c.Call.Return(_a0)
	return _c
}

// EDeliver provides a mock function with given fields:
func (_m *MockENotifier_Prototype) EDeliver() bool {
	ret := _m.mock.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockENotifier_EDeliver_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EDeliver'
type MockENotifier_EDeliver_Call struct {
	*mock.Call
}

// EDeliver is a helper method to define mock.On call
func (_e *MockENotifier_Expecter) EDeliver() *MockENotifier_EDeliver_Call {
	return &MockENotifier_EDeliver_Call{Call: _e.mock.On("EDeliver")}
}

func (_c *MockENotifier_EDeliver_Call) Run(run func()) *MockENotifier_EDeliver_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockENotifier_EDeliver_Call) Return(_a0 bool) *MockENotifier_EDeliver_Call {
	_c.Call.Return(_a0)
	return _c
}

// ENotify provides a mock function with given fields: _a0
func (_m *MockENotifier_Prototype) ENotify(_a0 ENotification) {
	_m.mock.Called(_a0)
}

// MockENotifier_ENotify_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ENotify'
type MockENotifier_ENotify_Call struct {
	*mock.Call
}

// ENotify is a helper method to define mock.On call
//   - _a0 ENotification
func (_e *MockENotifier_Expecter) ENotify(_a0 interface{}) *MockENotifier_ENotify_Call {
	return &MockENotifier_ENotify_Call{Call: _e.mock.On("ENotify", _a0)}
}

func (_c *MockENotifier_ENotify_Call) Run(run func(_a0 ENotification)) *MockENotifier_ENotify_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(ENotification))
	})
	return _c
}

func (_c *MockENotifier_ENotify_Call) Return() *MockENotifier_ENotify_Call {
	_c.Call.Return()
	return _c
}

// ESetDeliver provides a mock function with given fields: _a0
func (_m *MockENotifier_Prototype) ESetDeliver(_a0 bool) {
	_m.mock.Called(_a0)
}

// MockENotifier_ESetDeliver_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ESetDeliver'
type MockENotifier_ESetDeliver_Call struct {
	*mock.Call
}

// ESetDeliver is a helper method to define mock.On call
//   - _a0 bool
func (_e *MockENotifier_Expecter) ESetDeliver(_a0 interface{}) *MockENotifier_ESetDeliver_Call {
	return &MockENotifier_ESetDeliver_Call{Call: _e.mock.On("ESetDeliver", _a0)}
}

func (_c *MockENotifier_ESetDeliver_Call) Run(run func(_a0 bool)) *MockENotifier_ESetDeliver_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(bool))
	})
	return _c
}

func (_c *MockENotifier_ESetDeliver_Call) Return() *MockENotifier_ESetDeliver_Call {
	_c.Call.Return()
	return _c
}

type mockConstructorTestingTNewMockENotifier interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockENotifier creates a new instance of MockENotifier_Prototype. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockENotifier(t mockConstructorTestingTNewMockENotifier) *MockENotifier {
	mock := &MockENotifier{}
	mock.SetMock(&mock.Mock)
	mock.Mock.Test(t)
	t.Cleanup(func() { mock.AssertExpectations(t) })
	return mock
}
