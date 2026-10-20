// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package xml

import (
	"strings"
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// An unknown element, nested content included, must not cost the known siblings after it.
func TestXMLDecoderUnknownElementKeepsLaterSiblings(t *testing.T) {
	factory := ecore.GetFactory()
	pkg := factory.CreateEPackage()
	pkg.SetName("lenient")
	pkg.SetNsPrefix("lenient")
	pkg.SetNsURI("urn:lenient")
	rootClass, itemClass := factory.CreateEClass(), factory.CreateEClass()
	rootClass.SetName("Root")
	itemClass.SetName("Item")
	name := factory.CreateEAttribute()
	name.SetName("name")
	name.SetEType(ecore.GetPackage().GetEString())
	itemClass.GetEStructuralFeatures().Add(name)
	items := factory.CreateEReference()
	items.SetName("items")
	items.SetEType(itemClass)
	items.SetContainment(true)
	items.SetUpperBound(-1)
	rootClass.GetEStructuralFeatures().Add(items)
	pkg.GetEClassifiers().AddAll(ecore.NewImmutableEList([]any{rootClass, itemClass}))

	document := `<?xml version="1.0" encoding="UTF-8"?>
<lenient:Root xmlns:lenient="urn:lenient">
  <items name="first"/>
  <unknown><nested><deeper/></nested></unknown>
  <items name="second"/>
</lenient:Root>`
	processor := NewProcessor(ProcessorPackages([]ecore.EPackage{pkg}))
	resource := processor.GetResourceSet().CreateResource(ecore.NewURI("lenient.xml"))
	resource.LoadWithReader(strings.NewReader(document), nil)

	require.Equal(t, 1, resource.GetContents().Size())
	root := resource.GetContents().Get(0).(ecore.EObject)
	list := root.EGet(items).(ecore.EList)
	require.Equal(t, 2, list.Size(), "the item after the unknown element was dropped")
	assert.Equal(t, "second", list.Get(1).(ecore.EObject).EGet(name))
	assert.False(t, resource.GetErrors().Empty(), "the unknown element is still reported")
}
