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
	encodingxml "encoding/xml"
	"io"

	"github.com/masagroup/soft.go/codec/xml"
	"github.com/masagroup/soft.go/ecore"
)

const (
	xmiURI        = "http://www.omg.org/XMI"
	xmiNS         = "xmi"
	xmlNS         = "xmlns"
	typeAttrib    = "type"
	versionAttrib = "version"
	uuidAttrib    = "uuid"
)

// Decoder is the XMI decoder: an XML decoder reading the xmi:type, xmi:version and xmi:uuid
// attributes, without extended metadata.
type Decoder struct {
	*xml.Decoder
	xmiVersion string
}

func NewDecoder(resource ecore.EResource, r io.Reader, options map[string]any) *Decoder {
	l := new(Decoder)
	l.Decoder = xml.NewDecoder(resource, r, options)
	l.AddNotFeatures(
		encodingxml.Name{Space: xmiURI, Local: typeAttrib},
		encodingxml.Name{Space: xmiURI, Local: versionAttrib},
		encodingxml.Name{Space: xmiURI, Local: uuidAttrib},
	)
	l.SetExtendedMetaData(nil)
	l.SetInterfaces(l)
	return l
}

func (l *Decoder) GetXMIVersion() string {
	return l.xmiVersion
}

func (l *Decoder) GetXSIType() string {
	xsiType := l.Decoder.GetXSIType()
	if len(xsiType) == 0 {
		return l.GetAttributeValue(xmiURI, typeAttrib)
	}
	return xsiType
}

func (l *Decoder) HandleAttributes(object ecore.EObject) {
	version := l.GetAttributeValue(xmiURI, versionAttrib)
	if len(version) > 0 {
		l.xmiVersion = version
	}
	l.Decoder.HandleAttributes(object)
}
