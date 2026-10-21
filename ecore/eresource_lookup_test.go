package ecore

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEResourceImpl_GetEObject_Fragments(t *testing.T) {
	r := NewEResourceImpl()
	r.SetURI(NewURI("file:///a.xml"))
	ePackage := GetFactory().CreateEPackage()
	r.GetContents().Add(ePackage)
	assert.NotPanics(t, func() { assert.Nil(t, r.GetEObject("?")) })
	assert.Equal(t, ePackage, r.GetEObject("/"))
	assert.Equal(t, ePackage, r.GetEObject("/0"))
	// a root segment which is not an index designates no object
	assert.Nil(t, r.GetEObject("/abc"))
}

func TestEResourceImpl_NoURI(t *testing.T) {
	r := NewEResourceImpl()
	require.NotPanics(t, func() { r.LoadWithReader(strings.NewReader("x"), nil) })
	assert.Equal(t, 1, r.GetErrors().Size())
	var w strings.Builder
	require.NotPanics(t, func() { r.SaveWithWriter(&w, nil) })
	assert.Equal(t, 1, r.GetErrors().Size())
}

// a path segment designating no object gives no object
func TestGetEObject_MissingSegment(t *testing.T) {
	ePackage := GetFactory().CreateEPackage()
	assert.NotPanics(t, func() { assert.Nil(t, GetEObject(ePackage, "@eClassifiers.5")) })
}

// proxies resolving to each other are left unresolved instead of recursing forever
func TestResolveInResourceSet_ProxyCycle(t *testing.T) {
	f := GetFactory()
	rs := NewEResourceSetImpl()
	ra := rs.CreateResource(NewURI("memory:/a"))
	rb := rs.CreateResource(NewURI("memory:/b"))
	proxyA := f.CreateEClass()
	proxyA.(EObjectInternal).ESetProxyURI(NewURI("memory:/a#/"))
	proxyB := f.CreateEClass()
	proxyB.(EObjectInternal).ESetProxyURI(NewURI("memory:/b#/"))
	// proxyA resolves to proxyB (root of a), which resolves to proxyA (root of b)
	ra.GetContents().Add(proxyB)
	rb.GetContents().Add(proxyA)
	var resolved EObject
	require.NotPanics(t, func() { resolved = ResolveInResourceSet(proxyA, rs) })
	assert.Equal(t, proxyA, resolved)
}
