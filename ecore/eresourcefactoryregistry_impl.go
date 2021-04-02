// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import (
	"strings"
)

//EResourceFactoryRegistryImpl ...
type EResourceFactoryRegistryImpl struct {
	protocolToFactory  map[string]EResourceFactory
	extensionToFactory map[string]EResourceFactory
	delegate           EResourceFactoryRegistry
}

func NewEResourceFactoryRegistryImpl() *EResourceFactoryRegistryImpl {
	return &EResourceFactoryRegistryImpl{
		protocolToFactory:  make(map[string]EResourceFactory),
		extensionToFactory: make(map[string]EResourceFactory),
	}
}

func NewEResourceFactoryRegistryImplWithDelegate(delegate EResourceFactoryRegistry) *EResourceFactoryRegistryImpl {
	return &EResourceFactoryRegistryImpl{
		protocolToFactory:  make(map[string]EResourceFactory),
		extensionToFactory: make(map[string]EResourceFactory),
		delegate:           delegate,
	}
}

func (r *EResourceFactoryRegistryImpl) GetFactory(uri *URI) EResourceFactory {
	if factory, ok := r.protocolToFactory[uri.Scheme]; ok {
		return factory
	}

	ndx := strings.LastIndex(uri.Path, ".")
	if ndx != -1 {
		extension := uri.Path[ndx+1:]
		if factory, ok := r.extensionToFactory[extension]; ok {
			return factory
		}
	}
	if factory, ok := r.extensionToFactory[DEFAULT_EXTENSION]; ok {
		return factory
	}
	if r.delegate != nil {
		return r.delegate.GetFactory(uri)
	}
	return nil
}

func (r *EResourceFactoryRegistryImpl) GetProtocolToFactoryMap() map[string]EResourceFactory {
	return r.protocolToFactory
}

func (r *EResourceFactoryRegistryImpl) GetExtensionToFactoryMap() map[string]EResourceFactory {
	return r.extensionToFactory
}
