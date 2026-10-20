// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

// Package sqlite implements the SQLite codec and store of soft.go/ecore. Importing it registers
// the codec for the "sqlite" extension in the default codec registry.
package sqlite

import "github.com/masagroup/soft.go/ecore"

func init() {
	ecore.RegisterCodec("sqlite", &Codec{})
}
