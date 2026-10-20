// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package xmlcodec

import (
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExtendedMetatData_GetName(t *testing.T) {
	m := NewExtendedMetaData()
	mockElement := ecore.NewMockENamedElement(t)
	mockFeature := ecore.NewMockEStructuralFeature(t)
	mockAnnotation := ecore.NewMockEAnnotation(t)
	mockDetails := ecore.NewMockEMap(t)

	// no annotations
	mockElement.EXPECT().GetEAnnotation(annotationURI).Return(nil).Once()
	mockElement.EXPECT().GetName().Return("no annotations").Once()
	assert.Equal(t, "no annotations", m.GetName(mockElement))
	assert.Equal(t, "no annotations", m.GetName(mockElement))
	mock.AssertExpectationsForObjects(t, mockElement)

	// annotations
	mockFeature.EXPECT().GetEAnnotation(annotationURI).Return(mockAnnotation).Once()
	mockAnnotation.EXPECT().GetDetails().Return(mockDetails).Once()
	mockDetails.EXPECT().GetValue("name").Return("with annotations").Once()
	assert.Equal(t, "with annotations", m.GetName(mockFeature))
	assert.Equal(t, "with annotations", m.GetName(mockFeature))
	mock.AssertExpectationsForObjects(t, mockFeature, mockAnnotation, mockDetails)
}

func TestExtendedMetatData_GetType(t *testing.T) {
	m := NewExtendedMetaData()
	mockPackage := ecore.NewMockEPackage(t)
	mockClassifier1 := ecore.NewMockEClassifier(t)
	mockClassifier2 := ecore.NewMockEClassifier(t)
	mockAnnotation := ecore.NewMockEAnnotation(t)
	mockDetails := ecore.NewMockEMap(t)
	mockClassifiers := ecore.NewImmutableEList([]any{mockClassifier1, mockClassifier2})

	mockPackage.EXPECT().GetEClassifiers().Return(mockClassifiers).Once()
	mockClassifier1.EXPECT().GetEAnnotation(annotationURI).Return(nil).Once()
	mockClassifier1.EXPECT().GetName().Return("classifier1").Once()
	mockClassifier2.EXPECT().GetEAnnotation(annotationURI).Return(mockAnnotation).Once()
	mockAnnotation.EXPECT().GetDetails().Return(mockDetails).Once()
	mockDetails.EXPECT().GetValue("name").Return("classifier2").Once()

	assert.Equal(t, mockClassifier2, m.GetType(mockPackage, "classifier2"))
	assert.Equal(t, mockClassifier1, m.GetType(mockPackage, "classifier1"))
	mock.AssertExpectationsForObjects(t, mockPackage, mockClassifier1, mockClassifier2, mockAnnotation, mockDetails)
}

func TestExtendedMetatData_GetNamespace(t *testing.T) {
	m := NewExtendedMetaData()
	{
		mockFeature := ecore.NewMockEStructuralFeature(t)
		mockFeature.EXPECT().GetEAnnotation(annotationURI).Return(nil).Once()
		assert.Equal(t, "", m.GetNamespace(mockFeature))
		assert.Equal(t, "", m.GetNamespace(mockFeature))
		mock.AssertExpectationsForObjects(t, mockFeature)
	}
	{
		mockFeature := ecore.NewMockEStructuralFeature(t)
		mockAnnotation := ecore.NewMockEAnnotation(t)
		mockDetails := ecore.NewMockEMap(t)
		mockFeature.EXPECT().GetEAnnotation(annotationURI).Return(mockAnnotation).Once()
		mockAnnotation.EXPECT().GetDetails().Return(mockDetails).Once()
		mockDetails.EXPECT().GetValue("namespace").Return("namespace").Once()
		assert.Equal(t, "namespace", m.GetNamespace(mockFeature))
		assert.Equal(t, "namespace", m.GetNamespace(mockFeature))
		mock.AssertExpectationsForObjects(t, mockFeature, mockAnnotation, mockDetails)
	}
	{
		mockFeature := ecore.NewMockEStructuralFeature(t)
		mockAnnotation := ecore.NewMockEAnnotation(t)
		mockDetails := ecore.NewMockEMap(t)
		mockClass := ecore.NewMockEClass(t)
		mockPackage := ecore.NewMockEPackage(t)
		mockFeature.EXPECT().GetEAnnotation(annotationURI).Return(mockAnnotation).Once()
		mockAnnotation.EXPECT().GetDetails().Return(mockDetails).Once()
		mockDetails.EXPECT().GetValue("namespace").Return("##targetNamespace").Once()
		mockFeature.EXPECT().GetEContainingClass().Return(mockClass).Once()
		mockClass.EXPECT().GetEPackage().Return(mockPackage).Once()
		mockPackage.EXPECT().GetNsURI().Return("uri").Once()
		assert.Equal(t, "uri", m.GetNamespace(mockFeature))
		assert.Equal(t, "uri", m.GetNamespace(mockFeature))
		mock.AssertExpectationsForObjects(t, mockFeature, mockAnnotation, mockDetails, mockClass, mockPackage)
	}
}

func TestExtendedMetatData_GetDocumentRoot(t *testing.T) {
	m := NewExtendedMetaData()
	{
		mockPackage := ecore.NewMockEPackage(t)
		mockClass1 := ecore.NewMockEClass(t)
		mockClass2 := ecore.NewMockEClass(t)
		mockAnnotation := ecore.NewMockEAnnotation(t)
		mockDetails := ecore.NewMockEMap(t)
		mockClassifiers := ecore.NewImmutableEList([]any{mockClass1, mockClass2})
		mockPackage.EXPECT().GetEClassifiers().Return(mockClassifiers).Once()
		mockClass1.EXPECT().GetEAnnotation(annotationURI).Return(nil).Once()
		mockClass1.EXPECT().GetName().Return("classifier1").Once()
		mockClass2.EXPECT().GetEAnnotation(annotationURI).Return(mockAnnotation).Once()
		mockAnnotation.EXPECT().GetDetails().Return(mockDetails).Once()
		mockDetails.EXPECT().GetValue("name").Return("").Once()
		assert.Equal(t, mockClass2, m.GetDocumentRoot(mockPackage))
		mock.AssertExpectationsForObjects(t, mockPackage, mockClass1, mockClass2, mockAnnotation, mockDetails)
	}
	{
		mockPackage := ecore.NewMockEPackage(t)
		mockClass1 := ecore.NewMockEClass(t)
		mockClass2 := ecore.NewMockEClass(t)
		mockClassifiers := ecore.NewImmutableEList([]any{mockClass1, mockClass2})
		mockPackage.EXPECT().GetEClassifiers().Return(mockClassifiers).Once()
		mockClass1.EXPECT().GetEAnnotation(annotationURI).Return(nil).Once()
		mockClass1.EXPECT().GetName().Return("classifier1").Once()
		mockClass2.EXPECT().GetEAnnotation(annotationURI).Return(nil).Once()
		mockClass2.EXPECT().GetName().Return("classifier2").Once()
		assert.Equal(t, nil, m.GetDocumentRoot(mockPackage))
		mock.AssertExpectationsForObjects(t, mockPackage, mockClass1, mockClass2)
	}
}
