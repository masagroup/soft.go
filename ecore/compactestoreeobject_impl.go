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
	"math/bits"
	"slices"
	"sync/atomic"

	"github.com/SokaDance/rmx"
)

type internalCompactEStoreEObjectImpl interface {
	GetEStore() EStore
	CreateList(feature EStructuralFeature, store EStore) EList
	CreateMap(feature EStructuralFeature, store EStore) EMap
}

// CompactEStoreEObjectImpl is an abstract compact reflective EObject implementation.
// It defines no store fields; concrete structs must implement GetEStore() EStore.
// It uses CompactEObjectContainer and a dynamic compact strong cache.
type CompactEStoreEObjectImpl struct {
	CompactEObjectContainer
	class EClass
	// cache holds the cached values of the feature IDs below 64. It is never changed, but replaced
	// under mutex: a value is read from it without lock, at each access to a feature.
	cache atomic.Pointer[compactCache]
	// cachedOverflow holds the cached values of the feature IDs from 64: without a store, the
	// cache is the only place the values of the object are kept. It is guarded by mutex.
	cachedOverflow map[int]any
	mutex          rmx.RecursiveMutex
}

// compactCache is the cached values of the feature IDs below 64 of a CompactEStoreEObjectImpl.
type compactCache struct {
	mask   uint64 // bitmask of cached feature IDs (0..63)
	values []any  // the values of the IDs of mask, in ascending ID order
}

func (o *CompactEStoreEObjectImpl) Initialize() {
	o.CompactEObjectContainer.Initialize()
	o.ESetInternalContainer(unitializedContainer, -1)
}

func (o *CompactEStoreEObjectImpl) asInternal() internalCompactEStoreEObjectImpl {
	return o.GetInterfaces().(internalCompactEStoreEObjectImpl)
}

func (o *CompactEStoreEObjectImpl) getStore() EStore {
	return o.asInternal().GetEStore()
}

func (o *CompactEStoreEObjectImpl) EClass() EClass {
	if o.class == nil {
		return o.AsEObjectInternal().EStaticClass()
	}
	return o.class
}

func (o *CompactEStoreEObjectImpl) SetEClass(class EClass) {
	o.class = class
}

func (o *CompactEStoreEObjectImpl) EStaticFeatureCount() int {
	return 0
}

func (o *CompactEStoreEObjectImpl) EDynamicProperties() EDynamicProperties {
	return o.GetInterfaces().(EDynamicProperties)
}

func (o *CompactEStoreEObjectImpl) Lock() {
	o.mutex.Lock()
}

func (o *CompactEStoreEObjectImpl) Unlock() {
	o.mutex.Unlock()
}

// ClearCache clears all cached feature values from the dynamic compact cache.
func (o *CompactEStoreEObjectImpl) ClearCache() {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.cache.Store(nil)
	o.cachedOverflow = nil
}

// isCached reports whether the value of featureID is cached. From 64, mutex must be held.
func (o *CompactEStoreEObjectImpl) isCached(featureID int) bool {
	if featureID >= 64 {
		_, cached := o.cachedOverflow[featureID]
		return cached
	}
	c := o.cache.Load()
	return c != nil && (c.mask&(1<<featureID)) != 0
}

// getCached returns the cached value of featureID, nil if none. From 64, mutex must be held;
// below, the value is read without lock.
func (o *CompactEStoreEObjectImpl) getCached(featureID int) any {
	if featureID >= 64 {
		return o.cachedOverflow[featureID]
	}
	c := o.cache.Load()
	if c == nil {
		return nil
	}
	bit := uint64(1) << featureID
	if (c.mask & bit) == 0 {
		return nil
	}
	return c.values[bits.OnesCount64(c.mask&(bit-1))]
}

// cachedValue returns the cached value of featureID, nil if none, holding mutex only for the
// IDs from 64.
func (o *CompactEStoreEObjectImpl) cachedValue(featureID int) any {
	if featureID < 64 {
		return o.getCached(featureID)
	}
	o.mutex.Lock()
	defer o.mutex.Unlock()
	return o.getCached(featureID)
}

// setCached caches the value of featureID. mutex must be held.
func (o *CompactEStoreEObjectImpl) setCached(featureID int, value any) {
	if featureID >= 64 {
		if o.cachedOverflow == nil {
			o.cachedOverflow = map[int]any{}
		}
		o.cachedOverflow[featureID] = value
		return
	}
	bit := uint64(1) << featureID
	var mask uint64
	var values []any
	if c := o.cache.Load(); c != nil {
		mask, values = c.mask, c.values
	}
	idx := bits.OnesCount64(mask & (bit - 1))
	if (mask & bit) == 0 {
		// a new slice, the old one may be read meanwhile
		newValues := make([]any, len(values)+1)
		copy(newValues[:idx], values[:idx])
		newValues[idx] = value
		copy(newValues[idx+1:], values[idx:])
		values = newValues
	} else {
		values = slices.Clone(values)
		values[idx] = value
	}
	o.cache.Store(&compactCache{mask: mask | bit, values: values})
}

// unsetCached removes the cached value of featureID. mutex must be held.
func (o *CompactEStoreEObjectImpl) unsetCached(featureID int) {
	if !o.isCached(featureID) {
		return
	}
	if featureID >= 64 {
		delete(o.cachedOverflow, featureID)
		return
	}
	bit := uint64(1) << featureID
	c := o.cache.Load()
	if c.mask == bit {
		o.cache.Store(nil)
		return
	}
	idx := bits.OnesCount64(c.mask & (bit - 1))
	// a new slice, the old one may be read meanwhile
	values := slices.Delete(slices.Clone(c.values), idx, idx+1)
	o.cache.Store(&compactCache{mask: c.mask &^ bit, values: values})
}

