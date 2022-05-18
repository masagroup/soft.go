// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

type treeIterator struct {
	object      interface{}
	data        []EIterator
	root        bool
	getChildren func(interface{}) EIterator
}

func newTreeIterator(object interface{}, root bool, getChildren func(interface{}) EIterator) *treeIterator {
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

func (it *treeIterator) Next() interface{} {
	if it.data == nil {
		// Yield that mapping, create a stack, and add it to the stack.
		current := it.getChildren(it.object)
		it.data = append(it.data, current)
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
		it.data = append(it.data, iterator)
	} else {
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
