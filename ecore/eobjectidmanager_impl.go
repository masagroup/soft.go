package ecore

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
	eChildren := eObject.EContents().(EObjectList).GetUnResolvedList()
	for it := eChildren.Iterator(); it.HasNext(); {
		eChild := it.Next().(EObject)
		m.Register(eChild)
	}
}

func (m *EObjectIDManagerImpl) SetID(eObject EObject, id interface{}) {
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
	}
}

func (m *EObjectIDManagerImpl) UnRegister(eObject EObject) {
	if id, isPresent := m.objectToID[eObject]; isPresent {
		delete(m.idToObject, id)
		delete(m.objectToID, eObject)
	}
	eChildren := eObject.EContents().(EObjectList).GetUnResolvedList()
	for it := eChildren.Iterator(); it.HasNext(); {
		eChild := it.Next().(EObject)
		m.UnRegister(eChild)
	}
}

func (m *EObjectIDManagerImpl) GetID(eObject EObject) interface{} {
	return m.objectToID[eObject]
}

func (m *EObjectIDManagerImpl) GetEObject(id interface{}) EObject {
	switch id.(type) {
	case string:
		return m.idToObject[id.(string)]
	default:
		return nil
	}
}
