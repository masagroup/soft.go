// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package xmi

import (
	"io"

	"github.com/masagroup/soft.go/codec/xml"
	"github.com/masagroup/soft.go/ecore"
)

// Encoder is the XMI encoder: an XML encoder writing the xmi namespace and version, without
// extended metadata.
type Encoder struct {
	*xml.Encoder
	xmiVersion string
}

func NewEncoder(resource ecore.EResource, w io.Writer, options map[string]any) *Encoder {
	s := new(Encoder)
	s.Encoder = xml.NewEncoder(resource, w, options)
	s.SetInterfaces(s)
	s.SetExtendedMetaData(nil)
	s.xmiVersion = "2.0"
	return s
}

func (s *Encoder) SetXMIVersion(xmiVersion string) {
	s.xmiVersion = xmiVersion
}

func (s *Encoder) SaveNamespaces() {
	s.AddAttribute(xmiNS+":"+versionAttrib, s.xmiVersion)
	s.AddAttribute(xmlNS+":"+xmiNS, xmiURI)
	s.Encoder.SaveNamespaces()
}
