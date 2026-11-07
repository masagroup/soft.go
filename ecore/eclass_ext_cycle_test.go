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
	"testing"

	"github.com/stretchr/testify/assert"
)

// newCycleClass returns a class with one attribute and one operation, both named after the class
func newCycleClass(name string) (*EClassExt, EAttribute, EOperation) {
	eClass := newEClassExt()
	eClass.SetName(name)
	eAttribute := newEAttributeExt()
	eAttribute.SetName(name + "Attribute")
	eClass.GetEStructuralFeatures().Add(eAttribute)
	eOperation := newEOperationExt()
	eOperation.SetName(name + "Operation")
	eClass.GetEOperations().Add(eOperation)
	return eClass, eAttribute, eOperation
}

// readDerived reads every value derived from the super types of eClass
func readDerived(eClass EClass) {
	eClass.GetEAllSuperTypes()
	eClass.GetEAllStructuralFeatures()
	eClass.GetEAllAttributes()
	eClass.GetEAllReferences()
	eClass.GetEAllContainments()
	eClass.GetEAllCrossReferences()
	eClass.GetEAllOperations()
	eClass.GetEContainmentFeatures()
	eClass.GetECrossReferenceFeatures()
	eClass.GetEIDAttribute()
	eClass.GetFeatureCount()
	eClass.GetOperationCount()
	eClass.GetEStructuralFeatureFromName("")
	eClass.IsSuperTypeOf(eClass)
}

// A class among its own super types is invalid, but neither its modification nor the reads of its
// derived values recurse without end: the class inherits from itself nothing more than it has.
func TestEClassSuperTypesCycle_Self(t *testing.T) {
	eClass, eAttribute, eOperation := newCycleClass("a")
	readDerived(eClass)

	eClass.GetESuperTypes().Add(eClass)
	readDerived(eClass)
	assert.Equal(t, []any{eClass}, eClass.GetEAllSuperTypes().ToArray())
	assert.Equal(t, []any{eAttribute}, eClass.GetEAllStructuralFeatures().ToArray())
	assert.Equal(t, []any{eAttribute}, eClass.GetEAllAttributes().ToArray())
	assert.Equal(t, []any{eOperation}, eClass.GetEAllOperations().ToArray())
	assert.Equal(t, 0, eClass.GetFeatureID(eAttribute))

	// modified in the cycle
	eAttribute2 := newEAttributeExt()
	eAttribute2.SetName("aAttribute2")
	eClass.GetEStructuralFeatures().Add(eAttribute2)
	eAttribute2.SetID(true)
	assert.Equal(t, []any{eAttribute, eAttribute2}, eClass.GetEAllAttributes().ToArray())
	assert.Equal(t, EAttribute(eAttribute2), eClass.GetEIDAttribute())

	// the cycle removed
	eClass.GetESuperTypes().Remove(eClass)
	assert.True(t, eClass.GetEAllSuperTypes().Empty())
	assert.Equal(t, []any{eAttribute, eAttribute2}, eClass.GetEAllStructuralFeatures().ToArray())
}

// a -> b -> a, and c -> a: every class of the cycle or above it is read and modified in any order
func TestEClassSuperTypesCycle_TwoClasses(t *testing.T) {
	for name, order := range map[string][]int{
		"a first": {0, 1, 2},
		"b first": {1, 0, 2},
		"c first": {2, 1, 0},
	} {
		t.Run(name, func(t *testing.T) {
			a, aAttribute, _ := newCycleClass("a")
			b, bAttribute, _ := newCycleClass("b")
			c, cAttribute, _ := newCycleClass("c")
			classes := []*EClassExt{a, b, c}
			c.GetESuperTypes().Add(a)
			a.GetESuperTypes().Add(b)
			for _, i := range order {
				readDerived(classes[i])
			}

			// the edge closing the cycle, added once the values are built
			b.GetESuperTypes().Add(a)
			for _, i := range order {
				readDerived(classes[i])
			}
			assert.True(t, a.GetEAllSuperTypes().Contains(b))
			assert.True(t, b.GetEAllSuperTypes().Contains(a))
			assert.True(t, c.GetEAllSuperTypes().Contains(a))
			assert.True(t, c.GetEAllSuperTypes().Contains(b))
			assert.True(t, c.GetEAllStructuralFeatures().Contains(cAttribute))

			// modified in the cycle: the sub classes are invalidated
			aAttribute2 := newEAttributeExt()
			aAttribute2.SetName("aAttribute2")
			a.GetEStructuralFeatures().Add(aAttribute2)
			assert.True(t, c.GetEAllAttributes().Contains(aAttribute2))
			assert.Equal(t, EStructuralFeature(aAttribute2), c.GetEStructuralFeatureFromName("aAttribute2"))

			// the cycle removed: the values are the ones of an acyclic hierarchy again
			b.GetESuperTypes().Remove(a)
			assert.Equal(t, []any{b, a}, c.GetEAllSuperTypes().ToArray())
			assert.Equal(t, []any{bAttribute, aAttribute, aAttribute2, cAttribute}, c.GetEAllStructuralFeatures().ToArray())
			assert.Equal(t, 3, c.GetFeatureID(cAttribute))
			assert.True(t, b.GetEAllSuperTypes().Empty())
		})
	}
}

// the cycle closed by the first add, before any value is built
func TestEClassSuperTypesCycle_BeforeReads(t *testing.T) {
	a, _, _ := newCycleClass("a")
	b, _, _ := newCycleClass("b")
	a.GetESuperTypes().Add(b)
	b.GetESuperTypes().Add(a)
	readDerived(b)
	readDerived(a)
	assert.True(t, a.IsSuperTypeOf(b))
	assert.True(t, b.IsSuperTypeOf(a))
}

// Each class of a lattice inherits from the two classes before it: a modification of the first one
// reaches the last one through an exponential number of paths, and invalidates each class once.
func TestEClassSuperTypesLattice(t *testing.T) {
	classes := []*EClassExt{}
	for i := range 64 {
		eClass := newEClassExt()
		for _, j := range []int{i - 1, i - 2} {
			if j >= 0 {
				eClass.GetESuperTypes().Add(classes[j])
			}
		}
		classes = append(classes, eClass)
	}
	last := classes[len(classes)-1]
	assert.Equal(t, 63, last.GetEAllSuperTypes().Size())

	eAttribute := newEAttributeExt()
	eAttribute.SetName("attribute")
	classes[0].GetEStructuralFeatures().Add(eAttribute)
	assert.Equal(t, []any{eAttribute}, last.GetEAllAttributes().ToArray())
}
