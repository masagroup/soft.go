package ecore

import (
	"errors"
	"fmt"
)

type EObjectIDManagerImpl struct {
	objectToID map[EObject]string
	idToObject map[string]EObject
}

func NewEObjectIDManagerImpl() EObjectIDManager {
	return &EObjectIDManagerImpl{
		objectToID: make(map[EObject]string),
		idToObject: make(map[string]EObject),
	}
}

func (m *EObjectIDManagerImpl) Clear() {
	m.objectToID = make(map[EObject]string)
	m.idToObject = make(map[string]EObject)
}

func (m *EObjectIDManagerImpl) Register(eObject EObject) {
	id := GetEObjectID(eObject)
	if len(id) > 0 {
		m.idToObject[id] = eObject
		m.objectToID[eObject] = id
	}
}

func (m *EObjectIDManagerImpl) SetID(eObject EObject, id interface{}) error {
	if id == nil {
		id = ""
	}
	if newID, isString := id.(string); isString {
		SetEObjectID(eObject, newID)

		oldID := m.objectToID[eObject]
		if len(newID) > 0 {
			m.objectToID[eObject] = newID
		} else {
			delete(m.objectToID, eObject)
		}

		if len(oldID) > 0 {
			delete(m.idToObject, oldID)
		}

		if len(newID) > 0 {
			m.idToObject[newID] = eObject
		}
		return nil
	}
	return errors.New(fmt.Sprintf("id :'%v' not supported by EObjectIDManager", id))
}

func (m *EObjectIDManagerImpl) UnRegister(eObject EObject) {
	if id, isPresent := m.objectToID[eObject]; isPresent {
		delete(m.idToObject, id)
		delete(m.objectToID, eObject)
	}
}

func (m *EObjectIDManagerImpl) GetID(eObject EObject) interface{} {
	if id, isPresent := m.objectToID[eObject]; isPresent {
		return id
	}
	return nil
}

func (m *EObjectIDManagerImpl) GetEObject(id interface{}) EObject {
	switch id.(type) {
	case string:
		return m.idToObject[id.(string)]
	default:
		return nil
	}
}

func (m *EObjectIDManagerImpl) GetDetachedID(eObject EObject) interface{} {
	return GetEObjectID(eObject)
}
