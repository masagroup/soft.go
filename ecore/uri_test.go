package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestURI_Constructor(t *testing.T) {
	assert.Equal(t, &URI{Scheme: "http"}, NewURI("http://"))
	assert.Equal(t, &URI{Scheme: "http", Host: "host"}, NewURI("http://host"))
	assert.Equal(t, &URI{Scheme: "http", Host: "host", Port: 10020}, NewURI("http://host:10020"))
	assert.Equal(t, &URI{Scheme: "http", Host: "host", Port: 10020, Path: "/path/path2"}, NewURI("http://host:10020/path/path2"))
	assert.Equal(t, &URI{Scheme: "http", Host: "host", Port: 10020, Path: "/path/path2", Query: "key1=foo&key2=&key3&=bar&=bar="}, NewURI("http://host:10020/path/path2?key1=foo&key2=&key3&=bar&=bar="))
	assert.Equal(t, &URI{Scheme: "http", Host: "host", Port: 10020, Path: "/path/path2", Fragment: "fragment"}, NewURI("http://host:10020/path/path2#fragment"))
	assert.Equal(t, &URI{Scheme: "file", Host: "file.txt", Query: "query", Fragment: "fragment"}, NewURI("file://file.txt?query#fragment"))
	assert.Equal(t, &URI{Scheme: "file", Path: "/file.txt", Query: "query", Fragment: "fragment"}, NewURI("file:///file.txt?query#fragment"))
	assert.Equal(t, &URI{Fragment: "fragment"}, NewURI("//#fragment"))
	assert.Equal(t, &URI{Path: "path"}, NewURI("path"))
	assert.Equal(t, &URI{Path: "./path"}, NewURI("./path"))
}

func TestURI_ParseURI(t *testing.T) {
	{
		uri, err := ParseURI("2http:///file.txt")
		assert.NotNil(t, err)
		assert.Nil(t, uri)
	}
}

func TestURI_IsAbsolute(t *testing.T) {
	assert.True(t, NewURI("http://toto").IsAbsolute())
	assert.False(t, NewURI("/toto").IsAbsolute())
}

func TestURI_IsOpaque(t *testing.T) {
	assert.True(t, NewURI("http://toto").IsOpaque())
	assert.False(t, NewURI("http://toto/").IsOpaque())
}

func TestURI_IsEmpty(t *testing.T) {
	{
		u := &URI{}
		assert.True(t, u.IsEmpty())
	}
	{
		u := &URI{Scheme: "t"}
		assert.False(t, u.IsEmpty())
	}
}

func TestURI_Copy(t *testing.T) {
	uri := &URI{
		Scheme:   "scheme",
		Username: "username",
		Password: "password",
		Host:     "host",
		Port:     10,
		Path:     "path",
		Query:    "query",
		Fragment: "fragment",
	}
	assert.Equal(t, uri, uri.Copy())
}

func TestURI_Normalize(t *testing.T) {
	assert.Equal(t, NewURI("http://host:10020"), NewURI("http://host:10020").Normalize())
	assert.Equal(t, NewURI("http://host:10020/"), NewURI("http://host:10020/").Normalize())
	assert.Equal(t, NewURI("http://host:10020/path"), NewURI("http://host:10020/path").Normalize())
	assert.Equal(t, NewURI("http://host:10020/path"), NewURI("http://host:10020/./path").Normalize())
	assert.Equal(t, NewURI("http://host:10020/path2"), NewURI("http://host:10020/path/../path2").Normalize())
	assert.Equal(t, NewURI("http://host:10020/path/path2"), NewURI("http://host:10020/path/./path2").Normalize())
}

func TestURI_Resolve(t *testing.T) {
	assert.Equal(t, NewURI("http://host:10020/path2/"), NewURI("http://host:10020/path/").Resolve(NewURI("http://host:10020/path2/")))
	assert.Equal(t, NewURI("http://host:10020/path2"), NewURI("http://host:10020/path/").Resolve(NewURI("../path2")))
	assert.Equal(t, NewURI("http://host:10020/path2"), NewURI("http://host:10020/path/").Resolve(NewURI("/path2")))
	assert.Equal(t, NewURI("http://host:10020/path/path2"), NewURI("http://host:10020/path/").Resolve(NewURI("./path2")))
	assert.Equal(t, NewURI("path/path3"), NewURI("path/path2").Resolve(NewURI("path3")))
}

func TestURI_Relativize(t *testing.T) {
	assert.Equal(t, NewURI("path2"), NewURI("http://host:10020/path/").Relativize(NewURI("http://host:10020/path/path2")))
	assert.Equal(t, NewURI("path1"), NewURI("testdata/path2").Relativize(NewURI("testdata/path1")))
}

func TestURI_ReplacePrefix(t *testing.T) {
	assert.Nil(t, NewURI("http://").ReplacePrefix(NewURI("file://"), nil))
	assert.Nil(t, NewURI("http://host").ReplacePrefix(NewURI("http://host2/path"), nil))
	assert.Nil(t, NewURI("http://host").ReplacePrefix(NewURI("http://host2/path"), nil))
	assert.Nil(t, NewURI("test/toto").ReplacePrefix(NewURI("info"), nil))
	{
		uri := NewURI("test:///toto").ReplacePrefix(&URI{Scheme: "test"}, &URI{Scheme: "file"})
		require.NotNil(t, uri)
		assert.Equal(t, "file", uri.Scheme)
	}
	{
		uri := NewURI("").ReplacePrefix(&URI{}, &URI{Path: "file"})
		require.NotNil(t, uri)
		assert.Equal(t, "file", uri.Path)
	}
	{
		uri := (&URI{Path: "toto"}).ReplacePrefix(&URI{}, &URI{Path: "test/"})
		require.NotNil(t, uri)
		assert.Equal(t, "test/toto", uri.Path)
	}
	{
		uri := (&URI{Path: "test/toto"}).ReplacePrefix(&URI{Path: "test/toto"}, &URI{Path: "test2"})
		require.NotNil(t, uri)
		assert.Equal(t, "test2", uri.Path)
	}
	{
		uri := (&URI{Path: "test/toto"}).ReplacePrefix(&URI{Path: "test"}, &URI{Path: "test2"})
		require.NotNil(t, uri)
		assert.Equal(t, "test2/toto", uri.Path)
	}
}
