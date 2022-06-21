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
	interfaces interface{}
	mapData    map[interface{}]interface{}
}

type eMapEntryFactory interface {
	newEntry(key interface{}, value interface{}) EMapEntry
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
	ml.m.doAdd(entry)
}

func (ml *basicEMapList) didSet(index int, newElem interface{}, oldElem interface{}) {
	newEntry := newElem.(EMapEntry)
	oldEntry := oldElem.(EMapEntry)
	ml.m.doRemove(oldEntry)
	ml.m.doAdd(newEntry)
}

func (ml *basicEMapList) didRemove(index int, oldElem interface{}) {
	oldEntry := oldElem.(EMapEntry)
	ml.m.doRemove(oldEntry)
}

func (ml *basicEMapList) didClear(oldObjects []interface{}) {
	ml.m.doClear()
}

func NewBasicEMap() *BasicEMap {
	basicEMap := &BasicEMap{}
	basicEMap.EList = newBasicEMapList(basicEMap)
	basicEMap.interfaces = basicEMap
	return basicEMap
}

func (m *BasicEMap) asEMapEntryFactory() eMapEntryFactory {
	return m.interfaces.(eMapEntryFactory)
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

func (m *BasicEMap) GetValue(key interface{}) interface{} {
	m.initDataMap()
	return m.mapData[key]
}

func (m *BasicEMap) Put(key interface{}, value interface{}) {
	if e := m.getEntryForKey(key); e != nil {
		e.SetValue(value)

		if m.mapData != nil {
			m.mapData[key] = value
		}
	} else {
		m.Add(m.asEMapEntryFactory().newEntry(key, value))
	}
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
	if e := m.getEntryForKey(key); e != nil {
		m.Remove(e)
		return e.GetValue()
	}
	return nil
}

func (m *BasicEMap) ContainsValue(value interface{}) bool {
	for it := m.Iterator(); it.HasNext(); {
		e := it.Next().(EMapEntry)
		if e.GetValue() == value {
			return true
		}
	}
	return false
}

func (m *BasicEMap) ContainsKey(key interface{}) bool {
	m.initDataMap()
	_, ok := m.mapData[key]
	return ok
}

func (m *BasicEMap) ToMap() map[interface{}]interface{} {
	m.initDataMap()
	return m.mapData
}

func (m *BasicEMap) initDataMap() {
	if m.mapData == nil {
		m.mapData = map[interface{}]interface{}{}
		for itEntry := m.Iterator(); itEntry.HasNext(); {
			entry := itEntry.Next().(EMapEntry)
			m.mapData[entry.GetKey()] = entry.GetValue()
		}
	}
}

func (m *BasicEMap) doAdd(e EMapEntry) {
	if m.mapData != nil {
		m.mapData[e.GetKey()] = e.GetValue()
	}
}

func (m *BasicEMap) doRemove(e EMapEntry) {
	delete(m.mapData, e.GetKey())
}

func (m *BasicEMap) doClear() {
	m.mapData = nil
}
