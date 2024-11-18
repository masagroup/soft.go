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

func TestBasicEListGet(t *testing.T) {
	arr := NewBasicEList([]any{3, 5, 7})
	assert.Equal(t, arr.Get(0), 3)
	assert.Equal(t, arr.Get(1), 5)
	assert.Equal(t, arr.Get(2), 7)
	assert.Panics(t, func() { arr.Get(3) })
}

func TestBasicSetInterafaces(t *testing.T) {
	arr := NewEmptyBasicEList()
	arr.SetInterfaces(arr)
	assert.Equal(t, arr, arr.interfaces)
}

func TestBasicEListSet(t *testing.T) {
	arr := NewBasicEList([]any{3, 5, 7})
	arr.Set(0, 4)
	arr.Set(1, 6)
	arr.Set(2, 8)
	assert.Equal(t, []any{4, 6, 8}, arr.ToArray())
	assert.Panics(t, func() { arr.Set(3, 1) })
}

func TestBasicEListSize(t *testing.T) {
	arr := NewBasicEList([]any{3, 5, 7})
	assert.Equal(t, arr.Size(), 3)
}

func TestBasicEListAddAll(t *testing.T) {
	arr := NewBasicEList([]any{3, 5, 7})
	arr2 := NewBasicEList([]any{2})
	assert.True(t, arr2.AddAll(arr))
	assert.Equal(t, []any{2, 3, 5, 7}, arr2.ToArray())
}

func TestBasicEListInsertPrepend(t *testing.T) {
	arr := NewBasicEList([]any{3})
	assert.True(t, arr.Insert(0, 2))
	assert.Equal(t, []any{2, 3}, arr.ToArray())
	assert.Panics(t, func() { arr.Insert(-1, 1) })
}

func TestBasicEListInsertAppend(t *testing.T) {
	arr := NewBasicEList([]any{3})
	assert.True(t, arr.Insert(1, 2))
	assert.Equal(t, []any{3, 2}, arr.ToArray())
	assert.Panics(t, func() { arr.Insert(3, 1) })
}

func TestBasicEListInsertMiddle(t *testing.T) {
	arr := NewBasicEList([]any{3, 5, 7})
	assert.True(t, arr.Insert(1, 2))
	assert.True(t, arr.Insert(2, 3))
	assert.Equal(t, []any{3, 2, 3, 5, 7}, arr.ToArray())
}

func TestBasicEListInsertAll(t *testing.T) {
	arr := NewBasicEList([]any{3, 5, 7})
	arr2 := NewBasicEList([]any{-3, -5, -7})
	arr2.InsertAll(1, arr)
	assert.Equal(t, []any{-3, 3, 5, 7, -5, -7}, arr2.ToArray())
	assert.Panics(t, func() { arr.InsertAll(-1, arr) })
	assert.Panics(t, func() { arr.InsertAll(6, arr) })
}

func TestBasicEListRemoveAll(t *testing.T) {
	arr := NewBasicEList([]any{3, 5, 7})
	arr2 := NewBasicEList([]any{3, 5})
	assert.True(t, arr.RemoveAll(arr2))
	assert.Equal(t, []any{7}, arr.ToArray())
}

func TestBasicEListMoveObjectBefore(t *testing.T) {
	arr := NewBasicEList([]any{2, 4, 6, 8, 10})
	arr.MoveObject(3, 4)
	assert.Equal(t, []any{2, 6, 8, 4, 10}, arr.ToArray())
}

func TestBasicEListMoveObjectAfter(t *testing.T) {
	arr := NewBasicEList([]any{2, 4, 6, 8, 10})
	arr.MoveObject(4, 4)
	assert.Equal(t, []any{2, 6, 8, 10, 4}, arr.ToArray())
}

func TestBasicEListMoveObjectEnd(t *testing.T) {
	arr := NewBasicEList([]any{2, 4, 6, 8, 10})
	arr.MoveObject(0, 4)
	assert.Equal(t, []any{4, 2, 6, 8, 10}, arr.ToArray())
}

func TestBasicEListMoveObjectSame(t *testing.T) {
	arr := NewBasicEList([]any{2, 4, 6, 8, 10})
	arr.MoveObject(1, 4)
	assert.Equal(t, []any{2, 4, 6, 8, 10}, arr.ToArray())
}

func TestBasicEListMoveObjectInvalid(t *testing.T) {
	arr := NewBasicEList([]any{2, 4, 6, 8, 10})
	assert.Panics(t, func() { arr.MoveObject(1, 3) })
}

func TestBasicEListMoveInvalid(t *testing.T) {
	arr := NewBasicEList([]any{2, 4, 6, 8, 10})
	assert.Panics(t, func() { arr.Move(1, 7) })
}

func TestBasicEListMoveBorders(t *testing.T) {
	arr := NewBasicEList([]any{2, 4})
	arr.Move(0, 1)
	assert.Equal(t, []any{4, 2}, arr.ToArray())
}

func TestBasicEListMoveBordersInverse(t *testing.T) {
	arr := NewBasicEList([]any{2, 4})
	arr.Move(1, 0)
	assert.Equal(t, []any{4, 2}, arr.ToArray())
}

func TestBasicEListMoveComplex(t *testing.T) {
	arr := NewBasicEList([]any{2, 4, 6, 8, 10})
	arr.Move(0, 3)
	assert.Equal(t, []any{4, 6, 8, 2, 10}, arr.ToArray())
}

func TestBasicEListMoveComplexEnd(t *testing.T) {
	arr := NewBasicEList([]any{2, 4, 6, 8, 10})
	arr.Move(0, 4)
	assert.Equal(t, []any{4, 6, 8, 10, 2}, arr.ToArray())
}

