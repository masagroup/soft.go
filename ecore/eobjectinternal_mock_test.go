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

func TestMockEObjectInternal_EDynamicProperties(t *testing.T) {
	o := NewMockEObjectInternal(t)
	p := NewMockEObjectProperties(t)
	m := NewMockRun(t)
	o.EXPECT().EDynamicProperties().Return(p).Run(func() { m.Run() }).Once()
	o.EXPECT().EDynamicProperties().Call.Return(func() EDynamicProperties {
		return p
	}).Once()
	assert.Equal(t, p, o.EDynamicProperties())
	assert.Equal(t, p, o.EDynamicProperties())
}

func TestMockEObjectInternal_EStaticClass(t *testing.T) {
	o := NewMockEObjectInternal(t)
	c := NewMockEClass(t)
	m := NewMockRun(t)
	o.EXPECT().EStaticClass().Return(c).Run(func() { m.Run() }).Once()
	o.EXPECT().EStaticClass().Call.Return(func() EClass {
		return c
	}).Once()
	assert.Equal(t, c, o.EStaticClass())
	assert.Equal(t, c, o.EStaticClass())
}

func TestMockEObjectInternal_EStaticFeatureCount(t *testing.T) {
	o := NewMockEObjectInternal(t)
	m := NewMockRun(t)
	o.EXPECT().EStaticFeatureCount().Return(1).Run(func() { m.Run() }).Once()
	o.EXPECT().EStaticFeatureCount().Call.Return(func() int {
		return 2
	}).Once()
	assert.Equal(t, 1, o.EStaticFeatureCount())
	assert.Equal(t, 2, o.EStaticFeatureCount())
}

func TestMockEObjectInternal_EInternalResource(t *testing.T) {
	o := NewMockEObjectInternal(t)
	r := NewMockEResource(t)
	m := NewMockRun(t)
	o.EXPECT().EInternalResource().Return(r).Run(func() { m.Run() }).Once()
	o.EXPECT().EInternalResource().Call.Return(func() EResource {
		return r
	}).Once()
	assert.Equal(t, r, o.EInternalResource())
	assert.Equal(t, r, o.EInternalResource())
	mock.AssertExpectationsForObjects(t, o, r)
}

func TestMockEObjectInternal_EInternalContainer(t *testing.T) {
	o := NewMockEObjectInternal(t)
	c := NewMockEObject(t)
	m := NewMockRun(t)
	o.EXPECT().EInternalContainer().Return(c).Run(func() { m.Run() }).Once()
	o.EXPECT().EInternalContainer().Call.Return(func() EObject {
		return c
	}).Once()
	assert.Equal(t, c, o.EInternalContainer())
	assert.Equal(t, c, o.EInternalContainer())
}

func TestMockEObjectInternal_ESetInternalContainer(t *testing.T) {
	o := NewMockEObjectInternal(t)
	c := NewMockEObject(t)
	m := NewMockRun(t, c, 1)
	o.EXPECT().ESetInternalContainer(c, 1).Return().Run(func(container EObject, containerFeatureID int) { m.Run(container, containerFeatureID) }).Once()
	o.ESetInternalContainer(c, 1)
}

func TestMockEObjectInternal_ESetInternalResource(t *testing.T) {
	o := NewMockEObjectInternal(t)
	r := NewMockEResource(t)
	m := NewMockRun(t, r)
	o.EXPECT().ESetInternalResource(r).Return().Run(func(resource EResource) { m.Run(resource) }).Once()
	o.ESetInternalResource(r)
}

func TestMockEObjectInternal_EInternalContainerFeatureID(t *testing.T) {
	o := NewMockEObjectInternal(t)
	m := NewMockRun(t)
	o.EXPECT().EInternalContainerFeatureID().Return(1).Run(func() { m.Run() }).Once()
	o.EXPECT().EInternalContainerFeatureID().Call.Return(func() int {
		return 2
	}).Once()
	assert.Equal(t, 1, o.EInternalContainerFeatureID())
	assert.Equal(t, 2, o.EInternalContainerFeatureID())
}

