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

func newNotifyingListTestFn(t *testing.T, factory func() *BasicENotifyingList) *eNotifyingListTest {
	l := new(eNotifyingListTest)
	l.BasicENotifyingList = factory()
	l.mockNotifier = NewMockENotifier(t)
	l.mockFeature = NewMockEStructuralFeature(t)
	l.mockAdapter = NewMockEAdapter(t)
	l.interfaces = l
	l.mockNotifier.EXPECT().EDeliver().Return(true).Maybe()
	l.mockNotifier.EXPECT().EAdapters().Return(NewImmutableEList([]any{l.mockAdapter})).Maybe()
	return l
}

func newNotifyingListTest(t *testing.T) *eNotifyingListTest {
	return newNotifyingListTestFn(t, NewBasicENotifyingList)
}

func newNotifyingListTestFromData(t *testing.T, data []any) *eNotifyingListTest {
	return newNotifyingListTestFn(t, func() *BasicENotifyingList { return newBasicENotifyingListFromData(data) })
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

func TestBasicNotifyingListAdd(t *testing.T) {
	l := newNotifyingListTest(t)
	l.mockNotifier.EXPECT().ENotify(mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == l.mockNotifier &&
			n.GetFeature() == l.mockFeature &&
			n.GetNewValue() == 3 &&
			n.GetOldValue() == nil &&
			n.GetEventType() == ADD &&
			n.GetPosition() == 0
	}))
	l.Add(3)

	l.mockNotifier.EXPECT().ENotify(mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == l.mockNotifier &&
			n.GetFeature() == l.mockFeature &&
			n.GetNewValue() == 4 &&
			n.GetOldValue() == nil &&
			n.GetEventType() == ADD &&
			n.GetPosition() == 1
	}))
	l.Add(4)
	assert.Equal(t, []any{3, 4}, l.ToArray())

}

func TestBasicNotifyingListAddAll(t *testing.T) {
	l := newNotifyingListTest(t)
	l.mockNotifier.EXPECT().ENotify(mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == l.mockNotifier &&
			n.GetFeature() == l.mockFeature &&
			reflect.DeepEqual(n.GetNewValue(), []any{2, 3}) &&
			n.GetOldValue() == nil &&
			n.GetEventType() == ADD_MANY &&
			n.GetPosition() == 0
	})).Once()
	l.AddAll(NewImmutableEList([]any{2, 3}))
	assert.Equal(t, []any{2, 3}, l.ToArray())

	l.mockNotifier.EXPECT().ENotify(mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == l.mockNotifier &&
			n.GetFeature() == l.mockFeature &&
			n.GetNewValue() == 4 &&
			n.GetOldValue() == nil &&
			n.GetEventType() == ADD &&
			n.GetPosition() == 2
	})).Once()
	l.AddAll(NewImmutableEList([]any{4}))
	assert.Equal(t, []any{2, 3, 4}, l.ToArray())
}

func TestBasicNotifyingListInsert(t *testing.T) {
	l := newNotifyingListTest(t)
	l.mockNotifier.EXPECT().ENotify(mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == l.mockNotifier &&
			n.GetFeature() == l.mockFeature &&
			n.GetNewValue() == 1 &&
			n.GetOldValue() == nil &&
			n.GetEventType() == ADD &&
			n.GetPosition() == 0
	})).Once()
	l.Insert(0, 1)
	assert.Equal(t, []any{1}, l.ToArray())

	l.mockNotifier.EXPECT().ENotify(mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == l.mockNotifier &&
			n.GetFeature() == l.mockFeature &&
			n.GetNewValue() == 2 &&
			n.GetOldValue() == nil &&
			n.GetEventType() == ADD &&
			n.GetPosition() == 0
	})).Once()
	l.Insert(0, 2)
	assert.Equal(t, []any{2, 1}, l.ToArray())

	l.mockNotifier.EXPECT().ENotify(mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == l.mockNotifier &&
			n.GetFeature() == l.mockFeature &&
			n.GetNewValue() == 3 &&
			n.GetOldValue() == nil &&
			n.GetEventType() == ADD &&
			n.GetPosition() == 1
	})).Once()
	l.Insert(1, 3)
	assert.Equal(t, []any{2, 3, 1}, l.ToArray())
}

