// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package binarycodec

import (
	"github.com/masagroup/soft.go/codec/xmlcodec"
	"github.com/masagroup/soft.go/ecore"
)

func diagnosticError(errors ecore.EList) string {
	if errors.Empty() {
		return ""
	} else {
		return errors.Get(0).(ecore.EDiagnostic).GetMessage()
	}
}

func loadPackage(packageFileName string) ecore.EPackage {
	xmiProcessor := xmlcodec.NewXMIProcessor()
	eResource := xmiProcessor.Load(ecore.NewURI("testdata/" + packageFileName))
	if eResource.IsLoaded() && eResource.GetContents().Size() > 0 {
		ePackage, _ := eResource.GetContents().Get(0).(ecore.EPackage)
		return ePackage
	} else {
		return nil
	}
}
