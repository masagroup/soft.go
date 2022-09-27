package ecore

type eNotifyingListInternal interface {
	ENotifyingList

	inverseAdd(object any, chain ENotificationChain) ENotificationChain

	inverseRemove(object any, chain ENotificationChain) ENotificationChain
}

// BasicENotifyingList ...
type BasicENotifyingList struct {
	BasicEList
}

// NewBasicENotifyingList ...
func NewBasicENotifyingList() *BasicENotifyingList {
	l := new(BasicENotifyingList)
	l.interfaces = l
	l.data = []any{}
	l.isUnique = true
	return l
}

func newBasicENotifyingListFromData(data []any) *BasicENotifyingList {
	l := new(BasicENotifyingList)
	l.interfaces = l
	l.data = data
	l.isUnique = true
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
	AbstractNotification
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

func (list *BasicENotifyingList) createNotification(eventType EventType, oldValue any, newValue any, position int) ENotification {
	n := new(notifyingListNotification)
	n.Initialize(n, eventType, oldValue, newValue, position)
	n.list = list
	return n
}

func (list *BasicENotifyingList) isNotificationRequired() bool {
	notifier := list.interfaces.(ENotifyingList).GetNotifier()
	return notifier != nil && notifier.EDeliver() && !notifier.EAdapters().Empty()
}

func (list *BasicENotifyingList) createAndAddNotification(ns ENotificationChain, eventType EventType, oldValue any, newValue any, position int) ENotificationChain {
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

func (list *BasicENotifyingList) createAndDispatchNotification(notifications ENotificationChain, eventType EventType, oldValue any, newValue any, position int) {
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

func (list *BasicENotifyingList) inverseAdd(object any, notifications ENotificationChain) ENotificationChain {
	return notifications
}

func (list *BasicENotifyingList) inverseRemove(object any, notifications ENotificationChain) ENotificationChain {
	return notifications
}

// AddWithNotification ...
func (list *BasicENotifyingList) AddWithNotification(object any, notifications ENotificationChain) ENotificationChain {
	index := list.Size()
	list.BasicEList.doAdd(object)
	return list.createAndAddNotification(notifications, ADD, nil, object, index)
}

// RemoveWithNotification ...
func (list *BasicENotifyingList) RemoveWithNotification(object any, notifications ENotificationChain) ENotificationChain {
	index := list.interfaces.(EList).IndexOf(object)
	if index != -1 {
		oldObject := list.BasicEList.doRemove(index)
		return list.createAndAddNotification(notifications, REMOVE, oldObject, nil, index)
	}
	return notifications
}

// SetWithNotification ...
func (list *BasicENotifyingList) SetWithNotification(index int, object any, notifications ENotificationChain) ENotificationChain {
	oldObject := list.BasicEList.doSet(index, object)
	return list.createAndAddNotification(notifications, SET, oldObject, object, index)
}

func (list *BasicENotifyingList) doAdd(object any) {
	index := list.interfaces.(EList).Size()
	list.BasicEList.doAdd(object)
	notifications := list.interfaces.(eNotifyingListInternal).inverseAdd(object, nil)
	list.createAndDispatchNotification(notifications, ADD, nil, object, index)
}

func (list *BasicENotifyingList) doAddAll(l EList) bool {
	return list.interfaces.(abstractEList).doInsertAll(list.Size(), l)
}

func (list *BasicENotifyingList) doInsert(index int, object any) {
	list.BasicEList.doInsert(index, object)
	notifications := list.interfaces.(eNotifyingListInternal).inverseAdd(object, nil)
	list.createAndDispatchNotification(notifications, ADD, nil, object, index)
}

func (list *BasicENotifyingList) doInsertAll(index int, l EList) bool {
	if l.Empty() {
		return false
	}

	result := list.BasicEList.doInsertAll(index, l)
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

func (list *BasicENotifyingList) doSet(index int, newObject any) any {
	oldObject := list.BasicEList.doSet(index, newObject)
	if newObject != oldObject {
		var notifications ENotificationChain
		notifications = list.interfaces.(eNotifyingListInternal).inverseRemove(oldObject, notifications)
		notifications = list.interfaces.(eNotifyingListInternal).inverseAdd(newObject, notifications)
		list.createAndDispatchNotification(notifications, SET, oldObject, newObject, index)
	}
	return oldObject
}

func (list *BasicENotifyingList) doClear() []any {
	oldData := list.BasicEList.doClear()
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
	return oldData
}

func (list *BasicENotifyingList) doMove(oldIndex, newIndex int) any {
	oldObject := list.BasicEList.doMove(oldIndex, newIndex)
	list.createAndDispatchNotification(nil, MOVE, oldIndex, oldObject, newIndex)
	return oldObject
}

// RemoveAt ...
func (list *BasicENotifyingList) doRemove(index int) any {
	oldObject := list.BasicEList.doRemove(index)
	var notifications ENotificationChain
	notifications = list.interfaces.(eNotifyingListInternal).inverseRemove(oldObject, notifications)
	list.createAndDispatchNotification(notifications, REMOVE, oldObject, nil, index)
	return oldObject
}

func (list *BasicENotifyingList) RemoveAll(collection EList) bool {
	return list.doRemoveAll(
		collection,
		func(index int, other any) bool {
			return list.doGet(index) == other
		})
}

func (list *BasicENotifyingList) doRemoveAll(collection EList, getAndCompare func(int, any) bool) bool {
	var positions []any
	var removed []any

	// compute positions and removed objects
	if !collection.Empty() {
		for i := 0; i < list.Size(); i++ {
			for j := 0; j < collection.Size(); j++ {
				object := collection.Get(j)
				if getAndCompare(i, object) {
					positions = append(positions, i)
					removed = append(removed, object)
					break
				}
			}
		}
	}

	// remove
	for i := len(positions) - 1; i >= 0; i-- {
		list.BasicEList.doRemove(positions[i].(int))
	}

	// inverse remove
	var notifications ENotificationChain
	for _, e := range removed {
		notifications = list.interfaces.(eNotifyingListInternal).inverseRemove(e, notifications)
	}

	// notifications
	if removed != nil {
		list.createAndDispatchNotificationFn(notifications,
			func() ENotification {
				if len(removed) == 1 {
					return list.createNotification(REMOVE, removed[0], nil, positions[0].(int))
				} else {
					return list.createNotification(REMOVE_MANY, removed, positions, positions[0].(int))
				}
			})
	} else if notifications != nil {
		notifications.Dispatch()
	}

	return removed != nil
}