func TestBasicEListMoveComplexInverse(t *testing.T) {
	arr := NewBasicEList([]any{2, 4, 6, 8, 10})
	arr.Move(3, 1)
	assert.Equal(t, []any{2, 8, 4, 6, 10}, arr.ToArray())
}

func TestBasicEListRemoveBegin(t *testing.T) {
	arr := NewBasicEList([]any{2, 4, 6})
	assert.Equal(t, arr.Size(), 3)
	arr.Remove(2)
	assert.Equal(t, []any{4, 6}, arr.ToArray())
}

func TestBasicEListRemoveEnd(t *testing.T) {
	arr := NewBasicEList([]any{2, 4, 6})
	assert.Equal(t, arr.Size(), 3)
	arr.Remove(6)
	assert.Equal(t, []any{2, 4}, arr.ToArray())
}

func TestBasicEListRemoveMiddle(t *testing.T) {
	arr := NewBasicEList([]any{2, 4, 6, 6})
	assert.Equal(t, arr.Size(), 4)
	arr.Remove(4)
	arr.Remove(6)
	assert.Equal(t, []any{2, 6}, arr.ToArray())
}

func TestBasicEListRemoveInvalid(t *testing.T) {
	arr := NewBasicEList([]any{2, 4, 6, 8, 10})
	assert.False(t, arr.Remove(7))
}

func TestBasicEListRemoveAtInvalid(t *testing.T) {
	arr := NewBasicEList([]any{2, 4, 6, 8, 10})
	assert.Panics(t, func() { arr.RemoveAt(7) })
}

func TestBasicEListRemoveRangeInvalid(t *testing.T) {
	arr := NewBasicEList([]any{2, 4, 6, 8, 10})
	assert.Panics(t, func() { arr.RemoveRange(8, 3) })
	assert.Panics(t, func() { arr.RemoveRange(3, 8) })
	assert.Panics(t, func() { arr.RemoveRange(3, 2) })
}

func TestBasicEListRemoveRangeStart(t *testing.T) {
	arr := NewBasicEList([]any{2, 4, 6, 8, 10})
	arr.RemoveRange(0, 2)
	assert.Equal(t, []any{6, 8, 10}, arr.ToArray())
}

func TestBasicEListRemoveRangeMiddle(t *testing.T) {
	arr := NewBasicEList([]any{2, 4, 6, 8, 10})
	arr.RemoveRange(1, 3)
	assert.Equal(t, []any{2, 8, 10}, arr.ToArray())
}

func TestBasicEListRemoveRangeEnd(t *testing.T) {
	arr := NewBasicEList([]any{2, 4, 6, 8, 10})
	arr.RemoveRange(3, 5)
	assert.Equal(t, []any{2, 4, 6}, arr.ToArray())
}

func TestBasicEListClear(t *testing.T) {
	arr := NewBasicEList([]any{3, 5})
	assert.Equal(t, arr.Size(), 2)
	arr.Clear()
	assert.True(t, arr.Empty())
}

func TestBasicEListEmptyTrue(t *testing.T) {
	arr := NewBasicEList([]any{3, 5})
	assert.False(t, arr.Empty())
}

func TestBasicEListEmptyFalse(t *testing.T) {
	arr := NewBasicEList([]any{})
	assert.True(t, arr.Empty())
}

func TestBasicEListContainsFalse(t *testing.T) {
	arr := NewBasicEList([]any{2})
	assert.False(t, arr.Contains(4))
}

func TestBasicEListContainsTrue(t *testing.T) {
	arr := NewBasicEList([]any{2})
	assert.True(t, arr.Contains(2))
}

func TestBasicEListIterate(t *testing.T) {
	arr := NewBasicEList([]any{0, 2, 4})
	i := 0
	for it := arr.Iterator(); it.HasNext(); {
		assert.Equal(t, it.Next(), i)
		i += 2
	}
	assert.Equal(t, 6, i)
}

func TestBasicEListAddAllUnique(t *testing.T) {
	arr := NewUniqueBasicEList([]any{3, 5, 7, 5})
	arr2 := NewUniqueBasicEList([]any{2, 7})
	arr2.AddAll(arr)
	assert.Equal(t, []any{2, 7, 3, 5}, arr2.ToArray())
	arr3 := NewBasicEList(nil)
	assert.False(t, arr2.AddAll(arr3))
}

func TestBasicEListAddUnique(t *testing.T) {
	arr := NewUniqueBasicEList([]any{})
	assert.True(t, arr.Add(2))
	assert.True(t, arr.Add(5))
	assert.False(t, arr.Add(2))
}

func TestBasicEListInsertUnique(t *testing.T) {
	arr := NewUniqueBasicEList([]any{3, 5, 7})
	assert.Equal(t, arr.Insert(1, 2), true)
	assert.Equal(t, arr.Insert(2, 9), true)
	assert.Equal(t, arr.Insert(2, 3), false)
	assert.Equal(t, []any{3, 2, 9, 5, 7}, arr.ToArray())
}

func TestBasicEListInsertAllUnique(t *testing.T) {
	arr := NewUniqueBasicEList([]any{3, 5, -5, 7})
	arr2 := NewUniqueBasicEList([]any{-3, -5, -7})
	arr2.InsertAll(1, arr)
	assert.Equal(t, arr2.ToArray(), []any{-3, 3, 5, 7, -5, -7})

	arr3 := NewBasicEList(nil)
	assert.False(t, arr2.InsertAll(1, arr3))
}

func TestBasicEListSetUnique(t *testing.T) {
	arr := NewUniqueBasicEList([]any{1, 2})
	assert.Panics(t, func() {
		arr.Set(0, 2)
	})
}
