// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

type EStore interface {
	Get(object EObject, feature EStructuralFeature, index int) any

	Set(object EObject, feature EStructuralFeature, index int, value any) any

	IsSet(object EObject, feature EStructuralFeature) bool

	UnSet(object EObject, feature EStructuralFeature)

	IsEmpty(object EObject, feature EStructuralFeature) bool

	Size(object EObject, feature EStructuralFeature) int

	Contains(object EObject, feature EStructuralFeature, value any) bool

	IndexOf(object EObject, feature EStructuralFeature, value any) int

	LastIndexOf(object EObject, feature EStructuralFeature, value any) int

	Add(object EObject, feature EStructuralFeature, index int, value any)

	Remove(object EObject, feature EStructuralFeature, index int) any

	Move(object EObject, feature EStructuralFeature, sourceIndex int, targetIndex int) any

	Clear(object EObject, feature EStructuralFeature)

	ToArray(object EObject, feature EStructuralFeature) []any
}
