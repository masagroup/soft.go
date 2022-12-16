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
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBasicNotifyingListAccessors(t *testing.T) {
	l := newBasicENotifyingListFromData([]any{})
	assert.Equal(t, nil, l.GetFeature())
	assert.Equal(t, -1, l.GetFeatureID())
	assert.Equal(t, nil, l.GetNotifier())
}

type eNotifyingListTest struct {
	*BasicENotifyingList
	mockNotifier *MockENotifier
	mockFeature  *MockEStructuralFeature
	mockAdapter  *MockEAdapter
}

func newNotifyingListTestFn(factory func() *BasicENotifyingList) *eNotifyingListTest {
	l := new(eNotifyingListTest)
	l.BasicENotifyingList = factory()
	l.mockNotifier = new(MockENotifier)
	l.mockFeature = new(MockEStructuralFeature)
	l.mockAdapter = new(MockEAdapter)
	l.interfaces = l
	l.mockNotifier.On("EDeliver").Return(true)
	l.mockNotifier.On("EAdapters").Return(NewImmutableEList([]any{l.mockAdapter}))
	return l
}

func newNotifyingListTest() *eNotifyingListTest {
	return newNotifyingListTestFn(NewBasicENotifyingList)
}

func newNotifyingListTestFromData(data []any) *eNotifyingListTest {
	return newNotifyingListTestFn(func() *BasicENotifyingList { return newBasicENotifyingListFromData(data) })
}

func (list *eNotifyingListTest) GetNotifier() ENotifier {
	return list.mockNotifier
}

func (list *eNotifyingListTest) GetFeature() EStructuralFeature {
	return list.mockFeature
}

func (list *eNotifyingListTest) GetFeatureID() int {
	return list.mockFeature.GetFeatureID()
}

func (list *eNotifyingListTest) assertExpectations(t *testing.T) {
	list.mockNotifier.AssertExpectations(t)
	list.mockFeature.AssertExpectations(t)
	list.mockAdapter.AssertExpectations(t)
}

func TestNotifyingListAdd(t *testing.T) {
	l := newNotifyingListTest()
	l.mockNotifier.On("ENotify", mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == l.mockNotifier &&
			n.GetFeature() == l.mockFeature &&
			n.GetNewValue() == 3 &&
			n.GetOldValue() == nil &&
			n.GetEventType() == ADD &&
			n.GetPosition() == 0
	}))
	l.Add(3)

	l.mockNotifier.On("ENotify", mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == l.mockNotifier &&
			n.GetFeature() == l.mockFeature &&
			n.GetNewValue() == 4 &&
			n.GetOldValue() == nil &&
			n.GetEventType() == ADD &&
			n.GetPosition() == 1
	}))
	l.Add(4)
	l.assertExpectations(t)
	assert.Equal(t, []any{3, 4}, l.ToArray())

}

func TestNotifyingListAddAll(t *testing.T) {
	l := newNotifyingListTest()
	l.mockNotifier.On("ENotify", mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == l.mockNotifier &&
			n.GetFeature() == l.mockFeature &&
			reflect.DeepEqual(n.GetNewValue(), []any{2, 3}) &&
			n.GetOldValue() == nil &&
			n.GetEventType() == ADD_MANY &&
			n.GetPosition() == 0
	})).Once()
	l.AddAll(NewImmutableEList([]any{2, 3}))
	l.assertExpectations(t)
	assert.Equal(t, []any{2, 3}, l.ToArray())

	l.mockNotifier.On("ENotify", mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == l.mockNotifier &&
			n.GetFeature() == l.mockFeature &&
			n.GetNewValue() == 4 &&
			n.GetOldValue() == nil &&
			n.GetEventType() == ADD &&
			n.GetPosition() == 2
	})).Once()
	l.AddAll(NewImmutableEList([]any{4}))
	l.assertExpectations(t)
	assert.Equal(t, []any{2, 3, 4}, l.ToArray())
}

func TestNotifyingListInsert(t *testing.T) {
	l := newNotifyingListTest()
	l.mockNotifier.On("ENotify", mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == l.mockNotifier &&
			n.GetFeature() == l.mockFeature &&
			n.GetNewValue() == 1 &&
			n.GetOldValue() == nil &&
			n.GetEventType() == ADD &&
			n.GetPosition() == 0
	})).Once()
	l.Insert(0, 1)
	l.assertExpectations(t)
	assert.Equal(t, []any{1}, l.ToArray())

	l.mockNotifier.On("ENotify", mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == l.mockNotifier &&
			n.GetFeature() == l.mockFeature &&
			n.GetNewValue() == 2 &&
			n.GetOldValue() == nil &&
			n.GetEventType() == ADD &&
			n.GetPosition() == 0
	})).Once()
	l.Insert(0, 2)
	l.assertExpectations(t)
	assert.Equal(t, []any{2, 1}, l.ToArray())

	l.mockNotifier.On("ENotify", mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == l.mockNotifier &&
			n.GetFeature() == l.mockFeature &&
			n.GetNewValue() == 3 &&
			n.GetOldValue() == nil &&
			n.GetEventType() == ADD &&
			n.GetPosition() == 1
	})).Once()
	l.Insert(1, 3)
	l.assertExpectations(t)
	assert.Equal(t, []any{2, 3, 1}, l.ToArray())
}

