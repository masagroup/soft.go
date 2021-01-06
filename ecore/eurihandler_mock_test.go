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
	"io"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMockEURIHandlerCanHandle(t *testing.T) {
	h := &MockEURIHandler{}
	uri, _ := url.Parse("test://file.t")
	h.On("CanHandle", uri).Return(true).Once()
	h.On("CanHandle", uri).Return(func(*url.URL) bool {
		return false
	}).Once()
	assert.True(t, h.CanHandle(uri))
	assert.False(t, h.CanHandle(uri))
	mock.AssertExpectationsForObjects(t, h)
}

func TestMockEURIHandlerCreateReader(t *testing.T) {
	h := &MockEURIHandler{}
	uri, _ := url.Parse("test://file.t")
	f, _ := os.Open(uri.String())
	h.On("CreateReader", uri).Return(f).Once()
	h.On("CreateReader", uri).Return(func(*url.URL) io.ReadCloser {
		return f
	}).Once()
	assert.Equal(t, f, h.CreateReader(uri))
	assert.Equal(t, f, h.CreateReader(uri))
	mock.AssertExpectationsForObjects(t, h)
}

func TestMockEURIHandlerCreateWriter(t *testing.T) {
	h := &MockEURIHandler{}
	uri, _ := url.Parse("test://file.t")
	f, _ := os.Create(uri.String())
	h.On("CreateWriter", uri).Return(f).Once()
	h.On("CreateWriter", uri).Return(func(*url.URL) io.WriteCloser {
		return f
	}).Once()
	assert.Equal(t, f, h.CreateWriter(uri))
	assert.Equal(t, f, h.CreateWriter(uri))
	mock.AssertExpectationsForObjects(t, h)
}
