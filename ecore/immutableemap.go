package ecore

type immutableEMap struct {
	immutableEList
	data map[any]any
}

func NewImmutableEMap(data map[any]any, factory func(k, v any) EMapEntry) *immutableEMap {
	m := &immutableEMap{
		data: data,
	}
	if factory != nil {
		entries := []any{}
		for k, v := range data {
			entries = append(entries, factory(k, v))
		}
		m.immutableEList.data = entries
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
