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
	"iter"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockEList_Add(t *testing.T) {
	l := NewMockEList(t)
	m := NewMockRun(t, 1)
	l.EXPECT().Add(1).Return(true).Run(func(e any) { m.Run(e) }).Once()
	l.EXPECT().Add(1).Call.Return(func(any) bool { return true }).Once()
	assert.True(t, l.Add(1))
	assert.True(t, l.Add(1))
}

func TestMockEList_Remove(t *testing.T) {
	l := NewMockEList(t)
	m := NewMockRun(t, 1)
	l.EXPECT().Remove(1).Return(true).Run(func(e any) { m.Run(e) }).Once()
	l.EXPECT().Remove(1).Call.Return(func(any) bool { return true }).Once()
	assert.True(t, l.Remove(1))
	assert.True(t, l.Remove(1))
}

func TestMockEList_RemoveAt(t *testing.T) {
	l := NewMockEList(t)
	m := NewMockRun(t, 1)
	l.EXPECT().RemoveAt(1).Return(1).Run(func(index int) { m.Run(index) }).Once()
	l.EXPECT().RemoveAt(1).Call.Return(func(int) any { return 2 }).Once()
	assert.Equal(t, 1, l.RemoveAt(1))
	assert.Equal(t, 2, l.RemoveAt(1))
}

func TestMockEList_Insert(t *testing.T) {
	l := NewMockEList(t)
	m := NewMockRun(t, 1, 2)
	l.EXPECT().Insert(1, 2).Return(true).Run(func(index int, e any) { m.Run(index, e) }).Once()
	l.EXPECT().Insert(1, 2).Call.Return(func(int, any) bool { return true }).Once()
	assert.True(t, l.Insert(1, 2))
	assert.True(t, l.Insert(1, 2))
}

func TestMockEList_MoveObject(t *testing.T) {
	l := NewMockEList(t)
	m := NewMockRun(t, 1, 2)
	l.EXPECT().MoveObject(1, 2).Return().Run(func(index int, e interface{}) { m.Run(index, e) }).Once()
	l.MoveObject(1, 2)
}

func TestMockEList_Move(t *testing.T) {
	l := NewMockEList(t)
	m := NewMockRun(t, 1, 2)
	l.EXPECT().Move(1, 2).Return(3).Run(func(oldIndex, newIndex int) { m.Run(oldIndex, newIndex) }).Once()
	l.EXPECT().Move(1, 2).Call.Return(func(int, int) any { return 3 }).Once()
	assert.Equal(t, 3, l.Move(1, 2))
	assert.Equal(t, 3, l.Move(1, 2))
}

func TestMockEList_Get(t *testing.T) {
	l := NewMockEList(t)
	m := NewMockRun(t, 1)
	l.EXPECT().Get(1).Return(1).Run(func(index int) { m.Run(index) }).Once()
	l.EXPECT().Get(1).Call.Return(func(int) any { return 0 }).Once()
	assert.Equal(t, 1, l.Get(1))
	assert.Equal(t, 0, l.Get(1))
}

func TestMockEList_Set(t *testing.T) {
	l := NewMockEList(t)
	m := NewMockRun(t, 1, 2)
	l.EXPECT().Set(1, 2).Return(3).Run(func(index int, e any) { m.Run(index, e) }).Once()
	l.EXPECT().Set(1, 2).Call.Return(func(int, any) any { return 4 }).Once()
	assert.Equal(t, 3, l.Set(1, 2))
	assert.Equal(t, 4, l.Set(1, 2))
}

func TestMockEList_AddAll(t *testing.T) {
	l := NewMockEList(t)
	c := NewMockEList(t)
	m := NewMockRun(t, c)
	l.EXPECT().AddAll(c).Return(true).Run(func(c Collection) { m.Run(c) }).Once()
	l.EXPECT().AddAll(c).Call.Return(func(Collection) bool { return true }).Once()
	assert.True(t, l.AddAll(c))
	assert.True(t, l.AddAll(c))
}

func TestMockEList_InsertAll(t *testing.T) {
	l := NewMockEList(t)
	c := NewMockEList(t)
	m := NewMockRun(t, 0, c)
	l.EXPECT().InsertAll(0, c).Return(true).Run(func(i int, c Collection) { m.Run(i, c) }).Once()
	l.EXPECT().InsertAll(0, c).Call.Return(func(int, Collection) bool { return true }).Once()
	assert.True(t, l.InsertAll(0, c))
	assert.True(t, l.InsertAll(0, c))
}

