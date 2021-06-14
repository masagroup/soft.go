package ecore

import (
	"io"
)

const (
	RESOURCE__RESOURCE_SET = 0

	RESOURCE__URI = 1

	RESOURCE__CONTENTS = 2

	RESOURCE__IS_LOADED = 4

	RESOURCE__ERRORS = 5

	RESOURCE__WARNINGS = 6
)

//EResource ...
type EResource interface {
	ENotifier

	GetResourceSet() EResourceSet

	GetURI() *URI
	SetURI(*URI)

	GetContents() EList
	GetAllContents() EIterator

	GetEObject(string) EObject
	GetURIFragment(EObject) string

	Attached(object EObject)
	Detached(object EObject)

	IsLoaded() bool

	Load()
	LoadWithOptions(options map[string]interface{})
	LoadWithReader(r io.Reader, options map[string]interface{})
	LoadWithDecoder(decoder EResourceDecoder)

	Unload()

	Save()
	SaveWithOptions(options map[string]interface{})
	SaveWithWriter(w io.Writer, options map[string]interface{})
	SaveWithEncoder(encoder EResourceEncoder)

	GetErrors() EList
	GetWarnings() EList

	SetObjectIDManager(EObjectIDManager)
	GetObjectIDManager() EObjectIDManager
}
