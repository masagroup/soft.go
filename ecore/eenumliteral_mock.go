// Code generated by soft.generator.go. DO NOT EDIT.

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
	"github.com/stretchr/testify/mock"
)

type MockEEnumLiteral struct {
	MockENamedElement
}

type MockEEnumLiteral_Expecter struct {
	MockENamedElement_Expecter
}

func (eEnumLiteral *MockEEnumLiteral) EXPECT() *MockEEnumLiteral_Expecter {
	e := &MockEEnumLiteral_Expecter{}
	e.Mock = &eEnumLiteral.Mock
	return e
}

// GetEEnum get the value of eEnum
func (eEnumLiteral *MockEEnumLiteral) GetEEnum() EEnum {
	ret := eEnumLiteral.Called()

	var r EEnum
	if rf, ok := ret.Get(0).(func() EEnum); ok {
		r = rf()
	} else {
		if ret.Get(0) != nil {
			r = ret.Get(0).(EEnum)
		}
	}

	return r
}

type MockEEnumLiteral_GetEEnum_Call struct {
	*mock.Call
}

func (e *MockEEnumLiteral_Expecter) GetEEnum() *MockEEnumLiteral_GetEEnum_Call {
	return &MockEEnumLiteral_GetEEnum_Call{Call: e.Mock.On("GetEEnum")}
}

func (c *MockEEnumLiteral_GetEEnum_Call) Run(run func()) *MockEEnumLiteral_GetEEnum_Call {
	c.Call.Run(func(mock.Arguments) {
		run()
	})
	return c
}

func (c *MockEEnumLiteral_GetEEnum_Call) Return(eEnum EEnum) *MockEEnumLiteral_GetEEnum_Call {
	c.Call.Return(eEnum)
	return c
}

// GetInstance get the value of instance
func (eEnumLiteral *MockEEnumLiteral) GetInstance() any {
	ret := eEnumLiteral.Called()

	var r any
	if rf, ok := ret.Get(0).(func() any); ok {
		r = rf()
	} else {
		if ret.Get(0) != nil {
			r = ret.Get(0).(any)
		}
	}

	return r
}

type MockEEnumLiteral_GetInstance_Call struct {
	*mock.Call
}

func (e *MockEEnumLiteral_Expecter) GetInstance() *MockEEnumLiteral_GetInstance_Call {
	return &MockEEnumLiteral_GetInstance_Call{Call: e.Mock.On("GetInstance")}
}

func (c *MockEEnumLiteral_GetInstance_Call) Run(run func()) *MockEEnumLiteral_GetInstance_Call {
	c.Call.Run(func(mock.Arguments) {
		run()
	})
	return c
}

func (c *MockEEnumLiteral_GetInstance_Call) Return(instance any) *MockEEnumLiteral_GetInstance_Call {
	c.Call.Return(instance)
	return c
}

// SetInstance provides mock implementation for setting the value of instance
func (eEnumLiteral *MockEEnumLiteral) SetInstance(newInstance any) {
	eEnumLiteral.Called(newInstance)
}

type MockEEnumLiteral_SetInstance_Call struct {
	*mock.Call
}

func (e *MockEEnumLiteral_Expecter) SetInstance(newInstance any) *MockEEnumLiteral_SetInstance_Call {
	return &MockEEnumLiteral_SetInstance_Call{Call: e.Mock.On("SetInstance", newInstance)}
}

func (c *MockEEnumLiteral_SetInstance_Call) Run(run func(newInstance any)) *MockEEnumLiteral_SetInstance_Call {
	c.Call.Run(func(args mock.Arguments) {
		run(args[0].(any))
	})
	return c
}

func (c *MockEEnumLiteral_SetInstance_Call) Return() *MockEEnumLiteral_SetInstance_Call {
	c.Call.Return()
	return c
}

// GetLiteral get the value of literal
func (eEnumLiteral *MockEEnumLiteral) GetLiteral() string {
	ret := eEnumLiteral.Called()

	var r string
	if rf, ok := ret.Get(0).(func() string); ok {
		r = rf()
	} else {
		if ret.Get(0) != nil {
			r = ret.Get(0).(string)
		}
	}

	return r
}

