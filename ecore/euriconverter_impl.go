package ecore

import (
	"fmt"
	"io"
	"net/url"
)

type EURIConverterImpl struct {
	uriHandlers EList
	uriMap      map[url.URL]url.URL
}

func NewEURIConverterImpl() *EURIConverterImpl {
	r := new(EURIConverterImpl)
	r.uriHandlers = NewImmutableEList([]interface{}{new(FileURIHandler)})
	r.uriMap = make(map[url.URL]url.URL)
	return r
}

func (r *EURIConverterImpl) CreateReader(uri *url.URL) (io.ReadCloser, error) {
	normalized := r.Normalize(uri)
	uriHandler := r.GetURIHandler(normalized)
	if uriHandler == nil {
		return nil, fmt.Errorf("URIHandler for URI '%s' not found", normalized.String())
	}
	return uriHandler.CreateReader(normalized)
}

func (r *EURIConverterImpl) CreateWriter(uri *url.URL) (io.WriteCloser, error) {
	normalized := r.Normalize(uri)
	uriHandler := r.GetURIHandler(normalized)
	if uriHandler == nil {
		return nil, fmt.Errorf("URIHandler for URI '%s' not found", normalized.String())
	}
	return uriHandler.CreateWriter(normalized)
}

func (r *EURIConverterImpl) GetURIMap() map[url.URL]url.URL {
	return r.uriMap
}

func (r *EURIConverterImpl) Normalize(uri *url.URL) *url.URL {
	normalized := r.getURIFromMap(uri)
	if normalized == uri || *normalized == *uri {
		return normalized
	}
	return r.Normalize(normalized)
}

func (r *EURIConverterImpl) getURIFromMap(uri *url.URL) *url.URL {
	for oldPrefix, newPrefix := range r.uriMap {
		if r := ReplacePrefixURI(uri, &oldPrefix, &newPrefix); r != nil {
			return r
		}
	}
	return uri
}

func (r *EURIConverterImpl) GetURIHandler(uri *url.URL) EURIHandler {
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
