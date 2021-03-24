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

func TestMockEMap_ContainsKey(t *testing.T) {
	l := &MockEMap{}
	l.On("ContainsKey", 1).Once().Return(true)
	l.On("ContainsKey", 2).Once().Return(func(interface{}) bool {
		return true
	})
	assert.True(t, l.ContainsKey(1))
	assert.True(t, l.ContainsKey(2))
	mock.AssertExpectationsForObjects(t, l)
}

func TestMockEMap_ContainsValue(t *testing.T) {
	l := &MockEMap{}
	l.On("ContainsValue", 1).Once().Return(true)
	l.On("ContainsValue", 2).Once().Return(func(interface{}) bool {
		return true
	})
	assert.True(t, l.ContainsValue(1))
	assert.True(t, l.ContainsValue(2))
	mock.AssertExpectationsForObjects(t, l)
}

func TestMockEMap_RemoveKey(t *testing.T) {
	l := &MockEMap{}
	l.On("RemoveKey", 1).Once().Return("1")
	l.On("RemoveKey", 2).Once().Return(func(interface{}) interface{} {
		return "2"
	})
	assert.Equal(t, "1", l.RemoveKey(1))
	assert.Equal(t, "2", l.RemoveKey(2))
	mock.AssertExpectationsForObjects(t, l)
}

func TestMockEMap_Put(t *testing.T) {
	l := &MockEMap{}
	l.On("Put", 1, "1")
	l.Put(1, "1")
	mock.AssertExpectationsForObjects(t, l)
}

func TestMockEMap_GetValue(t *testing.T) {
	l := &MockEMap{}
	l.On("GetValue", 1).Once().Return("1")
	l.On("GetValue", 2).Once().Return(func(interface{}) interface{} {
		return "2"
	})
	assert.Equal(t, "1", l.GetValue(1))
	assert.Equal(t, "2", l.GetValue(2))
	mock.AssertExpectationsForObjects(t, l)
}

func TestMockEMap_ToMap(t *testing.T) {
	l := &MockEMap{}
	m := map[interface{}]interface{}{}
	l.On("ToMap").Once().Return(m)
	l.On("ToMap").Once().Return(func() map[interface{}]interface{} {
		return m
	})
	assert.Equal(t, m, l.ToMap())
	assert.Equal(t, m, l.ToMap())
	mock.AssertExpectationsForObjects(t, l)
}
