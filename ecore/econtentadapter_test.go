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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestEContentAdapter_SetTarget_EObject(t *testing.T) {
	adapter := NewEContentAdapter()
	// create a hierarchy of objects
	nb := rand.Intn(10) + 1
	mockObjects := []*MockEObject{}
	mockLists := []*MockEList{}
	children := []interface{}{}
	for i := 0; i < nb; i++ {
		mockObject := new(MockEObject)
		mockAdapters := new(MockEList)
		mockLists = append(mockLists, mockAdapters)
		mockObjects = append(mockObjects, mockObject)
		children = append(children, mockObject)
	}
	mockChildren := NewImmutableEList(children)
	mockObject := new(MockEObject)

	// set adapter target -> this should recursively register adapter on all object children
	for i := 0; i < nb; i++ {
		mockObject := mockObjects[i]
		mockAdapters := mockLists[i]
		if i%2 == 0 {
			mockObject.On("EIsProxy").Return(false).Once()
			mockObject.On("EAdapters").Return(mockAdapters).Once()
			mockAdapters.On("Contains", adapter).Return(false).Once()
			mockAdapters.On("Add", adapter).Return(true).Once()
		} else {
			mockObject.On("EIsProxy").Return(true).Once()
		}
	}
	mockObject.On("EContents").Return(mockChildren).Once()
	adapter.SetTarget(mockObject)
	mock.AssertExpectationsForObjects(t, children...)
	mock.AssertExpectationsForObjects(t, mockObject)

	// unset adapter target -> this should recursively unregister adapter on all object children
	for i := 0; i < nb; i++ {
		mockObject := mockObjects[i]
		mockAdapters := mockLists[i]
		mockObject.On("EAdapters").Return(mockAdapters).Once()
		mockAdapters.On("Remove", adapter).Return(true).Once()
	}
	mockObject.On("EContents").Return(mockChildren).Once()
	adapter.UnSetTarget(mockObject)
	mock.AssertExpectationsForObjects(t, children...)
	mock.AssertExpectationsForObjects(t, mockObject)
}

func TestEContentAdapter_SetTarget_EObject_ResolveProxies(t *testing.T) {
	adapter := NewEContentAdapter()
	// create a hierarchy of objects
	nb := rand.Intn(10) + 1
	mockObjects := []*MockEObject{}
	mockLists := []*MockEList{}
	children := []interface{}{}
	for i := 0; i < nb; i++ {
		mockObject := new(MockEObject)
		mockAdapters := new(MockEList)
		mockLists = append(mockLists, mockAdapters)
		mockObjects = append(mockObjects, mockObject)
		children = append(children, mockObject)
	}
	mockChildren := NewImmutableEList(children)
	mockObject := new(MockEObject)

	// set adapter target -> this should recursively register adapter on all object children
	for i := 0; i < nb; i++ {
		mockObject := mockObjects[i]
		mockAdapters := mockLists[i]
		mockObject.On("EIsProxy").Return(false).Once()
		mockObject.On("EAdapters").Return(mockAdapters).Once()
		mockAdapters.On("Contains", adapter).Return(false).Once()
		mockAdapters.On("Add", adapter).Return(true).Once()
	}
	mockObject.On("EContents").Return(mockChildren).Once()
	adapter.SetTarget(mockObject)
	mock.AssertExpectationsForObjects(t, children...)
	mock.AssertExpectationsForObjects(t, mockObject)

	// unset adapter target -> this should recursively unregister adapter on all object children
	for i := 0; i < nb; i++ {
		mockObject := mockObjects[i]
		mockAdapters := mockLists[i]
		mockObject.On("EAdapters").Return(mockAdapters).Once()
		mockAdapters.On("Remove", adapter).Return(true).Once()
	}
	mockObject.On("EContents").Return(mockChildren).Once()
	adapter.UnSetTarget(mockObject)
	mock.AssertExpectationsForObjects(t, children...)
	mock.AssertExpectationsForObjects(t, mockObject)
}

