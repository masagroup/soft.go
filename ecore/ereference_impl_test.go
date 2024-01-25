// Code generated by soft.generator.go. DO NOT EDIT.

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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func discardEReference() {
	_ = assert.Equal
	_ = mock.Anything
	_ = testing.Coverage
	_ = NewMockEReference
}

func TestEReferenceAsEReference(t *testing.T) {
	o := newEReferenceImpl()
	assert.Equal(t, o, o.asEReference())
}

func TestEReferenceStaticClass(t *testing.T) {
	o := newEReferenceImpl()
	assert.Equal(t, GetPackage().GetEReference(), o.EStaticClass())
}

func TestEReferenceFeatureCount(t *testing.T) {
	o := newEReferenceImpl()
	assert.Equal(t, EREFERENCE_FEATURE_COUNT, o.EStaticFeatureCount())
}

func TestEReferenceEKeysGet(t *testing.T) {
	o := newEReferenceImpl()
	assert.NotNil(t, o.GetEKeys())
	assert.Panics(t, func() { _ = o.GetEKeys().Get(0).(EAttribute) })
}

func TestEReferenceEOppositeGet(t *testing.T) {
	o := newEReferenceImpl()

	// get default value
	assert.Nil(t, o.GetEOpposite())

	// initialize object with a mock value
	mockValue := NewMockEReference(t)
	o.eOpposite = mockValue

	// events
	mockAdapter := NewMockEAdapter(t)
	mockAdapter.EXPECT().SetTarget(o).Once()
	o.EAdapters().Add(mockAdapter)
	mock.AssertExpectationsForObjects(t, mockAdapter)

	// set object resource
	mockResourceSet := NewMockEResourceSet(t)
	mockResource := NewMockEResource(t)
	o.ESetInternalResource(mockResource)

	// get non resolved value
	mockValue.EXPECT().EIsProxy().Return(false).Once()
	assert.Equal(t, mockValue, o.GetEOpposite())
	mock.AssertExpectationsForObjects(t, mockValue, mockAdapter, mockResource, mockResourceSet)

	// get a resolved value
	mockURI := NewURI("test:///file.t")
	mockResolved := NewMockEReference(t)
	mockResolved.EXPECT().EProxyURI().Return(nil).Once()
	mockResource.EXPECT().GetResourceSet().Return(mockResourceSet).Once()
	mockResourceSet.EXPECT().GetEObject(mockURI, true).Return(mockResolved).Once()
	mockValue.EXPECT().EIsProxy().Return(true).Once()
	mockValue.EXPECT().EProxyURI().Return(mockURI).Twice()
	mockAdapter.EXPECT().NotifyChanged(mock.MatchedBy(func(notification ENotification) bool {
		return notification.GetEventType() == RESOLVE && notification.GetFeatureID() == EREFERENCE__EOPPOSITE && notification.GetOldValue() == mockValue && notification.GetNewValue() == mockResolved
	})).Once()
	assert.Equal(t, mockResolved, o.GetEOpposite())
	mock.AssertExpectationsForObjects(t, mockAdapter, mockValue, mockResolved, mockAdapter, mockResource, mockResourceSet)
}

func TestEReferenceEOppositeSet(t *testing.T) {
	o := newEReferenceImpl()
	v := NewMockEReference(t)
	mockAdapter := NewMockEAdapter(t)
	mockAdapter.EXPECT().SetTarget(o).Once()
	mockAdapter.EXPECT().NotifyChanged(mock.Anything).Once()
	o.EAdapters().Add(mockAdapter)
	o.SetEOpposite(v)
	mockAdapter.AssertExpectations(t)
}

func TestEReferenceEReferenceTypeGet(t *testing.T) {
	o := newEReferenceImpl()
	assert.Panics(t, func() { o.GetEReferenceType() })
}

func TestEReferenceContainerGet(t *testing.T) {
	o := newEReferenceImpl()
	assert.Panics(t, func() { o.IsContainer() })
}

func TestEReferenceContainmentGet(t *testing.T) {
	o := newEReferenceImpl()
	// get default value
	assert.Equal(t, bool(false), o.IsContainment())
	// get initialized value
	v := bool(true)
	o.isContainment = v
	assert.Equal(t, v, o.IsContainment())
}

func TestEReferenceContainmentSet(t *testing.T) {
	o := newEReferenceImpl()
	v := bool(true)
	mockAdapter := NewMockEAdapter(t)
	mockAdapter.EXPECT().SetTarget(o).Once()
	mockAdapter.EXPECT().NotifyChanged(mock.Anything).Once()
	o.EAdapters().Add(mockAdapter)
	o.SetContainment(v)
	mockAdapter.AssertExpectations(t)
}

func TestEReferenceResolveProxiesGet(t *testing.T) {
	o := newEReferenceImpl()
	// get default value
	assert.Equal(t, bool(true), o.IsResolveProxies())
	// get initialized value
	v := bool(true)
	o.isResolveProxies = v
	assert.Equal(t, v, o.IsResolveProxies())
}

