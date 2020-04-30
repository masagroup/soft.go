// *****************************************************************************
//
// This file is part of a MASA library or program.
// Refer to the included end-user license agreement for restrictions.
//
// Copyright (c) 2020 MASA Group
//
// *****************************************************************************

// *****************************************************************************
//
// Warning: This file was generated by soft.generator.go Generator
//
// *****************************************************************************

package ecore

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func discardMockEPackage() {
	_ = assert.Equal
	_ = testing.Coverage
}

// TestGetEClassifiers tests method GetEClassifiers
func TestGetEClassifiers(t *testing.T) {
	o := &MockEPackage{}
	l := &MockEList{}
	o.On("GetEClassifiers").Once().Return(l)
	assert.Equal(t, l, o.GetEClassifiers())

	o.On("GetEClassifiers").Once().Return(func() EList {
		return l
	})
	assert.Equal(t, l, o.GetEClassifiers())
	o.AssertExpectations(t)
}

// TestMockEPackageGetEFactoryInstance tests method GetEFactoryInstance
func TestMockEPackageGetEFactoryInstance(t *testing.T) {
	o := &MockEPackage{}
	r := &MockEFactory{}
	o.On("GetEFactoryInstance").Once().Return(r)
	assert.Equal(t, r, o.GetEFactoryInstance())

	o.On("GetEFactoryInstance").Once().Return(func() EFactory {
		return r
	})
	assert.Equal(t, r, o.GetEFactoryInstance())
	o.AssertExpectations(t)
}

// TestMockEPackageSetEFactoryInstance tests method SetEFactoryInstance
func TestMockEPackageSetEFactoryInstance(t *testing.T) {
	o := &MockEPackage{}
	v := &MockEFactory{}
	o.On("SetEFactoryInstance", v).Once()

	o.SetEFactoryInstance(v)
	o.AssertExpectations(t)
}

// TestGetESubPackages tests method GetESubPackages
func TestGetESubPackages(t *testing.T) {
	o := &MockEPackage{}
	l := &MockEList{}
	o.On("GetESubPackages").Once().Return(l)
	assert.Equal(t, l, o.GetESubPackages())

	o.On("GetESubPackages").Once().Return(func() EList {
		return l
	})
	assert.Equal(t, l, o.GetESubPackages())
	o.AssertExpectations(t)
}

// TestMockEPackageGetESuperPackage tests method GetESuperPackage
func TestMockEPackageGetESuperPackage(t *testing.T) {
	o := &MockEPackage{}
	r := &MockEPackage{}
	o.On("GetESuperPackage").Once().Return(r)
	assert.Equal(t, r, o.GetESuperPackage())

	o.On("GetESuperPackage").Once().Return(func() EPackage {
		return r
	})
	assert.Equal(t, r, o.GetESuperPackage())
	o.AssertExpectations(t)
}

// TestMockEPackageGetNsPrefix tests method GetNsPrefix
func TestMockEPackageGetNsPrefix(t *testing.T) {
	o := &MockEPackage{}
	r := "Test String"
	o.On("GetNsPrefix").Once().Return(r)
	assert.Equal(t, r, o.GetNsPrefix())

	o.On("GetNsPrefix").Once().Return(func() string {
		return r
	})
	assert.Equal(t, r, o.GetNsPrefix())
	o.AssertExpectations(t)
}

// TestMockEPackageSetNsPrefix tests method SetNsPrefix
func TestMockEPackageSetNsPrefix(t *testing.T) {
	o := &MockEPackage{}
	v := "Test String"
	o.On("SetNsPrefix", v).Once()

	o.SetNsPrefix(v)
	o.AssertExpectations(t)
}

// TestMockEPackageGetNsURI tests method GetNsURI
func TestMockEPackageGetNsURI(t *testing.T) {
	o := &MockEPackage{}
	r := "Test String"
	o.On("GetNsURI").Once().Return(r)
	assert.Equal(t, r, o.GetNsURI())

	o.On("GetNsURI").Once().Return(func() string {
		return r
	})
	assert.Equal(t, r, o.GetNsURI())
	o.AssertExpectations(t)
}

// TestMockEPackageSetNsURI tests method SetNsURI
func TestMockEPackageSetNsURI(t *testing.T) {
	o := &MockEPackage{}
	v := "Test String"
	o.On("SetNsURI", v).Once()

	o.SetNsURI(v)
	o.AssertExpectations(t)
}
