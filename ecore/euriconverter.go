package ecore

import (
	"io"
	"net/url"
)

//URIConverter ...
type EURIConverter interface {
	CreateReader(uri *url.URL) (io.ReadCloser, error)

	CreateWriter(uri *url.URL) (io.WriteCloser, error)

	GetURIMap() map[url.URL]url.URL

	Normalize(uri *url.URL) *url.URL

	GetURIHandler(uri *url.URL) EURIHandler

	GetURIHandlers() EList
}
