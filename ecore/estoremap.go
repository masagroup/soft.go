// *****************************************************************************
// Copyright(c) 2023 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

type eStoreEMapList struct {
	EStoreList
	m *EStoreMap
}

func (ml *eStoreEMapList) DidAdd(index int, elem any) {
	entry := elem.(EMapEntry)
	ml.m.doAdd(entry)
}

func (ml *eStoreEMapList) DidSet(index int, newElem any, oldElem any) {
	newEntry := newElem.(EMapEntry)
	oldEntry := oldElem.(EMapEntry)
	ml.m.doRemove(oldEntry)
	ml.m.doAdd(newEntry)
}

func (ml *eStoreEMapList) DidRemove(index int, oldElem any) {
	oldEntry := oldElem.(EMapEntry)
	ml.m.doRemove(oldEntry)
}

func (ml *eStoreEMapList) DidClear(oldObjects []any) {
	ml.m.doClear()
}

func newEStoreEMapList(m *EStoreMap, owner EObjectInternal, feature EStructuralFeature, store EStore) *eStoreEMapList {
	l := &eStoreEMapList{m: m}
	l.Initialize(owner, feature, store)
	l.interfaces = l
	return l
}

type EStoreMap struct {
	BasicEObjectMap
}

func NewEStoreMap(entryClass EClass, owner EObjectInternal, feature EStructuralFeature, store EStore) *EStoreMap {
	m := &EStoreMap{}
	m.EList = newEStoreEMapList(m, owner, feature, store)
	m.entryClass = entryClass
	m.interfaces = m
	return m
}
