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
)

func TestReflectiveEObjectImpl_EClass(t *testing.T) {
	o := NewReflectiveEObjectImpl()
	assert.Equal(t, GetPackage().GetEObject(), o.EClass())

	mockClass := NewMockEClass(t)
	o.SetEClass(mockClass)
	assert.Equal(t, mockClass, o.EClass())

}

func TestReflectiveEObjectImpl_EContainer(t *testing.T) {
	mockClass := NewMockEClass(t)
	mockContainer := NewMockEObjectInternal(t)
	mockResource := NewMockEResource(t)
	mockResourceSet := NewMockEResourceSet(t)

	o := NewReflectiveEObjectImpl()
	o.SetEClass(mockClass)
	o.ESetInternalContainer(mockContainer, 0)
	o.ESetInternalResource(mockResource)

	mockAdapter := NewMockEAdapter(t)
	mockAdapter.EXPECT().SetTarget(o).Once()
	o.EAdapters().Add(mockAdapter)
	mock.AssertExpectationsForObjects(t, mockAdapter)

	// non proxy
	mockContainer.EXPECT().EIsProxy().Return(false).Once()
	assert.Equal(t, mockContainer, o.EContainer())
	mock.AssertExpectationsForObjects(t, mockClass, mockContainer, mockAdapter)

	// proxy
	mockURI, _ := ParseURI("test://file.t")
	mockResolved := NewMockEObjectInternal(t)
	mockResolved.EXPECT().EProxyURI().Return(nil).Once()
	mockResource.EXPECT().GetResourceSet().Return(mockResourceSet).Once()
	mockResourceSet.EXPECT().GetEObject(mockURI, true).Return(mockResolved).Once()
	mockContainer.EXPECT().EIsProxy().Return(true).Once()
	mockContainer.EXPECT().EProxyURI().Return(mockURI).Twice()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(nil).Once()
	mockAdapter.EXPECT().NotifyChanged(mock.MatchedBy(func(notification ENotification) bool {
		return notification.GetEventType() == RESOLVE
	})).Once()
	assert.Equal(t, mockResolved, o.EContainer())
	mock.AssertExpectationsForObjects(t, mockClass, mockContainer, mockResolved, mockResource, mockResourceSet, mockAdapter)
}

func TestReflectiveEObjectImpl_EContainmentFeature(t *testing.T) {
	mockClass := NewMockEClass(t)
	mockContainer := NewMockEObjectInternal(t)
	mockReference := NewMockEReference(t)
	o := NewReflectiveEObjectImpl()
	o.SetEClass(mockClass)

	//
	o.ESetInternalContainer(mockContainer, -2)
	mockContainer.EXPECT().EClass().Return(mockClass)
	mockClass.EXPECT().GetEStructuralFeature(1).Return(mockReference).Once()
	assert.Equal(t, mockReference, o.EContainmentFeature())
	mock.AssertExpectationsForObjects(t, mockClass, mockContainer, mockReference)

	//
	o.ESetInternalContainer(mockContainer, 1)
	mockClass.EXPECT().GetEStructuralFeature(1).Return(mockReference).Once()
	assert.Equal(t, mockReference, o.EContainmentFeature())
	mock.AssertExpectationsForObjects(t, mockClass, mockContainer, mockReference)

	//
	o.ESetInternalContainer(mockContainer, 1)
	mockClass.EXPECT().GetEStructuralFeature(1).Return(nil).Once()
	assert.Panics(t, func() { o.EContainmentFeature() })

	o.ESetInternalContainer(nil, -1)
	assert.Nil(t, o.EContainmentFeature())
}

