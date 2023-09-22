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

import "github.com/stretchr/testify/mock"

// MockEGenericType is an mock type for the EGenericType type
type MockEGenericType struct {
	mock.Mock
	MockEGenericType_Prototype
}

// MockEGenericType_Prototype is the mock implementation of all EGenericType methods ( inherited and declared )
type MockEGenericType_Prototype struct {
	mock *mock.Mock
	MockEObjectInternal_Prototype
	MockEGenericType_Prototype_Methods
}

func (_mp *MockEGenericType_Prototype) SetMock(mock *mock.Mock) {
	_mp.mock = mock
	_mp.MockEObjectInternal_Prototype.SetMock(mock)
	_mp.MockEGenericType_Prototype_Methods.SetMock(mock)
}

// MockEGenericType_Expecter is the expecter implementation for all EGenericType methods ( inherited and declared )
type MockEGenericType_Expecter struct {
	MockEObjectInternal_Expecter
	MockEGenericType_Expecter_Methods
}

func (_me *MockEGenericType_Expecter) SetMock(mock *mock.Mock) {
	_me.MockEObjectInternal_Expecter.SetMock(mock)
	_me.MockEGenericType_Expecter_Methods.SetMock(mock)
}

func (e *MockEGenericType_Prototype) EXPECT() *MockEGenericType_Expecter {
	expecter := &MockEGenericType_Expecter{}
	expecter.SetMock(e.mock)
	return expecter
}

// MockEGenericType_Prototype_Methods is the mock implementation of EGenericType declared methods
type MockEGenericType_Prototype_Methods struct {
	mock *mock.Mock
}

func (_mdp *MockEGenericType_Prototype_Methods) SetMock(mock *mock.Mock) {
	_mdp.mock = mock
}

// MockEGenericType_Expecter_Methods is the expecter implementation of EGenericType declared methods
type MockEGenericType_Expecter_Methods struct {
	mock *mock.Mock
}

func (_mde *MockEGenericType_Expecter_Methods) SetMock(mock *mock.Mock) {
	_mde.mock = mock
}

// GetEClassifier get the value of eClassifier
func (e *MockEGenericType_Prototype_Methods) GetEClassifier() EClassifier {
	ret := e.mock.Called()

	var res EClassifier
	if rf, ok := ret.Get(0).(func() EClassifier); ok {
		res = rf()
	} else {
		if ret.Get(0) != nil {
			res = ret.Get(0).(EClassifier)
		}
	}

	return res
}

type MockEGenericType_GetEClassifier_Call struct {
	*mock.Call
}

func (e *MockEGenericType_Expecter_Methods) GetEClassifier() *MockEGenericType_GetEClassifier_Call {
	return &MockEGenericType_GetEClassifier_Call{Call: e.mock.On("GetEClassifier")}
}

func (c *MockEGenericType_GetEClassifier_Call) Run(run func()) *MockEGenericType_GetEClassifier_Call {
	c.Call.Run(func(mock.Arguments) {
		run()
	})
	return c
}

func (c *MockEGenericType_GetEClassifier_Call) Return(eClassifier EClassifier) *MockEGenericType_GetEClassifier_Call {
	c.Call.Return(eClassifier)
	return c
}

// SetEClassifier provides mock implementation for setting the value of eClassifier
func (e *MockEGenericType_Prototype_Methods) SetEClassifier(eClassifier EClassifier) {
	e.mock.Called(eClassifier)
}

type MockEGenericType_SetEClassifier_Call struct {
	*mock.Call
}

// SetEClassifier is a helper method to define mock.On call
// - eClassifier EClassifier
func (e *MockEGenericType_Expecter_Methods) SetEClassifier(eClassifier any) *MockEGenericType_SetEClassifier_Call {
	return &MockEGenericType_SetEClassifier_Call{Call: e.mock.On("SetEClassifier", eClassifier)}
}

func (c *MockEGenericType_SetEClassifier_Call) Run(run func(eClassifier EClassifier)) *MockEGenericType_SetEClassifier_Call {
	c.Call.Run(func(args mock.Arguments) {
		run(args[0].(EClassifier))
	})
	return c
}

func (c *MockEGenericType_SetEClassifier_Call) Return() *MockEGenericType_SetEClassifier_Call {
	c.Call.Return()
	return c
}

// GetELowerBound get the value of eLowerBound
func (e *MockEGenericType_Prototype_Methods) GetELowerBound() EGenericType {
	ret := e.mock.Called()

	var res EGenericType
	if rf, ok := ret.Get(0).(func() EGenericType); ok {
		res = rf()
	} else {
		if ret.Get(0) != nil {
			res = ret.Get(0).(EGenericType)
		}
	}

	return res
}

type MockEGenericType_GetELowerBound_Call struct {
	*mock.Call
}

func (e *MockEGenericType_Expecter_Methods) GetELowerBound() *MockEGenericType_GetELowerBound_Call {
	return &MockEGenericType_GetELowerBound_Call{Call: e.mock.On("GetELowerBound")}
}

