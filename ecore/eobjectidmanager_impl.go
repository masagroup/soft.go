// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import (
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

func (m *EObjectIDManagerImpl) SetID(eObject EObject, id any) error {
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
	return fmt.Errorf("id :'%v' not supported by EObjectIDManager", id)
}

func (m *EObjectIDManagerImpl) UnRegister(eObject EObject) {
	if id, isPresent := m.objectToID[eObject]; isPresent {
		delete(m.idToObject, id)
		delete(m.objectToID, eObject)
	}
}

func (m *EObjectIDManagerImpl) GetID(eObject EObject) any {
	if id, isPresent := m.objectToID[eObject]; isPresent {
		return id
	}
	return nil
}

func (m *EObjectIDManagerImpl) GetEObject(id any) EObject {
	switch i := id.(type) {
	case string:
		return m.idToObject[i]
	default:
		return nil
	}
}

func (m *EObjectIDManagerImpl) GetDetachedID(eObject EObject) any {
	return GetEObjectID(eObject)
}
