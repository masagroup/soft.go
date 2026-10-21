package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// a cycle of mappings stops instead of recursing forever
func TestEURIConverterImpl_NormalizeCycle(t *testing.T) {
	c := NewEURIConverterImpl()
	c.GetURIMap()[*NewURI("file:/a/")] = *NewURI("file:/b/")
	c.GetURIMap()[*NewURI("file:/b/")] = *NewURI("file:/a/")
	require.NotPanics(t, func() { c.Normalize(NewURI("file:/a/x.xml")) })
}

// the longest matching prefix maps the uri, whatever the iteration order of the map
func TestEURIConverterImpl_NormalizeLongestPrefix(t *testing.T) {
	for range 20 {
		c := NewEURIConverterImpl()
		c.GetURIMap()[*NewURI("file:/a/")] = *NewURI("file:/x/")
		c.GetURIMap()[*NewURI("file:/a/b/")] = *NewURI("file:/y/")
		assert.Equal(t, "file:/y/m.xml", c.Normalize(NewURI("file:/a/b/m.xml")).String())
	}
}

// a prefix matches whole path segments
func TestEURIConverterImpl_NormalizeSegmentBoundary(t *testing.T) {
	c := NewEURIConverterImpl()
	c.GetURIMap()[*NewURI("file:/a")] = *NewURI("file:/x")
	assert.Equal(t, "file:/ab/m.xml", c.Normalize(NewURI("file:/ab/m.xml")).String())
	assert.Equal(t, "file:/x/m.xml", c.Normalize(NewURI("file:/a/m.xml")).String())
}

// a supplied package is kept: the next lookups return the same package
func TestEPackageRegistryImpl_SupplierCached(t *testing.T) {
	r := NewEPackageRegistryImpl()
	calls := 0
	r.PutSupplier("http://p", func() EPackage {
		calls++
		return GetFactory().CreateEPackage()
	})
	first := r.GetPackage("http://p")
	assert.Equal(t, first, r.GetPackage("http://p"))
	assert.Equal(t, 1, calls)
}

// a REMOVE before the removals of a REMOVE_MANY keeps all their positions
func TestNotification_MergeRemoveBeforeRemoveMany(t *testing.T) {
	o := NewMockEObject(t)
	n := NewNotificationByFeatureID(o, REMOVE_MANY, 1, []any{"b", "c"}, []any{5, 6}, NO_INDEX)
	assert.True(t, n.Merge(NewNotificationByFeatureID(o, REMOVE, 1, "a", nil, 0)))
	assert.Equal(t, []any{"a", "b", "c"}, n.GetOldValue())
	assert.Equal(t, []any{0, 5, 6}, n.GetNewValue())
}
