// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
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
