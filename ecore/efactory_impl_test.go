// *****************************************************************************
//
// This file is part of a MASA library or program.
// Refer to the included end-user license agreement for restrictions.
//
// Copyright (c) 2020 MASA Group
//
// *****************************************************************************

// *****************************************************************************
//
// Warning: This file was generated by soft.generator.go Generator
//
// *****************************************************************************

package ecore

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func discardEFactory() {
	_ = assert.Equal
	_ = mock.Anything
	_ = testing.Coverage
}

func TestEFactoryAsEFactory(t *testing.T) {
	o := newEFactoryImpl()
	assert.Equal(t, o, o.asEFactory())
}

func TestEFactoryStaticClass(t *testing.T) {
	o := newEFactoryImpl()
	assert.Equal(t, GetPackage().GetEFactory(), o.EStaticClass())
}

func TestEFactoryFeatureCount(t *testing.T) {
	o := newEFactoryImpl()
	assert.Equal(t, EFACTORY_FEATURE_COUNT, o.EStaticFeatureCount())
}

func TestEFactoryEPackageGet(t *testing.T) {
	// default
	o := newEFactoryImpl()
	assert.Nil(t, o.GetEPackage())

	// set a mock container
	v := new(MockEPackage)
	o.ESetInternalContainer(v, EFACTORY__EPACKAGE)

	// no proxy
	v.On("EIsProxy").Return(false)
	assert.Equal(t, v, o.GetEPackage())
}

func TestEFactoryEPackageSet(t *testing.T) {
	// object
	o := newEFactoryImpl()

	// add listener
	mockAdapter := new(MockEAdapter)
	mockAdapter.On("SetTarget", o).Once()
	o.EAdapters().Add(mockAdapter)
	mock.AssertExpectationsForObjects(t, mockAdapter)

	// set with the mock value
	mockValue := new(MockEPackage)
	mockResource := new(MockEResource)
	mockValue.On("EInverseAdd", o, EPACKAGE__EFACTORY_INSTANCE, nil).Return(nil).Once()
	mockValue.On("EInternalResource").Return(mockResource).Once()
	mockResource.On("Attached", o).Once()
	mockAdapter.On("NotifyChanged", mock.Anything).Once()
	o.SetEPackage(mockValue)
	mock.AssertExpectationsForObjects(t, mockAdapter, mockValue, mockResource)

	// another value - in a different resource
	//mockNotifications := new(MockENotificationChain)
	mockValue2 := new(MockEPackage)
	mockResource2 := new(MockEResource)
	mockValue.On("EInverseRemove", o, EPACKAGE__EFACTORY_INSTANCE, nil).Return(nil).Once()
	mockValue.On("EInternalResource").Return(mockResource).Once()
	mockValue2.On("EInverseAdd", o, EPACKAGE__EFACTORY_INSTANCE, nil).Return(nil).Once()
	mockValue2.On("EInternalResource").Return(mockResource2).Once()
	mockResource.On("Detached", o).Once()
	mockResource2.On("Attached", o).Once()
	mockAdapter.On("NotifyChanged", mock.Anything).Once()
	o.SetEPackage(mockValue2)
	mock.AssertExpectationsForObjects(t, mockAdapter, mockValue, mockResource, mockValue2, mockResource2)
}

func TestEFactoryCreateOperation(t *testing.T) {
	o := newEFactoryImpl()
	assert.Panics(t, func() { o.Create(nil) })
}
func TestEFactoryCreateFromStringOperation(t *testing.T) {
	o := newEFactoryImpl()
	assert.Panics(t, func() { o.CreateFromString(nil, "") })
}
func TestEFactoryConvertToStringOperation(t *testing.T) {
	o := newEFactoryImpl()
	assert.Panics(t, func() { o.ConvertToString(nil, nil) })
}

func TestEFactoryEGetFromID(t *testing.T) {
	o := newEFactoryImpl()
	assert.Panics(t, func() { o.EGetFromID(-1, true) })
	assert.Equal(t, o.GetEPackage(), o.EGetFromID(EFACTORY__EPACKAGE, true))
}

func TestEFactoryESetFromID(t *testing.T) {
	o := newEFactoryImpl()
	assert.Panics(t, func() { o.ESetFromID(-1, nil) })
	{
		mockValue := new(MockEPackage)
		mockValue.On("EIsProxy").Return(false).Once()
		mockValue.On("EInternalResource").Return(nil).Once()
		mockValue.On("EInverseAdd", o, EPACKAGE__EFACTORY_INSTANCE, nil).Return(nil).Once()
		o.ESetFromID(EFACTORY__EPACKAGE, mockValue)
		assert.Equal(t, mockValue, o.EGetFromID(EFACTORY__EPACKAGE, false))
		mock.AssertExpectationsForObjects(t, mockValue)
	}

}

func TestEFactoryEIsSetFromID(t *testing.T) {
	o := newEFactoryImpl()
	assert.Panics(t, func() { o.EIsSetFromID(-1) })
	assert.False(t, o.EIsSetFromID(EFACTORY__EPACKAGE))
}

func TestEFactoryEUnsetFromID(t *testing.T) {
	o := newEFactoryImpl()
	assert.Panics(t, func() { o.EUnsetFromID(-1) })
	{
		o.EUnsetFromID(EFACTORY__EPACKAGE)
		assert.Nil(t, o.EGetFromID(EFACTORY__EPACKAGE, false))
	}
}

func TestEFactoryEInvokeFromID(t *testing.T) {
	o := newEFactoryImpl()
	assert.Panics(t, func() { o.EInvokeFromID(-1, nil) })
	assert.Panics(t, func() { o.EInvokeFromID(EFACTORY__CONVERT_TO_STRING_EDATATYPE_EJAVAOBJECT, nil) })
	assert.Panics(t, func() { o.EInvokeFromID(EFACTORY__CREATE_FROM_STRING_EDATATYPE_ESTRING, nil) })
	assert.Panics(t, func() { o.EInvokeFromID(EFACTORY__CREATE_ECLASS, nil) })
}
