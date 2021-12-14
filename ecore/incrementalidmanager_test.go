package ecore

import (
	"runtime"
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

	// register again and check it was detached
	m.Register(mockObject)
	assert.Equal(t, id, m.GetID(mockObject))
}

func TestIncrementalIDManagerDetachedObjectWithGC(t *testing.T) {
	m := NewIncrementalIDManager()

	// register object
	mockObject := &MockEObject{}
	m.Register(mockObject)
	id1 := m.GetID(mockObject)
	assert.NotNil(t, id1)

	// unregister object
	m.UnRegister(mockObject)
	assert.Nil(t, m.GetID(mockObject))

	// call GC
	mockObject = nil
	runtime.GC()

	mockObject = &MockEObject{}
	m.Register(mockObject)
	id2 := m.GetID(mockObject)
	assert.NotEqual(t, id1, id2)

}

func TestIncrementalIDManagerSetID(t *testing.T) {
	m := NewIncrementalIDManager()
	mockObject := &MockEObject{}

	assert.Nil(t, m.SetID(mockObject, 2))
	assert.Equal(t, int64(2), m.GetID(mockObject))

	assert.Nil(t, m.SetID(mockObject, nil))
	assert.Equal(t, nil, m.GetID(mockObject))

	assert.Nil(t, m.SetID(mockObject, "2"))
	assert.Equal(t, int64(2), m.GetID(mockObject))

	assert.NotNil(t, m.SetID(mockObject, mockObject))
}

func TestIncrementalIDManagerClear(t *testing.T) {
	m := NewIncrementalIDManager()
	mockObject := &MockEObject{}

	m.SetID(mockObject, 2)
	assert.Equal(t, int64(2), m.GetID(mockObject))

	m.Clear()
	assert.Equal(t, nil, m.GetID(mockObject))

	mockObject = nil
	runtime.GC()
}
