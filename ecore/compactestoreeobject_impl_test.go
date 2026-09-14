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

type testCompactEStoreEObject struct {
	CompactEStoreEObjectImpl
	store EStore
}

func newTestCompactEStoreEObject(store EStore) *testCompactEStoreEObject {
	o := &testCompactEStoreEObject{store: store}
	o.SetInterfaces(o)
	o.Initialize()
	return o
}

func (o *testCompactEStoreEObject) GetEStore() EStore {
	return o.store
}

func TestCompactEStoreEObjectImpl_GetAttribute_Transient(t *testing.T) {
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockStore := NewMockEStore(t)

	o := newTestCompactEStoreEObject(mockStore)
	o.SetEClass(mockClass)

	mockAttribute.EXPECT().IsMany().Return(false).Twice()
	mockAttribute.EXPECT().GetDefaultValue().Return(nil).Once()
	mockAttribute.EXPECT().IsTransient().Return(true).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Twice()

	assert.Nil(t, o.EGetFromID(0, false))

	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)
}

func TestCompactEStoreEObjectImpl_SetAttribute_Transient(t *testing.T) {
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockStore := NewMockEStore(t)

	o := newTestCompactEStoreEObject(mockStore)
	o.SetEClass(mockClass)

	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().IsTransient().Return(true).Twice()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Times(3)
	o.ESetFromID(0, "compact-transient-val")

	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)

	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Once()
	assert.Equal(t, "compact-transient-val", o.EGetFromID(0, false))

	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)
}

func TestCompactEStoreEObjectImpl_GetAttribute_PersistentCaching(t *testing.T) {
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockStore := NewMockEStore(t)

	o := newTestCompactEStoreEObject(mockStore)
	o.SetEClass(mockClass)

	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().IsTransient().Return(false).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Twice()
	mockStore.EXPECT().Get(o, mockAttribute, NO_INDEX).Return(999).Once()

	assert.Equal(t, 999, o.EGetFromID(0, false))

	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)

	// Second get: cache hit
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Once()
	assert.Equal(t, 999, o.EGetFromID(0, false))

	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)
}

func TestCompactEStoreEObjectImpl_SetAndUnset(t *testing.T) {
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockStore := NewMockEStore(t)

	o := newTestCompactEStoreEObject(mockStore)
	o.SetEClass(mockClass)

	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().IsTransient().Return(false).Twice()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Times(3)
	mockStore.EXPECT().Get(o, mockAttribute, NO_INDEX).Return(nil).Once()
	mockStore.EXPECT().Set(o, mockAttribute, NO_INDEX, "hello", false).Return(nil).Once()

	o.ESetFromID(0, "hello")

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
