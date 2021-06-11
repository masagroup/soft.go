package ecore

type XMIEncoder struct {
	*XMLEncoder
	xmiVersion string
}

func NewXMIEncoder(options map[string]interface{}) *XMIEncoder {
	s := new(XMIEncoder)
	s.XMLEncoder = NewXMLEncoder(options)
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
