package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// C extends A and B: the features of B have an ID in C (after the features of A) other than their
// ID in B. They are accessed, contained, and kept opposite with their ID in C.
type multipleInheritanceModel struct {
	ePackage EPackage
	eFactory EFactory
	a, b, c  EClass
	item     EClass
	fa, fb   EAttribute
	children EReference // B, many containment without opposite
	links    EReference // B, many, opposite of Item.owners
	owners   EReference // Item, many, opposite of B.links
	cLinksID int
	bLinksID int
}

func newMultipleInheritanceModel(t *testing.T) *multipleInheritanceModel {
	f := GetFactory()
	m := &multipleInheritanceModel{ePackage: f.CreateEPackage(), eFactory: f.CreateEFactory()}
	m.ePackage.SetEFactoryInstance(m.eFactory)
	newClass := func(name string) EClass {
		eClass := f.CreateEClass()
		eClass.SetName(name)
		m.ePackage.GetEClassifiers().Add(eClass)
		return eClass
	}
	newAttribute := func(eClass EClass, name string) EAttribute {
		eAttribute := f.CreateEAttribute()
		eAttribute.SetName(name)
		eAttribute.SetEType(GetPackage().GetEString())
		eClass.GetEStructuralFeatures().Add(eAttribute)
		return eAttribute
	}
	newReference := func(eClass EClass, name string, eType EClass) EReference {
		eReference := f.CreateEReference()
		eReference.SetName(name)
		eReference.SetEType(eType)
		eReference.SetUpperBound(UNBOUNDED_MULTIPLICITY)
		eClass.GetEStructuralFeatures().Add(eReference)
		return eReference
	}
	m.a, m.b, m.c, m.item = newClass("A"), newClass("B"), newClass("C"), newClass("Item")
	m.fa = newAttribute(m.a, "fa")
	m.fb = newAttribute(m.b, "fb")
	m.children = newReference(m.b, "children", m.item)
	m.children.SetContainment(true)
	m.links = newReference(m.b, "links", m.item)
	m.owners = newReference(m.item, "owners", m.b)
	m.links.SetEOpposite(m.owners)
	m.owners.SetEOpposite(m.links)
	m.c.GetESuperTypes().Add(m.a)
	m.c.GetESuperTypes().Add(m.b)
	m.bLinksID = m.b.GetFeatureID(m.links)
	m.cLinksID = m.c.GetFeatureID(m.links)
	require.NotEqual(t, m.bLinksID, m.cLinksID)
	return m
}

func TestMultipleInheritance_DynamicObject(t *testing.T) {
	m := newMultipleInheritanceModel(t)
	o := m.eFactory.Create(m.c)

	// containment without opposite: the child is contained by the feature of B
	child := m.eFactory.Create(m.item)
	o.EGet(m.children).(EList).Add(child)
	assert.Equal(t, m.children, child.EContainingFeature())
	assert.Equal(t, o, child.EContainer())

	// opposites, from both sides
	other := m.eFactory.Create(m.item)
	o.EGet(m.links).(EList).Add(other)
	assert.True(t, other.EGet(m.owners).(EList).Contains(o))
	third := m.eFactory.Create(m.item)
	third.EGet(m.owners).(EList).Add(o)
	assert.True(t, o.EGet(m.links).(EList).Contains(third))
	assert.Equal(t, 1, o.EGet(m.children).(EList).Size())
	third.EGet(m.owners).(EList).Remove(o)
	assert.False(t, o.EGet(m.links).(EList).Contains(third))

	// notification
	assert.Equal(t, m.cLinksID, NewNotificationByFeature(o, SET, m.links, nil, nil, NO_INDEX).GetFeatureID())
}

func TestMultipleInheritance_ReflectiveObject(t *testing.T) {
	m := newMultipleInheritanceModel(t)
	o := NewReflectiveEObjectImpl()
	o.SetEClass(m.c)
	assert.Equal(t, m.c.GetFeatureID(m.fb), o.EFeatureID(m.fb))
	o.ESet(m.fb, "b")
	o.ESet(m.fa, "a")
	assert.Equal(t, "b", o.EGet(m.fb))
	assert.Equal(t, "a", o.EGet(m.fa))
}
