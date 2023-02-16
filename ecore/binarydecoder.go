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
	"bytes"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

type binaryDecoderPackageData struct {
	ePackage   EPackage
	eClassData []*binaryDecoderClassData
}

type binaryDecoderClassData struct {
	eClass      EClass
	eFactory    EFactory
	featureData []*binaryDecoderFeatureData
}

type binaryDecoderFeatureData struct {
	eFeature    EStructuralFeature
	eFactory    EFactory
	eDataType   EDataType
	featureID   int
	featureKind binaryFeatureKind
}

type BinaryDecoder struct {
	resource         EResource
	r                io.Reader
	decoder          *msgpack.Decoder
	baseURI          *URI
	objects          []EObject
	uris             []*URI
	packageData      []*binaryDecoderPackageData
	enumLiterals     []string
	isResolveProxies bool
}

func NewBinaryDecoder(resource EResource, r io.Reader, options map[string]any) *BinaryDecoder {
	d := &BinaryDecoder{
		resource:     resource,
		r:            r,
		decoder:      msgpack.NewDecoder(r),
		objects:      []EObject{},
		uris:         []*URI{},
		packageData:  []*binaryDecoderPackageData{},
		enumLiterals: []string{},
	}
	if uri := resource.GetURI(); uri != nil {
		d.baseURI = uri
	}
	return d
}

func (d *BinaryDecoder) Decode() {
	if !binaryDebug {
		defer func() {
			if err, _ := recover().(error); err != nil {
				resourcePath := ""
				if d.resource.GetURI() != nil {
					resourcePath = d.resource.GetURI().String()
				}
				d.resource.GetErrors().Add(NewEDiagnosticImpl(err.Error(), resourcePath, 0, 0))
			}
		}()
	}
	d.decodeSignature()
	d.decodeVersion()

	// objects
	size := d.decodeInt()
	objects := make([]any, size)
	for i := 0; i < size; i++ {
		objects[i] = d.decodeObject()
	}

	// add objects to resource
	d.resource.GetContents().AddAll(NewImmutableEList(objects))
}

func (d *BinaryDecoder) DecodeObject() (eObject EObject, err error) {
	defer func() {
		if panicErr, _ := recover().(error); panicErr != nil {
			err = panicErr
		}
	}()
	d.decodeSignature()
	d.decodeVersion()
	eObject = d.decodeObject()
	return
}

func (e *BinaryDecoder) haltOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func (d *BinaryDecoder) decodeInt() int {
	i, err := d.decoder.DecodeInt()
	d.haltOnError(err)
	return i
}

func (d *BinaryDecoder) decodeInt64() int64 {
	i, err := d.decoder.DecodeInt64()
	d.haltOnError(err)
	return i
}

func (d *BinaryDecoder) decodeInt32() int32 {
	i, err := d.decoder.DecodeInt32()
	d.haltOnError(err)
	return i
}

func (d *BinaryDecoder) decodeInt16() int16 {
	i, err := d.decoder.DecodeInt16()
	d.haltOnError(err)
	return i
}

func (d *BinaryDecoder) decodeByte() byte {
	i, err := d.decoder.DecodeUint8()
	d.haltOnError(err)
	return i
}

func (d *BinaryDecoder) decodeBool() bool {
	b, err := d.decoder.DecodeBool()
	d.haltOnError(err)
	return b
}

func (d *BinaryDecoder) decodeString() string {
	s, err := d.decoder.DecodeString()
	d.haltOnError(err)
	return s
}

func (d *BinaryDecoder) decodeBytes() []byte {
	bytes, err := d.decoder.DecodeBytes()
	d.haltOnError(err)
	return bytes
}

func (d *BinaryDecoder) decodeDate() *time.Time {
	t, err := d.decoder.DecodeTime()
	d.haltOnError(err)
	// msgpack time is decoded as local
	// got to transform as UTC
	t = t.UTC()
	return &t
}

func (d *BinaryDecoder) decodeFloat64() float64 {
	f, err := d.decoder.DecodeFloat64()
	d.haltOnError(err)
	return f
}

func (d *BinaryDecoder) decodeFloat32() float32 {
	f, err := d.decoder.DecodeFloat32()
	d.haltOnError(err)
	return f
}

func (d *BinaryDecoder) decodeInterface() any {
	i, err := d.decoder.DecodeInterface()
	d.haltOnError(err)
	return i
}

func (d *BinaryDecoder) decodeSignature() {
	signature := d.decodeBytes()
	if !bytes.Equal(signature, binarySignature) {
		panic(errors.New("invalid signature for a binary emf serialization"))
	}
}

func (d *BinaryDecoder) decodeVersion() {
	if version := d.decodeInt(); version != binaryVersion {
		panic(errors.New("invalid version for binary emf serialization"))
	}
}

