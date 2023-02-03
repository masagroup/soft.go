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
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBasicEStoreList_Constructors(t *testing.T) {
	{
		mockOwner := NewMockEObject(t)
		mockFeature := NewMockEStructuralFeature(t)
		mockStore := NewMockEStore(t)
		list := NewBasicEStoreList(mockOwner, mockFeature, mockStore)
		assert.NotNil(t, list)
	}
	{
		mockOwner := NewMockEObject(t)
		mockReference := NewMockEReference(t)
		mockStore := NewMockEStore(t)
		mockReference.EXPECT().IsContainment().Return(true).Once()
		mockReference.EXPECT().IsResolveProxies().Return(false).Once()
		mockReference.EXPECT().IsUnsettable().Return(false).Once()
		mockReference.EXPECT().GetEOpposite().Return(nil).Once()
		list := NewBasicEStoreList(mockOwner, mockReference, mockStore)
		assert.NotNil(t, list)
	}
	{
		mockOwner := NewMockEObject(t)
		mockReference := NewMockEReference(t)
		mockOpposite := NewMockEReference(t)
		mockStore := NewMockEStore(t)
		mockReference.EXPECT().IsContainment().Return(true).Once()
		mockReference.EXPECT().IsResolveProxies().Return(false).Once()
		mockReference.EXPECT().IsUnsettable().Return(false).Once()
		mockReference.EXPECT().GetEOpposite().Return(mockOpposite).Once()
		list := NewBasicEStoreList(mockOwner, mockReference, mockStore)
		assert.NotNil(t, list)
	}
	{
		mockOwner := NewMockEObject(t)
		mockReference := NewMockEReference(t)
		mockStore := NewMockEStore(t)
		mockReference.EXPECT().IsContainment().Return(false).Once()
		mockReference.EXPECT().IsResolveProxies().Return(false).Once()
		mockReference.EXPECT().IsUnsettable().Return(false).Once()
		mockReference.EXPECT().GetEOpposite().Return(nil).Once()
		list := NewBasicEStoreList(mockOwner, mockReference, mockStore)
		assert.NotNil(t, list)
	}
	{
		mockOwner := NewMockEObject(t)
		mockReference := NewMockEReference(t)
		mockOpposite := NewMockEReference(t)
		mockStore := NewMockEStore(t)
		mockReference.EXPECT().IsContainment().Return(false).Once()
		mockReference.EXPECT().IsResolveProxies().Return(false).Once()
		mockReference.EXPECT().IsUnsettable().Return(false).Once()
		mockReference.EXPECT().GetEOpposite().Return(mockOpposite).Once()
		list := NewBasicEStoreList(mockOwner, mockReference, mockStore)
		assert.NotNil(t, list)
	}
}

func TestBasicEStoreList_Accessors(t *testing.T) {
	mockOwner := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockStore := NewMockEStore(t)
	list := NewBasicEStoreList(mockOwner, mockFeature, mockStore)

	assert.Equal(t, mockOwner, list.GetOwner())
	assert.Equal(t, mockFeature, list.GetFeature())
	assert.Equal(t, mockStore, list.GetStore())

	mockClass := NewMockEClass(t)
	mockClass.EXPECT().GetFeatureID(mockFeature).Return(0).Once()
	mockOwner.EXPECT().EClass().Return(mockClass).Once()
	assert.Equal(t, 0, list.GetFeatureID())
	mock.AssertExpectationsForObjects(t, mockOwner, mockClass, mockFeature, mockStore)
}

func TestBasicEStoreList_Add(t *testing.T) {
	mockOwner := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockStore := NewMockEStore(t)
	mockAdapter := NewMockEAdapter(t)
	list := NewBasicEStoreList(mockOwner, mockFeature, mockStore)

	// already present
	mockStore.EXPECT().Size(mockOwner, mockFeature).Return(0).Twice()
	mockStore.EXPECT().Contains(mockOwner, mockFeature, 1).Return(true).Once()
	assert.False(t, list.Add(1))
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)

	// add 1 to the list
	mockStore.EXPECT().Size(mockOwner, mockFeature).Return(0).Twice()
	mockStore.EXPECT().Contains(mockOwner, mockFeature, 1).Return(false).Once()
	mockStore.EXPECT().Add(mockOwner, mockFeature, 0, 1).Once()
	mockOwner.EXPECT().EDeliver().Return(false).Once()
	assert.True(t, list.Add(1))
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)

	// add 2 to the list
	mockStore.EXPECT().Size(mockOwner, mockFeature).Return(1).Twice()
	mockStore.EXPECT().Contains(mockOwner, mockFeature, 2).Return(false).Once()
	mockStore.EXPECT().Add(mockOwner, mockFeature, 1, 2).Once()
	mockOwner.EXPECT().EDeliver().Return(true).Once()
	mockOwner.EXPECT().EAdapters().Return(NewImmutableEList([]any{mockAdapter})).Once()
	mockOwner.EXPECT().ENotify(mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == mockOwner && n.GetFeature() == mockFeature && n.GetEventType() == ADD && n.GetNewValue() == 2
	})).Once()
	assert.True(t, list.Add(2))
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)
}

