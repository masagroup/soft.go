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

func TestMockENotifierImpl_Adapters(t *testing.T) {
	n := &ENotifierImpl{}
	n.SetInterfaces(n)
	assert.Nil(t, n.EBasicAdapters())
	assert.False(t, n.EBasicHasAdapters())

	assert.NotNil(t, n.EAdapters())
	assert.NotNil(t, n.EBasicAdapters())
	assert.False(t, n.EBasicHasAdapters())
}

func TestMockENotifierImpl_Deliver(t *testing.T) {
	n := &ENotifierImpl{}

	assert.False(t, n.EDeliver())

	n.ESetDeliver(true)
	assert.True(t, n.EDeliver())
}
