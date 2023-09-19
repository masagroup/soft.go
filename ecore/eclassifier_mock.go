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

// MockEClassifier is an mock type for the EClassifier type
type MockEClassifier struct {
	mock.Mock
	MockEClassifier_Prototype
}

// MockEClassifier_Prototype is the mock implementation of all EClassifier methods ( inherited and declared )
type MockEClassifier_Prototype struct {
	mock *mock.Mock
	MockENamedElement_Prototype
	MockEClassifier_Prototype_Methods
}

func (_mp *MockEClassifier_Prototype) SetMock(mock *mock.Mock) {
	_mp.mock = mock
	_mp.MockENamedElement_Prototype.SetMock(mock)
	_mp.MockEClassifier_Prototype_Methods.SetMock(mock)
}

// MockEClassifier_Expecter is the expecter implementation for all EClassifier methods ( inherited and declared )
type MockEClassifier_Expecter struct {
	MockENamedElement_Expecter
	MockEClassifier_Expecter_Methods
}

func (_me *MockEClassifier_Expecter) SetMock(mock *mock.Mock) {
	_me.MockENamedElement_Expecter.SetMock(mock)
	_me.MockEClassifier_Expecter_Methods.SetMock(mock)
}

func (e *MockEClassifier_Prototype) EXPECT() *MockEClassifier_Expecter {
	expecter := &MockEClassifier_Expecter{}
	expecter.SetMock(e.mock)
	return expecter
}

// MockEClassifier_Prototype_Methods is the mock implementation of EClassifier declared methods
type MockEClassifier_Prototype_Methods struct {
	mock *mock.Mock
}

func (_mdp *MockEClassifier_Prototype_Methods) SetMock(mock *mock.Mock) {
	_mdp.mock = mock
}

// MockEClassifier_Expecter_Methods is the expecter implementation of EClassifier declared methods
type MockEClassifier_Expecter_Methods struct {
	mock *mock.Mock
}

func (_mde *MockEClassifier_Expecter_Methods) SetMock(mock *mock.Mock) {
	_mde.mock = mock
}

// GetClassifierID get the value of classifierID
func (e *MockEClassifier_Prototype_Methods) GetClassifierID() int {
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

type MockEClassifier_GetClassifierID_Call struct {
	*mock.Call
}

func (e *MockEClassifier_Expecter_Methods) GetClassifierID() *MockEClassifier_GetClassifierID_Call {
	return &MockEClassifier_GetClassifierID_Call{Call: e.mock.On("GetClassifierID")}
}

func (c *MockEClassifier_GetClassifierID_Call) Run(run func()) *MockEClassifier_GetClassifierID_Call {
	c.Call.Run(func(mock.Arguments) {
		run()
	})
	return c
}

func (c *MockEClassifier_GetClassifierID_Call) Return(classifierID int) *MockEClassifier_GetClassifierID_Call {
	c.Call.Return(classifierID)
	return c
}

// SetClassifierID provides mock implementation for setting the value of classifierID
func (e *MockEClassifier_Prototype_Methods) SetClassifierID(classifierID int) {
	e.mock.Called(classifierID)
}

type MockEClassifier_SetClassifierID_Call struct {
	*mock.Call
}

// SetClassifierID is a helper method to define mock.On call
// - classifierID int
func (e *MockEClassifier_Expecter_Methods) SetClassifierID(classifierID any) *MockEClassifier_SetClassifierID_Call {
	return &MockEClassifier_SetClassifierID_Call{Call: e.mock.On("SetClassifierID", classifierID)}
}

func (c *MockEClassifier_SetClassifierID_Call) Run(run func(classifierID int)) *MockEClassifier_SetClassifierID_Call {
	c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int))
	})
	return c
}

func (c *MockEClassifier_SetClassifierID_Call) Return() *MockEClassifier_SetClassifierID_Call {
	c.Call.Return()
	return c
}

