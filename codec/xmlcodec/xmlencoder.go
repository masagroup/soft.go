// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package xmlcodec

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/masagroup/soft.go/ecore"
	"golang.org/x/net/html"
)

type xmlSaveFeatureKind int

const (
	xsfkTransient xmlSaveFeatureKind = iota
	xsfkDataTypeSingle
	xsfkDataTypeElementSingle
	xsfkDataTypeContentSingle
	xsfkDataTypeSingleNillable
	xsfkDataTypeMany
	xsfkObjectContainSingle
	xsfkObjectContainMany
	xsfkObjectHrefSingle
	xsfkObjectHrefMany
	xsfkObjectContainSingleUnsettable
	xsfkObjectContainManyUnsettable
	xsfkObjectHrefSingleUnsettable
	xsfkObjectHrefManyUnsettable
	xsfkObjectElementSingle
	xsfkObjectElementSingleUnsettable
	xsfkObjectElementMany
	xsfkObjectElementIDRefSingle
	xsfkObjectElementIDRefSingleUnsettable
	xsfkObjectElementIDRefMany
	xsfkAttributeFeatureMap
	xsfkElementFeatureMap
	xsfkObjectAttributeSingle
	xsfkObjectAttributeMany
	xsfkObjectAttributeIDRefSingle
	xsfkObjectAttributeIDRefMany
	xsfkDataTypeAttributeMany
)

type xmlEncoderInternal interface {
	saveNamespaces()
}

type XMLEncoder struct {
	interfaces       any
	w                io.Writer
	resource         ecore.EResource
	roots            ecore.EList
	prefixesToURI    map[string]string
	uriToPrefixes    map[string][]string
	packages         map[ecore.EPackage]string
	featureKinds     map[ecore.EStructuralFeature]xmlSaveFeatureKind
	extendedMetaData *ExtendedMetaData
	str              *xmlString
	errorFn          func(diagnostic ecore.EDiagnostic)
	idAttributeName  string
	xmlVersion       string
	encoding         string
	keepDefaults     bool
}

func NewXMLEncoder(resource ecore.EResource, w io.Writer, options map[string]any) *XMLEncoder {
	s := new(XMLEncoder)
	s.interfaces = s
	s.resource = resource
	s.w = w
	s.xmlVersion = "1.0"
	s.encoding = "UTF-8"
	s.str = newXmlString()
	s.packages = make(map[ecore.EPackage]string)
	s.uriToPrefixes = make(map[string][]string)
	s.prefixesToURI = make(map[string]string)
	s.featureKinds = make(map[ecore.EStructuralFeature]xmlSaveFeatureKind)
	if options != nil {
		s.idAttributeName, _ = options[XML_OPTION_ID_ATTRIBUTE_NAME].(string)
		s.roots, _ = options[XML_OPTION_ROOT_OBJECTS].(ecore.EList)
		if extendedMetaData := options[XML_OPTION_EXTENDED_META_DATA]; extendedMetaData != nil {
			s.extendedMetaData = extendedMetaData.(*ExtendedMetaData)
		}
	}
	if s.extendedMetaData == nil {
		s.extendedMetaData = NewExtendedMetaData()
	}
	return s
}

func (s *XMLEncoder) SetEncoding(encoding string) {
	s.encoding = encoding
}

func (s *XMLEncoder) SetXMLVersion(xmlVersion string) {
	s.xmlVersion = xmlVersion
}

func (s *XMLEncoder) EncodeResource() {
	s.errorFn = func(diagnostic ecore.EDiagnostic) {
		s.resource.GetErrors().Add(diagnostic)
	}
	contents := s.roots
	if contents == nil {
		contents = s.resource.GetContents()
	}
	if contents.Empty() {
		return
	}
	s.encodeTopObject(contents.Get(0).(ecore.EObject))
}

func (s *XMLEncoder) EncodeObject(eObject ecore.EObject) (err error) {
	s.errorFn = func(diagnostic ecore.EDiagnostic) {
		if err == nil {
			err = diagnostic
		}
	}
	s.encodeTopObject(eObject)
	return
}

