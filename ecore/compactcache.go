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
	"maps"
	"math/bits"
	"runtime"
	"slices"
	"sync/atomic"
)

// compactSlots is the number of feature IDs whose values a compactCache holds in its slots: the
// values of the IDs from compactSlots are in its overflow map.
const compactSlots = 63

// compactSealed is the bit of compactCache.claimed sealing the cache: no value is added to it in
// place any more.
const compactSealed = uint64(1) << compactSlots

// compactCache is the cached values of a CompactEStoreEObjectImpl. It is read without lock.
//
// The value of a feature ID below compactSlots is in the slot of the ID, published by mask. A value
// is added to a cache in place, without allocation: its writer claims the slot (claimed), writes
// it, then publishes it (mask). A published slot is never written again: a value replaced or
// removed, as the values of a cleared cache, makes a new cache, a copy, which replaces this one in
// the object with a compare-and-swap. A cache is sealed before it is read for a copy, so that a
// value added to it in place is not lost by the copy: a writer finding it sealed adds its value to
// the cache replacing it.
type compactCache struct {
	claimed atomic.Uint64 // the slots claimed by a writer, and compactSealed
	mask    atomic.Uint64 // the slots published
	values  []any         // the slots, by feature ID
	// overflow holds the cached values of the feature IDs from compactSlots: without a store, the
	// cache is the only place the values of the object are kept. It is never changed: a cache
	// with other overflow values is a copy.
	overflow map[int]any
	// epoch is the configuration of the object the values were cached with: a concrete type
	// that changes how its values are read (EStoreEObjectImpl: store, cache flag) replaces the
	// cache with a new epoch, so that a value read with the previous configuration fails to be
	// cached (cacheIfAbsentAt). It is 0 for the other types.
	epoch uint64
}

// compactChange is how a copy of a cache changes the value of a feature (copyWith).
type compactChange int

const (
	compactSet         compactChange = iota // the value replaces the cached one
	compactSetIfAbsent                      // the value is cached unless a value other than nil is
	compactUnset                            // the cached value is removed
)

// epochOf returns the epoch of c, 0 for none.
func (c *compactCache) epochOf() uint64 {
	if c == nil {
		return 0
	}
	return c.epoch
}

// emptyCache returns an empty cache of epoch: nil for 0.
func emptyCache(epoch uint64) *compactCache {
	if epoch == 0 {
		return nil
	}
	return &compactCache{epoch: epoch}
}

// get returns the value of featureID in c, and whether it is cached.
func (c *compactCache) get(featureID int) (any, bool) {
	if c == nil {
		return nil, false
	}
	if featureID >= compactSlots {
		value, cached := c.overflow[featureID]
		return value, cached
	}
	if c.mask.Load()&(uint64(1)<<featureID) == 0 {
		return nil, false
	}
	return c.values[featureID], true
}

// add caches value for featureID in place: it returns the value cached for featureID and true,
// the value added, or the value another writer cached first unless overwrite. It returns false
// when the value must go to a copy of c (replace): c is nil or sealed, it has no slot for
// featureID, or a value is cached for featureID and overwrite.
func (c *compactCache) add(featureID int, value any, overwrite bool) (any, bool) {
	if c == nil || featureID >= len(c.values) {
		return nil, false
	}
	bit := uint64(1) << featureID
	for {
		claimed := c.claimed.Load()
		if claimed&compactSealed != 0 {
			return nil, false
		}
		if claimed&bit != 0 {
			if overwrite {
				return nil, false
			}
			// cached by another writer: published as soon as written
			for c.mask.Load()&bit == 0 {
				runtime.Gosched()
			}
			return c.values[featureID], true
		}
		if c.claimed.CompareAndSwap(claimed, claimed|bit) {
			c.values[featureID] = value
			c.mask.Or(bit)
			return value, true
		}
	}
}

// seal seals c, and returns once the values added to it in place are published: c holds its
// final values.
func (c *compactCache) seal() {
	if c == nil {
		return
	}
	claimed := c.claimed.Or(compactSealed) &^ compactSealed
	for c.mask.Load() != claimed {
		runtime.Gosched()
	}
}

// copyWith returns a copy of c, sealed or not published, with the value of featureID changed by
// change, with slots for the feature IDs below slots at least, and the cached value of featureID.
func (c *compactCache) copyWith(featureID int, value any, change compactChange, slots int) (*compactCache, any) {
	var mask uint64
	result := &compactCache{}
	if c != nil {
		mask = c.mask.Load()
		result.epoch, result.overflow = c.epoch, c.overflow
		slots = max(slots, len(c.values))
	}
	if change != compactUnset && featureID < compactSlots {
		slots = max(slots, featureID+1)
	}
	if slots = min(slots, compactSlots); slots > 0 {
		result.values = make([]any, slots)
	}
	for m := mask; m != 0; m &= m - 1 {
		id := bits.TrailingZeros64(m)
		result.values[id] = c.values[id]
	}
	cached := value
	if current, _ := c.get(featureID); change == compactSetIfAbsent && current != nil {
		change, cached = compactSet, current
	}
	if featureID >= compactSlots {
		result.overflow = maps.Clone(result.overflow)
		if change == compactUnset {
			delete(result.overflow, featureID)
			if len(result.overflow) == 0 {
				result.overflow = nil
			}
		} else {
			if result.overflow == nil {
				result.overflow = map[int]any{}
			}
			result.overflow[featureID] = cached
		}
	} else {
		bit := uint64(1) << featureID
		if change == compactUnset {
			if mask&bit != 0 {
				result.values[featureID] = nil
			}
			mask &^= bit
		} else {
			result.values[featureID] = cached
			mask |= bit
		}
	}
	if mask == 0 && result.overflow == nil {
		return emptyCache(result.epoch), cached
	}
	result.mask.Store(mask)
	result.claimed.Store(mask)
	return result, cached
}

