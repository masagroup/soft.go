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

type mockEMapRun struct {
	mock.Mock
}

func (m *mockEMapRun) Run(args ...any) {
	m.Called(args...)
}

type mockConstructorTestingTMockEMapRun interface {
	mock.TestingT
	Cleanup(func())
}

// newMockEMapRun creates a new instance of MockEMap. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockEMapRun(t mockConstructorTestingTMockEMapRun, args ...any) *mockEMapRun {
	mock := &mockEMapRun{}
	mock.Test(t)
	mock.On("Run", args...).Once()
	t.Cleanup(func() { mock.AssertExpectations(t) })
	return mock
}

func TestMockEMap_ContainsKey(t *testing.T) {
	l := NewMockEMap(t)
	m := newMockEMapRun(t, 1)
	l.EXPECT().ContainsKey(1).Return(true).Run(func(key any) { m.Run(key) }).Once()
	l.EXPECT().ContainsKey(2).Call.Return(func(any) bool {
		return true
	}).Once()
	assert.True(t, l.ContainsKey(1))
	assert.True(t, l.ContainsKey(2))
}

func TestMockEMap_ContainsValue(t *testing.T) {
	l := NewMockEMap(t)
	m := newMockEMapRun(t, 1)
	l.EXPECT().ContainsValue(1).Return(true).Run(func(key any) { m.Run(key) }).Once()
	l.EXPECT().ContainsValue(2).Call.Return(func(any) bool {
		return true
	}).Once()
	assert.True(t, l.ContainsValue(1))
	assert.True(t, l.ContainsValue(2))
}

func TestMockEMap_RemoveKey(t *testing.T) {
	l := NewMockEMap(t)
	m := newMockEMapRun(t, 1)
	l.EXPECT().RemoveKey(1).Return("1").Run(func(key any) { m.Run(key) }).Once()
	l.EXPECT().RemoveKey(2).Call.Return(func(any) any {
		return "2"
	}).Once()
	assert.Equal(t, "1", l.RemoveKey(1))
	assert.Equal(t, "2", l.RemoveKey(2))
}

func TestMockEMap_Put(t *testing.T) {
	l := NewMockEMap(t)
	m := newMockEMapRun(t, 1, "1")
	l.EXPECT().Put(1, "1").Return().Run(func(key any, value any) { m.Run(key, value) }).Once()
	l.Put(1, "1")
}

func TestMockEMap_GetValue(t *testing.T) {
	l := NewMockEMap(t)
	m := newMockEMapRun(t, 1)
	l.EXPECT().GetValue(1).Return("1").Run(func(key any) { m.Run(key) }).Once()
	l.EXPECT().GetValue(2).Call.Return(func(any) any {
		return "2"
	}).Once()
	assert.Equal(t, "1", l.GetValue(1))
	assert.Equal(t, "2", l.GetValue(2))
}

func TestMockEMap_ToMap(t *testing.T) {
	l := NewMockEMap(t)
	m := map[any]any{}
	mr := newMockEMapRun(t)
	l.EXPECT().ToMap().Return(m).Run(func() { mr.Run() }).Once()
	l.EXPECT().ToMap().Call.Return(func() map[any]any {
		return m
	}).Once()
	assert.Equal(t, m, l.ToMap())
	assert.Equal(t, m, l.ToMap())
}
