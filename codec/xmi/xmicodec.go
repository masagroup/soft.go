// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package xmi

import (
	"io"
	"maps"

	"github.com/masagroup/soft.go/codec/xml"
	"github.com/masagroup/soft.go/ecore"
)

// Codec is the XMI codec: it encodes and decodes resources with Encoder and Decoder.
type Codec struct {
}

func (d *Codec) NewEncoder(resource ecore.EResource, w io.Writer, options map[string]any) ecore.EEncoder {
	return NewEncoder(resource, w, options)
}

func (d *Codec) NewDecoder(resource ecore.EResource, r io.Reader, options map[string]any) ecore.EDecoder {
	return NewDecoder(resource, r, options)
}

// EcoreCodec is the XMI codec of the .ecore resources: as EMF's EcoreResourceFactoryImpl, it
// saves the cross resource references in the encoded attribute style, unless the options of the
// save set xml.OPTION_USE_ENCODED_ATTRIBUTE_STYLE.
type EcoreCodec struct {
	Codec
}

func (d *EcoreCodec) NewEncoder(resource ecore.EResource, w io.Writer, options map[string]any) ecore.EEncoder {
	o := map[string]any{xml.OPTION_USE_ENCODED_ATTRIBUTE_STYLE: true}
	maps.Copy(o, options)
	return NewEncoder(resource, w, o)
}
