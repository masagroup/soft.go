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
	"strconv"
)

// FormatFloatXSD formats a float in the XSD lexical space: INF, -INF and NaN for the special
// values, decimal notation below 1e21 and the shortest exponent form from there.
func FormatFloatXSD(v float64, bitSize int) string {
	switch {
	case math.IsInf(v, 1):
		return "INF"
	case math.IsInf(v, -1):
		return "-INF"
	case math.IsNaN(v):
		return "NaN"
	case math.Abs(v) >= 1e21:
		return strconv.FormatFloat(v, 'g', -1, bitSize)
	default:
		return strconv.FormatFloat(v, 'f', -1, bitSize)
	}
}
