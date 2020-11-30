// *****************************************************************************
//
// This file is part of a MASA library or program.
// Refer to the included end-user license agreement for restrictions.
//
// Copyright (c) 2020 MASA Group
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
}
