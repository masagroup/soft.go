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
	d := NewMockEDiagnostic(t)
	v := "message"
	m := NewMockRun(t)
	d.EXPECT().GetMessage().Return(v).Run(func() { m.Run() }).Once()
	d.EXPECT().GetMessage().Call.Return(func() string { return v }).Once()
	assert.Equal(t, v, d.GetMessage())
	assert.Equal(t, v, d.GetMessage())
}

func TestMockEDiagnosticGetLocation(t *testing.T) {
	d := NewMockEDiagnostic(t)
	v := "location"
	m := NewMockRun(t)
	d.EXPECT().GetLocation().Return(v).Run(func() { m.Run() }).Once()
	d.EXPECT().GetLocation().Call.Return(func() string { return v }).Once()
	assert.Equal(t, v, d.GetLocation())
	assert.Equal(t, v, d.GetLocation())
}

func TestMockEDiagnosticGetLine(t *testing.T) {
	d := NewMockEDiagnostic(t)
	v := 1
	m := NewMockRun(t)
	d.EXPECT().GetLine().Return(v).Run(func() { m.Run() }).Once()
	d.EXPECT().GetLine().Call.Return(func() int { return v }).Once()
	assert.Equal(t, v, d.GetLine())
	assert.Equal(t, v, d.GetLine())
}

func TestMockEDiagnosticGetColumn(t *testing.T) {
	d := NewMockEDiagnostic(t)
	v := 1
	m := NewMockRun(t)
	d.EXPECT().GetColumn().Return(v).Run(func() { m.Run() }).Once()
	d.EXPECT().GetColumn().Call.Return(func() int { return v }).Once()
	assert.Equal(t, v, d.GetColumn())
	assert.Equal(t, v, d.GetColumn())
}

func TestMockEDiagnosticError(t *testing.T) {
	d := NewMockEDiagnostic(t)
	v := "error"
	m := NewMockRun(t)
	d.EXPECT().Error().Return(v).Run(func() { m.Run() }).Once()
	d.EXPECT().Error().Call.Return(func() string { return v }).Once()
	assert.Equal(t, v, d.Error())
	assert.Equal(t, v, d.Error())
}
