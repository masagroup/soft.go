// Code generated by soft.generator.go. DO NOT EDIT.

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
	"github.com/stretchr/testify/assert"
	"testing"
)

func discardMockEClass() {
	_ = assert.Equal
	_ = testing.Coverage
}

// TestMockEClassIsAbstract tests method IsAbstract
func TestMockEClassIsAbstract(t *testing.T) {
	o := NewMockEClass(t)
	r := bool(true)
	m := NewMockRun(t)
	o.EXPECT().IsAbstract().Return(r).Run(func() { m.Run() }).Once()
	o.EXPECT().IsAbstract().Call.Return(func() bool { return r }).Once()
	assert.Equal(t, r, o.IsAbstract())
	assert.Equal(t, r, o.IsAbstract())
}

// TestMockEClassSetAbstract tests method SetAbstract
func TestMockEClassSetAbstract(t *testing.T) {
	o := NewMockEClass(t)
	v := bool(true)
	m := NewMockRun(t, v)
	o.EXPECT().SetAbstract(v).Return().Run(func(_p0 bool) { m.Run(_p0) }).Once()
	o.SetAbstract(v)
}

// TestMockEClassGetEAllAttributes tests method GetEAllAttributes
func TestMockEClassGetEAllAttributes(t *testing.T) {
	o := NewMockEClass(t)
	l := NewMockEList(t)
	m := NewMockRun(t)
	o.EXPECT().GetEAllAttributes().Return(l).Run(func() { m.Run() }).Once()
	o.EXPECT().GetEAllAttributes().Call.Return(func() EList { return l }).Once()
	assert.Equal(t, l, o.GetEAllAttributes())
	assert.Equal(t, l, o.GetEAllAttributes())
}

// TestMockEClassGetEAllContainments tests method GetEAllContainments
func TestMockEClassGetEAllContainments(t *testing.T) {
	o := NewMockEClass(t)
	l := NewMockEList(t)
	m := NewMockRun(t)
	o.EXPECT().GetEAllContainments().Return(l).Run(func() { m.Run() }).Once()
	o.EXPECT().GetEAllContainments().Call.Return(func() EList { return l }).Once()
	assert.Equal(t, l, o.GetEAllContainments())
	assert.Equal(t, l, o.GetEAllContainments())
}

// TestMockEClassGetEAllCrossReferences tests method GetEAllCrossReferences
func TestMockEClassGetEAllCrossReferences(t *testing.T) {
	o := NewMockEClass(t)
	l := NewMockEList(t)
	m := NewMockRun(t)
	o.EXPECT().GetEAllCrossReferences().Return(l).Run(func() { m.Run() }).Once()
	o.EXPECT().GetEAllCrossReferences().Call.Return(func() EList { return l }).Once()
	assert.Equal(t, l, o.GetEAllCrossReferences())
	assert.Equal(t, l, o.GetEAllCrossReferences())
}

// TestMockEClassGetEAllOperations tests method GetEAllOperations
func TestMockEClassGetEAllOperations(t *testing.T) {
	o := NewMockEClass(t)
	l := NewMockEList(t)
	m := NewMockRun(t)
	o.EXPECT().GetEAllOperations().Return(l).Run(func() { m.Run() }).Once()
	o.EXPECT().GetEAllOperations().Call.Return(func() EList { return l }).Once()
	assert.Equal(t, l, o.GetEAllOperations())
	assert.Equal(t, l, o.GetEAllOperations())
}

// TestMockEClassGetEAllReferences tests method GetEAllReferences
func TestMockEClassGetEAllReferences(t *testing.T) {
	o := NewMockEClass(t)
	l := NewMockEList(t)
	m := NewMockRun(t)
	o.EXPECT().GetEAllReferences().Return(l).Run(func() { m.Run() }).Once()
	o.EXPECT().GetEAllReferences().Call.Return(func() EList { return l }).Once()
	assert.Equal(t, l, o.GetEAllReferences())
	assert.Equal(t, l, o.GetEAllReferences())
}

