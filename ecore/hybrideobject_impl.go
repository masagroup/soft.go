// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import "unsafe"

// HybridEObjectImpl is the root of the hybrid classes over EStoreEObjectImpl: classes whose values
// are kept in a store, as with reflective feature delegation, and whose getters read typed fields,
// a cache of the values of the store (soft.gen go, property hybridFields).
//
// The fields of a class are in a block of their own (Fields), not in the object: a field is
// written once, then published, and never written again, so that the getters read it without
// lock. A generated class declares its fields (NewLevel, Value, List) and its Layout, returned by
// its ELayout. The changes of the values drop the fields they change: a set publishes the value
// set, an unset drops its field, a new store or a cleared cache drops them all. A value changed
// in the store by another way is read again after ClearCache.
type HybridEObjectImpl struct {
	EStoreEObjectImpl
}

var _ LayoutProvider = (*HybridEObjectImpl)(nil)

// ELayout returns nil: the root has no field. A generated class with fields overrides it.
func (o *HybridEObjectImpl) ELayout() *Layout { return nil }

// layout returns the layout of the class of the object.
func (o *HybridEObjectImpl) layout() *Layout {
	if h, ok := o.GetInterfaces().(LayoutProvider); ok {
		return h.ELayout()
	}
	return nil
}

// EHybridBase returns the address of the root in its object (HybridRoot).
func (o *HybridEObjectImpl) EHybridBase() unsafe.Pointer { return unsafe.Pointer(o) }

// ELoadFeature returns the value of the single-valued feature featureID, feature, for its typed
// getter: the value of the store, the default value of the feature when it has none, or the
// value the object keeps without store.
func (o *HybridEObjectImpl) ELoadFeature(featureID int, feature EStructuralFeature) any {
	if o.GetEStore() == nil {
		// the reflective read of the object, without the fields, which would load again
		return o.EStoreEObjectImpl.EGetFromID(featureID, true)
	}
	value := o.EStoreEObjectImpl.EDynamicGet(featureID)
	if value == nil {
		return feature.GetDefaultValue()
	}
	if value == NIL {
		return nil
	}
	if object, _ := value.(EObject); object != nil && object.EIsProxy() && IsProxy(feature) {
		return o.EResolveProxy(object)
	}
	return value
}

// ELoadList returns the list of the many-valued feature featureID: the one the object keeps, made
// on first read.
func (o *HybridEObjectImpl) ELoadList(featureID int) EList {
	l, _ := o.EStoreEObjectImpl.EDynamicGet(featureID).(EList)
	return l
}

// ESetsUnchanged reports whether setting a feature to its current value changes nothing: the
// object has a store and nothing is notified.
func (o *HybridEObjectImpl) ESetsUnchanged() bool {
	return o.GetEStore() != nil && !o.ENotificationRequired()
}

// ESetField sets a plain feature for its typed setter: straight to EDynamicSet, without the
// reflective lookups and the read of the old value, unless the value is nil (stored by the rules
// of its feature) or the change is notified.
func (o *HybridEObjectImpl) ESetField(featureID int, value any) {
	if value == nil || o.ENotificationRequired() {
		o.EStoreEObjectImpl.ESetFromID(featureID, value)
		return
	}
	o.EDynamicSet(featureID, value)
}

// EBaseGetFromID is the reflective read of the object without the fields.
func (o *HybridEObjectImpl) EBaseGetFromID(featureID int, resolve bool) any {
	return o.EStoreEObjectImpl.EGetFromID(featureID, resolve)
}

// EBaseSetFromID is the reflective write of the object without the fields.
func (o *HybridEObjectImpl) EBaseSetFromID(featureID int, value any) {
	o.EStoreEObjectImpl.ESetFromID(featureID, value)
}

// EBaseIsSetFromID is the reflective isSet of the object without the fields.
func (o *HybridEObjectImpl) EBaseIsSetFromID(featureID int) bool {
	return o.EStoreEObjectImpl.EIsSetFromID(featureID)
}

// EGetFromID returns the value of the feature featureID, the one of its typed field if it has
// one. A generated class with fields overrides it.
func (o *HybridEObjectImpl) EGetFromID(featureID int, resolve bool) any {
	return o.layout().EGetFromID(o, featureID, resolve)
}

// ESetFromID sets the feature featureID, with its typed setter if it has one. A generated class
// with fields overrides it.
func (o *HybridEObjectImpl) ESetFromID(featureID int, value any) {
	o.layout().ESetFromID(o, featureID, value)
}

// EIsSetFromID reports whether the feature featureID is set, from its typed list if it has one. A
// generated class with fields overrides it.
func (o *HybridEObjectImpl) EIsSetFromID(featureID int) bool {
	return o.layout().EIsSetFromID(o, featureID)
}

// EDynamicGet reads a feature value: the one of its typed field for an object with a store.
func (o *HybridEObjectImpl) EDynamicGet(dynamicFeatureID int) any {
	if o.GetEStore() != nil {
		if value, ok := o.layout().Value(o, dynamicFeatureID); ok {
			return value
		}
	}
	return o.EStoreEObjectImpl.EDynamicGet(dynamicFeatureID)
}

// EDynamicSet sets a feature value, then publishes it in its typed field.
func (o *HybridEObjectImpl) EDynamicSet(dynamicFeatureID int, value any) {
	o.EStoreEObjectImpl.EDynamicSet(dynamicFeatureID, value)
	o.layout().Replace(o, dynamicFeatureID, value)
}

// EDynamicUnset unsets a feature value, then drops its typed field.
func (o *HybridEObjectImpl) EDynamicUnset(dynamicFeatureID int) {
	o.EStoreEObjectImpl.EDynamicUnset(dynamicFeatureID)
	o.layout().ClearFeatures(o, []int{dynamicFeatureID})
}

// SetEStore sets the store of the object, then drops the typed fields.
func (o *HybridEObjectImpl) SetEStore(store EStore) {
	o.EStoreEObjectImpl.SetEStore(store)
	o.layout().ClearAll(o)
}

// SetEStoreLocked sets the store of the object locked, then drops the typed fields.
func (o *HybridEObjectImpl) SetEStoreLocked(store EStore) {
	o.EStoreEObjectImpl.SetEStoreLocked(store)
	o.layout().ClearAll(o)
}

// SetCache sets whether the object caches its values, then drops the typed fields.
func (o *HybridEObjectImpl) SetCache(cache bool) {
	o.EStoreEObjectImpl.SetCache(cache)
	o.layout().ClearAll(o)
}

// ClearCache clears the cached values and the typed fields.
func (o *HybridEObjectImpl) ClearCache() {
	o.EStoreEObjectImpl.ClearCache()
	o.layout().ClearAll(o)
}
