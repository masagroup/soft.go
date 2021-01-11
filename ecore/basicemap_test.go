package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, map[interface{}]interface{}{2: "2"}, m.ToMap())
	assert.Equal(t, []interface{}{NewMapEntry(2, "2")}, m.ToArray())
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
	m.Add(NewMapEntry(2, "2"))
	assert.Equal(t, map[interface{}]interface{}{2: "2"}, m.ToMap())
	assert.Equal(t, []interface{}{NewMapEntry(2, "2")}, m.ToArray())
}

func TestBasicEMap_SetEntry(t *testing.T) {
	m := NewBasicEMap()
	m.Add(NewMapEntry(2, "2"))
	m.Set(0, NewMapEntry(3, "3"))
	assert.Equal(t, map[interface{}]interface{}{3: "3"}, m.ToMap())
	assert.Equal(t, []interface{}{NewMapEntry(3, "3")}, m.ToArray())
}

func TestBasicEMap_RemoveEntry(t *testing.T) {
	m := NewBasicEMap()
	m.Add(NewMapEntry(2, "2"))
	m.Add(NewMapEntry(3, "3"))
	assert.Equal(t, map[interface{}]interface{}{2: "2", 3: "3"}, m.ToMap())
	assert.Equal(t, []interface{}{NewMapEntry(2, "2"), NewMapEntry(3, "3")}, m.ToArray())

	m.RemoveAt(0)
	assert.Equal(t, map[interface{}]interface{}{3: "3"}, m.ToMap())
	assert.Equal(t, []interface{}{NewMapEntry(3, "3")}, m.ToArray())
}

func TestBasicEMap_Clear(t *testing.T) {
	m := NewBasicEMap()
	m.Add(NewMapEntry(2, "2"))
	m.Add(NewMapEntry(3, "3"))
	m.Clear()
	assert.Equal(t, map[interface{}]interface{}{}, m.ToMap())
	assert.Equal(t, []interface{}{}, m.ToArray())
}