func (s *XMLEncoder) encodeTopObject(eObject ecore.EObject) {
	s.saveHeader()

	// initialize prefixes if any in top
	if s.extendedMetaData != nil {
		eClass := eObject.EClass()
		if ePrefixMapFeature := s.extendedMetaData.GetXMLNSPrefixMapFeature(eClass); ePrefixMapFeature != nil {
			m := eObject.EGet(ePrefixMapFeature).(ecore.EMap)
			s.setPrefixToNamespace(m)
		}
	}

	s.saveTopObject(eObject)

	// namespaces
	s.str.resetToFirstElementMark()
	s.interfaces.(xmlEncoderInternal).saveNamespaces()

	// write result
	if err := s.str.write(s.w); err != nil {
		s.error(ecore.NewEDiagnosticImpl(err.Error(), s.resource.GetURI().String(), 0, 0))
	}
}

func (s *XMLEncoder) saveHeader() {
	s.str.add(fmt.Sprintf("<?xml version=\"%v\" encoding=\"%v\"?>", s.xmlVersion, s.encoding))
	s.str.addLine()
}

func (s *XMLEncoder) saveTopObject(eObject ecore.EObject) {
	eClass := eObject.EClass()
	if s.extendedMetaData == nil || s.extendedMetaData.GetDocumentRoot(eClass.GetEPackage()) != eClass {
		var name string
		if rootFeature := s.getRootFeature(eClass); rootFeature != nil {
			name = s.getFeatureQName(rootFeature)
		} else {
			name = s.getClassQName(eClass)
		}
		s.str.startElement(name)
		s.saveElementID(eObject)
		s.saveFeatures(eObject, false)
	} else {
		s.str.startElement("")
		s.saveFeatures(eObject, false)
	}
}

func (s *XMLEncoder) getRootFeature(eClassifier ecore.EClassifier) ecore.EStructuralFeature {
	if s.extendedMetaData != nil {
		for eClassifier != nil {
			if eClass := s.extendedMetaData.GetDocumentRoot(eClassifier.GetEPackage()); eClass != nil {
				for it := eClass.GetEStructuralFeatures().Iterator(); it.HasNext(); {
					eFeature := it.Next().(ecore.EStructuralFeature)
					if eFeature.GetEType() == eClassifier && eFeature.IsChangeable() {
						return eFeature
					}
				}
			}
			if eClass, _ := eClassifier.(ecore.EClass); eClass != nil {
				eSuperTypes := eClass.GetESuperTypes()
				if eSuperTypes.Empty() {
					eClassifier = nil
				} else {
					eClassifier = eSuperTypes.Get(0).(ecore.EClass)
				}
			} else {
				eClassifier = nil
			}
		}
	}
	return nil
}

func (s *XMLEncoder) saveNamespaces() {
	prefixes := make([]string, 0, len(s.prefixesToURI))
	for prefix := range s.prefixesToURI {
		prefixes = append(prefixes, prefix)
	}
	sort.Strings(prefixes)
	for _, prefix := range prefixes {
		attribute := "xmlns"
		if len(prefix) > 0 {
			attribute += ":" + prefix
		}
		s.str.addAttribute(attribute, s.prefixesToURI[prefix])
	}
}

func (s *XMLEncoder) saveElementID(eObject ecore.EObject) {
	if idManager := s.resource.GetObjectIDManager(); len(s.idAttributeName) > 0 && idManager != nil {
		if id := idManager.GetID(eObject); id != nil {
			s.str.addAttribute(s.idAttributeName, fmt.Sprintf("%v", id))
		}
	}
}

