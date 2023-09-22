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

import "github.com/stretchr/testify/assert"
import "testing"

func discardMockEPackage() {
	_ = assert.Equal
	_ = testing.Coverage
}

// TestMockEPackageGetEClassifiers tests method GetEClassifiers
func TestMockEPackageGetEClassifiers(t *testing.T) {
	o := NewMockEPackage(t)
	l := NewMockEList(t)
	m := NewMockRun(t)
	o.EXPECT().GetEClassifiers().Return(l).Run(func() { m.Run() }).Once()
	o.EXPECT().GetEClassifiers().Call.Return(func() EList { return l }).Once()
	assert.Equal(t, l, o.GetEClassifiers())
	assert.Equal(t, l, o.GetEClassifiers())
}

// TestMockEPackageGetEFactoryInstance tests method GetEFactoryInstance
func TestMockEPackageGetEFactoryInstance(t *testing.T) {
	o := NewMockEPackage(t)
	r := NewMockEFactory(t)
	m := NewMockRun(t)
	o.EXPECT().GetEFactoryInstance().Return(r).Run(func() { m.Run() }).Once()
	o.EXPECT().GetEFactoryInstance().Call.Return(func() EFactory { return r }).Once()
	assert.Equal(t, r, o.GetEFactoryInstance())
	assert.Equal(t, r, o.GetEFactoryInstance())
}

// TestMockEPackageSetEFactoryInstance tests method SetEFactoryInstance
func TestMockEPackageSetEFactoryInstance(t *testing.T) {
	o := NewMockEPackage(t)
	v := NewMockEFactory(t)
	m := NewMockRun(t, v)
	o.EXPECT().SetEFactoryInstance(v).Return().Run(func(_p0 EFactory) { m.Run(_p0) }).Once()
	o.SetEFactoryInstance(v)
}

// TestMockEPackageGetESubPackages tests method GetESubPackages
func TestMockEPackageGetESubPackages(t *testing.T) {
	o := NewMockEPackage(t)
	l := NewMockEList(t)
	m := NewMockRun(t)
	o.EXPECT().GetESubPackages().Return(l).Run(func() { m.Run() }).Once()
	o.EXPECT().GetESubPackages().Call.Return(func() EList { return l }).Once()
	assert.Equal(t, l, o.GetESubPackages())
	assert.Equal(t, l, o.GetESubPackages())
}

// TestMockEPackageGetESuperPackage tests method GetESuperPackage
func TestMockEPackageGetESuperPackage(t *testing.T) {
	o := NewMockEPackage(t)
	r := NewMockEPackage(t)
	m := NewMockRun(t)
	o.EXPECT().GetESuperPackage().Return(r).Run(func() { m.Run() }).Once()
	o.EXPECT().GetESuperPackage().Call.Return(func() EPackage { return r }).Once()
	assert.Equal(t, r, o.GetESuperPackage())
	assert.Equal(t, r, o.GetESuperPackage())
}

// TestMockEPackageGetNsPrefix tests method GetNsPrefix
func TestMockEPackageGetNsPrefix(t *testing.T) {
	o := NewMockEPackage(t)
	r := string("Test String")
	m := NewMockRun(t)
	o.EXPECT().GetNsPrefix().Return(r).Run(func() { m.Run() }).Once()
	o.EXPECT().GetNsPrefix().Call.Return(func() string { return r }).Once()
	assert.Equal(t, r, o.GetNsPrefix())
	assert.Equal(t, r, o.GetNsPrefix())
}

// TestMockEPackageSetNsPrefix tests method SetNsPrefix
func TestMockEPackageSetNsPrefix(t *testing.T) {
	o := NewMockEPackage(t)
	v := string("Test String")
	m := NewMockRun(t, v)
	o.EXPECT().SetNsPrefix(v).Return().Run(func(_p0 string) { m.Run(_p0) }).Once()
	o.SetNsPrefix(v)
}

// TestMockEPackageGetNsURI tests method GetNsURI
func TestMockEPackageGetNsURI(t *testing.T) {
	o := NewMockEPackage(t)
	r := string("Test String")
	m := NewMockRun(t)
	o.EXPECT().GetNsURI().Return(r).Run(func() { m.Run() }).Once()
	o.EXPECT().GetNsURI().Call.Return(func() string { return r }).Once()
	assert.Equal(t, r, o.GetNsURI())
	assert.Equal(t, r, o.GetNsURI())
}

// TestMockEPackageSetNsURI tests method SetNsURI
func TestMockEPackageSetNsURI(t *testing.T) {
	o := NewMockEPackage(t)
	v := string("Test String")
	m := NewMockRun(t, v)
	o.EXPECT().SetNsURI(v).Return().Run(func(_p0 string) { m.Run(_p0) }).Once()
	o.SetNsURI(v)
}

// TestMockEPackageGetEClassifier tests method GetEClassifier
func TestMockEPackageGetEClassifier(t *testing.T) {
	o := NewMockEPackage(t)
	name := string("Test String")
	m := NewMockRun(t, name)
	r := NewMockEClassifier(t)
	o.EXPECT().GetEClassifier(name).Return(r).Run(func(name string) { m.Run(name) }).Once()
	o.EXPECT().GetEClassifier(name).Call.Return(func() EClassifier {
		return r
	}).Once()
	assert.Equal(t, r, o.GetEClassifier(name))
	assert.Equal(t, r, o.GetEClassifier(name))
}
