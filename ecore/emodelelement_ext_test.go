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

func TestEModelElementGetAnnotation(t *testing.T) {
	m := newEModelElementExt()
	a1 := NewMockEAnnotation(t)
	a1.EXPECT().EInverseAdd(m, EANNOTATION__EMODEL_ELEMENT, nil).Return(nil)
	a1.EXPECT().GetSource().Return("a1")
	a2 := NewMockEAnnotation(t)
	a2.EXPECT().EInverseAdd(m, EANNOTATION__EMODEL_ELEMENT, nil).Return(nil)
	a2.EXPECT().GetSource().Return("a2")
	m.GetEAnnotations().Add(a1)
	m.GetEAnnotations().Add(a2)
	assert.Equal(t, a2, m.GetEAnnotation("a2"))
	assert.Equal(t, nil, m.GetEAnnotation("a"))
	mock.AssertExpectationsForObjects(t, a1, a2)
}
