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

var binaryVersion int

var binarySignature = []byte{'\211', 'e', 'm', 'f', '\n', '\r', '\032', '\n'}

type binaryEncoderPackageData struct {
	classData []*binaryEncoderClassData
	id        int
}

type binaryEncoderClassData struct {
	featureData []*binaryEncoderFeatureData
	packageID   int
	id          int
}

type binaryEncoderFeatureData struct {
	factory     EFactory
	dataType    EDataType
	name        string
	featureKind binaryFeatureKind
	isTransient bool
}

type BinaryEncoder struct {
	w                    io.Writer
	resource             EResource
	objectRoot           EObject
	encoder              *msgpack.Encoder
	baseURI              *URI
	objectToID           map[EObject]int
	classDataMap         map[EClass]*binaryEncoderClassData
	packageDataMap       map[EPackage]*binaryEncoderPackageData
	uriToIDMap           map[string]int
	enumLiteralToIDMap   map[string]int
	version              int
	isIDAttributeEncoded bool
}

func NewBinaryEncoder(resource EResource, w io.Writer, options map[string]any) *BinaryEncoder {
	return NewBinaryEncoderWithVersion(resource, w, options, binaryVersion)
}

func NewBinaryEncoderWithVersion(resource EResource, w io.Writer, options map[string]any, version int) *BinaryEncoder {
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

func (e *BinaryEncoder) EncodeResource() {
	var err error
	defer func() {
		if err != nil {
			// add error to resource errors
			resourcePath := ""
			if e.resource.GetURI() != nil {
				resourcePath = e.resource.GetURI().String()
			}
			e.resource.GetErrors().Add(NewEDiagnosticImpl(err.Error(), resourcePath, 0, 0))
		}
	}()

	if err = e.encodeSignature(); err != nil {
		return
	}
	if err = e.encodeVersion(); err != nil {
		return
	}
	err = e.encodeObjects(e.resource.GetContents(), checkContainer)
}

func (e *BinaryEncoder) EncodeObject(object EObject) error {
	e.objectRoot = object
	if err := e.encodeSignature(); err != nil {
		return err
	}
	if err := e.encodeVersion(); err != nil {
		return err
	}
	return e.encodeObject(object, checkContainer)
}

func (e *BinaryEncoder) encodeInt(i int) error {
	return e.encoder.EncodeInt(int64(i))
}

func (e *BinaryEncoder) encodeString(s string) error {
	return e.encoder.EncodeString(s)
}

func (e *BinaryEncoder) encodeBytes(bytes []byte) error {
	return e.encoder.EncodeBytes(bytes)
}

func (e *BinaryEncoder) encodeDate(date *time.Time) error {
	return e.encoder.EncodeTime(*date)
}

func (e *BinaryEncoder) encodeFloat64(f float64) error {
	return e.encoder.EncodeFloat64(f)
}

func (e *BinaryEncoder) encodeFloat32(f float32) error {
	return e.encoder.EncodeFloat32(f)
}

func (e *BinaryEncoder) encodeInt64(i int64) error {
	return e.encoder.EncodeInt64(i)
}

func (e *BinaryEncoder) encodeInt32(i int32) error {
	return e.encoder.EncodeInt32(i)
}

func (e *BinaryEncoder) encodeInt16(i int16) error {
	return e.encoder.EncodeInt16(i)
}

func (e *BinaryEncoder) encodeByte(b byte) error {
	return e.encoder.EncodeUint8(b)
}

func (e *BinaryEncoder) encode(i any) error {
	return e.encoder.Encode(i)
}

func (e *BinaryEncoder) encodeBool(b bool) error {
	return e.encoder.EncodeBool(b)
}

func (e *BinaryEncoder) encodeSignature() error {
	// Write a signature that will e obviously corrupt
	// if the binary contents end up eing UTF-8 encoded
	// or altered by line feed or carriage return changes.
	return e.encodeBytes(binarySignature)
}

func (e *BinaryEncoder) encodeVersion() error {
	return e.encodeInt(e.version)
}

func (e *BinaryEncoder) encodeObjects(objects EList, check checkType) error {
	if err := e.encodeInt(objects.Size()); err != nil {
		return err
	}
	for it := objects.Iterator(); it.HasNext(); {
		eObject, _ := it.Next().(EObject)
		if err := e.encodeObject(eObject, check); err != nil {
			return err
		}
	}
	return nil
}

func (e *BinaryEncoder) encodeObject(eObject EObject, check checkType) error {
	if eObject == nil {
		return e.encodeInt(-1)
	} else if id, isID := e.objectToID[eObject]; isID {
		return e.encodeInt(id)
	} else if eObjectInternal, _ := eObject.(EObjectInternal); eObjectInternal != nil {
		// object id
		objectID := len(e.objectToID)
		e.objectToID[eObject] = objectID
		if err := e.encodeInt(objectID); err != nil {
			return err
		}

		// object class
		eClass := eObject.EClass()
		eClassData, err := e.encodeClass(eClass)
		if err != nil {
			return err
		}

		// object uri if reference or proxy
		saveFeatureValues := true
		switch check {
		case checkDirectResource:
			if eObjectInternal.EIsProxy() {
				if err := e.encodeInt(-2); err != nil {
					return err
				}
				if err := e.encodeURI(eObjectInternal.EProxyURI()); err != nil {
					return err
				}
				saveFeatureValues = false
			} else if eResource := eObjectInternal.EInternalResource(); eResource != nil {
				if err := e.encodeInt(-2); err != nil {
					return err
				}
				if err := e.encodeURIWithFragment(eResource.GetURI(), eResource.GetURIFragment(eObjectInternal)); err != nil {
					return err
				}
				saveFeatureValues = false
			}
		case checkResource:
			if eObjectInternal.EIsProxy() {
				if err := e.encodeInt(-2); err != nil {
					return err
				}
				if err := e.encodeURI(eObjectInternal.EProxyURI()); err != nil {
					return err
				}
				saveFeatureValues = false
			} else if eResource := eObjectInternal.EResource(); eResource != nil &&
				(eResource != e.resource ||
					(e.objectRoot != nil && !IsAncestor(e.objectRoot, eObjectInternal))) {
				// encode object as uri and fragment if object is in a different resource
				// or if in the same resource and root object is not its ancestor
				if err := e.encodeInt(-2); err != nil {
					return err
				}
				if err := e.encodeURIWithFragment(eResource.GetURI(), eResource.GetURIFragment(eObjectInternal)); err != nil {
					return err
				}
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
					if err := e.encodeInt(-1); err != nil {
						return err
					}
					if err := e.encode(id); err != nil {
						return err
					}
				}
			}

			// features
			for featureID, featureData := range eClassData.featureData {
				if !featureData.isTransient && (check == checkContainer || featureData.featureKind != bfkObjectContainerProxy) {
					if err := e.encodeFeatureValue(eObjectInternal, featureID, featureData); err != nil {
						return err
					}
				}
			}

		}
		return e.encodeInt(0)
	}
	return nil
}

