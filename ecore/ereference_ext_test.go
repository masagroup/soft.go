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

func TestEReferenceExt_GetEReferenceType(t *testing.T) {
	r := newEReferenceExt()
	mockType := NewMockEClass(t)
	mockType.EXPECT().EIsProxy().Return(false).Once()
	r.SetEType(mockType)
	assert.Equal(t, mockType, r.GetEReferenceType())
	mock.AssertExpectationsForObjects(t, mockType)
}

func TestEReferenceExt_basicGetEReferenceType(t *testing.T) {
	r := newEReferenceExt()
	mockType := NewMockEClass(t)
	r.SetEType(mockType)
	assert.Equal(t, mockType, r.basicGetEReferenceType())
	mock.AssertExpectationsForObjects(t, mockType)
}
