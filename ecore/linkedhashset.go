package ecore

type linkedHashSet[K comparable] struct {
	m *linkedHashMap[K, struct{}]
}

func newLinkedHashSet[K comparable]() *linkedHashSet[K] {
	return &linkedHashSet[K]{
		m: newLinkedHashMap[K, struct{}](),
	}
}

func (h *linkedHashSet[K]) Add(k K) {
	h.m.Put(k, struct{}{})
}

func (h *linkedHashSet[K]) AddAll(collection Collection) {
	for it := collection.Iterator(); it.HasNext(); {
		h.Add(it.Next().(K))
	}
}

func (h *linkedHashSet[K]) Remove(k K) {
	h.m.Delete(k)
}

func (h *linkedHashSet[K]) RemoveAll(collection Collection) {
	if h.m.Len() >= collection.Size() {
		for it := collection.Iterator(); it.HasNext(); {
			h.m.Delete(it.Next().(K))
		}
	} else {
		for it := h.m.Iterator(); it.HasNext(); {
			k := it.Key()
			if collection.Contains(k) {
				h.m.Delete(k)
			}
		}
	}
}

func (h *linkedHashSet[K]) ToArray() []K {
	a := make([]K, h.m.Len())
	i := 0
	for it := h.m.Iterator(); it.HasNext(); {
		a[i] = it.Key()
		i++
	}
	return a
}
