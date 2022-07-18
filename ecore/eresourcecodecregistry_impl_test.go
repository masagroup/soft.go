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

func TestEResourceCodecRegistryGetCodecProtocol(t *testing.T) {
	mockCodec := new(MockEResourceCodec)

	rfr := NewEResourceCodecRegistryImpl()
	rfr.GetProtocolToCodecMap()["test"] = mockCodec

	assert.Equal(t, mockCodec, rfr.GetCodec(&URI{scheme: "test", Path: "//file.t"}))
	assert.Nil(t, rfr.GetCodec(&URI{scheme: "file", Path: "//file.t"}))
}

func TestEResourceCodecRegistryGetCodecExtension(t *testing.T) {
	mockCodec := new(MockEResourceCodec)

	rfr := NewEResourceCodecRegistryImpl()
	rfr.GetExtensionToCodecMap()["test"] = mockCodec

	assert.Equal(t, mockCodec, rfr.GetCodec(&URI{scheme: "file", Path: "//file.test"}))
	assert.Nil(t, rfr.GetCodec(&URI{scheme: "file", Path: "//file.t"}))
}
