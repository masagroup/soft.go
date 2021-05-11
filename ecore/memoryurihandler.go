package ecore

import (
	"io"
)

type MemoryURIHandler struct {
}

func (muh *MemoryURIHandler) CanHandle(uri *URI) bool {
	return uri.Scheme == "memory"
}

func (muh *MemoryURIHandler) CreateReader(uri *URI) (io.ReadCloser, error) {
	return nil, nil
}

func (muh *MemoryURIHandler) CreateWriter(uri *URI) (io.WriteCloser, error) {
	return nil, nil
}
