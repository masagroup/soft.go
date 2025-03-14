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
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestEStoreEObjectImpl_GetAttribute_Transient(t *testing.T) {

	// create mocks
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockStore := NewMockEStore(t)

	// create object
	o := NewEStoreEObjectImpl(true)
	o.SetEClass(mockClass)
	o.SetEStore(mockStore)

	// first get
	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().GetDefaultValue().Return(nil).Once()
	mockAttribute.EXPECT().IsTransient().Return(true).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Twice()
	assert.Nil(t, o.EGetFromID(0, false))
	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)
}

func TestEStoreEObjectImpl_GetAttribute_Caching(t *testing.T) {
	// create mocks
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockStore := NewMockEStore(t)

	// create object
	o := NewEStoreEObjectImpl(true)
	o.SetEClass(mockClass)
	o.SetEStore(mockStore)

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
	// create mocks
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockStore := NewMockEStore(t)

	// create object
	o := NewEStoreEObjectImpl(false)
	o.SetEClass(mockClass)
	o.SetEStore(mockStore)

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
	// create mocks
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockStore := NewMockEStore(t)

	// create object
	o := NewEStoreEObjectImpl(true)
	o.SetEClass(mockClass)
	o.SetEStore(mockStore)

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
	// create mocks
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockStore := NewMockEStore(t)

	// create object
	o := NewEStoreEObjectImpl(false)
	o.SetEClass(mockClass)
	o.SetEStore(mockStore)

	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().IsTransient().Return(false).Twice()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Times(3)
	mockStore.EXPECT().Get(o, mockAttribute, NO_INDEX).Return(nil).Once()
	mockStore.EXPECT().Set(o, mockAttribute, NO_INDEX, 2, false).Return(nil).Once()
	o.ESetFromID(0, 2)
	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)
}

func TestEStoreEObjectImpl_SetAttribute_Caching(t *testing.T) {
	// create mocks
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockStore := NewMockEStore(t)

	// create object
	o := NewEStoreEObjectImpl(true)
	o.SetEClass(mockClass)
	o.SetEStore(mockStore)

	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().IsTransient().Return(false).Twice()
	mockClass.EXPECT().GetFeatureCount().Return(1).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Times(3)
	mockStore.EXPECT().Get(o, mockAttribute, NO_INDEX).Return(nil).Once()
	mockStore.EXPECT().Set(o, mockAttribute, NO_INDEX, 2, false).Return(nil).Once()
	o.ESetFromID(0, 2)
	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)
}

func TestEStoreEObjectImpl_GetAttribute_Many_List(t *testing.T) {
	// create mocks
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockStore := NewMockEStore(t)

	// create object
	o := NewEStoreEObjectImpl(false)
	o.SetEClass(mockClass)
	o.SetEStore(mockStore)

	mockAttribute.EXPECT().IsUnique().Return(true).Once()
	mockAttribute.EXPECT().IsMany().Return(true).Once()
	mockAttribute.EXPECT().IsTransient().Return(false).Once()
	mockAttribute.EXPECT().GetEType().Return(nil).Once()
	mockClass.EXPECT().GetFeatureCount().Return(1).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Twice()
	mockStore.EXPECT().Size(o, mockAttribute).Return(0).Once()
	list := o.EGetFromID(0, true)
	require.NotNil(t, list)
	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)

	eobjectlist, _ := list.(EObjectList)
	assert.NotNil(t, eobjectlist)
	enotifyinglist, _ := list.(ENotifyingList)
	assert.NotNil(t, enotifyinglist)
}

func TestEStoreEObjectImpl_GetAttribute_Many_Map(t *testing.T) {
	// create mocks
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockStore := NewMockEStore(t)
	mockType := NewMockEClass(t)
	mockFeature := NewMockEStructuralFeature(t)

	// create object
	o := NewEStoreEObjectImpl(false)
	o.SetEClass(mockClass)
	o.SetEStore(mockStore)

	mockAttribute.EXPECT().IsUnique().Return(true).Once()
	mockAttribute.EXPECT().IsMany().Return(true).Once()
	mockAttribute.EXPECT().IsTransient().Return(false).Once()
	mockAttribute.EXPECT().GetEType().Return(mockType).Twice()
	mockType.EXPECT().GetInstanceTypeName().Return("java.util.Map.Entry").Once()
	mockType.EXPECT().GetEStructuralFeatureFromName("key").Return(mockFeature).Once()
	mockType.EXPECT().GetEStructuralFeatureFromName("value").Return(mockFeature).Once()
	mockClass.EXPECT().GetFeatureCount().Return(1).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Twice()
	mockStore.EXPECT().Size(o, mockAttribute).Return(0).Once()
	m := o.EGetFromID(0, true)
	require.NotNil(t, m)
	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore, mockFeature, mockType)

	emap, _ := m.(EMap)
	assert.NotNil(t, emap)
}

