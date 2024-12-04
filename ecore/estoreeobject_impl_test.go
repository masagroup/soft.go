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
	"github.com/stretchr/testify/require"
)

func TestEStoreEObjectImpl_GetAttribute_Transient(t *testing.T) {

	// create mocks
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockStore := NewMockEStore(t)

	// create object
	o := NewEStoreEObjectImpl(true)
	o.SetEClass(mockClass)
	o.SetEStore(mockStore)

	// first get
	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().GetDefaultValue().Return(nil).Once()
	mockAttribute.EXPECT().IsTransient().Return(true).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Twice()
	assert.Nil(t, o.EGetFromID(0, false))
	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)
}

func TestEStoreEObjectImpl_GetAttribute_Caching(t *testing.T) {
	// create mocks
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockStore := NewMockEStore(t)

	// create object
	o := NewEStoreEObjectImpl(true)
	o.SetEClass(mockClass)
	o.SetEStore(mockStore)

	// first get
	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().IsTransient().Return(false).Once()
	mockClass.EXPECT().GetFeatureCount().Return(1).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Twice()
	mockStore.EXPECT().Get(o, mockAttribute, NO_INDEX).Return(2).Once()
	assert.Equal(t, 2, o.EGetFromID(0, false))
	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)

	// second - test caching
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Once()
	assert.Equal(t, 2, o.EGetFromID(0, false))
	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)
}

func TestEStoreEObjectImpl_GetAttribute_NoCaching(t *testing.T) {
	// create mocks
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockStore := NewMockEStore(t)

	// create object
	o := NewEStoreEObjectImpl(false)
	o.SetEClass(mockClass)
	o.SetEStore(mockStore)

	for i := 0; i < 2; i++ {
		mockAttribute.EXPECT().IsMany().Return(false).Once()
		mockAttribute.EXPECT().IsTransient().Return(false).Once()
		mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Twice()
		mockStore.EXPECT().Get(o, mockAttribute, NO_INDEX).Return(2).Once()
		assert.Equal(t, 2, o.EGetFromID(0, false))
		mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)
	}
}

func TestEStoreEObjectImpl_SetAttribute_Transient(t *testing.T) {
	// create mocks
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockStore := NewMockEStore(t)

	// create object
	o := NewEStoreEObjectImpl(true)
	o.SetEClass(mockClass)
	o.SetEStore(mockStore)

	// set
	mockAttribute.EXPECT().IsTransient().Return(true).Twice()
	mockClass.EXPECT().GetFeatureCount().Return(1).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Times(3)
	o.ESetFromID(0, 2)
	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)

	// test get
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Once()
	assert.Equal(t, 2, o.EGetFromID(0, false))
	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)
}

func TestEStoreEObjectImpl_SetAttribute(t *testing.T) {
	// create mocks
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockStore := NewMockEStore(t)

	// create object
	o := NewEStoreEObjectImpl(false)
	o.SetEClass(mockClass)
	o.SetEStore(mockStore)

	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().IsTransient().Return(false).Twice()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Times(3)
	mockStore.EXPECT().Get(o, mockAttribute, NO_INDEX).Return(nil).Once()
	mockStore.EXPECT().Set(o, mockAttribute, NO_INDEX, 2).Return(nil).Once()
	o.ESetFromID(0, 2)
	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)
}

func TestEStoreEObjectImpl_SetAttribute_Caching(t *testing.T) {
	// create mocks
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockStore := NewMockEStore(t)

	// create object
	o := NewEStoreEObjectImpl(true)
	o.SetEClass(mockClass)
	o.SetEStore(mockStore)

	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().IsTransient().Return(false).Twice()
	mockClass.EXPECT().GetFeatureCount().Return(1).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Times(3)
	mockStore.EXPECT().Get(o, mockAttribute, NO_INDEX).Return(nil).Once()
	mockStore.EXPECT().Set(o, mockAttribute, NO_INDEX, 2).Return(nil).Once()
	o.ESetFromID(0, 2)
	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)
}

