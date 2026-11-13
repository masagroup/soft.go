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
	"iter"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// hybridTestStore is a store of values in maps, counting the reads of single values.
type hybridTestStore struct {
	EStore
	mu     sync.Mutex
	values map[EObject]map[EStructuralFeature]any
	gets   atomic.Int64
}

func newHybridTestStore() *hybridTestStore {
	return &hybridTestStore{values: map[EObject]map[EStructuralFeature]any{}}
}

func (s *hybridTestStore) of(object EObject) map[EStructuralFeature]any {
	m := s.values[object]
	if m == nil {
		m = map[EStructuralFeature]any{}
		s.values[object] = m
	}
	return m
}

func (s *hybridTestStore) list(object EObject, feature EStructuralFeature) []any {
	l, _ := s.of(object)[feature].([]any)
	return l
}

func (s *hybridTestStore) Get(object EObject, feature EStructuralFeature, index int) any {
	s.mu.Lock()
	defer s.mu.Unlock()
	if index == NO_INDEX {
		s.gets.Add(1)
		return s.of(object)[feature]
	}
	return s.list(object, feature)[index]
}

func (s *hybridTestStore) Set(object EObject, feature EStructuralFeature, index int, value any, _ bool) any {
	s.mu.Lock()
	defer s.mu.Unlock()
	if index == NO_INDEX {
		old := s.of(object)[feature]
		s.of(object)[feature] = value
		return old
	}
	l := s.list(object, feature)
	old := l[index]
	l[index] = value
	return old
}

func (s *hybridTestStore) IsSet(object EObject, feature EStructuralFeature) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.of(object)[feature]
	return ok
}

func (s *hybridTestStore) UnSet(object EObject, feature EStructuralFeature) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.of(object), feature)
}

func (s *hybridTestStore) Size(object EObject, feature EStructuralFeature) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.list(object, feature))
}

func (s *hybridTestStore) IsEmpty(object EObject, feature EStructuralFeature) bool {
	return s.Size(object, feature) == 0
}

func (s *hybridTestStore) Add(object EObject, feature EStructuralFeature, index int, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.of(object)[feature] = slices.Insert(s.list(object, feature), index, value)
}

func (s *hybridTestStore) IndexOf(object EObject, feature EStructuralFeature, value any) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return slices.Index(s.list(object, feature), value)
}

func (s *hybridTestStore) Contains(object EObject, feature EStructuralFeature, value any) bool {
	return s.IndexOf(object, feature, value) != -1
}

func (s *hybridTestStore) ToArray(object EObject, feature EStructuralFeature) []any {
	s.mu.Lock()
	defer s.mu.Unlock()
	return slices.Clone(s.list(object, feature))
}

func (s *hybridTestStore) All(object EObject, feature EStructuralFeature) iter.Seq[any] {
	return slices.Values(s.ToArray(object, feature))
}

func (s *hybridTestStore) GetContainer(EObject) (EObject, EStructuralFeature) { return nil, nil }

func (s *hybridTestStore) SetContainer(EObject, EObject, EStructuralFeature) {}

// hybridTestClass is a class of a name, a count and a list of items, with its features.
var hybridTestClass = sync.OnceValue(func() EClass {
	c := GetFactory().CreateEClass()
	c.SetName("HybridTest")
	name := GetFactory().CreateEAttribute()
	name.SetName("name")
	name.SetEType(GetPackage().GetEString())
	count := GetFactory().CreateEAttribute()
	count.SetName("count")
	count.SetEType(GetPackage().GetEInt())
	items := GetFactory().CreateEAttribute()
	items.SetName("items")
	items.SetUpperBound(-1)
	items.SetEType(GetPackage().GetEInt())
	ref := GetFactory().CreateEReference()
	ref.SetName("ref")
	ref.SetEType(GetPackage().GetEObject())
	c.GetEStructuralFeatures().AddAll(NewImmutableEList([]any{name, count, items, ref}))
	return c
})

