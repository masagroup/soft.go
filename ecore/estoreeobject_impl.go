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
	"github.com/SokaDance/rmx"
)

var unitializedContainer EObject = newEObjectImpl()

type EStoreEObjectImpl struct {
	ReflectiveEObjectImpl
	store EStore
	cache bool
	mutex rmx.RecursiveMutex
}

func NewEStoreEObjectImpl(cache bool) *EStoreEObjectImpl {
	o := new(EStoreEObjectImpl)
	o.cache = cache
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

func (o *EStoreEObjectImpl) Lock() {
	o.mutex.Lock()
	for _, v := range o.properties {
		if storeList, isStoreList := v.(*EStoreList); isStoreList {
			storeList.Lock()
		}
	}
}

func (o *EStoreEObjectImpl) Unlock() {
	for _, v := range o.properties {
		if storeList, isStoreList := v.(*EStoreList); isStoreList {
			storeList.Unlock()
		}
	}
	o.mutex.Unlock()
}

func (o *EStoreEObjectImpl) GetEStore() EStore {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	return o.store
}

func (o *EStoreEObjectImpl) SetEStore(newStore EStore) {
	o.mutex.Lock()
	if oldStore := o.store; newStore != oldStore {
		// set store to object and its children
		o.store = newStore
		if newStore == nil {
			// build cache with previous store
			if !o.cache {
				for featureID := range o.getProperties() {
					if eFeature := o.eDynamicFeature(featureID); !eFeature.IsTransient() && !eFeature.IsMany() {
						o.properties[featureID] = oldStore.Get(o.AsEObject(), eFeature, NO_INDEX)
					}
				}
			}
		} else {
			// set container in store
			o.setContainerInStore()

			// set children collection store
			for _, v := range o.properties {
				if storeList, isStoreList := v.(*EStoreList); isStoreList {
					storeList.SetEStore(newStore)
				}
			}

			// clear properties because we are not caching
			if !o.cache {
				o.properties = nil
			}
		}
	}
	o.mutex.Unlock()
}

func (o *EStoreEObjectImpl) SetCache(cache bool) {
	o.mutex.Lock()
	if o.cache != cache {
		o.cache = cache

		// set cache for all properties
		for _, v := range o.properties {
			if sc, _ := v.(ECacheProvider); sc != nil {
				sc.SetCache(cache)
			}
		}
		if o.store != nil {
			o.properties = nil
		}
	}
	o.mutex.Unlock()
}

func (o *EStoreEObjectImpl) IsCache() bool {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	return o.cache
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
		if o.store != nil {
			container, feature := o.store.GetContainer(o.AsEObject())
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
	if o.store != nil {
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
			o.store.SetContainer(o.AsEObject(), container, containerFeature)
		}
	}
}

// the store is read without holding the object : the store may wait for another goroutine
// that needs the object, as the store writer encoding it

func (o *EStoreEObjectImpl) EDynamicIsSet(dynamicFeatureID int) bool {
	o.mutex.Lock()
	if o.properties != nil && o.properties[dynamicFeatureID] != nil {
		o.mutex.Unlock()
		return true
	}
	o.mutex.Unlock()

	if store := o.AsEStoreEObject().GetEStore(); store != nil {
		eFeature := o.eDynamicFeature(dynamicFeatureID)
		return store.IsSet(o.AsEObject(), eFeature)
	}
	return false
}

func (o *EStoreEObjectImpl) EDynamicGet(dynamicFeatureID int) any {
	// retrieve value
	o.mutex.Lock()
	if o.properties != nil {
		if result := o.properties[dynamicFeatureID]; result != nil {
			o.mutex.Unlock()
			return result
		}
	}
	o.mutex.Unlock()

	// compute value
	var result any
	cache := false
	objectCache := o.IsCache()
	if eFeature := o.eDynamicFeature(dynamicFeatureID); !eFeature.IsTransient() {
		if eFeature.IsMany() {
			if IsMapType(eFeature) {
				result = o.createMap(eFeature, objectCache)
			} else {
				result = o.createList(eFeature, objectCache)
			}
			// feature is a container : create corresponding container type and cache it
			// in properties to avoid multiple allocation of the same container
			cache = true
		} else if store := o.AsEStoreEObject().GetEStore(); store != nil {
			// feature is not transient and we have a store
			result = store.Get(o.AsEObject(), eFeature, NO_INDEX)
			cache = objectCache
		}
	}

	// store value in properties, unless set meanwhile
	if cache {
		o.mutex.Lock()
		defer o.mutex.Unlock()
		properties := o.getProperties()
		if current := properties[dynamicFeatureID]; current != nil {
			return current
		}
		if o.cache != objectCache {
			// cache changed meanwhile
			if cacheProvider, _ := result.(ECacheProvider); cacheProvider != nil {
				cacheProvider.SetCache(o.cache)
			} else if !o.cache {
				return result
			}
		}
		properties[dynamicFeatureID] = result
	}
	return result
}

func (o *EStoreEObjectImpl) EDynamicSet(dynamicFeatureID int, value any) {
	o.mutex.Lock()
	// retrieve properties
	var properties []any
	eFeature := o.eDynamicFeature(dynamicFeatureID)
	store := o.AsEStoreEObject().GetEStore()
	if store != nil && !eFeature.IsTransient() {
		// store and feature is not transient
		store.Set(o.AsEObject(), eFeature, NO_INDEX, value, false)
		if o.cache {
			properties = o.getProperties()
		}
	} else {
		// no store or feature is transient
		// cache value in properties even if there is no cache
		properties = o.getProperties()
	}
	// store value in properties
	if properties != nil {
		properties[dynamicFeatureID] = value
	}
	o.mutex.Unlock()
}

func (o *EStoreEObjectImpl) EDynamicUnset(dynamicFeatureID int) {
	o.mutex.Lock()
	// retrieve previous value and unset properties
	if o.properties != nil {
		o.properties[dynamicFeatureID] = nil
	}
	eFeature := o.eDynamicFeature(dynamicFeatureID)
	store := o.AsEStoreEObject().GetEStore()
	if store != nil && !eFeature.IsTransient() {
		store.UnSet(o.AsEObject(), eFeature)
	}
	o.mutex.Unlock()
}

func (o *EStoreEObjectImpl) eDynamicFeature(dynamicFeatureID int) EStructuralFeature {
	return o.EClass().GetEStructuralFeature(o.EStaticFeatureCount() + dynamicFeatureID)
}

func (o *EStoreEObjectImpl) createList(eFeature EStructuralFeature, cache bool) EList {
	l := NewEStoreList(o.AsEObject(), eFeature, o.AsEStoreEObject().GetEStore())
	l.SetCache(cache)
	return l
}

func (o *EStoreEObjectImpl) createMap(eFeature EStructuralFeature, cache bool) EMap {
	eClass := eFeature.GetEType().(EClass)
	m := NewEStoreMap(eClass, o.AsEObject(), eFeature, o.AsEStoreEObject().GetEStore())
	m.SetCache(cache)
	return m
}
