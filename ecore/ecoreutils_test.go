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
	mockDataType.EXPECT().GetEPackage().Once().Return(mockPackage)
	mockPackage.EXPECT().GetEFactoryInstance().Once().Return(mockFactory)
	mockFactory.EXPECT().ConvertToString(mockDataType, mockObject).Once().Return("test")
	assert.Equal(t, "test", ConvertToString(mockDataType, mockObject))
	mock.AssertExpectationsForObjects(t, mockObject, mockDataType, mockPackage, mockFactory)
}

func TestEcoreUtilsCreateFromString(t *testing.T) {
	mockObject := NewMockEObject(t)
	mockDataType := NewMockEDataType(t)
	mockPackage := NewMockEPackage(t)
	mockFactory := NewMockEFactory(t)
	mockDataType.EXPECT().GetEPackage().Once().Return(mockPackage)
	mockPackage.EXPECT().GetEFactoryInstance().Once().Return(mockFactory)
	mockFactory.EXPECT().CreateFromString(mockDataType, "test").Once().Return(mockObject)
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

	mockObject.EXPECT().EClass().Return(mockClass).Once()
	mockClass.EXPECT().GetEIDAttribute().Return(nil).Once()
	assert.Equal(t, "", GetEObjectID(mockObject))
	mock.AssertExpectationsForObjects(t, mockObject, mockClass)

	mockObject.EXPECT().EClass().Return(mockClass).Once()
	mockObject.EXPECT().EIsSet(mockAttribute).Return(false).Once()
	mockClass.EXPECT().GetEIDAttribute().Return(mockAttribute).Once()
	assert.Equal(t, "", GetEObjectID(mockObject))
	mock.AssertExpectationsForObjects(t, mockObject, mockClass)

	mockObject.EXPECT().EClass().Return(mockClass).Once()
	mockObject.EXPECT().EIsSet(mockAttribute).Return(true).Once()
	mockObject.EXPECT().EGet(mockAttribute).Return(mockValue).Once()
	mockClass.EXPECT().GetEIDAttribute().Return(mockAttribute).Once()
	mockAttribute.EXPECT().GetEAttributeType().Return(mockDataType).Once()
	mockDataType.EXPECT().GetEPackage().Once().Return(mockPackage)
	mockPackage.EXPECT().GetEFactoryInstance().Once().Return(mockFactory)
	mockFactory.EXPECT().ConvertToString(mockDataType, mockValue).Once().Return("test")
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

	mockObject.EXPECT().EClass().Return(mockClass).Once()
	mockClass.EXPECT().GetEIDAttribute().Return(nil).Once()
	assert.Panics(t, func() { SetEObjectID(mockObject, "test") })
	mock.AssertExpectationsForObjects(t, mockObject, mockClass)

	mockObject.EXPECT().EClass().Return(mockClass).Once()
	mockObject.EXPECT().EUnset(mockAttribute).Once()
	mockClass.EXPECT().GetEIDAttribute().Return(mockAttribute).Once()
	SetEObjectID(mockObject, "")
	mock.AssertExpectationsForObjects(t, mockObject, mockClass, mockAttribute)

	mockObject.EXPECT().EClass().Return(mockClass).Once()
	mockObject.EXPECT().ESet(mockAttribute, mockValue).Once()
	mockClass.EXPECT().GetEIDAttribute().Return(mockAttribute).Once()
	mockAttribute.EXPECT().GetEAttributeType().Return(mockDataType).Once()
	mockDataType.EXPECT().GetEPackage().Once().Return(mockPackage)
	mockPackage.EXPECT().GetEFactoryInstance().Once().Return(mockFactory)
	mockFactory.EXPECT().CreateFromString(mockDataType, "test").Once().Return(mockValue)
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
	obj1.EXPECT().EIsProxy().Once().Return(true)
	obj1.EXPECT().EProxyURI().Once().Return(NewURI("test"))
	obj2.EXPECT().EProxyURI().Once().Return(NewURI("test"))
	assert.True(t, Equals(obj1, obj2))

	obj1.EXPECT().EIsProxy().Once().Return(true)
	obj1.EXPECT().EProxyURI().Once().Return(NewURI("test1"))
	obj2.EXPECT().EProxyURI().Once().Return(NewURI("test2"))
	assert.False(t, Equals(obj1, obj2))

	obj1.EXPECT().EIsProxy().Once().Return(true)
	obj1.EXPECT().EProxyURI().Once().Return(NewURI("test"))
	obj2.EXPECT().EProxyURI().Once().Return(nil)
	assert.False(t, Equals(obj1, obj2))

	obj1.EXPECT().EIsProxy().Once().Return(false)
	obj2.EXPECT().EIsProxy().Once().Return(true)
	assert.False(t, Equals(obj1, obj2))

	mock.AssertExpectationsForObjects(t, obj1, obj2)
}

