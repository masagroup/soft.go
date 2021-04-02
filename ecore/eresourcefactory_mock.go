// Code generated by mockery v1.0.0. DO NOT EDIT.

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
	mock "github.com/stretchr/testify/mock"
)

// MockEResourceFactory is an autogenerated mock type for the EResourceFactory type
type MockEResourceFactory struct {
	mock.Mock
}

// CreateResource provides a mock function with given fields: uri
func (_m *MockEResourceFactory) CreateResource(uri *URI) EResource {
	ret := _m.Called(uri)

	var r0 EResource
	if rf, ok := ret.Get(0).(func(*URI) EResource); ok {
		r0 = rf(uri)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(EResource)
		}
	}

	return r0
}
