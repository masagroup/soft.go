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
	basicEMap.mapData = make(map[interface{}]interface{})
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

func (m *BasicEMap) GetValue(value interface{}) interface{} {
	return m.mapData[value]
}

func (m *BasicEMap) Put(key interface{}, value interface{}) {
	m.mapData[key] = value
	if e := m.getEntryForKey(key); e != nil {
		e.SetValue(value)
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

func (m *BasicEMap) doAdd(e EMapEntry) {
	m.mapData[e.GetKey()] = e.GetValue()
}

func (m *BasicEMap) doRemove(e EMapEntry) {
	delete(m.mapData, e.GetKey())
}

func (m *BasicEMap) doClear() {
	m.mapData = make(map[interface{}]interface{})
}
