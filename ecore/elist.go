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
	Add(interface{}) bool

	AddAll(EList) bool

	Insert(int, interface{}) bool

	InsertAll(int, EList) bool

	MoveObject(int, interface{})

	Move(oldIndex int, newIndex int) interface{}

	Get(int) interface{}

	Set(int, interface{}) interface{}

	RemoveAt(int) interface{}

	Remove(interface{}) bool

	RemoveAll(EList) bool

	Size() int

	Clear()

	Empty() bool

	Contains(interface{}) bool

	IndexOf(interface{}) int

	Iterator() EIterator

	ToArray() []interface{}
}
