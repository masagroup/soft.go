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
	o := NewMockEObjectList(t)
	l := NewMockEList(t)
	m := NewMockRun(t)
	o.EXPECT().GetUnResolvedList().Return(l).Run(func() { m.Run() }).Once()
	o.EXPECT().GetUnResolvedList().Call.Return(func() EList { return l }).Once()
	assert.Equal(t, l, o.GetUnResolvedList())
	assert.Equal(t, l, o.GetUnResolvedList())
	mock.AssertExpectationsForObjects(t, o, l)
}