func (d *BinaryDecoder) decodeObject() EObject {
	id := d.decodeInt()
	switch id {
	case -1:
		return nil
	default:
		if len(d.objects) <= int(id) {
			var eResult EObject
			eClassData := d.decodeClass()
			eObject := eClassData.eFactory.Create(eClassData.eClass).(EObjectInternal)
			eResult = eObject
			featureID := d.decodeInt() - 1

			if featureID == -3 {
				// proxy object
				eProxyURI := d.decodeURI()
				eObject.ESetProxyURI(eProxyURI)
				if d.isResolveProxies {
					eResult = ResolveInResource(eObject, d.resource)
					d.objects = append(d.objects, eResult)
				} else {
					d.objects = append(d.objects, eObject)
				}
				featureID = d.decodeInt() - 1
			} else {
				// standard object
				d.objects = append(d.objects, eObject)
			}

			if featureID == -2 {
				// object id attribute
				objectID := d.decodeInterface()
				if objectIDManager := d.resource.GetObjectIDManager(); objectIDManager != nil {
					err := objectIDManager.SetID(eObject, objectID)
					if err != nil {
						panic(err)
					}
				}
				featureID = d.decodeInt() - 1
			}

			for ; featureID != -1; featureID = d.decodeInt() - 1 {
				eFeatureData := eClassData.featureData[featureID]
				if eFeatureData == nil {
					eFeatureData = d.newFeatureData(eClassData, featureID)
					eClassData.featureData[featureID] = eFeatureData
				}
				d.decodeFeatureValue(eObject, eFeatureData)
			}
			return eResult
		} else {
			return d.objects[id]
		}
	}
}

func (d *BinaryDecoder) decodeObjects(list EList) {
	size := d.decodeInt()
	objects := make([]any, size)
	for i := 0; i < size; i++ {
		objects[i] = d.decodeObject()
	}

	// If the list is empty, we need to add all the objects,
	// otherwise, the reference is bidirectional and the list is at least partially populated.
	existingSize := list.Size()
	if existingSize == 0 {
		list.AddAll(NewImmutableEList(objects))
	} else {
		indices := make([]int, existingSize)
		duplicateCount := 0
		existingObjects := make([]any, existingSize)
		copy(existingObjects, list.ToArray())
	LOOP:
		for i := 0; i < size; i++ {
			o := objects[i]
			count := duplicateCount
			for j := 0; j < existingSize; j++ {
				existing := existingObjects[j]
				if existing == o {
					if duplicateCount != count {
						list.Move(count, duplicateCount)
					}
					indices[duplicateCount] = i
					duplicateCount++
					existingObjects[j] = nil
					continue LOOP
				} else if existing != nil {
					count++
				}
			}
			objects[i-duplicateCount] = o
		}

		size -= existingSize
		list.AddAll(NewImmutableEList(objects))
		for i := 0; i < existingSize; i++ {
			newPosition := indices[i]
			oldPosition := size + i
			if newPosition != oldPosition {
				list.Move(oldPosition, newPosition)
			}
		}
	}
}

func (d *BinaryDecoder) decodeFeatureValue(eObject EObjectInternal, featureData *binaryDecoderFeatureData) {
	switch featureData.featureKind {
	case bfkObjectContainer:
		fallthrough
	case bfkObjectContainerProxy:
		fallthrough
	case bfkObject:
		fallthrough
	case bfkObjectProxy:
		fallthrough
	case bfkObjectContainment:
		fallthrough
	case bfkObjectContainmentProxy:
		eObject.ESetFromID(featureData.featureID, d.decodeObject())
	case bfkObjectList:
		fallthrough
	case bfkObjectListProxy:
		fallthrough
	case bfkObjectContainmentList:
		fallthrough
	case bfkObjectContainmentListProxy:
		l := eObject.EGetFromID(featureData.featureID, false).(EList)
		d.decodeObjects(l)
	case bfkData:
		valueStr := d.decodeString()
		value := featureData.eFactory.CreateFromString(featureData.eDataType, valueStr)
		eObject.ESetFromID(featureData.featureID, value)
	case bfkDataList:
		size := d.decodeInt()
		values := []any{}
		for i := 0; i < size; i++ {
			valueStr := d.decodeString()
			value := featureData.eFactory.CreateFromString(featureData.eDataType, valueStr)
			values = append(values, value)
		}
		l := eObject.EGetResolve(featureData.eFeature, false).(EList)
		l.AddAll(NewBasicEList(values))
	case bfkEnum:
		var valueStr string
		id := d.decodeInt()
		if len(d.enumLiterals) <= id {
			valueStr = d.decodeString()
			d.enumLiterals = append(d.enumLiterals, valueStr)
		} else {
			valueStr = d.enumLiterals[id]
		}
		value := featureData.eFactory.CreateFromString(featureData.eDataType, valueStr)
		eObject.ESetFromID(featureData.featureID, value)
	case bfkDate:
		eObject.ESetFromID(featureData.featureID, d.decodeDate())
	case bfkFloat64:
		eObject.ESetFromID(featureData.featureID, d.decodeFloat64())
	case bfkFloat32:
		eObject.ESetFromID(featureData.featureID, d.decodeFloat32())
	case bfkInt:
		eObject.ESetFromID(featureData.featureID, d.decodeInt())
	case bfkInt64:
		eObject.ESetFromID(featureData.featureID, d.decodeInt64())
	case bfkInt32:
		eObject.ESetFromID(featureData.featureID, d.decodeInt32())
	case bfkInt16:
		eObject.ESetFromID(featureData.featureID, d.decodeInt16())
	case bfkByte:
		eObject.ESetFromID(featureData.featureID, d.decodeByte())
	case bfkBool:
		eObject.ESetFromID(featureData.featureID, d.decodeBool())
	case bfkString:
		eObject.ESetFromID(featureData.featureID, d.decodeString())
	case bfkByteArray:
		eObject.ESetFromID(featureData.featureID, d.decodeBytes())
	}
}

