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
	"github.com/stretchr/testify/require"
)

func TestEResourceSetConstructor(t *testing.T) {
	rs := NewEResourceSetImpl()
	assert.Nil(t, rs.GetURIResourceMap())
}

func TestEResourceSetResourcesWithMock(t *testing.T) {
	rs := NewEResourceSetImpl()
	r := NewMockEResourceInternal(t)
	r.EXPECT().GetURI().Return(nil).Once()
	r.EXPECT().BasicSetResourceSet(rs, nil).Return(nil).Once()
	rs.GetResources().Add(r)
}

func TestEResourceSetResourcesNoMock(t *testing.T) {
	rs := NewEResourceSetImpl()
	r := NewEResourceImpl()

	rs.GetResources().Add(r)
	assert.Equal(t, rs, r.GetResourceSet())

	rs.GetResources().Remove(r)
	assert.Equal(t, nil, r.GetResourceSet())
}

func TestEResourceSetCreateResource(t *testing.T) {
	uri, _ := ParseURI("test://file.t")
	rs := NewEResourceSetImpl()
	r := rs.CreateResource(uri)
	require.NotNil(t, r)
	assert.Equal(t, rs, r.GetResourceSet())
}

func TestEResourceSetGetRegisteredResource(t *testing.T) {

	uri, _ := ParseURI("test://file.t")
	rs := NewEResourceSetImpl()

	// register resource
	mockResource := NewMockEResourceInternal(t)
	mockResource.EXPECT().GetURI().Return(uri)
	mockResource.EXPECT().BasicSetResourceSet(rs, nil).Return(nil).Once()
	rs.GetResources().Add(mockResource)

	// get registered resource - no loading
	assert.Equal(t, mockResource, rs.GetResource(uri, false))

	// get registered resource - loading
	mockResource.EXPECT().IsLoaded().Once().Return(false)
	mockResource.EXPECT().Load().Once()
	assert.Equal(t, mockResource, rs.GetResource(uri, true))
}

func TestEResourceSet_GetURIConverter(t *testing.T) {
	rs := NewEResourceSetImpl()
	assert.NotNil(t, rs.GetURIConverter())

	mockURIConverter := NewMockEURIConverter(t)
	rs.SetURIConverter(mockURIConverter)
	assert.Equal(t, mockURIConverter, rs.GetURIConverter())
}

func TestEResourceSet_GetPackageRegistry(t *testing.T) {
	rs := NewEResourceSetImpl()
	assert.NotNil(t, rs.GetPackageRegistry())

	mockPackageRegistry := NewMockEPackageRegistry(t)
	rs.SetPackageRegistry(mockPackageRegistry)
	assert.Equal(t, mockPackageRegistry, rs.GetPackageRegistry())
}

func TestEResourceSet_GetCodecRegistry(t *testing.T) {
	rs := NewEResourceSetImpl()
	assert.NotNil(t, rs.GetCodecRegistry())

	mockResourceCodecRegistry := NewMockECodecRegistry(t)
	rs.SetCodecRegistry(mockResourceCodecRegistry)
	assert.Equal(t, mockResourceCodecRegistry, rs.GetCodecRegistry())
}

func TestEResourceSet_GetURIResourceMap(t *testing.T) {
	rs := NewEResourceSetImpl()
	assert.Nil(t, rs.GetURIResourceMap())

	mockURIResourceMap := map[*URI]EResource{}
	rs.SetURIResourceMap(mockURIResourceMap)
	assert.Equal(t, mockURIResourceMap, rs.GetURIResourceMap())
}

func newCoexistingDynamicPackage(nsURI string) EPackage {
	ePackage := GetFactory().CreateEPackage()
	ePackage.SetName("coexisting")
	ePackage.SetNsPrefix("coexisting")
	ePackage.SetNsURI(nsURI)
	ePackage.SetEFactoryInstance(GetFactory().CreateEFactory())
	eClass := GetFactory().CreateEClass()
	eClass.SetName("Element")
	eAttribute := GetFactory().CreateEAttribute()
	eAttribute.SetName("name")
	eAttribute.SetEType(GetPackage().GetEString())
	eClass.GetEStructuralFeatures().Add(eAttribute)
	ePackage.GetEClassifiers().Add(eClass)
	eResource := NewEResourceImpl()
	eResource.SetURI(NewURI(nsURI))
	eResource.GetContents().Add(ePackage)
	return ePackage
}

func TestEResourceSet_CoexistingImplementations(t *testing.T) {
	nsURI := "http://masagroup.net/soft/test/coexisting"
	uri := NewURI(nsURI)
	package1 := newCoexistingDynamicPackage(nsURI)
	package2 := newCoexistingDynamicPackage(nsURI)
	require.NotSame(t, package1, package2)
	require.NotSame(t, package1.GetEFactoryInstance(), package2.GetEFactoryInstance())

	rs1 := CreateEResourceSet([]EPackage{package1})
	rs2 := NewEResourceSetImpl()
	rs2.GetPackageRegistry().RegisterPackage(GetPackage())
	rs2.GetPackageRegistry().RegisterPackage(package2)

	// each resource set resolves the nsURI to its own package
	assert.Same(t, package1, rs1.GetPackageRegistry().GetPackage(nsURI))
	assert.Same(t, package2, rs2.GetPackageRegistry().GetPackage(nsURI))
	assert.Same(t, package1.EResource(), rs1.GetResource(uri, false))
	assert.Same(t, package2.EResource(), rs2.GetResource(uri, false))
	assert.NotSame(t, rs1.GetResource(uri, false), rs2.GetResource(uri, false))

	// each resource set gets its own factory
	assert.Same(t, package1.GetEFactoryInstance(), rs1.GetPackageRegistry().GetFactory(nsURI))
	assert.Same(t, package2.GetEFactoryInstance(), rs2.GetPackageRegistry().GetFactory(nsURI))

	// both resource sets share ecore
	assert.Same(t, GetPackage(), rs1.GetPackageRegistry().GetPackage(NS_URI))
	assert.Same(t, GetPackage(), rs2.GetPackageRegistry().GetPackage(NS_URI))

	// global registry untouched
	assert.Nil(t, GetPackageRegistry().GetPackage(nsURI))
	assert.Nil(t, GetPackageRegistry().GetFactory(nsURI))
}
