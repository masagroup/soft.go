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
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockEResourceFactoryCreateResource(t *testing.T) {
	f := &MockEResourceFactory{}
	r := &MockEResource{}
	uri, _ := url.Parse("test://file.t")
	f.On("CreateResource", uri).Return(r).Once()
	f.On("CreateResource", uri).Return(func(*url.URL) EResource {
		return r
	}).Once()
	assert.Equal(t, r, f.CreateResource(uri))
	assert.Equal(t, r, f.CreateResource(uri))
	r.AssertExpectations(t)
}
