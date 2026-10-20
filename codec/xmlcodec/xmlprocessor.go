// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package xmlcodec

import (
	"github.com/masagroup/soft.go/ecore"
	"io"
	"maps"
	"strings"
)

type XMLProcessor struct {
	resourceSet    ecore.EResourceSet
	extendMetaData *ExtendedMetaData
	packages       []ecore.EPackage
}

type xmlProcessorOption interface {
	apply(*XMLProcessor)
}

type funcXmlProcessorOption struct {
	f func(*XMLProcessor)
}

func (fdo *funcXmlProcessorOption) apply(do *XMLProcessor) {
	fdo.f(do)
}

func newFuncXmlProcessorOption(f func(*XMLProcessor)) *funcXmlProcessorOption {
	return &funcXmlProcessorOption{
		f: f,
	}
}

func XMLProcessorPackages(packages []ecore.EPackage) xmlProcessorOption {
	return newFuncXmlProcessorOption(func(x *XMLProcessor) {
		x.packages = packages
	})
}

func XMLProcessorResourceSet(resourceSet ecore.EResourceSet) xmlProcessorOption {
	return newFuncXmlProcessorOption(func(x *XMLProcessor) {
		x.resourceSet = resourceSet
	})
}

func NewXMLProcessor(opts ...xmlProcessorOption) *XMLProcessor {
	p := &XMLProcessor{extendMetaData: NewExtendedMetaData()}
	for _, opt := range opts {
		opt.apply(p)
	}
	return p
}

func (p *XMLProcessor) GetResourceSet() ecore.EResourceSet {
	if p.resourceSet == nil {
		return ecore.CreateEResourceSet(p.packages)
	}
	return p.resourceSet
}

func (p *XMLProcessor) Load(uri *ecore.URI) ecore.EResource {
	return p.LoadWithOptions(uri, nil)
}

func (p *XMLProcessor) LoadWithOptions(uri *ecore.URI, options map[string]any) ecore.EResource {
	rs := p.GetResourceSet()
	r := rs.CreateResource(uri)
	o := map[string]any{XML_OPTION_EXTENDED_META_DATA: p.extendMetaData}
	maps.Copy(o, options)
	r.LoadWithOptions(o)
	return r
}

func (p *XMLProcessor) LoadWithReader(r io.Reader, options map[string]any) ecore.EResource {
	rs := p.GetResourceSet()
	rc := rs.CreateResource(ecore.NewURI("*.xml"))
	o := map[string]any{XML_OPTION_EXTENDED_META_DATA: p.extendMetaData}
	maps.Copy(o, options)
	rc.LoadWithReader(r, o)
	return rc
}

func (p *XMLProcessor) Save(resource ecore.EResource) {
	p.SaveWithOptions(resource, nil)
}

func (p *XMLProcessor) SaveWithOptions(resource ecore.EResource, options map[string]any) {
	o := map[string]any{XML_OPTION_EXTENDED_META_DATA: p.extendMetaData}
	maps.Copy(o, options)
	resource.SaveWithOptions(o)
}

func (p *XMLProcessor) SaveWithWriter(w io.Writer, resource ecore.EResource, options map[string]any) {
	o := map[string]any{XML_OPTION_EXTENDED_META_DATA: p.extendMetaData}
	maps.Copy(o, options)
	resource.SaveWithWriter(w, o)
}

func (p *XMLProcessor) SaveToString(resource ecore.EResource, options map[string]any) string {
	var strbuff strings.Builder
	p.SaveWithWriter(&strbuff, resource, options)
	return strbuff.String()
}

func (p *XMLProcessor) SaveObject(uri *ecore.URI, eObject ecore.EObject) ecore.EResource {
	rs := p.GetResourceSet()
	rc := rs.CreateResource(uri)
	if rc != nil {
		eCopy := ecore.Copy(eObject)
		rc.GetContents().Add(eCopy)
		rc.Save()
	}
	return rc
}
