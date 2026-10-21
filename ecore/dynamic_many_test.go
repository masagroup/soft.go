package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// a class with a many containment and a many reference resolving proxies
func newDynamicManyClasses() (eClass EClass, children EReference, refs EReference) {
	f := GetFactory()
	ePackage := f.CreateEPackage()
	ePackage.SetName("p")
	ePackage.SetNsURI("http://p")
	ePackage.SetEFactoryInstance(f.CreateEFactory())
	eClass = f.CreateEClass()
	eClass.SetName("C")
	ePackage.GetEClassifiers().Add(eClass)
	children = f.CreateEReference()
	children.SetName("children")
	children.SetEType(eClass)
	children.SetContainment(true)
	children.SetUpperBound(UNBOUNDED_MULTIPLICITY)
	eClass.GetEStructuralFeatures().Add(children)
	refs = f.CreateEReference()
	refs.SetName("refs")
	refs.SetEType(eClass)
	refs.SetUpperBound(UNBOUNDED_MULTIPLICITY)
	eClass.GetEStructuralFeatures().Add(refs)
	return
}

// setting a many feature of a dynamic object replaces the elements of its list: the objects set
// in a containment are contained
func TestDynamicEObject_SetMany(t *testing.T) {
	eClass, children, refs := newDynamicManyClasses()
	eFactory := eClass.GetEPackage().GetEFactoryInstance()
	parent := eFactory.Create(eClass)
	list := parent.EGet(children).(EList)
	child1 := eFactory.Create(eClass)
	child2 := eFactory.Create(eClass)

	parent.ESet(children, NewImmutableEList([]any{child1, child2}))
	assert.Same(t, list, parent.EGet(children))
	assert.Equal(t, []any{child1, child2}, list.ToArray())
	assert.Same(t, parent, child1.EContainer())
	assert.Same(t, parent, child2.EContainer())

	parent.ESet(children, NewImmutableEList([]any{child2}))
	assert.Equal(t, []any{child2}, list.ToArray())
	assert.Nil(t, child1.EContainer())

	parent.ESet(refs, NewImmutableEList([]any{child1}))
	parent.EGet(refs).(EList).Add(child2)
	assert.Equal(t, []any{child1, child2}, parent.EGet(refs).(EList).ToArray())
}

// the copy of a dynamic object contains the copies of its children, its references refer to them
func TestDynamicEObject_Copy(t *testing.T) {
	eClass, children, refs := newDynamicManyClasses()
	eFactory := eClass.GetEPackage().GetEFactoryInstance()
	parent := eFactory.Create(eClass)
	child := eFactory.Create(eClass)
	parent.EGet(children).(EList).Add(child)
	parent.EGet(refs).(EList).Add(child)

	eCopy := Copy(parent)
	copyChildren := eCopy.EGet(children).(EList)
	require.Equal(t, 1, copyChildren.Size())
	copyChild := copyChildren.Get(0).(EObject)
	assert.NotSame(t, child, copyChild)
	assert.Same(t, eCopy, copyChild.EContainer())
	assert.Equal(t, []any{copyChild}, eCopy.EGet(refs).(EList).ToArray())
	assert.True(t, Equals(parent, eCopy))
}

// EGetResolve of a many reference without resolve does not resolve its proxies, as the generated
// getters
func TestDynamicEObject_GetManyWithoutResolve(t *testing.T) {
	eClass, _, refs := newDynamicManyClasses()
	eFactory := eClass.GetEPackage().GetEFactoryInstance()
	target := eFactory.Create(eClass)
	rs := NewEResourceSetImpl()
	eResource := rs.CreateResource(NewURI("target.xml"))
	eResource.GetContents().Add(target)

	eObject := eFactory.Create(eClass)
	proxy := eFactory.Create(eClass).(EObjectInternal)
	proxy.ESetProxyURI(NewURI("target.xml#/0"))
	eObject.EGet(refs).(EList).Add(proxy)
	rs.CreateResource(NewURI("source.xml")).GetContents().Add(eObject)

	unresolved := eObject.EGetResolve(refs, false).(EList)
	assert.Same(t, proxy, unresolved.Get(0))
	assert.Equal(t, 0, unresolved.IndexOf(proxy))
	assert.Same(t, target, eObject.EGet(refs).(EList).Get(0))
}
