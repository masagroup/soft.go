// *****************************************************************************
//
// This file is part of a MASA library or program.
// Refer to the included end-user license agreement for restrictions.
//
// Copyright (c) 2020 MASA Group
//
// *****************************************************************************

package ecore

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEcoreUtilsConvertToString(t *testing.T) {
	mockObject := &MockEObject{}
	mockDataType := &MockEDataType{}
	mockPackage := &MockEPackage{}
	mockFactory := &MockEFactory{}
	mockDataType.On("GetEPackage").Once().Return(mockPackage)
	mockPackage.On("GetEFactoryInstance").Once().Return(mockFactory)
	mockFactory.On("ConvertToString", mockDataType, mockObject).Once().Return("test")
	assert.Equal(t, "test", ConvertToString(mockDataType, mockObject))
	mock.AssertExpectationsForObjects(t, mockObject, mockDataType, mockPackage, mockFactory)
}

func TestEcoreUtilsCreateFromString(t *testing.T) {
	mockObject := &MockEObject{}
	mockDataType := &MockEDataType{}
	mockPackage := &MockEPackage{}
	mockFactory := &MockEFactory{}
	mockDataType.On("GetEPackage").Once().Return(mockPackage)
	mockPackage.On("GetEFactoryInstance").Once().Return(mockFactory)
	mockFactory.On("CreateFromString", mockDataType, "test").Once().Return(mockObject)
	assert.Equal(t, mockObject, CreateFromString(mockDataType, "test"))
	mock.AssertExpectationsForObjects(t, mockObject, mockDataType, mockPackage, mockFactory)
}

func TestEcoreUtilsGetObjectID(t *testing.T) {
	mockObject := new(MockEObject)
	mockAttribute := new(MockEAttribute)
	mockClass := new(MockEClass)
	mockDataType := new(MockEDataType)
	mockPackage := new(MockEPackage)
	mockFactory := new(MockEFactory)
	mockValue := new(MockEObject)

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
	mockObject := new(MockEObject)
	mockAttribute := new(MockEAttribute)
	mockClass := new(MockEClass)
	mockDataType := new(MockEDataType)
	mockPackage := new(MockEPackage)
	mockFactory := new(MockEFactory)
	mockValue := new(MockEObject)

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
	assert.False(t, Equals(nil, &MockEObject{}))
	assert.False(t, Equals(&MockEObject{}, nil))
}

func TestEcoreUtilsEqualsProxy(t *testing.T) {

	obj1 := &MockEObjectInternal{}
	obj2 := &MockEObjectInternal{}
	obj1.On("EIsProxy").Once().Return(true)
	obj1.On("EProxyURI").Once().Return(&url.URL{Path: "test"})
	obj2.On("EProxyURI").Once().Return(&url.URL{Path: "test"})
	assert.True(t, Equals(obj1, obj2))

	obj1.On("EIsProxy").Once().Return(true)
	obj1.On("EProxyURI").Once().Return(&url.URL{Path: "test1"})
	obj2.On("EProxyURI").Once().Return(&url.URL{Path: "test2"})
	assert.False(t, Equals(obj1, obj2))

	obj1.On("EIsProxy").Once().Return(true)
	obj1.On("EProxyURI").Once().Return(&url.URL{Path: "test"})
	obj2.On("EProxyURI").Once().Return(nil)
	assert.False(t, Equals(obj1, obj2))

	obj1.On("EIsProxy").Once().Return(false)
	obj2.On("EIsProxy").Once().Return(true)
	assert.False(t, Equals(obj1, obj2))

	mock.AssertExpectationsForObjects(t, obj1, obj2)
}

func TestEcoreUtilsEqualsClass(t *testing.T) {
	obj1 := &MockEObjectInternal{}
	obj2 := &MockEObjectInternal{}
	obj1.On("EIsProxy").Once().Return(false)
	obj2.On("EIsProxy").Once().Return(false)
	obj1.On("EClass").Once().Return(&MockEClass{})
	obj2.On("EClass").Once().Return(&MockEClass{})
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
	eClass.GetEStructuralFeatures().AddAll(NewImmutableEList([]interface{}{eAttribute1, eAttribute2}))

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
	eClass.GetEStructuralFeatures().AddAll(NewImmutableEList([]interface{}{eAttribute1, eAttribute2}))

	// the model
	eObject1 := eFactory.Create(eClass)
	eObject1.ESet(eAttribute1, 2)
	eObject1.ESet(eAttribute2, "test")

	eObject2 := eFactory.Create(eClass)
	eObject2.ESet(eAttribute1, 3)
	eObject2.ESet(eAttribute2, "test3")

	list := NewImmutableEList([]interface{}{eObject1, eObject2})
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
	ePackage.GetEClassifiers().AddAll(NewImmutableEList([]interface{}{eClass1, eClass2}))

	eAttribute1 := ecoreFactory.CreateEAttribute()
	eAttribute1.SetName("attribute1")
	eAttribute1.SetEType(ecorePackage.GetEInt())
	eAttribute2 := ecoreFactory.CreateEAttribute()
	eAttribute2.SetName("attribute2")
	eAttribute2.SetEType(ecorePackage.GetEString())
	eClass2.GetEStructuralFeatures().AddAll(NewImmutableEList([]interface{}{eAttribute1, eAttribute2}))

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
	ePackage.GetEClassifiers().AddAll(NewImmutableEList([]interface{}{eClass1, eClass2}))

	eAttribute1 := ecoreFactory.CreateEAttribute()
	eAttribute1.SetName("attribute1")
	eAttribute1.SetEType(ecorePackage.GetEInt())
	eAttribute2 := ecoreFactory.CreateEAttribute()
	eAttribute2.SetName("attribute2")
	eAttribute2.SetEType(ecorePackage.GetEString())
	eClass2.GetEStructuralFeatures().AddAll(NewImmutableEList([]interface{}{eAttribute1, eAttribute2}))

	eReference1 := ecoreFactory.CreateEReference()
	eReference1.SetName("reference1")
	eReference1.SetContainment(true)
	eReference1.SetEType(eClass2)
	eReference2 := ecoreFactory.CreateEReference()
	eReference2.SetName("reference2")
	eReference2.SetContainment(false)
	eReference2.SetEType(eClass2)
	eClass1.GetEStructuralFeatures().AddAll(NewImmutableEList([]interface{}{eReference1, eReference2}))

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
	eObject.(EObjectInternal).ESetProxyURI(&url.URL{Path: "testPath"})

	eObjectCopy := Copy(eObject)
	assert.True(t, Equals(eObject, eObjectCopy))
}

func TestEcoreUtilsCopyReal(t *testing.T) {
	eClass := GetPackage().GetEClass()
	eClassCopy := Copy(eClass)
	assert.True(t, Equals(eClass, eClassCopy))
}
