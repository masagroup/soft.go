// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createMockEObjectWithID(t *testing.T, id string) *MockEObject {
	mockObject := NewMockEObject(t)
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockObject.EXPECT().EClass().Return(mockClass)
	mockClass.EXPECT().GetEIDAttribute().Return(mockAttribute)
	if len(id) > 0 {
		mockObject.EXPECT().EIsSet(mockAttribute).Return(true)

		mockDataType := NewMockEDataType(t)
		mockPackage := NewMockEPackage(t)
		mockFactory := NewMockEFactory(t)
		mockObject.EXPECT().EGet(mockAttribute).Return(id)
		mockAttribute.EXPECT().GetEAttributeType().Return(mockDataType)
		mockDataType.EXPECT().GetEPackage().Return(mockPackage)
		mockPackage.EXPECT().GetEFactoryInstance().Return(mockFactory)
		mockFactory.EXPECT().ConvertToString(mockDataType, id).Return(id)

	} else {
		mockObject.EXPECT().EIsSet(mockAttribute).Return(false)
	}
	return mockObject
}

func TestEObjectIDManagerImplRegisterNoID(t *testing.T) {
	m := NewEObjectIDManagerImpl()
	mockObject := createMockEObjectWithID(t, "")
	m.Register(mockObject)
	assert.Nil(t, m.GetID(mockObject))
}

func TestEObjectIDManagerImplRegisterWithID(t *testing.T) {
	m := NewEObjectIDManagerImpl()
	mockObject := createMockEObjectWithID(t, "id")
	m.Register(mockObject)
	assert.Equal(t, "id", m.GetID(mockObject))
}

func TestEObjectIDManagerImplUnRegisterWithID(t *testing.T) {
	m := NewEObjectIDManagerImpl()
	mockObject := createMockEObjectWithID(t, "id")
	m.Register(mockObject)
	assert.Equal(t, "id", m.GetID(mockObject))
	m.UnRegister(mockObject)
	assert.Nil(t, m.GetID(mockObject))
	assert.Equal(t, "id", m.GetDetachedID(mockObject))
}

func TestEObjectIDManagerSetID(t *testing.T) {
	m := NewEObjectIDManagerImpl()

	mockObject := NewMockEObject(t)
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockDataType := NewMockEDataType(t)
	mockPackage := NewMockEPackage(t)
	mockFactory := NewMockEFactory(t)
	mockIDValue := 0

	// set mock object id
	mockObject.EXPECT().EClass().Return(mockClass).Once()
	mockClass.EXPECT().GetEIDAttribute().Return(mockAttribute).Once()
	mockObject.EXPECT().ESet(mockAttribute, mockIDValue).Once()
	mockAttribute.EXPECT().GetEAttributeType().Return(mockDataType).Once()
	mockDataType.EXPECT().GetEPackage().Return(mockPackage).Once()
	mockPackage.EXPECT().GetEFactoryInstance().Return(mockFactory).Once()
	mockFactory.EXPECT().CreateFromString(mockDataType, "id1").Return(mockIDValue).Once()
	assert.Nil(t, m.SetID(mockObject, "id1"))

	assert.Equal(t, "id1", m.GetID(mockObject))
	assert.Equal(t, "id1", m.GetID(mockObject))
	assert.Equal(t, mockObject, m.GetEObject("id1"))
	mock.AssertExpectationsForObjects(t, mockObject, mockClass, mockAttribute, mockDataType, mockPackage, mockFactory)

	// reset mock object id
	mockObject.EXPECT().EClass().Return(mockClass).Once()
	mockClass.EXPECT().GetEIDAttribute().Return(mockAttribute).Once()
	mockObject.EXPECT().EUnset(mockAttribute).Once()
	assert.Nil(t, m.SetID(mockObject, nil))
	mock.AssertExpectationsForObjects(t, mockObject, mockClass, mockAttribute, mockDataType, mockPackage, mockFactory)

	// error
	assert.NotNil(t, m.SetID(mockObject, 1))

}

func TestEObjectIDManagerGetEObjectInvalidID(t *testing.T) {
	m := NewEObjectIDManagerImpl()
	assert.Nil(t, m.GetEObject(1))
}

func TestEObjectIDManagerClear(t *testing.T) {
	m := NewEObjectIDManagerImpl()
	o := NewMockEObject(t)
	m.Clear()
	assert.Nil(t, m.GetID(o))
}
