// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

type EMap interface {
	EList

	GetValue(key any) any

	Put(key any, value any)

	RemoveKey(key any) any

	ContainsValue(value any) bool

	ContainsKey(key any) bool

	ToMap() map[any]any
}
