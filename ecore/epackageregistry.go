// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

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
