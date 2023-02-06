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

func TestDynamicEObjectConstructor(t *testing.T) {
	o := NewDynamicEObjectImpl()
	assert.NotNil(t, o)
	assert.Equal(t, GetPackage().GetEObject(), o.EClass())
}

func TestDynamicEObject_EClass(t *testing.T) {
	o := NewDynamicEObjectImpl()
	mockClass := NewMockEClass(t)
	mockAdapters := NewMockEList(t)
	mockClass.EXPECT().GetFeatureCount().Return(0)
	mockClass.EXPECT().EAdapters().Return(mockAdapters)
	mockAdapters.EXPECT().Add(mock.Anything).Return(true).Once()
	o.SetEClass(mockClass)
	assert.Equal(t, mockClass, o.EClass())
}

func TestDynamicEObject_MockEClass(t *testing.T) {
	o := NewDynamicEObjectImpl()
	c1 := GetFactory().CreateEClass()
	c2 := GetFactory().CreateEClass()
	o.SetEClass(c1)
	assert.Equal(t, c1, o.EClass())

	o.SetEClass(c2)
	assert.Equal(t, c2, o.EClass())
}

func TestDynamicEObject_GetSet(t *testing.T) {
	o := NewDynamicEObjectImpl()
	c := GetFactory().CreateEClass()
	o.SetEClass(c)
	assert.Equal(t, c, o.EClass())

	f := GetFactory().CreateEAttribute()
	c.GetEStructuralFeatures().Add(f)
	assert.Nil(t, o.EGet(f))

	o.ESet(f, 1)
	assert.Equal(t, 1, o.EGet(f))
}

func TestDynamicEObject_Unset(t *testing.T) {
	o := NewDynamicEObjectImpl()
	c := GetFactory().CreateEClass()
	o.SetEClass(c)
	assert.Equal(t, c, o.EClass())

	f := GetFactory().CreateEAttribute()
	c.GetEStructuralFeatures().Add(f)
	assert.Nil(t, o.EGet(f))

	o.ESet(f, 1)
	assert.Equal(t, 1, o.EGet(f))

	o.EUnset(f)
	assert.Nil(t, o.EGet(f))
}

func TestDynamicEObject_Container(t *testing.T) {
	r1 := GetFactory().CreateEReference()
	r1.SetContainment(true)
	r1.SetName("ref")

	r2 := GetFactory().CreateEReference()
	r2.SetName("parent")
	r2.SetEOpposite(r1)

	c1 := GetFactory().CreateEClass()
	c1.GetEStructuralFeatures().Add(r1)

	c2 := GetFactory().CreateEClass()
	c2.GetEStructuralFeatures().Add(r2)

	o1 := NewDynamicEObjectImpl()
	o1.SetEClass(c1)

	o2 := NewDynamicEObjectImpl()
	o2.SetEClass(c2)

	assert.False(t, o2.EIsSet(r2))
	assert.False(t, o1.EIsSet(r1))

	o2.ESet(r2, o1)
	assert.Equal(t, o1, o2.EGet(r2))
	assert.Equal(t, o1, o2.EGetResolve(r2, false))
	assert.Equal(t, o2, o1.EGet(r1))
	assert.Equal(t, o2, o1.EGetResolve(r1, false))
	assert.True(t, o2.EIsSet(r2))
	assert.True(t, o1.EIsSet(r1))

	o2.EUnset(r2)
	assert.Equal(t, nil, o2.EGet(r2))
	assert.Equal(t, nil, o1.EGet(r1))
	assert.False(t, o2.EIsSet(r2))
	assert.False(t, o1.EIsSet(r1))
}

func TestDynamicEObject_Proxy(t *testing.T) {
	// meta model
	c1 := GetFactory().CreateEClass()
	c2 := GetFactory().CreateEClass()
	c3 := GetFactory().CreateEClass()

	r1 := GetFactory().CreateEReference()
	r1.SetContainment(true)
	r1.SetName("r1")
	r1.SetLowerBound(0)
	r1.SetUpperBound(-1)
	r1.SetEType(c2)

	r3 := GetFactory().CreateEReference()
	r3.SetName("r3")
	r3.SetEType(c2)
	r3.SetResolveProxies(true)

	c1.GetEStructuralFeatures().Add(r1)
	c1.SetName("c1")

	c2.SetName("c2")

	c3.GetEStructuralFeatures().Add(r3)
	c3.SetName("c3")

	// model - a container object with two children and another object
	// with one of these child reference
	o1 := NewDynamicEObjectImpl()
	o1.SetEClass(c1)

	o1c1 := NewDynamicEObjectImpl()
	o1c1.SetEClass(c2)

	o1c2 := NewDynamicEObjectImpl()
	o1c2.SetEClass(c2)

	o1cs, _ := o1.EGet(r1).(EList)
	assert.NotNil(t, o1cs)
	o1cs.AddAll(NewImmutableEList([]any{o1c1, o1c2}))

	o3 := NewDynamicEObjectImpl()
	o3.SetEClass(c3)

	// add to resource to enable proxy resolution
	resource := NewEResourceImpl()
	resource.SetURI(NewURIBuilder(nil).SetPath("r").URI())
	resource.GetContents().AddAll(NewImmutableEList([]any{o1, o3}))

	resourceSet := NewEResourceSetImpl()
	resourceSet.GetResources().Add(resource)

	oproxy := NewDynamicEObjectImpl()
	oproxy.ESetProxyURI(NewURIBuilder(nil).SetPath("r").SetFragment("//@r1.1").URI())
	assert.False(t, o3.EIsSet(r3))

	o3.ESet(r3, oproxy)
	assert.Equal(t, oproxy, o3.EGetResolve(r3, false))
	assert.Equal(t, o1c2, o3.EGetResolve(r3, true))
	assert.True(t, o3.EIsSet(r3))

	o3.EUnset(r3)
	assert.False(t, o3.EIsSet(r3))
	assert.Nil(t, o3.EGet(r3))
}

func TestDynamicEObject_Bidirectional(t *testing.T) {

	r1 := GetFactory().CreateEReference()
	r1.SetName("ref")

	r2 := GetFactory().CreateEReference()
	r2.SetName("parent")
	r2.SetEOpposite(r1)

	c1 := GetFactory().CreateEClass()
	c1.GetEStructuralFeatures().Add(r1)

	c2 := GetFactory().CreateEClass()
	c2.GetEStructuralFeatures().Add(r2)

	o1 := NewDynamicEObjectImpl()
	o1.SetEClass(c1)

	o2 := NewDynamicEObjectImpl()
	o2.SetEClass(c2)

	assert.False(t, o2.EIsSet(r2))
	assert.False(t, o1.EIsSet(r1))

	o2.ESet(r2, o1)
	assert.Equal(t, o1, o2.EGet(r2))
	assert.Equal(t, o1, o2.EGetResolve(r2, false))
	assert.Equal(t, o2, o1.EGet(r1))
	assert.Equal(t, o2, o1.EGetResolve(r1, false))
	assert.True(t, o2.EIsSet(r2))
	assert.True(t, o1.EIsSet(r1))

	o2.EUnset(r2)
	assert.Equal(t, nil, o2.EGet(r2))
	assert.Equal(t, nil, o1.EGet(r1))
	assert.False(t, o2.EIsSet(r2))
	assert.False(t, o1.EIsSet(r1))

}
