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

type MockEObjectInternal struct {
	MockEObjectInternal_Prototype
	mock.Mock
}

// MockEObjectInternal_Prototype is an mock type for the EObjectInternal type
type MockEObjectInternal_Prototype struct {
	mock *mock.Mock
	MockEObject_Prototype
	MockEObjectInternal_Prototype_Methods
}

func (_mp *MockEObjectInternal_Prototype) SetMock(mock *mock.Mock) {
	_mp.mock = mock
	_mp.MockEObject_Prototype.SetMock(mock)
	_mp.MockEObjectInternal_Prototype_Methods.SetMock(mock)
}

type MockEObjectInternal_Expecter struct {
	MockEObject_Expecter
	MockEObjectInternal_Expecter_Methods
}

func (_me *MockEObjectInternal_Expecter) SetMock(mock *mock.Mock) {
	_me.MockEObject_Expecter.SetMock(mock)
	_me.MockEObjectInternal_Expecter_Methods.SetMock(mock)
}

func (_m *MockEObjectInternal_Prototype) EXPECT() *MockEObjectInternal_Expecter {
	e := &MockEObjectInternal_Expecter{}
	e.SetMock(_m.mock)
	return e
}

type MockEObjectInternal_Prototype_Methods struct {
	mock *mock.Mock
}

func (_mdp *MockEObjectInternal_Prototype_Methods) SetMock(mock *mock.Mock) {
	_mdp.mock = mock
}

type MockEObjectInternal_Expecter_Methods struct {
	mock *mock.Mock
}

func (_mde *MockEObjectInternal_Expecter_Methods) SetMock(mock *mock.Mock) {
	_mde.mock = mock
}

// EBasicInverseAdd provides a mock function with given fields: otherEnd, featureID, notifications
func (_m *MockEObjectInternal_Prototype_Methods) EBasicInverseAdd(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain {
	ret := _m.mock.Called(otherEnd, featureID, notifications)

	var r0 ENotificationChain
	if rf, ok := ret.Get(0).(func(EObject, int, ENotificationChain) ENotificationChain); ok {
		r0 = rf(otherEnd, featureID, notifications)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ENotificationChain)
		}
	}

	return r0
}

// MockEObjectInternal_EBasicInverseAdd_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EBasicInverseAdd'
type MockEObjectInternal_EBasicInverseAdd_Call struct {
	*mock.Call
}

// EBasicInverseAdd is a helper method to define mock.On call
//   - otherEnd EObject
//   - featureID int
//   - notifications ENotificationChain
func (_e *MockEObjectInternal_Expecter_Methods) EBasicInverseAdd(otherEnd interface{}, featureID interface{}, notifications interface{}) *MockEObjectInternal_EBasicInverseAdd_Call {
	return &MockEObjectInternal_EBasicInverseAdd_Call{Call: _e.mock.On("EBasicInverseAdd", otherEnd, featureID, notifications)}
}

func (_c *MockEObjectInternal_EBasicInverseAdd_Call) Run(run func(otherEnd EObject, featureID int, notifications ENotificationChain)) *MockEObjectInternal_EBasicInverseAdd_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(EObject), args[1].(int), args[2].(ENotificationChain))
	})
	return _c
}

func (_c *MockEObjectInternal_EBasicInverseAdd_Call) Return(_a0 ENotificationChain) *MockEObjectInternal_EBasicInverseAdd_Call {
	_c.Call.Return(_a0)
	return _c
}

// EBasicInverseRemove provides a mock function with given fields: otherEnd, featureID, notifications
func (_m *MockEObjectInternal_Prototype_Methods) EBasicInverseRemove(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain {
	ret := _m.mock.Called(otherEnd, featureID, notifications)

	var r0 ENotificationChain
	if rf, ok := ret.Get(0).(func(EObject, int, ENotificationChain) ENotificationChain); ok {
		r0 = rf(otherEnd, featureID, notifications)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ENotificationChain)
		}
	}

	return r0
}

// MockEObjectInternal_EBasicInverseRemove_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EBasicInverseRemove'
type MockEObjectInternal_EBasicInverseRemove_Call struct {
	*mock.Call
}

// EBasicInverseRemove is a helper method to define mock.On call
//   - otherEnd EObject
//   - featureID int
//   - notifications ENotificationChain
func (_e *MockEObjectInternal_Expecter_Methods) EBasicInverseRemove(otherEnd interface{}, featureID interface{}, notifications interface{}) *MockEObjectInternal_EBasicInverseRemove_Call {
	return &MockEObjectInternal_EBasicInverseRemove_Call{Call: _e.mock.On("EBasicInverseRemove", otherEnd, featureID, notifications)}
}

