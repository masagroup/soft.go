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

// eAnnotationImpl is the implementation of the model object 'EAnnotation'
type eAnnotationImpl struct {
	eModelElementExt
	contents   EList
	details    EMap
	references EList
	source     string
}
type eAnnotationInitializers interface {
	initContents() EList
	initDetails() EMap
	initReferences() EList
}

type eAnnotationBasics interface {
	basicSetEModelElement(EModelElement, ENotificationChain) ENotificationChain
}

// newEAnnotationImpl is the constructor of a eAnnotationImpl
func newEAnnotationImpl() *eAnnotationImpl {
	eAnnotation := new(eAnnotationImpl)
	eAnnotation.SetInterfaces(eAnnotation)
	eAnnotation.Initialize()
	return eAnnotation
}

func (eAnnotation *eAnnotationImpl) Initialize() {
	eAnnotation.eModelElementExt.Initialize()
	eAnnotation.source = ""

}

func (eAnnotation *eAnnotationImpl) asEAnnotation() EAnnotation {
	return eAnnotation.GetInterfaces().(EAnnotation)
}

func (eAnnotation *eAnnotationImpl) asInitializers() eAnnotationInitializers {
	return eAnnotation.GetInterfaces().(eAnnotationInitializers)
}

func (eAnnotation *eAnnotationImpl) asBasics() eAnnotationBasics {
	return eAnnotation.GetInterfaces().(eAnnotationBasics)
}

func (eAnnotation *eAnnotationImpl) EStaticClass() EClass {
	return GetPackage().GetEAnnotationClass()
}

func (eAnnotation *eAnnotationImpl) EStaticFeatureCount() int {
	return EANNOTATION_FEATURE_COUNT
}

// GetContents get the value of contents
func (eAnnotation *eAnnotationImpl) GetContents() EList {
	if eAnnotation.contents == nil {
		eAnnotation.contents = eAnnotation.asInitializers().initContents()
	}
	return eAnnotation.contents
}

// GetDetails get the value of details
func (eAnnotation *eAnnotationImpl) GetDetails() EMap {
	if eAnnotation.details == nil {
		eAnnotation.details = eAnnotation.asInitializers().initDetails()
	}
	return eAnnotation.details
}

// GetEModelElement get the value of eModelElement
func (eAnnotation *eAnnotationImpl) GetEModelElement() EModelElement {
	if eAnnotation.EContainerFeatureID() == EANNOTATION__EMODEL_ELEMENT {
		return eAnnotation.EContainer().(EModelElement)
	}
	return nil
}

// SetEModelElement set the value of eModelElement
func (eAnnotation *eAnnotationImpl) SetEModelElement(newEModelElement EModelElement) {
	if newEModelElement != eAnnotation.EInternalContainer() || (newEModelElement != nil && eAnnotation.EContainerFeatureID() != EANNOTATION__EMODEL_ELEMENT) {
		var notifications ENotificationChain
		if eAnnotation.EInternalContainer() != nil {
			notifications = eAnnotation.EBasicRemoveFromContainer(notifications)
		}
		if newEModelElementInternal, _ := newEModelElement.(EObjectInternal); newEModelElementInternal != nil {
			notifications = newEModelElementInternal.EInverseAdd(eAnnotation.AsEObject(), EMODEL_ELEMENT__EANNOTATIONS, notifications)
		}
		notifications = eAnnotation.asBasics().basicSetEModelElement(newEModelElement, notifications)
		if notifications != nil {
			notifications.Dispatch()
		}
	} else if eAnnotation.ENotificationRequired() {
		eAnnotation.ENotify(NewNotificationByFeatureID(eAnnotation, SET, EANNOTATION__EMODEL_ELEMENT, newEModelElement, newEModelElement, NO_INDEX))
	}
}

func (eAnnotation *eAnnotationImpl) basicSetEModelElement(newEModelElement EModelElement, msgs ENotificationChain) ENotificationChain {
	return eAnnotation.EBasicSetContainer(newEModelElement, EANNOTATION__EMODEL_ELEMENT, msgs)
}

// GetReferences get the value of references
func (eAnnotation *eAnnotationImpl) GetReferences() EList {
	if eAnnotation.references == nil {
		eAnnotation.references = eAnnotation.asInitializers().initReferences()
	}
	return eAnnotation.references
}

// GetSource get the value of source
func (eAnnotation *eAnnotationImpl) GetSource() string {
	return eAnnotation.source
}

// SetSource set the value of source
func (eAnnotation *eAnnotationImpl) SetSource(newSource string) {
	oldSource := eAnnotation.source
	eAnnotation.source = newSource
	if eAnnotation.ENotificationRequired() {
		eAnnotation.ENotify(NewNotificationByFeatureID(eAnnotation.AsEObject(), SET, EANNOTATION__SOURCE, oldSource, newSource, NO_INDEX))
	}
}

func (eAnnotation *eAnnotationImpl) initContents() EList {
	return NewBasicEObjectList(eAnnotation.AsEObjectInternal(), EANNOTATION__CONTENTS, -1, true, true, false, false, false)
}