func TestMockEList_RemoveAll(t *testing.T) {
	l := NewMockEList(t)
	c := NewMockEList(t)
	m := NewMockRun(t, c)
	l.EXPECT().RemoveAll(c).Return(true).Run(func(c Collection) { m.Run(c) }).Once()
	l.EXPECT().RemoveAll(c).Call.Return(func(Collection) bool {
		return true
	}).Once()
	assert.True(t, l.RemoveAll(c))
	assert.True(t, l.RemoveAll(c))
}

func TestMockEList_RemoveRange(t *testing.T) {
	l := NewMockEList(t)
	m := NewMockRun(t, 1, 2)
	l.EXPECT().RemoveRange(1, 2).Return().Run(func(i, j int) { m.Run(i, j) }).Once()
	l.RemoveRange(1, 2)
}

func TestMockEList_Size(t *testing.T) {
	l := NewMockEList(t)
	m := NewMockRun(t)
	l.EXPECT().Size().Return(0).Run(func() { m.Run() }).Once()
	l.EXPECT().Size().Call.Return(func() int { return 1 }).Once()
	assert.Equal(t, 0, l.Size())
	assert.Equal(t, 1, l.Size())
}

func TestMockEList_Clear(t *testing.T) {
	l := NewMockEList(t)
	m := NewMockRun(t)
	l.EXPECT().Clear().Return().Run(func() { m.Run() }).Once()
	l.Clear()
}

func TestMockEList_Empty(t *testing.T) {
	l := NewMockEList(t)
	m := NewMockRun(t)
	l.EXPECT().Empty().Return(true).Run(func() { m.Run() }).Once()
	l.EXPECT().Empty().Call.Return(func() bool { return false }).Once()
	assert.True(t, l.Empty())
	assert.False(t, l.Empty())
}

func TestMockEList_Iterator(t *testing.T) {
	l := NewMockEList(t)
	m := NewMockRun(t)
	it := &MockEIterator{}
	l.EXPECT().Iterator().Return(it).Run(func() { m.Run() }).Once()
	l.EXPECT().Iterator().Call.Return(func() EIterator { return it }).Once()
	assert.Equal(t, it, l.Iterator())
	assert.Equal(t, it, l.Iterator())
}

func TestMockEList_ToArray(t *testing.T) {
	l := NewMockEList(t)
	m := NewMockRun(t)
	r := []any{}
	l.EXPECT().ToArray().Return(r).Run(func() { m.Run() }).Once()
	l.EXPECT().ToArray().Call.Return(func() []any { return r }).Once()
	assert.Equal(t, r, l.ToArray())
	assert.Equal(t, r, l.ToArray())
}

func TestMockEList_All(t *testing.T) {
	l := NewMockEList(t)
	m := NewMockRun(t)
	it := func(yield func(any) bool) {}
	l.EXPECT().All().Return(it).Run(func() { m.Run() }).Once()
	l.EXPECT().All().RunAndReturn(func() iter.Seq[any] { return it }).Once()
	l.EXPECT().All()
	require.NotNil(t, l.All())
	require.NotNil(t, l.All())
	require.Panics(t, func() {
		l.All()
	})
}

func TestMockEList_Contains(t *testing.T) {
	l := NewMockEList(t)
	m := NewMockRun(t, 1)
	l.EXPECT().Contains(1).Return(false).Run(func(e any) { m.Run(e) }).Once()
	l.EXPECT().Contains(2).Call.Return(func(any) bool { return true }).Once()
	assert.False(t, l.Contains(1))
	assert.True(t, l.Contains(2))
}

func TestMockEList_IndexOf(t *testing.T) {
	l := NewMockEList(t)
	m := NewMockRun(t, 1)
	l.EXPECT().IndexOf(1).Return(0).Run(func(e any) { m.Run(e) }).Once()
	l.EXPECT().IndexOf(2).Call.Return(func(any) int { return 1 }).Once()
	assert.Equal(t, 0, l.IndexOf(1))
	assert.Equal(t, 1, l.IndexOf(2))
}
