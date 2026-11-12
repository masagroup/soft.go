// *****************************************************************************
// Copyright(c) 2026 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

// derivedReferencesModel is a library that contains its items, of which a derived containment of
// books is a view, and references its favorite items, of which a derived reference of featured
// items is a view. The dynamic objects store the values of the derived references.
type derivedReferencesModel struct {
	libraryClass EClass
	itemClass    EClass
	items        EReference
	books        EReference
	favorites    EReference
	featured     EReference
	eFactory     EFactory
}

func newDerivedReferencesModel() *derivedReferencesModel {
	factory := GetFactory()
	ePackage := factory.CreateEPackage()
	ePackage.SetName("derived")
	ePackage.SetNsURI("http://derived/references")
	m := &derivedReferencesModel{}
	m.itemClass = factory.CreateEClass()
	m.itemClass.SetName("Item")
	m.libraryClass = factory.CreateEClass()
	m.libraryClass.SetName("Library")
	newReference := func(name string, containment, derived bool) EReference {
		reference := factory.CreateEReference()
		reference.SetName(name)
		reference.SetEType(m.itemClass)
		reference.SetContainment(containment)
		reference.SetDerived(derived)
		reference.SetUpperBound(UNBOUNDED_MULTIPLICITY)
		m.libraryClass.GetEStructuralFeatures().Add(reference)
		return reference
	}
	m.items = newReference("items", true, false)
	m.books = newReference("books", true, true)
	m.favorites = newReference("favorites", false, false)
	m.featured = newReference("featured", false, true)
	ePackage.GetEClassifiers().AddAll(NewImmutableEList([]any{m.itemClass, m.libraryClass}))
	ePackage.SetEFactoryInstance(factory.CreateEFactory())
	m.eFactory = ePackage.GetEFactoryInstance()
	return m
}

func TestEClassExt_FeaturesSubSet_Derived(t *testing.T) {
	// like EMF, the derived references are containment and cross reference features: the compact
	// objects read them
	m := newDerivedReferencesModel()
	assert.Equal(t, []any{m.items, m.books}, m.libraryClass.GetEContainmentFeatures().ToArray())
	assert.Equal(t, []any{m.favorites, m.featured}, m.libraryClass.GetECrossReferenceFeatures().ToArray())
}

func TestBasicEObjectImpl_EContents_Derived(t *testing.T) {
	m := newDerivedReferencesModel()
	library := m.eFactory.Create(m.libraryClass)
	item := m.eFactory.Create(m.itemClass)
	book := m.eFactory.Create(m.itemClass)
	library.EGet(m.items).(EList).Add(item)
	library.EGet(m.books).(EList).Add(book)

	assert.Equal(t, []any{item, book}, library.EContents().ToArray())
	assert.Equal(t, []EObject{item, book}, slices.Collect(AllContents(library)))
	contents := []EObject{}
	for it := library.EAllContents(); it.HasNext(); {
		contents = append(contents, it.Next().(EObject))
	}
	assert.Equal(t, []EObject{item, book}, contents)
}

func TestBasicEObjectImpl_ECrossReferences_Derived(t *testing.T) {
	m := newDerivedReferencesModel()
	library := m.eFactory.Create(m.libraryClass)
	favorite := m.eFactory.Create(m.itemClass)
	featured := m.eFactory.Create(m.itemClass)
	library.EGet(m.favorites).(EList).Add(favorite)
	library.EGet(m.featured).(EList).Add(featured)

	assert.Equal(t, []any{favorite, featured}, library.ECrossReferences().ToArray())

	// the usages of an object, found to delete it, are those of ECrossReferences
	references := map[EStructuralFeature]EObject{}
	for reference, object := range crossReferencesInObject(library, true) {
		references[reference] = object
	}
	assert.Equal(t, map[EStructuralFeature]EObject{m.favorites: favorite, m.featured: featured}, references)
}
