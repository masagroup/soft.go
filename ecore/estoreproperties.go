// *****************************************************************************
// Copyright(c) 2026 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import (
	"slices"
	"sync/atomic"
)

// eStoreProperties is the values of the dynamic features of a store object, by feature ID. The
// slice is never changed but replaced with a compare-and-swap: values are read without lock, and
// a value read from the store is kept without lock.
type eStoreProperties struct {
	values atomic.Pointer[[]any]
}

// get returns the value of featureID, nil if none.
func (p *eStoreProperties) get(featureID int) any {
	if values := p.values.Load(); values != nil && featureID < len(*values) {
		return (*values)[featureID]
	}
	return nil
}

// all returns the values: the slice must not be changed.
func (p *eStoreProperties) all() []any {
	if values := p.values.Load(); values != nil {
		return *values
	}
	return nil
}

// withProperty returns a copy of values with the value of featureID of class.
func withProperty(values *[]any, featureID int, value any, class EClass) *[]any {
	var result []any
	if values != nil {
		result = slices.Clone(*values)
	} else {
		result = make([]any, class.GetFeatureCount())
	}
	result[featureID] = value
	return &result
}

// set sets the value of featureID.
func (p *eStoreProperties) set(featureID int, value any, class EClass) {
	for {
		values := p.values.Load()
		if values == nil && value == nil {
			return
		}
		if p.values.CompareAndSwap(values, withProperty(values, featureID, value, class)) {
			return
		}
	}
}

// setIfAbsent sets the value of featureID unless one is set, and returns the one set.
func (p *eStoreProperties) setIfAbsent(featureID int, value any, class EClass) any {
	for {
		values := p.values.Load()
		if values != nil && (*values)[featureID] != nil {
			return (*values)[featureID]
		}
		if p.values.CompareAndSwap(values, withProperty(values, featureID, value, class)) {
			return value
		}
	}
}

// unset removes the value of featureID.
func (p *eStoreProperties) unset(featureID int) {
	for {
		values := p.values.Load()
		if values == nil || (*values)[featureID] == nil {
			return
		}
		if p.values.CompareAndSwap(values, withProperty(values, featureID, nil, nil)) {
			return
		}
	}
}

// lockStoreLists locks the store lists of values, and returns them to unlock.
func lockStoreLists(values []any) []*EStoreList {
	var lists []*EStoreList
	for _, v := range values {
		if storeList, isStoreList := v.(*EStoreList); isStoreList {
			storeList.Lock()
			lists = append(lists, storeList)
		}
	}
	return lists
}

func unlockStoreLists(lists []*EStoreList) {
	for _, storeList := range lists {
		storeList.Unlock()
	}
}
