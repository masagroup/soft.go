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
	"io"
	"time"

	"github.com/ugorji/go/codec"
)

type checkType int

const (
	checkNothing checkType = iota
	checkDirectResource
	checkResource
	checkContainer
)

const encoderVersion = 0

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
	featureKind int
	factory     EFactory
	dataType    EDataType
}

type BinaryEncoder struct {
	w              io.Writer
	resource       EResource
	encoder        *codec.Encoder
	baseURI        *URI
	version        int
	objectToID     map[EObject]int
	classDataMap   map[EClass]*binaryEncoderClassData
	packageDataMap map[EPackage]*binaryEncoderPackageData
	uriToIDMap     map[string]int
}

func NewBinaryEncoder(resource EResource, w io.Writer, options map[string]interface{}) *BinaryEncoder {
	return NewBinaryEncoderWithVersion(resource, w, options, encoderVersion)
}

func NewBinaryEncoderWithVersion(resource EResource, w io.Writer, options map[string]interface{}, version int) *BinaryEncoder {
	mh := &codec.MsgpackHandle{}
	e := &BinaryEncoder{
		w:              w,
		resource:       resource,
		encoder:        codec.NewEncoder(w, mh),
		version:        version,
		objectToID:     map[EObject]int{},
		classDataMap:   map[EClass]*binaryEncoderClassData{},
		packageDataMap: map[EPackage]*binaryEncoderPackageData{},
		uriToIDMap:     map[string]int{},
	}
	if uri := resource.GetURI(); uri != nil && uri.IsAbsolute() {
		e.baseURI = uri
	}
	e.encodeVersion()
	e.encodeSignature()
	return e
}

func (be *BinaryEncoder) Encode() {
	contents := be.resource.GetContents()
	be.encodeObjects(contents, checkContainer)
}

func (be *BinaryEncoder) EncodeObject(object EObject) error {
	return nil
}

func (be *BinaryEncoder) encodeSignature() {
	// Write a signature that will be obviously corrupt
	// if the binary contents end up being UTF-8 encoded
	// or altered by line feed or carriage return changes.
	be.encoder.Encode([]byte{'\211', 'e', 'm', 'f', '\n', '\r', '\032', '\n'})
}

func (be *BinaryEncoder) encodeVersion() {
	be.encoder.Encode(be.version)
}

func (be *BinaryEncoder) encodeObjects(objects EList, check checkType) {
	be.encoder.Encode(objects.Size())
	for it := objects.Iterator(); it.HasNext(); {
		switch eObject := it.Next().(type) {
		case nil:
			be.encodeObject(nil, check)
		case EObject:
			be.encodeObject(eObject, check)
		}
	}
}

func (be *BinaryEncoder) encodeObject(eObject EObject, check checkType) {
	if eObject == nil {
		be.encoder.Encode(nil)
	} else if id, isID := be.objectToID[eObject]; isID {
		be.encoder.Encode(id)
	} else {
		// object id
		var objectID interface{}
		if objectIDManager := be.resource.GetObjectIDManager(); objectIDManager != nil {
			objectID = objectIDManager.GetID(eObject)
		} else {
			id := len(be.objectToID)
			be.objectToID[eObject] = id
			objectID = id
		}
		be.encoder.Encode(objectID)
		// object class
		eClass := eObject.EClass()
		eClassData := be.encodeClass(eClass)

		saveFeatureValues := true
		eObjectInternal, _ := eObject.(EObjectInternal)
		if eObjectInternal == nil {
			return
		}
		// object uri if reference or proxy
		switch check {
		case checkDirectResource:
			if eResource := eObjectInternal.EInternalResource(); eResource != nil {
				be.encoder.Encode(0)
				be.encodeURIWithFragment(eResource.GetURI(), eResource.GetURIFragment(eObjectInternal))
				saveFeatureValues = false
			} else if eObjectInternal.EIsProxy() {
				be.encoder.Encode(0)
				be.encodeURI(eObjectInternal.EProxyURI())
				saveFeatureValues = false
			}
		case checkResource:
			if eResource := eObjectInternal.EResource(); eResource != nil && eResource != be.resource {
				be.encoder.Encode(0)
				be.encodeURIWithFragment(eResource.GetURI(), eResource.GetURIFragment(eObjectInternal))
				saveFeatureValues = false
			} else if eObjectInternal.EIsProxy() {
				be.encoder.Encode(0)
				be.encodeURI(eObjectInternal.EProxyURI())
				saveFeatureValues = false
			}
		case checkNothing:
		case checkContainer:
		}
		// object feature values
		for featureID, featureData := range eClassData.featureData {
			if saveFeatureValues && !featureData.isTransient && (check == checkContainer || featureData.featureKind != object_container_proxy) {
				be.encodeFeatureValue(eObjectInternal, featureID, featureData)
			}
		}
		be.encoder.Encode(1)
	}
}

