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
	"sync"
	"sync/atomic"
)

// EStoreList is a list whose elements are kept in a store, in a cache, or both.
//
// Its state (cached elements, size, store, cache flag) is never changed but replaced by the
// writers, which hold mutex: readers read it without lock. An element added at the end is written
// past the length of the published cache, which the readers of an older state don't read: adding
// is amortized O(1). The other changes of the cache copy it.
type EStoreList struct {
	BasicENotifyingList
	state       atomic.Pointer[eStoreListState]
	mutex       sync.Mutex
	owner       EObject
	feature     EStructuralFeature
	object      bool
	containment bool
	inverse     bool
	opposite    bool
	proxies     bool
	unset       bool
}

// eStoreListState is a state of an EStoreList. It is never changed once published.
type eStoreListState struct {
	// data is the cache of the elements, nil when not caching with a store. The elements up to
	// its length are never changed: the next state appends past it, or uses a copy.
	data  []any
	size  int
	store EStore
	cache bool
}

func (s *eStoreListState) clone() *eStoreListState {
	result := *s
	return &result
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
	state := &eStoreListState{store: store}
	if store != nil {
		// one store <=> no data
		state.size = store.Size(owner, feature)
	} else {
		// no store <=> we need data
		state.data = []any{}
	}
	list.state.Store(state)
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

// Lock prevents the writes to the list until Unlock. The list can be read meanwhile, by the
// goroutine holding the lock too.
func (list *EStoreList) Lock() {
	list.mutex.Lock()
}

func (list *EStoreList) Unlock() {
	list.mutex.Unlock()
}

func (list *EStoreList) GetEStore() EStore {
	return list.state.Load().store
}

func (list *EStoreList) SetEStore(newStore EStore) {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	list.setEStoreLocked(newStore)
}

// setEStoreLocked sets the store of the list: mutex must be held.
func (list *EStoreList) setEStoreLocked(newStore EStore) {
	current := list.state.Load()
	if oldStore := current.store; oldStore != newStore {
		state := current.clone()
		state.store = newStore
		if newStore == nil {
			// unbind previous store
			if !state.cache {
				// got to backup store if data is not existing
				state.data = slices.Clip(oldStore.ToArray(list.owner, list.feature))
			}
		} else {
			if !state.cache {
				state.data = nil
			}
		}
		list.state.Store(state)
	}
}

// Set object with a cache for its feature values
func (list *EStoreList) SetCache(cache bool) {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	current := list.state.Load()
	if current.cache != cache {
		state := current.clone()
		state.cache = cache
		var data []any
		if cache {
			// one cache
			if state.store != nil {
				// if there is a store , create data with store values
				data = slices.Clip(state.store.ToArray(list.owner, list.feature))
			} else {
				// data is the new cache
				data = state.data
			}
			state.data = data
			// if no data, set list.data as empty data
			if state.data == nil {
				state.data = []any{}
			}
		} else {
			// no cache
			data = state.data
			if state.store != nil {
				// if there is a store, we can remove data
				// otherwise keep it
				state.data = nil
			}
		}
		list.state.Store(state)

		// set cache for all list elements
		for _, v := range data {
			if sc, _ := v.(ECacheProvider); sc != nil {
				sc.SetCache(cache)
			}
		}
	}
}

// Returns true if object is caching feature values
func (list *EStoreList) IsCache() bool {
	return list.state.Load().cache
}

// insertData returns a copy of data with elements inserted at index. At the end, it appends to
// data: the elements past its length are not read.
func insertData(data []any, index int, elements ...any) []any {
	if index == len(data) {
		return append(data, elements...)
	}
	result := make([]any, 0, len(data)+len(elements))
	result = append(result, data[:index]...)
	result = append(result, elements...)
	return append(result, data[index:]...)
}

// removeData returns a copy of data without the elements from fromIndex to toIndex.
func removeData(data []any, fromIndex, toIndex int) []any {
	result := make([]any, 0, len(data)-(toIndex-fromIndex))
	result = append(result, data[:fromIndex]...)
	return append(result, data[toIndex:]...)
}

func (list *EStoreList) PerformAdd(object any) {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	current := list.state.Load()
	state := current.clone()
	index := state.size
	// add to data
	if state.data != nil {
		state.data = insertData(state.data, index, object)
	}
	// add to store
	if state.store != nil {
		state.store.Add(list.owner, list.feature, index, object)
	}
	// update size
	state.size++
	list.state.Store(state)
	// events
	listCallbacks := list.asEListCallbacks()
	listCallbacks.DidAdd(index, object)
	listCallbacks.DidChange()
}

func (list *EStoreList) PerformInsert(index int, object any) {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	state := list.state.Load().clone()
	// add to cache
	if state.data != nil {
		state.data = insertData(state.data, index, object)
	}
	// add to store
	if state.store != nil {
		state.store.Add(list.owner, list.feature, index, object)
	}
	// size
	state.size++
	list.state.Store(state)
	// events
	listCallbacks := list.asEListCallbacks()
	listCallbacks.DidAdd(index, object)
	listCallbacks.DidChange()
}

func (list *EStoreList) PerformInsertAll(index int, c Collection) bool {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	state := list.state.Load().clone()
	elements := c.ToArray()
	// add to cache
	if state.data != nil {
		if len(elements) == 0 {
			return false
		}
		state.data = insertData(state.data, index, elements...)
	}
	// add to store
	if state.store != nil {
		state.store.AddAll(list.owner, list.feature, index, c)
	}
	// size
	state.size += len(elements)
	list.state.Store(state)
	// events
	listCallbacks := list.asEListCallbacks()
	for i, object := range elements {
		listCallbacks.DidAdd(index+i, object)
		listCallbacks.DidChange()
	}
	return true
}

func (list *EStoreList) PerformClear() []any {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	state := list.state.Load().clone()
	var result []any

	//cache
	if state.data != nil {
		result = state.data
		state.data = []any{}
	}

	// store
	if state.store != nil {
		if result == nil {
			// clear with events
			result = state.store.ToArray(list.owner, list.feature)
		}
		state.store.Clear(list.owner, list.feature)
	}
	// size
	state.size = 0
	list.state.Store(state)
	list.asEListCallbacks().DidClear(result)
	return result
}

func (list *EStoreList) PerformRemove(index int) any {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	state := list.state.Load().clone()
	var result any

	//cache
	cached := state.data != nil
	if cached {
		result = state.data[index]
		state.data = removeData(state.data, index, index+1)
	}
	//store
	if state.store != nil {
		if !cached {
			// remove with events
			result = state.store.Remove(list.owner, list.feature, index, true)
		} else {
			// remove
			_ = state.store.Remove(list.owner, list.feature, index, false)
		}
	}
	// size
	state.size--
	list.state.Store(state)
	listCallbacks := list.asEListCallbacks()
	listCallbacks.DidRemove(index, result)
	listCallbacks.DidChange()
	return result
}

func (list *EStoreList) PerformRemoveRange(fromIndex int, toIndex int) []any {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	state := list.state.Load().clone()
	var result []any

	// cache
	if data := state.data; data != nil {
		result = slices.Clone(data[fromIndex:toIndex])
		state.data = removeData(data, fromIndex, toIndex)
		// store
		if state.store != nil {
			for i := fromIndex; i < toIndex; i++ {
				_ = state.store.Remove(list.owner, list.feature, fromIndex, false)
			}
		}
		state.size -= len(result)
		list.state.Store(state)
		listCallbacks := list.asEListCallbacks()
		for i := toIndex - 1; i >= fromIndex; i-- {
			listCallbacks.DidRemove(i, result[i-fromIndex])
		}
		listCallbacks.DidChange()
		return result
	}
	// store
	if state.store != nil {
		listCallbacks := list.asEListCallbacks()
		for i := fromIndex; i < toIndex; i++ {
			object := state.store.Remove(list.owner, list.feature, fromIndex, true)
			result = append(result, object)
			listCallbacks.DidRemove(fromIndex, result)
		}
		state.size -= len(result)
		list.state.Store(state)
		listCallbacks.DidChange()
	}
	return result
}

func (list *EStoreList) PerformSet(index int, object any) any {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	state := list.state.Load().clone()
	// cache
	var result any
	if state.data != nil {
		result = state.data[index]
		state.data = slices.Clone(state.data)
		state.data[index] = object
		list.state.Store(state)
		if state.store != nil {
			_ = state.store.Set(list.owner, list.feature, index, object, false)
		}
	} else if state.store != nil {
		// set with events
		result = state.store.Set(list.owner, list.feature, index, object, true)
	}
	listCallbacks := list.asEListCallbacks()
	listCallbacks.DidSet(index, object, result)
	listCallbacks.DidChange()
	return result
}

func (list *EStoreList) PerformMove(oldIndex, newIndex int) any {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	state := list.state.Load().clone()
	var result any
	if data := state.data; data != nil {
		result = data[oldIndex]
		if oldIndex != newIndex {
			data = slices.Clone(data)
			if newIndex < oldIndex {
				copy(data[newIndex+1:], data[newIndex:oldIndex])
			} else {
				copy(data[oldIndex:], data[oldIndex+1:newIndex+1])
			}
			data[newIndex] = result
			state.data = data
			list.state.Store(state)
			if state.store != nil {
				_ = state.store.Move(list.owner, list.feature, oldIndex, newIndex, false)
			}
			listCallbacks := list.asEListCallbacks()
			listCallbacks.DidMove(newIndex, result, oldIndex)
			listCallbacks.DidChange()
		} else if state.store != nil {
			_ = state.store.Move(list.owner, list.feature, oldIndex, newIndex, false)
		}
	} else if state.store != nil {
		result = state.store.Move(list.owner, list.feature, oldIndex, newIndex, true)
		listCallbacks := list.asEListCallbacks()
		listCallbacks.DidMove(newIndex, result, oldIndex)
		listCallbacks.DidChange()
	}
	return result
}

func (list *EStoreList) DoGet(index int) any {
	state := list.state.Load()
	return list.resolve(index, state.get(list, index))
}

func (s *eStoreListState) get(list *EStoreList, index int) any {
	if s.data != nil {
		return s.data[index]
	} else if s.store != nil {
		return s.store.Get(list.owner, list.feature, index)
	}
	return nil
}

func (list *EStoreList) resolve(index int, object any) any {
	if list.object && list.proxies {
		resolved := list.resolveProxy(object.(EObject))
		if resolved != object {
			// only update cache because only proxies are stored in a store
			if !list.cacheResolved(index, object, resolved) {
				// resolved by another reader meanwhile, which notified
				return resolved
			}
			var notifications ENotificationChain
			if list.containment {
				notifications = list.interfaces.(internalENotifyingList).InverseRemove(object, notifications)
				if resolvedInternal, _ := resolved.(EObjectInternal); resolvedInternal != nil && resolvedInternal.EInternalContainer() == nil {
					notifications = list.interfaces.(internalENotifyingList).InverseAdd(resolved, notifications)
				}
			}
			list.createAndDispatchNotification(notifications, RESOLVE, object, resolved, index)
		}
		return resolved
	}
	return object
}

// cacheResolved replaces object, at index in the cache, with resolved. It reports false when
// object is not there anymore, replaced by another reader or a writer.
func (list *EStoreList) cacheResolved(index int, object any, resolved any) bool {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	current := list.state.Load()
	if current.data == nil {
		return true
	}
	if index >= len(current.data) || current.data[index] != object {
		return false
	}
	state := current.clone()
	state.data = slices.Clone(state.data)
	state.data[index] = resolved
	list.state.Store(state)
	return true
}

func (list *EStoreList) resolveProxy(eObject EObject) EObject {
	if list.proxies && eObject.EIsProxy() {
		return list.owner.(EObjectInternal).EResolveProxy(eObject)
	}
	return eObject
}

func (list *EStoreList) All() iter.Seq[any] {
	return func(yield func(any) bool) {
		state := list.state.Load()
		if state.data != nil {
			for i, v := range state.data {
				if d := list.resolve(i, v); !yield(d) {
					break
				}
			}
		} else if state.store != nil {
			i := 0
			for v := range state.store.All(list.owner, list.feature) {
				if d := list.resolve(i, v); !yield(d) {
					break
				}
				i++
			}
		}
	}
}

// ToArray returns the elements of the list. The slice of a cache is returned: it must not be
// changed.
func (list *EStoreList) ToArray() []any {
	state := list.state.Load()
	if data := state.data; data != nil {
		if list.proxies {
			for i := range slices.Backward(data) {
				list.resolve(i, data[i])
			}
			// with the resolved elements
			if current := list.state.Load(); current.data != nil {
				data = current.data
			}
		}
		return slices.Clip(data)
	} else if state.store != nil {
		data := state.store.ToArray(list.owner, list.feature)
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
	return list.state.Load().size
}

func (list *EStoreList) Empty() bool {
	return list.state.Load().size == 0
}

func (list *EStoreList) IndexOf(element any) int {
	state := list.state.Load()
	if state.data != nil {
		for i, value := range state.data {
			if value == element || (list.resolve(i, value) == element) {
				return i
			}
		}
	} else if state.store != nil {
		result := state.store.IndexOf(list.owner, list.feature, element)
		if result >= 0 {
			return result
		}
		if list.object && list.proxies {
			for i := range state.size {
				eObject, _ := state.store.Get(list.owner, list.feature, i).(EObject)
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
	state := list.state.Load()
	if state.data != nil {
		return list.IndexOf(element) != -1
	} else if state.store != nil {
		if state.store.Contains(list.owner, list.feature, element) {
			return true
		} else if list.object && list.proxies {
			for i := range state.size {
				eObject, _ := state.store.Get(list.owner, list.feature, i).(EObject)
				eResolved := list.resolveProxy(eObject)
				if element == eResolved {
					return true
				}
			}
		}
	}
	return false
}

func (list *EStoreList) InverseAdd(object any, notifications ENotificationChain) ENotificationChain {
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

func (list *EStoreList) InverseRemove(object any, notifications ENotificationChain) ENotificationChain {
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

func (list *unResolvedEStoreList) DoGet(index int) any {
	return list.delegate.state.Load().get(list.delegate, index)
}

func (list *unResolvedEStoreList) IndexOf(elem any) int {
	return list.AbstractEList.IndexOf(elem)
}

func (list *unResolvedEStoreList) All() iter.Seq[any] {
	state := list.delegate.state.Load()
	if state.data != nil {
		return slices.Values(state.data)
	}
	if state.store != nil {
		return state.store.All(list.delegate.owner, list.delegate.feature)
	}
	return nil
}

// ToArray returns the elements of the list. The slice of a cache is returned: it must not be
// changed.
func (list *unResolvedEStoreList) ToArray() []any {
	state := list.delegate.state.Load()
	if state.data != nil {
		return slices.Clip(state.data)
	}
	if state.store != nil {
		return state.store.ToArray(list.delegate.owner, list.delegate.feature)
	}
	return nil
}

func (l *unResolvedEStoreList) GetUnResolvedList() EList {
	return l
}
