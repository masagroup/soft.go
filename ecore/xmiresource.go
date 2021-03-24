// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import "encoding/xml"

const (
	xmiURI        = "http://www.omg.org/XMI"
	xmiNS         = "xmi"
	versionAttrib = "version"
	uuidAttrib    = "uuid"
)

type xmiLoadImpl struct {
	*xmlLoadImpl
}

func newXMILoadImpl(options map[string]interface{}) *xmiLoadImpl {
	l := new(xmiLoadImpl)
	l.xmlLoadImpl = newXMLLoadImpl(options)
	l.notFeatures = append(l.notFeatures, xml.Name{Space: xmiURI, Local: typeAttrib}, xml.Name{Space: xmiURI, Local: versionAttrib}, xml.Name{Space: xmiURI, Local: uuidAttrib})
	l.extendedMetaData = nil
	l.interfaces = l
	return l
}

func (l *xmiLoadImpl) getXSIType() string {
	xsiType := l.xmlLoadImpl.getXSIType()
	if len(xsiType) == 0 && l.attributes != nil {
		return l.getAttributeValue(xmiURI, typeAttrib)
	}
	return xsiType
}

func (l *xmiLoadImpl) handleAttributes(object EObject) {
	version := l.getAttributeValue(xmiURI, versionAttrib)
	if len(version) > 0 {
		l.resource.(xmiResource).SetXMIVersion(version)
	}
	l.xmlLoadImpl.handleAttributes(object)
}

type xmiSaveImpl struct {
	*xmlSaveImpl
}

func newXMISaveImpl(options map[string]interface{}) *xmiSaveImpl {
	s := new(xmiSaveImpl)
	s.xmlSaveImpl = newXMLSaveImpl(options)
	s.interfaces = s
	s.extendedMetaData = nil
	return s
}

func (s *xmiSaveImpl) saveNamespaces() {
	s.str.addAttribute(xmiNS+":"+versionAttrib, s.resource.(xmiResource).GetXMIVersion())
	s.str.addAttribute(xmlNS+":"+xmiNS, xmiURI)
	s.xmlSaveImpl.saveNamespaces()
}

type xmiResource interface {
	xmlResource

	SetXMIVersion(version string)
	GetXMIVersion() string
}

type xmiResourceImpl struct {
	xmlResourceImpl
	xmiVersion string
}

func newXMIResourceImpl() *xmiResourceImpl {
	r := new(xmiResourceImpl)
	r.SetInterfaces(r)
	r.Initialize()
	return r
}

func (r *xmiResourceImpl) SetXMIVersion(xmiVersion string) {
	r.xmiVersion = xmiVersion
}

func (r *xmiResourceImpl) GetXMIVersion() string {
	return r.xmiVersion
}

func (r *xmiResourceImpl) createLoad(options map[string]interface{}) xmlLoad {
	return newXMILoadImpl(options)
}

func (r *xmiResourceImpl) createSave(options map[string]interface{}) xmlSave {
	return newXMISaveImpl(options)
}
