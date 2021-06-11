// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package ecore

import (
	io "io"
)

// MockEResourceInternal is an autogenerated mock type for the EResourceInternal type
type MockEResourceInternal struct {
	MockEResource
}

// IsLoading provides a mock function with given fields:
func (_m *MockEResourceInternal) IsLoading() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// BasicSetLoaded provides a mock function with given fields: _a0, _a1
func (_m *MockEResourceInternal) BasicSetLoaded(_a0 bool, _a1 ENotificationChain) ENotificationChain {
	ret := _m.Called(_a0, _a1)

	var r0 ENotificationChain
	if rf, ok := ret.Get(0).(func(bool, ENotificationChain) ENotificationChain); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ENotificationChain)
		}
	}

	return r0
}

// BasicSetResourceSet provides a mock function with given fields: _a0, _a1
func (_m *MockEResourceInternal) BasicSetResourceSet(_a0 EResourceSet, _a1 ENotificationChain) ENotificationChain {
	ret := _m.Called(_a0, _a1)

	var r0 ENotificationChain
	if rf, ok := ret.Get(0).(func(EResourceSet, ENotificationChain) ENotificationChain); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ENotificationChain)
		}
	}

	return r0
}

// DoAttached provides a mock function with given fields: o
func (_m *MockEResourceInternal) DoAttached(o EObject) {
	_m.Called(o)
}

// DoDetached provides a mock function with given fields: o
func (_m *MockEResourceInternal) DoDetached(o EObject) {
	_m.Called(o)
}

// DoLoad provides a mock function with given fields: rd, options
func (_m *MockEResourceInternal) DoLoad(rd io.Reader, d EResourceDecoder) {
	_m.Called(rd, d)
}

// DoSave provides a mock function with given fields: rd, options
func (_m *MockEResourceInternal) DoSave(rd io.Writer, e EResourceEncoder) {
	_m.Called(rd, e)
}

// DoUnload provides a mock function with given fields:
func (_m *MockEResourceInternal) DoUnload() {
	_m.Called()
}

// IsAttachedDetachedRequired provides a mock function with given fields:
func (_m *MockEResourceInternal) IsAttachedDetachedRequired() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
