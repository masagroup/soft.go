package ecore

import (
	"io"
	"net/url"
	"strings"
)

type XMLProcessor struct {
	extendMetaData *ExtendedMetaData
	packages       []EPackage
	factories      map[string]EResourceFactory
}

func NewXMLProcessor(packages []EPackage) *XMLProcessor {
	p := new(XMLProcessor)
	p.Initialize(packages)
	return p
}

func (p *XMLProcessor) Initialize(packages []EPackage) {
	p.extendMetaData = NewExtendedMetaData()
	p.packages = packages
}

func (p *XMLProcessor) Load(uri *url.URL) EResource {
	return p.LoadWithOptions(uri, nil)
}

func (p *XMLProcessor) LoadWithOptions(uri *url.URL, options map[string]interface{}) EResource {
	rs := p.CreateEResourceSet()
	r := rs.CreateResource(uri)
	o := map[string]interface{}{OPTION_EXTENDED_META_DATA: p.extendMetaData}
	if options != nil {
		for k, v := range options {
			o[k] = v
		}
	}
	r.LoadWithOptions(o)
	return r
}

func (p *XMLProcessor) LoadWithReader(r io.Reader, options map[string]interface{}) EResource {
	rs := p.CreateEResourceSet()
	rc := rs.CreateResource(&url.URL{Path: "*.xml"})
	o := map[string]interface{}{OPTION_EXTENDED_META_DATA: p.extendMetaData}
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
	o := map[string]interface{}{OPTION_EXTENDED_META_DATA: p.extendMetaData}
	if options != nil {
		for k, v := range options {
			o[k] = v
		}
	}
	resource.SaveWithOptions(o)
}

func (p *XMLProcessor) SaveWithWriter(w io.Writer, resource EResource, options map[string]interface{}) {
	o := map[string]interface{}{OPTION_EXTENDED_META_DATA: p.extendMetaData}
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

func (p *XMLProcessor) CreateEResourceSet() EResourceSet {
	rs := NewEResourceSetImpl()
	// packages
	packageRegistry := rs.GetPackageRegistry()
	packageRegistry.RegisterPackage(GetPackage())
	if p.packages != nil {
		for _, pack := range p.packages {
			packageRegistry.RegisterPackage(pack)
		}
	}
	// factories
	extensionToFactories := rs.GetResourceFactoryRegistry().GetExtensionToFactoryMap()
	extensionToFactories["ecore"] = &XMIResourceFactory{}
	extensionToFactories["xml"] = &XMLResourceFactory{}
	if p.factories != nil {
		for ext, factory := range p.factories {
			extensionToFactories[ext] = factory
		}
	}
	return rs
}
