package ecore

import (
	"zombiezen.com/go/sqlite"
)

type sqlBase struct {
	conn            *sqlite.Conn
	schema          *sqlSchema
	uri             *URI
	idAttributeName string
	idManager       EObjectIDManager
}
