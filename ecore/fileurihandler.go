package ecore

import (
	"io"
	"net/url"
	"os"
)

//URIHandler ...
type FileURIHandler struct {
}

func (fuh *FileURIHandler) CanHandle(uri *url.URL) bool {
	return uri.Scheme == "file" || (len(uri.Scheme) == 0 && len(uri.Host) == 0 && len(uri.RawQuery) == 0)
}

func (fuh *FileURIHandler) CreateReader(uri *url.URL) (io.ReadCloser, error) {
	fileName := uri.Path
	if fileName[0] == '/' {
		fileName = fileName[1:]
	}
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fuh *FileURIHandler) CreateWriter(uri *url.URL) (io.WriteCloser, error) {
	fileName := uri.Path
	if fileName[0] == '/' {
		fileName = fileName[1:]
	}
	f, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	return f, nil
}