func TestReflectiveEObjectImpl_ESetResource(t *testing.T) {
	mockClass := NewMockEClass(t)
	mockContainer := NewMockEObjectInternal(t)
	mockResource := NewMockEResource(t)
	mockNotifications := NewMockENotificationChain(t)
	mockReference := NewMockEReference(t)
	o := NewReflectiveEObjectImpl()
	o.SetEClass(mockClass)
	o.ESetInternalContainer(mockContainer, 0)

	// set resource with container feature which is a reference proxy
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockReference).Once()
	mockReference.EXPECT().IsResolveProxies().Return(true).Once()
	mockContainer.EXPECT().EResource().Return(mockResource).Once()
	mockResource.EXPECT().Detached(o).Once()
	assert.Equal(t, mockNotifications, o.ESetResource(mockResource, mockNotifications))
	mock.AssertExpectationsForObjects(t, mockResource, mockContainer, mockReference, mockClass, mockNotifications)

	// reset resource with container feature as a reference proxy
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockReference).Once()
	mockReference.EXPECT().IsResolveProxies().Return(true).Once()
	mockContainer.EXPECT().EResource().Return(mockResource).Once()
	mockResource.EXPECT().Attached(o).Once()
	assert.Equal(t, mockNotifications, o.ESetResource(nil, mockNotifications))
	mock.AssertExpectationsForObjects(t, mockResource, mockContainer, mockReference, mockClass, mockNotifications)

	// set new resource with a container as a non proxy reference
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockReference).Twice()
	mockReference.EXPECT().IsResolveProxies().Return(false).Once()
	mockReference.EXPECT().GetEOpposite().Return(nil).Once()
	mockContainer.EXPECT().EResource().Return(mockResource).Once()
	mockResource.EXPECT().Detached(o).Once()
	assert.Equal(t, mockNotifications, o.ESetResource(mockResource, mockNotifications))
	mock.AssertExpectationsForObjects(t, mockResource, mockContainer, mockReference, mockClass, mockNotifications)
}

func TestReflectiveEObjectImpl_EGetFromID(t *testing.T) {
	mockClass := NewMockEClass(t)
	o := NewReflectiveEObjectImpl()
	o.SetEClass(mockClass)

	mockClass.EXPECT().GetEStructuralFeature(0).Return(nil).Once()
	assert.Panics(t, func() { o.EGetFromID(0, false) })
	mock.AssertExpectationsForObjects(t, mockClass)
}

func TestReflectiveEObjectImpl_ESetFromID(t *testing.T) {
	mockClass := NewMockEClass(t)
	o := NewReflectiveEObjectImpl()
	o.SetEClass(mockClass)

	mockClass.EXPECT().GetEStructuralFeature(0).Return(nil).Once()
	assert.Panics(t, func() { o.ESetFromID(0, false) })
	mock.AssertExpectationsForObjects(t, mockClass)
}

func TestReflectiveEObjectImpl_EUnSetFromID(t *testing.T) {
	mockClass := NewMockEClass(t)
	o := NewReflectiveEObjectImpl()
	o.SetEClass(mockClass)

	mockClass.EXPECT().GetEStructuralFeature(0).Return(nil).Once()
	assert.Panics(t, func() { o.EUnsetFromID(0) })
	mock.AssertExpectationsForObjects(t, mockClass)
}

func TestReflectiveEObjectImpl_GetAttribute(t *testing.T) {
	mockAttribute := NewMockEAttribute(t)
	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().GetDefaultValue().Return(nil).Once()

	mockClass := NewMockEClass(t)
	mockClass.EXPECT().GetFeatureCount().Return(2).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Once()

	o := NewReflectiveEObjectImpl()
	o.SetEClass(mockClass)

	// get unitialized value
	assert.Nil(t, o.EGetFromID(0, true))

	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute)
}

func TestReflectiveEObjectImpl_GetAttribute_Default(t *testing.T) {
	mockDefault := NewMockEObject(t)

	mockAttribute := NewMockEAttribute(t)
	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().GetDefaultValue().Return(mockDefault).Once()

	mockClass := NewMockEClass(t)
	mockClass.EXPECT().GetFeatureCount().Return(2).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Once()

	o := NewReflectiveEObjectImpl()
	o.SetEClass(mockClass)

	// get default value
	assert.Equal(t, mockDefault, o.EGetFromID(0, true))
	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute)
}

func TestReflectiveEObjectImpl_GetAttribute_Many(t *testing.T) {
	mockAttribute := NewMockEAttribute(t)
	mockClass := NewMockEClass(t)
	o := NewReflectiveEObjectImpl()
	o.SetEClass(mockClass)

	// get unitialized value
	mockClass.EXPECT().GetFeatureCount().Return(2).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Once()
	mockAttribute.EXPECT().IsMany().Return(true).Once()
	mockAttribute.EXPECT().GetFeatureID().Return(2).Once()
	mockAttribute.EXPECT().IsUnique().Return(true).Once()
	mockAttribute.EXPECT().GetEType().Return(nil).Once()
	val := o.EGetFromID(0, false)
	assert.NotNil(t, val)
	l, _ := val.(EList)
	assert.NotNil(t, l)

	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute)
}

