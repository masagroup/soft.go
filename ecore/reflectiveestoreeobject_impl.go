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

type internalReflectiveEStoreEObjectImpl interface {
	GetEStore() EStore
}

// ReflectiveEStoreEObjectImpl is an abstract reflective EObject implementation.
// It defines no store fields; concrete structs must implement GetEStore() EStore.
// All values (persistent cache and transient values) are held as strong references.
type ReflectiveEStoreEObjectImpl struct {
	ReflectiveEObjectImpl
	mutex rmx.RecursiveMutex
}

func (o *ReflectiveEStoreEObjectImpl) Initialize() {
	o.ReflectiveEObjectImpl.Initialize()
	o.ESetInternalContainer(unitializedContainer, -1)
}

func (o *ReflectiveEStoreEObjectImpl) asInternal() internalReflectiveEStoreEObjectImpl {
	return o.GetInterfaces().(internalReflectiveEStoreEObjectImpl)
}

func (o *ReflectiveEStoreEObjectImpl) getStore() EStore {
	return o.asInternal().GetEStore()
}

func (o *ReflectiveEStoreEObjectImpl) Lock() {
	o.mutex.Lock()
	for _, v := range o.properties {
		if storeList, isStoreList := v.(*EStoreList); isStoreList {
			storeList.Lock()
		}
	}
}

func (o *ReflectiveEStoreEObjectImpl) Unlock() {
	for _, v := range o.properties {
		if storeList, isStoreList := v.(*EStoreList); isStoreList {
			storeList.Unlock()
		}
	}
	o.mutex.Unlock()
}

func (o *ReflectiveEStoreEObjectImpl) EDynamicGet(dynamicFeatureID int) any {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	// 1. Cache hit
	var result any
	if o.properties != nil {
		result = o.properties[dynamicFeatureID]
	}

	// 2. Cache miss
	if result == nil {
		var properties []any
		if feature := o.eDynamicFeature(dynamicFeatureID); !feature.IsTransient() {
			if feature.IsMany() {
				if IsMapType(feature) {
					result = o.createMap(feature, o.getStore())
				} else {
					result = o.createList(feature, o.getStore())
				}
				properties = o.getProperties()
			} else if store := o.getStore(); store != nil {
				result = store.Get(o.AsEObject(), feature, NO_INDEX)
				properties = o.getProperties()
			}
		} else if feature.IsMany() {
			if IsMapType(feature) {
				result = o.createMap(feature, nil)
			} else {
				result = o.createList(feature, nil)
			}
			properties = o.getProperties()
		}

		if properties != nil {
			properties[dynamicFeatureID] = result
		}
	}

	return result
}

func (o *ReflectiveEStoreEObjectImpl) EDynamicSet(dynamicFeatureID int, value any) {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	feature := o.eDynamicFeature(dynamicFeatureID)

	if !feature.IsTransient() {
		if store := o.getStore(); store != nil {
			store.Set(o.AsEObject(), feature, NO_INDEX, value, false)
		}
	}

	o.getProperties()[dynamicFeatureID] = value
}

func (o *ReflectiveEStoreEObjectImpl) EDynamicIsSet(dynamicFeatureID int) bool {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	if o.properties != nil && o.properties[dynamicFeatureID] != nil {
		return true
	}

	feature := o.eDynamicFeature(dynamicFeatureID)
	if !feature.IsTransient() {
		if store := o.getStore(); store != nil {
			return store.IsSet(o.AsEObject(), feature)
		}
	}
	return false
}

func (o *ReflectiveEStoreEObjectImpl) EDynamicUnset(dynamicFeatureID int) {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	if o.properties != nil {
		o.properties[dynamicFeatureID] = nil
	}

	feature := o.eDynamicFeature(dynamicFeatureID)
	if !feature.IsTransient() {
		if store := o.getStore(); store != nil {
			store.UnSet(o.AsEObject(), feature)
		}
	}
}

func (o *ReflectiveEStoreEObjectImpl) eDynamicFeature(dynamicFeatureID int) EStructuralFeature {
	return o.EClass().GetEStructuralFeature(o.EStaticFeatureCount() + dynamicFeatureID)
}

func (o *ReflectiveEStoreEObjectImpl) createList(feature EStructuralFeature, store EStore) EList {
	l := NewEStoreList(o.AsEObject(), feature, store)
	l.SetCache(true)
	return l
}

func (o *ReflectiveEStoreEObjectImpl) createMap(feature EStructuralFeature, store EStore) EMap {
	eClass := feature.GetEType().(EClass)
	return NewEStoreMap(eClass, o.AsEObject(), feature, store)
}

func (o *ReflectiveEStoreEObjectImpl) ESetInternalContainer(newContainer EObject, newContainerFeatureID int) {
	o.ReflectiveEObjectImpl.ESetInternalContainer(newContainer, newContainerFeatureID)
	o.setContainerInStore()
}

func (o *ReflectiveEStoreEObjectImpl) EInternalContainer() EObject {
	o.initializeContainerFromStore()
	return o.ReflectiveEObjectImpl.EInternalContainer()
}

func (o *ReflectiveEStoreEObjectImpl) EInternalContainerFeatureID() int {
	o.initializeContainerFromStore()
	return o.ReflectiveEObjectImpl.EInternalContainerFeatureID()
}

func (o *ReflectiveEStoreEObjectImpl) initializeContainerFromStore() {
	if o.ReflectiveEObjectImpl.EInternalContainer() == unitializedContainer {
		if store := o.getStore(); store != nil {
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

func (o *ReflectiveEStoreEObjectImpl) setContainerInStore() {
	if store := o.getStore(); store != nil {
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
