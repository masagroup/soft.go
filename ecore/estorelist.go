// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import "strconv"

type EStoreList struct {
	interfaces  any
	owner       EObject
	feature     EStructuralFeature
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
	list.interfaces = list
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

func (list *EStoreList) asListCallbacks() listCallBacks {
	return list.interfaces.(listCallBacks)
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

type basicEStoreListNotification struct {
	AbstractNotification
	list *EStoreList
}

func (notif *basicEStoreListNotification) GetNotifier() ENotifier {
	return notif.list.interfaces.(ENotifyingList).GetNotifier()
}

func (notif *basicEStoreListNotification) GetFeature() EStructuralFeature {
	return notif.list.interfaces.(ENotifyingList).GetFeature()
}

func (notif *basicEStoreListNotification) GetFeatureID() int {
	return notif.list.interfaces.(ENotifyingList).GetFeatureID()
}

func (list *EStoreList) createNotification(eventType EventType, oldValue any, newValue any, position int) ENotification {
	n := new(basicEStoreListNotification)
	n.Initialize(n, eventType, oldValue, newValue, position)
	n.list = list
	return n
}

func (list *EStoreList) isNotificationRequired() bool {
	notifier := list.interfaces.(ENotifyingList).GetNotifier()
	return notifier != nil && notifier.EDeliver() && !notifier.EAdapters().Empty()
}

func (list *EStoreList) createAndAddNotification(ns ENotificationChain, eventType EventType, oldValue any, newValue any, position int) ENotificationChain {
	notifications := ns
	if list.isNotificationRequired() {
		notification := list.createNotification(eventType, oldValue, newValue, position)
		if notifications != nil {
			notifications.Add(notification)
		} else {
			notifications = notification.(ENotificationChain)
		}
	}
	return notifications
}

func (list *EStoreList) createAndDispatchNotification(notifications ENotificationChain, eventType EventType, oldValue any, newValue any, position int) {
	list.createAndDispatchNotificationFn(notifications, func() ENotification {
		return list.createNotification(eventType, oldValue, newValue, position)
	})
}

func (list *EStoreList) createAndDispatchNotificationFn(notifications ENotificationChain, createNotification func() ENotification) {
	if list.isNotificationRequired() {
		notification := createNotification()
		if notifications != nil {
			notifications.Add(notification)
			notifications.Dispatch()
		} else {
			notifier := list.interfaces.(ENotifyingList).GetNotifier()
			if notifier != nil {
				notifier.ENotify(notification)
			}
		}
	} else {
		if notifications != nil {
			notifications.Dispatch()
		}
	}
}

func (list *EStoreList) Add(e any) bool {
	return list.Insert(list.Size(), e)
}

func (list *EStoreList) AddWithNotification(object any, notifications ENotificationChain) ENotificationChain {
	// add to store
	index := list.Size()
	list.store.Add(list.owner, list.feature, index, object)

	// callbacks
	listCallbacks := list.asListCallbacks()
	listCallbacks.DidAdd(index, object)
	listCallbacks.DidChange()

	// notifications
	return list.createAndAddNotification(notifications, ADD, nil, object, index)
}

func (list *EStoreList) AddAll(c EList) bool {
	return list.InsertAll(list.Size(), c)
}

func (list *EStoreList) Insert(index int, e any) bool {
	if index < 0 || index > list.Size() {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(list.Size()))
	}
	if list.Contains(e) {
		return false
	}

	// store operation
	list.store.Add(list.owner, list.feature, index, e)

	// callbacks
	listCallbacks := list.asListCallbacks()
	listCallbacks.DidAdd(index, e)
	listCallbacks.DidChange()

	// notifications
	notifications := list.interfaces.(eNotifyingListInternal).inverseAdd(e, nil)
	list.createAndDispatchNotification(notifications, ADD, nil, e, index)
	return true
}

func (list *EStoreList) InsertAll(index int, collection EList) bool {
	if index < 0 || index > list.Size() {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(list.Size()))
	}
	collection = getNonDuplicates(collection, list)
	if collection.Size() == 0 {
		return false
	}

	// add to the store && inverseAdd
	listCallbacks := list.asListCallbacks()
	i := index
	var notifications ENotificationChain = NewNotificationChain()
	var notifyingList eNotifyingListInternal = list.interfaces.(eNotifyingListInternal)
	for it := collection.Iterator(); it.HasNext(); i++ {
		element := it.Next()

		// store operation
		list.store.Add(list.owner, list.feature, i, element)

		// callback
		listCallbacks.DidAdd(i, element)

		// notifications
		notifications = notifyingList.inverseAdd(element, notifications)
	}

	// callbacks
	listCallbacks.DidChange()

	// notifications
	list.createAndDispatchNotificationFn(notifications, func() ENotification {
		if collection.Size() == 1 {
			return list.createNotification(ADD, nil, collection.Get(0), index)
		} else {
			return list.createNotification(ADD_MANY, nil, collection.ToArray(), index)
		}
	})
	return true
}

