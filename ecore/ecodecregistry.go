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

const (
	DEFAULT_EXTENSION = "*"
)

type ECodecRegistry interface {
	GetCodec(uri *URI) ECodec
	GetProtocolToCodecMap() map[string]ECodec
	GetExtensionToCodecMap() map[string]ECodec
}

var resourceCodecRegistryOnce sync.Once
var resourceCodecRegistryInstance ECodecRegistry

// GetCodecRegistry returns the process-wide default codec registry, built once on first use.
func GetCodecRegistry() ECodecRegistry {
	resourceCodecRegistryOnce.Do(func() {
		resourceCodecRegistryInstance = NewECodecRegistryImpl()
		// initialize with default codecs
		extensionToCodecs := resourceCodecRegistryInstance.GetExtensionToCodecMap()
		extensionToCodecs["ecore"] = &XMICodec{}
		extensionToCodecs["xmi"] = &XMICodec{}
		extensionToCodecs["xml"] = &XMLCodec{}
		extensionToCodecs["bin"] = &BinaryCodec{}
		extensionToCodecs["sqlite"] = &SQLCodec{}
		protocolToCodecs := resourceCodecRegistryInstance.GetProtocolToCodecMap()
		protocolToCodecs["memory"] = &NoCodec{}
	})
	return resourceCodecRegistryInstance
}
