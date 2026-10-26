// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package xml

import (
	"fmt"
	"io"
	"math"
	"sort"
	"strings"

	"github.com/masagroup/soft.go/ecore"
	"golang.org/x/net/html"
)

// whitespaceEscaper escapes the whitespaces an XML parser normalizes to spaces in an attribute value
// (html.EscapeString escapes \r already)
var whitespaceEscaper = strings.NewReplacer("\n", "&#xA;", "\t", "&#x9;")

// escapeAttribute escapes s as the value of an attribute
func escapeAttribute(s string) string {
	return whitespaceEscaper.Replace(html.EscapeString(s))
}

// escapeContent escapes s as the content of an element
func escapeContent(s string) string {
	return html.EscapeString(s)
}

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

// EncoderHooks are the methods of an encoder that an encoder embedding Encoder (as the XMI one)
// can override: Encoder calls them on the object set with SetInterfaces.
type EncoderHooks interface {
	SaveNamespaces()
}

type Encoder struct {
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
	useEncodedStyle  bool
}

func NewEncoder(resource ecore.EResource, w io.Writer, options map[string]any) *Encoder {
	s := new(Encoder)
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
		s.idAttributeName, _ = options[OPTION_ID_ATTRIBUTE_NAME].(string)
		s.roots, _ = options[OPTION_ROOT_OBJECTS].(ecore.EList)
		s.useEncodedStyle = options[OPTION_USE_ENCODED_ATTRIBUTE_STYLE] == true
		if extendedMetaData := options[OPTION_EXTENDED_META_DATA]; extendedMetaData != nil {
			s.extendedMetaData = extendedMetaData.(*ExtendedMetaData)
		}
	}
	if s.extendedMetaData == nil {
		s.extendedMetaData = NewExtendedMetaData()
	}
	return s
}

func (s *Encoder) SetEncoding(encoding string) {
	s.encoding = encoding
}

func (s *Encoder) SetXMLVersion(xmlVersion string) {
	s.xmlVersion = xmlVersion
}

