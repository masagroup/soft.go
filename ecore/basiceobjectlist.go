// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

type basicEObjectList struct {
	BasicENotifyingList
	owner            EObjectInternal
	featureID        int
	inverseFeatureID int
	containment      bool
	inverse          bool
	opposite         bool
	proxies          bool
	unset            bool
}

func NewBasicEObjectList(owner EObjectInternal, featureID int, inverseFeatureID int, containment, inverse, opposite, proxies, unset bool) *basicEObjectList {
	l := new(basicEObjectList)
	l.interfaces = l
	l.data = []any{}
	l.isUnique = true
	l.owner = owner
	l.featureID = featureID
	l.inverseFeatureID = inverseFeatureID
	l.containment = containment
	l.inverse = inverse
	l.opposite = opposite
	l.proxies = proxies
	l.unset = unset
	return l
}

// GetNotifier ...
func (list *basicEObjectList) GetNotifier() ENotifier {
	return list.owner
}

// GetFeature ...
func (list *basicEObjectList) GetFeature() EStructuralFeature {
	if list.owner != nil {
		return list.owner.EClass().GetEStructuralFeature(list.featureID)
	}
	return nil
}

// GetFeatureID ...
func (list *basicEObjectList) GetFeatureID() int {
	return list.featureID
}

// GetUnResolvedList ...
func (list *basicEObjectList) GetUnResolvedList() EList {
	if list.proxies {
		u := &unResolvedBasicEObjectList{}
		u.delegate = list
		return u
	}
	return list
}

func (list *basicEObjectList) IndexOf(elem any) int {
	for i, value := range list.data {
		if value == elem || (list.proxies && list.resolve(i, value) == elem) {
			return i
		}
	}
	return -1
}

func (list *basicEObjectList) doGet(index int) any {
	return list.resolve(index, list.BasicEList.doGet(index))
}

func (list *basicEObjectList) ToArray() []any {
	if list.proxies {
		for i := len(list.data) - 1; i >= 0; i-- {
			list.doGet(i)
		}
	}
	return list.data
}

func (list *basicEObjectList) resolve(index int, object any) any {
	if eObject, _ := object.(EObject); eObject != nil {
		resolved := list.resolveProxy(eObject)
		if resolved != object {
			list.BasicEList.doSet(index, resolved)
			var notifications ENotificationChain
			if list.containment {
				notifications = list.interfaces.(abstractENotifyingList).inverseRemove(object, notifications)
				if resolvedInternal, _ := resolved.(EObjectInternal); resolvedInternal != nil && resolvedInternal.EInternalContainer() == nil {
					notifications = list.interfaces.(abstractENotifyingList).inverseAdd(resolved, notifications)
				}
			}
			list.createAndDispatchNotification(notifications, RESOLVE, object, resolved, index)
		}
		return resolved
	}
	return object
}

func (list *basicEObjectList) resolveProxy(eObject EObject) EObject {
	if list.proxies && eObject.EIsProxy() {
		return list.owner.EResolveProxy(eObject)
	}
	return eObject
}

func (list *basicEObjectList) inverseAdd(object any, notifications ENotificationChain) ENotificationChain {
	internal, _ := object.(EObjectInternal)
	if internal != nil && list.inverse {
		if list.opposite {
			return internal.EInverseAdd(list.owner, list.inverseFeatureID, notifications)
		} else {
			return internal.EInverseAdd(list.owner, EOPPOSITE_FEATURE_BASE-list.featureID, notifications)
		}
	}
	return notifications
}

func (list *basicEObjectList) inverseRemove(object any, notifications ENotificationChain) ENotificationChain {
	internal, _ := object.(EObjectInternal)
	if internal != nil && list.inverse {
		if list.opposite {
			return internal.EInverseRemove(list.owner, list.inverseFeatureID, notifications)
		} else {
			return internal.EInverseRemove(list.owner, EOPPOSITE_FEATURE_BASE-list.featureID, notifications)
		}
	}
	return notifications
}

type unResolvedBasicEObjectList struct {
	AbstractDelegatingENotifyingList[*basicEObjectList]
}

func (list *unResolvedBasicEObjectList) doGet(index int) any {
	return list.delegate.data[index]
}

func (list *unResolvedBasicEObjectList) ToArray() []any {
	return list.delegate.data
}

func (l *unResolvedBasicEObjectList) GetUnResolvedList() EList {
	return l
}
