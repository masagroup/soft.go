package ecore

import (
	"io"
)

//URIHandler ...
type EURIHandler interface {
	CanHandle(uri *URI) bool

	CreateReader(uri *URI) (io.ReadCloser, error)

	CreateWriter(uri *URI) (io.WriteCloser, error)
}
