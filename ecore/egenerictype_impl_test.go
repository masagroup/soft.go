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

func discardEGenericType() {
	_ = assert.Equal
	_ = mock.Anything
	_ = testing.Coverage
}

func TestEGenericTypeAsEGenericType(t *testing.T) {
	o := newEGenericTypeImpl()
	assert.Equal(t, o, o.asEGenericType())
}

func TestEGenericTypeStaticClass(t *testing.T) {
	o := newEGenericTypeImpl()
	assert.Equal(t, GetPackage().GetEGenericType(), o.EStaticClass())
}

func TestEGenericTypeFeatureCount(t *testing.T) {
	o := newEGenericTypeImpl()
	assert.Equal(t, EGENERIC_TYPE_FEATURE_COUNT, o.EStaticFeatureCount())
}

func TestEGenericTypeEUpperBoundGet(t *testing.T) {
	o := newEGenericTypeImpl()
	// get default value
	assert.Nil(t, o.GetEUpperBound())

	// get initialized value
	v := new(MockEGenericType)
	o.eUpperBound = v
	assert.Equal(t, v, o.GetEUpperBound())
}

func TestEGenericTypeEUpperBoundSet(t *testing.T) {
	o := newEGenericTypeImpl()
	mockValue := new(MockEGenericType)
	mockValue.On("EInverseAdd", o, EOPPOSITE_FEATURE_BASE-EGENERIC_TYPE__EUPPER_BOUND, mock.Anything).Return(nil).Once()
	mockAdapter := new(MockEAdapter)
	mockAdapter.On("SetTarget", o).Once()
	mockAdapter.On("NotifyChanged", mock.Anything).Once()
	o.EAdapters().Add(mockAdapter)
	o.SetEUpperBound(mockValue)
	mock.AssertExpectationsForObjects(t, mockAdapter, mockValue)
}

func TestEGenericTypeETypeArgumentsGet(t *testing.T) {
	o := newEGenericTypeImpl()
	assert.NotNil(t, o.GetETypeArguments())
}

func TestEGenericTypeERawTypeGet(t *testing.T) {
	o := newEGenericTypeImpl()
	// get default value
	assert.Nil(t, o.GetERawType())

	// initialze object with a mock value
	mockValue := new(MockEClassifier)
	o.eRawType = mockValue

	// get non proxy value
	mockValue.On("EIsProxy").Return(false).Once()
	assert.Equal(t, mockValue, o.GetERawType())
	mock.AssertExpectationsForObjects(t, mockValue)

	// get a proxy value
	mockAdapter := new(MockEAdapter)
	mockAdapter.On("SetTarget", o).Once()
	o.EAdapters().Add(mockAdapter)
	mock.AssertExpectationsForObjects(t, mockAdapter)

	mockValue.On("EIsProxy").Return(true).Once()
	mockValue.On("EProxyURI").Return(nil).Once()
	assert.Equal(t, mockValue, o.GetERawType())
	mock.AssertExpectationsForObjects(t, mockAdapter, mockValue)
}

func TestEGenericTypeELowerBoundGet(t *testing.T) {
	o := newEGenericTypeImpl()
	// get default value
	assert.Nil(t, o.GetELowerBound())

	// get initialized value
	v := new(MockEGenericType)
	o.eLowerBound = v
	assert.Equal(t, v, o.GetELowerBound())
}

func TestEGenericTypeELowerBoundSet(t *testing.T) {
	o := newEGenericTypeImpl()
	mockValue := new(MockEGenericType)
	mockValue.On("EInverseAdd", o, EOPPOSITE_FEATURE_BASE-EGENERIC_TYPE__ELOWER_BOUND, mock.Anything).Return(nil).Once()
	mockAdapter := new(MockEAdapter)
	mockAdapter.On("SetTarget", o).Once()
	mockAdapter.On("NotifyChanged", mock.Anything).Once()
	o.EAdapters().Add(mockAdapter)
	o.SetELowerBound(mockValue)
	mock.AssertExpectationsForObjects(t, mockAdapter, mockValue)
}

func TestEGenericTypeETypeParameterGet(t *testing.T) {
	o := newEGenericTypeImpl()
	// get default value
	assert.Nil(t, o.GetETypeParameter())

	// get initialized value
	v := new(MockETypeParameter)
	o.eTypeParameter = v
	assert.Equal(t, v, o.GetETypeParameter())
}

