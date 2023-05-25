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
	io "io"

	mock "github.com/stretchr/testify/mock"
)

// MockECodec is mock type for the ECodec type
type MockECodec struct {
	mock.Mock
}

type MockECodec_Expecter struct {
	mock *mock.Mock
}

func (_m *MockECodec) EXPECT() *MockECodec_Expecter {
	return &MockECodec_Expecter{mock: &_m.Mock}
}

// NewDecoder provides a mock function with given fields: resource, r, options
func (_m *MockECodec) NewDecoder(resource EResource, r io.Reader, options map[string]interface{}) EDecoder {
	ret := _m.Called(resource, r, options)

	var r0 EDecoder
	if rf, ok := ret.Get(0).(func(EResource, io.Reader, map[string]interface{}) EDecoder); ok {
		r0 = rf(resource, r, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(EDecoder)
		}
	}

	return r0
}

// MockECodec_NewDecoder_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'NewDecoder'
type MockECodec_NewDecoder_Call struct {
	*mock.Call
}

// NewDecoder is a helper method to define mock.On call
//   - resource EResource
//   - r io.Reader
//   - options map[string]interface{}
func (_e *MockECodec_Expecter) NewDecoder(resource interface{}, r interface{}, options interface{}) *MockECodec_NewDecoder_Call {
	return &MockECodec_NewDecoder_Call{Call: _e.mock.On("NewDecoder", resource, r, options)}
}

func (_c *MockECodec_NewDecoder_Call) Run(run func(resource EResource, r io.Reader, options map[string]interface{})) *MockECodec_NewDecoder_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(EResource), args[1].(io.Reader), args[2].(map[string]interface{}))
	})
	return _c
}

func (_c *MockECodec_NewDecoder_Call) Return(_a0 EDecoder) *MockECodec_NewDecoder_Call {
	_c.Call.Return(_a0)
	return _c
}

// NewEncoder provides a mock function with given fields: resource, w, options
func (_m *MockECodec) NewEncoder(resource EResource, w io.Writer, options map[string]interface{}) EEncoder {
	ret := _m.Called(resource, w, options)

	var r0 EEncoder
	if rf, ok := ret.Get(0).(func(EResource, io.Writer, map[string]interface{}) EEncoder); ok {
		r0 = rf(resource, w, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(EEncoder)
		}
	}

	return r0
}

// MockECodec_NewEncoder_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'NewEncoder'
type MockECodec_NewEncoder_Call struct {
	*mock.Call
}

// NewEncoder is a helper method to define mock.On call
//   - resource EResource
//   - w io.Writer
//   - options map[string]interface{}
func (_e *MockECodec_Expecter) NewEncoder(resource interface{}, w interface{}, options interface{}) *MockECodec_NewEncoder_Call {
	return &MockECodec_NewEncoder_Call{Call: _e.mock.On("NewEncoder", resource, w, options)}
}

func (_c *MockECodec_NewEncoder_Call) Run(run func(resource EResource, w io.Writer, options map[string]interface{})) *MockECodec_NewEncoder_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(EResource), args[1].(io.Writer), args[2].(map[string]interface{}))
	})
	return _c
}

func (_c *MockECodec_NewEncoder_Call) Return(_a0 EEncoder) *MockECodec_NewEncoder_Call {
	_c.Call.Return(_a0)
	return _c
}

type mockConstructorTestingTNewMockECodec interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockECodec creates a new instance of MockECodec. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockECodec(t mockConstructorTestingTNewMockECodec) *MockECodec {
	mock := &MockECodec{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
