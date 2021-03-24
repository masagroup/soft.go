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

func TestMockEObjectList_GetUnResolvedList(t *testing.T) {
	o := &MockEObjectList{}
	l := &MockEList{}
	o.On("GetUnResolvedList").Once().Return(l)
	o.On("GetUnResolvedList").Once().Return(func() EList {
		return l
	})
	assert.Equal(t, l, o.GetUnResolvedList())
	assert.Equal(t, l, o.GetUnResolvedList())
	mock.AssertExpectationsForObjects(t, o, l)
}
