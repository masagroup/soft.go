package ecore

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBasicEObjectGetInterfaces(t *testing.T) {
	o := NewBasicEObject()
	assert.Equal(t, o, o.GetInterfaces())
}

func TestBasicEObjectGetEObject(t *testing.T) {
	o := NewBasicEObject()
	assert.Equal(t, o, o.AsEObject())
}

func TestBasicEObjectEClass(t *testing.T) {
	o := NewBasicEObject()
	assert.Equal(t, GetPackage().GetEObject(), o.EClass())
}

func TestBasicEObjectEIsProxy(t *testing.T) {
	o := NewBasicEObject()
	assert.False(t, o.EIsProxy())
	o.ESetProxyURI(&url.URL{})
	assert.True(t, o.EIsProxy())
}

func TestBasicEObjectContainer(t *testing.T) {
	// set the container
	o := NewBasicEObject()
	mockObject := new(MockEObject)
	mockResource := new(MockEResource)
	mockObject.On("EResource").Return(mockResource)
	mockObject.On("EIsProxy").Return(false)
	mockResource.On("Attached", o)
	mockNotifications := new(MockENotificationChain)
	assert.Equal(t, mockNotifications, o.EBasicSetContainer(mockObject, 1, mockNotifications))
	assert.Equal(t, mockObject, o.EContainer())
	assert.Equal(t, 1, o.EContainerFeatureID())
}

func TestBasicEObjectEBasicRemoveFromContainer(t *testing.T) {
	var o EObject = nil
	i, _ := o.(EObjectInternal)
	assert.Nil(t, i)
}

func TestBasicEObjectESetResource(t *testing.T) {
	// no container
	o := NewBasicEObject()
	mockResource := new(MockEResource)
	mockNotifications := new(MockENotificationChain)
	o.ESetResource(mockResource, mockNotifications)
	mock.AssertExpectationsForObjects(t, mockResource, mockNotifications)

	mockResource2 := new(MockEResource)
	mockContents := new(MockENotifyingList)
	mockResource.On("GetContents").Return(mockContents).Once()
	mockResource.On("Detached", o).Once()
	mockContents.On("RemoveWithNotification", o, mockNotifications).Return(mockNotifications).Once()
	o.ESetResource(mockResource2, mockNotifications)
	mock.AssertExpectationsForObjects(t, mockResource, mockResource2, mockNotifications)

	// container - tested with reflective object
}
