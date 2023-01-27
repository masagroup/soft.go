// *****************************************************************************
//
// This file is part of a MASA library or program.
// Refer to the included end-user license agreement for restrictions.
//
// Copyright (c) 2020 MASA Group
//
// *****************************************************************************

package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetELiterals tests method GetELiterals
func TestMockEDiagnosticGetMessage(t *testing.T) {
	o := NewMockEDiagnostic(t)
	v := "message"
	o.EXPECT().GetMessage().Once().Return(v)
	assert.Equal(t, v, o.GetMessage())

	o.EXPECT().GetMessage().Once().Return(func() string {
		return v
	})
	assert.Equal(t, v, o.GetMessage())
	o.AssertExpectations(t)
}

func TestMockEDiagnosticGetLocation(t *testing.T) {
	o := NewMockEDiagnostic(t)
	v := "location"
	o.EXPECT().GetLocation().Once().Return(v)
	assert.Equal(t, v, o.GetLocation())

	o.EXPECT().GetLocation().Once().Return(func() string {
		return v
	})
	assert.Equal(t, v, o.GetLocation())
	o.AssertExpectations(t)
}

func TestMockEDiagnosticGetLine(t *testing.T) {
	o := NewMockEDiagnostic(t)
	v := 1
	o.EXPECT().GetLine().Once().Return(v)
	assert.Equal(t, v, o.GetLine())

	o.EXPECT().GetLine().Once().Return(func() int {
		return v
	})
	assert.Equal(t, v, o.GetLine())
	o.AssertExpectations(t)
}

func TestMockEDiagnosticGetColumn(t *testing.T) {
	o := NewMockEDiagnostic(t)
	v := 1
	o.EXPECT().GetColumn().Once().Return(v)
	assert.Equal(t, v, o.GetColumn())

	o.EXPECT().GetColumn().Once().Return(func() int {
		return v
	})
	assert.Equal(t, v, o.GetColumn())
	o.AssertExpectations(t)
}

func TestMockEDiagnosticError(t *testing.T) {
	o := NewMockEDiagnostic(t)
	v := "error"
	o.EXPECT().Error().Once().Return(v)
	assert.Equal(t, v, o.Error())

	o.EXPECT().Error().Once().Return(func() string {
		return v
	})
	assert.Equal(t, v, o.Error())
	o.AssertExpectations(t)
}