// TestMockEClassGetEAllStructuralFeatures tests method GetEAllStructuralFeatures
func TestMockEClassGetEAllStructuralFeatures(t *testing.T) {
	o := NewMockEClass(t)
	l := NewMockEList(t)
	m := NewMockRun(t)
	o.EXPECT().GetEAllStructuralFeatures().Return(l).Run(func() { m.Run() }).Once()
	o.EXPECT().GetEAllStructuralFeatures().Call.Return(func() EList { return l }).Once()
	assert.Equal(t, l, o.GetEAllStructuralFeatures())
	assert.Equal(t, l, o.GetEAllStructuralFeatures())
}

// TestMockEClassGetEAllSuperTypes tests method GetEAllSuperTypes
func TestMockEClassGetEAllSuperTypes(t *testing.T) {
	o := NewMockEClass(t)
	l := NewMockEList(t)
	m := NewMockRun(t)
	o.EXPECT().GetEAllSuperTypes().Return(l).Run(func() { m.Run() }).Once()
	o.EXPECT().GetEAllSuperTypes().Call.Return(func() EList { return l }).Once()
	assert.Equal(t, l, o.GetEAllSuperTypes())
	assert.Equal(t, l, o.GetEAllSuperTypes())
}

// TestMockEClassGetEAttributes tests method GetEAttributes
func TestMockEClassGetEAttributes(t *testing.T) {
	o := NewMockEClass(t)
	l := NewMockEList(t)
	m := NewMockRun(t)
	o.EXPECT().GetEAttributes().Return(l).Run(func() { m.Run() }).Once()
	o.EXPECT().GetEAttributes().Call.Return(func() EList { return l }).Once()
	assert.Equal(t, l, o.GetEAttributes())
	assert.Equal(t, l, o.GetEAttributes())
}

// TestMockEClassGetEContainmentFeatures tests method GetEContainmentFeatures
func TestMockEClassGetEContainmentFeatures(t *testing.T) {
	o := NewMockEClass(t)
	l := NewMockEList(t)
	m := NewMockRun(t)
	o.EXPECT().GetEContainmentFeatures().Return(l).Run(func() { m.Run() }).Once()
	o.EXPECT().GetEContainmentFeatures().Call.Return(func() EList { return l }).Once()
	assert.Equal(t, l, o.GetEContainmentFeatures())
	assert.Equal(t, l, o.GetEContainmentFeatures())
}

// TestMockEClassGetECrossReferenceFeatures tests method GetECrossReferenceFeatures
func TestMockEClassGetECrossReferenceFeatures(t *testing.T) {
	o := NewMockEClass(t)
	l := NewMockEList(t)
	m := NewMockRun(t)
	o.EXPECT().GetECrossReferenceFeatures().Return(l).Run(func() { m.Run() }).Once()
	o.EXPECT().GetECrossReferenceFeatures().Call.Return(func() EList { return l }).Once()
	assert.Equal(t, l, o.GetECrossReferenceFeatures())
	assert.Equal(t, l, o.GetECrossReferenceFeatures())
}

// TestMockEClassGetEIDAttribute tests method GetEIDAttribute
func TestMockEClassGetEIDAttribute(t *testing.T) {
	o := NewMockEClass(t)
	r := NewMockEAttribute(t)
	m := NewMockRun(t)
	o.EXPECT().GetEIDAttribute().Return(r).Run(func() { m.Run() }).Once()
	o.EXPECT().GetEIDAttribute().Call.Return(func() EAttribute { return r }).Once()
	assert.Equal(t, r, o.GetEIDAttribute())
	assert.Equal(t, r, o.GetEIDAttribute())
}

// TestMockEClassGetEOperations tests method GetEOperations
func TestMockEClassGetEOperations(t *testing.T) {
	o := NewMockEClass(t)
	l := NewMockEList(t)
	m := NewMockRun(t)
	o.EXPECT().GetEOperations().Return(l).Run(func() { m.Run() }).Once()
	o.EXPECT().GetEOperations().Call.Return(func() EList { return l }).Once()
	assert.Equal(t, l, o.GetEOperations())
	assert.Equal(t, l, o.GetEOperations())
}

// TestMockEClassGetEReferences tests method GetEReferences
func TestMockEClassGetEReferences(t *testing.T) {
	o := NewMockEClass(t)
	l := NewMockEList(t)
	m := NewMockRun(t)
	o.EXPECT().GetEReferences().Return(l).Run(func() { m.Run() }).Once()
	o.EXPECT().GetEReferences().Call.Return(func() EList { return l }).Once()
	assert.Equal(t, l, o.GetEReferences())
	assert.Equal(t, l, o.GetEReferences())
}

