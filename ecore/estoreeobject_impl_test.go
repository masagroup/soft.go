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

type EStoreEObjectTest struct {
	*EStoreEObjectImpl
	mockStore *MockEStore
}

func newEStoreEObjectTest(t *testing.T, isCaching bool) *EStoreEObjectTest {
	o := new(EStoreEObjectTest)
	o.mockStore = NewMockEStore(t)
	o.EStoreEObjectImpl = NewEStoreEObjectImpl(isCaching)
	o.SetInterfaces(o)
	return o
}

func (o *EStoreEObjectTest) EStore() EStore {
	return o.mockStore
}

func TestEStoreEObjectImpl_GetAttribute_Transient(t *testing.T) {
	// create object
	o := newEStoreEObjectTest(t, true)

	// create mocks
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockStore := o.mockStore

	// initialise object with mock class
	o.SetEClass(mockClass)

	// first get
	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().GetDefaultValue().Return(nil).Once()
	mockAttribute.EXPECT().IsTransient().Return(true).Once()
	mockClass.EXPECT().GetFeatureCount().Return(1).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Twice()
	assert.Nil(t, o.EGetFromID(0, false))
	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)
}

func TestEStoreEObjectImpl_GetAttribute_Caching(t *testing.T) {
	// create object
	o := newEStoreEObjectTest(t, true)

	// create mocks
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockStore := o.mockStore

	// initialise object with mock class
	o.SetEClass(mockClass)

	// first get
	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().IsTransient().Return(false).Once()
	mockClass.EXPECT().GetFeatureCount().Return(1).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Twice()
	mockStore.EXPECT().Get(o, mockAttribute, NO_INDEX).Return(2).Once()
	assert.Equal(t, 2, o.EGetFromID(0, false))
	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)

	// second - test caching
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Once()
	assert.Equal(t, 2, o.EGetFromID(0, false))
	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)
}

func TestEStoreEObjectImpl_GetAttribute_NoCaching(t *testing.T) {
	// create object
	o := newEStoreEObjectTest(t, false)

	// create mocks
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockStore := o.mockStore

	// initialise object with mock class
	o.SetEClass(mockClass)

	mockClass.EXPECT().GetFeatureCount().Return(1).Once()
	for i := 0; i < 2; i++ {
		mockAttribute.EXPECT().IsMany().Return(false).Once()
		mockAttribute.EXPECT().IsTransient().Return(false).Once()
		mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Twice()
		mockStore.EXPECT().Get(o, mockAttribute, NO_INDEX).Return(2).Once()
		assert.Equal(t, 2, o.EGetFromID(0, false))
		mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)
	}
}

func TestEStoreEObjectImpl_SetAttribute_Transient(t *testing.T) {
	// create object
	o := newEStoreEObjectTest(t, true)

	// create mocks
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockStore := o.mockStore

	// initialise object with mock class
	o.SetEClass(mockClass)

	// set
	mockAttribute.EXPECT().IsTransient().Return(true).Twice()
	mockClass.EXPECT().GetFeatureCount().Return(1).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Times(3)
	o.ESetFromID(0, 2)
	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)

	// test get
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Once()
	assert.Equal(t, 2, o.EGetFromID(0, false))
	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)
}

func TestEStoreEObjectImpl_SetAttribute(t *testing.T) {

	// create object
	o := newEStoreEObjectTest(t, false)

	// create mocks
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockStore := o.mockStore

	// initialise object with mock class
	o.SetEClass(mockClass)

	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().IsTransient().Return(false).Twice()
	mockClass.EXPECT().GetFeatureCount().Return(1).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Times(3)
	mockStore.EXPECT().Get(o, mockAttribute, NO_INDEX).Return(nil).Once()
	mockStore.EXPECT().Set(o, mockAttribute, NO_INDEX, 2).Return(nil).Once()
	o.ESetFromID(0, 2)
	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)
}

func TestEStoreEObjectImpl_GetAttribute_Many(t *testing.T) {

	// create object
	o := newEStoreEObjectTest(t, false)

	// create mocks
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockStore := o.mockStore

	// initialise object with mock class
	o.SetEClass(mockClass)

	mockAttribute.EXPECT().IsMany().Return(true).Once()
	mockAttribute.EXPECT().IsTransient().Return(false).Once()
	mockClass.EXPECT().GetFeatureCount().Return(1).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Twice()
	list := o.EGetFromID(0, true)
	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)

	eobjectlist, _ := list.(EObjectList)
	assert.NotNil(t, eobjectlist)
	enotifyinglist, _ := list.(ENotifyingList)
	assert.NotNil(t, enotifyinglist)
}

func TestEStoreEObjectImpl_SetAttribute_Caching(t *testing.T) {

	// create object
	o := newEStoreEObjectTest(t, true)

	// create mocks
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockStore := o.mockStore

	// initialise object with mock class
	o.SetEClass(mockClass)

	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().IsTransient().Return(false).Twice()
	mockClass.EXPECT().GetFeatureCount().Return(1).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Times(3)
	mockStore.EXPECT().Get(o, mockAttribute, NO_INDEX).Return(nil).Once()
	mockStore.EXPECT().Set(o, mockAttribute, NO_INDEX, 2).Return(nil).Once()
	o.ESetFromID(0, 2)
	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)
}

func TestEStoreEObjectImpl_UnSetAttribute_Transient(t *testing.T) {
	// create object
	o := newEStoreEObjectTest(t, false)

	// create mocks
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)

	// initialise object with mock class
	o.SetEClass(mockClass)

	mockAttribute.EXPECT().IsTransient().Return(true).Twice()
	mockClass.EXPECT().GetFeatureCount().Return(1).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Times(3)
	o.EUnsetFromID(0)

}

func TestEStoreEObjectImpl_UnSetAttribute(t *testing.T) {

	// create object
	o := newEStoreEObjectTest(t, false)

	// create mocks
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockStore := o.mockStore

	// initialise object with mock class
	o.SetEClass(mockClass)

	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().IsTransient().Return(false).Twice()
	mockClass.EXPECT().GetFeatureCount().Return(1).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Times(3)
	mockStore.EXPECT().Get(o, mockAttribute, NO_INDEX).Return(2).Once()
	mockStore.EXPECT().UnSet(o, mockAttribute).Once()
	o.EUnsetFromID(0)
	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)
}