func TestBasicEStoreList_AddReferenceContainmentNoOpposite(t *testing.T) {
	mockOwner := NewMockEObject(t)
	mockReference := NewMockEReference(t)
	mockStore := NewMockEStore(t)
	mockReference.EXPECT().IsContainment().Return(true).Once()
	mockReference.EXPECT().IsResolveProxies().Return(false).Once()
	mockReference.EXPECT().IsUnsettable().Return(false).Once()
	mockReference.EXPECT().GetEOpposite().Return(nil).Once()
	list := NewBasicEStoreList(mockOwner, mockReference, mockStore)
	assert.NotNil(t, list)

	mockObject := NewMockEObjectInternal(t)
	mockStore.EXPECT().Size(mockOwner, mockReference).Return(0).Twice()
	mockStore.EXPECT().Contains(mockOwner, mockReference, mockObject).Return(false).Once()
	mockStore.EXPECT().Add(mockOwner, mockReference, 0, mockObject).Once()
	mockReference.EXPECT().GetFeatureID().Return(0).Once()
	mockObject.EXPECT().EInverseAdd(mockOwner, EOPPOSITE_FEATURE_BASE-0, nil).Return(nil).Once()
	mockOwner.EXPECT().EDeliver().Return(false).Once()
	assert.True(t, list.Add(mockObject))
	mock.AssertExpectationsForObjects(t, mockOwner, mockReference, mockStore, mockObject)
}

func TestBasicEStoreList_AddReferenceContainmentOpposite(t *testing.T) {
	mockOwner := NewMockEObject(t)
	mockReference := NewMockEReference(t)
	mockOpposite := NewMockEReference(t)
	mockStore := NewMockEStore(t)
	mockReference.EXPECT().IsContainment().Return(true).Once()
	mockReference.EXPECT().IsResolveProxies().Return(false).Once()
	mockReference.EXPECT().IsUnsettable().Return(false).Once()
	mockReference.EXPECT().GetEOpposite().Return(mockOpposite).Once()
	list := NewBasicEStoreList(mockOwner, mockReference, mockStore)
	assert.NotNil(t, list)

	mockObject := NewMockEObjectInternal(t)
	mockClass := NewMockEClass(t)
	mockStore.EXPECT().Size(mockOwner, mockReference).Return(0).Twice()
	mockStore.EXPECT().Contains(mockOwner, mockReference, mockObject).Return(false).Once()
	mockStore.EXPECT().Add(mockOwner, mockReference, 0, mockObject).Once()
	mockReference.EXPECT().GetEOpposite().Return(mockOpposite).Once()
	mockObject.EXPECT().EClass().Return(mockClass).Once()
	mockClass.EXPECT().GetFeatureID(mockOpposite).Return(1).Once()
	mockObject.EXPECT().EInverseAdd(mockOwner, 1, nil).Return(nil).Once()
	mockOwner.EXPECT().EDeliver().Return(false).Once()
	assert.True(t, list.Add(mockObject))
	mock.AssertExpectationsForObjects(t, mockOwner, mockReference, mockStore, mockObject, mockClass, mockOpposite)
}

func TestBasicEStoreList_AddWithNotification(t *testing.T) {
	mockOwner := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockStore := NewMockEStore(t)
	mockNotifications := NewMockENotificationChain(t)
	mockClass := NewMockEClass(t)
	mockAdapter := NewMockEAdapter(t)
	list := NewBasicEStoreList(mockOwner, mockFeature, mockStore)

	// add 1
	mockStore.EXPECT().Size(mockOwner, mockFeature).Return(1).Once()
	mockStore.EXPECT().Add(mockOwner, mockFeature, 1, 2).Once()
	mockOwner.EXPECT().EDeliver().Return(true).Once()
	mockOwner.EXPECT().EAdapters().Return(NewImmutableEList([]any{mockAdapter})).Once()
	mockOwner.EXPECT().EClass().Return(mockClass)
	mockClass.EXPECT().GetFeatureID(mockFeature).Return(1)
	mockNotifications.EXPECT().Add(mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == mockOwner && n.GetFeature() == mockFeature && n.GetFeatureID() == 1 && n.GetEventType() == ADD && n.GetNewValue() == 2
	})).Return(true).Once()
	assert.Equal(t, mockNotifications, list.AddWithNotification(2, mockNotifications))
	mock.AssertExpectationsForObjects(t, mockOwner, mockClass, mockFeature, mockStore, mockAdapter, mockNotifications)
}

