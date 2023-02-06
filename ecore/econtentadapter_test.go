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
	children := []any{}
	for i := 0; i < nb; i++ {
		mockObject := NewMockEObject(t)
		mockAdapters := NewMockEList(t)
		mockLists = append(mockLists, mockAdapters)
		mockObjects = append(mockObjects, mockObject)
		children = append(children, mockObject)
	}
	mockChildren := NewImmutableEList(children)
	mockObject := NewMockEObject(t)

	// set adapter target -> this should recursively register adapter on all object children
	for i := 0; i < nb; i++ {
		mockObject := mockObjects[i]
		mockAdapters := mockLists[i]
		if i%2 == 0 {
			mockObject.EXPECT().EIsProxy().Return(false).Once()
			mockObject.EXPECT().EAdapters().Return(mockAdapters).Once()
			mockAdapters.EXPECT().Contains(adapter).Return(false).Once()
			mockAdapters.EXPECT().Add(adapter).Return(true).Once()
		} else {
			mockObject.EXPECT().EIsProxy().Return(true).Once()
		}
	}
	mockObject.EXPECT().EContents().Return(mockChildren).Once()
	adapter.SetTarget(mockObject)
	mock.AssertExpectationsForObjects(t, children...)
	mock.AssertExpectationsForObjects(t, mockObject)

	// unset adapter target -> this should recursively unregister adapter on all object children
	for i := 0; i < nb; i++ {
		mockObject := mockObjects[i]
		mockAdapters := mockLists[i]
		mockObject.EXPECT().EAdapters().Return(mockAdapters).Once()
		mockAdapters.EXPECT().Remove(adapter).Return(true).Once()
	}
	mockObject.EXPECT().EContents().Return(mockChildren).Once()
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
	children := []any{}
	for i := 0; i < nb; i++ {
		mockObject := NewMockEObject(t)
		mockAdapters := NewMockEList(t)
		mockLists = append(mockLists, mockAdapters)
		mockObjects = append(mockObjects, mockObject)
		children = append(children, mockObject)
	}
	mockChildren := NewImmutableEList(children)
	mockObject := NewMockEObject(t)

	// set adapter target -> this should recursively register adapter on all object children
	for i := 0; i < nb; i++ {
		mockObject := mockObjects[i]
		mockAdapters := mockLists[i]
		mockObject.EXPECT().EIsProxy().Return(false).Once()
		mockObject.EXPECT().EAdapters().Return(mockAdapters).Once()
		mockAdapters.EXPECT().Contains(adapter).Return(false).Once()
		mockAdapters.EXPECT().Add(adapter).Return(true).Once()
	}
	mockObject.EXPECT().EContents().Return(mockChildren).Once()
	adapter.SetTarget(mockObject)
	mock.AssertExpectationsForObjects(t, children...)
	mock.AssertExpectationsForObjects(t, mockObject)

	// unset adapter target -> this should recursively unregister adapter on all object children
	for i := 0; i < nb; i++ {
		mockObject := mockObjects[i]
		mockAdapters := mockLists[i]
		mockObject.EXPECT().EAdapters().Return(mockAdapters).Once()
		mockAdapters.EXPECT().Remove(adapter).Return(true).Once()
	}
	mockObject.EXPECT().EContents().Return(mockChildren).Once()
	adapter.UnSetTarget(mockObject)
	mock.AssertExpectationsForObjects(t, children...)
	mock.AssertExpectationsForObjects(t, mockObject)
}

func TestEContentAdapter_SetTarget_EResourceSet(t *testing.T) {
	adapter := NewEContentAdapter()
	mockResourceSet := NewMockEResourceSet(t)
	mockResource := NewMockEResource(t)
	mockResourceAdapters := NewMockEList(t)

	// Set target
	mockResourceSet.EXPECT().GetResources().Return(NewImmutableEList([]any{mockResource})).Once()
	mockResource.EXPECT().EAdapters().Return(mockResourceAdapters).Once()
	mockResourceAdapters.EXPECT().Contains(adapter).Return(false).Once()
	mockResourceAdapters.EXPECT().Add(adapter).Return(true).Once()
	adapter.SetTarget(mockResourceSet)
	mock.AssertExpectationsForObjects(t, mockResourceSet, mockResource, mockResourceAdapters)

	// UnSet target
	mockResourceSet.EXPECT().GetResources().Return(NewImmutableEList([]any{mockResource})).Once()
	mockResource.EXPECT().EAdapters().Return(mockResourceAdapters).Once()
	mockResourceAdapters.EXPECT().Remove(adapter).Return(true).Once()
	adapter.UnSetTarget(mockResourceSet)
	mock.AssertExpectationsForObjects(t, mockResourceSet, mockResource, mockResourceAdapters)
}

