// *****************************************************************************
// Copyright(c) 2026 MASA Group
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
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newLockFreeTestList(t *testing.T) *EStoreList {
	attribute := GetFactory().CreateEAttribute()
	attribute.SetName("values")
	attribute.SetEType(GetPackage().GetEString())
	attribute.SetUpperBound(-1)
	attribute.SetUnique(false)
	owner := newTestCompactEStoreEObject(nil)
	list := NewEStoreList(owner, attribute, nil)
	list.SetCache(true)
	require.True(t, list.IsCache())
	return list
}

// runWithin fails the test if f does not return within a second: it would wait for itself.
func runWithin(t *testing.T, f func()) {
	done := make(chan struct{})
	go func() {
		defer close(done)
		f()
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		require.Fail(t, "deadlock")
	}
}

// Readers read a list without lock while it is changed: each iteration sees a consistent state,
// whose elements are distinct, and never fewer elements than a state seen before.
func TestEStoreList_ConcurrentReadsDuringWrites(t *testing.T) {
	list := newLockFreeTestList(t)
	const count = 2000
	var done atomic.Bool
	var group sync.WaitGroup
	for range 4 {
		group.Go(func() {
			seen := 0
			for !done.Load() {
				values := map[any]bool{}
				for v := range list.All() {
					if values[v] {
						t.Errorf("%v seen twice", v)
						return
					}
					values[v] = true
				}
				if len(values) < seen {
					t.Errorf("%d elements, %d before", len(values), seen)
					return
				}
				seen = len(values)
				if size := list.Size(); size > 0 {
					if v := list.Get(size - 1); v == nil {
						t.Error("nil element")
						return
					}
					if index := list.IndexOf("v0"); index < 0 {
						t.Error("v0 not found")
						return
					}
				}
			}
		})
	}
	for i := range count {
		list.Add(fmt.Sprintf("v%d", i))
		switch i % 100 {
		case 50:
			list.Move(0, i)
			list.Move(i, 0)
		case 75:
			list.Set(i, list.Get(i))
		}
	}
	done.Store(true)
	group.Wait()
	require.Equal(t, count, list.Size())
	for i, v := range list.ToArray() {
		require.Equal(t, fmt.Sprintf("v%d", i), v)
	}
}

// An element is appended past the length of the published cache, without copying it; the other
// changes copy it, so that an append never writes an element a reader of an older state reads.
func TestEStoreList_AppendWithoutCopy(t *testing.T) {
	list := newLockFreeTestList(t)
	for i := range 100 {
		list.Add(i)
	}
	before := list.state.Load()
	list.Add(100)
	after := list.state.Load()
	require.Len(t, before.data, 100)
	require.Len(t, after.data, 101)
	if cap(before.data) > 100 {
		// same array
		assert.Same(t, &before.data[0], &after.data[0])
	}

	// a removal copies: the next append does not write the element of index 100 of after
	list.RemoveAt(100)
	list.Add(-1)
	assert.Equal(t, 100, after.data[100])
	assert.Equal(t, -1, list.Get(100))

	// appends cost amortized O(1): a state and, from time to time, a larger array
	allocs := testing.AllocsPerRun(1000, func() {
		list.Add(0)
	})
	assert.LessOrEqual(t, allocs, 2.0)
}

// The goroutine holding the lock of a list reads it.
func TestEStoreList_ReadsWhileLocked(t *testing.T) {
	list := newLockFreeTestList(t)
	list.Add("a")
	list.Add("b")
	list.Lock()
	runWithin(t, func() {
		assert.Equal(t, 2, list.Size())
		assert.Equal(t, "b", list.Get(1))
		assert.Equal(t, 1, list.IndexOf("b"))
		assert.True(t, list.Contains("a"))
		assert.Equal(t, []any{"a", "b"}, list.ToArray())
		assert.Equal(t, []any{"a", "b"}, list.GetUnResolvedList().ToArray())
		n := 0
		for range list.All() {
			n++
		}
		assert.Equal(t, 2, n)
	})
	list.Unlock()
}

// The goroutine holding the lock of an object (as the SQL encoder does) reads it, values read
// from the store or created on a miss included.
func TestEStoreObjects_ReadsWhileLocked(t *testing.T) {
	m := newReadBenchModel()
	reflective := newTestReflectiveEStoreEObject(nil)
	reflective.SetEClass(m.class)
	store := NewEStoreEObjectImpl(true)
	store.SetEClass(m.class)
	for _, o := range []interface {
		EObject
		ELockProvider
	}{reflective, store} {
		o.ESet(m.attributes[0], "a")
		o.Lock()
		runWithin(t, func() {
			assert.Equal(t, "a", o.EGet(m.attributes[0]))
			assert.True(t, o.EIsSet(m.attributes[0]))
			// created on a miss
			list := o.EGet(m.refs).(EList)
			assert.Equal(t, 0, list.Size())
			assert.Same(t, list, o.EGet(m.refs))
		})
		o.Unlock()
	}
}

