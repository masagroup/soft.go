// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

type EDataTypeInternal interface {
	EDataType
	SetDefaultValue(defaultValue any)
}

type EDataTypeExt struct {
	EDataTypeImpl
	defaultValue any
}

func newEDataTypeExt() *EDataTypeExt {
	eDataType := new(EDataTypeExt)
	eDataType.SetInterfaces(eDataType)
	eDataType.Initialize()
	return eDataType
}

func (eDataType *EDataTypeExt) SetDefaultValue(newDefaultValue any) {
	oldDefaultValue := eDataType.defaultValue
	eDataType.defaultValue = newDefaultValue
	if eDataType.ENotificationRequired() {
		eDataType.ENotify(NewNotificationByFeatureID(eDataType.AsEObject(), SET, EDATA_TYPE__DEFAULT_VALUE, oldDefaultValue, newDefaultValue, NO_INDEX))
	}
}

func (eDataType *EDataTypeExt) GetDefaultValue() any {
	return eDataType.defaultValue
}
