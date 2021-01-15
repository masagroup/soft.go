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
