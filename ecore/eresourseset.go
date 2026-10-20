// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

const (
	RESOURCE_SET__RESOURCES = 0
)

// EResourceSet ...
type EResourceSet interface {
	ENotifier

	GetResources() EList
	GetResource(uri *URI, loadOnDemand bool) EResource
	CreateResource(uri *URI) EResource

	GetEObject(uri *URI, loadOnDemand bool) EObject

	GetURIConverter() EURIConverter
	SetURIConverter(uriConverter EURIConverter)

	GetPackageRegistry() EPackageRegistry
	SetPackageRegistry(packageregistry EPackageRegistry)

	GetCodecRegistry() ECodecRegistry
	SetCodecRegistry(resourceCodecRegistry ECodecRegistry)

	SetURIResourceMap(uriMap map[*URI]EResource)
	GetURIResourceMap() map[*URI]EResource
}

// CreateEResourceSet creates a resource set whose package registry holds the
// ecore package and the given packages (one package per nsURI).
//
// This is how a resource set chooses an implementation when several generated
// implementations of the same metamodel (same nsURI, same IDs) are linked in a
// program: pass the package of the implementation this resource set uses.
// Resource sets created with different implementations of the same metamodel
// coexist; each resolves the nsURI to its own package and factory, and the
// global registry is left untouched.
func CreateEResourceSet(packages []EPackage) EResourceSet {
	rs := NewEResourceSetImpl()
	// packages
	packageRegistry := rs.GetPackageRegistry()
	packageRegistry.RegisterPackage(GetPackage())
	for _, pack := range packages {
		packageRegistry.RegisterPackage(pack)
	}
	return rs
}
