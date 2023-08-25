// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

type basicEObjectImplProperties struct {
	proxyURI        *URI
	contents        *eContentsList
	crossReferences *eContentsList
}

type BasicEObjectImpl struct {
	AbstractEObject
	adapters           EList
	resource           EResource
	container          EObject
	properties         *basicEObjectImplProperties
	containerFeatureID int
	deliver            bool
	proxy              bool
}

func (o *BasicEObjectImpl) Initialize() {
	o.deliver = true
	o.containerFeatureID = -1
}

func (o *BasicEObjectImpl) getObjectProperties() *basicEObjectImplProperties {
	if o.properties == nil {
		o.properties = new(basicEObjectImplProperties)
	}
	return o.properties
}

func (o *BasicEObjectImpl) EDeliver() bool {
	return o.deliver
}

func (o *BasicEObjectImpl) ESetDeliver(deliver bool) {
	o.deliver = deliver
}

func (o *BasicEObjectImpl) EAdapters() EList {
	if o.adapters == nil {
		o.adapters = newNotifierAdapterList(&o.AbstractENotifier)
	}
	return o.adapters
}

func (o *BasicEObjectImpl) EBasicHasAdapters() bool {
	return o.adapters != nil
}

func (o *BasicEObjectImpl) EBasicAdapters() EList {
	return o.adapters
}

// EIsProxy ...
func (o *BasicEObjectImpl) EIsProxy() bool {
	return o.proxy
}

// EProxyURI ...
func (o *BasicEObjectImpl) EProxyURI() *URI {
	if o.proxy {
		return o.getObjectProperties().proxyURI
	} else {
		return nil
	}
}

// ESetProxyURI ...
func (o *BasicEObjectImpl) ESetProxyURI(uri *URI) {
	o.proxy = uri != nil
	o.getObjectProperties().proxyURI = uri
}

// EContents ...
func (o *BasicEObjectImpl) EContents() EList {
	properties := o.getObjectProperties()
	if properties.contents == nil {
		eObject := o.AsEObject()
		properties.contents = newEContentsList(eObject, eObject.EClass().GetEAllContainments(), true)
	}
	return properties.contents
}

// ECrossReferences ...
func (o *BasicEObjectImpl) ECrossReferences() EList {
	properties := o.getObjectProperties()
	if properties.crossReferences == nil {
		eObject := o.AsEObject()
		properties.crossReferences = newEContentsList(eObject, eObject.EClass().GetECrossReferenceFeatures(), true)
	}
	return properties.crossReferences
}

// ESetContainer ...
func (o *BasicEObjectImpl) ESetInternalContainer(newContainer EObject, newContainerFeatureID int) {
	o.container = newContainer
	o.containerFeatureID = newContainerFeatureID
}

func (o *BasicEObjectImpl) EInternalContainer() EObject {
	return o.container
}

func (o *BasicEObjectImpl) EInternalContainerFeatureID() int {
	return o.containerFeatureID
}

// EInternalResource ...
func (o *BasicEObjectImpl) EInternalResource() EResource {
	return o.resource
}

// ESetInternalResource ...
func (o *BasicEObjectImpl) ESetInternalResource(resource EResource) {
	o.resource = resource
}
