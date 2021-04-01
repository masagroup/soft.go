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
)

func TestXMIResourceFactoryCreate(t *testing.T) {
	f := &XMIResourceFactory{}
	uri := &URI{}
	r := f.CreateResource(uri)
	assert.NotNil(t, r)
	assert.Equal(t, uri, r.GetURI())
}
