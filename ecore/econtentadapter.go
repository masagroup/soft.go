// *****************************************************************************
//
// This file is part of a MASA library or program.
// Refer to the included end-user license agreement for restrictions.
//
// Copyright (c) 2019 MASA Group
//
// *****************************************************************************

package ecore

// An adapter that maintains itself as an adapter for all contained objects
// as they come and go.
// It can be installed for an {@link EObject}, a {@link Resource}, or a {@link ResourceSet}.
type EContentAdapter struct {
	*Adapter
	interfaces interface{}
}

func NewEContentAdapter() *EContentAdapter {
	ca := new(EContentAdapter)
	ca.Adapter = NewAdapter()
	ca.interfaces = ca
	return ca
}

func (adapter *EContentAdapter) SetInterfaces(interfaces interface{}) {
	adapter.interfaces = interfaces
}

func (adapter *EContentAdapter) NotifyChanged(notification ENotification) {
	adapter.selfAdapt(notification)
}

func (adapter *EContentAdapter) SetTarget(notifier ENotifier) {
	adapter.Adapter.SetTarget(notifier)
	var it EIterator
	switch t := notifier.(type) {
	case EObject:
		it = t.EContents().Iterator()
	case EResource:
		it = t.GetContents().Iterator()
	case EResourceSet:
		it = t.GetResources().Iterator()
	}
	for it.HasNext() {
		n := it.Next().(ENotifier)
		adapter.addAdapter(n)
	}
}

func (adapter *EContentAdapter) UnSetTarget(notifier ENotifier) {
	adapter.Adapter.UnSetTarget(notifier)
	switch t := notifier.(type) {
	case EObject:
		for it := t.EContents().Iterator(); it.HasNext(); {
			notifier := it.Next().(ENotifier)
			adapter.removeAdapterWithChecks(notifier, false, true)
		}
	case EResource:
		for it := t.GetContents().Iterator(); it.HasNext(); {
			notifier := it.Next().(ENotifier)
			adapter.removeAdapterWithChecks(notifier, true, false)
		}
	case EResourceSet:
		for it := t.GetResources().Iterator(); it.HasNext(); {
			notifier := it.Next().(ENotifier)
			adapter.removeAdapterWithChecks(notifier, false, false)
		}
	}
}

func (adapter *EContentAdapter) selfAdapt(notification ENotification) {
	notifier := notification.GetNotifier()
	switch notifier.(type) {
	case EObject:
		feature := notification.GetFeature()
		if reference, _ := feature.(EReference); reference != nil {
			if reference.IsContainment() {
				adapter.handleContainment(notification)
			}
		}
	case EResource:
		if notification.GetFeatureID() == RESOURCE__CONTENTS {
			adapter.handleContainment(notification)
		}
	case EResourceSet:
		if notification.GetFeatureID() == RESOURCE_SET__RESOURCES {
			adapter.handleContainment(notification)
		}
	}
}

func (adapter *EContentAdapter) handleContainment(notification ENotification) {
	switch notification.GetEventType() {
	case RESOLVE:
		// We need to be careful that the proxy may be resolved while we are attaching this adapter.
		// We need to avoid attaching the adapter during the resolve
		// and also attaching it again as we walk the eContents() later.
		// Checking here avoids having to check during addAdapter.
		//
		oldNotifier := notification.GetOldValue().(ENotifier)
		if oldNotifier.EAdapters().Contains(adapter) {
			adapter.removeAdapter(oldNotifier)
			adapter.addAdapter(notification.GetNewValue().(ENotifier))
		}
	case UNSET:
		adapter.removeAdapterWithChecks(notification.GetOldValue().(ENotifier), false, true)
		adapter.addAdapter(notification.GetNewValue().(ENotifier))
	case SET:
		adapter.removeAdapterWithChecks(notification.GetOldValue().(ENotifier), false, true)
		adapter.addAdapter(notification.GetNewValue().(ENotifier))
	case ADD:
		adapter.addAdapter(notification.GetNewValue().(ENotifier))
	case ADD_MANY:
		newValues := notification.GetNewValue().([]interface{})
		for _, notifier := range newValues {
			adapter.addAdapter(notifier.(ENotifier))
		}
	case REMOVE:
		oldValue := notification.GetOldValue().(ENotifier)
		if oldValue != nil {
			_, checkContainer := notification.GetNotifier().(EResource)
			checkResource := notification.GetFeature() != nil
			adapter.removeAdapterWithChecks(oldValue, checkContainer, checkResource)
		}
	case REMOVE_MANY:
		_, checkContainer := notification.GetNotifier().(EResource)
		checkResource := notification.GetFeature() != nil
		oldValues := notification.GetOldValue().([]interface{})
		for _, notifier := range oldValues {
			adapter.removeAdapterWithChecks(notifier.(ENotifier), checkContainer, checkResource)
		}
	}
}

func (adapter *EContentAdapter) addAdapter(notifier ENotifier) {
	if notifier != nil {
		eAdapters := notifier.EAdapters()
		if !eAdapters.Contains(adapter.interfaces) {
			eAdapters.Add(adapter.interfaces)
		}
	}
}

func (adapter *EContentAdapter) removeAdapter(notifier ENotifier) {
	if notifier != nil {
		notifier.EAdapters().Remove(adapter.interfaces)
	}
}

func (adapter *EContentAdapter) removeAdapterWithChecks(notifier ENotifier, checkContainer bool, checkResource bool) {
	if notifier != nil {
		if checkContainer || checkResource {
			if internalEObject, _ := notifier.(EObjectInternal); internalEObject != nil {
				if checkResource {
					if internalResource := internalEObject.EInternalResource(); internalResource != nil && internalResource.EAdapters().Contains(adapter) {
						return
					}
				}
				if checkContainer {
					if internalContainer := internalEObject.EInternalContainer(); internalContainer != nil && internalContainer.EAdapters().Contains(adapter) {
						return
					}
				}
			}
		}
		adapter.removeAdapter(notifier)
	}
}
