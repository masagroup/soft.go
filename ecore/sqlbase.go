package ecore

type sqlBase struct {
	codecVersion    int64
	schema          *sqlSchema
	uri             *URI
	idAttributeName string
	idManager       EObjectIDManager
}
