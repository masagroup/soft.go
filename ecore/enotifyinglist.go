// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

// ENotifyingList ...
type ENotifyingList interface {
	EList

	GetNotifier() ENotifier

	GetFeature() EStructuralFeature

	GetFeatureID() int

	AddWithNotification(object any, notifications ENotificationChain) ENotificationChain

	RemoveWithNotification(object any, notifications ENotificationChain) ENotificationChain

	SetWithNotification(index int, object any, notifications ENotificationChain) ENotificationChain
}

// EUnsettableList is the list of a many-valued feature, which tells whether the feature is set and
// unsets it, as EMF InternalEList.Unsettable. The list of a feature that is not unsettable is set
// when it is not empty, and Unset clears it. The list of an unsettable feature is set from its first
// change, a clear of an empty list included, until Unset, which clears it and notifies UNSET.
type EUnsettableList interface {
	EList

	IsSet() bool

	Unset()
}

// isSetList returns true when the list of a many-valued feature is set: as an EUnsettableList
// tells, else when it is not empty.
func isSetList(list EList) bool {
	if unsettable, _ := list.(EUnsettableList); unsettable != nil {
		return unsettable.IsSet()
	}
	return !list.Empty()
}

// unsetList unsets the list of a many-valued feature: as an EUnsettableList does, else clearing it.
func unsetList(list EList) {
	if unsettable, _ := list.(EUnsettableList); unsettable != nil {
		unsettable.Unset()
	} else {
		list.Clear()
	}
}
