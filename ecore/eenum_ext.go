// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

type eEnumExt struct {
	eEnumImpl
}

func newEEnumExt() *eEnumExt {
	eEnum := new(eEnumExt)
	eEnum.SetInterfaces(eEnum)
	eEnum.Initialize()
	return eEnum
}

// GetEEnumLiteralByName default implementation
func (eEnum *eEnumExt) GetEEnumLiteralByName(name string) EEnumLiteral {
	for it := eEnum.GetELiterals().Iterator(); it.HasNext(); {
		eLiteral := it.Next().(EEnumLiteral)
		if eLiteral.GetName() == name {
			return eLiteral
		}
	}
	return nil
}

// GetEEnumLiteralByValue default implementation
func (eEnum *eEnumExt) GetEEnumLiteralByValue(value int) EEnumLiteral {
	for it := eEnum.GetELiterals().Iterator(); it.HasNext(); {
		eLiteral := it.Next().(EEnumLiteral)
		if eLiteral.GetValue() == value {
			return eLiteral
		}
	}
	return nil
}

// GetEEnumLiteralByLiteral default implementation
func (eEnum *eEnumExt) GetEEnumLiteralByLiteral(literal string) EEnumLiteral {
	for it := eEnum.GetELiterals().Iterator(); it.HasNext(); {
		eLiteral := it.Next().(EEnumLiteral)
		if eLiteral.GetLiteral() == literal {
			return eLiteral
		}
	}
	return nil
}
