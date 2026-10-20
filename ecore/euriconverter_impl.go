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
	// the uri is mapped until it does not change, or until a uri already met: a cycle of mappings
	var visited map[string]struct{}
	current := uri
	for {
		normalized := r.getURIFromMap(current)
		if current.Equals(normalized) {
			return normalized
		}
		if visited == nil {
			visited = map[string]struct{}{}
		}
		visited[current.String()] = struct{}{}
		if _, isVisited := visited[normalized.String()]; isVisited {
			return current
		}
		current = normalized
	}
}

// getURIFromMap maps uri with the longest prefix of the uri map matching it, as EMF does: the map
// iteration order does not choose between overlapping prefixes. A prefix matches whole path
// segments: it ends with '/', or the path of uri is the prefix or continues with '/'.
func (r *EURIConverterImpl) getURIFromMap(uri *URI) *URI {
	var result *URI
	bestLength := -1
	bestKey := ""
	for oldPrefix, newPrefix := range r.uriMap {
		oldPath := oldPrefix.Path()
		if len(oldPath) > 0 && oldPath[len(oldPath)-1] != '/' && len(uri.Path()) > len(oldPath) && uri.Path()[len(oldPath)] != '/' {
			continue
		}
		if mapped := uri.ReplacePrefix(&oldPrefix, &newPrefix); mapped != nil {
			if key := oldPrefix.String(); len(oldPath) > bestLength || (len(oldPath) == bestLength && key < bestKey) {
				result, bestLength, bestKey = mapped, len(oldPath), key
			}
		}
	}
	if result == nil {
		return uri
	}
	return result
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