func TestEReferenceResolveProxiesSet(t *testing.T) {
	o := newEReferenceImpl()
	v := bool(true)
	mockAdapter := NewMockEAdapter(t)
	mockAdapter.EXPECT().SetTarget(o).Once()
	mockAdapter.EXPECT().NotifyChanged(mock.Anything).Once()
	o.EAdapters().Add(mockAdapter)
	o.SetResolveProxies(v)
	mockAdapter.AssertExpectations(t)
}

func TestEReferenceEGetFromID(t *testing.T) {
	o := newEReferenceImpl()
	assert.Panics(t, func() { o.EGetFromID(-1, true) })
	assert.Panics(t, func() { o.EGetFromID(EREFERENCE__CONTAINER, true) })
	assert.Panics(t, func() { o.EGetFromID(EREFERENCE__CONTAINER, false) })
	assert.Equal(t, o.IsContainment(), o.EGetFromID(EREFERENCE__CONTAINMENT, true))
	assert.Equal(t, o.GetEKeys(), o.EGetFromID(EREFERENCE__EKEYS, true))
	assert.Equal(t, o.GetEKeys().(EObjectList).GetUnResolvedList(), o.EGetFromID(EREFERENCE__EKEYS, false))
	assert.Equal(t, o.GetEOpposite(), o.EGetFromID(EREFERENCE__EOPPOSITE, true))
	assert.Panics(t, func() { o.EGetFromID(EREFERENCE__EREFERENCE_TYPE, true) })
	assert.Panics(t, func() { o.EGetFromID(EREFERENCE__EREFERENCE_TYPE, false) })
	assert.Equal(t, o.IsResolveProxies(), o.EGetFromID(EREFERENCE__RESOLVE_PROXIES, true))
}

func TestEReferenceESetFromID(t *testing.T) {
	o := newEReferenceImpl()
	assert.Panics(t, func() { o.ESetFromID(-1, nil) })
	{
		v := bool(true)
		o.ESetFromID(EREFERENCE__CONTAINMENT, v)
		assert.Equal(t, v, o.EGetFromID(EREFERENCE__CONTAINMENT, false))
	}
	{
		// list with a value
		mockValue := NewMockEAttribute(t)
		l := NewImmutableEList([]any{mockValue})
		mockValue.EXPECT().EIsProxy().Return(false).Once()

		// set list with new contents
		o.ESetFromID(EREFERENCE__EKEYS, l)
		// checks
		assert.Equal(t, 1, o.GetEKeys().Size())
		assert.Equal(t, mockValue, o.GetEKeys().Get(0))
		mock.AssertExpectationsForObjects(t, mockValue)
	}
	{
		v := NewMockEReference(t)
		o.ESetFromID(EREFERENCE__EOPPOSITE, v)
		assert.Equal(t, v, o.EGetFromID(EREFERENCE__EOPPOSITE, false))

		o.ESetFromID(EREFERENCE__EOPPOSITE, nil)
	}
	{
		v := bool(true)
		o.ESetFromID(EREFERENCE__RESOLVE_PROXIES, v)
		assert.Equal(t, v, o.EGetFromID(EREFERENCE__RESOLVE_PROXIES, false))
	}

}

func TestEReferenceEIsSetFromID(t *testing.T) {
	o := newEReferenceImpl()
	assert.Panics(t, func() { o.EIsSetFromID(-1) })
	assert.Panics(t, func() { o.EIsSetFromID(EREFERENCE__CONTAINER) })
	assert.False(t, o.EIsSetFromID(EREFERENCE__CONTAINMENT))
	assert.False(t, o.EIsSetFromID(EREFERENCE__EKEYS))
	assert.False(t, o.EIsSetFromID(EREFERENCE__EOPPOSITE))
	assert.Panics(t, func() { o.EIsSetFromID(EREFERENCE__EREFERENCE_TYPE) })
	assert.False(t, o.EIsSetFromID(EREFERENCE__RESOLVE_PROXIES))
}

func TestEReferenceEUnsetFromID(t *testing.T) {
	o := newEReferenceImpl()
	assert.Panics(t, func() { o.EUnsetFromID(-1) })
	{
		o.EUnsetFromID(EREFERENCE__CONTAINMENT)
		v := o.EGetFromID(EREFERENCE__CONTAINMENT, false)
		assert.Equal(t, bool(false), v)
	}
	{
		o.EUnsetFromID(EREFERENCE__EKEYS)
		v := o.EGetFromID(EREFERENCE__EKEYS, false)
		assert.NotNil(t, v)
		l := v.(EList)
		assert.True(t, l.Empty())
	}
	{
		o.EUnsetFromID(EREFERENCE__EOPPOSITE)
		assert.Nil(t, o.EGetFromID(EREFERENCE__EOPPOSITE, false))
	}
	{
		o.EUnsetFromID(EREFERENCE__RESOLVE_PROXIES)
		v := o.EGetFromID(EREFERENCE__RESOLVE_PROXIES, false)
		assert.Equal(t, bool(true), v)
	}
}
