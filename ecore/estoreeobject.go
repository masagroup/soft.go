package ecore

type EStoreEObject interface {
	EObject

	EStore() EStore
}
