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
)

func TestMockECodecRegistry_GetCodec(t *testing.T) {
	r := NewMockECodecRegistry(t)
	c := NewMockECodec(t)
	uri := NewURI("test:///file.t")
	m := NewMockRun(t, uri)
	r.EXPECT().GetCodec(uri).Return(c).Run(func(uri *URI) { m.Run(uri) }).Once()
	r.EXPECT().GetCodec(uri).Call.Return(func(*URI) ECodec {
		return c
	}).Once()
	assert.Equal(t, c, r.GetCodec(uri))
	assert.Equal(t, c, r.GetCodec(uri))
}

func TestMockECodecRegistryGetProtocolToCodecMap(t *testing.T) {
	r := NewMockECodecRegistry(t)
	m := make(map[string]ECodec)
	mr := NewMockRun(t)
	r.EXPECT().GetProtocolToCodecMap().Return(m).Run(func() { mr.Run() }).Once()
	r.EXPECT().GetProtocolToCodecMap().Call.Return(func() map[string]ECodec {
		return m
	}).Once()
	assert.Equal(t, m, r.GetProtocolToCodecMap())
	assert.Equal(t, m, r.GetProtocolToCodecMap())
}

func TestMockECodecRegistryGetExtensionToCodecMap(t *testing.T) {
	r := NewMockECodecRegistry(t)
	m := make(map[string]ECodec)
	mr := NewMockRun(t)
	r.EXPECT().GetExtensionToCodecMap().Return(m).Run(func() { mr.Run() }).Once()
	r.EXPECT().GetExtensionToCodecMap().Call.Return(func() map[string]ECodec {
		return m
	}).Once()
	assert.Equal(t, m, r.GetExtensionToCodecMap())
	assert.Equal(t, m, r.GetExtensionToCodecMap())
}