func TestBasicEStoreList_Insert(t *testing.T) {
	mockOwner := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockStore := NewMockEStore(t)
	list := NewBasicEStoreList(mockOwner, mockFeature, mockStore)

	mockStore.EXPECT().Size(mockOwner, mockFeature).Return(1).Once()
	assert.Panics(t, func() {
		list.Insert(-1, 0)
	})
	mock.AssertExpectationsForObjects(t, mockStore)

	mockStore.EXPECT().Size(mockOwner, mockFeature).Return(1).Twice()
	assert.Panics(t, func() {
		list.Insert(2, 0)
	})
	mock.AssertExpectationsForObjects(t, mockStore)
}

func TestBasicEStoreList_AddAll(t *testing.T) {
	mockOwner := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockStore := NewMockEStore(t)
	list := NewBasicEStoreList(mockOwner, mockFeature, mockStore)

	mockStore.EXPECT().Size(mockOwner, mockFeature).Return(0).Twice()
	mockStore.EXPECT().Contains(mockOwner, mockFeature, 1).Return(false).Once()
	mockStore.EXPECT().Add(mockOwner, mockFeature, 0, 1).Once()
	mockOwner.EXPECT().EDeliver().Return(false).Once()
	assert.True(t, list.AddAll(NewImmutableEList([]any{1})))
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)
}

func TestBasicEStoreList_InsertAll(t *testing.T) {
	mockOwner := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockStore := NewMockEStore(t)
	mockAdapter := NewMockEAdapter(t)
	list := NewBasicEStoreList(mockOwner, mockFeature, mockStore)

	// invalid index
	mockStore.EXPECT().Size(mockOwner, mockFeature).Return(1).Once()
	assert.Panics(t, func() {
		list.InsertAll(-1, NewImmutableEList([]any{}))
	})
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)

	// already present element
	mockStore.EXPECT().Size(mockOwner, mockFeature).Return(0).Once()
	mockStore.EXPECT().Contains(mockOwner, mockFeature, 1).Return(true).Once()
	assert.False(t, list.InsertAll(0, NewImmutableEList([]any{1})))
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)

	// single element inserted
	mockStore.EXPECT().Size(mockOwner, mockFeature).Return(0).Once()
	mockStore.EXPECT().Contains(mockOwner, mockFeature, 1).Return(false).Once()
	mockStore.EXPECT().Add(mockOwner, mockFeature, 0, 1).Once()
	mockOwner.EXPECT().EDeliver().Return(true).Once()
	mockOwner.EXPECT().EAdapters().Return(NewImmutableEList([]any{mockAdapter})).Once()
	mockOwner.EXPECT().ENotify(mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == mockOwner && n.GetFeature() == mockFeature && n.GetEventType() == ADD && n.GetNewValue() == 1
	})).Once()
	assert.True(t, list.InsertAll(0, NewImmutableEList([]any{1})))
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)

	mockStore.EXPECT().Size(mockOwner, mockFeature).Return(0).Once()
	mockStore.EXPECT().Contains(mockOwner, mockFeature, 1).Return(false).Once()
	mockStore.EXPECT().Contains(mockOwner, mockFeature, 2).Return(false).Once()
	mockStore.EXPECT().Add(mockOwner, mockFeature, 0, 1).Once()
	mockStore.EXPECT().Add(mockOwner, mockFeature, 1, 2).Once()
	mockOwner.EXPECT().EDeliver().Return(true).Once()
	mockOwner.EXPECT().EAdapters().Return(NewImmutableEList([]any{mockAdapter})).Once()
	mockOwner.EXPECT().ENotify(mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == mockOwner && n.GetFeature() == mockFeature && n.GetEventType() == ADD_MANY
	})).Once()
	assert.True(t, list.InsertAll(0, NewImmutableEList([]any{1, 2})))
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)

}

