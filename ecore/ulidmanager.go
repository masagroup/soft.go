package ecore

import (
	"fmt"
	"sync"
	"unsafe"

	"github.com/oklog/ulid/v2"
)

// Universally Unique Lexicographically Sortable Identifier
// https://pkg.go.dev/github.com/oklog/ulid
type ULIDManager struct {
	mutex        sync.RWMutex
	detachedToID map[uintptr]string
	objectToID   map[EObject]string
	idToObject   map[string]EObject
}

func NewULIDManager() *ULIDManager {
	return &ULIDManager{
		detachedToID: make(map[uintptr]string),
		objectToID:   make(map[EObject]string),
		idToObject:   make(map[string]EObject),
	}
}

func (m *ULIDManager) newID() string {
	return ulid.Make().String()
}

func (m *ULIDManager) Clear() {
	m.mutex.Lock()
	m.detachedToID = make(map[uintptr]string)
	m.objectToID = make(map[EObject]string)
	m.idToObject = make(map[string]EObject)
	m.mutex.Unlock()
}

func (m *ULIDManager) Register(eObject EObject) {
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

func (m *ULIDManager) setID(eObject EObject, id any) error {
	if id == nil {
		id = ""
	}
	if newID, isString := id.(string); isString {
		oldID, isOldID := m.objectToID[eObject]
		if len(newID) > 0 {
			m.objectToID[eObject] = newID
		} else {
			delete(m.objectToID, eObject)
		}

		if isOldID {
			delete(m.idToObject, oldID)
		}

		if len(newID) > 0 {
			m.idToObject[newID] = eObject
		}
		return nil
	}
	return fmt.Errorf("id:'%v' not supported by ULIDManager", id)
}

func (m *ULIDManager) SetID(eObject EObject, id any) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.setID(eObject, id)
}

func (m *ULIDManager) UnRegister(eObject EObject) {
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

func (m *ULIDManager) GetID(eObject EObject) any {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if id, isPresent := m.objectToID[eObject]; isPresent {
		return id
	}
	return nil
}

func (m *ULIDManager) GetEObject(id any) EObject {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	switch i := id.(type) {
	case string:
		return m.idToObject[i]
	default:
		return nil
	}
}

func (m *ULIDManager) GetDetachedID(eObject EObject) any {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	objectHash := m.getHash(eObject)
	if id, isDetached := m.detachedToID[objectHash]; isDetached {
		return id
	}
	return nil
}

func (m *ULIDManager) getHash(eObject EObject) uintptr {
	i := (*[2]uintptr)(unsafe.Pointer(&eObject))
	return *(*uintptr)(unsafe.Pointer(&i[1]))
}
