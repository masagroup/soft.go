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
	o := &MockEDiagnostic{}
	v := "message"
	o.On("GetMessage").Once().Return(v)
	assert.Equal(t, v, o.GetMessage())

	o.On("GetMessage").Once().Return(func() string {
		return v
	})
	assert.Equal(t, v, o.GetMessage())
	o.AssertExpectations(t)
}

func TestMockEDiagnosticGetLocation(t *testing.T) {
	o := &MockEDiagnostic{}
	v := "location"
	o.On("GetLocation").Once().Return(v)
	assert.Equal(t, v, o.GetLocation())

	o.On("GetLocation").Once().Return(func() string {
		return v
	})
	assert.Equal(t, v, o.GetLocation())
	o.AssertExpectations(t)
}

func TestMockEDiagnosticGetLine(t *testing.T) {
	o := &MockEDiagnostic{}
	v := 1
	o.On("GetLine").Once().Return(v)
	assert.Equal(t, v, o.GetLine())

	o.On("GetLine").Once().Return(func() int {
		return v
	})
	assert.Equal(t, v, o.GetLine())
	o.AssertExpectations(t)
}

func TestMockEDiagnosticGetColumn(t *testing.T) {
	o := &MockEDiagnostic{}
	v := 1
	o.On("GetColumn").Once().Return(v)
	assert.Equal(t, v, o.GetColumn())

	o.On("GetColumn").Once().Return(func() int {
		return v
	})
	assert.Equal(t, v, o.GetColumn())
	o.AssertExpectations(t)
}

func TestMockEDiagnosticError(t *testing.T) {
	o := &MockEDiagnostic{}
	v := "error"
	o.On("Error").Once().Return(v)
	assert.Equal(t, v, o.Error())

	o.On("Error").Once().Return(func() string {
		return v
	})
	assert.Equal(t, v, o.Error())
	o.AssertExpectations(t)
}
