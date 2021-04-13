package ecore

import (
	"crypto/rand"
	"encoding/base64"
)

// Adapted from https://neilmadden.blog/2018/08/30/moving-away-from-uuids/
// Adapted from https://neilmadden.blog/2018/08/30/moving-away-from-uuids/

type UniqueIDManager struct {
	detachedToID map[EObject]string
	objectToID   map[EObject]string
	idToObject   map[string]EObject
	size         int
}

func NewUniqueIDManager(size int) *UniqueIDManager {
	return &UniqueIDManager{
		detachedToID: make(map[EObject]string),
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
	m.detachedToID = make(map[EObject]string)
	m.objectToID = make(map[EObject]string)
	m.idToObject = make(map[string]EObject)
}

func (m *UniqueIDManager) Register(eObject EObject) {
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

func (m *UniqueIDManager) SetID(eObject EObject, id interface{}) {
	if id == nil {
		id = ""
	}
	if newID, isInt := id.(string); isInt {
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
	}
}

func (m *UniqueIDManager) UnRegister(eObject EObject) {
	if id, isPresent := m.objectToID[eObject]; isPresent {
		delete(m.idToObject, id)
		delete(m.objectToID, eObject)
		m.detachedToID[eObject] = id
	}
}

func (m *UniqueIDManager) GetID(eObject EObject) interface{} {
	if id, isPresent := m.objectToID[eObject]; isPresent {
		return id
	}
	return nil
}

func (m *UniqueIDManager) GetDetachedID(eObject EObject) interface{} {
	if id, isDetached := m.detachedToID[eObject]; isDetached {
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
