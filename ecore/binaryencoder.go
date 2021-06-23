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

	"github.com/ugorji/go/codec"
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
	errorFn        func(diagnostic EDiagnostic)
}

func NewBinaryEncoder(resource EResource, w io.Writer, options map[string]interface{}) *BinaryEncoder {
	return NewBinaryEncoderWithVersion(resource, w, options, binaryVersion)
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
	return e
}

func (e *BinaryEncoder) Encode() {
	defer func() {
		if err, _ := recover().(error); err != nil {
			resourcePath := ""
			if e.resource.GetURI() != nil {
				resourcePath = e.resource.GetURI().String()
			}
			e.resource.GetErrors().Add(NewEDiagnosticImpl(err.Error(), resourcePath, 0, 0))
		}
	}()
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
	e.encodeSignature()
	e.encodeVersion()
	e.encodeObject(object, checkContainer)
	return
}

func (e *BinaryEncoder) encode(v interface{}) {
	if err := e.encoder.Encode(v); err != nil {
		panic(err)
	}
}

func (e *BinaryEncoder) encodeSignature() {
	// Write a signature that will e obviously corrupt
	// if the binary contents end up eing UTF-8 encoded
	// or altered by line feed or carriage return changes.
	e.encode(binarySignature)
}

func (e *BinaryEncoder) encodeVersion() {
	e.encode(e.version)
}

func (e *BinaryEncoder) encodeObjects(objects EList, check checkType) {
	e.encode(objects.Size())
	for it := objects.Iterator(); it.HasNext(); {
		switch eObject := it.Next().(type) {
		case nil:
			e.encodeObject(nil, check)
		case EObject:
			e.encodeObject(eObject, check)
		}
	}
}

func (e *BinaryEncoder) encodeObject(eObject EObject, check checkType) {
	if eObject == nil {
		e.encode(nil)
	} else if id, isID := e.objectToID[eObject]; isID {
		e.encode(id)
	} else {
		// object id
		var objectID interface{}
		if objectIDManager := e.resource.GetObjectIDManager(); objectIDManager != nil {
			objectID = objectIDManager.GetID(eObject)
		} else {
			id := len(e.objectToID)
			e.objectToID[eObject] = id
			objectID = id
		}
		e.encode(objectID)
		// object class
		eClass := eObject.EClass()
		eClassData := e.encodeClass(eClass)

		saveFeatureValues := true
		eObjectInternal, _ := eObject.(EObjectInternal)
		if eObjectInternal == nil {
			return
		}
		// object uri if reference or proxy
		switch check {
		case checkDirectResource:
			if eResource := eObjectInternal.EInternalResource(); eResource != nil {
				e.encode(-1)
				e.encodeURIWithFragment(eResource.GetURI(), eResource.GetURIFragment(eObjectInternal))
				saveFeatureValues = false
			} else if eObjectInternal.EIsProxy() {
				e.encode(-1)
				e.encodeURI(eObjectInternal.EProxyURI())
				saveFeatureValues = false
			}
		case checkResource:
			if eResource := eObjectInternal.EResource(); eResource != nil && eResource != e.resource {
				e.encode(-1)
				e.encodeURIWithFragment(eResource.GetURI(), eResource.GetURIFragment(eObjectInternal))
				saveFeatureValues = false
			} else if eObjectInternal.EIsProxy() {
				e.encode(-1)
				e.encodeURI(eObjectInternal.EProxyURI())
				saveFeatureValues = false
			}
		case checkNothing:
		case checkContainer:
		}
		// object feature values
		for featureID, featureData := range eClassData.featureData {
			if saveFeatureValues && !featureData.isTransient && (check == checkContainer || featureData.featureKind != object_container_proxy) {
				e.encodeFeatureValue(eObjectInternal, featureID, featureData)
			}
		}
		e.encode(0)
	}
}

func (e *BinaryEncoder) encodeFeatureValue(eObject EObjectInternal, featureID int, featureData *binaryEncoderFeatureData) {
	if eObject.EIsSetFromID(featureID) {
		e.encode(featureID + 1)
		if len(featureData.name) > 0 {
			e.encode(featureData.name)
			featureData.name = ""
		}
		value := eObject.EGetFromID(featureID, false)
		switch featureData.featureKind {
		case object:
			fallthrough
		case object_containment:
			e.encodeObject(value.(EObject), checkNothing)
		case object_container_proxy:
			e.encodeObject(value.(EObject), checkResource)
		case object_containment_proxy:
			e.encodeObject(value.(EObject), checkDirectResource)
		case object_proxy:
			e.encodeObject(value.(EObject), checkResource)
		case object_list:
			fallthrough
		case object_containment_list:
			e.encodeObjects(value.(EList), checkNothing)
		case object_containment_list_proxy:
			e.encodeObjects(value.(EList), checkDirectResource)
		case object_list_proxy:
			e.encodeObjects(value.(EList), checkResource)
		case data:
			valueStr := featureData.factory.ConvertToString(featureData.dataType, value)
			e.encode(valueStr)
		case data_list:
			valuesStr := []string{}
			for it := value.(EList).Iterator(); it.HasNext(); {
				value := it.Next()
				valueStr := featureData.factory.ConvertToString(featureData.dataType, value)
				valuesStr = append(valuesStr, valueStr)
			}
			e.encode(valuesStr)
		case enum:
			e.encode(value)
		case date:
			e.encode(value)
		case primitive:
			e.encode(value)
		}
	}
}

func (e *BinaryEncoder) encodeClass(eClass EClass) *binaryEncoderClassData {
	eClassData, _ := e.classDataMap[eClass]
	if eClassData != nil {
		e.encode(eClassData.packageID)
		e.encode(eClassData.id)
	} else {
		eClassData = e.newClassData(eClass)
		e.encode(eClassData.id)
		e.encode(eClass.GetName())
		e.classDataMap[eClass] = eClassData
	}
	return eClassData
}

func (e *BinaryEncoder) encodePackage(ePackage EPackage) *binaryEncoderPackageData {
	ePackageData, _ := e.packageDataMap[ePackage]
	if ePackageData != nil {
		e.encode(ePackageData.id)
	} else {
		ePackageData = e.newPackageData(ePackage)
		e.encode(ePackageData.id)
		e.encode(ePackage.GetNsURI())
		e.encodeURI(GetURI(ePackage))
		e.packageDataMap[ePackage] = ePackageData
	}
	return ePackageData
}

func (e *BinaryEncoder) encodeURI(uri *URI) {
	if uri == nil {
		e.encode(nil)
	} else {
		e.encodeURIWithFragment(uri.TrimFragment(), uri.Fragment)
	}
}

func (e *BinaryEncoder) encodeURIWithFragment(uri *URI, fragment string) {
	if uri == nil {
		e.encode(nil)
	} else {
		uriPath := uri.String()
		if id, isID := e.uriToIDMap[uriPath]; isID {
			e.encode(id)
		} else {
			id := len(e.uriToIDMap)
			e.uriToIDMap[uriPath] = id
			e.encode(id)
			e.encode(e.relativizeURI(uri).String())
		}
		e.encode(fragment)
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