func TestEcoreUtilsEqualsClass(t *testing.T) {
	obj1 := NewMockEObjectInternal(t)
	obj2 := NewMockEObjectInternal(t)
	obj1.EXPECT().EIsProxy().Once().Return(false)
	obj2.EXPECT().EIsProxy().Once().Return(false)
	obj1.EXPECT().EClass().Once().Return(NewMockEClass(t))
	obj2.EXPECT().EClass().Once().Return(NewMockEClass(t))
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
	mockEObject.EXPECT().EIsProxy().Return(true).Once()
	mockEObject.EXPECT().EProxyURI().Return(mockURI).Once()
	assert.Equal(t, mockURI, GetURI(mockEObject))
}

func TestEcoreUtils_Remove(t *testing.T) {
	mockObject := NewMockEObjectInternal(t)
	mockReference := NewMockEReference(t)
	mockContainer := NewMockEObject(t)

	// resource - container - feature single
	mockObject.EXPECT().EInternalContainer().Return(mockContainer).Once()
	mockObject.EXPECT().EContainmentFeature().Return(mockReference).Once()
	mockReference.EXPECT().IsMany().Return(false).Once()
	mockContainer.EXPECT().EUnset(mockReference).Once()
	mockObject.EXPECT().EInternalResource().Return(nil).Once()
	Remove(mockObject)
	mock.AssertExpectationsForObjects(t, mockObject, mockReference, mockContainer)

	// resource - container - feature many
	mockList := NewMockEList(t)
	mockObject.EXPECT().EInternalContainer().Return(mockContainer).Once()
	mockObject.EXPECT().EContainmentFeature().Return(mockReference).Once()
	mockReference.EXPECT().IsMany().Return(true).Once()
	mockContainer.EXPECT().EGet(mockReference).Return(mockList).Once()
	mockList.EXPECT().Remove(mockObject).Return(true).Once()
	mockObject.EXPECT().EInternalResource().Return(nil).Once()
	Remove(mockObject)
	mock.AssertExpectationsForObjects(t, mockObject, mockReference, mockContainer)

	// resource - no container
	mockResource := NewMockEResource(t)
	mockObject.EXPECT().EInternalContainer().Return(nil).Once()
	mockObject.EXPECT().EInternalResource().Return(mockResource).Once()
	mockResource.EXPECT().GetContents().Return(mockList).Once()
	mockList.EXPECT().Remove(mockObject).Return(true).Once()
	Remove(mockObject)
	mock.AssertExpectationsForObjects(t, mockObject, mockReference, mockContainer, mockResource)
}

func TestEcoreUtils_GetAncestor(t *testing.T) {

	mockObject0 := NewMockEObject(t)
	mockObject1 := NewMockEObject(t)
	mockObject2 := NewMockEObject(t)
	mockClass := NewMockEClass(t)
	mockOtherClass := NewMockEClass(t)

	mockObject0.EXPECT().EClass().Return(mockOtherClass).Once()
	mockObject0.EXPECT().EContainer().Return(mockObject1).Once()
	mockObject1.EXPECT().EClass().Return(mockOtherClass).Once()
	mockObject1.EXPECT().EContainer().Return(mockObject2).Once()
	mockObject2.EXPECT().EClass().Return(mockClass).Once()
	assert.Equal(t, mockObject2, GetAncestor(mockObject0, mockClass))
	mock.AssertExpectationsForObjects(t, mockObject0, mockObject1, mockObject2, mockClass, mockOtherClass)

	mockObject0.EXPECT().EClass().Return(mockOtherClass).Once()
	mockObject0.EXPECT().EContainer().Return(mockObject1).Once()
	mockObject1.EXPECT().EClass().Return(mockOtherClass).Once()
	mockObject1.EXPECT().EContainer().Return(mockObject2).Once()
	mockObject2.EXPECT().EClass().Return(mockOtherClass).Once()
	mockObject2.EXPECT().EContainer().Return(nil).Once()
	assert.Equal(t, nil, GetAncestor(mockObject0, mockClass))
	mock.AssertExpectationsForObjects(t, mockObject0, mockObject1, mockObject2, mockClass, mockOtherClass)

	assert.Equal(t, nil, GetAncestor(nil, mockClass))
	mock.AssertExpectationsForObjects(t, mockClass)

	mockObject0.EXPECT().EClass().Return(mockClass).Once()
	assert.Equal(t, mockObject0, GetAncestor(mockObject0, mockClass))
	mock.AssertExpectationsForObjects(t, mockObject0, mockClass)
}

func TestEcoreUtils_IsAncestor(t *testing.T) {

	mockObject0 := NewMockEObject(t)
	mockObject1 := NewMockEObject(t)
	mockObject2 := NewMockEObject(t)

	assert.True(t, IsAncestor(nil, nil))

	mockObject0.EXPECT().EContainer().Return(mockObject1).Once()
	mockObject1.EXPECT().EContainer().Return(mockObject2).Once()
	assert.True(t, IsAncestor(mockObject2, mockObject0))
	mock.AssertExpectationsForObjects(t, mockObject0, mockObject1, mockObject2)

	mockObject0.EXPECT().EContainer().Return(mockObject1).Once()
	mockObject1.EXPECT().EContainer().Return(nil).Once()
	assert.False(t, IsAncestor(mockObject2, mockObject0))
	mock.AssertExpectationsForObjects(t, mockObject0, mockObject1, mockObject2)
}
