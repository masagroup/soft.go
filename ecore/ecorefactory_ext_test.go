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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFactoryDate(t *testing.T) {
	factory := newEcoreFactoryExt()
	{
		mockEDataType := &MockEDataType{}
		mockEDataType.On("GetClassifierID").Return(EDATE)
		{

			date := factory.CreateFromString(mockEDataType, "2020-05-12T17:33:10.77Z")
			expected := time.Date(2020, time.May, 12, 17, 33, 10, 770000000, time.UTC)
			assert.Equal(t, &expected, date)
		}
		{
			date := factory.CreateFromString(mockEDataType, "2007-06-02T10:26:13.000Z")
			expected := time.Date(2007, time.June, 2, 10, 26, 13, 0, time.UTC)
			assert.Equal(t, &expected, date)
		}
		{
			date := time.Date(2020, time.May, 12, 17, 33, 10, 770000000, time.UTC)
			dateStr := factory.ConvertToString(mockEDataType, &date)
			expected := "2020-05-12T17:33:10.77Z"
			assert.Equal(t, expected, dateStr)
		}
	}

}

func TestFactoryBoolean(t *testing.T) {
	factory := newEcoreFactoryExt()
	assert.Equal(t, factory.convertEBooleanToString(nil, factory.createEBooleanFromString(nil, "true")), "true")
	assert.Equal(t, factory.convertEBooleanToString(nil, factory.createEBooleanFromString(nil, "false")), "false")
}

func TestFactoryBooleanObject(t *testing.T) {
	factory := newEcoreFactoryExt()
	assert.Equal(t, factory.convertEBooleanObjectToString(nil, factory.createEBooleanObjectFromString(nil, "true")), "true")
	assert.Equal(t, factory.convertEBooleanObjectToString(nil, factory.createEBooleanObjectFromString(nil, "false")), "false")
}

func TestFactoryByte(t *testing.T) {
	factory := newEcoreFactoryExt()
	assert.Equal(t, "golang\u0000", factory.createEByteFromString(nil, ""))
	assert.Equal(t, factory.convertEByteToString(nil, factory.createEByteFromString(nil, "a")), "a")
}

func TestFactoryByteObject(t *testing.T) {
	factory := newEcoreFactoryExt()
	assert.Equal(t, factory.convertEByteObjectToString(nil, factory.createEByteObjectFromString(nil, "a")), "a")
}

func TestFactoryByteArray(t *testing.T) {
	factory := newEcoreFactoryExt()
	assert.Equal(t, []byte("ab"), factory.createEByteArrayFromString(nil, "ab"))
	assert.Equal(t, factory.convertEByteArrayToString(nil, factory.createEByteArrayFromString(nil, "ab")), "ab")
}

func TestFactoryChar(t *testing.T) {
	factory := newEcoreFactoryExt()
	assert.Equal(t, factory.convertECharToString(nil, factory.createECharFromString(nil, "e")), "e")
}

func TestFactoryCharObject(t *testing.T) {
	factory := newEcoreFactoryExt()
	assert.Equal(t, factory.convertECharacterObjectToString(nil, factory.createECharacterObjectFromString(nil, "e")), "e")
}

func TestFactoryDouble(t *testing.T) {
	factory := newEcoreFactoryExt()
	assert.Equal(t, factory.convertEDoubleToString(nil, factory.createEDoubleFromString(nil, "4.987453")), "4.987453")
}

func TestFactoryDoubleObject(t *testing.T) {
	factory := newEcoreFactoryExt()
	assert.Equal(t, factory.convertEDoubleObjectToString(nil, factory.createEDoubleObjectFromString(nil, "4.987453")), "4.987453")
}

func TestFactoryFloat(t *testing.T) {
	factory := newEcoreFactoryExt()
	assert.Equal(t, factory.convertEFloatToString(nil, factory.createEFloatFromString(nil, "4.987453")), "4.987453")
}

func TestFactoryFloatObject(t *testing.T) {
	factory := newEcoreFactoryExt()
	assert.Equal(t, factory.convertEFloatObjectToString(nil, factory.createEFloatObjectFromString(nil, "4.987453")), "4.987453")
}

func TestFactoryInt(t *testing.T) {
	factory := newEcoreFactoryExt()
	assert.Equal(t, factory.convertEIntToString(nil, factory.createEIntFromString(nil, "50000000")), "50000000")
}

func TestFactoryIntObject(t *testing.T) {
	factory := newEcoreFactoryExt()
	assert.Equal(t, factory.convertEIntegerObjectToString(nil, factory.createEIntegerObjectFromString(nil, "50000000")), "50000000")
}

func TestFactoryLong(t *testing.T) {
	factory := newEcoreFactoryExt()
	assert.Equal(t, factory.convertELongToString(nil, factory.createELongFromString(nil, "5000000000000")), "5000000000000")
}

func TestFactoryLongObject(t *testing.T) {
	factory := newEcoreFactoryExt()
	assert.Equal(t, factory.convertELongObjectToString(nil, factory.createELongObjectFromString(nil, "5000000000000")), "5000000000000")
}

func TestFactoryShort(t *testing.T) {
	factory := newEcoreFactoryExt()
	assert.Equal(t, "5000", factory.convertEShortToString(nil, factory.createEShortFromString(nil, "5000")))
}

func TestFactoryShortObject(t *testing.T) {
	factory := newEcoreFactoryExt()
	assert.Equal(t, "5000", factory.convertEShortObjectToString(nil, factory.createEShortObjectFromString(nil, "5000")))
}

func TestFactoryString(t *testing.T) {
	factory := newEcoreFactoryExt()
	assert.Equal(t, factory.convertEStringToString(nil, factory.createEStringFromString(nil, "Hi I'm a string")), "Hi I'm a string")
}

func TestFactoryBigDecimal(t *testing.T) {
	factory := newEcoreFactoryExt()
	assert.Equal(t, factory.convertEBigDecimalToString(nil, factory.createEBigDecimalFromString(nil, "10.0")), "10")
	assert.Equal(t, factory.convertEBigDecimalToString(nil, factory.createEBigDecimalFromString(nil, "10.1")), "10.1")
}

func TestFactoryBigInteger(t *testing.T) {
	factory := newEcoreFactoryExt()
	assert.Equal(t, factory.convertEBigIntegerToString(nil, factory.createEBigIntegerFromString(nil, "10")), "10")
}
