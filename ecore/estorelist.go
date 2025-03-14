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
	"iter"
	"slices"

	"github.com/SokaDance/rmx"
)

type EStoreList struct {
	BasicENotifyingList
	mutex       rmx.RecursiveMutex
	owner       EObject
	feature     EStructuralFeature
	cache       bool
	size        int
	store       EStore
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
	list.isUnique = feature.IsUnique()
	list.owner = owner
	list.feature = feature
	list.store = store
	list.cache = false
	if list.store != nil {
		// one store <=> no data
		list.data = nil
		list.size = list.store.Size(owner, feature)
	} else {
		// no store <=> we need data
		list.data = []any{}
		list.size = 0
	}
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

func (list *EStoreList) Lock() {
	list.mutex.Lock()
}

func (list *EStoreList) Unlock() {
	list.mutex.Unlock()
}

func (list *EStoreList) GetEStore() EStore {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	return list.store
}

func (list *EStoreList) SetEStore(newStore EStore) {
	list.mutex.Lock()
	if oldStore := list.store; oldStore != newStore {
		list.store = newStore
		if newStore == nil {
			// unbind previous store
			if !list.cache {
				// got to backup store if data is not existing
				list.data = oldStore.ToArray(list.owner, list.feature)
			}
		} else {
			if !list.cache {
				list.data = nil
			}
		}
	}
	list.mutex.Unlock()
}

// Set object with a cache for its feature values
func (list *EStoreList) SetCache(cache bool) {
	list.mutex.Lock()
	if list.cache != cache {
		list.cache = cache
		var data []any
		if cache {
			// one cache
			if list.store != nil {
				// if there is a store , create data with store values
				data = list.store.ToArray(list.owner, list.feature)
			} else {
				// data is the new cache
				data = list.data
			}
			list.data = data
			// if no data, set list.data as empty data
			if list.data == nil {
				list.data = []any{}
			}
		} else {
			// no cache
			data = list.data
			if list.store != nil {
				// if there is a store, we can remove data
				// otherwise keep it
				list.data = nil
			}
		}

		// set cache for all list elements
		for _, v := range data {
			if sc, _ := v.(ECacheProvider); sc != nil {
				sc.SetCache(cache)
			}
		}
	}
	list.mutex.Unlock()
}

// Returns true if object is caching feature values
func (list *EStoreList) IsCache() bool {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	return list.cache
}

func (list *EStoreList) performAdd(object any) {
	list.mutex.Lock()
	// add to cache
	if list.data != nil {
		list.BasicENotifyingList.performAdd(object)
	}
	// add to store
	if list.store != nil {
		list.store.Add(list.owner, list.feature, list.size, object)
	}
	// size
	list.size++
	list.mutex.Unlock()
}

func (list *EStoreList) performInsert(index int, object any) {
	list.mutex.Lock()
	// add to cache
	if list.data != nil {
		list.BasicENotifyingList.performInsert(index, object)
	}
	// add to store
	if list.store != nil {
		list.store.Add(list.owner, list.feature, index, object)
	}
	// size
	list.size++
	list.mutex.Unlock()
}

func (list *EStoreList) performInsertAll(index int, c Collection) bool {
	list.mutex.Lock()
	defer list.mutex.Unlock()

	// add to cache
	if list.data != nil {
		if !list.BasicENotifyingList.performInsertAll(index, c) {
			return false
		}
	}
	// add to store
	if list.store != nil {
		list.store.AddAll(list.owner, list.feature, index, c)
	}
	// size
	list.size += c.Size()
	return true
}

func (list *EStoreList) performClear() []any {
	list.mutex.Lock()
	var result []any
	//cache
	if list.data != nil {
		result = list.BasicENotifyingList.performClear()
	}
	// store
	if list.store != nil {
		if list.data == nil {
			result = list.store.ToArray(list.owner, list.feature)
		}
		list.store.Clear(list.owner, list.feature)

	}
	// size
	list.size = 0
	list.mutex.Unlock()
	return result
}

func (list *EStoreList) performRemove(index int) any {
	list.mutex.Lock()
	var result any
	//cache
	if list.data != nil {
		result = list.BasicENotifyingList.performRemove(index)
	}
	//store
	if list.store != nil {
		if list.data == nil {
			result = list.store.Remove(list.owner, list.feature, index, true)
		} else {
			_ = list.store.Remove(list.owner, list.feature, index, false)
		}
	}
	// size
	list.size--
	list.mutex.Unlock()
	return result
}

func (list *EStoreList) performRemoveRange(fromIndex int, toIndex int) []any {
	list.mutex.Lock()
	var result []any
	// cache
	if list.data != nil {
		result = list.BasicENotifyingList.performRemoveRange(fromIndex, toIndex)
	}
	// store
	if list.store != nil {
		for i := fromIndex; i < toIndex; i++ {
			if list.data == nil {
				result = append(result, list.store.Remove(list.owner, list.feature, i, true))
			} else {
				_ = list.store.Remove(list.owner, list.feature, i, false)
			}

		}
	}
	// size
	list.size -= len(result)
	list.mutex.Unlock()
	return result
}

func (list *EStoreList) performSet(index int, object any) any {
	list.mutex.Lock()
	// cache
	var result any
	if list.data != nil {
		result = list.BasicENotifyingList.performSet(index, object)
	}
	// store
	if list.store != nil {
		if list.data == nil {
			result = list.store.Set(list.owner, list.feature, index, object, true)
		} else {
			_ = list.store.Set(list.owner, list.feature, index, object, false)
		}
	}
	list.mutex.Unlock()
	return result
}

func (list *EStoreList) performMove(oldIndex, newIndex int) any {
	list.mutex.Lock()
	var result any
	if list.data != nil {
		result = list.BasicENotifyingList.performMove(oldIndex, newIndex)
	}
	if list.store != nil {
		if list.data == nil {
			result = list.store.Move(list.owner, list.feature, oldIndex, newIndex, true)
		} else {
			_ = list.store.Move(list.owner, list.feature, oldIndex, newIndex, false)
		}
	}
	list.mutex.Unlock()
	return result
}

func (list *EStoreList) doGet(index int) any {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	return list.resolve(index, list.get(index))
}

func (list *EStoreList) get(index int) any {
	if list.data != nil {
		return list.data[index]
	} else if list.store != nil {
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

func (list *EStoreList) All() iter.Seq[any] {
	return func(yield func(any) bool) {
		list.mutex.Lock()
		if list.data != nil {
			for i, v := range list.data {
				if d := list.resolve(i, v); !yield(d) {
					break
				}
			}
		} else if list.store != nil {
			i := 0
			for v := range list.store.All(list.owner, list.feature) {
				if d := list.resolve(i, v); !yield(d) {
					break
				}
				i++
			}
		}
		list.mutex.Unlock()
	}
}

func (list *EStoreList) ToArray() []any {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	if list.data != nil {
		if list.proxies {
			for i := len(list.data) - 1; i >= 0; i-- {
				list.doGet(i)
			}
		}
		return list.data
	} else if list.store != nil {
		data := list.store.ToArray(list.owner, list.feature)
		if list.proxies {
			for i := len(data) - 1; i >= 0; i-- {
				data[i] = list.resolve(i, data[i])
			}
		}
		return data
	}
	return nil
}

func (list *EStoreList) Size() int {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	return list.size
}

func (list *EStoreList) Empty() bool {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	return list.size == 0
}

func (list *EStoreList) IndexOf(element any) int {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	if list.data != nil {
		for i, value := range list.data {
			if value == element || (list.resolve(i, value) == element) {
				return i
			}
		}
	} else if list.store != nil {
		result := list.store.IndexOf(list.owner, list.feature, element)
		if result >= 0 {
			return result
		}
		if list.object && list.proxies {
			for i := range list.size {
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

func (list *EStoreList) Contains(element any) bool {
	if list.data != nil {
		return list.BasicENotifyingList.Contains(element)
	} else if list.store != nil {
		list.mutex.Lock()
		defer list.mutex.Unlock()
		if list.store.Contains(list.owner, list.feature, element) {
			return true
		} else if list.object && list.proxies {
			for i := range list.size {
				eObject, _ := list.store.Get(list.owner, list.feature, i).(EObject)
				eResolved := list.resolveProxy(eObject)
				if element == eResolved {
					return true
				}
			}
		}
	}
	return false
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
		return newUnResolvedEStoreList(list)
	}
	return list
}

type unResolvedEStoreList struct {
	AbstractDelegatingENotifyingList[*EStoreList]
}

func newUnResolvedEStoreList(delegate *EStoreList) *unResolvedEStoreList {
	l := &unResolvedEStoreList{}
	l.interfaces = l
	l.delegate = delegate
	l.isUnique = true
	return l
}

func (list *unResolvedEStoreList) doGet(index int) any {
	list.delegate.mutex.Lock()
	defer list.delegate.mutex.Unlock()
	return list.delegate.get(index)
}

func (list *unResolvedEStoreList) IndexOf(elem any) int {
	return list.AbstractEList.IndexOf(elem)
}

func (list *unResolvedEStoreList) All() iter.Seq[any] {
	list.delegate.mutex.Lock()
	defer list.delegate.mutex.Unlock()
	if list.delegate.data != nil {
		return slices.Values(list.delegate.data)
	}
	if list.delegate.store != nil {
		return list.delegate.store.All(list.delegate.owner, list.delegate.feature)
	}
	return nil
}

func (list *unResolvedEStoreList) ToArray() []any {
	list.delegate.mutex.Lock()
	defer list.delegate.mutex.Unlock()
	if list.delegate.data != nil {
		return list.delegate.data
	}
	if list.delegate.store != nil {
		return list.delegate.store.ToArray(list.delegate.owner, list.delegate.feature)
	}
	return nil
}

func (l *unResolvedEStoreList) GetUnResolvedList() EList {
	return l
}
