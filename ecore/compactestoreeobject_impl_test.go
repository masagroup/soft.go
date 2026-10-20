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
	"sync"
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

func TestCompactEStoreEObjectImpl_Cache_FeatureIDsFrom64(t *testing.T) {
	o := newTestCompactEStoreEObject(nil)
	o.setCached(3, "a")
	o.setCached(64, "b")
	o.setCached(70, "c")
	assert.True(t, o.isCached(3))
	assert.True(t, o.isCached(64))
	assert.True(t, o.isCached(70))
	assert.False(t, o.isCached(65))
	assert.Equal(t, "a", o.getCached(3))
	assert.Equal(t, "b", o.getCached(64))
	assert.Equal(t, "c", o.getCached(70))

	o.setCached(70, "d")
	assert.Equal(t, "d", o.getCached(70))

	o.unsetCached(64)
	assert.False(t, o.isCached(64))
	assert.Nil(t, o.getCached(64))
	assert.Equal(t, "d", o.getCached(70))

	o.ClearCache()
	assert.False(t, o.isCached(3))
	assert.False(t, o.isCached(70))
}

func TestCompactEStoreEObjectImpl_SetAttribute_NoStore_FeatureIDFrom64(t *testing.T) {
	// without a store, the cache keeps the values of the object: a feature ID from 64 is kept too
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	o := newTestCompactEStoreEObject(nil)
	o.SetEClass(mockClass)

	mockClass.EXPECT().GetEStructuralFeature(70).Return(mockAttribute)
	mockAttribute.EXPECT().IsTransient().Return(false)
	o.EDynamicSet(70, "value")
	assert.Equal(t, "value", o.EDynamicGet(70))
	assert.True(t, o.EDynamicIsSet(70))
	o.EDynamicUnset(70)
	assert.False(t, o.EDynamicIsSet(70))
}

// The cached values are read without lock while others are set: a reader sees a value set
// before or after, never a partial cache.
func TestCompactEStoreEObjectImpl_Cache_ConcurrentReads(t *testing.T) {
	o := newTestCompactEStoreEObject(nil)
	o.setCached(1, "a")

	var group sync.WaitGroup
	for range 8 {
		group.Go(func() {
			for range 1000 {
				if o.cachedValue(1) != "a" {
					t.Error("a cached value was not found")
					return
				}
				o.cachedValue(2)
			}
		})
	}
	for i := range 1000 {
		o.setCached(2+i%10, i)
		if i%3 == 0 {
			o.unsetCached(2 + i%10)
		}
	}
	group.Wait()
	assert.Equal(t, "a", o.cachedValue(1))
}
