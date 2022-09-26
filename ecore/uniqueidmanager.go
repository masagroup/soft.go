package ecore

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sync"
	"unsafe"
)

// Adapted from https://neilmadden.blog/2018/08/30/moving-away-from-uuids/

type UniqueIDManager struct {
	mutex        sync.RWMutex
	detachedToID map[uintptr]string
	objectToID   map[EObject]string
	idToObject   map[string]EObject
	size         int
}

func NewUniqueIDManager(size int) *UniqueIDManager {
	return &UniqueIDManager{
		detachedToID: make(map[uintptr]string),
		objectToID:   make(map[EObject]string),
		idToObject:   make(map[string]EObject),
		size:         size,
	}
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}
	return b, nil
}

func generateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func (m *UniqueIDManager) newID() string {
	id, error := generateRandomString(m.size)
	for error != nil {
		id, error = generateRandomString(m.size)
	}
	return id
}

func (m *UniqueIDManager) Clear() {
	m.mutex.Lock()
	m.detachedToID = make(map[uintptr]string)
	m.objectToID = make(map[EObject]string)
	m.idToObject = make(map[string]EObject)
	m.mutex.Unlock()
}

func (m *UniqueIDManager) Register(eObject EObject) {
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

func (m *UniqueIDManager) setID(eObject EObject, id any) error {
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
	return fmt.Errorf("id:'%v' not supported by UniqueIDManager", id)
}

func (m *UniqueIDManager) SetID(eObject EObject, id any) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.setID(eObject, id)
}

func (m *UniqueIDManager) UnRegister(eObject EObject) {
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

func (m *UniqueIDManager) GetID(eObject EObject) any {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	if id, isPresent := m.objectToID[eObject]; isPresent {
		return id
	}
	return nil
}

func (m *UniqueIDManager) GetEObject(id any) EObject {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	switch i := id.(type) {
	case string:
		return m.idToObject[i]
	default:
		return nil
	}
}

func (m *UniqueIDManager) GetDetachedID(eObject EObject) any {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	objectHash := m.getHash(eObject)
	if id, isDetached := m.detachedToID[objectHash]; isDetached {
		return id
	}
	return nil
}

func (m *UniqueIDManager) getHash(eObject EObject) uintptr {
	i := (*[2]uintptr)(unsafe.Pointer(&eObject))
	return *(*uintptr)(unsafe.Pointer(&i[1]))
}