func (_c *MockEObjectInternal_EBasicInverseRemove_Call) Run(run func(otherEnd EObject, featureID int, notifications ENotificationChain)) *MockEObjectInternal_EBasicInverseRemove_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(EObject), args[1].(int), args[2].(ENotificationChain))
	})
	return _c
}

func (_c *MockEObjectInternal_EBasicInverseRemove_Call) Return(_a0 ENotificationChain) *MockEObjectInternal_EBasicInverseRemove_Call {
	_c.Call.Return(_a0)
	return _c
}

// EDerivedFeatureID provides a mock function with given fields: container, featureID
func (_m *MockEObjectInternal_Prototype_Methods) EDerivedFeatureID(container EObject, featureID int) int {
	ret := _m.mock.Called(container, featureID)

	var r0 int
	if rf, ok := ret.Get(0).(func(EObject, int) int); ok {
		r0 = rf(container, featureID)
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// MockEObjectInternal_EDerivedFeatureID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EDerivedFeatureID'
type MockEObjectInternal_EDerivedFeatureID_Call struct {
	*mock.Call
}

// EDerivedFeatureID is a helper method to define mock.On call
//   - container EObject
//   - featureID int
func (_e *MockEObjectInternal_Expecter_Methods) EDerivedFeatureID(container interface{}, featureID interface{}) *MockEObjectInternal_EDerivedFeatureID_Call {
	return &MockEObjectInternal_EDerivedFeatureID_Call{Call: _e.mock.On("EDerivedFeatureID", container, featureID)}
}

func (_c *MockEObjectInternal_EDerivedFeatureID_Call) Run(run func(container EObject, featureID int)) *MockEObjectInternal_EDerivedFeatureID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(EObject), args[1].(int))
	})
	return _c
}

func (_c *MockEObjectInternal_EDerivedFeatureID_Call) Return(_a0 int) *MockEObjectInternal_EDerivedFeatureID_Call {
	_c.Call.Return(_a0)
	return _c
}

// EDerivedOperationID provides a mock function with given fields: container, operationID
func (_m *MockEObjectInternal_Prototype_Methods) EDerivedOperationID(container EObject, operationID int) int {
	ret := _m.mock.Called(container, operationID)

	var r0 int
	if rf, ok := ret.Get(0).(func(EObject, int) int); ok {
		r0 = rf(container, operationID)
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// MockEObjectInternal_EDerivedOperationID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EDerivedOperationID'
type MockEObjectInternal_EDerivedOperationID_Call struct {
	*mock.Call
}

// EDerivedOperationID is a helper method to define mock.On call
//   - container EObject
//   - operationID int
func (_e *MockEObjectInternal_Expecter_Methods) EDerivedOperationID(container interface{}, operationID interface{}) *MockEObjectInternal_EDerivedOperationID_Call {
	return &MockEObjectInternal_EDerivedOperationID_Call{Call: _e.mock.On("EDerivedOperationID", container, operationID)}
}

func (_c *MockEObjectInternal_EDerivedOperationID_Call) Run(run func(container EObject, operationID int)) *MockEObjectInternal_EDerivedOperationID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(EObject), args[1].(int))
	})
	return _c
}

func (_c *MockEObjectInternal_EDerivedOperationID_Call) Return(_a0 int) *MockEObjectInternal_EDerivedOperationID_Call {
	_c.Call.Return(_a0)
	return _c
}

// EDynamicProperties provides a mock function with given fields:
func (_m *MockEObjectInternal_Prototype_Methods) EDynamicProperties() EDynamicProperties {
	ret := _m.mock.Called()

	var r0 EDynamicProperties
	if rf, ok := ret.Get(0).(func() EDynamicProperties); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(EDynamicProperties)
		}
	}

	return r0
}

// MockEObjectInternal_EDynamicProperties_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EDynamicProperties'
type MockEObjectInternal_EDynamicProperties_Call struct {
	*mock.Call
}

// EDynamicProperties is a helper method to define mock.On call
func (_e *MockEObjectInternal_Expecter_Methods) EDynamicProperties() *MockEObjectInternal_EDynamicProperties_Call {
	return &MockEObjectInternal_EDynamicProperties_Call{Call: _e.mock.On("EDynamicProperties")}
}

func (_c *MockEObjectInternal_EDynamicProperties_Call) Run(run func()) *MockEObjectInternal_EDynamicProperties_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockEObjectInternal_EDynamicProperties_Call) Return(_a0 EDynamicProperties) *MockEObjectInternal_EDynamicProperties_Call {
	_c.Call.Return(_a0)
	return _c
}

