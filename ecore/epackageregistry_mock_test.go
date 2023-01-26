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

func TestMockEPackageRegistryRegisterPackage(t *testing.T) {
	rp := &MockEPackageRegistry{}
	p := NewMockEPackage(t)
	rp.On("RegisterPackage", p).Once()
	rp.RegisterPackage(p)
	mock.AssertExpectationsForObjects(t, rp, p)
}

func TestMockEPackageRegistryUnRegisterPackage(t *testing.T) {
	rp := &MockEPackageRegistry{}
	p := NewMockEPackage(t)
	rp.On("UnregisterPackage", p).Once()
	rp.UnregisterPackage(p)
	mock.AssertExpectationsForObjects(t, rp, p)
}

func TestMockEPackageRegistryPutPackage(t *testing.T) {
	rp := &MockEPackageRegistry{}
	p := NewMockEPackage(t)
	rp.On("PutPackage", "nsURI", p).Once()
	rp.PutPackage("nsURI", p)
	mock.AssertExpectationsForObjects(t, rp, p)
}

func TestMockEPackageRegistryPutSupplier(t *testing.T) {
	rp := &MockEPackageRegistry{}
	rp.On("PutSupplier", "nsURI", mock.AnythingOfType("func() ecore.EPackage")).Once()
	rp.PutSupplier("nsURI", func() EPackage {
		return nil
	})
	mock.AssertExpectationsForObjects(t, rp)
}

func TestMockEPackageRegistryRemove(t *testing.T) {
	rp := &MockEPackageRegistry{}
	rp.On("Remove", "nsURI").Once()
	rp.Remove("nsURI")
	mock.AssertExpectationsForObjects(t, rp)
}

func TestMockEPackageRegistryGetPackage(t *testing.T) {
	rp := &MockEPackageRegistry{}
	p := NewMockEPackage(t)
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
	f := NewMockEFactory(t)
	rp.On("GetFactory", "f").Return(f).Once()
	rp.On("GetFactory", "f").Return(func(string) EFactory {
		return f
	}).Once()
	assert.Equal(t, f, rp.GetFactory("f"))
	assert.Equal(t, f, rp.GetFactory("f"))
	mock.AssertExpectationsForObjects(t, rp, f)
}
