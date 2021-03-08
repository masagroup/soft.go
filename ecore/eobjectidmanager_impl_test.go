package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createMockEObjectWithID(id string) *MockEObject {
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

func TestEObjectIDManagerImplRegisterNoID(t *testing.T) {

	m := NewEObjectIDManagerImpl()

	mockObject := createMockEObjectWithID("")
	mockChild1 := createMockEObjectWithID("")
	mockChild2 := createMockEObjectWithID("")
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

func TestEObjectIDManagerImplRegisterWithID(t *testing.T) {

	m := NewEObjectIDManagerImpl()

	mockObject := createMockEObjectWithID("id")
	mockChild1 := createMockEObjectWithID("id1")
	mockChild2 := createMockEObjectWithID("id2")
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

func TestEObjectIDManagerImplUnRegisterWithID(t *testing.T) {

	m := NewEObjectIDManagerImpl()

	mockObject := createMockEObjectWithID("id")
	mockChild1 := createMockEObjectWithID("id1")
	mockChild2 := createMockEObjectWithID("id2")
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

func TestEObjectIDManagerSetID(t *testing.T) {
	m := NewEObjectIDManagerImpl()

	mockObject := &MockEObject{}
	mockClass := &MockEClass{}
	mockAttribute := &MockEAttribute{}
	mockDataType := &MockEDataType{}
	mockPackage := &MockEPackage{}
	mockFactory := &MockEFactory{}
	mockIDValue := 0

	// set mock object id
	mockObject.On("EClass").Return(mockClass).Once()
	mockClass.On("GetEIDAttribute").Return(mockAttribute).Once()
	mockObject.On("ESet", mockAttribute, mockIDValue).Once()
	mockAttribute.On("GetEAttributeType").Return(mockDataType).Once()
	mockDataType.On("GetEPackage").Return(mockPackage).Once()
	mockPackage.On("GetEFactoryInstance").Return(mockFactory).Once()
	mockFactory.On("CreateFromString", mockDataType, "id1").Return(mockIDValue).Once()

	m.SetID(mockObject, "id1")

	assert.Equal(t, "id1", m.GetID(mockObject))
	assert.Equal(t, "id1", m.GetID(mockObject))
	assert.Equal(t, mockObject, m.GetEObject("id1"))
	mock.AssertExpectationsForObjects(t, mockObject, mockClass, mockAttribute, mockDataType, mockPackage, mockFactory)

	// reset mock object id
	mockObject.On("EClass").Return(mockClass).Once()
	mockClass.On("GetEIDAttribute").Return(mockAttribute).Once()
	mockObject.On("EUnset", mockAttribute).Once()
	m.SetID(mockObject, nil)
	mock.AssertExpectationsForObjects(t, mockObject, mockClass, mockAttribute, mockDataType, mockPackage, mockFactory)
}
