// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import (
	"strconv"
)

type listCallBacks interface {
	DidAdd(index int, elem any)

	DidSet(index int, newElem any, oldElem any)

	DidRemove(index int, old any)

	DidClear(oldObjects []any)

	DidMove(newIndex int, movedObject any, oldIndex int)

	DidChange()
}

type abstractEList interface {
	listCallBacks

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

// basicEList is an array of a dynamic size
type BasicEList struct {
	interfaces any
	data       []any
	isUnique   bool
}

// NewEmptyBasicEList return a new ArrayEList
func NewEmptyBasicEList() *BasicEList {
	a := new(BasicEList)
	a.interfaces = a
	a.data = []any{}
	a.isUnique = false
	return a
}

// NewBasicEList return a new ArrayEList
func NewBasicEList(data []any) *BasicEList {
	a := new(BasicEList)
	a.interfaces = a
	a.data = data
	a.isUnique = false
	return a
}

// NewUniqueBasicEList return a new ArrayEList with isUnique set as true
func NewUniqueBasicEList(data []any) *BasicEList {
	a := new(BasicEList)
	a.interfaces = a
	a.data = data
	a.isUnique = true
	return a
}

// Remove all elements from list that are not in ref list
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

func (list *BasicEList) SetInterfaces(interfaces any) {
	list.interfaces = interfaces
}

func (list *BasicEList) asEList() EList {
	return list.interfaces.(EList)
}

func (list *BasicEList) asAbstractEList() abstractEList {
	return list.interfaces.(abstractEList)
}

func (list *BasicEList) Add(elem any) bool {
	if list.isUnique && list.Contains(elem) {
		return false
	}
	list.asAbstractEList().doAdd(elem)
	return true
}

// Add a new element to the array
func (list *BasicEList) doAdd(e any) {
	size := len(list.data)
	list.data = append(list.data, e)
	// events
	abstractEList := list.asAbstractEList()
	abstractEList.DidAdd(size, e)
	abstractEList.DidChange()
}

// AddAll elements of an array in the current one
func (list *BasicEList) AddAll(collection EList) bool {
	if list.isUnique {
		collection = getNonDuplicates(collection, list)
		if collection.Size() == 0 {
			return false
		}
	}
	list.asAbstractEList().doAddAll(collection)
	return true
}

func (list *BasicEList) doAddAll(collection EList) bool {
	data := collection.ToArray()
	list.data = append(list.data, data...)
	abstractEList := list.asAbstractEList()
	// events
	for i, element := range data {
		abstractEList.DidAdd(i, element)
	}
	abstractEList.DidChange()
	return len(data) != 0
}

// Insert an element in the array
func (list *BasicEList) Insert(index int, elem any) bool {
	if index < 0 || index > list.Size() {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(list.Size()))
	}
	if list.isUnique && list.Contains(elem) {
		return false
	}
	list.asAbstractEList().doInsert(index, elem)
	return true
}

func (list *BasicEList) doInsert(index int, e any) {
	list.data = append(list.data, nil)
	copy(list.data[index+1:], list.data[index:])
	list.data[index] = e
	// events
	abstractEList := list.asAbstractEList()
	abstractEList.DidAdd(index, e)
	abstractEList.DidChange()
}

// InsertAll element of an array at a given position
func (list *BasicEList) InsertAll(index int, collection EList) bool {
	if index < 0 || index > list.Size() {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(list.Size()))
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

func (list *BasicEList) doInsertAll(index int, collection EList) bool {
	data := collection.ToArray()
	list.data = append(list.data[:index], append(data, list.data[index:]...)...)
	// events
	abstractEList := list.asAbstractEList()
	for i, element := range data {
		abstractEList.DidAdd(i+index, element)
		abstractEList.DidChange()
	}
	return len(data) != 0
}

// Move an element to the given index
func (list *BasicEList) MoveObject(newIndex int, elem any) {
	oldIndex := list.asEList().IndexOf(elem)
	if oldIndex == -1 {
		panic("Object not found")
	}
	list.asAbstractEList().doMove(oldIndex, newIndex)
}

// Swap move an element from oldIndex to newIndex
func (list *BasicEList) Move(oldIndex, newIndex int) any {
	return list.asAbstractEList().doMove(oldIndex, newIndex)
}

func (list *BasicEList) doMove(oldIndex, newIndex int) any {
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
		abstractEList := list.asAbstractEList()
		abstractEList.DidMove(newIndex, object, oldIndex)
		abstractEList.DidChange()
	}
	return object
}

