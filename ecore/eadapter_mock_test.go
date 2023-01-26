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

func TestMockEAdapter_GetTarget(t *testing.T) {
	mockNotifier := NewMockENotifier(t)
	mockAdapter := NewMockEAdapter(t)
	m := NewMockRun(t)
	mockAdapter.EXPECT().GetTarget().Return(mockNotifier).Run(func() { m.Run() }).Once()
	mockAdapter.EXPECT().GetTarget().Once().Return(func() ENotifier { return mockNotifier })
	assert.Equal(t, mockNotifier, mockAdapter.GetTarget())
	assert.Equal(t, mockNotifier, mockAdapter.GetTarget())
}

func TestMockEAdapter_NotifyChanged(t *testing.T) {
	mockNotification := NewMockENotification(t)
	mockAdapter := NewMockEAdapter(t)
	m := NewMockRun(t, mockNotification)
	mockAdapter.EXPECT().NotifyChanged(mockNotification).Return().Run(func(notification ENotification) { m.Run(notification) }).Once()
	mockAdapter.NotifyChanged(mockNotification)
}

func TestMockEAdapter_SetTarget(t *testing.T) {
	mockNotifier := NewMockENotifier(t)
	mockAdapter := NewMockEAdapter(t)
	m := NewMockRun(t, mockNotifier)
	mockAdapter.EXPECT().SetTarget(mockNotifier).Return().Run(func(_a0 ENotifier) { m.Run(_a0) }).Once()
	mockAdapter.SetTarget(mockNotifier)
}

func TestMockEAdapter_UnSetTarget(t *testing.T) {
	mockNotifier := NewMockENotifier(t)
	mockAdapter := NewMockEAdapter(t)
	m := NewMockRun(t, mockNotifier)
	mockAdapter.EXPECT().UnSetTarget(mockNotifier).Return().Run(func(_a0 ENotifier) { m.Run(_a0) }).Once()
	mockAdapter.UnSetTarget(mockNotifier)
}
