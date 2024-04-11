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
	isCaching bool
	store     EStore
}

func NewEStoreEObjectImpl(isCaching bool) *EStoreEObjectImpl {
	o := new(EStoreEObjectImpl)
	o.isCaching = isCaching
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
		for _, v := range o.getProperties() {
			if storeProvider, _ := v.(EStoreProvider); storeProvider != nil {
				storeProvider.SetEStore(newStore)
			}
		}

		if newStore == nil {
			// build cache with previous store
			if !o.isCaching {
				eClass := o.AsEObject().EClass()
				properties := o.getProperties()
				for featureID := 0; featureID < len(properties); featureID++ {
					eFeature := eClass.GetEStructuralFeature(featureID)
					properties[featureID] = oldStore.Get(o.AsEObject(), eFeature, NO_INDEX)
				}
			}
		} else {
			// initialize store with cache
			eClass := o.AsEObject().EClass()
			properties := o.getProperties()
			for featureID, value := range properties {
				if o.EDynamicIsSet(featureID) {
					eFeature := eClass.GetEStructuralFeature(featureID)
					newStore.Set(o.AsEObject(), eFeature, NO_INDEX, value)
				}
			}

			// clear properties because we are not caching
			if !o.isCaching {
				o.clearProperties()
			}
		}

	}
}

func (o *EStoreEObjectImpl) SetCaching(isCaching bool) {
	o.isCaching = isCaching
}

func (o *EStoreEObjectImpl) IsCaching() bool {
	return o.isCaching
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
				if o.isCaching {
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
		if o.isCaching {
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
	return NewEStoreList(o.AsEObject(), eFeature, o.AsEStoreEObject().GetEStore())
}

func (o *EStoreEObjectImpl) createMap(eFeature EStructuralFeature) EMap {
	eClass := eFeature.GetEType().(EClass)
	return NewEStoreMap(eClass, o.AsEObject(), eFeature, o.AsEStoreEObject().GetEStore())
}
