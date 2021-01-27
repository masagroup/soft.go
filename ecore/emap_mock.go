package ecore

type MockEMap struct {
	MockEList
}

// ContainsKey provides a mock function with given fields: key
func (_m *MockEMap) ContainsKey(key interface{}) bool {
	ret := _m.Called(key)

	var r0 bool
	if rf, ok := ret.Get(0).(func(interface{}) bool); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ContainsValue provides a mock function with given fields: value
func (_m *MockEMap) ContainsValue(value interface{}) bool {
	ret := _m.Called(value)

	var r0 bool
	if rf, ok := ret.Get(0).(func(interface{}) bool); ok {
		r0 = rf(value)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// GetValue provides a mock function with given fields: value
func (_m *MockEMap) GetValue(value interface{}) interface{} {
	ret := _m.Called(value)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(interface{}) interface{}); ok {
		r0 = rf(value)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// Put provides a mock function with given fields: key, value
func (_m *MockEMap) Put(key interface{}, value interface{}) {
	_m.Called(key, value)
}

// RemoveKey provides a mock function with given fields: key
func (_m *MockEMap) RemoveKey(key interface{}) interface{} {
	ret := _m.Called(key)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(interface{}) interface{}); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// ToMap provides a mock function with given fields:
func (_m *MockEMap) ToMap() map[interface{}]interface{} {
	ret := _m.Called()

	var r0 map[interface{}]interface{}
	if rf, ok := ret.Get(0).(func() map[interface{}]interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[interface{}]interface{})
		}
	}

	return r0
}
