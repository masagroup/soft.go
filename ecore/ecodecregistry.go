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
	"sync"
)

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
// It only maps the memory protocol: the codec packages (codec/xml, codec/xmi, codec/binary,
// codec/sql, codec/json) register their extensions when they are imported, and codec/defaults
// imports the default ones.
func GetCodecRegistry() ECodecRegistry {
	resourceCodecRegistryOnce.Do(func() {
		resourceCodecRegistryInstance = NewECodecRegistryImpl()
		protocolToCodecs := resourceCodecRegistryInstance.GetProtocolToCodecMap()
		protocolToCodecs["memory"] = &NoCodec{}
	})
	return resourceCodecRegistryInstance
}

// RegisterCodec maps extension to codec in the default codec registry. It is meant to be called
// from the init function of a codec package: the registry is not safe for concurrent writes.
func RegisterCodec(extension string, codec ECodec) {
	GetCodecRegistry().GetExtensionToCodecMap()[extension] = codec
}

// RegisterProtocolCodec maps protocol to codec in the default codec registry. Like
// RegisterCodec, it is meant to be called from an init function.
func RegisterProtocolCodec(protocol string, codec ECodec) {
	GetCodecRegistry().GetProtocolToCodecMap()[protocol] = codec
}

// codecPackages are the packages registering the codecs of the well-known extensions.
var codecPackages = map[string]string{
	"ecore":  "github.com/masagroup/soft.go/codec/xmi",
	"xmi":    "github.com/masagroup/soft.go/codec/xmi",
	"xml":    "github.com/masagroup/soft.go/codec/xml",
	"bin":    "github.com/masagroup/soft.go/codec/binary",
	"sqlite": "github.com/masagroup/soft.go/codec/sql",
	"json":   "github.com/masagroup/soft.go/codec/json",
}

// codecNotFoundMessage is the error of a resource without codec: for a well-known extension, it
// names the package to import.
func codecNotFoundMessage(uri *URI) string {
	message := "Unable to find codec for '" + uri.String() + "'"
	p := uri.Path()
	if ndx := strings.LastIndex(p, "."); ndx != -1 {
		if pkg, ok := codecPackages[p[ndx+1:]]; ok {
			message += ": import \"" + pkg + "\""
		}
	}
	return message
}
