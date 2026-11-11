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
	"reflect"
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
// It uses CompactEObjectImpl, a container published atomically and a dynamic compact strong cache.
type CompactEStoreEObjectImpl struct {
	CompactEObjectImpl
	class EClass
	// cache holds the cached values. It is never changed, but replaced with a compare-and-swap:
	// a value is read from it without lock, at each access to a feature, and a value read from
	// the store is cached without lock. The object takes no lock: a concrete type that needs one
	// (ELockProvider for the SQL encoder, ordered writes) provides it.
	cache atomic.Pointer[compactCache]
	// container is the container of the object and its feature ID, nil while not known: it is
	// then read from the store by the first reader, and published with a compare-and-swap, as
	// readers do not take a lock. An object without store keeps it unknown until its container
	// is set, so that a later store is asked, and not given a nil container.
	container atomic.Pointer[compactContainer]
}

// compactContainer is a container of a CompactEStoreEObjectImpl. It is never changed.
type compactContainer struct {
	container EObject
	featureID int
}

var noContainer = &compactContainer{featureID: -1}

func (o *CompactEStoreEObjectImpl) Initialize() {
	o.CompactEObjectImpl.Initialize()
	// the container is not known yet (nil)
	o.container.Store(nil)
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

// setFeatureValue returns values with the value of featureID, without it if nil.
func setFeatureValue(values []cachedFeatureValue, featureID int, value any) []cachedFeatureValue {
	index := slices.IndexFunc(values, func(v cachedFeatureValue) bool { return v.featureID == featureID })
	switch {
	case index >= 0 && value == nil:
		return slices.Delete(values, index, index+1)
	case index >= 0:
		values[index].value = value
		return values
	case value != nil:
		return append(values, cachedFeatureValue{featureID, value})
	}
	return values
}

// cachedValues returns the cached values.
func (o *CompactEStoreEObjectImpl) cachedValues() []any {
	featureValues := o.cache.Load().featureValues()
	values := make([]any, len(featureValues))
	for i, v := range featureValues {
		values[i] = v.value
	}
	return values
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
	if cached := o.cachedValue(dynamicFeatureID); cached != nil {
		return isSetValue(o.EClass().GetEStructuralFeature(dynamicFeatureID), cached)
	}

	feature := o.EClass().GetEStructuralFeature(dynamicFeatureID)
	if !feature.IsTransient() {
		if store := o.getStore(); store != nil {
			return store.IsSet(o.AsEObject(), feature)
		}
	}
	return false
}

// isSetValue reports whether the cached value of a feature is set, as a store reports it once
// written: a value of an attribute, not unsettable, equal to its default is not set. A value read
// from the store and cached is so not set by its read.
func isSetValue(feature EStructuralFeature, value any) bool {
	if value == nil {
		return false
	}
	if value == NIL {
		return true
	}
	if _, isAttribute := feature.(EAttribute); !isAttribute || feature.IsMany() || feature.IsUnsettable() {
		return true
	}
	return !reflect.DeepEqual(value, feature.GetDefaultValue())
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
	o.container.Store(&compactContainer{container: newContainer, featureID: newContainerFeatureID})
	o.setContainerInStore()
}

func (o *CompactEStoreEObjectImpl) EInternalContainer() EObject {
	return o.loadContainer().container
}

func (o *CompactEStoreEObjectImpl) EInternalContainerFeatureID() int {
	return o.loadContainer().featureID
}

// loadContainer returns the container of the object, read from the store while not known. It
// writes nothing but the published container, with a compare-and-swap: a container set
// meanwhile wins.
func (o *CompactEStoreEObjectImpl) loadContainer() *compactContainer {
	if c := o.container.Load(); c != nil {
		return c
	}
	store := o.getStore()
	if store == nil {
		// unknown until set, or until a store knows it
		return noContainer
	}
	loaded := noContainer
	if container, feature := store.GetContainer(o.AsEObject()); container != nil && feature != nil {
		featureID := EOPPOSITE_FEATURE_BASE - container.EClass().GetFeatureID(feature)
		if reference, _ := feature.(EReference); reference != nil {
			if opposite := reference.GetEOpposite(); opposite != nil {
				featureID = o.AsEObject().EClass().GetFeatureID(opposite)
			}
		}
		loaded = &compactContainer{container: container, featureID: featureID}
	}
	if o.container.CompareAndSwap(nil, loaded) {
		return loaded
	}
	return o.container.Load()
}

// setContainerInStore writes the container of the object in its store, when it is known: an
// unknown container is the one the store has.
func (o *CompactEStoreEObjectImpl) setContainerInStore() {
	if store := o.getStore(); store != nil {
		c := o.container.Load()
		if c == nil {
			return
		}
		var containerFeature EStructuralFeature
		if container := c.container; container != nil {
			if c.featureID <= EOPPOSITE_FEATURE_BASE {
				containerFeature = container.EClass().GetEStructuralFeature(EOPPOSITE_FEATURE_BASE - c.featureID)
			} else {
				containerFeature = o.AsEObject().EClass().GetEStructuralFeature(c.featureID).(EReference).GetEOpposite()
			}
		}
		store.SetContainer(o.AsEObject(), c.container, containerFeature)
	}
}
