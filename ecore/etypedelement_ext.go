// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

// ETypedElementExt is the extension of the model object 'ETypedElement'
type ETypedElementExt struct {
	ETypedElementImpl
}

func newETypedElementExt() *ETypedElementExt {
	eTypedElement := new(ETypedElementExt)
	eTypedElement.SetInterfaces(eTypedElement)
	eTypedElement.Initialize()
	return eTypedElement
}

// IsMany get the value of isMany
func (eTypedElement *ETypedElementExt) IsMany() bool {
	upper := eTypedElement.GetUpperBound()
	return upper > 1 || upper == UNBOUNDED_MULTIPLICITY
}

// IsRequired get the value of isRequired
func (eTypedElement *ETypedElementExt) IsRequired() bool {
	lower := eTypedElement.GetLowerBound()
	return lower >= 1
}
