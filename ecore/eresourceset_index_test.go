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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEResourceSetIndex_AddLookup(t *testing.T) {
	rs := NewEResourceSetImpl()
	assert.Empty(t, rs.resourceIndex)

	uri1 := NewURI("test://file1.t")
	uri2 := NewURI("test://file2.t")
	r1 := rs.CreateResource(uri1)
	r2 := NewEResourceImpl()
	r2.SetURI(uri2)
	rs.GetResources().Add(r2)

	// indexed on add
	assert.Len(t, rs.resourceIndex, 2)
	assert.Equal(t, r1, rs.resourceIndex[uri1.String()])
	assert.Equal(t, r2, rs.resourceIndex[uri2.String()])

	// found through the index, with the same or an equal uri
	assert.Equal(t, r1, rs.GetResource(uri1, false))
	assert.Equal(t, r2, rs.GetResource(NewURI("test://file2.t"), false))
	assert.Nil(t, rs.GetResource(NewURI("test://file3.t"), false))

	// a resource without uri is not indexed
	rs.GetResources().Add(NewEResourceImpl())
	assert.Len(t, rs.resourceIndex, 2)
}

func TestEResourceSetIndex_URIChange(t *testing.T) {
	rs := NewEResourceSetImpl()
	oldURI := NewURI("test://old.t")
	newURI := NewURI("test://new.t")
	r := rs.CreateResource(oldURI)
	require.Equal(t, r, rs.GetResource(oldURI, false))

	r.SetURI(newURI)
	assert.Len(t, rs.resourceIndex, 1)
	assert.Equal(t, r, rs.resourceIndex[newURI.String()])
	assert.Nil(t, rs.GetResource(oldURI, false))
	assert.Equal(t, r, rs.GetResource(newURI, false))

	r.SetURI(nil)
	assert.Empty(t, rs.resourceIndex)
	assert.Nil(t, rs.GetResource(newURI, false))

	// the uri of a resource outside the set does not touch the index
	other := NewEResourceImpl()
	other.SetURI(oldURI)
	assert.Empty(t, rs.resourceIndex)
	assert.Nil(t, rs.GetResource(oldURI, false))
}

func TestEResourceSetIndex_Remove(t *testing.T) {
	rs := NewEResourceSetImpl()
	uri1 := NewURI("test://file1.t")
	uri2 := NewURI("test://file2.t")
	uri3 := NewURI("test://file3.t")
	r1 := rs.CreateResource(uri1)
	r2 := rs.CreateResource(uri2)
	r3 := rs.CreateResource(uri3)

	rs.GetResources().Remove(r1)
	assert.Len(t, rs.resourceIndex, 2)
	assert.Nil(t, rs.GetResource(uri1, false))
	assert.Equal(t, r2, rs.GetResource(uri2, false))

	// a resource moved to another set is not found anymore
	rs2 := NewEResourceSetImpl()
	rs2.GetResources().Add(r2)
	assert.Nil(t, rs.GetResource(uri2, false))
	assert.Equal(t, r2, rs2.GetResource(uri2, false))
	assert.Len(t, rs.resourceIndex, 1)

	// set
	r4 := NewEResourceImpl()
	r4.SetURI(NewURI("test://file4.t"))
	rs.GetResources().Set(0, r4)
	assert.Nil(t, rs.GetResource(uri3, false))
	assert.Equal(t, r4, rs.GetResource(r4.GetURI(), false))
	assert.Len(t, rs.resourceIndex, 1)

	// clear
	rs.GetResources().Add(r3)
	rs.GetResources().Clear()
	assert.Empty(t, rs.resourceIndex)
	assert.Nil(t, rs.GetResource(uri3, false))
	assert.Nil(t, rs.GetResource(r4.GetURI(), false))
}

func TestEResourceSetIndex_DuplicateURI(t *testing.T) {
	rs := NewEResourceSetImpl()
	uri := NewURI("test://file.t")
	r1 := rs.CreateResource(uri)
	r2 := rs.CreateResource(uri)
	// the first resource of the list wins
	assert.Equal(t, r1, rs.GetResource(uri, false))

	// once it is removed, the second one is found by scanning the list, and indexed
	rs.GetResources().Remove(r1)
	assert.Empty(t, rs.resourceIndex)
	assert.Equal(t, r2, rs.GetResource(uri, false))
	assert.Equal(t, r2, rs.resourceIndex[uri.String()])
}