func TestReflectiveEObjectImpl_GetReference_Many(t *testing.T) {
	mockReference := NewMockEReference(t)
	mockClass := NewMockEClass(t)
	o := NewReflectiveEObjectImpl()
	o.SetEClass(mockClass)

	// get unitialized value
	mockClass.EXPECT().GetFeatureCount().Return(2).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockReference).Once()
	mockReference.EXPECT().IsMany().Return(true).Once()
	mockReference.EXPECT().GetEType().Return(nil).Once()
	mockReference.EXPECT().GetEOpposite().Return(nil).Twice()
	mockReference.EXPECT().IsContainment().Return(false).Once()
	mockReference.EXPECT().GetFeatureID().Return(0).Once()
	mockReference.EXPECT().EIsProxy().Return(false).Once()
	mockReference.EXPECT().IsUnsettable().Return(false).Once()
	val := o.EGetFromID(0, false)
	assert.NotNil(t, val)

	// check its is an object list
	l, _ := val.(EObjectList)
	assert.NotNil(t, l)

	mock.AssertExpectationsForObjects(t, mockClass, mockReference)
}

func TestReflectiveEObjectImpl_SetAttribute(t *testing.T) {
	mockAttribute := NewMockEAttribute(t)

	mockClass := NewMockEClass(t)
	mockClass.EXPECT().GetFeatureCount().Return(2).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Twice()

	mockObject := NewMockEObject(t)

	o := NewReflectiveEObjectImpl()
	o.SetEClass(mockClass)

	// set
	o.ESetFromID(0, mockObject)

	// check that value is well set
	assert.Equal(t, mockObject, o.EGetFromID(0, true))

	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute)
}

func TestReflectiveEObjectImpl_UnsetAttribute(t *testing.T) {
	mockAttribute := NewMockEAttribute(t)
	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().GetDefaultValue().Return(nil).Once()

	mockClass := NewMockEClass(t)
	mockClass.EXPECT().GetFeatureCount().Return(2).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Times(3)

	mockObject := NewMockEObject(t)

	o := NewReflectiveEObjectImpl()
	o.SetEClass(mockClass)

	// set - unset
	o.ESetFromID(0, mockObject)
	o.EUnsetFromID(0)

	// check that value is unset
	assert.Nil(t, o.EGetFromID(0, true))

	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute)
}

func TestReflectiveEObjectImpl_GetContainer(t *testing.T) {
	mockOpposite := NewMockEReference(t)
	mockReference := NewMockEReference(t)
	mockClass := NewMockEClass(t)

	o := NewReflectiveEObjectImpl()
	o.SetEClass(mockClass)

	// get non initialized container
	mockOpposite.EXPECT().IsContainment().Return(true).Once()
	mockReference.EXPECT().GetEOpposite().Return(mockOpposite).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockReference).Once()
	mockClass.EXPECT().GetFeatureID(mockReference).Return(0).Once()
	assert.Nil(t, o.EGetFromID(0, true))
	mock.AssertExpectationsForObjects(t, mockClass, mockReference, mockOpposite)
}