func (e *BinaryEncoder) encodeFeatureValue(eObject EObjectInternal, featureID int, featureData *binaryEncoderFeatureData) error {
	if eObject.EIsSetFromID(featureID) {
		if err := e.encodeInt(featureID + 1); err != nil {
			return err
		}
		if len(featureData.name) > 0 {
			if err := e.encodeString(featureData.name); err != nil {
				return err
			}
			featureData.name = ""
		}
		value := eObject.EGetFromID(featureID, false)
		switch featureData.featureKind {
		case bfkObject:
			fallthrough
		case bfkObjectContainment:
			return e.encodeObject(value.(EObject), checkNothing)
		case bfkObjectContainerProxy:
			return e.encodeObject(value.(EObject), checkResource)
		case bfkObjectContainmentProxy:
			return e.encodeObject(value.(EObject), checkDirectResource)
		case bfkObjectProxy:
			return e.encodeObject(value.(EObject), checkResource)
		case bfkObjectList:
			fallthrough
		case bfkObjectContainmentList:
			return e.encodeObjects(value.(EList), checkNothing)
		case bfkObjectContainmentListProxy:
			return e.encodeObjects(value.(EList), checkDirectResource)
		case bfkObjectListProxy:
			return e.encodeObjects(value.(EList), checkResource)
		case bfkData:
			valueStr := featureData.factory.ConvertToString(featureData.dataType, value)
			return e.encodeString(valueStr)
		case bfkDataList:
			l := value.(EList)
			if err := e.encodeInt(l.Size()); err != nil {
				return err
			}
			for it := l.Iterator(); it.HasNext(); {
				value := it.Next()
				valueStr := featureData.factory.ConvertToString(featureData.dataType, value)
				if err := e.encodeString(valueStr); err != nil {
					return err
				}
			}
		case bfkEnum:
			literalStr := featureData.factory.ConvertToString(featureData.dataType, value)
			if enumID, isID := e.enumLiteralToIDMap[literalStr]; isID {
				return e.encodeInt(enumID)
			} else {
				enumID := len(e.enumLiteralToIDMap)
				e.enumLiteralToIDMap[literalStr] = enumID
				if err := e.encodeInt(enumID); err != nil {
					return err
				}
				return e.encodeString(literalStr)
			}
		case bfkDate:
			return e.encodeDate(value.(*time.Time))
		case bfkFloat64:
			return e.encodeFloat64(value.(float64))
		case bfkFloat32:
			return e.encodeFloat32(value.(float32))
		case bfkInt:
			return e.encodeInt(value.(int))
		case bfkInt64:
			return e.encodeInt64(value.(int64))
		case bfkInt32:
			return e.encodeInt32(value.(int32))
		case bfkInt16:
			return e.encodeInt16(value.(int16))
		case bfkByte:
			return e.encodeByte(value.(byte))
		case bfkBool:
			return e.encodeBool(value.(bool))
		case bfkString:
			return e.encodeString(value.(string))
		case bfkByteArray:
			return e.encodeBytes(value.([]byte))
		default:
			return fmt.Errorf("feature with feature kind '%v' is not supported", featureData.featureKind)
		}
	}
	return nil
}