func (c *MockEGenericType_GetELowerBound_Call) Run(run func()) *MockEGenericType_GetELowerBound_Call {
	c.Call.Run(func(mock.Arguments) {
		run()
	})
	return c
}

func (c *MockEGenericType_GetELowerBound_Call) Return(eLowerBound EGenericType) *MockEGenericType_GetELowerBound_Call {
	c.Call.Return(eLowerBound)
	return c
}

// SetELowerBound provides mock implementation for setting the value of eLowerBound
func (e *MockEGenericType_Prototype_Methods) SetELowerBound(eLowerBound EGenericType) {
	e.mock.Called(eLowerBound)
}

type MockEGenericType_SetELowerBound_Call struct {
	*mock.Call
}

// SetELowerBound is a helper method to define mock.On call
// - eLowerBound EGenericType
func (e *MockEGenericType_Expecter_Methods) SetELowerBound(eLowerBound any) *MockEGenericType_SetELowerBound_Call {
	return &MockEGenericType_SetELowerBound_Call{Call: e.mock.On("SetELowerBound", eLowerBound)}
}

func (c *MockEGenericType_SetELowerBound_Call) Run(run func(eLowerBound EGenericType)) *MockEGenericType_SetELowerBound_Call {
	c.Call.Run(func(args mock.Arguments) {
		run(args[0].(EGenericType))
	})
	return c
}

func (c *MockEGenericType_SetELowerBound_Call) Return() *MockEGenericType_SetELowerBound_Call {
	c.Call.Return()
	return c
}

// GetERawType get the value of eRawType
func (e *MockEGenericType_Prototype_Methods) GetERawType() EClassifier {
	ret := e.mock.Called()

	var res EClassifier
	if rf, ok := ret.Get(0).(func() EClassifier); ok {
		res = rf()
	} else {
		if ret.Get(0) != nil {
			res = ret.Get(0).(EClassifier)
		}
	}

	return res
}

type MockEGenericType_GetERawType_Call struct {
	*mock.Call
}

func (e *MockEGenericType_Expecter_Methods) GetERawType() *MockEGenericType_GetERawType_Call {
	return &MockEGenericType_GetERawType_Call{Call: e.mock.On("GetERawType")}
}

func (c *MockEGenericType_GetERawType_Call) Run(run func()) *MockEGenericType_GetERawType_Call {
	c.Call.Run(func(mock.Arguments) {
		run()
	})
	return c
}

func (c *MockEGenericType_GetERawType_Call) Return(eRawType EClassifier) *MockEGenericType_GetERawType_Call {
	c.Call.Return(eRawType)
	return c
}

// GetETypeArguments get the value of eTypeArguments
func (e *MockEGenericType_Prototype_Methods) GetETypeArguments() EList {
	ret := e.mock.Called()

	var res EList
	if rf, ok := ret.Get(0).(func() EList); ok {
		res = rf()
	} else {
		if ret.Get(0) != nil {
			res = ret.Get(0).(EList)
		}
	}

	return res
}

type MockEGenericType_GetETypeArguments_Call struct {
	*mock.Call
}

func (e *MockEGenericType_Expecter_Methods) GetETypeArguments() *MockEGenericType_GetETypeArguments_Call {
	return &MockEGenericType_GetETypeArguments_Call{Call: e.mock.On("GetETypeArguments")}
}

func (c *MockEGenericType_GetETypeArguments_Call) Run(run func()) *MockEGenericType_GetETypeArguments_Call {
	c.Call.Run(func(mock.Arguments) {
		run()
	})
	return c
}

func (c *MockEGenericType_GetETypeArguments_Call) Return(eTypeArguments EList) *MockEGenericType_GetETypeArguments_Call {
	c.Call.Return(eTypeArguments)
	return c
}

// GetETypeParameter get the value of eTypeParameter
func (e *MockEGenericType_Prototype_Methods) GetETypeParameter() ETypeParameter {
	ret := e.mock.Called()

	var res ETypeParameter
	if rf, ok := ret.Get(0).(func() ETypeParameter); ok {
		res = rf()
	} else {
		if ret.Get(0) != nil {
			res = ret.Get(0).(ETypeParameter)
		}
	}

	return res
}

type MockEGenericType_GetETypeParameter_Call struct {
	*mock.Call
}

func (e *MockEGenericType_Expecter_Methods) GetETypeParameter() *MockEGenericType_GetETypeParameter_Call {
	return &MockEGenericType_GetETypeParameter_Call{Call: e.mock.On("GetETypeParameter")}
}

func (c *MockEGenericType_GetETypeParameter_Call) Run(run func()) *MockEGenericType_GetETypeParameter_Call {
	c.Call.Run(func(mock.Arguments) {
		run()
	})
	return c
}

func (c *MockEGenericType_GetETypeParameter_Call) Return(eTypeParameter ETypeParameter) *MockEGenericType_GetETypeParameter_Call {
	c.Call.Return(eTypeParameter)
	return c
}

