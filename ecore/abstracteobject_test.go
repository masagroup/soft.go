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
)

// func TestAbstractEObjectGetInterfaces(t *testing.T) {
// 	o := NewAbstractEObject()
// 	assert.Equal(t, o, o.GetInterfaces())
// }

// func TestAbstractEObjectGetEObject(t *testing.T) {
// 	o := NewAbstractEObject()
// 	assert.Equal(t, o, o.AsEObject())
// }

// func TestAbstractEObjectEClass(t *testing.T) {
// 	o := NewAbstractEObject()
// 	assert.Equal(t, GetPackage().GetEObject(), o.EClass())
// }

// func TestAbstractEObjectEIsProxy(t *testing.T) {
// 	o := NewAbstractEObject()
// 	assert.False(t, o.EIsProxy())
// 	o.ESetProxyURI(&URI{})
// 	assert.True(t, o.EIsProxy())
// }

func TestAbstractEObjectContainer(t *testing.T) {
	// set the container
	// o := NewAbstractEObject()
	// mockObject := new(MockEObject)
	// mockResource := new(MockEResource)
	// mockObject.On("EResource").Return(mockResource)
	// mockObject.On("EIsProxy").Return(false)
	// mockResource.On("Attached", o)
	// mockNotifications := new(MockENotificationChain)
	// assert.Equal(t, mockNotifications, o.EBasicSetContainer(mockObject, 1, mockNotifications))
	// assert.Equal(t, mockObject, o.EContainer())
	//assert.Equal(t, 1, o.EContainerFeatureID())
}

// func TestAbstractEObjectEBasicRemoveFromContainer(t *testing.T) {
// 	var o EObject = nil
// 	i, _ := o.(EObjectInternal)
// 	assert.Nil(t, i)
// }

// func TestAbstractEObjectESetResource(t *testing.T) {
// 	// no container
// 	o := NewAbstractEObject()
// 	mockResource := new(MockEResource)
// 	mockNotifications := new(MockENotificationChain)
// 	o.ESetResource(mockResource, mockNotifications)
// 	mock.AssertExpectationsForObjects(t, mockResource, mockNotifications)

// 	mockResource2 := new(MockEResource)
// 	mockContents := new(MockENotifyingList)
// 	mockResource.On("GetContents").Return(mockContents).Once()
// 	mockResource.On("Detached", o).Once()
// 	mockContents.On("RemoveWithNotification", o, mockNotifications).Return(mockNotifications).Once()
// 	o.ESetResource(mockResource2, mockNotifications)
// 	mock.AssertExpectationsForObjects(t, mockResource, mockResource2, mockNotifications)

// 	// container - tested with reflective object
// }

// func TestAbstractEObject_EStaticFeatureCount(t *testing.T) {
// 	o := NewAbstractEObject()
// 	assert.Equal(t, 0, o.EStaticFeatureCount())
// }

// func TestAbstractEObject_EDynamicProperties(t *testing.T) {
// 	o := NewAbstractEObject()
// 	assert.Nil(t, o.EDynamicProperties())
// }
