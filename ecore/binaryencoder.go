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
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

type checkType int

const (
	checkNothing checkType = iota
	checkDirectResource
	checkResource
	checkContainer
)

var binaryVersion int = 0

var binarySignature []byte = []byte{'\211', 'e', 'm', 'f', '\n', '\r', '\032', '\n'}

type binaryEncoderPackageData struct {
	id        int
	classData []*binaryEncoderClassData
}

type binaryEncoderClassData struct {
	packageID   int
	id          int
	featureData []*binaryEncoderFeatureData
}

type binaryEncoderFeatureData struct {
	name        string
	isTransient bool
	featureKind binaryFeatureKind
	factory     EFactory
	dataType    EDataType
}

type BinaryEncoder struct {
	w                    io.Writer
	resource             EResource
	encoder              *msgpack.Encoder
	objectRoot           EObject
	baseURI              *URI
	version              int
	objectToID           map[EObject]int
	classDataMap         map[EClass]*binaryEncoderClassData
	packageDataMap       map[EPackage]*binaryEncoderPackageData
	uriToIDMap           map[string]int
	enumLiteralToIDMap   map[string]int
	isIDAttributeEncoded bool
}

func NewBinaryEncoder(resource EResource, w io.Writer, options map[string]interface{}) *BinaryEncoder {
	return NewBinaryEncoderWithVersion(resource, w, options, binaryVersion)
}

func NewBinaryEncoderWithVersion(resource EResource, w io.Writer, options map[string]interface{}, version int) *BinaryEncoder {
	e := &BinaryEncoder{
		w:                  w,
		resource:           resource,
		encoder:            msgpack.NewEncoder(w),
		version:            version,
		objectToID:         map[EObject]int{},
		classDataMap:       map[EClass]*binaryEncoderClassData{},
		packageDataMap:     map[EPackage]*binaryEncoderPackageData{},
		uriToIDMap:         map[string]int{},
		enumLiteralToIDMap: map[string]int{},
	}
	if uri := resource.GetURI(); uri != nil {
		e.baseURI = uri
	}
	if options != nil {
		e.isIDAttributeEncoded = options[BINARY_OPTION_ID_ATTRIBUTE] == true
	}
	return e
}

func (e *BinaryEncoder) Encode() {
	if !binaryDebug {
		defer func() {
			if err, _ := recover().(error); err != nil {
				resourcePath := ""
				if e.resource.GetURI() != nil {
					resourcePath = e.resource.GetURI().String()
				}
				e.resource.GetErrors().Add(NewEDiagnosticImpl(err.Error(), resourcePath, 0, 0))
			}
		}()
	}
	e.encodeSignature()
	e.encodeVersion()
	e.encodeObjects(e.resource.GetContents(), checkContainer)
}

func (e *BinaryEncoder) EncodeObject(object EObject) (err error) {
	defer func() {
		if panicErr, _ := recover().(error); panicErr != nil {
			err = panicErr
		}
	}()
	e.objectRoot = object
	e.encodeSignature()
	e.encodeVersion()
	e.encodeObject(object, checkContainer)
	return
}

func (e *BinaryEncoder) haltOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func (e *BinaryEncoder) encodeInt(i int) {
	e.haltOnError(e.encoder.EncodeInt(int64(i)))
}

func (e *BinaryEncoder) encodeString(s string) {
	e.haltOnError(e.encoder.EncodeString(s))
}

func (e *BinaryEncoder) encodeBytes(bytes []byte) {
	e.haltOnError(e.encoder.EncodeBytes(bytes))
}

func (e *BinaryEncoder) encodeDate(date *time.Time) {
	e.haltOnError(e.encoder.EncodeTime(*date))
}

func (e *BinaryEncoder) encodeFloat64(f float64) {
	e.haltOnError(e.encoder.EncodeFloat64(f))
}

func (e *BinaryEncoder) encodeFloat32(f float32) {
	e.haltOnError(e.encoder.EncodeFloat32(f))
}

func (e *BinaryEncoder) encodeInt64(i int64) {
	e.haltOnError(e.encoder.EncodeInt64(i))
}

func (e *BinaryEncoder) encodeInt32(i int32) {
	e.haltOnError(e.encoder.EncodeInt32(i))
}

func (e *BinaryEncoder) encodeInt16(i int16) {
	e.haltOnError(e.encoder.EncodeInt16(i))
}

func (e *BinaryEncoder) encodeInt8(i int8) {
	e.haltOnError(e.encoder.EncodeInt8(i))
}

func (e *BinaryEncoder) encodeByte(b byte) {
	e.haltOnError(e.encoder.EncodeUint8(b))
}

func (e *BinaryEncoder) encode(i interface{}) {
	e.haltOnError(e.encoder.Encode(i))
}

func (e *BinaryEncoder) encodeBool(b bool) {
	e.haltOnError(e.encoder.EncodeBool(b))
}

