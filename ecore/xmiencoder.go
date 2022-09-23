// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import "io"

type XMIEncoder struct {
	*XMLEncoder
	xmiVersion string
}

func NewXMIEncoder(resource EResource, w io.Writer, options map[string]any) *XMIEncoder {
	s := new(XMIEncoder)
	s.XMLEncoder = NewXMLEncoder(resource, w, options)
	s.interfaces = s
	s.extendedMetaData = nil
	s.xmiVersion = "2.0"
	return s
}

func (s *XMIEncoder) SetXMIVersion(xmiVersion string) {
	s.xmiVersion = xmiVersion
}

func (s *XMIEncoder) saveNamespaces() {
	s.str.addAttribute(xmiNS+":"+versionAttrib, s.xmiVersion)
	s.str.addAttribute(xmlNS+":"+xmiNS, xmiURI)
	s.XMLEncoder.saveNamespaces()
}
