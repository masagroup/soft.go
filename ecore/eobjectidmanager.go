package ecore

type EObjectIDManager interface {
	Clear()

	Register(EObject)
	UnRegister(EObject)

	SetID(EObject, interface{})

	GetID(EObject) interface{}
	GetEObject(interface{}) EObject

	GetDetachedID(EObject) interface{}
}
