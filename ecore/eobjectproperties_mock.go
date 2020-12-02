// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package ecore

import mock "github.com/stretchr/testify/mock"

// MockEObjectProperties is an autogenerated mock type for the EObjectProperties type
type MockEObjectProperties struct {
	mock.Mock
}

// EDynamicGet provides a mock function with given fields: dynamicFeatureID
func (_m *MockEObjectProperties) EDynamicGet(dynamicFeatureID int) interface{} {
	ret := _m.Called(dynamicFeatureID)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(int) interface{}); ok {
		r0 = rf(dynamicFeatureID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// EDynamicSet provides a mock function with given fields: dynamicFeatureID, newValue
func (_m *MockEObjectProperties) EDynamicSet(dynamicFeatureID int, newValue interface{}) {
	_m.Called(dynamicFeatureID, newValue)
}

// EDynamicUnset provides a mock function with given fields: dynamicFeatureID
func (_m *MockEObjectProperties) EDynamicUnset(dynamicFeatureID int) {
	_m.Called(dynamicFeatureID)
}