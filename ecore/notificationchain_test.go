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

func TestNotificationChainAdd(t *testing.T) {
	chain := NewNotificationChain()
	assert.False(t, chain.Add(nil))

	mockNotif1 := &MockENotification{}
	assert.True(t, chain.Add(mockNotif1))

	mockNotif2 := &MockENotification{}
	mockNotif1.On("Merge", mockNotif2).Return(true).Once()
	assert.False(t, chain.Add(mockNotif2))

	mockNotif3 := &MockENotification{}
	mockNotif1.On("Merge", mockNotif3).Return(false).Once()
	assert.True(t, chain.Add(mockNotif3))

	mockNotif1.AssertExpectations(t)
	mockNotif2.AssertExpectations(t)
	mockNotif3.AssertExpectations(t)
}
