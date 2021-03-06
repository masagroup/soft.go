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
	"net/url"
	"testing"
)

func discardEOperation() {
	_ = assert.Equal
	_ = mock.Anything
	_ = testing.Coverage
	_ = url.Parse
}

func TestEOperationAsEOperation(t *testing.T) {
	o := newEOperationImpl()
	assert.Equal(t, o, o.asEOperation())
}

func TestEOperationStaticClass(t *testing.T) {
	o := newEOperationImpl()
	assert.Equal(t, GetPackage().GetEOperation(), o.EStaticClass())
}

func TestEOperationFeatureCount(t *testing.T) {
	o := newEOperationImpl()
	assert.Equal(t, EOPERATION_FEATURE_COUNT, o.EStaticFeatureCount())
}

func TestEOperationEContainingClassGet(t *testing.T) {
	// default
	o := newEOperationImpl()
	assert.Nil(t, o.GetEContainingClass())

	// set a mock container
	v := new(MockEClass)
	o.ESetInternalContainer(v, EOPERATION__ECONTAINING_CLASS)

	// no proxy
	v.On("EIsProxy").Return(false)
	assert.Equal(t, v, o.GetEContainingClass())
}

func TestEOperationEExceptionsGet(t *testing.T) {
	o := newEOperationImpl()
	assert.NotNil(t, o.GetEExceptions())
}

func TestEOperationEExceptionsUnSet(t *testing.T) {
	o := newEOperationImpl()
	o.UnsetEExceptions()
	l := o.GetEExceptions()
	assert.True(t, l.Empty())
}

func TestEOperationEParametersGet(t *testing.T) {
	o := newEOperationImpl()
	assert.NotNil(t, o.GetEParameters())
}

func TestEOperationOperationIDGet(t *testing.T) {
	o := newEOperationImpl()
	// get default value
	assert.Equal(t, int(-1), o.GetOperationID())
	// get initialized value
	v := int(45)
	o.operationID = v
	assert.Equal(t, v, o.GetOperationID())
}

func TestEOperationOperationIDSet(t *testing.T) {
	o := newEOperationImpl()
	v := int(45)
	mockAdapter := new(MockEAdapter)
	mockAdapter.On("SetTarget", o).Once()
	mockAdapter.On("NotifyChanged", mock.Anything).Once()
	o.EAdapters().Add(mockAdapter)
	o.SetOperationID(v)
	mockAdapter.AssertExpectations(t)
}

func TestEOperationIsOverrideOfOperation(t *testing.T) {
	o := newEOperationImpl()
	assert.Panics(t, func() { o.IsOverrideOf(nil) })
}

func TestEOperationEGetFromID(t *testing.T) {
	o := newEOperationImpl()
	assert.Panics(t, func() { o.EGetFromID(-1, true) })
	assert.Equal(t, o.GetEContainingClass(), o.EGetFromID(EOPERATION__ECONTAINING_CLASS, true))
	assert.Equal(t, o.GetEExceptions(), o.EGetFromID(EOPERATION__EEXCEPTIONS, true))
	assert.Equal(t, o.GetEExceptions().(EObjectList).GetUnResolvedList(), o.EGetFromID(EOPERATION__EEXCEPTIONS, false))
	assert.Equal(t, o.GetEParameters(), o.EGetFromID(EOPERATION__EPARAMETERS, true))
	assert.Equal(t, o.GetEParameters().(EObjectList).GetUnResolvedList(), o.EGetFromID(EOPERATION__EPARAMETERS, false))
	assert.Equal(t, o.GetOperationID(), o.EGetFromID(EOPERATION__OPERATION_ID, true))
}

func TestEOperationESetFromID(t *testing.T) {
	o := newEOperationImpl()
	assert.Panics(t, func() { o.ESetFromID(-1, nil) })
	{
		// list with a value
		mockValue := new(MockEClassifier)
		l := NewImmutableEList([]interface{}{mockValue})
		mockValue.On("EIsProxy").Return(false).Once()

		// set list with new contents
		o.ESetFromID(EOPERATION__EEXCEPTIONS, l)
		// checks
		assert.Equal(t, 1, o.GetEExceptions().Size())
		assert.Equal(t, mockValue, o.GetEExceptions().Get(0))
		mock.AssertExpectationsForObjects(t, mockValue)
	}
	{
		// list with a value
		mockValue := new(MockEParameter)
		l := NewImmutableEList([]interface{}{mockValue})
		mockValue.On("EInverseAdd", o, EPARAMETER__EOPERATION, mock.Anything).Return(nil).Once()

		// set list with new contents
		o.ESetFromID(EOPERATION__EPARAMETERS, l)
		// checks
		assert.Equal(t, 1, o.GetEParameters().Size())
		assert.Equal(t, mockValue, o.GetEParameters().Get(0))
		mock.AssertExpectationsForObjects(t, mockValue)
	}
	{
		v := int(45)
		o.ESetFromID(EOPERATION__OPERATION_ID, v)
		assert.Equal(t, v, o.EGetFromID(EOPERATION__OPERATION_ID, false))
	}

}

