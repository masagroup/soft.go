// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

type EStoreList struct {
	BasicENotifyingList
	owner       EObject
	feature     EStructuralFeature
	store       EStore
	cache       bool
	object      bool
	containment bool
	inverse     bool
	opposite    bool
	proxies     bool
	unset       bool
}

func NewEStoreList(owner EObject, feature EStructuralFeature, store EStore) *EStoreList {
	list := &EStoreList{}
	list.Initialize(owner, feature, store)
	list.SetInterfaces(list)
	return list
}

func (list *EStoreList) Initialize(owner EObject, feature EStructuralFeature, store EStore) {
	list.owner = owner
	list.feature = feature
	list.store = store
	list.object = false
	list.containment = false
	list.inverse = false
	list.opposite = false
	list.proxies = false
	list.unset = false
	if ref, _ := feature.(EReference); ref != nil {
		list.object = true
		list.containment = ref.IsContainment()
		list.proxies = ref.IsResolveProxies()
		list.unset = ref.IsUnsettable()
		reverseFeature := ref.GetEOpposite()
		if list.containment {
			if reverseFeature != nil {
				list.inverse = true
				list.opposite = true
			} else {
				list.inverse = true
				list.opposite = false
			}
		} else {
			if reverseFeature != nil {
				list.inverse = true
				list.opposite = true
			} else {
				list.inverse = false
				list.opposite = false
			}
		}
	}
}

func (list *EStoreList) GetOwner() EObject {
	return list.owner
}

func (list *EStoreList) GetNotifier() ENotifier {
	return list.owner
}

func (list *EStoreList) GetFeature() EStructuralFeature {
	return list.feature
}

func (list *EStoreList) GetFeatureID() int {
	return list.owner.EClass().GetFeatureID(list.feature)
}

func (list *EStoreList) GetEStore() EStore {
	return list.store
}

func (list *EStoreList) SetEStore(store EStore) {
	list.store = store
}

// Set object with a cache for its feature values
func (list *EStoreList) SetCache(cache bool) {
	list.cache = cache
}

// Returns true if object is caching feature values
func (list *EStoreList) IsCache() bool {
	return list.cache
}

// Clear object feature values cache
func (list *EStoreList) ClearCache() {
	list.data = nil
}

func (list *EStoreList) performAdd(object any) {
	// add to cache
	if list.data != nil {
		list.BasicENotifyingList.performAdd(object)
	}

	// add to store
	index := list.Size()
	if list.store != nil {
		list.store.Add(list.owner, list.feature, index, object)
	}
}

func (list *EStoreList) performAddAll(c EList) {
	// add to cache
	if list.data != nil {
		list.BasicENotifyingList.performAddAll(c)
	}

	// add to store
	if list.store != nil {
		index := list.Size()
		for it := c.Iterator(); it.HasNext(); index++ {
			list.store.Add(list.owner, list.feature, index, it.Next())
		}
	}
}

func (list *EStoreList) performInsert(index int, object any) {
	// add to cache
	if list.data != nil {
		list.BasicENotifyingList.performInsert(index, object)
	}

	// add to store
	if list.store != nil {
		list.store.Add(list.owner, list.feature, index, object)
	}
}

func (list *EStoreList) performInsertAll(index int, c EList) bool {
	// add to cache
	if list.data != nil {
		if !list.BasicENotifyingList.performInsertAll(index, c) {
			return false
		}
	}

	// add to store
	if list.store != nil {
		i := index
		for it := c.Iterator(); it.HasNext(); i++ {
			list.store.Add(list.owner, list.feature, i, it.Next())
		}
	}

	return true
}

func (list *EStoreList) performClear() []any {
	var result []any
	//cache
	if list.data != nil {
		result = list.BasicENotifyingList.performClear()
	}
	// store
	if list.store != nil {
		if result == nil {
			result = list.store.ToArray(list.owner, list.feature)
		}
		list.store.Clear(list.owner, list.feature)
	}
	return result
}

func (list *EStoreList) performRemove(index int) any {
	var result any
	//cache
	if list.data != nil {
		result = list.BasicENotifyingList.performRemove(index)
	}
	//store
	if list.store != nil {
		result = list.store.Remove(list.owner, list.feature, index)
	}
	return result
}

func (list *EStoreList) performRemoveRange(fromIndex int, toIndex int) []any {
	var result []any
	if list.data != nil {
		result = list.BasicENotifyingList.performRemoveRange(fromIndex, toIndex)
	}
	if list.store != nil {
		var objects []any
		for i := fromIndex; i < toIndex; i++ {
			// store
			object := list.store.Remove(list.owner, list.feature, i)
			objects = append(objects, object)
		}
		result = objects
	}
	return result
}

func (list *EStoreList) performSet(index int, object any) any {
	var result any
	if list.data != nil {
		result = list.BasicENotifyingList.performSet(index, object)
	}
	if list.store != nil {
		result = list.store.Set(list.owner, list.feature, index, object)
	}
	return result
}

func (list *EStoreList) performMove(oldIndex, newIndex int) any {
	var result any
	if list.data != nil {
		result = list.BasicENotifyingList.performMove(oldIndex, newIndex)
	}
	if list.store != nil {
		result = list.store.Move(list.owner, list.feature, oldIndex, newIndex)
	}
	return result
}

