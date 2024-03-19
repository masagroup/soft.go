package ecore

import (
	"fmt"
	"reflect"
	"strconv"
	"sync"

	"github.com/google/uuid"
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

func NewUniqueIDManager[ID comparable](newID func() ID, isValidID func(ID) bool, getID func(any) (ID, error), setID func(ID)) *UniqueIDManager[ID] {
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

type UUIDManager = UniqueIDManager[uuid.UUID]

func NewUUIDManager() *UUIDManager {
	return NewUniqueIDManager(
		func() uuid.UUID {
			return uuid.New()
		},
		func(u uuid.UUID) bool {
			return true
		},
		func(a any) (uuid.UUID, error) {
			switch v := a.(type) {
			case string:
				return uuid.Parse(v)
			case []byte:
				return uuid.FromBytes(v)
			case uuid.UUID:
				return v, nil
			default:
				return uuid.UUID{}, fmt.Errorf("id:'%v' not supported by UUIDManager", a)
			}
		},
		func(s uuid.UUID) {
		},
	)
}

type ULIDManager = UniqueIDManager[ulid.ULID]

func NewULIDManager() *ULIDManager {
	return NewUniqueIDManager(
		func() ulid.ULID {
			return ulid.Make()
		},
		func(u ulid.ULID) bool {
			return true
		},
		func(a any) (ulid.ULID, error) {
			switch v := a.(type) {
			case string:
				return ulid.Parse(v)
			case []byte:
				u := ulid.ULID{}
				err := u.UnmarshalBinary(v)
				return u, err
			case ulid.ULID:
				return v, nil
			default:
				return ulid.ULID{}, fmt.Errorf("id:'%v' not supported by UUIDManager", a)
			}
		},
		func(s ulid.ULID) {
		},
	)
}

type IncrementalIDManager = UniqueIDManager[int64]

func NewIncrementalIDManager() *IncrementalIDManager {
	currentID := int64(0)
	return NewUniqueIDManager(
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