func TestEContentAdapter_SetTarget_EResource(t *testing.T) {
	adapter := NewEContentAdapter()
	mockResource := NewMockEResource(t)
	mockObject := NewMockEObject(t)
	mockObjectAdapters := NewMockEList(t)

	// Set target
	mockResource.EXPECT().GetContents().Return(NewImmutableEList([]any{mockObject})).Once()
	mockObject.EXPECT().EAdapters().Return(mockObjectAdapters).Once()
	mockObjectAdapters.EXPECT().Contains(adapter).Return(false).Once()
	mockObjectAdapters.EXPECT().Add(adapter).Return(true).Once()
	adapter.SetTarget(mockResource)
	mock.AssertExpectationsForObjects(t, mockResource, mockObject, mockObjectAdapters)

	// UnSet target
	mockResource.EXPECT().GetContents().Return(NewImmutableEList([]any{mockObject})).Once()
	mockObject.EXPECT().EAdapters().Return(mockObjectAdapters).Once()
	mockObjectAdapters.EXPECT().Remove(adapter).Return(true).Once()
	adapter.UnSetTarget(mockResource)
	mock.AssertExpectationsForObjects(t, mockResource, mockObject, mockObjectAdapters)
}

func TestEContentAdapterNotifyChanged(t *testing.T) {
	adapter := NewEContentAdapter()
	mockNotification := NewMockENotification(t)
	mockObject := NewMockEObject(t)

	mockAttribute := NewMockEAttribute(t)
	mockNotification.EXPECT().GetNotifier().Once().Return(mockObject)
	mockNotification.EXPECT().GetFeature().Once().Return(mockAttribute)
	adapter.NotifyChanged(mockNotification)

	mockReference := NewMockEReference(t)
	mockReference.EXPECT().IsContainment().Once().Return(false)
	mockNotification.EXPECT().GetNotifier().Once().Return(mockObject)
	mockNotification.EXPECT().GetFeature().Once().Return(mockReference)
	adapter.NotifyChanged(mockNotification)

	mock.AssertExpectationsForObjects(t, mockNotification, mockObject, mockAttribute, mockReference)
}

func TestEContentAdapterNotifyChanged_Resolve(t *testing.T) {
	adapter := NewEContentAdapter()
	mockNotification := NewMockENotification(t)
	mockObject := NewMockEObject(t)
	mockReference := NewMockEReference(t)
	mockOldObject := NewMockEObject(t)
	mockOldAdapters := NewMockEList(t)

	mockReference.EXPECT().IsContainment().Once().Return(true)
	mockOldObject.EXPECT().EAdapters().Once().Return(mockOldAdapters)
	mockOldAdapters.EXPECT().Contains(adapter).Once().Return(false)
	mockNotification.EXPECT().GetNotifier().Once().Return(mockObject)
	mockNotification.EXPECT().GetEventType().Once().Return(RESOLVE)
	mockNotification.EXPECT().GetFeature().Once().Return(mockReference)
	mockNotification.EXPECT().GetOldValue().Once().Return(mockOldObject)
	adapter.NotifyChanged(mockNotification)
	mock.AssertExpectationsForObjects(t, mockNotification, mockObject, mockReference, mockOldObject, mockOldAdapters)
}

