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
	mockResource.EXPECT().BasicSetResourceSet(rs, nil).Return(nil).Once()
	rs.GetResources().Add(mockResource)

	// get registered resource - no loading
	mockResource.EXPECT().GetURI().Return(uri).Once()
	assert.Equal(t, mockResource, rs.GetResource(uri, false))

	// get registered resource - loading
	mockResource.EXPECT().GetURI().Return(uri).Once()
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

func TestEResourceSet_GetResourceCodecRegistry(t *testing.T) {
	rs := NewEResourceSetImpl()
	assert.NotNil(t, rs.GetResourceCodecRegistry())

	mockResourceCodecRegistry := NewMockEResourceCodecRegistry(t)
	rs.SetResourceCodecRegistry(mockResourceCodecRegistry)
	assert.Equal(t, mockResourceCodecRegistry, rs.GetResourceCodecRegistry())
}

func TestEResourceSet_GetURIResourceMap(t *testing.T) {
	rs := NewEResourceSetImpl()
	assert.Nil(t, rs.GetURIResourceMap())

	mockURIResourceMap := map[*URI]EResource{}
	rs.SetURIResourceMap(mockURIResourceMap)
	assert.Equal(t, mockURIResourceMap, rs.GetURIResourceMap())
}
