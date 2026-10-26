// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package xmi

import (
	"os"
	"strings"
	"testing"

	"github.com/masagroup/soft.go/codec/xml"
	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// the options of a save of the libraries in the style of the XMI codec: the cross resource references
// are href elements and the lines are not wrapped
var xmiStyle = map[string]any{xml.OPTION_USE_ENCODED_ATTRIBUTE_STYLE: false, xml.OPTION_LINE_WIDTH: 0}

func TestXMIEncoderLibrarySimple(t *testing.T) {
	// load/save
	xmiProcessor := NewProcessor()
	resource := xmiProcessor.Load(ecore.NewURI("testdata/library.simple.ecore"))
	require.NotNil(t, resource)
	result := xmiProcessor.SaveToString(resource, xmiStyle)
	// check
	bytes, err := os.ReadFile("testdata/library.simple.ecore")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(result, "\r\n", "\n"))
}

func TestXMIEncoderLibraryNoRoot(t *testing.T) {
	// load/save
	xmiProcessor := NewProcessor()
	resource := xmiProcessor.Load(ecore.NewURI("testdata/library.noroot.ecore"))
	require.NotNil(t, resource)
	result := xmiProcessor.SaveToString(resource, xmiStyle)
	// check
	bytes, err := os.ReadFile("testdata/library.noroot.ecore")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(result, "\r\n", "\n"))
}

func TestXMIEncoderLibraryComplex(t *testing.T) {
	// load/save
	xmiProcessor := NewProcessor()
	resource := xmiProcessor.Load(ecore.NewURI("testdata/library.complex.ecore"))
	require.NotNil(t, resource)
	result := xmiProcessor.SaveToString(resource, xmiStyle)
	// check
	bytes, err := os.ReadFile("testdata/library.complex.ecore")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(result, "\r\n", "\n"))
}

func BenchmarkXMIEncoderLibrarySimple(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		resource := ecore.NewEResourceImpl()
		resource.SetURI(ecore.NewURI("testdata/library.simple.ecore"))
		resource.Load()

		var strbuff strings.Builder
		resource.SaveWithWriter(&strbuff, nil)
		resource = nil
	}
}

func BenchmarkXMIEncoderLibraryNoRoot(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		resource := ecore.NewEResourceImpl()
		resource.SetURI(ecore.NewURI("testdata/library.noroot.ecore"))
		resource.Load()

		var strbuff strings.Builder
		resource.SaveWithWriter(&strbuff, nil)
		resource = nil
	}
}

func TestXMIEncoderLibraryEncodedAttributeStyle(t *testing.T) {
	// load/save with the cross resource references as attributes, as EMF saves .ecore files
	xmiProcessor := NewProcessor()
	resource := xmiProcessor.Load(ecore.NewURI("testdata/library.encoded.ecore"))
	require.NotNil(t, resource)
	require.True(t, resource.GetErrors().Empty())
	bytes, err := os.ReadFile("testdata/library.encoded.ecore")
	assert.Nil(t, err)
	expected := strings.ReplaceAll(string(bytes), "\r\n", "\n")
	// the default style of a .ecore resource, and with the options
	for _, options := range []map[string]any{nil, {xml.OPTION_USE_ENCODED_ATTRIBUTE_STYLE: true, xml.OPTION_LINE_WIDTH: 80}} {
		result := xmiProcessor.SaveToString(resource, options)
		assert.Equal(t, expected, strings.ReplaceAll(result, "\r\n", "\n"))
	}

	// without line width, an element on one line
	result := xmiProcessor.SaveToString(resource, map[string]any{xml.OPTION_LINE_WIDTH: 0})
	assert.Contains(t, result, `<eStructuralFeatures xsi:type="ecore:EReference" name="books" upperBound="-1" eType="#//Novel" containment="true"/>`)

	// with the options off, as href elements
	result = xmiProcessor.SaveToString(resource, xmiStyle)
	assert.Contains(t, result, `<eType xsi:type="ecore:EClass" href="library.simple.ecore#//Book"/>`)
	assert.Contains(t, result, `<eSuperTypes href="library.simple.ecore#//Book"/>`)

	// the default style of a .xmi resource is href elements
	resource.SetURI(ecore.NewURI("testdata/library.encoded.xmi"))
	result = xmiProcessor.SaveToString(resource, nil)
	assert.Contains(t, result, `<eSuperTypes href="library.simple.ecore#//Book"/>`)
}
