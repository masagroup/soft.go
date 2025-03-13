// Code generated by mockery v2.53.2. DO NOT EDIT.

package ecore

import mock "github.com/stretchr/testify/mock"

// MockEObjectProperties is an autogenerated mock type for the EDynamicProperties type
type MockEObjectProperties struct {
	mock.Mock
}

type MockEObjectProperties_Expecter struct {
	mock *mock.Mock
}

func (_m *MockEObjectProperties) EXPECT() *MockEObjectProperties_Expecter {
	return &MockEObjectProperties_Expecter{mock: &_m.Mock}
}

// EDynamicGet provides a mock function with given fields: dynamicFeatureID
func (_m *MockEObjectProperties) EDynamicGet(dynamicFeatureID int) interface{} {
	ret := _m.Called(dynamicFeatureID)

	if len(ret) == 0 {
		panic("no return value specified for EDynamicGet")
	}

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

// MockEObjectProperties_EDynamicGet_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EDynamicGet'
type MockEObjectProperties_EDynamicGet_Call struct {
	*mock.Call
}

// EDynamicGet is a helper method to define mock.On call
//   - dynamicFeatureID int
func (_e *MockEObjectProperties_Expecter) EDynamicGet(dynamicFeatureID interface{}) *MockEObjectProperties_EDynamicGet_Call {
	return &MockEObjectProperties_EDynamicGet_Call{Call: _e.mock.On("EDynamicGet", dynamicFeatureID)}
}

func (_c *MockEObjectProperties_EDynamicGet_Call) Run(run func(dynamicFeatureID int)) *MockEObjectProperties_EDynamicGet_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int))
	})
	return _c
}

func (_c *MockEObjectProperties_EDynamicGet_Call) Return(_a0 interface{}) *MockEObjectProperties_EDynamicGet_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockEObjectProperties_EDynamicGet_Call) RunAndReturn(run func(int) interface{}) *MockEObjectProperties_EDynamicGet_Call {
	_c.Call.Return(run)
	return _c
}

// EDynamicIsSet provides a mock function with given fields: dynamicFeatureID
func (_m *MockEObjectProperties) EDynamicIsSet(dynamicFeatureID int) bool {
	ret := _m.Called(dynamicFeatureID)

	if len(ret) == 0 {
		panic("no return value specified for EDynamicIsSet")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(int) bool); ok {
		r0 = rf(dynamicFeatureID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockEObjectProperties_EDynamicIsSet_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EDynamicIsSet'
type MockEObjectProperties_EDynamicIsSet_Call struct {
	*mock.Call
}

// EDynamicIsSet is a helper method to define mock.On call
//   - dynamicFeatureID int
func (_e *MockEObjectProperties_Expecter) EDynamicIsSet(dynamicFeatureID interface{}) *MockEObjectProperties_EDynamicIsSet_Call {
	return &MockEObjectProperties_EDynamicIsSet_Call{Call: _e.mock.On("EDynamicIsSet", dynamicFeatureID)}
}

func (_c *MockEObjectProperties_EDynamicIsSet_Call) Run(run func(dynamicFeatureID int)) *MockEObjectProperties_EDynamicIsSet_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int))
	})
	return _c
}

func (_c *MockEObjectProperties_EDynamicIsSet_Call) Return(_a0 bool) *MockEObjectProperties_EDynamicIsSet_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockEObjectProperties_EDynamicIsSet_Call) RunAndReturn(run func(int) bool) *MockEObjectProperties_EDynamicIsSet_Call {
	_c.Call.Return(run)
	return _c
}

// EDynamicSet provides a mock function with given fields: dynamicFeatureID, newValue
func (_m *MockEObjectProperties) EDynamicSet(dynamicFeatureID int, newValue interface{}) {
	_m.Called(dynamicFeatureID, newValue)
}

// MockEObjectProperties_EDynamicSet_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EDynamicSet'
type MockEObjectProperties_EDynamicSet_Call struct {
	*mock.Call
}

// EDynamicSet is a helper method to define mock.On call
//   - dynamicFeatureID int
//   - newValue interface{}
func (_e *MockEObjectProperties_Expecter) EDynamicSet(dynamicFeatureID interface{}, newValue interface{}) *MockEObjectProperties_EDynamicSet_Call {
	return &MockEObjectProperties_EDynamicSet_Call{Call: _e.mock.On("EDynamicSet", dynamicFeatureID, newValue)}
}

func (_c *MockEObjectProperties_EDynamicSet_Call) Run(run func(dynamicFeatureID int, newValue interface{})) *MockEObjectProperties_EDynamicSet_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int), args[1].(interface{}))
	})
	return _c
}

func (_c *MockEObjectProperties_EDynamicSet_Call) Return() *MockEObjectProperties_EDynamicSet_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockEObjectProperties_EDynamicSet_Call) RunAndReturn(run func(int, interface{})) *MockEObjectProperties_EDynamicSet_Call {
	_c.Run(run)
	return _c
}

// EDynamicUnset provides a mock function with given fields: dynamicFeatureID
func (_m *MockEObjectProperties) EDynamicUnset(dynamicFeatureID int) {
	_m.Called(dynamicFeatureID)
}

// MockEObjectProperties_EDynamicUnset_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EDynamicUnset'
type MockEObjectProperties_EDynamicUnset_Call struct {
	*mock.Call
}

// EDynamicUnset is a helper method to define mock.On call
//   - dynamicFeatureID int
func (_e *MockEObjectProperties_Expecter) EDynamicUnset(dynamicFeatureID interface{}) *MockEObjectProperties_EDynamicUnset_Call {
	return &MockEObjectProperties_EDynamicUnset_Call{Call: _e.mock.On("EDynamicUnset", dynamicFeatureID)}
}

func (_c *MockEObjectProperties_EDynamicUnset_Call) Run(run func(dynamicFeatureID int)) *MockEObjectProperties_EDynamicUnset_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int))
	})
	return _c
}

func (_c *MockEObjectProperties_EDynamicUnset_Call) Return() *MockEObjectProperties_EDynamicUnset_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockEObjectProperties_EDynamicUnset_Call) RunAndReturn(run func(int)) *MockEObjectProperties_EDynamicUnset_Call {
	_c.Run(run)
	return _c
}

// NewMockEObjectProperties creates a new instance of MockEObjectProperties. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockEObjectProperties(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockEObjectProperties {
	mock := &MockEObjectProperties{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