func (s *XMLEncoder) saveFeatures(eObject ecore.EObject, attributesOnly bool) bool {
	eClass := eObject.EClass()
	eAllFeatures := eClass.GetEAllStructuralFeatures()
	var elementFeatures []int
	elementCount := 0
	i := 0
	for it := eAllFeatures.Iterator(); it.HasNext(); i++ {
		// current feature
		eFeature := it.Next().(ecore.EStructuralFeature)
		// compute feature kind
		kind, ok := s.featureKinds[eFeature]
		if !ok {
			kind = s.getSaveFeatureKind(eFeature)
			s.featureKinds[eFeature] = kind
		}

		if kind != xsfkTransient && s.shouldSaveFeature(eObject, eFeature) {
			switch kind {
			case xsfkDataTypeSingle:
				s.saveDataTypeSingle(eObject, eFeature)
				continue
			case xsfkDataTypeSingleNillable:
				if !s.isNil(eObject, eFeature) {
					s.saveDataTypeSingle(eObject, eFeature)
					continue
				}
			case xsfkObjectContainManyUnsettable:
				fallthrough
			case xsfkDataTypeMany:
				if s.isEmpty(eObject, eFeature) {
					s.saveManyEmpty(eObject, eFeature)
					continue
				}
			case xsfkObjectContainSingleUnsettable:
			case xsfkObjectContainSingle:
			case xsfkObjectContainMany:
			case xsfkObjectHrefSingleUnsettable:
				if !s.isNil(eObject, eFeature) {
					switch s.getSaveResourceKindSingle(eObject, eFeature) {
					case cross:
					case same:
						s.saveIDRefSingle(eObject, eFeature)
						continue
					default:
						continue
					}
				}
			case xsfkObjectHrefSingle:
				switch s.getSaveResourceKindSingle(eObject, eFeature) {
				case cross:
				case same:
					s.saveIDRefSingle(eObject, eFeature)
					continue
				default:
					continue
				}
			case xsfkObjectHrefManyUnsettable:
				if s.isEmpty(eObject, eFeature) {
					s.saveManyEmpty(eObject, eFeature)
					continue
				} else {
					switch s.getSaveResourceKindMany(eObject, eFeature) {
					case cross:
					case same:
						s.saveIDRefMany(eObject, eFeature)
						continue
					default:
						continue
					}
				}

			case xsfkObjectHrefMany:
				switch s.getSaveResourceKindMany(eObject, eFeature) {
				case cross:
				case same:
					s.saveIDRefMany(eObject, eFeature)
					continue
				default:
					continue
				}
			default:
				continue
			}
			if attributesOnly {
				continue
			}
			if elementFeatures == nil {
				elementFeatures = make([]int, eAllFeatures.Size())
			}
			elementFeatures[elementCount] = i
			elementCount++
		}
	}
	if elementFeatures == nil {
		s.str.endEmptyElement()
		return false
	}
	for i := 0; i < elementCount; i++ {
		eFeature := eAllFeatures.Get(elementFeatures[i]).(ecore.EStructuralFeature)
		kind := s.featureKinds[eFeature]
		switch kind {
		case xsfkDataTypeSingleNillable:
			s.saveNil(eObject, eFeature)
		case xsfkDataTypeMany:
			s.saveDataTypeMany(eObject, eFeature)
		case xsfkObjectContainSingleUnsettable:
			if s.isNil(eObject, eFeature) {
				s.saveNil(eObject, eFeature)
			} else {
				s.saveContainedSingle(eObject, eFeature)
			}
		case xsfkObjectContainSingle:
			s.saveContainedSingle(eObject, eFeature)
		case xsfkObjectContainManyUnsettable:
			fallthrough
		case xsfkObjectContainMany:
			s.saveContainedMany(eObject, eFeature)
		case xsfkObjectHrefSingleUnsettable:
			if s.isNil(eObject, eFeature) {
				s.saveNil(eObject, eFeature)
			} else {
				s.saveHRefSingle(eObject, eFeature)
			}
		case xsfkObjectHrefSingle:
			s.saveHRefSingle(eObject, eFeature)
		case xsfkObjectHrefManyUnsettable:
			fallthrough
		case xsfkObjectHrefMany:
			s.saveHRefMany(eObject, eFeature)
		}
	}

	s.str.endElement()
	return true
}

