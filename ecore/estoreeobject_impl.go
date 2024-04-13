// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

type EStoreEObjectImpl struct {
	ReflectiveEObjectImpl
	cache bool
	store EStore
}

func NewEStoreEObjectImpl(cache bool) *EStoreEObjectImpl {
	o := new(EStoreEObjectImpl)
	o.cache = cache
	o.SetInterfaces(o)
	o.Initialize()
	return o
}

func (o *EStoreEObjectImpl) AsEStoreEObject() EStoreEObject {
	return o.GetInterfaces().(EStoreEObject)
}

func (o *EStoreEObjectImpl) GetEStore() EStore {
	return o.store
}

func (o *EStoreEObjectImpl) SetEStore(newStore EStore) {
	if oldStore := o.store; newStore != oldStore {
		// set store to object and its children
		o.store = newStore
		if newStore == nil {
			// build cache with previous store
			if !o.cache {
				for featureID := 0; featureID < len(o.properties); featureID++ {
					if eFeature := o.eDynamicFeature(featureID); !eFeature.IsTransient() {
						o.properties[featureID] = oldStore.Get(o.AsEObject(), eFeature, NO_INDEX)
					}
				}
			}
		} else {
			// set children store
			for _, v := range o.getProperties() {
				if sp, _ := v.(EStoreProvider); sp != nil {
					sp.SetEStore(newStore)
				}
			}

			// clear properties because we are not caching
			if !o.cache {
				o.clearProperties()
			}
		}

	}
}

func (o *EStoreEObjectImpl) SetCache(cache bool) {
	o.cache = cache
}

func (o *EStoreEObjectImpl) IsCache() bool {
	return o.cache
}

func (o *EStoreEObjectImpl) ClearCache() {
	o.clearProperties()
}

func (o *EStoreEObjectImpl) EDynamicGet(dynamicFeatureID int) any {
	result := o.getProperties()[dynamicFeatureID]
	if result == nil {
		eFeature := o.eDynamicFeature(dynamicFeatureID)
		if !eFeature.IsTransient() {
			if eFeature.IsMany() {
				if IsMapType(eFeature) {
					result = o.createMap(eFeature)
				} else {
					result = o.createList(eFeature)
				}
				o.getProperties()[dynamicFeatureID] = result
			} else if store := o.AsEStoreEObject().GetEStore(); store != nil {
				result = store.Get(o.AsEObject(), eFeature, NO_INDEX)
				if o.cache {
					o.getProperties()[dynamicFeatureID] = result
				}
			}
		}
	}
	return result
}

func (o *EStoreEObjectImpl) EDynamicSet(dynamicFeatureID int, value any) {
	eFeature := o.eDynamicFeature(dynamicFeatureID)
	if eFeature.IsTransient() {
		o.getProperties()[dynamicFeatureID] = value
	} else if store := o.AsEStoreEObject().GetEStore(); store != nil {
		store.Set(o.AsEObject(), eFeature, NO_INDEX, value)
		if o.cache {
			o.getProperties()[dynamicFeatureID] = value
		}
	} else {
		o.getProperties()[dynamicFeatureID] = value
	}
}

func (o *EStoreEObjectImpl) EDynamicUnset(dynamicFeatureID int) {
	o.getProperties()[dynamicFeatureID] = nil
	if eFeature := o.eDynamicFeature(dynamicFeatureID); !eFeature.IsTransient() {
		if store := o.AsEStoreEObject().GetEStore(); store != nil {
			store.UnSet(o.AsEObject(), eFeature)
		}
	}
}

func (o *EStoreEObjectImpl) eDynamicFeature(dynamicFeatureID int) EStructuralFeature {
	return o.EClass().GetEStructuralFeature(o.EStaticFeatureCount() + dynamicFeatureID)
}

func (o *EStoreEObjectImpl) createList(eFeature EStructuralFeature) EList {
	l := NewEStoreList(o.AsEObject(), eFeature, o.AsEStoreEObject().GetEStore())
	l.SetCache(o.cache)
	l.SetEStore(o.store)
	return l
}

func (o *EStoreEObjectImpl) createMap(eFeature EStructuralFeature) EMap {
	eClass := eFeature.GetEType().(EClass)
	m := NewEStoreMap(eClass, o.AsEObject(), eFeature, o.AsEStoreEObject().GetEStore())
	m.SetCache(o.cache)
	m.SetEStore(o.store)
	return m
}
