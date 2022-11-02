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

	assert.Equal(t, int64(0), m.GetID(mockObject))
	assert.Nil(t, m.GetID(mockOther))

	assert.Equal(t, mockObject, m.GetEObject(0))
	assert.Nil(t, m.GetEObject(""))
	assert.Nil(t, m.GetEObject(3))

	m.Register(mockOther)
	assert.Equal(t, int64(1), m.GetID(mockOther))
}

func TestIncrementalIDManagerUnRegister(t *testing.T) {
	m := NewIncrementalIDManager()

	// register object
	mockObject := &MockEObject{}
	m.Register(mockObject)
	id := m.GetID(mockObject)
	assert.Equal(t, int64(0), id)
	assert.Nil(t, m.GetDetachedID(mockObject))

	// unregister object
	m.UnRegister(mockObject)
	assert.Nil(t, m.GetID(mockObject))
	assert.Nil(t, m.GetDetachedID(mockObject))

	// register again and check it was detached
	m.Register(mockObject)
	id = m.GetID(mockObject)
	assert.Equal(t, int64(1), id)
	assert.False(t, m.KeepIDs(true))

	// unregister object
	m.UnRegister(mockObject)
	assert.Nil(t, m.GetID(mockObject))
	assert.Equal(t, id, m.GetDetachedID(mockObject))

	// register object
	m.Register(mockObject)
	assert.Equal(t, id, m.GetID(mockObject))
	assert.Nil(t, m.GetDetachedID(mockObject))
}

func TestIncrementalIDManagerSetID(t *testing.T) {
	m := NewIncrementalIDManager()
	mockObject := &MockEObject{}
	mockOther := &MockEObject{}
	assert.Nil(t, m.SetID(mockObject, 2))
	assert.Equal(t, int64(2), m.GetID(mockObject))

	assert.Nil(t, m.SetID(mockObject, nil))
	assert.Equal(t, nil, m.GetID(mockObject))

	assert.Nil(t, m.SetID(mockObject, "2"))
	assert.Equal(t, int64(2), m.GetID(mockObject))

	assert.NotNil(t, m.SetID(mockObject, mockObject))

	m.Register(mockOther)
	assert.Equal(t, int64(3), m.GetID(mockOther))
}

func TestIncrementalIDManagerClear(t *testing.T) {
	m := NewIncrementalIDManager()
	mockObject := &MockEObject{}

	m.SetID(mockObject, 2)
	assert.Equal(t, int64(2), m.GetID(mockObject))

	m.Clear()
	assert.Equal(t, nil, m.GetID(mockObject))
}

func TestUUIDManagerRegister(t *testing.T) {
	m := NewUUIDManager(20)
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

func TestUUIDManagerUnRegister(t *testing.T) {
	m := NewUUIDManager(20)

	// register object
	mockObject := &MockEObject{}
	m.Register(mockObject)

	id := m.GetID(mockObject)
	assert.NotNil(t, id)
	assert.Nil(t, m.GetDetachedID(mockObject))

	// unregister
	assert.False(t, m.KeepIDs(true))
	m.UnRegister(mockObject)
	assert.Nil(t, m.GetID(mockObject))
	assert.Equal(t, id, m.GetDetachedID(mockObject))

	// register again and check it was detached
	m.Register(mockObject)
	assert.Equal(t, id, m.GetID(mockObject))
}

func TestULIDManagerRegister(t *testing.T) {
	m := NewULIDManager()
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

func TestULIDManagerUnRegister(t *testing.T) {
	m := NewULIDManager()

	// register object
	mockObject := &MockEObject{}
	m.Register(mockObject)

	id := m.GetID(mockObject)
	assert.NotNil(t, id)
	assert.Nil(t, m.GetDetachedID(mockObject))

	// unregister
	assert.False(t, m.KeepIDs(true))
	m.UnRegister(mockObject)
	assert.Nil(t, m.GetID(mockObject))
	assert.Equal(t, id, m.GetDetachedID(mockObject))

	// register again and check it was detached
	m.Register(mockObject)
	assert.Equal(t, id, m.GetID(mockObject))
}
