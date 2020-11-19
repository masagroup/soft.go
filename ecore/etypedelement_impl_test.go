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

func discardETypedElement() {
	_ = assert.Equal
	_ = mock.Anything
	_ = testing.Coverage
}

func TestETypedElementAsETypedElement(t *testing.T) {
	o := newETypedElementImpl()
	assert.Equal(t, o, o.asETypedElement())
}

func TestETypedElementStaticClass(t *testing.T) {
	o := newETypedElementImpl()
	assert.Equal(t, GetPackage().GetETypedElement(), o.EStaticClass())
}

func TestETypedElementFeatureCount(t *testing.T) {
	o := newETypedElementImpl()
	assert.Equal(t, ETYPED_ELEMENT_FEATURE_COUNT, o.EStaticFeatureCount())
}

func TestETypedElementOrderedGet(t *testing.T) {
	o := newETypedElementImpl()
	// get default value
	assert.Equal(t, true, o.IsOrdered())
	// get initialized value
	v := true
	o.isOrdered = v
	assert.Equal(t, v, o.IsOrdered())
}

func TestETypedElementOrderedSet(t *testing.T) {
	o := newETypedElementImpl()
	v := true
	mockAdapter := new(MockEAdapter)
	mockAdapter.On("SetTarget", o).Once()
	mockAdapter.On("NotifyChanged", mock.Anything).Once()
	o.EAdapters().Add(mockAdapter)
	o.SetOrdered(v)
	mockAdapter.AssertExpectations(t)
}

func TestETypedElementUniqueGet(t *testing.T) {
	o := newETypedElementImpl()
	// get default value
	assert.Equal(t, true, o.IsUnique())
	// get initialized value
	v := true
	o.isUnique = v
	assert.Equal(t, v, o.IsUnique())
}

func TestETypedElementUniqueSet(t *testing.T) {
	o := newETypedElementImpl()
	v := true
	mockAdapter := new(MockEAdapter)
	mockAdapter.On("SetTarget", o).Once()
	mockAdapter.On("NotifyChanged", mock.Anything).Once()
	o.EAdapters().Add(mockAdapter)
	o.SetUnique(v)
	mockAdapter.AssertExpectations(t)
}

func TestETypedElementLowerBoundGet(t *testing.T) {
	o := newETypedElementImpl()
	// get default value
	assert.Equal(t, 0, o.GetLowerBound())
	// get initialized value
	v := 45
	o.lowerBound = v
	assert.Equal(t, v, o.GetLowerBound())
}

func TestETypedElementLowerBoundSet(t *testing.T) {
	o := newETypedElementImpl()
	v := 45
	mockAdapter := new(MockEAdapter)
	mockAdapter.On("SetTarget", o).Once()
	mockAdapter.On("NotifyChanged", mock.Anything).Once()
	o.EAdapters().Add(mockAdapter)
	o.SetLowerBound(v)
	mockAdapter.AssertExpectations(t)
}

func TestETypedElementUpperBoundGet(t *testing.T) {
	o := newETypedElementImpl()
	// get default value
	assert.Equal(t, 1, o.GetUpperBound())
	// get initialized value
	v := 45
	o.upperBound = v
	assert.Equal(t, v, o.GetUpperBound())
}

func TestETypedElementUpperBoundSet(t *testing.T) {
	o := newETypedElementImpl()
	v := 45
	mockAdapter := new(MockEAdapter)
	mockAdapter.On("SetTarget", o).Once()
	mockAdapter.On("NotifyChanged", mock.Anything).Once()
	o.EAdapters().Add(mockAdapter)
	o.SetUpperBound(v)
	mockAdapter.AssertExpectations(t)
}

func TestETypedElementManyGet(t *testing.T) {
	o := newETypedElementImpl()
	assert.Panics(t, func() { o.IsMany() })
}

func TestETypedElementRequiredGet(t *testing.T) {
	o := newETypedElementImpl()
	assert.Panics(t, func() { o.IsRequired() })
}

func TestETypedElementETypeGet(t *testing.T) {
	o := newETypedElementImpl()
	// get default value
	assert.Nil(t, o.GetEType())

	// initialze object with a mock value
	mockValue := new(MockEClassifier)
	o.eType = mockValue

	// get non proxy value
	mockValue.On("EIsProxy").Return(false).Once()
	assert.Equal(t, mockValue, o.GetEType())
	mock.AssertExpectationsForObjects(t, mockValue)

	// get a proxy value
	mockAdapter := new(MockEAdapter)
	mockAdapter.On("SetTarget", o).Once()
	o.EAdapters().Add(mockAdapter)
	mock.AssertExpectationsForObjects(t, mockAdapter)

	mockValue.On("EIsProxy").Return(true).Once()
	mockValue.On("EProxyURI").Return(nil).Once()
	assert.Equal(t, mockValue, o.GetEType())
	mock.AssertExpectationsForObjects(t, mockAdapter, mockValue)
}