func TestEStoreEObjectImpl_UnSetAttribute_Transient(t *testing.T) {
	// create mocks
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockStore := NewMockEStore(t)

	// create object
	o := NewEStoreEObjectImpl(false)
	o.SetEClass(mockClass)
	o.SetEStore(mockStore)

	mockAttribute.EXPECT().IsTransient().Return(true).Twice()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Times(3)
	o.EUnsetFromID(0)

}

func TestEStoreEObjectImpl_UnSetAttribute(t *testing.T) {
	// create mocks
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	mockStore := NewMockEStore(t)

	// create object
	o := NewEStoreEObjectImpl(false)
	o.SetEClass(mockClass)
	o.SetEStore(mockStore)

	// create cache
	mockClass.EXPECT().GetFeatureCount().Return(1).Once()
	o.getProperties()

	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().IsTransient().Return(false).Twice()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Times(3)
	mockStore.EXPECT().Get(o, mockAttribute, NO_INDEX).Return(2).Once()
	mockStore.EXPECT().UnSet(o, mockAttribute).Once()
	o.EUnsetFromID(0)
	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)
}

func TestEStoreEObjectImpl_GetContainer(t *testing.T) {
	mockClass := NewMockEClass(t)
	mockStore := NewMockEStore(t)
	mockObject := NewMockEObject(t)
	mockObjectClass := NewMockEClass(t)
	mockReference := NewMockEReference(t)
	mockOpposite := NewMockEReference(t)

	// create object
	o := NewEStoreEObjectImpl(false)
	o.SetEClass(mockClass)
	o.SetEStore(mockStore)

	mockStore.EXPECT().GetContainer(o.AsEObject()).Return(mockObject, mockReference).Once()
	mockObject.EXPECT().EClass().Return(mockObjectClass).Once()
	mockObjectClass.EXPECT().GetFeatureID(mockReference).Return(1).Once()
	mockReference.EXPECT().GetEOpposite().Return(mockOpposite).Once()
	mockClass.EXPECT().GetFeatureID(mockOpposite).Return(1).Once()
	mockObject.EXPECT().EIsProxy().Return(false).Once()
	require.Equal(t, mockObject, o.EContainer())
	require.Equal(t, 1, o.EContainerFeatureID())
}

func TestEStoreEObjectImpl_IsSet_NoCache_NoStore(t *testing.T) {
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	o := NewEStoreEObjectImpl(false)
	o.SetEClass(mockClass)
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Once()
	require.False(t, o.EIsSetFromID(0))
}

func TestEStoreEObjectImpl_IsSet_NoCache(t *testing.T) {
	mockClass := NewMockEClass(t)
	mockStore := NewMockEStore(t)
	mockAttribute := NewMockEAttribute(t)
	o := NewEStoreEObjectImpl(false)
	o.SetEClass(mockClass)
	o.SetEStore(mockStore)
	mockStore.EXPECT().IsSet(o, mockAttribute).Return(true).Once()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute).Twice()
	require.True(t, o.EIsSetFromID(0))
}

func TestEStoreEObjectImpl_IsSet_WithCache(t *testing.T) {
	mockClass := NewMockEClass(t)
	mockAttribute := NewMockEAttribute(t)
	o := NewEStoreEObjectImpl(true)
	o.SetEClass(mockClass)
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockAttribute)
	mockAttribute.EXPECT().IsTransient().Return(false).Once()
	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockClass.EXPECT().GetFeatureCount().Return(1).Once()
	o.ESetFromID(0, 2)
	require.True(t, o.EIsSetFromID(0))
}

type EStoreEObjectImplWithCacheTestSuite struct {
	suite.Suite
	mockClass   *MockEClass
	mockStore   *MockEStore
	mockFeature *MockEAttribute
	o           *EStoreEObjectImpl
}

