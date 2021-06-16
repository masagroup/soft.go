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

func TestMockEResourceEncoder_Encode(t *testing.T) {
	mockEncoder := &MockEResourceEncoder{}
	mockEncoder.On("Encode").Once()
	mockEncoder.Encode()
}

func TestMockEResourceEncoder_EncodeObject(t *testing.T) {
	mockEncoder := &MockEResourceEncoder{}
	mockObject := &MockEObject{}
	mockEncoder.On("EncodeObject", mockObject).Return(nil).Once()
	assert.Nil(t, mockEncoder.EncodeObject(mockObject))
}
