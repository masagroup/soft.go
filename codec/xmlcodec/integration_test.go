// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package xmlcodec

import (
	"testing"

	"github.com/masagroup/soft.go/ecore"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertResource(t *testing.T, resource ecore.EResource) {
	require.NotNil(t, resource)
	require.True(t, resource.IsLoaded())
	require.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	require.True(t, resource.GetWarnings().Empty(), diagnosticError(resource.GetWarnings()))
}

func TestEAllContents(t *testing.T) {
	xmiProcessor := NewXMIProcessor()
	shopEcoreResource := xmiProcessor.Load(ecore.CreateFileURI("testdata/shop.ecore"))
	assertResource(t, shopEcoreResource)
	shopPackage, _ := shopEcoreResource.GetContents().Get(0).(ecore.EPackage)
	require.NotNil(t, shopPackage)
	ordersEcoreResource := xmiProcessor.Load(ecore.CreateFileURI("testdata/orders.ecore"))
	assertResource(t, ordersEcoreResource)
	ordersPackage, _ := ordersEcoreResource.GetContents().Get(0).(ecore.EPackage)
	require.NotNil(t, ordersPackage)

	xmlProcessor := NewXMLProcessor(XMLProcessorPackages([]ecore.EPackage{shopPackage, ordersPackage}))
	ordersResource := xmlProcessor.Load(ecore.CreateFileURI("testdata/orders.xml"))
	assertResource(t, ordersResource)
	ordersRoot, _ := ordersResource.GetContents().Get(0).(ecore.EObject)
	require.NotNil(t, ordersRoot)

	i := 0
	for it := ordersResource.GetAllContents(); it.HasNext(); {
		eObject, _ := it.Next().(ecore.EObject)
		require.NotNil(t, eObject)
		i++
	}
	assert.Equal(t, 11, i)
}

func TestGetEObjectInResource(t *testing.T) {
	xmiProcessor := NewXMIProcessor()
	shopEcoreResource := xmiProcessor.Load(ecore.CreateFileURI("testdata/shop.ecore"))
	assertResource(t, shopEcoreResource)
	eProductClass := shopEcoreResource.GetEObject("//Product")
	require.NotNil(t, eProductClass)
}

func TestResourceLoadInvalid(t *testing.T) {
	r := ecore.NewEResourceImpl()
	r.SetURI(ecore.NewURI("testdata/invalid.xml"))
	r.Load()
	assert.False(t, r.IsLoaded())
	assert.False(t, r.GetErrors().Empty())
}

// The core tests build the models of delete.xmi and eallcontents.xml in code: check here that
// the files decode.
func TestLoadCoreTestModels(t *testing.T) {
	deletePackage := loadPackage("delete.ecore")
	require.NotNil(t, deletePackage)
	deleteResource := NewXMIProcessor(XMIProcessorPackages([]ecore.EPackage{deletePackage})).Load(ecore.NewURI("testdata/delete.xmi"))
	assertResource(t, deleteResource)
	deleteRoot := deleteResource.GetContents().Get(0).(ecore.EObject)
	deleteRootClass := deletePackage.GetEClassifier("Root").(ecore.EClass)
	for name, size := range map[string]int{"a": 3, "b": 3, "o": 1, "c": 3} {
		assert.Equal(t, size, deleteRoot.EGet(deleteRootClass.GetEStructuralFeatureFromName(name)).(ecore.EList).Size(), name)
	}

	eAllContentsPackage := loadPackage("eallcontents.ecore")
	require.NotNil(t, eAllContentsPackage)
	eAllContentsResource := NewXMLProcessor(XMLProcessorPackages([]ecore.EPackage{eAllContentsPackage})).Load(ecore.NewURI("testdata/eallcontents.xml"))
	assertResource(t, eAllContentsResource)
	i := 0
	for it := eAllContentsResource.GetAllContents(); it.HasNext(); it.Next() {
		i++
	}
	// root, theater, 3 parties, 4 formations, 6 automats and 5 units
	assert.Equal(t, 20, i)
}