func (s *XMLEncoder) saveDataTypeSingle(eObject ecore.EObject, eFeature ecore.EStructuralFeature) {
	val := eObject.EGetResolve(eFeature, false)
	str, ok := s.getDataType(val, eFeature, true)
	if ok {
		s.str.addAttribute(s.getFeatureQName(eFeature), html.EscapeString(str))
	}
}

func (s *XMLEncoder) saveDataTypeMany(eObject ecore.EObject, eFeature ecore.EStructuralFeature) {
	l := eObject.EGetResolve(eFeature, false).(ecore.EList)
	d := eFeature.GetEType().(ecore.EDataType)
	p := d.GetEPackage()
	f := p.GetEFactoryInstance()
	name := s.getFeatureQName(eFeature)
	for it := l.Iterator(); it.HasNext(); {
		value := it.Next()
		if value == nil {
			s.str.startElement(name)
			s.str.addAttribute("xsi:nil", "true")
			s.str.endEmptyElement()
			s.uriToPrefixes[xsiURI] = []string{xsiNS}
			s.prefixesToURI[xsiNS] = xsiURI
		} else {
			str := f.ConvertToString(d, value)
			s.str.addContent(name, str)
		}
	}
}

func (s *XMLEncoder) saveManyEmpty(_ ecore.EObject, eFeature ecore.EStructuralFeature) {
	s.str.addAttribute(s.getFeatureQName(eFeature), "")
}

//lint:ignore U1000 func (*XMLEncoder).saveEObjectSingle is unused
func (s *XMLEncoder) saveEObjectSingle(eObject ecore.EObject, eFeature ecore.EStructuralFeature) {
	value, _ := eObject.EGetResolve(eFeature, false).(ecore.EObject)
	if value != nil {
		id := s.getHRef(value)
		s.str.addAttribute(s.getFeatureQName(eFeature), id)
	}
}

//lint:ignore U1000 func (*XMLEncoder).saveEObjectMany is unused
func (s *XMLEncoder) saveEObjectMany(eObject ecore.EObject, eFeature ecore.EStructuralFeature) {
	l := eObject.EGetResolve(eFeature, false).(ecore.EList)
	failure := false
	var buffer strings.Builder
	for it := l.Iterator(); ; {
		value, _ := it.Next().(ecore.EObject)
		if value != nil {
			id := s.getHRef(value)
			if id == "" {
				failure = true
				if !it.HasNext() {
					break
				}
			} else {
				buffer.WriteString(id)
				if it.HasNext() {
					buffer.WriteString(" ")
				} else {
					break
				}
			}
		}
	}
	if !failure && buffer.Len() > 0 {
		s.str.addAttribute(s.getFeatureQName(eFeature), buffer.String())
	}
}

func (s *XMLEncoder) saveNil(_ ecore.EObject, eFeature ecore.EStructuralFeature) {
	s.str.addNil(s.getFeatureQName(eFeature))
}

func (s *XMLEncoder) saveContainedSingle(eObject ecore.EObject, eFeature ecore.EStructuralFeature) {
	value, _ := eObject.EGetResolve(eFeature, false).(ecore.EObjectInternal)
	if value != nil {
		s.saveEObjectInternal(value, eFeature)
	}
}

func (s *XMLEncoder) saveContainedMany(eObject ecore.EObject, eFeature ecore.EStructuralFeature) {
	l := eObject.EGetResolve(eFeature, false).(ecore.EList)
	for it := l.Iterator(); it.HasNext(); {
		value, _ := it.Next().(ecore.EObjectInternal)
		if value != nil {
			s.saveEObjectInternal(value, eFeature)
		}
	}
}

func (s *XMLEncoder) saveEObjectInternal(o ecore.EObjectInternal, f ecore.EStructuralFeature) {
	if o.EInternalResource() != nil || o.EIsProxy() {
		s.saveHRef(o, f)
	} else {
		s.saveEObject(o, f)
	}
}

