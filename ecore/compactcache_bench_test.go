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
	"runtime"
	"testing"
)

// Benchmarks of the cache of a store object (CompactEStoreEObjectImpl): the first read of each
// feature of an object caches its value (cacheIfAbsent), as EDynamicGet does, then the reads hit
// the cache.

// newCacheBenchClass returns a class of n attributes.
func newCacheBenchClass(n int) EClass {
	class := GetFactory().CreateEClass()
	class.SetName(fmt.Sprintf("Class%d", n))
	for i := range n {
		attribute := GetFactory().CreateEAttribute()
		attribute.SetName(fmt.Sprintf("attribute%d", i))
		attribute.SetEType(GetPackage().GetEString())
		class.GetEStructuralFeatures().Add(attribute)
	}
	return class
}

// cacheBenchCases are classes of features, of which the first read ones are cached.
var cacheBenchCases = []struct {
	features, read int
}{
	{1, 1},
	{4, 1},
	{4, 4},
	{8, 2},
	{8, 8},
	{16, 4},
	{16, 16},
	{40, 40},
}

// fillCache caches the values of the first read features of o.
func fillCache(o *testCompactEStoreEObject, read int) {
	for featureID := range read {
		o.cacheIfAbsent(featureID, "value")
	}
}

// BenchmarkCompactCache_Fill caches the values of the features of a new object, one after the
// other: the work and the allocations of the first read of each.
func BenchmarkCompactCache_Fill(b *testing.B) {
	for _, c := range cacheBenchCases {
		b.Run(fmt.Sprintf("features=%d/read=%d", c.features, c.read), func(b *testing.B) {
			class := newCacheBenchClass(c.features)
			b.ReportAllocs()
			for b.Loop() {
				o := newTestCompactEStoreEObject(nil)
				o.SetEClass(class)
				fillCache(o, c.read)
			}
		})
	}
}

// BenchmarkCompactCache_Live caches the values of 100,000 objects and reports the bytes each
// object holds once filled (its cache, with the object itself), after a GC.
func BenchmarkCompactCache_Live(b *testing.B) {
	const objects = 100_000
	for _, c := range cacheBenchCases {
		b.Run(fmt.Sprintf("features=%d/read=%d", c.features, c.read), func(b *testing.B) {
			class := newCacheBenchClass(c.features)
			var perObject float64
			for b.Loop() {
				var before, after runtime.MemStats
				runtime.GC()
				runtime.ReadMemStats(&before)
				held := make([]*testCompactEStoreEObject, objects)
				for i := range held {
					o := newTestCompactEStoreEObject(nil)
					o.SetEClass(class)
					fillCache(o, c.read)
					held[i] = o
				}
				runtime.GC()
				runtime.ReadMemStats(&after)
				perObject = float64(after.HeapAlloc-before.HeapAlloc) / objects
				runtime.KeepAlive(held)
			}
			b.ReportMetric(perObject, "B/object")
		})
	}
}

// BenchmarkCompactCache_Read reads the cached values of a filled object.
func BenchmarkCompactCache_Read(b *testing.B) {
	for _, c := range cacheBenchCases {
		b.Run(fmt.Sprintf("features=%d/read=%d", c.features, c.read), func(b *testing.B) {
			o := newTestCompactEStoreEObject(nil)
			o.SetEClass(newCacheBenchClass(c.features))
			fillCache(o, c.read)
			b.ReportAllocs()
			for b.Loop() {
				for featureID := range c.read {
					if o.cachedValue(featureID) == nil {
						b.Fatal("not cached")
					}
				}
			}
		})
	}
}
