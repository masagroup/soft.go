// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import "sync"

// AbstractEAdapter is a abstract implementation of EAdapter interface
type AbstractEAdapter struct {
	// mutex guards target: one adapter is attached to many notifiers, so concurrent
	// loads write this field from several goroutines. Never held across an attach
	// walk, which re-enters SetTarget for every child.
	mutex  sync.Mutex
	target ENotifier
}

func NewAbstractAdapter() *AbstractEAdapter {
	return &AbstractEAdapter{}
}

// GetTarget Returns the target from which the AbstractEAdapter receives notification.
func (a *AbstractEAdapter) GetTarget() ENotifier {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.target
}

// SetTarget Sets the target from which the AbstractEAdapter will receive notification.
func (a *AbstractEAdapter) SetTarget(notifier ENotifier) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.target = notifier
}

func (a *AbstractEAdapter) UnSetTarget(notifier ENotifier) {
	// Compare and clear under one acquisition; SetTarget would re-enter the lock.
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if notifier == a.target {
		a.target = nil
	}
}
