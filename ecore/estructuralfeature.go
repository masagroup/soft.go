// *****************************************************************************
//
// This file is part of a MASA library or program.
// Refer to the included end-user license agreement for restrictions.
//
// Copyright (c) 2020 MASA Group
//
// *****************************************************************************

// *****************************************************************************
//
// Warning: This file was generated by soft.generator.go Generator
//
// *****************************************************************************

package ecore

import "reflect"

// EStructuralFeature is the representation of the model object 'EStructuralFeature'
type EStructuralFeature interface {
	ETypedElement

	GetContainerClass() reflect.Type

	IsChangeable() bool
	SetChangeable(bool)

	IsVolatile() bool
	SetVolatile(bool)

	IsTransient() bool
	SetTransient(bool)

	GetDefaultValueLiteral() string
	SetDefaultValueLiteral(string)

	GetDefaultValue() interface{}
	SetDefaultValue(interface{})

	IsUnsettable() bool
	SetUnsettable(bool)

	IsDerived() bool
	SetDerived(bool)

	GetFeatureID() int
	SetFeatureID(int)

	GetEContainingClass() EClass

	// Start of user code EStructuralFeature
	// End of user code
}
