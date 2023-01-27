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

func TestBasicEMap_Constructor(t *testing.T) {
	m := NewBasicEMap()
	assert.NotNil(t, m)

	var mp EMap = m
	assert.NotNil(t, mp)

	var ml EList = m
	assert.NotNil(t, ml)
}

func TestBasicEMap_Put(t *testing.T) {
	m := NewBasicEMap()
	m.Put(2, "2")
	assert.Equal(t, map[any]any{2: "2"}, m.ToMap())
}

func TestBasicEMap_GetValue(t *testing.T) {
	m := NewBasicEMap()
	assert.Nil(t, m.GetValue(2))

	m.Put(2, "2")
	assert.Equal(t, "2", m.GetValue(2))
}

func TestBasicEMap_RemoveKey(t *testing.T) {
	m := NewBasicEMap()
	m.Put(2, "2")

	assert.Equal(t, "2", m.RemoveKey(2))
	assert.Nil(t, m.GetValue(2))

	assert.Nil(t, m.RemoveKey(2))
}

func TestBasicEMap_PutOverwrite(t *testing.T) {
	m := NewBasicEMap()
	assert.Nil(t, m.GetValue(2))
	m.Put(2, "3")
	assert.Equal(t, "3", m.GetValue(2))

	m.Put(2, "2")
	assert.Equal(t, "2", m.GetValue(2))
	assert.Equalf(t, 1, m.Size(), "Don't store old cell.")
}

func TestBasicEMap_ContainsKey(t *testing.T) {
	m := NewBasicEMap()
	assert.False(t, m.ContainsKey(2))

	m.Put(2, "2")
	assert.True(t, m.ContainsKey(2))

	m.RemoveKey(2)
	assert.False(t, m.ContainsKey(2))
}

func TestBasicEMap_ContainsValue(t *testing.T) {
	m := NewBasicEMap()
	assert.False(t, m.ContainsValue("2"))

	m.Put(2, "2")
	assert.True(t, m.ContainsValue("2"))

	m.RemoveKey(2)
	assert.False(t, m.ContainsValue("2"))
}

func TestBasicEMap_AddEntry(t *testing.T) {
	m := NewBasicEMap()
	mockEntry := NewMockEMapEntry(t)
	m.Add(mockEntry)
	mock.AssertExpectationsForObjects(t, mockEntry)

	mockEntry.EXPECT().GetKey().Once().Return(2)
	mockEntry.EXPECT().GetValue().Once().Return("2")
	assert.Equal(t, map[any]any{2: "2"}, m.ToMap())
	mock.AssertExpectationsForObjects(t, mockEntry)
}

func TestBasicEMap_SetEntry(t *testing.T) {
	m := NewBasicEMap()
	mockEntry := NewMockEMapEntry(t)
	m.Add(mockEntry)

	mockOtherEntry := NewMockEMapEntry(t)
	mockEntry.EXPECT().GetKey().Once().Return(2)
	mockOtherEntry.EXPECT().GetKey().Once().Return(3)
	mockOtherEntry.EXPECT().GetValue().Once().Return("3")
	m.Set(0, mockOtherEntry)
	assert.Equal(t, map[any]any{3: "3"}, m.ToMap())
	mock.AssertExpectationsForObjects(t, mockEntry, mockOtherEntry)
}

func TestBasicEMap_RemoveEntry(t *testing.T) {
	m := NewBasicEMap()
	mockEntry1 := NewMockEMapEntry(t)
	mockEntry1.EXPECT().GetKey().Once().Return(2)
	mockEntry1.EXPECT().GetValue().Once().Return("2")
	mockEntry2 := NewMockEMapEntry(t)
	mockEntry2.EXPECT().GetKey().Once().Return(3)
	mockEntry2.EXPECT().GetValue().Once().Return("3")
	m.Add(mockEntry1)
	m.Add(mockEntry2)
	assert.Equal(t, map[any]any{2: "2", 3: "3"}, m.ToMap())
	mock.AssertExpectationsForObjects(t, mockEntry1, mockEntry2)

	mockEntry1.EXPECT().GetKey().Once().Return(2)
	m.RemoveAt(0)
	assert.Equal(t, map[any]any{3: "3"}, m.ToMap())
	mock.AssertExpectationsForObjects(t, mockEntry1, mockEntry2)
}

func TestBasicEMap_Clear(t *testing.T) {
	m := NewBasicEMap()
	mockEntry1 := NewMockEMapEntry(t)
	mockEntry2 := NewMockEMapEntry(t)
	m.Add(mockEntry1)
	m.Add(mockEntry2)

	m.Clear()
	assert.Equal(t, map[any]any{}, m.ToMap())
	assert.Equal(t, []any{}, m.ToArray())
	mock.AssertExpectationsForObjects(t, mockEntry1, mockEntry2)
}

func TestBasicEMap_UpdateEntry(t *testing.T) {
	m := NewBasicEMap()
	m.Put(2, "2")
	e := m.Get(0).(EMapEntry)
	e.SetKey(3)
	e.SetValue("3")
	assert.Equal(t, map[any]any{3: "3"}, m.ToMap())
	e.SetKey(2)
	e.SetValue("2")
	assert.Equal(t, map[any]any{3: "3"}, m.ToMap())
}
