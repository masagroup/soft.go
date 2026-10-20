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

// Processor loads and saves XMI resources: an XML processor whose resources are .ecore ones.
type Processor struct {
	*xml.Processor
}

// NewProcessor returns an XMI processor configured with the XML processor options, as
// xml.ProcessorPackages or xml.ProcessorResourceSet.
func NewProcessor(opts ...xml.ProcessorOption) *Processor {
	return &Processor{xml.NewProcessor(opts...)}
}

func (p *Processor) LoadWithReader(r io.Reader, options map[string]any) ecore.EResource {
	rs := p.GetResourceSet()
	rc := rs.CreateResource(ecore.NewURI("*.ecore"))
	o := map[string]any{xml.OPTION_EXTENDED_META_DATA: p.GetExtendedMetaData()}
	maps.Copy(o, options)
	rc.LoadWithReader(r, o)
	return rc
}