func (suite *EStoreEObjectImplWithCacheTestSuite) SetupTest() {
	t := suite.T()
	suite.mockClass = NewMockEClass(t)
	suite.mockStore = NewMockEStore(t)
	suite.mockFeature = NewMockEAttribute(t)
	suite.o = NewEStoreEObjectImpl(true)
	suite.o.SetEClass(suite.mockClass)
	suite.mockFeature.EXPECT().IsUnique().Return(true).Once()
	suite.mockFeature.EXPECT().IsMany().Return(true).Once()
	suite.mockFeature.EXPECT().IsTransient().Return(false).Once()
	suite.mockFeature.EXPECT().GetEType().Return(nil).Once()
	suite.mockClass.EXPECT().GetFeatureCount().Return(1).Once()
	suite.mockClass.EXPECT().GetEStructuralFeature(0).Return(suite.mockFeature).Twice()
	list := suite.o.EGetFromID(0, true)
	require.NotNil(t, list)
}

func (suite *EStoreEObjectImplWithCacheTestSuite) TestLockUnlock() {
	suite.o.Lock()
	defer suite.o.Unlock()
	require.Equal(suite.T(), suite.mockClass, suite.o.EClass())
}

func (suite *EStoreEObjectImplWithCacheTestSuite) TestSetStore() {
	t := suite.T()
	mockStore := NewMockEStore(t)
	require.Equal(t, nil, suite.o.GetEStore())
	suite.o.SetEStore(mockStore)
	require.Equal(t, mockStore, suite.o.GetEStore())
	suite.o.SetEStore(nil)
	require.Equal(t, nil, suite.o.GetEStore())
}

func (suite *EStoreEObjectImplWithCacheTestSuite) TestSetCache() {
	t := suite.T()
	require.True(t, suite.o.IsCache())
	suite.o.SetCache(false)
	require.False(t, suite.o.IsCache())
	suite.o.SetCache(true)
	require.True(t, suite.o.IsCache())
}

func TestEStoreEObjectImplWithCache(t *testing.T) {
	suite.Run(t, &EStoreEObjectImplWithCacheTestSuite{})
}

type EStoreEObjectImplNoCacheTestSuite struct {
	suite.Suite
	mockClass       *MockEClass
	mockStore       *MockEStore
	mockFeatureList *MockEAttribute
	mockAttribute   *MockEAttribute
	o               *EStoreEObjectImpl
}

func (suite *EStoreEObjectImplNoCacheTestSuite) SetupTest() {
	t := suite.T()
	suite.mockClass = NewMockEClass(t)
	suite.mockStore = NewMockEStore(t)
	suite.mockFeatureList = NewMockEAttribute(t)
	suite.mockAttribute = NewMockEAttribute(t)
	suite.o = NewEStoreEObjectImpl(false)
	suite.o.SetEClass(suite.mockClass)
	suite.mockFeatureList.EXPECT().IsUnique().Return(true).Once()
	suite.mockFeatureList.EXPECT().IsMany().Return(true).Once()
	suite.mockFeatureList.EXPECT().IsTransient().Return(false).Once()
	suite.mockFeatureList.EXPECT().GetEType().Return(nil).Once()
	suite.mockClass.EXPECT().GetFeatureCount().Return(2).Once()
	suite.mockClass.EXPECT().GetEStructuralFeature(0).Return(suite.mockFeatureList).Twice()
	list := suite.o.EGetFromID(0, true)
	require.NotNil(t, list)
}

func (suite *EStoreEObjectImplNoCacheTestSuite) TestSetStore() {
	t := suite.T()
	mockStore := NewMockEStore(t)
	require.Equal(t, nil, suite.o.GetEStore())
	suite.o.SetEStore(mockStore)
	require.Equal(t, mockStore, suite.o.GetEStore())

	suite.mockClass.EXPECT().GetFeatureCount().Return(2).Once()
	suite.mockClass.EXPECT().GetEStructuralFeature(0).Return(suite.mockFeatureList).Once()
	suite.mockFeatureList.EXPECT().IsTransient().Return(false).Once()
	suite.mockFeatureList.EXPECT().IsMany().Return(true).Once()
	suite.mockClass.EXPECT().GetEStructuralFeature(1).Return(suite.mockAttribute).Once()
	suite.mockAttribute.EXPECT().IsTransient().Return(false).Once()
	suite.mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockStore.EXPECT().Get(suite.o, suite.mockAttribute, NO_INDEX).Return(2).Once()
	suite.o.SetEStore(nil)
	require.Equal(t, nil, suite.o.GetEStore())
}

func (suite *EStoreEObjectImplNoCacheTestSuite) TestSetCache() {
	t := suite.T()
	require.False(t, suite.o.IsCache())
	suite.o.SetCache(true)
	require.True(t, suite.o.IsCache())
	suite.o.SetCache(false)
	require.False(t, suite.o.IsCache())
}

func TestEStoreEObjectImplNoCache(t *testing.T) {
	suite.Run(t, &EStoreEObjectImplNoCacheTestSuite{})
}
