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
	DEFAULT_EXTENSION = "*"
)

//EResourceFactoryRegistry ...
type EResourceFactoryRegistry interface {
	GetFactory(uri *URI) EResourceFactory
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
