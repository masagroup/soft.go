// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package xml

import (
	"github.com/masagroup/soft.go/ecore"
	"io"
	"maps"
	"strings"
)

type Processor struct {
	resourceSet    ecore.EResourceSet
	extendMetaData *ExtendedMetaData
	packages       []ecore.EPackage
}

type ProcessorOption interface {
	apply(*Processor)
}

type funcXmlProcessorOption struct {
	f func(*Processor)
}

func (fdo *funcXmlProcessorOption) apply(do *Processor) {
	fdo.f(do)
}

func newFuncXmlProcessorOption(f func(*Processor)) *funcXmlProcessorOption {
	return &funcXmlProcessorOption{
		f: f,
	}
}

func ProcessorPackages(packages []ecore.EPackage) ProcessorOption {
	return newFuncXmlProcessorOption(func(x *Processor) {
		x.packages = packages
	})
}

func ProcessorResourceSet(resourceSet ecore.EResourceSet) ProcessorOption {
	return newFuncXmlProcessorOption(func(x *Processor) {
		x.resourceSet = resourceSet
	})
}

func NewProcessor(opts ...ProcessorOption) *Processor {
	p := &Processor{extendMetaData: NewExtendedMetaData()}
	for _, opt := range opts {
		opt.apply(p)
	}
	return p
}

func (p *Processor) GetResourceSet() ecore.EResourceSet {
	if p.resourceSet == nil {
		return ecore.CreateEResourceSet(p.packages)
	}
	return p.resourceSet
}

func (p *Processor) Load(uri *ecore.URI) ecore.EResource {
	return p.LoadWithOptions(uri, nil)
}

func (p *Processor) LoadWithOptions(uri *ecore.URI, options map[string]any) ecore.EResource {
	rs := p.GetResourceSet()
	r := rs.CreateResource(uri)
	o := map[string]any{OPTION_EXTENDED_META_DATA: p.extendMetaData}
	maps.Copy(o, options)
	r.LoadWithOptions(o)
	return r
}

func (p *Processor) LoadWithReader(r io.Reader, options map[string]any) ecore.EResource {
	rs := p.GetResourceSet()
	rc := rs.CreateResource(ecore.NewURI("*.xml"))
	o := map[string]any{OPTION_EXTENDED_META_DATA: p.extendMetaData}
	maps.Copy(o, options)
	rc.LoadWithReader(r, o)
	return rc
}

func (p *Processor) Save(resource ecore.EResource) {
	p.SaveWithOptions(resource, nil)
}

func (p *Processor) SaveWithOptions(resource ecore.EResource, options map[string]any) {
	o := map[string]any{OPTION_EXTENDED_META_DATA: p.extendMetaData}
	maps.Copy(o, options)
	resource.SaveWithOptions(o)
}

func (p *Processor) SaveWithWriter(w io.Writer, resource ecore.EResource, options map[string]any) {
	o := map[string]any{OPTION_EXTENDED_META_DATA: p.extendMetaData}
	maps.Copy(o, options)
	resource.SaveWithWriter(w, o)
}

func (p *Processor) SaveToString(resource ecore.EResource, options map[string]any) string {
	var strbuff strings.Builder
	p.SaveWithWriter(&strbuff, resource, options)
	return strbuff.String()
}

func (p *Processor) SaveObject(uri *ecore.URI, eObject ecore.EObject) ecore.EResource {
	rs := p.GetResourceSet()
	rc := rs.CreateResource(uri)
	if rc != nil {
		eCopy := ecore.Copy(eObject)
		rc.GetContents().Add(eCopy)
		rc.Save()
	}
	return rc
}

// GetExtendedMetaData returns the extended metadata the processor loads and saves with.
func (p *Processor) GetExtendedMetaData() *ExtendedMetaData {
	return p.extendMetaData
}
