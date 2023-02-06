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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockEObjectIDManagerClear(t *testing.T) {
	rm := NewMockEObjectIDManager(t)
	m := NewMockRun(t)
	rm.EXPECT().Clear().Return().Run(func() { m.Run() }).Once()
	rm.Clear()
}

func TestMockEObjectIDManagerRegister(t *testing.T) {
	rm := NewMockEObjectIDManager(t)
	o := NewMockEObject(t)
	m := NewMockRun(t, o)
	rm.EXPECT().Register(o).Return().Run(func(_a0 EObject) { m.Run(_a0) }).Once()
	rm.Register(o)
}

func TestMockEObjectIDManagerUnRegister(t *testing.T) {
	rm := NewMockEObjectIDManager(t)
	o := NewMockEObject(t)
	m := NewMockRun(t, o)
	rm.EXPECT().UnRegister(o).Return().Run(func(_a0 EObject) { m.Run(_a0) }).Once()
	rm.UnRegister(o)
}

func TestMockEObjectIDManagerGetID(t *testing.T) {
	rm := NewMockEObjectIDManager(t)
	o := NewMockEObject(t)
	m := NewMockRun(t, o)
	rm.EXPECT().GetID(o).Return("id1").Run(func(_a0 EObject) { m.Run(_a0) }).Once()
	rm.EXPECT().GetID(o).Call.Return(func(EObject) any {
		return "id2"
	}).Once()
	assert.Equal(t, "id1", rm.GetID(o))
	assert.Equal(t, "id2", rm.GetID(o))
}

func TestMockEObjectIDManagerSetID(t *testing.T) {
	rm := NewMockEObjectIDManager(t)
	o := NewMockEObject(t)
	m := NewMockRun(t, o, "id1")
	rm.EXPECT().SetID(o, "id1").Return(errors.New("error")).Run(func(_a0 EObject, _a1 interface{}) { m.Run(_a0, _a1) }).Once()
	rm.EXPECT().SetID(o, "id2").Call.Return(func(EObject, interface{}) error { return errors.New("error") }).Once()
	assert.NotNil(t, rm.SetID(o, "id1"))
	assert.NotNil(t, rm.SetID(o, "id2"))
}

func TestMockEObjectIDManagerGetDetachedID(t *testing.T) {
	rm := NewMockEObjectIDManager(t)
	o := NewMockEObject(t)
	m := NewMockRun(t, o)
	rm.EXPECT().GetDetachedID(o).Return("id1").Run(func(_a0 EObject) { m.Run(_a0) }).Once()
	rm.EXPECT().GetDetachedID(o).Call.Return(func(EObject) any {
		return "id2"
	}).Once()
	assert.Equal(t, "id1", rm.GetDetachedID(o))
	assert.Equal(t, "id2", rm.GetDetachedID(o))
}

func TestMockEObjectIDManagerGetEObject(t *testing.T) {
	rm := NewMockEObjectIDManager(t)
	o := NewMockEObject(t)
	m := NewMockRun(t, "id1")
	rm.EXPECT().GetEObject("id1").Return(o).Run(func(_a0 interface{}) { m.Run(_a0) }).Once()
	rm.EXPECT().GetEObject("id2").Call.Return(func(interface{}) EObject {
		return o
	}).Once()
	assert.Equal(t, o, rm.GetEObject("id1"))
	assert.Equal(t, o, rm.GetEObject("id2"))
}
