package ecore

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"reflect"
	"strconv"
	"sync"

	"github.com/oklog/ulid/v2"
)

type UniqueIDManager[ID comparable] struct {
	detachedToID map[uintptr]ID
	objectToID   map[EObject]ID
	idToObject   map[ID]EObject
	newID        func() ID
	isValidID    func(ID) bool
	getID        func(any) (ID, error)
	setID        func(ID)
	mutex        sync.RWMutex
}

func newUniqueIDManager[ID comparable](newID func() ID, isValidID func(ID) bool, getID func(any) (ID, error), setID func(ID)) *UniqueIDManager[ID] {
	return &UniqueIDManager[ID]{
		detachedToID: map[uintptr]ID{},
		objectToID:   map[EObject]ID{},
		idToObject:   map[ID]EObject{},
		newID:        newID,
		isValidID:    isValidID,
		getID:        getID,
		setID:        setID,
	}
}

func (m *UniqueIDManager[ID]) Clear() {
	m.mutex.Lock()
	m.detachedToID = map[uintptr]ID{}
	m.objectToID = map[EObject]ID{}
	m.idToObject = map[ID]EObject{}
	m.mutex.Unlock()
}

func (m *UniqueIDManager[ID]) setObjectID(eObject EObject, newID ID) {
	if oldID, isOldID := m.objectToID[eObject]; isOldID {
		delete(m.idToObject, oldID)
	}

	if m.isValidID(newID) {
		m.setID(newID)
		m.objectToID[eObject] = newID
		m.idToObject[newID] = eObject
	} else {
		delete(m.objectToID, eObject)
	}
}

func (m *UniqueIDManager[ID]) getPointer(eObject EObject) uintptr {
	return reflect.ValueOf(eObject).Pointer()
}

func (m *UniqueIDManager[ID]) Register(eObject EObject) {
	m.mutex.Lock()
	if _, isID := m.objectToID[eObject]; !isID {
		ptr := m.getPointer(eObject)
		newID, isOldID := m.detachedToID[ptr]
		if isOldID {
			delete(m.detachedToID, ptr)
		} else {
			newID = m.newID()
		}
		m.setObjectID(eObject, newID)
	}
	m.mutex.Unlock()
}

func (m *UniqueIDManager[ID]) UnRegister(eObject EObject) {
	m.mutex.Lock()
	if id, isID := m.objectToID[eObject]; isID {
		delete(m.idToObject, id)
		delete(m.objectToID, eObject)
		ptr := m.getPointer(eObject)
		m.detachedToID[ptr] = id
	}
	m.mutex.Unlock()
}

func (m *UniqueIDManager[ID]) SetID(eObject EObject, id any) error {
	m.mutex.Lock()
	newID, err := m.getID(id)
	if err == nil {
		m.setObjectID(eObject, newID)
	}
	m.mutex.Unlock()
	return err
}

func (m *UniqueIDManager[ID]) GetID(eObject EObject) any {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	if id, isPresent := m.objectToID[eObject]; isPresent {
		return id
	}
	return nil
}

func (m *UniqueIDManager[ID]) GetEObject(id any) EObject {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	if v, err := m.getID(id); err == nil {
		return m.idToObject[v]
	}
	return nil
}

func (m *UniqueIDManager[ID]) GetDetachedID(eObject EObject) any {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	ptr := m.getPointer(eObject)
	if id, isDetached := m.detachedToID[ptr]; isDetached {
		return id
	}
	return nil
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

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

type UUIDManager = UniqueIDManager[string]

func NewUUIDManager(size int) *UUIDManager {
	return newUniqueIDManager(
		func() string {
			for {
				id, error := generateRandomString(size)
				if error == nil {
					return id
				}
			}
		},
		func(s string) bool {
			return len(s) > 0
		},
		func(a any) (string, error) {
			if id, isString := a.(string); isString {
				return id, nil
			}
			return "", fmt.Errorf("id:'%v' not supported by UUIDManager", a)
		},
		func(s string) {
		},
	)
}

type ULIDManager = UniqueIDManager[string]

func NewULIDManager() *ULIDManager {
	return newUniqueIDManager(
		func() string {
			return ulid.Make().String()
		},
		func(s string) bool {
			return len(s) > 0
		},
		func(v any) (string, error) {
			if id, isString := v.(string); isString {
				return id, nil
			}
			return "", fmt.Errorf("id:'%v' not supported by ULIDManager", v)
		},
		func(s string) {
		},
	)
}

type IncrementalIDManager = UniqueIDManager[int64]

func NewIncrementalIDManager() *IncrementalIDManager {
	currentID := int64(0)
	return newUniqueIDManager(
		func() int64 {
			id := currentID
			currentID++
			return id
		},
		func(i int64) bool {
			return i >= 0
		},
		func(id any) (int64, error) {
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
		},
		func(id int64) {
			currentID = max(id+1, currentID)
		},
	)
}
