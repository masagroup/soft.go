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
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMockEURIConverterCreateReader(t *testing.T) {
	h := &MockEURIConverter{}
	uri, _ := url.Parse("test://file.t")
	f, _ := os.Open(uri.String())
	h.On("CreateReader", uri).Return(f, nil).Once()
	h.On("CreateReader", uri).Return(func(*url.URL) (io.ReadCloser, error) {
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
	mock.AssertExpectationsForObjects(t, h)
}

func TestMockEURIConverterCreateWriter(t *testing.T) {
	h := &MockEURIConverter{}
	uri, _ := url.Parse("test://file.t")
	f, _ := os.Create(uri.String())
	h.On("CreateWriter", uri).Return(f, nil).Once()
	h.On("CreateWriter", uri).Return(func(*url.URL) (io.WriteCloser, error) {
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
	mock.AssertExpectationsForObjects(t, h)
}

func TestMockEURIConverterNormalize(t *testing.T) {
	h := &MockEURIConverter{}
	uri1, _ := url.Parse("test://file.t")
	uri2, _ := url.Parse("test://file.t")
	h.On("Normalize", uri1).Return(uri2).Once()
	h.On("Normalize", uri1).Return(func(*url.URL) *url.URL {
		return uri2
	}).Once()
	assert.Equal(t, uri2, h.Normalize(uri1))
	assert.Equal(t, uri2, h.Normalize(uri1))
	mock.AssertExpectationsForObjects(t, h)
}

func TestMockEURIConverterGetURIHandler(t *testing.T) {
	h := &MockEURIConverter{}
	u := &MockEURIHandler{}
	uri, _ := url.Parse("test://file.t")
	h.On("GetURIHandler", uri).Return(u).Once()
	h.On("GetURIHandler", uri).Return(func(*url.URL) EURIHandler {
		return u
	}).Once()
	assert.Equal(t, u, h.GetURIHandler(uri))
	assert.Equal(t, u, h.GetURIHandler(uri))
	mock.AssertExpectationsForObjects(t, h, u)
}

func TestMockEURIConverterGetURIHandlers(t *testing.T) {
	h := &MockEURIConverter{}
	l := &MockEList{}
	h.On("GetURIHandlers").Return(l).Once()
	h.On("GetURIHandlers").Return(func() EList {
		return l
	}).Once()
	assert.Equal(t, l, h.GetURIHandlers())
	assert.Equal(t, l, h.GetURIHandlers())
	mock.AssertExpectationsForObjects(t, h)
}
