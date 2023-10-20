package ecore

import "database/sql"

type sqlBase struct {
	db              *sql.DB
	schema          *sqlSchema
	uri             *URI
	idAttributeName string
	idManager       EObjectIDManager
}
