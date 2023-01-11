// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

// EAttributeExt is the extension of the model object 'EAttribute'
type EAttributeExt struct {
	EAttributeImpl
}

func newEAttributeExt() *EAttributeExt {
	eAttribute := new(EAttributeExt)
	eAttribute.SetInterfaces(eAttribute)
	eAttribute.Initialize()
	return eAttribute
}

func (eAttribute *EAttributeExt) GetEAttributeType() EDataType {
	return eAttribute.GetEType().(EDataType)
}

func (eAttribute *EAttributeExt) basicGetEAttributeType() EDataType {
	return eAttribute.basicGetEType().(EDataType)
}

func (eAttribute *EAttributeExt) SetID(newIsID bool) {
	eAttribute.EAttributeImpl.SetID(newIsID)
	eClass := eAttribute.GetEContainingClass()
	if eClass != nil {
		classExt := eClass.(EClassInternal)
		classExt.setModified(ECLASS__EATTRIBUTES)
	}
}
