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

type EStoreEObjectImpl struct {
	ReflectiveEObjectImpl
	// state is the store, the cache flag and the values of the object. It is never changed but
	// replaced with a compare-and-swap: it is read without lock, and a value read from the store
	// is kept without lock, unless the store or the cache flag changed meanwhile.
	state atomic.Pointer[eStoreEObjectState]
	// mutex serializes the writers of the object, and is held by Lock while the object is
	// encoded. Readers never take it.
	mutex sync.Mutex
	// lockedLists are the store lists locked by Lock, guarded by mutex
	lockedLists []*EStoreList
}

type eStoreEObjectState struct {
	store      EStore
	cache      bool
	properties []any
}

func (s *eStoreEObjectState) getStore() EStore {
	if s == nil {
		return nil
	}
	return s.store
}

func (s *eStoreEObjectState) isCache() bool {
	return s != nil && s.cache
}

func (s *eStoreEObjectState) getProperties() []any {
	if s == nil {
		return nil
	}
	return s.properties
}

func (s *eStoreEObjectState) getProperty(featureID int) any {
	if s == nil || s.properties == nil {
		return nil
	}
	return s.properties[featureID]
}

// clone returns a copy of s to change: it is read meanwhile.
func (s *eStoreEObjectState) clone() *eStoreEObjectState {
	if s == nil {
		return &eStoreEObjectState{}
	}
	result := *s
	return &result
}

// withProperty returns a copy of s with the value of featureID of class.
func (s *eStoreEObjectState) withProperty(featureID int, value any, class EClass) *eStoreEObjectState {
	result := s.clone()
	if result.properties != nil {
		result.properties = slices.Clone(result.properties)
	} else {
		result.properties = make([]any, class.GetFeatureCount())
	}
	result.properties[featureID] = value
	return result
}

func NewEStoreEObjectImpl(cache bool) *EStoreEObjectImpl {
	o := new(EStoreEObjectImpl)
	o.state.Store(&eStoreEObjectState{cache: cache})
	o.SetInterfaces(o)
	o.Initialize()
	o.ESetInternalContainer(unitializedContainer, -1)
	return o
}

func (o *EStoreEObjectImpl) Initialize() {
	o.ReflectiveEObjectImpl.Initialize()
	o.ESetInternalContainer(unitializedContainer, -1)
}

func (o *EStoreEObjectImpl) AsEStoreEObject() EStoreEObject {
	return o.GetInterfaces().(EStoreEObject)
}

// Lock prevents the writes to the object and to its store lists until Unlock. The object can be
// read meanwhile, by the goroutine holding the lock too.
func (o *EStoreEObjectImpl) Lock() {
	o.mutex.Lock()
	o.lockedLists = lockStoreLists(o.state.Load().getProperties())
}

func (o *EStoreEObjectImpl) Unlock() {
	unlockStoreLists(o.lockedLists)
	o.lockedLists = nil
	o.mutex.Unlock()
}

func (o *EStoreEObjectImpl) GetEStore() EStore {
	return o.state.Load().getStore()
}

func (o *EStoreEObjectImpl) SetEStore(newStore EStore) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.setEStoreLocked(newStore)
}

// setEStoreLocked sets the store of the object: mutex must be held, by Lock or SetEStore.
func (o *EStoreEObjectImpl) setEStoreLocked(newStore EStore) {
	for {
		current := o.state.Load()
		oldStore := current.getStore()
		if oldStore == newStore {
			return
		}
		// set store to object and its children
		state := current.clone()
		state.store = newStore
		if newStore == nil {
			// build cache with previous store
			if !state.cache {
				properties := slices.Clone(state.properties)
				if properties == nil {
					properties = make([]any, o.EClass().GetFeatureCount())
				}
				for featureID := range properties {
					if eFeature := o.eDynamicFeature(featureID); !eFeature.IsTransient() && !eFeature.IsMany() {
						properties[featureID] = oldStore.Get(o.AsEObject(), eFeature, NO_INDEX)
					}
				}
				state.properties = properties
			}
		} else if !state.cache {
			// clear properties because we are not caching
			state.properties = nil
		}
		if !o.state.CompareAndSwap(current, state) {
			// a value read from the store was kept meanwhile
			continue
		}
		if newStore != nil {
			// set container in store
			o.setContainerInStore()

			// set children collection store
			for _, v := range current.getProperties() {
				if storeList, isStoreList := v.(*EStoreList); isStoreList {
					if slices.Contains(o.lockedLists, storeList) {
						// locked with the object by Lock
						storeList.setEStoreLocked(newStore)
					} else {
						storeList.SetEStore(newStore)
					}
				}
			}
		}
		return
	}
}

