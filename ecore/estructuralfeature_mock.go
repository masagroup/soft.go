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
import "reflect"

// MockEStructuralFeature is an mock type for the EStructuralFeature type
type MockEStructuralFeature struct {
	mock.Mock
	MockEStructuralFeature_Prototype
}

// MockEStructuralFeature_Prototype is the mock implementation of all EStructuralFeature methods ( inherited and declared )
type MockEStructuralFeature_Prototype struct {
	mock *mock.Mock
	MockETypedElement_Prototype
	MockEStructuralFeature_Prototype_Methods
}

func (_mp *MockEStructuralFeature_Prototype) SetMock(mock *mock.Mock) {
	_mp.mock = mock
	_mp.MockETypedElement_Prototype.SetMock(mock)
	_mp.MockEStructuralFeature_Prototype_Methods.SetMock(mock)
}

// MockEStructuralFeature_Expecter is the expecter implementation for all EStructuralFeature methods ( inherited and declared )
type MockEStructuralFeature_Expecter struct {
	MockETypedElement_Expecter
	MockEStructuralFeature_Expecter_Methods
}

func (_me *MockEStructuralFeature_Expecter) SetMock(mock *mock.Mock) {
	_me.MockETypedElement_Expecter.SetMock(mock)
	_me.MockEStructuralFeature_Expecter_Methods.SetMock(mock)
}

func (e *MockEStructuralFeature_Prototype) EXPECT() *MockEStructuralFeature_Expecter {
	expecter := &MockEStructuralFeature_Expecter{}
	expecter.SetMock(e.mock)
	return expecter
}

// MockEStructuralFeature_Prototype_Methods is the mock implementation of EStructuralFeature declared methods
type MockEStructuralFeature_Prototype_Methods struct {
	mock *mock.Mock
}

func (_mdp *MockEStructuralFeature_Prototype_Methods) SetMock(mock *mock.Mock) {
	_mdp.mock = mock
}

// MockEStructuralFeature_Expecter_Methods is the expecter implementation of EStructuralFeature declared methods
type MockEStructuralFeature_Expecter_Methods struct {
	mock *mock.Mock
}

func (_mde *MockEStructuralFeature_Expecter_Methods) SetMock(mock *mock.Mock) {
	_mde.mock = mock
}

// IsChangeable get the value of isChangeable
func (e *MockEStructuralFeature_Prototype_Methods) IsChangeable() bool {
	ret := e.mock.Called()

	var r bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r = rf()
	} else {
		if ret.Get(0) != nil {
			r = ret.Get(0).(bool)
		}
	}

	return r
}

type MockEStructuralFeature_IsChangeable_Call struct {
	*mock.Call
}

func (e *MockEStructuralFeature_Expecter_Methods) IsChangeable() *MockEStructuralFeature_IsChangeable_Call {
	return &MockEStructuralFeature_IsChangeable_Call{Call: e.mock.On("IsChangeable")}
}

func (c *MockEStructuralFeature_IsChangeable_Call) Run(run func()) *MockEStructuralFeature_IsChangeable_Call {
	c.Call.Run(func(mock.Arguments) {
		run()
	})
	return c
}

func (c *MockEStructuralFeature_IsChangeable_Call) Return(isChangeable bool) *MockEStructuralFeature_IsChangeable_Call {
	c.Call.Return(isChangeable)
	return c
}

// SetChangeable provides mock implementation for setting the value of isChangeable
func (e *MockEStructuralFeature_Prototype_Methods) SetChangeable(isChangeable bool) {
	e.mock.Called(isChangeable)
}

type MockEStructuralFeature_SetChangeable_Call struct {
	*mock.Call
}

// SetChangeable is a helper method to define mock.On call
// - isChangeable bool
func (e *MockEStructuralFeature_Expecter_Methods) SetChangeable(isChangeable any) *MockEStructuralFeature_SetChangeable_Call {
	return &MockEStructuralFeature_SetChangeable_Call{Call: e.mock.On("SetChangeable", isChangeable)}
}

func (c *MockEStructuralFeature_SetChangeable_Call) Run(run func(isChangeable bool)) *MockEStructuralFeature_SetChangeable_Call {
	c.Call.Run(func(args mock.Arguments) {
		run(args[0].(bool))
	})
	return c
}

func (c *MockEStructuralFeature_SetChangeable_Call) Return() *MockEStructuralFeature_SetChangeable_Call {
	c.Call.Return()
	return c
}

