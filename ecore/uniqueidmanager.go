package ecore

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"runtime"
	"unsafe"
)

// Adapted from https://neilmadden.blog/2018/08/30/moving-away-from-uuids/

type UniqueIDManager struct {
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
	m.detachedToID = make(map[uintptr]string)
	m.objectToID = make(map[EObject]string)
	m.idToObject = make(map[string]EObject)
}

func (m *UniqueIDManager) Register(eObject EObject) {
	if _, isID := m.objectToID[eObject]; !isID {
		// if object is detached, retrieve its id
		// remove it from detached map
		// remove its finalizer
		objectHash := m.getHash(eObject)
		newID, isOldID := m.detachedToID[objectHash]
		if isOldID {
			delete(m.detachedToID, objectHash)
			runtime.SetFinalizer(eObject, nil)
		} else {
			newID = m.newID()
		}
		m.SetID(eObject, newID)
	}
}

func (m *UniqueIDManager) SetID(eObject EObject, id interface{}) error {
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
	return errors.New(fmt.Sprintf("id:'%v' not supported by UniqueIDManager", id))
}

func (m *UniqueIDManager) UnRegister(eObject EObject) {
	if id, isPresent := m.objectToID[eObject]; isPresent {
		delete(m.idToObject, id)
		delete(m.objectToID, eObject)
		// register as detached
		// add to detached map
		// register a finalizer
		objectHash := m.getHash(eObject)
		m.detachedToID[objectHash] = id
		runtime.SetFinalizer(eObject, func(_ interface{}) {
			delete(m.detachedToID, objectHash)
		})
	}
}

func (m *UniqueIDManager) GetID(eObject EObject) interface{} {
	if id, isPresent := m.objectToID[eObject]; isPresent {
		return id
	}
	return nil
}

func (m *UniqueIDManager) GetEObject(id interface{}) EObject {
	switch id.(type) {
	case string:
		return m.idToObject[id.(string)]
	default:
		return nil
	}
}

func (m *UniqueIDManager) getHash(eObject EObject) uintptr {
	i := (*[2]uintptr)(unsafe.Pointer(&eObject))
	return *(*uintptr)(unsafe.Pointer(&i[1]))
}
