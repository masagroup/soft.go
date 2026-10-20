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

type testContentAdapter struct {
	ecore.EContentAdapter
}

func BenchmarkEContentAdapterWithBigModel(b *testing.B) {
	for b.Loop() {
		// load package
		ePackage := loadPackage("library.complex.ecore")
		assert.NotNil(b, ePackage)

		// load resource
		xmlProcessor := NewXMLProcessor(XMLProcessorPackages([]ecore.EPackage{ePackage}))
		eResource := xmlProcessor.Load(ecore.NewURI("testdata/library.complex.big.xml"))
		require.NotNil(b, eResource)
		assert.True(b, eResource.IsLoaded())
		assert.True(b, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
		assert.True(b, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

		adapter := new(testContentAdapter)
		adapter.SetInterfaces(adapter)

		eResource.EAdapters().Add(adapter)
		eResource.EAdapters().Remove(adapter)
	}
}

func BenchmarkEContentAdapterWithTreeModel(b *testing.B) {
	for b.Loop() {
		// load package
		ePackage := loadPackage("tree.ecore")
		assert.NotNil(b, ePackage)

		// load resource
		xmlProcessor := NewXMLProcessor(XMLProcessorPackages([]ecore.EPackage{ePackage}))
		eResource := xmlProcessor.Load(ecore.NewURI("testdata/tree.xml"))
		require.NotNil(b, eResource)
		assert.True(b, eResource.IsLoaded())
		assert.True(b, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
		assert.True(b, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

		adapter := new(testContentAdapter)
		adapter.SetInterfaces(adapter)

		eResource.EAdapters().Add(adapter)
		eResource.EAdapters().Remove(adapter)
	}
}
