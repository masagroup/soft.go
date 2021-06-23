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

	"github.com/ugorji/go/codec"
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
	featureID   int
	eFeature    EStructuralFeature
	featureKind int
	eFactory    EFactory
	eDataType   EDataType
}

type BinaryDecoder struct {
	resource         EResource
	r                io.Reader
	decoder          *codec.Decoder
	baseURI          *URI
	objects          []EObject
	uris             []*URI
	packageData      []*binaryDecoderPackageData
	isResolveProxies bool
}

func NewBinaryDecoder(resource EResource, r io.Reader, options map[string]interface{}) *BinaryDecoder {
	mh := &codec.MsgpackHandle{}
	d := &BinaryDecoder{
		resource:    resource,
		r:           r,
		decoder:     codec.NewDecoder(r, mh),
		objects:     []EObject{},
		uris:        []*URI{},
		packageData: []*binaryDecoderPackageData{},
	}
	if uri := resource.GetURI(); uri != nil && uri.IsAbsolute() {
		d.baseURI = uri
	}
	return d
}

func (d *BinaryDecoder) Decode() {
	defer func() {
		if err, _ := recover().(error); err != nil {
			resourcePath := ""
			if d.resource.GetURI() != nil {
				resourcePath = d.resource.GetURI().String()
			}
			d.resource.GetErrors().Add(NewEDiagnosticImpl(err.Error(), resourcePath, 0, 0))
		}
	}()
	d.decodeSignature()
	d.decodeVersion()

	// objects
	size := d.decodeInt()
	objects := make([]interface{}, size)
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

func (d *BinaryDecoder) decode(v interface{}) {
	if err := d.decoder.Decode(v); err != nil {
		panic(err)
	}
}

func (d *BinaryDecoder) decodeInt() int {
	var i int
	d.decode(&i)
	return i
}

func (d *BinaryDecoder) decodeString() string {
	var s string
	d.decode(&s)
	return s
}

func (d *BinaryDecoder) decodeSignature() {
	signature := make([]byte, 8)
	d.decode(signature)
	if bytes.Compare(signature, binarySignature) != 0 {
		panic(errors.New("Invalid signature for a binary EMF serialization"))
	}
}

func (d *BinaryDecoder) decodeVersion() {
	if version := d.decodeInt(); version != binaryVersion {
		panic(errors.New("Invalid version for binary EMF serialization"))
	}
}

func (d *BinaryDecoder) decodeObject() EObject {
	var iid interface{}
	d.decode(&iid)
	switch id := iid.(type) {
	case int:
		if len(d.objects) <= id {
			var eResult EObject
			eClassData := d.decodeClass()
			eObject := eClassData.eFactory.Create(eClassData.eClass).(EObjectInternal)
			eResult = eObject
			featureID := d.decodeInt() - 1
			if featureID == -2 {
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
	default:
		return nil
	}
}

func (d *BinaryDecoder) decodeObjects(list EList) {
	size := d.decodeInt()
	objects := make([]interface{}, size)
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
		existingObjects := make([]interface{}, existingSize)
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
	case object_container:
		fallthrough
	case object_container_proxy:
		fallthrough
	case object:
		fallthrough
	case object_proxy:
		fallthrough
	case object_containment:
		fallthrough
	case object_containment_proxy:
		eObject.ESetFromID(featureData.featureID, d.decodeObject())
	case object_list:
		fallthrough
	case object_list_proxy:
		fallthrough
	case object_containment_list:
		fallthrough
	case object_containment_list_proxy:
		l := eObject.EGetFromID(featureData.featureID, false).(EList)
		d.decodeObjects(l)
	case data:
		valueStr := d.decodeString()
		value := featureData.eFactory.CreateFromString(featureData.eDataType, valueStr)
		eObject.ESetFromID(featureData.featureID, value)
	case data_list:
		values := []interface{}{}
		var valuesStr []string
		d.decode(&valuesStr)
		for _, valueStr := range valuesStr {
			value := featureData.eFactory.CreateFromString(featureData.eDataType, valueStr)
			values = append(values, value)
		}
		eObject.ESetFromID(featureData.featureID, values)
	case enum:
		eObject.ESetFromID(featureData.featureID, d.decodeInt())
	case date:
		var t time.Time
		d.decode(&t)
		eObject.ESetFromID(featureData.featureID, t)
	case primitive:
		var i interface{}
		d.decode(&i)
		eObject.ESetFromID(featureData.featureID, i)
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
		if d.resource.GetResourceSet() != nil {
			packageRegistry = d.resource.GetResourceSet().GetPackageRegistry()
		}
		ePackage := packageRegistry.GetPackage(nsURI)
		if ePackage == nil {
			ePackage, _ = d.resource.GetResourceSet().GetEObject(uri, true).(EPackage)
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
	var iid interface{}
	d.decode(&iid)
	switch id := iid.(type) {
	case int:
		var uri *URI
		if len(d.uris) <= id {
			// build uri
			uriStr := d.decodeString()
			uri = d.resolveURI(NewURI(uriStr))
			// add it to the uri array
			d.uris = append(d.uris, uri)
		} else {
			uri = d.uris[id]
		}
		d.decode(&uri.Fragment)
		return uri
	default:
		return nil
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
		panic(errors.New(fmt.Sprintf("Unable to find class '%v' in package '%v'", className, ePackage.GetNsURI())))
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
		panic(errors.New(fmt.Sprintf("Unable to find feature '%v' in '%v' EClass", eFeatureName, eClassData.eClass.GetName())))
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
