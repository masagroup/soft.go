// *****************************************************************************
//
// This file is part of a MASA library or program.
// Refer to the included end-user license agreement for restrictions.
//
// Copyright (c) 2020 MASA Group
//
// *****************************************************************************

package ecore

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMockEResourceGetResourceSet tests method GetResourceSet
func TestMockEResourceGetResourceSet(t *testing.T) {
	r := NewMockEResource(t)
	set := NewMockEResourceSet(t)
	m := NewMockRun(t)
	r.EXPECT().GetResourceSet().Return(set).Run(func() { m.Run() }).Once()
	r.EXPECT().GetResourceSet().Call.Return(func() EResourceSet {
		return set
	}).Once()
	assert.Equal(t, set, r.GetResourceSet())
	assert.Equal(t, set, r.GetResourceSet())
}

// TestMockEResourceGetURI tests method GetURI
func TestMockEResourceGetURI(t *testing.T) {
	r := NewMockEResource(t)
	uri := &URI{}
	m := NewMockRun(t)
	r.EXPECT().GetURI().Return(uri).Run(func() { m.Run() }).Once()
	r.EXPECT().GetURI().Call.Return(func() *URI {
		return uri
	}).Once()
	assert.Equal(t, uri, r.GetURI())
	assert.Equal(t, uri, r.GetURI())
}

// TestMockEResourceSetURI tests method SetURI
func TestMockEResourceSetURI(t *testing.T) {
	r := NewMockEResource(t)
	uri := &URI{}
	m := NewMockRun(t, uri)
	r.EXPECT().SetURI(uri).Return().Run(func(_a0 *URI) { m.Run(uri) }).Once()
	r.SetURI(uri)
}

// TestMockEResourceGetContents tests method GetContents
func TestMockEResourceGetContents(t *testing.T) {
	r := NewMockEResource(t)
	c := NewMockEList(t)
	m := NewMockRun(t)
	r.EXPECT().GetContents().Return(c).Run(func() { m.Run() }).Once()
	r.EXPECT().GetContents().Call.Return(func() EList {
		return c
	}).Once()
	assert.Equal(t, c, r.GetContents())
	assert.Equal(t, c, r.GetContents())
}

// TestMockEResourceGetAllContents tests method GetAllContents
func TestMockEResourceGetAllContents(t *testing.T) {
	r := NewMockEResource(t)
	i := NewMockEIterator(t)
	m := NewMockRun(t)
	r.EXPECT().GetAllContents().Return(i).Run(func() { m.Run() }).Once()
	r.EXPECT().GetAllContents().Call.Return(func() EIterator {
		return i
	}).Once()
	assert.Equal(t, i, r.GetAllContents())
	assert.Equal(t, i, r.GetAllContents())
}

// TestMockEResourceAttached tests method Attached
func TestMockEResourceAttached(t *testing.T) {
	r := NewMockEResource(t)
	o := NewMockEObject(t)
	m := NewMockRun(t, o)
	r.EXPECT().Attached(o).Return().Run(func(object EObject) { m.Run(object) }).Once()
	r.Attached(o)
}

// TestMockEResourceAttached tests method Detached
func TestMockEResourceDetached(t *testing.T) {
	r := NewMockEResource(t)
	o := NewMockEObject(t)
	m := NewMockRun(t, o)
	r.EXPECT().Detached(o).Return().Run(func(object EObject) { m.Run(object) }).Once()
	r.Detached(o)
}

// TestMockEResourceGetEObject tests method GetEObject
func TestMockEResourceGetEObject(t *testing.T) {
	r := NewMockEResource(t)
	o := NewMockEObject(t)
	m := NewMockRun(t, "test")
	r.EXPECT().GetEObject("test").Return(o).Run(func(_a0 string) { m.Run(_a0) }).Once()
	r.EXPECT().GetEObject("test").Call.Return(func(string) EObject {
		return o
	}).Once()
	assert.Equal(t, o, r.GetEObject("test"))
	assert.Equal(t, o, r.GetEObject("test"))
}

// TestMockEResourceGetURIFragment tests method GetURIFragment
func TestMockEResourceGetURIFragment(t *testing.T) {
	r := NewMockEResource(t)
	o := NewMockEObject(t)
	m := NewMockRun(t, o)
	r.EXPECT().GetURIFragment(o).Return("fragment1").Run(func(_a0 EObject) { m.Run(_a0) }).Once()
	r.EXPECT().GetURIFragment(o).Call.Return(func(EObject) string {
		return "fragment2"
	}).Once()
	assert.Equal(t, "fragment1", r.GetURIFragment(o))
	assert.Equal(t, "fragment2", r.GetURIFragment(o))
}

// TestMockEResourceLoad tests method Load
func TestMockEResourceLoad(t *testing.T) {
	r := NewMockEResource(t)
	m := NewMockRun(t)
	r.EXPECT().Load().Return().Run(func() { m.Run() }).Once()
	r.Load()
}

func TestMockEResourceLoadWithOptions(t *testing.T) {
	r := NewMockEResource(t)
	options := make(map[string]any)
	m := NewMockRun(t, options)
	r.EXPECT().LoadWithOptions(options).Return().Run(func(options map[string]interface{}) { m.Run(options) }).Once()
	r.LoadWithOptions(options)
}

