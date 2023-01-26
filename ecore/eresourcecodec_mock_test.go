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
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMockEResourceCodec_NewDecoder(t *testing.T) {
	mockCodec := &MockEResourceCodec{}
	mockDecoder := &MockEResourceDecoder{}
	mockResource := NewMockEResource(t)
	options := map[string]any{}
	var r io.Reader = nil
	mockCodec.On("NewDecoder", mockResource, r, options).Return(mockDecoder).Once()
	mockCodec.On("NewDecoder", mockResource, r, options).Return(func(EResource, io.Reader, map[string]any) EResourceDecoder { return mockDecoder }).Once()
	assert.Equal(t, mockDecoder, mockCodec.NewDecoder(mockResource, r, options))
	assert.Equal(t, mockDecoder, mockCodec.NewDecoder(mockResource, r, options))
	mock.AssertExpectationsForObjects(t, mockCodec)
}

func TestMockEResourceCodec_NewEncoder(t *testing.T) {
	mockCodec := &MockEResourceCodec{}
	mockEncoder := &MockEResourceEncoder{}
	mockResource := NewMockEResource(t)
	options := map[string]any{}
	var w io.Writer = nil
	mockCodec.On("NewEncoder", mockResource, w, options).Return(mockEncoder).Once()
	mockCodec.On("NewEncoder", mockResource, w, options).Return(func(EResource, io.Writer, map[string]any) EResourceEncoder { return mockEncoder }).Once()
	assert.Equal(t, mockEncoder, mockCodec.NewEncoder(mockResource, w, options))
	assert.Equal(t, mockEncoder, mockCodec.NewEncoder(mockResource, w, options))
	mock.AssertExpectationsForObjects(t, mockCodec)
}
