// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

type eEnumLiteralExt struct {
	eEnumLiteralImpl
}

func newEEnumLiteralExt() *eEnumLiteralExt {
	eEnumLiteral := new(eEnumLiteralExt)
	eEnumLiteral.SetInterfaces(eEnumLiteral)
	eEnumLiteral.Initialize()
	return eEnumLiteral
}

func (eEnumLiteral *eEnumLiteralExt) GetLiteral() string {
	if len(eEnumLiteral.literal) == 0 {
		return eEnumLiteral.GetName()
	}
	return eEnumLiteral.literal
}