// GetDefaultValue get the value of defaultValue
func (e *MockEStructuralFeature_Prototype_Methods) GetDefaultValue() any {
	ret := e.mock.Called()

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

type MockEStructuralFeature_GetDefaultValue_Call struct {
	*mock.Call
}

func (e *MockEStructuralFeature_Expecter_Methods) GetDefaultValue() *MockEStructuralFeature_GetDefaultValue_Call {
	return &MockEStructuralFeature_GetDefaultValue_Call{Call: e.mock.On("GetDefaultValue")}
}

func (c *MockEStructuralFeature_GetDefaultValue_Call) Run(run func()) *MockEStructuralFeature_GetDefaultValue_Call {
	c.Call.Run(func(mock.Arguments) {
		run()
	})
	return c
}

func (c *MockEStructuralFeature_GetDefaultValue_Call) Return(defaultValue any) *MockEStructuralFeature_GetDefaultValue_Call {
	c.Call.Return(defaultValue)
	return c
}

// SetDefaultValue provides mock implementation for setting the value of defaultValue
func (e *MockEStructuralFeature_Prototype_Methods) SetDefaultValue(defaultValue any) {
	e.mock.Called(defaultValue)
}

type MockEStructuralFeature_SetDefaultValue_Call struct {
	*mock.Call
}

// SetDefaultValue is a helper method to define mock.On call
// - defaultValue any
func (e *MockEStructuralFeature_Expecter_Methods) SetDefaultValue(defaultValue any) *MockEStructuralFeature_SetDefaultValue_Call {
	return &MockEStructuralFeature_SetDefaultValue_Call{Call: e.mock.On("SetDefaultValue", defaultValue)}
}

func (c *MockEStructuralFeature_SetDefaultValue_Call) Run(run func(defaultValue any)) *MockEStructuralFeature_SetDefaultValue_Call {
	c.Call.Run(func(args mock.Arguments) {
		run(args[0])
	})
	return c
}

func (c *MockEStructuralFeature_SetDefaultValue_Call) Return() *MockEStructuralFeature_SetDefaultValue_Call {
	c.Call.Return()
	return c
}

// GetDefaultValueLiteral get the value of defaultValueLiteral
func (e *MockEStructuralFeature_Prototype_Methods) GetDefaultValueLiteral() string {
	ret := e.mock.Called()

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

type MockEStructuralFeature_GetDefaultValueLiteral_Call struct {
	*mock.Call
}

func (e *MockEStructuralFeature_Expecter_Methods) GetDefaultValueLiteral() *MockEStructuralFeature_GetDefaultValueLiteral_Call {
	return &MockEStructuralFeature_GetDefaultValueLiteral_Call{Call: e.mock.On("GetDefaultValueLiteral")}
}

func (c *MockEStructuralFeature_GetDefaultValueLiteral_Call) Run(run func()) *MockEStructuralFeature_GetDefaultValueLiteral_Call {
	c.Call.Run(func(mock.Arguments) {
		run()
	})
	return c
}

func (c *MockEStructuralFeature_GetDefaultValueLiteral_Call) Return(defaultValueLiteral string) *MockEStructuralFeature_GetDefaultValueLiteral_Call {
	c.Call.Return(defaultValueLiteral)
	return c
}

// SetDefaultValueLiteral provides mock implementation for setting the value of defaultValueLiteral
func (e *MockEStructuralFeature_Prototype_Methods) SetDefaultValueLiteral(defaultValueLiteral string) {
	e.mock.Called(defaultValueLiteral)
}

type MockEStructuralFeature_SetDefaultValueLiteral_Call struct {
	*mock.Call
}

// SetDefaultValueLiteral is a helper method to define mock.On call
// - defaultValueLiteral string
func (e *MockEStructuralFeature_Expecter_Methods) SetDefaultValueLiteral(defaultValueLiteral any) *MockEStructuralFeature_SetDefaultValueLiteral_Call {
	return &MockEStructuralFeature_SetDefaultValueLiteral_Call{Call: e.mock.On("SetDefaultValueLiteral", defaultValueLiteral)}
}

func (c *MockEStructuralFeature_SetDefaultValueLiteral_Call) Run(run func(defaultValueLiteral string)) *MockEStructuralFeature_SetDefaultValueLiteral_Call {
	c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return c
}

func (c *MockEStructuralFeature_SetDefaultValueLiteral_Call) Return() *MockEStructuralFeature_SetDefaultValueLiteral_Call {
	c.Call.Return()
	return c
}

// IsDerived get the value of isDerived
func (e *MockEStructuralFeature_Prototype_Methods) IsDerived() bool {
	ret := e.mock.Called()

	var r bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r = rf()
	} else {
		if ret.Get(0) != nil {
			r = ret.Get(0).(bool)
		}
	}

	return r
}

type MockEStructuralFeature_IsDerived_Call struct {
	*mock.Call
}

func (e *MockEStructuralFeature_Expecter_Methods) IsDerived() *MockEStructuralFeature_IsDerived_Call {
	return &MockEStructuralFeature_IsDerived_Call{Call: e.mock.On("IsDerived")}
}

func (c *MockEStructuralFeature_IsDerived_Call) Run(run func()) *MockEStructuralFeature_IsDerived_Call {
	c.Call.Run(func(mock.Arguments) {
		run()
	})
	return c
}

func (c *MockEStructuralFeature_IsDerived_Call) Return(isDerived bool) *MockEStructuralFeature_IsDerived_Call {
	c.Call.Return(isDerived)
	return c
}

// SetDerived provides mock implementation for setting the value of isDerived
func (e *MockEStructuralFeature_Prototype_Methods) SetDerived(isDerived bool) {
	e.mock.Called(isDerived)
}

type MockEStructuralFeature_SetDerived_Call struct {
	*mock.Call
}

// SetDerived is a helper method to define mock.On call
// - isDerived bool
func (e *MockEStructuralFeature_Expecter_Methods) SetDerived(isDerived any) *MockEStructuralFeature_SetDerived_Call {
	return &MockEStructuralFeature_SetDerived_Call{Call: e.mock.On("SetDerived", isDerived)}
}

func (c *MockEStructuralFeature_SetDerived_Call) Run(run func(isDerived bool)) *MockEStructuralFeature_SetDerived_Call {
	c.Call.Run(func(args mock.Arguments) {
		run(args[0].(bool))
	})
	return c
}

func (c *MockEStructuralFeature_SetDerived_Call) Return() *MockEStructuralFeature_SetDerived_Call {
	c.Call.Return()
	return c
}

// GetEContainingClass get the value of eContainingClass
func (e *MockEStructuralFeature_Prototype_Methods) GetEContainingClass() EClass {
	ret := e.mock.Called()

	var r EClass
	if rf, ok := ret.Get(0).(func() EClass); ok {
		r = rf()
	} else {
		if ret.Get(0) != nil {
			r = ret.Get(0).(EClass)
		}
	}

	return r
}

type MockEStructuralFeature_GetEContainingClass_Call struct {
	*mock.Call
}

func (e *MockEStructuralFeature_Expecter_Methods) GetEContainingClass() *MockEStructuralFeature_GetEContainingClass_Call {
	return &MockEStructuralFeature_GetEContainingClass_Call{Call: e.mock.On("GetEContainingClass")}
}

func (c *MockEStructuralFeature_GetEContainingClass_Call) Run(run func()) *MockEStructuralFeature_GetEContainingClass_Call {
	c.Call.Run(func(mock.Arguments) {
		run()
	})
	return c
}

func (c *MockEStructuralFeature_GetEContainingClass_Call) Return(eContainingClass EClass) *MockEStructuralFeature_GetEContainingClass_Call {
	c.Call.Return(eContainingClass)
	return c
}

// GetFeatureID get the value of featureID
func (e *MockEStructuralFeature_Prototype_Methods) GetFeatureID() int {
	ret := e.mock.Called()

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

type MockEStructuralFeature_GetFeatureID_Call struct {
	*mock.Call
}

func (e *MockEStructuralFeature_Expecter_Methods) GetFeatureID() *MockEStructuralFeature_GetFeatureID_Call {
	return &MockEStructuralFeature_GetFeatureID_Call{Call: e.mock.On("GetFeatureID")}
}

func (c *MockEStructuralFeature_GetFeatureID_Call) Run(run func()) *MockEStructuralFeature_GetFeatureID_Call {
	c.Call.Run(func(mock.Arguments) {
		run()
	})
	return c
}

func (c *MockEStructuralFeature_GetFeatureID_Call) Return(featureID int) *MockEStructuralFeature_GetFeatureID_Call {
	c.Call.Return(featureID)
	return c
}

// SetFeatureID provides mock implementation for setting the value of featureID
func (e *MockEStructuralFeature_Prototype_Methods) SetFeatureID(featureID int) {
	e.mock.Called(featureID)
}

type MockEStructuralFeature_SetFeatureID_Call struct {
	*mock.Call
}

// SetFeatureID is a helper method to define mock.On call
// - featureID int
func (e *MockEStructuralFeature_Expecter_Methods) SetFeatureID(featureID any) *MockEStructuralFeature_SetFeatureID_Call {
	return &MockEStructuralFeature_SetFeatureID_Call{Call: e.mock.On("SetFeatureID", featureID)}
}

func (c *MockEStructuralFeature_SetFeatureID_Call) Run(run func(featureID int)) *MockEStructuralFeature_SetFeatureID_Call {
	c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int))
	})
	return c
}

