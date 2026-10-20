// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/masagroup/soft.go/codec/xmi"
	"github.com/masagroup/soft.go/ecore"
)

func diagnosticError(errors ecore.EList) string {
	if errors.Empty() {
		return ""
	} else {
		return errors.Get(0).(ecore.EDiagnostic).GetMessage()
	}
}

func loadPackage(packageFileName string) ecore.EPackage {
	xmiProcessor := xmi.NewProcessor()
	eResource := xmiProcessor.Load(ecore.NewURI("testdata/" + packageFileName))
	if eResource.IsLoaded() && eResource.GetContents().Size() > 0 {
		ePackage, _ := eResource.GetContents().Get(0).(ecore.EPackage)
		return ePackage
	} else {
		return nil
	}
}

func loadTestPackage(t *testing.T, resourceSet ecore.EResourceSet, packageURI *ecore.URI) (ecore.EResource, ecore.EPackage) {
	// load package
	r := resourceSet.CreateResource(packageURI)
	r.SetObjectIDManager(ecore.NewIncrementalIDManager())
	r.Load()
	assert.True(t, r.IsLoaded())
	assert.True(t, r.GetErrors().Empty(), diagnosticError(r.GetErrors()))
	assert.True(t, r.GetWarnings().Empty(), diagnosticError(r.GetWarnings()))

	// retrieve package
	ePackage, _ := r.GetContents().Get(0).(ecore.EPackage)
	require.NotNil(t, ePackage)
	resourceSet.GetPackageRegistry().RegisterPackage(ePackage)
	return r, ePackage
}

func loadTestModel(t *testing.T, resourceSet ecore.EResourceSet, modelURI *ecore.URI) (ecore.EResource, ecore.EObject) {
	// load package
	r := resourceSet.CreateResource(modelURI)
	r.SetObjectIDManager(ecore.NewIncrementalIDManager())
	r.Load()
	require.True(t, r.IsLoaded())
	require.True(t, r.GetErrors().Empty(), diagnosticError(r.GetErrors()))
	require.True(t, r.GetWarnings().Empty(), diagnosticError(r.GetWarnings()))
	require.Equal(t, 1, r.GetContents().Size())

	// retrieve root object
	return r, r.GetContents().Get(0).(ecore.EObject)
}

