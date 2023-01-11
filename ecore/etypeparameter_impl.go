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

// ETypeParameterImpl is the implementation of the model object 'ETypeParameter'
type ETypeParameterImpl struct {
	ENamedElementImpl
	eBounds EList
}
type eTypeParameterInitializers interface {
	initEBounds() EList
}

// newETypeParameterImpl is the constructor of a ETypeParameterImpl
func newETypeParameterImpl() *ETypeParameterImpl {
	eTypeParameter := new(ETypeParameterImpl)
	eTypeParameter.SetInterfaces(eTypeParameter)
	eTypeParameter.Initialize()
	return eTypeParameter
}

func (eTypeParameter *ETypeParameterImpl) Initialize() {
	eTypeParameter.ENamedElementImpl.Initialize()

}

func (eTypeParameter *ETypeParameterImpl) asETypeParameter() ETypeParameter {
	return eTypeParameter.GetInterfaces().(ETypeParameter)
}

func (eTypeParameter *ETypeParameterImpl) asInitializers() eTypeParameterInitializers {
	return eTypeParameter.GetInterfaces().(eTypeParameterInitializers)
}

func (eTypeParameter *ETypeParameterImpl) EStaticClass() EClass {
	return GetPackage().GetETypeParameter()
}

func (eTypeParameter *ETypeParameterImpl) EStaticFeatureCount() int {
	return ETYPE_PARAMETER_FEATURE_COUNT
}

// GetEBounds get the value of eBounds
func (eTypeParameter *ETypeParameterImpl) GetEBounds() EList {
	if eTypeParameter.eBounds == nil {
		eTypeParameter.eBounds = eTypeParameter.asInitializers().initEBounds()
	}
	return eTypeParameter.eBounds
}

func (eTypeParameter *ETypeParameterImpl) initEBounds() EList {
	return NewBasicEObjectList(eTypeParameter.AsEObjectInternal(), ETYPE_PARAMETER__EBOUNDS, -1, true, true, false, false, false)
}

func (eTypeParameter *ETypeParameterImpl) EGetFromID(featureID int, resolve bool) any {
	switch featureID {
	case ETYPE_PARAMETER__EBOUNDS:
		return eTypeParameter.asETypeParameter().GetEBounds()
	default:
		return eTypeParameter.ENamedElementImpl.EGetFromID(featureID, resolve)
	}
}

func (eTypeParameter *ETypeParameterImpl) ESetFromID(featureID int, newValue any) {
	switch featureID {
	case ETYPE_PARAMETER__EBOUNDS:
		list := eTypeParameter.asETypeParameter().GetEBounds()
		list.Clear()
		list.AddAll(newValue.(EList))
	default:
		eTypeParameter.ENamedElementImpl.ESetFromID(featureID, newValue)
	}
}

func (eTypeParameter *ETypeParameterImpl) EUnsetFromID(featureID int) {
	switch featureID {
	case ETYPE_PARAMETER__EBOUNDS:
		eTypeParameter.asETypeParameter().GetEBounds().Clear()
	default:
		eTypeParameter.ENamedElementImpl.EUnsetFromID(featureID)
	}
}

func (eTypeParameter *ETypeParameterImpl) EIsSetFromID(featureID int) bool {
	switch featureID {
	case ETYPE_PARAMETER__EBOUNDS:
		return eTypeParameter.eBounds != nil && eTypeParameter.eBounds.Size() != 0
	default:
		return eTypeParameter.ENamedElementImpl.EIsSetFromID(featureID)
	}
}

func (eTypeParameter *ETypeParameterImpl) EBasicInverseRemove(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain {
	switch featureID {
	case ETYPE_PARAMETER__EBOUNDS:
		list := eTypeParameter.GetEBounds().(ENotifyingList)
		return list.RemoveWithNotification(otherEnd, notifications)
	default:
		return eTypeParameter.ENamedElementImpl.EBasicInverseRemove(otherEnd, featureID, notifications)
	}
}
