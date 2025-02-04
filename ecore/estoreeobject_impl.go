// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

var unitializedContainer EObject = newEObjectImpl()

type EStoreEObjectImpl struct {
	ReflectiveEObjectImpl
	store EStore
	cache bool
}

func NewEStoreEObjectImpl(cache bool) *EStoreEObjectImpl {
	o := new(EStoreEObjectImpl)
	o.cache = cache
	o.SetInterfaces(o)
	o.Initialize()
	o.ESetInternalContainer(unitializedContainer, -1)
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
			for _, v := range o.properties {
				if sp, _ := v.(EStoreProvider); sp != nil {
					sp.SetEStore(newStore)
				}
			}

			// clear properties because we are not caching
			if !o.cache {
				o.properties = nil
			}
		}
	}
}

func (o *EStoreEObjectImpl) SetCache(cache bool) {
	if o.cache != cache {
		o.cache = cache

		// set cache for all properties
		for _, v := range o.properties {
			if sc, _ := v.(ECacheProvider); sc != nil {
				sc.SetCache(cache)
			}
		}
		o.properties = nil
	}
}

func (o *EStoreEObjectImpl) IsCache() bool {
	return o.cache
}

func (o *EStoreEObjectImpl) EInternalContainer() EObject {
	o.initializeContainer()
	return o.ReflectiveEObjectImpl.EInternalContainer()
}

func (o *EStoreEObjectImpl) EInternalContainerFeatureID() int {
	o.initializeContainer()
	return o.ReflectiveEObjectImpl.EInternalContainerFeatureID()
}

func (o *EStoreEObjectImpl) ESetInternalContainer(newContainer EObject, newContainerFeatureID int) {
	o.ReflectiveEObjectImpl.ESetInternalContainer(newContainer, newContainerFeatureID)
	if o.store != nil {
		var containerFeature EStructuralFeature
		if newContainerFeatureID <= EOPPOSITE_FEATURE_BASE {
			containerFeature = newContainer.EClass().GetEStructuralFeature(EOPPOSITE_FEATURE_BASE - newContainerFeatureID)
		} else {
			containerFeature = o.AsEObject().EClass().GetEStructuralFeature(newContainerFeatureID).(EReference).GetEOpposite()
		}
		o.store.SetContainer(o.AsEObject(), newContainer, containerFeature)
	}
}

func (o *EStoreEObjectImpl) initializeContainer() {
	if o.ReflectiveEObjectImpl.EInternalContainer() == unitializedContainer && o.store != nil {
		container, feature := o.store.GetContainer(o.AsEObject())
		if container != nil && feature != nil {
			featureID := EOPPOSITE_FEATURE_BASE - container.EClass().GetFeatureID(feature)
			if reference, _ := feature.(EReference); reference != nil {
				if opposite := reference.GetEOpposite(); opposite != nil {
					featureID = o.AsEObject().EClass().GetFeatureID(opposite)
				}
			}
			o.ReflectiveEObjectImpl.ESetInternalContainer(container, featureID)
		}
	}
}

func (o *EStoreEObjectImpl) EDynamicIsSet(dynamicFeatureID int) bool {
	if o.properties != nil && o.properties[dynamicFeatureID] != nil {
		return true
	}
	if store := o.AsEStoreEObject().GetEStore(); store != nil {
		eFeature := o.eDynamicFeature(dynamicFeatureID)
		return store.IsSet(o.AsEObject(), eFeature)
	}
	return false
}

func (o *EStoreEObjectImpl) EDynamicGet(dynamicFeatureID int) any {
	// retrieve value
	var result any
	if o.properties != nil {
		result = o.properties[dynamicFeatureID]
	}
	// compute value
	if result == nil {
		var properties []any
		if eFeature := o.eDynamicFeature(dynamicFeatureID); !eFeature.IsTransient() {
			if eFeature.IsMany() {
				if IsMapType(eFeature) {
					result = o.createMap(eFeature)
				} else {
					result = o.createList(eFeature)
				}
				// feature is a container : create corresponding container type and cache it
				// in properties to avoid multiple allocation of the same container
				properties = o.getProperties()
			} else if store := o.AsEStoreEObject().GetEStore(); store != nil {
				// feature is not transient and we have a store
				result = store.Get(o.AsEObject(), eFeature, NO_INDEX)
				if o.cache {
					properties = o.getProperties()
				}

			}
		}
		// store value in properties
		if properties != nil {
			properties[dynamicFeatureID] = result
		}
	}
	return result
}

func (o *EStoreEObjectImpl) EDynamicSet(dynamicFeatureID int, value any) {
	var properties []any
	eFeature := o.eDynamicFeature(dynamicFeatureID)
	store := o.AsEStoreEObject().GetEStore()
	if store != nil && !eFeature.IsTransient() {
		// store and feature is not transient
		if o.cache {
			properties = o.getProperties()
		}
		// we can execute an async operation if we have a cache and an async store
		if asyncStore, _ := store.(EStoreAsync); asyncStore != nil && o.cache {
			asyncStore.AsyncOperation(o, func() {
				store.Set(o.AsEObject(), eFeature, NO_INDEX, value, false)
			})
		} else {
			store.Set(o.AsEObject(), eFeature, NO_INDEX, value, false)
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
}

func (o *EStoreEObjectImpl) EDynamicUnset(dynamicFeatureID int) {
	if o.properties != nil {
		o.properties[dynamicFeatureID] = nil
	}
	eFeature := o.eDynamicFeature(dynamicFeatureID)
	store := o.AsEStoreEObject().GetEStore()
	if store != nil && !eFeature.IsTransient() {
		if asyncStore, _ := store.(EStoreAsync); asyncStore != nil && o.cache {
			asyncStore.AsyncOperation(o, func() {
				store.UnSet(o.AsEObject(), eFeature)
			})
		} else {
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
	return l
}

func (o *EStoreEObjectImpl) createMap(eFeature EStructuralFeature) EMap {
	eClass := eFeature.GetEType().(EClass)
	m := NewEStoreMap(eClass, o.AsEObject(), eFeature, o.AsEStoreEObject().GetEStore())
	m.SetCache(o.cache)
	return m
}
