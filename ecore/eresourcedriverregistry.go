package ecore

type EResourceDriverRegistry interface {
	GetDriver(uri *URI) EResourceDriver
	GetProtocolToDriverMap() map[string]EResourceDriver
	GetExtensionToDriverMap() map[string]EResourceDriver
}

var resourceDriverRegistryInstance EResourceDriverRegistry

func GetResourceDriverRegistry() EResourceDriverRegistry {
	if resourceDriverRegistryInstance == nil {
		resourceDriverRegistryInstance = NewEResourceDriverRegistryImpl()
		// initialize with default factories
		extensionToDrivers := resourceDriverRegistryInstance.GetExtensionToDriverMap()
		extensionToDrivers["ecore"] = &XMIDriver{}
		extensionToDrivers["xml"] = &XMLDriver{}
	}
	return resourceDriverRegistryInstance
}
