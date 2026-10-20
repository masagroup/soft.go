// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package xmlcodec

import "github.com/masagroup/soft.go/ecore"

const (
	annotationURI = "http:///org/eclipse/emf/ecore/util/ExtendedMetaData"
)

type ExtendedMetaData struct {
	metaData map[any]any
}

type ENamedElementExtendedMetaData interface {
	getName() string
}

type ENamedElementExtendedMetaDataImpl struct {
	emd      *ExtendedMetaData
	eElement ecore.ENamedElement
	name     string
}

func (m *ENamedElementExtendedMetaDataImpl) getName() string {
	if m.name == "unitialized" {
		m.name = m.emd.basicGetName(m.eElement)
	}
	return m.name
}

type EPackageExtentedMetaData interface {
	getType(string) ecore.EClassifier
}

type EPackageExtentedMetaDataImpl struct {
	emd                 *ExtendedMetaData
	ePackage            ecore.EPackage
	nameToClassifierMap map[string]ecore.EClassifier
}

func (m *EPackageExtentedMetaDataImpl) getType(name string) ecore.EClassifier {
	var eResult ecore.EClassifier = nil
	if m.nameToClassifierMap != nil {
		eResult = m.nameToClassifierMap[name]
	}
	if eResult == nil {
		eClassifiers := m.ePackage.GetEClassifiers()
		if m.nameToClassifierMap == nil || len(m.nameToClassifierMap) != eClassifiers.Size() {
			nameToClassifierMap := make(map[string]ecore.EClassifier)
			for it := eClassifiers.Iterator(); it.HasNext(); {
				eClassifier := it.Next().(ecore.EClassifier)
				eClassifierName := m.emd.GetName(eClassifier)
				nameToClassifierMap[eClassifierName] = eClassifier
				if eClassifierName == name {
					eResult = eClassifier
					break
				}
			}
			m.nameToClassifierMap = nameToClassifierMap
		}
	}
	return eResult
}

type EStructuralFeatureExtentedMetaData interface {
	ENamedElementExtendedMetaData
	getNamespace() string
}

type EStructuralFeatureExtentedMetaDataImpl struct {
	ENamedElementExtendedMetaDataImpl
	namespace string
}

func (m *EStructuralFeatureExtentedMetaDataImpl) getNamespace() string {
	if m.namespace == "unitialized" {
		m.namespace = m.emd.basicGetNamespace(m.eElement.(ecore.EStructuralFeature))
	}
	return m.namespace
}

func NewExtendedMetaData() *ExtendedMetaData {
	return &ExtendedMetaData{metaData: make(map[any]any)}
}

func (emd *ExtendedMetaData) getENamedElementExtendedMetaData(eElement ecore.ENamedElement) ENamedElementExtendedMetaData {
	var result ENamedElementExtendedMetaData
	if i, ok := emd.metaData[eElement]; ok {
		result = i.(ENamedElementExtendedMetaData)
	} else {
		if eFeature, _ := eElement.(ecore.EStructuralFeature); eFeature != nil {
			result = &EStructuralFeatureExtentedMetaDataImpl{ENamedElementExtendedMetaDataImpl: ENamedElementExtendedMetaDataImpl{emd: emd, eElement: eElement, name: "unitialized"}, namespace: "unitialized"}
		} else {
			result = &ENamedElementExtendedMetaDataImpl{emd: emd, eElement: eElement, name: "unitialized"}
		}
		emd.metaData[eElement] = result
	}
	return result
}

func (emd *ExtendedMetaData) getEStructuralFeatureExtentedMetaData(eFeature ecore.EStructuralFeature) EStructuralFeatureExtentedMetaData {
	var result EStructuralFeatureExtentedMetaData
	if i, ok := emd.metaData[eFeature]; ok {
		result = i.(EStructuralFeatureExtentedMetaData)
	} else {
		result = &EStructuralFeatureExtentedMetaDataImpl{ENamedElementExtendedMetaDataImpl: ENamedElementExtendedMetaDataImpl{emd: emd, eElement: eFeature, name: "unitialized"}, namespace: "unitialized"}
		emd.metaData[eFeature] = result
	}
	return result
}

func (emd *ExtendedMetaData) getEPackageExtentedMetaData(ePackage ecore.EPackage) EPackageExtentedMetaData {
	var result EPackageExtentedMetaData
	if i, ok := emd.metaData[ePackage]; ok {
		result = i.(EPackageExtentedMetaData)
	} else {
		result = &EPackageExtentedMetaDataImpl{emd: emd, ePackage: ePackage}
		emd.metaData[ePackage] = result
	}
	return result
}

func (emd *ExtendedMetaData) GetType(ePackage ecore.EPackage, name string) ecore.EClassifier {
	return emd.getEPackageExtentedMetaData(ePackage).getType(name)
}

func (emd *ExtendedMetaData) GetName(eElement ecore.ENamedElement) string {
	return emd.getENamedElementExtendedMetaData(eElement).getName()
}

func (emd *ExtendedMetaData) GetNamespace(eFeature ecore.EStructuralFeature) string {
	return emd.getEStructuralFeatureExtentedMetaData(eFeature).getNamespace()
}

func (emd *ExtendedMetaData) GetDocumentRoot(ePackage ecore.EPackage) ecore.EClass {
	eClassifier := emd.GetType(ePackage, "")
	if eClassifier != nil {
		return eClassifier.(ecore.EClass)
	}
	return nil
}

func (emd *ExtendedMetaData) GetXMLNSPrefixMapFeature(eClass ecore.EClass) ecore.EReference {
	for it := eClass.GetEAllReferences().Iterator(); it.HasNext(); {
		eReference := it.Next().(ecore.EReference)
		if emd.GetName(eReference) == "xmlns:prefix" {
			return eReference
		}
	}
	return nil
}

func (emd *ExtendedMetaData) GetXSISchemaLocationMapFeature(eClass ecore.EClass) ecore.EReference {
	for it := eClass.GetEAllReferences().Iterator(); it.HasNext(); {
		eReference := it.Next().(ecore.EReference)
		if emd.GetName(eReference) == "xsi:schemaLocation" {
			return eReference
		}
	}
	return nil
}

func (emd *ExtendedMetaData) basicGetName(eElement ecore.ENamedElement) string {
	if annotation := eElement.GetEAnnotation(annotationURI); annotation != nil {
		if name := annotation.GetDetails().GetValue("name"); name != nil {
			return name.(string)
		}
	}
	return eElement.GetName()
}

func (emd *ExtendedMetaData) basicGetNamespace(eFeature ecore.EStructuralFeature) string {
	if annotation := eFeature.GetEAnnotation(annotationURI); annotation != nil {
		if value := annotation.GetDetails().GetValue("namespace"); value != nil {
			namespace := value.(string)
			if namespace == "##targetNamespace" {
				eContainingClass := eFeature.GetEContainingClass()
				if eContainingClass != nil {
					ePackage := eContainingClass.GetEPackage()
					if ePackage != nil {
						return ePackage.GetNsURI()
					}
				}
			} else {
				return namespace
			}
		}
	}
	return ""
}