func TestMockEObjectInternal_ESetResource(t *testing.T) {
	o := NewMockEObjectInternal(t)
	r := NewMockEResource(t)
	n := NewMockENotificationChain(t)
	m := NewMockRun(t, r, n)
	o.EXPECT().ESetResource(r, n).Return(n).Run(func(r EResource, n ENotificationChain) { m.Run(r, n) }).Once()
	o.EXPECT().ESetResource(r, n).Call.Return(func(r EResource, n ENotificationChain) ENotificationChain {
		return n
	}).Once()
	assert.Equal(t, n, o.ESetResource(r, n))
	assert.Equal(t, n, o.ESetResource(r, n))
}

func TestMockEObjectInternal_EInverseAdd(t *testing.T) {
	o := NewMockEObjectInternal(t)
	obj := NewMockEObject(t)
	n := NewMockENotificationChain(t)
	m := NewMockRun(t, obj, 1, n)
	o.EXPECT().EInverseAdd(obj, 1, n).Return(n).Run(func(o EObject, f int, n ENotificationChain) { m.Run(o, f, n) }).Once()
	o.EXPECT().EInverseAdd(obj, 1, n).Call.Return(func(o EObject, f int, n ENotificationChain) ENotificationChain {
		return n
	}).Once()
	assert.Equal(t, n, o.EInverseAdd(obj, 1, n))
	assert.Equal(t, n, o.EInverseAdd(obj, 1, n))
}

func TestMockEObjectInternal_EInverseRemove(t *testing.T) {
	o := NewMockEObjectInternal(t)
	obj := NewMockEObject(t)
	n := NewMockENotificationChain(t)
	m := NewMockRun(t, obj, 1, n)
	o.EXPECT().EInverseRemove(obj, 1, n).Return(n).Run(func(o EObject, f int, n ENotificationChain) { m.Run(o, f, n) }).Once()
	o.EXPECT().EInverseRemove(obj, 1, n).Call.Return(func(o EObject, f int, n ENotificationChain) ENotificationChain {
		return n
	}).Once()
	assert.Equal(t, n, o.EInverseRemove(obj, 1, n))
	assert.Equal(t, n, o.EInverseRemove(obj, 1, n))
}

func TestMockEObjectInternal_EBasicInverseAdd(t *testing.T) {
	o := NewMockEObjectInternal(t)
	obj := NewMockEObject(t)
	n := NewMockENotificationChain(t)
	m := NewMockRun(t, obj, 1, n)
	o.EXPECT().EBasicInverseAdd(obj, 1, n).Return(n).Run(func(o EObject, f int, n ENotificationChain) { m.Run(o, f, n) }).Once()
	o.EXPECT().EBasicInverseAdd(obj, 1, n).Call.Return(func(o EObject, f int, n ENotificationChain) ENotificationChain {
		return n
	}).Once()
	assert.Equal(t, n, o.EBasicInverseAdd(obj, 1, n))
	assert.Equal(t, n, o.EBasicInverseAdd(obj, 1, n))
}

func TestMockEObjectInternal_EBasicInverseRemove(t *testing.T) {
	o := NewMockEObjectInternal(t)
	obj := NewMockEObject(t)
	n := NewMockENotificationChain(t)
	m := NewMockRun(t, obj, 1, n)
	o.EXPECT().EBasicInverseRemove(obj, 1, n).Return(n).Run(func(o EObject, f int, n ENotificationChain) { m.Run(o, f, n) }).Once()
	o.EXPECT().EBasicInverseRemove(obj, 1, n).Call.Return(func(o EObject, f int, n ENotificationChain) ENotificationChain {
		return n
	}).Once()
	assert.Equal(t, n, o.EBasicInverseRemove(obj, 1, n))
	assert.Equal(t, n, o.EBasicInverseRemove(obj, 1, n))
}