func TestBasicEStoreList_MoveObject(t *testing.T) {
	mockOwner := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockStore := NewMockEStore(t)
	mockObject := NewMockEObjectInternal(t)
	mockAdapter := NewMockEAdapter(t)
	list := NewBasicEStoreList(mockOwner, mockFeature, mockStore)

	mockStore.EXPECT().IndexOf(mockOwner, mockFeature, 1).Return(-1).Once()
	assert.Panics(t, func() {
		list.MoveObject(1, 1)
	})
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore, mockObject)

	mockStore.EXPECT().IndexOf(mockOwner, mockFeature, 1).Return(0).Once()
	mockStore.EXPECT().Size(mockOwner, mockFeature).Return(1).Twice()
	mockStore.EXPECT().Move(mockOwner, mockFeature, 1, 0).Return(mockObject).Once()
	mockOwner.EXPECT().EDeliver().Return(true).Once()
	mockOwner.EXPECT().EAdapters().Return(NewImmutableEList([]any{mockAdapter})).Once()
	mockOwner.EXPECT().ENotify(mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == mockOwner && n.GetFeature() == mockFeature && n.GetEventType() == MOVE
	})).Once()
	list.MoveObject(1, 1)
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore, mockObject)
}

func TestBasicEStoreList_Move(t *testing.T) {
	mockOwner := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockStore := NewMockEStore(t)
	mockObject := NewMockEObjectInternal(t)
	mockAdapter := NewMockEAdapter(t)
	list := NewBasicEStoreList(mockOwner, mockFeature, mockStore)

	mockStore.EXPECT().Size(mockOwner, mockFeature).Return(1).Once()
	assert.Panics(t, func() {
		list.Move(-1, 1)
	})
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore, mockObject)

	mockStore.EXPECT().Size(mockOwner, mockFeature).Return(1).Twice()
	mockStore.EXPECT().Move(mockOwner, mockFeature, 1, 0).Return(mockObject).Once()
	mockOwner.EXPECT().EDeliver().Return(true).Once()
	mockOwner.EXPECT().EAdapters().Return(NewImmutableEList([]any{mockAdapter})).Once()
	mockOwner.EXPECT().ENotify(mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == mockOwner && n.GetFeature() == mockFeature && n.GetEventType() == MOVE
	})).Once()
	list.Move(0, 1)
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore, mockObject)
}

func TestBasicEStoreList_Get_NoProxy(t *testing.T) {
	mockOwner := NewMockEObject(t)
	mockReference := NewMockEReference(t)
	mockOpposite := NewMockEReference(t)
	mockStore := NewMockEStore(t)
	mockObject := NewMockEObject(t)
	mockReference.EXPECT().IsContainment().Return(false).Once()
	mockReference.EXPECT().IsResolveProxies().Return(false).Once()
	mockReference.EXPECT().IsUnsettable().Return(false).Once()
	mockReference.EXPECT().GetEOpposite().Return(mockOpposite).Once()
	list := NewBasicEStoreList(mockOwner, mockReference, mockStore)
	assert.NotNil(t, list)

	mockStore.EXPECT().Get(list.owner, list.feature, 0).Return(mockObject).Once()
	assert.Equal(t, mockObject, list.Get(0))
}

func TestBasicEStoreList_Get_Proxy(t *testing.T) {
	mockOwner := NewMockEObjectInternal(t)
	mockReference := NewMockEReference(t)
	mockOpposite := NewMockEReference(t)
	mockStore := NewMockEStore(t)
	mockObject := NewMockEObjectInternal(t)
	mockResolved := NewMockEObjectInternal(t)
	mockReference.EXPECT().IsContainment().Return(true).Once()
	mockReference.EXPECT().IsResolveProxies().Return(true).Once()
	mockReference.EXPECT().IsUnsettable().Return(false).Once()
	mockReference.EXPECT().GetEOpposite().Return(mockOpposite).Once()
	list := NewBasicEStoreList(mockOwner, mockReference, mockStore)
	assert.NotNil(t, list)
	mock.AssertExpectationsForObjects(t, mockOwner, mockReference, mockStore, mockObject, mockOpposite)

	// no proxy object
	mockStore.EXPECT().Get(list.owner, list.feature, 0).Return(mockObject).Once()
	mockObject.EXPECT().EIsProxy().Return(false).Once()
	assert.Equal(t, mockObject, list.Get(0))
	mock.AssertExpectationsForObjects(t, mockOwner, mockReference, mockStore, mockObject, mockOpposite, mockResolved)

	// proxy object
	mockClass := NewMockEClass(t)
	mockStore.EXPECT().Get(list.owner, list.feature, 0).Return(mockObject).Once()
	mockObject.EXPECT().EIsProxy().Return(true).Once()
	mockOwner.EXPECT().EResolveProxy(mockObject).Return(mockResolved).Once()
	mockStore.EXPECT().Set(list.owner, list.feature, 0, mockResolved).Return(mockResolved).Return(mockObject).Once()
	mockReference.EXPECT().GetEOpposite().Return(mockOpposite).Once()
	mockClass.EXPECT().GetFeatureID(mockOpposite).Return(0).Once()
	mockObject.EXPECT().EClass().Return(mockClass).Once()
	mockObject.EXPECT().EInverseRemove(mockOwner, 0, nil).Return(nil).Once()
	mockResolved.EXPECT().EInternalContainer().Return(nil).Once()
	mockReference.EXPECT().GetEOpposite().Return(mockOpposite).Once()
	mockClass.EXPECT().GetFeatureID(mockOpposite).Return(0).Once()
	mockResolved.EXPECT().EClass().Return(mockClass).Once()
	mockResolved.EXPECT().EInverseAdd(mockOwner, 0, nil).Return(nil).Once()
	mockOwner.EXPECT().EDeliver().Return(false).Once()
	assert.Equal(t, mockResolved, list.Get(0))
	mock.AssertExpectationsForObjects(t, mockOwner, mockReference, mockStore, mockObject, mockOpposite, mockResolved, mockClass)
}

