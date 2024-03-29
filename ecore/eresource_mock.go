// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import (
	io "io"

	mock "github.com/stretchr/testify/mock"
)

// MockEResource is an mock type for the EResource type
type MockEResource struct {
	MockEResource_Prototype
	mock.Mock
}

// MockEResource_Prototype is the mock implementation of all EResource methods ( inherited and declared )
type MockEResource_Prototype struct {
	mock *mock.Mock
	MockENotifier_Prototype
	MockEResourceListener_Prototype
	MockEResource_Prototype_Methods
}

func (_mp *MockEResource_Prototype) SetMock(mock *mock.Mock) {
	_mp.mock = mock
	_mp.MockENotifier_Prototype.SetMock(mock)
	_mp.MockEResourceListener_Prototype.SetMock(mock)
	_mp.MockEResource_Prototype_Methods.SetMock(mock)
}

// MockEResource_Expecter is the expecter implementation for all EResource methods ( inherited and declared )
type MockEResource_Expecter struct {
	MockENotifier_Expecter
	MockEResourceListener_Expecter
	MockEResource_Expecter_Methods
}

func (_me *MockEResource_Expecter) SetMock(mock *mock.Mock) {
	_me.MockENotifier_Expecter.SetMock(mock)
	_me.MockEResourceListener_Expecter.SetMock(mock)
	_me.MockEResource_Expecter_Methods.SetMock(mock)
}

func (_m *MockEResource_Prototype) EXPECT() *MockEResource_Expecter {
	expecter := &MockEResource_Expecter{}
	expecter.SetMock(_m.mock)
	return expecter
}

// MockEResource_Prototype_Methods is the mock implementation of EResource declared methods
type MockEResource_Prototype_Methods struct {
	mock *mock.Mock
}

func (_mdp *MockEResource_Prototype_Methods) SetMock(mock *mock.Mock) {
	_mdp.mock = mock
}

// MockEResource_Expecter_Methods is the expecter implementation of EResource declared methods
type MockEResource_Expecter_Methods struct {
	mock *mock.Mock
}

func (_mde *MockEResource_Expecter_Methods) SetMock(mock *mock.Mock) {
	_mde.mock = mock
}

// GetAllContents provides a mock function with given fields:
func (_m *MockEResource_Prototype_Methods) GetAllContents() EIterator {
	ret := _m.mock.Called()

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

// MockEResource_GetAllContents_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllContents'
type MockEResource_GetAllContents_Call struct {
	*mock.Call
}

// GetAllContents is a helper method to define mock.On call
func (_e *MockEResource_Expecter_Methods) GetAllContents() *MockEResource_GetAllContents_Call {
	return &MockEResource_GetAllContents_Call{Call: _e.mock.On("GetAllContents")}
}

func (_c *MockEResource_GetAllContents_Call) Run(run func()) *MockEResource_GetAllContents_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockEResource_GetAllContents_Call) Return(_a0 EIterator) *MockEResource_GetAllContents_Call {
	_c.Call.Return(_a0)
	return _c
}

// GetContents provides a mock function with given fields:
func (_m *MockEResource_Prototype_Methods) GetContents() EList {
	ret := _m.mock.Called()

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

// MockEResource_GetContents_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetContents'
type MockEResource_GetContents_Call struct {
	*mock.Call
}

// GetContents is a helper method to define mock.On call
func (_e *MockEResource_Expecter_Methods) GetContents() *MockEResource_GetContents_Call {
	return &MockEResource_GetContents_Call{Call: _e.mock.On("GetContents")}
}

func (_c *MockEResource_GetContents_Call) Run(run func()) *MockEResource_GetContents_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockEResource_GetContents_Call) Return(_a0 EList) *MockEResource_GetContents_Call {
	_c.Call.Return(_a0)
	return _c
}

// GetEObject provides a mock function with given fields: _a0
func (_m *MockEResource_Prototype_Methods) GetEObject(_a0 string) EObject {
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

// MockEResource_GetEObject_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetEObject'
type MockEResource_GetEObject_Call struct {
	*mock.Call
}

// GetEObject is a helper method to define mock.On call
//   - _a0 string
func (_e *MockEResource_Expecter_Methods) GetEObject(_a0 interface{}) *MockEResource_GetEObject_Call {
	return &MockEResource_GetEObject_Call{Call: _e.mock.On("GetEObject", _a0)}
}

func (_c *MockEResource_GetEObject_Call) Run(run func(_a0 string)) *MockEResource_GetEObject_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockEResource_GetEObject_Call) Return(_a0 EObject) *MockEResource_GetEObject_Call {
	_c.Call.Return(_a0)
	return _c
}

