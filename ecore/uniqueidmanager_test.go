package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUniqueIDManagerRegister(t *testing.T) {
	m := NewUniqueIDManager(20)
	mockObject := &MockEObject{}
	mockOther := &MockEObject{}
	m.Register(mockObject)

	id := m.GetID(mockObject)
	assert.NotNil(t, id)
	assert.Nil(t, m.GetID(mockOther))

	assert.Equal(t, mockObject, m.GetEObject(id))
	assert.Nil(t, m.GetEObject(""))
	assert.Nil(t, m.GetEObject(3))

}

func TestUniqueIDManagerUnRegister(t *testing.T) {
	m := NewUniqueIDManager(20)

	// register object
	mockObject := &MockEObject{}
	m.Register(mockObject)

	id := m.GetID(mockObject)
	assert.NotNil(t, id)

	// unregister
	m.UnRegister(mockObject)
	assert.Nil(t, m.GetID(mockObject))
	assert.Equal(t, id, m.GetDetachedID(mockObject))

	// register again and check it was detached
	m.Register(mockObject)
	assert.Equal(t, id, m.GetID(mockObject))
}