// the store is read without holding the object : the store may wait for another goroutine
// that needs the object, as the store writer encoding it

func (o *CompactEStoreEObjectImpl) EDynamicGet(dynamicFeatureID int) any {
	// 1. Cache hit
	result := o.cachedValue(dynamicFeatureID)
	if result != nil {
		return result
	}

	// 2. Cache miss
	shouldCache := false
	feature := o.EClass().GetEStructuralFeature(dynamicFeatureID)
	if !feature.IsTransient() {
		if feature.IsMany() {
			if IsMapType(feature) {
				result = o.asInternal().CreateMap(feature, o.getStore())
			} else {
				result = o.asInternal().CreateList(feature, o.getStore())
			}
			shouldCache = true
		} else if store := o.getStore(); store != nil {
			result = store.Get(o.AsEObject(), feature, NO_INDEX)
			shouldCache = true
		}
	} else if feature.IsMany() {
		if IsMapType(feature) {
			result = o.asInternal().CreateMap(feature, nil)
		} else {
			result = o.asInternal().CreateList(feature, nil)
		}
		shouldCache = true
	}

	if shouldCache && result != nil {
		o.mutex.Lock()
		defer o.mutex.Unlock()
		// unless set meanwhile
		if current := o.getCached(dynamicFeatureID); current != nil {
			return current
		}
		o.setCached(dynamicFeatureID, result)
	}
	return result
}

func (o *CompactEStoreEObjectImpl) EDynamicSet(dynamicFeatureID int, value any) {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	feature := o.EClass().GetEStructuralFeature(dynamicFeatureID)

	if !feature.IsTransient() {
		if store := o.getStore(); store != nil {
			store.Set(o.AsEObject(), feature, NO_INDEX, value, false)
		}
	}

	o.setCached(dynamicFeatureID, value)
}

func (o *CompactEStoreEObjectImpl) EDynamicIsSet(dynamicFeatureID int) bool {
	cached := o.cachedValue(dynamicFeatureID) != nil
	if cached {
		return true
	}

	feature := o.EClass().GetEStructuralFeature(dynamicFeatureID)
	if !feature.IsTransient() {
		if store := o.getStore(); store != nil {
			return store.IsSet(o.AsEObject(), feature)
		}
	}
	return false
}

func (o *CompactEStoreEObjectImpl) EDynamicUnset(dynamicFeatureID int) {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	o.unsetCached(dynamicFeatureID)

	feature := o.EClass().GetEStructuralFeature(dynamicFeatureID)
	if !feature.IsTransient() {
		if store := o.getStore(); store != nil {
			store.UnSet(o.AsEObject(), feature)
		}
	}
}

func (o *CompactEStoreEObjectImpl) CreateList(feature EStructuralFeature, store EStore) EList {
	l := NewEStoreList(o.AsEObject(), feature, store)
	l.SetCache(true)
	return l
}

func (o *CompactEStoreEObjectImpl) CreateMap(feature EStructuralFeature, store EStore) EMap {
	eClass := feature.GetEType().(EClass)
	return NewEStoreMap(eClass, o.AsEObject(), feature, store)
}

func (o *CompactEStoreEObjectImpl) ESetInternalContainer(newContainer EObject, newContainerFeatureID int) {
	o.CompactEObjectContainer.ESetInternalContainer(newContainer, newContainerFeatureID)
	o.setContainerInStore()
}

func (o *CompactEStoreEObjectImpl) EInternalContainer() EObject {
	o.initializeContainerFromStore()
	return o.CompactEObjectContainer.EInternalContainer()
}

func (o *CompactEStoreEObjectImpl) EInternalContainerFeatureID() int {
	o.initializeContainerFromStore()
	return o.CompactEObjectContainer.EInternalContainerFeatureID()
}

func (o *CompactEStoreEObjectImpl) initializeContainerFromStore() {
	if o.CompactEObjectContainer.EInternalContainer() == unitializedContainer {
		if store := o.getStore(); store != nil {
			container, feature := store.GetContainer(o.AsEObject())
			if container != nil && feature != nil {
				featureID := EOPPOSITE_FEATURE_BASE - container.EClass().GetFeatureID(feature)
				if reference, _ := feature.(EReference); reference != nil {
					if opposite := reference.GetEOpposite(); opposite != nil {
						featureID = o.AsEObject().EClass().GetFeatureID(opposite)
					}
				}
				o.CompactEObjectContainer.ESetInternalContainer(container, featureID)
			} else {
				o.CompactEObjectContainer.ESetInternalContainer(nil, -1)
			}
		} else {
			o.CompactEObjectContainer.ESetInternalContainer(nil, -1)
		}
	}
}

func (o *CompactEStoreEObjectImpl) setContainerInStore() {
	if store := o.getStore(); store != nil {
		container := o.CompactEObjectContainer.EInternalContainer()
		containerFeatureID := o.CompactEObjectContainer.EInternalContainerFeatureID()
		if container != unitializedContainer {
			var containerFeature EStructuralFeature
			if container != nil {
				if containerFeatureID <= EOPPOSITE_FEATURE_BASE {
					containerFeature = container.EClass().GetEStructuralFeature(EOPPOSITE_FEATURE_BASE - containerFeatureID)
				} else {
					containerFeature = o.AsEObject().EClass().GetEStructuralFeature(containerFeatureID).(EReference).GetEOpposite()
				}
			}
			store.SetContainer(o.AsEObject(), container, containerFeature)
		}
	}
}
