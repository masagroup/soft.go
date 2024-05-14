package ecore

type SQLObject interface {
	EObject
	SetSqlID(int64)
	GetSqlID() int64
}
