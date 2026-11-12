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
	"github.com/stretchr/testify/require"
)

func collectIterator(it EIterator) []EObject {
	objects := []EObject{}
	for it.HasNext() {
		objects = append(objects, it.Next().(EObject))
	}
	return objects
}

func TestEAllContents_TreeOrder(t *testing.T) {
	p := newEAllContentsBenchPackage()
	root, count := p.newTree(3, 3, 2)
	expected := collectIterator(newTreeEAllContentsIterator(root))
	require.Len(t, expected, count)
	assert.Equal(t, expected, collectIterator(root.EAllContents()))
	assert.Equal(t, expected, slices.Collect(AllContents(root)))
}

func TestEAllContents_Empty(t *testing.T) {
	p := newEAllContentsBenchPackage()
	root, _ := p.newTree(0, 0, 0)
	it := root.EAllContents()
	assert.False(t, it.HasNext())
	assert.Panics(t, func() { it.Next() })
	assert.Empty(t, slices.Collect(AllContents(root)))
}

func TestEAllContents_FromEContents(t *testing.T) {
	// an object not walked from its features: its contents are read from EContents
	child := NewMockEObject(t)
	child.EXPECT().EContents().Return(NewImmutableEList(nil))
	o := NewMockEObject(t)
	o.EXPECT().EContents().Return(NewImmutableEList([]any{child}))
	assert.Equal(t, []EObject{child}, collectIterator(newEAllContentsIterator(o)))
	assert.Equal(t, []EObject{child}, slices.Collect(AllContents(o)))
}

func TestAllContents_Break(t *testing.T) {
	p := newEAllContentsBenchPackage()
	root, _ := p.newTree(3, 3, 2)
	expected := collectIterator(root.EAllContents())[:5]
	objects := []EObject{}
	for o := range AllContents(root) {
		objects = append(objects, o)
		if len(objects) == 5 {
			break
		}
	}
	assert.Equal(t, expected, objects)
}

// pruneExpected returns the contents of root, without the contents of the objects pruned.
func pruneExpected(root EObject, pruned func(EObject) bool) []EObject {
	objects := []EObject{}
	for it := newTreeEAllContentsIterator(root).(*treeIterator); it.HasNext(); {
		o := it.Next().(EObject)
		objects = append(objects, o)
		if pruned(o) {
			it.Prune()
		}
	}
	return objects
}

func TestEAllContents_Prune(t *testing.T) {
	p := newEAllContentsBenchPackage()
	root, _ := p.newTree(3, 3, 2)
	// prune the nodes of the many containment, leaves included: a leaf has no contents
	pruned := func(o EObject) bool { return o.EContainingFeature() == p.children || o.EClass() == p.leafClass }
	expected := pruneExpected(root, pruned)
	require.Less(t, len(expected), len(collectIterator(root.EAllContents())))

	t.Run("after Next", func(t *testing.T) {
		objects := []EObject{}
		for it := root.EAllContents().(*eAllContentsIterator); it.HasNext(); {
			o := it.Next().(EObject)
			objects = append(objects, o)
			if pruned(o) {
				it.Prune()
			}
		}
		assert.Equal(t, expected, objects)
	})
	t.Run("after HasNext", func(t *testing.T) {
		objects := []EObject{}
		for it := root.EAllContents().(*eAllContentsIterator); it.HasNext(); {
			o := it.Next().(EObject)
			objects = append(objects, o)
			if pruned(o) {
				it.HasNext()
				it.Prune()
				it.Prune() // once only
			}
		}
		assert.Equal(t, expected, objects)
	})
}

// newAllContentsResource returns a resource of two trees and the objects GetAllContents iterates.
func newAllContentsResource(p *eAllContentsBenchPackage, depth, fanOut, leaves int) (EResource, []EObject) {
	resource := NewEResourceImpl()
	for range 2 {
		root, _ := p.newTree(depth, fanOut, leaves)
		resource.GetContents().Add(root)
	}
	return resource, collectIterator(resource.GetAllContents())
}

func TestAllContents(t *testing.T) {
	p := newEAllContentsBenchPackage()
	resource, expected := newAllContentsResource(p, 3, 3, 2)
	assert.Equal(t, expected, slices.Collect(AllContents(resource)))
	assert.Equal(t, expected, slices.Collect(allContentsResolve(resource, false)))
	assert.Empty(t, slices.Collect(AllContents(NewEResourceImpl())))

	objects := []EObject{}
	for o := range AllContents(resource) {
		objects = append(objects, o)
		if len(objects) == 5 {
			break
		}
	}
	assert.Equal(t, expected[:5], objects)
}

func TestEResourceImpl_AttachedDetached_Contents(t *testing.T) {
	p := newEAllContentsBenchPackage()
	root, count := p.newTree(2, 3, 1)
	resource := NewEResourceImpl()
	idManager := NewIncrementalIDManager()
	resource.SetObjectIDManager(idManager)
	resource.GetContents().Add(root)
	for o := range AllContents(root) {
		assert.NotEmpty(t, idManager.GetID(o), "%v", o)
	}
	assert.Equal(t, count+1, len(slices.Collect(AllContents(resource))))
	resource.GetContents().Remove(root)
	for o := range AllContents(root) {
		assert.Nil(t, idManager.GetID(o))
	}
}

func TestAllContents_ResourceSet(t *testing.T) {
	p := newEAllContentsBenchPackage()
	resourceSet := NewEResourceSetImpl()
	expected := []EObject{}
	for range 2 {
		resource, objects := newAllContentsResource(p, 2, 2, 1)
		resourceSet.GetResources().Add(resource)
		expected = append(expected, objects...)
	}
	assert.Equal(t, expected, slices.Collect(AllContents(resourceSet)))
	assert.Equal(t, expected, slices.Collect(allContentsResolve(resourceSet, false)))
	assert.Equal(t, expected[:3], slices.Collect(func(yield func(EObject) bool) {
		n := 0
		for o := range AllContents(resourceSet) {
			if n++; n > 3 || !yield(o) {
				return
			}
		}
	}))
}

func TestAllContents_OtherNotifier(t *testing.T) {
	assert.Empty(t, slices.Collect(AllContents(NewMockENotifier(t))))
}
