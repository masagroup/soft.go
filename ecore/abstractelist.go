package ecore

import (
	"iter"
	"strconv"

	"github.com/ugurcsen/gods-generic/sets/linkedhashset"
)

type internalAbstractEList interface {
	EList

	DoGet(index int) any

	DoSet(index int, elem any) any

	DoAdd(elem any)

	DoAddAll(list Collection) bool

	DoInsert(index int, elem any)

	DoInsertAll(index int, list Collection) bool

	DoClear() []any

	DoMove(oldIndex, newIndex int) any

	DoRemove(index int) any

	DoRemoveRange(fromIndex, toIndex int) []any
}

type AbstractEList struct {
	interfaces any
	isUnique   bool
}

func (list *AbstractEList) SetInterfaces(interfaces any) {
	list.interfaces = interfaces
}

func (list *AbstractEList) asEList() EList {
	return list.interfaces.(EList)
}

func (list *AbstractEList) asInternal() internalAbstractEList {
	return list.interfaces.(internalAbstractEList)
}

func (list *AbstractEList) getNonDuplicates(collection Collection) Collection {
	// initialize hashset with collection
	hashSet := linkedhashset.New[any]()
	for it := collection.Iterator(); it.HasNext(); {
		hashSet.Add(it.Next())
	}
	// remove all elements in list
	collection = list.asEList()
	if hashSet.Size() >= collection.Size() {
		for it := collection.Iterator(); it.HasNext(); {
			hashSet.Remove(it.Next())
		}
	} else {
		for it := hashSet.Iterator(); it.Next(); {
			v := it.Value()
			if collection.Contains(v) {
				hashSet.Remove(v)
			}
		}
	}
	array := make([]any, hashSet.Size())
	hashSet.Each(func(index int, value any) {
		array[index] = value
	})
	return NewImmutableEList(array)
}

func (list *AbstractEList) Add(elem any) bool {
	l := list.asInternal()
	if list.isUnique && l.Contains(elem) {
		return false
	}
	l.DoAdd(elem)
	return true
}

func (list *AbstractEList) AddAll(collection Collection) bool {
	if list.isUnique {
		collection = list.getNonDuplicates(collection)
		if collection.Size() == 0 {
			return false
		}
	}
	list.asInternal().DoAddAll(collection)
	return true
}

func (list *AbstractEList) Insert(index int, elem any) bool {
	l := list.asInternal()
	if size := l.Size(); index < 0 || index > size {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(size))
	}
	if list.isUnique && l.Contains(elem) {
		return false
	}
	l.DoInsert(index, elem)
	return true
}

func (list *AbstractEList) InsertAll(index int, collection Collection) bool {
	l := list.asInternal()
	if size := l.Size(); index < 0 || index > size {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(size))
	}
	if list.isUnique {
		collection = list.getNonDuplicates(collection)
		if collection.Size() == 0 {
			return false
		}
	}
	l.DoInsertAll(index, collection)
	return true
}

func (list *AbstractEList) MoveObject(newIndex int, elem any) {
	l := list.asInternal()
	oldIndex := l.IndexOf(elem)
	if oldIndex == -1 {
		panic("Object not found")
	}
	l.DoMove(oldIndex, newIndex)
}

// Swap move an element from oldIndex to newIndex
func (list *AbstractEList) Move(oldIndex, newIndex int) any {
	l := list.asInternal()
	if size := l.Size(); oldIndex < 0 || oldIndex >= size || newIndex < 0 || newIndex > size {
		panic("Index out of bounds: oldIndex=" + strconv.Itoa(oldIndex) + " newIndex=" + strconv.Itoa(newIndex) + " size=" + strconv.Itoa(size))
	}
	return l.DoMove(oldIndex, newIndex)
}

// RemoveAt remove an element at a given position
func (list *AbstractEList) RemoveAt(index int) any {
	l := list.asInternal()
	if size := l.Size(); index < 0 || index >= size {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(size))
	}
	return l.DoRemove(index)
}

// Remove an element in an array
func (list *AbstractEList) Remove(elem any) bool {
	l := list.asInternal()
	index := l.IndexOf(elem)
	if index == -1 {
		return false
	}
	l.DoRemove(index)
	return true
}

func (list *AbstractEList) RemoveRange(fromIndex int, toIndex int) {
	l := list.asInternal()
	size := l.Size()
	if fromIndex < 0 || fromIndex >= size {
		panic("Index out of bounds: fromIndex=" + strconv.Itoa(fromIndex) + " size=" + strconv.Itoa(size))
	}
	if toIndex < 0 || toIndex > size {
		panic("Index out of bounds: toIndex=" + strconv.Itoa(toIndex) + " size=" + strconv.Itoa(size))
	}
	if fromIndex > toIndex {
		panic("Indexes invalid: fromIndex=" + strconv.Itoa(fromIndex) + "must be less than toIndex=" + strconv.Itoa(toIndex))
	}
	l.DoRemoveRange(fromIndex, toIndex)
}

func (list *AbstractEList) RemoveAll(collection Collection) bool {
	modified := false
	l := list.asInternal()
	for i := l.Size() - 1; i >= 0; i-- {
		if collection.Contains(l.DoGet(i)) {
			l.RemoveAt(i)
			modified = true
		}
	}
	return modified
}

// Get an element of the array
func (list *AbstractEList) Get(index int) any {
	l := list.asInternal()
	if size := l.Size(); index < 0 || index >= size {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(size))
	}
	return l.DoGet(index)
}

// Set an element of the array
func (list *AbstractEList) Set(index int, elem any) any {
	l := list.asInternal()
	if size := l.Size(); index < 0 || index >= size {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(size))
	}
	if list.isUnique {
		currIndex := l.IndexOf(elem)
		if currIndex >= 0 && currIndex != index {
			panic("element already in list")
		}
	}
	return l.DoSet(index, elem)
}

func (list *AbstractEList) Clear() {
	list.asInternal().DoClear()
}

func (list *AbstractEList) Size() int {
	panic("not implemented")
}

func (list *AbstractEList) Empty() bool {
	return list.asEList().Size() == 0
}

// Contains return if an array contains or not an element
func (list *AbstractEList) Contains(elem any) bool {
	return list.asEList().IndexOf(elem) != -1
}

func (list *AbstractEList) IndexOf(elem any) int {
	i := 0
	for it := list.asEList().Iterator(); it.HasNext(); i++ {
		if value := it.Next(); value == elem {
			return i
		}
	}
	return -1
}

func (list *AbstractEList) Iterator() EIterator {
	return &listIterator{list: list.asEList()}
}

func (list *AbstractEList) All() iter.Seq[any] {
	return func(yield func(any) bool) {
		l := list.asEList()
		for i := 0; i < l.Size(); i++ {
			if !yield(l.Get(i)) {
				return
			}
		}
	}
}

func (list *AbstractEList) ToArray() []any {
	panic("not implemented")
}
