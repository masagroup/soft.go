package sqlite

import "github.com/masagroup/soft.go/ecore"

type Object interface {
	ecore.EObject
	SetSQLID(int64)
	GetSQLID() int64
}
