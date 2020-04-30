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
	"github.com/stretchr/testify/mock"
)

// TestGetELiterals tests method GetELiterals
func TestMockEIteratorHasNext(t *testing.T) {
	o := &MockEIterator{}
	v := true
	o.On("HasNext").Once().Return(v)
	assert.Equal(t, v, o.HasNext())

	o.On("HasNext").Once().Return(func() bool {
		return v
	})
	assert.Equal(t, v, o.HasNext())
	o.AssertExpectations(t)
}

func TestMockEIteratorNext(t *testing.T) {
	o := &MockEIterator{}
	v := &mock.Mock{}
	o.On("Next").Once().Return(v)
	assert.Equal(t, v, o.Next())

	o.On("Next").Once().Return(func() interface{} {
		return v
	})
	assert.Equal(t, v, o.Next())
	o.AssertExpectations(t)
}