func TestNotifyingListInsertAll(t *testing.T) {
	l := newNotifyingListTest()

	assert.False(t, l.doInsertAll(0, NewImmutableEList([]any{})))

	l.mockNotifier.On("ENotify", mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == l.mockNotifier &&
			n.GetFeature() == l.mockFeature &&
			reflect.DeepEqual(n.GetNewValue(), []any{1, 2, 3}) &&
			n.GetOldValue() == nil &&
			n.GetEventType() == ADD_MANY &&
			n.GetPosition() == 0
	})).Once()
	assert.True(t, l.InsertAll(0, NewImmutableEList([]any{1, 2, 3})))
	l.assertExpectations(t)
	assert.Equal(t, []any{1, 2, 3}, l.ToArray())

	l.mockNotifier.On("ENotify", mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == l.mockNotifier &&
			n.GetFeature() == l.mockFeature &&
			reflect.DeepEqual(n.GetNewValue(), []any{4, 5}) &&
			n.GetOldValue() == nil &&
			n.GetEventType() == ADD_MANY &&
			n.GetPosition() == 1
	})).Once()
	assert.True(t, l.InsertAll(1, NewImmutableEList([]any{4, 5})))
	l.assertExpectations(t)
	assert.Equal(t, []any{1, 4, 5, 2, 3}, l.ToArray())
}

func TestNotifyingListSet(t *testing.T) {
	l := newNotifyingListTestFromData([]any{1, 2})
	l.mockNotifier.On("ENotify", mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == l.mockNotifier &&
			n.GetFeature() == l.mockFeature &&
			n.GetNewValue() == 3 &&
			n.GetOldValue() == 2 &&
			n.GetEventType() == SET &&
			n.GetPosition() == 1
	})).Once()
	l.Set(1, 3)
	l.assertExpectations(t)
	assert.Equal(t, []any{1, 3}, l.ToArray())
}

func TestNotifyingListRemoveAt(t *testing.T) {
	l := newNotifyingListTestFromData([]any{1, 2, 3})
	l.mockNotifier.On("ENotify", mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == l.mockNotifier &&
			n.GetFeature() == l.mockFeature &&
			n.GetNewValue() == nil &&
			n.GetOldValue() == 2 &&
			n.GetEventType() == REMOVE &&
			n.GetPosition() == 1
	})).Once()
	l.RemoveAt(1)
	l.assertExpectations(t)
	assert.Equal(t, []any{1, 3}, l.ToArray())
}

func TestNotifyingListAddWithNotification(t *testing.T) {
	l := newNotifyingListTest()

	// no notifications
	l.AddWithNotification(1, nil)
	l.assertExpectations(t)

	// with notifications
	mockChain := new(MockENotificationChain)
	mockChain.On("Add", mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == l.mockNotifier &&
			n.GetFeature() == l.mockFeature &&
			n.GetNewValue() == 2 &&
			n.GetOldValue() == nil &&
			n.GetEventType() == ADD &&
			n.GetPosition() == 1
	})).Once().Return(true)
	l.AddWithNotification(2, mockChain)
	l.assertExpectations(t)
	mockChain.AssertExpectations(t)

}

func TestNotifyingListRemoveWithNotification(t *testing.T) {
	l := newNotifyingListTestFromData([]any{1})
	mockChain := new(MockENotificationChain)
	mockChain.On("Add", mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == l.mockNotifier &&
			n.GetFeature() == l.mockFeature &&
			n.GetNewValue() == nil &&
			n.GetOldValue() == 1 &&
			n.GetEventType() == REMOVE &&
			n.GetPosition() == 0
	})).Once().Return(true)
	l.RemoveWithNotification(1, mockChain)
	l.assertExpectations(t)
	mockChain.AssertExpectations(t)
}

func TestNotifyingListSetWithNotification(t *testing.T) {
	l := newNotifyingListTestFromData([]any{1})
	mockChain := new(MockENotificationChain)
	mockChain.On("Add", mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == l.mockNotifier &&
			n.GetFeature() == l.mockFeature &&
			n.GetNewValue() == 2 &&
			n.GetOldValue() == 1 &&
			n.GetEventType() == SET &&
			n.GetPosition() == 0
	})).Once().Return(true)
	l.SetWithNotification(0, 2, mockChain)
	l.assertExpectations(t)
	mockChain.AssertExpectations(t)
}

