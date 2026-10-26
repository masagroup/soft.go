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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// the EIsSet and EUnset of the dynamic features, as EMF, for the in-memory objects (which keep
// NIL) and a store-backed object without store

type unsetTestObject interface {
	EObjectInternal
	SetEClass(EClass)
}

type unsetTestKind struct {
	name string
	new  func() unsetTestObject
	// keepsNil is true when the object tells a feature set to nil from an unset one
	keepsNil bool
}

var unsetTestKinds = []unsetTestKind{
	{"Dynamic", func() unsetTestObject { return NewDynamicEObjectImpl() }, true},
	{"Reflective", func() unsetTestObject { return NewReflectiveEObjectImpl() }, true},
	{"EStore", func() unsetTestObject { return NewEStoreEObjectImpl(true) }, false},
}

// unsetTestModel is a class C, with a single-valued attribute of each kind, a many-valued
// unsettable attribute, an unsettable reference, a containment and a bidirectional reference (to
// the friendOf of a class Child)
type unsetTestModel struct {
	c, child                                     EClass
	attribute, unsettableAttribute, intAttribute EAttribute
	// a date has no default value: it can be set to nil
	unsettableDate                EAttribute
	unsettableMany                EAttribute
	unsettableReference, contents EReference
	friends, friendOf             EReference
}

func newUnsetTestModel() *unsetTestModel {
	factory := GetFactory()
	m := &unsetTestModel{}
	m.c = factory.CreateEClass()
	m.c.SetName("C")
	m.child = factory.CreateEClass()
	m.child.SetName("Child")

	newAttribute := func(name string, eType EDataType) EAttribute {
		a := factory.CreateEAttribute()
		a.SetName(name)
		a.SetEType(eType)
		return a
	}
	m.attribute = newAttribute("attribute", GetPackage().GetEString())
	m.attribute.SetDefaultValueLiteral("default")
	m.unsettableAttribute = newAttribute("unsettableAttribute", GetPackage().GetEString())
	m.unsettableAttribute.SetDefaultValueLiteral("default")
	m.unsettableAttribute.SetUnsettable(true)
	m.intAttribute = newAttribute("intAttribute", GetPackage().GetEInt())
	m.unsettableDate = newAttribute("unsettableDate", GetPackage().GetEDate())
	m.unsettableDate.SetUnsettable(true)
	m.unsettableMany = newAttribute("unsettableMany", GetPackage().GetEString())
	m.unsettableMany.SetUpperBound(UNBOUNDED_MULTIPLICITY)
	m.unsettableMany.SetUnsettable(true)

	newReference := func(name string, eType EClass, many bool) EReference {
		r := factory.CreateEReference()
		r.SetName(name)
		r.SetEType(eType)
		if many {
			r.SetUpperBound(UNBOUNDED_MULTIPLICITY)
		}
		return r
	}
	m.unsettableReference = newReference("unsettableReference", m.child, false)
	m.unsettableReference.SetUnsettable(true)
	m.contents = newReference("contents", m.child, true)
	m.contents.SetContainment(true)
	m.friends = newReference("friends", m.child, true)
	m.friendOf = newReference("friendOf", m.c, true)
	m.friends.SetEOpposite(m.friendOf)
	m.friendOf.SetEOpposite(m.friends)

	m.c.GetEStructuralFeatures().AddAll(NewImmutableEList([]any{
		m.attribute, m.unsettableAttribute, m.intAttribute, m.unsettableDate, m.unsettableMany,
		m.unsettableReference, m.contents, m.friends,
	}))
	m.child.GetEStructuralFeatures().Add(m.friendOf)
	return m
}

// unsetTestAdapter records the notifications of its target
type unsetTestAdapter struct {
	AbstractEAdapter
	notifications []ENotification
}

func (a *unsetTestAdapter) NotifyChanged(notification ENotification) {
	a.notifications = append(a.notifications, notification)
}

func runUnsetTest(t *testing.T, test func(t *testing.T, kind unsetTestKind, m *unsetTestModel, newC, newChild func() unsetTestObject)) {
	for _, kind := range unsetTestKinds {
		t.Run(kind.name, func(t *testing.T) {
			m := newUnsetTestModel()
			newObject := func(class EClass) func() unsetTestObject {
				return func() unsetTestObject {
					o := kind.new()
					o.SetEClass(class)
					return o
				}
			}
			test(t, kind, m, newObject(m.c), newObject(m.child))
		})
	}
}