func TestEContentAdapter_SetTarget_EResourceSet(t *testing.T) {
	adapter := NewEContentAdapter()
	mockResourceSet := &MockEResourceSet{}
	mockResource := &MockEResource{}
	mockResourceAdapters := new(MockEList)

	// Set target
	mockResourceSet.On("GetResources").Return(NewImmutableEList([]interface{}{mockResource})).Once()
	mockResource.On("EAdapters").Return(mockResourceAdapters).Once()
	mockResourceAdapters.On("Contains", adapter).Return(false).Once()
	mockResourceAdapters.On("Add", adapter).Return(true).Once()
	adapter.SetTarget(mockResourceSet)
	mock.AssertExpectationsForObjects(t, mockResourceSet, mockResource, mockResourceAdapters)

	// UnSet target
	mockResourceSet.On("GetResources").Return(NewImmutableEList([]interface{}{mockResource})).Once()
	mockResource.On("EAdapters").Return(mockResourceAdapters).Once()
	mockResourceAdapters.On("Remove", adapter).Return(true).Once()
	adapter.UnSetTarget(mockResourceSet)
	mock.AssertExpectationsForObjects(t, mockResourceSet, mockResource, mockResourceAdapters)
}

func TestEContentAdapter_SetTarget_EResource(t *testing.T) {
	adapter := NewEContentAdapter()
	mockResource := &MockEResource{}
	mockObject := &MockEObject{}
	mockObjectAdapters := new(MockEList)

	// Set target
	mockResource.On("GetContents").Return(NewImmutableEList([]interface{}{mockObject})).Once()
	mockObject.On("EAdapters").Return(mockObjectAdapters).Once()
	mockObjectAdapters.On("Contains", adapter).Return(false).Once()
	mockObjectAdapters.On("Add", adapter).Return(true).Once()
	adapter.SetTarget(mockResource)
	mock.AssertExpectationsForObjects(t, mockResource, mockObject, mockObjectAdapters)

	// UnSet target
	mockResource.On("GetContents").Return(NewImmutableEList([]interface{}{mockObject})).Once()
	mockObject.On("EAdapters").Return(mockObjectAdapters).Once()
	mockObjectAdapters.On("Remove", adapter).Return(true).Once()
	adapter.UnSetTarget(mockResource)
	mock.AssertExpectationsForObjects(t, mockResource, mockObject, mockObjectAdapters)
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

	mockReference.On("IsContainment").Once().Return(true)
	mockNotification.On("GetNotifier").Once().Return(mockObject)
	mockNotification.On("GetEventType").Once().Return(UNSET)
	mockNotification.On("GetFeature").Once().Return(mockReference)
	mockNotification.On("GetOldValue").Once().Return(nil)
	mockNotification.On("GetNewValue").Once().Return(nil)
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

	mockReference.On("IsContainment").Once().Return(true)
	mockNotification.On("GetNotifier").Once().Return(mockObject)
	mockNotification.On("GetEventType").Once().Return(SET)
	mockNotification.On("GetFeature").Once().Return(mockReference)
	mockNotification.On("GetOldValue").Once().Return(nil)
	mockNotification.On("GetNewValue").Once().Return(nil)
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

	mockReference.On("IsContainment").Once().Return(true)
	mockNotification.On("GetNotifier").Once().Return(mockObject)
	mockNotification.On("GetEventType").Once().Return(ADD)
	mockNotification.On("GetFeature").Once().Return(mockReference)
	mockNotification.On("GetNewValue").Once().Return(nil)
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

	mockReference.On("IsContainment").Once().Return(true)
	mockNotification.On("GetNotifier").Once().Return(mockObject)
	mockNotification.On("GetEventType").Once().Return(ADD_MANY)
	mockNotification.On("GetFeature").Once().Return(mockReference)
	mockNotification.On("GetNewValue").Once().Return(nil)
	adapter.NotifyChanged(mockNotification)
	mock.AssertExpectationsForObjects(t, mockNotification, mockObject, mockReference)

	mockReference.On("IsContainment").Once().Return(true)
	mockNotification.On("GetNotifier").Once().Return(mockObject)
	mockNotification.On("GetEventType").Once().Return(ADD_MANY)
	mockNotification.On("GetFeature").Once().Return(mockReference)
	mockNotification.On("GetNewValue").Once().Return(struct{}{})
	adapter.NotifyChanged(mockNotification)
	mock.AssertExpectationsForObjects(t, mockNotification, mockObject, mockReference)

	mockReference.On("IsContainment").Once().Return(true)
	mockNotification.On("GetNotifier").Once().Return(mockObject)
	mockNotification.On("GetEventType").Once().Return(ADD_MANY)
	mockNotification.On("GetFeature").Once().Return(mockReference)
	mockNotification.On("GetNewValue").Once().Return([]interface{}{struct{}{}, nil})
	adapter.NotifyChanged(mockNotification)
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

func TestEContentAdapterIntegration(t *testing.T) {
	mockAdapter := NewMockContentAdapter()

	rs := NewEResourceSetImpl()
	rs.EAdapters().Add(mockAdapter)

	// add a new resource to resource set & check that mockAdapter is called
	r := NewEResourceImpl()
	mockAdapter.On("NotifyChanged", mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == rs &&
			n.GetFeatureID() == RESOURCE_SET__RESOURCES &&
			n.GetNewValue() == r &&
			n.GetEventType() == ADD &&
			n.GetPosition() == 0
	})).Once()
	rs.GetResources().Add(r)
	mock.AssertExpectationsForObjects(t, mockAdapter)

	// add a new object to resource & check that mockAdapter is called
	o := NewEObjectImpl()
	mockAdapter.On("NotifyChanged", mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == r &&
			n.GetFeatureID() == RESOURCE__IS_LOADED &&
			n.GetNewValue() == true &&
			n.GetOldValue() == false &&
			n.GetEventType() == SET &&
			n.GetPosition() == -1
	})).Once()
	mockAdapter.On("NotifyChanged", mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == r &&
			n.GetFeatureID() == RESOURCE__CONTENTS &&
			n.GetNewValue() == o &&
			n.GetEventType() == ADD &&
			n.GetPosition() == 0
	})).Once()
	r.GetContents().Add(o)
	mock.AssertExpectationsForObjects(t, mockAdapter)

	// remove object from resource & check that mockAdpater is called
	mockAdapter.On("NotifyChanged", mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == r &&
			n.GetFeatureID() == RESOURCE__IS_LOADED &&
			n.GetNewValue() == false &&
			n.GetOldValue() == true &&
			n.GetEventType() == SET &&
			n.GetPosition() == -1
	})).Once()
	mockAdapter.On("NotifyChanged", mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == r &&
			n.GetFeatureID() == RESOURCE__CONTENTS &&
			n.GetOldValue() == o &&
			n.GetEventType() == REMOVE &&
			n.GetPosition() == 0
	})).Once()
	mockAdapter.On("NotifyChanged", mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == o &&
			n.GetOldValue() == mockAdapter &&
			n.GetEventType() == REMOVING_ADAPTER &&
			n.GetPosition() == 0
	})).Once()
	r.GetContents().Remove(o)
	mock.AssertExpectationsForObjects(t, mockAdapter)

	// remove resource from resource set & check that mockAdapter is called & correctly removed
	mockAdapter.On("NotifyChanged", mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == r &&
			n.GetFeatureID() == RESOURCE__RESOURCE_SET &&
			n.GetOldValue() == rs &&
			n.GetNewValue() == nil &&
			n.GetEventType() == SET &&
			n.GetPosition() == -1
	})).Once()
	mockAdapter.On("NotifyChanged", mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == rs &&
			n.GetFeatureID() == RESOURCE_SET__RESOURCES &&
			n.GetOldValue() == r &&
			n.GetEventType() == REMOVE &&
			n.GetPosition() == 0
	})).Once()
	mockAdapter.On("NotifyChanged", mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == r &&
			n.GetOldValue() == mockAdapter &&
			n.GetEventType() == REMOVING_ADAPTER &&
			n.GetPosition() == 0
	})).Once()
	rs.GetResources().Remove(r)
	mock.AssertExpectationsForObjects(t, mockAdapter)
}

