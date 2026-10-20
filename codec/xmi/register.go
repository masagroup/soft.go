// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

// Package xmi implements the XMI codec of soft.go/ecore, built on the XML codec of package xml.
// Importing it registers the codec for the "ecore" and "xmi" extensions in the default codec
// registry.
package xmi

import "github.com/masagroup/soft.go/ecore"

func init() {
	ecore.RegisterCodec("ecore", &Codec{})
	ecore.RegisterCodec("xmi", &Codec{})
}
