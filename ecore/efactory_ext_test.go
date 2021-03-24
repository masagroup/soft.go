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

func TestEFactoryExtCreate(t *testing.T) {
	f := NewEFactoryExt()
	mockClass := &MockEClass{}
	mockPackage := &MockEPackage{}
	mockClass.On("GetEPackage").Return(mockPackage).Once()
	assert.Panics(t, func() { f.Create(mockClass) })
	mock.AssertExpectationsForObjects(t, mockClass, mockPackage)
}
