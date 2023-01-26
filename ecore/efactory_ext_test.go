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

func TestEFactoryExtCreate(t *testing.T) {
	f := newEFactoryExt()
	mockClass := NewMockEClass(t)
	mockPackage := NewMockEPackage(t)
	mockClass.EXPECT().GetEPackage().Return(mockPackage).Once()
	mockClass.EXPECT().GetName().Return("mockClass").Once()
	assert.Panics(t, func() { f.Create(mockClass) })
}
