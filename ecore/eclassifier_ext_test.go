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

func TestEClassifierExtClassifierID(t *testing.T) {
	c := newEClassifierExt()
	assert.Equal(t, -1, c.GetClassifierID())

	mockPackage := &MockEPackage{}
	mockClassifiers := &MockEList{}
	c.ESetInternalContainer(mockPackage, ECLASSIFIER__EPACKAGE)
	mockPackage.On("GetEClassifiers").Return(mockClassifiers).Once()
	mockPackage.On("EIsProxy").Return(false).Once()
	mockClassifiers.On("IndexOf", c).Return(0).Once()
	assert.Equal(t, 0, c.GetClassifierID())
	mock.AssertExpectationsForObjects(t, mockPackage, mockClassifiers)
}

func TestEClassifierExtGetDefaultValue(t *testing.T) {
	c := newEClassifierExt()
	assert.Nil(t, c.GetDefaultValue())
}
