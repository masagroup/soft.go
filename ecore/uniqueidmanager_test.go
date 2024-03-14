package ecore

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIncrementalIDManagerRegister(t *testing.T) {
	m := NewIncrementalIDManager()
	mockObject := NewMockEObject(t)
	mockOther := NewMockEObject(t)
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
	mockObject := NewMockEObject(t)
	m.Register(mockObject)
	id := m.GetID(mockObject)
	assert.Equal(t, int64(0), id)
	assert.Nil(t, m.GetDetachedID(mockObject))

	// unregister object
	m.UnRegister(mockObject)
	assert.Nil(t, m.GetID(mockObject))
	assert.Equal(t, id, m.GetDetachedID(mockObject))

	// register again and check it was detached
	m.Register(mockObject)
	assert.Equal(t, id, m.GetID(mockObject))
}

func TestIncrementalIDManagerSetID(t *testing.T) {
	m := NewIncrementalIDManager()
	mockObject := NewMockEObject(t)

	for _, test := range []struct {
		inputID     any
		expectedErr error
		expectedID  any
	}{
		{nil, nil, nil},
		{int(1), nil, int64(1)},
		{int8(1), nil, int64(1)},
		{int16(1), nil, int64(1)},
		{int32(1), nil, int64(1)},
		{int64(1), nil, int64(1)},
		{uint(1), nil, int64(1)},
		{uint8(1), nil, int64(1)},
		{uint16(1), nil, int64(1)},
		{uint32(1), nil, int64(1)},
		{uint64(1), nil, int64(1)},
		{"1", nil, int64(1)},
		{true, errors.New("id:'true' not supported by IncrementalIDManager"), int64(1)},
	} {
		assert.Equal(t, test.expectedErr, m.SetID(mockObject, test.inputID))
		assert.Equal(t, test.expectedID, m.GetID(mockObject))
	}

	mockOther := NewMockEObject(t)
	m.Register(mockOther)
	assert.Equal(t, int64(2), m.GetID(mockOther))
}

func TestIncrementalIDManagerClear(t *testing.T) {
	m := NewIncrementalIDManager()
	mockObject := NewMockEObject(t)

	err := m.SetID(mockObject, 2)
	assert.Nil(t, err)
	assert.Equal(t, int64(2), m.GetID(mockObject))

	m.Clear()
	assert.Equal(t, nil, m.GetID(mockObject))
}

func TestUUIDManagerRegister(t *testing.T) {
	m := NewUUIDManager(20)
	mockObject := NewMockEObject(t)
	mockOther := NewMockEObject(t)
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
	mockObject := NewMockEObject(t)
	m.Register(mockObject)

	id := m.GetID(mockObject)
	assert.NotNil(t, id)
	assert.Nil(t, m.GetDetachedID(mockObject))

	// unregister
	m.UnRegister(mockObject)
	assert.Nil(t, m.GetID(mockObject))
	assert.Equal(t, id, m.GetDetachedID(mockObject))

	// register again and check it was detached
	m.Register(mockObject)
	assert.Equal(t, id, m.GetID(mockObject))
}
