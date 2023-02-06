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
)

func TestMockEURIHandlerCanHandle(t *testing.T) {
	h := NewMockEURIHandler(t)
	uri := NewURI("test:///file.t")
	m := NewMockRun(t, uri)
	h.EXPECT().CanHandle(uri).Return(true).Run(func(uri *URI) { m.Run(uri) }).Once()
	h.EXPECT().CanHandle(uri).Call.Return(func(*URI) bool {
		return false
	}).Once()
	assert.True(t, h.CanHandle(uri))
	assert.False(t, h.CanHandle(uri))
}

func TestMockEURIHandlerCreateReader(t *testing.T) {
	h := NewMockEURIHandler(t)
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

func TestMockEURIHandlerCreateWriter(t *testing.T) {
	h := NewMockEURIHandler(t)
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
