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
	"slices"
	"sync"
	"sync/atomic"
)

var unitializedContainer EObject = newEObjectImpl()

// EStoreEObjectImpl is an EObject whose values are kept in a store set on it (SetEStore), and
// cached when it caches (SetCache). Its values and the store lists of its many-valued features
// are kept in the compact cache of CompactEStoreEObjectImpl, read without lock.
type EStoreEObjectImpl struct {
	CompactEStoreEObjectImpl
	// config is the store and the cache flag of the object. It is never changed but replaced
	// under mutex, after the cache: a value read with a configuration is cached only while the
	// cache has its epoch.
	config atomic.Pointer[eStoreEObjectConfig]
	// mutex serializes the writers of the object, and is held by Lock while the object is
	// encoded. Readers never take it.
	mutex sync.Mutex
	// lockedLists are the store lists locked by Lock, guarded by mutex: a pointer, as it is
	// set only while the object is locked
	lockedLists *[]*EStoreList
}

// eStoreEObjectConfig is a configuration of an EStoreEObjectImpl. It is never changed.
type eStoreEObjectConfig struct {
	store EStore
	cache bool
	epoch uint64
}

func (c *eStoreEObjectConfig) getStore() EStore {
	if c == nil {
		return nil
	}
	return c.store
}

func (c *eStoreEObjectConfig) isCache() bool {
	return c != nil && c.cache
}

func (c *eStoreEObjectConfig) getEpoch() uint64 {
	if c == nil {
		return 0
	}
	return c.epoch
}

func NewEStoreEObjectImpl(cache bool) *EStoreEObjectImpl {
	o := new(EStoreEObjectImpl)
	o.config.Store(&eStoreEObjectConfig{cache: cache})
	o.SetInterfaces(o)
	o.Initialize()
	return o
}

func (o *EStoreEObjectImpl) AsEStoreEObject() EStoreEObject {
	return o.GetInterfaces().(EStoreEObject)
}

// Lock prevents the writes to the object and to its store lists until Unlock. The object can be
// read meanwhile, by the goroutine holding the lock too.
func (o *EStoreEObjectImpl) Lock() {
	o.mutex.Lock()
	if lists := lockStoreLists(o.cachedValues()); lists != nil {
		o.lockedLists = &lists
	}
}

func (o *EStoreEObjectImpl) Unlock() {
	if o.lockedLists != nil {
		unlockStoreLists(*o.lockedLists)
	}
	o.lockedLists = nil
	o.mutex.Unlock()
}

func (o *EStoreEObjectImpl) GetEStore() EStore {
	return o.config.Load().getStore()
}

// reconfigure sets the store and the cache flag of the object, and replaces its cached values
// with the ones values returns from the current ones, which it returns. mutex must be held.
func (o *EStoreEObjectImpl) reconfigure(store EStore, cache bool, values func(current []cachedFeatureValue) []cachedFeatureValue) []cachedFeatureValue {
	next := &eStoreEObjectConfig{store: store, cache: cache, epoch: o.config.Load().getEpoch() + 1}
	for {
		c := o.cache.Load()
		current := c.featureValues()
		replaced := emptyCache(next.epoch)
		for _, v := range values(current) {
			replaced = replaced.with(v.featureID, v.value)
		}
		// unless a value read from the store was cached meanwhile
		if o.cache.CompareAndSwap(c, replaced) {
			// the cache of the new epoch first: from now on, a value read with the previous
			// configuration fails to be cached
			o.config.Store(next)
			return current
		}
	}
}

func (o *EStoreEObjectImpl) SetEStore(newStore EStore) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.SetEStoreLocked(newStore)
}

// SetEStoreLocked sets the store of the object: mutex must be held, by Lock or SetEStore.
func (o *EStoreEObjectImpl) SetEStoreLocked(newStore EStore) {
	current := o.config.Load()
	oldStore := current.getStore()
	if oldStore == newStore {
		return
	}
	cache := current.isCache()
	old := o.reconfigure(newStore, cache, func(values []cachedFeatureValue) []cachedFeatureValue {
		if newStore == nil {
			if !cache {
				// build cache with previous store
				values = slices.Clone(values)
				for featureID := range o.EClass().GetFeatureCount() {
					if eFeature := o.eDynamicFeature(featureID); !eFeature.IsTransient() && !eFeature.IsMany() {
						values = setFeatureValue(values, featureID, oldStore.Get(o.AsEObject(), eFeature, NO_INDEX))
					}
				}
			}
			return values
		}
		if !cache {
			// clear properties because we are not caching
			return nil
		}
		return values
	})
	if newStore != nil {
		// set container in store
		o.setContainerInStore()

		// set children collection store
		for _, v := range old {
			if storeList, isStoreList := v.value.(*EStoreList); isStoreList {
				if o.lockedLists != nil && slices.Contains(*o.lockedLists, storeList) {
					// locked with the object by Lock
					storeList.SetEStoreLocked(newStore)
				} else {
					storeList.SetEStore(newStore)
				}
			}
		}
	}
}

