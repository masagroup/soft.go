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

type ECodecRegistry interface {
	GetCodec(uri *URI) ECodec
	GetProtocolToCodecMap() map[string]ECodec
	GetExtensionToCodecMap() map[string]ECodec
}

var resourceCodecRegistryInstance ECodecRegistry

func GetCodecRegistry() ECodecRegistry {
	if resourceCodecRegistryInstance == nil {
		resourceCodecRegistryInstance = NewECodecRegistryImpl()
		// initialize with default codecs
		extensionToCodecs := resourceCodecRegistryInstance.GetExtensionToCodecMap()
		extensionToCodecs["ecore"] = &XMICodec{}
		extensionToCodecs["xml"] = &XMLCodec{}
		extensionToCodecs["bin"] = &BinaryCodec{}
		extensionToCodecs["sqlite"] = &SQLCodec{}
		protocolToCodecs := resourceCodecRegistryInstance.GetProtocolToCodecMap()
		protocolToCodecs["memory"] = &NoCodec{}
	}
	return resourceCodecRegistryInstance
}
