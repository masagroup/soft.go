package ecore

type EMapEntry interface {
	GetKey() interface{}

	SetKey(interface{})

	GetValue() interface{}

	SetValue(interface{})
}