func TestEObjectUnset_ManyContainment(t *testing.T) {
	runUnsetTest(t, func(t *testing.T, kind unsetTestKind, m *unsetTestModel, newC, newChild func() unsetTestObject) {
		o := newC()
		child := newChild()
		list := o.EGet(m.contents).(EList)
		list.Add(child)
		assert.Equal(t, o, child.EContainer())
		assert.True(t, o.EIsSet(m.contents))

		o.EUnset(m.contents)
		assert.False(t, o.EIsSet(m.contents))
		assert.Nil(t, child.EContainer())
		// the list is kept, cleared
		assert.Equal(t, 0, list.Size())
		assert.Same(t, list, o.EGet(m.contents))
	})
}

func TestEObjectUnset_ManyBidirectional(t *testing.T) {
	runUnsetTest(t, func(t *testing.T, kind unsetTestKind, m *unsetTestModel, newC, newChild func() unsetTestObject) {
		o := newC()
		child := newChild()
		o.EGet(m.friends).(EList).Add(child)
		assert.Equal(t, 1, child.EGet(m.friendOf).(EList).Size())

		o.EUnset(m.friends)
		assert.False(t, o.EIsSet(m.friends))
		assert.Equal(t, 0, child.EGet(m.friendOf).(EList).Size())
	})
}

func TestEObjectUnset_ManyUnsettable(t *testing.T) {
	runUnsetTest(t, func(t *testing.T, kind unsetTestKind, m *unsetTestModel, newC, newChild func() unsetTestObject) {
		o := newC()
		// read: not set
		list := o.EGet(m.unsettableMany).(EList)
		assert.False(t, o.EIsSet(m.unsettableMany))

		// changed, even to empty: set
		list.Add("a")
		list.Clear()
		assert.True(t, o.EIsSet(m.unsettableMany))

		// unset: not set, the list notifies UNSET
		adapter := &unsetTestAdapter{}
		o.EAdapters().Add(adapter)
		o.EUnset(m.unsettableMany)
		assert.False(t, o.EIsSet(m.unsettableMany))
		if assert.NotEmpty(t, adapter.notifications) {
			last := adapter.notifications[len(adapter.notifications)-1]
			assert.Equal(t, UNSET, last.GetEventType())
			assert.Equal(t, m.unsettableMany, last.GetFeature())
			assert.Equal(t, true, last.GetOldValue())
			assert.Equal(t, false, last.GetNewValue())
		}

		// set again by a change
		list.Add("b")
		assert.True(t, o.EIsSet(m.unsettableMany))
	})
}

func TestEObjectUnset_SetNil(t *testing.T) {
	runUnsetTest(t, func(t *testing.T, kind unsetTestKind, m *unsetTestModel, newC, newChild func() unsetTestObject) {
		// an unsettable attribute whose type has no default value set to nil is set when the object
		// keeps NIL, else unset
		o := newC()
		o.ESet(m.unsettableDate, nil)
		assert.Equal(t, kind.keepsNil, o.EIsSet(m.unsettableDate))
		assert.Nil(t, o.EGet(m.unsettableDate))
		o.EUnset(m.unsettableDate)
		assert.False(t, o.EIsSet(m.unsettableDate))

		// a string has a default value, "": a string attribute set to nil is unset, unsettable or not
		for _, attribute := range []EAttribute{m.unsettableAttribute, m.attribute} {
			o := newC()
			o.ESet(attribute, "a")
			o.ESet(attribute, nil)
			assert.False(t, o.EIsSet(attribute), attribute.GetName())
			assert.Equal(t, "default", o.EGet(attribute), attribute.GetName())
		}

		// an unsettable reference set to nil
		o = newC()
		o.ESet(m.unsettableReference, nil)
		assert.Equal(t, kind.keepsNil, o.EIsSet(m.unsettableReference))
		assert.Nil(t, o.EGet(m.unsettableReference))
		child := newChild()
		o.ESet(m.unsettableReference, child)
		assert.Equal(t, child, o.EGet(m.unsettableReference))
		o.ESet(m.unsettableReference, nil)
		assert.Equal(t, kind.keepsNil, o.EIsSet(m.unsettableReference))
		o.EUnset(m.unsettableReference)
		assert.False(t, o.EIsSet(m.unsettableReference))

		// an attribute whose type has a default value set to nil is unset
		o = newC()
		o.ESet(m.intAttribute, 2)
		o.ESet(m.intAttribute, nil)
		assert.False(t, o.EIsSet(m.intAttribute))
		assert.Equal(t, 0, o.EGet(m.intAttribute))
	})
}