func (be *BinaryEncoder) encodeFeatureValue(eObject EObjectInternal, featureID int, featureData *binaryEncoderFeatureData) {
	if eObject.EIsSetFromID(featureID) {
		be.encoder.Encode(featureID)
		if len(featureData.name) > 0 {
			be.encoder.Encode(featureData.name)
			featureData.name = ""
		}
		value := eObject.EGetFromID(featureID, false)
		switch featureData.featureKind {
		case object:
			fallthrough
		case object_containment:
			be.encodeObject(value.(EObject), checkNothing)
		case object_container_proxy:
			be.encodeObject(value.(EObject), checkResource)
		case object_containment_proxy:
			be.encodeObject(value.(EObject), checkDirectResource)
		case object_proxy:
			be.encodeObject(value.(EObject), checkResource)
		case object_list:
			fallthrough
		case object_containment_list:
			be.encodeObjects(value.(EList), checkNothing)
		case object_containment_list_proxy:
			be.encodeObjects(value.(EList), checkDirectResource)
		case object_list_proxy:
			be.encodeObjects(value.(EList), checkResource)
		case data:
			valueStr := featureData.factory.ConvertToString(featureData.dataType, value)
			be.encoder.Encode(valueStr)
		case data_list:
			valuesStr := []string{}
			for it := value.(EList).Iterator(); it.HasNext(); {
				value := it.Next()
				valueStr := featureData.factory.ConvertToString(featureData.dataType, value)
				valuesStr = append(valuesStr, valueStr)
			}
			be.encoder.Encode(valuesStr)
		case enum:
			be.encoder.Encode(value)
		case date:
			t := value.(*time.Time)
			be.encoder.Encode(t.Unix())
		case primitive:
			be.encoder.Encode(value)
		}
	}
}

func (be *BinaryEncoder) encodeClass(eClass EClass) *binaryEncoderClassData {
	eClassData, _ := be.classDataMap[eClass]
	if eClassData != nil {
		be.encoder.Encode(eClassData.packageID)
		be.encoder.Encode(eClassData.id)
	} else {
		eClassData = be.newClassData(eClass)
		be.encoder.Encode(eClassData.id)
		be.encoder.Encode(eClass.GetName())
		be.classDataMap[eClass] = eClassData
	}
	return eClassData
}

func (be *BinaryEncoder) encodePackage(ePackage EPackage) *binaryEncoderPackageData {
	ePackageData, _ := be.packageDataMap[ePackage]
	if ePackageData != nil {
		be.encoder.Encode(ePackageData.id)
	} else {
		ePackageData = be.newPackageData(ePackage)
		be.encoder.Encode(ePackageData.id)
		be.encoder.Encode(ePackage.GetNsURI())
		be.encodeURI(GetURI(ePackage))
		be.packageDataMap[ePackage] = ePackageData
	}
	return ePackageData
}

func (be *BinaryEncoder) encodeURI(uri *URI) {
	if uri == nil {
		be.encoder.Encode(nil)
	} else {
		be.encodeURIWithFragment(uri.TrimFragment(), uri.Fragment)
	}
}

func (be *BinaryEncoder) encodeURIWithFragment(uri *URI, fragment string) {
	if uri == nil {
		be.encoder.Encode(nil)
	} else {
		uriPath := uri.String()
		if id, isID := be.uriToIDMap[uriPath]; isID {
			be.encoder.Encode(id)
		} else {
			id := len(be.uriToIDMap)
			be.uriToIDMap[uriPath] = id
			be.encoder.Encode(id)
			be.encoder.Encode(be.relativizeURI(uri).String())
		}
		be.encoder.Encode(fragment)
	}
}

func (be *BinaryEncoder) relativizeURI(uri *URI) *URI {
	if be.baseURI != nil {
		return be.baseURI.Relativize(uri)
	}
	return uri
}

func (be *BinaryEncoder) newPackageData(ePackage EPackage) *binaryEncoderPackageData {
	return &binaryEncoderPackageData{
		id:        len(be.packageDataMap),
		classData: make([]*binaryEncoderClassData, ePackage.GetEClassifiers().Size()),
	}
}

func (be *BinaryEncoder) newClassID(ePackageData *binaryEncoderPackageData) int {
	for i, c := range ePackageData.classData {
		if c == nil {
			return i
		}
	}
	return -1
}

func (be *BinaryEncoder) newClassData(eClass EClass) *binaryEncoderClassData {
	ePackageData := be.encodePackage(eClass.GetEPackage())
	eClassData := &binaryEncoderClassData{
		packageID:   ePackageData.id,
		id:          be.newClassID(ePackageData),
		featureData: []*binaryEncoderFeatureData{},
	}
	for it := eClass.GetEAllStructuralFeatures().Iterator(); it.HasNext(); {
		eFeature := it.Next().(EStructuralFeature)
		eClassData.featureData = append(eClassData.featureData, be.newFeatureData(eFeature))
	}
	return eClassData
}

func (be *BinaryEncoder) newFeatureData(eFeature EStructuralFeature) *binaryEncoderFeatureData {
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