func TestEStoreEObjectImpl_GetAttribute_Many_List(t *testing.T) {
	// create mocks
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockStore := NewMockEStore(t)

	// create object
	o := NewEStoreEObjectImpl(false)
	o.SetEClass(mockClass)
	o.SetEStore(mockStore)

	mockAttribute.EXPECT().IsUnique().Return(true).Once()
	mockAttribute.EXPECT().IsMany().Return(true).Once()
	mockAttribute.EXPECT().IsTransient().Return(false).Once()
	mockAttribute.EXPECT().GetEType().Return(nil).Once()
	mockClass.EXPECT().GetFeatureCount().Return(1).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Twice()
	mockStore.EXPECT().Size(o, mockAttribute).Return(0).Once()
	list := o.EGetFromID(0, true)
	require.NotNil(t, list)
	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)

	eobjectlist, _ := list.(EObjectList)
	assert.NotNil(t, eobjectlist)
	enotifyinglist, _ := list.(ENotifyingList)
	assert.NotNil(t, enotifyinglist)
}

func TestEStoreEObjectImpl_GetAttribute_Many_Map(t *testing.T) {
	// create mocks
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockStore := NewMockEStore(t)
	mockType := NewMockEClass(t)
	mockFeature := NewMockEStructuralFeature(t)

	// create object
	o := NewEStoreEObjectImpl(false)
	o.SetEClass(mockClass)
	o.SetEStore(mockStore)

	mockAttribute.EXPECT().IsUnique().Return(true).Once()
	mockAttribute.EXPECT().IsMany().Return(true).Once()
	mockAttribute.EXPECT().IsTransient().Return(false).Once()
	mockAttribute.EXPECT().GetEType().Return(mockType).Twice()
	mockType.EXPECT().GetInstanceTypeName().Return("java.util.Map.Entry").Once()
	mockType.EXPECT().GetEStructuralFeatureFromName("key").Return(mockFeature).Once()
	mockType.EXPECT().GetEStructuralFeatureFromName("value").Return(mockFeature).Once()
	mockClass.EXPECT().GetFeatureCount().Return(1).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Twice()
	mockStore.EXPECT().Size(o, mockAttribute).Return(0).Once()
	m := o.EGetFromID(0, true)
	require.NotNil(t, m)
	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore, mockFeature, mockType)

	emap, _ := m.(EMap)
	assert.NotNil(t, emap)
}

func TestEStoreEObjectImpl_UnSetAttribute_Transient(t *testing.T) {
	// create mocks
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockStore := NewMockEStore(t)

	// create object
	o := NewEStoreEObjectImpl(false)
	o.SetEClass(mockClass)
	o.SetEStore(mockStore)

	mockAttribute.EXPECT().IsTransient().Return(true).Twice()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Times(3)
	o.EUnsetFromID(0)

}

func TestEStoreEObjectImpl_UnSetAttribute(t *testing.T) {
	// create mocks
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockStore := NewMockEStore(t)

	// create object
	o := NewEStoreEObjectImpl(false)
	o.SetEClass(mockClass)
	o.SetEStore(mockStore)

	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().IsTransient().Return(false).Twice()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Times(3)
	mockStore.EXPECT().Get(o, mockAttribute, NO_INDEX).Return(2).Once()
	mockStore.EXPECT().UnSet(o, mockAttribute).Once()
	o.EUnsetFromID(0)
	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)
}

func TestEStoreEObjectImpl_GetContainer(t *testing.T) {
	mockClass := NewMockEClass(t)
	mockStore := NewMockEStore(t)
	mockObject := NewMockEObject(t)
	mockObjectClass := NewMockEClass(t)
	mockReference := NewMockEReference(t)
	mockOpposite := NewMockEReference(t)

	// create object
	o := NewEStoreEObjectImpl(false)
	o.SetEClass(mockClass)
	o.SetEStore(mockStore)

	mockStore.EXPECT().GetContainer(o.AsEObject()).Return(mockObject, mockReference).Once()
	mockObject.EXPECT().EClass().Return(mockObjectClass).Once()
	mockObjectClass.EXPECT().GetFeatureID(mockReference).Return(1).Once()
	mockReference.EXPECT().GetEOpposite().Return(mockOpposite).Once()
	mockClass.EXPECT().GetFeatureID(mockOpposite).Return(1).Once()
	mockObject.EXPECT().EIsProxy().Return(false).Once()
	require.Equal(t, mockObject, o.EContainer())
}