const (
	hybridTestName = iota
	hybridTestCount
	hybridTestItems
	hybridTestRef
)

// hybridTestObject is a class as soft.gen generates it with the property hybridFields.
type hybridTestObject struct {
	HybridEObjectImpl
	fields Fields[hybridTestFields]
}

type hybridTestFields struct {
	FieldsHeader
	name  string
	count int
	items EList
	ref   EObject
}

var (
	hybridTestLevel = NewLevel[hybridTestFields](hybridTestClass,
		unsafe.Offsetof(hybridTestObject{}.fields)-unsafe.Offsetof(hybridTestObject{}.HybridEObjectImpl))
	hybridTestNameField  = Value[string](hybridTestLevel, hybridTestName, 0, unsafe.Offsetof(hybridTestFields{}.name), Plain)
	hybridTestCountField = Value[int](hybridTestLevel, hybridTestCount, 1, unsafe.Offsetof(hybridTestFields{}.count), Plain)
	hybridTestItemsField = List(hybridTestLevel, hybridTestItems, 2, unsafe.Offsetof(hybridTestFields{}.items))
	hybridTestRefField   = Value[EObject](hybridTestLevel, hybridTestRef, 3, unsafe.Offsetof(hybridTestFields{}.ref), Reference|Plain)
	hybridTestLayout     = NewLayout(nil, hybridTestLevel)
)

func newHybridTestObject(store EStore) *hybridTestObject {
	o := new(hybridTestObject)
	o.SetInterfaces(o)
	o.Initialize()
	o.SetEClass(hybridTestClass())
	if store != nil {
		o.SetEStore(store)
	}
	return o
}

func (o *hybridTestObject) ELayout() *Layout { return hybridTestLayout }

func (o *hybridTestObject) EGetFromID(featureID int, resolve bool) any {
	switch featureID {
	case hybridTestName:
		return o.GetName()
	case hybridTestCount:
		return o.GetCount()
	case hybridTestItems:
		return o.GetItems()
	}
	return o.HybridEObjectImpl.EGetFromID(featureID, resolve)
}

func (o *hybridTestObject) ESetFromID(featureID int, value any) {
	hybridTestLayout.ESetFromID(&o.HybridEObjectImpl, featureID, value)
}

func (o *hybridTestObject) EIsSetFromID(featureID int) bool {
	return hybridTestLayout.EIsSetFromID(&o.HybridEObjectImpl, featureID)
}

func (o *hybridTestObject) GetName() string {
	if block := o.fields.Published(1 << 0); block != nil {
		return block.name
	}
	return hybridTestNameField.Load(&o.HybridEObjectImpl, o.fields.Ref())
}

func (o *hybridTestObject) SetName(name string) {
	hybridTestNameField.Set(&o.HybridEObjectImpl, o.fields.Ref(), name)
}

func (o *hybridTestObject) GetCount() int {
	if block := o.fields.Published(1 << 1); block != nil {
		return block.count
	}
	return hybridTestCountField.Load(&o.HybridEObjectImpl, o.fields.Ref())
}

func (o *hybridTestObject) SetCount(count int) {
	hybridTestCountField.Set(&o.HybridEObjectImpl, o.fields.Ref(), count)
}

func (o *hybridTestObject) GetItems() EList {
	if block := o.fields.Published(1 << 2); block != nil {
		return block.items
	}
	return hybridTestItemsField.Load(&o.HybridEObjectImpl, o.fields.Ref())
}

func (o *hybridTestObject) GetRef() EObject {
	if block := o.fields.Published(1 << 3); block != nil && (block.ref == nil || !block.ref.EIsProxy()) {
		return block.ref
	}
	return hybridTestRefField.Load(&o.HybridEObjectImpl, o.fields.Ref())
}

func (o *hybridTestObject) SetRef(ref EObject) {
	hybridTestRefField.Set(&o.HybridEObjectImpl, o.fields.Ref(), ref)
}

