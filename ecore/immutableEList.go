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
	"iter"
	"strconv"
)

type emptyImmutableEList struct {
}

func (l *emptyImmutableEList) Add(elem any) bool {
	panic("Immutable list can't be modified")
}

func (l *emptyImmutableEList) AddAll(list Collection) bool {
	panic("Immutable list can't be modified")
}

func (l *emptyImmutableEList) Insert(index int, elem any) bool {
	panic("Immutable list can't be modified")
}

func (l *emptyImmutableEList) InsertAll(index int, list Collection) bool {
	panic("Immutable list can't be modified")
}

func (l *emptyImmutableEList) MoveObject(newIndex int, elem any) {
	panic("Immutable list can't be modified")
}

func (l *emptyImmutableEList) Move(oldIndex, newIndex int) any {
	panic("Immutable list can't be modified")
}

// Get an element of the array
func (l *emptyImmutableEList) Get(index int) any {
	panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(l.Size()))
}

func (l *emptyImmutableEList) Set(index int, elem any) any {
	panic("Immutable list can't be modified")
}

func (l *emptyImmutableEList) RemoveAt(index int) any {
	panic("Immutable list can't be modified")
}

func (l *emptyImmutableEList) Remove(elem any) bool {
	panic("Immutable list can't be modified")
}

func (l *emptyImmutableEList) RemoveAll(collection Collection) bool {
	panic("Immutable list can't be modified")
}

func (l *emptyImmutableEList) RemoveRange(fromIndex, toIndex int) {
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
func (l *emptyImmutableEList) Contains(elem any) bool {
	return false
}

// IndexOf return the index on an element in an array, else return -1
func (l *emptyImmutableEList) IndexOf(elem any) int {
	return -1
}

// Iterator through the array
func (l *emptyImmutableEList) Iterator() EIterator {
	return &listIterator{list: l}
}

func (l *emptyImmutableEList) All() iter.Seq[any] {
	return func(yield func(any) bool) {}
}

// ToArray convert to array
func (l *emptyImmutableEList) ToArray() []any {
	return []any{}
}

func (l *emptyImmutableEList) GetUnResolvedList() EList {
	return l
}

func NewEmptyImmutableEList() *emptyImmutableEList {
	return &emptyImmutableEList{}
}

type immutableEList struct {
	emptyImmutableEList
	data []any
}

// NewImmutableEList return a new ImmutableEList
func NewImmutableEList(data []any) *immutableEList {
	return &immutableEList{data: data}
}

// Get an element of the array
func (l *immutableEList) Get(index int) any {
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
func (l *immutableEList) Contains(elem any) bool {
	return l.IndexOf(elem) != -1
}

// IndexOf return the index on an element in an array, else return -1
func (l *immutableEList) IndexOf(elem any) int {
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

func (l *immutableEList) All() iter.Seq[any] {
	return func(yield func(any) bool) {
		for _, value := range l.data {
			if !yield(value) {
				return
			}
		}
	}
}

// ToArray convert to array
func (l *immutableEList) ToArray() []any {
	return l.data
}

func (l *immutableEList) GetUnResolvedList() EList {
	return l
}
