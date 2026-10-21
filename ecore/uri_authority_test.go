package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// an empty authority ("file:///path") is kept by the uris built from a uri
func TestURI_EmptyAuthority(t *testing.T) {
	base := NewURI("file:///a/library.xml")
	assert.Equal(t, "file:///a/other.xml#//@books.0", base.Resolve(NewURI("other.xml#//@books.0")).String())
	assert.Equal(t, "file:///a/library.xml#//@books.0", base.Resolve(NewURI("#//@books.0")).String())
	assert.Equal(t, "file:///b/other.xml", base.Resolve(NewURI("/b/other.xml")).String())
	assert.Equal(t, "file:///a/library.xml", NewURI("file:///a/library.xml#f").TrimFragment().String())
	assert.Equal(t, "file:///a/library.xml", NewURI("file:///a/library.xml?q").TrimQuery().String())
	assert.Equal(t, "file:///b/library.xml", NewURI("file:///a/../b/library.xml").Normalize().String())
	assert.Equal(t, "file:///a/library.xml#f", NewURIBuilder(base).SetFragment("f").URI().String())
	assert.Equal(t, "file:///c/library.xml", base.ReplacePrefix(NewURI("file:///a"), NewURI("file:///c")).String())
	assert.True(t, base.Resolve(NewURI("#f")).TrimFragment().Equals(base))

	// without authority, and an authority given by the reference
	noAuthority := NewURI("file:/a/library.xml")
	assert.Equal(t, "file:/a/other.xml", noAuthority.Resolve(NewURI("other.xml")).String())
	assert.Equal(t, "file:/a/library.xml", NewURI("file:/a/library.xml#f").TrimFragment().String())
	assert.Equal(t, "http://h/b", NewURI("http://g/a").Resolve(NewURI("//h/b")).String())
	assert.Equal(t, "a/b", NewURIFromComponents("", "", "", "", "", "a/b", "", "").String())
}
