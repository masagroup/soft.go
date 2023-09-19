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

import "github.com/stretchr/testify/assert"
import "github.com/stretchr/testify/mock"
import "testing"

func discardEEnumLiteral() {
	_ = assert.Equal
	_ = mock.Anything
	_ = testing.Coverage
	_ = NewMockEEnumLiteral
}

func TestEEnumLiteralAsEEnumLiteral(t *testing.T) {
	o := newEEnumLiteralImpl()
	assert.Equal(t, o, o.asEEnumLiteral())
}

func TestEEnumLiteralStaticClass(t *testing.T) {
	o := newEEnumLiteralImpl()
	assert.Equal(t, GetPackage().GetEEnumLiteral(), o.EStaticClass())
}

func TestEEnumLiteralFeatureCount(t *testing.T) {
	o := newEEnumLiteralImpl()
	assert.Equal(t, EENUM_LITERAL_FEATURE_COUNT, o.EStaticFeatureCount())
}

func TestEEnumLiteralEEnumGet(t *testing.T) {
	// default
	o := newEEnumLiteralImpl()
	assert.Nil(t, o.GetEEnum())

	// set a mock container
	v := NewMockEEnum(t)
	o.ESetInternalContainer(v, EENUM_LITERAL__EENUM)

	// no proxy
	v.EXPECT().EIsProxy().Return(false).Once()
	assert.Equal(t, v, o.GetEEnum())
}

func TestEEnumLiteralInstanceGet(t *testing.T) {
	o := newEEnumLiteralImpl()
	// get default value
	assert.Nil(t, o.GetInstance())
	// get initialized value
	v := any(nil)
	o.instance = v
	assert.Equal(t, v, o.GetInstance())
}

func TestEEnumLiteralInstanceSet(t *testing.T) {
	o := newEEnumLiteralImpl()
	v := any(nil)
	mockAdapter := NewMockEAdapter(t)
	mockAdapter.EXPECT().SetTarget(o).Once()
	mockAdapter.EXPECT().NotifyChanged(mock.Anything).Once()
	o.EAdapters().Add(mockAdapter)
	o.SetInstance(v)
	mockAdapter.AssertExpectations(t)
}

func TestEEnumLiteralLiteralGet(t *testing.T) {
	o := newEEnumLiteralImpl()
	// get default value
	assert.Equal(t, string(""), o.GetLiteral())
	// get initialized value
	v := string("Test String")
	o.literal = v
	assert.Equal(t, v, o.GetLiteral())
}

func TestEEnumLiteralLiteralSet(t *testing.T) {
	o := newEEnumLiteralImpl()
	v := string("Test String")
	mockAdapter := NewMockEAdapter(t)
	mockAdapter.EXPECT().SetTarget(o).Once()
	mockAdapter.EXPECT().NotifyChanged(mock.Anything).Once()
	o.EAdapters().Add(mockAdapter)
	o.SetLiteral(v)
	mockAdapter.AssertExpectations(t)
}

func TestEEnumLiteralValueGet(t *testing.T) {
	o := newEEnumLiteralImpl()
	// get default value
	assert.Equal(t, int(0), o.GetValue())
	// get initialized value
	v := int(45)
	o.value = v
	assert.Equal(t, v, o.GetValue())
}

func TestEEnumLiteralValueSet(t *testing.T) {
	o := newEEnumLiteralImpl()
	v := int(45)
	mockAdapter := NewMockEAdapter(t)
	mockAdapter.EXPECT().SetTarget(o).Once()
	mockAdapter.EXPECT().NotifyChanged(mock.Anything).Once()
	o.EAdapters().Add(mockAdapter)
	o.SetValue(v)
	mockAdapter.AssertExpectations(t)
}

func TestEEnumLiteralEGetFromID(t *testing.T) {
	o := newEEnumLiteralImpl()
	assert.Panics(t, func() { o.EGetFromID(-1, true) })
	assert.Equal(t, o.GetEEnum(), o.EGetFromID(EENUM_LITERAL__EENUM, true))
	assert.Equal(t, o.GetInstance(), o.EGetFromID(EENUM_LITERAL__INSTANCE, true))
	assert.Equal(t, o.GetLiteral(), o.EGetFromID(EENUM_LITERAL__LITERAL, true))
	assert.Equal(t, o.GetValue(), o.EGetFromID(EENUM_LITERAL__VALUE, true))
}

