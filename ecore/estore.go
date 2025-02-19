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
)

type EStore interface {
	AddRoot(object EObject)

	RemoveRoot(object EObject)

	Get(object EObject, feature EStructuralFeature, index int) any

	Set(object EObject, feature EStructuralFeature, index int, value any, oldValue bool) any

	IsSet(object EObject, feature EStructuralFeature) bool

	UnSet(object EObject, feature EStructuralFeature)

	IsEmpty(object EObject, feature EStructuralFeature) bool

	Size(object EObject, feature EStructuralFeature) int

	Contains(object EObject, feature EStructuralFeature, value any) bool

	IndexOf(object EObject, feature EStructuralFeature, value any) int

	LastIndexOf(object EObject, feature EStructuralFeature, value any) int

	Add(object EObject, feature EStructuralFeature, index int, value any)

	AddAll(object EObject, feature EStructuralFeature, index int, collection Collection)

	Remove(object EObject, feature EStructuralFeature, index int) any

	Move(object EObject, feature EStructuralFeature, sourceIndex int, targetIndex int) any

	Clear(object EObject, feature EStructuralFeature)

	GetContainer(object EObject) (EObject, EStructuralFeature)

	All(object EObject, feature EStructuralFeature) iter.Seq[any]

	ToArray(object EObject, feature EStructuralFeature) []any
}

type EStoreAsync interface {
	EStore
	TaskManager
	Close() error
}

type EStoreProvider interface {
	SetEStore(store EStore)

	GetEStore() EStore
}
