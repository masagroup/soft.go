// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

// notificationChain is an implementation of ENotificationChain interface
type notificationChain struct {
	notifications *BasicEList
}

// NewNotificationChain ...
func NewNotificationChain() *notificationChain {
	return &notificationChain{notifications: NewEmptyBasicEList()}
}

// Add Adds a notification to the chain.
func (chain *notificationChain) Add(newNotif ENotification) bool {
	if newNotif == nil {
		return false
	}
	for it := chain.notifications.Iterator(); it.HasNext(); {
		if it.Next().(ENotification).Merge(newNotif) {
			return false
		}
	}
	chain.notifications.Add(newNotif)
	return true
}

// Dispatch Dispatches each notification to the appropriate notifier via notifier.ENotify method
func (chain *notificationChain) Dispatch() {
	for it := chain.notifications.Iterator(); it.HasNext(); {
		value := it.Next().(ENotification)
		notifier := value.GetNotifier()
		if notifier != nil && value.GetEventType() != -1 {
			notifier.ENotify(value)
		}
	}
}
