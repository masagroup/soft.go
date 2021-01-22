package ecore

type EPackageRegistry interface {
	RegisterPackage(pack EPackage)
	RegisterPackageWithURI(pack EPackage, nsURI string)
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
