// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

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
	interfaces       interface{}
	w                io.Writer
	resource         EResource
	str              *xmlString
	packages         map[EPackage]string
	uriToPrefixes    map[string][]string
	prefixesToURI    map[string]string
	featureKinds     map[EStructuralFeature]xmlSaveFeatureKind
	extendedMetaData *ExtendedMetaData
	keepDefaults     bool
	idAttributeName  string
	roots            EList
	xmlVersion       string
	encoding         string
	errorFn          func(diagnostic EDiagnostic)
}

func NewXMLEncoder(resource EResource, w io.Writer, options map[string]interface{}) *XMLEncoder {
	s := new(XMLEncoder)
	s.interfaces = s
	s.resource = resource
	s.w = w
	s.xmlVersion = "1.0"
	s.encoding = "UTF-8"
	s.str = newXmlString()
	s.packages = make(map[EPackage]string)
	s.uriToPrefixes = make(map[string][]string)
	s.prefixesToURI = make(map[string]string)
	s.featureKinds = make(map[EStructuralFeature]xmlSaveFeatureKind)
	if options != nil {
		s.idAttributeName, _ = options[XML_OPTION_ID_ATTRIBUTE_NAME].(string)
		s.roots, _ = options[XML_OPTION_ROOT_OBJECTS].(EList)
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

func (s *XMLEncoder) Encode() {
	s.errorFn = func(diagnostic EDiagnostic) {
		s.resource.GetErrors().Add(diagnostic)
	}
	contents := s.roots
	if contents == nil {
		contents = s.resource.GetContents()
	}
	if contents.Empty() {
		return
	}
	s.encodeTopObject(contents.Get(0).(EObject))
}

func (s *XMLEncoder) EncodeObject(eObject EObject) (err error) {
	s.errorFn = func(diagnostic EDiagnostic) {
		if err == nil {
			err = diagnostic
		}
	}
	s.encodeTopObject(eObject)
	return
}

func (s *XMLEncoder) encodeTopObject(eObject EObject) {
	s.saveHeader()

	// initialize prefixes if any in top
	if s.extendedMetaData != nil {
		eClass := eObject.EClass()
		if ePrefixMapFeature := s.extendedMetaData.GetXMLNSPrefixMapFeature(eClass); ePrefixMapFeature != nil {
			m := eObject.EGet(ePrefixMapFeature).(EMap)
			s.setPrefixToNamespace(m)
		}
	}

	s.saveTopObject(eObject)

	// namespaces
	s.str.resetToFirstElementMark()
	s.interfaces.(xmlEncoderInternal).saveNamespaces()

	// write result
	s.str.write(s.w)
}

func (s *XMLEncoder) saveHeader() {
	s.str.add(fmt.Sprintf("<?xml version=\"%v\" encoding=\"%v\"?>", s.xmlVersion, s.encoding))
	s.str.addLine()
}

func (s *XMLEncoder) saveTopObject(eObject EObject) {
	eClass := eObject.EClass()
	if s.extendedMetaData == nil || s.extendedMetaData.GetDocumentRoot(eClass.GetEPackage()) != eClass {
		var name string
		if rootFeature := s.getRootFeature(eClass); rootFeature != nil {
			name = s.getFeatureQName(rootFeature)
		} else {
			name = s.getClassQName(eClass)
		}
		s.str.startElement(name)
	} else {
		s.str.startElement("")
	}
	s.saveElementID(eObject)
	s.saveFeatures(eObject, false)
}

func (s *XMLEncoder) getRootFeature(eClassifier EClassifier) EStructuralFeature {
	if s.extendedMetaData != nil {
		for eClassifier != nil {
			if eClass := s.extendedMetaData.GetDocumentRoot(eClassifier.GetEPackage()); eClass != nil {
				for it := eClass.GetEStructuralFeatures().Iterator(); it.HasNext(); {
					eFeature := it.Next().(EStructuralFeature)
					if eFeature.GetEType() == eClassifier && eFeature.IsChangeable() {
						return eFeature
					}
				}
			}
			if eClass, _ := eClassifier.(EClass); eClass != nil {
				eSuperTypes := eClass.GetESuperTypes()
				if eSuperTypes.Empty() {
					eClassifier = nil
				} else {
					eClassifier = eSuperTypes.Get(0).(EClass)
				}
			} else {
				eClassifier = nil
			}
		}
	}
	return nil
}

func (s *XMLEncoder) saveNamespaces() {
	var prefixes []string
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

func (s *XMLEncoder) saveElementID(eObject EObject) {
	if idManager := s.resource.GetObjectIDManager(); len(s.idAttributeName) > 0 && idManager != nil {
		id := idManager.GetID(eObject)
		var objectID string
		switch id.(type) {
		case nil:
			objectID = ""
		case int:
			objectID = strconv.Itoa(id.(int))
		case string:
			objectID = id.(string)
		}
		if len(objectID) > 0 {
			s.str.addAttribute(s.idAttributeName, objectID)
		}
	}

}

func (s *XMLEncoder) saveFeatures(eObject EObject, attributesOnly bool) bool {
	eClass := eObject.EClass()
	eAllFeatures := eClass.GetEAllStructuralFeatures()
	var elementFeatures []int
	elementCount := 0
	i := 0
	for it := eAllFeatures.Iterator(); it.HasNext(); i++ {
		// current feature
		eFeature := it.Next().(EStructuralFeature)
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
		eFeature := eAllFeatures.Get(elementFeatures[i]).(EStructuralFeature)
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

func (s *XMLEncoder) saveDataTypeSingle(eObject EObject, eFeature EStructuralFeature) {
	val := eObject.EGetResolve(eFeature, false)
	str, ok := s.getDataType(val, eFeature, true)
	if ok {
		s.str.addAttribute(s.getFeatureQName(eFeature), html.EscapeString(str))
	}
}

func (s *XMLEncoder) saveDataTypeMany(eObject EObject, eFeature EStructuralFeature) {
	l := eObject.EGetResolve(eFeature, false).(EList)
	d := eFeature.GetEType().(EDataType)
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

func (s *XMLEncoder) saveManyEmpty(eObject EObject, eFeature EStructuralFeature) {
	s.str.addAttribute(s.getFeatureQName(eFeature), "")
}

func (s *XMLEncoder) saveEObjectSingle(eObject EObject, eFeature EStructuralFeature) {
	value, _ := eObject.EGetResolve(eFeature, false).(EObject)
	if value != nil {
		id := s.getHRef(value)
		s.str.addAttribute(s.getFeatureQName(eFeature), id)
	}
}

func (s *XMLEncoder) saveEObjectMany(eObject EObject, eFeature EStructuralFeature) {
	l := eObject.EGetResolve(eFeature, false).(EList)
	failure := false
	var buffer strings.Builder
	for it := l.Iterator(); ; {
		value, _ := it.Next().(EObject)
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

func (s *XMLEncoder) saveNil(eObject EObject, eFeature EStructuralFeature) {
	s.str.addNil(s.getFeatureQName(eFeature))
}

func (s *XMLEncoder) saveContainedSingle(eObject EObject, eFeature EStructuralFeature) {
	value, _ := eObject.EGetResolve(eFeature, false).(EObjectInternal)
	if value != nil {
		s.saveEObjectInternal(value, eFeature)
	}
}

func (s *XMLEncoder) saveContainedMany(eObject EObject, eFeature EStructuralFeature) {
	l := eObject.EGetResolve(eFeature, false).(EList)
	for it := l.Iterator(); it.HasNext(); {
		value, _ := it.Next().(EObjectInternal)
		if value != nil {
			s.saveEObjectInternal(value, eFeature)
		}
	}
}

func (s *XMLEncoder) saveEObjectInternal(o EObjectInternal, f EStructuralFeature) {
	if o.EInternalResource() != nil || o.EIsProxy() {
		s.saveHRef(o, f)
	} else {
		s.saveEObject(o, f)
	}
}

func (s *XMLEncoder) saveEObject(o EObject, f EStructuralFeature) {
	s.str.startElement(s.getFeatureQName(f))
	eClass := o.EClass()
	eType := f.GetEType()
	if eType != eClass && eType != GetPackage().GetEObject() {
		s.saveTypeAttribute(eClass)
	}
	s.saveElementID(o)
	s.saveFeatures(o, false)
}

func (s *XMLEncoder) saveTypeAttribute(eClass EClass) {
	s.str.addAttribute("xsi:type", s.getClassQName(eClass))
	s.uriToPrefixes[xsiURI] = []string{xsiNS}
	s.prefixesToURI[xsiNS] = xsiURI
}

func (s *XMLEncoder) saveHRefSingle(eObject EObject, eFeature EStructuralFeature) {
	value, _ := eObject.EGetResolve(eFeature, false).(EObject)
	if value != nil {
		s.saveHRef(value, eFeature)
	}
}

func (s *XMLEncoder) saveHRefMany(eObject EObject, eFeature EStructuralFeature) {
	l := eObject.EGetResolve(eFeature, false).(EList)
	for it := l.Iterator(); it.HasNext(); {
		value, _ := it.Next().(EObject)
		if value != nil {
			s.saveHRef(value, eFeature)
		}
	}
}

func (s *XMLEncoder) saveHRef(eObject EObject, eFeature EStructuralFeature) {
	href := s.getHRef(eObject)
	if href != "" {
		s.str.startElement(s.getFeatureQName(eFeature))
		eClass := eObject.EClass()
		eType, _ := eFeature.GetEType().(EClass)
		if eType != eClass && eType != nil && eType.IsAbstract() {
			s.saveTypeAttribute(eClass)
		}
		s.str.addAttribute("href", href)
		s.str.endEmptyElement()
	}
}

func (s *XMLEncoder) saveIDRefSingle(eObject EObject, eFeature EStructuralFeature) {
	value, _ := eObject.EGetResolve(eFeature, false).(EObject)
	if value != nil {
		id := s.getIDRef(value)
		if id != "" {
			s.str.addAttribute(s.getFeatureQName(eFeature), id)
		}
	}
}

func (s *XMLEncoder) saveIDRefMany(eObject EObject, eFeature EStructuralFeature) {
	l := eObject.EGetResolve(eFeature, false).(EList)
	failure := false
	var buffer strings.Builder
	for it := l.Iterator(); ; {
		value, _ := it.Next().(EObject)
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

func (s *XMLEncoder) isNil(eObject EObject, eFeature EStructuralFeature) bool {
	return eObject.EGetResolve(eFeature, false) == nil
}

func (s *XMLEncoder) isEmpty(eObject EObject, eFeature EStructuralFeature) bool {
	return eObject.EGetResolve(eFeature, false).(EList).Empty()
}

func (s *XMLEncoder) shouldSaveFeature(o EObject, f EStructuralFeature) bool {
	return o.EIsSet(f) || (s.keepDefaults && f.GetDefaultValueLiteral() != "")
}

func (s *XMLEncoder) getSaveFeatureKind(f EStructuralFeature) xmlSaveFeatureKind {
	if f.IsTransient() {
		return xsfkTransient
	}

	isMany := f.IsMany()
	isUnsettable := f.IsUnsettable()

	if r, _ := f.(EReference); r != nil {
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
		d, _ := f.GetEType().(EDataType)
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

func (s *XMLEncoder) getSaveResourceKindSingle(eObject EObject, eFeature EStructuralFeature) int {
	value, _ := eObject.EGetResolve(eFeature, false).(EObjectInternal)
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

func (s *XMLEncoder) getSaveResourceKindMany(eObject EObject, eFeature EStructuralFeature) int {
	list, _ := eObject.EGetResolve(eFeature, false).(EList)
	if list == nil || list.Empty() {
		return skip
	}
	for it := list.Iterator(); it.HasNext(); {
		o, _ := it.Next().(EObjectInternal)
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

func (s *XMLEncoder) getClassQName(eClass EClass) string {
	return s.getElementQName(eClass.GetEPackage(), s.getXmlName(eClass), false)
}

func (s *XMLEncoder) getFeatureQName(eFeature EStructuralFeature) string {
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

func (s *XMLEncoder) getElementQName(ePackage EPackage, name string, mustHavePrefix bool) string {
	nsPrefix := s.getPrefix(ePackage, mustHavePrefix)
	if nsPrefix == "" {
		return name
	} else if len(name) == 0 {
		return nsPrefix
	} else {
		return nsPrefix + ":" + name
	}
}

func (s *XMLEncoder) getXmlName(eElement ENamedElement) string {
	if s.extendedMetaData != nil {
		return s.extendedMetaData.GetName(eElement)
	}
	return eElement.GetName()
}

func (s *XMLEncoder) getPrefix(ePackage EPackage, mustHavePrefix bool) string {
	nsPrefix, ok := s.packages[ePackage]
	if !ok {
		nsURI := ePackage.GetNsURI()
		found := false
		prefixes := s.uriToPrefixes[nsURI]
		if prefixes != nil {
			for _, prefix := range prefixes {
				nsPrefix = prefix
				if !mustHavePrefix || len(nsPrefix) > 0 {
					found = true
					break
				}
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

func (s *XMLEncoder) setPrefixToNamespace(prefixToNamespaceMap EMap) {
	for it := prefixToNamespaceMap.Iterator(); it.HasNext(); {
		entry := it.Next().(EMapEntry)
		prefix := entry.GetKey().(string)
		nsURI := entry.GetValue().(string)
		if ePackage := s.getPackageForSpace(nsURI); ePackage != nil {
			s.packages[ePackage] = prefix
		}
		s.prefixesToURI[prefix] = nsURI
		s.uriToPrefixes[nsURI] = append(s.uriToPrefixes[nsURI], prefix)
	}
}

func (s *XMLEncoder) getPackageForSpace(nsURI string) EPackage {
	packageRegistry := GetPackageRegistry()
	if s.resource.GetResourceSet() != nil {
		packageRegistry = s.resource.GetResourceSet().GetPackageRegistry()
	}
	return packageRegistry.GetPackage(nsURI)
}

func (s *XMLEncoder) getDataType(value interface{}, f EStructuralFeature, isAttribute bool) (string, bool) {
	if value == nil {
		return "", false
	} else {
		d := f.GetEType().(EDataType)
		p := d.GetEPackage()
		f := p.GetEFactoryInstance()
		s := f.ConvertToString(d, value)
		return s, true
	}
}

func (s *XMLEncoder) getHRef(eObject EObject) string {
	eInternal, _ := eObject.(EObjectInternal)
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

func (s *XMLEncoder) getResourceHRef(resource EResource, object EObject) *URI {
	uri := resource.GetURI().Copy()
	uri.Fragment = resource.GetURIFragment(object)
	return uri
}

func (s *XMLEncoder) getIDRef(eObject EObject) string {
	if s.resource == nil {
		return ""
	} else {
		return "#" + s.resource.GetURIFragment(eObject)
	}
}

func (s *XMLEncoder) handleDanglingHREF(eObject EObject) {
	s.error(NewEDiagnosticImpl(fmt.Sprintf("Object '%p' is not contained in a resource.", eObject), s.resource.GetURI().String(), 0, 0))
}

func (s *XMLEncoder) error(diagnostic EDiagnostic) {
	s.errorFn(diagnostic)
}
