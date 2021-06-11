// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import (
	"fmt"
	"strconv"
)

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
		panic(fmt.Sprintf("The datatype '%v' is not a valid classifier", eDataType.GetName()))
	}

	if eEnum := eDataType.(EEnum); eEnum != nil {
		result := eEnum.GetEEnumLiteralByLiteral(literalValue)
		if result == nil {
			panic(fmt.Sprintf("The value '%v' is not a valid enumerator of '%v'", literalValue, eDataType.GetName()))

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
	if eFactory.GetEPackage() != eDataType.GetEPackage() {
		panic(fmt.Sprintf("The datatype '%v' is not a valid classifier", eDataType.GetName()))
	}

	if eEnum := eDataType.(EEnum); eEnum != nil {
		result := eEnum.GetEEnumLiteralByValue(instanceValue.(int))
		if result == nil {
			panic(fmt.Sprintf("The value '%v' is not a valid enumerator of '%v'", instanceValue, eDataType.GetName()))
		}
		return result.GetLiteral()
	}

	switch eDataType.GetInstanceTypeName() {
	case "float64":
		v, _ := instanceValue.(float64)
		return strconv.FormatFloat(v, 'f', -1, 64)
	case "float32":
		v, _ := instanceValue.(float64)
		return strconv.FormatFloat(v, 'f', -1, 32)
	case "int":
		v, _ := instanceValue.(int)
		return strconv.Itoa(v)
	case "int64":
		v, _ := instanceValue.(int64)
		return strconv.FormatInt(v, 10)
	case "int32":
		v, _ := instanceValue.(int32)
		return strconv.FormatInt(int64(v), 10)
	case "int16":
		v, _ := instanceValue.(int16)
		return strconv.FormatInt(int64(v), 10)
	case "int8":
		v, _ := instanceValue.(int8)
		return strconv.FormatInt(int64(v), 10)
	case "bool":
		v, _ := instanceValue.(bool)
		return strconv.FormatBool(v)
	case "string":
		return instanceValue.(string)
	}

	panic("ConvertToString not implemented")
}
