package ecore

// EObjectList is a list of EObject
type EObjectList interface {
	EList

	GetUnResolvedList() EList
}
