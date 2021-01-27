// *****************************************************************************
//
// This file is part of a MASA library or program.
// Refer to the included end-user license agreement for restrictions.
//
// Copyright (c) 2020 MASA Group
//
// *****************************************************************************

package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMockEPackageRegistryRegisterPackage(t *testing.T) {
	rp := &MockEPackageRegistry{}
	p := &MockEPackage{}
	rp.On("RegisterPackage", p).Once()
	rp.RegisterPackage(p)
	mock.AssertExpectationsForObjects(t, rp, p)
}

func TestMockEPackageRegistryRegisterPackageWithURI(t *testing.T) {
	rp := &MockEPackageRegistry{}
	p := &MockEPackage{}
	rp.On("RegisterPackageWithURI", p, "nsURI").Once()
	rp.RegisterPackageWithURI(p, "nsURI")
	mock.AssertExpectationsForObjects(t, rp, p)
}

func TestMockEPackageRegistryUnRegisterPackage(t *testing.T) {
	rp := &MockEPackageRegistry{}
	p := &MockEPackage{}
	rp.On("UnregisterPackage", p).Once()
	rp.UnregisterPackage(p)
	mock.AssertExpectationsForObjects(t, rp, p)
}

func TestMockEPackageRegistryGetPackage(t *testing.T) {
	rp := &MockEPackageRegistry{}
	p := &MockEPackage{}
	rp.On("GetPackage", "p").Return(p).Once()
	rp.On("GetPackage", "p").Return(func(string) EPackage {
		return p
	}).Once()
	assert.Equal(t, p, rp.GetPackage("p"))
	assert.Equal(t, p, rp.GetPackage("p"))
	mock.AssertExpectationsForObjects(t, rp, p)
}

func TestMockEPackageRegistryGetFactory(t *testing.T) {
	rp := &MockEPackageRegistry{}
	f := &MockEFactory{}
	rp.On("GetFactory", "f").Return(f).Once()
	rp.On("GetFactory", "f").Return(func(string) EFactory {
		return f
	}).Once()
	assert.Equal(t, f, rp.GetFactory("f"))
	assert.Equal(t, f, rp.GetFactory("f"))
	mock.AssertExpectationsForObjects(t, rp, f)
}
