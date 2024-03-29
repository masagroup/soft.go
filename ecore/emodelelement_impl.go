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

// EModelElementImpl is the implementation of the model object 'EModelElement'
type EModelElementImpl struct {
	CompactEObjectContainer
	eAnnotations EList
}
type eModelElementInitializers interface {
	initEAnnotations() EList
}

// newEModelElementImpl is the constructor of a EModelElementImpl
func newEModelElementImpl() *EModelElementImpl {
	e := new(EModelElementImpl)
	e.SetInterfaces(e)
	e.Initialize()
	return e
}

func (e *EModelElementImpl) Initialize() {
	e.CompactEObjectContainer.Initialize()

}

func (e *EModelElementImpl) asEModelElement() EModelElement {
	return e.GetInterfaces().(EModelElement)
}

func (e *EModelElementImpl) asInitializers() eModelElementInitializers {
	return e.GetInterfaces().(eModelElementInitializers)
}

func (e *EModelElementImpl) EStaticClass() EClass {
	return GetPackage().GetEModelElement()
}

func (e *EModelElementImpl) EStaticFeatureCount() int {
	return EMODEL_ELEMENT_FEATURE_COUNT
}

// GetEAnnotation default implementation
func (e *EModelElementImpl) GetEAnnotation(string) EAnnotation {
	panic("GetEAnnotation not implemented")
}

// GetEAnnotations get the value of eAnnotations
func (e *EModelElementImpl) GetEAnnotations() EList {
	if e.eAnnotations == nil {
		e.eAnnotations = e.asInitializers().initEAnnotations()
	}
	return e.eAnnotations
}

func (e *EModelElementImpl) initEAnnotations() EList {
	return NewBasicEObjectList(e.AsEObjectInternal(), EMODEL_ELEMENT__EANNOTATIONS, EANNOTATION__EMODEL_ELEMENT, true, true, true, false, false)
}

func (e *EModelElementImpl) EGetFromID(featureID int, resolve bool) any {
	switch featureID {
	case EMODEL_ELEMENT__EANNOTATIONS:
		return e.asEModelElement().GetEAnnotations()
	default:
		return e.CompactEObjectContainer.EGetFromID(featureID, resolve)
	}
}

func (e *EModelElementImpl) ESetFromID(featureID int, newValue any) {
	switch featureID {
	case EMODEL_ELEMENT__EANNOTATIONS:
		list := e.asEModelElement().GetEAnnotations()
		list.Clear()
		list.AddAll(newValue.(EList))
	default:
		e.CompactEObjectContainer.ESetFromID(featureID, newValue)
	}
}

func (e *EModelElementImpl) EUnsetFromID(featureID int) {
	switch featureID {
	case EMODEL_ELEMENT__EANNOTATIONS:
		e.asEModelElement().GetEAnnotations().Clear()
	default:
		e.CompactEObjectContainer.EUnsetFromID(featureID)
	}
}

func (e *EModelElementImpl) EIsSetFromID(featureID int) bool {
	switch featureID {
	case EMODEL_ELEMENT__EANNOTATIONS:
		return e.eAnnotations != nil && e.eAnnotations.Size() != 0
	default:
		return e.CompactEObjectContainer.EIsSetFromID(featureID)
	}
}

func (e *EModelElementImpl) EInvokeFromID(operationID int, arguments EList) any {
	switch operationID {
	case EMODEL_ELEMENT__GET_EANNOTATION_ESTRING:
		return e.asEModelElement().GetEAnnotation(arguments.Get(0).(string))
	default:
		return e.CompactEObjectContainer.EInvokeFromID(operationID, arguments)
	}
}

func (e *EModelElementImpl) EBasicInverseAdd(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain {
	switch featureID {
	case EMODEL_ELEMENT__EANNOTATIONS:
		list := e.GetEAnnotations().(ENotifyingList)
		return list.AddWithNotification(otherEnd, notifications)
	default:
		return e.CompactEObjectContainer.EBasicInverseAdd(otherEnd, featureID, notifications)
	}
}

func (e *EModelElementImpl) EBasicInverseRemove(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain {
	switch featureID {
	case EMODEL_ELEMENT__EANNOTATIONS:
		list := e.GetEAnnotations().(ENotifyingList)
		return list.RemoveWithNotification(otherEnd, notifications)
	default:
		return e.CompactEObjectContainer.EBasicInverseRemove(otherEnd, featureID, notifications)
	}
}