func (e *BinaryEncoder) encodeSignature() {
	// Write a signature that will e obviously corrupt
	// if the binary contents end up eing UTF-8 encoded
	// or altered by line feed or carriage return changes.
	e.encodeBytes(binarySignature)
}

func (e *BinaryEncoder) encodeVersion() {
	e.encodeInt(e.version)
}

func (e *BinaryEncoder) encodeObjects(objects EList, check checkType) {
	e.encodeInt(objects.Size())
	for it := objects.Iterator(); it.HasNext(); {
		eObject, _ := it.Next().(EObject)
		e.encodeObject(eObject, check)
	}
}

func (e *BinaryEncoder) encodeObject(eObject EObject, check checkType) {
	if eObject == nil {
		e.encodeInt(-1)
	} else if id, isID := e.objectToID[eObject]; isID {
		e.encodeInt(id)
	} else if eObjectInternal, _ := eObject.(EObjectInternal); eObjectInternal != nil {
		// object id
		objectID := len(e.objectToID)
		e.objectToID[eObject] = objectID
		e.encodeInt(objectID)

		// object class
		eClass := eObject.EClass()
		eClassData := e.encodeClass(eClass)

		// object uri if reference or proxy
		saveFeatureValues := true
		switch check {
		case checkDirectResource:
			if eObjectInternal.EIsProxy() {
				e.encodeInt(-2)
				e.encodeURI(eObjectInternal.EProxyURI())
				saveFeatureValues = false
			} else if eResource := eObjectInternal.EInternalResource(); eResource != nil {
				e.encodeInt(-2)
				e.encodeURIWithFragment(eResource.GetURI(), eResource.GetURIFragment(eObjectInternal))
				saveFeatureValues = false
			}
		case checkResource:
			if eObjectInternal.EIsProxy() {
				e.encodeInt(-2)
				e.encodeURI(eObjectInternal.EProxyURI())
				saveFeatureValues = false
			} else if eResource := eObjectInternal.EResource(); eResource != nil &&
				(eResource != e.resource ||
					(e.objectRoot != nil && !IsAncestor(e.objectRoot, eObjectInternal))) {
				// encode object as uri and fragment if object is in a different resource
				// or if in the same resource and root object is not its ancestor
				e.encodeInt(-2)
				e.encodeURIWithFragment(eResource.GetURI(), eResource.GetURIFragment(eObjectInternal))
				saveFeatureValues = false
			}
		case checkNothing:
		case checkContainer:
		}
		// object feature values
		if saveFeatureValues {

			// id attribute
			if objectIDManager := e.resource.GetObjectIDManager(); e.isIDAttributeEncoded && objectIDManager != nil {
				if id := objectIDManager.GetID(eObject); id != nil {
					e.encodeInt(-1)
					e.encode(id)
				}
			}

			// features
			for featureID, featureData := range eClassData.featureData {
				if !featureData.isTransient && (check == checkContainer || featureData.featureKind != bfkObjectContainerProxy) {
					e.encodeFeatureValue(eObjectInternal, featureID, featureData)
				}
			}

		}

		e.encodeInt(0)
	}
}

func (e *BinaryEncoder) encodeFeatureValue(eObject EObjectInternal, featureID int, featureData *binaryEncoderFeatureData) {
	if eObject.EIsSetFromID(featureID) {
		e.encodeInt(featureID + 1)
		if len(featureData.name) > 0 {
			e.encodeString(featureData.name)
			featureData.name = ""
		}
		value := eObject.EGetFromID(featureID, false)
		switch featureData.featureKind {
		case bfkObject:
			fallthrough
		case bfkObjectContainment:
			e.encodeObject(value.(EObject), checkNothing)
		case bfkObjectContainerProxy:
			e.encodeObject(value.(EObject), checkResource)
		case bfkObjectContainmentProxy:
			e.encodeObject(value.(EObject), checkDirectResource)
		case bfkObjectProxy:
			e.encodeObject(value.(EObject), checkResource)
		case bfkObjectList:
			fallthrough
		case bfkObjectContainmentList:
			e.encodeObjects(value.(EList), checkNothing)
		case bfkObjectContainmentListProxy:
			e.encodeObjects(value.(EList), checkDirectResource)
		case bfkObjectListProxy:
			e.encodeObjects(value.(EList), checkResource)
		case bfkData:
			valueStr := featureData.factory.ConvertToString(featureData.dataType, value)
			e.encodeString(valueStr)
		case bfkDataList:
			l := value.(EList)
			e.encodeInt(l.Size())
			for it := l.Iterator(); it.HasNext(); {
				value := it.Next()
				valueStr := featureData.factory.ConvertToString(featureData.dataType, value)
				e.encodeString(valueStr)
			}
		case bfkEnum:
			literalStr := featureData.factory.ConvertToString(featureData.dataType, value)
			if enumID, isID := e.enumLiteralToIDMap[literalStr]; isID {
				e.encodeInt(enumID)
			} else {
				enumID := len(e.enumLiteralToIDMap)
				e.enumLiteralToIDMap[literalStr] = enumID
				e.encodeInt(enumID)
				e.encodeString(literalStr)
			}
		case bfkDate:
			e.encodeDate(value.(*time.Time))
		case bfkFloat64:
			e.encodeFloat64(value.(float64))
		case bfkFloat32:
			e.encodeFloat32(value.(float32))
		case bfkInt:
			e.encodeInt(value.(int))
		case bfkInt64:
			e.encodeInt64(value.(int64))
		case bfkInt32:
			e.encodeInt32(value.(int32))
		case bfkInt16:
			e.encodeInt16(value.(int16))
		case bfkByte:
			e.encodeByte(value.(byte))
		case bfkBool:
			e.encodeBool(value.(bool))
		case bfkString:
			e.encodeString(value.(string))
		case bfkByteArray:
			e.encodeBytes(value.([]byte))
		default:
			panic(fmt.Sprintf("feature with feature kind '%v' is not supported", featureData.featureKind))
		}
	}
}