// GetDefaultValue get the value of defaultValue
func (e *MockEClassifier_Prototype_Methods) GetDefaultValue() any {
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

type MockEClassifier_GetDefaultValue_Call struct {
	*mock.Call
}

func (e *MockEClassifier_Expecter_Methods) GetDefaultValue() *MockEClassifier_GetDefaultValue_Call {
	return &MockEClassifier_GetDefaultValue_Call{Call: e.mock.On("GetDefaultValue")}
}

func (c *MockEClassifier_GetDefaultValue_Call) Run(run func()) *MockEClassifier_GetDefaultValue_Call {
	c.Call.Run(func(mock.Arguments) {
		run()
	})
	return c
}

func (c *MockEClassifier_GetDefaultValue_Call) Return(defaultValue any) *MockEClassifier_GetDefaultValue_Call {
	c.Call.Return(defaultValue)
	return c
}

// GetEPackage get the value of ePackage
func (e *MockEClassifier_Prototype_Methods) GetEPackage() EPackage {
	ret := e.mock.Called()

	var r EPackage
	if rf, ok := ret.Get(0).(func() EPackage); ok {
		r = rf()
	} else {
		if ret.Get(0) != nil {
			r = ret.Get(0).(EPackage)
		}
	}

	return r
}

type MockEClassifier_GetEPackage_Call struct {
	*mock.Call
}

func (e *MockEClassifier_Expecter_Methods) GetEPackage() *MockEClassifier_GetEPackage_Call {
	return &MockEClassifier_GetEPackage_Call{Call: e.mock.On("GetEPackage")}
}

func (c *MockEClassifier_GetEPackage_Call) Run(run func()) *MockEClassifier_GetEPackage_Call {
	c.Call.Run(func(mock.Arguments) {
		run()
	})
	return c
}

func (c *MockEClassifier_GetEPackage_Call) Return(ePackage EPackage) *MockEClassifier_GetEPackage_Call {
	c.Call.Return(ePackage)
	return c
}

// GetInstanceClass get the value of instanceClass
func (e *MockEClassifier_Prototype_Methods) GetInstanceClass() reflect.Type {
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

type MockEClassifier_GetInstanceClass_Call struct {
	*mock.Call
}

func (e *MockEClassifier_Expecter_Methods) GetInstanceClass() *MockEClassifier_GetInstanceClass_Call {
	return &MockEClassifier_GetInstanceClass_Call{Call: e.mock.On("GetInstanceClass")}
}

func (c *MockEClassifier_GetInstanceClass_Call) Run(run func()) *MockEClassifier_GetInstanceClass_Call {
	c.Call.Run(func(mock.Arguments) {
		run()
	})
	return c
}

func (c *MockEClassifier_GetInstanceClass_Call) Return(instanceClass reflect.Type) *MockEClassifier_GetInstanceClass_Call {
	c.Call.Return(instanceClass)
	return c
}

// SetInstanceClass provides mock implementation for setting the value of instanceClass
func (e *MockEClassifier_Prototype_Methods) SetInstanceClass(instanceClass reflect.Type) {
	e.mock.Called(instanceClass)
}

type MockEClassifier_SetInstanceClass_Call struct {
	*mock.Call
}

// SetInstanceClass is a helper method to define mock.On call
// - instanceClass reflect.Type
func (e *MockEClassifier_Expecter_Methods) SetInstanceClass(instanceClass any) *MockEClassifier_SetInstanceClass_Call {
	return &MockEClassifier_SetInstanceClass_Call{Call: e.mock.On("SetInstanceClass", instanceClass)}
}

func (c *MockEClassifier_SetInstanceClass_Call) Run(run func(instanceClass reflect.Type)) *MockEClassifier_SetInstanceClass_Call {
	c.Call.Run(func(args mock.Arguments) {
		run(args[0].(reflect.Type))
	})
	return c
}

func (c *MockEClassifier_SetInstanceClass_Call) Return() *MockEClassifier_SetInstanceClass_Call {
	c.Call.Return()
	return c
}

// GetInstanceClassName get the value of instanceClassName
func (e *MockEClassifier_Prototype_Methods) GetInstanceClassName() string {
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

type MockEClassifier_GetInstanceClassName_Call struct {
	*mock.Call
}

func (e *MockEClassifier_Expecter_Methods) GetInstanceClassName() *MockEClassifier_GetInstanceClassName_Call {
	return &MockEClassifier_GetInstanceClassName_Call{Call: e.mock.On("GetInstanceClassName")}
}

func (c *MockEClassifier_GetInstanceClassName_Call) Run(run func()) *MockEClassifier_GetInstanceClassName_Call {
	c.Call.Run(func(mock.Arguments) {
		run()
	})
	return c
}

func (c *MockEClassifier_GetInstanceClassName_Call) Return(instanceClassName string) *MockEClassifier_GetInstanceClassName_Call {
	c.Call.Return(instanceClassName)
	return c
}

// SetInstanceClassName provides mock implementation for setting the value of instanceClassName
func (e *MockEClassifier_Prototype_Methods) SetInstanceClassName(instanceClassName string) {
	e.mock.Called(instanceClassName)
}

type MockEClassifier_SetInstanceClassName_Call struct {
	*mock.Call
}

// SetInstanceClassName is a helper method to define mock.On call
// - instanceClassName string
func (e *MockEClassifier_Expecter_Methods) SetInstanceClassName(instanceClassName any) *MockEClassifier_SetInstanceClassName_Call {
	return &MockEClassifier_SetInstanceClassName_Call{Call: e.mock.On("SetInstanceClassName", instanceClassName)}
}

func (c *MockEClassifier_SetInstanceClassName_Call) Run(run func(instanceClassName string)) *MockEClassifier_SetInstanceClassName_Call {
	c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return c
}

func (c *MockEClassifier_SetInstanceClassName_Call) Return() *MockEClassifier_SetInstanceClassName_Call {
	c.Call.Return()
	return c
}

// GetInstanceTypeName get the value of instanceTypeName
func (e *MockEClassifier_Prototype_Methods) GetInstanceTypeName() string {
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

type MockEClassifier_GetInstanceTypeName_Call struct {
	*mock.Call
}

func (e *MockEClassifier_Expecter_Methods) GetInstanceTypeName() *MockEClassifier_GetInstanceTypeName_Call {
	return &MockEClassifier_GetInstanceTypeName_Call{Call: e.mock.On("GetInstanceTypeName")}
}

func (c *MockEClassifier_GetInstanceTypeName_Call) Run(run func()) *MockEClassifier_GetInstanceTypeName_Call {
	c.Call.Run(func(mock.Arguments) {
		run()
	})
	return c
}

func (c *MockEClassifier_GetInstanceTypeName_Call) Return(instanceTypeName string) *MockEClassifier_GetInstanceTypeName_Call {
	c.Call.Return(instanceTypeName)
	return c
}

// SetInstanceTypeName provides mock implementation for setting the value of instanceTypeName
func (e *MockEClassifier_Prototype_Methods) SetInstanceTypeName(instanceTypeName string) {
	e.mock.Called(instanceTypeName)
}

type MockEClassifier_SetInstanceTypeName_Call struct {
	*mock.Call
}

// SetInstanceTypeName is a helper method to define mock.On call
// - instanceTypeName string
func (e *MockEClassifier_Expecter_Methods) SetInstanceTypeName(instanceTypeName any) *MockEClassifier_SetInstanceTypeName_Call {
	return &MockEClassifier_SetInstanceTypeName_Call{Call: e.mock.On("SetInstanceTypeName", instanceTypeName)}
}

func (c *MockEClassifier_SetInstanceTypeName_Call) Run(run func(instanceTypeName string)) *MockEClassifier_SetInstanceTypeName_Call {
	c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return c
}

func (c *MockEClassifier_SetInstanceTypeName_Call) Return() *MockEClassifier_SetInstanceTypeName_Call {
	c.Call.Return()
	return c
}

// IsInstance provides mock implementation
func (e *MockEClassifier_Prototype_Methods) IsInstance(object any) bool {
	ret := e.mock.Called(object)

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

type MockEClassifier_IsInstance_Call struct {
	*mock.Call
}

// IsInstance is a helper method to define mock.On call
// - object any
func (e *MockEClassifier_Expecter_Methods) IsInstance(object any) *MockEClassifier_IsInstance_Call {
	return &MockEClassifier_IsInstance_Call{Call: e.mock.On("IsInstance", object)}
}

func (c *MockEClassifier_IsInstance_Call) Run(run func(any)) *MockEClassifier_IsInstance_Call {
	c.Call.Run(func(_args mock.Arguments) {
		run(_args[0])
	})
	return c
}

func (c *MockEClassifier_IsInstance_Call) Return(_a0 bool) *MockEClassifier_IsInstance_Call {
	c.Call.Return(_a0)
	return c
}

type mockConstructorTestingTNewMockEClassifier interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockEClassifier creates a new instance of MockEClassifier. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockEClassifier(t mockConstructorTestingTNewMockEClassifier) *MockEClassifier {
	mock := &MockEClassifier{}
	mock.SetMock(&mock.Mock)
	mock.Mock.Test(t)
	t.Cleanup(func() { mock.AssertExpectations(t) })
	return mock
}
