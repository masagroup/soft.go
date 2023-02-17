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
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html/charset"
)

const (
	href                            = "href"
	typeAttrib                      = "type"
	schemaLocationAttrib            = "schemaLocation"
	noNamespaceSchemaLocationAttrib = "noNamespaceSchemaLocation"
	xsiURI                          = "http://www.w3.org/2001/XMLSchema-instance"
	xsiNS                           = "xsi"
	xmlNS                           = "xmlns"
)

type xmlLoadFeatureKind int

const (
	xlfkSingle xmlLoadFeatureKind = iota
	xlfkMany
	xlfkManyAdd
	xlfkManyMove
	xlfkOther
)

const (
	load_object_type = "object"
	load_error_type  = "error"
)

type reference struct {
	object  EObject
	feature EStructuralFeature
	id      string
	pos     int
}

type xmlDecoderInternal interface {
	getXSIType() string
	handleAttributes(object EObject)
}

type XMLDecoder struct {
	resource               EResource
	interfaces             any
	spacesToFactories      map[string]EFactory
	attachFn               func(object EObject)
	errorFn                func(diagnostic EDiagnostic)
	extendedMetaData       *ExtendedMetaData
	prefixesToURI          map[string]string
	decoder                *xml.Decoder
	textBuilder            *strings.Builder
	namespaces             *xmlNamespaces
	xmlVersion             string
	encoding               string
	idAttributeName        string
	references             []reference
	attributes             []xml.Attr
	sameDocumentProxies    []EObject
	notFeatures            []xml.Name
	types                  []any
	objects                []EObject
	deferred               []EObject
	elements               []string
	isSuppressDocumentRoot bool
	isResolveDeferred      bool
}

func NewXMLDecoder(resource EResource, r io.Reader, options map[string]any) *XMLDecoder {
	l := new(XMLDecoder)
	l.interfaces = l
	l.resource = resource
	l.decoder = xml.NewDecoder(r)
	l.decoder.CharsetReader = charset.NewReaderLabel
	l.namespaces = newXMLNamespaces()
	l.prefixesToURI = make(map[string]string)
	l.spacesToFactories = make(map[string]EFactory)
	l.notFeatures = append(l.notFeatures, xml.Name{Space: xsiURI, Local: typeAttrib}, xml.Name{Space: xsiURI, Local: schemaLocationAttrib}, xml.Name{Space: xsiURI, Local: noNamespaceSchemaLocationAttrib})
	if options != nil {
		l.idAttributeName, _ = options[XML_OPTION_ID_ATTRIBUTE_NAME].(string)
		l.isResolveDeferred = options[XML_OPTION_DEFERRED_REFERENCE_RESOLUTION] == true
		l.isSuppressDocumentRoot = options[XML_OPTION_SUPPRESS_DOCUMENT_ROOT] == true
		if extendedMetaData := options[XML_OPTION_EXTENDED_META_DATA]; extendedMetaData != nil {
			l.extendedMetaData = extendedMetaData.(*ExtendedMetaData)
		}
		if options[XML_OPTION_DEFERRED_ROOT_ATTACHMENT] == true {
			l.deferred = []EObject{}
		}
	}
	if l.extendedMetaData == nil {
		l.extendedMetaData = NewExtendedMetaData()
	}
	return l
}

func (l *XMLDecoder) GetXMLVersion() string {
	return l.xmlVersion
}

func (l *XMLDecoder) GetEncoding() string {
	return l.encoding
}

func (l *XMLDecoder) Decode() {
	l.attachFn = func(object EObject) {
		l.resource.GetContents().Add(object)
	}
	l.errorFn = func(diagnostic EDiagnostic) {
		l.resource.GetErrors().Add(diagnostic)
	}
	l.decodeTopObject()
}

func (l *XMLDecoder) DecodeObject() (eObject EObject, err error) {
	l.attachFn = func(o EObject) {
		if eObject == nil {
			eObject = o
		}
	}
	l.errorFn = func(diagnostic EDiagnostic) {
		if err == nil {
			err = diagnostic
		}
	}
	l.decodeTopObject()
	return
}