func TestBasicNotifyingListInsertAll(t *testing.T) {
	l := newNotifyingListTest(t)

	assert.False(t, l.doInsertAll(0, NewImmutableEList([]any{})))

	l.mockNotifier.EXPECT().ENotify(mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == l.mockNotifier &&
			n.GetFeature() == l.mockFeature &&
			reflect.DeepEqual(n.GetNewValue(), []any{1, 2, 3}) &&
			n.GetOldValue() == nil &&
			n.GetEventType() == ADD_MANY &&
			n.GetPosition() == 0
	})).Once()
	assert.True(t, l.InsertAll(0, NewImmutableEList([]any{1, 2, 3})))
	assert.Equal(t, []any{1, 2, 3}, l.ToArray())

	l.mockNotifier.EXPECT().ENotify(mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == l.mockNotifier &&
			n.GetFeature() == l.mockFeature &&
			reflect.DeepEqual(n.GetNewValue(), []any{4, 5}) &&
			n.GetOldValue() == nil &&
			n.GetEventType() == ADD_MANY &&
			n.GetPosition() == 1
	})).Once()
	assert.True(t, l.InsertAll(1, NewImmutableEList([]any{4, 5})))
	assert.Equal(t, []any{1, 4, 5, 2, 3}, l.ToArray())
}

func TestBasicNotifyingListSet(t *testing.T) {
	l := newNotifyingListTestFromData(t, []any{1, 2})
	l.mockNotifier.EXPECT().ENotify(mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == l.mockNotifier &&
			n.GetFeature() == l.mockFeature &&
			n.GetNewValue() == 3 &&
			n.GetOldValue() == 2 &&
			n.GetEventType() == SET &&
			n.GetPosition() == 1
	})).Once()
	l.Set(1, 3)
	assert.Equal(t, []any{1, 3}, l.ToArray())
}

func TestBasicNotifyingListRemoveAt(t *testing.T) {
	l := newNotifyingListTestFromData(t, []any{1, 2, 3})
	l.mockNotifier.EXPECT().ENotify(mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == l.mockNotifier &&
			n.GetFeature() == l.mockFeature &&
			n.GetNewValue() == nil &&
			n.GetOldValue() == 2 &&
			n.GetEventType() == REMOVE &&
			n.GetPosition() == 1
	})).Once()
	l.RemoveAt(1)
	assert.Equal(t, []any{1, 3}, l.ToArray())
}

func TestBasicNotifyingListAddWithNotification(t *testing.T) {
	l := newNotifyingListTest(t)

	// no notifications
	l.AddWithNotification(1, nil)

	// with notifications
	mockChain := NewMockENotificationChain(t)
	mockChain.EXPECT().Add(mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == l.mockNotifier &&
			n.GetFeature() == l.mockFeature &&
			n.GetNewValue() == 2 &&
			n.GetOldValue() == nil &&
			n.GetEventType() == ADD &&
			n.GetPosition() == 1
	})).Once().Return(true)
	l.AddWithNotification(2, mockChain)

}

func TestBasicNotifyingListRemoveWithNotification(t *testing.T) {
	l := newNotifyingListTestFromData(t, []any{1})
	mockChain := NewMockENotificationChain(t)
	mockChain.EXPECT().Add(mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == l.mockNotifier &&
			n.GetFeature() == l.mockFeature &&
			n.GetNewValue() == nil &&
			n.GetOldValue() == 1 &&
			n.GetEventType() == REMOVE &&
			n.GetPosition() == 0
	})).Once().Return(true)
	l.RemoveWithNotification(1, mockChain)
}

func TestBasicNotifyingListSetWithNotification(t *testing.T) {
	l := newNotifyingListTestFromData(t, []any{1})
	mockChain := NewMockENotificationChain(t)
	mockChain.EXPECT().Add(mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == l.mockNotifier &&
			n.GetFeature() == l.mockFeature &&
			n.GetNewValue() == 2 &&
			n.GetOldValue() == 1 &&
			n.GetEventType() == SET &&
			n.GetPosition() == 0
	})).Once().Return(true)
	l.SetWithNotification(0, 2, mockChain)
}