func TestMockEObjectInternal_EFeatureID(t *testing.T) {
	mockObject := NewMockEObjectInternal(t)
	mockFeature := NewMockEStructuralFeature(t)
	m := NewMockRun(t, mockFeature)
	mockObject.EXPECT().EFeatureID(mockFeature).Return(1).Run(func(f EStructuralFeature) { m.Run(f) }).Once()
	mockObject.EXPECT().EFeatureID(mockFeature).Call.Return(func(EStructuralFeature) int {
		return 2
	}).Once()
	assert.Equal(t, 1, mockObject.EFeatureID(mockFeature))
	assert.Equal(t, 2, mockObject.EFeatureID(mockFeature))
}

func TestMockEObjectInternal_EDerivedFeatureID(t *testing.T) {
	o := NewMockEObjectInternal(t)
	obj := NewMockEObject(t)
	m := NewMockRun(t, obj, 1)
	o.EXPECT().EDerivedFeatureID(obj, 1).Return(1).Run(func(container EObject, featureID int) { m.Run(container, featureID) }).Once()
	o.EXPECT().EDerivedFeatureID(obj, 1).Call.Return(func(container EObject, featureID int) int {
		return 2
	}).Once()
	assert.Equal(t, 1, o.EDerivedFeatureID(obj, 1))
	assert.Equal(t, 2, o.EDerivedFeatureID(obj, 1))
}

func TestMockEObjectInternal_EOperationID(t *testing.T) {
	mockObject := NewMockEObjectInternal(t)
	mockOperation := NewMockEOperation(t)
	m := NewMockRun(t, mockOperation)
	mockObject.EXPECT().EOperationID(mockOperation).Return(1).Run(func(operation EOperation) { m.Run(operation) }).Once()
	mockObject.EXPECT().EOperationID(mockOperation).Call.Return(func(EOperation) int {
		return 2
	}).Once()
	assert.Equal(t, 1, mockObject.EOperationID(mockOperation))
	assert.Equal(t, 2, mockObject.EOperationID(mockOperation))
}

func TestMockEObjectInternal_EDerivedOperationID(t *testing.T) {
	o := NewMockEObjectInternal(t)
	obj := NewMockEObject(t)
	m := NewMockRun(t, obj, 1)
	o.EXPECT().EDerivedOperationID(obj, 1).Return(1).Run(func(container EObject, featureID int) { m.Run(container, featureID) }).Once()
	o.EXPECT().EDerivedOperationID(obj, 1).Call.Return(func(container EObject, featureID int) int {
		return 2
	}).Once()
	assert.Equal(t, 1, o.EDerivedOperationID(obj, 1))
	assert.Equal(t, 2, o.EDerivedOperationID(obj, 1))
}

func TestMockEObjectInternal_EGetFromID(t *testing.T) {
	o := NewMockEObjectInternal(t)
	obj := NewMockEObject(t)
	m := NewMockRun(t, 1, false)
	o.EXPECT().EGetFromID(1, false).Return(obj).Run(func(featureID int, resolve bool) { m.Run(featureID, resolve) }).Once()
	o.EXPECT().EGetFromID(1, true).Call.Return(func(featureID int, resolve bool) any {
		return obj
	}).Once()
	assert.Equal(t, obj, o.EGetFromID(1, false))
	assert.Equal(t, obj, o.EGetFromID(1, true))
}

func TestMockEObjectInternal_EIsSetFromID(t *testing.T) {
	o := NewMockEObjectInternal(t)
	m := NewMockRun(t, 1)
	o.EXPECT().EIsSetFromID(1).Return(false).Run(func(featureID int) { m.Run(featureID) }).Once()
	o.EXPECT().EIsSetFromID(1).Call.Return(func(featureID int) bool {
		return true
	}).Once()
	assert.False(t, o.EIsSetFromID(1))
	assert.True(t, o.EIsSetFromID(1))
}

func TestMockEObjectInternal_ESetFromID(t *testing.T) {
	o := NewMockEObjectInternal(t)
	obj := NewMockEObject(t)
	m := NewMockRun(t, 1, obj)
	o.EXPECT().ESetFromID(1, obj).Return().Run(func(featureID int, newValue interface{}) { m.Run(featureID, newValue) }).Once()
	o.ESetFromID(1, obj)
}

