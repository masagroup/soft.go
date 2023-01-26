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
)

func TestEAttributeEClass(t *testing.T) {
	assert.Equal(t, GetPackage().GetEAttribute(), GetFactory().CreateEAttribute().EClass())
}

func TestEAttribute_GetEAttributeType(t *testing.T) {
	mockType := NewMockEDataType(t)
	a := newEAttributeExt()
	a.SetEType(mockType)

	mockType.On("EIsProxy").Return(false).Once()
	assert.Equal(t, mockType, a.GetEAttributeType())
	mockType.AssertExpectations(t)
}

func TestEAttribute_BasicGetEAttributeType(t *testing.T) {
	mockType := NewMockEDataType(t)
	a := newEAttributeExt()
	a.SetEType(mockType)
	assert.Equal(t, mockType, a.basicGetEAttributeType())
}
