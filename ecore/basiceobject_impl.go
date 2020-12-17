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
	adapters           *notifierAdapterList
	resource           EResource
	container          EObject
	containerFeatureID int
	properties         *basicEObjectImplProperties
}

func NewBasicEObjectImpl() *BasicEObjectImpl {
	o := new(BasicEObjectImpl)
	o.deliver = true
	o.containerFeatureID = -1
	return o
}

func (o *BasicEObjectImpl) getObjectProperties() *basicEObjectImplProperties {
	if o.properties == nil {
		o.properties = new(basicEObjectImplProperties)
	}
	return o.properties
}

// EIsProxy ...
func (o *BasicEObjectImpl) EIsProxy() bool {
	return o.getObjectProperties().proxyURI != nil
}

// EProxyURI ...
func (o *BasicEObjectImpl) EProxyURI() *url.URL {
	return o.getObjectProperties().proxyURI
}

// ESetProxyURI ...
func (o *BasicEObjectImpl) ESetProxyURI(uri *url.URL) {
	o.getObjectProperties().proxyURI = uri
}

// EContents ...
func (o *BasicEObjectImpl) EContents() EList {
	if o.contents == nil {
		o.contents = newContentsListAdapter(o, func(eClass EClass) EList { return eClass.GetEContainmentFeatures() })
	}
	return o.contents.GetList()
}

// ECrossReferences ...
func (o *BasicEObjectImpl) ECrossReferences() EList {
	if o.crossReferenceS == nil {
		o.crossReferenceS = newContentsListAdapter(o, func(eClass EClass) EList { return eClass.GetECrossReferenceFeatures() })
	}
	return o.crossReferenceS.GetList()
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
