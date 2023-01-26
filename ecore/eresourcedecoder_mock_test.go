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
	mockDecoder.On("Decode").Once()
	mockDecoder.Decode()
}

func TestMockEResourceDecoder_DecodeObject(t *testing.T) {
	mockDecoder := &MockEResourceDecoder{}
	mockObject := NewMockEObject(t)
	mockDecoder.On("DecodeObject").Return(mockObject, nil).Once()
	mockDecoder.On("DecodeObject").Return(func() (EObject, error) {
		return mockObject, nil
	}).Once()
	{
		obj, err := mockDecoder.DecodeObject()
		assert.Equal(t, mockObject, obj)
		assert.Nil(t, err)
	}
	{
		obj, err := mockDecoder.DecodeObject()
		assert.Equal(t, mockObject, obj)
		assert.Nil(t, err)
	}
}
