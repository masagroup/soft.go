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

// eAllContentsBenchPackage is a package of a class Node with a many and a single containment of
// nodes, and a leaf class Leaf with an attribute only.
type eAllContentsBenchPackage struct {
	nodeClass EClass
	leafClass EClass
	children  EReference
	single    EReference
	leaves    EReference
}

func newEAllContentsBenchPackage() *eAllContentsBenchPackage {
	factory := GetFactory()
	ePackage := factory.CreateEPackage()
	ePackage.SetName("bench")
	ePackage.SetNsURI("http://bench/eallcontents")
	ePackage.SetNsPrefix("bench")

	leafClass := factory.CreateEClass()
	leafClass.SetName("Leaf")
	name := factory.CreateEAttribute()
	name.SetName("name")
	name.SetEType(GetPackage().GetEString())
	leafClass.GetEStructuralFeatures().Add(name)

	nodeClass := factory.CreateEClass()
	nodeClass.SetName("Node")
	newContainment := func(name string, eType EClass, many bool) EReference {
		reference := factory.CreateEReference()
		reference.SetName(name)
		reference.SetEType(eType)
		reference.SetContainment(true)
		if many {
			reference.SetUpperBound(UNBOUNDED_MULTIPLICITY)
		}
		nodeClass.GetEStructuralFeatures().Add(reference)
		return reference
	}
	children := newContainment("children", nodeClass, true)
	single := newContainment("single", nodeClass, false)
	leaves := newContainment("leaves", leafClass, true)

	ePackage.GetEClassifiers().AddAll(NewImmutableEList([]any{leafClass, nodeClass}))
	ePackage.SetEFactoryInstance(factory.CreateEFactory())
	return &eAllContentsBenchPackage{
		nodeClass: nodeClass,
		leafClass: leafClass,
		children:  children,
		single:    single,
		leaves:    leaves,
	}
}

// newTree returns a tree of nodes of the given depth: each node holds fanOut nodes, a single node
// every other level, and leaves leaves.
func (p *eAllContentsBenchPackage) newTree(depth, fanOut, leaves int) (EObject, int) {
	factory := p.nodeClass.GetEPackage().GetEFactoryInstance()
	count := 0
	var build func(level int) EObject
	build = func(level int) EObject {
		node := factory.Create(p.nodeClass)
		count++
		leafList := node.EGet(p.leaves).(EList)
		for range leaves {
			leafList.Add(factory.Create(p.leafClass))
			count++
		}
		if level < depth {
			children := node.EGet(p.children).(EList)
			for range fanOut {
				children.Add(build(level + 1))
			}
			if level%2 == 0 {
				node.ESet(p.single, build(level+1))
			}
		}
		return node
	}
	root := build(0)
	return root, count - 1
}

// newTreeEAllContentsIterator is the iterator EAllContents was: a tree iterator of the iterators
// of EContents.
func newTreeEAllContentsIterator(object EObject) EIterator {
	return newTreeIterator(object, false, func(o any) EIterator {
		return o.(EObject).EContents().Iterator()
	})
}

// newTreeResourceAllContentsIterator is the iterator GetAllContents was: a tree iterator of the
// roots of the resource and of the iterators of EContents.
func newTreeResourceAllContentsIterator(resource EResource) EIterator {
	return newTreeIterator(resource, false, func(o any) EIterator {
		if r, _ := o.(EResource); r != nil {
			return r.GetContents().Iterator()
		}
		return o.(EObject).EContents().Iterator()
	})
}

// BenchmarkEAllContents walks the contents of trees with the tree iterator EAllContents was
// (tree), the iterator of EAllContents (iterator) and the sequence AllContents (seq).
func BenchmarkEAllContents(b *testing.B) {
	p := newEAllContentsBenchPackage()
	for _, c := range []struct{ depth, fanOut, leaves int }{
		{4, 8, 0},
		{4, 8, 4},
		{10, 2, 1},
	} {
		root, count := p.newTree(c.depth, c.fanOut, c.leaves)
		for _, v := range []struct {
			name string
			walk func(EObject) int
		}{
			{"tree", func(o EObject) int {
				n := 0
				for it := newTreeEAllContentsIterator(o); it.HasNext(); {
					it.Next()
					n++
				}
				return n
			}},
			{"iterator", func(o EObject) int {
				n := 0
				for it := o.EAllContents(); it.HasNext(); {
					it.Next()
					n++
				}
				return n
			}},
			{"seq", func(o EObject) int {
				n := 0
				for range AllContents(o) {
					n++
				}
				return n
			}},
		} {
			b.Run(fmt.Sprintf("depth=%d/fanout=%d/leaves=%d/impl=%s", c.depth, c.fanOut, c.leaves, v.name), func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					if n := v.walk(root); n != count {
						b.Fatalf("%d objects, expected %d", n, count)
					}
				}
				b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N*count), "ns/object")
			})
		}
	}
}

// BenchmarkAllContents_Resource walks the contents of a resource with the tree iterator
// GetAllContents was (tree), GetAllContents (iterator) and AllContents (seq).
func BenchmarkAllContents_Resource(b *testing.B) {
	p := newEAllContentsBenchPackage()
	resource, expected := newAllContentsResource(p, 4, 8, 4)
	count := len(expected)
	for _, v := range []struct {
		name string
		walk func(EResource) int
	}{
		{"tree", func(r EResource) int {
			n := 0
			for it := newTreeResourceAllContentsIterator(r); it.HasNext(); {
				it.Next()
				n++
			}
			return n
		}},
		{"iterator", func(r EResource) int {
			n := 0
			for it := r.GetAllContents(); it.HasNext(); {
				it.Next()
				n++
			}
			return n
		}},
		{"seq", func(r EResource) int {
			n := 0
			for range AllContents(r) {
				n++
			}
			return n
		}},
	} {
		b.Run("impl="+v.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if n := v.walk(resource); n != count {
					b.Fatalf("%d objects, expected %d", n, count)
				}
			}
			b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N*count), "ns/object")
		})
	}
}

// BenchmarkEAllContentsWithClass walks the instances of a class in the contents of a tree, the
// leaves, at its bottom, and the nodes, everywhere: with EAllContentsWithClass (class), which
// builds its table, EAllContentsWithTable (table), and their sequences (seqclass, seqtable).
func BenchmarkEAllContentsWithClass(b *testing.B) {
	p := newEAllContentsBenchPackage()
	root, _ := p.newTree(4, 8, 4)
	for _, end := range []EClass{p.leafClass, p.nodeClass} {
		count := len(instancesInContents(root, end))
		table := NewEClassTransitionsTable(root.EClass(), end)
		for _, v := range []struct {
			name string
			walk func() int
		}{
			{"class", func() int {
				n := 0
				for it := EAllContentsWithClass(root, end); it.HasNext(); {
					it.Next()
					n++
				}
				return n
			}},
			{"table", func() int {
				n := 0
				for it := EAllContentsWithTable(root, table); it.HasNext(); {
					it.Next()
					n++
				}
				return n
			}},
			{"seqclass", func() int {
				n := 0
				for range AllContentsWithClass(root, end) {
					n++
				}
				return n
			}},
			{"seqtable", func() int {
				n := 0
				for range AllContentsWithTable(root, table) {
					n++
				}
				return n
			}},
		} {
			b.Run(fmt.Sprintf("end=%s/impl=%s", end.GetName(), v.name), func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					if n := v.walk(); n != count {
						b.Fatalf("%d objects, expected %d", n, count)
					}
				}
				b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N*count), "ns/object")
			})
		}
	}
}
