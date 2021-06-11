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
)

func TestMockEResourceDecoder_Decode(t *testing.T) {
	mockDecoder := &MockEResourceDecoder{}
	mockResource := &MockEResource{}
	var mockReader io.Reader = nil
	mockDecoder.On("Decode", mockResource, mockReader).Once()
	mockDecoder.Decode(mockResource, mockReader)
}
