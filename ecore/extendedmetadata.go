// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

const (
	annotationURI = "http:///org/eclipse/emf/ecore/util/ExtendedMetaData"
)

type ExtendedMetaData struct {
	metaData map[interface{}]interface{}
}

type ENamedElementExtendedMetaData interface {
	getName() string
}

type ENamedElementExtendedMetaDataImpl struct {
	emd      *ExtendedMetaData
	eElement ENamedElement
	name     string
}

func (m *ENamedElementExtendedMetaDataImpl) getName() string {
	if m.name == "unitialized" {
		m.name = m.emd.basicGetName(m.eElement)
	}
	return m.name
}

type EPackageExtentedMetaData interface {
	getType(string) EClassifier
}

type EPackageExtentedMetaDataImpl struct {
	emd                 *ExtendedMetaData
	ePackage            EPackage
	nameToClassifierMap map[string]EClassifier
}

func (m *EPackageExtentedMetaDataImpl) getType(name string) EClassifier {
	var eResult EClassifier = nil
	if m.nameToClassifierMap != nil {
		eResult = m.nameToClassifierMap[name]
	}
	if eResult == nil {
		eClassifiers := m.ePackage.GetEClassifiers()
		if m.nameToClassifierMap == nil || len(m.nameToClassifierMap) != eClassifiers.Size() {
			nameToClassifierMap := make(map[string]EClassifier)
			for it := eClassifiers.Iterator(); it.HasNext(); {
				eClassifier := it.Next().(EClassifier)
				eClassifierName := m.emd.GetName(eClassifier)
				nameToClassifierMap[eClassifierName] = eClassifier
				if eClassifierName == name {
					eResult = eClassifier
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
		m.namespace = m.emd.basicGetNamespace(m.eElement.(EStructuralFeature))
	}
	return m.namespace
}

func NewExtendedMetaData() *ExtendedMetaData {
	return &ExtendedMetaData{metaData: make(map[interface{}]interface{})}
}

func (emd *ExtendedMetaData) getENamedElementExtendedMetaData(eElement ENamedElement) ENamedElementExtendedMetaData {
	var result ENamedElementExtendedMetaData
	if i, ok := emd.metaData[eElement]; ok {
		result = i.(ENamedElementExtendedMetaData)
	} else {
		if eFeature, _ := eElement.(EStructuralFeature); eFeature != nil {
			result = &EStructuralFeatureExtentedMetaDataImpl{ENamedElementExtendedMetaDataImpl: ENamedElementExtendedMetaDataImpl{emd: emd, eElement: eElement, name: "unitialized"}, namespace: "unitialized"}
		} else {
			result = &ENamedElementExtendedMetaDataImpl{emd: emd, eElement: eElement, name: "unitialized"}
		}
		emd.metaData[eElement] = result
	}
	return result
}

func (emd *ExtendedMetaData) getEStructuralFeatureExtentedMetaData(eFeature EStructuralFeature) EStructuralFeatureExtentedMetaData {
	var result EStructuralFeatureExtentedMetaData
	if i, ok := emd.metaData[eFeature]; ok {
		result = i.(EStructuralFeatureExtentedMetaData)
	} else {
		result = &EStructuralFeatureExtentedMetaDataImpl{ENamedElementExtendedMetaDataImpl: ENamedElementExtendedMetaDataImpl{emd: emd, eElement: eFeature, name: "unitialized"}, namespace: "unitialized"}
		emd.metaData[eFeature] = result
	}
	return result
}

func (emd *ExtendedMetaData) getEPackageExtentedMetaData(ePackage EPackage) EPackageExtentedMetaData {
	var result EPackageExtentedMetaData
	if i, ok := emd.metaData[ePackage]; ok {
		result = i.(EPackageExtentedMetaData)
	} else {
		result = &EPackageExtentedMetaDataImpl{emd: emd, ePackage: ePackage}
		emd.metaData[ePackage] = result
	}
	return result
}

func (emd *ExtendedMetaData) GetType(ePackage EPackage, name string) EClassifier {
	return emd.getEPackageExtentedMetaData(ePackage).getType(name)
}

func (emd *ExtendedMetaData) GetName(eElement ENamedElement) string {
	return emd.getENamedElementExtendedMetaData(eElement).getName()
}

func (emd *ExtendedMetaData) GetNamespace(eFeature EStructuralFeature) string {
	return emd.getEStructuralFeatureExtentedMetaData(eFeature).getNamespace()
}

func (emd *ExtendedMetaData) GetDocumentRoot(ePackage EPackage) EClass {
	eClassifier := emd.GetType(ePackage, "")
	if eClassifier != nil {
		return eClassifier.(EClass)
	}
	return nil
}

func (emd *ExtendedMetaData) GetXMLNSPrefixMapFeature(eClass EClass) EReference {
	for it := eClass.GetEAllReferences().Iterator(); it.HasNext(); {
		eReference := it.Next().(EReference)
		if emd.GetName(eReference) == "xmlns:prefix" {
			return eReference
		}
	}
	return nil
}

func (emd *ExtendedMetaData) GetXSISchemaLocationMapFeature(eClass EClass) EReference {
	for it := eClass.GetEAllReferences().Iterator(); it.HasNext(); {
		eReference := it.Next().(EReference)
		if emd.GetName(eReference) == "xsi:schemaLocation" {
			return eReference
		}
	}
	return nil
}

func (emd *ExtendedMetaData) basicGetName(eElement ENamedElement) string {
	if annotation := eElement.GetEAnnotation(annotationURI); annotation != nil {
		if name := annotation.GetDetails().GetValue("name"); name != nil {
			return name.(string)
		}
	}
	return eElement.GetName()
}

func (emd *ExtendedMetaData) basicGetNamespace(eFeature EStructuralFeature) string {
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
