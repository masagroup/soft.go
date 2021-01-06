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
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockEResourceFactoryCreateResource(t *testing.T) {
	f := &MockEResourceFactory{}
	r := &MockEResource{}
	uri, _ := url.Parse("test://file.t")
	f.On("CreateResource", uri).Return(r).Once()
	f.On("CreateResource", uri).Return(func(*url.URL) EResource {
		return r
	}).Once()
	assert.Equal(t, r, f.CreateResource(uri))
	assert.Equal(t, r, f.CreateResource(uri))
	r.AssertExpectations(t)
}