func (c *MockEStructuralFeature_SetFeatureID_Call) Return() *MockEStructuralFeature_SetFeatureID_Call {
	c.Call.Return()
	return c
}

// IsTransient get the value of isTransient
func (e *MockEStructuralFeature_Prototype_Methods) IsTransient() bool {
	ret := e.mock.Called()

	var r bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r = rf()
	} else {
		if ret.Get(0) != nil {
			r = ret.Get(0).(bool)
		}
	}

	return r
}

type MockEStructuralFeature_IsTransient_Call struct {
	*mock.Call
}

func (e *MockEStructuralFeature_Expecter_Methods) IsTransient() *MockEStructuralFeature_IsTransient_Call {
	return &MockEStructuralFeature_IsTransient_Call{Call: e.mock.On("IsTransient")}
}

func (c *MockEStructuralFeature_IsTransient_Call) Run(run func()) *MockEStructuralFeature_IsTransient_Call {
	c.Call.Run(func(mock.Arguments) {
		run()
	})
	return c
}

func (c *MockEStructuralFeature_IsTransient_Call) Return(isTransient bool) *MockEStructuralFeature_IsTransient_Call {
	c.Call.Return(isTransient)
	return c
}

// SetTransient provides mock implementation for setting the value of isTransient
func (e *MockEStructuralFeature_Prototype_Methods) SetTransient(isTransient bool) {
	e.mock.Called(isTransient)
}

