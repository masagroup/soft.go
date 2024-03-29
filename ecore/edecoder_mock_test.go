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

func TestMockEDecoder_DecodeResource(t *testing.T) {
	mockDecoder := NewMockEDecoder(t)
	m := NewMockRun(t)
	mockDecoder.EXPECT().DecodeResource().Return().Run(func() { m.Run() }).Once()
	mockDecoder.DecodeResource()
}

func TestMockEDecoder_DecodeObject(t *testing.T) {
	mockDecoder := NewMockEDecoder(t)
	mockObject := NewMockEObject(t)
	m := NewMockRun(t)
	mockDecoder.EXPECT().DecodeObject().Return(mockObject, nil).Run(func() { m.Run() }).Once()
	mockDecoder.EXPECT().DecodeObject().Call.Return(func() (EObject, error) {
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
