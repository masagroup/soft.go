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

func TestEcoreUtilsConvertToString(t *testing.T) {
	mockObject := NewMockEObject(t)
	mockDataType := NewMockEDataType(t)
	mockPackage := NewMockEPackage(t)
	mockFactory := NewMockEFactory(t)
	mockDataType.On("GetEPackage").Once().Return(mockPackage)
	mockPackage.On("GetEFactoryInstance").Once().Return(mockFactory)
	mockFactory.On("ConvertToString", mockDataType, mockObject).Once().Return("test")
	assert.Equal(t, "test", ConvertToString(mockDataType, mockObject))
	mock.AssertExpectationsForObjects(t, mockObject, mockDataType, mockPackage, mockFactory)
}

func TestEcoreUtilsCreateFromString(t *testing.T) {
	mockObject := NewMockEObject(t)
	mockDataType := NewMockEDataType(t)
	mockPackage := NewMockEPackage(t)
	mockFactory := NewMockEFactory(t)
	mockDataType.On("GetEPackage").Once().Return(mockPackage)
	mockPackage.On("GetEFactoryInstance").Once().Return(mockFactory)
	mockFactory.On("CreateFromString", mockDataType, "test").Once().Return(mockObject)
	assert.Equal(t, mockObject, CreateFromString(mockDataType, "test"))
	mock.AssertExpectationsForObjects(t, mockObject, mockDataType, mockPackage, mockFactory)
}

func TestEcoreUtilsGetObjectID(t *testing.T) {
	mockObject := NewMockEObject(t)
	mockAttribute := NewMockEAttribute(t)
	mockClass := NewMockEClass(t)
	mockDataType := NewMockEDataType(t)
	mockPackage := NewMockEPackage(t)
	mockFactory := NewMockEFactory(t)
	mockValue := NewMockEObject(t)

	mockObject.On("EClass").Return(mockClass).Once()
	mockClass.On("GetEIDAttribute").Return(nil).Once()
	assert.Equal(t, "", GetEObjectID(mockObject))
	mock.AssertExpectationsForObjects(t, mockObject, mockClass)

	mockObject.On("EClass").Return(mockClass).Once()
	mockObject.On("EIsSet", mockAttribute).Return(false).Once()
	mockClass.On("GetEIDAttribute").Return(mockAttribute).Once()
	assert.Equal(t, "", GetEObjectID(mockObject))
	mock.AssertExpectationsForObjects(t, mockObject, mockClass)

	mockObject.On("EClass").Return(mockClass).Once()
	mockObject.On("EIsSet", mockAttribute).Return(true).Once()
	mockObject.On("EGet", mockAttribute).Return(mockValue).Once()
	mockClass.On("GetEIDAttribute").Return(mockAttribute).Once()
	mockAttribute.On("GetEAttributeType").Return(mockDataType).Once()
	mockDataType.On("GetEPackage").Once().Return(mockPackage)
	mockPackage.On("GetEFactoryInstance").Once().Return(mockFactory)
	mockFactory.On("ConvertToString", mockDataType, mockValue).Once().Return("test")
	assert.Equal(t, "test", GetEObjectID(mockObject))
	mock.AssertExpectationsForObjects(t, mockObject, mockClass, mockDataType, mockPackage, mockFactory, mockValue)
}