func TestEGenericTypeETypeParameterSet(t *testing.T) {
	o := newEGenericTypeImpl()
	v := new(MockETypeParameter)
	mockAdapter := new(MockEAdapter)
	mockAdapter.On("SetTarget", o).Once()
	mockAdapter.On("NotifyChanged", mock.Anything).Once()
	o.EAdapters().Add(mockAdapter)
	o.SetETypeParameter(v)
	mockAdapter.AssertExpectations(t)
}

func TestEGenericTypeEClassifierGet(t *testing.T) {
	o := newEGenericTypeImpl()
	// get default value
	assert.Nil(t, o.GetEClassifier())

	// initialze object with a mock value
	mockValue := new(MockEClassifier)
	o.eClassifier = mockValue

	// get non proxy value
	mockValue.On("EIsProxy").Return(false).Once()
	assert.Equal(t, mockValue, o.GetEClassifier())
	mock.AssertExpectationsForObjects(t, mockValue)

	// get a proxy value
	mockAdapter := new(MockEAdapter)
	mockAdapter.On("SetTarget", o).Once()
	o.EAdapters().Add(mockAdapter)
	mock.AssertExpectationsForObjects(t, mockAdapter)

	mockValue.On("EIsProxy").Return(true).Once()
	mockValue.On("EProxyURI").Return(nil).Once()
	assert.Equal(t, mockValue, o.GetEClassifier())
	mock.AssertExpectationsForObjects(t, mockAdapter, mockValue)
}

func TestEGenericTypeEClassifierSet(t *testing.T) {
	o := newEGenericTypeImpl()
	v := new(MockEClassifier)
	mockAdapter := new(MockEAdapter)
	mockAdapter.On("SetTarget", o).Once()
	mockAdapter.On("NotifyChanged", mock.Anything).Once()
	o.EAdapters().Add(mockAdapter)
	o.SetEClassifier(v)
	mockAdapter.AssertExpectations(t)
}

func TestEGenericTypeIsInstanceOperation(t *testing.T) {
	o := newEGenericTypeImpl()
	assert.Panics(t, func() { o.IsInstance(nil) })
}

func TestEGenericTypeEGetFromID(t *testing.T) {
	o := newEGenericTypeImpl()
	assert.Panics(t, func() { o.EGetFromID(-1, true) })
	assert.Equal(t, o.GetEUpperBound(), o.EGetFromID(EGENERIC_TYPE__EUPPER_BOUND, true))
	assert.Equal(t, o.GetETypeArguments(), o.EGetFromID(EGENERIC_TYPE__ETYPE_ARGUMENTS, true))
	assert.Equal(t, o.GetETypeArguments().(EObjectList).GetUnResolvedList(), o.EGetFromID(EGENERIC_TYPE__ETYPE_ARGUMENTS, false))
	assert.Equal(t, o.GetERawType(), o.EGetFromID(EGENERIC_TYPE__ERAW_TYPE, true))
	assert.Equal(t, o.GetELowerBound(), o.EGetFromID(EGENERIC_TYPE__ELOWER_BOUND, true))
	assert.Equal(t, o.GetETypeParameter(), o.EGetFromID(EGENERIC_TYPE__ETYPE_PARAMETER, true))
	assert.Equal(t, o.GetEClassifier(), o.EGetFromID(EGENERIC_TYPE__ECLASSIFIER, true))
}

