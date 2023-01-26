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

func TestMockENotifyingListGetNotifier(t *testing.T) {
	l := NewMockENotifyingList(t)
	n := NewMockENotifier(t)
	l.On("GetNotifier").Return(n).Once()
	l.On("GetNotifier").Return(func() ENotifier {
		return n
	}).Once()
	assert.Equal(t, n, l.GetNotifier())
	assert.Equal(t, n, l.GetNotifier())
	mock.AssertExpectationsForObjects(t, l, n)
}

func TestMockEPackageRegistryGetFeature(t *testing.T) {
	l := NewMockENotifyingList(t)
	f := NewMockEStructuralFeature(t)
	l.On("GetFeature").Return(f).Once()
	l.On("GetFeature").Return(func() EStructuralFeature {
		return f
	}).Once()
	assert.Equal(t, f, l.GetFeature())
	assert.Equal(t, f, l.GetFeature())
	mock.AssertExpectationsForObjects(t, l, f)
}

func TestMockEPackageRegistryGetFeatureID(t *testing.T) {
	l := NewMockENotifyingList(t)
	l.On("GetFeatureID").Return(1).Once()
	l.On("GetFeatureID").Return(func() int {
		return 2
	}).Once()
	assert.Equal(t, 1, l.GetFeatureID())
	assert.Equal(t, 2, l.GetFeatureID())
	mock.AssertExpectationsForObjects(t, l)
}

func TestMockEPackageRegistryAddWithNotification(t *testing.T) {
	l := NewMockENotifyingList(t)
	n := NewMockENotificationChain(t)
	v := NewMockEObject(t)
	l.On("AddWithNotification", v, n).Return(n).Once()
	l.On("AddWithNotification", v, n).Return(func(object any, notifications ENotificationChain) ENotificationChain {
		return n
	}).Once()
	assert.Equal(t, n, l.AddWithNotification(v, n))
	assert.Equal(t, n, l.AddWithNotification(v, n))
	mock.AssertExpectationsForObjects(t, l, n, v)
}

func TestMockEPackageRegistryRemoveWithNotification(t *testing.T) {
	l := NewMockENotifyingList(t)
	n := NewMockENotificationChain(t)
	v := NewMockEObject(t)
	l.On("RemoveWithNotification", v, n).Return(n).Once()
	l.On("RemoveWithNotification", v, n).Return(func(object any, notifications ENotificationChain) ENotificationChain {
		return n
	}).Once()
	assert.Equal(t, n, l.RemoveWithNotification(v, n))
	assert.Equal(t, n, l.RemoveWithNotification(v, n))
	mock.AssertExpectationsForObjects(t, l, n, v)
}

func TestMockEPackageRegistrySetWithNotification(t *testing.T) {
	l := NewMockENotifyingList(t)
	n := NewMockENotificationChain(t)
	v := NewMockEObject(t)
	l.On("SetWithNotification", 0, v, n).Return(n).Once()
	l.On("SetWithNotification", 1, v, n).Return(func(index int, object any, notifications ENotificationChain) ENotificationChain {
		return n
	}).Once()
	assert.Equal(t, n, l.SetWithNotification(0, v, n))
	assert.Equal(t, n, l.SetWithNotification(1, v, n))
	mock.AssertExpectationsForObjects(t, l, n, v)
}
