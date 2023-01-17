package ecore

type EResourceListener interface {
	Attached(object EObject)
	Detached(object EObject)
}
