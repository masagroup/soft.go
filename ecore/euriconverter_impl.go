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
	return r
}

func (r *EURIConverterImpl) CreateReader(uri *url.URL) (io.ReadCloser, error) {
	uriHandler := r.GetURIHandler(uri)
	if uriHandler == nil {
		return nil, fmt.Errorf("URIHandler for URI '%s' not found", uri.String())
	}
	return uriHandler.CreateReader(uri)
}

func (r *EURIConverterImpl) CreateWriter(uri *url.URL) (io.WriteCloser, error) {
	uriHandler := r.GetURIHandler(uri)
	if uriHandler == nil {
		return nil, fmt.Errorf("URIHandler for URI '%s' not found", uri.String())
	}
	return uriHandler.CreateWriter(uri)
}

func (r *EURIConverterImpl) GetURIMap() map[url.URL]url.URL {
	return r.uriMap
}

func (r *EURIConverterImpl) Normalize(uri *url.URL) *url.URL {
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