func (l *XMLDecoder) decodeTopObject() {
	for {
		t, tokenErr := l.decoder.Token()
		if tokenErr != nil {
			if tokenErr != io.EOF {
				l.error(NewEDiagnosticImpl(tokenErr.Error(), l.resource.GetURI().String(), int(l.decoder.InputOffset()), 0))
			}
			break
		}
		switch t := t.(type) {
		case xml.StartElement:
			l.startElement(t)
		case xml.EndElement:
			l.endElement(t)
		case xml.CharData:
			l.text(string([]byte(t)))
		case xml.Comment:
			l.comment(string([]byte(t)))
		case xml.ProcInst:
			l.processingInstruction(t)
		case xml.Directive:
			l.directive(string([]byte(t)))
		}
	}
}

func (l *XMLDecoder) startElement(e xml.StartElement) {
	l.elements = append(l.elements, e.Name.Local)
	l.setAttributes(e.Attr)
	l.namespaces.pushContext()
	l.handlePrefixMapping()
	if len(l.objects) == 0 {
		l.handleSchemaLocation()
	}
	l.processElement(e.Name.Space, e.Name.Local)
}

func (l *XMLDecoder) endElement(e xml.EndElement) {
	if len(l.elements) > 0 {
		l.elements = l.elements[:len(l.elements)-1]
	}

	// remove last object
	var eRoot EObject
	var eObject EObject
	if len(l.objects) > 0 {
		eRoot = l.objects[0]
		eObject = l.objects[len(l.objects)-1]
		l.objects = l.objects[:len(l.objects)-1]
	}

	// remove last type
	var eType any
	if len(l.types) > 0 {
		eType = l.types[len(l.types)-1]
		l.types = l.types[:len(l.types)-1]
	}
	if l.textBuilder != nil {
		if eType == load_object_type {
			if l.textBuilder.Len() > 0 {
				l.handleProxy(eObject, l.textBuilder.String())
			}
		} else if eType != load_error_type {
			if eObject == nil && len(l.objects) > 0 {
				eObject = l.objects[len(l.objects)-1]
			}
			l.setFeatureValue(eObject, eType.(EStructuralFeature), l.textBuilder.String(), -1)
		}
	}
	l.textBuilder = nil

	// end of the document
	if len(l.elements) == 0 {
		for _, deferred := range l.deferred {
			l.resource.GetContents().Add(deferred)
		}
		l.handleReferences()
		l.recordSchemaLocations(eRoot)
	}

	context := l.namespaces.popContext()
	for _, p := range context {
		delete(l.spacesToFactories, p[1].(string))
	}

}

func (l *XMLDecoder) setAttributes(attributes []xml.Attr) []xml.Attr {
	old := l.attributes
	l.attributes = attributes
	return old
}

func (l *XMLDecoder) getAttributeValue(uri string, local string) string {
	if l.attributes != nil {
		for _, attr := range l.attributes {
			if attr.Name.Space == uri && attr.Name.Local == local {
				return attr.Value
			}
		}
	}
	return ""
}

func (l *XMLDecoder) getXSIType() string {
	return l.getAttributeValue(xsiURI, typeAttrib)
}

func (l *XMLDecoder) handleSchemaLocation() {
	xsiSchemaLocation := l.getAttributeValue(xsiURI, schemaLocationAttrib)
	if len(xsiSchemaLocation) > 0 {
		l.handleXSISchemaLocation(xsiSchemaLocation)
	}

	xsiNoNamespaceSchemaLocation := l.getAttributeValue(xsiURI, noNamespaceSchemaLocationAttrib)
	if len(xsiNoNamespaceSchemaLocation) > 0 {
		l.handleXSINoNamespaceSchemaLocation(xsiNoNamespaceSchemaLocation)
	}
}

func (l *XMLDecoder) handleXSISchemaLocation(loc string) {
}

func (l *XMLDecoder) handleXSINoNamespaceSchemaLocation(loc string) {
}

func (l *XMLDecoder) handlePrefixMapping() {
	if l.attributes != nil {
		for _, attr := range l.attributes {
			if attr.Name.Space == xmlNS {
				l.startPrefixMapping(attr.Name.Local, attr.Value)
			} else if attr.Name.Space == "" && attr.Name.Local == xmlNS {
				l.startPrefixMapping("", attr.Value)
			}
		}
	}
}

func (l *XMLDecoder) startPrefixMapping(prefix string, uri string) {
	l.namespaces.declarePrefix(prefix, uri)
	if _, exists := l.prefixesToURI[prefix]; exists {
		index := 1
		for _, exists = l.prefixesToURI[prefix+"_"+fmt.Sprintf("%d", index)]; exists; {
			index++
		}
		prefix += "_" + fmt.Sprintf("%d", index)
	}
	l.prefixesToURI[prefix] = uri
	delete(l.spacesToFactories, uri)
}

