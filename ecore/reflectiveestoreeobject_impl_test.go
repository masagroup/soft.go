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

type testReflectiveEStoreEObject struct {
	ReflectiveEStoreEObjectImpl
	store EStore
}

func newTestReflectiveEStoreEObject(store EStore) *testReflectiveEStoreEObject {
	o := &testReflectiveEStoreEObject{store: store}
	o.SetInterfaces(o)
	o.Initialize()
	return o
}

func (o *testReflectiveEStoreEObject) GetEStore() EStore {
	return o.store
}

func TestReflectiveEStoreEObjectImpl_GetAttribute_Transient(t *testing.T) {
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockStore := NewMockEStore(t)

	o := newTestReflectiveEStoreEObject(mockStore)
	o.SetEClass(mockClass)

	mockAttribute.EXPECT().IsMany().Return(false).Twice()
	mockAttribute.EXPECT().GetDefaultValue().Return(nil).Once()
	mockAttribute.EXPECT().IsTransient().Return(true).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Twice()

	assert.Nil(t, o.EGetFromID(0, false))

	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)
}

func TestReflectiveEStoreEObjectImpl_SetAttribute_Transient(t *testing.T) {
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockStore := NewMockEStore(t)

	o := newTestReflectiveEStoreEObject(mockStore)
	o.SetEClass(mockClass)

	// Set
	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().IsTransient().Return(true).Twice()
	mockClass.EXPECT().GetFeatureCount().Return(1).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Times(3)
	o.ESetFromID(0, "transient-val")

	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)

	// Get (cache hit)
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Once()
	assert.Equal(t, "transient-val", o.EGetFromID(0, false))

	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)
}

func TestReflectiveEStoreEObjectImpl_GetAttribute_PersistentCaching(t *testing.T) {
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockStore := NewMockEStore(t)

	o := newTestReflectiveEStoreEObject(mockStore)
	o.SetEClass(mockClass)

	// First get
	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().IsTransient().Return(false).Once()
	mockClass.EXPECT().GetFeatureCount().Return(1).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Twice()
	mockStore.EXPECT().Get(o, mockAttribute, NO_INDEX).Return(42).Once()

	assert.Equal(t, 42, o.EGetFromID(0, false))

	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)

	// Second get: cache hit in properties
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Once()
	assert.Equal(t, 42, o.EGetFromID(0, false))

	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)
}

func TestReflectiveEStoreEObjectImpl_SetAndUnset(t *testing.T) {
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockStore := NewMockEStore(t)

	o := newTestReflectiveEStoreEObject(mockStore)
	o.SetEClass(mockClass)

	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().IsTransient().Return(false).Twice()
	mockClass.EXPECT().GetFeatureCount().Return(1).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Times(3)
	mockStore.EXPECT().Get(o, mockAttribute, NO_INDEX).Return(nil).Once()
	mockStore.EXPECT().Set(o, mockAttribute, NO_INDEX, 100, false).Return(nil).Once()

	o.ESetFromID(0, 100)

	// IsSet
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Once()
	assert.True(t, o.EIsSetFromID(0))

	// Unset
	mockAttribute.EXPECT().IsTransient().Return(false).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Twice()
	mockStore.EXPECT().UnSet(o, mockAttribute).Return().Once()

	o.EUnsetFromID(0)

	// IsSet after unset
	mockAttribute.EXPECT().IsTransient().Return(false).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Twice()
	mockStore.EXPECT().IsSet(o, mockAttribute).Return(false).Once()
	assert.False(t, o.EIsSetFromID(0))

	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)
}
