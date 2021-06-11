// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package ecore

import mock "github.com/stretchr/testify/mock"

// MockEResourceCodec is an autogenerated mock type for the EResourceCodec type
type MockEResourceCodec struct {
	mock.Mock
}

// NewDecoder provides a mock function with given fields: options
func (_m *MockEResourceCodec) NewDecoder(options map[string]interface{}) EResourceDecoder {
	ret := _m.Called(options)

	var r0 EResourceDecoder
	if rf, ok := ret.Get(0).(func(map[string]interface{}) EResourceDecoder); ok {
		r0 = rf(options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(EResourceDecoder)
		}
	}

	return r0
}

// NewEncoder provides a mock function with given fields: options
func (_m *MockEResourceCodec) NewEncoder(options map[string]interface{}) EResourceEncoder {
	ret := _m.Called(options)

	var r0 EResourceEncoder
	if rf, ok := ret.Get(0).(func(map[string]interface{}) EResourceEncoder); ok {
		r0 = rf(options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(EResourceEncoder)
		}
	}

	return r0
}