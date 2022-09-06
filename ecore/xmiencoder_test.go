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
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestXMIEncoderLibrarySimple(t *testing.T) {
	// load/save
	xmiProcessor := NewXMIProcessor()
	resource := xmiProcessor.Load(NewURI("testdata/library.simple.ecore"))
	require.NotNil(t, resource)
	result := xmiProcessor.SaveToString(resource, nil)
	// check
	bytes, err := os.ReadFile("testdata/library.simple.ecore")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(result, "\r\n", "\n"))
}

func TestXMIEncoderLibraryNoRoot(t *testing.T) {
	// load/save
	xmiProcessor := NewXMIProcessor()
	resource := xmiProcessor.Load(NewURI("testdata/library.noroot.ecore"))
	require.NotNil(t, resource)
	result := xmiProcessor.SaveToString(resource, nil)
	// check
	bytes, err := os.ReadFile("testdata/library.noroot.ecore")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(result, "\r\n", "\n"))
}

func TestXMIEncoderLibraryComplex(t *testing.T) {
	// load/save
	xmiProcessor := NewXMIProcessor()
	resource := xmiProcessor.Load(NewURI("testdata/library.complex.ecore"))
	require.NotNil(t, resource)
	result := xmiProcessor.SaveToString(resource, nil)
	// check
	bytes, err := os.ReadFile("testdata/library.complex.ecore")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(result, "\r\n", "\n"))
}

func BenchmarkXMIEncoderLibrarySimple(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		resource := NewEResourceImpl()
		resource.SetURI(NewURI("testdata/library.simple.ecore"))
		resource.Load()

		var strbuff strings.Builder
		resource.SaveWithWriter(&strbuff, nil)
		resource = nil
	}
}

func BenchmarkXMIEncoderLibraryNoRoot(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		resource := NewEResourceImpl()
		resource.SetURI(NewURI("testdata/library.noroot.ecore"))
		resource.Load()

		var strbuff strings.Builder
		resource.SaveWithWriter(&strbuff, nil)
		resource = nil
	}
}
