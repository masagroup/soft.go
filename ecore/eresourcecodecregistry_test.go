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
	"github.com/stretchr/testify/require"
)

func TestResoureCodecRegistrySingleton(t *testing.T) {
	r := GetResourceCodecRegistry()
	require.NotNil(t, r)
	assert.NotNil(t, r.GetExtensionToCodecMap()["ecore"])
	assert.NotNil(t, r.GetExtensionToCodecMap()["xml"])
}

func TestResoureCodecRegistrySingletonGetCodec(t *testing.T) {
	r := GetResourceCodecRegistry()
	assert.NotNil(t, r.GetCodec(NewURI("*.xml")))
	assert.NotNil(t, r.GetCodec(NewURI("*.ecore")))
}
