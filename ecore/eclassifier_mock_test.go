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
	"reflect"
	"testing"
)

func discardMockEClassifier() {
	_ = assert.Equal
	_ = testing.Coverage
}

// TestMockEClassifierGetClassifierID tests method GetClassifierID
func TestMockEClassifierGetClassifierID(t *testing.T) {
	o := NewMockEClassifier(t)
	r := int(45)
	m := NewMockRun(t)
	o.EXPECT().GetClassifierID().Return(r).Run(func() { m.Run() }).Once()
	o.EXPECT().GetClassifierID().Call.Return(func() int { return r }).Once()
	assert.Equal(t, r, o.GetClassifierID())
	assert.Equal(t, r, o.GetClassifierID())
}

// TestMockEClassifierSetClassifierID tests method SetClassifierID
func TestMockEClassifierSetClassifierID(t *testing.T) {
	o := NewMockEClassifier(t)
	v := int(45)
	m := NewMockRun(t, v)
	o.EXPECT().SetClassifierID(v).Return().Run(func(_p0 int) { m.Run(_p0) }).Once()
	o.SetClassifierID(v)
}

// TestMockEClassifierGetDefaultValue tests method GetDefaultValue
func TestMockEClassifierGetDefaultValue(t *testing.T) {
	o := NewMockEClassifier(t)
	r := any(nil)
	m := NewMockRun(t)
	o.EXPECT().GetDefaultValue().Return(r).Run(func() { m.Run() }).Once()
	o.EXPECT().GetDefaultValue().Call.Return(func() any { return r }).Once()
	assert.Equal(t, r, o.GetDefaultValue())
	assert.Equal(t, r, o.GetDefaultValue())
}

// TestMockEClassifierGetEPackage tests method GetEPackage
func TestMockEClassifierGetEPackage(t *testing.T) {
	o := NewMockEClassifier(t)
	r := new(MockEPackage)
	m := NewMockRun(t)
	o.EXPECT().GetEPackage().Return(r).Run(func() { m.Run() }).Once()
	o.EXPECT().GetEPackage().Call.Return(func() EPackage { return r }).Once()
	assert.Equal(t, r, o.GetEPackage())
	assert.Equal(t, r, o.GetEPackage())
}

// TestMockEClassifierGetInstanceClass tests method GetInstanceClass
func TestMockEClassifierGetInstanceClass(t *testing.T) {
	o := NewMockEClassifier(t)
	r := reflect.Type(reflect.TypeOf(""))
	m := NewMockRun(t)
	o.EXPECT().GetInstanceClass().Return(r).Run(func() { m.Run() }).Once()
	o.EXPECT().GetInstanceClass().Call.Return(func() reflect.Type { return r }).Once()
	assert.Equal(t, r, o.GetInstanceClass())
	assert.Equal(t, r, o.GetInstanceClass())
}

// TestMockEClassifierSetInstanceClass tests method SetInstanceClass
func TestMockEClassifierSetInstanceClass(t *testing.T) {
	o := NewMockEClassifier(t)
	v := reflect.Type(reflect.TypeOf(""))
	m := NewMockRun(t, v)
	o.EXPECT().SetInstanceClass(v).Return().Run(func(_p0 reflect.Type) { m.Run(_p0) }).Once()
	o.SetInstanceClass(v)
}

// TestMockEClassifierGetInstanceClassName tests method GetInstanceClassName
func TestMockEClassifierGetInstanceClassName(t *testing.T) {
	o := NewMockEClassifier(t)
	r := string("Test String")
	m := NewMockRun(t)
	o.EXPECT().GetInstanceClassName().Return(r).Run(func() { m.Run() }).Once()
	o.EXPECT().GetInstanceClassName().Call.Return(func() string { return r }).Once()
	assert.Equal(t, r, o.GetInstanceClassName())
	assert.Equal(t, r, o.GetInstanceClassName())
}

// TestMockEClassifierSetInstanceClassName tests method SetInstanceClassName
func TestMockEClassifierSetInstanceClassName(t *testing.T) {
	o := NewMockEClassifier(t)
	v := string("Test String")
	m := NewMockRun(t, v)
	o.EXPECT().SetInstanceClassName(v).Return().Run(func(_p0 string) { m.Run(_p0) }).Once()
	o.SetInstanceClassName(v)
}

// TestMockEClassifierGetInstanceTypeName tests method GetInstanceTypeName
func TestMockEClassifierGetInstanceTypeName(t *testing.T) {
	o := NewMockEClassifier(t)
	r := string("Test String")
	m := NewMockRun(t)
	o.EXPECT().GetInstanceTypeName().Return(r).Run(func() { m.Run() }).Once()
	o.EXPECT().GetInstanceTypeName().Call.Return(func() string { return r }).Once()
	assert.Equal(t, r, o.GetInstanceTypeName())
	assert.Equal(t, r, o.GetInstanceTypeName())
}

// TestMockEClassifierSetInstanceTypeName tests method SetInstanceTypeName
func TestMockEClassifierSetInstanceTypeName(t *testing.T) {
	o := NewMockEClassifier(t)
	v := string("Test String")
	m := NewMockRun(t, v)
	o.EXPECT().SetInstanceTypeName(v).Return().Run(func(_p0 string) { m.Run(_p0) }).Once()
	o.SetInstanceTypeName(v)
}

// TestMockEClassifierIsInstance tests method IsInstance
func TestMockEClassifierIsInstance(t *testing.T) {
	o := NewMockEClassifier(t)
	object := any(nil)
	m := NewMockRun(t, object)
	r := bool(true)
	o.EXPECT().IsInstance(object).Return(r).Run(func(object any) { m.Run(object) }).Once()
	o.EXPECT().IsInstance(object).Call.Return(func() bool {
		return r
	}).Once()
	assert.Equal(t, r, o.IsInstance(object))
	assert.Equal(t, r, o.IsInstance(object))
	o.AssertExpectations(t)
}
