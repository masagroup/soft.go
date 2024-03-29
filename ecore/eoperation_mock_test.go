// Code generated by soft.generator.go. DO NOT EDIT.

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
	"github.com/stretchr/testify/assert"
	"testing"
)

func discardMockEOperation() {
	_ = assert.Equal
	_ = testing.Coverage
}

// TestMockEOperationGetEContainingClass tests method GetEContainingClass
func TestMockEOperationGetEContainingClass(t *testing.T) {
	o := NewMockEOperation(t)
	r := NewMockEClass(t)
	m := NewMockRun(t)
	o.EXPECT().GetEContainingClass().Return(r).Run(func() { m.Run() }).Once()
	o.EXPECT().GetEContainingClass().Call.Return(func() EClass { return r }).Once()
	assert.Equal(t, r, o.GetEContainingClass())
	assert.Equal(t, r, o.GetEContainingClass())
}

// TestMockEOperationGetEExceptions tests method GetEExceptions
func TestMockEOperationGetEExceptions(t *testing.T) {
	o := NewMockEOperation(t)
	l := NewMockEList(t)
	m := NewMockRun(t)
	o.EXPECT().GetEExceptions().Return(l).Run(func() { m.Run() }).Once()
	o.EXPECT().GetEExceptions().Call.Return(func() EList { return l }).Once()
	assert.Equal(t, l, o.GetEExceptions())
	assert.Equal(t, l, o.GetEExceptions())
}

// TestMockEOperationUnsetEExceptions tests method UnsetEExceptions
func TestMockEOperationUnsetEExceptions(t *testing.T) {
	o := NewMockEOperation(t)
	m := NewMockRun(t)
	o.EXPECT().UnsetEExceptions().Return().Run(func() { m.Run() }).Once()
	o.UnsetEExceptions()
}

// TestMockEOperationGetEParameters tests method GetEParameters
func TestMockEOperationGetEParameters(t *testing.T) {
	o := NewMockEOperation(t)
	l := NewMockEList(t)
	m := NewMockRun(t)
	o.EXPECT().GetEParameters().Return(l).Run(func() { m.Run() }).Once()
	o.EXPECT().GetEParameters().Call.Return(func() EList { return l }).Once()
	assert.Equal(t, l, o.GetEParameters())
	assert.Equal(t, l, o.GetEParameters())
}

// TestMockEOperationGetOperationID tests method GetOperationID
func TestMockEOperationGetOperationID(t *testing.T) {
	o := NewMockEOperation(t)
	r := int(45)
	m := NewMockRun(t)
	o.EXPECT().GetOperationID().Return(r).Run(func() { m.Run() }).Once()
	o.EXPECT().GetOperationID().Call.Return(func() int { return r }).Once()
	assert.Equal(t, r, o.GetOperationID())
	assert.Equal(t, r, o.GetOperationID())
}

// TestMockEOperationSetOperationID tests method SetOperationID
func TestMockEOperationSetOperationID(t *testing.T) {
	o := NewMockEOperation(t)
	v := int(45)
	m := NewMockRun(t, v)
	o.EXPECT().SetOperationID(v).Return().Run(func(_p0 int) { m.Run(_p0) }).Once()
	o.SetOperationID(v)
}

// TestMockEOperationIsOverrideOf tests method IsOverrideOf
func TestMockEOperationIsOverrideOf(t *testing.T) {
	o := NewMockEOperation(t)
	someOperation := NewMockEOperation(t)
	m := NewMockRun(t, someOperation)
	r := bool(true)
	o.EXPECT().IsOverrideOf(someOperation).Return(r).Run(func(someOperation EOperation) { m.Run(someOperation) }).Once()
	o.EXPECT().IsOverrideOf(someOperation).Call.Return(func() bool {
		return r
	}).Once()
	assert.Equal(t, r, o.IsOverrideOf(someOperation))
	assert.Equal(t, r, o.IsOverrideOf(someOperation))
}
