package ecore

type internalENotifyingList interface {
	ENotifyingList
	internalAbstractEList

	PerformAdd(object any)

	PerformInsert(index int, object any)

	PerformInsertAll(index int, list Collection) bool

	PerformClear() []any

	PerformRemove(index int) any

	PerformRemoveRange(fromIndex int, toIndex int) []any

	PerformSet(index int, object any) any

	PerformMove(oldIndex, newIndex int) any

	InverseAdd(object any, chain ENotificationChain) ENotificationChain

	InverseRemove(object any, chain ENotificationChain) ENotificationChain
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

func (list *BasicENotifyingList) asENotifyingList() ENotifyingList {
	return list.interfaces.(ENotifyingList)
}

func (list *BasicENotifyingList) asInternal() internalENotifyingList {
	return list.interfaces.(internalENotifyingList)
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
	return notif.list.asENotifyingList().GetNotifier()
}

func (notif *notifyingListNotification) GetFeature() EStructuralFeature {
	return notif.list.asENotifyingList().GetFeature()
}

func (notif *notifyingListNotification) GetFeatureID() int {
	return notif.list.asENotifyingList().GetFeatureID()
}

func (list *BasicENotifyingList) createNotification(eventType EventType, oldValue any, newValue any, position int) ENotification {
	n := new(notifyingListNotification)
	n.Initialize(n, eventType, oldValue, newValue, position)
	n.list = list
	return n
}

func (list *BasicENotifyingList) isNotificationRequired() bool {
	notifier := list.asENotifyingList().GetNotifier()
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
			notifier := list.asENotifyingList().GetNotifier()
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

func (list *BasicENotifyingList) InverseAdd(object any, notifications ENotificationChain) ENotificationChain {
	return notifications
}

func (list *BasicENotifyingList) InverseRemove(object any, notifications ENotificationChain) ENotificationChain {
	return notifications
}

// AddWithNotification ...
func (list *BasicENotifyingList) AddWithNotification(object any, notifications ENotificationChain) ENotificationChain {
	notifyingList := list.asInternal()
	index := notifyingList.Size()
	notifyingList.PerformAdd(object)
	return list.createAndAddNotification(notifications, ADD, nil, object, index)
}

// RemoveWithNotification ...
func (list *BasicENotifyingList) RemoveWithNotification(object any, notifications ENotificationChain) ENotificationChain {
	notifyingList := list.asInternal()
	index := notifyingList.IndexOf(object)
	if index != -1 {

		oldObject := notifyingList.PerformRemove(index)
		return list.createAndAddNotification(notifications, REMOVE, oldObject, nil, index)
	}
	return notifications
}

// SetWithNotification ...
func (list *BasicENotifyingList) SetWithNotification(index int, object any, notifications ENotificationChain) ENotificationChain {
	notifyingList := list.asInternal()
	oldObject := notifyingList.PerformSet(index, object)
	return list.createAndAddNotification(notifications, SET, oldObject, object, index)
}

func (list *BasicENotifyingList) DoAdd(object any) {
	index := list.asEList().Size()
	notifyingList := list.asInternal()
	notifyingList.PerformAdd(object)
	notifications := notifyingList.InverseAdd(object, nil)
	list.createAndDispatchNotification(notifications, ADD, nil, object, index)
}

func (list *BasicENotifyingList) DoAddAll(collection Collection) bool {
	notifyingList := list.asInternal()
	return list.DoInsertAll(notifyingList.Size(), collection)
}

func (list *BasicENotifyingList) DoInsert(index int, object any) {
	notifyingList := list.asInternal()
	notifyingList.PerformInsert(index, object)
	notifications := notifyingList.InverseAdd(object, nil)
	list.createAndDispatchNotification(notifications, ADD, nil, object, index)
}

func (list *BasicENotifyingList) DoInsertAll(index int, l Collection) bool {
	if l.Empty() {
		return false
	}
	notifyingList := list.asInternal()
	result := notifyingList.PerformInsertAll(index, l)
	var notifications ENotificationChain = NewNotificationChain()
	for it := l.Iterator(); it.HasNext(); {
		notifications = notifyingList.InverseAdd(it.Next(), notifications)
	}
	list.createAndDispatchNotificationFn(notifications, func() ENotification {
		if l.Size() == 1 {
			return list.createNotification(ADD, nil, l.Iterator().Next(), index)
		} else {
			return list.createNotification(ADD_MANY, nil, l.ToArray(), index)
		}
	})
	return result
}

func (list *BasicENotifyingList) DoSet(index int, newObject any) any {
	notifyingList := list.asInternal()
	oldObject := notifyingList.PerformSet(index, newObject)
	if newObject != oldObject {
		var notifications ENotificationChain
		notifications = notifyingList.InverseRemove(oldObject, notifications)
		notifications = notifyingList.InverseAdd(newObject, notifications)
		list.createAndDispatchNotification(notifications, SET, oldObject, newObject, index)
	}
	return oldObject
}

func (list *BasicENotifyingList) DoClear() []any {
	notifyingList := list.asInternal()
	oldData := notifyingList.PerformClear()
	if len(oldData) == 0 {
		list.createAndDispatchNotification(nil, REMOVE_MANY, []any{}, nil, -1)
	} else {
		var notifications ENotificationChain = NewNotificationChain()
		for _, e := range oldData {
			notifications = notifyingList.InverseRemove(e, notifications)
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

func (list *BasicENotifyingList) DoMove(oldIndex, newIndex int) any {
	notifyingList := list.asInternal()
	oldObject := notifyingList.PerformMove(oldIndex, newIndex)
	list.createAndDispatchNotification(nil, MOVE, oldIndex, oldObject, newIndex)
	return oldObject
}

func (list *BasicENotifyingList) DoRemove(index int) any {
	notifyingList := list.asInternal()
	oldObject := notifyingList.PerformRemove(index)
	// inverse remove
	var notifications ENotificationChain
	notifications = notifyingList.InverseRemove(oldObject, notifications)
	// notifications
	list.createAndDispatchNotification(notifications, REMOVE, oldObject, nil, index)
	return oldObject
}

func (list *BasicENotifyingList) DoRemoveRange(fromIndex int, toIndex int) []any {
	notifyingList := list.asInternal()
	objects := notifyingList.PerformRemoveRange(fromIndex, toIndex)
	if len(objects) > 0 {
		// inverse remove
		var notifications ENotificationChain
		for _, object := range objects {
			notifications = notifyingList.InverseRemove(object, notifications)
		}
		// notifications
		list.createAndDispatchNotificationFn(notifications,
			func() ENotification {
				if len(objects) == 1 {
					return list.createNotification(REMOVE, objects[0], nil, fromIndex)
				} else {
					positions := make([]any, len(objects))
					for i := range objects {
						positions[i] = fromIndex + i
					}
					return list.createNotification(REMOVE_MANY, objects, positions, fromIndex)
				}
			})
	}
	return objects
}

func (list *BasicENotifyingList) RemoveAll(collection Collection) bool {
	return list.doRemoveAll(
		collection,
		func(index int, other any) bool {
			return list.asInternal().DoGet(index) == other
		})
}

func (list *BasicENotifyingList) doRemoveAll(collection Collection, getAndCompare func(int, any) bool) bool {
	var positions []any
	var removed []any
	notifyingList := list.asInternal()

	// compute positions and removed objects
	if !collection.Empty() {
		for i := 0; i < notifyingList.Size(); i++ {
			j := 0
			for itCollection := collection.Iterator(); itCollection.HasNext(); j++ {
				object := itCollection.Next()
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
		notifyingList.PerformRemove(positions[i].(int))
	}

	// inverse remove
	var notifications ENotificationChain
	for _, e := range removed {
		notifications = notifyingList.InverseRemove(e, notifications)
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

func (list *BasicENotifyingList) PerformAdd(object any) {
	list.BasicEList.DoAdd(object)
}

func (list *BasicENotifyingList) PerformInsert(index int, object any) {
	list.BasicEList.DoInsert(index, object)
}

func (list *BasicENotifyingList) PerformInsertAll(index int, l Collection) bool {
	return list.BasicEList.DoInsertAll(index, l)
}

func (list *BasicENotifyingList) PerformClear() []any {
	return list.BasicEList.DoClear()
}

func (list *BasicENotifyingList) PerformRemove(index int) any {
	return list.BasicEList.DoRemove(index)
}

func (list *BasicENotifyingList) PerformRemoveRange(fromIndex int, toIndex int) []any {
	return list.BasicEList.DoRemoveRange(fromIndex, toIndex)
}

func (list *BasicENotifyingList) PerformSet(index int, object any) any {
	return list.BasicEList.DoSet(index, object)
}

func (list *BasicENotifyingList) PerformMove(oldIndex, newIndex int) any {
	return list.BasicEList.DoMove(oldIndex, newIndex)
}
