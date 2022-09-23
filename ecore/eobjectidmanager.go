package ecore

type EObjectIDManager interface {
	Clear()

	Register(EObject)
	UnRegister(EObject)

	SetID(EObject, any) error

	GetID(EObject) any
	GetEObject(any) EObject

	GetDetachedID(EObject) any
}
