// *****************************************************************************
//
// This file is part of a MASA library or program.
// Refer to the included end-user license agreement for restrictions.
//
// Copyright (c) 2020 MASA Group
//
// *****************************************************************************

package ecore

import (
	url "net/url"
)

// MockEObjectInternal is an autogenerated mock type for the EObjectInternal type
type MockEObjectInternal struct {
	MockEObject
}

// EBasicInverseAdd provides a mock function with given fields: otherEnd, featureID, notifications
func (_m *MockEObjectInternal) EBasicInverseAdd(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain {
	ret := _m.Called(otherEnd, featureID, notifications)

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

// EBasicInverseRemove provides a mock function with given fields: otherEnd, featureID, notifications
func (_m *MockEObjectInternal) EBasicInverseRemove(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain {
	ret := _m.Called(otherEnd, featureID, notifications)

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

// EDerivedFeatureID provides a mock function with given fields: container, featureID
func (_m *MockEObjectInternal) EDerivedFeatureID(container EObject, featureID int) int {
	ret := _m.Called(container, featureID)

	var r0 int
	if rf, ok := ret.Get(0).(func(EObject, int) int); ok {
		r0 = rf(container, featureID)
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// EDerivedOperationID provides a mock function with given fields: container, operationID
func (_m *MockEObjectInternal) EDerivedOperationID(container EObject, operationID int) int {
	ret := _m.Called(container, operationID)

	var r0 int
	if rf, ok := ret.Get(0).(func(EObject, int) int); ok {
		r0 = rf(container, operationID)
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// EGetFromID provides a mock function with given fields: featureID, resolve, core
func (_m *MockEObjectInternal) EGetFromID(featureID int, resolve bool) interface{} {
	ret := _m.Called(featureID, resolve)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(int, bool) interface{}); ok {
		r0 = rf(featureID, resolve)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// EInternalContainer provides a mock function with given fields:
func (_m *MockEObjectInternal) EInternalContainer() EObject {
	ret := _m.Called()

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

// EInternalResource provides a mock function with given fields:
func (_m *MockEObjectInternal) EInternalResource() EResource {
	ret := _m.Called()

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

// EInverseAdd provides a mock function with given fields: otherEnd, featureID, notifications
func (_m *MockEObjectInternal) EInverseAdd(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain {
	ret := _m.Called(otherEnd, featureID, notifications)

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

// EInverseRemove provides a mock function with given fields: otherEnd, featureID, notifications
func (_m *MockEObjectInternal) EInverseRemove(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain {
	ret := _m.Called(otherEnd, featureID, notifications)

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

// EInvokeFromID provides a mock function with given fields: operationID, arguments
func (_m *MockEObjectInternal) EInvokeFromID(operationID int, arguments EList) interface{} {
	ret := _m.Called(operationID, arguments)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(int, EList) interface{}); ok {
		r0 = rf(operationID, arguments)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// EIsSetFromID provides a mock function with given fields: featureID
func (_m *MockEObjectInternal) EIsSetFromID(featureID int) bool {
	ret := _m.Called(featureID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(int) bool); ok {
		r0 = rf(featureID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// EObjectForFragmentSegment provides a mock function with given fields: _a0
func (_m *MockEObjectInternal) EObjectForFragmentSegment(_a0 string) EObject {
	ret := _m.Called(_a0)

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

// EProxyURI provides a mock function with given fields:
func (_m *MockEObjectInternal) EProxyURI() *url.URL {
	ret := _m.Called()

	var r0 *url.URL
	if rf, ok := ret.Get(0).(func() *url.URL); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*url.URL)
		}
	}

	return r0
}

// EProxyURI provides a mock function with given fields:
func (_m *MockEObjectInternal) EInternalContainerFeatureID() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(int)
		}
	}

	return r0
}

// EResolveProxy provides a mock function with given fields: proxy
func (_m *MockEObjectInternal) EResolveProxy(proxy EObject) EObject {
	ret := _m.Called(proxy)

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

// ESetFromID provides a mock function with given fields: featureID, newValue
func (_m *MockEObjectInternal) ESetFromID(featureID int, newValue interface{}) {
	_m.Called(featureID, newValue)
}

// ESetProxyURI provides a mock function with given fields: uri
func (_m *MockEObjectInternal) ESetProxyURI(uri *url.URL) {
	_m.Called(uri)
}

// ESetInternalContainer provides a mock function with given fields: container, containerFeatureID
func (_m *MockEObjectInternal) ESetInternalContainer(container EObject, containerFeatureID int) {
	_m.Called(container, containerFeatureID)
}

// ESetInternalResource provides a mock function with given fields: resource
func (_m *MockEObjectInternal) ESetInternalResource(resource EResource) {
	_m.Called(resource)
}

// ESetResource provides a mock function with given fields: resource, notifications
func (_m *MockEObjectInternal) ESetResource(resource EResource, notifications ENotificationChain) ENotificationChain {
	ret := _m.Called(resource, notifications)

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

// EStaticClass provides a mock function with given fields:
func (_m *MockEObjectInternal) EStaticClass() EClass {
	ret := _m.Called()

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

// EStaticFeatureCount provides a mock function with given fields: featureID
func (_m *MockEObjectInternal) EStaticFeatureCount() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// EProperties provides a mock function with given fields: featureID
func (_m *MockEObjectInternal) EDynamicProperties() EDynamicProperties {
	ret := _m.Called()

	var r0 EDynamicProperties
	if rf, ok := ret.Get(0).(func() EDynamicProperties); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(EDynamicProperties)
	}

	return r0
}

// EURIFragmentSegment provides a mock function with given fields: _a0, _a1
func (_m *MockEObjectInternal) EURIFragmentSegment(_a0 EStructuralFeature, _a1 EObject) string {
	ret := _m.Called(_a0, _a1)

	var r0 string
	if rf, ok := ret.Get(0).(func(EStructuralFeature, EObject) string); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// EUnsetFromID provides a mock function with given fields: featureID
func (_m *MockEObjectInternal) EUnsetFromID(featureID int) {
	_m.Called(featureID)
}
