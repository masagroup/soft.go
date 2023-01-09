// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

type EEnumExt struct {
	EEnumImpl
}

func newEEnumExt() *EEnumExt {
	eEnum := new(EEnumExt)
	eEnum.SetInterfaces(eEnum)
	eEnum.Initialize()
	return eEnum
}

func (eEnum *EEnumExt) GetDefaultValue() any {
	if eLiterals := eEnum.GetELiterals(); !eLiterals.Empty() {
		eLiteral := eLiterals.Get(0).(EEnumLiteral)
		return eLiteral.GetValue()
	}
	return nil
}

// GetEEnumLiteralByName default implementation
func (eEnum *EEnumExt) GetEEnumLiteralByName(name string) EEnumLiteral {
	for it := eEnum.GetELiterals().Iterator(); it.HasNext(); {
		eLiteral := it.Next().(EEnumLiteral)
		if eLiteral.GetName() == name {
			return eLiteral
		}
	}
	return nil
}

// GetEEnumLiteralByValue default implementation
func (eEnum *EEnumExt) GetEEnumLiteralByValue(value int) EEnumLiteral {
	for it := eEnum.GetELiterals().Iterator(); it.HasNext(); {
		eLiteral := it.Next().(EEnumLiteral)
		if eLiteral.GetValue() == value {
			return eLiteral
		}
	}
	return nil
}

// GetEEnumLiteralByLiteral default implementation
func (eEnum *EEnumExt) GetEEnumLiteralByLiteral(literal string) EEnumLiteral {
	for it := eEnum.GetELiterals().Iterator(); it.HasNext(); {
		eLiteral := it.Next().(EEnumLiteral)
		if eLiteral.GetLiteral() == literal {
			return eLiteral
		}
	}
	return nil
}
