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

func TestMockEResourceDecoder_DecodeResource(t *testing.T) {
	mockDecoder := &MockEResourceDecoder{}
	mockResource := &MockEResource{}
	mockDecoder.On("DecodeResource", mockResource).Once()
	mockDecoder.DecodeResource(mockResource)
}

func TestMockEResourceDecoder_DecodeObject(t *testing.T) {
	mockDecoder := &MockEResourceDecoder{}
	mockResource := &MockEResource{}
	var mockObject EObject
	mockDecoder.On("DecodeObject", &mockObject, mockResource).Return(nil).Once()
	mockDecoder.On("DecodeObject", &mockObject, mockResource).Return(func(*EObject, EResource) error {
		return nil
	}).Once()
	assert.Nil(t, mockDecoder.DecodeObject(&mockObject, mockResource))
	assert.Nil(t, mockDecoder.DecodeObject(&mockObject, mockResource))
}
