package ecore

type ECacheProvider interface {
	// Set object with a cache for its feature values if set to true
	SetCache(bool)

	// Returns true if object is caching values
	IsCache() bool
}