// TestHybridEObjectImpl_UnresolvedProxyIsNotKept checks that a reference read while its proxy does
// not resolve is resolved again by the next read: the proxy is not published in its field.
func TestHybridEObjectImpl_UnresolvedProxyIsNotKept(t *testing.T) {
	for _, withStore := range []bool{false, true} {
		rs := NewEResourceSetImpl()
		a := rs.CreateResource(NewURI("a.xml"))
		var store EStore
		if withStore {
			store = newHybridTestStore()
		}
		o := newHybridTestObject(store)
		a.GetContents().Add(o)
		b := rs.CreateResource(NewURI("b.xml"))
		b.(interface {
			BasicSetLoaded(bool, ENotificationChain) ENotificationChain
		}).BasicSetLoaded(true, nil)

		proxy := newHybridTestObject(nil)
		proxy.ESetProxyURI(NewURI("b.xml#/0"))
		o.SetRef(proxy)
		assert.Same(t, proxy, o.GetRef(), "store %v: b.xml has no object yet", withStore)
		assert.Nil(t, o.fields.Published(1<<3), "store %v: the unresolved proxy is not published", withStore)

		target := newHybridTestObject(nil)
		b.GetContents().Add(target)
		assert.Same(t, target, o.GetRef(), "store %v: resolved by the next read", withStore)
		assert.NotNil(t, o.fields.Published(1<<3), "store %v: the resolved object is published", withStore)
	}
}

// TestHybridEObjectImpl_PublishedObjectBecomesAProxy checks that an object published in a
// reference field, then made a proxy, as the XML decoder does with an element giving its proxy URI
// after the object is set, is resolved by the next reads, and replaced in its field once resolved.
func TestHybridEObjectImpl_PublishedObjectBecomesAProxy(t *testing.T) {
	for _, withStore := range []bool{false, true} {
		rs := NewEResourceSetImpl()
		a := rs.CreateResource(NewURI("a.xml"))
		var store EStore
		if withStore {
			store = newHybridTestStore()
		}
		o := newHybridTestObject(store)
		a.GetContents().Add(o)
		b := rs.CreateResource(NewURI("b.xml"))
		b.(interface {
			BasicSetLoaded(bool, ENotificationChain) ENotificationChain
		}).BasicSetLoaded(true, nil)

		object := newHybridTestObject(nil)
		o.SetRef(object)
		require.NotNil(t, o.fields.Published(1<<3), "store %v: the object set is published", withStore)
		object.ESetProxyURI(NewURI("b.xml#/0"))
		assert.Same(t, object, o.GetRef(), "store %v: b.xml has no object yet", withStore)
		assert.Same(t, object, o.EGet(hybridTestClass().GetEStructuralFeature(hybridTestRef)), "store %v: reflective", withStore)

		target := newHybridTestObject(nil)
		b.GetContents().Add(target)
		assert.Same(t, target, o.GetRef(), "store %v: resolved by the next read", withStore)
		assert.Same(t, target, o.fields.Published(1<<3).ref, "store %v: the field holds the object resolved", withStore)
		assert.Same(t, target, o.EGet(hybridTestClass().GetEStructuralFeature(hybridTestRef)), "store %v: reflective", withStore)
	}
}

func TestHybridEObjectImpl_GetterReadsTheStoreOnce(t *testing.T) {
	store := newHybridTestStore()
	o := newHybridTestObject(store)
	store.values[o] = map[EStructuralFeature]any{hybridTestClass().GetEStructuralFeature(hybridTestName): "a"}

	assert.Equal(t, "a", o.GetName())
	assert.Equal(t, "a", o.GetName())
	assert.Equal(t, int64(1), store.gets.Load(), "the second read is the field")
	assert.Equal(t, 0, o.GetCount(), "the default value of a value the store does not have")
}

