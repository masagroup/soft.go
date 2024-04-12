package ecore

type EStoreCache interface {
	// Set object with a cache for its feature values
	SetCache(bool)

	// Returns true if object is caching feature values
	IsCache() bool

	// Clear object feature values cache
	ClearCache()
}
