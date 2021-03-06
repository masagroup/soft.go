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

func discardENamedElement() {
	_ = assert.Equal
	_ = mock.Anything
	_ = testing.Coverage
	_ = url.Parse
}

func TestENamedElementAsENamedElement(t *testing.T) {
	o := newENamedElementImpl()
	assert.Equal(t, o, o.asENamedElement())
}

func TestENamedElementStaticClass(t *testing.T) {
	o := newENamedElementImpl()
	assert.Equal(t, GetPackage().GetENamedElement(), o.EStaticClass())
}

func TestENamedElementFeatureCount(t *testing.T) {
	o := newENamedElementImpl()
	assert.Equal(t, ENAMED_ELEMENT_FEATURE_COUNT, o.EStaticFeatureCount())
}

func TestENamedElementNameGet(t *testing.T) {
	o := newENamedElementImpl()
	// get default value
	assert.Equal(t, string(""), o.GetName())
	// get initialized value
	v := string("Test String")
	o.name = v
	assert.Equal(t, v, o.GetName())
}

func TestENamedElementNameSet(t *testing.T) {
	o := newENamedElementImpl()
	v := string("Test String")
	mockAdapter := new(MockEAdapter)
	mockAdapter.On("SetTarget", o).Once()
	mockAdapter.On("NotifyChanged", mock.Anything).Once()
	o.EAdapters().Add(mockAdapter)
	o.SetName(v)
	mockAdapter.AssertExpectations(t)
}

func TestENamedElementEGetFromID(t *testing.T) {
	o := newENamedElementImpl()
	assert.Panics(t, func() { o.EGetFromID(-1, true) })
	assert.Equal(t, o.GetName(), o.EGetFromID(ENAMED_ELEMENT__NAME, true))
}

func TestENamedElementESetFromID(t *testing.T) {
	o := newENamedElementImpl()
	assert.Panics(t, func() { o.ESetFromID(-1, nil) })
	{
		v := string("Test String")
		o.ESetFromID(ENAMED_ELEMENT__NAME, v)
		assert.Equal(t, v, o.EGetFromID(ENAMED_ELEMENT__NAME, false))
	}

}

func TestENamedElementEIsSetFromID(t *testing.T) {
	o := newENamedElementImpl()
	assert.Panics(t, func() { o.EIsSetFromID(-1) })
	assert.False(t, o.EIsSetFromID(ENAMED_ELEMENT__NAME))
}

func TestENamedElementEUnsetFromID(t *testing.T) {
	o := newENamedElementImpl()
	assert.Panics(t, func() { o.EUnsetFromID(-1) })
	{
		o.EUnsetFromID(ENAMED_ELEMENT__NAME)
		v := o.EGetFromID(ENAMED_ELEMENT__NAME, false)
		assert.Equal(t, string(""), v)
	}
}
