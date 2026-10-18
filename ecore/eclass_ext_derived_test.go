// *****************************************************************************
// Copyright(c) 2026 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A build that raced an invalidation is never published: it is published with a
// compare-and-swap against the state it started from, which the invalidation replaced.
func TestEClassDerivedValue_BuildRacingInvalidation(t *testing.T) {
	var d derivedValue[int]
	start := d.state.Load()
	d.invalidate()
	assert.False(t, d.publish(start, 1))
	_, ok := d.load()
	assert.False(t, ok)

	// a reader whose build is invalidated builds again, and does not return the stale value
	builds := 0
	value := d.get(func() int {
		builds++
		if builds == 1 {
			d.invalidate()
		}
		return builds
	})
	assert.Equal(t, 2, value)
	value, ok = d.load()
	assert.True(t, ok)
	assert.Equal(t, 2, value)
	assert.Equal(t, uint64(2), d.state.Load().generation)

	// a build started before cannot replace a published one
	assert.False(t, d.publish(start, 3))
	value, _ = d.load()
	assert.Equal(t, 2, value)

	// invalidated: built again
	d.invalidate()
	assert.Equal(t, 4, d.get(func() int { return 4 }))
}

type derivedReadsModel struct {
	super, class                 *EClassExt
	superAttribute, attribute    *EAttributeExt
	superContainment, crossRef   *EReferenceExt
	superOperation, subOperation *EOperationExt
}

func newDerivedReadsModel() *derivedReadsModel {
	m := &derivedReadsModel{
		super:            newEClassExt(),
		class:            newEClassExt(),
		superAttribute:   newEAttributeExt(),
		attribute:        newEAttributeExt(),
		superContainment: newEReferenceExt(),
		crossRef:         newEReferenceExt(),
		superOperation:   newEOperationExt(),
		subOperation:     newEOperationExt(),
	}
	m.superAttribute.SetName("superAttribute")
	m.superAttribute.SetID(true)
	m.attribute.SetName("attribute")
	m.superContainment.SetName("superContainment")
	m.superContainment.SetContainment(true)
	m.crossRef.SetName("crossRef")
	m.superOperation.SetName("operation")
	m.subOperation.SetName("operation")
	m.super.GetEStructuralFeatures().AddAll(NewImmutableEList([]any{m.superAttribute, m.superContainment}))
	m.super.GetEOperations().Add(m.superOperation)
	m.class.GetESuperTypes().Add(m.super)
	m.class.GetEStructuralFeatures().AddAll(NewImmutableEList([]any{m.attribute, m.crossRef}))
	m.class.GetEOperations().Add(m.subOperation)
	return m
}

// checkDerivedReads reports whether all the derived values of the class of m match the model.
func (m *derivedReadsModel) checkDerivedReads(t *testing.T) bool {
	c := m.class
	checks := []struct {
		name     string
		actual   []any
		expected []any
	}{
		{"all features", c.GetEAllStructuralFeatures().ToArray(), []any{m.superAttribute, m.superContainment, m.attribute, m.crossRef}},
		{"all attributes", c.GetEAllAttributes().ToArray(), []any{m.superAttribute, m.attribute}},
		{"attributes", c.GetEAttributes().ToArray(), []any{m.attribute}},
		{"all references", c.GetEAllReferences().ToArray(), []any{m.superContainment, m.crossRef}},
		{"references", c.GetEReferences().ToArray(), []any{m.crossRef}},
		{"all containments", c.GetEAllContainments().ToArray(), []any{m.superContainment}},
		{"all cross references", c.GetEAllCrossReferences().ToArray(), []any{m.crossRef}},
		{"containment features", c.GetEContainmentFeatures().ToArray(), []any{m.superContainment}},
		{"cross reference features", c.GetECrossReferenceFeatures().ToArray(), []any{m.crossRef}},
		{"all operations", c.GetEAllOperations().ToArray(), []any{m.superOperation, m.subOperation}},
		{"all super types", c.GetEAllSuperTypes().ToArray(), []any{m.super}},
	}
	for _, check := range checks {
		if !assert.Equal(t, check.expected, check.actual, check.name) {
			return false
		}
	}
	return assert.Equal(t, EAttribute(m.superAttribute), c.GetEIDAttribute()) &&
		assert.Equal(t, EStructuralFeature(m.crossRef), c.GetEStructuralFeatureFromName("crossRef")) &&
		assert.Equal(t, EOperation(m.subOperation), c.GetOverride(m.superOperation)) &&
		assert.Equal(t, 3, c.GetFeatureID(m.crossRef))
}

// Readers keep reading the derived values of a class, without lock, while they are invalidated in
// the class and in its super class.
func TestEClassDerivedReadsConcurrentWithSetModified(t *testing.T) {
	m := newDerivedReadsModel()
	require.True(t, m.checkDerivedReads(t))

	featureIDs := []int{ECLASS__ESTRUCTURAL_FEATURES, ECLASS__EATTRIBUTES, ECLASS__EREFERENCES, ECLASS__EOPERATIONS, ECLASS__ESUPER_TYPES, ECLASS__NAME}
	var stop atomic.Bool
	var group sync.WaitGroup
	for range 8 {
		group.Go(func() {
			for !stop.Load() {
				if !m.checkDerivedReads(t) {
					return
				}
			}
		})
	}
	for i := range 2000 {
		featureID := featureIDs[i%len(featureIDs)]
		if i%2 == 0 {
			m.super.setModified(featureID)
		} else {
			m.class.setModified(featureID)
		}
	}
	stop.Store(true)
	group.Wait()
	assert.True(t, m.checkDerivedReads(t))
}

