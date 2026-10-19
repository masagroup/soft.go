//go:build !ecore_check

// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

// CheckSameImplementation checks nothing without the build tag ecore_check: the calls of an empty
// function are inlined away (see eobjectfactory_check.go).
func CheckSameImplementation[T EObject](owner T, value any) {}

// checkSameImplementation checks nothing without the build tag ecore_check (see
// eobjectfactory_check.go).
func checkSameImplementation(owner any, value any) {}

// checkSameImplementation checks nothing without the build tag ecore_check (see
// eobjectfactory_check.go).
func (rc *resourceContents) checkSameImplementation(eObject any) {}