func TestHybridEObjectImpl_SetPublishesTheValue(t *testing.T) {
	store := newHybridTestStore()
	o := newHybridTestObject(store)

	o.SetCount(3)
	assert.Equal(t, 3, store.values[o][hybridTestClass().GetEStructuralFeature(hybridTestCount)], "in the store")
	gets := store.gets.Load()
	assert.Equal(t, 3, o.GetCount())
	assert.Equal(t, gets, store.gets.Load(), "read from the field the set published")

	o.SetCount(3)
	o.SetCount(4)
	assert.Equal(t, 4, o.GetCount())
	assert.Equal(t, 4, store.values[o][hybridTestClass().GetEStructuralFeature(hybridTestCount)])
}

func TestHybridEObjectImpl_ClearCacheReadsTheStoreAgain(t *testing.T) {
	store := newHybridTestStore()
	o := newHybridTestObject(store)
	o.SetName("a")
	store.values[o][hybridTestClass().GetEStructuralFeature(hybridTestName)] = "changed in the store"
	assert.Equal(t, "a", o.GetName(), "the field until the cache is cleared")

	o.ClearCache()
	assert.Equal(t, "changed in the store", o.GetName())
}

func TestHybridEObjectImpl_UnsetDropsTheField(t *testing.T) {
	store := newHybridTestStore()
	o := newHybridTestObject(store)
	o.SetName("a")
	o.EUnset(hybridTestClass().GetEStructuralFeature(hybridTestName))
	assert.Equal(t, "", o.GetName())
	assert.False(t, o.EIsSet(hybridTestClass().GetEStructuralFeature(hybridTestName)))
}

func TestHybridEObjectImpl_Reflective(t *testing.T) {
	store := newHybridTestStore()
	o := newHybridTestObject(store)
	class := hybridTestClass()
	name, count, items := class.GetEStructuralFeature(hybridTestName), class.GetEStructuralFeature(hybridTestCount), class.GetEStructuralFeature(hybridTestItems)

	o.ESet(name, "b")
	o.ESet(count, 7)
	assert.Equal(t, "b", o.GetName(), "ESet goes to the typed setter")
	assert.Equal(t, "b", o.EGet(name))
	assert.Equal(t, 7, o.EGet(count))
	assert.Equal(t, "b", store.values[o][name])

	assert.False(t, o.EIsSet(items), "an empty list is not set")
	o.GetItems().Add(1)
	assert.True(t, o.EIsSet(items))
	assert.Same(t, o.GetItems(), o.EGet(items), "one list")
	assert.Equal(t, []any{1}, store.ToArray(o, items))
}

func TestHybridEObjectImpl_WithoutStore(t *testing.T) {
	o := newHybridTestObject(nil)
	assert.Equal(t, "", o.GetName())
	o.SetName("a")
	assert.Equal(t, "a", o.GetName())
	o.fields.Ref().clear()
	assert.Equal(t, "a", o.GetName(), "without store, read again from the values the object keeps")

	store := newHybridTestStore()
	o.SetEStore(store)
	assert.Equal(t, "", o.GetName(), "a new store, read again")
}

func TestHybridEObjectImpl_ConcurrentReadsAndSets(t *testing.T) {
	store := newHybridTestStore()
	o := newHybridTestObject(store)
	o.SetCount(0)
	var wg sync.WaitGroup
	for range 4 {
		wg.Go(func() {
			for range 1000 {
				_ = o.GetCount()
				_ = o.GetName()
				_ = o.EGet(hybridTestClass().GetEStructuralFeature(hybridTestCount))
			}
		})
	}
	wg.Go(func() {
		for i := range 1000 {
			o.SetCount(i)
			if i%100 == 0 {
				o.ClearCache()
			}
		}
	})
	wg.Wait()
	require.Equal(t, 999, o.GetCount())
	assert.Equal(t, 999, store.values[o][hybridTestClass().GetEStructuralFeature(hybridTestCount)])
}