func (e *BinaryEncoder) encodeClass(eClass EClass) *binaryEncoderClassData {
	eClassData, _ := e.classDataMap[eClass]
	if eClassData != nil {
		e.encodeInt(eClassData.packageID)
		e.encodeInt(eClassData.id)
	} else {
		eClassData = e.newClassData(eClass)
		e.encodeInt(eClassData.id)
		e.encodeString(eClass.GetName())
		e.classDataMap[eClass] = eClassData
	}
	return eClassData
}

func (e *BinaryEncoder) encodePackage(ePackage EPackage) *binaryEncoderPackageData {
	ePackageData, _ := e.packageDataMap[ePackage]
	if ePackageData != nil {
		e.encodeInt(ePackageData.id)
	} else {
		ePackageData = e.newPackageData(ePackage)
		e.encodeInt(ePackageData.id)
		e.encodeString(ePackage.GetNsURI())
		e.encodeURI(GetURI(ePackage))
		e.packageDataMap[ePackage] = ePackageData
	}
	return ePackageData
}

func (e *BinaryEncoder) encodeURI(uri *URI) {
	if uri == nil {
		e.encodeInt(-1)
	} else {
		e.encodeURIWithFragment(uri.TrimFragment(), uri.Fragment)
	}
}

func (e *BinaryEncoder) encodeURIWithFragment(uri *URI, fragment string) {
	if uri == nil {
		e.encodeInt(-1)
	} else {
		uriPath := uri.String()
		if id, isID := e.uriToIDMap[uriPath]; isID {
			e.encodeInt(id)
		} else {
			id := len(e.uriToIDMap)
			e.uriToIDMap[uriPath] = id
			e.encodeInt(id)
			e.encodeString(e.relativizeURI(uri).String())
		}
		e.encodeString(fragment)
	}
}

func (e *BinaryEncoder) relativizeURI(uri *URI) *URI {
	if e.baseURI != nil {
		return e.baseURI.Relativize(uri)
	}
	return uri
}

func (e *BinaryEncoder) newPackageData(ePackage EPackage) *binaryEncoderPackageData {
	return &binaryEncoderPackageData{
		id:        len(e.packageDataMap),
		classData: make([]*binaryEncoderClassData, ePackage.GetEClassifiers().Size()),
	}
}

func (e *BinaryEncoder) newClassID(ePackageData *binaryEncoderPackageData) int {
	for i, c := range ePackageData.classData {
		if c == nil {
			return i
		}
	}
	return -1
}

func (e *BinaryEncoder) newClassData(eClass EClass) *binaryEncoderClassData {
	ePackageData := e.encodePackage(eClass.GetEPackage())
	eClassData := &binaryEncoderClassData{
		packageID:   ePackageData.id,
		id:          e.newClassID(ePackageData),
		featureData: []*binaryEncoderFeatureData{},
	}
	ePackageData.classData[eClassData.id] = eClassData
	for it := eClass.GetEAllStructuralFeatures().Iterator(); it.HasNext(); {
		eFeature := it.Next().(EStructuralFeature)
		eClassData.featureData = append(eClassData.featureData, e.newFeatureData(eFeature))
	}
	return eClassData
}

func (e *BinaryEncoder) newFeatureData(eFeature EStructuralFeature) *binaryEncoderFeatureData {
	eFeatureData := &binaryEncoderFeatureData{
		name:        eFeature.GetName(),
		featureKind: getBinaryCodecFeatureKind(eFeature),
	}
	if eReference, _ := eFeature.(EReference); eReference != nil {
		eFeatureData.isTransient = eReference.IsTransient() || (eReference.IsContainer() && !eReference.IsResolveProxies())
	} else if eAttribute, _ := eFeature.(EAttribute); eAttribute != nil {
		eDataType := eAttribute.GetEAttributeType()
		eFeatureData.isTransient = eAttribute.IsTransient()
		eFeatureData.dataType = eDataType
		eFeatureData.factory = eDataType.GetEPackage().GetEFactoryInstance()
	}
	return eFeatureData
}
