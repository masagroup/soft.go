// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

// EClassifierExt is the extesnion of the model object 'EClassifier'
type EClassifierExt struct {
	EClassifierImpl
}

func newEClassifierExt() *EClassifierExt {
	eClassifier := new(EClassifierExt)
	eClassifier.SetInterfaces(eClassifier)
	eClassifier.Initialize()
	return eClassifier
}

func (eClassifier *EClassifierExt) initClassifierID() int {
	if ePackage := eClassifier.GetEPackage(); ePackage != nil {
		return ePackage.GetEClassifiers().IndexOf(eClassifier.asEClassifier())
	}
	return -1
}

func (eClassifier *EClassifierExt) GetDefaultValue() any {
	return nil
}

func (eClassifier *EClassifierExt) GetInstanceTypeName() string {
	return eClassifier.GetInstanceClassName()
}

// SetInstanceTypeName set the value of instanceTypeName
func (eClassifier *EClassifierExt) SetInstanceTypeName(newInstanceTypeName string) {
	eClassifier.SetInstanceClassName(newInstanceTypeName)
}
