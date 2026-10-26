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

func TestXMIEncoderLibrarySimple(t *testing.T) {
	// load/save
	xmiProcessor := NewProcessor()
	resource := xmiProcessor.Load(ecore.NewURI("testdata/library.simple.ecore"))
	require.NotNil(t, resource)
	result := xmiProcessor.SaveToString(resource, nil)
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
	result := xmiProcessor.SaveToString(resource, nil)
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
	result := xmiProcessor.SaveToString(resource, nil)
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
	result := xmiProcessor.SaveToString(resource, map[string]any{xml.OPTION_USE_ENCODED_ATTRIBUTE_STYLE: true})
	// check
	bytes, err := os.ReadFile("testdata/library.encoded.ecore")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(result, "\r\n", "\n"))

	// without the option, as href elements
	result = xmiProcessor.SaveToString(resource, nil)
	assert.Contains(t, result, `<eType xsi:type="ecore:EClass" href="library.simple.ecore#//Book"/>`)
	assert.Contains(t, result, `<eSuperTypes href="library.simple.ecore#//Book"/>`)
}
