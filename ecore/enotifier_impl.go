// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

type ENotifierImpl struct {
	AbstractENotifier
	adapters EList
	deliver  bool
}

func (notifier *ENotifierImpl) Initialize() {
	notifier.deliver = true
	notifier.adapters = nil
}

func (notifier *ENotifierImpl) EAdapters() EList {
	if notifier.adapters == nil {
		notifier.adapters = newNotifierAdapterList(&notifier.AbstractENotifier)
	}
	return notifier.adapters
}

func (notifier *ENotifierImpl) EBasicAdapters() EList {
	return notifier.adapters
}

func (notifier *ENotifierImpl) EDeliver() bool {
	return notifier.deliver
}

func (notifier *ENotifierImpl) ESetDeliver(deliver bool) {
	notifier.deliver = deliver
}