func TestBasicEStoreList_Set(t *testing.T) {
	mockOwner := NewMockEObjectInternal(t)
	mockReference := NewMockEReference(t)
	mockStore := NewMockEStore(t)
	mockNewObject := NewMockEObjectInternal(t)
	mockOldObject := NewMockEObjectInternal(t)
	mockReference.EXPECT().IsContainment().Return(true).Once()
	mockReference.EXPECT().IsResolveProxies().Return(false).Once()
	mockReference.EXPECT().IsUnsettable().Return(false).Once()
	mockReference.EXPECT().GetEOpposite().Return(nil).Once()
	list := NewBasicEStoreList(mockOwner, mockReference, mockStore)
	assert.NotNil(t, list)
	mock.AssertExpectationsForObjects(t, mockOwner, mockReference, mockStore)

	mockStore.EXPECT().Size(mockOwner, mockReference).Return(1).Once()
	assert.Panics(t, func() {
		list.Set(-1, mockNewObject)
	})
	mock.AssertExpectationsForObjects(t, mockOwner, mockReference, mockStore, mockNewObject)

	mockStore.EXPECT().Size(list.owner, list.feature).Return(2).Once()
	mockStore.EXPECT().IndexOf(list.owner, list.feature, mockNewObject).Return(1).Once()
	assert.Panics(t, func() {
		list.Set(0, mockNewObject)
	})
	mock.AssertExpectationsForObjects(t, mockOwner, mockReference, mockStore, mockNewObject)

	mockStore.EXPECT().Size(list.owner, list.feature).Return(1).Once()
	mockStore.EXPECT().IndexOf(list.owner, list.feature, mockNewObject).Return(-1).Once()
	mockStore.EXPECT().Set(list.owner, list.feature, 0, mockNewObject).Return(mockOldObject).Once()
	mockReference.EXPECT().GetFeatureID().Return(0).Once()
	mockOldObject.EXPECT().EInverseRemove(mockOwner, EOPPOSITE_FEATURE_BASE-0, nil).Return(nil).Once()
	mockReference.EXPECT().GetFeatureID().Return(0).Once()
	mockNewObject.EXPECT().EInverseAdd(mockOwner, EOPPOSITE_FEATURE_BASE-0, nil).Return(nil).Once()
	mockOwner.EXPECT().EDeliver().Return(false).Once()
	assert.Equal(t, mockOldObject, list.Set(0, mockNewObject))
	mock.AssertExpectationsForObjects(t, mockOwner, mockReference, mockStore, mockNewObject, mockOldObject)

}