func TestNotifyingListClear(t *testing.T) {
	{
		l := newNotifyingListTestFromData([]any{})
		l.mockNotifier.On("ENotify", mock.MatchedBy(func(n ENotification) bool {
			return n.GetNotifier() == l.mockNotifier &&
				n.GetFeature() == l.mockFeature &&
				n.GetNewValue() == nil &&
				n.GetEventType() == REMOVE_MANY &&
				n.GetPosition() == NO_INDEX
		})).Once()
		l.Clear()
		l.assertExpectations(t)
	}
	{
		l := newNotifyingListTestFromData([]any{1})
		l.mockNotifier.On("ENotify", mock.MatchedBy(func(n ENotification) bool {
			return n.GetNotifier() == l.mockNotifier &&
				n.GetFeature() == l.mockFeature &&
				n.GetNewValue() == nil &&
				n.GetOldValue() == 1 &&
				n.GetEventType() == REMOVE &&
				n.GetPosition() == 0
		})).Once()
		l.Clear()
		l.assertExpectations(t)
	}
	{
		l := newNotifyingListTestFromData([]any{1, 2})
		l.mockNotifier.On("ENotify", mock.MatchedBy(func(n ENotification) bool {
			return n.GetNotifier() == l.mockNotifier &&
				n.GetFeature() == l.mockFeature &&
				n.GetNewValue() == nil &&
				n.GetEventType() == REMOVE_MANY &&
				n.GetPosition() == NO_INDEX
		})).Once()
		l.Clear()
		l.assertExpectations(t)
	}

}

func TestNotifyingRemoveAllClear(t *testing.T) {
	{
		l := newNotifyingListTestFromData([]any{})
		l.RemoveAll(NewImmutableEList([]any{}))
	}
	{
		l := newNotifyingListTestFromData([]any{1})
		l.mockNotifier.On("ENotify", mock.MatchedBy(func(n ENotification) bool {
			return n.GetNotifier() == l.mockNotifier &&
				n.GetFeature() == l.mockFeature &&
				n.GetOldValue() == 1 &&
				n.GetNewValue() == nil &&
				n.GetEventType() == REMOVE &&
				n.GetPosition() == 0
		})).Once()
		l.RemoveAll(NewImmutableEList([]any{1}))
		l.assertExpectations(t)
	}
	{
		l := newNotifyingListTestFromData([]any{1, 2, 3})
		l.mockNotifier.On("ENotify", mock.MatchedBy(func(n ENotification) bool {
			return n.GetNotifier() == l.mockNotifier &&
				n.GetFeature() == l.mockFeature &&
				reflect.DeepEqual(n.GetOldValue(), []any{2, 3}) &&
				reflect.DeepEqual(n.GetNewValue(), []any{1, 2}) &&
				n.GetEventType() == REMOVE_MANY &&
				n.GetPosition() == 1
		})).Once()
		l.RemoveAll(NewImmutableEList([]any{3, 2}))
		l.assertExpectations(t)
		assert.Equal(t, []any{1}, l.ToArray())
	}
}

func TestNotifyingRemoveRange(t *testing.T) {
	{
		l := newNotifyingListTestFromData([]any{1, 2, 3})
		l.mockNotifier.On("ENotify", mock.MatchedBy(func(n ENotification) bool {
			return n.GetNotifier() == l.mockNotifier &&
				n.GetFeature() == l.mockFeature &&
				n.GetOldValue() == 1 &&
				n.GetNewValue() == nil &&
				n.GetEventType() == REMOVE &&
				n.GetPosition() == 0
		})).Once()
		l.RemoveRange(0, 1)
		l.assertExpectations(t)
		assert.Equal(t, []any{2, 3}, l.ToArray())
	}
	{
		l := newNotifyingListTestFromData([]any{1, 2, 3})
		l.mockNotifier.On("ENotify", mock.MatchedBy(func(n ENotification) bool {
			return n.GetNotifier() == l.mockNotifier &&
				n.GetFeature() == l.mockFeature &&
				reflect.DeepEqual(n.GetOldValue(), []any{1, 2}) &&
				reflect.DeepEqual(n.GetNewValue(), []any{0, 1}) &&
				n.GetEventType() == REMOVE_MANY &&
				n.GetPosition() == 0
		})).Once()
		l.RemoveRange(0, 2)
		l.assertExpectations(t)
		assert.Equal(t, []any{3}, l.ToArray())
	}
}
