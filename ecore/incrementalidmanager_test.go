package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIncrementalIDManagerRegister(t *testing.T) {
	m := NewIncrementalIDManager()
	mockObject := &MockEObject{}
	mockOther := &MockEObject{}
	m.Register(mockObject)

	assert.Equal(t, 0, m.GetID(mockObject))
	assert.Nil(t, m.GetID(mockOther))

	assert.Equal(t, mockObject, m.GetEObject(0))
	assert.Nil(t, m.GetEObject(""))
	assert.Nil(t, m.GetEObject(3))
}

func TestIncrementalIDManagerUnRegister(t *testing.T) {
	m := NewIncrementalIDManager()

	// register object
	mockObject := &MockEObject{}
	m.Register(mockObject)
	id := m.GetID(mockObject)
	assert.NotNil(t, id)

	// unregister object
	m.UnRegister(mockObject)
	assert.Nil(t, m.GetID(mockObject))
	assert.Equal(t, id, m.GetDetachedID(mockObject))

	// register again and check it was detached
	m.Register(mockObject)
	assert.Equal(t, id, m.GetID(mockObject))
	assert.Nil(t, m.GetDetachedID(mockObject))
}

func TestIncrementalIDManagerSetID(t *testing.T) {
	m := NewIncrementalIDManager()
	mockObject := &MockEObject{}

	m.SetID(mockObject, 2)
	assert.Equal(t, 2, m.GetID(mockObject))

	m.SetID(mockObject, nil)
	assert.Equal(t, nil, m.GetID(mockObject))

	m.SetID(mockObject, "2")
	assert.Equal(t, 2, m.GetID(mockObject))
}

func TestIncrementalIDManagerClear(t *testing.T) {
	m := NewIncrementalIDManager()
	mockObject := &MockEObject{}

	m.SetID(mockObject, 2)
	assert.Equal(t, 2, m.GetID(mockObject))

	m.Clear()
	assert.Equal(t, nil, m.GetID(mockObject))
}
