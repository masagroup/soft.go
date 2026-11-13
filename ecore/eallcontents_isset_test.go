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
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// allContentsBoth returns the contents of o walked by EAllContents and by AllContents, which
// must be the same.
func allContentsBoth(t *testing.T, o EObject) []EObject {
	t.Helper()
	fromIterator := collectIterator(o.EAllContents())
	fromSeq := slices.Collect(AllContents(o))
	if len(fromSeq) == 0 {
		fromSeq = []EObject{}
	}
	assert.Equal(t, fromIterator, fromSeq, "EAllContents and AllContents differ")
	return fromIterator
}

// A single-valued containment is walked when it holds an object, whether it was set or unset
// before, and skipped otherwise.
func TestAllContents_SingleContainment(t *testing.T) {
	p := newEAllContentsBenchPackage()
	factory := p.nodeClass.GetEPackage().GetEFactoryInstance()
	newNode := func() EObject { return factory.Create(p.nodeClass) }

	root := newNode()
	assert.Empty(t, allContentsBoth(t, root), "not set")

	child := newNode()
	root.ESet(p.single, child)
	assert.Equal(t, []EObject{child}, allContentsBoth(t, root), "set")

	grandChild := newNode()
	child.ESet(p.single, grandChild)
	assert.Equal(t, []EObject{child, grandChild}, allContentsBoth(t, root), "contents of the single child")

	root.ESet(p.single, nil)
	assert.Empty(t, allContentsBoth(t, root), "set to nil")

	root.ESet(p.single, child)
	root.EUnset(p.single)
	assert.Empty(t, allContentsBoth(t, root), "unset")
}

// An unsettable single-valued containment set to nil is set, with no object: it is skipped.
func TestAllContents_UnsettableSingleContainmentSetToNil(t *testing.T) {
	factory := GetFactory()
	ePackage := factory.CreateEPackage()
	ePackage.SetName("unsettable")
	ePackage.SetNsURI("http://unsettable/eallcontents")
	nodeClass := factory.CreateEClass()
	nodeClass.SetName("Node")
	single := factory.CreateEReference()
	single.SetName("single")
	single.SetEType(nodeClass)
	single.SetContainment(true)
	single.SetUnsettable(true)
	nodeClass.GetEStructuralFeatures().Add(single)
	ePackage.GetEClassifiers().Add(nodeClass)
	ePackage.SetEFactoryInstance(factory.CreateEFactory())

	instances := ePackage.GetEFactoryInstance()
	root := instances.Create(nodeClass)
	child := instances.Create(nodeClass)
	root.ESet(single, child)
	assert.Equal(t, []EObject{child}, allContentsBoth(t, root))

	root.ESet(single, nil)
	assert.True(t, root.EIsSet(single), "an unsettable feature set to nil is set")
	assert.Empty(t, allContentsBoth(t, root))
}

// newMockSingleContentsObject returns a mock object whose class has one single-valued
// containment holding value, nil for none. The mock has no expectation of EIsSet: a walk asking
// it fails the test.
func newMockSingleContentsObject(t *testing.T, value EObject) *MockEObject {
	eReference := NewMockEReference(t)
	eReference.EXPECT().IsMany().Return(false).Maybe()
	eClass := NewMockEClass(t)
	eClass.EXPECT().GetEContainmentFeatures().Return(NewImmutableEList([]any{eReference})).Maybe()
	o := NewMockEObject(t)
	o.EXPECT().EClass().Return(eClass).Maybe()
	o.EXPECT().EGetResolve(eReference, mock.Anything).Return(value).Maybe()
	return o
}

// newMockLeafObject returns a mock object whose class has no containment.
func newMockLeafObject(t *testing.T) *MockEObject {
	eClass := NewMockEClass(t)
	eClass.EXPECT().GetEContainmentFeatures().Return(NewImmutableEList([]any{})).Maybe()
	o := NewMockEObject(t)
	o.EXPECT().EClass().Return(eClass).Maybe()
	return o
}

// A walk reads a single-valued containment without asking EIsSet.
func TestAllContents_SingleContainmentNotAskedIsSet(t *testing.T) {
	leaf := newMockLeafObject(t)
	o := newMockSingleContentsObject(t, leaf)
	assert.Equal(t, []EObject{leaf}, collectIterator(newEAllContentsIterator(newEAllContentsFrame(o, nil), nil)))
	assert.Equal(t, []EObject{leaf}, slices.Collect(AllContents(o)))

	empty := newMockSingleContentsObject(t, nil)
	assert.Empty(t, collectIterator(newEAllContentsIterator(newEAllContentsFrame(empty, nil), nil)))
	assert.Empty(t, slices.Collect(AllContents(empty)))
}

// A walk asks EIsSet of a many-valued containment, and does not read it when it is not set: its
// list is not made.
func TestAllContents_ManyContainmentNotSetNotRead(t *testing.T) {
	eReference := NewMockEReference(t)
	eReference.EXPECT().IsMany().Return(true).Maybe()
	eClass := NewMockEClass(t)
	eClass.EXPECT().GetEContainmentFeatures().Return(NewImmutableEList([]any{eReference})).Maybe()
	o := NewMockEObject(t)
	o.EXPECT().EClass().Return(eClass).Maybe()
	o.EXPECT().EIsSet(eReference).Return(false).Times(2)
	// no expectation of EGetResolve: reading the feature fails the test
	assert.Empty(t, collectIterator(newEAllContentsIterator(newEAllContentsFrame(o, nil), nil)))
	assert.Empty(t, slices.Collect(AllContents(o)))
}
