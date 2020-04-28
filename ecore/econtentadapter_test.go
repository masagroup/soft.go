// *****************************************************************************
//
// This file is part of a MASA library or program.
// Refer to the included end-user license agreement for restrictions.
//
// Copyright (c) 2019 MASA Group
//
// *****************************************************************************

package ecore

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestEContentAdapterSetTarget(t *testing.T) {
	adapter := NewEContentAdapter()
	nb := rand.Intn(10) + 1
	children := []interface{}{}
	for i := 0; i < nb; i++ {
		mockObject := new(MockEObject)
		mockAdapters := new(MockEList)

		mockObject.On("EAdapters").Return(mockAdapters)
		mockAdapters.On("Contains", adapter).Return(false)
		mockAdapters.On("Add", adapter).Return(true)
		mockAdapters.On("Remove", adapter).Return(true)
		children = append(children, mockObject)
	}

	mockChildren := NewImmutableEList(children)
	mockObject := new(MockEObject)
	mockAdapters := new(MockEList)
	mockObject.On("EContents").Return(mockChildren)

	// set adapter target -> this should recursively register adapter on all object children
	adapter.SetTarget(mockObject)

	// unset adapter target -> this should recursively unregister adapter on all object children
	adapter.UnSetTarget(mockObject)

	mock.AssertExpectationsForObjects(t, children...)
	mock.AssertExpectationsForObjects(t, mockObject, mockAdapters)
}

func TestEContentAdapterNotifyChanged(t *testing.T) {
	adapter := NewEContentAdapter()
	mockNotification := new(MockENotification)
	mockObject := new(MockEObject)

	mockAttribute := new(MockEAttribute)
	mockNotification.On("GetNotifier").Once().Return(mockObject)
	mockNotification.On("GetFeature").Once().Return(mockAttribute)
	adapter.NotifyChanged(mockNotification)

	mockReference := new(MockEReference)
	mockReference.On("IsContainment").Once().Return(false)
	mockNotification.On("GetNotifier").Once().Return(mockObject)
	mockNotification.On("GetFeature").Once().Return(mockReference)
	adapter.NotifyChanged(mockNotification)

	mock.AssertExpectationsForObjects(t, mockNotification, mockObject, mockAttribute, mockReference)
}

func TestEContentAdapterNotifyChanged_Resolve(t *testing.T) {
	adapter := NewEContentAdapter()
	mockNotification := new(MockENotification)
	mockObject := new(MockEObject)
	mockReference := new(MockEReference)
	mockOldObject := new(MockEObject)
	mockOldAdapters := new(MockEList)

	mockReference.On("IsContainment").Once().Return(true)
	mockOldObject.On("EAdapters").Once().Return(mockOldAdapters)
	mockOldAdapters.On("Contains", adapter).Once().Return(false)

	mockNotification.On("GetNotifier").Once().Return(mockObject)
	mockNotification.On("GetEventType").Once().Return(RESOLVE)
	mockNotification.On("GetFeature").Once().Return(mockReference)
	mockNotification.On("GetOldValue").Once().Return(mockOldObject)

	adapter.NotifyChanged(mockNotification)

	mock.AssertExpectationsForObjects(t, mockNotification, mockObject, mockReference, mockOldObject, mockOldAdapters)
}

func TestEContentAdapterNotifyChanged_Resolve_Contains(t *testing.T) {
	adapter := NewEContentAdapter()
	mockNotification := new(MockENotification)
	mockObject := new(MockEObject)
	mockReference := new(MockEReference)
	mockOldObject := new(MockEObject)
	mockOldAdapters := new(MockEList)
	mockNewObject := new(MockEObject)
	mockNewAdapters := new(MockEList)

	mockReference.On("IsContainment").Once().Return(true)

	mockOldObject.On("EAdapters").Twice().Return(mockOldAdapters)
	mockOldAdapters.On("Contains", adapter).Once().Return(true)
	mockOldAdapters.On("Remove", adapter).Once().Return(true)

	mockNewObject.On("EAdapters").Once().Return(mockNewAdapters)
	mockNewAdapters.On("Contains", adapter).Once().Return(false)
	mockNewAdapters.On("Add", adapter).Once().Return(true)

	mockNotification.On("GetNotifier").Once().Return(mockObject)
	mockNotification.On("GetEventType").Once().Return(RESOLVE)
	mockNotification.On("GetFeature").Once().Return(mockReference)
	mockNotification.On("GetOldValue").Once().Return(mockOldObject)
	mockNotification.On("GetNewValue").Once().Return(mockNewObject)

	adapter.NotifyChanged(mockNotification)

	mock.AssertExpectationsForObjects(t, mockNotification, mockObject, mockReference, mockOldObject, mockOldAdapters, mockNewObject, mockNewAdapters)
}

