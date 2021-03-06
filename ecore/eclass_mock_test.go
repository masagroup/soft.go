// *****************************************************************************
//
// This file is part of a MASA library or program.
// Refer to the included end-user license agreement for restrictions.
//
// Copyright (c) 2020 MASA Group
//
// *****************************************************************************

// *****************************************************************************
//
// Warning: This file was generated by soft.generator.go Generator
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
	o := &MockEClass{}
	r := bool(true)
	o.On("IsAbstract").Once().Return(r)
	o.On("IsAbstract").Once().Return(func() bool {
		return r
	})
	assert.Equal(t, r, o.IsAbstract())
	assert.Equal(t, r, o.IsAbstract())
	o.AssertExpectations(t)
}

// TestMockEClassSetAbstract tests method SetAbstract
func TestMockEClassSetAbstract(t *testing.T) {
	o := &MockEClass{}
	v := bool(true)
	o.On("SetAbstract", v).Once()
	o.SetAbstract(v)
	o.AssertExpectations(t)
}

// TestMockEClassGetEAllAttributes tests method GetEAllAttributes
func TestMockEClassGetEAllAttributes(t *testing.T) {
	o := &MockEClass{}
	l := &MockEList{}
	// return a value
	o.On("GetEAllAttributes").Once().Return(l)
	o.On("GetEAllAttributes").Once().Return(func() EList {
		return l
	})
	assert.Equal(t, l, o.GetEAllAttributes())
	assert.Equal(t, l, o.GetEAllAttributes())
	o.AssertExpectations(t)
}

// TestMockEClassGetEAllContainments tests method GetEAllContainments
func TestMockEClassGetEAllContainments(t *testing.T) {
	o := &MockEClass{}
	l := &MockEList{}
	// return a value
	o.On("GetEAllContainments").Once().Return(l)
	o.On("GetEAllContainments").Once().Return(func() EList {
		return l
	})
	assert.Equal(t, l, o.GetEAllContainments())
	assert.Equal(t, l, o.GetEAllContainments())
	o.AssertExpectations(t)
}

// TestMockEClassGetEAllOperations tests method GetEAllOperations
func TestMockEClassGetEAllOperations(t *testing.T) {
	o := &MockEClass{}
	l := &MockEList{}
	// return a value
	o.On("GetEAllOperations").Once().Return(l)
	o.On("GetEAllOperations").Once().Return(func() EList {
		return l
	})
	assert.Equal(t, l, o.GetEAllOperations())
	assert.Equal(t, l, o.GetEAllOperations())
	o.AssertExpectations(t)
}

// TestMockEClassGetEAllReferences tests method GetEAllReferences
func TestMockEClassGetEAllReferences(t *testing.T) {
	o := &MockEClass{}
	l := &MockEList{}
	// return a value
	o.On("GetEAllReferences").Once().Return(l)
	o.On("GetEAllReferences").Once().Return(func() EList {
		return l
	})
	assert.Equal(t, l, o.GetEAllReferences())
	assert.Equal(t, l, o.GetEAllReferences())
	o.AssertExpectations(t)
}

// TestMockEClassGetEAllStructuralFeatures tests method GetEAllStructuralFeatures
func TestMockEClassGetEAllStructuralFeatures(t *testing.T) {
	o := &MockEClass{}
	l := &MockEList{}
	// return a value
	o.On("GetEAllStructuralFeatures").Once().Return(l)
	o.On("GetEAllStructuralFeatures").Once().Return(func() EList {
		return l
	})
	assert.Equal(t, l, o.GetEAllStructuralFeatures())
	assert.Equal(t, l, o.GetEAllStructuralFeatures())
	o.AssertExpectations(t)
}

// TestMockEClassGetEAllSuperTypes tests method GetEAllSuperTypes
func TestMockEClassGetEAllSuperTypes(t *testing.T) {
	o := &MockEClass{}
	l := &MockEList{}
	// return a value
	o.On("GetEAllSuperTypes").Once().Return(l)
	o.On("GetEAllSuperTypes").Once().Return(func() EList {
		return l
	})
	assert.Equal(t, l, o.GetEAllSuperTypes())
	assert.Equal(t, l, o.GetEAllSuperTypes())
	o.AssertExpectations(t)
}

// TestMockEClassGetEAttributes tests method GetEAttributes
func TestMockEClassGetEAttributes(t *testing.T) {
	o := &MockEClass{}
	l := &MockEList{}
	// return a value
	o.On("GetEAttributes").Once().Return(l)
	o.On("GetEAttributes").Once().Return(func() EList {
		return l
	})
	assert.Equal(t, l, o.GetEAttributes())
	assert.Equal(t, l, o.GetEAttributes())
	o.AssertExpectations(t)
}

