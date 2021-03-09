package ecore

import (
	"net/url"
	"path/filepath"
)

func TrimURIFragment(uri *url.URL) *url.URL {
	return &url.URL{
		Scheme:     uri.Scheme,
		User:       uri.User,
		Host:       uri.Host,
		Path:       uri.Path,
		ForceQuery: uri.ForceQuery,
		RawPath:    uri.RawPath,
		RawQuery:   uri.RawQuery,
	}
}

func CloneURI(uri *url.URL) *url.URL {
	return &url.URL{
		Scheme:     uri.Scheme,
		User:       uri.User,
		Host:       uri.Host,
		Path:       uri.Path,
		Fragment:   uri.Fragment,
		ForceQuery: uri.ForceQuery,
		RawPath:    uri.RawPath,
		RawQuery:   uri.RawQuery,
	}
}

func CreateFileURI(path string) *url.URL {
	p := filepath.ToSlash(path)
	if len(p) == 0 {
		return &url.URL{}
	} else {
		if isAbsolute := filepath.IsAbs(p); isAbsolute {
			return &url.URL{Scheme: "file", Path: p}
		} else {
			return &url.URL{Path: p}
		}
	}
}
