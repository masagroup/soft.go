// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

// Package xmlcodec implements the XML and XMI codecs of soft.go/ecore. Importing it registers the
// XMI codec for the "ecore" and "xmi" extensions and the XML codec for the "xml" extension in the
// default codec registry.
package xmlcodec

import "github.com/masagroup/soft.go/ecore"

func init() {
	ecore.RegisterCodec("ecore", &XMICodec{})
	ecore.RegisterCodec("xmi", &XMICodec{})
	ecore.RegisterCodec("xml", &XMLCodec{})
}
