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

func TestMockEObjectProperties_EDynamicGet(t *testing.T) {
	o := NewMockEObjectProperties(t)
	obj := NewMockEObject(t)
	m := NewMockRun(t, 1)
	o.EXPECT().EDynamicGet(1).Return(obj).Run(func(dynamicFeatureID int) { m.Run(dynamicFeatureID) }).Once()
	o.EXPECT().EDynamicGet(1).Call.Return(func(dynamicFeatureID int) any {
		return obj
	}).Once()
	assert.Equal(t, obj, o.EDynamicGet(1))
	assert.Equal(t, obj, o.EDynamicGet(1))
}

func TestMockEObjectProperties_EDynamicSet(t *testing.T) {
	o := NewMockEObjectProperties(t)
	obj := NewMockEObject(t)
	m := NewMockRun(t, 1, obj)
	o.EXPECT().EDynamicSet(1, obj).Return().Run(func(dynamicFeatureID int, newValue interface{}) { m.Run(dynamicFeatureID, newValue) }).Once()
	o.EDynamicSet(1, obj)
}

func TestMockEObjectProperties_EDynamicUnset(t *testing.T) {
	o := NewMockEObjectProperties(t)
	m := NewMockRun(t, 1)
	o.EXPECT().EDynamicUnset(1).Return().Run(func(dynamicFeatureID int) { m.Run(dynamicFeatureID) }).Once()
	o.EDynamicUnset(1)
	mock.AssertExpectationsForObjects(t, o)
}