func TestEContentAdapterNotifyChanged_Resolve_Contains(t *testing.T) {
	adapter := NewEContentAdapter()
	mockNotification := NewMockENotification(t)
	mockObject := NewMockEObject(t)
	mockReference := NewMockEReference(t)
	mockOldObject := NewMockEObject(t)
	mockOldAdapters := NewMockEList(t)
	mockNewObject := NewMockEObject(t)
	mockNewAdapters := NewMockEList(t)

	mockReference.EXPECT().IsContainment().Once().Return(true)
	mockOldObject.EXPECT().EAdapters().Twice().Return(mockOldAdapters)
	mockOldAdapters.EXPECT().Contains(adapter).Once().Return(true)
	mockOldAdapters.EXPECT().Remove(adapter).Once().Return(true)
	mockNewObject.EXPECT().EAdapters().Once().Return(mockNewAdapters)
	mockNewAdapters.EXPECT().Contains(adapter).Once().Return(false)
	mockNewAdapters.EXPECT().Add(adapter).Once().Return(true)
	mockNotification.EXPECT().GetNotifier().Once().Return(mockObject)
	mockNotification.EXPECT().GetEventType().Once().Return(RESOLVE)
	mockNotification.EXPECT().GetFeature().Once().Return(mockReference)
	mockNotification.EXPECT().GetOldValue().Once().Return(mockOldObject)
	mockNotification.EXPECT().GetNewValue().Once().Return(mockNewObject)
	adapter.NotifyChanged(mockNotification)
	mock.AssertExpectationsForObjects(t, mockNotification, mockObject, mockReference, mockOldObject, mockOldAdapters, mockNewObject, mockNewAdapters)
}

func TestEContentAdapterNotifyChanged_UnSet(t *testing.T) {

	adapter := NewEContentAdapter()
	mockNotification := NewMockENotification(t)
	mockObject := NewMockEObjectInternal(t)
	mockReference := NewMockEReference(t)
	mockOldObject := NewMockEObjectInternal(t)
	mockOldAdapters := NewMockEList(t)
	mockNewObject := NewMockEObjectInternal(t)
	mockNewAdapters := NewMockEList(t)

	mockReference.EXPECT().IsContainment().Once().Return(true)
	mockOldObject.EXPECT().EAdapters().Return(mockOldAdapters)
	mockOldObject.EXPECT().EInternalResource().Return(nil)
	mockOldAdapters.EXPECT().Remove(adapter).Once().Return(true)
	mockNewObject.EXPECT().EAdapters().Once().Return(mockNewAdapters)
	mockNewAdapters.EXPECT().Contains(adapter).Once().Return(false)
	mockNewAdapters.EXPECT().Add(adapter).Once().Return(true)
	mockNotification.EXPECT().GetNotifier().Once().Return(mockObject)
	mockNotification.EXPECT().GetEventType().Once().Return(UNSET)
	mockNotification.EXPECT().GetFeature().Once().Return(mockReference)
	mockNotification.EXPECT().GetOldValue().Once().Return(mockOldObject)
	mockNotification.EXPECT().GetNewValue().Once().Return(mockNewObject)
	adapter.NotifyChanged(mockNotification)
	mock.AssertExpectationsForObjects(t, mockNotification, mockObject, mockReference, mockOldObject, mockOldAdapters, mockNewObject, mockNewAdapters)

	mockReference.EXPECT().IsContainment().Once().Return(true)
	mockNotification.EXPECT().GetNotifier().Once().Return(mockObject)
	mockNotification.EXPECT().GetEventType().Once().Return(UNSET)
	mockNotification.EXPECT().GetFeature().Once().Return(mockReference)
	mockNotification.EXPECT().GetOldValue().Once().Return(nil)
	mockNotification.EXPECT().GetNewValue().Once().Return(nil)
	adapter.NotifyChanged(mockNotification)
	mock.AssertExpectationsForObjects(t, mockNotification, mockObject, mockReference, mockOldObject, mockOldAdapters, mockNewObject, mockNewAdapters)

}

