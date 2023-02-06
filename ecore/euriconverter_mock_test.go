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
	"errors"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMockEURIConverterCreateReader(t *testing.T) {
	h := NewMockEURIConverter(t)
	uri := NewURI("test:///file.t")
	f, _ := os.Open(uri.String())
	m := NewMockRun(t, uri)
	h.EXPECT().CreateReader(uri).Return(f, nil).Run(func(uri *URI) { m.Run(uri) }).Once()
	h.EXPECT().CreateReader(uri).Call.Return(func(*URI) (io.ReadCloser, error) {
		return nil, errors.New("error")
	}).Once()
	{
		r, err := h.CreateReader(uri)
		assert.Equal(t, f, r)
		assert.Nil(t, err)
	}
	{
		r, err := h.CreateReader(uri)
		assert.Equal(t, nil, r)
		assert.NotNil(t, err)
		assert.Equal(t, "error", err.Error())
	}
}

func TestMockEURIConverterCreateWriter(t *testing.T) {
	h := NewMockEURIConverter(t)
	uri := NewURI("test:///file.t")
	f, _ := os.Create(uri.String())
	m := NewMockRun(t, uri)
	h.EXPECT().CreateWriter(uri).Return(f, nil).Run(func(uri *URI) { m.Run(uri) }).Once()
	h.EXPECT().CreateWriter(uri).Call.Return(func(*URI) (io.WriteCloser, error) {
		return nil, errors.New("error")
	}).Once()
	{
		r, err := h.CreateWriter(uri)
		assert.Equal(t, f, r)
		assert.Nil(t, err)
	}
	{
		r, err := h.CreateWriter(uri)
		assert.Equal(t, nil, r)
		assert.NotNil(t, err)
		assert.Equal(t, "error", err.Error())
	}
}

func TestMockEURIConverterGetURIMap(t *testing.T) {
	h := NewMockEURIConverter(t)
	mp := map[URI]URI{*NewURI("toto"): *NewURI("tata")}
	m := NewMockRun(t)
	h.EXPECT().GetURIMap().Return(mp).Run(func() { m.Run() }).Once()
	h.EXPECT().GetURIMap().Call.Return(func() map[URI]URI {
		return mp
	}).Once()
	assert.Equal(t, mp, h.GetURIMap())
	assert.Equal(t, mp, h.GetURIMap())
}

func TestMockEURIConverterNormalize(t *testing.T) {
	h := NewMockEURIConverter(t)
	uri1 := NewURI("test:///file.t")
	uri2 := NewURI("test:///file.t")
	m := NewMockRun(t, uri1)
	h.EXPECT().Normalize(uri1).Return(uri2).Run(func(uri *URI) { m.Run(uri) }).Once()
	h.EXPECT().Normalize(uri1).Call.Return(func(*URI) *URI { return uri2 }).Once()
	assert.Equal(t, uri2, h.Normalize(uri1))
	assert.Equal(t, uri2, h.Normalize(uri1))
}

func TestMockEURIConverterGetURIHandler(t *testing.T) {
	h := NewMockEURIConverter(t)
	u := NewMockEURIHandler(t)
	uri := NewURI("test:///file.t")
	m := NewMockRun(t, uri)
	h.EXPECT().GetURIHandler(uri).Return(u).Run(func(uri *URI) { m.Run(uri) }).Once()
	h.EXPECT().GetURIHandler(uri).Call.Return(func(*URI) EURIHandler {
		return u
	}).Once()
	assert.Equal(t, u, h.GetURIHandler(uri))
	assert.Equal(t, u, h.GetURIHandler(uri))
}

func TestMockEURIConverterGetURIHandlers(t *testing.T) {
	h := NewMockEURIConverter(t)
	l := NewMockEList(t)
	m := NewMockRun(t)
	h.EXPECT().GetURIHandlers().Return(l).Run(func() { m.Run() }).Once()
	h.EXPECT().GetURIHandlers().Call.Return(func() EList {
		return l
	}).Once()
	assert.Equal(t, l, h.GetURIHandlers())
	assert.Equal(t, l, h.GetURIHandlers())
	mock.AssertExpectationsForObjects(t, h)
}
