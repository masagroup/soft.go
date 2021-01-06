package ecore

import "net/url"

type basicEObjectImplProperties struct {
	proxyURI        *url.URL
	contents        *contentsListAdapter
	crossReferenceS *contentsListAdapter
}

type BasicEObjectImpl struct {
	AbstractEObject
	deliver            bool
	proxy              bool
	adapters           *notifierAdapterList
	resource           EResource
	container          EObject
	containerFeatureID int
	properties         *basicEObjectImplProperties
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
func (o *BasicEObjectImpl) EProxyURI() *url.URL {
	if o.proxy {
		return o.getObjectProperties().proxyURI
	} else {
		return nil
	}
}

// ESetProxyURI ...
func (o *BasicEObjectImpl) ESetProxyURI(uri *url.URL) {
	o.proxy = uri != nil
	o.getObjectProperties().proxyURI = uri
}

// EContents ...
func (o *BasicEObjectImpl) EContents() EList {
	if o.getObjectProperties().contents == nil {
		o.getObjectProperties().contents = newContentsListAdapter(&o.AbstractEObject, func(eClass EClass) EList { return eClass.GetEContainmentFeatures() })
	}
	return o.getObjectProperties().contents.GetList()
}

// ECrossReferences ...
func (o *BasicEObjectImpl) ECrossReferences() EList {
	if o.getObjectProperties().crossReferenceS == nil {
		o.getObjectProperties().crossReferenceS = newContentsListAdapter(&o.AbstractEObject, func(eClass EClass) EList { return eClass.GetECrossReferenceFeatures() })
	}
	return o.getObjectProperties().crossReferenceS.GetList()
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