// GetErrors provides a mock function with given fields:
func (_m *MockEResource_Prototype_Methods) GetErrors() EList {
	ret := _m.mock.Called()

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

// MockEResource_GetErrors_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetErrors'
type MockEResource_GetErrors_Call struct {
	*mock.Call
}

// GetErrors is a helper method to define mock.On call
func (_e *MockEResource_Expecter_Methods) GetErrors() *MockEResource_GetErrors_Call {
	return &MockEResource_GetErrors_Call{Call: _e.mock.On("GetErrors")}
}

func (_c *MockEResource_GetErrors_Call) Run(run func()) *MockEResource_GetErrors_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockEResource_GetErrors_Call) Return(_a0 EList) *MockEResource_GetErrors_Call {
	_c.Call.Return(_a0)
	return _c
}

// GetObjectIDManager provides a mock function with given fields:
func (_m *MockEResource_Prototype_Methods) GetObjectIDManager() EObjectIDManager {
	ret := _m.mock.Called()

	var r0 EObjectIDManager
	if rf, ok := ret.Get(0).(func() EObjectIDManager); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(EObjectIDManager)
		}
	}

	return r0
}

// MockEResource_GetObjectIDManager_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetObjectIDManager'
type MockEResource_GetObjectIDManager_Call struct {
	*mock.Call
}

// GetObjectIDManager is a helper method to define mock.On call
func (_e *MockEResource_Expecter_Methods) GetObjectIDManager() *MockEResource_GetObjectIDManager_Call {
	return &MockEResource_GetObjectIDManager_Call{Call: _e.mock.On("GetObjectIDManager")}
}

func (_c *MockEResource_GetObjectIDManager_Call) Run(run func()) *MockEResource_GetObjectIDManager_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockEResource_GetObjectIDManager_Call) Return(_a0 EObjectIDManager) *MockEResource_GetObjectIDManager_Call {
	_c.Call.Return(_a0)
	return _c
}

// GetResourceListeners provides a mock function with given fields:
func (_m *MockEResource_Prototype_Methods) GetResourceListeners() EList {
	ret := _m.mock.Called()

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

// MockEResource_GetResourceListeners_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetResourceListeners'
type MockEResource_GetResourceListeners_Call struct {
	*mock.Call
}

// GetResourceListeners is a helper method to define mock.On call
func (_e *MockEResource_Expecter_Methods) GetResourceListeners() *MockEResource_GetResourceListeners_Call {
	return &MockEResource_GetResourceListeners_Call{Call: _e.mock.On("GetResourceListeners")}
}

func (_c *MockEResource_GetResourceListeners_Call) Run(run func()) *MockEResource_GetResourceListeners_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockEResource_GetResourceListeners_Call) Return(_a0 EList) *MockEResource_GetResourceListeners_Call {
	_c.Call.Return(_a0)
	return _c
}

// GetResourceSet provides a mock function with given fields:
func (_m *MockEResource_Prototype_Methods) GetResourceSet() EResourceSet {
	ret := _m.mock.Called()

	var r0 EResourceSet
	if rf, ok := ret.Get(0).(func() EResourceSet); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(EResourceSet)
		}
	}

	return r0
}

// MockEResource_GetResourceSet_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetResourceSet'
type MockEResource_GetResourceSet_Call struct {
	*mock.Call
}

// GetResourceSet is a helper method to define mock.On call
func (_e *MockEResource_Expecter_Methods) GetResourceSet() *MockEResource_GetResourceSet_Call {
	return &MockEResource_GetResourceSet_Call{Call: _e.mock.On("GetResourceSet")}
}

func (_c *MockEResource_GetResourceSet_Call) Run(run func()) *MockEResource_GetResourceSet_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockEResource_GetResourceSet_Call) Return(_a0 EResourceSet) *MockEResource_GetResourceSet_Call {
	_c.Call.Return(_a0)
	return _c
}

