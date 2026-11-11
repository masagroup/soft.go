// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// compactTestFeatures are feature IDs in the first slots, in the slots of a grown cache, and in
// the overflow map.
var compactTestFeatures = []int{0, 3, 7, 8, 15, 16, 31, 40, 62, 63, 64, 70}

// TestCompactCache_ConcurrentCacheIfAbsent caches each feature from several goroutines at once:
// one value wins per feature, every goroutine gets it, and it stays cached.
func TestCompactCache_ConcurrentCacheIfAbsent(t *testing.T) {
	for round := range 200 {
		o := newTestCompactEStoreEObject(nil)
		const writers = 8
		got := make([][]any, writers)
		var group sync.WaitGroup
		for w := range writers {
			group.Go(func() {
				for _, featureID := range compactTestFeatures {
					got[w] = append(got[w], o.cacheIfAbsent(featureID, fmt.Sprintf("%d-%d", featureID, w)))
				}
			})
		}
		group.Wait()
		for i, featureID := range compactTestFeatures {
			cached := o.cachedValue(featureID)
			require.NotNil(t, cached, "round %d feature %d", round, featureID)
			for w := range writers {
				assert.Equal(t, cached, got[w][i], "round %d feature %d writer %d", round, featureID, w)
			}
		}
	}
}

// TestCompactCache_ConcurrentSetCached sets distinct features from several goroutines at once,
// some of them replaced and unset meanwhile: no value set is lost by a copy of the cache, the
// cache growing or a writer replacing it.
func TestCompactCache_ConcurrentSetCached(t *testing.T) {
	for round := range 200 {
		o := newTestCompactEStoreEObject(nil)
		var group sync.WaitGroup
		for w := range 4 {
			group.Go(func() {
				for i, featureID := range compactTestFeatures {
					if i%4 != w {
						continue
					}
					o.setCached(featureID, "first")
					o.unsetCached(featureID)
					o.setCached(featureID, fmt.Sprintf("value %d", featureID))
				}
			})
		}
		group.Wait()
		for _, featureID := range compactTestFeatures {
			assert.Equal(t, fmt.Sprintf("value %d", featureID), o.cachedValue(featureID), "round %d feature %d", round, featureID)
		}
	}
}

// TestCompactCache_ClearWhileCaching clears the cache while values are cached: the cache is empty
// after the last clear, and a value cached after it stays.
func TestCompactCache_ClearWhileCaching(t *testing.T) {
	o := newTestCompactEStoreEObject(nil)
	var group sync.WaitGroup
	group.Go(func() {
		for i := range 2000 {
			o.cacheIfAbsent(compactTestFeatures[i%len(compactTestFeatures)], i)
		}
	})
	for range 200 {
		o.ClearCache()
	}
	group.Wait()
	o.ClearCache()
	assert.Empty(t, o.cachedValues())
	o.setCached(5, "after")
	assert.Equal(t, "after", o.cachedValue(5))
}

// TestCompactCache_InPlace caches the values of an object in one block with a slot for each
// feature of its class, in place: the cache of the object stays the same. A feature ID beyond the
// class grows it.
func TestCompactCache_InPlace(t *testing.T) {
	const features = 12
	class := GetFactory().CreateEClass()
	for i := range features {
		attribute := GetFactory().CreateEAttribute()
		attribute.SetName(fmt.Sprintf("a%d", i))
		class.GetEStructuralFeatures().Add(attribute)
	}
	o := newTestCompactEStoreEObject(nil)
	o.SetEClass(class)

	o.cacheIfAbsent(0, "a")
	c := o.cache.Load()
	require.Len(t, c.values, features, "a slot for each feature of the class")
	for featureID := 1; featureID < features; featureID++ {
		o.cacheIfAbsent(featureID, featureID)
	}
	assert.Same(t, c, o.cache.Load(), "values added in place")
	o.cacheIfAbsent(features, "beyond")
	assert.Len(t, o.cache.Load().values, 2*features)
	o.setCached(0, "b")
	assert.Equal(t, "b", o.cachedValue(0))
	for featureID := 1; featureID < features; featureID++ {
		assert.Equal(t, featureID, o.cachedValue(featureID))
	}
	assert.Equal(t, "beyond", o.cachedValue(features))
}
