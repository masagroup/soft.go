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
	Collection

	// Insert inserts element at specified index
	// returns true if element was added
	Insert(index int, element any) bool

	// Insert inserts collection at specified index
	// returns true if element was added
	InsertAll(index int, element Collection) bool

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

	// IndexOf returns the index of element in the list if present, -1 otherwise
	IndexOf(any) int
}

// BinarySearch searches for target in a sorted list and returns the earliest
// position where target is found, or the position where target would appear
// in the sort order. The list must be sorted in increasing order, where "increasing"
// is defined by cmp. cmp should return 0 if the list element matches
// the target, a negative number if the list element precedes the target,
// or a positive number if the list element follows the target.
// cmp must implement the same ordering as the list, such that if
// cmp(a, t) < 0 and cmp(b, t) >= 0, then a must precede b in the list.
func BinarySearch(l EList, target any, cmp func(any, any) int) (int, bool) {
	n := l.Size()
	// Define cmp(l[-1], target) < 0 and cmp(l[n], target) >= 0 .
	// Invariant: cmp(l[i - 1], target) < 0, cmp(l[j], target) >= 0.
	i, j := 0, n
	for i < j {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h
		// i ≤ h < j
		if cmp(l.Get(h), target) < 0 {
			i = h + 1 // preserves cmp(l[i - 1], target) < 0
		} else {
			j = h // preserves cmp(l[j], target) >= 0
		}
	}
	// i == j, cmp(l[i-1], target) < 0, and cmp(l[j], target) (= cmp(l[i], target)) >= 0  =>  answer is i.
	return i, i < n && cmp(l.Get(i), target) == 0
}
