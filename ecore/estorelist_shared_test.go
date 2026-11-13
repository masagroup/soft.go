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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newSharedStoreList(t *testing.T, owner EObject, feature EStructuralFeature, shared *EStoreListSharedState) *EStoreList {
	t.Helper()
	list := &EStoreList{}
	list.InitializeShared(owner, feature, shared)
	list.SetInterfaces(list)
	return list
}

func TestEStoreListSharedState_NilStore(t *testing.T) {
	assert.Panics(t, func() { NewEStoreListSharedState(nil) })
}

// A list initialized with a shared state does not ask its store for its size when it is made,
// asks it at each Size, and reads its elements from its store.
func TestEStoreList_Shared_Read(t *testing.T) {
	mockOwner := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockStore := NewMockEStore(t)
	mockFeature.EXPECT().IsUnique().Return(true).Once()
	shared := NewEStoreListSharedState(mockStore)
	list := newSharedStoreList(t, mockOwner, mockFeature, shared)
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)

	assert.Same(t, mockStore, list.GetEStore())
	assert.False(t, list.IsCache())

	mockStore.EXPECT().Size(mockOwner, mockFeature).Return(2) // Size, Empty and the index check of Get
	assert.Equal(t, 2, list.Size())
	assert.False(t, list.Empty())

	mockStore.EXPECT().Get(mockOwner, mockFeature, 1).Return("b").Once()
	assert.Equal(t, "b", list.Get(1))

	mockStore.EXPECT().ToArray(mockOwner, mockFeature).Return([]any{"a", "b"}).Once()
	assert.Equal(t, []any{"a", "b"}, list.ToArray())
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)
}

// The first change of a list sharing a state gives it a state of its own, with the size its
// store gives: the other lists sharing the state, and the state, are unchanged.
func TestEStoreList_Shared_AddOwnState(t *testing.T) {
	mockOwner := NewMockEObject(t)
	mockOther := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockStore := NewMockEStore(t)
	mockFeature.EXPECT().IsUnique().Return(false).Twice()
	shared := NewEStoreListSharedState(mockStore)
	list := newSharedStoreList(t, mockOwner, mockFeature, shared)
	other := newSharedStoreList(t, mockOther, mockFeature, shared)

	mockStore.EXPECT().Size(mockOwner, mockFeature).Return(3) // Add checks the size, then its state is made
	mockStore.EXPECT().Add(mockOwner, mockFeature, 3, "d").Once()
	mockOwner.EXPECT().EDeliver().Return(false).Once()
	assert.True(t, list.Add("d"))
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)

	// its own state knows its size: no more store Size
	assert.Equal(t, 4, list.Size())
	assert.NotSame(t, &shared.state, list.state.Load())
	assert.Equal(t, int32(-1), shared.state.size)

	// the other list still shares the state
	assert.Same(t, &shared.state, other.state.Load())
	mockStore.EXPECT().Size(mockOther, mockFeature).Return(1).Once()
	assert.Equal(t, 1, other.Size())
	mock.AssertExpectationsForObjects(t, mockOther, mockFeature, mockStore)
}

// SetCache on a list sharing a state caches the elements of its store in a state of its own.
func TestEStoreList_Shared_SetCache(t *testing.T) {
	mockOwner := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockStore := NewMockEStore(t)
	mockFeature.EXPECT().IsUnique().Return(true).Once()
	shared := NewEStoreListSharedState(mockStore)
	list := newSharedStoreList(t, mockOwner, mockFeature, shared)

	// no cache already: nothing to do
	list.SetCache(false)
	assert.Same(t, &shared.state, list.state.Load())

	mockStore.EXPECT().Size(mockOwner, mockFeature).Return(2).Once()
	mockStore.EXPECT().ToArray(mockOwner, mockFeature).Return([]any{"a", "b"}).Once()
	list.SetCache(true)
	assert.True(t, list.IsCache())
	assert.Equal(t, 2, list.Size())
	assert.Equal(t, "b", list.Get(1))
	assert.Equal(t, int32(-1), shared.state.size)
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)
}

// SetEStore on a list sharing a state gives it a state of its own, with the elements of the
// former store when the new one is nil.
func TestEStoreList_Shared_SetEStoreNil(t *testing.T) {
	mockOwner := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockStore := NewMockEStore(t)
	mockFeature.EXPECT().IsUnique().Return(true).Once()
	shared := NewEStoreListSharedState(mockStore)
	list := newSharedStoreList(t, mockOwner, mockFeature, shared)

	mockStore.EXPECT().Size(mockOwner, mockFeature).Return(1).Once()
	mockStore.EXPECT().ToArray(mockOwner, mockFeature).Return([]any{"a"}).Once()
	list.SetEStore(nil)
	require.Nil(t, list.GetEStore())
	assert.Equal(t, 1, list.Size())
	assert.Equal(t, "a", list.Get(0))
	assert.Same(t, mockStore, shared.state.store)
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)
}

// Remove, Clear and the loops on the size of IndexOf and Contains work on a list sharing a state.
func TestEStoreList_Shared_RemoveClearIndexOf(t *testing.T) {
	mockOwner := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockStore := NewMockEStore(t)
	mockFeature.EXPECT().IsUnique().Return(true).Once()
	shared := NewEStoreListSharedState(mockStore)
	list := newSharedStoreList(t, mockOwner, mockFeature, shared)

	mockStore.EXPECT().IndexOf(mockOwner, mockFeature, "b").Return(1).Once()
	assert.Equal(t, 1, list.IndexOf("b"))

	mockStore.EXPECT().Size(mockOwner, mockFeature).Return(2) // RemoveAt checks the index, then its state is made
	mockStore.EXPECT().Remove(mockOwner, mockFeature, 0, true).Return("a").Once()
	mockOwner.EXPECT().EDeliver().Return(false).Once()
	assert.Equal(t, "a", list.RemoveAt(0))
	assert.Equal(t, 1, list.Size())

	mockStore.EXPECT().ToArray(mockOwner, mockFeature).Return([]any{"b"}).Once()
	mockStore.EXPECT().Clear(mockOwner, mockFeature).Once()
	mockOwner.EXPECT().EDeliver().Return(false).Once()
	list.Clear()
	assert.Equal(t, 0, list.Size())
	assert.True(t, list.Empty())
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)
}