func (list *EStoreList) doGet(index int) any {
	return list.resolve(index, list.get(index))
}

func (list *EStoreList) get(index int) any {
	if list.data != nil {
		return list.data[index]
	}
	if list.store != nil {
		return list.store.Get(list.owner, list.feature, index)
	}
	return nil
}

func (list *EStoreList) resolve(index int, object any) any {
	if list.object && list.proxies {
		resolved := list.resolveProxy(object.(EObject))
		if resolved != object {
			// only update cache because only proxies are stored in a store
			if list.data != nil {
				list.data[index] = resolved
			}
			var notifications ENotificationChain
			if list.containment {
				notifications = list.interfaces.(abstractENotifyingList).inverseRemove(object, notifications)
				if resolvedInternal, _ := resolved.(EObjectInternal); resolvedInternal != nil && resolvedInternal.EInternalContainer() == nil {
					notifications = list.interfaces.(abstractENotifyingList).inverseAdd(resolved, notifications)
				}
			}
			list.createAndDispatchNotification(notifications, RESOLVE, object, resolved, index)
		}
		return resolved
	}
	return object
}

func (list *EStoreList) resolveProxy(eObject EObject) EObject {
	if list.proxies && eObject.EIsProxy() {
		return list.owner.(EObjectInternal).EResolveProxy(eObject)
	}
	return eObject
}

func (list *EStoreList) ToArray() []any {
	if list.data != nil {
		if list.proxies {
			for i := len(list.data) - 1; i >= 0; i-- {
				list.doGet(i)
			}
		}
		return list.data
	}
	if list.store != nil {
		data := list.store.ToArray(list.owner, list.feature)
		if list.proxies {
			for i := len(list.data) - 1; i >= 0; i-- {
				data[i] = list.resolve(i, data[i])
			}
		}
		return data
	}
	return nil
}

func (list *EStoreList) Size() int {
	if list.data != nil {
		return len(list.data)
	}
	if list.store != nil {
		return list.store.Size(list.owner, list.feature)
	}
	return 0
}

func (list *EStoreList) Empty() bool {
	if list.data != nil {
		return len(list.data) == 0
	}
	if list.store != nil {
		return list.store.IsEmpty(list.owner, list.feature)
	}
	return true
}

func (list *EStoreList) Contains(element any) bool {
	if list.data != nil {
		for i, value := range list.data {
			if value == element || (list.resolve(i, value) == element) {
				return true
			}
		}
	}
	if list.store != nil {
		result := list.store.Contains(list.owner, list.feature, element)
		if !result && list.object && list.proxies {
			for i := 0; i < list.Size(); i++ {
				eObject, _ := list.store.Get(list.owner, list.feature, i).(EObject)
				eResolved := list.resolveProxy(eObject)
				if element == eResolved {
					return true
				}
			}
		}
		return result
	}
	return false
}

func (list *EStoreList) IndexOf(element any) int {
	if list.data != nil {
		for i, value := range list.data {
			if value == element || (list.resolve(i, value) == element) {
				return i
			}
		}
	}
	if list.store != nil {
		result := list.store.IndexOf(list.owner, list.feature, element)
		if result >= 0 {
			return result
		}
		if list.object && list.proxies {
			for i := 0; i < list.Size(); i++ {
				eObject, _ := list.store.Get(list.owner, list.feature, i).(EObject)
				eResolved := list.resolveProxy(eObject)
				if element == eResolved {
					return i
				}
			}
		}
	}
	return -1
}

func (list *EStoreList) inverseAdd(object any, notifications ENotificationChain) ENotificationChain {
	internal, _ := object.(EObjectInternal)
	if internal != nil && list.inverse {
		if list.opposite {
			inverseReference := list.feature.(EReference).GetEOpposite()
			inverseFeatureID := internal.EClass().GetFeatureID(inverseReference)
			return internal.EInverseAdd(list.owner, inverseFeatureID, notifications)
		} else {
			featureID := list.feature.GetFeatureID()
			return internal.EInverseAdd(list.owner, EOPPOSITE_FEATURE_BASE-featureID, notifications)
		}
	}
	return notifications
}

func (list *EStoreList) inverseRemove(object any, notifications ENotificationChain) ENotificationChain {
	internal, _ := object.(EObjectInternal)
	if internal != nil && list.inverse {
		if list.opposite {
			inverseReference := list.feature.(EReference).GetEOpposite()
			inverseFeatureID := internal.EClass().GetFeatureID(inverseReference)
			return internal.EInverseRemove(list.owner, inverseFeatureID, notifications)
		} else {
			featureID := list.feature.GetFeatureID()
			return internal.EInverseRemove(list.owner, EOPPOSITE_FEATURE_BASE-featureID, notifications)
		}
	}
	return notifications
}

func (list *EStoreList) GetUnResolvedList() EList {
	if list.object && list.proxies {
		l := &unResolvedEStoreList{}
		l.delegate = list
		return l
	}
	return list
}

type unResolvedEStoreList struct {
	AbstractDelegatingENotifyingList[*EStoreList]
}

func (list *unResolvedEStoreList) doGet(index int) any {
	return list.delegate.get(index)
}

func (list *unResolvedEStoreList) ToArray() []any {
	return list.delegate.data
}

func (l *unResolvedEStoreList) GetUnResolvedList() EList {
	return l
}
