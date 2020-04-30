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

func discardMockEAnnotation() {
	_ = assert.Equal
	_ = testing.Coverage
}

// TestGetContents tests method GetContents
func TestGetContents(t *testing.T) {
	o := &MockEAnnotation{}
	l := &MockEList{}
	o.On("GetContents").Once().Return(l)
	assert.Equal(t, l, o.GetContents())

	o.On("GetContents").Once().Return(func() EList {
		return l
	})
	assert.Equal(t, l, o.GetContents())
	o.AssertExpectations(t)
}

// TestGetDetails tests method GetDetails
func TestGetDetails(t *testing.T) {
	o := &MockEAnnotation{}
	l := &MockEList{}
	o.On("GetDetails").Once().Return(l)
	assert.Equal(t, l, o.GetDetails())

	o.On("GetDetails").Once().Return(func() EList {
		return l
	})
	assert.Equal(t, l, o.GetDetails())
	o.AssertExpectations(t)
}

// TestMockEAnnotationGetEModelElement tests method GetEModelElement
func TestMockEAnnotationGetEModelElement(t *testing.T) {
	o := &MockEAnnotation{}
	r := &MockEModelElement{}
	o.On("GetEModelElement").Once().Return(r)
	assert.Equal(t, r, o.GetEModelElement())

	o.On("GetEModelElement").Once().Return(func() EModelElement {
		return r
	})
	assert.Equal(t, r, o.GetEModelElement())
	o.AssertExpectations(t)
}

// TestMockEAnnotationSetEModelElement tests method SetEModelElement
func TestMockEAnnotationSetEModelElement(t *testing.T) {
	o := &MockEAnnotation{}
	v := &MockEModelElement{}
	o.On("SetEModelElement", v).Once()

	o.SetEModelElement(v)
	o.AssertExpectations(t)
}

// TestGetReferences tests method GetReferences
func TestGetReferences(t *testing.T) {
	o := &MockEAnnotation{}
	l := &MockEList{}
	o.On("GetReferences").Once().Return(l)
	assert.Equal(t, l, o.GetReferences())

	o.On("GetReferences").Once().Return(func() EList {
		return l
	})
	assert.Equal(t, l, o.GetReferences())
	o.AssertExpectations(t)
}

// TestMockEAnnotationGetSource tests method GetSource
func TestMockEAnnotationGetSource(t *testing.T) {
	o := &MockEAnnotation{}
	r := "Test String"
	o.On("GetSource").Once().Return(r)
	assert.Equal(t, r, o.GetSource())

	o.On("GetSource").Once().Return(func() string {
		return r
	})
	assert.Equal(t, r, o.GetSource())
	o.AssertExpectations(t)
}

// TestMockEAnnotationSetSource tests method SetSource
func TestMockEAnnotationSetSource(t *testing.T) {
	o := &MockEAnnotation{}
	v := "Test String"
	o.On("SetSource", v).Once()

	o.SetSource(v)
	o.AssertExpectations(t)
}
