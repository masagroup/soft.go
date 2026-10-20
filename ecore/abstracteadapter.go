// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

// AbstractEAdapter is the base implementation of an EAdapter: it holds the target of the adapter.
// An adapter embeds it and implements NotifyChanged.
type AbstractEAdapter struct {
	target ENotifier
}

// NewAbstractAdapter returns an AbstractEAdapter without target.
func NewAbstractAdapter() *AbstractEAdapter {
	return &AbstractEAdapter{}
}

// GetTarget returns the notifier the adapter receives notifications from, nil when it has none.
func (a *AbstractEAdapter) GetTarget() ENotifier {
	return a.target
}

// SetTarget sets the notifier the adapter receives notifications from: the notifier calls it when
// the adapter is added to its adapters (ENotifier.EAdapters).
func (a *AbstractEAdapter) SetTarget(notifier ENotifier) {
	a.target = notifier
}

// UnSetTarget unsets the target of the adapter when it is notifier: the notifier calls it when the
// adapter is removed from its adapters, after notifying it of the removal (REMOVING_ADAPTER) when it
// delivers notifications (EDeliver). The target is kept when it is another notifier.
func (a *AbstractEAdapter) UnSetTarget(notifier ENotifier) {
	if notifier == a.target {
		a.SetTarget(nil)
	}
}