type MockEStructuralFeature_SetTransient_Call struct {
	*mock.Call
}

// SetTransient is a helper method to define mock.On call
// - isTransient bool
func (e *MockEStructuralFeature_Expecter_Methods) SetTransient(isTransient any) *MockEStructuralFeature_SetTransient_Call {
	return &MockEStructuralFeature_SetTransient_Call{Call: e.mock.On("SetTransient", isTransient)}
}

func (c *MockEStructuralFeature_SetTransient_Call) Run(run func(isTransient bool)) *MockEStructuralFeature_SetTransient_Call {
	c.Call.Run(func(args mock.Arguments) {
		run(args[0].(bool))
	})
	return c
}

func (c *MockEStructuralFeature_SetTransient_Call) Return() *MockEStructuralFeature_SetTransient_Call {
	c.Call.Return()
	return c
}

// IsUnsettable get the value of isUnsettable
func (e *MockEStructuralFeature_Prototype_Methods) IsUnsettable() bool {
	ret := e.mock.Called()

	var r bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r = rf()
	} else {
		if ret.Get(0) != nil {
			r = ret.Get(0).(bool)
		}
	}

	return r
}

type MockEStructuralFeature_IsUnsettable_Call struct {
	*mock.Call
}

func (e *MockEStructuralFeature_Expecter_Methods) IsUnsettable() *MockEStructuralFeature_IsUnsettable_Call {
	return &MockEStructuralFeature_IsUnsettable_Call{Call: e.mock.On("IsUnsettable")}
}

func (c *MockEStructuralFeature_IsUnsettable_Call) Run(run func()) *MockEStructuralFeature_IsUnsettable_Call {
	c.Call.Run(func(mock.Arguments) {
		run()
	})
	return c
}

func (c *MockEStructuralFeature_IsUnsettable_Call) Return(isUnsettable bool) *MockEStructuralFeature_IsUnsettable_Call {
	c.Call.Return(isUnsettable)
	return c
}

// SetUnsettable provides mock implementation for setting the value of isUnsettable
func (e *MockEStructuralFeature_Prototype_Methods) SetUnsettable(isUnsettable bool) {
	e.mock.Called(isUnsettable)
}

type MockEStructuralFeature_SetUnsettable_Call struct {
	*mock.Call
}

// SetUnsettable is a helper method to define mock.On call
// - isUnsettable bool
func (e *MockEStructuralFeature_Expecter_Methods) SetUnsettable(isUnsettable any) *MockEStructuralFeature_SetUnsettable_Call {
	return &MockEStructuralFeature_SetUnsettable_Call{Call: e.mock.On("SetUnsettable", isUnsettable)}
}

