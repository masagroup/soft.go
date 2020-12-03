package ecore

import "strconv"

type BasicEStoreList struct {
	interfaces  interface{}
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

func NewBasicEStoreList(owner EObject, feature EStructuralFeature, store EStore) *BasicEStoreList {
	list := new(BasicEStoreList)
	list.interfaces = list
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
	return list
}

func (list *BasicEStoreList) GetOwner() EObject {
	return list.owner
}

func (list *BasicEStoreList) GetFeature() EStructuralFeature {
	return list.feature
}

func (list *BasicEStoreList) GetFeatureID() int {
	return list.owner.EClass().GetFeatureID(list.feature)
}

type basicEStoreListNotification struct {
	*abstractNotification
	list *BasicEStoreList
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

func (list *BasicEStoreList) createNotification(eventType EventType, oldValue interface{}, newValue interface{}, position int) ENotification {
	n := new(basicEStoreListNotification)
	n.abstractNotification = NewAbstractNotification(eventType, oldValue, newValue, position)
	n.interfaces = n
	n.list = list
	return n
}

func (list *BasicEStoreList) isNotificationRequired() bool {
	notifier := list.interfaces.(ENotifyingList).GetNotifier()
	return notifier != nil && notifier.EDeliver() && !notifier.EAdapters().Empty()
}

func (list *BasicEStoreList) createAndAddNotification(ns ENotificationChain, eventType EventType, oldValue interface{}, newValue interface{}, position int) ENotificationChain {
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

func (list *BasicEStoreList) createAndDispatchNotification(notifications ENotificationChain, eventType EventType, oldValue interface{}, newValue interface{}, position int) {
	list.createAndDispatchNotificationFn(notifications, func() ENotification {
		return list.createNotification(eventType, oldValue, newValue, position)
	})
}

func (list *BasicEStoreList) createAndDispatchNotificationFn(notifications ENotificationChain, createNotification func() ENotification) {
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

func (list *BasicEStoreList) Add(e interface{}) bool {
	return list.Insert(list.Size(), e)
}

func (list *BasicEStoreList) AddWithNotification(object interface{}, notifications ENotificationChain) ENotificationChain {
	index := list.Size()
	list.store.Add(list.owner, list.feature, index, object)
	return list.createAndAddNotification(notifications, ADD, nil, object, index)
}

func (list *BasicEStoreList) AddAll(c EList) bool {
	return list.InsertAll(list.Size(), c)
}

func (list *BasicEStoreList) Insert(index int, e interface{}) bool {
	if index < 0 || index > list.Size() {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(list.Size()))
	}
	if list.Contains(e) {
		return false
	}
	// add to the store && inversAdd
	list.store.Add(list.owner, list.feature, index, e)
	notifications := list.interfaces.(eNotifyingListInternal).inverseAdd(e, nil)
	// notifications
	list.createAndDispatchNotification(notifications, ADD, nil, e, index)
	return true
}

func (list *BasicEStoreList) InsertAll(index int, collection EList) bool {
	if index < 0 || index > list.Size() {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(list.Size()))
	}
	collection = getNonDuplicates(collection, list)
	if collection.Size() == 0 {
		return false
	}
	// add to the store && inverseAdd
	var i int = index
	var notifications ENotificationChain = NewNotificationChain()
	var notifyingList eNotifyingListInternal = list.interfaces.(eNotifyingListInternal)
	for it := collection.Iterator(); it.HasNext(); i++ {
		element := it.Next()
		list.store.Add(list.owner, list.feature, i, element)
		notifications = notifyingList.inverseAdd(element, notifications)
	}
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

func (list *BasicEStoreList) MoveObject(newIndex int, elem interface{}) {
	oldIndex := list.IndexOf(elem)
	if oldIndex == -1 {
		panic("Object not found")
	}
	list.Move(oldIndex, newIndex)
}

func (list *BasicEStoreList) Move(oldIndex int, newIndex int) interface{} {
	if oldIndex < 0 || oldIndex >= list.Size() ||
		newIndex < 0 || newIndex > list.Size() {
		panic("Index out of bounds: oldIndex=" + strconv.Itoa(oldIndex) + " newIndex=" + strconv.Itoa(newIndex) + " size=" + strconv.Itoa(list.Size()))
	}
	oldObject := list.store.Move(list.owner, list.feature, newIndex, oldIndex)
	list.createAndDispatchNotification(nil, MOVE, oldIndex, oldObject, newIndex)
	return oldObject
}

func (list *BasicEStoreList) Get(index int) interface{} {
	return list.resolve(index, list.store.Get(list.owner, list.feature, index))
}

func (list *BasicEStoreList) resolve(index int, object interface{}) interface{} {
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

func (list *BasicEStoreList) resolveProxy(eObject EObject) EObject {
	if list.proxies && eObject.EIsProxy() {
		return list.owner.(EObjectInternal).EResolveProxy(eObject)
	}
	return eObject
}

func (list *BasicEStoreList) Set(index int, newObject interface{}) interface{} {
	if index < 0 || index >= list.Size() {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(list.Size()))
	}
	currIndex := list.IndexOf(newObject)
	if currIndex >= 0 && currIndex != index {
		panic("element already in list")
	}

	oldObject := list.store.Set(list.owner, list.feature, index, newObject)
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
func (list *BasicEStoreList) SetWithNotification(index int, object interface{}, notifications ENotificationChain) ENotificationChain {
	oldObject := list.store.Set(list.owner, list.feature, index, object)
	return list.createAndAddNotification(notifications, SET, oldObject, object, index)
}

func (list *BasicEStoreList) RemoveAt(index int) interface{} {
	if index < 0 || index >= list.Size() {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(list.Size()))
	}
	oldObject := list.store.Remove(list.owner, list.feature, index)
	var notifications ENotificationChain
	var notifyingList eNotifyingListInternal = list.interfaces.(eNotifyingListInternal)
	notifications = notifyingList.inverseRemove(oldObject, notifications)
	list.createAndDispatchNotification(notifications, REMOVE, oldObject, nil, index)
	return oldObject
}

func (list *BasicEStoreList) Remove(element interface{}) bool {
	index := list.IndexOf(element)
	if index == -1 {
		return false
	}
	list.RemoveAt(index)
	return true
}

// RemoveWithNotification ...
func (list *BasicEStoreList) RemoveWithNotification(object interface{}, notifications ENotificationChain) ENotificationChain {
	index := list.IndexOf(object)
	if index != -1 {
		oldObject := list.store.Remove(list.owner, list.feature, index)
		return list.createAndAddNotification(notifications, REMOVE, oldObject, nil, index)
	}
	return notifications
}

func (list *BasicEStoreList) RemoveAll(collection EList) bool {
	modified := false
	for i := list.Size(); i-1 >= 0; i-- {
		element := list.store.Get(list.owner, list.feature, i)
		if collection.Contains(element) {
			list.Remove(i)
			modified = true
		}
	}
	return modified
}

func (list *BasicEStoreList) Size() int {
	return list.store.Size(list.owner, list.feature)
}

func (list *BasicEStoreList) Clear() {
	oldData := list.store.ToArray(list.owner, list.feature)
	list.store.Clear(list.owner, list.feature)
	if len(oldData) == 0 {
		list.createAndDispatchNotification(nil, REMOVE_MANY, []interface{}{}, nil, -1)
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

func (list *BasicEStoreList) Empty() bool {
	return list.store.IsEmpty(list.owner, list.feature)
}

func (list *BasicEStoreList) Contains(element interface{}) bool {
	result := list.store.Contains(list.owner, list.feature, element)
	if !result && list.object && list.proxies {
		for i := 0; i < list.Size(); i++ {
			eObject, _ := list.store.Get(list.owner, list.feature, i).(EObject)
			eResolved := list.resolveProxy(eObject)
			if eObject == eResolved {
				return true
			}
		}
	}
	return result
}

func (list *BasicEStoreList) IndexOf(element interface{}) int {
	result := list.store.IndexOf(list.owner, list.feature, element)
	if result >= 0 {
		return result
	}
	if list.object && list.proxies {
		for i := 0; i < list.Size(); i++ {
			eObject, _ := list.store.Get(list.owner, list.feature, i).(EObject)
			eResolved := list.resolveProxy(eObject)
			if eObject == eResolved {
				return i
			}
		}
	}
	return -1
}

func (list *BasicEStoreList) Iterator() EIterator {
	return &listIterator{list: list}
}

func (list *BasicEStoreList) ToArray() []interface{} {
	return list.store.ToArray(list.owner, list.feature)
}

func (list *BasicEStoreList) inverseAdd(object interface{}, notifications ENotificationChain) ENotificationChain {
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

func (list *BasicEStoreList) inverseRemove(object interface{}, notifications ENotificationChain) ENotificationChain {
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

func (list *BasicEStoreList) GetUnResolvedList() EList {
	l := NewBasicEStoreList(list.owner, list.feature, list.store)
	l.proxies = false
	return l
}
