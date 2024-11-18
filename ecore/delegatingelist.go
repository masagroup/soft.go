package ecore

type AbstractDelegatingEList[T abstractEList] struct {
	AbstractEList
	delegate T
}

// Add a new element to the array
func (list *AbstractDelegatingEList[T]) doAdd(e any) {
	list.delegate.doAdd(e)
}

func (list *AbstractDelegatingEList[T]) doAddAll(c Collection) bool {
	return list.delegate.doAddAll(c)
}

func (list *AbstractDelegatingEList[T]) doInsert(index int, e any) {
	list.delegate.doInsert(index, e)
}

func (list *AbstractDelegatingEList[T]) doInsertAll(index int, collection Collection) bool {
	return list.delegate.doInsertAll(index, collection)
}

func (list *AbstractDelegatingEList[T]) doMove(oldIndex, newIndex int) any {
	return list.delegate.doMove(oldIndex, newIndex)
}

func (list *AbstractDelegatingEList[T]) doRemove(index int) any {
	return list.delegate.doRemove(index)
}

func (list *AbstractDelegatingEList[T]) doRemoveRange(fromIndex int, toIndex int) []any {
	return list.delegate.doRemoveRange(fromIndex, toIndex)
}

func (list *AbstractDelegatingEList[T]) doGet(index int) any {
	return list.delegate.doGet(index)
}

func (list *AbstractDelegatingEList[T]) doSet(index int, elem any) any {
	return list.delegate.doSet(index, elem)
}

func (list *AbstractDelegatingEList[T]) doClear() []any {
	return list.delegate.doClear()
}

// Size count the number of element in the array
func (list *AbstractDelegatingEList[T]) Size() int {
	return list.delegate.Size()
}

func (list *AbstractDelegatingEList[T]) ToArray() []any {
	return list.delegate.ToArray()
}

type AbstractDelegatingENotifyingList[T abstractENotifyingList] struct {
	AbstractDelegatingEList[T]
}

func (l *AbstractDelegatingENotifyingList[T]) GetNotifier() ENotifier {
	return l.delegate.GetNotifier()
}

func (l *AbstractDelegatingENotifyingList[T]) GetFeature() EStructuralFeature {
	return l.delegate.GetFeature()
}

func (l *AbstractDelegatingENotifyingList[T]) GetFeatureID() int {
	return l.delegate.GetFeatureID()
}

func (l *AbstractDelegatingENotifyingList[T]) AddWithNotification(object any, notifications ENotificationChain) ENotificationChain {
	return l.delegate.AddWithNotification(object, notifications)
}

func (l *AbstractDelegatingENotifyingList[T]) RemoveWithNotification(object any, notifications ENotificationChain) ENotificationChain {
	return l.delegate.RemoveWithNotification(object, notifications)
}

func (l *AbstractDelegatingENotifyingList[T]) SetWithNotification(index int, object any, notifications ENotificationChain) ENotificationChain {
	return l.delegate.SetWithNotification(index, object, notifications)
}
