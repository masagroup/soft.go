package json

import (
	"encoding/json"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/masagroup/soft.go/ecore"
)

// jsonReference is a reference to an object of the document, resolved once it is decoded.
type jsonReference struct {
	object   ecore.EObject
	feature  ecore.EStructuralFeature
	fragment string
	position int // index in a many-valued feature, -1 for a single-valued one
}

// Decoder decodes the JSON written by Encoder: one root object whose "eClass" is first, its
// features by name, and the cross references as {"eClass", "eRef"} objects. It reads the input as
// a stream of tokens.
type Decoder struct {
	resource        ecore.EResource
	reader          io.Reader
	decoder         *json.Decoder
	packageRegistry ecore.EPackageRegistry
	featureKinds    map[ecore.EStructuralFeature]jsonFeatureKind
	references      []jsonReference
	attachFn        func(object ecore.EObject)
	errorFn         func(diagnostic ecore.EDiagnostic)
	idAttributeName string
	failed          bool // a syntax or read error stops the decoding
}

func NewDecoder(resource ecore.EResource, r io.Reader, options map[string]any) *Decoder {
	d := &Decoder{
		resource:     resource,
		reader:       r,
		featureKinds: map[ecore.EStructuralFeature]jsonFeatureKind{},
	}
	if options != nil {
		d.idAttributeName, _ = options[OPTION_ID_ATTRIBUTE_NAME].(string)
	}
	return d
}

func (d *Decoder) DecodeResource() {
	d.attachFn = func(object ecore.EObject) {
		d.resource.GetContents().Add(object)
	}
	d.errorFn = func(diagnostic ecore.EDiagnostic) {
		d.resource.GetErrors().Add(diagnostic)
	}
	d.decodeTopObject()
}

func (d *Decoder) DecodeObject() (eObject ecore.EObject, err error) {
	d.attachFn = func(object ecore.EObject) {
		if eObject == nil {
			eObject = object
		}
	}
	d.errorFn = func(diagnostic ecore.EDiagnostic) {
		if err == nil {
			err = diagnostic
		}
	}
	d.decodeTopObject()
	return
}

func (d *Decoder) decodeTopObject() {
	d.decoder = json.NewDecoder(d.reader)
	d.decoder.UseNumber()
	d.packageRegistry = ecore.GetPackageRegistry()
	if resourceSet := d.resource.GetResourceSet(); resourceSet != nil {
		d.packageRegistry = resourceSet.GetPackageRegistry()
	}
	if eObject := d.decodeObject(nil); eObject != nil {
		d.attachFn(eObject)
	}
	if !d.failed {
		if _, err := d.decoder.Token(); err != io.EOF {
			d.error("Unexpected content after the root object")
		}
	}
	d.resolveReferences()
}

// decodeObject decodes an object, of the class named by its "eClass" key or else of eType.
func (d *Decoder) decodeObject(eType ecore.EClass) ecore.EObject {
	if !d.expectDelim('{') {
		return nil
	}
	var eObject ecore.EObject
	var eClass ecore.EClass
	skip := false
	for first := true; d.decoder.More(); first = false {
		key, ok := d.readString()
		if !ok {
			return nil
		}
		if first {
			if key == "eClass" {
				name, ok := d.readString()
				if !ok {
					return nil
				}
				if eClass = d.getClass(name); eClass == nil {
					skip = true
				}
				if !skip {
					eObject = d.createObject(eClass)
				}
				continue
			}
			if eClass = d.getDefaultClass(eType); eClass == nil {
				skip = true
			} else {
				eObject = d.createObject(eClass)
			}
		}
		if skip {
			d.skipValue()
		} else {
			d.decodeKey(eObject, eClass, key)
		}
		if d.failed {
			return nil
		}
	}
	if !d.expectDelim('}') {
		return nil
	}
	if eObject == nil && !skip {
		// empty object
		if eClass = d.getDefaultClass(eType); eClass != nil {
			eObject = d.createObject(eClass)
		}
	}
	return eObject
}

