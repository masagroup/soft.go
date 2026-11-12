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
	c.GetEStructuralFeatures().AddAll(NewImmutableEList([]any{name, count, items}))
	return c
})

const (
	hybridTestName = iota
	hybridTestCount
	hybridTestItems
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
}

var (
	hybridTestLevel = NewLevel[hybridTestFields](hybridTestClass,
		unsafe.Offsetof(hybridTestObject{}.fields)-unsafe.Offsetof(hybridTestObject{}.HybridEObjectImpl))
	hybridTestNameField  = Value[string](hybridTestLevel, hybridTestName, 0, unsafe.Offsetof(hybridTestFields{}.name), Plain)
	hybridTestCountField = Value[int](hybridTestLevel, hybridTestCount, 1, unsafe.Offsetof(hybridTestFields{}.count), Plain)
	hybridTestItemsField = List(hybridTestLevel, hybridTestItems, 2, unsafe.Offsetof(hybridTestFields{}.items))
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