// SetETypeParameter provides mock implementation for setting the value of eTypeParameter
func (e *MockEGenericType_Prototype_Methods) SetETypeParameter(eTypeParameter ETypeParameter) {
	e.mock.Called(eTypeParameter)
}

type MockEGenericType_SetETypeParameter_Call struct {
	*mock.Call
}

// SetETypeParameter is a helper method to define mock.On call
// - eTypeParameter ETypeParameter
func (e *MockEGenericType_Expecter_Methods) SetETypeParameter(eTypeParameter any) *MockEGenericType_SetETypeParameter_Call {
	return &MockEGenericType_SetETypeParameter_Call{Call: e.mock.On("SetETypeParameter", eTypeParameter)}
}

func (c *MockEGenericType_SetETypeParameter_Call) Run(run func(eTypeParameter ETypeParameter)) *MockEGenericType_SetETypeParameter_Call {
	c.Call.Run(func(args mock.Arguments) {
		run(args[0].(ETypeParameter))
	})
	return c
}

func (c *MockEGenericType_SetETypeParameter_Call) Return() *MockEGenericType_SetETypeParameter_Call {
	c.Call.Return()
	return c
}

// GetEUpperBound get the value of eUpperBound
func (e *MockEGenericType_Prototype_Methods) GetEUpperBound() EGenericType {
	ret := e.mock.Called()

	var res EGenericType
	if rf, ok := ret.Get(0).(func() EGenericType); ok {
		res = rf()
	} else {
		if ret.Get(0) != nil {
			res = ret.Get(0).(EGenericType)
		}
	}

	return res
}

type MockEGenericType_GetEUpperBound_Call struct {
	*mock.Call
}

func (e *MockEGenericType_Expecter_Methods) GetEUpperBound() *MockEGenericType_GetEUpperBound_Call {
	return &MockEGenericType_GetEUpperBound_Call{Call: e.mock.On("GetEUpperBound")}
}

func (c *MockEGenericType_GetEUpperBound_Call) Run(run func()) *MockEGenericType_GetEUpperBound_Call {
	c.Call.Run(func(mock.Arguments) {
		run()
	})
	return c
}

func (c *MockEGenericType_GetEUpperBound_Call) Return(eUpperBound EGenericType) *MockEGenericType_GetEUpperBound_Call {
	c.Call.Return(eUpperBound)
	return c
}

// SetEUpperBound provides mock implementation for setting the value of eUpperBound
func (e *MockEGenericType_Prototype_Methods) SetEUpperBound(eUpperBound EGenericType) {
	e.mock.Called(eUpperBound)
}

type MockEGenericType_SetEUpperBound_Call struct {
	*mock.Call
}

// SetEUpperBound is a helper method to define mock.On call
// - eUpperBound EGenericType
func (e *MockEGenericType_Expecter_Methods) SetEUpperBound(eUpperBound any) *MockEGenericType_SetEUpperBound_Call {
	return &MockEGenericType_SetEUpperBound_Call{Call: e.mock.On("SetEUpperBound", eUpperBound)}
}

func (c *MockEGenericType_SetEUpperBound_Call) Run(run func(eUpperBound EGenericType)) *MockEGenericType_SetEUpperBound_Call {
	c.Call.Run(func(args mock.Arguments) {
		run(args[0].(EGenericType))
	})
	return c
}

func (c *MockEGenericType_SetEUpperBound_Call) Return() *MockEGenericType_SetEUpperBound_Call {
	c.Call.Return()
	return c
}

// IsInstance provides mock implementation
func (e *MockEGenericType_Prototype_Methods) IsInstance(object any) bool {
	ret := e.mock.Called(object)

	var res bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		res = rf()
	} else {
		if ret.Get(0) != nil {
			res = ret.Get(0).(bool)
		}
	}

	return res
}

type MockEGenericType_IsInstance_Call struct {
	*mock.Call
}

// IsInstance is a helper method to define mock.On call
// - object any
func (e *MockEGenericType_Expecter_Methods) IsInstance(object any) *MockEGenericType_IsInstance_Call {
	return &MockEGenericType_IsInstance_Call{Call: e.mock.On("IsInstance", object)}
}

func (c *MockEGenericType_IsInstance_Call) Run(run func(any)) *MockEGenericType_IsInstance_Call {
	c.Call.Run(func(_args mock.Arguments) {
		run(_args[0])
	})
	return c
}

func (c *MockEGenericType_IsInstance_Call) Return(_a0 bool) *MockEGenericType_IsInstance_Call {
	c.Call.Return(_a0)
	return c
}

type mockConstructorTestingTNewMockEGenericType interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockEGenericType creates a new instance of MockEGenericType. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockEGenericType(t mockConstructorTestingTNewMockEGenericType) *MockEGenericType {
	mock := &MockEGenericType{}
	mock.SetMock(&mock.Mock)
	mock.Mock.Test(t)
	t.Cleanup(func() { mock.AssertExpectations(t) })
	return mock
}