// EFeatureID provides a mock function with given fields: feature
func (_m *MockEObjectInternal_Prototype_Methods) EFeatureID(feature EStructuralFeature) int {
	ret := _m.mock.Called(feature)

	var r0 int
	if rf, ok := ret.Get(0).(func(EStructuralFeature) int); ok {
		r0 = rf(feature)
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// MockEObjectInternal_EFeatureID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EFeatureID'
type MockEObjectInternal_EFeatureID_Call struct {
	*mock.Call
}

// EFeatureID is a helper method to define mock.On call
//   - feature EStructuralFeature
func (_e *MockEObjectInternal_Expecter_Methods) EFeatureID(feature interface{}) *MockEObjectInternal_EFeatureID_Call {
	return &MockEObjectInternal_EFeatureID_Call{Call: _e.mock.On("EFeatureID", feature)}
}

func (_c *MockEObjectInternal_EFeatureID_Call) Run(run func(feature EStructuralFeature)) *MockEObjectInternal_EFeatureID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(EStructuralFeature))
	})
	return _c
}

func (_c *MockEObjectInternal_EFeatureID_Call) Return(_a0 int) *MockEObjectInternal_EFeatureID_Call {
	_c.Call.Return(_a0)
	return _c
}

// EGetFromID provides a mock function with given fields: featureID, resolve
func (_m *MockEObjectInternal_Prototype_Methods) EGetFromID(featureID int, resolve bool) interface{} {
	ret := _m.mock.Called(featureID, resolve)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(int, bool) interface{}); ok {
		r0 = rf(featureID, resolve)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0)
		}
	}

	return r0
}

// MockEObjectInternal_EGetFromID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EGetFromID'
type MockEObjectInternal_EGetFromID_Call struct {
	*mock.Call
}

// EGetFromID is a helper method to define mock.On call
//   - featureID int
//   - resolve bool
func (_e *MockEObjectInternal_Expecter_Methods) EGetFromID(featureID interface{}, resolve interface{}) *MockEObjectInternal_EGetFromID_Call {
	return &MockEObjectInternal_EGetFromID_Call{Call: _e.mock.On("EGetFromID", featureID, resolve)}
}

func (_c *MockEObjectInternal_EGetFromID_Call) Run(run func(featureID int, resolve bool)) *MockEObjectInternal_EGetFromID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int), args[1].(bool))
	})
	return _c
}

func (_c *MockEObjectInternal_EGetFromID_Call) Return(_a0 interface{}) *MockEObjectInternal_EGetFromID_Call {
	_c.Call.Return(_a0)
	return _c
}

// EInternalContainer provides a mock function with given fields:
func (_m *MockEObjectInternal_Prototype_Methods) EInternalContainer() EObject {
	ret := _m.mock.Called()

	var r0 EObject
	if rf, ok := ret.Get(0).(func() EObject); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(EObject)
		}
	}

	return r0
}

// MockEObjectInternal_EInternalContainer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EInternalContainer'
type MockEObjectInternal_EInternalContainer_Call struct {
	*mock.Call
}

// EInternalContainer is a helper method to define mock.On call
func (_e *MockEObjectInternal_Expecter_Methods) EInternalContainer() *MockEObjectInternal_EInternalContainer_Call {
	return &MockEObjectInternal_EInternalContainer_Call{Call: _e.mock.On("EInternalContainer")}
}

func (_c *MockEObjectInternal_EInternalContainer_Call) Run(run func()) *MockEObjectInternal_EInternalContainer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockEObjectInternal_EInternalContainer_Call) Return(_a0 EObject) *MockEObjectInternal_EInternalContainer_Call {
	_c.Call.Return(_a0)
	return _c
}

// EInternalContainerFeatureID provides a mock function with given fields:
func (_m *MockEObjectInternal_Prototype_Methods) EInternalContainerFeatureID() int {
	ret := _m.mock.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// MockEObjectInternal_EInternalContainerFeatureID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EInternalContainerFeatureID'
type MockEObjectInternal_EInternalContainerFeatureID_Call struct {
	*mock.Call
}

// EInternalContainerFeatureID is a helper method to define mock.On call
func (_e *MockEObjectInternal_Expecter_Methods) EInternalContainerFeatureID() *MockEObjectInternal_EInternalContainerFeatureID_Call {
	return &MockEObjectInternal_EInternalContainerFeatureID_Call{Call: _e.mock.On("EInternalContainerFeatureID")}
}

func (_c *MockEObjectInternal_EInternalContainerFeatureID_Call) Run(run func()) *MockEObjectInternal_EInternalContainerFeatureID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockEObjectInternal_EInternalContainerFeatureID_Call) Return(_a0 int) *MockEObjectInternal_EInternalContainerFeatureID_Call {
	_c.Call.Return(_a0)
	return _c
}