func TestEcoreUtilsSetObjectID(t *testing.T) {
	mockObject := NewMockEObject(t)
	mockAttribute := NewMockEAttribute(t)
	mockClass := NewMockEClass(t)
	mockDataType := NewMockEDataType(t)
	mockPackage := NewMockEPackage(t)
	mockFactory := NewMockEFactory(t)
	mockValue := NewMockEObject(t)

	mockObject.On("EClass").Return(mockClass).Once()
	mockClass.On("GetEIDAttribute").Return(nil).Once()
	assert.Panics(t, func() { SetEObjectID(mockObject, "test") })
	mock.AssertExpectationsForObjects(t, mockObject, mockClass)

	mockObject.On("EClass").Return(mockClass).Once()
	mockObject.On("EUnset", mockAttribute).Once()
	mockClass.On("GetEIDAttribute").Return(mockAttribute).Once()
	SetEObjectID(mockObject, "")
	mock.AssertExpectationsForObjects(t, mockObject, mockClass, mockAttribute)

	mockObject.On("EClass").Return(mockClass).Once()
	mockObject.On("ESet", mockAttribute, mockValue).Once()
	mockClass.On("GetEIDAttribute").Return(mockAttribute).Once()
	mockAttribute.On("GetEAttributeType").Return(mockDataType).Once()
	mockDataType.On("GetEPackage").Once().Return(mockPackage)
	mockPackage.On("GetEFactoryInstance").Once().Return(mockFactory)
	mockFactory.On("CreateFromString", mockDataType, "test").Once().Return(mockValue)
	SetEObjectID(mockObject, "test")
	mock.AssertExpectationsForObjects(t, mockObject, mockClass, mockDataType, mockPackage, mockFactory, mockValue)
}

func TestEcoreUtilsCopyNil(t *testing.T) {
	assert.Nil(t, Copy(nil))
}

func TestEcoreUtilsEqualsNil(t *testing.T) {
	assert.True(t, Equals(nil, nil))
	assert.False(t, Equals(nil, NewMockEObject(t)))
	assert.False(t, Equals(NewMockEObject(t), nil))
}

func TestEcoreUtilsEqualsProxy(t *testing.T) {

	obj1 := NewMockEObjectInternal(t)
	obj2 := NewMockEObjectInternal(t)
	obj1.On("EIsProxy").Once().Return(true)
	obj1.On("EProxyURI").Once().Return(NewURI("test"))
	obj2.On("EProxyURI").Once().Return(NewURI("test"))
	assert.True(t, Equals(obj1, obj2))

	obj1.On("EIsProxy").Once().Return(true)
	obj1.On("EProxyURI").Once().Return(NewURI("test1"))
	obj2.On("EProxyURI").Once().Return(NewURI("test2"))
	assert.False(t, Equals(obj1, obj2))

	obj1.On("EIsProxy").Once().Return(true)
	obj1.On("EProxyURI").Once().Return(NewURI("test"))
	obj2.On("EProxyURI").Once().Return(nil)
	assert.False(t, Equals(obj1, obj2))

	obj1.On("EIsProxy").Once().Return(false)
	obj2.On("EIsProxy").Once().Return(true)
	assert.False(t, Equals(obj1, obj2))

	mock.AssertExpectationsForObjects(t, obj1, obj2)
}

func TestEcoreUtilsEqualsClass(t *testing.T) {
	obj1 := NewMockEObjectInternal(t)
	obj2 := NewMockEObjectInternal(t)
	obj1.On("EIsProxy").Once().Return(false)
	obj2.On("EIsProxy").Once().Return(false)
	obj1.On("EClass").Once().Return(NewMockEClass(t))
	obj2.On("EClass").Once().Return(NewMockEClass(t))
	assert.False(t, Equals(obj1, obj2))
	mock.AssertExpectationsForObjects(t, obj1, obj2)
}

func TestEcoreUtilsCopyAttribute(t *testing.T) {
	ecoreFactory := GetFactory()
	ecorePackage := GetPackage()

	// the meta model
	ePackage := ecoreFactory.CreateEPackage()
	eFactory := ecoreFactory.CreateEFactory()
	eClass := ecoreFactory.CreateEClass()
	ePackage.SetEFactoryInstance(eFactory)
	ePackage.GetEClassifiers().Add(eClass)
	eAttribute1 := ecoreFactory.CreateEAttribute()
	eAttribute1.SetName("attribute1")
	eAttribute1.SetEType(ecorePackage.GetEInt())
	eAttribute2 := ecoreFactory.CreateEAttribute()
	eAttribute2.SetName("attribute2")
	eAttribute2.SetEType(ecorePackage.GetEString())
	eClass.GetEStructuralFeatures().AddAll(NewImmutableEList([]any{eAttribute1, eAttribute2}))

	// the model
	eObject := eFactory.Create(eClass)
	eObject.ESet(eAttribute1, 2)
	eObject.ESet(eAttribute2, "test")

	eObjectCopy := Copy(eObject)
	assert.True(t, Equals(eObject, eObjectCopy))

	eObject.ESet(eAttribute2, "test2")
	assert.False(t, Equals(eObject, eObjectCopy))
}

