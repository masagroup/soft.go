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

//EResourceSet ...
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

	GetResourceFactoryRegistry() EResourceFactoryRegistry
	SetResourceFactoryRegistry(resourceFactoryRegistry EResourceFactoryRegistry)

	SetURIResourceMap(uriMap map[*URI]EResource)
	GetURIResourceMap() map[*URI]EResource
}