// Concurrent misses of a list feature create one list: the first cached is returned to all.
func TestEStoreObjects_ConcurrentMisses(t *testing.T) {
	m := newReadBenchModel()
	for range 20 {
		compact := newTestCompactEStoreEObject(nil)
		compact.SetEClass(m.class)
		reflective := newTestReflectiveEStoreEObject(nil)
		reflective.SetEClass(m.class)
		store := NewEStoreEObjectImpl(true)
		store.SetEClass(m.class)
		for _, o := range []EObject{compact, reflective, store} {
			lists := make([]any, 8)
			var group sync.WaitGroup
			for i := range lists {
				group.Go(func() {
					lists[i] = o.EGet(m.refs)
				})
			}
			group.Wait()
			for _, l := range lists {
				require.Same(t, lists[0], l)
			}
		}
	}
}

// Readers read the values of objects without lock while they are set and unset: they read a
// value set, and the last value set wins.
func TestEStoreObjects_ConcurrentReadsDuringWrites(t *testing.T) {
	m := newReadBenchModel()
	compact := newTestCompactEStoreEObject(nil)
	compact.SetEClass(m.class)
	reflective := newTestReflectiveEStoreEObject(nil)
	reflective.SetEClass(m.class)
	store := NewEStoreEObjectImpl(true)
	store.SetEClass(m.class)
	for _, o := range []EObject{compact, reflective, store} {
		o.ESet(m.attributes[1], "fixed")
		var done atomic.Bool
		var group sync.WaitGroup
		for range 4 {
			group.Go(func() {
				for !done.Load() {
					if v := o.EGet(m.attributes[1]); v != "fixed" {
						t.Errorf("%v instead of fixed", v)
						return
					}
					if v := o.EGet(m.attributes[0]); v != nil {
						if _, isString := v.(string); !isString {
							t.Errorf("%v is not a value set", v)
							return
						}
					}
					o.EGet(m.refs)
				}
			})
		}
		for i := range 1000 {
			o.ESet(m.attributes[0], fmt.Sprint(i))
			if i%3 == 0 {
				o.EUnset(m.attributes[0])
			}
		}
		o.ESet(m.attributes[0], "last")
		done.Store(true)
		group.Wait()
		assert.Equal(t, "last", o.EGet(m.attributes[0]))
	}
}

// A list created on a miss racing a change of the cache flag of the object is created again: the
// list kept has the cache flag of the object.
func TestEStoreEObjectImpl_MissRacingSetCache(t *testing.T) {
	m := newReadBenchModel()
	for range 50 {
		o := NewEStoreEObjectImpl(true)
		o.SetEClass(m.class)
		var group sync.WaitGroup
		for range 4 {
			group.Go(func() {
				o.EGet(m.refs)
			})
		}
		for i := range 20 {
			o.SetCache(i%2 == 0)
		}
		group.Wait()
		list := o.EGet(m.refs).(*EStoreList)
		assert.Equal(t, o.IsCache(), list.IsCache())
	}
}

// A value read from the store while the store of the object changes is not cached: the reader
// reads it again from the new store.
func TestEStoreEObjectImpl_MissRacingSetEStore(t *testing.T) {
	m := newReadBenchModel()
	attribute := m.attributes[0]
	o := NewEStoreEObjectImpl(true)
	o.SetEClass(m.class)
	oldStore := NewMockEStore(t)
	newStore := NewMockEStore(t)
	o.SetEStore(oldStore)

	reading := make(chan struct{})
	release := make(chan struct{})
	oldStore.EXPECT().Get(o, attribute, NO_INDEX).RunAndReturn(func(EObject, EStructuralFeature, int) any {
		close(reading)
		<-release
		return "old"
	}).Once()
	newStore.EXPECT().Get(o, attribute, NO_INDEX).Return("new").Once()

	var value any
	var group sync.WaitGroup
	group.Go(func() {
		value = o.EGet(attribute)
	})
	<-reading
	o.SetEStore(newStore)
	close(release)
	group.Wait()
	assert.Equal(t, "new", value)
	assert.Equal(t, "new", o.EGet(attribute))
}

// The epoch of the cache: a value read with a configuration replaced meanwhile fails to be
// cached, and ClearCache keeps the epoch.
func TestEStoreEObjectImpl_CacheEpoch(t *testing.T) {
	m := newReadBenchModel()
	o := NewEStoreEObjectImpl(true)
	o.SetEClass(m.class)
	epoch := o.config.Load().getEpoch()
	o.SetCache(false)
	_, cached := o.cacheIfAbsentAt(epoch, 1, "stale")
	assert.False(t, cached)
	assert.Nil(t, o.getCached(1))

	current := o.config.Load().getEpoch()
	assert.Equal(t, epoch+1, current)
	value, cached := o.cacheIfAbsentAt(current, 1, "fresh")
	assert.True(t, cached)
	assert.Equal(t, "fresh", value)
	o.ClearCache()
	assert.Nil(t, o.getCached(1))
	assert.Equal(t, current, o.cache.Load().epochOf())
	_, cached = o.cacheIfAbsentAt(current, 1, "again")
	assert.True(t, cached)
}
