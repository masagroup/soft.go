package ecore

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLinkedHashMap(t *testing.T) {
	hm := newLinkedHashMap[int, int]()
	require.Equal(t, 0, hm.Len())

	_, exists := hm.Get(0)
	require.False(t, exists, "shouldn't have found the value")

	_, _, exists = hm.Oldest()
	require.False(t, exists, "shouldn't have found a value")

	_, _, exists = hm.Newest()
	require.False(t, exists, "shouldn't have found a value")

	hm.Put(0, 0)
	require.Equal(t, 1, hm.Len(), "wrong hashmap length")

	val0, exists := hm.Get(0)
	require.True(t, exists, "should have found the value")
	require.Zero(t, val0, "wrong value")

	rkey0, val0, exists := hm.Oldest()
	require.True(t, exists, "should have found the value")
	require.Equal(t, 0, rkey0, "wrong key")
	require.Zero(t, val0, "wrong value")

	rkey0, val0, exists = hm.Newest()
	require.True(t, exists, "should have found the value")
	require.Equal(t, 0, rkey0, "wrong key")
	require.Zero(t, val0, "wrong value")

	hm.Put(1, 1)
	require.Equal(t, 2, hm.Len(), "wrong hashmap length")

	val1, exists := hm.Get(1)
	require.True(t, exists, "should have found the value")
	require.Equal(t, 1, val1, "wrong value")

	rkey0, val0, exists = hm.Oldest()
	require.True(t, exists, "should have found the value")
	require.Equal(t, 0, rkey0, "wrong key")
	require.Zero(t, val0, "wrong value")

	rkey1, val1, exists := hm.Newest()
	require.True(t, exists, "should have found the value")
	require.Equal(t, 1, rkey1, "wrong key")
	require.Equal(t, 1, val1, "wrong value")

	hm.Delete(0)
	require.Equal(t, 1, hm.Len(), "wrong hashmap length")

	_, exists = hm.Get(0)
	require.False(t, exists, "shouldn't have found the value")

	rkey1, val1, exists = hm.Oldest()
	require.True(t, exists, "should have found the value")
	require.Equal(t, rkey1, 1, "wrong key")
	require.Equal(t, 1, val1, "wrong value")

	rkey1, val1, exists = hm.Newest()
	require.True(t, exists, "should have found the value")
	require.Equal(t, 1, rkey1, "wrong key")
	require.Equal(t, 1, val1, "wrong value")

	hm.Put(0, 0)
	require.Equal(t, 2, hm.Len(), "wrong hashmap length")

	hm.Put(1, 1)
	require.Equal(t, 2, hm.Len(), "wrong hashmap length")

	rkey0, val0, exists = hm.Oldest()
	require.True(t, exists, "should have found the value")
	require.Equal(t, 0, rkey0, "wrong key")
	require.Zero(t, val0, "wrong value")

	rkey1, val1, exists = hm.Newest()
	require.True(t, exists, "should have found the value")
	require.Equal(t, 1, rkey1, "wrong key")
	require.Equal(t, 1, val1, "wrong value")
}

func TestLinkedHashMap_Iterator(t *testing.T) {

	// Case: No elements
	{
		lh := newLinkedHashMap[int, int]()
		it := lh.Iterator()
		require.NotNil(t, it)
		// Should immediately be exhausted
		require.False(t, it.HasNext())
		require.False(t, it.HasNext())
		// Should be empty
		require.Zero(t, it.Key())
		require.Zero(t, it.Value())
	}

	// Case: 1 element
	{
		lh := newLinkedHashMap[int, int]()
		it := lh.Iterator()
		require.NotNil(t, it)
		lh.Put(1, 1)
		require.True(t, it.HasNext())
		require.Equal(t, 1, it.Key())
		require.Equal(t, 1, it.Value())
		// Should be empty
		require.False(t, it.HasNext())
		// Re-assign 1 --> 10
		lh.Put(1, 10)
		it = lh.Iterator() // New iterator
		require.True(t, it.HasNext())
		require.Equal(t, 1, it.Key())
		require.Equal(t, 10, it.Value())
		// Should be empty
		require.False(t, it.HasNext())
		// Delete 1
		lh.Delete(1)
		it = lh.Iterator()
		require.NotNil(t, it)
		// Should immediately be exhausted
		require.False(t, it.HasNext())
	}

	// Case: Multiple elements
	{
		lh := newLinkedHashMap[int, int]()
		lh.Put(1, 1)
		lh.Put(2, 2)
		lh.Put(3, 3)
		it := lh.Iterator()
		// Should give back all 3 elements
		require.True(t, it.HasNext())
		require.Equal(t, 1, it.Key())
		require.Equal(t, 1, it.Value())
		require.True(t, it.HasNext())
		require.Equal(t, 2, it.Key())
		require.Equal(t, 2, it.Value())
		require.True(t, it.HasNext())
		require.Equal(t, 3, it.Key())
		require.Equal(t, 3, it.Value())
		// Should be exhausted
		require.False(t, it.HasNext())
	}

	// Case: Delete element that has been iterated over
	{
		lh := newLinkedHashMap[int, int]()
		lh.Put(1, 1)
		lh.Put(2, 2)
		lh.Put(3, 3)
		it := lh.Iterator()
		require.True(t, it.HasNext())
		require.True(t, it.HasNext())
		lh.Delete(1)
		lh.Delete(2)
		require.True(t, it.HasNext())
		require.Equal(t, 3, it.Key())
		require.Equal(t, 3, it.Value())
		// Should be exhausted
		require.False(t, it.HasNext())
	}
}
