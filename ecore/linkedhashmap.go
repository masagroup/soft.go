package ecore

import "iter"

func zero[T any]() T {
	return *new(T)
}

type keyValue[K, V any] struct {
	key   K
	value V
}

// linkedHashMap provides an ordered O(1) mapping from keys to values.
//
// Entries are tracked by insertion order.
type linkedHashMap[K comparable, V any] struct {
	entryMap  map[K]*linkedListElement[keyValue[K, V]]
	entryList *linkedList[keyValue[K, V]]
	freeList  []*linkedListElement[keyValue[K, V]]
}

func newLinkedHashMap[K comparable, V any]() *linkedHashMap[K, V] {
	return newLinkedHashMapWithSize[K, V](0)
}

func newLinkedHashMapWithSize[K comparable, V any](initialSize int) *linkedHashMap[K, V] {
	lh := &linkedHashMap[K, V]{
		entryMap:  make(map[K]*linkedListElement[keyValue[K, V]], initialSize),
		entryList: newLinkedList[keyValue[K, V]](),
		freeList:  make([]*linkedListElement[keyValue[K, V]], initialSize),
	}
	for i := range lh.freeList {
		lh.freeList[i] = &linkedListElement[keyValue[K, V]]{}
	}
	return lh
}

func (lh *linkedHashMap[K, V]) Put(key K, value V) {
	if e, ok := lh.entryMap[key]; ok {
		lh.entryList.MoveToBack(e)
		e.Value = keyValue[K, V]{
			key:   key,
			value: value,
		}
		return
	}

	var e *linkedListElement[keyValue[K, V]]
	if numFree := len(lh.freeList); numFree > 0 {
		numFree--
		e = lh.freeList[numFree]
		lh.freeList = lh.freeList[:numFree]
	} else {
		e = &linkedListElement[keyValue[K, V]]{}
	}

	e.Value = keyValue[K, V]{
		key:   key,
		value: value,
	}
	lh.entryMap[key] = e
	lh.entryList.PushBack(e)
}

func (lh *linkedHashMap[K, V]) Get(key K) (V, bool) {
	if e, ok := lh.entryMap[key]; ok {
		return e.Value.value, true
	}
	return zero[V](), false
}

func (lh *linkedHashMap[K, V]) Delete(key K) bool {
	e, ok := lh.entryMap[key]
	if ok {
		lh.remove(e)
	}
	return ok
}

func (lh *linkedHashMap[K, V]) Clear() {
	for _, e := range lh.entryMap {
		lh.remove(e)
	}
}

// remove assumes that [e] is currently in the Hashmap.
func (lh *linkedHashMap[K, V]) remove(e *linkedListElement[keyValue[K, V]]) {
	delete(lh.entryMap, e.Value.key)
	lh.entryList.Remove(e)
	e.Value = keyValue[K, V]{} // Free the key value pair
	lh.freeList = append(lh.freeList, e)
}

func (lh *linkedHashMap[K, V]) Len() int {
	return len(lh.entryMap)
}

func (lh *linkedHashMap[K, V]) Oldest() (K, V, bool) {
	if e := lh.entryList.Front(); e != nil {
		return e.Value.key, e.Value.value, true
	}
	return zero[K](), zero[V](), false
}

func (lh *linkedHashMap[K, V]) Newest() (K, V, bool) {
	if e := lh.entryList.Back(); e != nil {
		return e.Value.key, e.Value.value, true
	}
	return zero[K](), zero[V](), false
}

func (lh *linkedHashMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for e := lh.entryList.Front(); e != nil; {
			next := e.Next()
			if !yield(e.Value.key, e.Value.value) {
				return
			}
			e = next
		}
	}
}

func (lh *linkedHashMap[K, V]) Iterator() *Iterator[K, V] {
	return &Iterator[K, V]{lh: lh}
}

// Iterates over the keys and values in a LinkedHashmap from oldest to newest.
// Assumes the underlying LinkedHashmap is not modified while the iterator is in
// use, except to delete elements that have already been iterated over.
type Iterator[K comparable, V any] struct {
	lh                     *linkedHashMap[K, V]
	key                    K
	value                  V
	next                   *linkedListElement[keyValue[K, V]]
	initialized, exhausted bool
}

func (it *Iterator[K, V]) HasNext() bool {
	// If the iterator has been exhausted, there is no next value.
	if it.exhausted {
		it.key = zero[K]()
		it.value = zero[V]()
		it.next = nil
		return false
	}

	// If the iterator was not yet initialized, do it now.
	if !it.initialized {
		it.initialized = true
		oldest := it.lh.entryList.Front()
		if oldest == nil {
			it.exhausted = true
			it.key = zero[K]()
			it.value = zero[V]()
			it.next = nil
			return false
		}
		it.next = oldest
	}

	// It's important to ensure that [it.next] is not nil
	// by not deleting elements that have not yet been iterated
	// over from [it.lh]
	it.key = it.next.Value.key
	it.value = it.next.Value.value
	it.next = it.next.Next() // Next time, return next element
	it.exhausted = it.next == nil
	return true
}

func (it *Iterator[K, V]) Key() K {
	return it.key
}

func (it *Iterator[K, V]) Value() V {
	return it.value
}