func TestEcoreUtilsCopyAllAttribute(t *testing.T) {
	ecoreFactory := GetFactory()
	ecorePackage := GetPackage()

	// the meta model
	ePackage := ecoreFactory.CreateEPackage()
	eFactory := ecoreFactory.CreateEFactory()
	eClass := ecoreFactory.CreateEClass()
	ePackage.SetEFactoryInstance(eFactory)
	ePackage.GetEClassifiers().Add(eClass)
	eAttribute1 := ecoreFactory.CreateEAttribute()
	eAttribute1.SetName("attribute1")
	eAttribute1.SetEType(ecorePackage.GetEInt())
	eAttribute2 := ecoreFactory.CreateEAttribute()
	eAttribute2.SetName("attribute2")
	eAttribute2.SetEType(ecorePackage.GetEString())
	eClass.GetEStructuralFeatures().AddAll(NewImmutableEList([]any{eAttribute1, eAttribute2}))

	// the model
	eObject1 := eFactory.Create(eClass)
	eObject1.ESet(eAttribute1, 2)
	eObject1.ESet(eAttribute2, "test")

	eObject2 := eFactory.Create(eClass)
	eObject2.ESet(eAttribute1, 3)
	eObject2.ESet(eAttribute2, "test3")

	list := NewImmutableEList([]any{eObject1, eObject2})
	listCopy := CopyAll(list)
	assert.True(t, EqualsAll(list, listCopy))
}

func TestEcoreUtilsCopyContainment(t *testing.T) {
	ecoreFactory := GetFactory()
	ecorePackage := GetPackage()

	// the meta model
	ePackage := ecoreFactory.CreateEPackage()
	eFactory := ecoreFactory.CreateEFactory()
	eClass1 := ecoreFactory.CreateEClass()
	eClass2 := ecoreFactory.CreateEClass()
	ePackage.SetEFactoryInstance(eFactory)
	ePackage.GetEClassifiers().AddAll(NewImmutableEList([]any{eClass1, eClass2}))

	eAttribute1 := ecoreFactory.CreateEAttribute()
	eAttribute1.SetName("attribute1")
	eAttribute1.SetEType(ecorePackage.GetEInt())
	eAttribute2 := ecoreFactory.CreateEAttribute()
	eAttribute2.SetName("attribute2")
	eAttribute2.SetEType(ecorePackage.GetEString())
	eClass2.GetEStructuralFeatures().AddAll(NewImmutableEList([]any{eAttribute1, eAttribute2}))

	eReference1 := ecoreFactory.CreateEReference()
	eReference1.SetName("reference1")
	eReference1.SetContainment(true)
	eReference1.SetEType(eClass2)
	eClass1.GetEStructuralFeatures().Add(eReference1)

	// the model
	eObject1 := eFactory.Create(eClass1)
	eObject2 := eFactory.Create(eClass2)
	eObject2.ESet(eAttribute1, 2)
	eObject2.ESet(eAttribute2, "test")
	eObject1.ESet(eReference1, eObject2)

	eObject1Copy := Copy(eObject1)
	assert.True(t, Equals(eObject1, eObject1Copy))

	eObject2.ESet(eAttribute2, "test2")
	assert.False(t, Equals(eObject1, eObject1Copy))
}