// hybridRootObject is a hybrid class with fields but without the generated reflective methods:
// the reflective methods of HybridEObjectImpl find its fields through its layout.
type hybridRootObject struct {
	HybridEObjectImpl
	fields Fields[hybridTestFields]
}

var (
	hybridRootLevel = NewLevel[hybridTestFields](hybridTestClass,
		unsafe.Offsetof(hybridRootObject{}.fields)-unsafe.Offsetof(hybridRootObject{}.HybridEObjectImpl))
	_                 = Value[string](hybridRootLevel, hybridTestName, 0, unsafe.Offsetof(hybridTestFields{}.name), Plain)
	_                 = Value[int](hybridRootLevel, hybridTestCount, 1, unsafe.Offsetof(hybridTestFields{}.count), Plain)
	_                 = List(hybridRootLevel, hybridTestItems, 2, unsafe.Offsetof(hybridTestFields{}.items))
	hybridRootLayout  = NewLayout(nil, hybridRootLevel)
	hybridChildLayout = NewLayout(hybridRootLayout, NewLevel[struct{ FieldsHeader }](hybridTestClass, 0))
)

func (o *hybridRootObject) ELayout() *Layout { return hybridRootLayout }

func newHybridRootObject(store EStore) *hybridRootObject {
	o := new(hybridRootObject)
	o.SetInterfaces(o)
	o.Initialize()
	o.SetEClass(hybridTestClass())
	o.SetEStore(store)
	return o
}

func TestHybridEObjectImpl_RootReflectiveMethods(t *testing.T) {
	store := newHybridTestStore()
	o := newHybridRootObject(store)
	class := hybridTestClass()
	name, count, items := class.GetEStructuralFeature(hybridTestName), class.GetEStructuralFeature(hybridTestCount), class.GetEStructuralFeature(hybridTestItems)

	o.ESet(name, "a")
	o.ESet(count, 2)
	gets := store.gets.Load()
	assert.Equal(t, "a", o.EGet(name))
	assert.Equal(t, 2, o.EGet(count))
	assert.Equal(t, gets, store.gets.Load(), "read from the fields the sets published")
	assert.Equal(t, "a", o.fields.Published(1<<0).name)

	assert.False(t, o.EIsSet(items))
	list := o.EGet(items).(EList)
	list.Add(5)
	assert.True(t, o.EIsSet(items))
	assert.Same(t, list, o.EGet(items), "one list")
	assert.Equal(t, []any{5}, store.ToArray(o, items))

	o.ESet(count, "not an int")
	assert.Equal(t, "not an int", store.values[o][count], "a value not of the type of the field goes to the store")
	assert.Nil(t, o.fields.Published(1<<1), "and drops the field")
}

func TestHybridEObjectImpl_WithoutLayout(t *testing.T) {
	store := newHybridTestStore()
	o := new(HybridEObjectImpl)
	o.SetInterfaces(o)
	o.Initialize()
	o.SetEClass(hybridTestClass())
	o.SetEStore(store)
	name := hybridTestClass().GetEStructuralFeature(hybridTestName)

	assert.False(t, o.EIsSet(name))
	o.ESet(name, "a")
	assert.Equal(t, "a", o.EGet(name), "the reflective methods of soft.go")
	assert.True(t, o.EIsSet(name))
	o.EUnset(name)
	assert.False(t, o.EIsSet(name))
}

func TestHybridEObjectImpl_SetKeepsTheOtherFields(t *testing.T) {
	store := newHybridTestStore()
	o := newHybridTestObject(store)
	items := o.GetItems()
	_ = o.GetName()

	// the field of the name is claimed by its load: the set replaces the block, copying the others
	o.SetName("b")
	assert.Equal(t, "b", o.GetName())
	assert.Same(t, items, o.GetItems(), "the list is copied to the new block")
	assert.NotNil(t, o.fields.Published(1<<2), "still published")
}

