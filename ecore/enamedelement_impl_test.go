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

func discardENamedElement() {
	_ = assert.Equal
	_ = mock.Anything
	_ = testing.Coverage
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
