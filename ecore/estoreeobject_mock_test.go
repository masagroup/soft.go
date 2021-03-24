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

func TestMockEStoreEObjectEStore(t *testing.T) {
	o := &MockEStoreEObject{}
	mockStore := &MockEStore{}
	o.On("EStore").Return(mockStore).Once()
	o.On("EStore").Return(func() EStore {
		return mockStore
	}).Once()
	assert.Equal(t, mockStore, o.EStore())
	assert.Equal(t, mockStore, o.EStore())
	o.AssertExpectations(t)
}