func TestHybridEObjectImpl_ConcurrentListLoads(t *testing.T) {
	for range 100 {
		store := newHybridTestStore()
		o := newHybridTestObject(store)
		lists := make([]EList, 8)
		var wg sync.WaitGroup
		for i := range lists {
			wg.Go(func() { lists[i] = o.GetItems() })
		}
		wg.Wait()
		for _, l := range lists {
			require.Same(t, lists[0], l, "every reader gets the same list")
		}
	}
}

func TestHybridEObjectImpl_SetCacheDropsTheFields(t *testing.T) {
	store := newHybridTestStore()
	o := newHybridTestObject(store)
	o.SetName("a")
	store.values[o][hybridTestClass().GetEStructuralFeature(hybridTestName)] = "changed in the store"

	o.SetCache(!o.IsCache())
	assert.Equal(t, "changed in the store", o.GetName())
}

func TestHybridEObjectImpl_SetEStoreLockedDropsTheFields(t *testing.T) {
	o := newHybridTestObject(newHybridTestStore())
	o.SetName("a")

	other := newHybridTestStore()
	other.values[o] = map[EStructuralFeature]any{hybridTestClass().GetEStructuralFeature(hybridTestName): "other"}
	o.Lock()
	o.SetEStoreLocked(other)
	o.Unlock()
	assert.Equal(t, "other", o.GetName())
}

func TestNewLayout_Parent(t *testing.T) {
	store := newHybridTestStore()
	o := newHybridRootObject(store)
	o.ESet(hybridTestClass().GetEStructuralFeature(hybridTestName), "a")
	value, ok := hybridChildLayout.Value(o, hybridTestName)
	assert.True(t, ok, "a field of the parent level")
	assert.Equal(t, "a", value)
	_, ok = hybridChildLayout.Value(o, 10)
	assert.False(t, ok, "no field")
}

func TestNewLevel_HeaderFirst(t *testing.T) {
	type wrong struct {
		name string
		FieldsHeader
	}
	assert.Panics(t, func() { NewLevel[wrong](hybridTestClass, 0) })
}

func (s *hybridTestStore) AddAll(object EObject, feature EStructuralFeature, index int, c Collection) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.of(object)[feature] = slices.Insert(s.list(object, feature), index, c.ToArray()...)
}

func (s *hybridTestStore) Clear(object EObject, feature EStructuralFeature) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.of(object), feature)
}

type hybridTestAdapter struct {
	AbstractEAdapter
	notifications []ENotification
}

func (a *hybridTestAdapter) NotifyChanged(n ENotification) {
	a.notifications = append(a.notifications, n)
}

func TestHybridEObjectImpl_SetNotifies(t *testing.T) {
	o := newHybridTestObject(newHybridTestStore())
	a := &hybridTestAdapter{}
	o.EAdapters().Add(a)
	o.SetName("a")
	o.SetName("a")
	o.SetName("b")
	var values [][2]any
	for _, n := range a.notifications {
		values = append(values, [2]any{n.GetOldValue(), n.GetNewValue()})
	}
	assert.Equal(t, [][2]any{{"", "a"}, {"a", "a"}, {"a", "b"}}, values, "the same value too")
	assert.Equal(t, "b", o.GetName())

	o.EUnset(hybridTestClass().GetEStructuralFeature(hybridTestName))
	assert.Len(t, a.notifications, 4)
	assert.Equal(t, "", o.GetName())
}

