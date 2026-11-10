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

// colorEnum returns an enum of a package with literals of values 1, 2, and an attribute of it
// without default literal.
func colorEnum() (EPackage, EEnum, EAttribute) {
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
	return ePackage, eEnum, eAttribute
}

// The default value of an enum is kept while its factory and its first literal do not change, and
// read again once they do.
func TestEEnumExt_GetDefaultValueKept(t *testing.T) {
	ePackage, eEnum, _ := colorEnum()
	ePackage.SetEFactoryInstance(GetFactory().CreateEFactory())
	assert.Equal(t, 1, eEnum.GetDefaultValue())
	assert.Equal(t, 1, eEnum.GetDefaultValue(), "kept")

	ePackage.SetEFactoryInstance(&colorFactory{EFactoryExt: newEFactoryExt()})
	assert.Equal(t, color(1), eEnum.GetDefaultValue(), "another factory")

	first := eEnum.GetELiterals().Get(0).(EEnumLiteral)
	first.SetValue(5)
	assert.Equal(t, color(5), eEnum.GetDefaultValue(), "the value of the first literal changed")

	eLiteral := GetFactory().CreateEEnumLiteral()
	eLiteral.SetName("blue")
	eLiteral.SetValue(3)
	eEnum.GetELiterals().Insert(0, eLiteral)
	assert.Equal(t, color(3), eEnum.GetDefaultValue(), "another first literal")
}

// BenchmarkEEnumExt_GetDefaultValue reads the default value of an attribute of an enum without
// default literal, as each access to an unset attribute of the enum does.
func BenchmarkEEnumExt_GetDefaultValue(b *testing.B) {
	for _, generated := range []bool{false, true} {
		b.Run(map[bool]string{false: "dynamic", true: "generated"}[generated], func(b *testing.B) {
			ePackage, _, eAttribute := colorEnum()
			if generated {
				ePackage.SetEFactoryInstance(&colorFactory{EFactoryExt: newEFactoryExt()})
			} else {
				ePackage.SetEFactoryInstance(GetFactory().CreateEFactory())
			}
			b.ReportAllocs()
			for b.Loop() {
				if eAttribute.GetDefaultValue() == nil {
					b.Fatal("no default value")
				}
			}
		})
	}
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
