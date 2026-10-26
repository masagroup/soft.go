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
	"strings"
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The contents of an annotation, a feature typed EObject, are saved with their type, as EMF saves
// them, and load back as the objects they were.
func TestXMIEncoderAnnotationContentsType(t *testing.T) {
	dir := writeDocuments(t, map[string]string{
		"notes.ecore": ecorePackage("notes", "", noteClass),
		"document.ecore": ecorePackage("document", `xmlns:notes="urn:notes" xsi:schemaLocation="urn:notes notes.ecore"`,
			`  <eAnnotations source="meta">
    <contents xsi:type="notes:Note" text="a note"/>
    <contents xsi:type="ecore:EAnnotation" source="nested"/>
  </eAnnotations>`),
	})
	resource := loadDocument(dir, "document.ecore")
	require.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))

	saved := NewProcessor().SaveToString(resource, nil)
	assert.Contains(t, saved, `<contents xsi:type="notes:Note" text="a note"/>`)
	assert.Contains(t, saved, `<contents xsi:type="ecore:EAnnotation" source="nested"/>`)

	// the saved document loads with the notes package registered, without its schema location
	resourceSet := ecore.NewEResourceSetImpl()
	notes := resource.GetResourceSet().GetPackageRegistry().GetPackage("urn:notes")
	require.NotNil(t, notes)
	resourceSet.GetPackageRegistry().RegisterPackage(notes)
	reloaded := resourceSet.CreateResource(ecore.NewURI("document.ecore"))
	reloaded.LoadWithReader(strings.NewReader(saved), nil)
	note := annotationContent(t, reloaded, 0)
	assert.Equal(t, notes, note.EClass().GetEPackage())
	assert.Equal(t, "a note", note.EGet(note.EClass().GetEStructuralFeatureFromName("text")))
	nested, _ := annotationContent(t, reloaded, 1).(ecore.EAnnotation)
	require.NotNil(t, nested)
	assert.Equal(t, "nested", nested.GetSource())
}
