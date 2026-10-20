package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type eAllContentsClassModel struct {
	eFactory EFactory
	root     EClass
	item     EClass // abstract super type of book
	book     EClass
	books    EReference // root -> book
	more     EReference // root -> book, the last transition of root
	items    EReference // root -> item
}

func newEAllContentsClassModel() *eAllContentsClassModel {
	f := GetFactory()
	ePackage := f.CreateEPackage()
	m := &eAllContentsClassModel{eFactory: f.CreateEFactory()}
	ePackage.SetEFactoryInstance(m.eFactory)
	newClass := func(name string, abstract bool) EClass {
		eClass := f.CreateEClass()
		eClass.SetName(name)
		eClass.SetAbstract(abstract)
		ePackage.GetEClassifiers().Add(eClass)
		return eClass
	}
	newContainment := func(eClass EClass, name string, eType EClass) EReference {
		eReference := f.CreateEReference()
		eReference.SetName(name)
		eReference.SetEType(eType)
		eReference.SetContainment(true)
		eReference.SetUpperBound(UNBOUNDED_MULTIPLICITY)
		eClass.GetEStructuralFeatures().Add(eReference)
		return eReference
	}
	m.root, m.item, m.book = newClass("Root", false), newClass("Item", true), newClass("Book", false)
	m.book.GetESuperTypes().Add(m.item)
	m.books = newContainment(m.root, "books", m.book)
	m.more = newContainment(m.root, "more", m.book)
	m.items = newContainment(m.root, "items", m.item)
	return m
}

func countWithClass(eObject EObject, eClass EClass) int {
	n := 0
	for it := EAllContentsWithClass(eObject, eClass); it.HasNext(); it.Next() {
		n++
	}
	return n
}

// an unset containment does not end the exploration of the object: its other containments are tried
func TestEAllContentsWithClass_UnsetTransition(t *testing.T) {
	m := newEAllContentsClassModel()
	root := m.eFactory.Create(m.root)
	root.EGet(m.books).(EList).Add(m.eFactory.Create(m.book))
	assert.Equal(t, 1, countWithClass(root, m.book))
}

// an instance of the class is found whatever the type of the reference holding it
func TestEAllContentsWithClass_SubClass(t *testing.T) {
	m := newEAllContentsClassModel()
	root := m.eFactory.Create(m.root)
	root.EGet(m.items).(EList).Add(m.eFactory.Create(m.book))
	root.EGet(m.more).(EList).Add(m.eFactory.Create(m.book))
	assert.Equal(t, 2, countWithClass(root, m.book))
	// the instances of a sub class are instances of the class
	assert.Equal(t, 2, countWithClass(root, m.item))
}