func TestEContentAdapterNotifyChanged_Set(t *testing.T) {

	adapter := NewEContentAdapter()
	mockNotification := NewMockENotification(t)
	mockObject := NewMockEObjectInternal(t)
	mockReference := NewMockEReference(t)
	mockOldObject := NewMockEObjectInternal(t)
	mockOldAdapters := NewMockEList(t)
	mockNewObject := NewMockEObjectInternal(t)
	mockNewAdapters := NewMockEList(t)

	mockReference.EXPECT().IsContainment().Once().Return(true)
	mockOldObject.EXPECT().EAdapters().Return(mockOldAdapters)
	mockOldObject.EXPECT().EInternalResource().Return(nil)
	mockOldAdapters.EXPECT().Remove(adapter).Once().Return(true)
	mockNewObject.EXPECT().EAdapters().Once().Return(mockNewAdapters)
	mockNewAdapters.EXPECT().Contains(adapter).Once().Return(false)
	mockNewAdapters.EXPECT().Add(adapter).Once().Return(true)
	mockNotification.EXPECT().GetNotifier().Once().Return(mockObject)
	mockNotification.EXPECT().GetEventType().Once().Return(SET)
	mockNotification.EXPECT().GetFeature().Once().Return(mockReference)
	mockNotification.EXPECT().GetOldValue().Once().Return(mockOldObject)
	mockNotification.EXPECT().GetNewValue().Once().Return(mockNewObject)
	adapter.NotifyChanged(mockNotification)
	mock.AssertExpectationsForObjects(t, mockNotification, mockObject, mockReference, mockOldObject, mockOldAdapters, mockNewObject, mockNewAdapters)

	mockReference.EXPECT().IsContainment().Once().Return(true)
	mockNotification.EXPECT().GetNotifier().Once().Return(mockObject)
	mockNotification.EXPECT().GetEventType().Once().Return(SET)
	mockNotification.EXPECT().GetFeature().Once().Return(mockReference)
	mockNotification.EXPECT().GetOldValue().Once().Return(nil)
	mockNotification.EXPECT().GetNewValue().Once().Return(nil)
	adapter.NotifyChanged(mockNotification)
	mock.AssertExpectationsForObjects(t, mockNotification, mockObject, mockReference, mockOldObject, mockOldAdapters, mockNewObject, mockNewAdapters)

}

func TestEContentAdapterNotifyChanged_Add(t *testing.T) {

	adapter := NewEContentAdapter()
	mockNotification := NewMockENotification(t)
	mockObject := NewMockEObjectInternal(t)
	mockReference := NewMockEReference(t)
	mockNewObject := NewMockEObjectInternal(t)
	mockNewAdapters := NewMockEList(t)

	mockReference.EXPECT().IsContainment().Once().Return(true)
	mockNewObject.EXPECT().EAdapters().Once().Return(mockNewAdapters)
	mockNewAdapters.EXPECT().Contains(adapter).Once().Return(false)
	mockNewAdapters.EXPECT().Add(adapter).Once().Return(true)
	mockNotification.EXPECT().GetNotifier().Once().Return(mockObject)
	mockNotification.EXPECT().GetEventType().Once().Return(ADD)
	mockNotification.EXPECT().GetFeature().Once().Return(mockReference)
	mockNotification.EXPECT().GetNewValue().Once().Return(mockNewObject)
	adapter.NotifyChanged(mockNotification)
	mock.AssertExpectationsForObjects(t, mockNotification, mockObject, mockReference, mockNewObject, mockNewAdapters)

	mockReference.EXPECT().IsContainment().Once().Return(true)
	mockNotification.EXPECT().GetNotifier().Once().Return(mockObject)
	mockNotification.EXPECT().GetEventType().Once().Return(ADD)
	mockNotification.EXPECT().GetFeature().Once().Return(mockReference)
	mockNotification.EXPECT().GetNewValue().Once().Return(nil)
	adapter.NotifyChanged(mockNotification)
	mock.AssertExpectationsForObjects(t, mockNotification, mockObject, mockReference, mockNewObject, mockNewAdapters)

}

