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

func discardEEnum() {
	_ = assert.Equal
	_ = mock.Anything
	_ = testing.Coverage
}

func TestEEnumAsEEnum(t *testing.T) {
	o := newEEnumImpl()
	assert.Equal(t, o, o.asEEnum())
}

func TestEEnumStaticClass(t *testing.T) {
	o := newEEnumImpl()
	assert.Equal(t, GetPackage().GetEEnum(), o.EStaticClass())
}

func TestEEnumFeatureCount(t *testing.T) {
	o := newEEnumImpl()
	assert.Equal(t, EENUM_FEATURE_COUNT, o.EStaticFeatureCount())
}

func TestEEnumELiteralsGet(t *testing.T) {
	o := newEEnumImpl()
	assert.NotNil(t, o.GetELiterals())
}

func TestEEnumGetEEnumLiteralByNameOperation(t *testing.T) {
	o := newEEnumImpl()
	assert.Panics(t, func() { o.GetEEnumLiteralByName("") })
}
func TestEEnumGetEEnumLiteralByValueOperation(t *testing.T) {
	o := newEEnumImpl()
	assert.Panics(t, func() { o.GetEEnumLiteralByValue(0) })
}
func TestEEnumGetEEnumLiteralByLiteralOperation(t *testing.T) {
	o := newEEnumImpl()
	assert.Panics(t, func() { o.GetEEnumLiteralByLiteral("") })
}

func TestEEnumEGetFromID(t *testing.T) {
	o := newEEnumImpl()
	assert.Panics(t, func() { o.EGetFromID(-1, true) })
	assert.Equal(t, o.GetELiterals(), o.EGetFromID(EENUM__ELITERALS, true))
	assert.Equal(t, o.GetELiterals().(EObjectList).GetUnResolvedList(), o.EGetFromID(EENUM__ELITERALS, false))
}

func TestEEnumESetFromID(t *testing.T) {
	o := newEEnumImpl()
	assert.Panics(t, func() { o.ESetFromID(-1, nil) })
	{
		// list with a value
		mockValue := new(MockEEnumLiteral)
		l := NewImmutableEList([]interface{}{mockValue})
		// expectations
		mockValue.On("EInverseAdd", o, EENUM_LITERAL__EENUM, mock.Anything).Return(nil).Once()
		// set list with new contents
		o.ESetFromID(EENUM__ELITERALS, l)
		// checks
		assert.Equal(t, 1, o.GetELiterals().Size())
		assert.Equal(t, mockValue, o.GetELiterals().Get(0))
		mock.AssertExpectationsForObjects(t, mockValue)
	}

}

func TestEEnumEIsSetFromID(t *testing.T) {
	o := newEEnumImpl()
	assert.Panics(t, func() { o.EIsSetFromID(-1) })
	assert.False(t, o.EIsSetFromID(EENUM__ELITERALS))
}

func TestEEnumEUnsetFromID(t *testing.T) {
	o := newEEnumImpl()
	assert.Panics(t, func() { o.EUnsetFromID(-1) })
	{
		o.EUnsetFromID(EENUM__ELITERALS)
		v := o.EGetFromID(EENUM__ELITERALS, false)
		assert.NotNil(t, v)
		l := v.(EList)
		assert.True(t, l.Empty())
	}
}

func TestEEnumEInvokeFromID(t *testing.T) {
	o := newEEnumImpl()
	assert.Panics(t, func() { o.EInvokeFromID(-1, nil) })
	assert.Panics(t, func() { o.EInvokeFromID(EENUM__GET_EENUM_LITERAL_BY_LITERAL_ESTRING, nil) })
	assert.Panics(t, func() { o.EInvokeFromID(EENUM__GET_EENUM_LITERAL_EINT, nil) })
	assert.Panics(t, func() { o.EInvokeFromID(EENUM__GET_EENUM_LITERAL_ESTRING, nil) })
}
