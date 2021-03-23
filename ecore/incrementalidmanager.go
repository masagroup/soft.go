package ecore

import "strconv"

type IncrementalIDManager struct {
	detachedToID map[EObject]int
	objectToID   map[EObject]int
	idToObject   map[int]EObject
	currentID    int
}

func NewIncrementalIDManager() *IncrementalIDManager {
	return &IncrementalIDManager{
		detachedToID: make(map[EObject]int),
		objectToID:   make(map[EObject]int),
		idToObject:   make(map[int]EObject),
		currentID:    0,
	}
}

func (m *IncrementalIDManager) newID() int {
	id := m.currentID
	m.currentID++
	return id
}

func (m *IncrementalIDManager) getID(id interface{}) int {
	switch v := id.(type) {
	case string:
		newID, err := strconv.Atoi(v)
		if err == nil {
			return newID
		}
	case int:
		return v
	}
	return -1
}

func (m *IncrementalIDManager) Clear() {
	m.detachedToID = make(map[EObject]int)
	m.objectToID = make(map[EObject]int)
	m.idToObject = make(map[int]EObject)
}

func (m *IncrementalIDManager) Register(eObject EObject) {
	if _, isID := m.objectToID[eObject]; !isID {
		newID, isOldID := m.detachedToID[eObject]
		if isOldID {
			delete(m.detachedToID, eObject)
		} else {
			newID = m.newID()
		}
		m.SetID(eObject, newID)
	}
	// register children
	eChildren := eObject.EContents().(EObjectList).GetUnResolvedList()
	for it := eChildren.Iterator(); it.HasNext(); {
		eChild := it.Next().(EObject)
		m.Register(eChild)
	}
}

func (m *IncrementalIDManager) SetID(eObject EObject, id interface{}) {
	newID := m.getID(id)
	oldID, isOldID := m.objectToID[eObject]
	if newID >= 0 {
		m.objectToID[eObject] = newID
	} else {
		delete(m.objectToID, eObject)
	}

	if isOldID {
		delete(m.idToObject, oldID)
	}

	if newID >= 0 {
		m.idToObject[newID] = eObject
	}
}

func (m *IncrementalIDManager) UnRegister(eObject EObject) {
	if id, isPresent := m.objectToID[eObject]; isPresent {
		delete(m.idToObject, id)
		delete(m.objectToID, eObject)
		m.detachedToID[eObject] = id
	}
	// unregister children
	eChildren := eObject.EContents().(EObjectList).GetUnResolvedList()
	for it := eChildren.Iterator(); it.HasNext(); {
		eChild := it.Next().(EObject)
		m.UnRegister(eChild)
	}
}

func (m *IncrementalIDManager) GetID(eObject EObject) interface{} {
	if id, isPresent := m.objectToID[eObject]; isPresent {
		return id
	}
	return nil
}

func (m *IncrementalIDManager) GetEObject(id interface{}) EObject {
	if v := m.getID(id); v >= 0 {
		return m.idToObject[v]
	}
	return nil
}
