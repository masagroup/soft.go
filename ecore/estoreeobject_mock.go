// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package ecore

import mock "github.com/stretchr/testify/mock"

// MockEStoreEObject is an autogenerated mock type for the EStoreEObject type
type MockEStoreEObject struct {
	mock.Mock
}

// EAdapters provides a mock function with given fields:
func (_m *MockEStoreEObject) EAdapters() EList {
	ret := _m.Called()

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

// EAllContents provides a mock function with given fields:
func (_m *MockEStoreEObject) EAllContents() EIterator {
	ret := _m.Called()

	var r0 EIterator
	if rf, ok := ret.Get(0).(func() EIterator); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(EIterator)
		}
	}

	return r0
}

// EClass provides a mock function with given fields:
func (_m *MockEStoreEObject) EClass() EClass {
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

// EContainer provides a mock function with given fields:
func (_m *MockEStoreEObject) EContainer() EObject {
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

// EContainingFeature provides a mock function with given fields:
func (_m *MockEStoreEObject) EContainingFeature() EStructuralFeature {
	ret := _m.Called()

	var r0 EStructuralFeature
	if rf, ok := ret.Get(0).(func() EStructuralFeature); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(EStructuralFeature)
		}
	}

	return r0
}

// EContainmentFeature provides a mock function with given fields:
func (_m *MockEStoreEObject) EContainmentFeature() EReference {
	ret := _m.Called()

	var r0 EReference
	if rf, ok := ret.Get(0).(func() EReference); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(EReference)
		}
	}

	return r0
}

// EContents provides a mock function with given fields:
func (_m *MockEStoreEObject) EContents() EList {
	ret := _m.Called()

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

// ECrossReferences provides a mock function with given fields:
func (_m *MockEStoreEObject) ECrossReferences() EList {
	ret := _m.Called()

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

// EDeliver provides a mock function with given fields:
func (_m *MockEStoreEObject) EDeliver() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// EGet provides a mock function with given fields: _a0
func (_m *MockEStoreEObject) EGet(_a0 EStructuralFeature) interface{} {
	ret := _m.Called(_a0)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(EStructuralFeature) interface{}); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// EGetResolve provides a mock function with given fields: _a0, _a1
func (_m *MockEStoreEObject) EGetResolve(_a0 EStructuralFeature, _a1 bool) interface{} {
	ret := _m.Called(_a0, _a1)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(EStructuralFeature, bool) interface{}); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// EInvoke provides a mock function with given fields: _a0, _a1
func (_m *MockEStoreEObject) EInvoke(_a0 EOperation, _a1 EList) interface{} {
	ret := _m.Called(_a0, _a1)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(EOperation, EList) interface{}); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// EIsProxy provides a mock function with given fields:
func (_m *MockEStoreEObject) EIsProxy() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// EIsSet provides a mock function with given fields: _a0
func (_m *MockEStoreEObject) EIsSet(_a0 EStructuralFeature) bool {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(EStructuralFeature) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ENotificationRequired provides a mock function with given fields:
func (_m *MockEStoreEObject) ENotificationRequired() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ENotify provides a mock function with given fields: _a0
func (_m *MockEStoreEObject) ENotify(_a0 ENotification) {
	_m.Called(_a0)
}

// EResource provides a mock function with given fields:
func (_m *MockEStoreEObject) EResource() EResource {
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

// ESet provides a mock function with given fields: _a0, _a1
func (_m *MockEStoreEObject) ESet(_a0 EStructuralFeature, _a1 interface{}) {
	_m.Called(_a0, _a1)
}

// ESetDeliver provides a mock function with given fields: _a0
func (_m *MockEStoreEObject) ESetDeliver(_a0 bool) {
	_m.Called(_a0)
}

// EStore provides a mock function with given fields:
func (_m *MockEStoreEObject) EStore() EStore {
	ret := _m.Called()

	var r0 EStore
	if rf, ok := ret.Get(0).(func() EStore); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(EStore)
		}
	}

	return r0
}

// EUnset provides a mock function with given fields: _a0
func (_m *MockEStoreEObject) EUnset(_a0 EStructuralFeature) {
	_m.Called(_a0)
}
