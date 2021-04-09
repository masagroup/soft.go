package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUniqueIDManagerRegister(t *testing.T) {
	m := NewUniqueIDManager(20)
	mockObject := &MockEObject{}
	mockChild1 := &MockEObject{}
	mockChild2 := &MockEObject{}
	mockOther := &MockEObject{}
	mockObject.On("EContents").Return(NewImmutableEList([]interface{}{mockChild1, mockChild2})).Once()
	mockChild1.On("EContents").Return(NewEmptyImmutableList()).Once()
	mockChild2.On("EContents").Return(NewEmptyImmutableList()).Once()
	m.Register(mockObject)
	mock.AssertExpectationsForObjects(t, mockObject, mockChild1, mockChild2, mockOther)

	id := m.GetID(mockObject)
	assert.NotNil(t, id)
	id1 := m.GetID(mockChild1)
	assert.NotNil(t, id1)
	id2 := m.GetID(mockChild2)
	assert.NotNil(t, id2)
	assert.Nil(t, m.GetID(mockOther))

	assert.Equal(t, mockObject, m.GetEObject(id))
	assert.Equal(t, mockChild1, m.GetEObject(id1))
	assert.Equal(t, mockChild2, m.GetEObject(id2))
	assert.Nil(t, m.GetEObject(""))
	assert.Nil(t, m.GetEObject(3))

}

func TestUniqueIDManagerUnRegister(t *testing.T) {
	m := NewUniqueIDManager(20)

	// register object and its children
	mockObject := &MockEObject{}
	mockChild1 := &MockEObject{}
	mockChild2 := &MockEObject{}
	mockObject.On("EContents").Return(NewImmutableEList([]interface{}{mockChild1, mockChild2})).Once()
	mockChild1.On("EContents").Return(NewEmptyImmutableList()).Once()
	mockChild2.On("EContents").Return(NewEmptyImmutableList()).Once()
	m.Register(mockObject)
	mock.AssertExpectationsForObjects(t, mockObject, mockChild1, mockChild2)
	id2 := m.GetID(mockChild2)
	assert.NotNil(t, id2)

	// unregister one child
	mockChild2.On("EContents").Return(NewEmptyImmutableList()).Once()
	m.UnRegister(mockChild2)
	assert.Nil(t, m.GetID(mockChild2))
	mock.AssertExpectationsForObjects(t, mockObject, mockChild1, mockChild2)

	// register this child again and check it was detached
	mockChild2.On("EContents").Return(NewEmptyImmutableList()).Once()
	m.Register(mockChild2)
	mock.AssertExpectationsForObjects(t, mockObject, mockChild1, mockChild2)
	assert.Equal(t, id2, m.GetID(mockChild2))
}
