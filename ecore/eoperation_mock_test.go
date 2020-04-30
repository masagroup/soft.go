// *****************************************************************************
//
// This file is part of a MASA library or program.
// Refer to the included end-user license agreement for restrictions.
//
// Copyright (c) 2020 MASA Group
//
// *****************************************************************************

// *****************************************************************************
//
// Warning: This file was generated by soft.generator.go Generator
//
// *****************************************************************************

package ecore

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func discardMockEOperation() {
	_ = assert.Equal
	_ = testing.Coverage
}

// TestMockEOperationGetEContainingClass tests method GetEContainingClass
func TestMockEOperationGetEContainingClass(t *testing.T) {
	o := &MockEOperation{}
	r := &MockEClass{}
	o.On("GetEContainingClass").Once().Return(r)
	o.On("GetEContainingClass").Once().Return(func() EClass {
		return r
	})
	assert.Equal(t, r, o.GetEContainingClass())
	assert.Equal(t, r, o.GetEContainingClass())
	o.AssertExpectations(t)
}

// TestMockEOperationGetEExceptions tests method GetEExceptions
func TestMockEOperationGetEExceptions(t *testing.T) {
	o := &MockEOperation{}
	l := &MockEList{}
	// return a value
	o.On("GetEExceptions").Once().Return(l)
	o.On("GetEExceptions").Once().Return(func() EList {
		return l
	})
	assert.Equal(t, l, o.GetEExceptions())
	assert.Equal(t, l, o.GetEExceptions())
	o.AssertExpectations(t)
}

// TestMockEOperationUnsetEExceptions tests method UnsetEExceptions
func TestMockEOperationUnsetEExceptions(t *testing.T) {
	o := &MockEOperation{}
	o.On("UnsetEExceptions").Once()
	o.UnsetEExceptions()
	o.AssertExpectations(t)
}

// TestMockEOperationGetEParameters tests method GetEParameters
func TestMockEOperationGetEParameters(t *testing.T) {
	o := &MockEOperation{}
	l := &MockEList{}
	// return a value
	o.On("GetEParameters").Once().Return(l)
	o.On("GetEParameters").Once().Return(func() EList {
		return l
	})
	assert.Equal(t, l, o.GetEParameters())
	assert.Equal(t, l, o.GetEParameters())
	o.AssertExpectations(t)
}

// TestMockEOperationGetOperationID tests method GetOperationID
func TestMockEOperationGetOperationID(t *testing.T) {
	o := &MockEOperation{}
	r := 45
	o.On("GetOperationID").Once().Return(r)
	o.On("GetOperationID").Once().Return(func() int {
		return r
	})
	assert.Equal(t, r, o.GetOperationID())
	assert.Equal(t, r, o.GetOperationID())
	o.AssertExpectations(t)
}

// TestMockEOperationSetOperationID tests method SetOperationID
func TestMockEOperationSetOperationID(t *testing.T) {
	o := &MockEOperation{}
	v := 45
	o.On("SetOperationID", v).Once()
	o.SetOperationID(v)
	o.AssertExpectations(t)
}

// TestMockEOperationIsOverrideOf tests method IsOverrideOf
func TestMockEOperationIsOverrideOf(t *testing.T) {
	o := &MockEOperation{}
	someOperation := &MockEOperation{}
	r := true
	o.On("IsOverrideOf", someOperation).Return(r).Once()
	o.On("IsOverrideOf", someOperation).Return(func() bool {
		return r
	}).Once()
	assert.Equal(t, r, o.IsOverrideOf(someOperation))
	assert.Equal(t, r, o.IsOverrideOf(someOperation))
	o.AssertExpectations(t)
}