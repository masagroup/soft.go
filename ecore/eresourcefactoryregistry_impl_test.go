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

func TestEResourceFactoryRegistryGetFactoryProtocol(t *testing.T) {
	mockFactory := new(MockEResourceFactory)

	rfr := NewEResourceFactoryRegistryImpl()
	rfr.GetProtocolToFactoryMap()["test"] = mockFactory

	assert.Equal(t, mockFactory, rfr.GetFactory(&url.URL{Scheme: "test", Path: "//file.t"}))
	assert.Nil(t, rfr.GetFactory(&url.URL{Scheme: "file", Path: "//file.t"}))
}

func TestEResourceFactoryRegistryGetFactoryExtension(t *testing.T) {
	mockFactory := new(MockEResourceFactory)

	rfr := NewEResourceFactoryRegistryImpl()
	rfr.GetExtensionToFactoryMap()["test"] = mockFactory

	assert.Equal(t, mockFactory, rfr.GetFactory(&url.URL{Scheme: "file", Path: "//file.test"}))
	assert.Nil(t, rfr.GetFactory(&url.URL{Scheme: "file", Path: "//file.t"}))
}
