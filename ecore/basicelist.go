// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
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

// Remove all elements from list that are not in ref list
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

func (list *basicEList) SetInterfaces(interfaces interface{}) {
	list.interfaces = interfaces
}

func (list *basicEList) Add(elem interface{}) bool {
	if list.isUnique && list.Contains(elem) {
		return false
	}
	list.interfaces.(abstractEList).doAdd(elem)
	return true
}

// Add a new element to the array
func (list *basicEList) doAdd(e interface{}) {
	size := len(list.data)
	list.data = append(list.data, e)
	// events
	interfaces := list.interfaces.(abstractEList)
	interfaces.didAdd(size, e)
	interfaces.didChange()
}

// AddAll elements of an array in the current one
func (list *basicEList) AddAll(collection EList) bool {
	if list.isUnique {
		collection = getNonDuplicates(collection, list)
		if collection.Size() == 0 {
			return false
		}
	}
	list.interfaces.(abstractEList).doAddAll(collection)
	return true
}

func (list *basicEList) doAddAll(collection EList) bool {
	data := collection.ToArray()
	list.data = append(list.data, data...)
	interfaces := list.interfaces.(abstractEList)
	// events
	for i, element := range data {
		interfaces.didAdd(i, element)
		interfaces.didChange()
	}
	return len(data) != 0
}

// Insert an element in the array
func (list *basicEList) Insert(index int, elem interface{}) bool {
	if index < 0 || index > list.Size() {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(list.Size()))
	}
	if list.isUnique && list.Contains(elem) {
		return false
	}
	list.interfaces.(abstractEList).doInsert(index, elem)
	return true
}

func (list *basicEList) doInsert(index int, e interface{}) {
	list.data = append(list.data, nil)
	copy(list.data[index+1:], list.data[index:])
	list.data[index] = e
	// events
	interfaces := list.interfaces.(abstractEList)
	interfaces.didAdd(index, e)
	interfaces.didChange()
}

// InsertAll element of an array at a given position
func (list *basicEList) InsertAll(index int, collection EList) bool {
	if index < 0 || index > list.Size() {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(list.Size()))
	}
	if list.isUnique {
		collection = getNonDuplicates(collection, list)
		if collection.Size() == 0 {
			return false
		}
	}
	list.interfaces.(abstractEList).doInsertAll(index, collection)
	return true
}

func (list *basicEList) doInsertAll(index int, collection EList) bool {
	data := collection.ToArray()
	list.data = append(list.data[:index], append(data, list.data[index:]...)...)
	// events
	interfaces := list.interfaces.(abstractEList)
	for i, element := range data {
		interfaces.didAdd(i+index, element)
		interfaces.didChange()
	}
	return len(data) != 0
}

// Move an element to the given index
func (list *basicEList) MoveObject(newIndex int, elem interface{}) {
	oldIndex := list.interfaces.(EList).IndexOf(elem)
	if oldIndex == -1 {
		panic("Object not found")
	}
	list.interfaces.(abstractEList).doMove(oldIndex, newIndex)
}

// Swap move an element from oldIndex to newIndex
func (list *basicEList) Move(oldIndex, newIndex int) interface{} {
	return list.interfaces.(abstractEList).doMove(oldIndex, newIndex)
}

func (list *basicEList) doMove(oldIndex, newIndex int) interface{} {
	if oldIndex < 0 || oldIndex >= list.Size() ||
		newIndex < 0 || newIndex > list.Size() {
		panic("Index out of bounds: oldIndex=" + strconv.Itoa(oldIndex) + " newIndex=" + strconv.Itoa(newIndex) + " size=" + strconv.Itoa(list.Size()))
	}

	object := list.data[oldIndex]
	if oldIndex != newIndex {
		if newIndex < oldIndex {
			copy(list.data[newIndex+1:], list.data[newIndex:oldIndex])
		} else {
			copy(list.data[oldIndex:], list.data[oldIndex+1:newIndex+1])
		}
		list.data[newIndex] = object

		// events
		interfaces := list.interfaces.(abstractEList)
		interfaces.didMove(newIndex, object, oldIndex)
		interfaces.didChange()
	}
	return object
}

