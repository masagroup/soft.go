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

func TestEEnumExt_GetEEnumLiteralByName(t *testing.T) {
	mockEnumLiteral := NewMockEEnumLiteral(t)
	eEnum := newEEnumExt()

	mockEnumLiteral.On("EInverseAdd", eEnum, EENUM_LITERAL__EENUM, nil).Return(nil).Once()
	eEnum.GetELiterals().Add(mockEnumLiteral)
	mock.AssertExpectationsForObjects(t, mockEnumLiteral)

	mockEnumLiteral.On("GetName").Return("mockEnumLiteral").Once()
	assert.Nil(t, eEnum.GetEEnumLiteralByName("test"))
	mock.AssertExpectationsForObjects(t, mockEnumLiteral)

	mockEnumLiteral.On("GetName").Return("mockEnumLiteral").Once()
	assert.Equal(t, mockEnumLiteral, eEnum.GetEEnumLiteralByName("mockEnumLiteral"))
	mock.AssertExpectationsForObjects(t, mockEnumLiteral)
}

func TestEEnumExt_GetEEnumLiteralByValue(t *testing.T) {
	mockEnumLiteral := NewMockEEnumLiteral(t)
	eEnum := newEEnumExt()

	mockEnumLiteral.On("EInverseAdd", eEnum, EENUM_LITERAL__EENUM, nil).Return(nil).Once()
	eEnum.GetELiterals().Add(mockEnumLiteral)
	mock.AssertExpectationsForObjects(t, mockEnumLiteral)

	mockEnumLiteral.On("GetValue").Return(0).Once()
	assert.Nil(t, eEnum.GetEEnumLiteralByValue(1))
	mock.AssertExpectationsForObjects(t, mockEnumLiteral)

	mockEnumLiteral.On("GetValue").Return(0).Once()
	assert.Equal(t, mockEnumLiteral, eEnum.GetEEnumLiteralByValue(0))
	mock.AssertExpectationsForObjects(t, mockEnumLiteral)
}

func TestEEnumExt_GetEEnumLiteralByLiteral(t *testing.T) {
	mockEnumLiteral := NewMockEEnumLiteral(t)
	eEnum := newEEnumExt()

	mockEnumLiteral.On("EInverseAdd", eEnum, EENUM_LITERAL__EENUM, nil).Return(nil).Once()
	eEnum.GetELiterals().Add(mockEnumLiteral)
	mock.AssertExpectationsForObjects(t, mockEnumLiteral)

	mockEnumLiteral.On("GetLiteral").Return("no literal").Once()
	assert.Nil(t, eEnum.GetEEnumLiteralByLiteral("literal"))
	mock.AssertExpectationsForObjects(t, mockEnumLiteral)

	mockEnumLiteral.On("GetLiteral").Return("literal").Once()
	assert.Equal(t, mockEnumLiteral, eEnum.GetEEnumLiteralByLiteral("literal"))
	mock.AssertExpectationsForObjects(t, mockEnumLiteral)
}

func TestEEnumExt_GetDefaultValue(t *testing.T) {
	eEnum := newEEnumExt()

	mockEnumLiteral := NewMockEEnumLiteral(t)
	mockEnumLiteral.On("EInverseAdd", eEnum, EENUM_LITERAL__EENUM, nil).Return(nil).Once()
	eEnum.GetELiterals().Add(mockEnumLiteral)
	mock.AssertExpectationsForObjects(t, mockEnumLiteral)

	mockValue := 2
	mockEnumLiteral.On("GetValue").Return(mockValue).Once()
	assert.Equal(t, mockValue, eEnum.GetDefaultValue())
	mock.AssertExpectationsForObjects(t, mockEnumLiteral)
}
