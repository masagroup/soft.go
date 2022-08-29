package ecore

import (
	"io"
)

// URIConverter ...
type EURIConverter interface {
	CreateReader(uri *URI) (io.ReadCloser, error)

	CreateWriter(uri *URI) (io.WriteCloser, error)

	GetURIMap() map[URI]URI

	Normalize(uri *URI) *URI

	GetURIHandler(uri *URI) EURIHandler

	GetURIHandlers() EList
}
