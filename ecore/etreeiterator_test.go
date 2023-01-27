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
	emptyList := NewImmutableEList(nil)
	mockObject := NewMockEObject(t)
	mockObject.EXPECT().EContents().Return(emptyList)
	it := newEAllContentsIterator(mockObject)
	assert.False(t, it.HasNext())
}

func TestEAllContentsIteratorNotEmpty(t *testing.T) {
	emptyList := NewImmutableEList(nil)
	mockObject := NewMockEObject(t)
	mockChild1 := NewMockEObject(t)
	mockGrandChild1 := NewMockEObject(t)
	mockGrandChild2 := NewMockEObject(t)
	mockChild2 := NewMockEObject(t)
	mockObject.EXPECT().EContents().Return(NewImmutableEList([]any{mockChild1, mockChild2}))
	mockChild1.EXPECT().EContents().Return(NewImmutableEList([]any{mockGrandChild1, mockGrandChild2}))
	mockGrandChild1.EXPECT().EContents().Return(emptyList)
	mockGrandChild2.EXPECT().EContents().Return(emptyList)
	mockChild2.EXPECT().EContents().Return(emptyList)

	var result []any
	for it := newEAllContentsIterator(mockObject); it.HasNext(); {
		result = append(result, it.Next())
	}
	assert.Equal(t, []any{mockChild1, mockGrandChild1, mockGrandChild2, mockChild2}, result)
}
