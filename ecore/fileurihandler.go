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

func (fuh *FileURIHandler) CreateReader(uri *url.URL) io.ReadCloser {
	fileName := uri.Path
	if fileName[0] == '/' {
		fileName = fileName[1:]
	}
	f, error := os.Open(fileName)
	if error != nil {
		return nil
	}
	return f
}

func (fuh *FileURIHandler) CreateWriter(uri *url.URL) io.WriteCloser {
	fileName := uri.Path
	if fileName[0] == '/' {
		fileName = fileName[1:]
	}
	f, error := os.Create(fileName)
	if error != nil {
		return nil
	}
	return f
}
