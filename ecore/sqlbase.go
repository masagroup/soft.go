package ecore

import "go.uber.org/zap"

type sqlBase struct {
	codecVersion    int64
	schema          *sqlSchema
	uri             *URI
	objectIDName    string
	objectIDManager EObjectIDManager
	isObjectID      bool
	isContainerID   bool
	logger          *zap.Logger
}