func TestReflectiveEObjectImpl_SetContainer(t *testing.T) {
	mockObject := NewMockEObjectInternal(t)
	mockObjectClass := NewMockEClass(t)
	mockOpposite := NewMockEReference(t)
	mockReference := NewMockEReference(t)
	mockClass := NewMockEClass(t)

	o := NewReflectiveEObjectImpl()
	o.SetEClass(mockClass)

	// set reference as mockObject
	mockObject.EXPECT().EResource().Return(nil).Once()
	mockObject.EXPECT().EInverseAdd(o, 0, nil).Return(nil).Once()
	mockObject.EXPECT().EClass().Return(mockObjectClass).Once()
	mockOpposite.EXPECT().IsContainment().Return(true).Once()
	mockReference.EXPECT().GetEOpposite().Return(mockOpposite).Twice()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockReference).Once()
	mockClass.EXPECT().GetFeatureID(mockReference).Return(0).Once()
	mockObjectClass.EXPECT().GetFeatureID(mockOpposite).Return(0).Once()
	o.ESetFromID(0, mockObject)
	mock.AssertExpectationsForObjects(t, mockClass, mockReference, mockOpposite, mockObject, mockObjectClass)

	// get unresolved
	mockOpposite.EXPECT().IsContainment().Return(true).Once()
	mockReference.EXPECT().GetEOpposite().Return(mockOpposite).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockReference).Once()
	mockClass.EXPECT().GetFeatureID(mockReference).Return(0).Once()
	assert.Equal(t, mockObject, o.EGetFromID(0, false))
	mock.AssertExpectationsForObjects(t, mockClass, mockReference, mockOpposite, mockObject)

	// get resolved
	mockObject.EXPECT().EIsProxy().Return(false).Once()
	mockOpposite.EXPECT().IsContainment().Return(true).Once()
	mockReference.EXPECT().GetEOpposite().Return(mockOpposite).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockReference).Once()
	mockClass.EXPECT().GetFeatureID(mockReference).Return(0).Once()
	assert.Equal(t, mockObject, o.EGetFromID(0, true))
	mock.AssertExpectationsForObjects(t, mockClass, mockReference, mockOpposite, mockObject)

	// set reference as nil
	mockObject.EXPECT().EInverseRemove(o.GetInterfaces(), 0, nil).Return(nil).Once()
	mockObject.EXPECT().EResource().Return(nil).Once()
	mockOpposite.EXPECT().IsContainment().Return(true).Once()
	mockOpposite.EXPECT().GetFeatureID().Return(0).Once()
	mockReference.EXPECT().GetEOpposite().Return(mockOpposite).Twice()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockReference).Twice()
	mockClass.EXPECT().GetFeatureID(mockReference).Return(0).Once()
	o.ESetFromID(0, nil)
	mock.AssertExpectationsForObjects(t, mockClass, mockReference, mockOpposite, mockObject)

	// get unresolved
	mockOpposite.EXPECT().IsContainment().Return(true).Once()
	mockReference.EXPECT().GetEOpposite().Return(mockOpposite).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockReference).Once()
	mockClass.EXPECT().GetFeatureID(mockReference).Return(0).Once()
	assert.Nil(t, o.EGetFromID(0, false))
	mock.AssertExpectationsForObjects(t, mockClass, mockReference, mockOpposite, mockObject)
}

func TestReflectiveEObjectImpl_UnSetContainer(t *testing.T) {

	mockObject := NewMockEObjectInternal(t)
	mockObjectClass := NewMockEClass(t)
	mockOpposite := NewMockEReference(t)
	mockReference := NewMockEReference(t)
	mockClass := NewMockEClass(t)

	o := NewReflectiveEObjectImpl()
	o.SetEClass(mockClass)

	// set reference as mockObject
	mockObject.EXPECT().EResource().Return(nil).Once()
	mockObject.EXPECT().EInverseAdd(o, 0, nil).Return(nil).Once()
	mockObject.EXPECT().EClass().Return(mockObjectClass).Once()
	mockOpposite.EXPECT().IsContainment().Return(true).Once()
	mockReference.EXPECT().GetEOpposite().Return(mockOpposite).Twice()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockReference).Once()
	mockClass.EXPECT().GetFeatureID(mockReference).Return(0).Once()
	mockObjectClass.EXPECT().GetFeatureID(mockOpposite).Return(0).Once()
	o.ESetFromID(0, mockObject)
	mock.AssertExpectationsForObjects(t, mockClass, mockReference, mockOpposite, mockObject, mockObjectClass)

	// unset
	mockObject.EXPECT().EInverseRemove(o.GetInterfaces(), 0, nil).Return(nil).Once()
	mockObject.EXPECT().EResource().Return(nil).Once()
	mockOpposite.EXPECT().IsContainment().Return(true).Once()
	mockOpposite.EXPECT().GetFeatureID().Return(0).Once()
	mockReference.EXPECT().GetEOpposite().Return(mockOpposite).Twice()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockReference).Twice()
	mockClass.EXPECT().GetFeatureID(mockReference).Return(0).Once()
	o.EUnsetFromID(0)

	// get unresolved
	mockOpposite.EXPECT().IsContainment().Return(true).Once()
	mockReference.EXPECT().GetEOpposite().Return(mockOpposite).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockReference).Once()
	mockClass.EXPECT().GetFeatureID(mockReference).Return(0).Once()
	assert.Nil(t, o.EGetFromID(0, false))
	mock.AssertExpectationsForObjects(t, mockClass, mockReference, mockOpposite, mockObject)
}