func (s *Encoder) EncodeResource() {
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

func (s *Encoder) EncodeObject(eObject ecore.EObject) (err error) {
	s.errorFn = func(diagnostic ecore.EDiagnostic) {
		if err == nil {
			err = diagnostic
		}
	}
	s.encodeTopObject(eObject)
	return
}

func (s *Encoder) encodeTopObject(eObject ecore.EObject) {
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
	s.interfaces.(EncoderHooks).SaveNamespaces()

	// write result
	if err := s.str.write(s.w); err != nil {
		s.error(ecore.NewEDiagnosticImpl(err.Error(), s.location(), 0, 0))
	}
}

func (s *Encoder) saveHeader() {
	s.str.add(fmt.Sprintf("<?xml version=\"%v\" encoding=\"%v\"?>", s.xmlVersion, s.encoding))
	s.str.addLine()
}

func (s *Encoder) saveTopObject(eObject ecore.EObject) {
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

func (s *Encoder) getRootFeature(eClassifier ecore.EClassifier) ecore.EStructuralFeature {
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

func (s *Encoder) SaveNamespaces() {
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
		s.str.addAttribute(attribute, escapeAttribute(s.prefixesToURI[prefix]))
	}
}

func (s *Encoder) saveElementID(eObject ecore.EObject) {
	if idManager := s.resource.GetObjectIDManager(); len(s.idAttributeName) > 0 && idManager != nil {
		if id := idManager.GetID(eObject); id != nil {
			s.str.addAttribute(s.idAttributeName, escapeAttribute(fmt.Sprintf("%v", id)))
		}
	}
}

func (s *Encoder) saveFeatures(eObject ecore.EObject, attributesOnly bool) bool {
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
			case xsfkObjectAttributeSingle:
				// An attribute reference: a full href across resources, an IDREF otherwise.
				switch s.getSaveResourceKindSingle(eObject, eFeature) {
				case cross:
					s.saveEObjectSingle(eObject, eFeature)
				case same:
					s.saveIDRefSingle(eObject, eFeature)
				}
				continue
			case xsfkObjectAttributeMany:
				switch s.getSaveResourceKindMany(eObject, eFeature) {
				case cross:
					s.saveEObjectMany(eObject, eFeature)
				case same:
					s.saveIDRefMany(eObject, eFeature)
				}
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
			case xsfkDataTypeContentSingle:
			case xsfkDataTypeElementSingle:
			case xsfkObjectHrefSingleUnsettable:
				if !s.isNil(eObject, eFeature) {
					switch s.getSaveResourceKindSingle(eObject, eFeature) {
					case cross:
						if s.useEncodedStyle {
							s.saveEncodedHRefSingle(eObject, eFeature)
							continue
						}
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
					if s.useEncodedStyle {
						s.saveEncodedHRefSingle(eObject, eFeature)
						continue
					}
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
						if s.useEncodedStyle {
							s.saveEncodedHRefMany(eObject, eFeature)
							continue
						}
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
					if s.useEncodedStyle {
						s.saveEncodedHRefMany(eObject, eFeature)
						continue
					}
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
		case xsfkDataTypeContentSingle:
			s.saveDataTypeContent(eObject, eFeature)
		case xsfkDataTypeElementSingle:
			s.saveDataTypeElementSingle(eObject, eFeature)
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

func (s *Encoder) saveDataTypeSingle(eObject ecore.EObject, eFeature ecore.EStructuralFeature) {
	val := eObject.EGetResolve(eFeature, false)
	str, ok := s.getDataType(val, eFeature, true)
	if ok {
		s.str.addAttribute(s.getFeatureQName(eFeature), escapeAttribute(str))
	}
}

// saveDataTypeElementSingle writes a kind="element" scalar attribute as a child element.
func (s *Encoder) saveDataTypeElementSingle(eObject ecore.EObject, eFeature ecore.EStructuralFeature) {
	val := eObject.EGetResolve(eFeature, false)
	str, ok := s.getDataType(val, eFeature, true)
	if ok {
		s.str.addContent(s.getFeatureQName(eFeature), escapeContent(str))
	}
}

// saveDataTypeContent writes a kind="simple" feature as the element's text.
func (s *Encoder) saveDataTypeContent(eObject ecore.EObject, eFeature ecore.EStructuralFeature) {
	val := eObject.EGetResolve(eFeature, false)
	str, ok := s.getDataType(val, eFeature, true)
	if ok {
		s.str.addContentValue(escapeContent(str))
	}
}

func (s *Encoder) saveDataTypeMany(eObject ecore.EObject, eFeature ecore.EStructuralFeature) {
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
			str := convertToString(f, d, value)
			s.str.addContent(name, escapeContent(str))
		}
	}
}

func (s *Encoder) saveManyEmpty(_ ecore.EObject, eFeature ecore.EStructuralFeature) {
	s.str.addAttribute(s.getFeatureQName(eFeature), "")
}

func (s *Encoder) saveEObjectSingle(eObject ecore.EObject, eFeature ecore.EStructuralFeature) {
	value, _ := eObject.EGetResolve(eFeature, false).(ecore.EObject)
	if value != nil {
		id := s.getHRef(value)
		s.str.addAttribute(s.getFeatureQName(eFeature), escapeAttribute(id))
	}
}

func (s *Encoder) saveEObjectMany(eObject ecore.EObject, eFeature ecore.EStructuralFeature) {
	s.saveRefsAttribute(eObject, eFeature, s.getHRef)
}

func (s *Encoder) saveNil(_ ecore.EObject, eFeature ecore.EStructuralFeature) {
	s.str.addNil(s.getFeatureQName(eFeature))
}

func (s *Encoder) saveContainedSingle(eObject ecore.EObject, eFeature ecore.EStructuralFeature) {
	value, _ := eObject.EGetResolve(eFeature, false).(ecore.EObjectInternal)
	if value != nil {
		s.saveEObjectInternal(value, eFeature)
	}
}

func (s *Encoder) saveContainedMany(eObject ecore.EObject, eFeature ecore.EStructuralFeature) {
	l := eObject.EGetResolve(eFeature, false).(ecore.EList)
	for it := l.Iterator(); it.HasNext(); {
		value, _ := it.Next().(ecore.EObjectInternal)
		if value != nil {
			s.saveEObjectInternal(value, eFeature)
		}
	}
}

func (s *Encoder) saveEObjectInternal(o ecore.EObjectInternal, f ecore.EStructuralFeature) {
	if o.EInternalResource() != nil || o.EIsProxy() {
		s.saveHRef(o, f)
	} else {
		s.saveEObject(o, f)
	}
}

func (s *Encoder) saveEObject(o ecore.EObject, f ecore.EStructuralFeature) {
	s.str.startElement(s.getFeatureQName(f))
	eClass := o.EClass()
	eType := f.GetEType()
	if eType != eClass && eType != ecore.GetPackage().GetEObject() {
		s.saveTypeAttribute(eClass)
	}
	s.saveElementID(o)
	s.saveFeatures(o, false)
}

func (s *Encoder) saveTypeAttribute(eClass ecore.EClass) {
	s.str.addAttribute("xsi:type", s.getClassQName(eClass))
	s.uriToPrefixes[xsiURI] = []string{xsiNS}
	s.prefixesToURI[xsiNS] = xsiURI
}

func (s *Encoder) saveHRefSingle(eObject ecore.EObject, eFeature ecore.EStructuralFeature) {
	value, _ := eObject.EGetResolve(eFeature, false).(ecore.EObject)
	if value != nil {
		s.saveHRef(value, eFeature)
	}
}

func (s *Encoder) saveHRefMany(eObject ecore.EObject, eFeature ecore.EStructuralFeature) {
	l := eObject.EGetResolve(eFeature, false).(ecore.EList)
	for it := l.Iterator(); it.HasNext(); {
		value, _ := it.Next().(ecore.EObject)
		if value != nil {
			s.saveHRef(value, eFeature)
		}
	}
}

func (s *Encoder) saveHRef(eObject ecore.EObject, eFeature ecore.EStructuralFeature) {
	href := s.getHRef(eObject)
	if href != "" {
		s.str.startElement(s.getFeatureQName(eFeature))
		eClass := eObject.EClass()
		eType, _ := eFeature.GetEType().(ecore.EClass)
		if eType != eClass && eType != nil && eType.IsAbstract() {
			s.saveTypeAttribute(eClass)
		}
		s.str.addAttribute("href", escapeAttribute(href))
		s.str.endEmptyElement()
	}
}

// saveEncodedHRefSingle saves a cross resource reference in the encoded attribute style, as EMF
// does for .ecore files: eType="ecore:EClass other.ecore#//X"
func (s *Encoder) saveEncodedHRefSingle(eObject ecore.EObject, eFeature ecore.EStructuralFeature) {
	value, _ := eObject.EGetResolve(eFeature, false).(ecore.EObject)
	if value != nil {
		if href := s.getEncodedHRef(value, eFeature); href != "" {
			s.str.addAttribute(s.getFeatureQName(eFeature), escapeAttribute(href))
		}
	}
}

func (s *Encoder) saveEncodedHRefMany(eObject ecore.EObject, eFeature ecore.EStructuralFeature) {
	s.saveRefsAttribute(eObject, eFeature, func(value ecore.EObject) string {
		return s.getEncodedHRef(value, eFeature)
	})
}

// getEncodedHRef is the href of eObject prefixed with the qualified name of its class when the
// type of eFeature does not tell it, as the xsi:type of an href element would be.
func (s *Encoder) getEncodedHRef(eObject ecore.EObject, eFeature ecore.EStructuralFeature) string {
	href := s.getHRef(eObject)
	if href != "" {
		eClass := eObject.EClass()
		eType, _ := eFeature.GetEType().(ecore.EClass)
		if eType != eClass && eType != nil && eType.IsAbstract() {
			return s.getClassQName(eClass) + " " + href
		}
	}
	return href
}

func (s *Encoder) saveIDRefSingle(eObject ecore.EObject, eFeature ecore.EStructuralFeature) {
	value, _ := eObject.EGetResolve(eFeature, false).(ecore.EObject)
	if value != nil {
		id := s.getIDRef(value)
		if id != "" {
			s.str.addAttribute(s.getFeatureQName(eFeature), escapeAttribute(id))
		}
	}
}

func (s *Encoder) saveIDRefMany(eObject ecore.EObject, eFeature ecore.EStructuralFeature) {
	s.saveRefsAttribute(eObject, eFeature, s.getIDRef)
}

// saveRefsAttribute saves the references of a many feature as an attribute of their space
// separated ids, nothing if an id is missing. The nil elements are skipped.
func (s *Encoder) saveRefsAttribute(eObject ecore.EObject, eFeature ecore.EStructuralFeature, getID func(ecore.EObject) string) {
	l := eObject.EGetResolve(eFeature, false).(ecore.EList)
	var buffer strings.Builder
	for it := l.Iterator(); it.HasNext(); {
		value, _ := it.Next().(ecore.EObject)
		if value == nil {
			continue
		}
		id := getID(value)
		if id == "" {
			return
		}
		if buffer.Len() > 0 {
			buffer.WriteString(" ")
		}
		buffer.WriteString(id)
	}
	if buffer.Len() > 0 {
		s.str.addAttribute(s.getFeatureQName(eFeature), escapeAttribute(buffer.String()))
	}
}

func (s *Encoder) isNil(eObject ecore.EObject, eFeature ecore.EStructuralFeature) bool {
	return eObject.EGetResolve(eFeature, false) == nil
}

func (s *Encoder) isEmpty(eObject ecore.EObject, eFeature ecore.EStructuralFeature) bool {
	return eObject.EGetResolve(eFeature, false).(ecore.EList).Empty()
}

func (s *Encoder) shouldSaveFeature(o ecore.EObject, f ecore.EStructuralFeature) bool {
	return o.EIsSet(f) || (s.keepDefaults && f.GetDefaultValueLiteral() != "")
}

func (s *Encoder) getSaveFeatureKind(f ecore.EStructuralFeature) xmlSaveFeatureKind {
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
		// A kind="attribute" reference is an attribute, as the decoder reads it.
		if s.extendedMetaData != nil && s.extendedMetaData.GetFeatureKind(f) == "attribute" {
			if isMany {
				return xsfkObjectAttributeMany
			}
			return xsfkObjectAttributeSingle
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
		// A kind="simple" attribute is the element text, not an attribute.
		if !isMany && s.extendedMetaData != nil && s.extendedMetaData.GetFeatureKind(f) == "simple" {
			return xsfkDataTypeContentSingle
		}
		// A kind="element" scalar attribute is a child element, as the decoder reads it.
		if !isMany && s.extendedMetaData != nil && s.extendedMetaData.GetFeatureKind(f) == "element" {
			return xsfkDataTypeElementSingle
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

func (s *Encoder) getSaveResourceKindSingle(eObject ecore.EObject, eFeature ecore.EStructuralFeature) int {
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

func (s *Encoder) getSaveResourceKindMany(eObject ecore.EObject, eFeature ecore.EStructuralFeature) int {
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

func (s *Encoder) getClassQName(eClass ecore.EClass) string {
	return s.getElementQName(eClass.GetEPackage(), s.getXmlName(eClass), false)
}

func (s *Encoder) getFeatureQName(eFeature ecore.EStructuralFeature) string {
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

func (s *Encoder) getElementQName(ePackage ecore.EPackage, name string, mustHavePrefix bool) string {
	nsPrefix := s.getPrefix(ePackage, mustHavePrefix)
	if nsPrefix == "" {
		return name
	} else if len(name) == 0 {
		return nsPrefix
	} else {
		return nsPrefix + ":" + name
	}
}

func (s *Encoder) getXmlName(eElement ecore.ENamedElement) string {
	if s.extendedMetaData != nil {
		return s.extendedMetaData.GetName(eElement)
	}
	return eElement.GetName()
}

func (s *Encoder) getPrefix(ePackage ecore.EPackage, mustHavePrefix bool) string {
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
				for _, exists = s.prefixesToURI[nsPrefix+"_"+fmt.Sprintf("%d", index)]; exists; _, exists = s.prefixesToURI[nsPrefix+"_"+fmt.Sprintf("%d", index)] {
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

func (s *Encoder) setPrefixToNamespace(prefixToNamespaceMap ecore.EMap) {
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

func (s *Encoder) getPackageForSpace(nsURI string) ecore.EPackage {
	packageRegistry := ecore.GetPackageRegistry()
	if s.resource.GetResourceSet() != nil {
		packageRegistry = s.resource.GetResourceSet().GetPackageRegistry()
	}
	return packageRegistry.GetPackage(nsURI)
}

func (s *Encoder) getDataType(value any, f ecore.EStructuralFeature, _ bool) (string, bool) {
	if value == nil {
		return "", false
	} else {
		d := f.GetEType().(ecore.EDataType)
		p := d.GetEPackage()
		f := p.GetEFactoryInstance()
		s := convertToString(f, d, value)
		return s, true
	}
}

// convertToString is the lexical form of a data value: the one of its factory, but for the
// infinite and NaN floats, written in the XSD lexical space of xsd:double and xsd:float (INF,
// -INF, NaN) instead of the Go forms (+Inf, -Inf). The decoder reads both.
func convertToString(eFactory ecore.EFactory, eDataType ecore.EDataType, value any) string {
	var v float64
	switch value := value.(type) {
	case float64:
		v = value
	case float32:
		v = float64(value)
	default:
		return eFactory.ConvertToString(eDataType, value)
	}
	switch {
	case math.IsInf(v, 1):
		return "INF"
	case math.IsInf(v, -1):
		return "-INF"
	case math.IsNaN(v):
		return "NaN"
	}
	return eFactory.ConvertToString(eDataType, value)
}

func (s *Encoder) getHRef(eObject ecore.EObject) string {
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

func (s *Encoder) getResourceHRef(resource ecore.EResource, object ecore.EObject) *ecore.URI {
	return ecore.NewURIBuilder(resource.GetURI()).SetFragment(resource.GetURIFragment(object)).URI()
}

func (s *Encoder) getIDRef(eObject ecore.EObject) string {
	if s.resource == nil {
		return ""
	} else {
		return "#" + s.resource.GetURIFragment(eObject)
	}
}

func (s *Encoder) handleDanglingHREF(eObject ecore.EObject) {
	s.error(ecore.NewEDiagnosticImpl(fmt.Sprintf("Object '%p' is not contained in a resource.", eObject), s.location(), 0, 0))
}

func (s *Encoder) error(diagnostic ecore.EDiagnostic) {
	s.errorFn(diagnostic)
}

// location is the uri of the resource in the diagnostics, empty for a resource without uri.
func (s *Encoder) location() string {
	if uri := s.resource.GetURI(); uri != nil {
		return uri.String()
	}
	return ""
}

// SetInterfaces sets the object whose EncoderHooks the encoder calls: an encoder embedding
// Encoder sets itself to override them.
func (s *Encoder) SetInterfaces(interfaces any) {
	s.interfaces = interfaces
}

// AddAttribute adds an attribute to the start tag being saved: it is meant to be called from
// SaveNamespaces.
func (s *Encoder) AddAttribute(name string, value string) {
	s.str.addAttribute(name, value)
}

// SetExtendedMetaData sets the extended metadata of the encoder, nil for none.
func (s *Encoder) SetExtendedMetaData(extendedMetaData *ExtendedMetaData) {
	s.extendedMetaData = extendedMetaData
}