// TestMockEClassGetEContainmentFeatures tests method GetEContainmentFeatures
func TestMockEClassGetEContainmentFeatures(t *testing.T) {
	o := &MockEClass{}
	l := &MockEList{}
	// return a value
	o.On("GetEContainmentFeatures").Once().Return(l)
	o.On("GetEContainmentFeatures").Once().Return(func() EList {
		return l
	})
	assert.Equal(t, l, o.GetEContainmentFeatures())
	assert.Equal(t, l, o.GetEContainmentFeatures())
	o.AssertExpectations(t)
}

// TestMockEClassGetECrossReferenceFeatures tests method GetECrossReferenceFeatures
func TestMockEClassGetECrossReferenceFeatures(t *testing.T) {
	o := &MockEClass{}
	l := &MockEList{}
	// return a value
	o.On("GetECrossReferenceFeatures").Once().Return(l)
	o.On("GetECrossReferenceFeatures").Once().Return(func() EList {
		return l
	})
	assert.Equal(t, l, o.GetECrossReferenceFeatures())
	assert.Equal(t, l, o.GetECrossReferenceFeatures())
	o.AssertExpectations(t)
}

// TestMockEClassGetEIDAttribute tests method GetEIDAttribute
func TestMockEClassGetEIDAttribute(t *testing.T) {
	o := &MockEClass{}
	r := new(MockEAttribute)
	o.On("GetEIDAttribute").Once().Return(r)
	o.On("GetEIDAttribute").Once().Return(func() EAttribute {
		return r
	})
	assert.Equal(t, r, o.GetEIDAttribute())
	assert.Equal(t, r, o.GetEIDAttribute())
	o.AssertExpectations(t)
}

// TestMockEClassGetEOperations tests method GetEOperations
func TestMockEClassGetEOperations(t *testing.T) {
	o := &MockEClass{}
	l := &MockEList{}
	// return a value
	o.On("GetEOperations").Once().Return(l)
	o.On("GetEOperations").Once().Return(func() EList {
		return l
	})
	assert.Equal(t, l, o.GetEOperations())
	assert.Equal(t, l, o.GetEOperations())
	o.AssertExpectations(t)
}

// TestMockEClassGetEReferences tests method GetEReferences
func TestMockEClassGetEReferences(t *testing.T) {
	o := &MockEClass{}
	l := &MockEList{}
	// return a value
	o.On("GetEReferences").Once().Return(l)
	o.On("GetEReferences").Once().Return(func() EList {
		return l
	})
	assert.Equal(t, l, o.GetEReferences())
	assert.Equal(t, l, o.GetEReferences())
	o.AssertExpectations(t)
}

// TestMockEClassGetEStructuralFeatures tests method GetEStructuralFeatures
func TestMockEClassGetEStructuralFeatures(t *testing.T) {
	o := &MockEClass{}
	l := &MockEList{}
	// return a value
	o.On("GetEStructuralFeatures").Once().Return(l)
	o.On("GetEStructuralFeatures").Once().Return(func() EList {
		return l
	})
	assert.Equal(t, l, o.GetEStructuralFeatures())
	assert.Equal(t, l, o.GetEStructuralFeatures())
	o.AssertExpectations(t)
}

// TestMockEClassGetESuperTypes tests method GetESuperTypes
func TestMockEClassGetESuperTypes(t *testing.T) {
	o := &MockEClass{}
	l := &MockEList{}
	// return a value
	o.On("GetESuperTypes").Once().Return(l)
	o.On("GetESuperTypes").Once().Return(func() EList {
		return l
	})
	assert.Equal(t, l, o.GetESuperTypes())
	assert.Equal(t, l, o.GetESuperTypes())
	o.AssertExpectations(t)
}

// TestMockEClassIsInterface tests method IsInterface
func TestMockEClassIsInterface(t *testing.T) {
	o := &MockEClass{}
	r := bool(true)
	o.On("IsInterface").Once().Return(r)
	o.On("IsInterface").Once().Return(func() bool {
		return r
	})
	assert.Equal(t, r, o.IsInterface())
	assert.Equal(t, r, o.IsInterface())
	o.AssertExpectations(t)
}

// TestMockEClassSetInterface tests method SetInterface
func TestMockEClassSetInterface(t *testing.T) {
	o := &MockEClass{}
	v := bool(true)
	o.On("SetInterface", v).Once()
	o.SetInterface(v)
	o.AssertExpectations(t)
}

// TestMockEClassGetEOperation tests method GetEOperation
func TestMockEClassGetEOperation(t *testing.T) {
	o := &MockEClass{}
	operationID := int(45)
	r := new(MockEOperation)
	o.On("GetEOperation", operationID).Return(r).Once()
	o.On("GetEOperation", operationID).Return(func() EOperation {
		return r
	}).Once()
	assert.Equal(t, r, o.GetEOperation(operationID))
	assert.Equal(t, r, o.GetEOperation(operationID))
	o.AssertExpectations(t)
}