func TestEOperationEIsSetFromID(t *testing.T) {
	o := newEOperationImpl()
	assert.Panics(t, func() { o.EIsSetFromID(-1) })
	assert.False(t, o.EIsSetFromID(EOPERATION__ECONTAINING_CLASS))
	assert.False(t, o.EIsSetFromID(EOPERATION__EEXCEPTIONS))
	assert.False(t, o.EIsSetFromID(EOPERATION__EPARAMETERS))
	assert.False(t, o.EIsSetFromID(EOPERATION__OPERATION_ID))
}

func TestEOperationEUnsetFromID(t *testing.T) {
	o := newEOperationImpl()
	assert.Panics(t, func() { o.EUnsetFromID(-1) })
	{
		o.EUnsetFromID(EOPERATION__EEXCEPTIONS)
		v := o.EGetFromID(EOPERATION__EEXCEPTIONS, false)
		assert.NotNil(t, v)
		l := v.(EList)
		assert.True(t, l.Empty())
	}
	{
		o.EUnsetFromID(EOPERATION__EPARAMETERS)
		v := o.EGetFromID(EOPERATION__EPARAMETERS, false)
		assert.NotNil(t, v)
		l := v.(EList)
		assert.True(t, l.Empty())
	}
	{
		o.EUnsetFromID(EOPERATION__OPERATION_ID)
		v := o.EGetFromID(EOPERATION__OPERATION_ID, false)
		assert.Equal(t, int(-1), v)
	}
}

func TestEOperationEInvokeFromID(t *testing.T) {
	o := newEOperationImpl()
	assert.Panics(t, func() { o.EInvokeFromID(-1, nil) })
	assert.Panics(t, func() { o.EInvokeFromID(EOPERATION__IS_OVERRIDE_OF_EOPERATION, nil) })
}

func TestEOperationEBasicInverseAdd(t *testing.T) {
	o := newEOperationImpl()
	{
		mockObject := new(MockEObject)
		mockNotifications := new(MockENotificationChain)
		assert.Equal(t, mockNotifications, o.EBasicInverseAdd(mockObject, -1, mockNotifications))
	}
	{
		mockObject := new(MockEClass)
		mockObject.On("EResource").Return(nil).Once()
		mockObject.On("EIsProxy").Return(false).Once()
		o.EBasicInverseAdd(mockObject, EOPERATION__ECONTAINING_CLASS, nil)
		assert.Equal(t, mockObject, o.GetEContainingClass())
		mock.AssertExpectationsForObjects(t, mockObject)

		mockOther := new(MockEClass)
		mockOther.On("EResource").Return(nil).Once()
		mockOther.On("EIsProxy").Return(false).Once()
		mockObject.On("EResource").Return(nil).Once()
		mockObject.On("EInverseRemove", o, ECLASS__EOPERATIONS, nil).Return(nil).Once()
		o.EBasicInverseAdd(mockOther, EOPERATION__ECONTAINING_CLASS, nil)
		assert.Equal(t, mockOther, o.GetEContainingClass())
		mock.AssertExpectationsForObjects(t, mockObject, mockOther)
	}
	{
		mockObject := new(MockEParameter)
		o.EBasicInverseAdd(mockObject, EOPERATION__EPARAMETERS, nil)
		l := o.GetEParameters()
		assert.True(t, l.Contains(mockObject))
		mock.AssertExpectationsForObjects(t, mockObject)
	}

}

func TestEOperationEBasicInverseRemove(t *testing.T) {
	o := newEOperationImpl()
	{
		mockObject := new(MockEObject)
		mockNotifications := new(MockENotificationChain)
		assert.Equal(t, mockNotifications, o.EBasicInverseRemove(mockObject, -1, mockNotifications))
	}
	{
		mockObject := new(MockEClass)
		o.EBasicInverseRemove(mockObject, EOPERATION__ECONTAINING_CLASS, nil)
		mock.AssertExpectationsForObjects(t, mockObject)
	}
	{
		// initialize list with a mock object
		mockObject := new(MockEParameter)
		mockObject.On("EInverseAdd", o, EPARAMETER__EOPERATION, mock.Anything).Return(nil).Once()

		l := o.GetEParameters()
		l.Add(mockObject)

		// basic inverse remove
		o.EBasicInverseRemove(mockObject, EOPERATION__EPARAMETERS, nil)

		// check it was removed
		assert.False(t, l.Contains(mockObject))
		mock.AssertExpectationsForObjects(t, mockObject)
	}

}
