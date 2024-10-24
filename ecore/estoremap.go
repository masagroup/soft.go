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

func newEStoreEMapList(m *EStoreMap, owner EObject, feature EStructuralFeature, store EStore) *eStoreEMapList {
	l := &eStoreEMapList{m: m}
	l.Initialize(owner, feature, store)
	l.SetInterfaces(l)
	return l
}

type EStoreMap struct {
	BasicEObjectMap
	store EStore
}

func NewEStoreMap(entryClass EClass, owner EObject, feature EStructuralFeature, store EStore) *EStoreMap {
	m := &EStoreMap{store: store}
	m.EList = newEStoreEMapList(m, owner, feature, store)
	m.entryClass = entryClass
	m.interfaces = m
	return m
}

func (m *EStoreMap) GetEStore() EStore {
	return m.store
}

func (m *EStoreMap) SetEStore(store EStore) {
	m.store = store
	sp := m.EList.(EStoreProvider)
	sp.SetEStore(store)
}

// Set object with a cache for its feature values
func (m *EStoreMap) SetCache(cache bool) {
	// reset map data
	m.mapData = nil
	// set internal list cache
	sc := m.EList.(ECacheProvider)
	sc.SetCache(cache)
}

// Returns true if object is caching feature values
func (m *EStoreMap) IsCache() bool {
	sc := m.EList.(ECacheProvider)
	return sc.IsCache()
}