func (d *Decoder) decodeKey(eObject ecore.EObject, eClass ecore.EClass, key string) {
	if len(d.idAttributeName) > 0 && key == d.idAttributeName {
		if id, ok := d.readString(); ok {
			if idManager := d.resource.GetObjectIDManager(); idManager != nil {
				if err := idManager.SetID(eObject, id); err != nil {
					d.error(err.Error())
				}
			}
		}
		return
	}
	eFeature := eClass.GetEStructuralFeatureFromName(key)
	if eFeature == nil {
		d.error("Unknown feature '" + key + "' in class '" + eClass.GetName() + "'")
		d.skipValue()
		return
	}
	kind, ok := d.featureKinds[eFeature]
	if !ok {
		kind = getJSONCodecFeatureKind(eFeature)
		d.featureKinds[eFeature] = kind
	}
	switch kind {
	case jfkData:
		if str, ok := d.readString(); ok {
			eObject.ESet(eFeature, d.createData(eFeature, str))
		}
	case jfkDataList:
		list := eObject.EGetResolve(eFeature, false).(ecore.EList)
		d.decodeArray(func(int) {
			if str, ok := d.readString(); ok {
				list.Add(d.createData(eFeature, str))
			}
		})
	case jfkObject:
		if child := d.decodeObject(featureClass(eFeature)); child != nil {
			eObject.ESet(eFeature, child)
		}
	case jfkObjectList:
		list := eObject.EGetResolve(eFeature, false).(ecore.EList)
		d.decodeArray(func(int) {
			if child := d.decodeObject(featureClass(eFeature)); child != nil {
				list.Add(child)
			}
		})
	case jfkObjectReference:
		d.decodeReference(eObject, eFeature, -1)
	case jfkObjectReferenceList:
		d.decodeArray(func(position int) {
			d.decodeReference(eObject, eFeature, position)
		})
	default:
		// transient features are not written
		d.skipValue()
	}
}

// decodeArray calls decodeElement for each element of an array, with its index.
func (d *Decoder) decodeArray(decodeElement func(position int)) {
	if !d.expectDelim('[') {
		return
	}
	for position := 0; d.decoder.More(); position++ {
		decodeElement(position)
		if d.failed {
			return
		}
	}
	d.expectDelim(']')
}

// decodeReference decodes a {"eClass", "eRef"} object: a reference to an object of the document
// ("#fragment") is resolved once the document is decoded, any other is a proxy.
func (d *Decoder) decodeReference(eObject ecore.EObject, eFeature ecore.EStructuralFeature, position int) {
	if !d.expectDelim('{') {
		return
	}
	className := ""
	reference := ""
	for d.decoder.More() {
		key, ok := d.readString()
		if !ok {
			return
		}
		switch key {
		case "eClass":
			className, _ = d.readString()
		case "eRef":
			reference, _ = d.readString()
		default:
			d.error("Unknown key '" + key + "' in the reference of feature '" + eFeature.GetName() + "'")
			d.skipValue()
		}
		if d.failed {
			return
		}
	}
	if !d.expectDelim('}') {
		return
	}

	if fragment, ok := strings.CutPrefix(reference, "#"); ok {
		d.references = append(d.references, jsonReference{object: eObject, feature: eFeature, fragment: fragment, position: position})
		return
	}

	eClass := featureClass(eFeature)
	if len(className) > 0 {
		if eClass = d.getClass(className); eClass == nil {
			return
		}
	}
	if eClass == nil || eClass.IsAbstract() || eClass.IsInterface() {
		d.error("Unable to create the proxy of reference '" + reference + "'")
		return
	}
	uri, err := ecore.ParseURI(reference)
	if err != nil {
		d.error("Invalid reference '" + reference + "': " + err.Error())
		return
	}
	if resourceURI := d.resource.GetURI(); resourceURI != nil && !uri.IsAbsolute() {
		uri = resourceURI.Resolve(uri)
	}
	eProxy := d.createObject(eClass)
	eProxy.(ecore.EObjectInternal).ESetProxyURI(uri)
	if position < 0 {
		eObject.ESet(eFeature, eProxy)
	} else {
		eObject.EGetResolve(eFeature, false).(ecore.EList).Add(eProxy)
	}
}

// resolveReferences sets the references to the objects of the document, at their position in
// many-valued features.
func (d *Decoder) resolveReferences() {
	for _, reference := range d.references {
		eTarget := d.getEObject(reference.fragment)
		if eTarget == nil {
			d.error("Unresolved reference '" + reference.fragment + "'")
			continue
		}
		if reference.position < 0 {
			reference.object.ESet(reference.feature, eTarget)
			continue
		}
		list := reference.object.EGetResolve(reference.feature, false).(ecore.EList)
		// the opposite of a previous reference may have added the target already
		if index := list.IndexOf(eTarget); index != -1 {
			if index != reference.position && reference.position < list.Size() {
				list.Move(index, reference.position)
			}
		} else if reference.position <= list.Size() {
			list.Insert(reference.position, eTarget)
		} else {
			list.Add(eTarget)
		}
	}
	d.references = nil
}

