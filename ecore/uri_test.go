// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestURI_Constructor(t *testing.T) {
	assert.Equal(t, &URI{scheme: "http"}, NewURI("http://"))
	assert.Equal(t, &URI{scheme: "http", host: "host"}, NewURI("http://host"))
	assert.Equal(t, &URI{scheme: "http", host: "host", port: 10020}, NewURI("http://host:10020"))
	assert.Equal(t, &URI{scheme: "http", host: "host", port: 10020, Path: "/path/path2"}, NewURI("http://host:10020/path/path2"))
	assert.Equal(t, &URI{scheme: "http", host: "host", port: 10020, Path: "/path/path2", query: "key1=foo&key2=&key3&=bar&=bar="}, NewURI("http://host:10020/path/path2?key1=foo&key2=&key3&=bar&=bar="))
	assert.Equal(t, &URI{scheme: "http", host: "host", port: 10020, Path: "/path/path2", fragment: "fragment"}, NewURI("http://host:10020/path/path2#fragment"))
	assert.Equal(t, &URI{scheme: "file", host: "file.txt", query: "query", fragment: "fragment"}, NewURI("file://file.txt?query#fragment"))
	assert.Equal(t, &URI{scheme: "file", Path: "/file.txt", query: "query", fragment: "fragment"}, NewURI("file:///file.txt?query#fragment"))
	assert.Equal(t, &URI{fragment: "fragment"}, NewURI("//#fragment"))
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
		u := &URI{scheme: "t"}
		assert.False(t, u.IsEmpty())
	}
}

func TestURI_Copy(t *testing.T) {
	uri := &URI{
		scheme:   "scheme",
		username: "username",
		password: "password",
		host:     "host",
		port:     10,
		Path:     "path",
		query:    "query",
		fragment: "fragment",
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
		uri := NewURI("test:///toto").ReplacePrefix(&URI{scheme: "test"}, &URI{scheme: "file"})
		require.NotNil(t, uri)
		assert.Equal(t, "file", uri.scheme)
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

func TestCreateFileURI(t *testing.T) {
	assert.Equal(t, &URI{}, CreateFileURI(""))
	if runtime.GOOS == "windows" {
		assert.Equal(t, &URI{Path: "test/toto"}, CreateFileURI("test\\toto"))
		assert.Equal(t, &URI{scheme: "file", Path: "D:/test/toto"}, CreateFileURI("D:\\test\\toto"))
	} else {
		assert.Equal(t, &URI{Path: "test/toto"}, CreateFileURI("test/toto"))
		assert.Equal(t, &URI{scheme: "file", Path: "/test/toto"}, CreateFileURI("/test/toto"))
	}
}

func TestCreateMemoryURI(t *testing.T) {
	assert.Nil(t, CreateMemoryURI(""))
	assert.Equal(t, &URI{scheme: "memory", Path: "path"}, CreateMemoryURI("path"))
	assert.Equal(t, "memory:path", CreateMemoryURI("path").String())
}
