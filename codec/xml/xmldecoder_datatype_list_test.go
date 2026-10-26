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
	"fmt"
	"strings"
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type tagListModel struct {
	pkg  ecore.EPackage
	tags ecore.EAttribute
}

func newTagListModel() *tagListModel {
	f := ecore.GetFactory()
	m := &tagListModel{pkg: f.CreateEPackage()}
	m.pkg.SetName("tags")
	m.pkg.SetNsPrefix("tags")
	m.pkg.SetNsURI("urn:tags")
	root := f.CreateEClass()
	root.SetName("Root")
	m.tags = f.CreateEAttribute()
	m.tags.SetName("tags")
	m.tags.SetEType(ecore.GetPackage().GetEString())
	m.tags.SetUpperBound(-1)
	root.GetEStructuralFeatures().Add(m.tags)
	m.pkg.GetEClassifiers().Add(root)
	return m
}

func tagListDocument(n int) string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n" + `<tags:Root xmlns:tags="urn:tags">` + "\n")
	for i := range n {
		fmt.Fprintf(&b, "  <tags>tag%d</tags>\n", i)
	}
	b.WriteString("</tags:Root>\n")
	return b.String()
}

func (m *tagListModel) load(tb testing.TB, document string) ecore.EList {
	processor := NewProcessor(ProcessorPackages([]ecore.EPackage{m.pkg}))
	resource := processor.GetResourceSet().CreateResource(ecore.NewURI("tags.xml"))
	resource.LoadWithReader(strings.NewReader(document), nil)
	require.True(tb, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	return resource.GetContents().Get(0).(ecore.EObject).EGet(m.tags).(ecore.EList)
}

// A long many-valued data type feature loads every value, in document order.
func TestXMLDecoderLoadsALongDataTypeListInOrder(t *testing.T) {
	m := newTagListModel()
	tags := m.load(t, tagListDocument(2000))
	require.Equal(t, 2000, tags.Size())
	assert.Equal(t, "tag0", tags.Get(0))
	assert.Equal(t, "tag1999", tags.Get(1999))
}

func BenchmarkXMLDecoderLongDataTypeList(b *testing.B) {
	m := newTagListModel()
	document := tagListDocument(20000)
	for b.Loop() {
		m.load(b, document)
	}
}
