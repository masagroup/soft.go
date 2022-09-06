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
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewXmlProcessor(t *testing.T) {
	p := NewXMLProcessor(nil)
	require.NotNil(t, p)
	assert.NotNil(t, p.GetResourceSet())
}

func TestNewSharedXmlProcessor(t *testing.T) {
	mockResourceSet := &MockEResourceSet{}
	p := NewSharedXMLProcessor(mockResourceSet)
	require.NotNil(t, p)
	assert.Equal(t, mockResourceSet, p.GetResourceSet())
}

func xmlProcessorLoad(t *testing.T, xmlProcessor *XMLProcessor, path string) EObject {
	resource := xmlProcessor.Load(NewURI(path))
	require.NotNil(t, resource)
	require.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	eObject, _ := resource.GetContents().Get(0).(EObject)
	require.NotNil(t, eObject)
	return eObject
}

func TestSaveObject(t *testing.T) {

	testsCases := []struct {
		name  string
		meta  string
		model string
	}{
		{"shop", "shop.ecore", "shop.xml"},
		{"orders", "orders.ecore", "orders.xml"},
	}

	resourceSet := NewEResourceSetImpl()
	packageRegistry := resourceSet.GetPackageRegistry()
	xmlProcessor := NewSharedXMLProcessor(resourceSet)
	for _, testCase := range testsCases {
		t.Run(testCase.name, func(t *testing.T) {
			ePackage, _ := xmlProcessorLoad(t, xmlProcessor, "testdata/"+testCase.meta).(EPackage)
			require.NotNil(t, ePackage)
			packageRegistry.RegisterPackage(ePackage)
			eModel := xmlProcessorLoad(t, xmlProcessor, "testdata/"+testCase.model)
			require.NotNil(t, eModel)

			resultName := "testdata/" + testCase.name + ".result.xml"
			xmlProcessor.SaveObject(NewURI(resultName), eModel)

			// src
			src, err := os.ReadFile("testdata/" + testCase.model)
			assert.Nil(t, err)

			// result
			result, err := os.ReadFile(resultName)
			assert.Nil(t, err)
			assert.Equal(t, strings.ReplaceAll(string(src), "\r\n", "\n"), strings.ReplaceAll(string(result), "\r\n", "\n"))
		})
	}

}

type nodeFactory struct {
	eFactory        EFactory
	eNodeClass      EClass
	eNameAttribute  EAttribute
	eNodesReference EReference
}

func newNodeFactory(ePackage EPackage) *nodeFactory {
	f := new(nodeFactory)
	f.eFactory = ePackage.GetEFactoryInstance()
	f.eNodeClass, _ = ePackage.GetEClassifier("Node").(EClass)
	f.eNameAttribute, _ = f.eNodeClass.GetEStructuralFeatureFromName("name").(EAttribute)
	f.eNodesReference, _ = f.eNodeClass.GetEStructuralFeatureFromName("nodes").(EReference)
	return f
}

func (f *nodeFactory) newNode(name string, depth int) EObject {
	n := f.eFactory.Create(f.eNodeClass)
	n.ESet(f.eNameAttribute, name)
	if depth > 0 {
		children := n.EGet(f.eNodesReference).(EList)
		for i := 0; i < rand.Intn(5); i++ {
			c := f.newNode(fmt.Sprintf("%v.%v", name, i), depth-1)
			children.Add(c)
		}
	}
	return n
}

func TestSerializationTree(t *testing.T) {
	// load package
	ePackage := loadPackage("tree.ecore")
	assert.NotNil(t, ePackage)

	f := newNodeFactory(ePackage)
	eNode := f.newNode("0", 5)

	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	xmlProcessor.SaveObject(NewURI("testdata/tree.xml"), eNode)
}
