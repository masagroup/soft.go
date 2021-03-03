package ecore

import (
	"io"
	"net/url"
)

//URIHandler ...
type EURIHandler interface {
	CanHandle(uri *url.URL) bool

	CreateReader(uri *url.URL) (io.ReadCloser, error)

	CreateWriter(uri *url.URL) (io.WriteCloser, error)
}
