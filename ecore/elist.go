// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

// EList is the interface for dynamic containers
type EList interface {
	// Add adds element at the end of the list
	// returns true if element was added
	Add(element any) bool

	// Add adds collection at the end of the list
	// returns true if collection was added to the list
	AddAll(collection EList) bool

	// Insert inserts element at specified index
	// returns true if element was added
	Insert(index int, element any) bool

	// Insert inserts collection at specified index
	// returns true if element was added
	InsertAll(index int, element EList) bool

	// MoveObject moves eleemnt to specified index
	MoveObject(index int, element any)

	// Move moves element from oldIndex to newIndex
	Move(oldIndex int, newIndex int) any

	// Get returns element at specified index
	Get(index int) any

	// Set replaces the element at the specified index in this list with the
	// specified element (optional operation).
	// returns old element
	Set(index int, element any) any

	// RemoveAt removes element at specified index and returns old element
	RemoveAt(index int) any

	// RemoveRange removes elements in [fromIndex,toIndex) range
	RemoveRange(fromIndex int, toIndex int)

	// Remove removes element from the list
	// returns true if element was effectively removed
	Remove(element any) bool

	// RemoveAll removes collection from the list
	// returns true if collection was effectively removed
	RemoveAll(collection EList) bool

	// Size returns the size of the list
	Size() int

	// Clear clears the list
	Clear()

	// Empty returns true if list is empty , false otherwise
	Empty() bool

	// Contains returns true if list constians element, false otherwise
	Contains(element any) bool

	// IndexOf returns the index of element in the list if present, -1 otherwise
	IndexOf(any) int

	// Iterator returns the list iterator
	Iterator() EIterator

	// ToArray returns the list array
	ToArray() []any
}
