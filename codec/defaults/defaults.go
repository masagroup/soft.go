// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

// Package defaults registers the default codecs of soft.go/ecore: XMI for the "ecore" and "xmi"
// extensions, XML for "xml", binary for "bin" and SQLite for "sqlite". Import it for its side
// effects:
//
//	import _ "github.com/masagroup/soft.go/codec/defaults"
package defaults

import (
	_ "github.com/masagroup/soft.go/codec/binary"
	_ "github.com/masagroup/soft.go/codec/sqlite"
	_ "github.com/masagroup/soft.go/codec/xmi"
	_ "github.com/masagroup/soft.go/codec/xml"
)
