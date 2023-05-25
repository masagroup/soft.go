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

func TestMockEResourceSetGetResources(t *testing.T) {
	rs := NewMockEResourceSet(t)
	l := NewMockEList(t)
	m := NewMockRun(t)
	rs.EXPECT().GetResources().Return(l).Run(func() { m.Run() }).Once()
	rs.EXPECT().GetResources().Call.Return(func() EList {
		return l
	}).Once()
	assert.Equal(t, l, rs.GetResources())
	assert.Equal(t, l, rs.GetResources())
}

func TestMockEResourceSetGetResource(t *testing.T) {
	rs := NewMockEResourceSet(t)
	r := NewMockEResource(t)
	uri, _ := ParseURI("test://file.t")
	m := NewMockRun(t, uri, false)
	rs.EXPECT().GetResource(uri, false).Return(r).Run(func(uri *URI, loadOnDemand bool) { m.Run(uri, loadOnDemand) }).Once()
	rs.EXPECT().GetResource(uri, true).Call.Return(func(uri *URI, loadOnDemand bool) EResource {
		return r
	}).Once()
	assert.Equal(t, r, rs.GetResource(uri, false))
	assert.Equal(t, r, rs.GetResource(uri, true))
}

func TestMockEResourceSetCreateResource(t *testing.T) {
	rs := NewMockEResourceSet(t)
	r := NewMockEResource(t)
	uri, _ := ParseURI("test://file.t")
	m := NewMockRun(t, uri)
	rs.EXPECT().CreateResource(uri).Return(r).Run(func(uri *URI) { m.Run(uri) }).Once()
	rs.EXPECT().CreateResource(uri).Call.Return(func(uri *URI) EResource {
		return r
	}).Once()
	assert.Equal(t, r, rs.CreateResource(uri))
	assert.Equal(t, r, rs.CreateResource(uri))
}

func TestMockEResourceSetGetEObject(t *testing.T) {
	rs := NewMockEResourceSet(t)
	o := NewMockEObject(t)
	uri, _ := ParseURI("test://file.t")
	m := NewMockRun(t, uri, false)
	rs.EXPECT().GetEObject(uri, false).Return(o).Run(func(uri *URI, loadOnDemand bool) { m.Run(uri, loadOnDemand) }).Once()
	rs.EXPECT().GetEObject(uri, true).Call.Return(func(uri *URI, loadOnDemand bool) EObject {
		return o
	}).Once()
	assert.Equal(t, o, rs.GetEObject(uri, false))
	assert.Equal(t, o, rs.GetEObject(uri, true))
}

func TestMockEResourceSetGetURIConverter(t *testing.T) {
	rs := NewMockEResourceSet(t)
	c := NewMockEURIConverter(t)
	m := NewMockRun(t)
	rs.EXPECT().GetURIConverter().Return(c).Run(func() { m.Run() }).Once()
	rs.EXPECT().GetURIConverter().Call.Return(func() EURIConverter {
		return c
	}).Once()
	assert.Equal(t, c, rs.GetURIConverter())
	assert.Equal(t, c, rs.GetURIConverter())
}

func TestMockEResourceSetSetURIConverter(t *testing.T) {
	rs := NewMockEResourceSet(t)
	c := NewMockEURIConverter(t)
	m := NewMockRun(t, c)
	rs.EXPECT().SetURIConverter(c).Return().Run(func(uriConverter EURIConverter) { m.Run(uriConverter) }).Once()
	rs.SetURIConverter(c)
}

func TestMockEResourceSetGetPackageRegistry(t *testing.T) {
	rs := NewMockEResourceSet(t)
	pr := NewMockEPackageRegistry(t)
	m := NewMockRun(t)
	rs.EXPECT().GetPackageRegistry().Return(pr).Run(func() { m.Run() }).Once()
	rs.EXPECT().GetPackageRegistry().Call.Return(func() EPackageRegistry {
		return pr
	}).Once()
	assert.Equal(t, pr, rs.GetPackageRegistry())
	assert.Equal(t, pr, rs.GetPackageRegistry())
}

func TestMockEResourceSetSetPackageRegistry(t *testing.T) {
	rs := NewMockEResourceSet(t)
	pr := NewMockEPackageRegistry(t)
	m := NewMockRun(t, pr)
	rs.EXPECT().SetPackageRegistry(pr).Return().Run(func(packageregistry EPackageRegistry) { m.Run(packageregistry) }).Once()
	rs.SetPackageRegistry(pr)
	mock.AssertExpectationsForObjects(t, rs, pr)
}

func TestMockEResourceSetGetCodecRegistry(t *testing.T) {
	rs := NewMockEResourceSet(t)
	pr := NewMockECodecRegistry(t)
	m := NewMockRun(t)
	rs.EXPECT().GetCodecRegistry().Return(pr).Run(func() { m.Run() }).Once()
	rs.EXPECT().GetCodecRegistry().Call.Return(func() ECodecRegistry {
		return pr
	}).Once()
	assert.Equal(t, pr, rs.GetCodecRegistry())
	assert.Equal(t, pr, rs.GetCodecRegistry())
}

func TestMockEResourceSetSetResourceCodecRegistry(t *testing.T) {
	rs := NewMockEResourceSet(t)
	pr := NewMockECodecRegistry(t)
	m := NewMockRun(t, pr)
	rs.EXPECT().SetResourceCodecRegistry(pr).Return().Run(func(packageregistry ECodecRegistry) { m.Run(packageregistry) }).Once()
	rs.SetResourceCodecRegistry(pr)
}

func TestMockEResourceSetGetURIResourceMap(t *testing.T) {
	rs := NewMockEResourceSet(t)
	pr := make(map[*URI]EResource)
	m := NewMockRun(t)
	rs.EXPECT().GetURIResourceMap().Return(pr).Run(func() { m.Run() }).Once()
	rs.EXPECT().GetURIResourceMap().Call.Return(func() map[*URI]EResource {
		return pr
	}).Once()
	assert.Equal(t, pr, rs.GetURIResourceMap())
	assert.Equal(t, pr, rs.GetURIResourceMap())
}

func TestMockEResourceSetSetURIResourceMap(t *testing.T) {
	rs := NewMockEResourceSet(t)
	pr := make(map[*URI]EResource)
	m := NewMockRun(t, pr)
	rs.EXPECT().SetURIResourceMap(pr).Return().Run(func(mp map[*URI]EResource) { m.Run(mp) }).Once()
	rs.SetURIResourceMap(pr)
	mock.AssertExpectationsForObjects(t, rs)
}
