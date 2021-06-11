package ecore

import "strings"

type EResourceCodecRegistryImpl struct {
	protocolToCodec  map[string]EResourceCodec
	extensionToCodec map[string]EResourceCodec
	delegate         EResourceCodecRegistry
}

func NewEResourceCodecRegistryImpl() *EResourceCodecRegistryImpl {
	return &EResourceCodecRegistryImpl{
		protocolToCodec:  make(map[string]EResourceCodec),
		extensionToCodec: make(map[string]EResourceCodec),
	}
}

func NewEResourceCodecRegistryImplWithDelegate(delegate EResourceCodecRegistry) *EResourceCodecRegistryImpl {
	return &EResourceCodecRegistryImpl{
		protocolToCodec:  make(map[string]EResourceCodec),
		extensionToCodec: make(map[string]EResourceCodec),
		delegate:         delegate,
	}
}

func (r *EResourceCodecRegistryImpl) GetCodec(uri *URI) EResourceCodec {
	if factory, ok := r.protocolToCodec[uri.Scheme]; ok {
		return factory
	}

	ndx := strings.LastIndex(uri.Path, ".")
	if ndx != -1 {
		extension := uri.Path[ndx+1:]
		if factory, ok := r.extensionToCodec[extension]; ok {
			return factory
		}
	}
	if factory, ok := r.extensionToCodec[DEFAULT_EXTENSION]; ok {
		return factory
	}
	if r.delegate != nil {
		return r.delegate.GetCodec(uri)
	}
	return nil
}

func (r *EResourceCodecRegistryImpl) GetProtocolToCodecMap() map[string]EResourceCodec {
	return r.protocolToCodec
}

func (r *EResourceCodecRegistryImpl) GetExtensionToCodecMap() map[string]EResourceCodec {
	return r.extensionToCodec
}