func TestEEnumLiteralESetFromID(t *testing.T) {
	o := newEEnumLiteralImpl()
	assert.Panics(t, func() { o.ESetFromID(-1, nil) })
	{
		v := any(nil)
		o.ESetFromID(EENUM_LITERAL__INSTANCE, v)
		assert.Equal(t, v, o.EGetFromID(EENUM_LITERAL__INSTANCE, false))
	}
	{
		v := string("Test String")
		o.ESetFromID(EENUM_LITERAL__LITERAL, v)
		assert.Equal(t, v, o.EGetFromID(EENUM_LITERAL__LITERAL, false))
	}
	{
		v := int(45)
		o.ESetFromID(EENUM_LITERAL__VALUE, v)
		assert.Equal(t, v, o.EGetFromID(EENUM_LITERAL__VALUE, false))
	}

}

func TestEEnumLiteralEIsSetFromID(t *testing.T) {
	o := newEEnumLiteralImpl()
	assert.Panics(t, func() { o.EIsSetFromID(-1) })
	assert.False(t, o.EIsSetFromID(EENUM_LITERAL__EENUM))
	assert.False(t, o.EIsSetFromID(EENUM_LITERAL__INSTANCE))
	assert.False(t, o.EIsSetFromID(EENUM_LITERAL__LITERAL))
	assert.False(t, o.EIsSetFromID(EENUM_LITERAL__VALUE))
}

func TestEEnumLiteralEUnsetFromID(t *testing.T) {
	o := newEEnumLiteralImpl()
	assert.Panics(t, func() { o.EUnsetFromID(-1) })
	{
		o.EUnsetFromID(EENUM_LITERAL__INSTANCE)
		v := o.EGetFromID(EENUM_LITERAL__INSTANCE, false)
		assert.Nil(t, v)
	}
	{
		o.EUnsetFromID(EENUM_LITERAL__LITERAL)
		v := o.EGetFromID(EENUM_LITERAL__LITERAL, false)
		assert.Equal(t, string(""), v)
	}
	{
		o.EUnsetFromID(EENUM_LITERAL__VALUE)
		v := o.EGetFromID(EENUM_LITERAL__VALUE, false)
		assert.Equal(t, int(0), v)
	}
}

func TestEEnumLiteralEBasicInverseAdd(t *testing.T) {
	o := newEEnumLiteralImpl()
	{
		mockObject := NewMockEObject(t)
		mockNotifications := NewMockENotificationChain(t)
		assert.Equal(t, mockNotifications, o.EBasicInverseAdd(mockObject, -1, mockNotifications))
	}
	{
		mockObject := NewMockEEnum(t)
		mockObject.EXPECT().EResource().Return(nil).Once()
		mockObject.EXPECT().EIsProxy().Return(false).Once()
		o.EBasicInverseAdd(mockObject, EENUM_LITERAL__EENUM, nil)
		assert.Equal(t, mockObject, o.GetEEnum())
		mock.AssertExpectationsForObjects(t, mockObject)

		mockOther := NewMockEEnum(t)
		mockOther.EXPECT().EResource().Return(nil).Once()
		mockOther.EXPECT().EIsProxy().Return(false).Once()
		mockObject.EXPECT().EResource().Return(nil).Once()
		mockObject.EXPECT().EInverseRemove(o, EENUM__ELITERALS, nil).Return(nil).Once()
		o.EBasicInverseAdd(mockOther, EENUM_LITERAL__EENUM, nil)
		assert.Equal(t, mockOther, o.GetEEnum())
		mock.AssertExpectationsForObjects(t, mockObject, mockOther)
	}

}

func TestEEnumLiteralEBasicInverseRemove(t *testing.T) {
	o := newEEnumLiteralImpl()
	{
		mockObject := NewMockEObject(t)
		mockNotifications := NewMockENotificationChain(t)
		assert.Equal(t, mockNotifications, o.EBasicInverseRemove(mockObject, -1, mockNotifications))
	}
	{
		mockObject := NewMockEEnum(t)
		o.EBasicInverseRemove(mockObject, EENUM_LITERAL__EENUM, nil)
		mock.AssertExpectationsForObjects(t, mockObject)
	}

}
