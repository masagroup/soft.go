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

// TestGetELiterals tests method GetELiterals
func TestMockEIteratorHasNext(t *testing.T) {
	o := &MockEIterator{}
	v := true
	o.On("HasNext").Once().Return(v)
	assert.Equal(t, v, o.HasNext())

	o.On("HasNext").Once().Return(func() bool {
		return v
	})
	assert.Equal(t, v, o.HasNext())
	o.AssertExpectations(t)
}

func TestMockEIteratorNext(t *testing.T) {
	o := &MockEIterator{}
	v := &mock.Mock{}
	o.On("Next").Once().Return(v)
	assert.Equal(t, v, o.Next())

	o.On("Next").Once().Return(func() interface{} {
		return v
	})
	assert.Equal(t, v, o.Next())
	o.AssertExpectations(t)
}