func TestReflectiveEObjectImpl_GetReferenceProxy(t *testing.T) {
	mockObject := NewMockEObjectInternal(t)
	mockClass := NewMockEClass(t)
	mockReference := NewMockEReference(t)
	mockResource := NewMockEResource(t)
	mockResourceSet := NewMockEResourceSet(t)
	mockURI, _ := ParseURI("test://file.t")
	mockResolved := NewMockEObjectInternal(t)

	o := NewReflectiveEObjectImpl()
	o.SetEClass(mockClass)
	o.ESetInternalResource(mockResource)

	// add listener
	mockAdapter := NewMockEAdapter(t)
	mockAdapter.EXPECT().SetTarget(o).Once()
	o.EAdapters().Add(mockAdapter)
	mock.AssertExpectationsForObjects(t, mockAdapter)

	// simple set
	mockClass.EXPECT().GetFeatureCount().Return(1).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockReference).Once()
	mockReference.EXPECT().GetEOpposite().Return(nil).Twice()
	mockReference.EXPECT().IsContainment().Return(false).Once()
	mockAdapter.EXPECT().NotifyChanged(mock.MatchedBy(func(notification ENotification) bool {
		return notification.GetEventType() == SET && notification.GetOldValue() == nil && notification.GetNewValue() == mockObject
	})).Once()
	o.ESetFromID(0, mockObject)
	mock.AssertExpectationsForObjects(t, mockObject, mockClass, mockReference, mockAdapter)

	// get with resolution - no contains
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockReference).Once()
	mockReference.EXPECT().GetEOpposite().Return(nil).Once()
	mockReference.EXPECT().IsContainment().Return(false).Once()
	mockReference.EXPECT().IsResolveProxies().Return(true).Once()
	mockObject.EXPECT().EIsProxy().Return(true).Once()
	mockResolved.EXPECT().EProxyURI().Return(nil).Once()
	mockResource.EXPECT().GetResourceSet().Return(mockResourceSet).Once()
	mockResourceSet.EXPECT().GetEObject(mockURI, true).Return(mockResolved).Once()
	mockObject.EXPECT().EProxyURI().Return(mockURI).Twice()
	mockAdapter.EXPECT().NotifyChanged(mock.MatchedBy(func(notification ENotification) bool {
		return notification.GetEventType() == RESOLVE
	})).Once()
	assert.Equal(t, mockResolved, o.EGetFromID(0, true))
	mock.AssertExpectationsForObjects(t, mockObject, mockClass, mockReference, mockAdapter, mockResource, mockResourceSet)
}

func TestReflectiveEObjectImpl_GetReferenceProxyContainment(t *testing.T) {
	mockObject := NewMockEObjectInternal(t)
	mockClass := NewMockEClass(t)
	mockReference := NewMockEReference(t)
	mockResource := NewMockEResource(t)
	mockResourceSet := NewMockEResourceSet(t)
	mockURI, _ := ParseURI("test://file.t")
	mockResolved := NewMockEObjectInternal(t)

	o := NewReflectiveEObjectImpl()
	o.SetEClass(mockClass)
	o.ESetInternalResource(mockResource)

	// add listener
	mockAdapter := NewMockEAdapter(t)
	mockAdapter.EXPECT().SetTarget(o).Once()
	o.EAdapters().Add(mockAdapter)
	mock.AssertExpectationsForObjects(t, mockAdapter)

	// simple set
	mockClass.EXPECT().GetFeatureCount().Return(1).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockReference).Once()
	mockReference.EXPECT().GetEOpposite().Return(nil).Twice()
	mockReference.EXPECT().IsContainment().Return(false).Once()
	mockAdapter.EXPECT().NotifyChanged(mock.MatchedBy(func(notification ENotification) bool {
		return notification.GetEventType() == SET && notification.GetOldValue() == nil && notification.GetNewValue() == mockObject
	})).Once()
	o.ESetFromID(0, mockObject)
	mock.AssertExpectationsForObjects(t, mockObject, mockClass, mockReference, mockAdapter)

	// get with resolution and containment
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockReference).Once()
	mockClass.EXPECT().GetFeatureID(mockReference).Return(0).Once()
	mockReference.EXPECT().GetEOpposite().Return(nil).Twice()
	mockReference.EXPECT().IsContainment().Return(true).Once()
	mockReference.EXPECT().IsResolveProxies().Return(true).Once()
	mockObject.EXPECT().EIsProxy().Return(true).Once()
	mockResolved.EXPECT().EProxyURI().Return(nil).Once()
	mockResource.EXPECT().GetResourceSet().Return(mockResourceSet).Once()
	mockResourceSet.EXPECT().GetEObject(mockURI, true).Return(mockResolved).Once()
	mockObject.EXPECT().EProxyURI().Return(mockURI).Twice()
	mockObject.EXPECT().EInverseRemove(o, -1, nil).Return(nil).Once()
	mockResolved.EXPECT().EInverseAdd(o, -1, nil).Return(nil).Once()
	mockAdapter.EXPECT().NotifyChanged(mock.MatchedBy(func(notification ENotification) bool {
		return notification.GetEventType() == RESOLVE
	})).Once()
	assert.Equal(t, mockResolved, o.EGetFromID(0, true))
	mock.AssertExpectationsForObjects(t, mockObject, mockClass, mockReference, mockAdapter, mockResource, mockResourceSet)
}