func (o *EStoreEObjectImpl) SetCache(cache bool) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	current := o.config.Load()
	if current.isCache() == cache {
		return
	}
	store := current.getStore()
	old := o.reconfigure(store, cache, func(values []cachedFeatureValue) []cachedFeatureValue {
		if store != nil {
			return nil
		}
		return values
	})
	// set cache for all properties
	for _, v := range old {
		if sc, _ := v.value.(ECacheProvider); sc != nil {
			sc.SetCache(cache)
		}
	}
}

func (o *EStoreEObjectImpl) IsCache() bool {
	return o.config.Load().isCache()
}

// ClearCache clears the cached values.
func (o *EStoreEObjectImpl) ClearCache() {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.CompactEStoreEObjectImpl.ClearCache()
}

// the store is read without holding the object : the store may wait for another goroutine
// that needs the object, as the store writer encoding it

func (o *EStoreEObjectImpl) EDynamicIsSet(dynamicFeatureID int) bool {
	if cached := o.getCached(dynamicFeatureID); cached != nil {
		return isSetValue(o.eDynamicFeature(dynamicFeatureID), cached)
	}
	if store := o.config.Load().getStore(); store != nil {
		eFeature := o.eDynamicFeature(dynamicFeatureID)
		return store.IsSet(o.AsEObject(), eFeature)
	}
	return false
}

func (o *EStoreEObjectImpl) EDynamicGet(dynamicFeatureID int) any {
	// retrieve value
	if result := o.getCached(dynamicFeatureID); result != nil {
		return result
	}
	for {
		config := o.config.Load()
		// compute value
		var result any
		cache := false
		store := config.getStore()
		objectCache := config.isCache()
		if eFeature := o.eDynamicFeature(dynamicFeatureID); !eFeature.IsTransient() {
			if eFeature.IsMany() {
				if IsMapType(eFeature) {
					result = o.createMap(eFeature, store, objectCache)
				} else {
					result = o.createList(eFeature, store, objectCache)
				}
				// feature is a container : create corresponding container type and cache it
				// in properties to avoid multiple allocation of the same container
				cache = true
			} else if store != nil {
				// feature is not transient and we have a store
				result = store.Get(o.AsEObject(), eFeature, NO_INDEX)
				cache = objectCache
			}
		}
		if !cache || result == nil {
			return result
		}
		// store value in properties, unless set meanwhile
		if value, cached := o.cacheIfAbsentAt(config.getEpoch(), dynamicFeatureID, result); cached {
			return value
		}
		// store or cache changed meanwhile: compute again
	}
}

func (o *EStoreEObjectImpl) EDynamicSet(dynamicFeatureID int, value any) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	eFeature := o.eDynamicFeature(dynamicFeatureID)
	config := o.config.Load()
	if store := config.getStore(); store != nil && !eFeature.IsTransient() {
		// store and feature is not transient
		store.Set(o.AsEObject(), eFeature, NO_INDEX, value, false)
		if !config.isCache() {
			// no cache : the value is in the store only
			return
		}
	}
	// cached, or no store or transient feature: the cache is the only place of the value
	o.setCached(dynamicFeatureID, value)
}

func (o *EStoreEObjectImpl) EDynamicUnset(dynamicFeatureID int) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	// retrieve previous value and unset properties
	o.unsetCached(dynamicFeatureID)
	eFeature := o.eDynamicFeature(dynamicFeatureID)
	if store := o.config.Load().getStore(); store != nil && !eFeature.IsTransient() {
		store.UnSet(o.AsEObject(), eFeature)
	}
}

func (o *EStoreEObjectImpl) eDynamicFeature(dynamicFeatureID int) EStructuralFeature {
	return o.EClass().GetEStructuralFeature(o.EStaticFeatureCount() + dynamicFeatureID)
}

func (o *EStoreEObjectImpl) createList(eFeature EStructuralFeature, store EStore, cache bool) EList {
	l := NewEStoreList(o.AsEObject(), eFeature, store)
	l.SetCache(cache)
	return l
}

func (o *EStoreEObjectImpl) createMap(eFeature EStructuralFeature, store EStore, cache bool) EMap {
	eClass := eFeature.GetEType().(EClass)
	m := NewEStoreMap(eClass, o.AsEObject(), eFeature, store)
	m.SetCache(cache)
	return m
}