// GetURI provides a mock function with given fields:
func (_m *MockEResource_Prototype_Methods) GetURI() *URI {
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

// MockEResource_GetURI_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetURI'
type MockEResource_GetURI_Call struct {
	*mock.Call
}

// GetURI is a helper method to define mock.On call
func (_e *MockEResource_Expecter_Methods) GetURI() *MockEResource_GetURI_Call {
	return &MockEResource_GetURI_Call{Call: _e.mock.On("GetURI")}
}

func (_c *MockEResource_GetURI_Call) Run(run func()) *MockEResource_GetURI_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockEResource_GetURI_Call) Return(_a0 *URI) *MockEResource_GetURI_Call {
	_c.Call.Return(_a0)
	return _c
}

// GetURIFragment provides a mock function with given fields: _a0
func (_m *MockEResource_Prototype_Methods) GetURIFragment(_a0 EObject) string {
	ret := _m.mock.Called(_a0)

	var r0 string
	if rf, ok := ret.Get(0).(func(EObject) string); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockEResource_GetURIFragment_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetURIFragment'
type MockEResource_GetURIFragment_Call struct {
	*mock.Call
}

// GetURIFragment is a helper method to define mock.On call
//   - _a0 EObject
func (_e *MockEResource_Expecter_Methods) GetURIFragment(_a0 interface{}) *MockEResource_GetURIFragment_Call {
	return &MockEResource_GetURIFragment_Call{Call: _e.mock.On("GetURIFragment", _a0)}
}

func (_c *MockEResource_GetURIFragment_Call) Run(run func(_a0 EObject)) *MockEResource_GetURIFragment_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(EObject))
	})
	return _c
}

func (_c *MockEResource_GetURIFragment_Call) Return(_a0 string) *MockEResource_GetURIFragment_Call {
	_c.Call.Return(_a0)
	return _c
}

// GetWarnings provides a mock function with given fields:
func (_m *MockEResource_Prototype_Methods) GetWarnings() EList {
	ret := _m.mock.Called()

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

// MockEResource_GetWarnings_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetWarnings'
type MockEResource_GetWarnings_Call struct {
	*mock.Call
}

// GetWarnings is a helper method to define mock.On call
func (_e *MockEResource_Expecter_Methods) GetWarnings() *MockEResource_GetWarnings_Call {
	return &MockEResource_GetWarnings_Call{Call: _e.mock.On("GetWarnings")}
}

func (_c *MockEResource_GetWarnings_Call) Run(run func()) *MockEResource_GetWarnings_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockEResource_GetWarnings_Call) Return(_a0 EList) *MockEResource_GetWarnings_Call {
	_c.Call.Return(_a0)
	return _c
}

// IsLoaded provides a mock function with given fields:
func (_m *MockEResource_Prototype_Methods) IsLoaded() bool {
	ret := _m.mock.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockEResource_IsLoaded_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsLoaded'
type MockEResource_IsLoaded_Call struct {
	*mock.Call
}

// IsLoaded is a helper method to define mock.On call
func (_e *MockEResource_Expecter_Methods) IsLoaded() *MockEResource_IsLoaded_Call {
	return &MockEResource_IsLoaded_Call{Call: _e.mock.On("IsLoaded")}
}

func (_c *MockEResource_IsLoaded_Call) Run(run func()) *MockEResource_IsLoaded_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockEResource_IsLoaded_Call) Return(_a0 bool) *MockEResource_IsLoaded_Call {
	_c.Call.Return(_a0)
	return _c
}

// IsLoading provides a mock function with given fields:
func (_m *MockEResource_Prototype_Methods) IsLoading() bool {
	ret := _m.mock.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockEResource_IsLoading_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsLoading'
type MockEResource_IsLoading_Call struct {
	*mock.Call
}

// IsLoading is a helper method to define mock.On call
func (_e *MockEResource_Expecter_Methods) IsLoading() *MockEResource_IsLoading_Call {
	return &MockEResource_IsLoading_Call{Call: _e.mock.On("IsLoading")}
}

func (_c *MockEResource_IsLoading_Call) Run(run func()) *MockEResource_IsLoading_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockEResource_IsLoading_Call) Return(_a0 bool) *MockEResource_IsLoading_Call {
	_c.Call.Return(_a0)
	return _c
}

// Load provides a mock function with given fields:
func (_m *MockEResource_Prototype_Methods) Load() {
	_m.mock.Called()
}

