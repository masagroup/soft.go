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
	"github.com/stretchr/testify/mock"
)

func TestMockEResourceFactoryRegistryGetFactory(t *testing.T) {
	f := &MockEResourceFactoryRegistry{}
	r := &MockEResourceFactory{}
	uri, _ := url.Parse("test://file.t")
	f.On("GetFactory", uri).Return(r).Once()
	f.On("GetFactory", uri).Return(func(*url.URL) EResourceFactory {
		return r
	}).Once()
	assert.Equal(t, r, f.GetFactory(uri))
	assert.Equal(t, r, f.GetFactory(uri))
	mock.AssertExpectationsForObjects(t, f, r)
}

func TestMockEResourceFactoryRegistryGetProtocolToFactoryMap(t *testing.T) {
	f := &MockEResourceFactoryRegistry{}
	m := make(map[string]EResourceFactory)
	f.On("GetProtocolToFactoryMap").Return(m).Once()
	f.On("GetProtocolToFactoryMap").Return(func() map[string]EResourceFactory {
		return m
	}).Once()
	assert.Equal(t, m, f.GetProtocolToFactoryMap())
	assert.Equal(t, m, f.GetProtocolToFactoryMap())
	f.AssertExpectations(t)
}

func TestMockEResourceFactoryRegistryGetExtensionToFactoryMap(t *testing.T) {
	f := &MockEResourceFactoryRegistry{}
	m := make(map[string]EResourceFactory)
	f.On("GetExtensionToFactoryMap").Return(m).Once()
	f.On("GetExtensionToFactoryMap").Return(func() map[string]EResourceFactory {
		return m
	}).Once()
	assert.Equal(t, m, f.GetExtensionToFactoryMap())
	assert.Equal(t, m, f.GetExtensionToFactoryMap())
	f.AssertExpectations(t)
}
