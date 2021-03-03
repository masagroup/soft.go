package ecore

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateFileURI(t *testing.T) {
	assert.Equal(t, &url.URL{Scheme: "file", Path: "C:/path/file.xml"}, CreateFileURI("C:/path/file.xml"))
	assert.Equal(t, &url.URL{Scheme: "file", Path: "C:/path/file.xml"}, CreateFileURI("C:\\path\\file.xml"))
	assert.Equal(t, &url.URL{Path: "path/file.xml"}, CreateFileURI("path/file.xml"))
	assert.Equal(t, &url.URL{Path: "path/file.xml"}, CreateFileURI("path\\file.xml"))
}
