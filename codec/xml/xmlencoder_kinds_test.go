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

type kindsModel struct {
	pkg                            ecore.EPackage
	root, target, note             ecore.EClass
	ref, refs, children, notes     ecore.EReference
	targetName, noteText, category ecore.EAttribute
}

func withKind(element ecore.ENamedElement, kind, name string) {
	annotation := ecore.GetFactory().CreateEAnnotation()
	annotation.SetSource("http:///org/eclipse/emf/ecore/util/ExtendedMetaData")
	annotation.GetDetails().Put("kind", kind)
	if name != "" {
		annotation.GetDetails().Put("name", name)
	}
	element.GetEAnnotations().Add(annotation)
}

func newKindsModel() *kindsModel {
	f := ecore.GetFactory()
	m := &kindsModel{pkg: f.CreateEPackage()}
	m.pkg.SetName("kinds")
	m.pkg.SetNsPrefix("kinds")
	m.pkg.SetNsURI("urn:kinds")
	newClass := func(name string) ecore.EClass {
		c := f.CreateEClass()
		c.SetName(name)
		m.pkg.GetEClassifiers().Add(c)
		return c
	}
	newAttribute := func(owner ecore.EClass, name string) ecore.EAttribute {
		a := f.CreateEAttribute()
		a.SetName(name)
		a.SetEType(ecore.GetPackage().GetEString())
		owner.GetEStructuralFeatures().Add(a)
		return a
	}
	newReference := func(owner ecore.EClass, name string, target ecore.EClass, containment, many bool) ecore.EReference {
		r := f.CreateEReference()
		r.SetName(name)
		r.SetEType(target)
		r.SetContainment(containment)
		if many {
			r.SetUpperBound(-1)
		}
		owner.GetEStructuralFeatures().Add(r)
		return r
	}
	m.root, m.target, m.note = newClass("Root"), newClass("Target"), newClass("Note")
	m.targetName = newAttribute(m.target, "name")
	m.noteText = newAttribute(m.note, "text")
	withKind(m.noteText, "simple", ":0")
	m.category = newAttribute(m.root, "category")
	withKind(m.category, "element", "")
	m.ref = newReference(m.root, "ref", m.target, false, false)
	withKind(m.ref, "attribute", "")
	m.refs = newReference(m.root, "refs", m.target, false, true)
	withKind(m.refs, "attribute", "")
	m.children = newReference(m.root, "children", m.target, true, true)
	m.notes = newReference(m.root, "notes", m.note, true, true)
	return m
}

func (m *kindsModel) roundTrip(t *testing.T, root ecore.EObject) (string, ecore.EObject) {
	t.Helper()
	processor := NewProcessor(ProcessorPackages([]ecore.EPackage{m.pkg}))
	options := map[string]any{OPTION_EXTENDED_META_DATA: NewExtendedMetaData()}
	resource := processor.GetResourceSet().CreateResource(ecore.NewURI("root.xml"))
	resource.GetContents().Add(root)
	encoded := processor.SaveToString(resource, options)
	require.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	reloaded := processor.GetResourceSet().CreateResource(ecore.NewURI("reloaded.xml"))
	reloaded.LoadWithReader(strings.NewReader(encoded), options)
	require.True(t, reloaded.GetErrors().Empty(), diagnosticError(reloaded.GetErrors()))
	require.Equal(t, 1, reloaded.GetContents().Size())
	return encoded, reloaded.GetContents().Get(0).(ecore.EObject)
}

// A reference declared kind="attribute" is written as an attribute and read back as one.
func TestXMLEncoderKindAttributeReferenceToAnotherResource(t *testing.T) {
	m := newKindsModel()
	proxy := func(uri string) ecore.EObject {
		o := m.pkg.GetEFactoryInstance().Create(m.target)
		o.(ecore.EObjectInternal).ESetProxyURI(ecore.NewURI(uri))
		return o
	}
	root := m.pkg.GetEFactoryInstance().Create(m.root)
	root.ESet(m.ref, proxy("other.xml#//first"))
	root.EGet(m.refs).(ecore.EList).AddAll(ecore.NewImmutableEList([]any{proxy("other.xml#//first"), proxy("other.xml#//second")}))

	encoded, reloaded := m.roundTrip(t, root)
	assert.Contains(t, encoded, ` ref="other.xml#//first"`)
	assert.Contains(t, encoded, ` refs="other.xml#//first other.xml#//second"`)
	assert.NotContains(t, encoded, "href=")
	assert.Equal(t, "other.xml#//first", reloaded.EGetResolve(m.ref, false).(ecore.EObjectInternal).EProxyURI().String())
	assert.Equal(t, 2, reloaded.EGetResolve(m.refs, false).(ecore.EList).Size())
}

// A kind="attribute" reference to an object in no resource is saved as an IDREF, not an error.
func TestXMLEncoderKindAttributeReferenceToADetachedObject(t *testing.T) {
	m := newKindsModel()
	root := m.pkg.GetEFactoryInstance().Create(m.root)
	root.ESet(m.ref, m.pkg.GetEFactoryInstance().Create(m.target))
	processor := NewProcessor(ProcessorPackages([]ecore.EPackage{m.pkg}))
	resource := processor.GetResourceSet().CreateResource(ecore.NewURI("root.xml"))
	resource.GetContents().Add(root)
	encoded := processor.SaveToString(resource, map[string]any{OPTION_EXTENDED_META_DATA: NewExtendedMetaData()})
	require.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	assert.Contains(t, encoded, ` ref="#`)
	assert.NotContains(t, encoded, "href=")
}

// A kind="simple" attribute is the element's text content.
func TestXMLEncoderKindSimpleContentIsElementText(t *testing.T) {
	m := newKindsModel()
	root := m.pkg.GetEFactoryInstance().Create(m.root)
	note := m.pkg.GetEFactoryInstance().Create(m.note)
	note.ESet(m.noteText, "en g/m3 & more")
	root.EGet(m.notes).(ecore.EList).Add(note)

	encoded, reloaded := m.roundTrip(t, root)
	assert.Contains(t, encoded, ">en g/m3 &amp; more</notes>")
	assert.Equal(t, "en g/m3 & more", reloaded.EGet(m.notes).(ecore.EList).Get(0).(ecore.EObject).EGet(m.noteText))
}