func TestEResourceSetIndex_URIConverterPathPrefix(t *testing.T) {
	rs := NewEResourceSetImpl()
	rs.GetURIConverter().GetURIMap()[*NewURI("test://h/old/")] = *NewURI("test://h/new/")

	// indexed by uri and by normalized uri
	r := rs.CreateResource(NewURI("test://h/old/file.t"))
	assert.Len(t, rs.resourceIndex, 2)
	assert.Equal(t, r, rs.resourceIndex["test://h/old/file.t"])
	assert.Equal(t, r, rs.resourceIndex["test://h/new/file.t"])
	assert.Equal(t, r, rs.GetResource(NewURI("test://h/old/file.t"), false))
	assert.Equal(t, r, rs.GetResource(NewURI("test://h/new/file.t"), false))

	// a mapping added after the resource was indexed : found by scanning the list, and re-indexed
	rs.GetURIConverter().GetURIMap()[*NewURI("test://h/new/")] = *NewURI("test://h/newer/")
	assert.Equal(t, r, rs.GetResource(NewURI("test://h/newer/file.t"), false))
	assert.Equal(t, r, rs.resourceIndex["test://h/newer/file.t"])
	assert.Equal(t, r, rs.GetResource(NewURI("test://h/old/file.t"), false))
	assert.Equal(t, r, rs.GetResource(NewURI("test://h/new/file.t"), false))

	// a stale index entry is not trusted
	rs.GetURIConverter().GetURIMap()[*NewURI("test://h/other/")] = *NewURI("test://h/new/")
	other := rs.CreateResource(NewURI("test://unrelated.t"))
	rs.resourceIndex["test://h/newer/file.t"] = other
	assert.Equal(t, r, rs.GetResource(NewURI("test://h/other/file.t"), false))
	assert.Equal(t, r, rs.resourceIndex["test://h/newer/file.t"])

	// a new converter rebuilds the index
	rs.SetURIConverter(NewEURIConverterImpl())
	assert.Len(t, rs.resourceIndex, 2)
	assert.Equal(t, r, rs.resourceIndex["test://h/old/file.t"])
	assert.Equal(t, other, rs.resourceIndex["test://unrelated.t"])
	assert.Equal(t, r, rs.GetResource(NewURI("test://h/old/file.t"), false))
	assert.Nil(t, rs.GetResource(NewURI("test://h/new/file.t"), false))
}

func TestEResourceSetIndex_URIConverterScheme(t *testing.T) {
	// a scheme mapped to a directory, as sw.recorder maps physical: to the directory of the database
	rs := NewEResourceSetImpl()
	schemeURI := NewURIBuilder(nil).SetScheme("physical").URI()
	dirURI := NewURIBuilder(nil).SetScheme("file").SetPath("/data/physical").URI()
	rs.GetURIConverter().GetURIMap()[*schemeURI] = *dirURI

	docURI := NewURIBuilder(nil).SetScheme("physical").SetPath("/doc.xml").URI()
	fileURI := rs.GetURIConverter().Normalize(docURI)
	require.Equal(t, "file:/data/physical/doc.xml", fileURI.String())

	r := rs.CreateResource(docURI)
	assert.Len(t, rs.resourceIndex, 2)
	assert.Equal(t, r, rs.resourceIndex[docURI.String()])
	assert.Equal(t, r, rs.resourceIndex[fileURI.String()])
	assert.Equal(t, r, rs.GetResource(docURI, false))
	assert.Equal(t, r, rs.GetResource(NewURI("physical:/doc.xml"), false))
	assert.Equal(t, r, rs.GetResource(fileURI, false))
	assert.Nil(t, rs.GetResource(NewURI("physical:/other.xml"), false))

	// through GetEObject, as proxies are resolved
	root := GetFactory().CreateEPackage()
	r.GetContents().Add(root)
	assert.Equal(t, root, rs.GetEObject(NewURI("physical:/doc.xml#/"), false))
	assert.Equal(t, root, rs.GetEObject(NewURI("file:/data/physical/doc.xml#/"), false))
}

func TestEResourceSetIndex_LoadOnDemand(t *testing.T) {
	uri := NewURI("test://file.t")
	rs := NewEResourceSetImpl()

	// an indexed resource is loaded on demand
	mockResource := NewMockEResourceInternal(t)
	mockResource.EXPECT().GetURI().Return(uri)
	mockResource.EXPECT().BasicSetResourceSet(rs, nil).Return(nil).Once()
	rs.GetResources().Add(mockResource)
	assert.Equal(t, mockResource, rs.resourceIndex[uri.String()])

	assert.Equal(t, mockResource, rs.GetResource(uri, false))

	mockResource.EXPECT().IsLoaded().Once().Return(false)
	mockResource.EXPECT().Load().Once()
	assert.Equal(t, mockResource, rs.GetResource(uri, true))

	mockResource.EXPECT().IsLoaded().Once().Return(true)
	assert.Equal(t, mockResource, rs.GetResource(uri, true))

	// an unknown uri creates the resource, which is indexed
	created := rs.GetResource(NewURI("memory://created.t"), true)
	require.NotNil(t, created)
	assert.Equal(t, created, rs.resourceIndex["memory://created.t"])
	assert.Equal(t, created, rs.GetResource(NewURI("memory://created.t"), true))
}
