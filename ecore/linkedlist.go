package ecore

// linkedListElement is an element of a linked list.
type linkedListElement[T any] struct {
	next, prev *linkedListElement[T]
	list       *linkedList[T]
	Value      T
}

// Next returns the next element or nil.
func (e *linkedListElement[T]) Next() *linkedListElement[T] {
	if p := e.next; e.list != nil && p != &e.list.sentinel {
		return p
	}
	return nil
}

// Prev returns the previous element or nil.
func (e *linkedListElement[T]) Prev() *linkedListElement[T] {
	if p := e.prev; e.list != nil && p != &e.list.sentinel {
		return p
	}
	return nil
}

// linkedList implements a doubly linked list with a sentinel node.
//
// See: https://en.wikipedia.org/wiki/Doubly_linked_list
//
// This datastructure is designed to be an almost complete drop-in replacement
// for the standard library's "container/list".
//
// The primary design change is to remove all memory allocations from the list
// definition. This allows these lists to be used in performance critical paths.
// Additionally the zero value is not useful. Lists must be created with the
// NewList method.
type linkedList[T any] struct {
	// sentinel is only used as a placeholder to avoid complex nil checks.
	// sentinel.Value is never used.
	sentinel linkedListElement[T]
	length   int
}

// NewList creates a new doubly linked list.
func newLinkedList[T any]() *linkedList[T] {
	l := &linkedList[T]{}
	l.sentinel.next = &l.sentinel
	l.sentinel.prev = &l.sentinel
	l.sentinel.list = l
	return l
}

// Len returns the number of elements in l.
func (l *linkedList[_]) Len() int {
	return l.length
}

// Front returns the element at the front of l.
// If l is empty, nil is returned.
func (l *linkedList[T]) Front() *linkedListElement[T] {
	if l.length == 0 {
		return nil
	}
	return l.sentinel.next
}

// Back returns the element at the back of l.
// If l is empty, nil is returned.
func (l *linkedList[T]) Back() *linkedListElement[T] {
	if l.length == 0 {
		return nil
	}
	return l.sentinel.prev
}

// Remove removes e from l if e is in l.
func (l *linkedList[T]) Remove(e *linkedListElement[T]) {
	if e.list != l {
		return
	}

	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil
	e.prev = nil
	e.list = nil
	l.length--
}

// PushFront inserts e at the front of l.
// If e is already in a list, l is not modified.
func (l *linkedList[T]) PushFront(e *linkedListElement[T]) {
	l.insertAfter(e, &l.sentinel)
}

// PushBack inserts e at the back of l.
// If e is already in a list, l is not modified.
func (l *linkedList[T]) PushBack(e *linkedListElement[T]) {
	l.insertAfter(e, l.sentinel.prev)
}

// InsertBefore inserts e immediately before location.
// If e is already in a list, l is not modified.
// If location is not in l, l is not modified.
func (l *linkedList[T]) InsertBefore(e *linkedListElement[T], location *linkedListElement[T]) {
	if location.list == l {
		l.insertAfter(e, location.prev)
	}
}

// InsertAfter inserts e immediately after location.
// If e is already in a list, l is not modified.
// If location is not in l, l is not modified.
func (l *linkedList[T]) InsertAfter(e *linkedListElement[T], location *linkedListElement[T]) {
	if location.list == l {
		l.insertAfter(e, location)
	}
}

// MoveToFront moves e to the front of l.
// If e is not in l, l is not modified.
func (l *linkedList[T]) MoveToFront(e *linkedListElement[T]) {
	// If e is already at the front of l, there is nothing to do.
	if e != l.sentinel.next {
		l.moveAfter(e, &l.sentinel)
	}
}

// MoveToBack moves e to the back of l.
// If e is not in l, l is not modified.
func (l *linkedList[T]) MoveToBack(e *linkedListElement[T]) {
	l.moveAfter(e, l.sentinel.prev)
}

// MoveBefore moves e immediately before location.
// If the elements are equal or not in l, the list is not modified.
func (l *linkedList[T]) MoveBefore(e, location *linkedListElement[T]) {
	// Don't introduce a cycle by moving an element before itself.
	if e != location {
		l.moveAfter(e, location.prev)
	}
}

// MoveAfter moves e immediately after location.
// If the elements are equal or not in l, the list is not modified.
func (l *linkedList[T]) MoveAfter(e, location *linkedListElement[T]) {
	l.moveAfter(e, location)
}

func (l *linkedList[T]) insertAfter(e, location *linkedListElement[T]) {
	if e.list != nil {
		// Don't insert an element that is already in a list
		return
	}

	e.prev = location
	e.next = location.next
	e.prev.next = e
	e.next.prev = e
	e.list = l
	l.length++
}

func (l *linkedList[T]) moveAfter(e, location *linkedListElement[T]) {
	if e.list != l || location.list != l || e == location {
		// Don't modify an element that is in a different list.
		// Don't introduce a cycle by moving an element after itself.
		return
	}

	e.prev.next = e.next
	e.next.prev = e.prev

	e.prev = location
	e.next = location.next
	e.prev.next = e
	e.next.prev = e
}

// PushFront inserts v into a new element at the front of l.
func PushFront[T any](l *linkedList[T], v T) {
	l.PushFront(&linkedListElement[T]{
		Value: v,
	})
}

// PushBack inserts v into a new element at the back of l.
func PushBack[T any](l *linkedList[T], v T) {
	l.PushBack(&linkedListElement[T]{
		Value: v,
	})
}

// InsertBefore inserts v into a new element immediately before location.
// If location is not in l, l is not modified.
func InsertBefore[T any](l *linkedList[T], v T, location *linkedListElement[T]) {
	l.InsertBefore(
		&linkedListElement[T]{
			Value: v,
		},
		location,
	)
}

// InsertAfter inserts v into a new element immediately after location.
// If location is not in l, l is not modified.
func InsertAfter[T any](l *linkedList[T], v T, location *linkedListElement[T]) {
	l.InsertAfter(
		&linkedListElement[T]{
			Value: v,
		},
		location,
	)
}
