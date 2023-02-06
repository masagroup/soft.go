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

func TestMockENotificationChainAdd(t *testing.T) {
	nc := NewMockENotificationChain(t)
	n := NewMockENotification(t)
	m := NewMockRun(t, n)
	nc.EXPECT().Add(n).Return(true).Run(func(n ENotification) { m.Run(n) }).Once()
	nc.EXPECT().Add(n).Call.Return(func(ENotification) bool { return false }).Once()
	assert.True(t, nc.Add(n))
	assert.False(t, nc.Add(n))
}

func TestMockENotificationChainDispatch(t *testing.T) {
	nc := NewMockENotificationChain(t)
	m := NewMockRun(t)
	nc.EXPECT().Dispatch().Return().Run(func() { m.Run() }).Once()
	nc.Dispatch()
}
