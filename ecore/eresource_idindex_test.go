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
	"fmt"
	"reflect"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type idModel struct {
	pkg      EPackage
	item     EClass
	id       EAttribute
	resource EResource
}

func newIDModel(n int) *idModel {
	f := GetFactory()
	m := &idModel{pkg: f.CreateEPackage()}
	m.pkg.SetName("ids")
	m.pkg.SetNsPrefix("ids")
	m.pkg.SetNsURI("urn:ids")
	m.item = f.CreateEClass()
	m.item.SetName("Item")
	m.id = f.CreateEAttribute()
	m.id.SetName("id")
	m.id.SetEType(GetPackage().GetEString())
	m.id.SetID(true)
	m.item.GetEStructuralFeatures().Add(m.id)
	m.pkg.GetEClassifiers().Add(m.item)
	m.resource = NewEResourceImpl()
	for i := range n {
		m.add(fmt.Sprintf("id%d", i))
	}
	return m
}

func (m *idModel) add(id string) EObject {
	o := m.pkg.GetEFactoryInstance().Create(m.item)
	o.ESet(m.id, id)
	m.resource.GetContents().Add(o)
	return o
}

// Lookups follow ids that change, and objects attached or detached after the first lookup.
func TestEResourceGetObjectByIDFollowsTheModel(t *testing.T) {
	m := newIDModel(3)
	first := m.resource.GetContents().Get(0).(EObject)
	require.Same(t, first, m.resource.GetEObject("id0"))

	first.ESet(m.id, "renamed")
	assert.Same(t, first, m.resource.GetEObject("renamed"))
	assert.Nil(t, m.resource.GetEObject("id0"))

	added := m.add("added")
	assert.Same(t, added, m.resource.GetEObject("added"))
	m.resource.GetContents().Remove(added)
	assert.Nil(t, m.resource.GetEObject("added"))
	assert.Nil(t, m.resource.GetEObject("unknown"))
}

// Concurrent first lookups build one index.
func TestEResourceGetObjectByIDConcurrent(t *testing.T) {
	m := newIDModel(100)
	var group sync.WaitGroup
	for i := range 8 {
		group.Add(1)
		go func() {
			defer group.Done()
			id := fmt.Sprintf("id%d", i*10)
			assert.Equal(t, id, m.resource.GetEObject(id).EGet(m.id))
		}()
	}
	group.Wait()
}

func BenchmarkEResourceGetObjectByID(b *testing.B) {
	m := newIDModel(5000)
	for b.Loop() {
		for i := range 5000 {
			m.resource.GetEObject(fmt.Sprintf("id%d", i))
		}
	}
}

// A tree whose classes declare no id attribute is scanned once: a document referring by name
// misses on every lookup, and rescanning for each made resolving it quadratic. An object with an
// id attached later is still found, and so is its new id after a rename.
func TestEResourceGetObjectByIDScansAnIDlessTreeOnce(t *testing.T) {
	m := newIDModel(0)
	f := GetFactory()
	named := f.CreateEClass()
	named.SetName("Named")
	name := f.CreateEAttribute()
	name.SetName("name")
	name.SetEType(GetPackage().GetEString())
	named.GetEStructuralFeatures().Add(name)
	m.pkg.GetEClassifiers().Add(named)
	for i := range 3 {
		o := m.pkg.GetEFactoryInstance().Create(named)
		o.ESet(name, fmt.Sprintf("n%d", i))
		m.resource.GetContents().Add(o)
	}
	resource := m.resource.(*EResourceImpl)

	assert.Nil(t, m.resource.GetEObject("n1"))
	assert.True(t, resource.idless, "a scan meeting no id attribute marks the tree")
	scanned := reflect.ValueOf(resource.idIndex).Pointer()
	assert.Nil(t, m.resource.GetEObject("n2"))
	assert.Equal(t, scanned, reflect.ValueOf(resource.idIndex).Pointer(), "a miss did not rescan")

	late := m.add("late")
	assert.False(t, resource.idless, "attaching an object with an id attribute clears the mark")
	assert.Same(t, late, m.resource.GetEObject("late"))
	late.ESet(m.id, "renamed")
	assert.Same(t, late, m.resource.GetEObject("renamed"))
}

func BenchmarkEResourceGetObjectByIDMissOnIDlessTree(b *testing.B) {
	f := GetFactory()
	pkg := f.CreateEPackage()
	pkg.SetName("names")
	pkg.SetNsURI("urn:names")
	named := f.CreateEClass()
	named.SetName("Named")
	pkg.GetEClassifiers().Add(named)
	resource := NewEResourceImpl()
	for range 5000 {
		resource.GetContents().Add(pkg.GetEFactoryInstance().Create(named))
	}
	for b.Loop() {
		for i := range 1000 {
			resource.GetEObject(fmt.Sprintf("n%d", i))
		}
	}
}
