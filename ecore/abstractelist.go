package ecore

import "strconv"

type abstractEList interface {
	EList

	doGet(index int) any

	doSet(index int, elem any) any

	doAdd(elem any)

	doAddAll(list EList) bool

	doInsert(index int, elem any)

	doInsertAll(index int, list EList) bool

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

func (list *AbstractEList) Add(elem any) bool {
	if list.isUnique && list.Contains(elem) {
		return false
	}
	list.asAbstractEList().doAdd(elem)
	return true
}

func (list *AbstractEList) AddAll(collection EList) bool {
	if list.isUnique {
		collection = getNonDuplicates(collection, list)
		if collection.Size() == 0 {
			return false
		}
	}
	list.asAbstractEList().doAddAll(collection)
	return true
}

func (list *AbstractEList) Insert(index int, elem any) bool {
	if size := list.asEList().Size(); index < 0 || index > size {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(size))
	}
	if list.isUnique && list.Contains(elem) {
		return false
	}
	list.asAbstractEList().doInsert(index, elem)
	return true
}

func (list *AbstractEList) InsertAll(index int, collection EList) bool {
	if size := list.asEList().Size(); index < 0 || index > size {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(size))
	}
	if list.isUnique {
		collection = getNonDuplicates(collection, list)
		if collection.Size() == 0 {
			return false
		}
	}
	list.asAbstractEList().doInsertAll(index, collection)
	return true
}

func (list *AbstractEList) MoveObject(newIndex int, elem any) {
	oldIndex := list.asEList().IndexOf(elem)
	if oldIndex == -1 {
		panic("Object not found")
	}
	list.asAbstractEList().doMove(oldIndex, newIndex)
}

// Swap move an element from oldIndex to newIndex
func (list *AbstractEList) Move(oldIndex, newIndex int) any {
	if size := list.asEList().Size(); oldIndex < 0 || oldIndex >= size || newIndex < 0 || newIndex > size {
		panic("Index out of bounds: oldIndex=" + strconv.Itoa(oldIndex) + " newIndex=" + strconv.Itoa(newIndex) + " size=" + strconv.Itoa(size))
	}
	return list.asAbstractEList().doMove(oldIndex, newIndex)
}

// RemoveAt remove an element at a given position
func (list *AbstractEList) RemoveAt(index int) any {
	return list.asAbstractEList().doRemove(index)
}

// Remove an element in an array
func (list *AbstractEList) Remove(elem any) bool {
	index := list.asEList().IndexOf(elem)
	if index == -1 {
		return false
	}
	list.asAbstractEList().doRemove(index)
	return true
}

func (list *AbstractEList) RemoveRange(fromIndex int, toIndex int) {
	size := list.asEList().Size()
	if fromIndex < 0 || fromIndex >= size {
		panic("Index out of bounds: fromIndex=" + strconv.Itoa(fromIndex) + " size=" + strconv.Itoa(size))
	}
	if toIndex < 0 || toIndex > list.asEList().Size() {
		panic("Index out of bounds: toIndex=" + strconv.Itoa(toIndex) + " size=" + strconv.Itoa(size))
	}
	if fromIndex > toIndex {
		panic("Indexes invalid: fromIndex=" + strconv.Itoa(fromIndex) + "must be less than toIndex=" + strconv.Itoa(toIndex))
	}
	list.asAbstractEList().doRemoveRange(fromIndex, toIndex)
}

func (list *AbstractEList) RemoveAll(collection EList) bool {
	modified := false
	l := list.asEList()
	for i := l.Size() - 1; i >= 0; i-- {
		if collection.Contains(l.Get(i)) {
			list.asEList().RemoveAt(i)
			modified = true
		}
	}
	return modified
}

// Get an element of the array
func (list *AbstractEList) Get(index int) any {
	if size := list.asEList().Size(); index < 0 || index >= size {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(size))
	}
	return list.asAbstractEList().doGet(index)
}

// Set an element of the array
func (list *AbstractEList) Set(index int, elem any) any {
	if size := list.asEList().Size(); index < 0 || index >= size {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(size))
	}
	if list.isUnique {
		currIndex := list.asEList().IndexOf(elem)
		if currIndex >= 0 && currIndex != index {
			panic("element already in list")
		}
	}
	return list.asAbstractEList().doSet(index, elem)
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

func (list *AbstractEList) ToArray() []any {
	panic("not implemented")
}

func getNonDuplicates(list EList, ref EList) *BasicEList {
	newList := NewBasicEList([]any{})
	for it := list.Iterator(); it.HasNext(); {
		value := it.Next()
		if !newList.Contains(value) && !ref.Contains(value) {
			newList.Add(value)
		}
	}
	return newList
}