func TestBasicEStoreList_SetWithNotification(t *testing.T) {
	mockOwner := NewMockEObjectInternal(t)
	mockReference := NewMockEReference(t)
	mockOpposite := NewMockEReference(t)
	mockStore := NewMockEStore(t)
	mockNewObject := NewMockEObjectInternal(t)
	mockOldObject := NewMockEObjectInternal(t)
	mockAdapter := NewMockEAdapter(t)
	mockReference.EXPECT().IsContainment().Return(true).Once()
	mockReference.EXPECT().IsResolveProxies().Return(true).Once()
	mockReference.EXPECT().IsUnsettable().Return(false).Once()
	mockReference.EXPECT().GetEOpposite().Return(mockOpposite).Once()
	list := NewBasicEStoreList(mockOwner, mockReference, mockStore)
	assert.NotNil(t, list)
	mock.AssertExpectationsForObjects(t, mockOwner, mockReference, mockStore, mockOpposite)

	mockStore.EXPECT().Set(list.owner, list.feature, 0, mockNewObject).Return(mockOldObject).Once()
	mockOwner.EXPECT().EDeliver().Return(true).Once()
	mockOwner.EXPECT().EAdapters().Return(NewImmutableEList([]any{mockAdapter})).Once()
	assert.NotNil(t, list.SetWithNotification(0, mockNewObject, nil))
	mock.AssertExpectationsForObjects(t, mockOwner, mockReference, mockStore, mockOpposite, mockNewObject, mockOldObject)
}

func TestBasicEStoreList_RemoveAt(t *testing.T) {
	mockOwner := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockStore := NewMockEStore(t)
	list := NewBasicEStoreList(mockOwner, mockFeature, mockStore)

	mockStore.EXPECT().Size(mockOwner, mockFeature).Return(1).Once()
	assert.Panics(t, func() {
		list.RemoveAt(-1)
	})

	mockStore.EXPECT().Size(mockOwner, mockFeature).Return(1).Once()
	mockStore.EXPECT().Remove(mockOwner, mockFeature, 0).Return(1).Once()
	mockOwner.EXPECT().EDeliver().Return(false).Once()
	assert.Equal(t, 1, list.RemoveAt(0))
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)
}

func TestBasicEStoreList_Remove(t *testing.T) {
	mockOwner := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockStore := NewMockEStore(t)
	list := NewBasicEStoreList(mockOwner, mockFeature, mockStore)

	mockStore.EXPECT().IndexOf(mockOwner, mockFeature, 1).Return(-1).Once()
	assert.False(t, list.Remove(1))
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)

	mockStore.EXPECT().IndexOf(mockOwner, mockFeature, 1).Return(0).Once()
	mockStore.EXPECT().Size(mockOwner, mockFeature).Return(1).Once()
	mockStore.EXPECT().Remove(mockOwner, mockFeature, 0).Return(1).Once()
	mockOwner.EXPECT().EDeliver().Return(false).Once()
	assert.True(t, list.Remove(1))
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)
}

func TestBasicEStoreList_RemoveWithNotification(t *testing.T) {
	mockOwner := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockStore := NewMockEStore(t)
	list := NewBasicEStoreList(mockOwner, mockFeature, mockStore)
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)

	mockNotifications := NewMockENotificationChain(t)
	mockStore.EXPECT().IndexOf(mockOwner, mockFeature, 1).Return(-1).Once()
	assert.Equal(t, mockNotifications, list.RemoveWithNotification(1, mockNotifications))
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore, mockNotifications)

	mockStore.EXPECT().IndexOf(mockOwner, mockFeature, 1).Return(0).Once()
	mockStore.EXPECT().Remove(mockOwner, mockFeature, 0).Return(1).Once()
	mockOwner.EXPECT().EDeliver().Return(false).Once()
	assert.Equal(t, mockNotifications, list.RemoveWithNotification(1, mockNotifications))
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore, mockNotifications)
}

func TestBasicEStoreList_RemoveAll(t *testing.T) {
	mockOwner := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockStore := NewMockEStore(t)
	list := NewBasicEStoreList(mockOwner, mockFeature, mockStore)
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)

	mockStore.EXPECT().Size(mockOwner, mockFeature).Return(1).Once()
	mockStore.EXPECT().Get(mockOwner, mockFeature, 0).Return(1).Once()
	mockStore.EXPECT().Size(mockOwner, mockFeature).Return(1).Once()
	mockStore.EXPECT().Remove(mockOwner, mockFeature, 0).Return(1).Once()
	mockOwner.EXPECT().EDeliver().Return(false)
	list.RemoveAll(NewImmutableEList([]any{1}))
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)

}

func TestBasicEStoreList_Size(t *testing.T) {
	mockOwner := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockStore := NewMockEStore(t)
	list := NewBasicEStoreList(mockOwner, mockFeature, mockStore)
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)

	mockStore.EXPECT().Size(mockOwner, mockFeature).Return(1).Once()
	assert.Equal(t, 1, list.Size())
}

