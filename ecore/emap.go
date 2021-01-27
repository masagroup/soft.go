package ecore

type EMap interface {
	EList

	GetValue(value interface{}) interface{}

	Put(key interface{}, value interface{})

	RemoveKey(key interface{}) interface{}

	ContainsValue(value interface{}) bool

	ContainsKey(key interface{}) bool

	ToMap() map[interface{}]interface{}
}
