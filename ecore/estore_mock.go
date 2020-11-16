// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package ecore

import mock "github.com/stretchr/testify/mock"

// MockEStore is an autogenerated mock type for the EStore type
type MockEStore struct {
	mock.Mock
}

// Add provides a mock function with given fields: object, feature, index, value
func (_m *MockEStore) Add(object EObject, feature EStructuralFeature, index int, value interface{}) {
	_m.Called(object, feature, index, value)
}

// Clear provides a mock function with given fields: object, feature
func (_m *MockEStore) Clear(object EObject, feature EStructuralFeature) {
	_m.Called(object, feature)
}

// Contains provides a mock function with given fields: object, feature, value
func (_m *MockEStore) Contains(object EObject, feature EStructuralFeature, value interface{}) bool {
	ret := _m.Called(object, feature, value)

	var r0 bool
	if rf, ok := ret.Get(0).(func(EObject, EStructuralFeature, interface{}) bool); ok {
		r0 = rf(object, feature, value)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Create provides a mock function with given fields: eClass
func (_m *MockEStore) Create(eClass EClass) EObject {
	ret := _m.Called(eClass)

	var r0 EObject
	if rf, ok := ret.Get(0).(func(EClass) EObject); ok {
		r0 = rf(eClass)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(EObject)
		}
	}

	return r0
}

// Get provides a mock function with given fields: object, feature, index
func (_m *MockEStore) Get(object EObject, feature EStructuralFeature, index int) interface{} {
	ret := _m.Called(object, feature, index)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(EObject, EStructuralFeature, int) interface{}); ok {
		r0 = rf(object, feature, index)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// GetContainer provides a mock function with given fields: object
func (_m *MockEStore) GetContainer(object EObject) EObject {
	ret := _m.Called(object)

	var r0 EObject
	if rf, ok := ret.Get(0).(func(EObject) EObject); ok {
		r0 = rf(object)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(EObject)
		}
	}

	return r0
}

// GetContainingFeature provides a mock function with given fields: object
func (_m *MockEStore) GetContainingFeature(object EObject) EStructuralFeature {
	ret := _m.Called(object)

	var r0 EStructuralFeature
	if rf, ok := ret.Get(0).(func(EObject) EStructuralFeature); ok {
		r0 = rf(object)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(EStructuralFeature)
		}
	}

	return r0
}

// IndexOf provides a mock function with given fields: object, feature, value
func (_m *MockEStore) IndexOf(object EObject, feature EStructuralFeature, value interface{}) int {
	ret := _m.Called(object, feature, value)

	var r0 int
	if rf, ok := ret.Get(0).(func(EObject, EStructuralFeature, interface{}) int); ok {
		r0 = rf(object, feature, value)
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// IsEmpty provides a mock function with given fields: object, feature
func (_m *MockEStore) IsEmpty(object EObject, feature EStructuralFeature) {
	_m.Called(object, feature)
}

// IsSet provides a mock function with given fields: object, feature
func (_m *MockEStore) IsSet(object EObject, feature EStructuralFeature) bool {
	ret := _m.Called(object, feature)

	var r0 bool
	if rf, ok := ret.Get(0).(func(EObject, EStructuralFeature) bool); ok {
		r0 = rf(object, feature)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// LastIndexOf provides a mock function with given fields: object, feature, value
func (_m *MockEStore) LastIndexOf(object EObject, feature EStructuralFeature, value interface{}) int {
	ret := _m.Called(object, feature, value)

	var r0 int
	if rf, ok := ret.Get(0).(func(EObject, EStructuralFeature, interface{}) int); ok {
		r0 = rf(object, feature, value)
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// Move provides a mock function with given fields: object, feature, targetIndex, sourceIndex
func (_m *MockEStore) Move(object EObject, feature EStructuralFeature, targetIndex int, sourceIndex int) interface{} {
	ret := _m.Called(object, feature, targetIndex, sourceIndex)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(EObject, EStructuralFeature, int, int) interface{}); ok {
		r0 = rf(object, feature, targetIndex, sourceIndex)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// Remove provides a mock function with given fields: object, feature, index
func (_m *MockEStore) Remove(object EObject, feature EStructuralFeature, index int) interface{} {
	ret := _m.Called(object, feature, index)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(EObject, EStructuralFeature, int) interface{}); ok {
		r0 = rf(object, feature, index)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// Set provides a mock function with given fields: object, feature, index, value
func (_m *MockEStore) Set(object EObject, feature EStructuralFeature, index int, value interface{}) {
	_m.Called(object, feature, index, value)
}

// Size provides a mock function with given fields: object, feature
func (_m *MockEStore) Size(object EObject, feature EStructuralFeature) int {
	ret := _m.Called(object, feature)

	var r0 int
	if rf, ok := ret.Get(0).(func(EObject, EStructuralFeature) int); ok {
		r0 = rf(object, feature)
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// ToArray provides a mock function with given fields: object, feature
func (_m *MockEStore) ToArray(object EObject, feature EStructuralFeature) []interface{} {
	ret := _m.Called(object, feature)

	var r0 []interface{}
	if rf, ok := ret.Get(0).(func(EObject, EStructuralFeature) []interface{}); ok {
		r0 = rf(object, feature)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]interface{})
		}
	}

	return r0
}

// UnSet provides a mock function with given fields: object, feature
func (_m *MockEStore) UnSet(object EObject, feature EStructuralFeature) {
	_m.Called(object, feature)
}