func (s *XMLEncoder) saveEObject(o ecore.EObject, f ecore.EStructuralFeature) {
	s.str.startElement(s.getFeatureQName(f))
	eClass := o.EClass()
	eType := f.GetEType()
	if eType != eClass && eType != ecore.GetPackage().GetEObject() {
		s.saveTypeAttribute(eClass)
	}
	s.saveElementID(o)
	s.saveFeatures(o, false)
}

func (s *XMLEncoder) saveTypeAttribute(eClass ecore.EClass) {
	s.str.addAttribute("xsi:type", s.getClassQName(eClass))
	s.uriToPrefixes[xsiURI] = []string{xsiNS}
	s.prefixesToURI[xsiNS] = xsiURI
}

func (s *XMLEncoder) saveHRefSingle(eObject ecore.EObject, eFeature ecore.EStructuralFeature) {
	value, _ := eObject.EGetResolve(eFeature, false).(ecore.EObject)
	if value != nil {
		s.saveHRef(value, eFeature)
	}
}

func (s *XMLEncoder) saveHRefMany(eObject ecore.EObject, eFeature ecore.EStructuralFeature) {
	l := eObject.EGetResolve(eFeature, false).(ecore.EList)
	for it := l.Iterator(); it.HasNext(); {
		value, _ := it.Next().(ecore.EObject)
		if value != nil {
			s.saveHRef(value, eFeature)
		}
	}
}

func (s *XMLEncoder) saveHRef(eObject ecore.EObject, eFeature ecore.EStructuralFeature) {
	href := s.getHRef(eObject)
	if href != "" {
		s.str.startElement(s.getFeatureQName(eFeature))
		eClass := eObject.EClass()
		eType, _ := eFeature.GetEType().(ecore.EClass)
		if eType != eClass && eType != nil && eType.IsAbstract() {
			s.saveTypeAttribute(eClass)
		}
		s.str.addAttribute("href", href)
		s.str.endEmptyElement()
	}
}

func (s *XMLEncoder) saveIDRefSingle(eObject ecore.EObject, eFeature ecore.EStructuralFeature) {
	value, _ := eObject.EGetResolve(eFeature, false).(ecore.EObject)
	if value != nil {
		id := s.getIDRef(value)
		if id != "" {
			s.str.addAttribute(s.getFeatureQName(eFeature), id)
		}
	}
}

func (s *XMLEncoder) saveIDRefMany(eObject ecore.EObject, eFeature ecore.EStructuralFeature) {
	l := eObject.EGetResolve(eFeature, false).(ecore.EList)
	failure := false
	var buffer strings.Builder
	for it := l.Iterator(); ; {
		value, _ := it.Next().(ecore.EObject)
		if value != nil {
			id := s.getIDRef(value)
			if id == "" {
				failure = true
				if !it.HasNext() {
					break
				}
			} else {
				buffer.WriteString(id)
				if it.HasNext() {
					buffer.WriteString(" ")
				} else {
					break
				}
			}
		}
	}
	if !failure && buffer.Len() > 0 {
		s.str.addAttribute(s.getFeatureQName(eFeature), buffer.String())
	}
}

func (s *XMLEncoder) isNil(eObject ecore.EObject, eFeature ecore.EStructuralFeature) bool {
	return eObject.EGetResolve(eFeature, false) == nil
}

func (s *XMLEncoder) isEmpty(eObject ecore.EObject, eFeature ecore.EStructuralFeature) bool {
	return eObject.EGetResolve(eFeature, false).(ecore.EList).Empty()
}

func (s *XMLEncoder) shouldSaveFeature(o ecore.EObject, f ecore.EStructuralFeature) bool {
	return o.EIsSet(f) || (s.keepDefaults && f.GetDefaultValueLiteral() != "")
}