func (d *BinaryDecoder) decodeClass() *binaryDecoderClassData {
	ePackageData := d.decodePackage()
	id := d.decodeInt()
	eClassData := ePackageData.eClassData[id]
	if eClassData == nil {
		eClassData = d.newClassData(ePackageData)
		ePackageData.eClassData[id] = eClassData
	}
	return eClassData
}

func (d *BinaryDecoder) decodePackage() *binaryDecoderPackageData {
	id := d.decodeInt()
	if len(d.packageData) <= id {
		// decode package parameters
		nsURI := d.decodeString()
		uri := d.decodeURI()

		// retrieve package
		packageRegistry := GetPackageRegistry()
		resourceSet := d.resource.GetResourceSet()
		if resourceSet != nil {
			packageRegistry = resourceSet.GetPackageRegistry()
		}
		ePackage := packageRegistry.GetPackage(nsURI)
		if ePackage == nil && resourceSet != nil {
			ePackage, _ = resourceSet.GetEObject(uri, true).(EPackage)
		}
		if ePackage == nil {
			panic(fmt.Errorf("unable to find package '%s'", nsURI))
		}
		// create new package data
		ePackageData := d.newPackageData(ePackage)
		d.packageData = append(d.packageData, ePackageData)
		return ePackageData

	} else {
		return d.packageData[id]
	}
}

func (d *BinaryDecoder) decodeURI() *URI {
	id := d.decodeInt()
	switch id {
	case -1:
		return nil
	default:
		var uri *URI
		if len(d.uris) <= int(id) {
			// build uri
			uriStr := d.decodeString()
			if uriStr == "" {
				uri = d.baseURI
			} else {
				uri = d.resolveURI(NewURI(uriStr))
			}
			// add it to the uri array
			d.uris = append(d.uris, uri)
		} else {
			uri = d.uris[id]
		}
		return NewURIBuilder(uri).SetFragment(d.decodeString()).URI()
	}
}

func (d *BinaryDecoder) resolveURI(uri *URI) *URI {
	if d.baseURI != nil {
		return d.baseURI.Resolve(uri)
	}
	return uri
}

func (d *BinaryDecoder) newPackageData(ePackage EPackage) *binaryDecoderPackageData {
	return &binaryDecoderPackageData{
		ePackage:   ePackage,
		eClassData: make([]*binaryDecoderClassData, ePackage.GetEClassifiers().Size()),
	}
}

func (d *BinaryDecoder) newClassData(ePackageData *binaryDecoderPackageData) *binaryDecoderClassData {
	className := d.decodeString()
	ePackage := ePackageData.ePackage
	eClass, _ := ePackage.GetEClassifier(className).(EClass)
	if eClass == nil {
		panic(fmt.Errorf("unable to find class '%v' in package '%v'", className, ePackage.GetNsURI()))
	}
	return &binaryDecoderClassData{
		eClass:      eClass,
		eFactory:    ePackage.GetEFactoryInstance(),
		featureData: make([]*binaryDecoderFeatureData, eClass.GetFeatureCount()),
	}
}

func (d *BinaryDecoder) newFeatureData(eClassData *binaryDecoderClassData, featureID int) *binaryDecoderFeatureData {
	eFeatureName := d.decodeString()
	eFeature := eClassData.eClass.GetEStructuralFeatureFromName(eFeatureName)
	if eFeature == nil {
		panic(fmt.Errorf("unable to find feature '%v' in '%v' EClass", eFeatureName, eClassData.eClass.GetName()))
	}
	eFeatureData := &binaryDecoderFeatureData{
		eFeature:    eFeature,
		featureID:   featureID,
		featureKind: getBinaryCodecFeatureKind(eFeature),
	}
	if eAttribute, _ := eFeature.(EAttribute); eAttribute != nil {
		eFeatureData.eDataType = eAttribute.GetEAttributeType()
		eFeatureData.eFactory = eFeatureData.eDataType.GetEPackage().GetEFactoryInstance()
	}
	return eFeatureData
}