func TestEContentAdapterNotifyChanged_UnSet(t *testing.T) {

	adapter := NewEContentAdapter()
	mockNotification := new(MockENotification)
	mockObject := new(MockEObjectInternal)
	mockReference := new(MockEReference)
	mockOldObject := new(MockEObjectInternal)
	mockOldAdapters := new(MockEList)
	mockNewObject := new(MockEObjectInternal)
	mockNewAdapters := new(MockEList)

	mockReference.On("IsContainment").Once().Return(true)

	mockOldObject.On("EAdapters").Return(mockOldAdapters)
	mockOldObject.On("EInternalResource").Return(nil)
	mockOldAdapters.On("Remove", adapter).Once().Return(true)

	mockNewObject.On("EAdapters").Once().Return(mockNewAdapters)
	mockNewAdapters.On("Contains", adapter).Once().Return(false)
	mockNewAdapters.On("Add", adapter).Once().Return(true)

	mockNotification.On("GetNotifier").Once().Return(mockObject)
	mockNotification.On("GetEventType").Once().Return(UNSET)
	mockNotification.On("GetFeature").Once().Return(mockReference)
	mockNotification.On("GetOldValue").Once().Return(mockOldObject)
	mockNotification.On("GetNewValue").Once().Return(mockNewObject)

	adapter.NotifyChanged(mockNotification)

	mock.AssertExpectationsForObjects(t, mockNotification, mockObject, mockReference, mockOldObject, mockOldAdapters, mockNewObject, mockNewAdapters)
}

func TestEContentAdapterNotifyChanged_Set(t *testing.T) {

	adapter := NewEContentAdapter()
	mockNotification := new(MockENotification)
	mockObject := new(MockEObjectInternal)
	mockReference := new(MockEReference)
	mockOldObject := new(MockEObjectInternal)
	mockOldAdapters := new(MockEList)
	mockNewObject := new(MockEObjectInternal)
	mockNewAdapters := new(MockEList)

	mockReference.On("IsContainment").Once().Return(true)

	mockOldObject.On("EAdapters").Return(mockOldAdapters)
	mockOldObject.On("EInternalResource").Return(nil)
	mockOldAdapters.On("Remove", adapter).Once().Return(true)

	mockNewObject.On("EAdapters").Once().Return(mockNewAdapters)
	mockNewAdapters.On("Contains", adapter).Once().Return(false)
	mockNewAdapters.On("Add", adapter).Once().Return(true)

	mockNotification.On("GetNotifier").Once().Return(mockObject)
	mockNotification.On("GetEventType").Once().Return(SET)
	mockNotification.On("GetFeature").Once().Return(mockReference)
	mockNotification.On("GetOldValue").Once().Return(mockOldObject)
	mockNotification.On("GetNewValue").Once().Return(mockNewObject)

	adapter.NotifyChanged(mockNotification)

	mock.AssertExpectationsForObjects(t, mockNotification, mockObject, mockReference, mockOldObject, mockOldAdapters, mockNewObject, mockNewAdapters)

}

