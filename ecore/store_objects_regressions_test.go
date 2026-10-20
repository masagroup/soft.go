package ecore

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newStoreObjectsModel() (EFactory, EClass, EAttribute, EAttribute) {
	f := GetFactory()
	ePackage := f.CreateEPackage()
	ePackage.SetEFactoryInstance(f.CreateEFactory())
	eClass := f.CreateEClass()
	eClass.SetName("C")
	ePackage.GetEClassifiers().Add(eClass)
	count := f.CreateEAttribute()
	count.SetName("count")
	count.SetEType(GetPackage().GetEInt())
	eClass.GetEStructuralFeatures().Add(count)
	names := f.CreateEAttribute()
	names.SetName("names")
	names.SetEType(GetPackage().GetEString())
	names.SetUpperBound(UNBOUNDED_MULTIPLICITY)
	names.SetUnique(false)
	eClass.GetEStructuralFeatures().Add(names)
	return ePackage.GetEFactoryInstance(), eClass, count, names
}

// the store gets the elements of AddAll, not the collection of the caller, which may change
// before an asynchronous store writes them
func TestEStoreList_AddAllCopiesTheCollection(t *testing.T) {
	eFactory, eClass, _, names := newStoreObjectsModel()
	owner := eFactory.Create(eClass)
	store := NewMockEStore(t)
	store.EXPECT().Size(owner, names).Return(0).Once()
	var received Collection
	store.EXPECT().AddAll(owner, names, 0, mock.Anything).Run(func(_ EObject, _ EStructuralFeature, _ int, c Collection) { received = c }).Once()
	l := NewEStoreList(owner, names, store)
	source := NewBasicEList([]any{"a", "b"})
	l.AddAll(source)
	source.Clear()
	assert.Equal(t, []any{"a", "b"}, received.ToArray())
}

// reading a stored default value does not make it set
func TestEStoreEObjectImpl_ReadDoesNotSet(t *testing.T) {
	_, eClass, count, _ := newStoreObjectsModel()
	store := NewMockEStore(t)
	store.EXPECT().SetContainer(mock.Anything, mock.Anything, mock.Anything).Maybe()
	store.EXPECT().GetContainer(mock.Anything).Return(nil, nil).Maybe()
	store.EXPECT().Get(mock.Anything, count, NO_INDEX).Return(0).Maybe()
	store.EXPECT().IsSet(mock.Anything, count).Return(false).Maybe()
	o := NewEStoreEObjectImpl(true)
	o.SetEClass(eClass)
	o.SetEStore(store)
	assert.False(t, o.EIsSet(count))
	assert.Equal(t, 0, o.EGet(count))
	assert.False(t, o.EIsSet(count))
}

// a container read before the store is set is not written as nil in the store
func TestEStoreEObjectImpl_ContainerReadBeforeStore(t *testing.T) {
	_, eClass, _, _ := newStoreObjectsModel()
	o := NewEStoreEObjectImpl(false)
	o.SetEClass(eClass)
	assert.Nil(t, o.EInternalContainer())
	container := NewEStoreEObjectImpl(false)
	containment := GetFactory().CreateEReference()
	store := NewMockEStore(t)
	store.EXPECT().GetContainer(o).Return(container, containment).Once()
	o.SetEStore(store) // no SetContainer: the store keeps its container
	assert.Equal(t, container, o.EInternalContainer())
}

// the contents and the container of an object are read by concurrent goroutines (checked with -race)
func TestCompactObjects_ConcurrentReads(t *testing.T) {
	_, eClass, _, _ := newStoreObjectsModel()
	reflective := NewReflectiveEObjectImpl()
	reflective.SetEClass(eClass)
	store := NewEStoreEObjectImpl(true)
	store.SetEClass(eClass)
	for _, o := range []EObject{reflective, store} {
		var wg sync.WaitGroup
		for range 8 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = o.EContents().Size()
				_ = o.ECrossReferences().Size()
				_ = o.EContainer()
			}()
		}
		wg.Wait()
	}
}
