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
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompactEObjectConstructor(t *testing.T) {
	o := &CompactEObjectImpl{}
	o.Initialize()
	assert.True(t, o.EDeliver())
}

func TestCompactEObject_EDeliver(t *testing.T) {
	o := &CompactEObjectImpl{}
	assert.False(t, o.EDeliver())

	o.ESetDeliver(true)
	assert.True(t, o.EDeliver())

	o.ESetDeliver(false)
	assert.False(t, o.EDeliver())
}

func TestCompactEObject_NoFields(t *testing.T) {
	o := &CompactEObjectImpl{}
	assert.False(t, o.hasField(container_flag))
	assert.Nil(t, o.getField(container_flag))
}

func TestCompactEObject_OneField(t *testing.T) {
	c := NewMockEObject(t)
	o := &CompactEObjectImpl{}
	o.setField(container_flag, c)
	assert.True(t, o.hasField(container_flag))
	assert.Equal(t, c, o.getField(container_flag))

	c2 := NewMockEObject(t)
	o.setField(container_flag, c2)
	assert.True(t, o.hasField(container_flag))
	assert.Equal(t, c2, o.getField(container_flag))
}

func TestCompactEObject_TwoFields_Begin(t *testing.T) {
	c := NewMockEObject(t)
	r := NewMockEResource(t)
	o := &CompactEObjectImpl{}
	o.setField(resource_flag, r)
	o.setField(container_flag, c)
	assert.True(t, o.hasField(container_flag))
	assert.True(t, o.hasField(resource_flag))
	assert.False(t, o.hasField(properties_flag))
	assert.Equal(t, c, o.getField(container_flag))
	assert.Equal(t, r, o.getField(resource_flag))

	c2 := NewMockEObject(t)
	o.setField(container_flag, c2)
	assert.True(t, o.hasField(container_flag))
	assert.Equal(t, c2, o.getField(container_flag))
}

func TestCompactEObject_TwoFields_End(t *testing.T) {
	c := NewMockEObject(t)
	r := NewMockEResource(t)
	o := &CompactEObjectImpl{}
	o.setField(container_flag, c)
	o.setField(resource_flag, r)
	assert.True(t, o.hasField(container_flag))
	assert.True(t, o.hasField(resource_flag))
	assert.False(t, o.hasField(properties_flag))
	assert.Equal(t, c, o.getField(container_flag))
	assert.Equal(t, r, o.getField(resource_flag))
}

func TestCompactEObject_ManyFields(t *testing.T) {
	c := NewMockEObject(t)
	r := NewMockEResource(t)
	p := []any{}
	o := &CompactEObjectImpl{}
	o.setField(container_flag, c)
	o.setField(properties_flag, p)
	o.setField(resource_flag, r)
	assert.True(t, o.hasField(container_flag))
	assert.True(t, o.hasField(resource_flag))
	assert.True(t, o.hasField(properties_flag))
	assert.Equal(t, c, o.getField(container_flag))
	assert.Equal(t, r, o.getField(resource_flag))
	assert.Equal(t, p, o.getField(properties_flag))
}

func TestCompactEObject_RemoveNoField(t *testing.T) {
	o := &CompactEObjectImpl{}
	o.setField(container_flag, nil)
}

func TestCompactEObject_RemoveOneField(t *testing.T) {
	c := NewMockEObject(t)
	o := &CompactEObjectImpl{}
	o.setField(container_flag, c)
	assert.True(t, o.hasField(container_flag))
	o.setField(container_flag, nil)
	assert.False(t, o.hasField(container_flag))
}

func TestCompactEObject_RemoveManyFields(t *testing.T) {
	c := NewMockEObject(t)
	r := NewMockEResource(t)
	p := []any{}
	o := &CompactEObjectImpl{}
	o.setField(container_flag, c)
	o.setField(properties_flag, p)
	o.setField(resource_flag, r)
	o.setField(container_flag, nil)
	o.setField(resource_flag, nil)
	assert.False(t, o.hasField(container_flag))
	assert.False(t, o.hasField(resource_flag))
	assert.True(t, o.hasField(properties_flag))
	assert.Equal(t, p, o.getField(properties_flag))

	o.setField(resource_flag, r)
	o.setField(properties_flag, nil)
	assert.False(t, o.hasField(properties_flag))
	assert.Nil(t, o.getField(properties_flag))
}

func TestCompactEObject_EClass(t *testing.T) {
	c := NewMockEClass(t)
	o := &CompactEObjectImpl{}
	o.SetInterfaces(o)
	assert.NotNil(t, o.EClass())
	o.SetEClass(c)
	assert.Equal(t, c, o.EClass())
}

func TestCompactEObject_EAdapters(t *testing.T) {
	o := &CompactEObjectImpl{}
	o.SetInterfaces(o)
	assert.False(t, o.EBasicHasAdapters())
	assert.Nil(t, o.EBasicAdapters())

	adapters := o.EAdapters()
	assert.NotNil(t, adapters)
	assert.Equal(t, adapters, o.EBasicAdapters())
	assert.True(t, o.EBasicHasAdapters())
}

func TestCompactEObject_EProxy(t *testing.T) {
	o := &CompactEObjectImpl{}
	o.SetInterfaces(o)
	assert.False(t, o.EIsProxy())
	assert.Nil(t, o.EProxyURI())

	mockURI := NewURI("test:///file.t")
	o.ESetProxyURI(mockURI)
	assert.True(t, o.EIsProxy())
	assert.Equal(t, mockURI, o.EProxyURI())
}

