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
	EFactoryImpl
}

// newEFactoryExt ...
func newEFactoryExt() *EFactoryExt {
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
	if IsMapEntry(eClass) {
		eEntry := NewDynamicEMapEntryImpl()
		eEntry.SetEClass(eClass)
		return eEntry
	} else {
		eObject := NewDynamicEObjectImpl()
		eObject.SetEClass(eClass)
		return eObject
	}
}

// CreateFromString default implementation
func (eFactory *EFactoryExt) CreateFromString(eDataType EDataType, literalValue string) any {
	if eFactory.GetEPackage() != eDataType.GetEPackage() {
		panic(fmt.Sprintf("The datatype '%v' is not a valid classifier", eDataType.GetName()))
	}

	if eEnum, _ := eDataType.(EEnum); eEnum != nil {
		result := eEnum.GetEEnumLiteralByLiteral(literalValue)
		if result == nil {
			panic(fmt.Sprintf("The value '%v' is not a valid enumerator of '%v'", literalValue, eDataType.GetName()))

		}
		return result.GetValue()
	}

	switch getInstanceTypeName(eDataType) {
	case "float64", "java.lang.Double", "double":
		value, _ := strconv.ParseFloat(literalValue, 64)
		return value
	case "float32", "java.lang.Float", "float":
		value, _ := strconv.ParseFloat(literalValue, 32)
		return float32(value)
	case "int", "java.lang.Integer":
		value, _ := strconv.Atoi(literalValue)
		return value
	case "uint64", "com.google.common.primitives.UnsignedLong":
		value, _ := strconv.ParseUint(literalValue, 10, 64)
		return value
	case "int64", "java.lang.Long", "java.math.BigInteger", "long":
		value, _ := strconv.ParseInt(literalValue, 10, 64)
		return value
	case "int32":
		value, _ := strconv.ParseInt(literalValue, 10, 32)
		return int32(value)
	case "int16", "java.lang.Short", "short":
		value, _ := strconv.ParseInt(literalValue, 10, 16)
		return int16(value)
	case "int8":
		value, _ := strconv.ParseInt(literalValue, 10, 8)
		return int8(value)
	case "bool", "java.lang.Boolean", "boolean":
		value, _ := strconv.ParseBool(literalValue)
		return value
	case "string", "java.lang.String":
		return literalValue
	case "byte[]", "[]byte":
		return []byte(literalValue)
	}

	panic("CreateFromString not implemented")
}

func (eFactory *EFactoryExt) ConvertToString(eDataType EDataType, instanceValue any) string {
	if eFactory.GetEPackage() != eDataType.GetEPackage() {
		panic(fmt.Sprintf("The datatype '%v' is not a valid classifier", eDataType.GetName()))
	}

	if eEnum, _ := eDataType.(EEnum); eEnum != nil {
		result := eEnum.GetEEnumLiteralByValue(instanceValue.(int))
		if result == nil {
			panic(fmt.Sprintf("The value '%v' is not a valid enumerator of '%v'", instanceValue, eDataType.GetName()))
		}
		return result.GetLiteral()
	}

	switch getInstanceTypeName(eDataType) {
	case "float64", "java.lang.Double", "double":
		v, _ := instanceValue.(float64)
		return strconv.FormatFloat(v, 'f', -1, 64)
	case "float32", "java.lang.Float", "float":
		v, _ := instanceValue.(float32)
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case "int", "java.lang.Integer":
		v, _ := instanceValue.(int)
		return strconv.Itoa(v)
	case "uint64", "com.google.common.primitives.UnsignedLong":
		v, _ := instanceValue.(uint64)
		return strconv.FormatUint(v, 10)
	case "int64", "java.lang.Long", "java.math.BigInteger", "long":
		v, _ := instanceValue.(int64)
		return strconv.FormatInt(v, 10)
	case "int32":
		v, _ := instanceValue.(int32)
		return strconv.FormatInt(int64(v), 10)
	case "int16", "java.lang.Short", "short":
		v, _ := instanceValue.(int16)
		return strconv.FormatInt(int64(v), 10)
	case "int8":
		v, _ := instanceValue.(int8)
		return strconv.FormatInt(int64(v), 10)
	case "bool", "java.lang.Boolean", "boolean":
		v, _ := instanceValue.(bool)
		return strconv.FormatBool(v)
	case "string", "java.lang.String":
		return instanceValue.(string)
	}

	panic("ConvertToString not implemented")
}
