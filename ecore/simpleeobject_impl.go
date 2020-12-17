package ecore

import "net/url"

type simpleEObjectImplProperties struct {
	resource        EResource
	proxyURI        *url.URL
	contents        *contentsListAdapter
	crossReferenceS *contentsListAdapter
}

type SimpleEObjectImpl struct {
	BasicEObject
	deliver            bool
	adapters           *basicNotifierAdapterList
	container          EObject
	containerFeatureID int
	properties         *simpleEObjectImplProperties
}
