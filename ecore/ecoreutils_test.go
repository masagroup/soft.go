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
	"github.com/stretchr/testify/require"
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

func TestEcoreUtils_Delete(t *testing.T) {
	// load package
	ePackage := loadPackage("delete.ecore")
	assert.NotNil(t, ePackage)

	// load resource
	xmiProcessor := NewXMIProcessor(XMIProcessorPackages([]EPackage{ePackage}))
	eResource := xmiProcessor.Load(NewURI("testdata/delete.xmi"))
	require.NotNil(t, eResource)
	require.True(t, eResource.IsLoaded())
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	require.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	root := eResource.GetContents().Get(0).(EObject)
	classRoot := ePackage.GetEClassifier("Root").(EClass)

	refA := classRoot.GetEStructuralFeatureFromName("a")
	listA, _ := root.EGet(refA).(EList)
	require.NotNil(t, listA)
	a1, _ := listA.Get(0).(EObject)
	require.NotNil(t, a1)

	refO := classRoot.GetEStructuralFeatureFromName("o")
	listO, _ := root.EGet(refO).(EList)
	require.NotNil(t, listO)
	o1, _ := listO.Get(0).(EObject)
	require.NotNil(t, o1)

	// check that a1 is referenced by o1
	classO, _ := ePackage.GetEClassifier("O").(EClass)
	require.NotNil(t, classO)
	classO_A := classO.GetEStructuralFeatureFromName("a")
	require.NotNil(t, classO_A)
	classO_Name := classO.GetEStructuralFeatureFromName("name")
	require.NotNil(t, classO_Name)
	o1_aList := o1.EGet(classO_A).(EList)
	require.Equal(t, 2, o1_aList.Size())

	// delete a1
	Delete(a1)

	// check that a1 is not referenced anymore
	require.Equal(t, 1, o1_aList.Size())

	// C
	classC, _ := ePackage.GetEClassifier("C").(EClass)
	require.NotNil(t, classO)
	classC_M := classC.GetEStructuralFeatureFromName("m")
	require.NotNil(t, classC_M)
	classC_Name := classC.GetEStructuralFeatureFromName("name")
	require.NotNil(t, classC_Name)

	refC := classRoot.GetEStructuralFeatureFromName("c")
	listC, _ := root.EGet(refC).(EList)
	require.NotNil(t, listC)
	require.Equal(t, 3, listC.Size())

	// retrieve c1 & c2 & check origin values
	c1, _ := listC.Get(0).(EObject)
	require.NotNil(t, c1)
	require.Equal(t, "c1", c1.EGet(classC_Name))
	require.Equal(t, nil, c1.EGet(classC_M))

	c2, _ := listC.Get(1).(EObject)
	require.NotNil(t, c2)
	require.Equal(t, c1, c2.EGet(classC_M))

	// delete c1
	Delete(c1)

	// check c2 ref no more c1
	require.Equal(t, nil, c2.EGet(classC_M))
}

func TestEcoreUtils_DeleteRecursive(t *testing.T) {
	// load package
	ePackage := loadPackage("delete.ecore")
	assert.NotNil(t, ePackage)

	// load resource
	xmiProcessor := NewXMIProcessor(XMIProcessorPackages([]EPackage{ePackage}))
	eResource := xmiProcessor.Load(NewURI("testdata/delete.xmi"))
	require.NotNil(t, eResource)
	require.True(t, eResource.IsLoaded())
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	require.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))

	root := eResource.GetContents().Get(0).(EObject)
	classRoot := ePackage.GetEClassifier("Root").(EClass)

	refA := classRoot.GetEStructuralFeatureFromName("a")
	listA, _ := root.EGet(refA).(EList)
	require.NotNil(t, listA)
	a1, _ := listA.Get(0).(EObject)
	require.NotNil(t, a1)

	refB := classRoot.GetEStructuralFeatureFromName("b")
	listB, _ := root.EGet(refB).(EList)
	require.NotNil(t, listB)

	refO := classRoot.GetEStructuralFeatureFromName("o")
	listO, _ := root.EGet(refO).(EList)
	require.NotNil(t, listO)
	o1, _ := listO.Get(0).(EObject)
	require.NotNil(t, o1)

	// check that a1 is referenced by o1
	classO, _ := ePackage.GetEClassifier("O").(EClass)
	require.NotNil(t, classO)
	classO_A := classO.GetEStructuralFeatureFromName("a")
	require.NotNil(t, classO_A)
	classO_Name := classO.GetEStructuralFeatureFromName("name")
	require.NotNil(t, classO_Name)
	o1_aList := o1.EGet(classO_A).(EList)
	require.Equal(t, 2, o1_aList.Size())

	//
	classA, _ := ePackage.GetEClassifier("A").(EClass)
	require.NotNil(t, classA)
	classA_B := classA.GetEStructuralFeatureFromName("b")
	require.NotNil(t, classA_B)
	b1 := a1.EGet(classA_B).(EObject)
	require.NotNil(t, b1)
	require.True(t, listB.Contains(b1))
	require.Equal(t, 3, listA.Size())
	require.True(t, listA.Contains(a1))

	// delete a1
	DeleteRecursive(a1, true)

	// check that a1 and its children are no longer referenced or contained
	require.Equal(t, 1, o1_aList.Size())
	require.False(t, o1_aList.Contains(a1))
	require.Equal(t, 2, listB.Size())
	require.False(t, listB.Contains(b1))
	require.Equal(t, 2, listA.Size())
	require.False(t, listA.Contains(a1))
	require.Nil(t, a1.EContainer())
	require.Nil(t, a1.EResource())
}

