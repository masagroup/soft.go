package ecore

import (
	"errors"
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

// fileName returns the file of uri: on windows, "/C:/dir/file" is "C:/dir/file", and a path
// without drive keeps its leading '/'
func fileName(uri *URI) (string, error) {
	name := uri.Path()
	if len(name) == 0 {
		return "", errors.New("the uri '" + uri.String() + "' has no file path")
	}
	if runtime.GOOS == "windows" && len(name) >= 3 && name[0] == '/' && name[2] == ':' {
		name = name[1:]
	}
	return name, nil
}

func (fuh *FileURIHandler) CreateReader(uri *URI) (io.ReadCloser, error) {
	name, err := fileName(uri)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fuh *FileURIHandler) CreateWriter(uri *URI) (io.WriteCloser, error) {
	name, err := fileName(uri)
	if err != nil {
		return nil, err
	}
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	return f, nil
}