func TestEContentAdapterNotifyChanged_AddMany(t *testing.T) {

	adapter := NewEContentAdapter()
	mockNotification := NewMockENotification(t)
	mockObject := NewMockEObjectInternal(t)
	mockReference := NewMockEReference(t)

	mockReference.EXPECT().IsContainment().Once().Return(true)
	nb := rand.Intn(10) + 1
	mockChildren := []any{}
	for i := 0; i < nb; i++ {
		mockObject := NewMockEObject(t)
		mockAdapters := NewMockEList(t)

		mockObject.EXPECT().EAdapters().Return(mockAdapters)
		mockAdapters.EXPECT().Contains(adapter).Return(false)
		mockAdapters.EXPECT().Add(adapter).Return(true)
		mockChildren = append(mockChildren, mockObject)
	}
	mockNotification.EXPECT().GetNotifier().Once().Return(mockObject)
	mockNotification.EXPECT().GetEventType().Once().Return(ADD_MANY)
	mockNotification.EXPECT().GetFeature().Once().Return(mockReference)
	mockNotification.EXPECT().GetNewValue().Once().Return(mockChildren)
	adapter.NotifyChanged(mockNotification)
	mock.AssertExpectationsForObjects(t, mockChildren...)
	mock.AssertExpectationsForObjects(t, mockNotification, mockObject, mockReference)

	mockReference.EXPECT().IsContainment().Once().Return(true)
	mockNotification.EXPECT().GetNotifier().Once().Return(mockObject)
	mockNotification.EXPECT().GetEventType().Once().Return(ADD_MANY)
	mockNotification.EXPECT().GetFeature().Once().Return(mockReference)
	mockNotification.EXPECT().GetNewValue().Once().Return(nil)
	adapter.NotifyChanged(mockNotification)
	mock.AssertExpectationsForObjects(t, mockNotification, mockObject, mockReference)

	mockReference.EXPECT().IsContainment().Once().Return(true)
	mockNotification.EXPECT().GetNotifier().Once().Return(mockObject)
	mockNotification.EXPECT().GetEventType().Once().Return(ADD_MANY)
	mockNotification.EXPECT().GetFeature().Once().Return(mockReference)
	mockNotification.EXPECT().GetNewValue().Once().Return(struct{}{})
	adapter.NotifyChanged(mockNotification)
	mock.AssertExpectationsForObjects(t, mockNotification, mockObject, mockReference)

	mockReference.EXPECT().IsContainment().Once().Return(true)
	mockNotification.EXPECT().GetNotifier().Once().Return(mockObject)
	mockNotification.EXPECT().GetEventType().Once().Return(ADD_MANY)
	mockNotification.EXPECT().GetFeature().Once().Return(mockReference)
	mockNotification.EXPECT().GetNewValue().Once().Return([]any{struct{}{}, nil})
	adapter.NotifyChanged(mockNotification)
	mock.AssertExpectationsForObjects(t, mockNotification, mockObject, mockReference)

}

func TestEContentAdapterNotifyChanged_Remove(t *testing.T) {

	adapter := NewEContentAdapter()
	mockNotification := NewMockENotification(t)
	mockObject := NewMockEObjectInternal(t)
	mockReference := NewMockEReference(t)
	mockOldObject := NewMockEObjectInternal(t)
	mockOldAdapters := NewMockEList(t)

	mockReference.EXPECT().IsContainment().Once().Return(true)

	mockOldObject.EXPECT().EAdapters().Once().Return(mockOldAdapters)
	mockOldObject.EXPECT().EInternalResource().Return(nil)
	mockOldAdapters.EXPECT().Remove(adapter).Once().Return(true)

	mockNotification.EXPECT().GetNotifier().Return(mockObject)
	mockNotification.EXPECT().GetEventType().Return(REMOVE)
	mockNotification.EXPECT().GetFeature().Return(mockReference)
	mockNotification.EXPECT().GetOldValue().Return(mockOldObject)

	adapter.NotifyChanged(mockNotification)

	mock.AssertExpectationsForObjects(t, mockNotification, mockObject, mockReference, mockOldObject, mockOldAdapters)
}

func TestEContentAdapterNotifyChanged_RemoveMany(t *testing.T) {
	adapter := NewEContentAdapter()
	mockNotification := NewMockENotification(t)
	mockObject := NewMockEObjectInternal(t)
	mockReference := NewMockEReference(t)

	mockReference.EXPECT().IsContainment().Once().Return(true)

	nb := rand.Intn(10) + 1
	mockChildren := []any{}
	for i := 0; i < nb; i++ {
		mockObject := NewMockEObjectInternal(t)
		mockAdapters := NewMockEList(t)

		mockObject.EXPECT().EAdapters().Return(mockAdapters)
		mockObject.EXPECT().EInternalResource().Return(nil)
		mockAdapters.EXPECT().Remove(adapter).Return(true)
		mockChildren = append(mockChildren, mockObject)
	}

	mockNotification.EXPECT().GetNotifier().Twice().Return(mockObject)
	mockNotification.EXPECT().GetFeature().Twice().Return(mockReference)
	mockNotification.EXPECT().GetEventType().Once().Return(REMOVE_MANY)
	mockNotification.EXPECT().GetOldValue().Once().Return(mockChildren)

	adapter.NotifyChanged(mockNotification)

	mock.AssertExpectationsForObjects(t, mockChildren...)
	mock.AssertExpectationsForObjects(t, mockNotification, mockObject, mockReference)
}

