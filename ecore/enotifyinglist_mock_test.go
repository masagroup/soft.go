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

func TestMockENotifyingListGetNotifier(t *testing.T) {
	l := NewMockENotifyingList(t)
	n := NewMockENotifier(t)
	m := NewMockRun(t)
	l.EXPECT().GetNotifier().Return(n).Run(func() { m.Run() }).Once()
	l.EXPECT().GetNotifier().Call.Return(func() ENotifier {
		return n
	}).Once()
	assert.Equal(t, n, l.GetNotifier())
	assert.Equal(t, n, l.GetNotifier())
}

func TestMockEPackageRegistryGetFeature(t *testing.T) {
	l := NewMockENotifyingList(t)
	f := NewMockEStructuralFeature(t)
	m := NewMockRun(t)
	l.EXPECT().GetFeature().Return(f).Run(func() { m.Run() }).Once()
	l.EXPECT().GetFeature().Call.Return(func() EStructuralFeature {
		return f
	}).Once()
	assert.Equal(t, f, l.GetFeature())
	assert.Equal(t, f, l.GetFeature())
}

func TestMockEPackageRegistryGetFeatureID(t *testing.T) {
	l := NewMockENotifyingList(t)
	m := NewMockRun(t)
	l.EXPECT().GetFeatureID().Return(1).Run(func() { m.Run() }).Once()
	l.EXPECT().GetFeatureID().Call.Return(func() int {
		return 2
	}).Once()
	assert.Equal(t, 1, l.GetFeatureID())
	assert.Equal(t, 2, l.GetFeatureID())
}

func TestMockEPackageRegistryAddWithNotification(t *testing.T) {
	l := NewMockENotifyingList(t)
	n := NewMockENotificationChain(t)
	v := NewMockEObject(t)
	m := NewMockRun(t, v, n)
	l.EXPECT().AddWithNotification(v, n).Return(n).Run(func(object interface{}, notifications ENotificationChain) { m.Run(object, notifications) }).Once()
	l.EXPECT().AddWithNotification(v, n).Call.Return(func(object any, notifications ENotificationChain) ENotificationChain {
		return n
	}).Once()
	assert.Equal(t, n, l.AddWithNotification(v, n))
	assert.Equal(t, n, l.AddWithNotification(v, n))
}

func TestMockEPackageRegistryRemoveWithNotification(t *testing.T) {
	l := NewMockENotifyingList(t)
	n := NewMockENotificationChain(t)
	v := NewMockEObject(t)
	m := NewMockRun(t, v, n)
	l.EXPECT().RemoveWithNotification(v, n).Return(n).Run(func(object interface{}, notifications ENotificationChain) { m.Run(object, notifications) }).Once()
	l.EXPECT().RemoveWithNotification(v, n).Call.Return(func(object any, notifications ENotificationChain) ENotificationChain {
		return n
	}).Once()
	assert.Equal(t, n, l.RemoveWithNotification(v, n))
	assert.Equal(t, n, l.RemoveWithNotification(v, n))
}

func TestMockEPackageRegistrySetWithNotification(t *testing.T) {
	l := NewMockENotifyingList(t)
	n := NewMockENotificationChain(t)
	v := NewMockEObject(t)
	m := NewMockRun(t, 0, v, n)
	l.EXPECT().SetWithNotification(0, v, n).Return(n).Run(func(index int, object interface{}, notifications ENotificationChain) {
		m.Run(index, object, notifications)
	}).Once()
	l.EXPECT().SetWithNotification(1, v, n).Call.Return(func(index int, object any, notifications ENotificationChain) ENotificationChain {
		return n
	}).Once()
	assert.Equal(t, n, l.SetWithNotification(0, v, n))
	assert.Equal(t, n, l.SetWithNotification(1, v, n))
}
