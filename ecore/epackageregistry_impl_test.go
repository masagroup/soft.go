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
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMockEPackageRegistryImpl_RegisterPackage(t *testing.T) {
	rp := NewEPackageRegistryImpl()
	p := NewMockEPackage(t)
	p.EXPECT().GetNsURI().Return("uri").Once()
	rp.RegisterPackage(p)
	mock.AssertExpectationsForObjects(t, p)
}

func TestMockEPackageRegistryImpl_UnRegisterPackage(t *testing.T) {
	rp := NewEPackageRegistryImpl()
	p := NewMockEPackage(t)
	p.EXPECT().GetNsURI().Return("uri").Once()
	rp.UnregisterPackage(p)
	mock.AssertExpectationsForObjects(t, p)
}

func TestMockEPackageRegistryImpl_RegisterPutPackage(t *testing.T) {
	rp := NewEPackageRegistryImpl()
	p := NewMockEPackage(t)
	rp.PutPackage("nsURI", p)
	assert.Equal(t, p, rp.GetPackage("nsURI"))
}

func TestMockEPackageRegistryImpl_RegisterPutSupplier(t *testing.T) {
	rp := NewEPackageRegistryImpl()
	p := NewMockEPackage(t)
	f := func() EPackage {
		return p
	}
	rp.PutSupplier("nsURI", f)
	assert.Equal(t, p, rp.GetPackage("nsURI"))
}

func TestMockEPackageRegistryImpl_Remove(t *testing.T) {
	rp := NewEPackageRegistryImpl()
	p := NewMockEPackage(t)
	rp.PutPackage("nsURI", p)
	assert.Equal(t, p, rp.GetPackage("nsURI"))
	rp.Remove("nsURI")
	assert.Nil(t, rp.GetPackage("nsURI"))
}

func TestMockEPackageRegistryImpl_GetPackage(t *testing.T) {
	{
		rp := NewEPackageRegistryImpl()
		assert.Nil(t, rp.GetPackage("uri"))
	}
	{
		rp := NewEPackageRegistryImpl()
		p := NewMockEPackage(t)
		p.EXPECT().GetNsURI().Return("uri").Once()
		rp.RegisterPackage(p)
		mock.AssertExpectationsForObjects(t, p)
		assert.Equal(t, p, rp.GetPackage("uri"))
	}
	{
		delegate := &MockEPackageRegistry{}
		p := NewMockEPackage(t)
		rp := NewEPackageRegistryImplWithDelegate(delegate)
		delegate.EXPECT().GetPackage("uri").Return(p).Once()
		assert.Equal(t, p, rp.GetPackage("uri"))
		mock.AssertExpectationsForObjects(t, p, delegate)
	}
}

func TestMockEPackageRegistryImpl_GetFactory(t *testing.T) {
	{
		rp := NewEPackageRegistryImpl()
		assert.Nil(t, rp.GetFactory("uri"))
	}
	{
		rp := NewEPackageRegistryImpl()
		p := NewMockEPackage(t)
		p.EXPECT().GetNsURI().Return("uri").Once()
		rp.RegisterPackage(p)
		mock.AssertExpectationsForObjects(t, p)

		f := NewMockEFactory(t)
		p.EXPECT().GetEFactoryInstance().Return(f).Once()
		assert.Equal(t, f, rp.GetFactory("uri"))
		mock.AssertExpectationsForObjects(t, p, f)
	}
	{
		delegate := &MockEPackageRegistry{}
		p := NewMockEPackage(t)
		f := NewMockEFactory(t)
		rp := NewEPackageRegistryImplWithDelegate(delegate)
		delegate.EXPECT().GetFactory("uri").Return(f).Once()
		assert.Equal(t, f, rp.GetFactory("uri"))
		mock.AssertExpectationsForObjects(t, p, delegate)
	}
}

func TestEPackageRegistryImpl_PutRemoveFactory(t *testing.T) {
	p := NewMockEPackage(t)
	packageFactory := NewMockEFactory(t)
	f := NewMockEFactory(t)
	rp := NewEPackageRegistryImpl()
	rp.PutPackage("uri", p)

	// the registered factory, before the factory instance of the package
	rp.PutFactory("uri", f)
	assert.Equal(t, f, rp.GetFactory("uri"))

	// removed: the factory instance of the package
	rp.RemoveFactory("uri")
	p.EXPECT().GetEFactoryInstance().Return(packageFactory).Once()
	assert.Equal(t, packageFactory, rp.GetFactory("uri"))

	// a factory without its package registered
	rp.PutFactory("other", f)
	assert.Equal(t, f, rp.GetFactory("other"))
	assert.Nil(t, rp.GetPackage("other"))
}

func TestEPackageRegistryImpl_RegisterUnregisterFactory(t *testing.T) {
	p := NewMockEPackage(t)
	f := NewMockEFactory(t)
	f.EXPECT().GetEPackage().Return(p).Twice()
	p.EXPECT().GetNsURI().Return("uri").Twice()
	rp := NewEPackageRegistryImpl()

	// keyed by the nsURI of the package of the factory
	rp.RegisterFactory(f)
	assert.Equal(t, f, rp.GetFactory("uri"))

	rp.UnregisterFactory(f)
	assert.Nil(t, rp.GetFactory("uri"))

	// a factory without package
	noPackage := NewMockEFactory(t)
	noPackage.EXPECT().GetEPackage().Return(nil).Once()
	assert.PanicsWithValue(t, "the factory has no package", func() { rp.RegisterFactory(noPackage) })
}

func TestEPackageRegistryImpl_GetFactoryDelegate(t *testing.T) {
	f := NewMockEFactory(t)
	delegate := NewEPackageRegistryImpl()
	delegate.PutFactory("uri", f)
	// a factory of the delegate, without package nor factory in the registry
	rp := NewEPackageRegistryImplWithDelegate(delegate)
	assert.Equal(t, f, rp.GetFactory("uri"))
}

// Concurrent first callers must all get the one registry, built once.
func TestGetPackageRegistryConcurrentInitialization(t *testing.T) {
	previous := packageRegistryInstance
	t.Cleanup(func() {
		packageRegistryInstance = previous
		packageRegistryOnce = sync.Once{}
		packageRegistryOnce.Do(func() {})
	})
	packageRegistryInstance = nil
	packageRegistryOnce = sync.Once{}
	registries := make(chan EPackageRegistry, 100)
	var group sync.WaitGroup
	for range cap(registries) {
		group.Add(1)
		go func() {
			defer group.Done()
			registries <- GetPackageRegistry()
		}()
	}
	group.Wait()
	close(registries)

	want := <-registries
	assert.NotNil(t, want)
	for registry := range registries {
		assert.Same(t, want, registry)
	}
}
