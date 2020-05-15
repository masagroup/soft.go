package ecore

type EResourceIDManagerImpl struct {
	objectToID map[EObject]string
	idToObject map[string]EObject
}

func NewEResourceIDManagerImpl() EResourceIDManager {
	return &EResourceIDManagerImpl{
		objectToID: make(map[EObject]string),
		idToObject: make(map[string]EObject),
	}
}

func (m *EResourceIDManagerImpl) Register(eObject EObject) {
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

func (m *EResourceIDManagerImpl) UnRegister(eObject EObject) {
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

func (m *EResourceIDManagerImpl) GetID(eObject EObject) string {
	return m.objectToID[eObject]
}

func (m *EResourceIDManagerImpl) GetEObject(id string) EObject {
	return m.idToObject[id]
}
