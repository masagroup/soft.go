package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompactEObjectConstructor(t *testing.T) {
	o := &CompactEObjectImpl{}
	o.Initialize()
	assert.True(t, o.EDeliver())
}

func TestCompactEObject_EDeliver(t *testing.T) {
	o := &CompactEObjectImpl{}
	assert.False(t, o.EDeliver())

	o.ESetDeliver(true)
	assert.True(t, o.EDeliver())

	o.ESetDeliver(false)
	assert.False(t, o.EDeliver())
}

func TestCompactEObject_NoFields(t *testing.T) {
	o := &CompactEObjectImpl{}
	assert.False(t, o.hasField(container_flag))
	assert.Nil(t, o.getField(container_flag))
}

func TestCompactEObject_OneField(t *testing.T) {
	c := &MockEObject{}
	o := &CompactEObjectImpl{}
	o.setField(container_flag, c)
	assert.True(t, o.hasField(container_flag))
	assert.Equal(t, c, o.getField(container_flag))

	c2 := &MockEObject{}
	o.setField(container_flag, c2)
	assert.True(t, o.hasField(container_flag))
	assert.Equal(t, c2, o.getField(container_flag))
}

func TestCompactEObject_TwoFields_Begin(t *testing.T) {
	c := &MockEObject{}
	r := &MockEResource{}
	o := &CompactEObjectImpl{}
	o.setField(resource_flag, r)
	o.setField(container_flag, c)
	assert.True(t, o.hasField(container_flag))
	assert.True(t, o.hasField(resource_flag))
	assert.False(t, o.hasField(properties_flag))
	assert.Equal(t, c, o.getField(container_flag))
	assert.Equal(t, r, o.getField(resource_flag))

	c2 := &MockEObject{}
	o.setField(container_flag, c2)
	assert.True(t, o.hasField(container_flag))
	assert.Equal(t, c2, o.getField(container_flag))
}

func TestCompactEObject_TwoFields_End(t *testing.T) {
	c := &MockEObject{}
	r := &MockEResource{}
	o := &CompactEObjectImpl{}
	o.setField(container_flag, c)
	o.setField(resource_flag, r)
	assert.True(t, o.hasField(container_flag))
	assert.True(t, o.hasField(resource_flag))
	assert.False(t, o.hasField(properties_flag))
	assert.Equal(t, c, o.getField(container_flag))
	assert.Equal(t, r, o.getField(resource_flag))
}

func TestCompactEObject_ManyFields(t *testing.T) {
	c := &MockEObject{}
	r := &MockEResource{}
	p := []interface{}{}
	o := &CompactEObjectImpl{}
	o.setField(container_flag, c)
	o.setField(properties_flag, p)
	o.setField(resource_flag, r)
	assert.True(t, o.hasField(container_flag))
	assert.True(t, o.hasField(resource_flag))
	assert.True(t, o.hasField(properties_flag))
	assert.Equal(t, c, o.getField(container_flag))
	assert.Equal(t, r, o.getField(resource_flag))
	assert.Equal(t, p, o.getField(properties_flag))
}

func TestCompactEObject_RemoveNoField(t *testing.T) {
	o := &CompactEObjectImpl{}
	o.setField(container_flag, nil)
}

func TestCompactEObject_RemoveOneField(t *testing.T) {
	c := &MockEObject{}
	o := &CompactEObjectImpl{}
	o.setField(container_flag, c)
	assert.True(t, o.hasField(container_flag))
	o.setField(container_flag, nil)
	assert.False(t, o.hasField(container_flag))
}

func TestCompactEObject_RemoveManyFields(t *testing.T) {
	c := &MockEObject{}
	r := &MockEResource{}
	p := []interface{}{}
	o := &CompactEObjectImpl{}
	o.setField(container_flag, c)
	o.setField(properties_flag, p)
	o.setField(resource_flag, r)
	o.setField(container_flag, nil)
	o.setField(resource_flag, nil)
	assert.False(t, o.hasField(container_flag))
	assert.False(t, o.hasField(resource_flag))
	assert.True(t, o.hasField(properties_flag))
	assert.Equal(t, p, o.getField(properties_flag))

	o.setField(resource_flag, r)
	o.setField(properties_flag, nil)
	assert.False(t, o.hasField(properties_flag))
	assert.Nil(t, o.getField(properties_flag))
}