// Once a modification of a class returns, its readers see it: none gets a value built before.
func TestEClassDerivedReadsNotStaleAfterModification(t *testing.T) {
	m := newDerivedReadsModel()
	require.True(t, m.checkDerivedReads(t))

	// a feature added to the super class
	superAttribute2 := newEAttributeExt()
	superAttribute2.SetName("superAttribute2")
	m.super.GetEStructuralFeatures().Add(superAttribute2)
	assert.Equal(t, []any{m.superAttribute, superAttribute2, m.attribute}, m.class.GetEAllAttributes().ToArray())
	assert.Equal(t, EStructuralFeature(superAttribute2), m.class.GetEStructuralFeatureFromName("superAttribute2"))
	assert.Equal(t, 5, m.class.GetFeatureCount())
	assert.Equal(t, 4, m.class.GetFeatureID(m.crossRef))

	// a containment added to the class
	containment := newEReferenceExt()
	containment.SetName("containment")
	containment.SetContainment(true)
	m.class.GetEStructuralFeatures().Add(containment)
	assert.Equal(t, []any{m.superContainment, containment}, m.class.GetEAllContainments().ToArray())
	assert.Equal(t, []any{m.superContainment, containment}, m.class.GetEContainmentFeatures().ToArray())
	assert.Equal(t, []any{m.crossRef, containment}, m.class.GetEReferences().ToArray())

	// a feature removed
	m.class.GetEStructuralFeatures().Remove(m.crossRef)
	assert.Nil(t, m.class.GetEStructuralFeatureFromName("crossRef"))
	assert.Equal(t, []any{}, m.class.GetEAllCrossReferences().ToArray())
	assert.Equal(t, []any{}, m.class.GetECrossReferenceFeatures().ToArray())

	// an operation removed
	m.class.GetEOperations().Remove(m.subOperation)
	assert.Equal(t, []any{m.superOperation}, m.class.GetEAllOperations().ToArray())
	assert.Nil(t, m.class.GetOverride(m.superOperation))

	// a super type replaced
	otherSuper := newEClassExt()
	otherAttribute := newEAttributeExt()
	otherAttribute.SetName("otherAttribute")
	otherSuper.GetEStructuralFeatures().Add(otherAttribute)
	m.class.GetESuperTypes().Set(0, otherSuper)
	assert.Equal(t, []any{otherSuper}, m.class.GetEAllSuperTypes().ToArray())
	assert.Equal(t, []any{otherAttribute, m.attribute}, m.class.GetEAllAttributes().ToArray())
	assert.Nil(t, m.class.GetEIDAttribute())
	assert.Nil(t, m.class.GetEStructuralFeatureFromName("superAttribute"))
	assert.Equal(t, []any{}, m.class.GetEAllOperations().ToArray())
}

// Readers reading concurrently see each feature added to a class or its super class from the
// moment the addition returns. The writer waits for all the readers to see a feature before
// adding the next one: the lists owned by a class are not safe for concurrent changes.
func TestEClassDerivedReadsConcurrentNotStale(t *testing.T) {
	eClass := newEClassExt()
	eSuper := newEClassExt()
	eClass.GetESuperTypes().Add(eSuper)
	// the lists owned by the classes are created lazily, not safely for concurrent reads
	eClass.GetEStructuralFeatures()
	eSuper.GetEStructuralFeatures()
	const readers = 4
	const additions = 200
	var added atomic.Int64
	seen := make([]atomic.Int64, readers)
	for r := range seen {
		seen[r].Store(-1)
	}
	names := make([]string, additions+1)
	for i := range names {
		names[i] = fmt.Sprintf("attribute%d", i)
	}

	var group sync.WaitGroup
	for r := range readers {
		group.Go(func() {
			for {
				n := added.Load()
				if count := eClass.GetFeatureCount(); count < int(n) {
					t.Errorf("reader %d: %d features, %d added", r, count, n)
					return
				}
				if count := eClass.GetEAllAttributes().Size(); count < int(n) {
					t.Errorf("reader %d: %d attributes, %d added", r, count, n)
					return
				}
				if n > 0 && eClass.GetEStructuralFeatureFromName(names[n]) == nil {
					t.Errorf("reader %d: %s not found", r, names[n])
					return
				}
				seen[r].Store(n)
				if n == additions {
					return
				}
			}
		})
	}
	waitSeen := func(i int64) {
		for r := range readers {
			for seen[r].Load() != i && !t.Failed() {
				runtime.Gosched()
			}
		}
	}
	waitSeen(0)
	for i := 1; i <= additions && !t.Failed(); i++ {
		attribute := newEAttributeExt()
		attribute.SetName(names[i])
		if i%2 == 0 {
			eSuper.GetEStructuralFeatures().Add(attribute)
		} else {
			eClass.GetEStructuralFeatures().Add(attribute)
		}
		added.Store(int64(i))
		waitSeen(int64(i))
	}
	if t.Failed() {
		added.Store(additions)
	}
	group.Wait()
}
