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

// EResource ...
type EResource interface {
	ENotifier
	EResourceListener

	GetResourceSet() EResourceSet

	GetURI() *URI
	SetURI(*URI)

	GetContents() EList
	GetAllContents() EIterator

	GetEObject(string) EObject
	GetURIFragment(EObject) string

	IsLoading() bool
	IsLoaded() bool

	Load()
	LoadWithOptions(options map[string]any)
	LoadWithReader(r io.Reader, options map[string]any)

	Unload()

	Save()
	SaveWithOptions(options map[string]any)
	SaveWithWriter(w io.Writer, options map[string]any)

	GetErrors() EList
	GetWarnings() EList

	SetObjectIDManager(EObjectIDManager)
	GetObjectIDManager() EObjectIDManager

	GetResourceListeners() EList
}