func TestEGenericTypeESetFromID(t *testing.T) {
	o := newEGenericTypeImpl()
	assert.Panics(t, func() { o.ESetFromID(-1, nil) })
	{
		v := new(MockEClassifier)
		o.ESetFromID(EGENERIC_TYPE__ECLASSIFIER, v)
		assert.Equal(t, v, o.EGetFromID(EGENERIC_TYPE__ECLASSIFIER, false))
	}
	{
		mockValue := new(MockEGenericType)
		mockValue.On("EInverseAdd", o, EOPPOSITE_FEATURE_BASE-EGENERIC_TYPE__ELOWER_BOUND, mock.Anything).Return(nil).Once()
		o.ESetFromID(EGENERIC_TYPE__ELOWER_BOUND, mockValue)
		assert.Equal(t, mockValue, o.EGetFromID(EGENERIC_TYPE__ELOWER_BOUND, false))
		mock.AssertExpectationsForObjects(t, mockValue)
	}
	{
		// list with a value
		mockValue := new(MockEGenericType)
		l := NewImmutableEList([]interface{}{mockValue})
		// expectations
		mockValue.On("EInverseAdd", o, EOPPOSITE_FEATURE_BASE-EGENERIC_TYPE__ETYPE_ARGUMENTS, mock.Anything).Return(nil).Once()
		// set list with new contents
		o.ESetFromID(EGENERIC_TYPE__ETYPE_ARGUMENTS, l)
		// checks
		assert.Equal(t, 1, o.GetETypeArguments().Size())
		assert.Equal(t, mockValue, o.GetETypeArguments().Get(0))
		mock.AssertExpectationsForObjects(t, mockValue)
	}
	{
		v := new(MockETypeParameter)
		o.ESetFromID(EGENERIC_TYPE__ETYPE_PARAMETER, v)
		assert.Equal(t, v, o.EGetFromID(EGENERIC_TYPE__ETYPE_PARAMETER, false))
	}
	{
		mockValue := new(MockEGenericType)
		mockValue.On("EInverseAdd", o, EOPPOSITE_FEATURE_BASE-EGENERIC_TYPE__EUPPER_BOUND, mock.Anything).Return(nil).Once()
		o.ESetFromID(EGENERIC_TYPE__EUPPER_BOUND, mockValue)
		assert.Equal(t, mockValue, o.EGetFromID(EGENERIC_TYPE__EUPPER_BOUND, false))
		mock.AssertExpectationsForObjects(t, mockValue)
	}

}

func TestEGenericTypeEIsSetFromID(t *testing.T) {
	o := newEGenericTypeImpl()
	assert.Panics(t, func() { o.EIsSetFromID(-1) })
	assert.False(t, o.EIsSetFromID(EGENERIC_TYPE__ELOWER_BOUND))
	assert.False(t, o.EIsSetFromID(EGENERIC_TYPE__ECLASSIFIER))
	assert.False(t, o.EIsSetFromID(EGENERIC_TYPE__ERAW_TYPE))
	assert.False(t, o.EIsSetFromID(EGENERIC_TYPE__ETYPE_ARGUMENTS))
	assert.False(t, o.EIsSetFromID(EGENERIC_TYPE__ETYPE_PARAMETER))
	assert.False(t, o.EIsSetFromID(EGENERIC_TYPE__EUPPER_BOUND))
}

func TestEGenericTypeEUnsetFromID(t *testing.T) {
	o := newEGenericTypeImpl()
	assert.Panics(t, func() { o.EUnsetFromID(-1) })
	{
		o.EUnsetFromID(EGENERIC_TYPE__ELOWER_BOUND)
		assert.Nil(t, o.EGetFromID(EGENERIC_TYPE__ELOWER_BOUND, false))
	}
	{
		o.EUnsetFromID(EGENERIC_TYPE__ECLASSIFIER)
		assert.Nil(t, o.EGetFromID(EGENERIC_TYPE__ECLASSIFIER, false))
	}
	{
		o.EUnsetFromID(EGENERIC_TYPE__ETYPE_ARGUMENTS)
		v := o.EGetFromID(EGENERIC_TYPE__ETYPE_ARGUMENTS, false)
		assert.NotNil(t, v)
		l := v.(EList)
		assert.True(t, l.Empty())
	}
	{
		o.EUnsetFromID(EGENERIC_TYPE__ETYPE_PARAMETER)
		assert.Nil(t, o.EGetFromID(EGENERIC_TYPE__ETYPE_PARAMETER, false))
	}
	{
		o.EUnsetFromID(EGENERIC_TYPE__EUPPER_BOUND)
		assert.Nil(t, o.EGetFromID(EGENERIC_TYPE__EUPPER_BOUND, false))
	}
}

func TestEGenericTypeEInvokeFromID(t *testing.T) {
	o := newEGenericTypeImpl()
	assert.Panics(t, func() { o.EInvokeFromID(-1, nil) })
	assert.Panics(t, func() { o.EInvokeFromID(EGENERIC_TYPE__IS_INSTANCE_EJAVAOBJECT, nil) })
}
