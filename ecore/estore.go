package ecore

type EStore interface {
	Get(object EObject, feature EStructuralFeature, index int) interface{}

	Set(object EObject, feature EStructuralFeature, index int, value interface{})

	IsSet(object EObject, feature EStructuralFeature) bool

	UnSet(object EObject, feature EStructuralFeature)

	IsEmpty(object EObject, feature EStructuralFeature) bool

	Size(object EObject, feature EStructuralFeature) int

	Contains(object EObject, feature EStructuralFeature, value interface{}) bool

	IndexOf(object EObject, feature EStructuralFeature, value interface{}) int

	LastIndexOf(object EObject, feature EStructuralFeature, value interface{}) int

	Add(object EObject, feature EStructuralFeature, index int, value interface{})

	Remove(object EObject, feature EStructuralFeature, index int)

	Move(object EObject, feature EStructuralFeature, targetIndex int, sourceIndex int)

	Clear(object EObject, feature EStructuralFeature)

	ToArray(object EObject, feature EStructuralFeature) []interface{}

	GetContainer(object EObject) EObject

	GetContainingFeature(object EObject) EStructuralFeature

	Create(eClass EClass) EObject
}
