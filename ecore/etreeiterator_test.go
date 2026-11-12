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

func TestETreeIteratorWithRoot(t *testing.T) {
	emptyList := NewImmutableEList(nil)
	mockObject := NewMockEObject(t)
	it := newTreeIterator(mockObject, true, func(i any) EIterator {
		return emptyList.Iterator()
	})
	assert.True(t, it.HasNext())
	assert.Equal(t, mockObject, it.Next())
	assert.False(t, it.HasNext())
}

func TestEAllContentsIteratorEmpty(t *testing.T) {
	mockObject := newMockEContentsObject(t)
	it := newEAllContentsIterator(newEAllContentsFrame(mockObject, nil), nil)
	assert.False(t, it.HasNext())
}

func TestEAllContentsIteratorNotEmpty(t *testing.T) {
	mockGrandChild1 := newMockEContentsObject(t)
	mockGrandChild2 := newMockEContentsObject(t)
	mockChild1 := newMockEContentsObject(t, mockGrandChild1, mockGrandChild2)
	mockChild2 := newMockEContentsObject(t)
	mockObject := newMockEContentsObject(t, mockChild1, mockChild2)

	var result []any
	for it := newEAllContentsIterator(newEAllContentsFrame(mockObject, nil), nil); it.HasNext(); {
		result = append(result, it.Next())
	}
	assert.Equal(t, []any{mockChild1, mockGrandChild1, mockGrandChild2, mockChild2}, result)
}

// treeIterator is the iterator EAllContents and GetAllContents were: a tree iterator of the
// iterators of the children of each object. The tests and benchmarks of the walks compare them to
// it.
type treeIterator struct {
	object      any
	current     EIterator
	getChildren func(any) EIterator
	data        []EIterator
	root        bool
}

func newTreeIterator(object any, root bool, getChildren func(any) EIterator) *treeIterator {
	return &treeIterator{object: object, root: root, getChildren: getChildren}
}

func (it *treeIterator) HasNext() bool {
	if it.data == nil && !it.root {
		return it.hasAnyChildren()
	} else {
		return it.hasMoreChildren()
	}
}

func (it *treeIterator) hasAnyChildren() bool {
	current := it.getChildren(it.object)
	it.data = append(it.data, current)
	return current.HasNext()
}

func (it *treeIterator) hasMoreChildren() bool {
	return it.data == nil || len(it.data) != 0 && it.data[len(it.data)-1].HasNext()
}

func (it *treeIterator) Next() any {
	if it.data == nil {
		// Yield that mapping, create a stack, and add it to the stack.
		it.current = it.getChildren(it.object)
		it.data = append(it.data, it.current)
		if it.root {
			return it.object
		}
	}

	// Get the top iterator, retrieve it's result, and record it as the one to which
	// remove will be delegated.
	current := it.data[len(it.data)-1]
	result := current.Next()

	// If the result about to be returned has children...
	iterator := it.getChildren(result)
	if iterator.HasNext() {
		// Add iterator to the stack.
		it.current = iterator
		it.data = append(it.data, iterator)
	} else {
		//
		it.current = nil
		// While the current iterator has no next...
		for !current.HasNext() {
			// Pop it from the stack.
			it.data = it.data[:len(it.data)-1]

			// If the stack is empty, we're done.
			if len(it.data) == 0 {
				break
			}

			// Get the next one down and then test it for has next.
			current = it.data[len(it.data)-1]
		}
	}
	return result
}

func (it *treeIterator) Prune() {
	if it.current != nil {
		if len(it.data) > 0 && it.data[len(it.data)-1] == it.current {
			// Pop it from the stack.
			it.data = it.data[:len(it.data)-1]
			// then pop any iterators above it that have no next.
			for len(it.data) > 0 && !it.data[len(it.data)-1].HasNext() {
				it.data = it.data[:len(it.data)-1]
			}
		}
		// clear the current iterator to prevent further pruning until the next call to next.
		it.current = nil
	}
}
