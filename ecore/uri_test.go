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
	assert.Equal(t, NewURIBuilder(nil).SetScheme("http").URI(), NewURI("http:"))
	assert.Equal(t, NewURIBuilder(nil).SetScheme("http").SetHost("host").URI(), NewURI("http://host"))
	assert.Equal(t, NewURIBuilder(nil).SetScheme("http").SetHost("host").SetPort("10020").URI(), NewURI("http://host:10020"))
	assert.Equal(t, NewURIBuilder(nil).SetScheme("http").SetHost("host").SetPort("10020").SetPath("/path/path2").URI(), NewURI("http://host:10020/path/path2"))
	assert.Equal(t, NewURIBuilder(nil).SetScheme("http").SetHost("host").SetPort("10020").SetPath("/path/path2").SetQuery("key1=foo&key2=&key3&=bar&=bar=").URI(), NewURI("http://host:10020/path/path2?key1=foo&key2=&key3&=bar&=bar="))
	assert.Equal(t, NewURIBuilder(nil).SetScheme("http").SetHost("host").SetPort("10020").SetPath("/path/path2").SetFragment("fragment").URI(), NewURI("http://host:10020/path/path2#fragment"))
	assert.Equal(t, NewURIBuilder(nil).SetScheme("file").SetHost("file.txt").SetQuery("query").SetFragment("fragment").URI(), NewURI("file://file.txt?query#fragment"))
	assert.Equal(t, NewURIBuilder(nil).SetScheme("file").SetPath("/file.txt").SetQuery("query").SetFragment("fragment").URI(), NewURI("file:/file.txt?query#fragment"))
	assert.Equal(t, NewURIBuilder(nil).SetFragment("fragment").URI(), NewURI("#fragment"))
	assert.Equal(t, NewURIBuilder(nil).SetPath("path").URI(), NewURI("path"))
	assert.Equal(t, NewURIBuilder(nil).SetPath("./path").URI(), NewURI("./path"))
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
		assert.Equal(t, "file", uri.Scheme())
	}
	{
		uri := NewURI("").ReplacePrefix(NewURI(""), NewURI("file"))
		require.NotNil(t, uri)
		assert.Equal(t, "file", uri.Path())
	}
	{
		uri := NewURI("toto").ReplacePrefix(NewURI(""), NewURI("test/"))
		require.NotNil(t, uri)
		assert.Equal(t, "test/toto", uri.Path())
	}
	{
		uri := NewURI("test/toto").ReplacePrefix(NewURI("test/toto"), NewURI("test2"))
		require.NotNil(t, uri)
		assert.Equal(t, "test2", uri.Path())
	}
	{
		uri := NewURI("test/toto").ReplacePrefix(NewURI("test"), NewURI("test2"))
		require.NotNil(t, uri)
		assert.Equal(t, "test2/toto", uri.Path())
	}
}

func TestCreateFileURI(t *testing.T) {
	assert.Equal(t, &URI{}, CreateFileURI(""))
	if runtime.GOOS == "windows" {
		assert.Equal(t, NewURI("test/toto"), CreateFileURI("test\\toto"))
		assert.Equal(t, NewURI("file:D:/test/toto"), CreateFileURI("D:\\test\\toto"))
	} else {
		assert.Equal(t, NewURI("test/toto"), CreateFileURI("test/toto"))
		assert.Equal(t, NewURI("file:/test/toto"), CreateFileURI("/test/toto"))
	}
}

func TestCreateMemoryURI(t *testing.T) {
	assert.Nil(t, CreateMemoryURI(""))
	assert.Equal(t, NewURI("memory:path"), CreateMemoryURI("path"))
	assert.Equal(t, "memory:path", CreateMemoryURI("path").String())
}

func TestURI_Authority(t *testing.T) {
	assert.Equal(t, "", NewURI("http:///file.text").Authority())
	assert.Equal(t, "", NewURI("http:/file.text").Authority())
	assert.Equal(t, "host", NewURI("http://host/file.text").Authority())
	assert.Equal(t, "host:10", NewURI("http://host:10/file.text").Authority())
	assert.Equal(t, "userinfo@host:10", NewURI("http://userinfo@host:10/file.text").Authority())
}

// func TestURI_Cache(t *testing.T) {
// 	var p1, p2 uintptr
// 	{
// 		uri1 := NewURI("http://toto")
// 		uri2 := NewURI("http://toto")
// 		assert.Same(t, uri1, uri2)
// 		p1 = uintptr(unsafe.Pointer(uri1))
// 	}
// 	runtime.GC()
// 	{
// 		uri1 := NewURI("http://toto")
// 		p2 = uintptr(unsafe.Pointer(uri1))
// 	}
// 	assert.NotEqual(t, p1, p2)
// }