func TestEObjectUnset_Notifications(t *testing.T) {
	runUnsetTest(t, func(t *testing.T, kind unsetTestKind, m *unsetTestModel, newC, newChild func() unsetTestObject) {
		o := newC()
		adapter := &unsetTestAdapter{}
		o.EAdapters().Add(adapter)
		last := func() ENotification {
			return adapter.notifications[len(adapter.notifications)-1]
		}

		// an attribute is unset to its default value: SET, UNSET when unsettable
		o.ESet(m.attribute, "a")
		o.EUnset(m.attribute)
		assert.Equal(t, SET, last().GetEventType())
		assert.Equal(t, "a", last().GetOldValue())
		assert.Equal(t, "default", last().GetNewValue())

		o.ESet(m.unsettableAttribute, "a")
		o.EUnset(m.unsettableAttribute)
		assert.Equal(t, UNSET, last().GetEventType())
		assert.Equal(t, "a", last().GetOldValue())
		assert.Equal(t, "default", last().GetNewValue())

		// the old value of an unset attribute is its default value
		o.ESet(m.attribute, "b")
		assert.Equal(t, SET, last().GetEventType())
		assert.Equal(t, "default", last().GetOldValue())
		assert.Equal(t, "b", last().GetNewValue())

		// an attribute whose type has a default value set to nil notifies the default value
		o.ESet(m.intAttribute, 2)
		o.ESet(m.intAttribute, nil)
		assert.Equal(t, 2, last().GetOldValue())
		assert.Equal(t, 0, last().GetNewValue())

		// a reference unset without value still notifies
		count := len(adapter.notifications)
		o.EUnset(m.unsettableReference)
		assert.Equal(t, count+1, len(adapter.notifications))
		assert.Equal(t, UNSET, last().GetEventType())
	})
}

func TestBasicENotifyingList_Unsettable(t *testing.T) {
	l := NewBasicENotifyingList()
	assert.False(t, l.IsSet())
	l.Add(1)
	assert.True(t, l.IsSet())
	l.Unset()
	assert.False(t, l.IsSet())
	assert.True(t, l.Empty())

	l.unsettable = true
	assert.False(t, l.IsSet())
	l.Clear()
	assert.True(t, l.IsSet())
	l.Unset()
	assert.False(t, l.IsSet())
}

func TestEStoreList_Unsettable(t *testing.T) {
	m := newUnsetTestModel()
	owner := NewDynamicEObjectImpl()
	owner.SetEClass(m.c)

	// without store: the list records its set state
	l := NewEStoreList(owner, m.unsettableMany, nil)
	assert.False(t, l.IsSet())
	l.Clear()
	assert.True(t, l.IsSet())
	l.Unset()
	assert.False(t, l.IsSet())

	// with store: the store tells, and unset unsets the feature in the store
	mockStore := NewMockEStore(t)
	mockStore.EXPECT().Size(owner, m.unsettableMany).Return(0).Once()
	l = NewEStoreList(owner, m.unsettableMany, mockStore)
	mockStore.EXPECT().IsSet(owner, m.unsettableMany).Return(true).Once()
	assert.True(t, l.IsSet())
	mockStore.EXPECT().ToArray(owner, m.unsettableMany).Return([]any{}).Once()
	mockStore.EXPECT().Clear(owner, m.unsettableMany).Once()
	mockStore.EXPECT().IsSet(owner, m.unsettableMany).Return(true).Once()
	mockStore.EXPECT().UnSet(owner, m.unsettableMany).Once()
	l.Unset()
	mock.AssertExpectationsForObjects(t, mockStore)
}
