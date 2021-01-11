package ecore

type BasicEMap struct {
	EList
	mapData map[interface{}]interface{}
}

type mapList struct {
	basicEList
	m *BasicEMap
}

func newMapList(m *BasicEMap) *mapList {
	l := new(mapList)
	l.interfaces = l
	l.data = []interface{}{}
	l.isUnique = true
	l.m = m
	return l
}

func (ml *mapList) didAdd(index int, elem interface{}) {
	entry := elem.(EMapEntry)
	ml.m.mapData[entry.GetKey()] = entry.GetValue()
}

func (ml *mapList) didSet(index int, newElem interface{}, oldElem interface{}) {
	newEntry := newElem.(EMapEntry)
	oldEntry := oldElem.(EMapEntry)
	delete(ml.m.mapData, oldEntry.GetKey())
	ml.m.mapData[newEntry.GetKey()] = newEntry.GetValue()
}

func (ml *mapList) didRemove(index int, oldElem interface{}) {
	oldEntry := oldElem.(EMapEntry)
	delete(ml.m.mapData, oldEntry.GetKey())
}

func (ml *mapList) didClear(oldObjects []interface{}) {
	ml.m.mapData = make(map[interface{}]interface{})
}

type mapEntry struct {
	key   interface{}
	value interface{}
}

func NewMapEntry(key, value interface{}) *mapEntry {
	return &mapEntry{key: key, value: value}
}

func (entry *mapEntry) GetKey() interface{} {
	return entry.key
}

func (entry *mapEntry) GetValue() interface{} {
	return entry.value
}

func NewBasicEMap() *BasicEMap {
	basicEMap := &BasicEMap{}
	basicEMap.EList = newMapList(basicEMap)
	basicEMap.mapData = make(map[interface{}]interface{})
	return basicEMap
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
	m.Add(&mapEntry{key: key, value: value})
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
