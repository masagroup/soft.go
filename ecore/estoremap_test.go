package ecore

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestEStoreMap_Constructor(t *testing.T) {
	mockClass := NewMockEClass(t)
	mockOwner := NewMockEObjectInternal(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockStore := NewMockEStore(t)
	mockFeature.EXPECT().IsUnique().Return(true).Once()
	mockStore.EXPECT().Size(mockOwner, mockFeature).Return(0).Once()
	m := NewEStoreMap(mockClass, mockOwner, mockFeature, mockStore)
	require.NotNil(t, m)
	require.Equal(t, mockStore, m.GetEStore())
	require.Equal(t, false, m.IsCache())
}

func TestEStoreMap_SetEStore(t *testing.T) {
	mockClass := NewMockEClass(t)
	mockOwner := NewMockEObjectInternal(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockStore := NewMockEStore(t)
	mockFeature.EXPECT().IsUnique().Return(true).Once()
	mockStore.EXPECT().Size(mockOwner, mockFeature).Return(0).Once()
	m := NewEStoreMap(mockClass, mockOwner, mockFeature, mockStore)
	require.NotNil(t, m)
	require.Equal(t, mockStore, m.GetEStore())
	require.Equal(t, false, m.IsCache())

	mockStore.EXPECT().ToArray(mockOwner, mockFeature).Return([]any{}).Once()
	m.SetEStore(nil)
	require.Equal(t, nil, m.GetEStore())
}

type MockEObjectEMapEntryWithCache struct {
	mock.Mock
	MockEObjectEMapEntryWithCache_Prototype
}

type MockEObjectEMapEntryWithCache_Prototype struct {
	mock *mock.Mock
	MockEObject_Prototype
	MockEMapEntry_Prototype
	MockECacheProvider_Prototype
}

func (_mp *MockEObjectEMapEntryWithCache_Prototype) SetMock(mock *mock.Mock) {
	_mp.mock = mock
	_mp.MockEObject_Prototype.SetMock(mock)
	_mp.MockEMapEntry_Prototype.SetMock(mock)
	_mp.MockECacheProvider_Prototype.SetMock(mock)
}

type MockEObjectEMapEntryWithCache_Expecter struct {
	MockEObject_Expecter
	MockEMapEntry_Expecter
	MockECacheProvider_Expecter
}

func (_me *MockEObjectEMapEntryWithCache_Expecter) SetMock(mock *mock.Mock) {
	_me.MockEObject_Expecter.SetMock(mock)
	_me.MockEMapEntry_Expecter.SetMock(mock)
	_me.MockECacheProvider_Expecter.SetMock(mock)
}

func (eMapEntry *MockEObjectEMapEntryWithCache_Prototype) EXPECT() *MockEObjectEMapEntryWithCache_Expecter {
	e := &MockEObjectEMapEntryWithCache_Expecter{}
	e.SetMock(eMapEntry.mock)
	return e
}

type mockConstructorTestingTNewMockEObjectEMapEntryWithCache interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockENotifier creates a new instance of MockENotifier_Prototype. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockEObjectEMapEntryWithCache(t mockConstructorTestingTNewMockEObjectEMapEntryWithCache) *MockEObjectEMapEntryWithCache {
	mock := &MockEObjectEMapEntryWithCache{}
	mock.SetMock(&mock.Mock)
	mock.Mock.Test(t)
	t.Cleanup(func() { mock.AssertExpectations(t) })
	return mock
}

func TestEStoreMap_Put(t *testing.T) {
	mockClass := NewMockEClass(t)
	mockOwner := NewMockEObjectInternal(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockStore := NewMockEStore(t)
	mockFeature.EXPECT().IsUnique().Return(true).Once()
	mockStore.EXPECT().Size(mockOwner, mockFeature).Return(0).Once()
	m := NewEStoreMap(mockClass, mockOwner, mockFeature, mockStore)
	require.NotNil(t, m)

	mockEntry := NewMockEObjectEMapEntryWithCache(t)
	mockPackage := NewMockEPackage(t)
	mockFactory := NewMockEFactory(t)
	mockClass.EXPECT().GetEPackage().Return(mockPackage).Once()
	mockPackage.EXPECT().GetEFactoryInstance().Return(mockFactory).Once()
	mockFactory.EXPECT().Create(mockClass).Return(mockEntry).Once()
	mockStore.EXPECT().All(mockOwner, mockFeature).Return(slices.Values([]any{})).Once()
	mockEntry.EXPECT().SetCache(false).Once()
	mockEntry.EXPECT().SetKey(1).Once()
	mockEntry.EXPECT().SetValue(2).Once()
	mockStore.EXPECT().Contains(mockOwner, mockFeature, mockEntry).Return(false).Once()
	mockStore.EXPECT().Add(mockOwner, mockFeature, 0, mockEntry).Return().Once()
	mockOwner.EXPECT().EDeliver().Return(false).Once()
	m.Put(1, 2)
}

func TestEStoreMap_Add(t *testing.T) {
	mockClass := NewMockEClass(t)
	mockOwner := NewMockEObjectInternal(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockStore := NewMockEStore(t)
	mockFeature.EXPECT().IsUnique().Return(true).Once()
	mockStore.EXPECT().Size(mockOwner, mockFeature).Return(0).Once()
	m := NewEStoreMap(mockClass, mockOwner, mockFeature, mockStore)
	require.NotNil(t, m)

	mockEntry := NewMockEObjectEMapEntryWithCache(t)
	mockStore.EXPECT().Contains(mockOwner, mockFeature, mockEntry).Return(false).Once()
	mockStore.EXPECT().Add(mockOwner, mockFeature, 0, mockEntry).Once()
	mockOwner.EXPECT().EDeliver().Return(false).Once()
	m.Add(mockEntry)
}

func TestEStoreMap_Remove(t *testing.T) {
	mockClass := NewMockEClass(t)
	mockOwner := NewMockEObjectInternal(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockStore := NewMockEStore(t)
	mockFeature.EXPECT().IsUnique().Return(true).Once()
	mockStore.EXPECT().Size(mockOwner, mockFeature).Return(1).Once()
	m := NewEStoreMap(mockClass, mockOwner, mockFeature, mockStore)
	require.NotNil(t, m)

	mockEntry := NewMockEObjectEMapEntryWithCache(t)
	mockEntry.EXPECT().GetKey().Return("key").Once()
	mockStore.EXPECT().IndexOf(mockOwner, mockFeature, mockEntry).Return(0).Once()
	mockStore.EXPECT().Remove(mockOwner, mockFeature, 0, true).Return(mockEntry).Once()
	mockOwner.EXPECT().EDeliver().Return(false).Once()
	require.Equal(t, true, m.Remove(mockEntry))
}

func TestEStoreMap_Set(t *testing.T) {
	mockClass := NewMockEClass(t)
	mockOwner := NewMockEObjectInternal(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockStore := NewMockEStore(t)
	mockFeature.EXPECT().IsUnique().Return(true).Once()
	mockStore.EXPECT().Size(mockOwner, mockFeature).Return(1).Once()
	m := NewEStoreMap(mockClass, mockOwner, mockFeature, mockStore)
	require.NotNil(t, m)

	oldEntry := NewMockEObjectEMapEntryWithCache(t)
	newEntry := NewMockEObjectEMapEntryWithCache(t)
	mockStore.EXPECT().IndexOf(mockOwner, mockFeature, newEntry).Return(-1).Once()
	mockStore.EXPECT().Set(mockOwner, mockFeature, 0, newEntry, true).Return(oldEntry).Once()
	mockOwner.EXPECT().EDeliver().Return(false).Once()
	oldEntry.EXPECT().GetKey().Return("key").Once()
	require.Equal(t, oldEntry, m.Set(0, newEntry))
}

func TestEStoreMap_Clear(t *testing.T) {
	mockClass := NewMockEClass(t)
	mockOwner := NewMockEObjectInternal(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockStore := NewMockEStore(t)
	mockFeature.EXPECT().IsUnique().Return(true).Once()
	mockStore.EXPECT().Size(mockOwner, mockFeature).Return(1).Once()
	m := NewEStoreMap(mockClass, mockOwner, mockFeature, mockStore)
	require.NotNil(t, m)
	mockStore.EXPECT().ToArray(mockOwner, mockFeature).Return([]any{}).Once()
	mockStore.EXPECT().Clear(mockOwner, mockFeature).Once()
	mockOwner.EXPECT().EDeliver().Return(false).Once()
	m.Clear()
}
