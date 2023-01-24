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

type mockEListRun struct {
	mock.Mock
}

func (m *mockEListRun) Run(args ...any) {
	m.Called(args...)
}

type mockConstructorTestingTMockEListRun interface {
	mock.TestingT
	Cleanup(func())
}

// newMockEListRun creates a new instance of MockEList. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockEListRun(t mockConstructorTestingTMockEListRun, args ...any) *mockEListRun {
	mock := &mockEListRun{}
	mock.Test(t)
	mock.On("Run", args...).Once()
	t.Cleanup(func() { mock.AssertExpectations(t) })
	return mock
}

func TestMockEList_Add(t *testing.T) {
	l := NewMockEList(t)
	m := newMockEListRun(t, 1)
	l.EXPECT().Add(1).Run(func(e any) { m.Run(e) }).Return(true).Once()
	l.EXPECT().Add(1).Once().Return(func(any) bool { return true })
	assert.True(t, l.Add(1))
	assert.True(t, l.Add(1))
}

func TestMockEList_Remove(t *testing.T) {
	l := NewMockEList(t)
	m := newMockEListRun(t, 1)
	l.EXPECT().Remove(1).Run(func(e any) { m.Run(e) }).Return(true).Once()
	l.EXPECT().Remove(1).Once().Return(func(any) bool { return true })
	assert.True(t, l.Remove(1))
	assert.True(t, l.Remove(1))
}

func TestMockEList_RemoveAt(t *testing.T) {
	l := NewMockEList(t)
	m := newMockEListRun(t, 1)
	l.EXPECT().RemoveAt(1).Run(func(index int) { m.Run(index) }).Return(1).Once()
	l.EXPECT().RemoveAt(1).Once().Return(func(int) any { return 2 })
	assert.Equal(t, 1, l.RemoveAt(1))
	assert.Equal(t, 2, l.RemoveAt(1))
}

func TestMockEList_Insert(t *testing.T) {
	l := NewMockEList(t)
	m := newMockEListRun(t, 1, 2)
	l.EXPECT().Insert(1, 2).Run(func(index int, e any) { m.Run(index, e) }).Once().Return(true)
	l.EXPECT().Insert(1, 2).Once().Return(func(int, any) bool { return true })
	assert.True(t, l.Insert(1, 2))
	assert.True(t, l.Insert(1, 2))
}

func TestMockEList_MoveObject(t *testing.T) {
	l := NewMockEList(t)
	m := newMockEListRun(t, 1, 2)
	l.EXPECT().MoveObject(1, 2).Run(func(index int, e interface{}) { m.Run(index, e) }).Once()
	l.MoveObject(1, 2)
}

func TestMockEList_Move(t *testing.T) {
	l := NewMockEList(t)
	m := newMockEListRun(t, 1, 2)
	l.EXPECT().Move(1, 2).Run(func(oldIndex, newIndex int) { m.Run(oldIndex, newIndex) }).Once().Return(3)
	l.EXPECT().Move(1, 2).Once().Return(func(int, int) any { return 3 })
	assert.Equal(t, 3, l.Move(1, 2))
	assert.Equal(t, 3, l.Move(1, 2))
}

func TestMockEList_Get(t *testing.T) {
	l := NewMockEList(t)
	m := newMockEListRun(t, 1)
	l.EXPECT().Get(1).Run(func(index int) { m.Run(index) }).Once().Return(1)
	l.EXPECT().Get(1).Once().Return(func(int) any { return 0 })
	assert.Equal(t, 1, l.Get(1))
	assert.Equal(t, 0, l.Get(1))
}

func TestMockEList_Set(t *testing.T) {
	l := NewMockEList(t)
	m := newMockEListRun(t, 1, 2)
	l.EXPECT().Set(1, 2).Run(func(index int, e any) { m.Run(index, e) }).Once().Return(3)
	l.EXPECT().Set(1, 2).Once().Return(func(int, any) any { return 4 })
	assert.Equal(t, 3, l.Set(1, 2))
	assert.Equal(t, 4, l.Set(1, 2))
}

