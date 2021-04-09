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

type emptyImmutableEList struct {
}

func (l *emptyImmutableEList) Add(elem interface{}) bool {
	panic("Immutable list can't be modified")
}

func (l *emptyImmutableEList) AddAll(list EList) bool {
	panic("Immutable list can't be modified")
}

func (l *emptyImmutableEList) Insert(index int, elem interface{}) bool {
	panic("Immutable list can't be modified")
}

func (l *emptyImmutableEList) InsertAll(index int, list EList) bool {
	panic("Immutable list can't be modified")
}

func (l *emptyImmutableEList) MoveObject(newIndex int, elem interface{}) {
	panic("Immutable list can't be modified")
}

func (l *emptyImmutableEList) Move(oldIndex, newIndex int) interface{} {
	panic("Immutable list can't be modified")
}

// Get an element of the array
func (l *emptyImmutableEList) Get(index int) interface{} {
	panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(l.Size()))
}

func (l *emptyImmutableEList) Set(index int, elem interface{}) interface{} {
	panic("Immutable list can't be modified")
}

func (l *emptyImmutableEList) RemoveAt(index int) interface{} {
	panic("Immutable list can't be modified")
}

func (l *emptyImmutableEList) Remove(elem interface{}) bool {
	panic("Immutable list can't be modified")
}

func (l *emptyImmutableEList) RemoveAll(collection EList) bool {
	panic("Immutable list can't be modified")
}

// Size count the number of element in the array
func (l *emptyImmutableEList) Size() int {
	return 0
}

func (l *emptyImmutableEList) Clear() {
	panic("Immutable list can't be modified")
}

// Empty return true if the array contains 0 element
func (l *emptyImmutableEList) Empty() bool {
	return true
}

// Contains return if an array contains or not an element
func (l *emptyImmutableEList) Contains(elem interface{}) bool {
	return false
}

// IndexOf return the index on an element in an array, else return -1
func (l *emptyImmutableEList) IndexOf(elem interface{}) int {
	return -1
}

// Iterator through the array
func (l *emptyImmutableEList) Iterator() EIterator {
	return &listIterator{list: l}
}

// ToArray convert to array
func (l *emptyImmutableEList) ToArray() []interface{} {
	return []interface{}{}
}

func (l *emptyImmutableEList) GetUnResolvedList() EList {
	return l
}

func NewEmptyImmutableEList() *emptyImmutableEList {
	return &emptyImmutableEList{}
}

type immutableEList struct {
	emptyImmutableEList
	data []interface{}
}

// NewImmutableEList return a new ImmutableEList
func NewImmutableEList(data []interface{}) *immutableEList {
	return &immutableEList{data: data}
}

// Get an element of the array
func (l *immutableEList) Get(index int) interface{} {
	if index < 0 || index >= l.Size() {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(l.Size()))
	}
	return l.data[index]
}

// Size count the number of element in the array
func (l *immutableEList) Size() int {
	return len(l.data)
}

// Empty return true if the array contains 0 element
func (l *immutableEList) Empty() bool {
	return l.Size() == 0
}

// Contains return if an array contains or not an element
func (l *immutableEList) Contains(elem interface{}) bool {
	return l.IndexOf(elem) != -1
}

// IndexOf return the index on an element in an array, else return -1
func (l *immutableEList) IndexOf(elem interface{}) int {
	for i, value := range l.data {
		if value == elem {
			return i
		}
	}
	return -1
}

// Iterator through the array
func (l *immutableEList) Iterator() EIterator {
	return &listIterator{list: l}
}

// ToArray convert to array
func (l *immutableEList) ToArray() []interface{} {
	return l.data
}

func (l *immutableEList) GetUnResolvedList() EList {
	return l
}