func TestEcoreUtils_DeleteRecursive_False(t *testing.T) {
	ePackage := loadPackage("delete.ecore")
	require.NotNil(t, ePackage)

	xmiProcessor := NewXMIProcessor(XMIProcessorPackages([]EPackage{ePackage}))
	eResource := xmiProcessor.Load(NewURI("testdata/delete.xmi"))
	require.NotNil(t, eResource)
	require.True(t, eResource.IsLoaded())

	root := eResource.GetContents().Get(0).(EObject)
	classRoot := ePackage.GetEClassifier("Root").(EClass)

	refA := classRoot.GetEStructuralFeatureFromName("a")
	listA, _ := root.EGet(refA).(EList)
	require.NotNil(t, listA)
	a1, _ := listA.Get(0).(EObject)
	require.NotNil(t, a1)

	refB := classRoot.GetEStructuralFeatureFromName("b")
	listB, _ := root.EGet(refB).(EList)
	require.NotNil(t, listB)

	refO := classRoot.GetEStructuralFeatureFromName("o")
	listO, _ := root.EGet(refO).(EList)
	require.NotNil(t, listO)
	o1, _ := listO.Get(0).(EObject)
	require.NotNil(t, o1)

	classO, _ := ePackage.GetEClassifier("O").(EClass)
	classO_A := classO.GetEStructuralFeatureFromName("a")
	o1_aList := o1.EGet(classO_A).(EList)
	require.Equal(t, 2, o1_aList.Size())

	classA, _ := ePackage.GetEClassifier("A").(EClass)
	classA_B := classA.GetEStructuralFeatureFromName("b")
	b1 := a1.EGet(classA_B).(EObject)
	require.NotNil(t, b1)
	require.True(t, listB.Contains(b1))
	require.Equal(t, 3, listB.Size())

	// Non-recursive delete
	DeleteRecursive(a1, false)

	// a1 is removed from container and referencing feature o1.a
	require.Equal(t, 1, o1_aList.Size())
	require.False(t, o1_aList.Contains(a1))
	require.Equal(t, 2, listA.Size())
	require.False(t, listA.Contains(a1))
	// But contained child b1 was NOT in the deleted object set, so root.b still holds it
	require.Equal(t, 3, listB.Size())
	require.True(t, listB.Contains(b1))
}

func TestEcoreUtils_DeleteRecursive_NoResource(t *testing.T) {
	ePackage := loadPackage("delete.ecore")
	require.NotNil(t, ePackage)

	eFactory := ePackage.GetEFactoryInstance()
	classRoot := ePackage.GetEClassifier("Root").(EClass)
	classA := ePackage.GetEClassifier("A").(EClass)
	classB := ePackage.GetEClassifier("B").(EClass)
	classO := ePackage.GetEClassifier("O").(EClass)

	refRootA := classRoot.GetEStructuralFeatureFromName("a")
	refRootB := classRoot.GetEStructuralFeatureFromName("b")
	refRootO := classRoot.GetEStructuralFeatureFromName("o")
	refAB := classA.GetEStructuralFeatureFromName("b")
	refOA := classO.GetEStructuralFeatureFromName("a")

	root := eFactory.Create(classRoot)
	a1 := eFactory.Create(classA)
	b1 := eFactory.Create(classB)
	o1 := eFactory.Create(classO)

	a1.ESet(refAB, b1)
	root.EGet(refRootA).(EList).Add(a1)
	root.EGet(refRootB).(EList).Add(b1)
	o1.EGet(refOA).(EList).Add(a1)
	root.EGet(refRootO).(EList).Add(o1)

	require.Nil(t, root.EResource())
	require.Equal(t, 1, root.EGet(refRootA).(EList).Size())
	require.Equal(t, 1, root.EGet(refRootB).(EList).Size())
	require.Equal(t, 1, o1.EGet(refOA).(EList).Size())

	// Delete a1 recursively on standalone hierarchy without a resource
	DeleteRecursive(a1, true)

	require.Equal(t, 0, root.EGet(refRootA).(EList).Size())
	require.Equal(t, 0, root.EGet(refRootB).(EList).Size())
	require.Equal(t, 0, o1.EGet(refOA).(EList).Size())
	require.Nil(t, a1.EContainer())
}

