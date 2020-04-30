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

func TestEDiagnostic(t *testing.T) {
	d := NewEDiagnosticImpl("message", "location", 10, 20)
	assert.Equal(t, "message", d.GetMessage())
	assert.Equal(t, "location", d.GetLocation())
	assert.Equal(t, 10, d.GetLine())
	assert.Equal(t, 20, d.GetColumn())
}
