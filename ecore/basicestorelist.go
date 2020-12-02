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
	object := list.store.Move(list.owner, list.feature, newIndex, oldIndex)
	list.createAndDispatchNotification(nil, MOVE, oldIndex, object, newIndex)
	return object
}

func (list *BasicEStoreList) Get(int) interface{} {
	return nil
}

func (list *BasicEStoreList) Set(int, interface{}) interface{} {
	return nil
}

func (list *BasicEStoreList) RemoveAt(int) interface{} {
	return nil
}

func (list *BasicEStoreList) Remove(interface{}) bool {
	return false
}

func (list *BasicEStoreList) RemoveAll(EList) bool {
	return false
}

func (list *BasicEStoreList) Size() int {
	return 0
}

func (list *BasicEStoreList) Clear() {

}

func (list *BasicEStoreList) Empty() bool {
	return false
}

func (list *BasicEStoreList) Contains(interface{}) bool {
	return false
}

func (list *BasicEStoreList) IndexOf(interface{}) int {
	return 0
}

func (list *BasicEStoreList) Iterator() EIterator {
	return nil
}

func (list *BasicEStoreList) ToArray() []interface{} {
	return nil
}

func (list *BasicEStoreList) inverseAdd(object interface{}, notifications ENotificationChain) ENotificationChain {
	return notifications
}

func (list *BasicEStoreList) inverseRemove(object interface{}, notifications ENotificationChain) ENotificationChain {
	return notifications
}