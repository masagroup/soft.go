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

func TestMockEResourceCodec_NewDecoder(t *testing.T) {
	mockCodec := &MockEResourceCodec{}
	mockDecoder := &MockEResourceDecoder{}
	options := map[string]interface{}{}
	mockCodec.On("NewDecoder", options).Return(mockDecoder).Once()
	mockCodec.On("NewDecoder", options).Return(func(map[string]interface{}) EResourceDecoder { return mockDecoder }).Once()
	assert.Equal(t, mockDecoder, mockCodec.NewDecoder(options))
	assert.Equal(t, mockDecoder, mockCodec.NewDecoder(options))
	mock.AssertExpectationsForObjects(t, mockCodec)
}

func TestMockEResourceCodec_NewEncoder(t *testing.T) {
	mockCodec := &MockEResourceCodec{}
	mockEncoder := &MockEResourceEncoder{}
	options := map[string]interface{}{}
	mockCodec.On("NewEncoder", options).Return(mockEncoder).Once()
	mockCodec.On("NewEncoder", options).Return(func(map[string]interface{}) EResourceEncoder { return mockEncoder }).Once()
	assert.Equal(t, mockEncoder, mockCodec.NewEncoder(options))
	assert.Equal(t, mockEncoder, mockCodec.NewEncoder(options))
	mock.AssertExpectationsForObjects(t, mockCodec)
}
