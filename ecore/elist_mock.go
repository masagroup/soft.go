// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package ecore

import mock "github.com/stretchr/testify/mock"

// MockEList is an autogenerated mock type for the EList type
type MockEList struct {
	mock.Mock
}

// Add provides a mock function with given fields: _a0
func (_m *MockEList) Add(_a0 interface{}) bool {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(interface{}) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// AddAll provides a mock function with given fields: _a0
func (_m *MockEList) AddAll(_a0 EList) bool {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(EList) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Clear provides a mock function with given fields:
func (_m *MockEList) Clear() {
	_m.Called()
}

// Contains provides a mock function with given fields: _a0
func (_m *MockEList) Contains(_a0 interface{}) bool {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(interface{}) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Empty provides a mock function with given fields:
func (_m *MockEList) Empty() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Get provides a mock function with given fields: _a0
func (_m *MockEList) Get(_a0 int) interface{} {
	ret := _m.Called(_a0)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(int) interface{}); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// IndexOf provides a mock function with given fields: _a0
func (_m *MockEList) IndexOf(_a0 interface{}) int {
	ret := _m.Called(_a0)

	var r0 int
	if rf, ok := ret.Get(0).(func(interface{}) int); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// Insert provides a mock function with given fields: _a0, _a1
func (_m *MockEList) Insert(_a0 int, _a1 interface{}) bool {
	ret := _m.Called(_a0, _a1)

	var r0 bool
	if rf, ok := ret.Get(0).(func(int, interface{}) bool); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// InsertAll provides a mock function with given fields: _a0, _a1
func (_m *MockEList) InsertAll(_a0 int, _a1 EList) bool {
	ret := _m.Called(_a0, _a1)

	var r0 bool
	if rf, ok := ret.Get(0).(func(int, EList) bool); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Iterator provides a mock function with given fields:
func (_m *MockEList) Iterator() EIterator {
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

// Move provides a mock function with given fields: _a0, _a1
func (_m *MockEList) Move(_a0 int, _a1 int) interface{} {
	ret := _m.Called(_a0, _a1)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(int, int) interface{}); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// MoveObject provides a mock function with given fields: _a0, _a1
func (_m *MockEList) MoveObject(_a0 int, _a1 interface{}) {
	_m.Called(_a0, _a1)
}

// Remove provides a mock function with given fields: _a0
func (_m *MockEList) Remove(_a0 interface{}) bool {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(interface{}) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// RemoveAt provides a mock function with given fields: _a0
func (_m *MockEList) RemoveAt(_a0 int) interface{} {
	ret := _m.Called(_a0)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(int) interface{}); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// RemoveAll provides a mock function with given fields: _a0
func (_m *MockEList) RemoveAll(_a0 EList) bool {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(EList) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Set provides a mock function with given fields: _a0, _a1
func (_m *MockEList) Set(_a0 int, _a1 interface{}) interface{} {
	ret := _m.Called(_a0, _a1)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(int, interface{}) interface{}); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// Size provides a mock function with given fields:
func (_m *MockEList) Size() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// ToArray provides a mock function with given fields:
func (_m *MockEList) ToArray() []interface{} {
	ret := _m.Called()

	var r0 []interface{}
	if rf, ok := ret.Get(0).(func() []interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]interface{})
		}
	}

	return r0
}
