//go:build ecore_check

// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import "fmt"

// CheckSameImplementation panics when owner and value are objects of two different implementations
// of the same shared package: the objects of an implementation contain, reference and share a
// resource with objects of their own implementation only. Objects of different packages, and objects
// that are not of an implementation of a shared package, are not checked. Without the build tag
// ecore_check, it checks nothing (see eobjectfactory_nocheck.go).
func CheckSameImplementation[T EObject](owner T, value any) {
	checkSameImplementation(owner, value)
}

// checkSameImplementation is CheckSameImplementation for owners of any type: the objects of ecore
// pass their interfaces, the object they are part of, without type assertion.
func checkSameImplementation(owner any, value any) {
	o, _ := owner.(EObjectFactory)
	if o == nil {
		return
	}
	v, _ := value.(EObjectFactory)
	if v == nil {
		return
	}
	ownerFactory, valueFactory := o.EFactory(), v.EFactory()
	if ownerFactory == valueFactory {
		return
	}
	if ePackage := ownerFactory.GetEPackage(); ePackage == valueFactory.GetEPackage() {
		panic(fmt.Sprintf("cannot mix implementations of the package '%s': %T and %T", ePackage.GetNsURI(), owner, value))
	}
}

// checkSameImplementation panics when eObject is an object of an implementation of a shared package
// and the resource holds a root of another implementation of that package (see CheckSameImplementation).
// The roots of a package are of one implementation: eObject is compared with the first root of its
// package only.
func (rc *resourceContents) checkSameImplementation(eObject any) {
	objectFactory, _ := eObject.(EObjectFactory)
	if objectFactory == nil {
		return
	}
	ePackage := objectFactory.EFactory().GetEPackage()
	for _, root := range rc.data {
		if rootFactory, _ := root.(EObjectFactory); rootFactory != nil && root != eObject && rootFactory.EFactory().GetEPackage() == ePackage {
			checkSameImplementation(root, eObject)
			return
		}
	}
}
