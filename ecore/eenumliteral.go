// Code generated by soft.generator.go. DO NOT EDIT.

// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

// EEnumLiteral is the representation of the model object 'EEnumLiteral'
type EEnumLiteral interface {
	ENamedElement

	GetValue() int
	SetValue(int)

	GetInstance() any
	SetInstance(any)

	GetLiteral() string
	SetLiteral(string)

	GetEEnum() EEnum

	// Start of user code EEnumLiteral
	// End of user code
}
