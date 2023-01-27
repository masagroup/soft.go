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
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockEResourceCodec_NewDecoder(t *testing.T) {
	mockCodec := NewMockEResourceCodec(t)
	mockDecoder := NewMockEResourceDecoder(t)
	mockResource := NewMockEResource(t)
	options := map[string]any{}
	reader := strings.NewReader("")
	m := NewMockRun(t, mockResource, reader, options)
	mockCodec.EXPECT().NewDecoder(mockResource, reader, options).Return(mockDecoder).Run(func(resource EResource, r io.Reader, options map[string]interface{}) { m.Run(resource, r, options) }).Once()
	mockCodec.EXPECT().NewDecoder(mockResource, reader, options).Call.Return(func(EResource, io.Reader, map[string]any) EResourceDecoder { return mockDecoder }).Once()
	assert.Equal(t, mockDecoder, mockCodec.NewDecoder(mockResource, reader, options))
	assert.Equal(t, mockDecoder, mockCodec.NewDecoder(mockResource, reader, options))
}

func TestMockEResourceCodec_NewEncoder(t *testing.T) {
	mockCodec := &MockEResourceCodec{}
	mockEncoder := &MockEResourceEncoder{}
	mockResource := NewMockEResource(t)
	options := map[string]any{}
	writer := bytes.NewBufferString("")
	m := NewMockRun(t, mockResource, writer, options)
	mockCodec.EXPECT().NewEncoder(mockResource, writer, options).Return(mockEncoder).Run(func(resource EResource, r io.Writer, options map[string]interface{}) { m.Run(resource, r, options) }).Once()
	mockCodec.EXPECT().NewEncoder(mockResource, writer, options).Call.Return(func(EResource, io.Writer, map[string]any) EResourceEncoder { return mockEncoder }).Once()
	assert.Equal(t, mockEncoder, mockCodec.NewEncoder(mockResource, writer, options))
	assert.Equal(t, mockEncoder, mockCodec.NewEncoder(mockResource, writer, options))
}
