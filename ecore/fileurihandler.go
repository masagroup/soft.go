package ecore

import (
	"io"
	"os"
	"runtime"
)

// URIHandler ...
type FileURIHandler struct {
}

func (fuh *FileURIHandler) CanHandle(uri *URI) bool {
	return uri.scheme == "file" || (len(uri.scheme) == 0 && len(uri.host) == 0 && len(uri.query) == 0)
}

func (fuh *FileURIHandler) CreateReader(uri *URI) (io.ReadCloser, error) {
	fileName := uri.Path()
	if runtime.GOOS == "windows" && fileName[0] == '/' {
		fileName = fileName[1:]
	}
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fuh *FileURIHandler) CreateWriter(uri *URI) (io.WriteCloser, error) {
	fileName := uri.Path()
	if runtime.GOOS == "windows" && fileName[0] == '/' {
		fileName = fileName[1:]
	}
	f, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	return f, nil
}
