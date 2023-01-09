// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

type EEnumLiteralExt struct {
	EEnumLiteralImpl
}

func newEEnumLiteralExt() *EEnumLiteralExt {
	eEnumLiteral := new(EEnumLiteralExt)
	eEnumLiteral.SetInterfaces(eEnumLiteral)
	eEnumLiteral.Initialize()
	return eEnumLiteral
}

func (eEnumLiteral *EEnumLiteralExt) GetLiteral() string {
	if len(eEnumLiteral.literal) == 0 {
		return eEnumLiteral.GetName()
	}
	return eEnumLiteral.literal
}
