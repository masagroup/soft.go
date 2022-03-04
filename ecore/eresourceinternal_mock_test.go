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

func TestMockEResourceInternalDoLoad(t *testing.T) {
	r := &MockEResourceInternal{}
	mockDecoder := &MockEResourceDecoder{}
	r.On("DoLoad", mockDecoder).Once()
	r.DoLoad(mockDecoder)
	mock.AssertExpectationsForObjects(t, r)
}

func TestMockEResourceInternalDoSave(t *testing.T) {
	r := &MockEResourceInternal{}
	mockEncoder := &MockEResourceEncoder{}
	r.On("DoSave", mockEncoder).Once()
	r.DoSave(mockEncoder)
	mock.AssertExpectationsForObjects(t, r)
}

func TestMockEResourceInternalDoUnLoad(t *testing.T) {
	r := &MockEResourceInternal{}
	r.On("DoUnload").Once()
	r.DoUnload()
	mock.AssertExpectationsForObjects(t, r)
}

func TestMockEResourceInternalBasicSetLoaded(t *testing.T) {
	r := &MockEResourceInternal{}
	n1 := &MockENotificationChain{}
	n2 := &MockENotificationChain{}
	r.On("BasicSetLoaded", false, n1).Return(n2).Once()
	r.On("BasicSetLoaded", false, n1).Return(func(bool, ENotificationChain) ENotificationChain {
		return n2
	}).Once()
	assert.Equal(t, n2, r.BasicSetLoaded(false, n1))
	assert.Equal(t, n2, r.BasicSetLoaded(false, n1))
	mock.AssertExpectationsForObjects(t, r, n1, n2)
}

func TestMockEResourceInternalBasicSetResourceSet(t *testing.T) {
	r := &MockEResourceInternal{}
	rs := &MockEResourceSet{}
	n1 := &MockENotificationChain{}
	n2 := &MockENotificationChain{}
	r.On("BasicSetResourceSet", rs, n1).Return(n2).Once()
	r.On("BasicSetResourceSet", rs, n1).Return(func(EResourceSet, ENotificationChain) ENotificationChain {
		return n2
	}).Once()
	assert.Equal(t, n2, r.BasicSetResourceSet(rs, n1))
	assert.Equal(t, n2, r.BasicSetResourceSet(rs, n1))
	mock.AssertExpectationsForObjects(t, r, rs, n1, n2)
}

func TestMockEResourceInternalDoAttached(t *testing.T) {
	r := &MockEResourceInternal{}
	o := &MockEObject{}
	r.On("DoAttached", o)
	r.DoAttached(o)
	r.AssertExpectations(t)
}

func TestMockEResourceInternalDetached(t *testing.T) {
	r := &MockEResourceInternal{}
	o := &MockEObject{}
	r.On("DoDetached", o)
	r.DoDetached(o)
	r.AssertExpectations(t)
}

func TestMockEResourceInternalIsAttachedDetachedRequired(t *testing.T) {
	r := &MockEResourceInternal{}
	r.On("IsAttachedDetachedRequired").Return(true).Once()
	r.On("IsAttachedDetachedRequired").Return(func() bool {
		return false
	}).Once()
	assert.True(t, r.IsAttachedDetachedRequired())
	assert.False(t, r.IsAttachedDetachedRequired())
	mock.AssertExpectationsForObjects(t, r)
}