// TestMockEResourceLoadWithReader tests method LoadWithReader
func TestMockEResourceLoadWithReader(t *testing.T) {
	r := NewMockEResource(t)
	reader := strings.NewReader("")
	options := make(map[string]any)
	m := NewMockRun(t, reader, options)
	r.EXPECT().LoadWithReader(reader, options).Return().Run(func(r io.Reader, options map[string]interface{}) { m.Run(r, options) }).Once()
	r.LoadWithReader(reader, options)
}

// TestMockEResourceUnLoad tests method UnLoad
func TestMockEResourceUnLoad(t *testing.T) {
	r := NewMockEResource(t)
	m := NewMockRun(t)
	r.EXPECT().Unload().Return().Run(func() { m.Run() }).Once()
	r.Unload()
}

// TestMockEResourceSave tests method Save
func TestMockEResourceSave(t *testing.T) {
	r := NewMockEResource(t)
	m := NewMockRun(t)
	r.EXPECT().Save().Return().Run(func() { m.Run() }).Once()
	r.Save()
}

func TestMockEResourceSaveWithOptions(t *testing.T) {
	r := NewMockEResource(t)
	options := make(map[string]any)
	m := NewMockRun(t, options)
	r.EXPECT().SaveWithOptions(options).Return().Run(func(options map[string]interface{}) { m.Run(options) }).Once()
	r.SaveWithOptions(options)
}

// TestMockEResourceSaveWithReader tests method SaveWithReader
func TestMockEResourceSaveWithWriter(t *testing.T) {
	r := NewMockEResource(t)
	options := make(map[string]any)
	writer := bytes.NewBufferString("")
	m := NewMockRun(t, writer, options)
	r.EXPECT().SaveWithWriter(writer, options).Return().Run(func(w io.Writer, options map[string]interface{}) { m.Run(w, options) }).Once()
	r.SaveWithWriter(writer, options)
}

// TestMockEResourceGetErrors tests method GetErrors
func TestMockEResourceGetErrors(t *testing.T) {
	r := NewMockEResource(t)
	c := NewMockEList(t)
	m := NewMockRun(t)
	r.EXPECT().GetErrors().Return(c).Run(func() { m.Run() }).Once()
	r.EXPECT().GetErrors().Call.Return(func() EList {
		return c
	}).Once()
	assert.Equal(t, c, r.GetErrors())
	assert.Equal(t, c, r.GetErrors())
}

// TestMockEResourceGetWarnings tests method GetWarnings
func TestMockEResourceGetWarnings(t *testing.T) {
	r := NewMockEResource(t)
	c := NewMockEList(t)
	m := NewMockRun(t)
	r.EXPECT().GetWarnings().Return(c).Run(func() { m.Run() }).Once()
	r.EXPECT().GetWarnings().Call.Return(func() EList {
		return c
	}).Once()
	assert.Equal(t, c, r.GetWarnings())
	assert.Equal(t, c, r.GetWarnings())
}

// TestMockEResourceIsLoaded tests method IsLoaded
func TestMockEResourceIsLoaded(t *testing.T) {
	r := NewMockEResource(t)
	m := NewMockRun(t)
	r.EXPECT().IsLoaded().Return(true).Run(func() { m.Run() }).Once()
	r.EXPECT().IsLoaded().Call.Return(func() bool {
		return false
	}).Once()
	assert.True(t, r.IsLoaded())
	assert.False(t, r.IsLoaded())
}

// TestMockEResourceIsLoading tests method IsLoaded
func TestMockEResourceIsLoading(t *testing.T) {
	r := NewMockEResource(t)
	m := NewMockRun(t)
	r.EXPECT().IsLoading().Return(true).Run(func() { m.Run() }).Once()
	r.EXPECT().IsLoading().Call.Return(func() bool {
		return false
	}).Once()
	assert.True(t, r.IsLoading())
	assert.False(t, r.IsLoading())
}

func TestMockEResourceGetIDManager(t *testing.T) {
	r := NewMockEResource(t)
	id := NewMockEObjectIDManager(t)
	m := NewMockRun(t)
	r.EXPECT().GetObjectIDManager().Return(id).Run(func() { m.Run() }).Once()
	r.EXPECT().GetObjectIDManager().Call.Return(func() EObjectIDManager {
		return id
	}).Once()
	assert.Equal(t, id, r.GetObjectIDManager())
	assert.Equal(t, id, r.GetObjectIDManager())
}

func TestMockEResourceSetIDManager(t *testing.T) {
	r := NewMockEResource(t)
	id := NewMockEObjectIDManager(t)
	m := NewMockRun(t, id)
	r.EXPECT().SetObjectIDManager(id).Return().Run(func(_a0 EObjectIDManager) { m.Run(_a0) }).Once()
	r.SetObjectIDManager(id)
}

func TestMockEResourceGetResourceListeners(t *testing.T) {
	r := NewMockEResource(t)
	l := NewMockEList(t)
	m := NewMockRun(t)
	r.EXPECT().GetResourceListeners().Return(l).Run(func() { m.Run() }).Once()
	r.EXPECT().GetResourceListeners().Call.Return(func() EList {
		return l
	}).Once()
	assert.Equal(t, l, r.GetResourceListeners())
	assert.Equal(t, l, r.GetResourceListeners())
}
