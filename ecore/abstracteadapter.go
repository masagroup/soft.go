// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

// AbstractEAdapter is a abstract implementation of EAdapter interface
type AbstractEAdapter struct {
	target ENotifier
}

func NewAbstractAdapter() *AbstractEAdapter {
	return &AbstractEAdapter{}
}

// GetTarget Returns the target from which the AbstractEAdapter receives notification.
func (a *AbstractEAdapter) GetTarget() ENotifier {
	return a.target
}

// SetTarget Sets the target from which the AbstractEAdapter will receive notification.
func (a *AbstractEAdapter) SetTarget(notifier ENotifier) {
	a.target = notifier
}

func (a *AbstractEAdapter) UnSetTarget(notifier ENotifier) {
	if notifier == a.target {
		a.SetTarget(nil)
	}
}
