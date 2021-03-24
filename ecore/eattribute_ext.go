// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

// eAttributeExt is the extension of the model object 'EAttribute'
type eAttributeExt struct {
	eAttributeImpl
}

func newEAttributeExt() *eAttributeExt {
	eAttribute := new(eAttributeExt)
	eAttribute.SetInterfaces(eAttribute)
	eAttribute.Initialize()
	return eAttribute
}

func (eAttribute *eAttributeExt) GetEAttributeType() EDataType {
	return eAttribute.GetEType().(EDataType)
}

func (eAttribute *eAttributeExt) basicGetEAttributeType() EDataType {
	return eAttribute.basicGetEType().(EDataType)
}

func (eAttribute *eAttributeExt) SetID(newIsID bool) {
	eAttribute.eAttributeImpl.SetID(newIsID)
	eClass := eAttribute.GetEContainingClass()
	if eClass != nil {
		classExt := eClass.(*eClassExt)
		classExt.setModified(ECLASS__EATTRIBUTES)
	}

}
