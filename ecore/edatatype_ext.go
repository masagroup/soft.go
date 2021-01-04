// *****************************************************************************
//
// This file is part of a MASA library or program.
// Refer to the included end-user license agreement for restrictions.
//
// Copyright (c) 2019 MASA Group
//
// *****************************************************************************

package ecore

type EDataTypeInternal interface {
	EDataType
	SetDefaultValue(defaultValue interface{})
}

type eDataTypeExt struct {
	eDataTypeImpl
	defaultValue interface{}
}

func newEDataTypeExt() *eDataTypeExt {
	eDataType := new(eDataTypeExt)
	eDataType.SetInterfaces(eDataType)
	eDataType.Initialize()
	return eDataType
}

func (eDataType *eDataTypeExt) SetDefaultValue(newDefaultValue interface{}) {
	oldDefaultValue := eDataType.defaultValue
	eDataType.defaultValue = newDefaultValue
	if eDataType.ENotificationRequired() {
		eDataType.ENotify(NewNotificationByFeatureID(eDataType.AsEObject(), SET, EDATA_TYPE__DEFAULT_VALUE, oldDefaultValue, newDefaultValue, NO_INDEX))
	}
}
