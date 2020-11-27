package ecore

type eNotifyingListInternal interface {
	ENotifyingList

	inverseAdd(object interface{}, chain ENotificationChain) ENotificationChain

	inverseRemove(object interface{}, chain ENotificationChain) ENotificationChain
}

// BasicENotifyingList ...
type BasicENotifyingList struct {
	*basicEList
}

// NewBasicENotifyingList ...
func NewBasicENotifyingList() *BasicENotifyingList {
	l := new(BasicENotifyingList)
	l.basicEList = NewUniqueBasicEList([]interface{}{})
	l.interfaces = l
	return l
}

func newBasicENotifyingListFromData(data []interface{}) *BasicENotifyingList {
	l := new(BasicENotifyingList)
	l.basicEList = NewUniqueBasicEList(data)
	l.interfaces = l
	return l
}

// GetNotifier ...
func (list *BasicENotifyingList) GetNotifier() ENotifier {
	return nil
}

// GetFeature ...
func (list *BasicENotifyingList) GetFeature() EStructuralFeature {
	return nil
}

// GetFeatureID ...
func (list *BasicENotifyingList) GetFeatureID() int {
	return -1
}

type notifyingListNotification struct {
	*abstractNotification
	list *BasicENotifyingList
}

func (notif *notifyingListNotification) GetNotifier() ENotifier {
	return notif.list.interfaces.(ENotifyingList).GetNotifier()
}

func (notif *notifyingListNotification) GetFeature() EStructuralFeature {
	return notif.list.interfaces.(ENotifyingList).GetFeature()
}

func (notif *notifyingListNotification) GetFeatureID() int {
	return notif.list.interfaces.(ENotifyingList).GetFeatureID()
}

func (list *BasicENotifyingList) createNotification(eventType EventType, oldValue interface{}, newValue interface{}, position int) ENotification {
	n := new(notifyingListNotification)
	n.abstractNotification = NewAbstractNotification(eventType, oldValue, newValue, position)
	n.interfaces = n
	n.list = list
	return n
}

func (list *BasicENotifyingList) isNotificationRequired() bool {
	notifier := list.interfaces.(ENotifyingList).GetNotifier()
	return notifier != nil && notifier.EDeliver() && !notifier.EAdapters().Empty()
}

func (list *BasicENotifyingList) createAndAddNotification(ns ENotificationChain, eventType EventType, oldValue interface{}, newValue interface{}, position int) ENotificationChain {
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

func (list *BasicENotifyingList) createAndDispatchNotification(notifications ENotificationChain, eventType EventType, oldValue interface{}, newValue interface{}, position int) {
	list.createAndDispatchNotificationFn(notifications, func() ENotification {
		return list.createNotification(eventType, oldValue, newValue, position)
	})
}

func (list *BasicENotifyingList) createAndDispatchNotificationFn(notifications ENotificationChain, createNotification func() ENotification) {
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

func (list *BasicENotifyingList) inverseAdd(object interface{}, notifications ENotificationChain) ENotificationChain {
	return notifications
}

func (list *BasicENotifyingList) inverseRemove(object interface{}, notifications ENotificationChain) ENotificationChain {
	return notifications
}

// AddWithNotification ...
func (list *BasicENotifyingList) AddWithNotification(object interface{}, notifications ENotificationChain) ENotificationChain {
	index := list.Size()
	list.basicEList.doAdd(object)
	return list.createAndAddNotification(notifications, ADD, nil, object, index)
}

// RemoveWithNotification ...
func (list *BasicENotifyingList) RemoveWithNotification(object interface{}, notifications ENotificationChain) ENotificationChain {
	index := list.IndexOf(object)
	if index != -1 {
		oldObject := list.basicEList.doRemove(index)
		return list.createAndAddNotification(notifications, REMOVE, oldObject, nil, index)
	}
	return notifications
}

// SetWithNotification ...
func (list *BasicENotifyingList) SetWithNotification(index int, object interface{}, notifications ENotificationChain) ENotificationChain {
	oldObject := list.basicEList.doSet(index, object)
	return list.createAndAddNotification(notifications, SET, oldObject, object, index)
}

func (list *BasicENotifyingList) doAdd(object interface{}) {
	index := list.Size()
	list.basicEList.doAdd(object)
	notifications := list.interfaces.(eNotifyingListInternal).inverseAdd(object, nil)
	list.createAndDispatchNotification(notifications, ADD, nil, object, index)
}

func (list *BasicENotifyingList) doAddAll(l EList) bool {
	return list.doInsertAll(list.Size(), l)
}

func (list *BasicENotifyingList) doInsert(index int, object interface{}) {
	list.basicEList.doInsert(index, object)
	notifications := list.interfaces.(eNotifyingListInternal).inverseAdd(object, nil)
	list.createAndDispatchNotification(notifications, ADD, nil, object, index)
}

func (list *BasicENotifyingList) doInsertAll(index int, l EList) bool {
	if l.Empty() {
		return false
	}

	result := list.basicEList.doInsertAll(index, l)
	var notifications ENotificationChain = NewNotificationChain()
	for it := l.Iterator(); it.HasNext(); {
		notifications = list.interfaces.(eNotifyingListInternal).inverseAdd(it.Next(), notifications)
	}
	list.createAndDispatchNotificationFn(notifications, func() ENotification {
		if l.Size() == 1 {
			return list.createNotification(ADD, nil, l.Get(0), index)
		} else {
			return list.createNotification(ADD_MANY, nil, l.ToArray(), index)
		}
	})
	return result
}

func (list *BasicENotifyingList) doSet(index int, newObject interface{}) interface{} {
	oldObject := list.basicEList.doSet(index, newObject)
	if newObject != oldObject {
		var notifications ENotificationChain
		notifications = list.interfaces.(eNotifyingListInternal).inverseRemove(oldObject, notifications)
		notifications = list.interfaces.(eNotifyingListInternal).inverseAdd(newObject, notifications)
		list.createAndDispatchNotification(notifications, SET, oldObject, newObject, index)
	}
	return oldObject
}

func (list *BasicENotifyingList) doClear() []interface{} {
	oldData := list.basicEList.doClear()
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
	return oldData
}

func (list *BasicENotifyingList) doMove(oldIndex, newIndex int) interface{} {
	oldObject := list.basicEList.doMove(oldIndex, newIndex)
	list.createAndDispatchNotification(nil, MOVE, oldIndex, oldObject, newIndex)
	return oldObject
}

// RemoveAt ...
func (list *BasicENotifyingList) doRemove(index int) interface{} {
	oldObject := list.basicEList.doRemove(index)
	var notifications ENotificationChain
	notifications = list.interfaces.(eNotifyingListInternal).inverseRemove(oldObject, notifications)
	list.createAndDispatchNotification(notifications, REMOVE, oldObject, nil, index)
	return oldObject
}