// EInternalResource provides a mock function with given fields:
func (_m *MockEObjectInternal_Prototype_Methods) EInternalResource() EResource {
	ret := _m.mock.Called()

	var r0 EResource
	if rf, ok := ret.Get(0).(func() EResource); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(EResource)
		}
	}

	return r0
}

// MockEObjectInternal_EInternalResource_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EInternalResource'
type MockEObjectInternal_EInternalResource_Call struct {
	*mock.Call
}

// EInternalResource is a helper method to define mock.On call
func (_e *MockEObjectInternal_Expecter_Methods) EInternalResource() *MockEObjectInternal_EInternalResource_Call {
	return &MockEObjectInternal_EInternalResource_Call{Call: _e.mock.On("EInternalResource")}
}

func (_c *MockEObjectInternal_EInternalResource_Call) Run(run func()) *MockEObjectInternal_EInternalResource_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockEObjectInternal_EInternalResource_Call) Return(_a0 EResource) *MockEObjectInternal_EInternalResource_Call {
	_c.Call.Return(_a0)
	return _c
}

// EInverseAdd provides a mock function with given fields: otherEnd, featureID, notifications
func (_m *MockEObjectInternal_Prototype_Methods) EInverseAdd(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain {
	ret := _m.mock.Called(otherEnd, featureID, notifications)

	var r0 ENotificationChain
	if rf, ok := ret.Get(0).(func(EObject, int, ENotificationChain) ENotificationChain); ok {
		r0 = rf(otherEnd, featureID, notifications)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ENotificationChain)
		}
	}

	return r0
}

// MockEObjectInternal_EInverseAdd_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EInverseAdd'
type MockEObjectInternal_EInverseAdd_Call struct {
	*mock.Call
}

// EInverseAdd is a helper method to define mock.On call
//   - otherEnd EObject
//   - featureID int
//   - notifications ENotificationChain
func (_e *MockEObjectInternal_Expecter_Methods) EInverseAdd(otherEnd interface{}, featureID interface{}, notifications interface{}) *MockEObjectInternal_EInverseAdd_Call {
	return &MockEObjectInternal_EInverseAdd_Call{Call: _e.mock.On("EInverseAdd", otherEnd, featureID, notifications)}
}

func (_c *MockEObjectInternal_EInverseAdd_Call) Run(run func(otherEnd EObject, featureID int, notifications ENotificationChain)) *MockEObjectInternal_EInverseAdd_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(EObject), args[1].(int), args[2].(ENotificationChain))
	})
	return _c
}

func (_c *MockEObjectInternal_EInverseAdd_Call) Return(_a0 ENotificationChain) *MockEObjectInternal_EInverseAdd_Call {
	_c.Call.Return(_a0)
	return _c
}

// EInverseRemove provides a mock function with given fields: otherEnd, featureID, notifications
func (_m *MockEObjectInternal_Prototype_Methods) EInverseRemove(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain {
	ret := _m.mock.Called(otherEnd, featureID, notifications)

	var r0 ENotificationChain
	if rf, ok := ret.Get(0).(func(EObject, int, ENotificationChain) ENotificationChain); ok {
		r0 = rf(otherEnd, featureID, notifications)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ENotificationChain)
		}
	}

	return r0
}

// MockEObjectInternal_EInverseRemove_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EInverseRemove'
type MockEObjectInternal_EInverseRemove_Call struct {
	*mock.Call
}

// EInverseRemove is a helper method to define mock.On call
//   - otherEnd EObject
//   - featureID int
//   - notifications ENotificationChain
func (_e *MockEObjectInternal_Expecter_Methods) EInverseRemove(otherEnd interface{}, featureID interface{}, notifications interface{}) *MockEObjectInternal_EInverseRemove_Call {
	return &MockEObjectInternal_EInverseRemove_Call{Call: _e.mock.On("EInverseRemove", otherEnd, featureID, notifications)}
}

func (_c *MockEObjectInternal_EInverseRemove_Call) Run(run func(otherEnd EObject, featureID int, notifications ENotificationChain)) *MockEObjectInternal_EInverseRemove_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(EObject), args[1].(int), args[2].(ENotificationChain))
	})
	return _c
}

func (_c *MockEObjectInternal_EInverseRemove_Call) Return(_a0 ENotificationChain) *MockEObjectInternal_EInverseRemove_Call {
	_c.Call.Return(_a0)
	return _c
}

