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

	"github.com/stretchr/testify/require"
)

// crossReferencedModel is a dynamic model : a Root holds Nodes, a Node refers to a Node
type crossReferencedModel struct {
	ePackage  EPackage
	rootClass EClass
	nodeClass EClass
	nodes     EReference
	ref       EReference
}

func newCrossReferencedModel() *crossReferencedModel {
	m := &crossReferencedModel{}
	m.ePackage = GetFactory().CreateEPackage()
	m.ePackage.SetName("bench")
	m.ePackage.SetNsURI("http://bench")
	m.nodeClass = GetFactory().CreateEClass()
	m.nodeClass.SetName("Node")
	m.ref = GetFactory().CreateEReference()
	m.ref.SetName("ref")
	m.ref.SetEType(m.nodeClass)
	m.ref.SetResolveProxies(true)
	m.nodeClass.GetEStructuralFeatures().Add(m.ref)
	m.rootClass = GetFactory().CreateEClass()
	m.rootClass.SetName("Root")
	m.nodes = GetFactory().CreateEReference()
	m.nodes.SetName("nodes")
	m.nodes.SetEType(m.nodeClass)
	m.nodes.SetContainment(true)
	m.nodes.SetUpperBound(UNBOUNDED_MULTIPLICITY)
	m.rootClass.GetEStructuralFeatures().Add(m.nodes)
	m.ePackage.GetEClassifiers().Add(m.rootClass)
	m.ePackage.GetEClassifiers().Add(m.nodeClass)
	return m
}

// newCrossReferencedResourceSet builds a resource set of nbResources resources, each
// holding a root of nbNodes nodes. Every node refers, through a proxy, to a node of
// another resource, so that resolving the set performs nbResources*nbNodes resource
// lookups by URI.
func (m *crossReferencedModel) newResourceSet(nbResources, nbNodes int) EResourceSet {
	rs := NewEResourceSetImpl()
	for i := range nbResources {
		resource := NewEResourceImpl()
		resource.SetURI(NewURI(fmt.Sprintf("memory://bench/doc%d.xml", i)))
		root := m.ePackage.GetEFactoryInstance().Create(m.rootClass)
		nodes := root.EGet(m.nodes).(EList)
		// cross references : node j of resource i refers to node j of resource (i+1)%n
		target := (i + 1) % nbResources
		for j := range nbNodes {
			node := m.ePackage.GetEFactoryInstance().Create(m.nodeClass)
			proxy := m.ePackage.GetEFactoryInstance().Create(m.nodeClass)
			proxy.(EObjectInternal).ESetProxyURI(NewURI(fmt.Sprintf("memory://bench/doc%d.xml#//@nodes.%d", target, j)))
			node.ESet(m.ref, proxy)
			nodes.Add(node)
		}
		resource.GetContents().Add(root)
		rs.GetResources().Add(resource)
	}
	return rs
}

// countProxies returns the number of unresolved references of the set
func (m *crossReferencedModel) countProxies(rs EResourceSet) int {
	count := 0
	for r := range rs.GetResources().All() {
		root := r.(EResource).GetContents().Get(0).(EObject)
		for n := range root.EGet(m.nodes).(EList).All() {
			if n.(EObject).EGetResolve(m.ref, false).(EObject).EIsProxy() {
				count++
			}
		}
	}
	return count
}

func TestCrossReferencedResourceSetResolves(t *testing.T) {
	m := newCrossReferencedModel()
	rs := m.newResourceSet(4, 3)
	require.Equal(t, 12, m.countProxies(rs))
	ResolveAllInResourceSet(rs)
	require.Equal(t, 0, m.countProxies(rs))
	for i := range 4 {
		root := rs.GetResources().Get(i).(EResource).GetContents().Get(0).(EObject)
		targetRoot := rs.GetResources().Get((i + 1) % 4).(EResource).GetContents().Get(0).(EObject)
		for j := range 3 {
			node := root.EGet(m.nodes).(EList).Get(j).(EObject)
			require.Equal(t, targetRoot.EGet(m.nodes).(EList).Get(j), node.EGet(m.ref), "resource %d node %d", i, j)
		}
	}
}

// BenchmarkEResourceSet_ResolveAll measures ResolveAllInResourceSet on resource sets of
// a growing number of resources, with a constant total number of cross references
// (proxies) to resolve, so that the cost per resolution shows how GetResource grows
// with the number of resources in the set.
func BenchmarkEResourceSet_ResolveAll(b *testing.B) {
	const nbProxies = 4096
	m := newCrossReferencedModel()
	for _, nbResources := range []int{4, 16, 64, 256} {
		nbNodes := nbProxies / nbResources
		b.Run(fmt.Sprintf("resources=%d", nbResources), func(b *testing.B) {
			for b.Loop() {
				b.StopTimer()
				rs := m.newResourceSet(nbResources, nbNodes)
				b.StartTimer()
				ResolveAllInResourceSet(rs)
			}
			b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N)/float64(nbProxies), "ns/proxy")
		})
	}
}

// BenchmarkEResourceSet_GetResource measures a lookup of an already loaded resource
// by URI, for resource sets of a growing number of resources, with the uri of the
// resource (same) or with a uri the converter maps to it (mapped).
func BenchmarkEResourceSet_GetResource(b *testing.B) {
	for _, mapped := range []bool{false, true} {
		for _, nbResources := range []int{4, 16, 64, 256} {
			name := "same"
			if mapped {
				name = "mapped"
			}
			b.Run(fmt.Sprintf("%s/resources=%d", name, nbResources), func(b *testing.B) {
				rs := NewEResourceSetImpl()
				rs.GetURIConverter().GetURIMap()[*NewURI("bench:/")] = *NewURI("memory:/bench/")
				uris := make([]*URI, nbResources)
				for i := range nbResources {
					resourceURI := NewURI(fmt.Sprintf("memory:/bench/doc%d.xml", i))
					rs.CreateResource(resourceURI)
					if mapped {
						uris[i] = NewURI(fmt.Sprintf("bench:/doc%d.xml", i))
					} else {
						uris[i] = resourceURI
					}
				}
				i := 0
				for b.Loop() {
					if rs.GetResource(uris[i%nbResources], false) == nil {
						b.Fatal("resource not found")
					}
					i++
				}
			})
		}
	}
}
