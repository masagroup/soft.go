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

	"github.com/stretchr/testify/require"
)

func TestMockEObjectProperties_EDynamicGet(t *testing.T) {
	o := NewMockEObjectProperties(t)
	obj := NewMockEObject(t)
	m := NewMockRun(t, 1)
	o.EXPECT().EDynamicGet(1).Return(obj).Run(func(dynamicFeatureID int) { m.Run(dynamicFeatureID) }).Once()
	o.EXPECT().EDynamicGet(1).RunAndReturn(func(dynamicFeatureID int) any {
		return obj
	}).Once()
	o.EXPECT().EDynamicGet(1)
	require.Equal(t, obj, o.EDynamicGet(1))
	require.Equal(t, obj, o.EDynamicGet(1))
	require.Panics(t, func() {
		o.EDynamicGet(1)
	})
}

func TestMockEObjectProperties_EDynamicSet(t *testing.T) {
	o := NewMockEObjectProperties(t)
	obj := NewMockEObject(t)
	m := NewMockRun(t, 1, obj)
	o.EXPECT().EDynamicSet(1, obj).Return().Run(func(dynamicFeatureID int, newValue any) { m.Run(dynamicFeatureID, newValue) }).Once()
	o.EXPECT().EDynamicSet(1, obj).RunAndReturn(func(i1 int, i2 interface{}) {}).Once()
	o.EDynamicSet(1, obj)
	o.EDynamicSet(1, obj)
}

func TestMockEObjectProperties_EDynamicUnset(t *testing.T) {
	o := NewMockEObjectProperties(t)
	m := NewMockRun(t, 1)
	o.EXPECT().EDynamicUnset(1).Return().Run(func(dynamicFeatureID int) { m.Run(dynamicFeatureID) }).Once()
	o.EXPECT().EDynamicUnset(1).RunAndReturn(func(i int) {}).Once()
	o.EDynamicUnset(1)
	o.EDynamicUnset(1)
}

func TestMockEObjectProperties_EDynamicIsSet(t *testing.T) {
	o := NewMockEObjectProperties(t)
	m := NewMockRun(t, 1)
	o.EXPECT().EDynamicIsSet(1).Return(true).Run(func(dynamicFeatureID int) { m.Run(dynamicFeatureID) }).Once()
	o.EXPECT().EDynamicIsSet(1).RunAndReturn(func(dynamicFeatureID int) bool {
		return true
	}).Once()
	o.EXPECT().EDynamicIsSet(1)
	require.True(t, o.EDynamicIsSet(1))
	require.True(t, o.EDynamicIsSet(1))
	require.Panics(t, func() {
		o.EDynamicIsSet(1)
	})
}
