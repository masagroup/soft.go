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

func TestMockEURIHandlerCreateWriter(t *testing.T) {
	h := &MockEURIHandler{}
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
