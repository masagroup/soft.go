// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import "sync"

// EPackageRegistry maps a namespace URI to an EPackage (one package per nsURI).
//
// A program may link several complete generated implementations of the same
// metamodel: same nsURI, same classifier and feature IDs, but distinct EPackage
// instances, each with its own EFactory. The implementation a resource set uses
// is the package registered in that resource set's own registry, either with
// CreateEResourceSet(packages), XMLProcessorPackages(packages) or
// resourceSet.GetPackageRegistry().RegisterPackage(pack). Objects loaded in that
// resource set then all belong to that implementation; resources mixing
// implementations are not supported.
//
// The implementations of a shared package (one EPackage instance generated with
// the interfaces, see EObjectFactory) share the package and its classes: the
// implementation of a resource set is the factory registered with
// PutFactory(nsURI, factory), GetFactory returning it before the factory
// instance of the package.
type EPackageRegistry interface {
	PutPackage(nsURI string, pack EPackage)
	PutSupplier(nsURI string, supplier func() EPackage)
	Remove(nsURI string)

	RegisterPackage(pack EPackage)
	UnregisterPackage(pack EPackage)

	// PutFactory registers factory as the factory of nsURI: GetFactory returns it instead of the
	// factory instance of the package of nsURI. It chooses the implementation of a shared package
	// (see EObjectFactory) the objects of a resource set are created with.
	PutFactory(nsURI string, factory EFactory)
	RemoveFactory(nsURI string)

	GetPackage(nsURI string) EPackage
	GetFactory(nsURI string) EFactory
}

var packageRegistryOnce sync.Once
var packageRegistryInstance EPackageRegistry

// GetPackageRegistry returns the global package registry.
//
// The global registry holds the ecore package only and stays neutral: generated
// packages never register themselves in it. When several implementations of a
// metamodel coexist, register the one a resource set uses in that resource
// set's registry (resourceSet.GetPackageRegistry().RegisterPackage(pack)), not
// in the global one.
func GetPackageRegistry() EPackageRegistry {
	packageRegistryOnce.Do(func() {
		packageRegistryInstance = NewEPackageRegistryImpl()
		packageRegistryInstance.RegisterPackage(GetPackage())
	})
	return packageRegistryInstance
}
