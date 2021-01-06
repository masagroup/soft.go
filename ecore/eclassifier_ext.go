// *****************************************************************************
//
// This file is part of a MASA library or program.
// Refer to the included end-user license agreement for restrictions.
//
// Copyright (c) 2019 MASA Group
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
