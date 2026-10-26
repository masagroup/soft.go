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
	"path/filepath"
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ecorePackage is an .ecore document of a package named name with the namespace urn:<name>, its
// root element carrying the attributes attrs, holding body.
func ecorePackage(name, attrs, body string) string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<ecore:EPackage xmi:version="2.0" xmlns:xmi="http://www.omg.org/XMI" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
    xmlns:ecore="http://www.eclipse.org/emf/2002/Ecore" ` + attrs + `
    name="` + name + `" nsURI="urn:` + name + `" nsPrefix="` + name + `">
` + body + `
</ecore:EPackage>
`
}

// writeDocuments writes the documents by file name in a temporary directory, which it returns.
func writeDocuments(t *testing.T, documents map[string]string) string {
	dir := t.TempDir()
	for name, document := range documents {
		require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(document), 0o600))
	}
	return dir
}

// loadDocument loads the document name of dir in a new resource set.
func loadDocument(dir, name string) ecore.EResource {
	return NewProcessor().Load(ecore.CreateFileURI(filepath.Join(dir, name)))
}

// annotationContent is the content i of the annotation of the package of resource.
func annotationContent(t *testing.T, resource ecore.EResource, i int) ecore.EObject {
	require.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	ePackage := resource.GetContents().Get(0).(ecore.EPackage)
	annotation := ePackage.GetEAnnotation("meta")
	require.NotNil(t, annotation)
	return annotation.GetContents().Get(i).(ecore.EObject)
}

const noteClass = `  <eClassifiers xsi:type="ecore:EClass" name="Note">
    <eStructuralFeatures xsi:type="ecore:EAttribute" name="text" eType="ecore:EDataType http://www.eclipse.org/emf/2002/Ecore#//EString"/>
  </eClassifiers>`

// A package named by xsi:schemaLocation loads from its location, relative to the document, when the
// registry does not know its namespace, and is registered in the registry of the resource set.
func TestXMIDecoderSchemaLocation(t *testing.T) {
	dir := writeDocuments(t, map[string]string{
		"notes.ecore": ecorePackage("notes", "", noteClass),
		"document.ecore": ecorePackage("document", `xmlns:notes="urn:notes" xsi:schemaLocation="urn:notes notes.ecore"`,
			`  <eAnnotations source="meta"><contents xsi:type="notes:Note" text="a note"/></eAnnotations>`),
	})
	resource := loadDocument(dir, "document.ecore")
	note := annotationContent(t, resource, 0)
	assert.Equal(t, "Note", note.EClass().GetName())
	assert.Equal(t, "a note", note.EGet(note.EClass().GetEStructuralFeatureFromName("text")))

	resourceSet := resource.GetResourceSet()
	notes := resourceSet.GetPackageRegistry().GetPackage("urn:notes")
	require.NotNil(t, notes)
	assert.Equal(t, note.EClass().GetEPackage(), notes)
	assert.Equal(t, 2, resourceSet.GetResources().Size())
}

// The package of a schema location loads with its own schema locations.
func TestXMIDecoderSchemaLocationChained(t *testing.T) {
	dir := writeDocuments(t, map[string]string{
		"c.ecore": ecorePackage("c", "", noteClass),
		"b.ecore": ecorePackage("b", `xmlns:c="urn:c" xsi:schemaLocation="urn:c c.ecore"`, noteClass+`
  <eAnnotations source="meta"><contents xsi:type="c:Note" text="in b"/></eAnnotations>`),
		"a.ecore": ecorePackage("a", `xmlns:b="urn:b" xsi:schemaLocation="urn:b b.ecore"`,
			`  <eAnnotations source="meta"><contents xsi:type="b:Note" text="in a"/></eAnnotations>`),
	})
	resource := loadDocument(dir, "a.ecore")
	note := annotationContent(t, resource, 0)
	assert.Equal(t, "urn:b", note.EClass().GetEPackage().GetNsURI())

	b := resource.GetResourceSet().GetResource(ecore.CreateFileURI(filepath.Join(dir, "b.ecore")), false)
	require.NotNil(t, b)
	assert.Equal(t, "urn:c", annotationContent(t, b, 0).EClass().GetEPackage().GetNsURI())
}

// Two packages naming each other in xsi:schemaLocation load once each: the package being loaded is
// the one of the cycle.
func TestXMIDecoderSchemaLocationCycle(t *testing.T) {
	dir := writeDocuments(t, map[string]string{
		"a.ecore": ecorePackage("a", `xmlns:b="urn:b" xsi:schemaLocation="urn:b b.ecore"`, noteClass+`
  <eAnnotations source="meta"><contents xsi:type="b:Note" text="in a"/></eAnnotations>`),
		"b.ecore": ecorePackage("b", `xmlns:a="urn:a" xsi:schemaLocation="urn:a a.ecore"`, noteClass+`
  <eAnnotations source="meta"><contents xsi:type="a:Note" text="in b"/></eAnnotations>`),
	})
	resource := loadDocument(dir, "a.ecore")
	a := resource.GetContents().Get(0).(ecore.EPackage)
	note := annotationContent(t, resource, 0)
	assert.Equal(t, "urn:b", note.EClass().GetEPackage().GetNsURI())

	resourceSet := resource.GetResourceSet()
	assert.Equal(t, 2, resourceSet.GetResources().Size())
	b := resourceSet.GetResource(ecore.CreateFileURI(filepath.Join(dir, "b.ecore")), false)
	require.NotNil(t, b)
	assert.Equal(t, a, annotationContent(t, b, 0).EClass().GetEPackage())
}

// A schema location that does not load leaves the namespace unknown, the diagnostic naming the
// location, and is not kept in the resource set.
func TestXMIDecoderSchemaLocationMissing(t *testing.T) {
	dir := writeDocuments(t, map[string]string{
		"document.ecore": ecorePackage("document", `xmlns:notes="urn:notes" xsi:schemaLocation="urn:notes gone.ecore"`,
			`  <eAnnotations source="meta"><contents xsi:type="notes:Note" text="a note"/></eAnnotations>`),
	})
	resource := loadDocument(dir, "document.ecore")
	require.False(t, resource.GetErrors().Empty())
	message := diagnosticError(resource.GetErrors())
	assert.Contains(t, message, "Package {'notes'='urn:notes'} not found")
	assert.Contains(t, message, "gone.ecore")
	assert.Equal(t, 1, resource.GetResourceSet().GetResources().Size())
}
