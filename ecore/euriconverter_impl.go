package ecore

import (
	"fmt"
	"io"
)

type EURIConverterImpl struct {
	uriHandlers EList
	uriMap      map[URI]URI
}

func NewEURIConverterImpl() *EURIConverterImpl {
	r := new(EURIConverterImpl)
	r.uriHandlers = NewImmutableEList([]any{new(FileURIHandler), new(MemoryURIHandler)})
	r.uriMap = make(map[URI]URI)
	return r
}

func (r *EURIConverterImpl) CreateReader(uri *URI) (io.ReadCloser, error) {
	normalized := r.Normalize(uri)
	uriHandler := r.GetURIHandler(normalized)
	if uriHandler == nil {
		return nil, fmt.Errorf("URIHandler for URI '%s' not found", normalized.String())
	}
	return uriHandler.CreateReader(normalized)
}

func (r *EURIConverterImpl) CreateWriter(uri *URI) (io.WriteCloser, error) {
	normalized := r.Normalize(uri)
	uriHandler := r.GetURIHandler(normalized)
	if uriHandler == nil {
		return nil, fmt.Errorf("URIHandler for URI '%s' not found", normalized.String())
	}
	return uriHandler.CreateWriter(normalized)
}

func (r *EURIConverterImpl) GetURIMap() map[URI]URI {
	return r.uriMap
}

func (r *EURIConverterImpl) Normalize(uri *URI) *URI {
	normalized := r.getURIFromMap(uri)
	if uri.Equals(normalized) {
		return normalized
	}
	return r.Normalize(normalized)
}

func (r *EURIConverterImpl) getURIFromMap(uri *URI) *URI {
	for oldPrefix, newPrefix := range r.uriMap {
		if r := uri.ReplacePrefix(&oldPrefix, &newPrefix); r != nil {
			return r
		}
	}
	return uri
}

func (r *EURIConverterImpl) GetURIHandler(uri *URI) EURIHandler {
	if uri != nil {
		for it := r.uriHandlers.Iterator(); it.HasNext(); {
			uriHandler := it.Next().(EURIHandler)
			if uriHandler.CanHandle(uri) {
				return uriHandler
			}
		}
	}
	return nil
}

func (r *EURIConverterImpl) GetURIHandlers() EList {
	return r.uriHandlers
}