func (s *XMLEncoder) getSaveFeatureKind(f ecore.EStructuralFeature) xmlSaveFeatureKind {
	if f.IsTransient() {
		return xsfkTransient
	}

	isMany := f.IsMany()
	isUnsettable := f.IsUnsettable()

	if r, _ := f.(ecore.EReference); r != nil {
		if r.IsContainment() {
			if isMany {
				if isUnsettable {
					return xsfkObjectContainManyUnsettable
				} else {
					return xsfkObjectContainMany
				}
			} else {
				if isUnsettable {
					return xsfkObjectContainSingleUnsettable
				} else {
					return xsfkObjectContainSingle
				}
			}
		}
		opposite := r.GetEOpposite()
		if opposite != nil && opposite.IsContainment() {
			return xsfkTransient
		}
		if isMany {
			if isUnsettable {
				return xsfkObjectHrefManyUnsettable
			} else {
				return xsfkObjectHrefMany
			}
		} else {
			if isUnsettable {
				return xsfkObjectHrefSingleUnsettable
			} else {
				return xsfkObjectHrefSingle
			}
		}
	} else {
		// Attribute
		d := f.GetEType().(ecore.EDataType)
		if !d.IsSerializable() {
			return xsfkTransient
		}
		if isMany {
			return xsfkDataTypeMany
		}
		if isUnsettable {
			return xsfkDataTypeSingleNillable
		}
		return xsfkDataTypeSingle

	}

}

const (
	skip  = iota
	same  = iota
	cross = iota
)

func (s *XMLEncoder) getSaveResourceKindSingle(eObject ecore.EObject, eFeature ecore.EStructuralFeature) int {
	value, _ := eObject.EGetResolve(eFeature, false).(ecore.EObjectInternal)
	if value == nil {
		return skip
	} else if value.EIsProxy() {
		return cross
	} else {
		resource := value.EResource()
		if resource == s.resource || resource == nil {
			return same
		}
		return cross
	}
}

func (s *XMLEncoder) getSaveResourceKindMany(eObject ecore.EObject, eFeature ecore.EStructuralFeature) int {
	list, _ := eObject.EGetResolve(eFeature, false).(ecore.EList)
	if list == nil || list.Empty() {
		return skip
	}
	for it := list.Iterator(); it.HasNext(); {
		o, _ := it.Next().(ecore.EObjectInternal)
		if o == nil {
			return skip
		} else if o.EIsProxy() {
			return cross
		} else {
			r := o.EResource()
			if r != nil && r != s.resource {
				return cross
			}
		}

	}
	return same
}

func (s *XMLEncoder) getClassQName(eClass ecore.EClass) string {
	return s.getElementQName(eClass.GetEPackage(), s.getXmlName(eClass), false)
}

func (s *XMLEncoder) getFeatureQName(eFeature ecore.EStructuralFeature) string {
	if s.extendedMetaData != nil {
		name := s.extendedMetaData.GetName(eFeature)
		namespace := s.extendedMetaData.GetNamespace(eFeature)
		ePackage := s.getPackageForSpace(namespace)
		if ePackage != nil {
			return s.getElementQName(ePackage, name, false)
		} else {
			return name
		}
	} else {
		return eFeature.GetName()
	}
}

func (s *XMLEncoder) getElementQName(ePackage ecore.EPackage, name string, mustHavePrefix bool) string {
	nsPrefix := s.getPrefix(ePackage, mustHavePrefix)
	if nsPrefix == "" {
		return name
	} else if len(name) == 0 {
		return nsPrefix
	} else {
		return nsPrefix + ":" + name
	}
}

func (s *XMLEncoder) getXmlName(eElement ecore.ENamedElement) string {
	if s.extendedMetaData != nil {
		return s.extendedMetaData.GetName(eElement)
	}
	return eElement.GetName()
}

