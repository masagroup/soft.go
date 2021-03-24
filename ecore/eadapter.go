// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

// EAdapter is a receiver of notifications.
// An EAdapter is typically associated with a Notifier
type EAdapter interface {
	// NotifyChanged Notifies that a change to some feature has occurred.
	NotifyChanged(notification ENotification)

	// GetTarget Returns the target from which the adapter receives notification.
	GetTarget() ENotifier

	// SetTarget Sets the target from which the adapter will receive notification.
	SetTarget(ENotifier)

	// UnSetTarget Unsets the target from which the adapter will receive notification.
	UnSetTarget(ENotifier)
}
