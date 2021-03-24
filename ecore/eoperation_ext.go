// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

// eOperationExt is the extension of the model object 'EOperation'
type eOperationExt struct {
	eOperationImpl
}

func newEOperationExt() *eOperationExt {
	eOperation := new(eOperationExt)
	eOperation.SetInterfaces(eOperation)
	eOperation.Initialize()
	return eOperation
}

func (eOperation *eOperationExt) IsOverrideOf(otherOperation EOperation) bool {
	otherContainingClass := otherOperation.GetEContainingClass()
	if otherContainingClass != nil && otherContainingClass.IsSuperTypeOf(eOperation.GetEContainingClass()) && otherOperation.GetName() == eOperation.GetName() {
		parameters := eOperation.GetEParameters()
		otherParameters := otherOperation.GetEParameters()
		if parameters.Size() == otherParameters.Size() {
			for i := 0; i < parameters.Size(); i++ {
				parameter := parameters.Get(i).(EParameter)
				otherParameter := otherParameters.Get(i).(EParameter)
				if parameter.GetEType() != otherParameter.GetEType() {
					return false
				}
			}
			return true
		}
	}
	return false
}
