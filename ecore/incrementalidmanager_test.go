package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestIncrementalIDManagerRegister(t *testing.T) {
	m := NewIncrementalIDManager()
	mockObject := &MockEObject{}
	mockChild1 := &MockEObject{}
	mockChild2 := &MockEObject{}
	mockOther := &MockEObject{}
	mockObject.On("EContents").Return(NewImmutableEList([]interface{}{mockChild1, mockChild2})).Once()
	mockChild1.On("EContents").Return(NewEmptyImmutableEList()).Once()
	mockChild2.On("EContents").Return(NewEmptyImmutableEList()).Once()
	m.Register(mockObject)
	mock.AssertExpectationsForObjects(t, mockObject, mockChild1, mockChild2, mockOther)

	assert.Equal(t, 0, m.GetID(mockObject))
	assert.Equal(t, 1, m.GetID(mockChild1))
	assert.Equal(t, 2, m.GetID(mockChild2))
	assert.Nil(t, m.GetID(mockOther))

	assert.Equal(t, mockObject, m.GetEObject(0))
	assert.Equal(t, mockChild1, m.GetEObject(1))
	assert.Equal(t, mockChild2, m.GetEObject(2))
	assert.Nil(t, m.GetEObject(""))
	assert.Nil(t, m.GetEObject(3))

}

func TestIncrementalIDManagerUnRegister(t *testing.T) {
	m := NewIncrementalIDManager()

	// register object and its children
	mockObject := &MockEObject{}
	mockChild1 := &MockEObject{}
	mockChild2 := &MockEObject{}
	mockObject.On("EContents").Return(NewImmutableEList([]interface{}{mockChild1, mockChild2})).Once()
	mockChild1.On("EContents").Return(NewEmptyImmutableEList()).Once()
	mockChild2.On("EContents").Return(NewEmptyImmutableEList()).Once()
	m.Register(mockObject)
	mock.AssertExpectationsForObjects(t, mockObject, mockChild1, mockChild2)

	// unregister one child
	mockChild2.On("EContents").Return(NewEmptyImmutableEList()).Once()
	m.UnRegister(mockChild2)
	assert.Nil(t, m.GetID(mockChild2))
	mock.AssertExpectationsForObjects(t, mockObject, mockChild1, mockChild2)

	// register this child again and check it was detached
	mockChild2.On("EContents").Return(NewEmptyImmutableEList()).Once()
	m.Register(mockChild2)
	mock.AssertExpectationsForObjects(t, mockObject, mockChild1, mockChild2)
	assert.Equal(t, 2, m.GetID(mockChild2))
}
