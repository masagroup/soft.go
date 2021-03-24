// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

type BasicEMap struct {
	EList
	mapData map[interface{}]interface{}
}

type basicEMapList struct {
	basicEList
	m *BasicEMap
}

func newBasicEMapList(m *BasicEMap) *basicEMapList {
	l := new(basicEMapList)
	l.interfaces = l
	l.data = []interface{}{}
	l.isUnique = true
	l.m = m
	return l
}

func (ml *basicEMapList) didAdd(index int, elem interface{}) {
	entry := elem.(EMapEntry)
	ml.m.mapData[entry.GetKey()] = entry.GetValue()
}

func (ml *basicEMapList) didSet(index int, newElem interface{}, oldElem interface{}) {
	newEntry := newElem.(EMapEntry)
	oldEntry := oldElem.(EMapEntry)
	delete(ml.m.mapData, oldEntry.GetKey())
	ml.m.mapData[newEntry.GetKey()] = newEntry.GetValue()
}

func (ml *basicEMapList) didRemove(index int, oldElem interface{}) {
	oldEntry := oldElem.(EMapEntry)
	delete(ml.m.mapData, oldEntry.GetKey())
}

func (ml *basicEMapList) didClear(oldObjects []interface{}) {
	ml.m.mapData = make(map[interface{}]interface{})
}

func NewBasicEMap() *BasicEMap {
	basicEMap := &BasicEMap{}
	basicEMap.Initialize()
	return basicEMap
}

func (m *BasicEMap) Initialize() {
	m.EList = newBasicEMapList(m)
	m.mapData = make(map[interface{}]interface{})
}

func (m *BasicEMap) getEntryForKey(key interface{}) EMapEntry {
	for it := m.Iterator(); it.HasNext(); {
		e := it.Next().(EMapEntry)
		if e.GetKey() == key {
			return e
		}
	}
	return nil
}

func (m *BasicEMap) GetValue(value interface{}) interface{} {
	return m.mapData[value]
}

func (m *BasicEMap) Put(key interface{}, value interface{}) {
	m.mapData[key] = value
	m.Add(m.newEntry(key, value))
}

type eMapEntryImpl struct {
	key   interface{}
	value interface{}
}

func (e *eMapEntryImpl) GetKey() interface{} {
	return e.key
}

func (e *eMapEntryImpl) SetKey(key interface{}) {
	e.key = key
}

func (e *eMapEntryImpl) GetValue() interface{} {
	return e.value
}

func (e *eMapEntryImpl) SetValue(value interface{}) {
	e.value = value
}

func (m *BasicEMap) newEntry(key interface{}, value interface{}) EMapEntry {
	return &eMapEntryImpl{key: key, value: value}
}

func (m *BasicEMap) RemoveKey(key interface{}) interface{} {
	// remove from map data
	delete(m.mapData, key)

	// remove from list
	if e := m.getEntryForKey(key); e != nil {
		m.Remove(e)
		return e.GetValue()
	}
	return nil
}

func (m *BasicEMap) ContainsValue(value interface{}) bool {
	for _, v := range m.mapData {
		if v == value {
			return true
		}
	}
	return false
}

func (m *BasicEMap) ContainsKey(key interface{}) bool {
	_, ok := m.mapData[key]
	return ok
}

func (m *BasicEMap) ToMap() map[interface{}]interface{} {
	return m.mapData
}
