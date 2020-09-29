package ecore

type EResourceIDManager interface {
	Clear()

	Register(EObject)
	UnRegister(EObject)

	GetID(EObject) string
	GetEObject(string) EObject
}
