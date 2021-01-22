package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMockEPackageRegistryImpl_RegisterPackage(t *testing.T) {
	rp := NewEPackageRegistryImpl()
	p := &MockEPackage{}
	p.On("GetNsURI").Return("uri").Once()
	rp.RegisterPackage(p)
	mock.AssertExpectationsForObjects(t, p)
}

func TestMockEPackageRegistryImpl_RegisterPackageWithURI(t *testing.T) {
	rp := NewEPackageRegistryImpl()
	p := &MockEPackage{}
	rp.RegisterPackageWithURI(p, "nsURI")
}

func TestMockEPackageRegistryImpl_UnRegisterPackage(t *testing.T) {
	rp := NewEPackageRegistryImpl()
	p := &MockEPackage{}
	p.On("GetNsURI").Return("uri").Once()
	rp.UnregisterPackage(p)
	mock.AssertExpectationsForObjects(t, p)
}

func TestMockEPackageRegistryImpl_GetPackage(t *testing.T) {
	{
		rp := NewEPackageRegistryImpl()
		assert.Nil(t, rp.GetPackage("uri"))
	}
	{
		rp := NewEPackageRegistryImpl()
		p := &MockEPackage{}
		p.On("GetNsURI").Return("uri").Once()
		rp.RegisterPackage(p)
		mock.AssertExpectationsForObjects(t, p)
		assert.Equal(t, p, rp.GetPackage("uri"))
	}
	{
		delegate := &MockEPackageRegistry{}
		p := &MockEPackage{}
		rp := NewEPackageRegistryImplWithDelegate(delegate)
		delegate.On("GetPackage", "uri").Return(p).Once()
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
		p := &MockEPackage{}
		p.On("GetNsURI").Return("uri").Once()
		rp.RegisterPackage(p)
		mock.AssertExpectationsForObjects(t, p)

		f := &MockEFactory{}
		p.On("GetEFactoryInstance").Return(f).Once()
		assert.Equal(t, f, rp.GetFactory("uri"))
		mock.AssertExpectationsForObjects(t, p, f)
	}
	{
		delegate := &MockEPackageRegistry{}
		p := &MockEPackage{}
		f := &MockEFactory{}
		rp := NewEPackageRegistryImplWithDelegate(delegate)
		delegate.On("GetFactory", "uri").Return(f).Once()
		assert.Equal(t, f, rp.GetFactory("uri"))
		mock.AssertExpectationsForObjects(t, p, delegate)
	}
}
