package ecore

type EObjectIDManager interface {
	Clear()

	Register(EObject)
	UnRegister(EObject)

	GetID(EObject) string
	GetEObject(string) EObject
}
