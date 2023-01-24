// Code generated by soft.generator.go. DO NOT EDIT.

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
	"github.com/stretchr/testify/mock"
)

type MockETypeParameter struct {
	MockENamedElement
}

type MockETypeParameter_Expecter struct {
	MockENamedElement_Expecter
}

func (eTypeParameter *MockETypeParameter) EXPECT() *MockETypeParameter_Expecter {
	e := &MockETypeParameter_Expecter{}
	e.Mock = &eTypeParameter.Mock
	return e
}

// GetEBounds get the value of eBounds
func (eTypeParameter *MockETypeParameter) GetEBounds() EList {
	ret := eTypeParameter.Called()

	var r EList
	if rf, ok := ret.Get(0).(func() EList); ok {
		r = rf()
	} else {
		if ret.Get(0) != nil {
			r = ret.Get(0).(EList)
		}
	}

	return r
}

type MockETypeParameter_GetEBounds_Call struct {
	*mock.Call
}

func (e *MockETypeParameter_Expecter) GetEBounds() *MockETypeParameter_GetEBounds_Call {
	return &MockETypeParameter_GetEBounds_Call{Call: e.Mock.On("GetEBounds")}
}

func (c *MockETypeParameter_GetEBounds_Call) Run(run func()) *MockETypeParameter_GetEBounds_Call {
	c.Call.Run(func(mock.Arguments) {
		run()
	})
	return c
}

func (c *MockETypeParameter_GetEBounds_Call) Return(eBounds EList) *MockETypeParameter_GetEBounds_Call {
	c.Call.Return(eBounds)
	return c
}

type mockConstructorTestingTNewMockETypeParameter interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockETypeParameter creates a new instance of MockETypeParameter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockETypeParameter(t mockConstructorTestingTNewMockETypeParameter) *MockETypeParameter {
	mock := &MockETypeParameter{}
	mock.Mock.Test(t)
	t.Cleanup(func() { mock.AssertExpectations(t) })
	return mock
}
