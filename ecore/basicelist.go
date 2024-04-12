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

// basicEList is an array of a dynamic size
type BasicEList struct {
	AbstractEList
	data []any
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

func (list *BasicEList) asEListCallbacks() listCallBacks {
	return list.interfaces.(listCallBacks)
}

// Add a new element to the array
func (list *BasicEList) doAdd(e any) {
	size := len(list.data)
	list.data = append(list.data, e)
	// events
	listCallbacks := list.asEListCallbacks()
	listCallbacks.DidAdd(size, e)
	listCallbacks.DidChange()
}

func (list *BasicEList) doAddAll(collection EList) bool {
	data := collection.ToArray()
	list.data = append(list.data, data...)
	// events
	listCallbacks := list.asEListCallbacks()
	for i, element := range data {
		listCallbacks.DidAdd(i, element)
	}
	listCallbacks.DidChange()
	return len(data) != 0
}

func (list *BasicEList) doInsert(index int, e any) {
	list.data = append(list.data, nil)
	copy(list.data[index+1:], list.data[index:])
	list.data[index] = e
	// events
	listCallbacks := list.asEListCallbacks()
	listCallbacks.DidAdd(index, e)
	listCallbacks.DidChange()
}

func (list *BasicEList) doInsertAll(index int, collection EList) bool {
	data := collection.ToArray()
	list.data = append(list.data[:index], append(data, list.data[index:]...)...)
	// events
	listCallbacks := list.asEListCallbacks()
	for i, element := range data {
		listCallbacks.DidAdd(i+index, element)
		listCallbacks.DidChange()
	}
	return len(data) != 0
}

func (list *BasicEList) doMove(oldIndex, newIndex int) any {
	object := list.data[oldIndex]
	if oldIndex != newIndex {
		if newIndex < oldIndex {
			copy(list.data[newIndex+1:], list.data[newIndex:oldIndex])
		} else {
			copy(list.data[oldIndex:], list.data[oldIndex+1:newIndex+1])
		}
		list.data[newIndex] = object

		// events
		listCallbacks := list.asEListCallbacks()
		listCallbacks.DidMove(newIndex, object, oldIndex)
		listCallbacks.DidChange()
	}
	return object
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
	listCallbacks := list.asEListCallbacks()
	listCallbacks.DidRemove(index, object)
	listCallbacks.DidChange()
	return object
}

func (list *BasicEList) doRemoveRange(fromIndex int, toIndex int) []any {
	// backup old objects
	objects := append([]any{}, list.data[fromIndex:toIndex]...)
	// remove range
	list.data = append(list.data[:fromIndex], list.data[toIndex:]...)
	// events
	listCallbacks := list.asEListCallbacks()
	for i := toIndex - 1; i >= fromIndex; i-- {
		listCallbacks.DidRemove(i, objects[i-fromIndex])
	}
	listCallbacks.DidChange()
	return objects
}

func (list *BasicEList) doGet(index int) any {
	return list.data[index]
}

func (list *BasicEList) doSet(index int, elem any) any {
	old := list.data[index]
	list.data[index] = elem
	// events
	listCallbacks := list.asEListCallbacks()
	listCallbacks.DidSet(index, elem, old)
	listCallbacks.DidChange()
	return old
}

func (list *BasicEList) doClear() []any {
	oldData := list.data
	list.data = make([]any, 0)

	// events
	listCallbacks := list.asEListCallbacks()
	listCallbacks.DidClear(oldData)
	return oldData
}

// Size count the number of element in the array
func (list *BasicEList) Size() int {
	return len(list.data)
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
