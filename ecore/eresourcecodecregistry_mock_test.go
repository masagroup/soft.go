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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMockEResourceCodecRegistry_GetCodec(t *testing.T) {
	f := &MockEResourceCodecRegistry{}
	r := &MockEResourceCodec{}
	uri := NewURI("test:///file.t")
	f.On("GetCodec", uri).Return(r).Once()
	f.On("GetCodec", uri).Return(func(*URI) EResourceCodec {
		return r
	}).Once()
	assert.Equal(t, r, f.GetCodec(uri))
	assert.Equal(t, r, f.GetCodec(uri))
	mock.AssertExpectationsForObjects(t, f, r)
}

func TestMockEResourceCodecRegistryGetProtocolToCodecMap(t *testing.T) {
	f := &MockEResourceCodecRegistry{}
	m := make(map[string]EResourceCodec)
	f.On("GetProtocolToCodecMap").Return(m).Once()
	f.On("GetProtocolToCodecMap").Return(func() map[string]EResourceCodec {
		return m
	}).Once()
	assert.Equal(t, m, f.GetProtocolToCodecMap())
	assert.Equal(t, m, f.GetProtocolToCodecMap())
	f.AssertExpectations(t)
}

func TestMockEResourceCodecRegistryGetExtensionToCodecMap(t *testing.T) {
	f := &MockEResourceCodecRegistry{}
	m := make(map[string]EResourceCodec)
	f.On("GetExtensionToCodecMap").Return(m).Once()
	f.On("GetExtensionToCodecMap").Return(func() map[string]EResourceCodec {
		return m
	}).Once()
	assert.Equal(t, m, f.GetExtensionToCodecMap())
	assert.Equal(t, m, f.GetExtensionToCodecMap())
	f.AssertExpectations(t)
}
