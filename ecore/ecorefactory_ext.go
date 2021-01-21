// *****************************************************************************
//
// This file is part of a MASA library or program.
// Refer to the included end-user license agreement for restrictions.
//
// Copyright (c) 2019 MASA Group
//
// *****************************************************************************

package ecore

import (
	"fmt"
	"strconv"
	"time"
)

type ecoreFactoryExt struct {
	ecoreFactoryImpl
}

func newEcoreFactoryExt() *ecoreFactoryExt {
	factory := new(ecoreFactoryExt)
	factory.SetInterfaces(factory)
	factory.Initialize()
	return factory
}

func (factory *ecoreFactoryExt) createEBigDecimalFromString(eDataType EDataType, literalValue string) interface{} {
	value, _ := strconv.ParseFloat(literalValue, 64)
	return value
}

func (factory *ecoreFactoryExt) createEBigIntegerFromString(eDataType EDataType, literalValue string) interface{} {
	value, _ := strconv.ParseInt(literalValue, 10, 64)
	return value
}

func (factory *ecoreFactoryExt) createEBooleanFromString(dataType EDataType, literalValue string) interface{} {
	value, _ := strconv.ParseBool(literalValue)
	return value
}

func (factory *ecoreFactoryExt) createEBooleanObjectFromString(eDataType EDataType, literalValue string) interface{} {
	return factory.createEBooleanFromString(eDataType, literalValue)
}

func (factory *ecoreFactoryExt) createEByteFromString(eDataType EDataType, literalValue string) interface{} {
	if len(literalValue) == 0 {
		return "golang\u0000"
	} else {
		return []byte(literalValue)[0]
	}
}

func (factory *ecoreFactoryExt) createEByteObjectFromString(eDataType EDataType, literalValue string) interface{} {
	return factory.createEByteFromString(eDataType, literalValue)
}

func (factory *ecoreFactoryExt) createEByteArrayFromString(eDataType EDataType, literalValue string) interface{} {
	return []byte(literalValue)
}

func (factory *ecoreFactoryExt) createECharFromString(dataType EDataType, literalValue string) interface{} {
	if len(literalValue) == 0 {
		return "golang\u0000"
	} else {
		return literalValue[0]
	}
}

func (factory *ecoreFactoryExt) createECharacterObjectFromString(dataType EDataType, literalValue string) interface{} {
	return factory.createECharFromString(dataType, literalValue)
}

func (factory *ecoreFactoryExt) createEDateFromString(dataType EDataType, literalValue string) interface{} {
	t, _ := time.Parse(dateFormat, literalValue)
	return &t
}

func (factory *ecoreFactoryExt) createEDoubleFromString(dataType EDataType, literalValue string) interface{} {
	value, _ := strconv.ParseFloat(literalValue, 64)
	return value
}

func (factory *ecoreFactoryExt) createEDoubleObjectFromString(dataType EDataType, literalValue string) interface{} {
	return factory.createEDoubleFromString(dataType, literalValue)
}

func (factory *ecoreFactoryExt) createEFloatFromString(dataType EDataType, literalValue string) interface{} {
	value, _ := strconv.ParseFloat(literalValue, 32)
	return float32(value)
}

func (factory *ecoreFactoryExt) createEFloatObjectFromString(dataType EDataType, literalValue string) interface{} {
	return factory.createEFloatFromString(dataType, literalValue)
}

func (factory *ecoreFactoryExt) createEIntFromString(dataType EDataType, literalValue string) interface{} {
	value, _ := strconv.Atoi(literalValue)
	return int(value)
}

func (factory *ecoreFactoryExt) createEIntegerObjectFromString(dataType EDataType, literalValue string) interface{} {
	return factory.createEIntFromString(dataType, literalValue)
}

func (factory *ecoreFactoryExt) createELongFromString(dataType EDataType, literalValue string) interface{} {
	value, _ := strconv.ParseInt(literalValue, 10, 64)
	return value
}

func (factory *ecoreFactoryExt) createELongObjectFromString(dataType EDataType, literalValue string) interface{} {
	return factory.createELongFromString(dataType, literalValue)
}

func (factory *ecoreFactoryExt) createEShortFromString(dataType EDataType, literalValue string) interface{} {
	value, _ := strconv.ParseInt(literalValue, 10, 16)
	return int16(value)
}