// TestMockEClassGetEStructuralFeature tests method GetEStructuralFeature
func TestMockEClassGetEStructuralFeature(t *testing.T) {
	o := &MockEClass{}
	featureID := int(45)
	r := new(MockEStructuralFeature)
	o.On("GetEStructuralFeature", featureID).Return(r).Once()
	o.On("GetEStructuralFeature", featureID).Return(func() EStructuralFeature {
		return r
	}).Once()
	assert.Equal(t, r, o.GetEStructuralFeature(featureID))
	assert.Equal(t, r, o.GetEStructuralFeature(featureID))
	o.AssertExpectations(t)
}

// TestMockEClassGetEStructuralFeatureFromName tests method GetEStructuralFeatureFromName
func TestMockEClassGetEStructuralFeatureFromName(t *testing.T) {
	o := &MockEClass{}
	featureName := string("Test String")
	r := new(MockEStructuralFeature)
	o.On("GetEStructuralFeatureFromName", featureName).Return(r).Once()
	o.On("GetEStructuralFeatureFromName", featureName).Return(func() EStructuralFeature {
		return r
	}).Once()
	assert.Equal(t, r, o.GetEStructuralFeatureFromName(featureName))
	assert.Equal(t, r, o.GetEStructuralFeatureFromName(featureName))
	o.AssertExpectations(t)
}

// TestMockEClassGetFeatureCount tests method GetFeatureCount
func TestMockEClassGetFeatureCount(t *testing.T) {
	o := &MockEClass{}
	r := int(45)
	o.On("GetFeatureCount").Return(r).Once()
	o.On("GetFeatureCount").Return(func() int {
		return r
	}).Once()
	assert.Equal(t, r, o.GetFeatureCount())
	assert.Equal(t, r, o.GetFeatureCount())
	o.AssertExpectations(t)
}

// TestMockEClassGetFeatureID tests method GetFeatureID
func TestMockEClassGetFeatureID(t *testing.T) {
	o := &MockEClass{}
	feature := new(MockEStructuralFeature)
	r := int(45)
	o.On("GetFeatureID", feature).Return(r).Once()
	o.On("GetFeatureID", feature).Return(func() int {
		return r
	}).Once()
	assert.Equal(t, r, o.GetFeatureID(feature))
	assert.Equal(t, r, o.GetFeatureID(feature))
	o.AssertExpectations(t)
}

// TestMockEClassGetFeatureType tests method GetFeatureType
func TestMockEClassGetFeatureType(t *testing.T) {
	o := &MockEClass{}
	feature := new(MockEStructuralFeature)
	r := new(MockEClassifier)
	o.On("GetFeatureType", feature).Return(r).Once()
	o.On("GetFeatureType", feature).Return(func() EClassifier {
		return r
	}).Once()
	assert.Equal(t, r, o.GetFeatureType(feature))
	assert.Equal(t, r, o.GetFeatureType(feature))
	o.AssertExpectations(t)
}

// TestMockEClassGetOperationCount tests method GetOperationCount
func TestMockEClassGetOperationCount(t *testing.T) {
	o := &MockEClass{}
	r := int(45)
	o.On("GetOperationCount").Return(r).Once()
	o.On("GetOperationCount").Return(func() int {
		return r
	}).Once()
	assert.Equal(t, r, o.GetOperationCount())
	assert.Equal(t, r, o.GetOperationCount())
	o.AssertExpectations(t)
}

// TestMockEClassGetOperationID tests method GetOperationID
func TestMockEClassGetOperationID(t *testing.T) {
	o := &MockEClass{}
	operation := new(MockEOperation)
	r := int(45)
	o.On("GetOperationID", operation).Return(r).Once()
	o.On("GetOperationID", operation).Return(func() int {
		return r
	}).Once()
	assert.Equal(t, r, o.GetOperationID(operation))
	assert.Equal(t, r, o.GetOperationID(operation))
	o.AssertExpectations(t)
}

// TestMockEClassGetOverride tests method GetOverride
func TestMockEClassGetOverride(t *testing.T) {
	o := &MockEClass{}
	operation := new(MockEOperation)
	r := new(MockEOperation)
	o.On("GetOverride", operation).Return(r).Once()
	o.On("GetOverride", operation).Return(func() EOperation {
		return r
	}).Once()
	assert.Equal(t, r, o.GetOverride(operation))
	assert.Equal(t, r, o.GetOverride(operation))
	o.AssertExpectations(t)
}

// TestMockEClassIsSuperTypeOf tests method IsSuperTypeOf
func TestMockEClassIsSuperTypeOf(t *testing.T) {
	o := &MockEClass{}
	someClass := new(MockEClass)
	r := bool(true)
	o.On("IsSuperTypeOf", someClass).Return(r).Once()
	o.On("IsSuperTypeOf", someClass).Return(func() bool {
		return r
	}).Once()
	assert.Equal(t, r, o.IsSuperTypeOf(someClass))
	assert.Equal(t, r, o.IsSuperTypeOf(someClass))
	o.AssertExpectations(t)
}