// TestMockEClassGetEStructuralFeatures tests method GetEStructuralFeatures
func TestMockEClassGetEStructuralFeatures(t *testing.T) {
	o := NewMockEClass(t)
	l := NewMockEList(t)
	m := NewMockRun(t)
	o.EXPECT().GetEStructuralFeatures().Return(l).Run(func() { m.Run() }).Once()
	o.EXPECT().GetEStructuralFeatures().Call.Return(func() EList { return l }).Once()
	assert.Equal(t, l, o.GetEStructuralFeatures())
	assert.Equal(t, l, o.GetEStructuralFeatures())
}

// TestMockEClassGetESuperTypes tests method GetESuperTypes
func TestMockEClassGetESuperTypes(t *testing.T) {
	o := NewMockEClass(t)
	l := NewMockEList(t)
	m := NewMockRun(t)
	o.EXPECT().GetESuperTypes().Return(l).Run(func() { m.Run() }).Once()
	o.EXPECT().GetESuperTypes().Call.Return(func() EList { return l }).Once()
	assert.Equal(t, l, o.GetESuperTypes())
	assert.Equal(t, l, o.GetESuperTypes())
}

// TestMockEClassIsInterface tests method IsInterface
func TestMockEClassIsInterface(t *testing.T) {
	o := NewMockEClass(t)
	r := bool(true)
	m := NewMockRun(t)
	o.EXPECT().IsInterface().Return(r).Run(func() { m.Run() }).Once()
	o.EXPECT().IsInterface().Call.Return(func() bool { return r }).Once()
	assert.Equal(t, r, o.IsInterface())
	assert.Equal(t, r, o.IsInterface())
}

// TestMockEClassSetInterface tests method SetInterface
func TestMockEClassSetInterface(t *testing.T) {
	o := NewMockEClass(t)
	v := bool(true)
	m := NewMockRun(t, v)
	o.EXPECT().SetInterface(v).Return().Run(func(_p0 bool) { m.Run(_p0) }).Once()
	o.SetInterface(v)
}

// TestMockEClassGetEOperation tests method GetEOperation
func TestMockEClassGetEOperation(t *testing.T) {
	o := NewMockEClass(t)
	operationID := int(45)
	m := NewMockRun(t, operationID)
	r := NewMockEOperation(t)
	o.EXPECT().GetEOperation(operationID).Return(r).Run(func(operationID int) { m.Run(operationID) }).Once()
	o.EXPECT().GetEOperation(operationID).Call.Return(func() EOperation {
		return r
	}).Once()
	assert.Equal(t, r, o.GetEOperation(operationID))
	assert.Equal(t, r, o.GetEOperation(operationID))
}

// TestMockEClassGetEStructuralFeature tests method GetEStructuralFeature
func TestMockEClassGetEStructuralFeature(t *testing.T) {
	o := NewMockEClass(t)
	featureID := int(45)
	m := NewMockRun(t, featureID)
	r := NewMockEStructuralFeature(t)
	o.EXPECT().GetEStructuralFeature(featureID).Return(r).Run(func(featureID int) { m.Run(featureID) }).Once()
	o.EXPECT().GetEStructuralFeature(featureID).Call.Return(func() EStructuralFeature {
		return r
	}).Once()
	assert.Equal(t, r, o.GetEStructuralFeature(featureID))
	assert.Equal(t, r, o.GetEStructuralFeature(featureID))
}

// TestMockEClassGetEStructuralFeatureFromName tests method GetEStructuralFeatureFromName
func TestMockEClassGetEStructuralFeatureFromName(t *testing.T) {
	o := NewMockEClass(t)
	featureName := string("Test String")
	m := NewMockRun(t, featureName)
	r := NewMockEStructuralFeature(t)
	o.EXPECT().GetEStructuralFeatureFromName(featureName).Return(r).Run(func(featureName string) { m.Run(featureName) }).Once()
	o.EXPECT().GetEStructuralFeatureFromName(featureName).Call.Return(func() EStructuralFeature {
		return r
	}).Once()
	assert.Equal(t, r, o.GetEStructuralFeatureFromName(featureName))
	assert.Equal(t, r, o.GetEStructuralFeatureFromName(featureName))
}