// with returns a copy of c, sealed or not published, with the value of featureID.
func (c *compactCache) with(featureID int, value any) *compactCache {
	result, _ := c.copyWith(featureID, value, compactSet, 0)
	return result
}

// cachedFeatureValue is a value of a cache and its feature ID.
type cachedFeatureValue struct {
	featureID int
	value     any
}

// featureValues returns the values of c by ascending feature ID: those published when it is
// read, all of them once c is sealed.
func (c *compactCache) featureValues() []cachedFeatureValue {
	if c == nil {
		return nil
	}
	mask := c.mask.Load()
	values := make([]cachedFeatureValue, 0, bits.OnesCount64(mask)+len(c.overflow))
	for m := mask; m != 0; m &= m - 1 {
		id := bits.TrailingZeros64(m)
		values = append(values, cachedFeatureValue{id, c.values[id]})
	}
	for _, featureID := range slices.Sorted(maps.Keys(c.overflow)) {
		values = append(values, cachedFeatureValue{featureID, c.overflow[featureID]})
	}
	return values
}

// replace replaces c, the cache of the object, with a copy of c whose value of featureID is
// changed by change: c is sealed first. When featureID has no slot in c, the copy has a slot for
// each feature of the class of the object, so that the object caches all its values in place in
// one block, or twice the slots of c if more, for a feature ID beyond its class. It returns the
// value cached for featureID, and false if another writer replaced c meanwhile.
func (o *CompactEStoreEObjectImpl) replace(c *compactCache, featureID int, value any, change compactChange) (any, bool) {
	c.seal()
	slots := 0
	if change != compactUnset {
		if n := len(c.slotsOf()); featureID >= n {
			slots = max(o.EClass().GetFeatureCount(), 2*n)
		}
	}
	result, cached := c.copyWith(featureID, value, change, slots)
	return cached, o.cache.CompareAndSwap(c, result)
}

// slotsOf returns the slots of c, nil for none.
func (c *compactCache) slotsOf() []any {
	if c == nil {
		return nil
	}
	return c.values
}

// setCached caches the value of featureID.
func (o *CompactEStoreEObjectImpl) setCached(featureID int, value any) {
	for {
		c := o.cache.Load()
		if _, added := c.add(featureID, value, true); added {
			return
		}
		if _, replaced := o.replace(c, featureID, value, compactSet); replaced {
			return
		}
	}
}

// cacheIfAbsent caches the value of featureID unless one is cached, and returns the cached one.
func (o *CompactEStoreEObjectImpl) cacheIfAbsent(featureID int, value any) any {
	for {
		c := o.cache.Load()
		if current, _ := c.get(featureID); current != nil {
			return current
		}
		// a nil cached is replaced
		if current, added := c.add(featureID, value, false); added && current != nil {
			return current
		}
		if cached, replaced := o.replace(c, featureID, value, compactSetIfAbsent); replaced {
			return cached
		}
	}
}

// cacheIfAbsentAt caches the value of featureID unless one is cached, while the cache has epoch:
// it returns the cached value, and false if the cache has another epoch.
func (o *CompactEStoreEObjectImpl) cacheIfAbsentAt(epoch uint64, featureID int, value any) (any, bool) {
	for {
		c := o.cache.Load()
		if c.epochOf() != epoch {
			return nil, false
		}
		if current, _ := c.get(featureID); current != nil {
			return current, true
		}
		if current, added := c.add(featureID, value, false); added && current != nil {
			return current, true
		}
		if cached, replaced := o.replace(c, featureID, value, compactSetIfAbsent); replaced {
			return cached, true
		}
	}
}

// unsetCached removes the cached value of featureID.
func (o *CompactEStoreEObjectImpl) unsetCached(featureID int) {
	for {
		c := o.cache.Load()
		if _, cached := c.get(featureID); !cached {
			return
		}
		if _, replaced := o.replace(c, featureID, nil, compactUnset); replaced {
			return
		}
	}
}

// ClearCache clears all cached feature values from the dynamic compact cache. The cache is sealed
// first: a value added to it meanwhile is cleared too.
func (o *CompactEStoreEObjectImpl) ClearCache() {
	for {
		c := o.cache.Load()
		if c == nil {
			return
		}
		c.seal()
		if o.cache.CompareAndSwap(c, emptyCache(c.epoch)) {
			return
		}
	}
}