func TestBasicEStoreList_Clear(t *testing.T) {
	mockOwner := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockStore := NewMockEStore(t)
	mockAdapter := NewMockEAdapter(t)
	list := NewBasicEStoreList(mockOwner, mockFeature, mockStore)
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)

	// empty list
	mockStore.EXPECT().ToArray(mockOwner, mockFeature).Return([]any{}).Once()
	mockStore.EXPECT().Clear(mockOwner, mockFeature).Once()
	mockOwner.EXPECT().EDeliver().Return(true).Once()
	mockOwner.EXPECT().EAdapters().Return(NewImmutableEList([]any{mockAdapter})).Once()
	mockOwner.EXPECT().ENotify(mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == mockOwner && n.GetFeature() == mockFeature && n.GetEventType() == REMOVE_MANY && n.GetNewValue() == nil && len(n.GetOldValue().([]any)) == 0
	}))
	list.Clear()
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore, mockAdapter)

	// single element list
	mockStore.EXPECT().ToArray(mockOwner, mockFeature).Return([]any{1}).Once()
	mockStore.EXPECT().Clear(mockOwner, mockFeature).Once()
	mockOwner.EXPECT().EDeliver().Return(true).Once()
	mockOwner.EXPECT().EAdapters().Return(NewImmutableEList([]any{mockAdapter})).Once()
	mockOwner.EXPECT().ENotify(mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == mockOwner && n.GetFeature() == mockFeature && n.GetEventType() == REMOVE && n.GetNewValue() == nil && n.GetOldValue() == 1
	}))
	list.Clear()
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore, mockAdapter)

	// multi element list
	mockStore.EXPECT().ToArray(mockOwner, mockFeature).Return([]any{1, 2}).Once()
	mockStore.EXPECT().Clear(mockOwner, mockFeature).Once()
	mockOwner.EXPECT().EDeliver().Return(true).Once()
	mockOwner.EXPECT().EAdapters().Return(NewImmutableEList([]any{mockAdapter})).Once()
	mockOwner.EXPECT().ENotify(mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == mockOwner && n.GetFeature() == mockFeature && n.GetEventType() == REMOVE_MANY && n.GetNewValue() == nil && reflect.DeepEqual(n.GetOldValue(), []any{1, 2})
	}))
	list.Clear()
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore, mockAdapter)

}

func TestBasicEStoreList_Empty(t *testing.T) {
	mockOwner := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockStore := NewMockEStore(t)
	list := NewBasicEStoreList(mockOwner, mockFeature, mockStore)
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)

	mockStore.EXPECT().IsEmpty(mockOwner, mockFeature).Return(true).Once()
	assert.True(t, list.Empty())
}

func TestBasicEStoreList_Contains(t *testing.T) {
	{
		mockOwner := NewMockEObject(t)
		mockFeature := NewMockEStructuralFeature(t)
		mockStore := NewMockEStore(t)
		mockObject := NewMockEObject(t)
		list := NewBasicEStoreList(mockOwner, mockFeature, mockStore)
		mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)

		mockStore.EXPECT().Contains(mockOwner, mockFeature, mockObject).Return(true).Once()
		assert.True(t, list.Contains(mockObject))
		mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore, mockObject)
	}
	{
		mockOwner := NewMockEObjectInternal(t)
		mockReference := NewMockEReference(t)
		mockOpposite := NewMockEReference(t)
		mockStore := NewMockEStore(t)
		mockObject := NewMockEObjectInternal(t)
		mockResolved := NewMockEObjectInternal(t)
		mockReference.EXPECT().IsContainment().Return(false).Once()
		mockReference.EXPECT().IsResolveProxies().Return(true).Once()
		mockReference.EXPECT().IsUnsettable().Return(false).Once()
		mockReference.EXPECT().GetEOpposite().Return(nil).Once()
		list := NewBasicEStoreList(mockOwner, mockReference, mockStore)
		assert.NotNil(t, list)
		mock.AssertExpectationsForObjects(t, mockOwner, mockReference, mockStore, mockObject, mockOpposite)

		mockStore.EXPECT().Contains(mockOwner, mockReference, mockResolved).Return(false).Once()
		mockStore.EXPECT().Size(mockOwner, mockReference).Return(1)
		mockStore.EXPECT().Get(mockOwner, mockReference, 0).Return(mockObject).Once()
		mockObject.EXPECT().EIsProxy().Return(true).Once()
		mockOwner.EXPECT().EResolveProxy(mockObject).Return(mockResolved).Once()
		assert.True(t, list.Contains(mockResolved))
	}
}