// MockEResource_Load_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Load'
type MockEResource_Load_Call struct {
	*mock.Call
}

// Load is a helper method to define mock.On call
func (_e *MockEResource_Expecter_Methods) Load() *MockEResource_Load_Call {
	return &MockEResource_Load_Call{Call: _e.mock.On("Load")}
}

func (_c *MockEResource_Load_Call) Run(run func()) *MockEResource_Load_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockEResource_Load_Call) Return() *MockEResource_Load_Call {
	_c.Call.Return()
	return _c
}

// LoadWithOptions provides a mock function with given fields: options
func (_m *MockEResource_Prototype_Methods) LoadWithOptions(options map[string]interface{}) {
	_m.mock.Called(options)
}

// MockEResource_LoadWithOptions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'LoadWithOptions'
type MockEResource_LoadWithOptions_Call struct {
	*mock.Call
}

// LoadWithOptions is a helper method to define mock.On call
//   - options map[string]interface{}
func (_e *MockEResource_Expecter_Methods) LoadWithOptions(options interface{}) *MockEResource_LoadWithOptions_Call {
	return &MockEResource_LoadWithOptions_Call{Call: _e.mock.On("LoadWithOptions", options)}
}

func (_c *MockEResource_LoadWithOptions_Call) Run(run func(options map[string]interface{})) *MockEResource_LoadWithOptions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(map[string]interface{}))
	})
	return _c
}

func (_c *MockEResource_LoadWithOptions_Call) Return() *MockEResource_LoadWithOptions_Call {
	_c.Call.Return()
	return _c
}

// LoadWithReader provides a mock function with given fields: r, options
func (_m *MockEResource_Prototype_Methods) LoadWithReader(r io.Reader, options map[string]interface{}) {
	_m.mock.Called(r, options)
}

// MockEResource_LoadWithReader_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'LoadWithReader'
type MockEResource_LoadWithReader_Call struct {
	*mock.Call
}

// LoadWithReader is a helper method to define mock.On call
//   - r io.Reader
//   - options map[string]interface{}
func (_e *MockEResource_Expecter_Methods) LoadWithReader(r interface{}, options interface{}) *MockEResource_LoadWithReader_Call {
	return &MockEResource_LoadWithReader_Call{Call: _e.mock.On("LoadWithReader", r, options)}
}

func (_c *MockEResource_LoadWithReader_Call) Run(run func(r io.Reader, options map[string]interface{})) *MockEResource_LoadWithReader_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(io.Reader), args[1].(map[string]interface{}))
	})
	return _c
}

func (_c *MockEResource_LoadWithReader_Call) Return() *MockEResource_LoadWithReader_Call {
	_c.Call.Return()
	return _c
}

// Save provides a mock function with given fields:
func (_m *MockEResource_Prototype_Methods) Save() {
	_m.mock.Called()
}

// MockEResource_Save_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Save'
type MockEResource_Save_Call struct {
	*mock.Call
}

// Save is a helper method to define mock.On call
func (_e *MockEResource_Expecter_Methods) Save() *MockEResource_Save_Call {
	return &MockEResource_Save_Call{Call: _e.mock.On("Save")}
}

func (_c *MockEResource_Save_Call) Run(run func()) *MockEResource_Save_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockEResource_Save_Call) Return() *MockEResource_Save_Call {
	_c.Call.Return()
	return _c
}

// SaveWithOptions provides a mock function with given fields: options
func (_m *MockEResource_Prototype_Methods) SaveWithOptions(options map[string]interface{}) {
	_m.mock.Called(options)
}

// MockEResource_SaveWithOptions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SaveWithOptions'
type MockEResource_SaveWithOptions_Call struct {
	*mock.Call
}

// SaveWithOptions is a helper method to define mock.On call
//   - options map[string]interface{}
func (_e *MockEResource_Expecter_Methods) SaveWithOptions(options interface{}) *MockEResource_SaveWithOptions_Call {
	return &MockEResource_SaveWithOptions_Call{Call: _e.mock.On("SaveWithOptions", options)}
}

func (_c *MockEResource_SaveWithOptions_Call) Run(run func(options map[string]interface{})) *MockEResource_SaveWithOptions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(map[string]interface{}))
	})
	return _c
}

func (_c *MockEResource_SaveWithOptions_Call) Return() *MockEResource_SaveWithOptions_Call {
	_c.Call.Return()
	return _c
}

