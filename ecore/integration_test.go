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

func assertResource(t *testing.T, resource EResource) {
	require.NotNil(t, resource)
	require.True(t, resource.IsLoaded())
	require.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	require.True(t, resource.GetWarnings().Empty(), diagnosticError(resource.GetWarnings()))
}

func Test_EAllContents(t *testing.T) {
	xmiProcessor := NewXMIProcessor()
	shopEcoreResource := xmiProcessor.Load(CreateFileURI("testdata/shop.ecore"))
	assertResource(t, shopEcoreResource)
	shopPackage, _ := shopEcoreResource.GetContents().Get(0).(EPackage)
	require.NotNil(t, shopPackage)
	ordersEcoreResource := xmiProcessor.Load(CreateFileURI("testdata/orders.ecore"))
	assertResource(t, ordersEcoreResource)
	ordersPackage, _ := ordersEcoreResource.GetContents().Get(0).(EPackage)
	require.NotNil(t, ordersPackage)

	xmlProcessor := NewXMLProcessor([]EPackage{shopPackage, ordersPackage})
	ordersResource := xmlProcessor.Load(CreateFileURI("testdata/orders.xml"))
	assertResource(t, ordersResource)
	ordersRoot, _ := ordersResource.GetContents().Get(0).(EObject)
	require.NotNil(t, ordersRoot)

	i := 0
	for it := ordersResource.GetAllContents(); it.HasNext(); {
		eObject, _ := it.Next().(EObject)
		require.NotNil(t, eObject)
		i++
	}
	assert.Equal(t, 11, i)
}