func TestEcoreUtils_DeleteRecursive_ResourceSet(t *testing.T) {
	ePackage := loadPackage("delete.ecore")
	require.NotNil(t, ePackage)

	eFactory := ePackage.GetEFactoryInstance()
	classRoot := ePackage.GetEClassifier("Root").(EClass)
	classA := ePackage.GetEClassifier("A").(EClass)
	classB := ePackage.GetEClassifier("B").(EClass)
	classO := ePackage.GetEClassifier("O").(EClass)

	refRootA := classRoot.GetEStructuralFeatureFromName("a")
	refRootB := classRoot.GetEStructuralFeatureFromName("b")
	refAB := classA.GetEStructuralFeatureFromName("b")
	refOA := classO.GetEStructuralFeatureFromName("a")

	rs := NewEResourceSetImpl()
	r1 := rs.CreateResource(NewURI("memory://res1.xmi"))
	r2 := rs.CreateResource(NewURI("memory://res2.xmi"))

	root1 := eFactory.Create(classRoot)
	r1.GetContents().Add(root1)

	a1 := eFactory.Create(classA)
	b1 := eFactory.Create(classB)
	a1.ESet(refAB, b1)
	root1.EGet(refRootA).(EList).Add(a1)
	root1.EGet(refRootB).(EList).Add(b1)

	// Object in r2 referencing a1 in r1
	o2 := eFactory.Create(classO)
	o2.EGet(refOA).(EList).Add(a1)
	r2.GetContents().Add(o2)

	require.Equal(t, 1, o2.EGet(refOA).(EList).Size())
	require.Equal(t, 1, root1.EGet(refRootB).(EList).Size())

	// Delete a1 across ResourceSet
	DeleteRecursive(a1, true)

	require.Equal(t, 0, o2.EGet(refOA).(EList).Size())
	require.Equal(t, 0, root1.EGet(refRootA).(EList).Size())
	require.Equal(t, 0, root1.EGet(refRootB).(EList).Size())
}

func TestEcoreUtils_DeleteRecursive_DirectResourceChild(t *testing.T) {
	ePackage := loadPackage("delete.ecore")
	require.NotNil(t, ePackage)

	eFactory := ePackage.GetEFactoryInstance()
	classRoot := ePackage.GetEClassifier("Root").(EClass)
	classA := ePackage.GetEClassifier("A").(EClass)
	classB := ePackage.GetEClassifier("B").(EClass)

	refRootA := classRoot.GetEStructuralFeatureFromName("a")
	refAB := classA.GetEStructuralFeatureFromName("b")

	rs := NewEResourceSetImpl()
	r1 := rs.CreateResource(NewURI("memory://res1.xmi"))
	rChild := rs.CreateResource(NewURI("memory://child.xmi"))

	root := eFactory.Create(classRoot)
	r1.GetContents().Add(root)

	a1 := eFactory.Create(classA)
	root.EGet(refRootA).(EList).Add(a1)

	b1 := eFactory.Create(classB)
	a1.ESet(refAB, b1)

	// Set child b1 to have an internal resource
	b1.(EObjectInternal).ESetResource(rChild, nil)
	require.NotNil(t, b1.(EObjectInternal).EInternalResource())

	// Delete a1 recursively
	DeleteRecursive(a1, true)

	// Direct child should have been removed from containing feature
	require.Nil(t, a1.EGet(refAB))
	require.Equal(t, 0, root.EGet(refRootA).(EList).Size())
}

func TestEcoreUtils_DeleteRecursive_Root(t *testing.T) {
	ePackage := loadPackage("delete.ecore")
	require.NotNil(t, ePackage)

	xmiProcessor := NewXMIProcessor(XMIProcessorPackages([]EPackage{ePackage}))
	eResource := xmiProcessor.Load(NewURI("testdata/delete.xmi"))
	require.NotNil(t, eResource)
	require.True(t, eResource.IsLoaded())

	root := eResource.GetContents().Get(0).(EObject)
	require.Equal(t, 1, eResource.GetContents().Size())

	DeleteRecursive(root, true)

	require.Equal(t, 0, eResource.GetContents().Size())
	require.Nil(t, root.EResource())
}
