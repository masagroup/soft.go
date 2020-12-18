// *****************************************************************************
//
// This file is part of a MASA library or program.
// Refer to the included end-user license agreement for restrictions.
//
// Copyright (c) 2019 MASA Group
//
// *****************************************************************************

package ecore

// AbstractEAdapter is a abstract implementation of EAdapter interface
type AbstractEAdapter struct {
	target ENotifier
}

// NewAbstractEAdapter Constructor
func NewAbstractEAdapter() *AbstractEAdapter {
	return &AbstractEAdapter{target: nil}
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
