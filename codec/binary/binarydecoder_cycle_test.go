// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package binary

import (
	"bytes"
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// a package found by fuzzing: its class Item is among its own super types
var superTypesCycleInput = []byte("\xc4\b\x89emf\n\r\x1a\n\x00\x00\x00\xd9%http://www.eclipse.org/emf/2002/Ecore\x00\xd9%http://www.eclipse.org/emf/2002/Ecore\xa1/\x00\xa8EPackage\x02\xa4name\xa7derived\x03\xa5nsURI\xbahttp://example.com/derived\x04\xa8nsPrefix\xa7derived\x06\xaceClassifiers\x01\x01\x00\x01\xa6EClass\x02\xa4name\xa4Item\x0e\xabeSuperTypes\x01\x01\x00\x00\x10\x00e\xa7//Named\xe6\x00\a\xaceSubPackages\x01\x00\x00\x02\xa3sub\x03\xbehttp://example.com/derived/sub\x06\x01\x04\x00\x01\x02\xa4Part\x00\x00\xf4")

func TestBinaryDecoder_SuperTypesCycle_Fuzzed(t *testing.T) {
	eResource := ecore.NewEResourceSetImpl().CreateResource(ecore.NewURI("file:///cycle.bin"))
	eObject, err := NewDecoder(eResource, bytes.NewReader(superTypesCycleInput), nil).DecodeObject()
	assert.Nil(t, eObject)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cycle of super types")
}

// encodeSuperTypesCycle encodes a package of the classes a, b and c, with the super types given by
// their indices
func encodeSuperTypesCycle(t *testing.T, superTypes map[int][]int) []byte {
	ePackage := ecore.GetFactory().CreateEPackage()
	ePackage.SetName("cycle")
	ePackage.SetNsURI("http://example.com/cycle")
	classes := []ecore.EClass{}
	for _, name := range []string{"a", "b", "c"} {
		eClass := ecore.GetFactory().CreateEClass()
		eClass.SetName(name)
		eAttribute := ecore.GetFactory().CreateEAttribute()
		eAttribute.SetName(name + "Attribute")
		eAttribute.SetEType(ecore.GetPackage().GetEString())
		eClass.GetEStructuralFeatures().Add(eAttribute)
		ePackage.GetEClassifiers().Add(eClass)
		classes = append(classes, eClass)
	}
	for i, supers := range superTypes {
		for _, j := range supers {
			classes[i].GetESuperTypes().Add(classes[j])
		}
	}
	eResource := ecore.NewEResourceSetImpl().CreateResource(ecore.NewURI("file:///cycle.ecore"))
	eResource.GetContents().Add(ePackage)
	var buffer bytes.Buffer
	NewEncoder(eResource, &buffer, nil).EncodeResource()
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	return buffer.Bytes()
}

// a package whose super types close a cycle is refused, by DecodeResource and DecodeObject
func TestBinaryDecoder_SuperTypesCycle(t *testing.T) {
	for name, superTypes := range map[string]map[int][]int{
		"self":          {0: {0}},
		"two classes":   {0: {1}, 1: {0}},
		"three classes": {0: {1}, 1: {2}, 2: {0}},
		"above a cycle": {2: {0}, 0: {1}, 1: {0}},
	} {
		t.Run(name, func(t *testing.T) {
			data := encodeSuperTypesCycle(t, superTypes)

			eResource := ecore.NewEResourceSetImpl().CreateResource(ecore.NewURI("file:///cycle.bin"))
			NewDecoder(eResource, bytes.NewReader(data), nil).DecodeResource()
			require.Equal(t, 1, eResource.GetErrors().Size())
			assert.Contains(t, eResource.GetErrors().Get(0).(ecore.EDiagnostic).GetMessage(), "cycle of super types")
			assert.True(t, eResource.GetContents().Empty())
		})
	}
}

// a hierarchy without cycle, diamond included, is decoded
func TestBinaryDecoder_SuperTypesDiamond(t *testing.T) {
	data := encodeSuperTypesCycle(t, map[int][]int{1: {0}, 2: {0, 1}})
	eResource := ecore.NewEResourceSetImpl().CreateResource(ecore.NewURI("file:///diamond.bin"))
	NewDecoder(eResource, bytes.NewReader(data), nil).DecodeResource()
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	ePackage := eResource.GetContents().Get(0).(ecore.EPackage)
	c := ePackage.GetEClassifier("c").(ecore.EClass)
	assert.Equal(t, 3, c.GetEAllAttributes().Size())
}
