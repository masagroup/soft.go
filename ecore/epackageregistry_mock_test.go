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
	rp := NewMockEPackageRegistry(t)
	p := NewMockEPackage(t)
	m := NewMockRun(t, p)
	rp.EXPECT().RegisterPackage(p).Return().Run(func(pack EPackage) { m.Run(p) }).Once()
	rp.RegisterPackage(p)
}

func TestMockEPackageRegistryUnRegisterPackage(t *testing.T) {
	rp := NewMockEPackageRegistry(t)
	p := NewMockEPackage(t)
	m := NewMockRun(t, p)
	rp.EXPECT().UnregisterPackage(p).Return().Run(func(pack EPackage) { m.Run(p) }).Once()
	rp.UnregisterPackage(p)
}

func TestMockEPackageRegistryPutPackage(t *testing.T) {
	rp := NewMockEPackageRegistry(t)
	p := NewMockEPackage(t)
	m := NewMockRun(t, "nsURI", p)
	rp.EXPECT().PutPackage("nsURI", p).Return().Run(func(nsURI string, pack EPackage) { m.Run(nsURI, pack) }).Once()
	rp.PutPackage("nsURI", p)
}

func TestMockEPackageRegistryPutSupplier(t *testing.T) {
	rp := NewMockEPackageRegistry(t)
	m := NewMockRun(t, "nsURI", mock.AnythingOfType("func() ecore.EPackage"))
	rp.EXPECT().PutSupplier("nsURI", mock.AnythingOfType("func() ecore.EPackage")).Return().Run(func(nsURI string, supplier func() EPackage) { m.Run(nsURI, supplier) }).Once()
	rp.PutSupplier("nsURI", func() EPackage {
		return nil
	})
}

func TestMockEPackageRegistryRemove(t *testing.T) {
	rp := NewMockEPackageRegistry(t)
	m := NewMockRun(t, "nsURI")
	rp.EXPECT().Remove("nsURI").Return().Run(func(nsURI string) { m.Run(nsURI) }).Once()
	rp.Remove("nsURI")
}

func TestMockEPackageRegistryGetPackage(t *testing.T) {
	rp := NewMockEPackageRegistry(t)
	p := NewMockEPackage(t)
	m := NewMockRun(t, "p")
	rp.EXPECT().GetPackage("p").Return(p).Run(func(nsURI string) { m.Run(nsURI) }).Once()
	rp.EXPECT().GetPackage("p").Call.Return(func(string) EPackage {
		return p
	}).Once()
	assert.Equal(t, p, rp.GetPackage("p"))
	assert.Equal(t, p, rp.GetPackage("p"))
}

func TestMockEPackageRegistryGetFactory(t *testing.T) {
	rp := NewMockEPackageRegistry(t)
	f := NewMockEFactory(t)
	m := NewMockRun(t, "f")
	rp.EXPECT().GetFactory("f").Return(f).Run(func(nsURI string) { m.Run(nsURI) }).Once()
	rp.EXPECT().GetFactory("f").Call.Return(func(string) EFactory {
		return f
	}).Once()
	assert.Equal(t, f, rp.GetFactory("f"))
	assert.Equal(t, f, rp.GetFactory("f"))
}
