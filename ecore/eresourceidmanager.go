package ecore

type EResourceIDManager interface {
	Register(EObject)
	UnRegister(EObject)

	GetID(EObject) string
	GetEObject(string) EObject
}
