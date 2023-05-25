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

func TestMockEEncoder_Encode(t *testing.T) {
	m := NewMockRun(t)
	mockEncoder := NewMockEEncoder(t)
	mockEncoder.EXPECT().Encode().Return().Run(func() { m.Run() }).Once()
	mockEncoder.Encode()
}

func TestMockEEncoder_EncodeObject(t *testing.T) {
	mockEncoder := NewMockEEncoder(t)
	mockObject := NewMockEObject(t)
	m := NewMockRun(t, mockObject)
	mockEncoder.EXPECT().EncodeObject(mockObject).Return(nil).Run(func(object EObject) { m.Run(object) }).Once()
	mockEncoder.EXPECT().EncodeObject(mockObject).Call.Return(func(EObject) error { return nil }).Once()
	assert.Nil(t, mockEncoder.EncodeObject(mockObject))
	assert.Nil(t, mockEncoder.EncodeObject(mockObject))
}