type MockEEnumLiteral_GetLiteral_Call struct {
	*mock.Call
}

func (e *MockEEnumLiteral_Expecter) GetLiteral() *MockEEnumLiteral_GetLiteral_Call {
	return &MockEEnumLiteral_GetLiteral_Call{Call: e.Mock.On("GetLiteral")}
}

func (c *MockEEnumLiteral_GetLiteral_Call) Run(run func()) *MockEEnumLiteral_GetLiteral_Call {
	c.Call.Run(func(mock.Arguments) {
		run()
	})
	return c
}

func (c *MockEEnumLiteral_GetLiteral_Call) Return(literal string) *MockEEnumLiteral_GetLiteral_Call {
	c.Call.Return(literal)
	return c
}

// SetLiteral provides mock implementation for setting the value of literal
func (eEnumLiteral *MockEEnumLiteral) SetLiteral(newLiteral string) {
	eEnumLiteral.Called(newLiteral)
}

type MockEEnumLiteral_SetLiteral_Call struct {
	*mock.Call
}

func (e *MockEEnumLiteral_Expecter) SetLiteral(newLiteral string) *MockEEnumLiteral_SetLiteral_Call {
	return &MockEEnumLiteral_SetLiteral_Call{Call: e.Mock.On("SetLiteral", newLiteral)}
}

func (c *MockEEnumLiteral_SetLiteral_Call) Run(run func(newLiteral string)) *MockEEnumLiteral_SetLiteral_Call {
	c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return c
}

func (c *MockEEnumLiteral_SetLiteral_Call) Return() *MockEEnumLiteral_SetLiteral_Call {
	c.Call.Return()
	return c
}

// GetValue get the value of value
func (eEnumLiteral *MockEEnumLiteral) GetValue() int {
	ret := eEnumLiteral.Called()

	var r int
	if rf, ok := ret.Get(0).(func() int); ok {
		r = rf()
	} else {
		if ret.Get(0) != nil {
			r = ret.Get(0).(int)
		}
	}

	return r
}

type MockEEnumLiteral_GetValue_Call struct {
	*mock.Call
}

func (e *MockEEnumLiteral_Expecter) GetValue() *MockEEnumLiteral_GetValue_Call {
	return &MockEEnumLiteral_GetValue_Call{Call: e.Mock.On("GetValue")}
}

func (c *MockEEnumLiteral_GetValue_Call) Run(run func()) *MockEEnumLiteral_GetValue_Call {
	c.Call.Run(func(mock.Arguments) {
		run()
	})
	return c
}

func (c *MockEEnumLiteral_GetValue_Call) Return(value int) *MockEEnumLiteral_GetValue_Call {
	c.Call.Return(value)
	return c
}

// SetValue provides mock implementation for setting the value of value
func (eEnumLiteral *MockEEnumLiteral) SetValue(newValue int) {
	eEnumLiteral.Called(newValue)
}

type MockEEnumLiteral_SetValue_Call struct {
	*mock.Call
}

func (e *MockEEnumLiteral_Expecter) SetValue(newValue int) *MockEEnumLiteral_SetValue_Call {
	return &MockEEnumLiteral_SetValue_Call{Call: e.Mock.On("SetValue", newValue)}
}

func (c *MockEEnumLiteral_SetValue_Call) Run(run func(newValue int)) *MockEEnumLiteral_SetValue_Call {
	c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int))
	})
	return c
}

func (c *MockEEnumLiteral_SetValue_Call) Return() *MockEEnumLiteral_SetValue_Call {
	c.Call.Return()
	return c
}

type mockConstructorTestingTNewMockEEnumLiteral interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockEEnumLiteral creates a new instance of MockEEnumLiteral. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockEEnumLiteral(t mockConstructorTestingTNewMockEEnumLiteral) *MockEEnumLiteral {
	mock := &MockEEnumLiteral{}
	mock.Mock.Test(t)
	t.Cleanup(func() { mock.AssertExpectations(t) })
	return mock
}
