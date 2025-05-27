package ecore

import (
	"iter"
	"strconv"

	"github.com/ugurcsen/gods-generic/sets/linkedhashset"
)

type abstractEList interface {
	EList

	doGet(index int) any

	doSet(index int, elem any) any

	doAdd(elem any)

	doAddAll(list Collection) bool

	doInsert(index int, elem any)

	doInsertAll(index int, list Collection) bool

	doClear() []any

	doMove(oldIndex, newIndex int) any

	doRemove(index int) any

	doRemoveRange(fromIndex, toIndex int) []any
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

func (list *AbstractEList) asAbstractEList() abstractEList {
	return list.interfaces.(abstractEList)
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
	l := list.asAbstractEList()
	if list.isUnique && l.Contains(elem) {
		return false
	}
	l.doAdd(elem)
	return true
}

func (list *AbstractEList) AddAll(collection Collection) bool {
	if list.isUnique {
		collection = list.getNonDuplicates(collection)
		if collection.Size() == 0 {
			return false
		}
	}
	list.asAbstractEList().doAddAll(collection)
	return true
}

func (list *AbstractEList) Insert(index int, elem any) bool {
	l := list.asAbstractEList()
	if size := l.Size(); index < 0 || index > size {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(size))
	}
	if list.isUnique && l.Contains(elem) {
		return false
	}
	l.doInsert(index, elem)
	return true
}

func (list *AbstractEList) InsertAll(index int, collection Collection) bool {
	l := list.asAbstractEList()
	if size := l.Size(); index < 0 || index > size {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(size))
	}
	if list.isUnique {
		collection = list.getNonDuplicates(collection)
		if collection.Size() == 0 {
			return false
		}
	}
	l.doInsertAll(index, collection)
	return true
}

func (list *AbstractEList) MoveObject(newIndex int, elem any) {
	l := list.asAbstractEList()
	oldIndex := l.IndexOf(elem)
	if oldIndex == -1 {
		panic("Object not found")
	}
	l.doMove(oldIndex, newIndex)
}

// Swap move an element from oldIndex to newIndex
func (list *AbstractEList) Move(oldIndex, newIndex int) any {
	l := list.asAbstractEList()
	if size := l.Size(); oldIndex < 0 || oldIndex >= size || newIndex < 0 || newIndex > size {
		panic("Index out of bounds: oldIndex=" + strconv.Itoa(oldIndex) + " newIndex=" + strconv.Itoa(newIndex) + " size=" + strconv.Itoa(size))
	}
	return l.doMove(oldIndex, newIndex)
}

// RemoveAt remove an element at a given position
func (list *AbstractEList) RemoveAt(index int) any {
	l := list.asAbstractEList()
	if size := l.Size(); index < 0 || index >= size {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(size))
	}
	return l.doRemove(index)
}

// Remove an element in an array
func (list *AbstractEList) Remove(elem any) bool {
	l := list.asAbstractEList()
	index := l.IndexOf(elem)
	if index == -1 {
		return false
	}
	l.doRemove(index)
	return true
}

func (list *AbstractEList) RemoveRange(fromIndex int, toIndex int) {
	l := list.asAbstractEList()
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
	l.doRemoveRange(fromIndex, toIndex)
}

func (list *AbstractEList) RemoveAll(collection Collection) bool {
	modified := false
	l := list.asAbstractEList()
	for i := l.Size() - 1; i >= 0; i-- {
		if collection.Contains(l.doGet(i)) {
			l.RemoveAt(i)
			modified = true
		}
	}
	return modified
}

// Get an element of the array
func (list *AbstractEList) Get(index int) any {
	l := list.asAbstractEList()
	if size := l.Size(); index < 0 || index >= size {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(size))
	}
	return l.doGet(index)
}

// Set an element of the array
func (list *AbstractEList) Set(index int, elem any) any {
	l := list.asAbstractEList()
	if size := l.Size(); index < 0 || index >= size {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(size))
	}
	if list.isUnique {
		currIndex := l.IndexOf(elem)
		if currIndex >= 0 && currIndex != index {
			panic("element already in list")
		}
	}
	return l.doSet(index, elem)
}

func (list *AbstractEList) Clear() {
	list.asAbstractEList().doClear()
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