func TestEcoreUtilsCopyReferences(t *testing.T) {
	ecoreFactory := GetFactory()
	ecorePackage := GetPackage()

	// the meta model
	ePackage := ecoreFactory.CreateEPackage()
	eFactory := ecoreFactory.CreateEFactory()
	eClass1 := ecoreFactory.CreateEClass()
	eClass2 := ecoreFactory.CreateEClass()
	ePackage.SetEFactoryInstance(eFactory)
	ePackage.GetEClassifiers().AddAll(NewImmutableEList([]any{eClass1, eClass2}))

	eAttribute1 := ecoreFactory.CreateEAttribute()
	eAttribute1.SetName("attribute1")
	eAttribute1.SetEType(ecorePackage.GetEInt())
	eAttribute2 := ecoreFactory.CreateEAttribute()
	eAttribute2.SetName("attribute2")
	eAttribute2.SetEType(ecorePackage.GetEString())
	eClass2.GetEStructuralFeatures().AddAll(NewImmutableEList([]any{eAttribute1, eAttribute2}))

	eReference1 := ecoreFactory.CreateEReference()
	eReference1.SetName("reference1")
	eReference1.SetContainment(true)
	eReference1.SetEType(eClass2)
	eReference2 := ecoreFactory.CreateEReference()
	eReference2.SetName("reference2")
	eReference2.SetContainment(false)
	eReference2.SetEType(eClass2)
	eClass1.GetEStructuralFeatures().AddAll(NewImmutableEList([]any{eReference1, eReference2}))

	// the model
	eObject1 := eFactory.Create(eClass1)
	eObject2 := eFactory.Create(eClass2)
	eObject2.ESet(eAttribute1, 2)
	eObject2.ESet(eAttribute2, "test")
	eObject1.ESet(eReference1, eObject2)
	eObject1.ESet(eReference2, eObject2)

	eObject1Copy := Copy(eObject1)
	assert.True(t, Equals(eObject1, eObject1Copy))

	eObject2.ESet(eAttribute2, "test2")
	assert.False(t, Equals(eObject1, eObject1Copy))
}

func TestEcoreUtilsCopyProxy(t *testing.T) {
	// the meta model
	ecoreFactory := GetFactory()

	// the meta model
	ePackage := ecoreFactory.CreateEPackage()
	eFactory := ecoreFactory.CreateEFactory()
	eClass := ecoreFactory.CreateEClass()
	ePackage.SetEFactoryInstance(eFactory)
	ePackage.GetEClassifiers().Add(eClass)

	// the model
	eObject := eFactory.Create(eClass)
	eObject.(EObjectInternal).ESetProxyURI(NewURI("testPath"))

	eObjectCopy := Copy(eObject)
	assert.True(t, Equals(eObject, eObjectCopy))
}

func TestEcoreUtilsCopyReal(t *testing.T) {
	eClass := GetPackage().GetEClass()
	eClassCopy := Copy(eClass)
	assert.True(t, Equals(eClass, eClassCopy))
}

func TestEcoreUtils_GetURI(t *testing.T) {
	mockURI, _ := ParseURI("test://file.t")
	mockEObject := NewMockEObjectInternal(t)
	mockEObject.On("EIsProxy").Return(true).Once()
	mockEObject.On("EProxyURI").Return(mockURI).Once()
	assert.Equal(t, mockURI, GetURI(mockEObject))
}