func (s *XMLEncoder) getPrefix(ePackage ecore.EPackage, mustHavePrefix bool) string {
	nsPrefix, ok := s.packages[ePackage]
	if !ok {
		nsURI := ePackage.GetNsURI()
		found := false
		prefixes := s.uriToPrefixes[nsURI]
		for _, prefix := range prefixes {
			nsPrefix = prefix
			if !mustHavePrefix || len(nsPrefix) > 0 {
				found = true
				break
			}
		}
		if !found {
			nsPrefix = ePackage.GetNsPrefix()
			if len(nsPrefix) == 0 && mustHavePrefix {
				nsPrefix = "_"
			}

			if uri, exists := s.prefixesToURI[nsPrefix]; exists && uri != nsURI {
				index := 1
				for _, exists = s.prefixesToURI[nsPrefix+"_"+fmt.Sprintf("%d", index)]; exists; {
					index++
				}
				nsPrefix += "_" + fmt.Sprintf("%d", index)
			}
			s.prefixesToURI[nsPrefix] = nsURI
		}

		s.packages[ePackage] = nsPrefix
	}
	return nsPrefix
}

func (s *XMLEncoder) setPrefixToNamespace(prefixToNamespaceMap ecore.EMap) {
	for it := prefixToNamespaceMap.Iterator(); it.HasNext(); {
		entry := it.Next().(ecore.EMapEntry)
		prefix := entry.GetKey().(string)
		nsURI := entry.GetValue().(string)
		if ePackage := s.getPackageForSpace(nsURI); ePackage != nil {
			s.packages[ePackage] = prefix
		}
		s.prefixesToURI[prefix] = nsURI
		s.uriToPrefixes[nsURI] = append(s.uriToPrefixes[nsURI], prefix)
	}
}

func (s *XMLEncoder) getPackageForSpace(nsURI string) ecore.EPackage {
	packageRegistry := ecore.GetPackageRegistry()
	if s.resource.GetResourceSet() != nil {
		packageRegistry = s.resource.GetResourceSet().GetPackageRegistry()
	}
	return packageRegistry.GetPackage(nsURI)
}

func (s *XMLEncoder) getDataType(value any, f ecore.EStructuralFeature, _ bool) (string, bool) {
	if value == nil {
		return "", false
	} else {
		d := f.GetEType().(ecore.EDataType)
		p := d.GetEPackage()
		f := p.GetEFactoryInstance()
		s := f.ConvertToString(d, value)
		return s, true
	}
}

func (s *XMLEncoder) getHRef(eObject ecore.EObject) string {
	eInternal, _ := eObject.(ecore.EObjectInternal)
	if eInternal != nil {
		objectURI := eInternal.EProxyURI()
		if objectURI == nil {
			eOtherResource := eObject.EResource()
			if eOtherResource == nil {
				if s.resource != nil && s.resource.GetObjectIDManager() != nil && s.resource.GetObjectIDManager().GetID(eObject) != nil {
					objectURI = s.getResourceHRef(s.resource, eObject)
				} else {
					s.handleDanglingHREF(eObject)
					return ""
				}
			} else {
				objectURI = s.getResourceHRef(eOtherResource, eObject)
			}
		}
		objectURI = s.resource.GetURI().Relativize(objectURI)
		return objectURI.String()
	}
	return ""
}

func (s *XMLEncoder) getResourceHRef(resource ecore.EResource, object ecore.EObject) *ecore.URI {
	return ecore.NewURIBuilder(resource.GetURI()).SetFragment(resource.GetURIFragment(object)).URI()
}

func (s *XMLEncoder) getIDRef(eObject ecore.EObject) string {
	if s.resource == nil {
		return ""
	} else {
		return "#" + s.resource.GetURIFragment(eObject)
	}
}

func (s *XMLEncoder) handleDanglingHREF(eObject ecore.EObject) {
	s.error(ecore.NewEDiagnosticImpl(fmt.Sprintf("Object '%p' is not contained in a resource.", eObject), s.resource.GetURI().String(), 0, 0))
}

func (s *XMLEncoder) error(diagnostic ecore.EDiagnostic) {
	s.errorFn(diagnostic)
}