func (l *XMLDecoder) processElement(space string, local string) {
	if len(l.objects) == 0 {
		eObject := l.createTopObject(space, local)
		if eObject != nil {
			if l.deferred != nil {
				l.deferred = append(l.deferred, eObject)
			} else {
				l.attachFn(eObject)
			}
		}
	} else {
		l.handleFeature(space, local)
	}
}

func (l *XMLDecoder) validateObject(eObject EObject, space, typeName string) {
	if eObject == nil {
		l.error(NewEDiagnosticImpl("Class {'"+space+"':'"+typeName+"'} not found", l.resource.GetURI().String(), int(l.decoder.InputOffset()), 0))
	}
}

func (l *XMLDecoder) processObject(eObject EObject) {
	if eObject != nil {
		l.objects = append(l.objects, eObject)
		l.types = append(l.types, load_object_type)
	} else {
		l.types = append(l.types, load_error_type)
	}
}

func (l *XMLDecoder) createTopObject(space string, local string) EObject {
	eFactory := l.getFactoryForSpace(space)
	if eFactory != nil {
		ePackage := eFactory.GetEPackage()
		if l.extendedMetaData != nil && l.extendedMetaData.GetDocumentRoot(ePackage) != nil {
			eClass := l.extendedMetaData.GetDocumentRoot(ePackage)
			// add document root to object list & handle its features
			eObject := l.createObjectWithFactory(eFactory, eClass, false)
			l.processObject(eObject)
			l.handleFeature(space, local)
			if l.isSuppressDocumentRoot {
				// remove document root from object list
				l.objects = l.objects[1:]
				// remove type info from type list
				l.types = l.types[1:]

				// consider new child object as the future new root
				// remove it from document root
				if len(l.objects) > 0 {
					eObject = l.objects[0]
					// remove new object from its container ( document root )
					Remove(eObject)
				}
			}
			return eObject
		} else {
			eType := l.getType(ePackage, local)
			eObject := l.createObjectWithFactory(eFactory, eType, true)
			l.validateObject(eObject, space, local)
			l.processObject(eObject)
			return eObject
		}
	} else {
		prefix, _ := l.namespaces.getPrefix(space)
		l.handleUnknownPackage(prefix, space)
		return nil
	}
}

func (l *XMLDecoder) createObjectWithFactory(eFactory EFactory, eType EClassifier, handleAttributes bool) EObject {
	if eFactory != nil {
		eClass, _ := eType.(EClass)
		if eClass != nil && !eClass.IsAbstract() {
			eObject := eFactory.Create(eClass)
			if eObject != nil && handleAttributes {
				l.interfaces.(xmlDecoderInternal).handleAttributes(eObject)
			}
			return eObject
		}
	}
	return nil
}

func (l *XMLDecoder) createObjectFromFeatureType(eObject EObject, eFeature EStructuralFeature) EObject {
	var eResult EObject
	if eFeature != nil {
		if eType := eFeature.GetEType(); eType != nil {
			eFactory := eType.GetEPackage().GetEFactoryInstance()
			eResult = l.createObjectWithFactory(eFactory, eType, true)
		}
	}
	if eResult != nil {
		l.setFeatureValue(eObject, eFeature, eResult, -1)
	}
	l.processObject(eResult)
	return eResult
}

func (l *XMLDecoder) createObjectFromTypeName(eObject EObject, qname string, eFeature EStructuralFeature) EObject {
	prefix := ""
	local := qname
	if index := strings.Index(qname, ":"); index > 0 {
		prefix = qname[:index]
		local = qname[index+1:]
	}
	space, _ := l.namespaces.getURI(prefix)
	eFactory := l.getFactoryForSpace(space)
	if eFactory == nil {
		l.handleUnknownPackage(prefix, space)
		return nil
	}

	eType := l.getType(eFactory.GetEPackage(), local)
	eResult := l.createObjectWithFactory(eFactory, eType, true)
	l.validateObject(eResult, space, local)
	if eResult != nil {
		l.setFeatureValue(eObject, eFeature, eResult, -1)
	}
	l.processObject(eResult)
	return eResult
}

