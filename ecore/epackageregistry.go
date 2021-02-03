package ecore

type EPackageRegistry interface {
	PutPackage(nsURI string, pack EPackage)
	PutSupplier(nsURI string, supplier func() EPackage)
	Remove(nsURI string)

	RegisterPackage(pack EPackage)
	UnregisterPackage(pack EPackage)

	GetPackage(nsURI string) EPackage
	GetFactory(nsURI string) EFactory
}

var packageRegistryInstance EPackageRegistry

func GetPackageRegistry() EPackageRegistry {
	if packageRegistryInstance == nil {
		packageRegistryInstance = NewEPackageRegistryImpl()
		packageRegistryInstance.RegisterPackage(GetPackage())
	}
	return packageRegistryInstance
}