func TestEcoreUtils_Remove(t *testing.T) {
	mockObject := NewMockEObjectInternal(t)
	mockReference := NewMockEReference(t)
	mockContainer := NewMockEObject(t)

	// resource - container - feature single
	mockObject.On("EInternalContainer").Return(mockContainer).Once()
	mockObject.On("EContainmentFeature").Return(mockReference).Once()
	mockReference.On("IsMany").Return(false).Once()
	mockContainer.On("EUnset", mockReference).Once()
	mockObject.On("EInternalResource").Return(nil).Once()
	Remove(mockObject)
	mock.AssertExpectationsForObjects(t, mockObject, mockReference, mockContainer)

	// resource - container - feature many
	mockList := NewMockEList(t)
	mockObject.On("EInternalContainer").Return(mockContainer).Once()
	mockObject.On("EContainmentFeature").Return(mockReference).Once()
	mockReference.On("IsMany").Return(true).Once()
	mockContainer.On("EGet", mockReference).Return(mockList).Once()
	mockList.On("Remove", mockObject).Return(true).Once()
	mockObject.On("EInternalResource").Return(nil).Once()
	Remove(mockObject)
	mock.AssertExpectationsForObjects(t, mockObject, mockReference, mockContainer)

	// resource - no container
	mockResource := NewMockEResource(t)
	mockObject.On("EInternalContainer").Return(nil).Once()
	mockObject.On("EInternalResource").Return(mockResource).Once()
	mockResource.On("GetContents").Return(mockList).Once()
	mockList.On("Remove", mockObject).Return(true).Once()
	Remove(mockObject)
	mock.AssertExpectationsForObjects(t, mockObject, mockReference, mockContainer, mockResource)
}

func TestEcoreUtils_GetAncestor(t *testing.T) {

	mockObject0 := NewMockEObject(t)
	mockObject1 := NewMockEObject(t)
	mockObject2 := NewMockEObject(t)
	mockClass := NewMockEClass(t)
	mockOtherClass := NewMockEClass(t)

	mockObject0.On("EClass").Return(mockOtherClass).Once()
	mockObject0.On("EContainer").Return(mockObject1).Once()
	mockObject1.On("EClass").Return(mockOtherClass).Once()
	mockObject1.On("EContainer").Return(mockObject2).Once()
	mockObject2.On("EClass").Return(mockClass).Once()
	assert.Equal(t, mockObject2, GetAncestor(mockObject0, mockClass))
	mock.AssertExpectationsForObjects(t, mockObject0, mockObject1, mockObject2, mockClass, mockOtherClass)

	mockObject0.On("EClass").Return(mockOtherClass).Once()
	mockObject0.On("EContainer").Return(mockObject1).Once()
	mockObject1.On("EClass").Return(mockOtherClass).Once()
	mockObject1.On("EContainer").Return(mockObject2).Once()
	mockObject2.On("EClass").Return(mockOtherClass).Once()
	mockObject2.On("EContainer").Return(nil).Once()
	assert.Equal(t, nil, GetAncestor(mockObject0, mockClass))
	mock.AssertExpectationsForObjects(t, mockObject0, mockObject1, mockObject2, mockClass, mockOtherClass)

	assert.Equal(t, nil, GetAncestor(nil, mockClass))
	mock.AssertExpectationsForObjects(t, mockClass)

	mockObject0.On("EClass").Return(mockClass).Once()
	assert.Equal(t, mockObject0, GetAncestor(mockObject0, mockClass))
	mock.AssertExpectationsForObjects(t, mockObject0, mockClass)
}

func TestEcoreUtils_IsAncestor(t *testing.T) {

	mockObject0 := NewMockEObject(t)
	mockObject1 := NewMockEObject(t)
	mockObject2 := NewMockEObject(t)

	assert.True(t, IsAncestor(nil, nil))

	mockObject0.On("EContainer").Return(mockObject1).Once()
	mockObject1.On("EContainer").Return(mockObject2).Once()
	assert.True(t, IsAncestor(mockObject2, mockObject0))
	mock.AssertExpectationsForObjects(t, mockObject0, mockObject1, mockObject2)

	mockObject0.On("EContainer").Return(mockObject1).Once()
	mockObject1.On("EContainer").Return(nil).Once()
	assert.False(t, IsAncestor(mockObject2, mockObject0))
	mock.AssertExpectationsForObjects(t, mockObject0, mockObject1, mockObject2)
}