func TestMockEList_AddAll(t *testing.T) {
	l := NewMockEList(t)
	c := NewMockEList(t)
	m := newMockEListRun(t, c)
	l.EXPECT().AddAll(c).Run(func(c EList) { m.Run(c) }).Once().Return(true)
	l.EXPECT().AddAll(c).Once().Return(func(EList) bool { return true })
	assert.True(t, l.AddAll(c))
	assert.True(t, l.AddAll(c))
}

func TestMockEList_InsertAll(t *testing.T) {
	l := NewMockEList(t)
	c := NewMockEList(t)
	m := newMockEListRun(t, 0, c)
	l.EXPECT().InsertAll(0, c).Run(func(i int, c EList) { m.Run(i, c) }).Once().Return(true)
	l.EXPECT().InsertAll(0, c).Once().Return(func(int, EList) bool { return true })
	assert.True(t, l.InsertAll(0, c))
	assert.True(t, l.InsertAll(0, c))
}

func TestMockEList_RemoveAll(t *testing.T) {
	l := NewMockEList(t)
	c := NewMockEList(t)
	m := newMockEListRun(t, c)
	l.EXPECT().RemoveAll(c).Run(func(c EList) { m.Run(c) }).Once().Return(true)
	l.EXPECT().RemoveAll(c).Once().Return(func(EList) bool {
		return true
	})
	assert.True(t, l.RemoveAll(c))
	assert.True(t, l.RemoveAll(c))
}

func TestMockEList_RemoveRange(t *testing.T) {
	l := NewMockEList(t)
	m := newMockEListRun(t, 1, 2)
	l.EXPECT().RemoveRange(1, 2).Run(func(i, j int) { m.Run(i, j) }).Once()
	l.RemoveRange(1, 2)
}

func TestMockEList_Size(t *testing.T) {
	l := NewMockEList(t)
	m := newMockEListRun(t)
	l.EXPECT().Size().Run(func() { m.Run() }).Once().Return(0)
	l.EXPECT().Size().Once().Return(func() int { return 1 })
	assert.Equal(t, 0, l.Size())
	assert.Equal(t, 1, l.Size())
}

func TestMockEList_Clear(t *testing.T) {
	l := NewMockEList(t)
	m := newMockEListRun(t)
	l.EXPECT().Clear().Run(func() { m.Run() }).Once()
	l.Clear()
}

func TestMockEList_Empty(t *testing.T) {
	l := NewMockEList(t)
	m := newMockEListRun(t)
	l.EXPECT().Empty().Run(func() { m.Run() }).Once().Return(true)
	l.EXPECT().Empty().Once().Return(func() bool { return false })
	assert.True(t, l.Empty())
	assert.False(t, l.Empty())
}

func TestMockEList_Iterator(t *testing.T) {
	l := NewMockEList(t)
	m := newMockEListRun(t)
	it := &MockEIterator{}
	l.EXPECT().Iterator().Run(func() { m.Run() }).Once().Return(it)
	l.EXPECT().Iterator().Once().Return(func() EIterator { return it })
	assert.Equal(t, it, l.Iterator())
	assert.Equal(t, it, l.Iterator())
}

func TestMockEList_ToArray(t *testing.T) {
	l := NewMockEList(t)
	m := newMockEListRun(t)
	r := []any{}
	l.EXPECT().ToArray().Run(func() { m.Run() }).Once().Return(r)
	l.EXPECT().ToArray().Once().Return(func() []any { return r })
	assert.Equal(t, r, l.ToArray())
	assert.Equal(t, r, l.ToArray())
}

func TestMockEList_Contains(t *testing.T) {
	l := NewMockEList(t)
	m := newMockEListRun(t, 1)
	l.EXPECT().Contains(1).Run(func(e any) { m.Run(e) }).Once().Return(false)
	l.EXPECT().Contains(2).Once().Return(func(any) bool { return true })
	assert.False(t, l.Contains(1))
	assert.True(t, l.Contains(2))
}

func TestMockEList_IndexOf(t *testing.T) {
	l := NewMockEList(t)
	m := newMockEListRun(t, 1)
	l.EXPECT().IndexOf(1).Run(func(e any) { m.Run(e) }).Once().Return(0)
	l.EXPECT().IndexOf(2).Once().Return(func(any) int { return 1 })
	assert.Equal(t, 0, l.IndexOf(1))
	assert.Equal(t, 1, l.IndexOf(2))
}
