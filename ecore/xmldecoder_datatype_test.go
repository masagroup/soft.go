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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A data type contained by no package is reported as a load error instead of crashing the decode.
func TestXMLDecoderDataTypeWithNoPackageIsReportedNotDereferenced(t *testing.T) {
	factory := GetFactory()

	// Deliberately never added to the package's classifiers: this is the whole fixture.
	detached := factory.CreateEDataType()
	detached.SetName("Detached")
	detached.SetInstanceTypeName("string")

	single := factory.CreateEAttribute()
	single.SetName("value")
	single.SetEType(detached)

	many := factory.CreateEAttribute()
	many.SetName("values")
	many.SetEType(detached)
	many.SetUpperBound(UNBOUNDED_MULTIPLICITY)

	eClass := factory.CreateEClass()
	eClass.SetName("Root")
	eClass.GetEStructuralFeatures().Add(single)
	eClass.GetEStructuralFeatures().Add(many)

	ePackage := factory.CreateEPackage()
	ePackage.SetName("detachedtype")
	ePackage.SetNsPrefix("d")
	ePackage.SetNsURI("http://masa/test/detachedtype")
	ePackage.GetEClassifiers().Add(eClass)

	eResourceSet := NewEResourceSetImpl()
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)
	eResource := eResourceSet.CreateResource(NewURI("detachedtype.xml"))
	require.NotNil(t, eResource)

	document := `<?xml version="1.0" encoding="UTF-8"?>
<d:Root xmlns:d="http://masa/test/detachedtype" value="7" values="1 2 3"/>`
	require.NotPanics(t, func() {
		eResource.LoadWithReader(strings.NewReader(document), nil)
	})

	errors := eResource.GetErrors()
	require.False(t, errors.Empty(), "a metamodel the decoder cannot use has to be reported")
	reported := ""
	for it := errors.Iterator(); it.HasNext(); {
		reported += it.Next().(EDiagnostic).GetMessage() + "\n"
	}
	assert.Contains(t, reported, "Detached", "the diagnostic names the type")
	assert.Contains(t, reported, "value", "and the feature that led to it")
}
