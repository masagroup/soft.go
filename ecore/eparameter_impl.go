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

// EParameterImpl is the implementation of the model object 'EParameter'
type EParameterImpl struct {
	ETypedElementExt
}

// newEParameterImpl is the constructor of a EParameterImpl
func newEParameterImpl() *EParameterImpl {
	e := new(EParameterImpl)
	e.SetInterfaces(e)
	e.Initialize()
	return e
}

func (e *EParameterImpl) asEParameter() EParameter {
	return e.GetInterfaces().(EParameter)
}

func (e *EParameterImpl) EStaticClass() EClass {
	return GetPackage().GetEParameter()
}

func (e *EParameterImpl) EStaticFeatureCount() int {
	return EPARAMETER_FEATURE_COUNT
}

// GetEOperation get the value of eOperation
func (e *EParameterImpl) GetEOperation() EOperation {
	if e.EContainerFeatureID() == EPARAMETER__EOPERATION {
		return e.EContainer().(EOperation)
	}
	return nil
}

func (e *EParameterImpl) EGetFromID(featureID int, resolve bool) any {
	switch featureID {
	case EPARAMETER__EOPERATION:
		return e.asEParameter().GetEOperation()
	default:
		return e.ETypedElementExt.EGetFromID(featureID, resolve)
	}
}

func (e *EParameterImpl) EIsSetFromID(featureID int) bool {
	switch featureID {
	case EPARAMETER__EOPERATION:
		return e.asEParameter().GetEOperation() != nil
	default:
		return e.ETypedElementExt.EIsSetFromID(featureID)
	}
}

func (e *EParameterImpl) EBasicInverseAdd(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain {
	switch featureID {
	case EPARAMETER__EOPERATION:
		msgs := notifications
		if e.EInternalContainer() != nil {
			msgs = e.EBasicRemoveFromContainer(msgs)
		}
		return e.EBasicSetContainer(otherEnd, EPARAMETER__EOPERATION, msgs)
	default:
		return e.ETypedElementExt.EBasicInverseAdd(otherEnd, featureID, notifications)
	}
}

func (e *EParameterImpl) EBasicInverseRemove(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain {
	switch featureID {
	case EPARAMETER__EOPERATION:
		return e.EBasicSetContainer(nil, EPARAMETER__EOPERATION, notifications)
	default:
		return e.ETypedElementExt.EBasicInverseRemove(otherEnd, featureID, notifications)
	}
}
