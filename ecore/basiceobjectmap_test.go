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

func TestBasicEObjectMap_Constructor(t *testing.T) {
	mockClass := &MockEClass{}
	m := NewBasicEObjectMap(mockClass)
	assert.NotNil(t, m)

	var mp EMap = m
	assert.NotNil(t, mp)

	var ml EList = m
	assert.NotNil(t, ml)
}

type MockEObjectEMapEntry struct {
	MockEObject
	MockEMapEntry
}

func TestBasicEObjectMap_Put(t *testing.T) {
	mockClass := &MockEClass{}
	mockPackage := &MockEPackage{}
	mockFactory := &MockEFactory{}
	mockEntry := &MockEObjectEMapEntry{}
	m := NewBasicEObjectMap(mockClass)

	mockClass.On("GetEPackage").Once().Return(mockPackage)
	mockPackage.On("GetEFactoryInstance").Once().Return(mockFactory)
	mockFactory.On("Create", mockClass).Once().Return(mockEntry)
	mockEntry.On("SetKey", 2).Once()
	mockEntry.On("SetValue", "2").Once()
	mockEntry.On("GetKey").Once().Return(2)
	mockEntry.On("GetValue").Once().Return("2")
	m.Put(2, "2")
	mock.AssertExpectationsForObjects(t, mockClass, mockPackage, mockFactory, mockEntry)
}
