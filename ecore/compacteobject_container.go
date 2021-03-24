// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

type CompactEObjectContainer struct {
	CompactEObjectImpl
	container EObject
}

func (o *CompactEObjectContainer) ESetInternalContainer(newContainer EObject, newContainerFeatureID int) {
	o.container = newContainer
	o.flags = newContainerFeatureID<<16 | (o.flags & 0x00FF)
}

func (o *CompactEObjectContainer) EInternalContainer() EObject {
	return o.container
}
