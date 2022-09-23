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

func TestMockEMapEntry_GetKey(t *testing.T) {
	l := &MockEMapEntry{}
	l.On("GetKey").Once().Return("1")
	l.On("GetKey").Once().Return(func() any {
		return "2"
	})
	assert.Equal(t, "1", l.GetKey())
	assert.Equal(t, "2", l.GetKey())
	mock.AssertExpectationsForObjects(t, l)
}

func TestMockEMapEntry_SetKey(t *testing.T) {
	l := &MockEMapEntry{}
	l.On("SetKey", 1)
	l.SetKey(1)
	mock.AssertExpectationsForObjects(t, l)
}

func TestMockEMapEntry_GetValue(t *testing.T) {
	l := &MockEMapEntry{}
	l.On("GetValue").Once().Return("1")
	l.On("GetValue").Once().Return(func() any {
		return "2"
	})
	assert.Equal(t, "1", l.GetValue())
	assert.Equal(t, "2", l.GetValue())
	mock.AssertExpectationsForObjects(t, l)
}

func TestMockEMapEntry_SetValue(t *testing.T) {
	l := &MockEMapEntry{}
	l.On("SetValue", 1)
	l.SetValue(1)
	mock.AssertExpectationsForObjects(t, l)
}