// TestMockEClassGetFeatureCount tests method GetFeatureCount
func TestMockEClassGetFeatureCount(t *testing.T) {
	o := NewMockEClass(t)
	m := NewMockRun(t)
	r := int(45)
	o.EXPECT().GetFeatureCount().Return(r).Run(func() { m.Run() }).Once()
	o.EXPECT().GetFeatureCount().Call.Return(func() int {
		return r
	}).Once()
	assert.Equal(t, r, o.GetFeatureCount())
	assert.Equal(t, r, o.GetFeatureCount())
}

// TestMockEClassGetFeatureID tests method GetFeatureID
func TestMockEClassGetFeatureID(t *testing.T) {
	o := NewMockEClass(t)
	feature := NewMockEStructuralFeature(t)
	m := NewMockRun(t, feature)
	r := int(45)
	o.EXPECT().GetFeatureID(feature).Return(r).Run(func(feature EStructuralFeature) { m.Run(feature) }).Once()
	o.EXPECT().GetFeatureID(feature).Call.Return(func() int {
		return r
	}).Once()
	assert.Equal(t, r, o.GetFeatureID(feature))
	assert.Equal(t, r, o.GetFeatureID(feature))
}

// TestMockEClassGetFeatureType tests method GetFeatureType
func TestMockEClassGetFeatureType(t *testing.T) {
	o := NewMockEClass(t)
	feature := NewMockEStructuralFeature(t)
	m := NewMockRun(t, feature)
	r := NewMockEClassifier(t)
	o.EXPECT().GetFeatureType(feature).Return(r).Run(func(feature EStructuralFeature) { m.Run(feature) }).Once()
	o.EXPECT().GetFeatureType(feature).Call.Return(func() EClassifier {
		return r
	}).Once()
	assert.Equal(t, r, o.GetFeatureType(feature))
	assert.Equal(t, r, o.GetFeatureType(feature))
}

// TestMockEClassGetOperationCount tests method GetOperationCount
func TestMockEClassGetOperationCount(t *testing.T) {
	o := NewMockEClass(t)
	m := NewMockRun(t)
	r := int(45)
	o.EXPECT().GetOperationCount().Return(r).Run(func() { m.Run() }).Once()
	o.EXPECT().GetOperationCount().Call.Return(func() int {
		return r
	}).Once()
	assert.Equal(t, r, o.GetOperationCount())
	assert.Equal(t, r, o.GetOperationCount())
}

// TestMockEClassGetOperationID tests method GetOperationID
func TestMockEClassGetOperationID(t *testing.T) {
	o := NewMockEClass(t)
	operation := NewMockEOperation(t)
	m := NewMockRun(t, operation)
	r := int(45)
	o.EXPECT().GetOperationID(operation).Return(r).Run(func(operation EOperation) { m.Run(operation) }).Once()
	o.EXPECT().GetOperationID(operation).Call.Return(func() int {
		return r
	}).Once()
	assert.Equal(t, r, o.GetOperationID(operation))
	assert.Equal(t, r, o.GetOperationID(operation))
}

// TestMockEClassGetOverride tests method GetOverride
func TestMockEClassGetOverride(t *testing.T) {
	o := NewMockEClass(t)
	operation := NewMockEOperation(t)
	m := NewMockRun(t, operation)
	r := NewMockEOperation(t)
	o.EXPECT().GetOverride(operation).Return(r).Run(func(operation EOperation) { m.Run(operation) }).Once()
	o.EXPECT().GetOverride(operation).Call.Return(func() EOperation {
		return r
	}).Once()
	assert.Equal(t, r, o.GetOverride(operation))
	assert.Equal(t, r, o.GetOverride(operation))
}

// TestMockEClassIsSuperTypeOf tests method IsSuperTypeOf
func TestMockEClassIsSuperTypeOf(t *testing.T) {
	o := NewMockEClass(t)
	someClass := NewMockEClass(t)
	m := NewMockRun(t, someClass)
	r := bool(true)
	o.EXPECT().IsSuperTypeOf(someClass).Return(r).Run(func(someClass EClass) { m.Run(someClass) }).Once()
	o.EXPECT().IsSuperTypeOf(someClass).Call.Return(func() bool {
		return r
	}).Once()
	assert.Equal(t, r, o.IsSuperTypeOf(someClass))
	assert.Equal(t, r, o.IsSuperTypeOf(someClass))
}
