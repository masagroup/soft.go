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
	"io"
	"strings"
)

type XMLProcessor struct {
	extendMetaData *ExtendedMetaData
	packages       []EPackage
	resourceSet    EResourceSet
}

func NewXMLProcessor(packages []EPackage) *XMLProcessor {
	return &XMLProcessor{
		extendMetaData: NewExtendedMetaData(),
		packages:       packages,
	}
}

func NewSharedXMLProcessor(resourceSet EResourceSet) *XMLProcessor {
	return &XMLProcessor{
		extendMetaData: NewExtendedMetaData(),
		resourceSet:    resourceSet,
	}
}

func (p *XMLProcessor) GetResourceSet() EResourceSet {
	if p.resourceSet == nil {
		return CreateEResourceSet(p.packages)
	}
	return p.resourceSet
}

func (p *XMLProcessor) Load(uri *URI) EResource {
	return p.LoadWithOptions(uri, nil)
}

func (p *XMLProcessor) LoadWithOptions(uri *URI, options map[string]interface{}) EResource {
	rs := p.GetResourceSet()
	r := rs.CreateResource(uri)
	o := map[string]interface{}{XML_OPTION_EXTENDED_META_DATA: p.extendMetaData}
	if options != nil {
		for k, v := range options {
			o[k] = v
		}
	}
	r.LoadWithOptions(o)
	return r
}

func (p *XMLProcessor) LoadWithReader(r io.Reader, options map[string]interface{}) EResource {
	rs := p.GetResourceSet()
	rc := rs.CreateResource(&URI{Path: "*.xml"})
	o := map[string]interface{}{XML_OPTION_EXTENDED_META_DATA: p.extendMetaData}
	if options != nil {
		for k, v := range options {
			o[k] = v
		}
	}
	rc.LoadWithReader(r, o)
	return rc
}

func (p *XMLProcessor) Save(resource EResource) {
	p.SaveWithOptions(resource, nil)
}

func (p *XMLProcessor) SaveWithOptions(resource EResource, options map[string]interface{}) {
	o := map[string]interface{}{XML_OPTION_EXTENDED_META_DATA: p.extendMetaData}
	if options != nil {
		for k, v := range options {
			o[k] = v
		}
	}
	resource.SaveWithOptions(o)
}

func (p *XMLProcessor) SaveWithWriter(w io.Writer, resource EResource, options map[string]interface{}) {
	o := map[string]interface{}{XML_OPTION_EXTENDED_META_DATA: p.extendMetaData}
	if options != nil {
		for k, v := range options {
			o[k] = v
		}
	}
	resource.SaveWithWriter(w, o)
}

func (p *XMLProcessor) SaveToString(resource EResource, options map[string]interface{}) string {
	var strbuff strings.Builder
	p.SaveWithWriter(&strbuff, resource, options)
	return strbuff.String()
}

func (p *XMLProcessor) SaveObject(uri *URI, eObject EObject) EResource {
	rs := p.GetResourceSet()
	rc := rs.CreateResource(uri)
	if rc != nil {
		eCopy := Copy(eObject)
		rc.GetContents().Add(eCopy)
		rc.Save()
	}
	return rc
}