// RemoveAt remove an element at a given position
func (list *BasicEList) RemoveAt(index int) any {
	return list.asAbstractEList().doRemove(index)
}

// Remove an element in an array
func (list *BasicEList) Remove(elem any) bool {
	index := list.asEList().IndexOf(elem)
	if index == -1 {
		return false
	}
	list.asAbstractEList().doRemove(index)
	return true
}

func (list *BasicEList) doRemove(index int) any {
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
	abstractEList := list.asAbstractEList()
	abstractEList.DidRemove(index, object)
	abstractEList.DidChange()
	return object
}

func (list *BasicEList) RemoveRange(fromIndex int, toIndex int) {
	if fromIndex < 0 || fromIndex >= list.Size() {
		panic("Index out of bounds: fromIndex=" + strconv.Itoa(fromIndex) + " size=" + strconv.Itoa(list.Size()))
	}
	if toIndex < 0 || toIndex > list.Size() {
		panic("Index out of bounds: toIndex=" + strconv.Itoa(toIndex) + " size=" + strconv.Itoa(list.Size()))
	}
	if fromIndex > toIndex {
		panic("Indexes invalid: fromIndex=" + strconv.Itoa(fromIndex) + "must be less than toIndex=" + strconv.Itoa(toIndex))
	}
	list.asAbstractEList().doRemoveRange(fromIndex, toIndex)
}

func (list *BasicEList) doRemoveRange(fromIndex int, toIndex int) []any {
	// backup old objects
	objects := append([]any{}, list.data[fromIndex:toIndex]...)
	// remove range
	list.data = append(list.data[:fromIndex], list.data[toIndex:]...)
	// events
	abstractEList := list.asAbstractEList()
	for i := toIndex - 1; i >= fromIndex; i-- {
		abstractEList.DidRemove(i, objects[i-fromIndex])
	}
	abstractEList.DidChange()
	return objects
}

func (list *BasicEList) RemoveAll(collection EList) bool {
	modified := false
	for i := list.Size() - 1; i >= 0; i-- {
		if collection.Contains(list.Get(i)) {
			list.RemoveAt(i)
			modified = true
		}
	}
	return modified
}

// Get an element of the array
func (list *BasicEList) Get(index int) any {
	if index < 0 || index >= list.Size() {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(list.Size()))
	}
	return list.asAbstractEList().doGet(index)
}

func (list *BasicEList) doGet(index int) any {
	return list.data[index]
}

func (list *BasicEList) doSet(index int, elem any) any {
	old := list.data[index]
	list.data[index] = elem
	// events
	abstractEList := list.asAbstractEList()
	abstractEList.DidSet(index, elem, old)
	abstractEList.DidChange()
	return old
}

// Set an element of the array
func (list *BasicEList) Set(index int, elem any) any {
	if index < 0 || index >= list.Size() {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(list.Size()))
	}
	if list.isUnique {
		currIndex := list.IndexOf(elem)
		if currIndex >= 0 && currIndex != index {
			panic("element already in list")
		}
	}
	return list.asAbstractEList().doSet(index, elem)
}

// Size count the number of element in the array
func (list *BasicEList) Size() int {
	return len(list.data)
}

// Clear remove all elements of the array
func (list *BasicEList) Clear() {
	list.asAbstractEList().doClear()
}

func (list *BasicEList) doClear() []any {
	oldData := list.data
	list.data = make([]any, 0)

	// events
	abstractEList := list.asAbstractEList()
	abstractEList.DidClear(oldData)
	return oldData
}

// Empty return true if the array contains 0 element
func (list *BasicEList) Empty() bool {
	return list.Size() == 0
}

// Contains return if an array contains or not an element
func (list *BasicEList) Contains(elem any) bool {
	return list.asEList().IndexOf(elem) != -1
}

// IndexOf return the index on an element in an array, else return -1
func (list *BasicEList) IndexOf(elem any) int {
	for i, value := range list.data {
		if value == elem {
			return i
		}
	}
	return -1
}

// Iterator through the array
func (list *BasicEList) Iterator() EIterator {
	return &listIterator{list: list}
}

func (list *BasicEList) ToArray() []any {
	return list.data
}

func (list *BasicEList) DidAdd(index int, elem any) {

}

func (list *BasicEList) DidSet(index int, newElem any, oldElem any) {

}

func (list *BasicEList) DidRemove(index int, old any) {

}

func (list *BasicEList) DidClear(oldObjects []any) {

}

func (list *BasicEList) DidMove(newIndex int, movedObject any, oldIndex int) {

}

func (list *BasicEList) DidChange() {

}
