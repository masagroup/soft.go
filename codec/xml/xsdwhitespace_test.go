// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package xml

import (
	"strconv"
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A whitespace-padded number is not merely mis-parsed by the create*FromString family, it is
// silently zeroed - the strconv error is discarded and the zero value returned. Pin the reason the
// decoder no longer hands them the padding.
func TestXSDCollapse_PaddedNumberWouldOtherwiseBecomeZero(t *testing.T) {
	_, err := strconv.ParseFloat(" 10", 64)
	require.Error(t, err, "the premise: Go rejects a padded float outright, and the factory turns that into 0")

	p := ecore.GetPackage()
	assert.Equal(t, "10", xsdCollapse(p.GetEDouble(), " 10"))
	assert.Equal(t, "350", xsdCollapse(p.GetEFloat(), "\t350\n"))
	assert.Equal(t, "8", xsdCollapse(p.GetEInt(), " 8 "))
	assert.Equal(t, "true", xsdCollapse(p.GetEBoolean(), "\r\ntrue"))
}

// xs:string carries whiteSpace="preserve": its padding is part of the value and trimming it would
// be the very kind of silent edit this change exists to stop.
func TestXSDCollapse_LeavesAStringAlone(t *testing.T) {
	p := ecore.GetPackage()
	assert.Equal(t, " padded ", xsdCollapse(p.GetEString(), " padded "))
	assert.False(t, xsdCollapsesWhitespace(p.GetEString()))
	assert.True(t, xsdCollapsesWhitespace(p.GetEDouble()))
}

// An enumeration's lexical space is a set of NMTOKENs, so it collapses too - and a padded literal
// looked up with its padding does not merely mis-read, it panics in CreateFromString.
func TestXSDCollapse_CollapsesAnEnumLiteral(t *testing.T) {
	eEnum := ecore.GetFactory().CreateEEnum()
	eEnum.SetName("Level")
	assert.True(t, xsdCollapsesWhitespace(eEnum))
	assert.Equal(t, "detection", xsdCollapse(eEnum, " detection "))
}