func TestHybridEObjectImpl_DynamicSubclass(t *testing.T) {
	sub := GetFactory().CreateEClass()
	sub.SetName("Sub")
	sub.GetESuperTypes().Add(hybridTestClass())
	extra := GetFactory().CreateEAttribute()
	extra.SetName("extra")
	extra.SetEType(GetPackage().GetEString())
	sub.GetEStructuralFeatures().Add(extra)
	for _, withStore := range []bool{false, true} {
		var store EStore
		if withStore {
			store = newHybridTestStore()
		}
		o := newHybridTestObject(store)
		o.SetEClass(sub)
		o.SetName("n")
		o.ESet(extra, "e")
		assert.Equal(t, "n", o.GetName(), "store %v", withStore)
		assert.Equal(t, "n", o.EGet(sub.GetEStructuralFeatureFromName("name")), "store %v", withStore)
		assert.Equal(t, "e", o.EGet(extra), "store %v", withStore)
		assert.True(t, o.EIsSet(extra), "store %v", withStore)
		o.EUnset(extra)
		assert.False(t, o.EIsSet(extra), "store %v", withStore)
		assert.Equal(t, "n", o.GetName(), "store %v: the static field is kept", withStore)
	}
}

func TestHybridEObjectImpl_StoreRemoved(t *testing.T) {
	store := newHybridTestStore()
	o := newHybridTestObject(store)
	o.SetName("a")
	o.GetItems().Add(1)
	assert.Equal(t, "a", o.GetName())

	o.SetEStore(nil)
	assert.Equal(t, "a", o.GetName(), "kept when the store goes")
	assert.Equal(t, []any{1}, o.GetItems().ToArray())
	o.SetName("b")
	assert.Equal(t, "b", o.GetName())
	assert.Equal(t, "a", store.values[o][hybridTestClass().GetEStructuralFeature(hybridTestName)], "the store is left")
}

func TestHybridEObjectImpl_ListUnsetAndSet(t *testing.T) {
	o := newHybridTestObject(newHybridTestStore())
	items := hybridTestClass().GetEStructuralFeature(hybridTestItems)
	list := o.GetItems()
	list.Add(1)
	assert.True(t, o.EIsSet(items))
	o.EUnset(items)
	assert.False(t, o.EIsSet(items))
	assert.Same(t, list, o.GetItems(), "the same list")
	assert.Equal(t, 0, list.Size())
	o.ESet(items, NewImmutableEList([]any{4, 5}))
	assert.Equal(t, []any{4, 5}, o.GetItems().ToArray())
}

// TestHybridEObjectImpl_FieldsCacheWithoutObjectCache checks that the typed fields cache the
// values when the object does not: a value changed in the store is read after ClearCache.
func TestHybridEObjectImpl_FieldsCacheWithoutObjectCache(t *testing.T) {
	store := newHybridTestStore()
	o := newHybridTestObject(store)
	require.False(t, o.IsCache())
	o.SetName("a")
	assert.Equal(t, "a", o.GetName())
	store.values[o][hybridTestClass().GetEStructuralFeature(hybridTestName)] = "changed in the store"
	assert.Equal(t, "a", o.GetName(), "the field")
	o.ClearCache()
	assert.Equal(t, "changed in the store", o.GetName())
}

// hybridTestNodeClass is a class of a containment, child, its container reference, parent, and
// two opposite references, peer and peerOf.
var hybridTestNodeClass = sync.OnceValue(func() EClass {
	f := GetFactory()
	c := f.CreateEClass()
	c.SetName("Node")
	ref := func(name string) EReference {
		r := f.CreateEReference()
		r.SetName(name)
		r.SetEType(c)
		return r
	}
	child, parent, peer, peerOf := ref("child"), ref("parent"), ref("peer"), ref("peerOf")
	child.SetContainment(true)
	child.SetEOpposite(parent)
	parent.SetEOpposite(child)
	peer.SetEOpposite(peerOf)
	peerOf.SetEOpposite(peer)
	c.GetEStructuralFeatures().AddAll(NewImmutableEList([]any{child, parent, peer, peerOf}))
	return c
})

const (
	hybridTestNodeChild = iota
	hybridTestNodeParent
	hybridTestNodePeer
	hybridTestNodePeerOf
)

// hybridTestNode is hybridTestNodeClass as soft.gen generates it: a field per reference but the
// container one, parent, read from the container of the object.
type hybridTestNode struct {
	HybridEObjectImpl
	fields Fields[hybridTestNodeFields]
}

