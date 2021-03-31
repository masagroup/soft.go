package ecore

import (
	"errors"
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type URI struct {
	Scheme   string
	Username string
	Password string
	Host     string
	Port     int
	Path     string
	Query    string
	Fragment string
}

var uriRegExp = regexp.MustCompile(
	"(([a-zA-Z][a-zA-Z0-9+.-]*):)?" + // scheme:
		"([^?#]*)" + // authority and path
		"(?:\\?([^#]*))?" + //?query
		"(?:#(.*))?", // #fragment
)

var authorityAndPathRegExp = regexp.MustCompile("//([^/]*)(/.*)?")

var authorityRegExp = regexp.MustCompile(
	"(?:([^@:]*)(?::([^@]*))?@)?" + // username, password
		"(\\[[^\\]]*\\]|[^\\[:]*)" + // host (IP-literal (e.g. '['+IPv6+']',dotted-IPv4, or named host)
		"(?::(\\d*))?", // port
)

func NewURI(rawURI string) *URI {
	uri, _ := ParseURI(rawURI)
	return uri
}

func ParseURI(rawURI string) (*URI, error) {
	uriMatches := uriRegExp.FindStringSubmatch(rawURI)
	if uriMatches == nil {
		return nil, errors.New("invalid URI :'" + rawURI + "'")
	}
	uri := &URI{}
	uri.Scheme = strings.ToLower(uriMatches[2])
	authorityAndPath := uriMatches[3]
	autorityAndPathMatches := authorityAndPathRegExp.FindStringSubmatch(authorityAndPath)
	if autorityAndPathMatches == nil {
		uri.Path = authorityAndPath
	} else {
		authority := autorityAndPathMatches[1]
		autorityMatches := authorityRegExp.FindStringSubmatch(authority)
		if autorityMatches == nil {
			return nil, errors.New("invalid URI authority " + authority)
		}
		uri.Username = autorityMatches[1]
		uri.Password = autorityMatches[2]
		uri.Host = autorityMatches[3]
		if portStr := autorityMatches[4]; len(portStr) > 0 {
			uri.Port, _ = strconv.Atoi(portStr)
		}
		uri.Path = autorityAndPathMatches[2]
	}
	uri.Query = uriMatches[4]
	uri.Fragment = uriMatches[5]
	return uri, nil
}

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
