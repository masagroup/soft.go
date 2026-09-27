package ecore

type AbstractDelegatingEList[T internalAbstractEList] struct {
	AbstractEList
	delegate T
}

// Add a new element to the array
func (list *AbstractDelegatingEList[T]) DoAdd(e any) {
	list.delegate.DoAdd(e)
}

func (list *AbstractDelegatingEList[T]) DoAddAll(c Collection) bool {
	return list.delegate.DoAddAll(c)
}

func (list *AbstractDelegatingEList[T]) DoInsert(index int, e any) {
	list.delegate.DoInsert(index, e)
}

func (list *AbstractDelegatingEList[T]) DoInsertAll(index int, collection Collection) bool {
	return list.delegate.DoInsertAll(index, collection)
}

func (list *AbstractDelegatingEList[T]) DoMove(oldIndex, newIndex int) any {
	return list.delegate.DoMove(oldIndex, newIndex)
}

func (list *AbstractDelegatingEList[T]) DoRemove(index int) any {
	return list.delegate.DoRemove(index)
}

func (list *AbstractDelegatingEList[T]) DoRemoveRange(fromIndex int, toIndex int) []any {
	return list.delegate.DoRemoveRange(fromIndex, toIndex)
}

func (list *AbstractDelegatingEList[T]) DoGet(index int) any {
	return list.delegate.DoGet(index)
}

func (list *AbstractDelegatingEList[T]) DoSet(index int, elem any) any {
	return list.delegate.DoSet(index, elem)
}

func (list *AbstractDelegatingEList[T]) DoClear() []any {
	return list.delegate.DoClear()
}

// Size count the number of element in the array
func (list *AbstractDelegatingEList[T]) Size() int {
	return list.delegate.Size()
}

func (list *AbstractDelegatingEList[T]) ToArray() []any {
	return list.delegate.ToArray()
}

type AbstractDelegatingENotifyingList[T internalENotifyingList] struct {
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
