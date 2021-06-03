// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import "strconv"

// EFactoryExt is the extension of the model object 'EFactory'
type EFactoryExt struct {
	eFactoryImpl
}

// NewEFactoryExt ...
func NewEFactoryExt() *EFactoryExt {
	eFactory := new(EFactoryExt)
	eFactory.SetInterfaces(eFactory)
	eFactory.Initialize()
	return eFactory
}

// Create ...
func (eFactory *EFactoryExt) Create(eClass EClass) EObject {
	if eFactory.GetEPackage() != eClass.GetEPackage() || eClass.IsAbstract() {
		panic("The class '" + eClass.GetName() + "' is not a valid classifier")
	}
	eObject := NewDynamicEObjectImpl()
	eObject.SetEClass(eClass)
	return eObject
}

// CreateFromString default implementation
func (eFactory *EFactoryExt) CreateFromString(eDataType EDataType, literalValue string) interface{} {
	if eFactory.GetEPackage() != eDataType.GetEPackage() {
		panic("The datatype '" + eDataType.GetName() + "' is not a valid classifier")
	}

	if eEnum := eDataType.(EEnum); eEnum != nil {
		result := eEnum.GetEEnumLiteralByLiteral(literalValue)
		if result == nil {
			panic("The value '" + literalValue + "' is not a valid enumerator of '" + eDataType.GetName() + "'")
		}
		return result.GetValue()
	}

	switch eDataType.GetInstanceTypeName() {
	case "float64":
		value, _ := strconv.ParseFloat(literalValue, 64)
		return value
	case "float32":
		value, _ := strconv.ParseFloat(literalValue, 32)
		return float32(value)
	case "int":
		value, _ := strconv.Atoi(literalValue)
		return value
	case "int64":
		value, _ := strconv.ParseInt(literalValue, 10, 64)
		return value
	case "int32":
		value, _ := strconv.ParseInt(literalValue, 10, 32)
		return int32(value)
	case "int16":
		value, _ := strconv.ParseInt(literalValue, 10, 16)
		return int16(value)
	case "int8":
		value, _ := strconv.ParseInt(literalValue, 10, 8)
		return int8(value)
	case "bool":
		value, _ := strconv.ParseBool(literalValue)
		return value
	case "string":
		return literalValue
	}

	panic("CreateFromString not implemented")
}

func (eFactory *EFactoryExt) ConvertToString(eDataType EDataType, instanceValue interface{}) string {
	panic("ConvertToString not implemented")
}
