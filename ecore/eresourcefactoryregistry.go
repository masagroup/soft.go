package ecore

import "net/url"

const (
	DEFAULT_EXTENSION = "*"
)

//EResourceFactoryRegistry ...
type EResourceFactoryRegistry interface {
	GetFactory(uri *url.URL) EResourceFactory
	GetProtocolToFactoryMap() map[string]EResourceFactory
	GetExtensionToFactoryMap() map[string]EResourceFactory
}

var resourceFactoryRegistryInstance EResourceFactoryRegistry

func GetResourceFactoryRegistry() EResourceFactoryRegistry {
	if resourceFactoryRegistryInstance == nil {
		resourceFactoryRegistryInstance = NewEResourceFactoryRegistryImpl()
		// initialize with default factories
		extensionToFactories := resourceFactoryRegistryInstance.GetExtensionToFactoryMap()
		extensionToFactories["ecore"] = &XMIResourceFactory{}
		extensionToFactories["xml"] = &XMLResourceFactory{}
	}
	return resourceFactoryRegistryInstance
}
