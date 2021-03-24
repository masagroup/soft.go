// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

// EFactoryExt is the extension of the model object 'EFactory'
type EFactoryExt struct {
	eFactoryImpl
}

// NewEFactoryExt ...
func NewEFactoryExt() *EFactoryExt {
	eFactory := new(EFactoryExt)
	eFactory.SetInterfaces(eFactory)
	eFactory.Initialize()
	return eFactory
}

// Create ...
func (eFactory *EFactoryExt) Create(eClass EClass) EObject {
	if eFactory.GetEPackage() != eClass.GetEPackage() || eClass.IsAbstract() {
		panic("The class '" + eClass.GetName() + "' is not a valid classifier")
	}
	eObject := NewDynamicEObjectImpl()
	eObject.SetEClass(eClass)
	return eObject
}
