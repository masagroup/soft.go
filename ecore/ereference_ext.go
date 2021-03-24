// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

// eReferenceExt is the extension of the model object 'EReference'
type eReferenceExt struct {
	eReferenceImpl
	eReferenceType EClass
}

func newEReferenceExt() *eReferenceExt {
	eReference := new(eReferenceExt)
	eReference.SetInterfaces(eReference)
	eReference.Initialize()
	return eReference
}

func (eReference *eReferenceExt) IsContainer() bool {
	opposite := eReference.interfaces.(EReference).GetEOpposite()
	return opposite != nil && opposite.IsContainment()
}

func (eReference *eReferenceExt) GetEReferenceType() EClass {
	if eReference.eReferenceType == nil || eReference.eReferenceType.EIsProxy() {
		eType := eReference.GetEType()
		eReferenceType, _ := eType.(EClass)
		if eReferenceType != nil {
			eReference.eReferenceType = eReferenceType
		}
	}
	return eReference.eReferenceType
}

func (eReference *eReferenceExt) basicGetEReferenceType() EClass {
	if eReference.eReferenceType == nil {
		eType := eReference.basicGetEType()
		eReferenceType, _ := eType.(EClass)
		if eReferenceType != nil {
			eReference.eReferenceType = eReferenceType
		}
	}
	return eReference.eReferenceType
}
