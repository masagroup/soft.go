// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

type notifierAdapterList struct {
	basicEList
	notifier *AbstractENotifier
}

type notifierNotification struct {
	AbstractNotification
	notifier *AbstractENotifier
}

func (n *notifierNotification) GetNotifier() ENotifier {
	return n.notifier.interfaces.(ENotifier)
}

func (n *notifierNotification) GetFeature() EStructuralFeature {
	return nil
}

func (n *notifierNotification) GetFeatureID() int {
	return -1
}

func newNotifierNotification(notifier *AbstractENotifier, eventType EventType, oldValue interface{}, newValue interface{}, position int) *notifierNotification {
	n := new(notifierNotification)
	n.interfaces = n
	n.eventType = eventType
	n.oldValue = oldValue
	n.newValue = newValue
	n.position = position
	n.notifier = notifier
	return n
}

func newNotifierAdapterList(notifier *AbstractENotifier) *notifierAdapterList {
	l := new(notifierAdapterList)
	l.data = []interface{}{}
	l.isUnique = true
	l.interfaces = l
	l.notifier = notifier
	return l
}

func (l *notifierAdapterList) didAdd(index int, elem interface{}) {
	notifier := l.notifier.interfaces.(ENotifier)
	elem.(EAdapter).SetTarget(notifier)
}

func (l *notifierAdapterList) didRemove(index int, elem interface{}) {
	notifier := l.notifier.interfaces.(ENotifier)
	adapter := elem.(EAdapter)
	if notifier.EDeliver() {
		adapter.NotifyChanged(newNotifierNotification(l.notifier, REMOVING_ADAPTER, elem, nil, index))
	}
	adapter.UnSetTarget(notifier)
}

type AbstractENotifier struct {
	interfaces interface{}
}

type ENotifierInternal interface {
	ENotifier
	EBasicHasAdapters() bool
	EBasicAdapters() EList
}

func (notifier *AbstractENotifier) AsENotifier() ENotifier {
	return notifier.interfaces.(ENotifier)
}

func (notifier *AbstractENotifier) AsENotifierInternal() ENotifierInternal {
	return notifier.interfaces.(ENotifierInternal)
}

// SetInterfaces ...
func (notifier *AbstractENotifier) SetInterfaces(interfaces interface{}) {
	notifier.interfaces = interfaces
}

// GetInterfaces ...
func (notifier *AbstractENotifier) GetInterfaces() interface{} {
	return notifier.interfaces
}

func (notifier *AbstractENotifier) EBasicAdapters() EList {
	return nil
}

func (notifier *AbstractENotifier) EBasicHasAdapters() bool {
	adapters := notifier.AsENotifierInternal().EBasicAdapters()
	return adapters != nil && !adapters.Empty()
}

func (notifier *AbstractENotifier) EAdapters() EList {
	return NewEmptyImmutableList()
}

func (notifier *AbstractENotifier) EDeliver() bool {
	return false
}

func (notifier *AbstractENotifier) ESetDeliver(value bool) {
	panic("operation not supported")
}

func (notifier *AbstractENotifier) ENotify(notification ENotification) {
	n := notifier.AsENotifier()
	if adapters := n.EAdapters(); adapters != nil && n.EDeliver() {
		for it := adapters.Iterator(); it.HasNext(); {
			it.Next().(EAdapter).NotifyChanged(notification)
		}
	}
}

func (notifier *AbstractENotifier) ENotificationRequired() bool {
	n := notifier.interfaces.(ENotifierInternal)
	return n.EBasicHasAdapters() && n.EDeliver()
}
