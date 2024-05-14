package ecore

import (
	"container/list"
	"sync"
)

func zero[T any]() T {
	return *new(T)
}

// Iterates over the keys and values in a linkedHashMap
// from oldest to newest elements.
// Assumes the underlying linkedHashMap is not modified while
// the mapIteratorImpl is in use, except to delete elements that
// have already been iterated over.
type mapIterator[K, V any] interface {
	Next() bool
	Key() K
	Value() V
}

type linkedHashMapIterator[K comparable, V any] struct {
	lh                     *linkedHashMap[K, V]
	key                    K
	value                  V
	next                   *list.Element
	initialized, exhausted bool
}

func (it *linkedHashMapIterator[K, V]) Next() bool {
	// If the mapIteratorImpl has been exhausted, there is no next value.
	if it.exhausted {
		it.key = zero[K]()
		it.value = zero[V]()
		it.next = nil
		return false
	}

	it.lh.lock.RLock()
	defer it.lh.lock.RUnlock()

	// If the mapIteratorImpl was not yet initialized, do it now.
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
	kv := it.next.Value.(keyValue[K, V])
	it.key = kv.key
	it.value = kv.value
	it.next = it.next.Next() // Next time, return next element
	it.exhausted = it.next == nil
	return true
}

func (it *linkedHashMapIterator[K, V]) Key() K {
	return it.key
}

func (it *linkedHashMapIterator[K, V]) Value() V {
	return it.value
}

type keyValue[K, V any] struct {
	key   K
	value V
}

type linkedHashMap[K comparable, V any] struct {
	lock      sync.RWMutex
	entryMap  map[K]*list.Element
	entryList *list.List
}

func newLinkedHashMap[K comparable, V any]() *linkedHashMap[K, V] {
	return &linkedHashMap[K, V]{
		entryMap:  make(map[K]*list.Element),
		entryList: list.New(),
	}
}

func (lh *linkedHashMap[K, V]) Put(key K, val V) {
	lh.lock.Lock()
	defer lh.lock.Unlock()

	lh.put(key, val)
}

func (lh *linkedHashMap[K, V]) Get(key K) (V, bool) {
	lh.lock.RLock()
	defer lh.lock.RUnlock()

	return lh.get(key)
}

func (lh *linkedHashMap[K, V]) Delete(key K) {
	lh.lock.Lock()
	defer lh.lock.Unlock()

	lh.delete(key)
}

func (lh *linkedHashMap[K, V]) Len() int {
	lh.lock.RLock()
	defer lh.lock.RUnlock()

	return lh.len()
}

func (lh *linkedHashMap[K, V]) Oldest() (K, V, bool) {
	lh.lock.RLock()
	defer lh.lock.RUnlock()

	return lh.oldest()
}

func (lh *linkedHashMap[K, V]) Newest() (K, V, bool) {
	lh.lock.RLock()
	defer lh.lock.RUnlock()

	return lh.newest()
}

func (lh *linkedHashMap[K, V]) put(key K, value V) {
	if e, ok := lh.entryMap[key]; ok {
		lh.entryList.MoveToBack(e)
		e.Value = keyValue[K, V]{
			key:   key,
			value: value,
		}
	} else {
		lh.entryMap[key] = lh.entryList.PushBack(keyValue[K, V]{
			key:   key,
			value: value,
		})
	}
}

func (lh *linkedHashMap[K, V]) get(key K) (V, bool) {
	if e, ok := lh.entryMap[key]; ok {
		kv := e.Value.(keyValue[K, V])
		return kv.value, true
	}
	return zero[V](), false
}

func (lh *linkedHashMap[K, V]) delete(key K) {
	if e, ok := lh.entryMap[key]; ok {
		lh.entryList.Remove(e)
		delete(lh.entryMap, key)
	}
}

func (lh *linkedHashMap[K, V]) len() int {
	return len(lh.entryMap)
}

func (lh *linkedHashMap[K, V]) oldest() (K, V, bool) {
	if val := lh.entryList.Front(); val != nil {
		kv := val.Value.(keyValue[K, V])
		return kv.key, kv.value, true
	}
	return zero[K](), zero[V](), false
}

func (lh *linkedHashMap[K, V]) newest() (K, V, bool) {
	if val := lh.entryList.Back(); val != nil {
		kv := val.Value.(keyValue[K, V])
		return kv.key, kv.value, true
	}
	return zero[K](), zero[V](), false
}

func (lh *linkedHashMap[K, V]) Iterator() mapIterator[K, V] {
	return &linkedHashMapIterator[K, V]{lh: lh}
}