func (eAnnotation *eAnnotationImpl) initDetails() EMap {
	return NewBasicEObjectMap(GetPackage().GetEStringToStringMapEntry(), eAnnotation.AsEObjectInternal(), EANNOTATION__DETAILS, -1, false)
}

func (eAnnotation *eAnnotationImpl) initReferences() EList {
	return NewBasicEObjectList(eAnnotation.AsEObjectInternal(), EANNOTATION__REFERENCES, -1, false, false, false, true, false)
}

func (eAnnotation *eAnnotationImpl) EGetFromID(featureID int, resolve bool) any {
	switch featureID {
	case EANNOTATION__CONTENTS:
		return eAnnotation.asEAnnotation().GetContents()
	case EANNOTATION__DETAILS:
		return eAnnotation.asEAnnotation().GetDetails()
	case EANNOTATION__EMODEL_ELEMENT:
		return eAnnotation.asEAnnotation().GetEModelElement()
	case EANNOTATION__REFERENCES:
		list := eAnnotation.asEAnnotation().GetReferences()
		if !resolve {
			if objects, _ := list.(EObjectList); objects != nil {
				return objects.GetUnResolvedList()
			}
		}
		return list
	case EANNOTATION__SOURCE:
		return eAnnotation.asEAnnotation().GetSource()
	default:
		return eAnnotation.eModelElementExt.EGetFromID(featureID, resolve)
	}
}

func (eAnnotation *eAnnotationImpl) ESetFromID(featureID int, newValue any) {
	switch featureID {
	case EANNOTATION__CONTENTS:
		list := eAnnotation.asEAnnotation().GetContents()
		list.Clear()
		list.AddAll(newValue.(EList))
	case EANNOTATION__DETAILS:
		m := eAnnotation.asEAnnotation().GetDetails()
		m.Clear()
		m.AddAll(newValue.(EList))
	case EANNOTATION__EMODEL_ELEMENT:
		eAnnotation.asEAnnotation().SetEModelElement(newValue.(EModelElement))
	case EANNOTATION__REFERENCES:
		list := eAnnotation.asEAnnotation().GetReferences()
		list.Clear()
		list.AddAll(newValue.(EList))
	case EANNOTATION__SOURCE:
		eAnnotation.asEAnnotation().SetSource(newValue.(string))
	default:
		eAnnotation.eModelElementExt.ESetFromID(featureID, newValue)
	}
}

func (eAnnotation *eAnnotationImpl) EUnsetFromID(featureID int) {
	switch featureID {
	case EANNOTATION__CONTENTS:
		eAnnotation.asEAnnotation().GetContents().Clear()
	case EANNOTATION__DETAILS:
		eAnnotation.asEAnnotation().GetDetails().Clear()
	case EANNOTATION__EMODEL_ELEMENT:
		eAnnotation.asEAnnotation().SetEModelElement(nil)
	case EANNOTATION__REFERENCES:
		eAnnotation.asEAnnotation().GetReferences().Clear()
	case EANNOTATION__SOURCE:
		eAnnotation.asEAnnotation().SetSource("")
	default:
		eAnnotation.eModelElementExt.EUnsetFromID(featureID)
	}
}

func (eAnnotation *eAnnotationImpl) EIsSetFromID(featureID int) bool {
	switch featureID {
	case EANNOTATION__CONTENTS:
		return eAnnotation.contents != nil && eAnnotation.contents.Size() != 0
	case EANNOTATION__DETAILS:
		return eAnnotation.details != nil && eAnnotation.details.Size() != 0
	case EANNOTATION__EMODEL_ELEMENT:
		return eAnnotation.asEAnnotation().GetEModelElement() != nil
	case EANNOTATION__REFERENCES:
		return eAnnotation.references != nil && eAnnotation.references.Size() != 0
	case EANNOTATION__SOURCE:
		return eAnnotation.source != ""
	default:
		return eAnnotation.eModelElementExt.EIsSetFromID(featureID)
	}
}

func (eAnnotation *eAnnotationImpl) EBasicInverseAdd(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain {
	switch featureID {
	case EANNOTATION__EMODEL_ELEMENT:
		msgs := notifications
		if eAnnotation.EInternalContainer() != nil {
			msgs = eAnnotation.EBasicRemoveFromContainer(msgs)
		}
		return eAnnotation.asBasics().basicSetEModelElement(otherEnd.(EModelElement), msgs)
	default:
		return eAnnotation.eModelElementExt.EBasicInverseAdd(otherEnd, featureID, notifications)
	}
}

func (eAnnotation *eAnnotationImpl) EBasicInverseRemove(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain {
	switch featureID {
	case EANNOTATION__CONTENTS:
		list := eAnnotation.GetContents().(ENotifyingList)
		return list.RemoveWithNotification(otherEnd, notifications)
	case EANNOTATION__DETAILS:
		return notifications
	case EANNOTATION__EMODEL_ELEMENT:
		return eAnnotation.asBasics().basicSetEModelElement(nil, notifications)
	default:
		return eAnnotation.eModelElementExt.EBasicInverseRemove(otherEnd, featureID, notifications)
	}
}