// SaveWithWriter provides a mock function with given fields: w, options
func (_m *MockEResource_Prototype_Methods) SaveWithWriter(w io.Writer, options map[string]interface{}) {
	_m.mock.Called(w, options)
}

// MockEResource_SaveWithWriter_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SaveWithWriter'
type MockEResource_SaveWithWriter_Call struct {
	*mock.Call
}

// SaveWithWriter is a helper method to define mock.On call
//   - w io.Writer
//   - options map[string]interface{}
func (_e *MockEResource_Expecter_Methods) SaveWithWriter(w interface{}, options interface{}) *MockEResource_SaveWithWriter_Call {
	return &MockEResource_SaveWithWriter_Call{Call: _e.mock.On("SaveWithWriter", w, options)}
}

func (_c *MockEResource_SaveWithWriter_Call) Run(run func(w io.Writer, options map[string]interface{})) *MockEResource_SaveWithWriter_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(io.Writer), args[1].(map[string]interface{}))
	})
	return _c
}

func (_c *MockEResource_SaveWithWriter_Call) Return() *MockEResource_SaveWithWriter_Call {
	_c.Call.Return()
	return _c
}

// SetObjectIDManager provides a mock function with given fields: _a0
func (_m *MockEResource_Prototype_Methods) SetObjectIDManager(_a0 EObjectIDManager) {
	_m.mock.Called(_a0)
}

// MockEResource_SetObjectIDManager_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetObjectIDManager'
type MockEResource_SetObjectIDManager_Call struct {
	*mock.Call
}

// SetObjectIDManager is a helper method to define mock.On call
//   - _a0 EObjectIDManager
func (_e *MockEResource_Expecter_Methods) SetObjectIDManager(_a0 interface{}) *MockEResource_SetObjectIDManager_Call {
	return &MockEResource_SetObjectIDManager_Call{Call: _e.mock.On("SetObjectIDManager", _a0)}
}

func (_c *MockEResource_SetObjectIDManager_Call) Run(run func(_a0 EObjectIDManager)) *MockEResource_SetObjectIDManager_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(EObjectIDManager))
	})
	return _c
}

func (_c *MockEResource_SetObjectIDManager_Call) Return() *MockEResource_SetObjectIDManager_Call {
	_c.Call.Return()
	return _c
}

// SetURI provides a mock function with given fields: _a0
func (_m *MockEResource_Prototype_Methods) SetURI(_a0 *URI) {
	_m.mock.Called(_a0)
}

// MockEResource_SetURI_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetURI'
type MockEResource_SetURI_Call struct {
	*mock.Call
}

// SetURI is a helper method to define mock.On call
//   - _a0 *URI
func (_e *MockEResource_Expecter_Methods) SetURI(_a0 interface{}) *MockEResource_SetURI_Call {
	return &MockEResource_SetURI_Call{Call: _e.mock.On("SetURI", _a0)}
}

func (_c *MockEResource_SetURI_Call) Run(run func(_a0 *URI)) *MockEResource_SetURI_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*URI))
	})
	return _c
}

func (_c *MockEResource_SetURI_Call) Return() *MockEResource_SetURI_Call {
	_c.Call.Return()
	return _c
}

// Unload provides a mock function with given fields:
func (_m *MockEResource_Prototype_Methods) Unload() {
	_m.mock.Called()
}

// MockEResource_Unload_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Unload'
type MockEResource_Unload_Call struct {
	*mock.Call
}

// Unload is a helper method to define mock.On call
func (_e *MockEResource_Expecter_Methods) Unload() *MockEResource_Unload_Call {
	return &MockEResource_Unload_Call{Call: _e.mock.On("Unload")}
}

func (_c *MockEResource_Unload_Call) Run(run func()) *MockEResource_Unload_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockEResource_Unload_Call) Return() *MockEResource_Unload_Call {
	_c.Call.Return()
	return _c
}

type mockConstructorTestingTNewMockEResource interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockEResource creates a new instance of MockEResource_Prototype. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockEResource(t mockConstructorTestingTNewMockEResource) *MockEResource {
	mock := &MockEResource{}
	mock.SetMock(&mock.Mock)
	mock.Mock.Test(t)
	t.Cleanup(func() { mock.AssertExpectations(t) })
	return mock
}
