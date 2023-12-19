package ecore

type immutableEMap struct {
	immutableEList
	data map[any]any
}

func NewImmutableEMap(entries []any) *immutableEMap {
	m := &immutableEMap{
		immutableEList: immutableEList{
			data: entries,
		},
		data: map[any]any{},
	}
	for _, e := range entries {
		entry := e.(EMapEntry)
		m.data[entry.GetKey()] = entry.GetValue()
	}
	return m
}

func (m *immutableEMap) GetValue(key any) any {
	return m.data[key]
}

func (m *immutableEMap) Put(key any, value any) {
	panic("Immutable map can't be modified")
}

func (m *immutableEMap) RemoveKey(key any) any {
	panic("Immutable map can't be modified")
}

func (m *immutableEMap) ContainsValue(value any) bool {
	for _, v := range m.data {
		if v == value {
			return true
		}
	}
	return false
}

func (m *immutableEMap) ContainsKey(key any) bool {
	_, ok := m.data[key]
	return ok
}

func (m *immutableEMap) ToMap() map[any]any {
	return m.data
}
