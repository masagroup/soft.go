// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package xmi

import (
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestXMIDecoderLibrarySimple(t *testing.T) {
	xmiProcessor := NewProcessor()
	resource := xmiProcessor.Load(ecore.NewURI("testdata/library.simple.ecore"))
	require.NotNil(t, resource)
	assert.True(t, resource.IsLoaded())
	assert.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	assert.True(t, resource.GetWarnings().Empty(), diagnosticError(resource.GetWarnings()))

	contents := resource.GetContents()
	assert.Equal(t, 1, contents.Size())

	ePackage, _ := contents.Get(0).(ecore.EPackage)
	assert.NotNil(t, ePackage)
	assert.Equal(t, "library", ePackage.GetName())
	assert.Equal(t, "lib", ePackage.GetNsPrefix())
	assert.Equal(t, "http:///org/eclipse/emf/examples/library/library.simple.ecore/1.0.0", ePackage.GetNsURI())

	eClassifiers := ePackage.GetEClassifiers()
	assert.Equal(t, 2, eClassifiers.Size())

	eLibrary, _ := eClassifiers.Get(0).(ecore.EClassifier)
	assert.NotNil(t, eLibrary)
	assert.Equal(t, "Library", eLibrary.GetName())

	eLibraryClass, _ := eLibrary.(ecore.EClass)
	assert.NotNil(t, eLibraryClass)
	assert.Equal(t, 3, eLibraryClass.GetFeatureCount())

	eOwnerFeature := eLibraryClass.GetEStructuralFeature(0)
	assert.Equal(t, "owner", eOwnerFeature.GetName())
	eOwnerAttribute, _ := eOwnerFeature.(ecore.EAttribute)
	assert.NotNil(t, eOwnerAttribute)
	eOwnerAttributeType := eOwnerAttribute.GetEAttributeType()
	assert.Equal(t, "EString", eOwnerAttributeType.GetName())

	eLocationFeature := eLibraryClass.GetEStructuralFeature(1)
	assert.Equal(t, "location", eLocationFeature.GetName())
	eLocationAttribute, _ := eLocationFeature.(ecore.EAttribute)
	assert.NotNil(t, eLocationAttribute)
	eLocationType := eLocationAttribute.GetEAttributeType()
	assert.NotNil(t, eLocationType)

	eBooksFeature := eLibraryClass.GetEStructuralFeature(2)
	assert.Equal(t, "books", eBooksFeature.GetName())
	eBooksReference, _ := eBooksFeature.(ecore.EReference)
	assert.NotNil(t, eBooksReference)

	eBook := eClassifiers.Get(1).(ecore.EClassifier)
	assert.NotNil(t, eBook)
	assert.Equal(t, "Book", eBook.GetName())
	eBookClass, _ := eBook.(ecore.EClass)
	assert.NotNil(t, eBookClass)
	assert.Equal(t, 2, eBookClass.GetFeatureCount())

	eNameFeature := eBookClass.GetEStructuralFeature(0)
	assert.Equal(t, "name", eNameFeature.GetName())
	eNameAttribute, _ := eNameFeature.(ecore.EAttribute)
	assert.NotNil(t, eNameAttribute)

	eISBNFeature := eBookClass.GetEStructuralFeature(1)
	assert.Equal(t, "isbn", eISBNFeature.GetName())
	eISBNAttribute, _ := eISBNFeature.(ecore.EAttribute)
	assert.NotNil(t, eISBNAttribute)

	// check resolved reference
	assert.Equal(t, eBookClass, eBooksReference.GetEReferenceType())
}

func TestXMIDecoderLibraryNoRoot(t *testing.T) {
	xmiProcessor := NewProcessor()
	resource := xmiProcessor.Load(ecore.NewURI("testdata/library.noroot.ecore"))
	require.NotNil(t, resource)
	assert.True(t, resource.IsLoaded())
	assert.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	assert.True(t, resource.GetWarnings().Empty(), diagnosticError(resource.GetWarnings()))

	contents := resource.GetContents()
	assert.Equal(t, 1, contents.Size())

	ePackage, _ := contents.Get(0).(ecore.EPackage)
	assert.NotNil(t, ePackage)

	eClassifiers := ePackage.GetEClassifiers()
	eBook, _ := eClassifiers.Get(0).(ecore.EClassifier)
	assert.NotNil(t, eBook)
	assert.Equal(t, "Book", eBook.GetName())

	eBookClass, _ := eBook.(ecore.EClass)
	assert.NotNil(t, eBookClass)
	superTypes := eBookClass.GetESuperTypes()
	assert.Equal(t, 1, superTypes.Size())
	eCirculationItemClass := superTypes.Get(0).(ecore.EClass)
	assert.Equal(t, "CirculatingItem", eCirculationItemClass.GetName())

	eWriter, _ := eClassifiers.Get(2).(ecore.EClass)
	assert.NotNil(t, eWriter)
	assert.False(t, eWriter.GetEAnnotations().Empty())
	eAnnotation := eWriter.GetEAnnotation("http://net.masagroup/soft/2019/GenGo")
	assert.NotNil(t, eAnnotation)
	assert.Equal(t, "true", eAnnotation.GetDetails().GetValue("extension"))
}

func TestXMIDecoderLibraryComplex(t *testing.T) {
	xmiProcessor := NewProcessor()
	resource := xmiProcessor.Load(ecore.NewURI("testdata/library.complex.ecore"))
	require.NotNil(t, resource)
	assert.True(t, resource.IsLoaded())
	assert.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	assert.True(t, resource.GetWarnings().Empty(), diagnosticError(resource.GetWarnings()))

	contents := resource.GetContents()
	assert.Equal(t, 1, contents.Size())

	ePackage, _ := contents.Get(0).(ecore.EPackage)
	assert.NotNil(t, ePackage)

	eClassifiers := ePackage.GetEClassifiers()
	eDocumentRootClass, _ := eClassifiers.Get(0).(ecore.EClass)
	assert.NotNil(t, eDocumentRootClass)

	eXMNLSPrefixFeature, _ := eDocumentRootClass.GetEStructuralFeatureFromName("xMLNSPrefixMap").(ecore.EReference)
	assert.NotNil(t, eXMNLSPrefixFeature)

	eType := eXMNLSPrefixFeature.GetEType()
	assert.NotNil(t, eType)
	assert.Equal(t, "EStringToStringMapEntry", eType.GetName())
	assert.False(t, eType.EIsProxy())
}
