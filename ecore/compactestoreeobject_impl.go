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

	"github.com/SokaDance/rmx"
)

type internalCompactEStoreEObjectImpl interface {
	GetEStore() EStore
}

// CompactEStoreEObjectImpl is an abstract compact reflective EObject implementation.
// It defines no store fields; concrete structs must implement GetEStore() EStore.
// It uses CompactEObjectContainer and a dynamic compact strong cache.
type CompactEStoreEObjectImpl struct {
	CompactEObjectContainer
	class        EClass
	cachedValues any    // nil, single value, or []any
	cacheMask    uint64 // bitmask of cached feature IDs (0..63)
	mutex        rmx.RecursiveMutex
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

func (o *CompactEStoreEObjectImpl) isCached(featureID int) bool {
	return featureID < 64 && (o.cacheMask&(1<<featureID)) != 0
}

func (o *CompactEStoreEObjectImpl) getCached(featureID int) any {
	if !o.isCached(featureID) {
		return nil
	}
	if bits.OnesCount64(o.cacheMask) == 1 {
		return o.cachedValues
	}
	idx := bits.OnesCount64(o.cacheMask & ((1 << featureID) - 1))
	return o.cachedValues.([]any)[idx]
}

func (o *CompactEStoreEObjectImpl) setCached(featureID int, value any) {
	if featureID >= 64 {
		return
	}
	bit := uint64(1) << featureID
	if (o.cacheMask & bit) == 0 {
		count := bits.OnesCount64(o.cacheMask)
		if count == 0 {
			o.cachedValues = value
		} else if count == 1 {
			idx := bits.OnesCount64(o.cacheMask & (bit - 1))
			if idx == 0 {
				o.cachedValues = []any{value, o.cachedValues}
			} else {
				o.cachedValues = []any{o.cachedValues, value}
			}
		} else {
			oldSlice := o.cachedValues.([]any)
			idx := bits.OnesCount64(o.cacheMask & (bit - 1))
			newSlice := make([]any, count+1)
			copy(newSlice[:idx], oldSlice[:idx])
			newSlice[idx] = value
			copy(newSlice[idx+1:], oldSlice[idx:])
			o.cachedValues = newSlice
		}
		o.cacheMask |= bit
	} else {
		if bits.OnesCount64(o.cacheMask) == 1 {
			o.cachedValues = value
		} else {
			idx := bits.OnesCount64(o.cacheMask & (bit - 1))
			o.cachedValues.([]any)[idx] = value
		}
	}
}

func (o *CompactEStoreEObjectImpl) unsetCached(featureID int) {
	if !o.isCached(featureID) {
		return
	}
	bit := uint64(1) << featureID
	count := bits.OnesCount64(o.cacheMask)
	if count == 1 {
		o.cachedValues = nil
	} else if count == 2 {
		idx := bits.OnesCount64(o.cacheMask & (bit - 1))
		o.cachedValues = o.cachedValues.([]any)[1-idx]
	} else {
		idx := bits.OnesCount64(o.cacheMask & (bit - 1))
		s := o.cachedValues.([]any)
		o.cachedValues = append(s[:idx], s[idx+1:]...)
	}
	o.cacheMask &= ^bit
}

func (o *CompactEStoreEObjectImpl) EDynamicGet(dynamicFeatureID int) any {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	// 1. Cache hit
	result := o.getCached(dynamicFeatureID)

	// 2. Cache miss
	if result == nil {
		shouldCache := false
		feature := o.EClass().GetEStructuralFeature(dynamicFeatureID)
		if !feature.IsTransient() {
			if feature.IsMany() {
				if IsMapType(feature) {
					result = o.createMap(feature, o.getStore())
				} else {
					result = o.createList(feature, o.getStore())
				}
				shouldCache = true
			} else if store := o.getStore(); store != nil {
				result = store.Get(o.AsEObject(), feature, NO_INDEX)
				shouldCache = true
			}
		} else if feature.IsMany() {
			if IsMapType(feature) {
				result = o.createMap(feature, nil)
			} else {
				result = o.createList(feature, nil)
			}
			shouldCache = true
		}

		if shouldCache && result != nil {
			o.setCached(dynamicFeatureID, result)
		}
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
	o.mutex.Lock()
	defer o.mutex.Unlock()

	if o.getCached(dynamicFeatureID) != nil {
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

func (o *CompactEStoreEObjectImpl) createList(feature EStructuralFeature, store EStore) EList {
	l := NewEStoreList(o.AsEObject(), feature, store)
	l.SetCache(true)
	return l
}

func (o *CompactEStoreEObjectImpl) createMap(feature EStructuralFeature, store EStore) EMap {
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
