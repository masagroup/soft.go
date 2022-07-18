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

type EResourceCodecRegistryImpl struct {
	protocolToCodec  map[string]EResourceCodec
	extensionToCodec map[string]EResourceCodec
	delegate         EResourceCodecRegistry
}

func NewEResourceCodecRegistryImpl() *EResourceCodecRegistryImpl {
	return &EResourceCodecRegistryImpl{
		protocolToCodec:  make(map[string]EResourceCodec),
		extensionToCodec: make(map[string]EResourceCodec),
	}
}

func NewEResourceCodecRegistryImplWithDelegate(delegate EResourceCodecRegistry) *EResourceCodecRegistryImpl {
	return &EResourceCodecRegistryImpl{
		protocolToCodec:  make(map[string]EResourceCodec),
		extensionToCodec: make(map[string]EResourceCodec),
		delegate:         delegate,
	}
}

func (r *EResourceCodecRegistryImpl) GetCodec(uri *URI) EResourceCodec {
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

func (r *EResourceCodecRegistryImpl) GetProtocolToCodecMap() map[string]EResourceCodec {
	return r.protocolToCodec
}

func (r *EResourceCodecRegistryImpl) GetExtensionToCodecMap() map[string]EResourceCodec {
	return r.extensionToCodec
}