func (l *XMLDecoder) handleFeature(space string, local string) {
	var eObject EObject
	if len(l.objects) > 0 {
		eObject = l.objects[len(l.objects)-1]
	}
	if eObject != nil {
		eFeature := l.getFeature(eObject, space, local)
		if eFeature != nil {
			if featureKind := l.getLoadFeatureKind(eFeature); featureKind == xlfkSingle || featureKind == xlfkMany {
				l.textBuilder = &strings.Builder{}
				l.types = append(l.types, eFeature)
				l.objects = append(l.objects, nil)
			} else {
				xsiType := l.interfaces.(xmlDecoderInternal).getXSIType()
				if len(xsiType) > 0 {
					l.createObjectFromTypeName(eObject, xsiType, eFeature)
				} else {
					l.createObjectFromFeatureType(eObject, eFeature)
				}
			}
		} else {
			l.handleUnknownFeature(local)
		}
	} else {
		l.types = append(l.types, load_error_type)
		l.handleUnknownFeature(local)
	}

}

func (l *XMLDecoder) setFeatureValue(eObject EObject,
	eFeature EStructuralFeature,
	value any,
	position int) {
	kind := l.getLoadFeatureKind(eFeature)
	switch kind {
	case xlfkSingle:
		eClassifier := eFeature.GetEType()
		eDataType := eClassifier.(EDataType)
		eFactory := eDataType.GetEPackage().GetEFactoryInstance()
		if value == nil {
			eObject.ESet(eFeature, nil)
		} else {
			eObject.ESet(eFeature, eFactory.CreateFromString(eDataType, value.(string)))
		}

	case xlfkMany:
		eClassifier := eFeature.GetEType()
		eDataType := eClassifier.(EDataType)
		eFactory := eDataType.GetEPackage().GetEFactoryInstance()
		eList := eObject.EGetResolve(eFeature, false).(EList)
		if position == -2 {
		} else if value == nil {
			eList.Add(nil)
		} else {
			eList.Add(eFactory.CreateFromString(eDataType, value.(string)))
		}
	case xlfkManyAdd:
		fallthrough
	case xlfkManyMove:
		eList := eObject.EGetResolve(eFeature, false).(EList)
		if position == -1 {
			eList.Add(value)
		} else if position == -2 {
			eList.Clear()
		} else if eObject == value {
			index := eList.IndexOf(value)
			if index == -1 {
				eList.Insert(position, value)
			} else {
				eList.Move(position, index)
			}
		} else if kind == xlfkManyAdd {
			eList.Add(value)
		} else {
			eList.MoveObject(position, value)
		}
	default:
		eObject.ESet(eFeature, value)
	}
}

func (l *XMLDecoder) getLoadFeatureKind(eFeature EStructuralFeature) xmlLoadFeatureKind {
	eClassifier := eFeature.GetEType()
	if eDataType, _ := eClassifier.(EDataType); eDataType != nil {
		if eFeature.IsMany() {
			return xlfkMany
		}
		return xlfkSingle
	} else if eFeature.IsMany() {
		eReference := eFeature.(EReference)
		eOpposite := eReference.GetEOpposite()
		if eOpposite == nil || eOpposite.IsTransient() || !eOpposite.IsMany() {
			return xlfkManyAdd
		}
		return xlfkManyMove
	}
	return xlfkOther
}

func (l *XMLDecoder) handleAttributes(eObject EObject) {
	if l.attributes != nil {
		for _, attr := range l.attributes {
			name := attr.Name.Local
			uri := attr.Name.Space
			value := attr.Value
			if name == l.idAttributeName {
				if idManager := l.resource.GetObjectIDManager(); idManager != nil {
					if err := idManager.SetID(eObject, value); err != nil {
						l.error(NewEDiagnosticImpl(err.Error(), l.resource.GetURI().String(), int(l.decoder.InputOffset()), 0))
					}
				}
			} else if name == href {
				l.handleProxy(eObject, value)
			} else if name != xmlNS && uri != xmlNS && l.isUserAttribute(attr.Name) {
				l.setAttributeValue(eObject, name, value)
			}
		}
	}
}

func (l *XMLDecoder) isUserAttribute(name xml.Name) bool {
	for _, notFeature := range l.notFeatures {
		if notFeature == name {
			return false
		}
	}
	return true
}

func (l *XMLDecoder) getFactoryForSpace(space string) EFactory {
	factory := l.spacesToFactories[space]
	if factory == nil {
		packageRegistry := GetPackageRegistry()
		if l.resource.GetResourceSet() != nil {
			packageRegistry = l.resource.GetResourceSet().GetPackageRegistry()
		}
		factory = packageRegistry.GetFactory(space)
		if factory != nil {
			l.spacesToFactories[space] = factory
		}
	}
	return factory
}

