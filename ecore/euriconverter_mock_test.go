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
	h.EXPECT().CreateReader(uri).Once()
	h.EXPECT().CreateReader(uri).Return(f, nil).Run(func(uri *URI) { m.Run(uri) }).Once()
	h.EXPECT().CreateReader(uri).RunAndReturn(func(*URI) (io.ReadCloser, error) {
		return nil, errors.New("error")
	}).Once()
	h.EXPECT().CreateReader(uri).Call.Return(
		func(*URI) io.ReadCloser {
			return nil
		},
		func(*URI) error {
			return errors.New("error")
		},
	).Once()
	{
		assert.Panics(t, func() {
			_, _ = h.CreateReader(uri)
		})
	}
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
	h.EXPECT().CreateWriter(uri).Once()
	h.EXPECT().CreateWriter(uri).Return(f, nil).Run(func(uri *URI) { m.Run(uri) }).Once()
	h.EXPECT().CreateWriter(uri).RunAndReturn(func(*URI) (io.WriteCloser, error) {
		return nil, errors.New("error")
	}).Once()
	h.EXPECT().CreateWriter(uri).Call.Return(
		func(*URI) io.WriteCloser {
			return nil
		},
		func(*URI) error {
			return errors.New("error")
		},
	).Once()
	{
		assert.Panics(t, func() {
			_, _ = h.CreateWriter(uri)
		})
	}
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
	h.EXPECT().GetURIMap().Once()
	h.EXPECT().GetURIMap().Return(mp).Run(func() { m.Run() }).Once()
	h.EXPECT().GetURIMap().RunAndReturn(func() map[URI]URI {
		return mp
	}).Once()
	assert.Panics(t, func() {
		_ = h.GetURIMap()
	})
	assert.Equal(t, mp, h.GetURIMap())
	assert.Equal(t, mp, h.GetURIMap())
}

func TestMockEURIConverterNormalize(t *testing.T) {
	h := NewMockEURIConverter(t)
	uri1 := NewURI("test:///file.t")
	uri2 := NewURI("test:///file.t")
	m := NewMockRun(t, uri1)
	h.EXPECT().Normalize(uri1).Once()
	h.EXPECT().Normalize(uri1).Return(uri2).Run(func(uri *URI) { m.Run(uri) }).Once()
	h.EXPECT().Normalize(uri1).RunAndReturn(func(*URI) *URI { return uri2 }).Once()
	assert.Panics(t, func() {
		_ = h.Normalize(uri1)
	})
	assert.Equal(t, uri2, h.Normalize(uri1))
	assert.Equal(t, uri2, h.Normalize(uri1))
}

func TestMockEURIConverterGetURIHandler(t *testing.T) {
	h := NewMockEURIConverter(t)
	u := NewMockEURIHandler(t)
	uri := NewURI("test:///file.t")
	m := NewMockRun(t, uri)
	h.EXPECT().GetURIHandler(uri).Once()
	h.EXPECT().GetURIHandler(uri).Return(u).Run(func(uri *URI) { m.Run(uri) }).Once()
	h.EXPECT().GetURIHandler(uri).RunAndReturn(func(*URI) EURIHandler {
		return u
	}).Once()
	assert.Panics(t, func() {
		_ = h.GetURIHandler(uri)
	})
	assert.Equal(t, u, h.GetURIHandler(uri))
	assert.Equal(t, u, h.GetURIHandler(uri))
}

func TestMockEURIConverterGetURIHandlers(t *testing.T) {
	h := NewMockEURIConverter(t)
	l := NewMockEList(t)
	m := NewMockRun(t)
	h.EXPECT().GetURIHandlers().Once()
	h.EXPECT().GetURIHandlers().Return(l).Run(func() { m.Run() }).Once()
	h.EXPECT().GetURIHandlers().RunAndReturn(func() EList {
		return l
	}).Once()
	assert.Panics(t, func() {
		_ = h.GetURIHandlers()
	})
	assert.Equal(t, l, h.GetURIHandlers())
	assert.Equal(t, l, h.GetURIHandlers())
	mock.AssertExpectationsForObjects(t, h)
}
