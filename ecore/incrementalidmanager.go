package ecore

import (
	"fmt"
	"strconv"
	"sync"
	"unsafe"
)

type IncrementalIDManager struct {
	mutex        sync.RWMutex
	detachedToID map[uintptr]int64
	objectToID   map[EObject]int64
	idToObject   map[int64]EObject
	currentID    int64
}

func NewIncrementalIDManager() *IncrementalIDManager {
	return &IncrementalIDManager{
		detachedToID: make(map[uintptr]int64),
		objectToID:   make(map[EObject]int64),
		idToObject:   make(map[int64]EObject),
		currentID:    0,
	}
}

func (m *IncrementalIDManager) newID() int64 {
	id := m.currentID
	m.currentID++
	return id
}

func (m *IncrementalIDManager) getID(id any) (int64, error) {
	switch v := id.(type) {
	case nil:
		return -1, nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	case uint:
		return int64(v), nil
	case uint64:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	case int32:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int8:
		return int64(v), nil
	}
	return 0, fmt.Errorf("id:'%v' not supported by IncrementalIDManager", id)
}

func (m *IncrementalIDManager) Clear() {
	m.mutex.Lock()
	m.detachedToID = make(map[uintptr]int64)
	m.objectToID = make(map[EObject]int64)
	m.idToObject = make(map[int64]EObject)
	m.currentID = 0
	m.mutex.Unlock()
}

func (m *IncrementalIDManager) setID(eObject EObject, id any) error {
	if oldID, isOldID := m.objectToID[eObject]; isOldID {
		delete(m.idToObject, oldID)
	}
	newID, err := m.getID(id)
	if err == nil {
		if newID >= 0 {
			m.currentID = max(m.currentID, newID+1)
			m.objectToID[eObject] = newID
			m.idToObject[newID] = eObject
		} else {
			delete(m.objectToID, eObject)
		}
	}
	return err
}

func (m *IncrementalIDManager) Register(eObject EObject) {
	m.mutex.Lock()
	if _, isID := m.objectToID[eObject]; !isID {
		// if object is detached, retrieve its id
		// remove it from detached map
		objectHash := m.getHash(eObject)
		newID, isOldID := m.detachedToID[objectHash]
		if isOldID {
			delete(m.detachedToID, objectHash)
		} else {
			newID = m.newID()
		}
		m.setID(eObject, newID)
	}
	m.mutex.Unlock()
}

func max(x, y int64) int64 {
	if x < y {
		return y
	}
	return x
}

func (m *IncrementalIDManager) SetID(eObject EObject, id any) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.setID(eObject, id)
}

func (m *IncrementalIDManager) UnRegister(eObject EObject) {
	m.mutex.Lock()
	if id, isPresent := m.objectToID[eObject]; isPresent {
		delete(m.idToObject, id)
		delete(m.objectToID, eObject)
		// register as detached
		// add to detached map
		objectHash := m.getHash(eObject)
		m.detachedToID[objectHash] = id
	}
	m.mutex.Unlock()
}

func (m *IncrementalIDManager) GetID(eObject EObject) any {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if id, isPresent := m.objectToID[eObject]; isPresent {
		return id
	}
	return nil
}

func (m *IncrementalIDManager) GetEObject(id any) EObject {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if v, err := m.getID(id); err == nil {
		return m.idToObject[v]
	}
	return nil
}

func (m *IncrementalIDManager) GetDetachedID(eObject EObject) any {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	objectHash := m.getHash(eObject)
	if id, isDetached := m.detachedToID[objectHash]; isDetached {
		return id
	}
	return nil
}

func (m *IncrementalIDManager) getHash(eObject EObject) uintptr {
	i := (*[2]uintptr)(unsafe.Pointer(&eObject))
	return *(*uintptr)(unsafe.Pointer(&i[1]))
}
