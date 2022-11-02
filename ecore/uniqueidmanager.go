package ecore

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strconv"
	"sync"

	"github.com/oklog/ulid/v2"
)

type UniqueIDManager[ID comparable] struct {
	mutex        sync.RWMutex
	detachedToID map[EObject]ID
	objectToID   map[EObject]ID
	idToObject   map[ID]EObject
	newID        func() ID
	isValidID    func(ID) bool
	getID        func(any) (ID, error)
	setID        func(ID)
}

func NewUniqueIDManager[ID comparable](newID func() ID, isValidID func(ID) bool, getID func(any) (ID, error), setID func(ID)) *UniqueIDManager[ID] {
	return &UniqueIDManager[ID]{
		objectToID: make(map[EObject]ID),
		idToObject: make(map[ID]EObject),
		newID:      newID,
		isValidID:  isValidID,
		getID:      getID,
		setID:      setID,
	}
}

func (m *UniqueIDManager[ID]) Clear() {
	m.mutex.Lock()
	m.objectToID = make(map[EObject]ID)
	m.idToObject = make(map[ID]EObject)
	m.detachedToID = nil
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

func (m *UniqueIDManager[ID]) getDetachedID(eObject EObject) (id ID, isID bool) {
	if m.detachedToID == nil {
		var emptyID ID
		return emptyID, false
	}
	id, isID = m.detachedToID[eObject]
	return
}

func (m *UniqueIDManager[ID]) Register(eObject EObject) {
	m.mutex.Lock()
	if _, isID := m.objectToID[eObject]; !isID {
		newID, isOldID := m.getDetachedID(eObject)
		if isOldID {
			delete(m.detachedToID, eObject)
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
		if m.detachedToID != nil {
			m.detachedToID[eObject] = id
		}
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
	if m.detachedToID != nil {
		if id, isDetached := m.detachedToID[eObject]; isDetached {
			return id
		}
	}
	return nil
}

func (m *UniqueIDManager[ID]) KeepIDs(keepIDs bool) bool {
	detachedToID := m.detachedToID
	if keepIDs {
		if m.detachedToID == nil {
			m.detachedToID = make(map[EObject]ID)
		}
	} else {
		m.detachedToID = nil
	}
	return detachedToID != nil
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

func NewUUIDManager(size int) *UniqueIDManager[string] {
	return NewUniqueIDManager[string](
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

func NewULIDManager() *UniqueIDManager[string] {
	return NewUniqueIDManager[string](
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

func NewIncrementalIDManager() *UniqueIDManager[int64] {
	currentID := int64(0)
	return NewUniqueIDManager[int64](
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
