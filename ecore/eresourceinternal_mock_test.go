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
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMockEResourceInternalDoLoad(t *testing.T) {
	r := &MockEResourceInternal{}
	uri, _ := url.Parse("test://file.t")
	f, _ := os.Open(uri.String())
	m := make(map[string]interface{})
	r.On("DoLoad", f, m).Once()
	r.DoLoad(f, m)
	mock.AssertExpectationsForObjects(t, r)
}

func TestMockEResourceInternalDoSave(t *testing.T) {
	r := &MockEResourceInternal{}
	uri, _ := url.Parse("test://file.t")
	f, _ := os.Create(uri.String())
	m := make(map[string]interface{})
	r.On("DoSave", f, m).Once()
	r.DoSave(f, m)
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
	r.On("basicSetLoaded", false, n1).Return(n2).Once()
	r.On("basicSetLoaded", false, n1).Return(func(bool, ENotificationChain) ENotificationChain {
		return n2
	}).Once()
	assert.Equal(t, n2, r.basicSetLoaded(false, n1))
	assert.Equal(t, n2, r.basicSetLoaded(false, n1))
	mock.AssertExpectationsForObjects(t, r, n1, n2)
}

func TestMockEResourceInternalBasicSetResourceSet(t *testing.T) {
	r := &MockEResourceInternal{}
	rs := &MockEResourceSet{}
	n1 := &MockENotificationChain{}
	n2 := &MockENotificationChain{}
	r.On("basicSetResourceSet", rs, n1).Return(n2).Once()
	r.On("basicSetResourceSet", rs, n1).Return(func(EResourceSet, ENotificationChain) ENotificationChain {
		return n2
	}).Once()
	assert.Equal(t, n2, r.basicSetResourceSet(rs, n1))
	assert.Equal(t, n2, r.basicSetResourceSet(rs, n1))
	mock.AssertExpectationsForObjects(t, r, rs, n1, n2)
}