type hybridTestNodeFields struct {
	FieldsHeader
	child, peer, peerOf EObject
}

var (
	hybridTestNodeLevel = NewLevel[hybridTestNodeFields](hybridTestNodeClass,
		unsafe.Offsetof(hybridTestNode{}.fields)-unsafe.Offsetof(hybridTestNode{}.HybridEObjectImpl))
	hybridTestNodeFieldsByID = map[int]*ValueField[EObject]{
		hybridTestNodeChild:  Value[EObject](hybridTestNodeLevel, hybridTestNodeChild, 0, unsafe.Offsetof(hybridTestNodeFields{}.child), Reference),
		hybridTestNodePeer:   Value[EObject](hybridTestNodeLevel, hybridTestNodePeer, 1, unsafe.Offsetof(hybridTestNodeFields{}.peer), Reference),
		hybridTestNodePeerOf: Value[EObject](hybridTestNodeLevel, hybridTestNodePeerOf, 2, unsafe.Offsetof(hybridTestNodeFields{}.peerOf), Reference),
	}
	hybridTestNodeLayout = NewLayout(nil, hybridTestNodeLevel)
)

func newHybridTestNode(store EStore) *hybridTestNode {
	o := new(hybridTestNode)
	o.SetInterfaces(o)
	o.Initialize()
	o.SetEClass(hybridTestNodeClass())
	if store != nil {
		o.SetEStore(store)
	}
	return o
}

func (o *hybridTestNode) ELayout() *Layout { return hybridTestNodeLayout }

func (o *hybridTestNode) get(featureID int) EObject {
	return hybridTestNodeFieldsByID[featureID].Get(&o.HybridEObjectImpl, o.fields.Ref())
}

func (o *hybridTestNode) feature(featureID int) EStructuralFeature {
	return hybridTestNodeClass().GetEStructuralFeature(featureID)
}

func TestHybridEObjectImpl_ContainmentsAndOpposites(t *testing.T) {
	for _, withStore := range []bool{false, true} {
		var store EStore
		if withStore {
			store = newHybridTestStore()
		}
		a, b, c := newHybridTestNode(store), newHybridTestNode(store), newHybridTestNode(store)
		assert.Nil(t, a.get(hybridTestNodeChild))
		assert.Nil(t, b.get(hybridTestNodeChild))
		assert.Nil(t, b.get(hybridTestNodePeerOf))

		a.ESet(a.feature(hybridTestNodeChild), c)
		assert.Same(t, c, a.get(hybridTestNodeChild), "store %v", withStore)
		assert.Same(t, a, c.EGet(c.feature(hybridTestNodeParent)), "store %v", withStore)
		b.ESet(b.feature(hybridTestNodeChild), c)
		assert.Same(t, c, b.get(hybridTestNodeChild), "store %v: moved", withStore)
		assert.Same(t, b, c.EGet(c.feature(hybridTestNodeParent)), "store %v: moved", withStore)
		assert.Nil(t, a.get(hybridTestNodeChild), "store %v: moved away", withStore)
		assert.False(t, a.EIsSet(a.feature(hybridTestNodeChild)), "store %v", withStore)

		a.ESet(a.feature(hybridTestNodePeer), b)
		assert.Same(t, b, a.get(hybridTestNodePeer), "store %v", withStore)
		assert.Same(t, a, b.get(hybridTestNodePeerOf), "store %v: inverse", withStore)
		c.ESet(c.feature(hybridTestNodePeer), b)
		assert.Same(t, c, b.get(hybridTestNodePeerOf), "store %v: inverse replaced", withStore)
		assert.Nil(t, a.get(hybridTestNodePeer), "store %v: old end cleared", withStore)
		b.EUnset(b.feature(hybridTestNodePeerOf))
		assert.Nil(t, c.get(hybridTestNodePeer), "store %v: unset clears the other end", withStore)
	}
}