// EInvokeFromID provides a mock function with given fields: operationID, arguments
func (_m *MockEObjectInternal_Prototype_Methods) EInvokeFromID(operationID int, arguments EList) interface{} {
	ret := _m.mock.Called(operationID, arguments)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(int, EList) interface{}); ok {
		r0 = rf(operationID, arguments)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0)
		}
	}

	return r0
}

// MockEObjectInternal_EInvokeFromID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EInvokeFromID'
type MockEObjectInternal_EInvokeFromID_Call struct {
	*mock.Call
}

// EInvokeFromID is a helper method to define mock.On call
//   - operationID int
//   - arguments EList
func (_e *MockEObjectInternal_Expecter_Methods) EInvokeFromID(operationID interface{}, arguments interface{}) *MockEObjectInternal_EInvokeFromID_Call {
	return &MockEObjectInternal_EInvokeFromID_Call{Call: _e.mock.On("EInvokeFromID", operationID, arguments)}
}

func (_c *MockEObjectInternal_EInvokeFromID_Call) Run(run func(operationID int, arguments EList)) *MockEObjectInternal_EInvokeFromID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int), args[1].(EList))
	})
	return _c
}

func (_c *MockEObjectInternal_EInvokeFromID_Call) Return(_a0 interface{}) *MockEObjectInternal_EInvokeFromID_Call {
	_c.Call.Return(_a0)
	return _c
}

// EIsSetFromID provides a mock function with given fields: featureID
func (_m *MockEObjectInternal_Prototype_Methods) EIsSetFromID(featureID int) bool {
	ret := _m.mock.Called(featureID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(int) bool); ok {
		r0 = rf(featureID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockEObjectInternal_EIsSetFromID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EIsSetFromID'
type MockEObjectInternal_EIsSetFromID_Call struct {
	*mock.Call
}

// EIsSetFromID is a helper method to define mock.On call
//   - featureID int
func (_e *MockEObjectInternal_Expecter_Methods) EIsSetFromID(featureID interface{}) *MockEObjectInternal_EIsSetFromID_Call {
	return &MockEObjectInternal_EIsSetFromID_Call{Call: _e.mock.On("EIsSetFromID", featureID)}
}

func (_c *MockEObjectInternal_EIsSetFromID_Call) Run(run func(featureID int)) *MockEObjectInternal_EIsSetFromID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int))
	})
	return _c
}

func (_c *MockEObjectInternal_EIsSetFromID_Call) Return(_a0 bool) *MockEObjectInternal_EIsSetFromID_Call {
	_c.Call.Return(_a0)
	return _c
}

// EObjectForFragmentSegment provides a mock function with given fields: _a0
func (_m *MockEObjectInternal_Prototype_Methods) EObjectForFragmentSegment(_a0 string) EObject {
	ret := _m.mock.Called(_a0)

	var r0 EObject
	if rf, ok := ret.Get(0).(func(string) EObject); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(EObject)
		}
	}

	return r0
}

// MockEObjectInternal_EObjectForFragmentSegment_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EObjectForFragmentSegment'
type MockEObjectInternal_EObjectForFragmentSegment_Call struct {
	*mock.Call
}

// EObjectForFragmentSegment is a helper method to define mock.On call
//   - _a0 string
func (_e *MockEObjectInternal_Expecter_Methods) EObjectForFragmentSegment(_a0 interface{}) *MockEObjectInternal_EObjectForFragmentSegment_Call {
	return &MockEObjectInternal_EObjectForFragmentSegment_Call{Call: _e.mock.On("EObjectForFragmentSegment", _a0)}
}

func (_c *MockEObjectInternal_EObjectForFragmentSegment_Call) Run(run func(_a0 string)) *MockEObjectInternal_EObjectForFragmentSegment_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockEObjectInternal_EObjectForFragmentSegment_Call) Return(_a0 EObject) *MockEObjectInternal_EObjectForFragmentSegment_Call {
	_c.Call.Return(_a0)
	return _c
}

// EOperationID provides a mock function with given fields: operation
func (_m *MockEObjectInternal_Prototype_Methods) EOperationID(operation EOperation) int {
	ret := _m.mock.Called(operation)

	var r0 int
	if rf, ok := ret.Get(0).(func(EOperation) int); ok {
		r0 = rf(operation)
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// MockEObjectInternal_EOperationID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EOperationID'
type MockEObjectInternal_EOperationID_Call struct {
	*mock.Call
}

// EOperationID is a helper method to define mock.On call
//   - operation EOperation
func (_e *MockEObjectInternal_Expecter_Methods) EOperationID(operation interface{}) *MockEObjectInternal_EOperationID_Call {
	return &MockEObjectInternal_EOperationID_Call{Call: _e.mock.On("EOperationID", operation)}
}

func (_c *MockEObjectInternal_EOperationID_Call) Run(run func(operation EOperation)) *MockEObjectInternal_EOperationID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(EOperation))
	})
	return _c
}

