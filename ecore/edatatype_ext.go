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
	defaultValue      any
	defaultValueIsSet bool
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
	eDataType.defaultValueIsSet = true
	if eDataType.ENotificationRequired() {
		eDataType.ENotify(NewNotificationByFeatureID(eDataType.AsEObject(), SET, EDATA_TYPE__DEFAULT_VALUE, oldDefaultValue, newDefaultValue, NO_INDEX))
	}
}

func (eDataType *EDataTypeExt) GetDefaultValue() any {
	if !eDataType.defaultValueIsSet {
		switch getInstanceTypeName(eDataType) {
		case "float64", "java.lang.Double", "double":
			return float64(0)
		case "float32", "java.lang.Float", "float":
			return float32(0)
		case "int", "java.lang.Integer":
			return int(0)
		case "uint64", "com.google.common.primitives.UnsignedLong":
			return uint64(0)
		case "int64", "java.lang.Long", "java.math.BigInteger", "long":
			return int64(0)
		case "int32":
			return int32(0)
		case "int16", "java.lang.Short", "short":
			return int16(0)
		case "int8":
			return int8(0)
		case "bool", "java.lang.Boolean", "boolean":
			return false
		case "string", "java.lang.String":
			return ""
		}
		eDataType.defaultValueIsSet = true
	}
	return eDataType.defaultValue
}

func getInstanceTypeName(eDataType EDataType) string {
	if eAnnotation := eDataType.GetEAnnotation("http://net.masagroup/soft/2019/GenGo"); eAnnotation != nil {
		if instanceTypeName, _ := eAnnotation.GetDetails().GetValue("instanceTypeName").(string); instanceTypeName != "" {
			return instanceTypeName
		}
	}
	return eDataType.GetInstanceTypeName()
}
