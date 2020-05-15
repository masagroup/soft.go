package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createMockEObject(id string) *MockEObject {
	mockObject := &MockEObject{}
	mockClass := &MockEClass{}
	mockAttribute := &MockEAttribute{}
	mockObject.On("EClass").Return(mockClass)
	mockClass.On("GetEIDAttribute").Return(mockAttribute)
	if len(id) > 0 {
		mockObject.On("EIsSet", mockAttribute).Return(true)

		mockDataType := &MockEDataType{}
		mockPackage := &MockEPackage{}
		mockFactory := &MockEFactory{}
		mockObject.On("EGet", mockAttribute).Return(id)
		mockAttribute.On("GetEAttributeType").Return(mockDataType)
		mockDataType.On("GetEPackage").Return(mockPackage)
		mockPackage.On("GetEFactoryInstance").Return(mockFactory)
		mockFactory.On("ConvertToString", mockDataType, id).Return(id)

	} else {
		mockObject.On("EIsSet", mockAttribute).Return(false)
	}
	return mockObject
}

func TestEResourceIDManagerImplRegisterNoID(t *testing.T) {

	m := NewEResourceIDManagerImpl()

	mockObject := createMockEObject("")
	mockChild1 := createMockEObject("")
	mockChild2 := createMockEObject("")
	mockChildren := NewImmutableEList([]interface{}{mockChild1, mockChild2})
	mockObject.On("EContents").Return(mockChildren)
	mockChild1.On("EContents").Return(NewImmutableEList([]interface{}{}))
	mockChild2.On("EContents").Return(NewImmutableEList([]interface{}{}))

	m.Register(mockObject)

	assert.Equal(t, "", m.GetID(mockObject))
	assert.Equal(t, "", m.GetID(mockChild1))
	assert.Equal(t, "", m.GetID(mockChild2))
	mock.AssertExpectationsForObjects(t, mockObject)
	mock.AssertExpectationsForObjects(t, mockChildren.ToArray()...)

}

func TestEResourceIDManagerImplRegisterWithID(t *testing.T) {

	m := NewEResourceIDManagerImpl()

	mockObject := createMockEObject("id")
	mockChild1 := createMockEObject("id1")
	mockChild2 := createMockEObject("id2")
	mockChildren := NewImmutableEList([]interface{}{mockChild1, mockChild2})
	mockObject.On("EContents").Return(mockChildren)
	mockChild1.On("EContents").Return(NewImmutableEList([]interface{}{}))
	mockChild2.On("EContents").Return(NewImmutableEList([]interface{}{}))

	m.Register(mockObject)

	assert.Equal(t, "id", m.GetID(mockObject))
	assert.Equal(t, "id1", m.GetID(mockChild1))
	assert.Equal(t, "id2", m.GetID(mockChild2))

	assert.Equal(t, mockObject, m.GetEObject("id"))
	assert.Equal(t, mockChild1, m.GetEObject("id1"))
	assert.Equal(t, mockChild2, m.GetEObject("id2"))

	mock.AssertExpectationsForObjects(t, mockObject)
	mock.AssertExpectationsForObjects(t, mockChildren.ToArray()...)
}

func TestEResourceIDManagerImplUnRegisterWithID(t *testing.T) {

	m := NewEResourceIDManagerImpl()

	mockObject := createMockEObject("id")
	mockChild1 := createMockEObject("id1")
	mockChild2 := createMockEObject("id2")
	mockChildren := NewImmutableEList([]interface{}{mockChild1, mockChild2})
	mockObject.On("EContents").Return(mockChildren)
	mockChild1.On("EContents").Return(NewImmutableEList([]interface{}{}))
	mockChild2.On("EContents").Return(NewImmutableEList([]interface{}{}))

	m.Register(mockObject)

	assert.Equal(t, "id", m.GetID(mockObject))
	assert.Equal(t, "id1", m.GetID(mockChild1))
	assert.Equal(t, "id2", m.GetID(mockChild2))

	assert.Equal(t, mockObject, m.GetEObject("id"))
	assert.Equal(t, mockChild1, m.GetEObject("id1"))
	assert.Equal(t, mockChild2, m.GetEObject("id2"))

	m.UnRegister(mockObject)

	assert.Equal(t, "", m.GetID(mockObject))
	assert.Equal(t, "", m.GetID(mockChild1))
	assert.Equal(t, "", m.GetID(mockChild2))

	assert.Equal(t, nil, m.GetEObject("id"))
	assert.Equal(t, nil, m.GetEObject("id1"))
	assert.Equal(t, nil, m.GetEObject("id2"))

	mock.AssertExpectationsForObjects(t, mockObject)
	mock.AssertExpectationsForObjects(t, mockChildren.ToArray()...)
}