func (_c *MockEObjectInternal_EOperationID_Call) Return(_a0 int) *MockEObjectInternal_EOperationID_Call {
	_c.Call.Return(_a0)
	return _c
}

// EProxyURI provides a mock function with given fields:
func (_m *MockEObjectInternal_Prototype_Methods) EProxyURI() *URI {
	ret := _m.mock.Called()

	var r0 *URI
	if rf, ok := ret.Get(0).(func() *URI); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*URI)
		}
	}

	return r0
}

// MockEObjectInternal_EProxyURI_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EProxyURI'
type MockEObjectInternal_EProxyURI_Call struct {
	*mock.Call
}

// EProxyURI is a helper method to define mock.On call
func (_e *MockEObjectInternal_Expecter_Methods) EProxyURI() *MockEObjectInternal_EProxyURI_Call {
	return &MockEObjectInternal_EProxyURI_Call{Call: _e.mock.On("EProxyURI")}
}

func (_c *MockEObjectInternal_EProxyURI_Call) Run(run func()) *MockEObjectInternal_EProxyURI_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockEObjectInternal_EProxyURI_Call) Return(_a0 *URI) *MockEObjectInternal_EProxyURI_Call {
	_c.Call.Return(_a0)
	return _c
}

// EResolveProxy provides a mock function with given fields: proxy
func (_m *MockEObjectInternal_Prototype_Methods) EResolveProxy(proxy EObject) EObject {
	ret := _m.mock.Called(proxy)

	var r0 EObject
	if rf, ok := ret.Get(0).(func(EObject) EObject); ok {
		r0 = rf(proxy)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(EObject)
		}
	}

	return r0
}

// MockEObjectInternal_EResolveProxy_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EResolveProxy'
type MockEObjectInternal_EResolveProxy_Call struct {
	*mock.Call
}

// EResolveProxy is a helper method to define mock.On call
//   - proxy EObject
func (_e *MockEObjectInternal_Expecter_Methods) EResolveProxy(proxy interface{}) *MockEObjectInternal_EResolveProxy_Call {
	return &MockEObjectInternal_EResolveProxy_Call{Call: _e.mock.On("EResolveProxy", proxy)}
}

func (_c *MockEObjectInternal_EResolveProxy_Call) Run(run func(proxy EObject)) *MockEObjectInternal_EResolveProxy_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(EObject))
	})
	return _c
}

func (_c *MockEObjectInternal_EResolveProxy_Call) Return(_a0 EObject) *MockEObjectInternal_EResolveProxy_Call {
	_c.Call.Return(_a0)
	return _c
}

// ESetFromID provides a mock function with given fields: featureID, newValue
func (_m *MockEObjectInternal_Prototype_Methods) ESetFromID(featureID int, newValue interface{}) {
	_m.mock.Called(featureID, newValue)
}

// MockEObjectInternal_ESetFromID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ESetFromID'
type MockEObjectInternal_ESetFromID_Call struct {
	*mock.Call
}

// ESetFromID is a helper method to define mock.On call
//   - featureID int
//   - newValue interface{}
func (_e *MockEObjectInternal_Expecter_Methods) ESetFromID(featureID interface{}, newValue interface{}) *MockEObjectInternal_ESetFromID_Call {
	return &MockEObjectInternal_ESetFromID_Call{Call: _e.mock.On("ESetFromID", featureID, newValue)}
}

func (_c *MockEObjectInternal_ESetFromID_Call) Run(run func(featureID int, newValue interface{})) *MockEObjectInternal_ESetFromID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int), args[1])
	})
	return _c
}

func (_c *MockEObjectInternal_ESetFromID_Call) Return() *MockEObjectInternal_ESetFromID_Call {
	_c.Call.Return()
	return _c
}

// ESetInternalContainer provides a mock function with given fields: container, containerFeatureID
func (_m *MockEObjectInternal_Prototype_Methods) ESetInternalContainer(container EObject, containerFeatureID int) {
	_m.mock.Called(container, containerFeatureID)
}

// MockEObjectInternal_ESetInternalContainer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ESetInternalContainer'
type MockEObjectInternal_ESetInternalContainer_Call struct {
	*mock.Call
}

// ESetInternalContainer is a helper method to define mock.On call
//   - container EObject
//   - containerFeatureID int
func (_e *MockEObjectInternal_Expecter_Methods) ESetInternalContainer(container interface{}, containerFeatureID interface{}) *MockEObjectInternal_ESetInternalContainer_Call {
	return &MockEObjectInternal_ESetInternalContainer_Call{Call: _e.mock.On("ESetInternalContainer", container, containerFeatureID)}
}

