package ecore

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateFileURI(t *testing.T) {
	assert.Equal(t, &url.URL{Scheme: "file", Path: "C:/path/file.xml"}, CreateFileURI("C:/path/file.xml"))
	assert.Equal(t, &url.URL{Scheme: "file", Path: "C:/path/file.xml"}, CreateFileURI("C:\\path\\file.xml"))
	assert.Equal(t, &url.URL{Path: "path/file.xml"}, CreateFileURI("path/file.xml"))
	assert.Equal(t, &url.URL{Path: "path/file.xml"}, CreateFileURI("path\\file.xml"))
	assert.Equal(t, &url.URL{}, CreateFileURI(""))
}

func TestReplacePrefix(t *testing.T) {
	assert.Nil(t, ReplacePrefixURI(&url.URL{Scheme: "file"}, &url.URL{Path: "http"}, nil))
	assert.Nil(t, ReplacePrefixURI(&url.URL{Host: "host"}, &url.URL{Path: "http"}, nil))
	assert.Nil(t, ReplacePrefixURI(&url.URL{User: &url.Userinfo{}}, &url.URL{Path: "http"}, nil))
	assert.Nil(t, ReplacePrefixURI(&url.URL{User: url.UserPassword("username", "")}, &url.URL{User: url.UserPassword("username2", "")}, nil))
	assert.Nil(t, ReplacePrefixURI(&url.URL{Path: "test/toto"}, &url.URL{Path: "info"}, &url.URL{Path: "test2"}))
	{
		uri := ReplacePrefixURI(&url.URL{Scheme: "test", Path: "toto"}, &url.URL{Scheme: "test"}, &url.URL{Scheme: "file"})
		require.NotNil(t, uri)
		assert.Equal(t, "file", uri.Scheme)
		assert.Equal(t, "toto", uri.Path)
	}
	{
		uri := ReplacePrefixURI(&url.URL{Path: ""}, &url.URL{Path: ""}, &url.URL{Path: "test"})
		require.NotNil(t, uri)
		assert.Equal(t, "test", uri.Path)
	}
	{
		uri := ReplacePrefixURI(&url.URL{Path: "toto"}, &url.URL{Path: ""}, &url.URL{Path: "test/"})
		require.NotNil(t, uri)
		assert.Equal(t, "test/toto", uri.Path)
	}
	{
		uri := ReplacePrefixURI(&url.URL{Path: "test/toto"}, &url.URL{Path: "test/toto"}, &url.URL{Path: "test2"})
		require.NotNil(t, uri)
		assert.Equal(t, "test2", uri.Path)
	}
	{
		uri := ReplacePrefixURI(&url.URL{Path: "test/toto"}, &url.URL{Path: "test"}, &url.URL{Path: "test2"})
		require.NotNil(t, uri)
		assert.Equal(t, "test2/toto", uri.Path)
	}
}