func TestEContentAdapterNotifyChanged_Add(t *testing.T) {

	adapter := NewEContentAdapter()
	mockNotification := new(MockENotification)
	mockObject := new(MockEObjectInternal)
	mockReference := new(MockEReference)
	mockNewObject := new(MockEObjectInternal)
	mockNewAdapters := new(MockEList)

	mockReference.On("IsContainment").Once().Return(true)

	mockNewObject.On("EAdapters").Once().Return(mockNewAdapters)
	mockNewAdapters.On("Contains", adapter).Once().Return(false)
	mockNewAdapters.On("Add", adapter).Once().Return(true)

	mockNotification.On("GetNotifier").Once().Return(mockObject)
	mockNotification.On("GetEventType").Once().Return(ADD)
	mockNotification.On("GetFeature").Once().Return(mockReference)
	mockNotification.On("GetNewValue").Once().Return(mockNewObject)

	adapter.NotifyChanged(mockNotification)

	mock.AssertExpectationsForObjects(t, mockNotification, mockObject, mockReference, mockNewObject, mockNewAdapters)
}

func TestEContentAdapterNotifyChanged_AddMany(t *testing.T) {

	adapter := NewEContentAdapter()
	mockNotification := new(MockENotification)
	mockObject := new(MockEObjectInternal)
	mockReference := new(MockEReference)

	mockReference.On("IsContainment").Once().Return(true)

	nb := rand.Intn(10) + 1
	mockChildren := []interface{}{}
	for i := 0; i < nb; i++ {
		mockObject := new(MockEObject)
		mockAdapters := new(MockEList)

		mockObject.On("EAdapters").Return(mockAdapters)
		mockAdapters.On("Contains", adapter).Return(false)
		mockAdapters.On("Add", adapter).Return(true)
		mockChildren = append(mockChildren, mockObject)
	}

	mockNotification.On("GetNotifier").Once().Return(mockObject)
	mockNotification.On("GetEventType").Once().Return(ADD_MANY)
	mockNotification.On("GetFeature").Once().Return(mockReference)
	mockNotification.On("GetNewValue").Once().Return(mockChildren)

	adapter.NotifyChanged(mockNotification)

	mock.AssertExpectationsForObjects(t, mockChildren...)
	mock.AssertExpectationsForObjects(t, mockNotification, mockObject, mockReference)
}

func TestEContentAdapterNotifyChanged_Remove(t *testing.T) {

	adapter := NewEContentAdapter()
	mockNotification := new(MockENotification)
	mockObject := new(MockEObjectInternal)
	mockReference := new(MockEReference)
	mockOldObject := new(MockEObjectInternal)
	mockOldAdapters := new(MockEList)

	mockReference.On("IsContainment").Once().Return(true)

	mockOldObject.On("EAdapters").Once().Return(mockOldAdapters)
	mockOldObject.On("EInternalResource").Return(nil)
	mockOldAdapters.On("Remove", adapter).Once().Return(true)

	mockNotification.On("GetNotifier").Return(mockObject)
	mockNotification.On("GetEventType").Return(REMOVE)
	mockNotification.On("GetFeature").Return(mockReference)
	mockNotification.On("GetOldValue").Return(mockOldObject)

	adapter.NotifyChanged(mockNotification)

	mock.AssertExpectationsForObjects(t, mockNotification, mockObject, mockReference, mockOldObject, mockOldAdapters)
}

func TestEContentAdapterNotifyChanged_RemoveMany(t *testing.T) {

	adapter := NewEContentAdapter()
	mockNotification := new(MockENotification)
	mockObject := new(MockEObjectInternal)
	mockReference := new(MockEReference)

	mockReference.On("IsContainment").Once().Return(true)

	nb := rand.Intn(10) + 1
	mockChildren := []interface{}{}
	for i := 0; i < nb; i++ {
		mockObject := new(MockEObjectInternal)
		mockAdapters := new(MockEList)

		mockObject.On("EAdapters").Return(mockAdapters)
		mockObject.On("EInternalResource").Return(nil)
		mockAdapters.On("Remove", adapter).Return(true)
		mockChildren = append(mockChildren, mockObject)
	}

	mockNotification.On("GetNotifier").Twice().Return(mockObject)
	mockNotification.On("GetFeature").Twice().Return(mockReference)
	mockNotification.On("GetEventType").Once().Return(REMOVE_MANY)
	mockNotification.On("GetOldValue").Once().Return(mockChildren)

	adapter.NotifyChanged(mockNotification)

	mock.AssertExpectationsForObjects(t, mockChildren...)
	mock.AssertExpectationsForObjects(t, mockNotification, mockObject, mockReference)
}
