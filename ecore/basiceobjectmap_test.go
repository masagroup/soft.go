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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBasicEObjectMap_Constructor(t *testing.T) {
	mockClass := NewMockEClass(t)
	mockEObject := NewMockEObjectInternal(t)
	m := NewBasicEObjectMap(mockClass, mockEObject, 1, -1, false)
	assert.NotNil(t, m)

	var mp EMap = m
	assert.NotNil(t, mp)

	var ml EList = m
	assert.NotNil(t, ml)
}

type MockEObjectEMapEntry struct {
	MockEObjectEMapEntry_Prototype
	mock.Mock
}

type MockEObjectEMapEntry_Prototype struct {
	MockEObject_Prototype
	MockEMapEntry_Prototype
}

type MockEObjectEMapEntry_Expecter struct {
	MockEObject_Expecter
	MockEMapEntry_Expecter
}

func (eMapEntry *MockEObjectEMapEntry_Prototype) EXPECT() *MockEObjectEMapEntry_Expecter {
	e := &MockEObjectEMapEntry_Expecter{}
	e.MockEObject_Expecter.Mock = eMapEntry.Mock
	e.MockEMapEntry_Expecter.Mock = eMapEntry.Mock
	return e
}

type mockConstructorTestingTNewMockEObjectEMapEntry interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockENotifier creates a new instance of MockENotifier_Prototype. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockEObjectEMapEntry(t mockConstructorTestingTNewMockEObjectEMapEntry) *MockEObjectEMapEntry {
	mock := &MockEObjectEMapEntry{}
	mock.MockEObject_Prototype.Mock = &mock.Mock
	mock.MockEMapEntry_Prototype.Mock = &mock.Mock
	mock.Mock.Test(t)
	t.Cleanup(func() { mock.AssertExpectations(t) })
	return mock
}

func TestBasicEObjectMap_Add(t *testing.T) {
	mockClass := NewMockEClass(t)
	mockOwner := NewMockEObjectInternal(t)
	mockEntry := NewMockEObjectEMapEntry(t)
	m := NewBasicEObjectMap(mockClass, mockOwner, 1, -1, false)
	mockOwner.EXPECT().EDeliver().Return(false).Once()
	m.Add(mockEntry)
	mock.AssertExpectationsForObjects(t, mockClass, mockEntry, mockOwner)

	mockEntry.EXPECT().GetKey().Return(2).Once()
	mockEntry.EXPECT().GetValue().Return("2").Once()
	assert.Equal(t, "2", m.GetValue(2))
	mock.AssertExpectationsForObjects(t, mockClass, mockEntry, mockOwner)
}

func TestBasicEObjectMap_Put(t *testing.T) {
	mockClass := NewMockEClass(t)
	mockPackage := NewMockEPackage(t)
	mockFactory := NewMockEFactory(t)
	mockOwner := NewMockEObjectInternal(t)
	mockEntry := NewMockEObjectEMapEntry(t)
	m := NewBasicEObjectMap(mockClass, mockOwner, 1, -1, false)

	// put
	mockClass.EXPECT().GetEPackage().Once().Return(mockPackage)
	mockPackage.EXPECT().GetEFactoryInstance().Once().Return(mockFactory)
	mockFactory.EXPECT().Create(mockClass).Once().Return(mockEntry)
	mockEntry.EXPECT().SetKey(2).Once()
	mockEntry.EXPECT().SetValue("2").Once()
	mockOwner.EXPECT().EDeliver().Return(false).Once()
	m.Put(2, "2")
	mock.AssertExpectationsForObjects(t, mockClass, mockPackage, mockFactory, mockEntry, mockOwner)

	// check
	mockEntry.EXPECT().GetKey().Once().Return(2)
	mockEntry.EXPECT().GetValue().Once().Return("2")
	assert.Equal(t, "2", m.GetValue(2))
	mock.AssertExpectationsForObjects(t, mockClass, mockPackage, mockFactory, mockEntry, mockOwner)
}

func TestBasicEObjectMap_Put_WithNotification(t *testing.T) {
	mockClass := NewMockEClass(t)
	mockPackage := NewMockEPackage(t)
	mockFactory := NewMockEFactory(t)
	mockOwner := NewMockEObjectInternal(t)
	mockEntry := NewMockEObjectEMapEntry(t)
	mockAdapter := NewMockEAdapter(t)
	m := NewBasicEObjectMap(mockClass, mockOwner, 1, -1, false)

	mockClass.EXPECT().GetEPackage().Once().Return(mockPackage)
	mockPackage.EXPECT().GetEFactoryInstance().Once().Return(mockFactory)
	mockFactory.EXPECT().Create(mockClass).Once().Return(mockEntry)
	mockEntry.EXPECT().SetKey(2).Once()
	mockEntry.EXPECT().SetValue("2").Once()
	mockOwner.EXPECT().EDeliver().Return(true).Once()
	mockOwner.EXPECT().EAdapters().Return(NewImmutableEList([]any{mockAdapter})).Once()
	mockOwner.EXPECT().ENotify(mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == mockOwner && n.GetFeatureID() == 1 && n.GetEventType() == ADD && n.GetNewValue() == mockEntry
	})).Once()
	m.Put(2, "2")
	mock.AssertExpectationsForObjects(t, mockClass, mockPackage, mockFactory, mockEntry, mockOwner)
}

func TestBasicEObjectMap_DidSet(t *testing.T) {
	mockClass := NewMockEClass(t)
	mockOwner := NewMockEObjectInternal(t)
	mockEntry1 := NewMockEObjectEMapEntry(t)
	mockEntry2 := NewMockEObjectEMapEntry(t)
	m := NewBasicEObjectMap(mockClass, mockOwner, 1, -1, false)

	mockOwner.EXPECT().EDeliver().Once().Return(false)
	m.Add(mockEntry1)

	mockOwner.EXPECT().EDeliver().Once().Return(false)
	mockEntry1.EXPECT().GetKey().Once().Return("key1")
	m.Set(0, mockEntry2)
	mock.AssertExpectationsForObjects(t, mockClass, mockOwner, mockEntry1, mockEntry2)
}

func TestBasicEObjectMap_DidRemove(t *testing.T) {
	mockClass := NewMockEClass(t)
	mockOwner := NewMockEObjectInternal(t)
	mockEntry1 := NewMockEObjectEMapEntry(t)
	m := NewBasicEObjectMap(mockClass, mockOwner, 1, -1, false)

	mockOwner.EXPECT().EDeliver().Once().Return(false)
	m.Add(mockEntry1)

	mockOwner.EXPECT().EDeliver().Once().Return(false)
	mockEntry1.EXPECT().GetKey().Once().Return("key1")
	m.Remove(mockEntry1)
	mock.AssertExpectationsForObjects(t, mockClass, mockOwner, mockEntry1)
}