func TestBasicNotifyingListClear(t *testing.T) {
	{
		l := newNotifyingListTestFromData(t, []any{})
		l.mockNotifier.EXPECT().ENotify(mock.MatchedBy(func(n ENotification) bool {
			return n.GetNotifier() == l.mockNotifier &&
				n.GetFeature() == l.mockFeature &&
				n.GetNewValue() == nil &&
				n.GetEventType() == REMOVE_MANY &&
				n.GetPosition() == NO_INDEX
		})).Once()
		l.Clear()
	}
	{
		l := newNotifyingListTestFromData(t, []any{1})
		l.mockNotifier.EXPECT().ENotify(mock.MatchedBy(func(n ENotification) bool {
			return n.GetNotifier() == l.mockNotifier &&
				n.GetFeature() == l.mockFeature &&
				n.GetNewValue() == nil &&
				n.GetOldValue() == 1 &&
				n.GetEventType() == REMOVE &&
				n.GetPosition() == 0
		})).Once()
		l.Clear()
	}
	{
		l := newNotifyingListTestFromData(t, []any{1, 2})
		l.mockNotifier.EXPECT().ENotify(mock.MatchedBy(func(n ENotification) bool {
			return n.GetNotifier() == l.mockNotifier &&
				n.GetFeature() == l.mockFeature &&
				n.GetNewValue() == nil &&
				n.GetEventType() == REMOVE_MANY &&
				n.GetPosition() == NO_INDEX
		})).Once()
		l.Clear()
	}

}

func TestBasicNotifyingRemoveAllClear(t *testing.T) {
	{
		l := newNotifyingListTestFromData(t, []any{})
		l.RemoveAll(NewImmutableEList([]any{}))
	}
	{
		l := newNotifyingListTestFromData(t, []any{1})
		l.mockNotifier.EXPECT().ENotify(mock.MatchedBy(func(n ENotification) bool {
			return n.GetNotifier() == l.mockNotifier &&
				n.GetFeature() == l.mockFeature &&
				n.GetOldValue() == 1 &&
				n.GetNewValue() == nil &&
				n.GetEventType() == REMOVE &&
				n.GetPosition() == 0
		})).Once()
		l.RemoveAll(NewImmutableEList([]any{1}))
	}
	{
		l := newNotifyingListTestFromData(t, []any{1, 2, 3})
		l.mockNotifier.EXPECT().ENotify(mock.MatchedBy(func(n ENotification) bool {
			return n.GetNotifier() == l.mockNotifier &&
				n.GetFeature() == l.mockFeature &&
				reflect.DeepEqual(n.GetOldValue(), []any{2, 3}) &&
				reflect.DeepEqual(n.GetNewValue(), []any{1, 2}) &&
				n.GetEventType() == REMOVE_MANY &&
				n.GetPosition() == 1
		})).Once()
		l.RemoveAll(NewImmutableEList([]any{3, 2}))
		assert.Equal(t, []any{1}, l.ToArray())
	}
}

func TestBasicNotifyingRemoveRange(t *testing.T) {
	{
		l := newNotifyingListTestFromData(t, []any{1, 2, 3})
		l.mockNotifier.EXPECT().ENotify(mock.MatchedBy(func(n ENotification) bool {
			return n.GetNotifier() == l.mockNotifier &&
				n.GetFeature() == l.mockFeature &&
				n.GetOldValue() == 1 &&
				n.GetNewValue() == nil &&
				n.GetEventType() == REMOVE &&
				n.GetPosition() == 0
		})).Once()
		l.RemoveRange(0, 1)
		assert.Equal(t, []any{2, 3}, l.ToArray())
	}
	{
		l := newNotifyingListTestFromData(t, []any{1, 2, 3})
		l.mockNotifier.EXPECT().ENotify(mock.MatchedBy(func(n ENotification) bool {
			return n.GetNotifier() == l.mockNotifier &&
				n.GetFeature() == l.mockFeature &&
				reflect.DeepEqual(n.GetOldValue(), []any{1, 2}) &&
				reflect.DeepEqual(n.GetNewValue(), []any{0, 1}) &&
				n.GetEventType() == REMOVE_MANY &&
				n.GetPosition() == 0
		})).Once()
		l.RemoveRange(0, 2)
		assert.Equal(t, []any{3}, l.ToArray())
	}
}