func (e *BinaryEncoder) encodeClass(eClass EClass) (*binaryEncoderClassData, error) {
	if eClassData := e.classDataMap[eClass]; eClassData != nil {
		if err := e.encodeInt(eClassData.packageID); err != nil {
			return nil, err
		}
		if err := e.encodeInt(eClassData.id); err != nil {
			return nil, err
		}
		return eClassData, nil
	} else {
		eClassData, err := e.newClassData(eClass)
		if err != nil {
			return nil, err
		}
		if err := e.encodeInt(eClassData.id); err != nil {
			return nil, err
		}
		if err := e.encodeString(eClass.GetName()); err != nil {
			return nil, err
		}
		e.classDataMap[eClass] = eClassData
		return eClassData, nil
	}
}

func (e *BinaryEncoder) encodePackage(ePackage EPackage) (*binaryEncoderPackageData, error) {
	ePackageData := e.packageDataMap[ePackage]
	if ePackageData != nil {
		if err := e.encodeInt(ePackageData.id); err != nil {
			return nil, err
		}
	} else {
		ePackageData = e.newPackageData(ePackage)
		if err := e.encodeInt(ePackageData.id); err != nil {
			return nil, err
		}
		if err := e.encodeString(ePackage.GetNsURI()); err != nil {
			return nil, err
		}
		if err := e.encodeURI(GetURI(ePackage)); err != nil {
			return nil, err
		}
		e.packageDataMap[ePackage] = ePackageData
	}
	return ePackageData, nil
}

func (e *BinaryEncoder) encodeURI(uri *URI) error {
	if uri == nil {
		return e.encodeInt(-1)
	} else {
		return e.encodeURIWithFragment(uri.TrimFragment(), uri.Fragment())
	}
}

func (e *BinaryEncoder) encodeURIWithFragment(uri *URI, fragment string) error {
	if uri == nil {
		return e.encodeInt(-1)
	} else {
		uriPath := uri.String()
		if id, isID := e.uriToIDMap[uriPath]; isID {
			if err := e.encodeInt(id); err != nil {
				return err
			}
		} else {
			id := len(e.uriToIDMap)
			e.uriToIDMap[uriPath] = id
			if err := e.encodeInt(id); err != nil {
				return err
			}
			if err := e.encodeString(e.relativizeURI(uri).String()); err != nil {
				return err
			}
		}
		return e.encodeString(fragment)
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

func (e *BinaryEncoder) newClassData(eClass EClass) (*binaryEncoderClassData, error) {
	eFeatures := eClass.GetEAllStructuralFeatures()
	ePackageData, err := e.encodePackage(eClass.GetEPackage())
	if err != nil {
		return nil, err
	}
	eClassData := &binaryEncoderClassData{
		packageID:   ePackageData.id,
		id:          e.newClassID(ePackageData),
		featureData: make([]*binaryEncoderFeatureData, 0, eFeatures.Size()),
	}
	ePackageData.classData[eClassData.id] = eClassData
	for it := eFeatures.Iterator(); it.HasNext(); {
		eFeature := it.Next().(EStructuralFeature)
		eClassData.featureData = append(eClassData.featureData, e.newFeatureData(eFeature))
	}
	return eClassData, nil
}

const extendedMetaData = "http:///org/eclipse/emf/ecore/util/ExtendedMetaData"

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
	// if we have a xmlns prefix map or schema loction map, consider it as non transient to ensure
	// information is kept between binary encoder and xml encoder
	if eAnnotation := eFeature.GetEAnnotation(extendedMetaData); eAnnotation != nil {
		if name := eAnnotation.GetDetails().GetValue("name"); name == "xmlns:prefix" || name == "xsi:schemaLocation" {
			eFeatureData.isTransient = false
		}
	}
	return eFeatureData
}
