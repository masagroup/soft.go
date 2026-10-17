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
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatFloatXSD(t *testing.T) {
	cases := []struct {
		v    float64
		want string
	}{
		// Very large magnitudes: shortest exponent form, not a 309-digit integer.
		{math.MaxFloat64, "1.7976931348623157e+308"},
		{1e21, "1e+21"},
		{-5e21, "-5e+21"},
		// Small values stay decimal: a plain 'g' would write 1e-05.
		{0.00001, "0.00001"},
		{0.1, "0.1"},
		// Mid range: decimal, as before.
		{3.5, "3.5"},
		{100, "100"},
		{1e7, "10000000"},
		{1e20, "100000000000000000000"},
		{-2.5, "-2.5"},
		{0, "0"},
		// XSD special values.
		{math.Inf(1), "INF"},
		{math.Inf(-1), "-INF"},
		{math.NaN(), "NaN"},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, FormatFloatXSD(c.v, 64), "FormatFloatXSD(%g)", c.v)
	}
}

// The ecore factory writes EDouble, EFloat and EBigDecimal special values in XSD form.
func TestEcoreFactoryConvertsSpecialFloatsToXSD(t *testing.T) {
	factory, pkg := GetFactory(), GetPackage()
	assert.Equal(t, "INF", factory.ConvertToString(pkg.GetEDouble(), math.Inf(1)))
	assert.Equal(t, "-INF", factory.ConvertToString(pkg.GetEDouble(), math.Inf(-1)))
	assert.Equal(t, "NaN", factory.ConvertToString(pkg.GetEDouble(), math.NaN()))
	assert.Equal(t, "INF", factory.ConvertToString(pkg.GetEFloat(), float32(math.Inf(1))))
	assert.Equal(t, "-INF", factory.ConvertToString(pkg.GetEBigDecimal(), math.Inf(-1)))
	assert.Equal(t, "INF", factory.ConvertToString(pkg.GetEDoubleObject(), math.Inf(1)))
	assert.Equal(t, "NaN", factory.ConvertToString(pkg.GetEFloatObject(), float32(math.NaN())))
	assert.Equal(t, "1.7976931348623157e+308", factory.ConvertToString(pkg.GetEDouble(), math.MaxFloat64))
}