func (c *MockEStructuralFeature_SetUnsettable_Call) Run(run func(isUnsettable bool)) *MockEStructuralFeature_SetUnsettable_Call {
	c.Call.Run(func(args mock.Arguments) {
		run(args[0].(bool))
	})
	return c
}

func (c *MockEStructuralFeature_SetUnsettable_Call) Return() *MockEStructuralFeature_SetUnsettable_Call {
	c.Call.Return()
	return c
}

// IsVolatile get the value of isVolatile
func (e *MockEStructuralFeature_Prototype_Methods) IsVolatile() bool {
	ret := e.mock.Called()

	var r bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r = rf()
	} else {
		if ret.Get(0) != nil {
			r = ret.Get(0).(bool)
		}
	}

	return r
}

type MockEStructuralFeature_IsVolatile_Call struct {
	*mock.Call
}

func (e *MockEStructuralFeature_Expecter_Methods) IsVolatile() *MockEStructuralFeature_IsVolatile_Call {
	return &MockEStructuralFeature_IsVolatile_Call{Call: e.mock.On("IsVolatile")}
}

func (c *MockEStructuralFeature_IsVolatile_Call) Run(run func()) *MockEStructuralFeature_IsVolatile_Call {
	c.Call.Run(func(mock.Arguments) {
		run()
	})
	return c
}

func (c *MockEStructuralFeature_IsVolatile_Call) Return(isVolatile bool) *MockEStructuralFeature_IsVolatile_Call {
	c.Call.Return(isVolatile)
	return c
}

// SetVolatile provides mock implementation for setting the value of isVolatile
func (e *MockEStructuralFeature_Prototype_Methods) SetVolatile(isVolatile bool) {
	e.mock.Called(isVolatile)
}

type MockEStructuralFeature_SetVolatile_Call struct {
	*mock.Call
}

// SetVolatile is a helper method to define mock.On call
// - isVolatile bool
func (e *MockEStructuralFeature_Expecter_Methods) SetVolatile(isVolatile any) *MockEStructuralFeature_SetVolatile_Call {
	return &MockEStructuralFeature_SetVolatile_Call{Call: e.mock.On("SetVolatile", isVolatile)}
}

func (c *MockEStructuralFeature_SetVolatile_Call) Run(run func(isVolatile bool)) *MockEStructuralFeature_SetVolatile_Call {
	c.Call.Run(func(args mock.Arguments) {
		run(args[0].(bool))
	})
	return c
}

func (c *MockEStructuralFeature_SetVolatile_Call) Return() *MockEStructuralFeature_SetVolatile_Call {
	c.Call.Return()
	return c
}

// GetContainerClass provides mock implementation
func (e *MockEStructuralFeature_Prototype_Methods) GetContainerClass() reflect.Type {
	ret := e.mock.Called()

	var r reflect.Type
	if rf, ok := ret.Get(0).(func() reflect.Type); ok {
		r = rf()
	} else {
		if ret.Get(0) != nil {
			r = ret.Get(0).(reflect.Type)
		}
	}

	return r
}

type MockEStructuralFeature_GetContainerClass_Call struct {
	*mock.Call
}

// GetContainerClass is a helper method to define mock.On call
func (e *MockEStructuralFeature_Expecter_Methods) GetContainerClass() *MockEStructuralFeature_GetContainerClass_Call {
	return &MockEStructuralFeature_GetContainerClass_Call{Call: e.mock.On("GetContainerClass")}
}

func (c *MockEStructuralFeature_GetContainerClass_Call) Run(run func()) *MockEStructuralFeature_GetContainerClass_Call {
	c.Call.Run(func(_args mock.Arguments) {
		run()
	})
	return c
}

func (c *MockEStructuralFeature_GetContainerClass_Call) Return(_a0 reflect.Type) *MockEStructuralFeature_GetContainerClass_Call {
	c.Call.Return(_a0)
	return c
}

type mockConstructorTestingTNewMockEStructuralFeature interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockEStructuralFeature creates a new instance of MockEStructuralFeature. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockEStructuralFeature(t mockConstructorTestingTNewMockEStructuralFeature) *MockEStructuralFeature {
	mock := &MockEStructuralFeature{}
	mock.SetMock(&mock.Mock)
	mock.Mock.Test(t)
	t.Cleanup(func() { mock.AssertExpectations(t) })
	return mock
}
