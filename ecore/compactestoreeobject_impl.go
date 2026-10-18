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
	"maps"
	"math/bits"
	"slices"
	"sync/atomic"
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
	// cache holds the cached values. It is never changed, but replaced with a compare-and-swap:
	// a value is read from it without lock, at each access to a feature, and a value read from
	// the store is cached without lock. The object takes no lock: a concrete type that needs one
	// (ELockProvider for the SQL encoder, ordered writes) provides it.
	cache atomic.Pointer[compactCache]
}

// compactCache is the cached values of a CompactEStoreEObjectImpl.
type compactCache struct {
	mask   uint64 // bitmask of cached feature IDs (0..63)
	values []any  // the values of the IDs of mask, in ascending ID order
	// overflow holds the cached values of the feature IDs from 64: without a store, the cache is
	// the only place the values of the object are kept.
	overflow map[int]any
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

// ClearCache clears all cached feature values from the dynamic compact cache.
func (o *CompactEStoreEObjectImpl) ClearCache() {
	o.cache.Store(nil)
}

// get returns the value of featureID in c, and whether it is cached.
func (c *compactCache) get(featureID int) (any, bool) {
	if c == nil {
		return nil, false
	}
	if featureID >= 64 {
		value, cached := c.overflow[featureID]
		return value, cached
	}
	bit := uint64(1) << featureID
	if (c.mask & bit) == 0 {
		return nil, false
	}
	return c.values[bits.OnesCount64(c.mask&(bit-1))], true
}

// with returns a copy of c with the value of featureID: c is read meanwhile.
func (c *compactCache) with(featureID int, value any) *compactCache {
	result := &compactCache{}
	if c != nil {
		*result = *c
	}
	if featureID >= 64 {
		result.overflow = maps.Clone(result.overflow)
		if result.overflow == nil {
			result.overflow = map[int]any{}
		}
		result.overflow[featureID] = value
		return result
	}
	bit := uint64(1) << featureID
	idx := bits.OnesCount64(result.mask & (bit - 1))
	if (result.mask & bit) == 0 {
		values := make([]any, len(result.values)+1)
		copy(values[:idx], result.values[:idx])
		values[idx] = value
		copy(values[idx+1:], result.values[idx:])
		result.values = values
		result.mask |= bit
	} else {
		result.values = slices.Clone(result.values)
		result.values[idx] = value
	}
	return result
}

// without returns a copy of c without the value of featureID, nil if empty: c is read meanwhile.
func (c *compactCache) without(featureID int) *compactCache {
	result := *c
	if featureID >= 64 {
		result.overflow = maps.Clone(result.overflow)
		delete(result.overflow, featureID)
		if len(result.overflow) == 0 {
			result.overflow = nil
		}
	} else {
		bit := uint64(1) << featureID
		idx := bits.OnesCount64(result.mask & (bit - 1))
		result.values = slices.Delete(slices.Clone(result.values), idx, idx+1)
		result.mask &^= bit
	}
	if result.mask == 0 && result.overflow == nil {
		return nil
	}
	return &result
}

// isCached reports whether the value of featureID is cached.
func (o *CompactEStoreEObjectImpl) isCached(featureID int) bool {
	_, cached := o.cache.Load().get(featureID)
	return cached
}

// getCached returns the cached value of featureID, nil if none. It reads without lock.
func (o *CompactEStoreEObjectImpl) getCached(featureID int) any {
	value, _ := o.cache.Load().get(featureID)
	return value
}

// cachedValue returns the cached value of featureID, nil if none.
func (o *CompactEStoreEObjectImpl) cachedValue(featureID int) any {
	return o.getCached(featureID)
}

// setCached caches the value of featureID.
func (o *CompactEStoreEObjectImpl) setCached(featureID int, value any) {
	for {
		c := o.cache.Load()
		if o.cache.CompareAndSwap(c, c.with(featureID, value)) {
			return
		}
	}
}

// cacheIfAbsent caches the value of featureID unless one is cached, and returns the cached one.
func (o *CompactEStoreEObjectImpl) cacheIfAbsent(featureID int, value any) any {
	for {
		c := o.cache.Load()
		if current, _ := c.get(featureID); current != nil {
			return current
		}
		if o.cache.CompareAndSwap(c, c.with(featureID, value)) {
			return value
		}
	}
}

// unsetCached removes the cached value of featureID.
func (o *CompactEStoreEObjectImpl) unsetCached(featureID int) {
	for {
		c := o.cache.Load()
		if _, cached := c.get(featureID); !cached {
			return
		}
		if o.cache.CompareAndSwap(c, c.without(featureID)) {
			return
		}
	}
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
		// unless set meanwhile
		return o.cacheIfAbsent(dynamicFeatureID, result)
	}
	return result
}

// EDynamicSet writes the value to the store, then caches it. It takes no lock: two goroutines
// setting the same feature concurrently may leave the cache with the value the store got first.
// A concrete type that needs writes ordered serializes them itself.
func (o *CompactEStoreEObjectImpl) EDynamicSet(dynamicFeatureID int, value any) {
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

// EDynamicUnset unsets the value in the store, then removes it from the cache. It takes no lock,
// as EDynamicSet.
func (o *CompactEStoreEObjectImpl) EDynamicUnset(dynamicFeatureID int) {
	feature := o.EClass().GetEStructuralFeature(dynamicFeatureID)
	if !feature.IsTransient() {
		if store := o.getStore(); store != nil {
			store.UnSet(o.AsEObject(), feature)
		}
	}

	o.unsetCached(dynamicFeatureID)
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
