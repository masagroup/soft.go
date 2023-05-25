// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import "strings"

type ECodecRegistryImpl struct {
	protocolToCodec  map[string]ECodec
	extensionToCodec map[string]ECodec
	delegate         ECodecRegistry
}

func NewECodecRegistryImpl() *ECodecRegistryImpl {
	return &ECodecRegistryImpl{
		protocolToCodec:  make(map[string]ECodec),
		extensionToCodec: make(map[string]ECodec),
	}
}

func NewECodecRegistryImplWithDelegate(delegate ECodecRegistry) *ECodecRegistryImpl {
	return &ECodecRegistryImpl{
		protocolToCodec:  make(map[string]ECodec),
		extensionToCodec: make(map[string]ECodec),
		delegate:         delegate,
	}
}

func (r *ECodecRegistryImpl) GetCodec(uri *URI) ECodec {
	if factory, ok := r.protocolToCodec[uri.scheme]; ok {
		return factory
	}
	p := uri.Path()
	ndx := strings.LastIndex(p, ".")
	if ndx != -1 {
		extension := p[ndx+1:]
		if factory, ok := r.extensionToCodec[extension]; ok {
			return factory
		}
	}
	if factory, ok := r.extensionToCodec[DEFAULT_EXTENSION]; ok {
		return factory
	}
	if r.delegate != nil {
		return r.delegate.GetCodec(uri)
	}
	return nil
}

func (r *ECodecRegistryImpl) GetProtocolToCodecMap() map[string]ECodec {
	return r.protocolToCodec
}

func (r *ECodecRegistryImpl) GetExtensionToCodecMap() map[string]ECodec {
	return r.extensionToCodec
}
