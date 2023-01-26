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

func TestMockENotifierEAdapters(t *testing.T) {
	n := NewMockENotifier(t)
	a := NewMockEList(t)
	m := NewMockRun(t)
	n.EXPECT().EAdapters().Run(func() { m.Run() }).Return(a).Once()
	n.EXPECT().EAdapters().Once().Return(func() EList { return a })
	assert.Equal(t, a, n.EAdapters())
	assert.Equal(t, a, n.EAdapters())
}

func TestMockENotifierEDeliver(t *testing.T) {
	n := NewMockENotifier(t)
	m := NewMockRun(t)
	n.EXPECT().EDeliver().Run(func() { m.Run() }).Return(false).Once()
	n.EXPECT().EDeliver().Once().Return(func() bool { return true }).Once()
	assert.False(t, n.EDeliver())
	assert.True(t, n.EDeliver())
}

func TestMockENotifierESetDeliver(t *testing.T) {
	n := NewMockENotifier(t)
	m := NewMockRun(t, true)
	n.EXPECT().ESetDeliver(true).Return().Run(func(_a0 bool) { m.Run(_a0) }).Once()
	n.ESetDeliver(true)
}

func TestMockENotifierENotify(t *testing.T) {
	n := NewMockENotifier(t)
	notif := NewMockENotification(t)
	m := NewMockRun(t, notif)
	n.EXPECT().ENotify(notif).Return().Run(func(_a0 ENotification) { m.Run(_a0) }).Once()
	n.ENotify(notif)
}