func TestBasicEStoreList_IndexOf(t *testing.T) {
	{
		mockOwner := NewMockEObject(t)
		mockFeature := NewMockEStructuralFeature(t)
		mockStore := NewMockEStore(t)
		mockObject := NewMockEObject(t)
		list := NewBasicEStoreList(mockOwner, mockFeature, mockStore)
		mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)

		mockStore.EXPECT().IndexOf(mockOwner, mockFeature, mockObject).Return(1).Once()
		assert.Equal(t, 1, list.IndexOf(mockObject))
		mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore, mockObject)
	}
	{
		mockOwner := NewMockEObjectInternal(t)
		mockReference := NewMockEReference(t)
		mockOpposite := NewMockEReference(t)
		mockStore := NewMockEStore(t)
		mockObject := NewMockEObjectInternal(t)
		mockResolved := NewMockEObjectInternal(t)
		mockReference.EXPECT().IsContainment().Return(false).Once()
		mockReference.EXPECT().IsResolveProxies().Return(true).Once()
		mockReference.EXPECT().IsUnsettable().Return(false).Once()
		mockReference.EXPECT().GetEOpposite().Return(nil).Once()
		list := NewBasicEStoreList(mockOwner, mockReference, mockStore)
		assert.NotNil(t, list)
		mock.AssertExpectationsForObjects(t, mockOwner, mockReference, mockStore, mockObject, mockOpposite)

		mockStore.EXPECT().IndexOf(mockOwner, mockReference, mockResolved).Return(-1).Once()
		mockStore.EXPECT().Size(mockOwner, mockReference).Return(1)
		mockStore.EXPECT().Get(mockOwner, mockReference, 0).Return(mockObject).Once()
		mockObject.EXPECT().EIsProxy().Return(true).Once()
		mockOwner.EXPECT().EResolveProxy(mockObject).Return(mockResolved).Once()
		assert.Equal(t, 0, list.IndexOf(mockResolved))
	}
}

func TestBasicEStoreList_Iterator(t *testing.T) {
	mockOwner := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockStore := NewMockEStore(t)
	list := NewBasicEStoreList(mockOwner, mockFeature, mockStore)
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)
	assert.NotNil(t, list.Iterator())
}

func TestBasicEStoreList_ToArray(t *testing.T) {
	mockOwner := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockStore := NewMockEStore(t)
	list := NewBasicEStoreList(mockOwner, mockFeature, mockStore)
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)

	mockStore.EXPECT().ToArray(mockOwner, mockFeature).Return([]any{1, 2})
	assert.Equal(t, []any{1, 2}, list.ToArray())
}

func TestBasicEStoreList_GetUnResolvedList(t *testing.T) {
	mockOwner := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockStore := NewMockEStore(t)
	list := NewBasicEStoreList(mockOwner, mockFeature, mockStore)
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)
	assert.NotNil(t, list.GetUnResolvedList())
}

func TestBasicEStoreList_RemoveRange(t *testing.T) {
	mockOwner := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockStore := NewMockEStore(t)
	mockObject := &mock.Mock{}
	mockObject2 := &mock.Mock{}
	mockAdapter := NewMockEAdapter(t)
	list := NewBasicEStoreList(mockOwner, mockFeature, mockStore)
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)

	mockStore.EXPECT().Remove(mockOwner, mockFeature, 0).Return(mockObject).Once()
	mockOwner.EXPECT().EDeliver().Return(true).Once()
	mockOwner.EXPECT().EAdapters().Return(NewImmutableEList([]any{mockAdapter})).Once()
	mockOwner.EXPECT().ENotify(mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == mockOwner &&
			n.GetFeature() == mockFeature &&
			n.GetEventType() == REMOVE &&
			n.GetNewValue() == nil &&
			n.GetOldValue() == mockObject
	}))
	list.RemoveRange(0, 1)
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore, mockAdapter)

	mockStore.EXPECT().Remove(mockOwner, mockFeature, 0).Return(mockObject).Once()
	mockStore.EXPECT().Remove(mockOwner, mockFeature, 1).Return(mockObject2).Once()
	mockOwner.EXPECT().EDeliver().Return(true).Once()
	mockOwner.EXPECT().EAdapters().Return(NewImmutableEList([]any{mockAdapter})).Once()
	mockOwner.EXPECT().ENotify(mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == mockOwner &&
			n.GetFeature() == mockFeature &&
			n.GetEventType() == REMOVE_MANY &&
			reflect.DeepEqual(n.GetNewValue(), []int{0, 1}) &&
			reflect.DeepEqual(n.GetOldValue(), []any{mockObject, mockObject2})
	}))
	list.RemoveRange(0, 2)
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore, mockAdapter)
}
