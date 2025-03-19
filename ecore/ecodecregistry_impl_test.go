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

func TestECodecRegistryGetCodecProtocol(t *testing.T) {
	mockCodec := new(MockECodec)

	rfr := NewECodecRegistryImpl()
	rfr.GetProtocolToCodecMap()["test"] = mockCodec

	assert.Equal(t, mockCodec, rfr.GetCodec(NewURI("test:///file.t")))
	assert.Nil(t, rfr.GetCodec(NewURI("file:///file.t")))
}

func TestECodecRegistryGetCodecExtension(t *testing.T) {
	mockCodec := new(MockECodec)

	rfr := NewECodecRegistryImpl()
	rfr.GetExtensionToCodecMap()["test"] = mockCodec

	assert.Equal(t, mockCodec, rfr.GetCodec(NewURI("test:///file.test")))
	assert.Nil(t, rfr.GetCodec(NewURI("file:///file.t")))
}

func TestECodecRegistryGetCodecDefault(t *testing.T) {
	mockCodec := new(MockECodec)

	rfr := NewECodecRegistryImpl()
	rfr.GetExtensionToCodecMap()[DEFAULT_EXTENSION] = mockCodec

	assert.Equal(t, mockCodec, rfr.GetCodec(NewURI("test:///file.test")))
	assert.Equal(t, mockCodec, rfr.GetCodec(NewURI("file:///file.t")))
}
