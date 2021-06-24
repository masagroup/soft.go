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
	case int64:
		return int(v)
	case int32:
		return int(v)
	case int16:
		return int(v)
	case int8:
		return int(v)
	case int:
		return v
	}
	return -1
}

func (m *IncrementalIDManager) Clear() {
	m.detachedToID = make(map[EObject]int)
	m.objectToID = make(map[EObject]int)
	m.idToObject = make(map[int]EObject)
	m.currentID = 0
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
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func (m *IncrementalIDManager) SetID(eObject EObject, id interface{}) {
	if oldID, isOldID := m.objectToID[eObject]; isOldID {
		delete(m.idToObject, oldID)
	}

	newID := m.getID(id)
	if newID >= 0 {
		m.currentID = max(m.currentID, newID+1)
		m.objectToID[eObject] = newID
		m.idToObject[newID] = eObject
	} else {
		delete(m.objectToID, eObject)
	}
}

func (m *IncrementalIDManager) UnRegister(eObject EObject) {
	if id, isPresent := m.objectToID[eObject]; isPresent {
		delete(m.idToObject, id)
		delete(m.objectToID, eObject)
		m.detachedToID[eObject] = id
	}
}

func (m *IncrementalIDManager) GetID(eObject EObject) interface{} {
	if id, isPresent := m.objectToID[eObject]; isPresent {
		return id
	}
	return nil
}

func (m *IncrementalIDManager) GetDetachedID(eObject EObject) interface{} {
	if id, isDetached := m.detachedToID[eObject]; isDetached {
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
