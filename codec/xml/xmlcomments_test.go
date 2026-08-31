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
	"bytes"
	stdxml "encoding/xml"
	"io"
	"strings"
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A comment has to survive BOTH halves of a round trip, and the two halves are written apart: the
// decoder decides where a comment WAS, the encoder decides where it GOES. Neither can be tested
// alone - a decoder that mislabels a position and an encoder that mislabels it the same way would
// still put the comment somewhere wrong - so what is pinned here is the pair, on a document that
// carries a comment in every position one can occupy.
const commentedLibrary = `<?xml version="1.0" encoding="UTF-8"?>
<!-- a note before the root -->
<Library xmlns="http:///org/eclipse/emf/examples/library/library.ecore/1.0.0" name="My Library">
  <!-- before the first child -->
  <employees firstName="Ada" lastName="Lovelace"/>
  <!-- between two siblings -->
  <employees firstName="Alan" lastName="Turing"/>
  <branches name="Annex">
    <!-- inside a nested element -->
    <employees firstName="Grace" lastName="Hopper"/>
  </branches>
  <!-- after the last child -->
</Library>`

// documentEvents reduces a document to the sequence a comment's position is defined against:
// element names in document order, and comment bodies where they stand between them. Comparing two
// documents' event sequences answers "did every comment come back among the same siblings, in the
// same order" and nothing else - indentation, attribute order and line endings cannot affect it.
func documentEvents(t *testing.T, data []byte) []string {
	t.Helper()
	var out []string
	dec := stdxml.NewDecoder(bytes.NewReader(data))
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		switch tk := tok.(type) {
		case stdxml.StartElement:
			out = append(out, "<"+tk.Name.Local+">")
		case stdxml.EndElement:
			out = append(out, "</"+tk.Name.Local+">")
		case stdxml.Comment:
			out = append(out, "#"+strings.TrimSpace(string(tk)))
		}
	}
	return out
}

// commentedLibraryPackage is the part of the library model the commented document uses: a Library
// with a name, employees and nested branches, an Employee with a first and a last name.
func commentedLibraryPackage() ecore.EPackage {
	f := ecore.GetFactory()
	eString := ecore.GetPackage().GetEString()
	ePackage := f.CreateEPackage()
	ePackage.SetName("library")
	ePackage.SetNsPrefix("lib")
	ePackage.SetNsURI("http:///org/eclipse/emf/examples/library/library.ecore/1.0.0")
	attribute := func(name string) ecore.EAttribute {
		eAttribute := f.CreateEAttribute()
		eAttribute.SetName(name)
		eAttribute.SetEType(eString)
		return eAttribute
	}
	containment := func(name string, eType ecore.EClass) ecore.EReference {
		eReference := f.CreateEReference()
		eReference.SetName(name)
		eReference.SetEType(eType)
		eReference.SetContainment(true)
		eReference.SetUpperBound(-1)
		return eReference
	}
	employee := f.CreateEClass()
	employee.SetName("Employee")
	employee.GetEStructuralFeatures().AddAll(ecore.NewImmutableEList([]any{attribute("firstName"), attribute("lastName")}))
	library := f.CreateEClass()
	library.SetName("Library")
	library.GetEStructuralFeatures().AddAll(ecore.NewImmutableEList([]any{attribute("name"), containment("employees", employee), containment("branches", library)}))
	ePackage.GetEClassifiers().AddAll(ecore.NewImmutableEList([]any{library, employee}))
	return ePackage
}

func loadCommentedLibrary(t *testing.T) (*Processor, ecore.EResource) {
	t.Helper()
	processor := NewProcessor(ProcessorPackages([]ecore.EPackage{commentedLibraryPackage()}))
	resource := processor.GetResourceSet().CreateResource(ecore.NewURI("library.xml"))
	resource.LoadWithReader(strings.NewReader(commentedLibrary), nil)
	require.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	return processor, resource
}

func TestXMLCommentsSurviveARoundTrip(t *testing.T) {
	processor, resource := loadCommentedLibrary(t)

	var buffer bytes.Buffer
	processor.SaveWithWriter(&buffer, resource, map[string]any{OPTION_KEEP_COMMENTS: true})
	require.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))

	assert.Equal(t, documentEvents(t, []byte(commentedLibrary)), documentEvents(t, buffer.Bytes()))
	assert.True(t, resource.GetWarnings().Empty(), diagnosticError(resource.GetWarnings()))
}

// Without the option the encoder writes what it always wrote, byte for byte, although the decoder
// now collects the comments of every document.
func TestXMLCommentsAreNotWrittenUnlessAsked(t *testing.T) {
	processor, resource := loadCommentedLibrary(t)

	var buffer bytes.Buffer
	processor.SaveWithWriter(&buffer, resource, nil)
	require.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))

	assert.NotContains(t, buffer.String(), "<!--")
	assert.Contains(t, buffer.String(), "Lovelace") // the document itself was still written
}

// Collecting comments must not reach the object graph: they are not objects, not features and not
// text, and a decode that let them leak into any of the three would corrupt the model to preserve a
// note about it.
func TestXMLCommentsDoNotEnterTheGraph(t *testing.T) {
	_, resource := loadCommentedLibrary(t)

	library := resource.GetContents().Get(0).(ecore.EObject)
	employees := library.EGet(library.EClass().GetEStructuralFeatureFromName("employees")).(ecore.EList)
	assert.Equal(t, 2, employees.Size())
	assert.Equal(t, "My Library", library.EGet(library.EClass().GetEStructuralFeatureFromName("name")))
}

// This is the mechanism's limit, stated as a test rather than left to be discovered.
//
// A comment is addressed by the position it held - inside which element, after how many of its
// children - not by the object it was written about. So an edit does not carry a comment with its
// meaning: deleting <branches> takes away both the comment INSIDE it and the position of the
// comment that stood AFTER it, because Library now has two children where the trailing comment was
// pinned after three.
//
// A load and a write-back never move a position, which is why the round trip above is exact. An
// editor built on this needs re-anchoring, and until it has one, a comment that cannot be placed is
// said out loud: it becomes a warning rather than vanishing, and the rest of the document is still
// written, because refusing the save would be the worse failure.
func TestXMLCommentsAreAnchoredToAPositionNotAnObject(t *testing.T) {
	processor, resource := loadCommentedLibrary(t)

	library := resource.GetContents().Get(0).(ecore.EObject)
	branches := library.EGet(library.EClass().GetEStructuralFeatureFromName("branches")).(ecore.EList)
	branches.Clear()

	var buffer bytes.Buffer
	processor.SaveWithWriter(&buffer, resource, map[string]any{OPTION_KEEP_COMMENTS: true})

	require.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	require.Equal(t, 1, resource.GetWarnings().Size())
	assert.Contains(t, resource.GetWarnings().Get(0).(ecore.EDiagnostic).GetMessage(), "2 of 5 XML comments")
	assert.NotContains(t, buffer.String(), "inside a nested element")
	assert.NotContains(t, buffer.String(), "after the last child")
	// The comments whose position the edit did not disturb are written as they always were.
	assert.Contains(t, buffer.String(), "a note before the root")
	assert.Contains(t, buffer.String(), "before the first child")
	assert.Contains(t, buffer.String(), "between two siblings")
}
