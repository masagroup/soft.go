// *****************************************************************************
//
// This file is part of a MASA library or program.
// Refer to the included end-user license agreement for restrictions.
//
// Copyright (c) 2020 MASA Group
//
// *****************************************************************************

package ecore

import "strconv"

type abstractEList interface {
	doGet(index int) interface{}

	doSet(index int, elem interface{}) interface{}

	doAdd(elem interface{})

	doAddAll(list EList) bool

	doInsert(index int, elem interface{})

	doInsertAll(index int, list EList) bool

	doClear() []interface{}

	doMove(oldIndex, newIndew int) interface{}

	doRemove(index int) interface{}

	didAdd(index int, elem interface{})

	didSet(index int, newElem interface{}, oldElem interface{})

	didRemove(index int, old interface{})

	didClear(oldObjects []interface{})

	didMove(newIndex int, movedObject interface{}, oldIndex int)

	didChange()
}

// basicEList is an array of a dynamic size
type basicEList struct {
	interfaces interface{}
	data       []interface{}
	isUnique   bool
}

// NewEmptyBasicEList return a new ArrayEList
func NewEmptyBasicEList() *basicEList {
	a := new(basicEList)
	a.interfaces = a
	a.data = []interface{}{}
	a.isUnique = false
	return a
}

// NewBasicEList return a new ArrayEList
func NewBasicEList(data []interface{}) *basicEList {
	a := new(basicEList)
	a.interfaces = a
	a.data = data
	a.isUnique = false
	return a
}

// NewUniqueBasicEList return a new ArrayEList with isUnique set as true
func NewUniqueBasicEList(data []interface{}) *basicEList {
	a := new(basicEList)
	a.interfaces = a
	a.data = data
	a.isUnique = true
	return a
}

// Remove all elements in list that already are in ref
func getNonDuplicates(list EList, ref EList) *basicEList {
	newList := NewBasicEList([]interface{}{})
	for it := list.Iterator(); it.HasNext(); {
		value := it.Next()
		if !newList.Contains(value) && !ref.Contains(value) {
			newList.Add(value)
		}
	}
	return newList
}

func (arr *basicEList) Add(elem interface{}) bool {
	if arr.isUnique && arr.Contains(elem) {
		return false
	}
	arr.interfaces.(abstractEList).doAdd(elem)
	return true
}

// Add a new element to the array
func (arr *basicEList) doAdd(e interface{}) {
	size := len(arr.data)
	arr.data = append(arr.data, e)
	// events
	interfaces := arr.interfaces.(abstractEList)
	interfaces.didAdd(size, e)
	interfaces.didChange()
}

// AddAll elements of an array in the current one
func (arr *basicEList) AddAll(collection EList) bool {
	if arr.isUnique {
		collection = getNonDuplicates(collection, arr)
		if collection.Size() == 0 {
			return false
		}
	}
	arr.interfaces.(abstractEList).doAddAll(collection)
	return true
}

func (arr *basicEList) doAddAll(collection EList) bool {
	data := collection.ToArray()
	arr.data = append(arr.data, data...)
	interfaces := arr.interfaces.(abstractEList)
	// events
	for i, element := range data {
		interfaces.didAdd(i, element)
		interfaces.didChange()
	}
	return len(data) != 0
}

// Insert an element in the array
func (arr *basicEList) Insert(index int, elem interface{}) bool {
	if index < 0 || index > arr.Size() {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(arr.Size()))
	}
	if arr.isUnique && arr.Contains(elem) {
		return false
	}
	arr.interfaces.(abstractEList).doInsert(index, elem)
	return true
}

func (arr *basicEList) doInsert(index int, e interface{}) {
	arr.data = append(arr.data, nil)
	copy(arr.data[index+1:], arr.data[index:])
	arr.data[index] = e
	// events
	interfaces := arr.interfaces.(abstractEList)
	interfaces.didAdd(index, e)
	interfaces.didChange()
}

// InsertAll element of an array at a given position
func (arr *basicEList) InsertAll(index int, collection EList) bool {
	if index < 0 || index > arr.Size() {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(arr.Size()))
	}
	if arr.isUnique {
		collection = getNonDuplicates(collection, arr)
		if collection.Size() == 0 {
			return false
		}
	}
	arr.interfaces.(abstractEList).doInsertAll(index, collection)
	return true
}

func (arr *basicEList) doInsertAll(index int, collection EList) bool {
	data := collection.ToArray()
	arr.data = append(arr.data[:index], append(data, arr.data[index:]...)...)
	// events
	interfaces := arr.interfaces.(abstractEList)
	for i, element := range data {
		interfaces.didAdd(i+index, element)
		interfaces.didChange()
	}
	return len(data) != 0
}