func (l *XMLDecoder) setAttributeValue(eObject EObject, qname string, value string) {
	local := qname
	prefix := ""
	if index := strings.Index(qname, ":"); index > 0 {
		local = qname[index+1:]
		prefix = qname[:index-1]
	}
	space, _ := l.namespaces.getURI(prefix)
	eFeature := l.getFeature(eObject, space, local)
	if eFeature != nil {
		kind := l.getLoadFeatureKind(eFeature)
		if kind == xlfkSingle || kind == xlfkMany {
			l.setFeatureValue(eObject, eFeature, value, -2)
		} else {
			l.setValueFromId(eObject, eFeature.(EReference), value)
		}
	} else {
		l.handleUnknownFeature(local)
	}
}

func (l *XMLDecoder) setValueFromId(eObject EObject, eReference EReference, ids string) {
	mustAdd := l.isResolveDeferred
	mustAddOrNotOppositeIsMany := false
	isFirstID := true
	position := 0
	references := []reference{}
	tokens := strings.Split(ids, " ")
	qName := ""
	for _, token := range tokens {
		id := token
		if index := strings.Index(id, "#"); index != -1 {
			if index == 0 {
				id = id[1:]
			} else {
				oldAttributes := l.setAttributes(nil)
				var eProxy EObject
				if len(qName) == 0 {
					eProxy = l.createObjectFromFeatureType(eObject, eReference)
				} else {
					eProxy = l.createObjectFromTypeName(eObject, qName, eReference)
				}
				l.setAttributes(oldAttributes)
				if eProxy != nil {
					l.handleProxy(eProxy, id)
				}
				l.objects = l.objects[:len(l.objects)-1]
				qName = ""
				position++
				continue
			}
		} else if index := strings.Index(id, ":"); index != -1 {
			qName = id
			continue
		}

		if !l.isResolveDeferred {
			if isFirstID {
				eOpposite := eReference.GetEOpposite()
				if eOpposite != nil {
					mustAdd = eOpposite.IsTransient() || eReference.IsMany()
					mustAddOrNotOppositeIsMany = mustAdd || !eOpposite.IsMany()
				} else {
					mustAdd = true
					mustAddOrNotOppositeIsMany = true
				}
				isFirstID = false
			}

			if mustAddOrNotOppositeIsMany {
				resolved := l.resource.GetEObject(id)
				if resolved != nil {
					l.setFeatureValue(eObject, eReference, resolved, -1)
					qName = ""
					position++
					continue
				}
			}
		}

		if mustAdd {
			references = append(references, reference{object: eObject, feature: eReference, id: id, pos: position})
		}

		qName = ""
		position++
	}

	if position == 0 {
		l.setFeatureValue(eObject, eReference, nil, -2)
	} else {
		l.references = append(l.references, references...)
	}
}

func (l *XMLDecoder) handleProxy(eProxy EObject, id string) {
	resourceURI := l.resource.GetURI()
	uri, ok := ParseURI(id)
	if ok != nil || resourceURI == nil {
		return
	}
	// resolve reference uri
	if !uri.IsAbsolute() {
		uri = resourceURI.Resolve(uri)
	}

	// set object proxy uri
	eProxy.(EObjectInternal).ESetProxyURI(uri)

	if resourceURI.Equals(uri.TrimFragment()) {
		l.sameDocumentProxies = append(l.sameDocumentProxies, eProxy)
	}
}

