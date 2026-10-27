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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEEnumExt_GetEEnumLiteralByName(t *testing.T) {
	mockEnumLiteral := NewMockEEnumLiteral(t)
	eEnum := newEEnumExt()

	mockEnumLiteral.EXPECT().EInverseAdd(eEnum, EENUM_LITERAL__EENUM, nil).Return(nil).Once()
	eEnum.GetELiterals().Add(mockEnumLiteral)
	mock.AssertExpectationsForObjects(t, mockEnumLiteral)

	mockEnumLiteral.EXPECT().GetName().Return("mockEnumLiteral").Once()
	assert.Nil(t, eEnum.GetEEnumLiteralByName("test"))
	mock.AssertExpectationsForObjects(t, mockEnumLiteral)

	mockEnumLiteral.EXPECT().GetName().Return("mockEnumLiteral").Once()
	assert.Equal(t, mockEnumLiteral, eEnum.GetEEnumLiteralByName("mockEnumLiteral"))
	mock.AssertExpectationsForObjects(t, mockEnumLiteral)
}

func TestEEnumExt_GetEEnumLiteralByValue(t *testing.T) {
	mockEnumLiteral := NewMockEEnumLiteral(t)
	eEnum := newEEnumExt()

	mockEnumLiteral.EXPECT().EInverseAdd(eEnum, EENUM_LITERAL__EENUM, nil).Return(nil).Once()
	eEnum.GetELiterals().Add(mockEnumLiteral)
	mock.AssertExpectationsForObjects(t, mockEnumLiteral)

	mockEnumLiteral.EXPECT().GetValue().Return(0).Once()
	assert.Nil(t, eEnum.GetEEnumLiteralByValue(1))
	mock.AssertExpectationsForObjects(t, mockEnumLiteral)

	mockEnumLiteral.EXPECT().GetValue().Return(0).Once()
	assert.Equal(t, mockEnumLiteral, eEnum.GetEEnumLiteralByValue(0))
	mock.AssertExpectationsForObjects(t, mockEnumLiteral)
}

func TestEEnumExt_GetEEnumLiteralByLiteral(t *testing.T) {
	mockEnumLiteral := NewMockEEnumLiteral(t)
	eEnum := newEEnumExt()

	mockEnumLiteral.EXPECT().EInverseAdd(eEnum, EENUM_LITERAL__EENUM, nil).Return(nil).Once()
	eEnum.GetELiterals().Add(mockEnumLiteral)
	mock.AssertExpectationsForObjects(t, mockEnumLiteral)

	mockEnumLiteral.EXPECT().GetLiteral().Return("no literal").Once()
	assert.Nil(t, eEnum.GetEEnumLiteralByLiteral("literal"))
	mock.AssertExpectationsForObjects(t, mockEnumLiteral)

	mockEnumLiteral.EXPECT().GetLiteral().Return("literal").Once()
	assert.Equal(t, mockEnumLiteral, eEnum.GetEEnumLiteralByLiteral("literal"))
	mock.AssertExpectationsForObjects(t, mockEnumLiteral)
}

// color is the Go type a generated package gives an enum.
type color int

// colorFactory reads the literals of an enum as a generated factory does: as values of the Go
// type of the enum.
type colorFactory struct {
	*EFactoryExt
}

func (f *colorFactory) CreateFromString(eDataType EDataType, literalValue string) any {
	return color(f.EFactoryExt.CreateFromString(eDataType, literalValue).(int))
}

// The default value of an enum, and of an attribute of the enum, is its first literal as the
// factory of its package reads it: the Go type of the enum for a generated package, which the
// setters of its objects take, the int value of the literal for a dynamic package.
func TestEEnumExt_GetDefaultValueOfAPackage(t *testing.T) {
	ePackage := GetFactory().CreateEPackage()
	eEnum := GetFactory().CreateEEnum()
	eEnum.SetName("Color")
	for i, name := range []string{"red", "green"} {
		eLiteral := GetFactory().CreateEEnumLiteral()
		eLiteral.SetName(name)
		eLiteral.SetValue(i + 1)
		eEnum.GetELiterals().Add(eLiteral)
	}
	ePackage.GetEClassifiers().Add(eEnum)
	eAttribute := GetFactory().CreateEAttribute()
	eAttribute.SetEType(eEnum)

	ePackage.SetEFactoryInstance(GetFactory().CreateEFactory())
	assert.Equal(t, 1, eEnum.GetDefaultValue(), "a dynamic package")

	ePackage.SetEFactoryInstance(&colorFactory{EFactoryExt: newEFactoryExt()})
	assert.Equal(t, color(1), eEnum.GetDefaultValue(), "a generated package")
	assert.Equal(t, color(1), eAttribute.GetDefaultValue(), "an attribute without default literal")
	eAttribute.SetDefaultValueLiteral("red")
	assert.Equal(t, color(1), eAttribute.GetDefaultValue(), "an attribute with its default literal")
}

func TestEEnumExt_GetDefaultValue(t *testing.T) {
	eEnum := newEEnumExt()

	mockEnumLiteral := NewMockEEnumLiteral(t)
	mockEnumLiteral.EXPECT().EInverseAdd(eEnum, EENUM_LITERAL__EENUM, nil).Return(nil).Once()
	eEnum.GetELiterals().Add(mockEnumLiteral)
	mock.AssertExpectationsForObjects(t, mockEnumLiteral)

	mockValue := 2
	mockEnumLiteral.EXPECT().GetValue().Return(mockValue).Once()
	assert.Equal(t, mockValue, eEnum.GetDefaultValue())
	mock.AssertExpectationsForObjects(t, mockEnumLiteral)
}
