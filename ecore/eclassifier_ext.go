// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

// eClassifierExt is the extesnion of the model object 'EClassifier'
type eClassifierExt struct {
	eClassifierImpl
}

func newEClassifierExt() *eClassifierExt {
	eClassifier := new(eClassifierExt)
	eClassifier.SetInterfaces(eClassifier)
	eClassifier.Initialize()
	return eClassifier
}

func (eClassifier *eClassifierExt) initClassifierID() int {
	if ePackage := eClassifier.GetEPackage(); ePackage != nil {
		return ePackage.GetEClassifiers().IndexOf(eClassifier.asEClassifier())
	}
	return -1
}

func (eClassifier *eClassifierExt) GetDefaultValue() interface{} {
	return nil
}
