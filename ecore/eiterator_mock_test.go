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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockEIteratorRun struct {
	mock.Mock
}

func (m *mockEIteratorRun) Run(args ...any) {
	m.Called(args...)
}

type mockConstructorTestingTMockEIteratorRun interface {
	mock.TestingT
	Cleanup(func())
}

// newMockEIteratorRun creates a new instance of MockEIterator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockEIteratorRun(t mockConstructorTestingTMockEIteratorRun, args ...any) *mockEIteratorRun {
	mock := &mockEIteratorRun{}
	mock.Test(t)
	mock.On("Run", args...).Once()
	t.Cleanup(func() { mock.AssertExpectations(t) })
	return mock
}

// TestGetELiterals tests method GetELiterals
func TestMockEIteratorHasNext(t *testing.T) {
	o := NewMockEIterator(t)
	m := newMockEIteratorRun(t)
	o.EXPECT().HasNext().Return(true).Run(func() { m.Run() }).Once()
	o.EXPECT().HasNext().Call.Return(func() bool { return true }).Once()
	assert.True(t, o.HasNext())
	assert.True(t, o.HasNext())
}

func TestMockEIteratorNext(t *testing.T) {
	o := NewMockEIterator(t)
	v := &mock.Mock{}
	m := newMockEIteratorRun(t)
	o.EXPECT().Next().Return(v).Run(func() { m.Run() }).Once()
	o.EXPECT().Next().Call.Return(func() any { return v }).Once()
	assert.Equal(t, v, o.Next())
	assert.Equal(t, v, o.Next())
}