// Move an element to the given index
func (arr *basicEList) MoveObject(newIndex int, elem interface{}) {
	oldIndex := arr.IndexOf(elem)
	if oldIndex == -1 {
		panic("Object not found")
	}
	arr.interfaces.(abstractEList).doMove(oldIndex, newIndex)
}

// Swap move an element from oldIndex to newIndex
func (arr *basicEList) Move(oldIndex, newIndex int) interface{} {
	return arr.interfaces.(abstractEList).doMove(oldIndex, newIndex)
}

func (arr *basicEList) doMove(oldIndex, newIndex int) interface{} {
	if oldIndex < 0 || oldIndex >= arr.Size() ||
		newIndex < 0 || newIndex > arr.Size() {
		panic("Index out of bounds: oldIndex=" + strconv.Itoa(oldIndex) + " newIndex=" + strconv.Itoa(newIndex) + " size=" + strconv.Itoa(arr.Size()))
	}
	object := arr.data[oldIndex]
	copy(arr.data[oldIndex:], arr.data[oldIndex+1:])
	if newIndex > oldIndex {
		newIndex--
	}
	copy(arr.data[newIndex+1:], arr.data[newIndex:])
	arr.data[newIndex] = object
	// events
	interfaces := arr.interfaces.(abstractEList)
	interfaces.didMove(newIndex, object, oldIndex)
	interfaces.didChange()
	return object
}

// RemoveAt remove an element at a given position
func (arr *basicEList) RemoveAt(index int) interface{} {
	return arr.interfaces.(abstractEList).doRemove(index)
}

// Remove an element in an array
func (arr *basicEList) Remove(elem interface{}) bool {
	index := arr.IndexOf(elem)
	if index == -1 {
		return false
	}
	arr.interfaces.(abstractEList).doRemove(index)
	return true
}

func (arr *basicEList) doRemove(index int) interface{} {
	if index < 0 || index >= arr.Size() {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(arr.Size()))
	}
	object := arr.data[index]
	arr.data = append(arr.data[:index], arr.data[index+1:]...)
	// events
	interfaces := arr.interfaces.(abstractEList)
	interfaces.didRemove(index, object)
	interfaces.didChange()
	return object
}

// Get an element of the array
func (arr *basicEList) Get(index int) interface{} {
	if index < 0 || index >= arr.Size() {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(arr.Size()))
	}
	return arr.interfaces.(abstractEList).doGet(index)
}

func (arr *basicEList) doGet(index int) interface{} {
	return arr.data[index]
}

func (arr *basicEList) doSet(index int, elem interface{}) interface{} {
	old := arr.data[index]
	arr.data[index] = elem
	// events
	interfaces := arr.interfaces.(abstractEList)
	interfaces.didSet(index, elem, old)
	interfaces.didChange()
	return old
}

// Set an element of the array
func (arr *basicEList) Set(index int, elem interface{}) interface{} {
	if index < 0 || index >= arr.Size() {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(arr.Size()))
	}
	if arr.isUnique {
		currIndex := arr.IndexOf(elem)
		if currIndex >= 0 && currIndex != index {
			panic("element already in list")
		}
	}
	return arr.interfaces.(abstractEList).doSet(index, elem)
}

// Size count the number of element in the array
func (arr *basicEList) Size() int {
	return len(arr.data)
}

// Clear remove all elements of the array
func (arr *basicEList) Clear() {
	arr.interfaces.(abstractEList).doClear()
}

func (arr *basicEList) doClear() []interface{} {
	oldData := arr.data
	arr.data = make([]interface{}, 0)
	return oldData
}

// Empty return true if the array contains 0 element
func (arr *basicEList) Empty() bool {
	return arr.Size() == 0
}

// Contains return if an array contains or not an element
func (arr *basicEList) Contains(elem interface{}) bool {
	return arr.interfaces.(EList).IndexOf(elem) != -1
}

// IndexOf return the index on an element in an array, else return -1
func (arr *basicEList) IndexOf(elem interface{}) int {
	for i, value := range arr.data {
		if value == elem {
			return i
		}
	}
	return -1
}

// Iterator through the array
func (arr *basicEList) Iterator() EIterator {
	return &listIterator{list: arr}
}

func (arr *basicEList) ToArray() []interface{} {
	return arr.data
}

func (arr *basicEList) didAdd(index int, elem interface{}) {

}

func (arr *basicEList) didSet(index int, newElem interface{}, oldElem interface{}) {

}

func (arr *basicEList) didRemove(index int, old interface{}) {

}

func (arr *basicEList) didClear(oldObjects []interface{}) {

}

func (arr *basicEList) didMove(newIndex int, movedObject interface{}, oldIndex int) {

}

func (arr *basicEList) didChange() {

}
