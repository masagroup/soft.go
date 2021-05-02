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
	m.Register(mockObject)
	assert.Nil(t, m.GetID(mockObject))
	mock.AssertExpectationsForObjects(t, mockObject)
}

func TestEObjectIDManagerImplRegisterWithID(t *testing.T) {

	m := NewEObjectIDManagerImpl()

	mockObject := createMockEObjectWithID("id")
	m.Register(mockObject)
	assert.Equal(t, "id", m.GetID(mockObject))
	mock.AssertExpectationsForObjects(t, mockObject)
}

func TestEObjectIDManagerImplUnRegisterWithID(t *testing.T) {
	m := NewEObjectIDManagerImpl()
	mockObject := createMockEObjectWithID("id")
	m.Register(mockObject)
	assert.Equal(t, "id", m.GetID(mockObject))
	m.UnRegister(mockObject)
	assert.Nil(t, m.GetID(mockObject))
	assert.Equal(t, "id", m.GetDetachedID(mockObject))
	mock.AssertExpectationsForObjects(t, mockObject)
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

func TestEObjectIDManagerGetEObjectInvalidID(t *testing.T) {
	m := NewEObjectIDManagerImpl()
	assert.Nil(t, m.GetEObject(1))
}

func TestEObjectIDManagerClear(t *testing.T) {
	m := NewEObjectIDManagerImpl()
	m.Clear()
	assert.Nil(t, m.GetID(&MockEObject{}))
}
