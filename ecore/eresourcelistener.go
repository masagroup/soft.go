// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

// EResourceListener defines callbacks when object is attached/detached from the resource
type EResourceListener interface {
	// Attached is called when an new object is attached to the resource
	Attached(object EObject)
	// Detached is called when an new object is detached from the resource
	Detached(object EObject)
}