func (d *Decoder) getClass(name string) ecore.EClass {
	index := strings.LastIndex(name, "#//")
	if index == -1 {
		d.error("Invalid class '" + name + "'")
		return nil
	}
	nsURI, className := name[:index], name[index+3:]
	ePackage := d.packageRegistry.GetPackage(nsURI)
	if ePackage == nil {
		d.error("Unknown package '" + nsURI + "'")
		return nil
	}
	eClass, _ := ePackage.GetEClassifier(className).(ecore.EClass)
	if eClass == nil {
		d.error("Unknown class '" + className + "' in package '" + nsURI + "'")
		return nil
	}
	return eClass
}

// getDefaultClass returns eType for an object without "eClass", if it can be instantiated.
func (d *Decoder) getDefaultClass(eType ecore.EClass) ecore.EClass {
	if eType == nil || eType.IsAbstract() || eType.IsInterface() {
		d.error("Missing eClass")
		return nil
	}
	return eType
}

func (d *Decoder) createObject(eClass ecore.EClass) ecore.EObject {
	return ecore.GetRegistryFactory(d.packageRegistry, eClass.GetEPackage()).Create(eClass)
}

func (d *Decoder) createData(eFeature ecore.EStructuralFeature, str string) any {
	eDataType := eFeature.GetEType().(ecore.EDataType)
	return eDataType.GetEPackage().GetEFactoryInstance().CreateFromString(eDataType, str)
}

func featureClass(eFeature ecore.EStructuralFeature) ecore.EClass {
	eClass, _ := eFeature.GetEType().(ecore.EClass)
	return eClass
}

// readString reads a string; a number or a boolean is read as its literal.
func (d *Decoder) readString() (string, bool) {
	token, ok := d.token()
	if !ok {
		return "", false
	}
	switch v := token.(type) {
	case string:
		return v, true
	case json.Number:
		return v.String(), true
	case bool:
		return strconv.FormatBool(v), true
	}
	d.unexpected(token, "a string")
	return "", false
}

func (d *Decoder) expectDelim(delim json.Delim) bool {
	token, ok := d.token()
	if !ok {
		return false
	}
	if token != delim {
		d.unexpected(token, "'"+delim.String()+"'")
		return false
	}
	return true
}

// skipValue reads a value and drops it.
func (d *Decoder) skipValue() {
	depth := 0
	for {
		token, ok := d.token()
		if !ok {
			return
		}
		if delim, isDelim := token.(json.Delim); isDelim {
			if delim == '{' || delim == '[' {
				depth++
			} else {
				depth--
			}
		}
		if depth == 0 {
			return
		}
	}
}

// unexpected reports a token that is not the one expected: the decoding stops, as the stream is no
// longer where the decoder expects it to be.
func (d *Decoder) unexpected(token json.Token, expected string) {
	found := "null"
	switch v := token.(type) {
	case json.Delim:
		found = "'" + v.String() + "'"
	case string:
		found = "string '" + v + "'"
	case json.Number:
		found = "number " + v.String()
	case bool:
		found = strconv.FormatBool(v)
	}
	d.error("Expected " + expected + ", found " + found)
	d.failed = true
}

func (d *Decoder) token() (json.Token, bool) {
	if d.failed {
		return nil, false
	}
	token, err := d.decoder.Token()
	if err != nil {
		if errors.Is(err, io.EOF) {
			err = io.ErrUnexpectedEOF
		}
		d.error(err.Error())
		d.failed = true
		return nil, false
	}
	return token, true
}

func (d *Decoder) error(message string) {
	location := ""
	if uri := d.resource.GetURI(); uri != nil {
		location = uri.String()
	}
	offset := 0
	if d.decoder != nil {
		offset = int(d.decoder.InputOffset())
	}
	d.errorFn(ecore.NewEDiagnosticImpl(message, location, 0, offset))
}

// getEObject returns the object of the document at fragment, nil if there is none: a fragment
// naming an unknown feature makes the resource panic.
func (d *Decoder) getEObject(fragment string) (eObject ecore.EObject) {
	defer func() {
		if recover() != nil {
			eObject = nil
		}
	}()
	return d.resource.GetEObject(fragment)
}
