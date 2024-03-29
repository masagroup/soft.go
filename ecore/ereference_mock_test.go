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

func discardMockEReference() {
	_ = assert.Equal
	_ = testing.Coverage
}

// TestMockEReferenceIsContainer tests method IsContainer
func TestMockEReferenceIsContainer(t *testing.T) {
	o := NewMockEReference(t)
	r := bool(true)
	m := NewMockRun(t)
	o.EXPECT().IsContainer().Return(r).Run(func() { m.Run() }).Once()
	o.EXPECT().IsContainer().Call.Return(func() bool { return r }).Once()
	assert.Equal(t, r, o.IsContainer())
	assert.Equal(t, r, o.IsContainer())
}

// TestMockEReferenceIsContainment tests method IsContainment
func TestMockEReferenceIsContainment(t *testing.T) {
	o := NewMockEReference(t)
	r := bool(true)
	m := NewMockRun(t)
	o.EXPECT().IsContainment().Return(r).Run(func() { m.Run() }).Once()
	o.EXPECT().IsContainment().Call.Return(func() bool { return r }).Once()
	assert.Equal(t, r, o.IsContainment())
	assert.Equal(t, r, o.IsContainment())
}

// TestMockEReferenceSetContainment tests method SetContainment
func TestMockEReferenceSetContainment(t *testing.T) {
	o := NewMockEReference(t)
	v := bool(true)
	m := NewMockRun(t, v)
	o.EXPECT().SetContainment(v).Return().Run(func(_p0 bool) { m.Run(_p0) }).Once()
	o.SetContainment(v)
}

// TestMockEReferenceGetEKeys tests method GetEKeys
func TestMockEReferenceGetEKeys(t *testing.T) {
	o := NewMockEReference(t)
	l := NewMockEList(t)
	m := NewMockRun(t)
	o.EXPECT().GetEKeys().Return(l).Run(func() { m.Run() }).Once()
	o.EXPECT().GetEKeys().Call.Return(func() EList { return l }).Once()
	assert.Equal(t, l, o.GetEKeys())
	assert.Equal(t, l, o.GetEKeys())
}

// TestMockEReferenceGetEOpposite tests method GetEOpposite
func TestMockEReferenceGetEOpposite(t *testing.T) {
	o := NewMockEReference(t)
	r := NewMockEReference(t)
	m := NewMockRun(t)
	o.EXPECT().GetEOpposite().Return(r).Run(func() { m.Run() }).Once()
	o.EXPECT().GetEOpposite().Call.Return(func() EReference { return r }).Once()
	assert.Equal(t, r, o.GetEOpposite())
	assert.Equal(t, r, o.GetEOpposite())
}

// TestMockEReferenceSetEOpposite tests method SetEOpposite
func TestMockEReferenceSetEOpposite(t *testing.T) {
	o := NewMockEReference(t)
	v := NewMockEReference(t)
	m := NewMockRun(t, v)
	o.EXPECT().SetEOpposite(v).Return().Run(func(_p0 EReference) { m.Run(_p0) }).Once()
	o.SetEOpposite(v)
}

// TestMockEReferenceGetEReferenceType tests method GetEReferenceType
func TestMockEReferenceGetEReferenceType(t *testing.T) {
	o := NewMockEReference(t)
	r := NewMockEClass(t)
	m := NewMockRun(t)
	o.EXPECT().GetEReferenceType().Return(r).Run(func() { m.Run() }).Once()
	o.EXPECT().GetEReferenceType().Call.Return(func() EClass { return r }).Once()
	assert.Equal(t, r, o.GetEReferenceType())
	assert.Equal(t, r, o.GetEReferenceType())
}

// TestMockEReferenceIsResolveProxies tests method IsResolveProxies
func TestMockEReferenceIsResolveProxies(t *testing.T) {
	o := NewMockEReference(t)
	r := bool(true)
	m := NewMockRun(t)
	o.EXPECT().IsResolveProxies().Return(r).Run(func() { m.Run() }).Once()
	o.EXPECT().IsResolveProxies().Call.Return(func() bool { return r }).Once()
	assert.Equal(t, r, o.IsResolveProxies())
	assert.Equal(t, r, o.IsResolveProxies())
}

// TestMockEReferenceSetResolveProxies tests method SetResolveProxies
func TestMockEReferenceSetResolveProxies(t *testing.T) {
	o := NewMockEReference(t)
	v := bool(true)
	m := NewMockRun(t, v)
	o.EXPECT().SetResolveProxies(v).Return().Run(func(_p0 bool) { m.Run(_p0) }).Once()
	o.SetResolveProxies(v)
}