func TestCompactEObject_EResource(t *testing.T) {
	r := NewMockEResource(t)
	o := &CompactEObjectImpl{}
	o.SetInterfaces(o)
	assert.Nil(t, o.EInternalResource())

	o.ESetInternalResource(r)
	assert.Equal(t, r, o.EInternalResource())
}

func TestCompactEObject_EContainer(t *testing.T) {
	c := NewMockEObject(t)
	o := &CompactEObjectImpl{}
	o.SetInterfaces(o)
	o.Initialize()
	assert.Nil(t, o.EInternalContainer())
	assert.Equal(t, -1, o.EInternalContainerFeatureID())

	o.ESetInternalContainer(c, 20)
	assert.Equal(t, c, o.EInternalContainer())
	assert.Equal(t, 20, o.EInternalContainerFeatureID())

	o.ESetInternalContainer(nil, -1)
	assert.Equal(t, nil, o.EInternalContainer())
	assert.Equal(t, -1, o.EInternalContainerFeatureID())
}

func TestCompactEObject_EContents(t *testing.T) {
	o := &CompactEObjectImpl{}
	o.SetInterfaces(o)
	o.Initialize()
	assert.NotNil(t, o.EContents())
}

func TestCompactEObject_ECrossReferences(t *testing.T) {
	o := &CompactEObjectImpl{}
	o.SetInterfaces(o)
	o.Initialize()
	assert.NotNil(t, o.ECrossReferences())
}

// A lazily built list and the proxy URI share one packed slot: concurrent first readers must
// all get the same list, and never see the slot half promoted.
func TestCompactEObject_ConcurrentFirstReads(t *testing.T) {
	for range 200 {
		o := &CompactEObjectImpl{}
		o.Initialize()
		uri := NewURI("other.xml#//target")
		o.ESetProxyURI(uri)
		adapters := make([]EList, 8)
		var group sync.WaitGroup
		for i := range adapters {
			group.Add(2)
			go func() {
				defer group.Done()
				adapters[i] = o.EAdapters()
			}()
			go func() {
				defer group.Done()
				assert.Same(t, uri, o.EProxyURI())
			}()
		}
		group.Wait()
		for _, list := range adapters {
			require.Same(t, adapters[0], list)
		}
	}
// EClassExt is a compact object : eAnnotations is a containment inherited from EModelElement,
// EContents and EAllContents must visit it like BasicEObjectImpl does
func TestCompactEObject_EContents_Inherited(t *testing.T) {
	eClass := GetFactory().CreateEClass()
	eAnnotation := GetFactory().CreateEAnnotation()
	eClass.GetEAnnotations().Add(eAnnotation)
	assert.Equal(t, []any{eAnnotation}, eClass.EContents().ToArray())

	ePackage := GetFactory().CreateEPackage()
	ePackage.GetEClassifiers().Add(eClass)
	visited := []any{}
	for it := ePackage.EAllContents(); it.HasNext(); {
		visited = append(visited, it.Next())
	}
	assert.Equal(t, []any{ePackage.GetEFactoryInstance(), eClass, eAnnotation}, visited)
}

// EReferenceExt is a compact object : eType is a cross reference inherited from ETypedElement,
// ECrossReferences must contain it so that ResolveAll resolves it
func TestCompactEObject_ECrossReferences_Inherited(t *testing.T) {
	eFactory := GetFactory()
	eResourceSet := NewEResourceSetImpl()

	// the referenced class in its own resource
	eTargetPackage := eFactory.CreateEPackage()
	eTargetPackage.SetName("target")
	eTargetPackage.SetNsURI("http://target")
	eTargetClass := eFactory.CreateEClass()
	eTargetClass.SetName("Target")
	eTargetPackage.GetEClassifiers().Add(eTargetClass)
	eTargetResource := eResourceSet.CreateResource(NewURI("memory://target.ecore"))
	eTargetResource.GetContents().Add(eTargetPackage)

	// the referencing class with a proxy as eType in another resource
	ePackage := eFactory.CreateEPackage()
	ePackage.SetName("source")
	ePackage.SetNsURI("http://source")
	eClass := eFactory.CreateEClass()
	eClass.SetName("Source")
	ePackage.GetEClassifiers().Add(eClass)
	eReference := eFactory.CreateEReference()
	eReference.SetName("target")
	eProxy := eFactory.CreateEClass()
	eProxy.(EObjectInternal).ESetProxyURI(NewURI("memory://target.ecore#//Target"))
	eReference.SetEType(eProxy)
	eClass.GetEStructuralFeatures().Add(eReference)
	eResource := eResourceSet.CreateResource(NewURI("memory://source.ecore"))
	eResource.GetContents().Add(ePackage)

	// the proxy is not resolved until the cross references are visited
	eTypeFeature := GetPackage().GetETypedElement_EType()
	assert.Equal(t, eProxy, eReference.EGetResolve(eTypeFeature, false))

	ResolveAll(ePackage)
	assert.Equal(t, eTargetClass, eReference.EGetResolve(eTypeFeature, false))
	assert.True(t, eReference.ECrossReferences().Contains(eTargetClass))
}