func (list *EStoreList) MoveObject(newIndex int, elem any) {
	oldIndex := list.IndexOf(elem)
	if oldIndex == -1 {
		panic("Object not found")
	}
	list.Move(oldIndex, newIndex)
}

func (list *EStoreList) Move(oldIndex int, newIndex int) any {
	if oldIndex < 0 || oldIndex >= list.Size() ||
		newIndex < 0 || newIndex > list.Size() {
		panic("Index out of bounds: oldIndex=" + strconv.Itoa(oldIndex) + " newIndex=" + strconv.Itoa(newIndex) + " size=" + strconv.Itoa(list.Size()))
	}

	// store
	oldObject := list.store.Move(list.owner, list.feature, oldIndex, newIndex)

	// callbacks
	listCallbacks := list.asListCallbacks()
	listCallbacks.DidMove(newIndex, oldObject, oldIndex)
	listCallbacks.DidChange()

	// notifications
	list.createAndDispatchNotification(nil, MOVE, oldIndex, oldObject, newIndex)
	return oldObject
}

func (list *EStoreList) Get(index int) any {
	return list.resolve(index, list.store.Get(list.owner, list.feature, index))
}

func (list *EStoreList) resolve(index int, object any) any {
	if list.object && list.proxies {
		resolved := list.resolveProxy(object.(EObject))
		if resolved != object {
			list.store.Set(list.owner, list.feature, index, resolved)
			var notifications ENotificationChain
			if list.containment {
				notifications = list.interfaces.(eNotifyingListInternal).inverseRemove(object, notifications)
				if resolvedInternal, _ := resolved.(EObjectInternal); resolvedInternal != nil && resolvedInternal.EInternalContainer() == nil {
					notifications = list.interfaces.(eNotifyingListInternal).inverseAdd(resolved, notifications)
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

func (list *EStoreList) Set(index int, newObject any) any {
	if index < 0 || index >= list.Size() {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(list.Size()))
	}
	currIndex := list.IndexOf(newObject)
	if currIndex >= 0 && currIndex != index {
		panic("element already in list")
	}

	// store
	oldObject := list.store.Set(list.owner, list.feature, index, newObject)

	// callbacks
	listCallbacks := list.asListCallbacks()
	listCallbacks.DidSet(index, newObject, oldObject)
	listCallbacks.DidChange()

	// notifications
	if newObject != oldObject {
		var notifications ENotificationChain
		var notifyingList eNotifyingListInternal = list.interfaces.(eNotifyingListInternal)
		notifications = notifyingList.inverseRemove(oldObject, notifications)
		notifications = notifyingList.inverseAdd(newObject, notifications)
		list.createAndDispatchNotification(notifications, SET, oldObject, newObject, index)
	}
	return oldObject
}

// SetWithNotification ...
func (list *EStoreList) SetWithNotification(index int, newObject any, notifications ENotificationChain) ENotificationChain {
	// store
	oldObject := list.store.Set(list.owner, list.feature, index, newObject)

	// callbacks
	listCallbacks := list.asListCallbacks()
	listCallbacks.DidSet(index, newObject, oldObject)
	listCallbacks.DidChange()

	// notifications
	return list.createAndAddNotification(notifications, SET, oldObject, newObject, index)
}

func (list *EStoreList) RemoveAt(index int) any {
	if index < 0 || index >= list.Size() {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(list.Size()))
	}

	// store
	oldObject := list.store.Remove(list.owner, list.feature, index)

	// callbacks
	listCallbacks := list.asListCallbacks()
	listCallbacks.DidRemove(index, oldObject)
	listCallbacks.DidChange()

	// notifications
	var notifications ENotificationChain
	var notifyingList eNotifyingListInternal = list.interfaces.(eNotifyingListInternal)
	notifications = notifyingList.inverseRemove(oldObject, notifications)
	list.createAndDispatchNotification(notifications, REMOVE, oldObject, nil, index)
	return oldObject
}

func (list *EStoreList) Remove(element any) bool {
	index := list.IndexOf(element)
	if index == -1 {
		return false
	}
	list.RemoveAt(index)
	return true
}

// RemoveWithNotification ...
func (list *EStoreList) RemoveWithNotification(object any, notifications ENotificationChain) ENotificationChain {
	index := list.IndexOf(object)
	if index != -1 {
		// store
		oldObject := list.store.Remove(list.owner, list.feature, index)

		// callbacks
		listCallbacks := list.asListCallbacks()
		listCallbacks.DidRemove(index, oldObject)
		listCallbacks.DidChange()

		// notifications
		return list.createAndAddNotification(notifications, REMOVE, oldObject, nil, index)
	}
	return notifications
}

func (list *EStoreList) RemoveAll(collection EList) bool {
	modified := false
	for i := list.Size() - 1; i >= 0; i-- {
		element := list.store.Get(list.owner, list.feature, i)
		if collection.Contains(element) {
			list.RemoveAt(i)
			modified = true
		}
	}
	return modified
}

func (list *EStoreList) RemoveRange(fromIndex, toIndex int) {
	var objects []any
	var positions []int
	var notifications ENotificationChain = NewNotificationChain()
	listCallbacks := list.asListCallbacks()
	for i := fromIndex; i < toIndex; i++ {
		// store
		oldObject := list.store.Remove(list.owner, list.feature, i)

		// callbacks
		listCallbacks.DidRemove(i, oldObject)

		// notificatins
		notifications = list.interfaces.(eNotifyingListInternal).inverseRemove(oldObject, notifications)
		objects = append(objects, oldObject)
		positions = append(positions, i)
	}

	// callbacks
	listCallbacks.DidChange()

	// notifications
	if len(objects) > 0 {
		list.createAndDispatchNotificationFn(notifications,
			func() ENotification {
				if len(objects) == 1 {
					return list.createNotification(REMOVE, objects[0], nil, fromIndex)
				} else {
					return list.createNotification(REMOVE_MANY, objects, positions, fromIndex)
				}
			})
	}
}

func (list *EStoreList) Size() int {
	return list.store.Size(list.owner, list.feature)
}

func (list *EStoreList) Clear() {
	oldData := list.store.ToArray(list.owner, list.feature)

	// store
	list.store.Clear(list.owner, list.feature)

	// callbacks
	listCallbacks := list.asListCallbacks()
	listCallbacks.DidClear(oldData)
	listCallbacks.DidChange()

	// notifications
	if len(oldData) == 0 {
		list.createAndDispatchNotification(nil, REMOVE_MANY, []any{}, nil, -1)
	} else {
		var notifications ENotificationChain = NewNotificationChain()
		for _, e := range oldData {
			notifications = list.interfaces.(eNotifyingListInternal).inverseRemove(e, notifications)
		}

		list.createAndDispatchNotificationFn(notifications,
			func() ENotification {
				if len(oldData) == 1 {
					return list.createNotification(REMOVE, oldData[0], nil, 0)
				} else {
					return list.createNotification(REMOVE_MANY, oldData, nil, -1)
				}
			})
	}
}

func (list *EStoreList) Empty() bool {
	return list.store.IsEmpty(list.owner, list.feature)
}

func (list *EStoreList) Contains(element any) bool {
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

func (list *EStoreList) IndexOf(element any) int {
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
	return -1
}

func (list *EStoreList) Iterator() EIterator {
	return &listIterator{list: list}
}

func (list *EStoreList) ToArray() []any {
	return list.store.ToArray(list.owner, list.feature)
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
	l := NewEStoreList(list.owner, list.feature, list.store)
	l.proxies = false
	return l
}

func (list *EStoreList) DidAdd(index int, elem any) {

}

func (list *EStoreList) DidSet(index int, newElem any, oldElem any) {

}

func (list *EStoreList) DidRemove(index int, old any) {

}

func (list *EStoreList) DidClear(oldObjects []any) {

}

func (list *EStoreList) DidMove(newIndex int, movedObject any, oldIndex int) {

}

func (list *EStoreList) DidChange() {

}
