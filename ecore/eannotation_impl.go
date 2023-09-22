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

// EAnnotationImpl is the implementation of the model object 'EAnnotation'
type EAnnotationImpl struct {
	EModelElementExt
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

// newEAnnotationImpl is the constructor of a EAnnotationImpl
func newEAnnotationImpl() *EAnnotationImpl {
	e := new(EAnnotationImpl)
	e.SetInterfaces(e)
	e.Initialize()
	return e
}

func (e *EAnnotationImpl) Initialize() {
	e.EModelElementExt.Initialize()
	e.source = ""

}

func (e *EAnnotationImpl) asEAnnotation() EAnnotation {
	return e.GetInterfaces().(EAnnotation)
}

func (e *EAnnotationImpl) asInitializers() eAnnotationInitializers {
	return e.GetInterfaces().(eAnnotationInitializers)
}

func (e *EAnnotationImpl) asBasics() eAnnotationBasics {
	return e.GetInterfaces().(eAnnotationBasics)
}

func (e *EAnnotationImpl) EStaticClass() EClass {
	return GetPackage().GetEAnnotationClass()
}

func (e *EAnnotationImpl) EStaticFeatureCount() int {
	return EANNOTATION_FEATURE_COUNT
}

// GetContents get the value of contents
func (e *EAnnotationImpl) GetContents() EList {
	if e.contents == nil {
		e.contents = e.asInitializers().initContents()
	}
	return e.contents
}

// GetDetails get the value of details
func (e *EAnnotationImpl) GetDetails() EMap {
	if e.details == nil {
		e.details = e.asInitializers().initDetails()
	}
	return e.details
}

// GetEModelElement get the value of eModelElement
func (e *EAnnotationImpl) GetEModelElement() EModelElement {
	if e.EContainerFeatureID() == EANNOTATION__EMODEL_ELEMENT {
		return e.EContainer().(EModelElement)
	}
	return nil
}

// SetEModelElement set the value of eModelElement
func (e *EAnnotationImpl) SetEModelElement(newEModelElement EModelElement) {
	if newEModelElement != e.EInternalContainer() || (newEModelElement != nil && e.EContainerFeatureID() != EANNOTATION__EMODEL_ELEMENT) {
		var notifications ENotificationChain
		if e.EInternalContainer() != nil {
			notifications = e.EBasicRemoveFromContainer(notifications)
		}
		if newEModelElementInternal, _ := newEModelElement.(EObjectInternal); newEModelElementInternal != nil {
			notifications = newEModelElementInternal.EInverseAdd(e.AsEObject(), EMODEL_ELEMENT__EANNOTATIONS, notifications)
		}
		notifications = e.asBasics().basicSetEModelElement(newEModelElement, notifications)
		if notifications != nil {
			notifications.Dispatch()
		}
	} else if e.ENotificationRequired() {
		e.ENotify(NewNotificationByFeatureID(e, SET, EANNOTATION__EMODEL_ELEMENT, newEModelElement, newEModelElement, NO_INDEX))
	}
}

func (e *EAnnotationImpl) basicSetEModelElement(newEModelElement EModelElement, msgs ENotificationChain) ENotificationChain {
	return e.EBasicSetContainer(newEModelElement, EANNOTATION__EMODEL_ELEMENT, msgs)
}

// GetReferences get the value of references
func (e *EAnnotationImpl) GetReferences() EList {
	if e.references == nil {
		e.references = e.asInitializers().initReferences()
	}
	return e.references
}

// GetSource get the value of source
func (e *EAnnotationImpl) GetSource() string {
	return e.source
}

// SetSource set the value of source
func (e *EAnnotationImpl) SetSource(newSource string) {
	oldSource := e.source
	e.source = newSource
	if e.ENotificationRequired() {
		e.ENotify(NewNotificationByFeatureID(e.AsEObject(), SET, EANNOTATION__SOURCE, oldSource, newSource, NO_INDEX))
	}
}

func (e *EAnnotationImpl) initContents() EList {
	return NewBasicEObjectList(e.AsEObjectInternal(), EANNOTATION__CONTENTS, -1, true, true, false, false, false)
}

func (e *EAnnotationImpl) initDetails() EMap {
	return NewBasicEObjectMap(GetPackage().GetEStringToStringMapEntry(), e.AsEObjectInternal(), EANNOTATION__DETAILS, -1, false)
}

func (e *EAnnotationImpl) initReferences() EList {
	return NewBasicEObjectList(e.AsEObjectInternal(), EANNOTATION__REFERENCES, -1, false, false, false, true, false)
}

func (e *EAnnotationImpl) EGetFromID(featureID int, resolve bool) any {
	switch featureID {
	case EANNOTATION__CONTENTS:
		return e.asEAnnotation().GetContents()
	case EANNOTATION__DETAILS:
		return e.asEAnnotation().GetDetails()
	case EANNOTATION__EMODEL_ELEMENT:
		return e.asEAnnotation().GetEModelElement()
	case EANNOTATION__REFERENCES:
		list := e.asEAnnotation().GetReferences()
		if !resolve {
			if objects, _ := list.(EObjectList); objects != nil {
				return objects.GetUnResolvedList()
			}
		}
		return list
	case EANNOTATION__SOURCE:
		return e.asEAnnotation().GetSource()
	default:
		return e.EModelElementExt.EGetFromID(featureID, resolve)
	}
}

func (e *EAnnotationImpl) ESetFromID(featureID int, newValue any) {
	switch featureID {
	case EANNOTATION__CONTENTS:
		list := e.asEAnnotation().GetContents()
		list.Clear()
		list.AddAll(newValue.(EList))
	case EANNOTATION__DETAILS:
		m := e.asEAnnotation().GetDetails()
		m.Clear()
		m.AddAll(newValue.(EList))
	case EANNOTATION__EMODEL_ELEMENT:
		newValueOrNil, _ := newValue.(EModelElement)
		e.asEAnnotation().SetEModelElement(newValueOrNil)
	case EANNOTATION__REFERENCES:
		list := e.asEAnnotation().GetReferences()
		list.Clear()
		list.AddAll(newValue.(EList))
	case EANNOTATION__SOURCE:
		e.asEAnnotation().SetSource(newValue.(string))
	default:
		e.EModelElementExt.ESetFromID(featureID, newValue)
	}
}

func (e *EAnnotationImpl) EUnsetFromID(featureID int) {
	switch featureID {
	case EANNOTATION__CONTENTS:
		e.asEAnnotation().GetContents().Clear()
	case EANNOTATION__DETAILS:
		e.asEAnnotation().GetDetails().Clear()
	case EANNOTATION__EMODEL_ELEMENT:
		e.asEAnnotation().SetEModelElement(nil)
	case EANNOTATION__REFERENCES:
		e.asEAnnotation().GetReferences().Clear()
	case EANNOTATION__SOURCE:
		e.asEAnnotation().SetSource("")
	default:
		e.EModelElementExt.EUnsetFromID(featureID)
	}
}

func (e *EAnnotationImpl) EIsSetFromID(featureID int) bool {
	switch featureID {
	case EANNOTATION__CONTENTS:
		return e.contents != nil && e.contents.Size() != 0
	case EANNOTATION__DETAILS:
		return e.details != nil && e.details.Size() != 0
	case EANNOTATION__EMODEL_ELEMENT:
		return e.asEAnnotation().GetEModelElement() != nil
	case EANNOTATION__REFERENCES:
		return e.references != nil && e.references.Size() != 0
	case EANNOTATION__SOURCE:
		return e.source != ""
	default:
		return e.EModelElementExt.EIsSetFromID(featureID)
	}
}

func (e *EAnnotationImpl) EBasicInverseAdd(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain {
	switch featureID {
	case EANNOTATION__EMODEL_ELEMENT:
		msgs := notifications
		if e.EInternalContainer() != nil {
			msgs = e.EBasicRemoveFromContainer(msgs)
		}
		return e.asBasics().basicSetEModelElement(otherEnd.(EModelElement), msgs)
	default:
		return e.EModelElementExt.EBasicInverseAdd(otherEnd, featureID, notifications)
	}
}

func (e *EAnnotationImpl) EBasicInverseRemove(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain {
	switch featureID {
	case EANNOTATION__CONTENTS:
		list := e.GetContents().(ENotifyingList)
		return list.RemoveWithNotification(otherEnd, notifications)
	case EANNOTATION__DETAILS:
		return notifications
	case EANNOTATION__EMODEL_ELEMENT:
		return e.asBasics().basicSetEModelElement(nil, notifications)
	default:
		return e.EModelElementExt.EBasicInverseRemove(otherEnd, featureID, notifications)
	}
}
