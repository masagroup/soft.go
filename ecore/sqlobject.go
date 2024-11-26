package ecore

type SQLObject interface {
	EObject
	SetSQLID(int64)
	GetSQLID() int64
}
