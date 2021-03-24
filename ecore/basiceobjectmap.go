// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

type BasicEObjectMap struct {
	BasicEMap
	entryClass EClass
}

func NewBasicEObjectMap(entryClass EClass) *BasicEObjectMap {
	basicEObjectMap := &BasicEObjectMap{}
	basicEObjectMap.Initialize()
	basicEObjectMap.entryClass = entryClass
	return basicEObjectMap
}

func (m *BasicEObjectMap) Put(key interface{}, value interface{}) {
	m.mapData[key] = value
	m.Add(m.newEntry(key, value))
}

func (m *BasicEObjectMap) newEntry(key interface{}, value interface{}) EMapEntry {
	eFactory := m.entryClass.GetEPackage().GetEFactoryInstance()
	eEntry := eFactory.Create(m.entryClass).(EMapEntry)
	eEntry.SetKey(key)
	eEntry.SetValue(value)
	return eEntry
}
