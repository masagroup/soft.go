// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

// EEnumLiteralExt tracks whether its literal is set: as in EMF, where an
// unset literal is null, an unset literal gives the name while a literal
// set to "" stays "".
type EEnumLiteralExt struct {
	EEnumLiteralImpl
	isLiteralSet bool
}

func newEEnumLiteralExt() *EEnumLiteralExt {
	eEnumLiteral := new(EEnumLiteralExt)
	eEnumLiteral.SetInterfaces(eEnumLiteral)
	eEnumLiteral.Initialize()
	return eEnumLiteral
}

func (eEnumLiteral *EEnumLiteralExt) Initialize() {
	eEnumLiteral.EEnumLiteralImpl.Initialize()
	eEnumLiteral.isLiteralSet = false
}

func (eEnumLiteral *EEnumLiteralExt) GetLiteral() string {
	if !eEnumLiteral.isLiteralSet {
		return eEnumLiteral.GetName()
	}
	return eEnumLiteral.literal
}

func (eEnumLiteral *EEnumLiteralExt) SetLiteral(newLiteral string) {
	eEnumLiteral.isLiteralSet = true
	eEnumLiteral.EEnumLiteralImpl.SetLiteral(newLiteral)
}

func (eEnumLiteral *EEnumLiteralExt) EUnsetFromID(featureID int) {
	eEnumLiteral.EEnumLiteralImpl.EUnsetFromID(featureID)
	if featureID == EENUM_LITERAL__LITERAL {
		eEnumLiteral.isLiteralSet = false
	}
}

func (eEnumLiteral *EEnumLiteralExt) EIsSetFromID(featureID int) bool {
	if featureID == EENUM_LITERAL__LITERAL {
		return eEnumLiteral.isLiteralSet
	}
	return eEnumLiteral.EEnumLiteralImpl.EIsSetFromID(featureID)
}
