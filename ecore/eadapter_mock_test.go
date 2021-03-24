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
	"github.com/stretchr/testify/mock"
)

func TestMockEAdapter_GetTarget(t *testing.T) {
	mockNotifier := new(MockENotifier)
	mockAdapter := new(MockEAdapter)
	mockAdapter.On("GetTarget").Return(mockNotifier).Once()
	assert.Equal(t, mockNotifier, mockAdapter.GetTarget())
	mock.AssertExpectationsForObjects(t, mockNotifier, mockAdapter)

	mockAdapter.On("GetTarget").Return(func() ENotifier { return mockNotifier }).Once()
	assert.Equal(t, mockNotifier, mockAdapter.GetTarget())
	mock.AssertExpectationsForObjects(t, mockNotifier, mockAdapter)
}

func TestMockEAdapter_UnSetTarget(t *testing.T) {
	mockNotifier := new(MockENotifier)
	mockAdapter := new(MockEAdapter)
	mockAdapter.On("UnSetTarget", mockNotifier).Once()
	mockAdapter.UnSetTarget(mockNotifier)
	mock.AssertExpectationsForObjects(t, mockNotifier, mockAdapter)
}
