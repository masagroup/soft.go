// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

// ENotificationChain is an accumulator of notifications.
// As notifications are produced,they are accumulated in a chain,
// and possibly even merged, before finally being dispatched to the notifier.
type ENotificationChain interface {

	// Add Adds a notification to the chain.
	Add(ENotification) bool

	// Dispatch Dispatches each notification to the appropriate notifier via notifier.ENotify method
	Dispatch()
}