func TestETypedElementETypeSet(t *testing.T) {
	o := newETypedElementImpl()
	v := new(MockEClassifier)
	mockAdapter := new(MockEAdapter)
	mockAdapter.On("SetTarget", o).Once()
	mockAdapter.On("NotifyChanged", mock.Anything).Once()
	o.EAdapters().Add(mockAdapter)
	o.SetEType(v)
	mockAdapter.AssertExpectations(t)
}

func TestETypedElementEGetFromID(t *testing.T) {
	o := newETypedElementImpl()
	assert.Panics(t, func() { o.EGetFromID(-1, true) })
	assert.Equal(t, o.IsOrdered(), o.EGetFromID(ETYPED_ELEMENT__ORDERED, true))
	assert.Equal(t, o.IsUnique(), o.EGetFromID(ETYPED_ELEMENT__UNIQUE, true))
	assert.Equal(t, o.GetEType(), o.EGetFromID(ETYPED_ELEMENT__ETYPE, true))
	assert.Panics(t, func() { o.EGetFromID(ETYPED_ELEMENT__REQUIRED, true) })
	assert.Panics(t, func() { o.EGetFromID(ETYPED_ELEMENT__REQUIRED, false) })
	assert.Equal(t, o.GetLowerBound(), o.EGetFromID(ETYPED_ELEMENT__LOWER_BOUND, true))
	assert.Panics(t, func() { o.EGetFromID(ETYPED_ELEMENT__MANY, true) })
	assert.Panics(t, func() { o.EGetFromID(ETYPED_ELEMENT__MANY, false) })
	assert.Equal(t, o.GetUpperBound(), o.EGetFromID(ETYPED_ELEMENT__UPPER_BOUND, true))
}

func TestETypedElementESetFromID(t *testing.T) {
	o := newETypedElementImpl()
	assert.Panics(t, func() { o.ESetFromID(-1, nil) })
	{
		v := new(MockEClassifier)
		o.ESetFromID(ETYPED_ELEMENT__ETYPE, v)
		assert.Equal(t, v, o.EGetFromID(ETYPED_ELEMENT__ETYPE, false))
	}
	{
		v := 45
		o.ESetFromID(ETYPED_ELEMENT__LOWER_BOUND, v)
		assert.Equal(t, v, o.EGetFromID(ETYPED_ELEMENT__LOWER_BOUND, false))
	}
	{
		v := true
		o.ESetFromID(ETYPED_ELEMENT__ORDERED, v)
		assert.Equal(t, v, o.EGetFromID(ETYPED_ELEMENT__ORDERED, false))
	}
	{
		v := true
		o.ESetFromID(ETYPED_ELEMENT__UNIQUE, v)
		assert.Equal(t, v, o.EGetFromID(ETYPED_ELEMENT__UNIQUE, false))
	}
	{
		v := 45
		o.ESetFromID(ETYPED_ELEMENT__UPPER_BOUND, v)
		assert.Equal(t, v, o.EGetFromID(ETYPED_ELEMENT__UPPER_BOUND, false))
	}

}

func TestETypedElementEIsSetFromID(t *testing.T) {
	o := newETypedElementImpl()
	assert.Panics(t, func() { o.EIsSetFromID(-1) })
	assert.False(t, o.EIsSetFromID(ETYPED_ELEMENT__ORDERED))
	assert.False(t, o.EIsSetFromID(ETYPED_ELEMENT__UNIQUE))
	assert.False(t, o.EIsSetFromID(ETYPED_ELEMENT__ETYPE))
	assert.Panics(t, func() { o.EIsSetFromID(ETYPED_ELEMENT__REQUIRED) })
	assert.False(t, o.EIsSetFromID(ETYPED_ELEMENT__LOWER_BOUND))
	assert.Panics(t, func() { o.EIsSetFromID(ETYPED_ELEMENT__MANY) })
	assert.False(t, o.EIsSetFromID(ETYPED_ELEMENT__UPPER_BOUND))
}

func TestETypedElementEUnsetFromID(t *testing.T) {
	o := newETypedElementImpl()
	assert.Panics(t, func() { o.EUnsetFromID(-1) })
	{
		o.EUnsetFromID(ETYPED_ELEMENT__ORDERED)
		v := o.EGetFromID(ETYPED_ELEMENT__ORDERED, false)
		assert.Equal(t, true, v)
	}
	{
		o.EUnsetFromID(ETYPED_ELEMENT__UNIQUE)
		v := o.EGetFromID(ETYPED_ELEMENT__UNIQUE, false)
		assert.Equal(t, true, v)
	}
	{
		o.EUnsetFromID(ETYPED_ELEMENT__ETYPE)
		assert.Nil(t, o.EGetFromID(ETYPED_ELEMENT__ETYPE, false))
	}
	{
		o.EUnsetFromID(ETYPED_ELEMENT__LOWER_BOUND)
		v := o.EGetFromID(ETYPED_ELEMENT__LOWER_BOUND, false)
		assert.Equal(t, 0, v)
	}
	{
		o.EUnsetFromID(ETYPED_ELEMENT__UPPER_BOUND)
		v := o.EGetFromID(ETYPED_ELEMENT__UPPER_BOUND, false)
		assert.Equal(t, 1, v)
	}
}