func (l *XMLDecoder) handleReferences() {
	for _, eProxy := range l.sameDocumentProxies {
		for itRef := eProxy.EClass().GetEAllReferences().Iterator(); itRef.HasNext(); {
			eReference := itRef.Next().(EReference)
			eOpposite := eReference.GetEOpposite()
			if eOpposite != nil && eOpposite.IsChangeable() && eProxy.EIsSet(eReference) {
				resolvedObject := l.resource.GetEObject(eProxy.(EObjectInternal).EProxyURI().Fragment())
				if resolvedObject != nil {
					var proxyHolder EObject
					if eReference.IsMany() {
						value := eProxy.EGet(eReference)
						list := value.(EList)
						proxyHolder = list.Get(0).(EObject)
					} else {
						value := eProxy.EGet(eReference)
						proxyHolder = value.(EObject)
					}

					if eOpposite.IsMany() {
						value := proxyHolder.EGetResolve(eOpposite, false)
						holderContents := value.(EList)
						resolvedIndex := holderContents.IndexOf(resolvedObject)
						if resolvedIndex != -1 {
							proxyIndex := holderContents.IndexOf(eProxy)
							holderContents.Move(proxyIndex, resolvedIndex)
							if proxyIndex > resolvedIndex {
								holderContents.Remove(proxyIndex - 1)
							} else {
								holderContents.Remove(proxyIndex + 1)
							}
							break
						}
					}

					replace := false
					if eReference.IsMany() {
						value := resolvedObject.EGet(eReference)
						list := value.(EList)
						replace = !list.Contains(proxyHolder)
					} else {
						value := resolvedObject.EGet(eReference)
						object := value.(EObject)
						replace = object != proxyHolder
					}

					if replace {
						if eOpposite.IsMany() {
							value := proxyHolder.EGetResolve(eOpposite, false)
							list := value.(EList)
							ndx := list.IndexOf(eProxy)
							list.Set(ndx, resolvedObject)
						} else {
							proxyHolder.ESet(eOpposite, resolvedObject)
						}
					}
					break
				}
			}
		}
	}

	for _, reference := range l.references {
		eObject := l.resource.GetEObject(reference.id)
		if eObject != nil {
			l.setFeatureValue(reference.object, reference.feature, eObject, reference.pos)
		} else {
			l.error(NewEDiagnosticImpl(
				"Unresolved reference '"+reference.id+"'", l.resource.GetURI().String(), int(l.decoder.InputOffset()), 0))
		}
	}
}

func (l *XMLDecoder) recordSchemaLocations(eObject EObject) {
	if l.extendedMetaData != nil && eObject != nil {
		eClass := eObject.EClass()
		if xmlnsPrefixMapFeature := l.extendedMetaData.GetXMLNSPrefixMapFeature(eClass); xmlnsPrefixMapFeature != nil {
			m := eObject.EGet(xmlnsPrefixMapFeature).(EMap)
			for prefix, nsURI := range l.prefixesToURI {
				m.Put(prefix, nsURI)
			}
		}
	}
}

func (l *XMLDecoder) getFeature(eObject EObject, space, name string) EStructuralFeature {
	eClass := eObject.EClass()
	eFeature := eClass.GetEStructuralFeatureFromName(name)
	if eFeature == nil && l.extendedMetaData != nil {
		features := eClass.GetEAllStructuralFeatures()
		for it := features.Iterator(); it.HasNext(); {
			feature := it.Next().(EStructuralFeature)
			if name == l.extendedMetaData.GetName(feature) {
				return feature
			}
		}
	}
	return eFeature
}

func (l *XMLDecoder) getType(ePackage EPackage, name string) EClassifier {
	if l.extendedMetaData != nil {
		return l.extendedMetaData.GetType(ePackage, name)
	}
	return ePackage.GetEClassifier(name)
}

func (l *XMLDecoder) handleUnknownFeature(name string) {
	l.error(NewEDiagnosticImpl("Feature "+name+" not found", l.resource.GetURI().String(), int(l.decoder.InputOffset()), 0))
}

func (l *XMLDecoder) handleUnknownPackage(name, space string) {
	l.error(NewEDiagnosticImpl("Package {'"+name+"'='"+space+"'} not found", l.resource.GetURI().String(), int(l.decoder.InputOffset()), 0))
}

func (l *XMLDecoder) error(diagnostic EDiagnostic) {
	l.errorFn(diagnostic)
}

func (l *XMLDecoder) text(data string) {
	if l.textBuilder != nil {
		l.textBuilder.WriteString(data)
	}
}

func (l *XMLDecoder) comment(comment string) {
}

func (l *XMLDecoder) processingInstruction(procInst xml.ProcInst) {
	if procInst.Target == "xml" {
		content := string(procInst.Inst)
		if ver := l.procInst("version", content); ver != "" {
			l.xmlVersion = ver
		}
		if encoding := l.procInst("encoding", content); encoding != "" {
			l.encoding = encoding
		}
	}
}

func (l *XMLDecoder) procInst(param, s string) string {
	// TODO: this parsing is somewhat lame and not exact.
	// It works for all actual cases, though.
	param = param + "="
	idx := strings.Index(s, param)
	if idx == -1 {
		return ""
	}
	v := s[idx+len(param):]
	if v == "" {
		return ""
	}
	if v[0] != '\'' && v[0] != '"' {
		return ""
	}
	idx = strings.IndexRune(v[1:], rune(v[0]))
	if idx == -1 {
		return ""
	}
	return v[1 : idx+1]
}

func (l *XMLDecoder) directive(directive string) {
}