func (o *EStoreEObjectImpl) SetCache(cache bool) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	for {
		current := o.state.Load()
		if current.isCache() == cache {
			return
		}
		state := current.clone()
		state.cache = cache
		if state.store != nil {
			state.properties = nil
		}
		if !o.state.CompareAndSwap(current, state) {
			// a value read from the store was kept meanwhile
			continue
		}
		// set cache for all properties
		for _, v := range current.getProperties() {
			if sc, _ := v.(ECacheProvider); sc != nil {
				sc.SetCache(cache)
			}
		}
		return
	}
}

func (o *EStoreEObjectImpl) IsCache() bool {
	return o.state.Load().isCache()
}

func (o *EStoreEObjectImpl) ESetInternalContainer(newContainer EObject, newContainerFeatureID int) {
	o.ReflectiveEObjectImpl.ESetInternalContainer(newContainer, newContainerFeatureID)
	o.setContainerInStore()
}

func (o *EStoreEObjectImpl) EInternalContainer() EObject {
	o.initializeContainerFromStore()
	return o.ReflectiveEObjectImpl.EInternalContainer()
}

func (o *EStoreEObjectImpl) EInternalContainerFeatureID() int {
	o.initializeContainerFromStore()
	return o.ReflectiveEObjectImpl.EInternalContainerFeatureID()
}

func (o *EStoreEObjectImpl) initializeContainerFromStore() {
	if o.ReflectiveEObjectImpl.EInternalContainer() == unitializedContainer {
		if store := o.state.Load().getStore(); store != nil {
			container, feature := store.GetContainer(o.AsEObject())
			if container != nil && feature != nil {
				featureID := EOPPOSITE_FEATURE_BASE - container.EClass().GetFeatureID(feature)
				if reference, _ := feature.(EReference); reference != nil {
					if opposite := reference.GetEOpposite(); opposite != nil {
						featureID = o.AsEObject().EClass().GetFeatureID(opposite)
					}
				}
				o.ReflectiveEObjectImpl.ESetInternalContainer(container, featureID)
			} else {
				o.ReflectiveEObjectImpl.ESetInternalContainer(nil, -1)
			}
		} else {
			o.ReflectiveEObjectImpl.ESetInternalContainer(nil, -1)
		}
	}
}

func (o *EStoreEObjectImpl) setContainerInStore() {
	if store := o.state.Load().getStore(); store != nil {
		container := o.ReflectiveEObjectImpl.EInternalContainer()
		containerFeatureID := o.ReflectiveEObjectImpl.EInternalContainerFeatureID()
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

// the store is read without holding the object : the store may wait for another goroutine
// that needs the object, as the store writer encoding it

func (o *EStoreEObjectImpl) EDynamicIsSet(dynamicFeatureID int) bool {
	state := o.state.Load()
	if state.getProperty(dynamicFeatureID) != nil {
		return true
	}
	if store := state.getStore(); store != nil {
		eFeature := o.eDynamicFeature(dynamicFeatureID)
		return store.IsSet(o.AsEObject(), eFeature)
	}
	return false
}

func (o *EStoreEObjectImpl) EDynamicGet(dynamicFeatureID int) any {
	// retrieve value
	state := o.state.Load()
	if result := state.getProperty(dynamicFeatureID); result != nil {
		return result
	}

	for {
		// compute value
		var result any
		cache := false
		store := state.getStore()
		objectCache := state.isCache()
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
		for {
			current := o.state.Load()
			if value := current.getProperty(dynamicFeatureID); value != nil {
				return value
			}
			if current.getStore() != store || current.isCache() != objectCache {
				// store or cache changed meanwhile: compute again
				state = current
				break
			}
			if o.state.CompareAndSwap(current, current.withProperty(dynamicFeatureID, result, o.EClass())) {
				return result
			}
		}
	}
}

func (o *EStoreEObjectImpl) EDynamicSet(dynamicFeatureID int, value any) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	eFeature := o.eDynamicFeature(dynamicFeatureID)
	store := o.state.Load().getStore()
	persistent := store != nil && !eFeature.IsTransient()
	if persistent {
		// store and feature is not transient
		store.Set(o.AsEObject(), eFeature, NO_INDEX, value, false)
	}
	for {
		current := o.state.Load()
		if persistent && !current.isCache() {
			// no cache : the value is in the store only
			return
		}
		// no store or feature is transient: cache value in properties even if there is no cache
		if o.state.CompareAndSwap(current, current.withProperty(dynamicFeatureID, value, o.EClass())) {
			return
		}
	}
}

func (o *EStoreEObjectImpl) EDynamicUnset(dynamicFeatureID int) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	// retrieve previous value and unset properties
	for {
		current := o.state.Load()
		if current.getProperty(dynamicFeatureID) == nil {
			break
		}
		if o.state.CompareAndSwap(current, current.withProperty(dynamicFeatureID, nil, nil)) {
			break
		}
	}
	eFeature := o.eDynamicFeature(dynamicFeatureID)
	store := o.state.Load().getStore()
	if store != nil && !eFeature.IsTransient() {
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
