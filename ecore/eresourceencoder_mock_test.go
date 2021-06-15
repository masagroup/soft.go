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

func TestMockEResourceEncoder_EncodeResource(t *testing.T) {
	mockEncoder := &MockEResourceEncoder{}
	mockResource := &MockEResource{}
	mockEncoder.On("EncodeResource", mockResource).Once()
	mockEncoder.EncodeResource(mockResource)
}

func TestMockEResourceEncoder_EncodeObject(t *testing.T) {
	mockEncoder := &MockEResourceEncoder{}
	mockObject := &MockEObject{}
	mockResource := &MockEResource{}
	mockEncoder.On("EncodeObject", mockObject, mockResource).Return(nil).Once()
	assert.Nil(t, mockEncoder.EncodeObject(mockObject, mockResource))
}
