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

func ReplacePrefixURI(uri *url.URL, oldPrefix *url.URL, newPrefix *url.URL) *url.URL {
	if uri.Scheme != oldPrefix.Scheme ||
		uri.Host != oldPrefix.Host ||
		(uri.User != oldPrefix.User || uri.User != nil && oldPrefix.User != nil && *uri.User == *oldPrefix.User) {
		return nil
	}
	if oldLen := len(oldPrefix.Path); len(uri.Path) >= oldLen && uri.Path[0:oldLen] == oldPrefix.Path {
		return &url.URL{Scheme: newPrefix.Scheme,
			User:     newPrefix.User,
			Host:     newPrefix.Host,
			Path:     newPrefix.Path + uri.Path[oldLen:],
			RawQuery: uri.RawQuery,
			Fragment: uri.Fragment}
	}
	return nil

}