func TestReflectiveEObjectImpl_GetReferenceProxyContainmentBidirectional(t *testing.T) {
	mockClass := NewMockEClass(t)
	mockObject := NewMockEObjectInternal(t)
	mockObjectClass := NewMockEClass(t)
	mockReference := NewMockEReference(t)
	mockOpposite := NewMockEReference(t)
	mockResource := NewMockEResource(t)
	mockResourceSet := NewMockEResourceSet(t)
	mockURI, _ := ParseURI("test://file.t")
	mockResolved := NewMockEObjectInternal(t)
	mockResolvedClass := NewMockEClass(t)

	o := NewReflectiveEObjectImpl()
	o.SetEClass(mockClass)
	o.ESetInternalResource(mockResource)

	// add listener
	mockAdapter := NewMockEAdapter(t)
	mockAdapter.EXPECT().SetTarget(o).Once()
	o.EAdapters().Add(mockAdapter)
	mock.AssertExpectationsForObjects(t, mockAdapter)

	// simple set
	mockClass.EXPECT().GetFeatureCount().Return(1).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockReference).Once()
	mockReference.EXPECT().GetEOpposite().Return(nil).Twice()
	mockReference.EXPECT().IsContainment().Return(false).Once()
	mockAdapter.EXPECT().NotifyChanged(mock.MatchedBy(func(notification ENotification) bool {
		return notification.GetEventType() == SET && notification.GetOldValue() == nil && notification.GetNewValue() == mockObject
	})).Once()
	o.ESetFromID(0, mockObject)
	mock.AssertExpectationsForObjects(t, mockObject, mockClass, mockReference, mockAdapter)

	// get with resolution and containment
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockReference).Once()
	mockReference.EXPECT().GetEOpposite().Return(mockOpposite).Times(3)
	mockReference.EXPECT().IsContainment().Return(true).Once()
	mockOpposite.EXPECT().IsContainment().Return(false).Once()
	mockReference.EXPECT().IsResolveProxies().Return(true).Once()
	mockObject.EXPECT().EIsProxy().Return(true).Once()
	mockResolved.EXPECT().EProxyURI().Return(nil).Once()
	mockResource.EXPECT().GetResourceSet().Return(mockResourceSet).Once()
	mockResourceSet.EXPECT().GetEObject(mockURI, true).Return(mockResolved).Once()
	mockObject.EXPECT().EProxyURI().Return(mockURI).Twice()
	mockObject.EXPECT().EClass().Return(mockObjectClass).Once()
	mockObjectClass.EXPECT().GetFeatureID(mockOpposite).Return(1).Once()
	mockResolved.EXPECT().EClass().Return(mockResolvedClass).Once()
	mockResolvedClass.EXPECT().GetFeatureID(mockOpposite).Return(2).Once()
	mockObject.EXPECT().EInverseRemove(o, 1, nil).Return(nil).Once()
	mockResolved.EXPECT().EInverseAdd(o, 2, nil).Return(nil).Once()
	mockAdapter.EXPECT().NotifyChanged(mock.MatchedBy(func(notification ENotification) bool {
		return notification.GetEventType() == RESOLVE
	})).Once()
	assert.Equal(t, mockResolved, o.EGetFromID(0, true))
	mock.AssertExpectationsForObjects(t, mockObject, mockClass, mockReference, mockAdapter, mockResource, mockResourceSet)
}

func TestReflectiveEObjectImpl_EContents(t *testing.T) {
	o := NewReflectiveEObjectImpl()
	assert.NotNil(t, o.EContents())
}

func TestReflectiveEObjectImpl_EAllContents(t *testing.T) {
	o := NewReflectiveEObjectImpl()
	assert.NotNil(t, o.EAllContents())
}

func TestReflectiveEObjectImpl_ECrossReferences(t *testing.T) {
	o := NewReflectiveEObjectImpl()
	assert.NotNil(t, o.ECrossReferences())
}
