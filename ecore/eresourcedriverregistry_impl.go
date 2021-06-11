package ecore

import "strings"

type EResourceDriverRegistryImpl struct {
	protocolToDriver  map[string]EResourceDriver
	extensionToDriver map[string]EResourceDriver
	delegate          EResourceDriverRegistry
}

func NewEResourceDriverRegistryImpl() *EResourceDriverRegistryImpl {
	return &EResourceDriverRegistryImpl{
		protocolToDriver:  make(map[string]EResourceDriver),
		extensionToDriver: make(map[string]EResourceDriver),
	}
}

func NewEResourceDriverRegistryImplWithDelegate(delegate EResourceDriverRegistry) *EResourceDriverRegistryImpl {
	return &EResourceDriverRegistryImpl{
		protocolToDriver:  make(map[string]EResourceDriver),
		extensionToDriver: make(map[string]EResourceDriver),
		delegate:          delegate,
	}
}

func (r *EResourceDriverRegistryImpl) GetDriver(uri *URI) EResourceDriver {
	if factory, ok := r.protocolToDriver[uri.Scheme]; ok {
		return factory
	}

	ndx := strings.LastIndex(uri.Path, ".")
	if ndx != -1 {
		extension := uri.Path[ndx+1:]
		if factory, ok := r.extensionToDriver[extension]; ok {
			return factory
		}
	}
	if factory, ok := r.extensionToDriver[DEFAULT_EXTENSION]; ok {
		return factory
	}
	if r.delegate != nil {
		return r.delegate.GetDriver(uri)
	}
	return nil
}

func (r *EResourceDriverRegistryImpl) GetProtocolToDriverMap() map[string]EResourceDriver {
	return r.protocolToDriver
}

func (r *EResourceDriverRegistryImpl) GetExtensionToDriverMap() map[string]EResourceDriver {
	return r.extensionToDriver
}