func TestMockEObjectInternal_EUnsetFromID(t *testing.T) {
	o := NewMockEObjectInternal(t)
	m := NewMockRun(t, 1)
	o.EXPECT().EUnsetFromID(1).Return().Run(func(featureID int) { m.Run(featureID) }).Once()
	o.EUnsetFromID(1)
}

func TestMockEObjectInternal_EInvokeFromID(t *testing.T) {
	o := NewMockEObjectInternal(t)
	obj := NewMockEObject(t)
	args := NewMockEList(t)
	m := NewMockRun(t, 1, args)
	o.EXPECT().EInvokeFromID(1, args).Return(obj).Run(func(operationID int, arguments EList) { m.Run(operationID, arguments) }).Once()
	o.EXPECT().EInvokeFromID(1, args).Call.Return(func(operationID int, arguments EList) any {
		return obj
	}).Once()
	assert.Equal(t, obj, o.EInvokeFromID(1, args))
	assert.Equal(t, obj, o.EInvokeFromID(1, args))
}

func TestMockEObjectInternal_EObjectForFragmentSegment(t *testing.T) {
	o := NewMockEObjectInternal(t)
	obj := NewMockEObject(t)
	m := NewMockRun(t, "uri")
	o.EXPECT().EObjectForFragmentSegment("uri").Return(obj).Run(func(_a0 string) { m.Run(_a0) }).Once()
	o.EXPECT().EObjectForFragmentSegment("uri").Call.Return(func(string) EObject {
		return obj
	}).Once()
	assert.Equal(t, obj, o.EObjectForFragmentSegment("uri"))
	assert.Equal(t, obj, o.EObjectForFragmentSegment("uri"))
}

func TestMockEObjectInternal_EURIFragmentSegment(t *testing.T) {
	o := NewMockEObjectInternal(t)
	f := NewMockEStructuralFeature(t)
	obj := NewMockEObject(t)
	m := NewMockRun(t, f, obj)
	o.EXPECT().EURIFragmentSegment(f, obj).Return("uri").Run(func(_a0 EStructuralFeature, _a1 EObject) { m.Run(_a0, _a1) }).Once()
	o.EXPECT().EURIFragmentSegment(f, obj).Call.Return(func(EStructuralFeature, EObject) string {
		return "uri"
	}).Once()
	assert.Equal(t, "uri", o.EURIFragmentSegment(f, obj))
	assert.Equal(t, "uri", o.EURIFragmentSegment(f, obj))
}

func TestMockEObjectInternal_EProxyURI(t *testing.T) {
	o := NewMockEObjectInternal(t)
	uri := NewURI("test:///file.t")
	m := NewMockRun(t)
	o.EXPECT().EProxyURI().Return(uri).Run(func() { m.Run() }).Once()
	o.EXPECT().EProxyURI().Call.Return(func() *URI {
		return uri
	}).Once()
	assert.Equal(t, uri, o.EProxyURI())
	assert.Equal(t, uri, o.EProxyURI())
}

func TestMockEObjectInternal_ESetProxyURI(t *testing.T) {
	o := NewMockEObjectInternal(t)
	mockURI := NewURI("test:///file.t")
	m := NewMockRun(t, mockURI)
	o.EXPECT().ESetProxyURI(mockURI).Return().Run(func(uri *URI) { m.Run(uri) }).Once()
	o.ESetProxyURI(mockURI)
}

func TestMockEObjectInternal_EResolveProxy(t *testing.T) {
	o := NewMockEObjectInternal(t)
	obj := NewMockEObject(t)
	result := NewMockEObject(t)
	m := NewMockRun(t, obj)
	o.EXPECT().EResolveProxy(obj).Return(result).Run(func(proxy EObject) { m.Run(proxy) }).Once()
	o.EXPECT().EResolveProxy(obj).Call.Return(func(proxy EObject) EObject {
		return result
	}).Once()
	assert.Equal(t, result, o.EResolveProxy(obj))
	assert.Equal(t, result, o.EResolveProxy(obj))
}
