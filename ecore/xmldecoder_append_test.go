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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AddUnique appends without the duplicate check Add makes on a unique list.
func TestAbstractEListAddUnique(t *testing.T) {
	l := NewUniqueBasicEList(nil)
	l.Add(1)
	assert.False(t, l.Add(1))
	l.AddUnique(2)
	assert.Equal(t, []any{1, 2}, l.ToArray())
}

type longListModel struct {
	pkg   EPackage
	items EReference
	name  EAttribute
}

func newLongListModel() *longListModel {
	f := GetFactory()
	m := &longListModel{pkg: f.CreateEPackage()}
	m.pkg.SetName("long")
	m.pkg.SetNsPrefix("long")
	m.pkg.SetNsURI("urn:long")
	root, item := f.CreateEClass(), f.CreateEClass()
	root.SetName("Root")
	item.SetName("Item")
	m.name = f.CreateEAttribute()
	m.name.SetName("name")
	m.name.SetEType(GetPackage().GetEString())
	item.GetEStructuralFeatures().Add(m.name)
	m.items = f.CreateEReference()
	m.items.SetName("items")
	m.items.SetEType(item)
	m.items.SetContainment(true)
	m.items.SetUpperBound(-1)
	root.GetEStructuralFeatures().Add(m.items)
	m.pkg.GetEClassifiers().AddAll(NewImmutableEList([]any{root, item}))
	return m
}

func longListDocument(n int) string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n" + `<long:Root xmlns:long="urn:long">` + "\n")
	for i := range n {
		fmt.Fprintf(&b, "  <items name=\"item%d\"/>\n", i)
	}
	b.WriteString("</long:Root>\n")
	return b.String()
}

func (m *longListModel) load(tb testing.TB, document string) EList {
	processor := NewXMLProcessor(XMLProcessorPackages([]EPackage{m.pkg}))
	resource := processor.GetResourceSet().CreateResource(NewURI("long.xml"))
	resource.LoadWithReader(strings.NewReader(document), nil)
	require.True(tb, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	return resource.GetContents().Get(0).(EObject).EGet(m.items).(EList)
}

// A long many-valued feature loads every item, in document order.
func TestXMLDecoderLoadsALongListInOrder(t *testing.T) {
	m := newLongListModel()
	items := m.load(t, longListDocument(2000))
	require.Equal(t, 2000, items.Size())
	assert.Equal(t, "item0", items.Get(0).(EObject).EGet(m.name))
	assert.Equal(t, "item1999", items.Get(1999).(EObject).EGet(m.name))
}

func BenchmarkXMLDecoderLongList(b *testing.B) {
	m := newLongListModel()
	document := longListDocument(20000)
	for b.Loop() {
		m.load(b, document)
	}
}