func (_c *MockEObjectInternal_ESetInternalContainer_Call) Run(run func(container EObject, containerFeatureID int)) *MockEObjectInternal_ESetInternalContainer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(EObject), args[1].(int))
	})
	return _c
}

func (_c *MockEObjectInternal_ESetInternalContainer_Call) Return() *MockEObjectInternal_ESetInternalContainer_Call {
	_c.Call.Return()
	return _c
}

// ESetInternalResource provides a mock function with given fields: resource
func (_m *MockEObjectInternal_Prototype_Methods) ESetInternalResource(resource EResource) {
	_m.mock.Called(resource)
}

// MockEObjectInternal_ESetInternalResource_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ESetInternalResource'
type MockEObjectInternal_ESetInternalResource_Call struct {
	*mock.Call
}

// ESetInternalResource is a helper method to define mock.On call
//   - resource EResource
func (_e *MockEObjectInternal_Expecter_Methods) ESetInternalResource(resource interface{}) *MockEObjectInternal_ESetInternalResource_Call {
	return &MockEObjectInternal_ESetInternalResource_Call{Call: _e.mock.On("ESetInternalResource", resource)}
}

func (_c *MockEObjectInternal_ESetInternalResource_Call) Run(run func(resource EResource)) *MockEObjectInternal_ESetInternalResource_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(EResource))
	})
	return _c
}

func (_c *MockEObjectInternal_ESetInternalResource_Call) Return() *MockEObjectInternal_ESetInternalResource_Call {
	_c.Call.Return()
	return _c
}

// ESetProxyURI provides a mock function with given fields: uri
func (_m *MockEObjectInternal_Prototype_Methods) ESetProxyURI(uri *URI) {
	_m.mock.Called(uri)
}

// MockEObjectInternal_ESetProxyURI_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ESetProxyURI'
type MockEObjectInternal_ESetProxyURI_Call struct {
	*mock.Call
}

// ESetProxyURI is a helper method to define mock.On call
//   - uri *URI
func (_e *MockEObjectInternal_Expecter_Methods) ESetProxyURI(uri interface{}) *MockEObjectInternal_ESetProxyURI_Call {
	return &MockEObjectInternal_ESetProxyURI_Call{Call: _e.mock.On("ESetProxyURI", uri)}
}

func (_c *MockEObjectInternal_ESetProxyURI_Call) Run(run func(uri *URI)) *MockEObjectInternal_ESetProxyURI_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*URI))
	})
	return _c
}

func (_c *MockEObjectInternal_ESetProxyURI_Call) Return() *MockEObjectInternal_ESetProxyURI_Call {
	_c.Call.Return()
	return _c
}

// ESetResource provides a mock function with given fields: resource, notifications
func (_m *MockEObjectInternal_Prototype_Methods) ESetResource(resource EResource, notifications ENotificationChain) ENotificationChain {
	ret := _m.mock.Called(resource, notifications)

	var r0 ENotificationChain
	if rf, ok := ret.Get(0).(func(EResource, ENotificationChain) ENotificationChain); ok {
		r0 = rf(resource, notifications)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ENotificationChain)
		}
	}

	return r0
}

// MockEObjectInternal_ESetResource_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ESetResource'
type MockEObjectInternal_ESetResource_Call struct {
	*mock.Call
}

// ESetResource is a helper method to define mock.On call
//   - resource EResource
//   - notifications ENotificationChain
func (_e *MockEObjectInternal_Expecter_Methods) ESetResource(resource interface{}, notifications interface{}) *MockEObjectInternal_ESetResource_Call {
	return &MockEObjectInternal_ESetResource_Call{Call: _e.mock.On("ESetResource", resource, notifications)}
}

func (_c *MockEObjectInternal_ESetResource_Call) Run(run func(resource EResource, notifications ENotificationChain)) *MockEObjectInternal_ESetResource_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(EResource), args[1].(ENotificationChain))
	})
	return _c
}

func (_c *MockEObjectInternal_ESetResource_Call) Return(_a0 ENotificationChain) *MockEObjectInternal_ESetResource_Call {
	_c.Call.Return(_a0)
	return _c
}

// EStaticClass provides a mock function with given fields:
func (_m *MockEObjectInternal_Prototype_Methods) EStaticClass() EClass {
	ret := _m.mock.Called()

	var r0 EClass
	if rf, ok := ret.Get(0).(func() EClass); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(EClass)
		}
	}

	return r0
}

// MockEObjectInternal_EStaticClass_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EStaticClass'
type MockEObjectInternal_EStaticClass_Call struct {
	*mock.Call
}