type testContentAdapter struct {
	EContentAdapter
}

func BenchmarkEContentAdapterWithBigModel(b *testing.B) {
	for n := 0; n < b.N; n++ {
		// load package
		ePackage := loadPackage("library.complex.ecore")
		assert.NotNil(b, ePackage)

		// load resource
		xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
		eResource := xmlProcessor.Load(&URI{Path: "testdata/library.complex.big.xml"})
		require.NotNil(b, eResource)
		assert.True(b, eResource.IsLoaded())
		assert.True(b, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
		assert.True(b, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

		adapter := new(testContentAdapter)
		adapter.SetInterfaces(adapter)

		eResource.EAdapters().Add(adapter)
		eResource.EAdapters().Remove(adapter)
	}
}

func BenchmarkEContentAdapterWithTreeModel(b *testing.B) {
	for n := 0; n < b.N; n++ {
		// load package
		ePackage := loadPackage("tree.ecore")
		assert.NotNil(b, ePackage)

		// load resource
		xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
		eResource := xmlProcessor.Load(&URI{Path: "testdata/tree.xml"})
		require.NotNil(b, eResource)
		assert.True(b, eResource.IsLoaded())
		assert.True(b, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
		assert.True(b, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

		adapter := new(testContentAdapter)
		adapter.SetInterfaces(adapter)

		eResource.EAdapters().Add(adapter)
		eResource.EAdapters().Remove(adapter)
	}
}
