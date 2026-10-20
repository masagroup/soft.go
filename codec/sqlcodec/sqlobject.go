package sqlcodec

import "github.com/masagroup/soft.go/ecore"

type SQLObject interface {
	ecore.EObject
	SetSQLID(int64)
	GetSQLID() int64
}
