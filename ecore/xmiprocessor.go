// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import "io"

type XMIProcessor struct {
	XMLProcessor
}

func NewXMIProcessor() *XMIProcessor {
	return &XMIProcessor{XMLProcessor{
		extendMetaData: NewExtendedMetaData(),
	}}
}

func NewSharedXMIProcessor(resourceSet EResourceSet) *XMIProcessor {
	return &XMIProcessor{XMLProcessor{
		extendMetaData: NewExtendedMetaData(),
		resourceSet:    resourceSet,
	}}
}

func (p *XMIProcessor) LoadWithReader(r io.Reader, options map[string]any) EResource {
	rs := p.GetResourceSet()
	rc := rs.CreateResource(NewURI("*.ecore"))
	o := map[string]any{XML_OPTION_EXTENDED_META_DATA: p.extendMetaData}
	for k, v := range options {
		o[k] = v
	}
	rc.LoadWithReader(r, o)
	return rc
}