func (factory *ecoreFactoryExt) createEShortObjectFromString(dataType EDataType, literalValue string) interface{} {
	return factory.createEShortFromString(dataType, literalValue)
}

func (factory *ecoreFactoryExt) createEStringFromString(dataType EDataType, literalValue string) interface{} {
	return literalValue
}

func (factory *ecoreFactoryExt) convertEBigDecimalToString(dataType EDataType, instanceValue interface{}) string {
	v, _ := instanceValue.(float64)
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func (factory *ecoreFactoryExt) convertEBigIntegerToString(dataType EDataType, instanceValue interface{}) string {
	v, _ := instanceValue.(int64)
	return strconv.FormatInt(v, 10)
}

func (factory *ecoreFactoryExt) convertEBooleanToString(dataType EDataType, instanceValue interface{}) string {
	v, _ := instanceValue.(bool)
	return strconv.FormatBool(v)
}

func (factory *ecoreFactoryExt) convertEBooleanObjectToString(dataType EDataType, instanceValue interface{}) string {
	return factory.convertEBooleanToString(dataType, instanceValue)
}

func (factory *ecoreFactoryExt) convertEByteToString(eDataType EDataType, instanceValue interface{}) string {
	b := instanceValue.(byte)
	return string([]byte{b})
}

func (factory *ecoreFactoryExt) convertEByteObjectToString(eDataType EDataType, instanceValue interface{}) string {
	return factory.convertEByteToString(eDataType, instanceValue)
}

func (factory *ecoreFactoryExt) convertEByteArrayToString(eDataType EDataType, instanceValue interface{}) string {
	b := instanceValue.([]byte)
	return string(b)
}

func (factory *ecoreFactoryExt) convertECharToString(dataType EDataType, instanceValue interface{}) string {
	return fmt.Sprintf("%c", instanceValue)
}

func (factory *ecoreFactoryExt) convertECharacterObjectToString(eDataType EDataType, instanceValue interface{}) string {
	return factory.convertECharToString(eDataType, instanceValue)
}

const (
	dateFormat string = "2006-01-02T15:04:05.999Z"
)

func (factory *ecoreFactoryExt) convertEDateToString(dataType EDataType, instanceValue interface{}) string {
	t := instanceValue.(*time.Time)
	return t.Format(dateFormat)
}

func (factory *ecoreFactoryExt) convertEDoubleToString(dataType EDataType, instanceValue interface{}) string {
	v, _ := instanceValue.(float64)
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func (factory *ecoreFactoryExt) convertEDoubleObjectToString(eDataType EDataType, instanceValue interface{}) string {
	return factory.convertEDoubleToString(eDataType, instanceValue)
}

func (factory *ecoreFactoryExt) convertEFloatToString(dataType EDataType, instanceValue interface{}) string {
	v, _ := instanceValue.(float32)
	return strconv.FormatFloat(float64(v), 'f', -1, 32)
}

func (factory *ecoreFactoryExt) convertEFloatObjectToString(eDataType EDataType, instanceValue interface{}) string {
	return factory.convertEFloatToString(eDataType, instanceValue)
}

func (factory *ecoreFactoryExt) convertEIntToString(dataType EDataType, instanceValue interface{}) string {
	return fmt.Sprintf("%d", instanceValue)
}

func (factory *ecoreFactoryExt) convertEIntegerObjectToString(eDataType EDataType, instanceValue interface{}) string {
	return factory.convertEIntToString(eDataType, instanceValue)
}

func (factory *ecoreFactoryExt) convertELongToString(dataType EDataType, instanceValue interface{}) string {
	v, _ := instanceValue.(int64)
	return strconv.FormatInt(v, 10)
}

func (factory *ecoreFactoryExt) convertELongObjectToString(eDataType EDataType, instanceValue interface{}) string {
	return factory.convertELongToString(eDataType, instanceValue)
}

func (factory *ecoreFactoryExt) convertEShortToString(dataType EDataType, instanceValue interface{}) string {
	v, _ := instanceValue.(int16)
	return strconv.FormatInt(int64(v), 10)
}

func (factory *ecoreFactoryExt) convertEShortObjectToString(eDataType EDataType, instanceValue interface{}) string {
	return factory.convertEShortToString(eDataType, instanceValue)
}

func (factory *ecoreFactoryExt) convertEStringToString(dataType EDataType, instanceValue interface{}) string {
	v, _ := instanceValue.(string)
	return v
}
