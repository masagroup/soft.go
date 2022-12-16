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
)

func TestMockEList_Add(t *testing.T) {
	l := NewMockEList(t)
	l.On("Add", 1).Once().Return(true)
	l.On("Add", 1).Once().Return(func(any) bool {
		return true
	})
	assert.True(t, l.Add(1))
	assert.True(t, l.Add(1))
}

func TestMockEList_Remove(t *testing.T) {
	l := NewMockEList(t)
	l.On("Remove", 1).Once().Return(true)
	l.On("Remove", 1).Once().Return(func(any) bool {
		return true
	})
	assert.True(t, l.Remove(1))
	assert.True(t, l.Remove(1))
}

func TestMockEList_RemoveAt(t *testing.T) {
	l := NewMockEList(t)
	l.On("RemoveAt", 1).Once().Return(1)
	l.On("RemoveAt", 1).Once().Return(func(int) any {
		return 2
	})
	assert.Equal(t, 1, l.RemoveAt(1))
	assert.Equal(t, 2, l.RemoveAt(1))
}

func TestMockEList_Insert(t *testing.T) {
	l := NewMockEList(t)
	l.On("Insert", 1, 2).Once().Return(true)
	l.On("Insert", 1, 2).Once().Return(func(int, any) bool {
		return true
	})
	assert.True(t, l.Insert(1, 2))
	assert.True(t, l.Insert(1, 2))
}

func TestMockEList_MoveObject(t *testing.T) {
	l := NewMockEList(t)
	l.On("MoveObject", 1, 2).Once()
	l.MoveObject(1, 2)
}

func TestMockEList_Move(t *testing.T) {
	l := NewMockEList(t)
	l.On("Move", 1, 2).Once().Return(3)
	l.On("Move", 1, 2).Once().Return(func(int, int) any {
		return 3
	})
	assert.Equal(t, 3, l.Move(1, 2))
	assert.Equal(t, 3, l.Move(1, 2))
}

func TestMockEList_Get(t *testing.T) {
	l := NewMockEList(t)
	l.On("Get", 1).Once().Return(1)
	l.On("Get", 1).Once().Return(func(int) any {
		return 0
	})
	assert.Equal(t, 1, l.Get(1))
	assert.Equal(t, 0, l.Get(1))
}

func TestMockEList_Set(t *testing.T) {
	l := NewMockEList(t)
	l.On("Set", 1, 2).Once().Return(3)
	l.On("Set", 1, 2).Once().Return(func(int, any) any {
		return 4
	})
	assert.Equal(t, 3, l.Set(1, 2))
	assert.Equal(t, 4, l.Set(1, 2))
}

func TestMockEList_AddAll(t *testing.T) {
	l := NewMockEList(t)
	c := NewMockEList(t)
	l.On("AddAll", c).Once().Return(true)
	l.On("AddAll", c).Once().Return(func(EList) bool {
		return true
	})
	assert.True(t, l.AddAll(c))
	assert.True(t, l.AddAll(c))
}

func TestMockEList_InsertAll(t *testing.T) {
	l := NewMockEList(t)
	c := NewMockEList(t)
	l.On("InsertAll", 0, c).Once().Return(true)
	l.On("InsertAll", 0, c).Once().Return(func(int, EList) bool {
		return true
	})
	assert.True(t, l.InsertAll(0, c))
	assert.True(t, l.InsertAll(0, c))
}

func TestMockEList_RemoveAll(t *testing.T) {
	l := NewMockEList(t)
	c := NewMockEList(t)
	l.On("RemoveAll", c).Once().Return(true)
	l.On("RemoveAll", c).Once().Return(func(EList) bool {
		return true
	})
	assert.True(t, l.RemoveAll(c))
	assert.True(t, l.RemoveAll(c))
}

func TestMockEList_RemoveRange(t *testing.T) {
	l := NewMockEList(t)
	l.On("RemoveRange", 1, 2).Once()
	l.RemoveRange(1, 2)
}

func TestMockEList_Size(t *testing.T) {
	l := NewMockEList(t)
	l.On("Size").Once().Return(0)
	l.On("Size").Once().Return(func() int {
		return 1
	})
	assert.Equal(t, 0, l.Size())
	assert.Equal(t, 1, l.Size())
}

func TestMockEList_Clear(t *testing.T) {
	l := NewMockEList(t)
	l.On("Clear").Once()
	l.Clear()
}

func TestMockEList_Empty(t *testing.T) {
	l := NewMockEList(t)
	l.On("Empty").Once().Return(true)
	l.On("Empty").Once().Return(func() bool {
		return false
	})
	assert.True(t, l.Empty())
	assert.False(t, l.Empty())
}

func TestMockEList_Iterator(t *testing.T) {
	l := NewMockEList(t)
	it := &MockEIterator{}
	l.On("Iterator").Once().Return(it)
	l.On("Iterator").Once().Return(func() EIterator {
		return it
	})
	assert.Equal(t, it, l.Iterator())
	assert.Equal(t, it, l.Iterator())
}

func TestMockEList_ToArray(t *testing.T) {
	l := NewMockEList(t)
	r := []any{}
	l.On("ToArray").Once().Return(r)
	l.On("ToArray").Once().Return(func() []any {
		return r
	})
	assert.Equal(t, r, l.ToArray())
	assert.Equal(t, r, l.ToArray())
}

func TestMockEList_Contains(t *testing.T) {
	l := NewMockEList(t)
	l.On("Contains", 1).Once().Return(false)
	l.On("Contains", 2).Once().Return(func(any) bool {
		return true
	})
	assert.False(t, l.Contains(1))
	assert.True(t, l.Contains(2))
}

func TestMockEList_IndexOf(t *testing.T) {
	l := NewMockEList(t)
	l.On("IndexOf", 1).Once().Return(0)
	l.On("IndexOf", 2).Once().Return(func(any) int {
		return 1
	})
	assert.Equal(t, 0, l.IndexOf(1))
	assert.Equal(t, 1, l.IndexOf(2))
}