// RemoveAt remove an element at a given position
func (list *basicEList) RemoveAt(index int) interface{} {
	return list.interfaces.(abstractEList).doRemove(index)
}

// Remove an element in an array
func (list *basicEList) Remove(elem interface{}) bool {
	index := list.interfaces.(EList).IndexOf(elem)
	if index == -1 {
		return false
	}
	list.interfaces.(abstractEList).doRemove(index)
	return true
}

func (list *basicEList) doRemove(index int) interface{} {
	if index < 0 || index >= list.Size() {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(list.Size()))
	}
	// retrieve removed object
	object := list.data[index]

	// remove index
	copy(list.data[index:], list.data[index+1:])
	list.data[len(list.data)-1] = nil
	list.data = list.data[:len(list.data)-1]

	// events
	interfaces := list.interfaces.(abstractEList)
	interfaces.didRemove(index, object)
	interfaces.didChange()
	return object
}

func (list *basicEList) RemoveAll(collection EList) bool {
	modified := false
	for i := list.Size() - 1; i >= 0; {
		if collection.Contains(list.Get(i)) {
			list.RemoveAt(i)
			modified = true
		}
		i--
	}
	return modified
}

// Get an element of the array
func (list *basicEList) Get(index int) interface{} {
	if index < 0 || index >= list.Size() {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(list.Size()))
	}
	return list.interfaces.(abstractEList).doGet(index)
}

func (list *basicEList) doGet(index int) interface{} {
	return list.data[index]
}

func (list *basicEList) doSet(index int, elem interface{}) interface{} {
	old := list.data[index]
	list.data[index] = elem
	// events
	interfaces := list.interfaces.(abstractEList)
	interfaces.didSet(index, elem, old)
	interfaces.didChange()
	return old
}

// Set an element of the array
func (list *basicEList) Set(index int, elem interface{}) interface{} {
	if index < 0 || index >= list.Size() {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(list.Size()))
	}
	if list.isUnique {
		currIndex := list.IndexOf(elem)
		if currIndex >= 0 && currIndex != index {
			panic("element already in list")
		}
	}
	return list.interfaces.(abstractEList).doSet(index, elem)
}

// Size count the number of element in the array
func (list *basicEList) Size() int {
	return len(list.data)
}

// Clear remove all elements of the array
func (list *basicEList) Clear() {
	list.interfaces.(abstractEList).doClear()
}

func (list *basicEList) doClear() []interface{} {
	oldData := list.data
	list.data = make([]interface{}, 0)

	// events
	interfaces := list.interfaces.(abstractEList)
	interfaces.didClear(oldData)
	return oldData
}

// Empty return true if the array contains 0 element
func (list *basicEList) Empty() bool {
	return list.Size() == 0
}

// Contains return if an array contains or not an element
func (list *basicEList) Contains(elem interface{}) bool {
	return list.interfaces.(EList).IndexOf(elem) != -1
}

// IndexOf return the index on an element in an array, else return -1
func (list *basicEList) IndexOf(elem interface{}) int {
	for i, value := range list.data {
		if value == elem {
			return i
		}
	}
	return -1
}

// Iterator through the array
func (list *basicEList) Iterator() EIterator {
	return &listIterator{list: list}
}

func (list *basicEList) ToArray() []interface{} {
	return list.data
}

func (list *basicEList) didAdd(index int, elem interface{}) {

}

func (list *basicEList) didSet(index int, newElem interface{}, oldElem interface{}) {

}

func (list *basicEList) didRemove(index int, old interface{}) {

}

func (list *basicEList) didClear(oldObjects []interface{}) {

}

func (list *basicEList) didMove(newIndex int, movedObject interface{}, oldIndex int) {

}

func (list *basicEList) didChange() {

}
