package ecore

type sqlBase struct {
	codecVersion    int64
	schema          *sqlSchema
	uri             *URI
	objectIDName    string
	objectIDManager EObjectIDManager
	isContainerID   bool
}
