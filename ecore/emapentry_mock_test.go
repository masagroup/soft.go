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

type mockEMapEntryRun struct {
	mock.Mock
}

func (m *mockEMapEntryRun) Run(args ...any) {
	m.Called(args...)
}

type mockConstructorTestingTMockMockEMapEntryRun interface {
	mock.TestingT
	Cleanup(func())
}

// newMockEListRun creates a new instance of MockEList. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockEMapEntryRun(t mockConstructorTestingTMockMockEMapEntryRun, args ...any) *mockEMapEntryRun {
	mock := &mockEMapEntryRun{}
	mock.Test(t)
	mock.On("Run", args...).Once()
	t.Cleanup(func() { mock.AssertExpectations(t) })
	return mock
}

func TestMockEMapEntry_GetKey(t *testing.T) {
	l := NewMockEMapEntry(t)
	m := newMockEMapEntryRun(t)
	l.EXPECT().GetKey().Return("1").Run(func() { m.Run() }).Once()
	l.EXPECT().GetKey().Call.Return(func() any {
		return "2"
	}).Once()
	assert.Equal(t, "1", l.GetKey())
	assert.Equal(t, "2", l.GetKey())
}

func TestMockEMapEntry_SetKey(t *testing.T) {
	l := NewMockEMapEntry(t)
	m := newMockEMapEntryRun(t, 1)
	l.EXPECT().SetKey(1).Return().Run(func(_a0 interface{}) { m.Run(_a0) }).Once()
	l.SetKey(1)
}

func TestMockEMapEntry_GetValue(t *testing.T) {
	l := NewMockEMapEntry(t)
	m := newMockEMapEntryRun(t)
	l.EXPECT().GetValue().Return("1").Run(func() { m.Run() }).Once()
	l.EXPECT().GetValue().Call.Return(func() any {
		return "2"
	}).Once()
	assert.Equal(t, "1", l.GetValue())
	assert.Equal(t, "2", l.GetValue())
}

func TestMockEMapEntry_SetValue(t *testing.T) {
	l := NewMockEMapEntry(t)
	m := newMockEMapEntryRun(t, 1)
	l.EXPECT().SetValue(1).Return().Run(func(_a0 interface{}) { m.Run(_a0) }).Once()
	l.SetValue(1)
}