func TestEContentAdapterIntegration(t *testing.T) {
	mockAdapter := NewMockEContentAdapter(t)

	rs := NewEResourceSetImpl()
	rs.EAdapters().Add(mockAdapter)

	// add a new resource to resource set & check that mockAdapter is called
	r := NewEResourceImpl()
	mockAdapter.EXPECT().NotifyChanged(mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == rs &&
			n.GetFeatureID() == RESOURCE_SET__RESOURCES &&
			n.GetNewValue() == r &&
			n.GetEventType() == ADD &&
			n.GetPosition() == 0
	})).Once()
	rs.GetResources().Add(r)
	mock.AssertExpectationsForObjects(t, mockAdapter)

	// add a new object to resource & check that mockAdapter is called
	o := newEObjectImpl()
	mockAdapter.EXPECT().NotifyChanged(mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == r &&
			n.GetFeatureID() == RESOURCE__IS_LOADED &&
			n.GetNewValue() == true &&
			n.GetOldValue() == false &&
			n.GetEventType() == SET &&
			n.GetPosition() == -1
	})).Once()
	mockAdapter.EXPECT().NotifyChanged(mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == r &&
			n.GetFeatureID() == RESOURCE__CONTENTS &&
			n.GetNewValue() == o &&
			n.GetEventType() == ADD &&
			n.GetPosition() == 0
	})).Once()
	r.GetContents().Add(o)
	mock.AssertExpectationsForObjects(t, mockAdapter)

	// remove object from resource & check that mockAdpater is called
	mockAdapter.EXPECT().NotifyChanged(mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == r &&
			n.GetFeatureID() == RESOURCE__IS_LOADED &&
			n.GetNewValue() == false &&
			n.GetOldValue() == true &&
			n.GetEventType() == SET &&
			n.GetPosition() == -1
	})).Once()
	mockAdapter.EXPECT().NotifyChanged(mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == r &&
			n.GetFeatureID() == RESOURCE__CONTENTS &&
			n.GetOldValue() == o &&
			n.GetEventType() == REMOVE &&
			n.GetPosition() == 0
	})).Once()
	mockAdapter.EXPECT().NotifyChanged(mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == o &&
			n.GetOldValue() == mockAdapter &&
			n.GetEventType() == REMOVING_ADAPTER &&
			n.GetPosition() == 0
	})).Once()
	r.GetContents().Remove(o)
	mock.AssertExpectationsForObjects(t, mockAdapter)

	// remove resource from resource set & check that mockAdapter is called & correctly removed
	mockAdapter.EXPECT().NotifyChanged(mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == r &&
			n.GetFeatureID() == RESOURCE__RESOURCE_SET &&
			n.GetOldValue() == rs &&
			n.GetNewValue() == nil &&
			n.GetEventType() == SET &&
			n.GetPosition() == -1
	})).Once()
	mockAdapter.EXPECT().NotifyChanged(mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == rs &&
			n.GetFeatureID() == RESOURCE_SET__RESOURCES &&
			n.GetOldValue() == r &&
			n.GetEventType() == REMOVE &&
			n.GetPosition() == 0
	})).Once()
	mockAdapter.EXPECT().NotifyChanged(mock.MatchedBy(func(n ENotification) bool {
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
		xmlProcessor := NewXMLProcessor(XMLProcessorPackages([]EPackage{ePackage}))
		eResource := xmlProcessor.Load(NewURI("testdata/library.complex.big.xml"))
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
		xmlProcessor := NewXMLProcessor(XMLProcessorPackages([]EPackage{ePackage}))
		eResource := xmlProcessor.Load(NewURI("testdata/tree.xml"))
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