// EStaticClass is a helper method to define mock.On call
func (_e *MockEObjectInternal_Expecter_Methods) EStaticClass() *MockEObjectInternal_EStaticClass_Call {
	return &MockEObjectInternal_EStaticClass_Call{Call: _e.mock.On("EStaticClass")}
}

func (_c *MockEObjectInternal_EStaticClass_Call) Run(run func()) *MockEObjectInternal_EStaticClass_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockEObjectInternal_EStaticClass_Call) Return(_a0 EClass) *MockEObjectInternal_EStaticClass_Call {
	_c.Call.Return(_a0)
	return _c
}

// EStaticFeatureCount provides a mock function with given fields:
func (_m *MockEObjectInternal_Prototype_Methods) EStaticFeatureCount() int {
	ret := _m.mock.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// MockEObjectInternal_EStaticFeatureCount_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EStaticFeatureCount'
type MockEObjectInternal_EStaticFeatureCount_Call struct {
	*mock.Call
}

// EStaticFeatureCount is a helper method to define mock.On call
func (_e *MockEObjectInternal_Expecter_Methods) EStaticFeatureCount() *MockEObjectInternal_EStaticFeatureCount_Call {
	return &MockEObjectInternal_EStaticFeatureCount_Call{Call: _e.mock.On("EStaticFeatureCount")}
}

func (_c *MockEObjectInternal_EStaticFeatureCount_Call) Run(run func()) *MockEObjectInternal_EStaticFeatureCount_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockEObjectInternal_EStaticFeatureCount_Call) Return(_a0 int) *MockEObjectInternal_EStaticFeatureCount_Call {
	_c.Call.Return(_a0)
	return _c
}

// EURIFragmentSegment provides a mock function with given fields: _a0, _a1
func (_m *MockEObjectInternal_Prototype_Methods) EURIFragmentSegment(_a0 EStructuralFeature, _a1 EObject) string {
	ret := _m.mock.Called(_a0, _a1)

	var r0 string
	if rf, ok := ret.Get(0).(func(EStructuralFeature, EObject) string); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockEObjectInternal_EURIFragmentSegment_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EURIFragmentSegment'
type MockEObjectInternal_EURIFragmentSegment_Call struct {
	*mock.Call
}

// EURIFragmentSegment is a helper method to define mock.On call
//   - _a0 EStructuralFeature
//   - _a1 EObject
func (_e *MockEObjectInternal_Expecter_Methods) EURIFragmentSegment(_a0 interface{}, _a1 interface{}) *MockEObjectInternal_EURIFragmentSegment_Call {
	return &MockEObjectInternal_EURIFragmentSegment_Call{Call: _e.mock.On("EURIFragmentSegment", _a0, _a1)}
}

func (_c *MockEObjectInternal_EURIFragmentSegment_Call) Run(run func(_a0 EStructuralFeature, _a1 EObject)) *MockEObjectInternal_EURIFragmentSegment_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(EStructuralFeature), args[1].(EObject))
	})
	return _c
}

func (_c *MockEObjectInternal_EURIFragmentSegment_Call) Return(_a0 string) *MockEObjectInternal_EURIFragmentSegment_Call {
	_c.Call.Return(_a0)
	return _c
}

// EUnsetFromID provides a mock function with given fields: featureID
func (_m *MockEObjectInternal_Prototype_Methods) EUnsetFromID(featureID int) {
	_m.mock.Called(featureID)
}

// MockEObjectInternal_EUnsetFromID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EUnsetFromID'
type MockEObjectInternal_EUnsetFromID_Call struct {
	*mock.Call
}

// EUnsetFromID is a helper method to define mock.On call
//   - featureID int
func (_e *MockEObjectInternal_Expecter_Methods) EUnsetFromID(featureID interface{}) *MockEObjectInternal_EUnsetFromID_Call {
	return &MockEObjectInternal_EUnsetFromID_Call{Call: _e.mock.On("EUnsetFromID", featureID)}
}

func (_c *MockEObjectInternal_EUnsetFromID_Call) Run(run func(featureID int)) *MockEObjectInternal_EUnsetFromID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int))
	})
	return _c
}

func (_c *MockEObjectInternal_EUnsetFromID_Call) Return() *MockEObjectInternal_EUnsetFromID_Call {
	_c.Call.Return()
	return _c
}

type mockConstructorTestingTNewMockEObjectInternal interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockEObjectInternal creates a new instance of MockEObjectInternal_Prototype. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockEObjectInternal(t mockConstructorTestingTNewMockEObjectInternal) *MockEObjectInternal {
	mock := &MockEObjectInternal{}
	mock.SetMock(&mock.Mock)
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
