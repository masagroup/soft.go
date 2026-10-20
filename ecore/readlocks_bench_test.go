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
	"testing"
)

// Benchmarks of the read paths of classes, store objects and store lists, serial and parallel:
// the locks they take show in parallel runs.

type readBenchModel struct {
	super      EClass
	class      EClass
	attributes []EAttribute
	refs       EReference
	contained  EReference
	operation  EOperation
}

func newReadBenchModel() *readBenchModel {
	m := &readBenchModel{}
	factory := GetFactory()
	m.super = factory.CreateEClass()
	m.super.SetName("Super")
	m.class = factory.CreateEClass()
	m.class.SetName("Class")
	m.class.GetESuperTypes().Add(m.super)
	for i := range 6 {
		attribute := factory.CreateEAttribute()
		attribute.SetName(fmt.Sprintf("attribute%d", i))
		attribute.SetEType(GetPackage().GetEString())
		m.attributes = append(m.attributes, attribute)
		if i%2 == 0 {
			m.super.GetEStructuralFeatures().Add(attribute)
		} else {
			m.class.GetEStructuralFeatures().Add(attribute)
		}
	}
	m.refs = factory.CreateEReference()
	m.refs.SetName("refs")
	m.refs.SetEType(m.class)
	m.refs.SetUpperBound(-1)
	m.refs.SetResolveProxies(true)
	m.class.GetEStructuralFeatures().Add(m.refs)
	m.contained = factory.CreateEReference()
	m.contained.SetName("contained")
	m.contained.SetEType(m.class)
	m.contained.SetUpperBound(-1)
	m.contained.SetContainment(true)
	m.super.GetEStructuralFeatures().Add(m.contained)
	m.operation = factory.CreateEOperation()
	m.operation.SetName("operation")
	m.super.GetEOperations().Add(m.operation)
	return m
}

func runSerialAndParallel(b *testing.B, op func()) {
	b.Run("serial", func(b *testing.B) {
		for b.Loop() {
			op()
		}
	})
	b.Run("parallel", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				op()
			}
		})
	})
}

func BenchmarkReadLocks_EClass(b *testing.B) {
	m := newReadBenchModel()
	class := m.class
	var sink int
	runSerialAndParallel(b, func() {
		sink += class.GetEAllAttributes().Size()
		sink += class.GetEAllReferences().Size()
		sink += class.GetEAllContainments().Size()
		sink += class.GetEAllCrossReferences().Size()
		sink += class.GetEAllSuperTypes().Size()
		sink += class.GetEAllOperations().Size()
		sink += class.GetEContainmentFeatures().Size()
		sink += class.GetECrossReferenceFeatures().Size()
		sink += class.GetFeatureCount()
		if class.GetEStructuralFeatureFromName("attribute3") != nil {
			sink++
		}
		if class.GetOverride(m.operation) != nil {
			sink++
		}
	})
	_ = sink
}

func benchmarkObjectEGet(b *testing.B, m *readBenchModel, o EObject) {
	for i, attribute := range m.attributes {
		o.ESet(attribute, fmt.Sprintf("value%d", i))
	}
	o.EGet(m.refs)
	var sink int
	runSerialAndParallel(b, func() {
		for _, attribute := range m.attributes {
			if o.EGet(attribute) != nil {
				sink++
			}
		}
		sink += o.EGet(m.refs).(EList).Size()
	})
	_ = sink
}

func BenchmarkReadLocks_CompactEStoreEObjectEGet(b *testing.B) {
	m := newReadBenchModel()
	o := newTestCompactEStoreEObject(nil)
	o.SetEClass(m.class)
	benchmarkObjectEGet(b, m, o)
}

func BenchmarkReadLocks_ReflectiveEStoreEObjectEGet(b *testing.B) {
	m := newReadBenchModel()
	o := newTestReflectiveEStoreEObject(nil)
	o.SetEClass(m.class)
	benchmarkObjectEGet(b, m, o)
}

func BenchmarkReadLocks_EStoreEObjectEGet(b *testing.B) {
	m := newReadBenchModel()
	o := NewEStoreEObjectImpl(true)
	o.SetEClass(m.class)
	benchmarkObjectEGet(b, m, o)
}

func newReadBenchList(m *readBenchModel, size int) (EList, []EObject) {
	owner := newTestCompactEStoreEObject(nil)
	owner.SetEClass(m.class)
	list := owner.EGet(m.refs).(EList)
	targets := make([]EObject, size)
	for i := range targets {
		target := newTestCompactEStoreEObject(nil)
		target.SetEClass(m.class)
		targets[i] = target
		list.Add(target)
	}
	return list, targets
}

func BenchmarkReadLocks_EStoreListGet(b *testing.B) {
	m := newReadBenchModel()
	list, _ := newReadBenchList(m, 100)
	var sink int
	runSerialAndParallel(b, func() {
		for i := 0; i < list.Size(); i++ {
			if list.Get(i) != nil {
				sink++
			}
		}
	})
	_ = sink
}

func BenchmarkReadLocks_EStoreListAll(b *testing.B) {
	m := newReadBenchModel()
	list, targets := newReadBenchList(m, 100)
	last := targets[len(targets)-1]
	var sink int
	runSerialAndParallel(b, func() {
		for v := range list.All() {
			if v != nil {
				sink++
			}
		}
		sink += list.IndexOf(last)
	})
	_ = sink
}

// Appends to a list: the cost per element must not grow with the size of the list.
func BenchmarkReadLocks_EStoreListAdd(b *testing.B) {
	m := newReadBenchModel()
	attribute := GetFactory().CreateEAttribute()
	attribute.SetName("values")
	attribute.SetEType(GetPackage().GetEString())
	attribute.SetUpperBound(-1)
	attribute.SetUnique(false)
	m.class.GetEStructuralFeatures().Add(attribute)
	for _, size := range []int{100, 10000} {
		b.Run(fmt.Sprint(size), func(b *testing.B) {
			for b.Loop() {
				list := NewEStoreList(newTestCompactEStoreEObject(nil), attribute, nil)
				list.SetCache(true)
				for range size {
					list.Add("v")
				}
			}
		})
	}
}
