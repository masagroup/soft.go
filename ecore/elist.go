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
	Add(any) bool

	AddAll(EList) bool

	Insert(int, any) bool

	InsertAll(int, EList) bool

	MoveObject(int, any)

	Move(oldIndex int, newIndex int) any

	Get(int) any

	Set(int, any) any

	RemoveAt(int) any

	Remove(any) bool

	RemoveAll(EList) bool

	Size() int

	Clear()

	Empty() bool

	Contains(any) bool

	IndexOf(any) int

	Iterator() EIterator

	ToArray() []any
}
