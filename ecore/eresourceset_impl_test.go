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
	r.On("BasicSetResourceSet", rs, nil).Return(nil)
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
	mockResource.On("BasicSetResourceSet", rs, nil).Return(nil)
	rs.GetResources().Add(mockResource)

	// get registered resource - no loading
	mockResource.On("GetURI").Return(uri)
	assert.Equal(t, mockResource, rs.GetResource(uri, false))

	// get registered resource - loading
	mockResource.On("IsLoaded").Once().Return(false)
	mockResource.On("Load").Once()
	assert.Equal(t, mockResource, rs.GetResource(uri, true))
}
