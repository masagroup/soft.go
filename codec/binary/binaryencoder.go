// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package binary

import (
	"fmt"
	"io"
	"time"

	"github.com/masagroup/soft.go/ecore"
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
	factory     ecore.EFactory
	dataType    ecore.EDataType
	name        string
	featureKind binaryFeatureKind
	isTransient bool
}

type Encoder struct {
	w                           io.Writer
	resource                    ecore.EResource
	objectRoot                  ecore.EObject
	encoder                     *msgpack.Encoder
	baseURI                     *ecore.URI
	objectToID                  map[ecore.EObject]int
	classDataMap                map[ecore.EClass]*binaryEncoderClassData
	packageDataMap              map[ecore.EPackage]*binaryEncoderPackageData
	uriToIDMap                  map[string]int
	enumLiteralToIDMap          map[string]int
	version                     int
	isIDAttributeEncoded        bool
	isNamespaceAttributeEncoded bool
}

func NewEncoder(resource ecore.EResource, w io.Writer, options map[string]any) *Encoder {
	return NewEncoderWithVersion(resource, w, options, binaryVersion)
}

func NewEncoderWithVersion(resource ecore.EResource, w io.Writer, options map[string]any, version int) *Encoder {
	e := &Encoder{
		w:                           w,
		resource:                    resource,
		encoder:                     msgpack.NewEncoder(w),
		version:                     version,
		objectToID:                  map[ecore.EObject]int{},
		classDataMap:                map[ecore.EClass]*binaryEncoderClassData{},
		packageDataMap:              map[ecore.EPackage]*binaryEncoderPackageData{},
		uriToIDMap:                  map[string]int{},
		enumLiteralToIDMap:          map[string]int{},
		isNamespaceAttributeEncoded: true,
	}
	if uri := resource.GetURI(); uri != nil {
		e.baseURI = uri
	}
	if options != nil {
		e.isIDAttributeEncoded = options[OPTION_ID_ATTRIBUTE] == true
		isNamespaceAttributeEncoded, isNamespaceAttibuteDefined := options[OPTION_NAMESPACE_ATTRIBUTE]
		e.isNamespaceAttributeEncoded = (isNamespaceAttibuteDefined && isNamespaceAttributeEncoded == true) || !isNamespaceAttibuteDefined
	}
	return e
}

func (e *Encoder) EncodeResource() {
	var err error
	defer func() {
		if err != nil {
			// add error to resource errors
			resourcePath := ""
			if e.resource.GetURI() != nil {
				resourcePath = e.resource.GetURI().String()
			}
			e.resource.GetErrors().Add(ecore.NewEDiagnosticImpl(err.Error(), resourcePath, 0, 0))
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

func (e *Encoder) EncodeObject(object ecore.EObject) error {
	e.objectRoot = object
	if err := e.encodeSignature(); err != nil {
		return err
	}
	if err := e.encodeVersion(); err != nil {
		return err
	}
	return e.encodeObject(object, checkContainer)
}

func (e *Encoder) encodeInt(i int) error {
	return e.encoder.EncodeInt(int64(i))
}

func (e *Encoder) encodeString(s string) error {
	return e.encoder.EncodeString(s)
}

func (e *Encoder) encodeBytes(bytes []byte) error {
	return e.encoder.EncodeBytes(bytes)
}

func (e *Encoder) encodeDate(date *time.Time) error {
	return e.encoder.EncodeTime(date.UTC())
}

func (e *Encoder) encodeFloat64(f float64) error {
	return e.encoder.EncodeFloat64(f)
}

func (e *Encoder) encodeFloat32(f float32) error {
	return e.encoder.EncodeFloat32(f)
}

func (e *Encoder) encodeInt64(i int64) error {
	return e.encoder.EncodeInt64(i)
}

func (e *Encoder) encodeInt32(i int32) error {
	return e.encoder.EncodeInt32(i)
}

func (e *Encoder) encodeInt16(i int16) error {
	return e.encoder.EncodeInt16(i)
}

func (e *Encoder) encodeByte(b byte) error {
	return e.encoder.EncodeUint8(b)
}

func (e *Encoder) encode(i any) error {
	return e.encoder.Encode(i)
}

func (e *Encoder) encodeBool(b bool) error {
	return e.encoder.EncodeBool(b)
}

func (e *Encoder) encodeSignature() error {
	// Write a signature that will e obviously corrupt
	// if the binary contents end up eing UTF-8 encoded
	// or altered by line feed or carriage return changes.
	return e.encodeBytes(binarySignature)
}

func (e *Encoder) encodeVersion() error {
	return e.encodeInt(e.version)
}

func (e *Encoder) encodeObjects(objects ecore.EList, check checkType) error {
	if err := e.encodeInt(objects.Size()); err != nil {
		return err
	}
	for it := objects.Iterator(); it.HasNext(); {
		eObject, _ := it.Next().(ecore.EObject)
		if err := e.encodeObject(eObject, check); err != nil {
			return err
		}
	}
	return nil
}

func (e *Encoder) encodeObject(eObject ecore.EObject, check checkType) error {
	if eObject == nil {
		return e.encodeInt(-1)
	} else if id, isID := e.objectToID[eObject]; isID {
		return e.encodeInt(id)
	} else if eObjectInternal, _ := eObject.(ecore.EObjectInternal); eObjectInternal != nil {
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
					(e.objectRoot != nil && !ecore.IsAncestor(e.objectRoot, eObjectInternal))) {
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

func (e *Encoder) encodeFeatureValue(eObject ecore.EObjectInternal, featureID int, featureData *binaryEncoderFeatureData) error {
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
		// nil for an unsettable reference set to nil : encodeObject writes it as -1
		eValue, _ := value.(ecore.EObject)
		switch featureData.featureKind {
		case bfkObject:
			fallthrough
		case bfkObjectContainment:
			return e.encodeObject(eValue, checkNothing)
		case bfkObjectContainerProxy:
			return e.encodeObject(eValue, checkResource)
		case bfkObjectContainmentProxy:
			return e.encodeObject(eValue, checkDirectResource)
		case bfkObjectProxy:
			return e.encodeObject(eValue, checkResource)
		case bfkObjectList:
			fallthrough
		case bfkObjectContainmentList:
			return e.encodeObjects(value.(ecore.EList), checkNothing)
		case bfkObjectContainmentListProxy:
			return e.encodeObjects(value.(ecore.EList), checkDirectResource)
		case bfkObjectListProxy:
			return e.encodeObjects(value.(ecore.EList), checkResource)
		case bfkData:
			valueStr := featureData.factory.ConvertToString(featureData.dataType, value)
			return e.encodeString(valueStr)
		case bfkDataList:
			l := value.(ecore.EList)
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

func (e *Encoder) encodeClass(eClass ecore.EClass) (*binaryEncoderClassData, error) {
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

func (e *Encoder) encodePackage(ePackage ecore.EPackage) (*binaryEncoderPackageData, error) {
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
		if err := e.encodeURI(ecore.GetURI(ePackage)); err != nil {
			return nil, err
		}
		e.packageDataMap[ePackage] = ePackageData
	}
	return ePackageData, nil
}

func (e *Encoder) encodeURI(uri *ecore.URI) error {
	if uri == nil {
		return e.encodeInt(-1)
	} else {
		return e.encodeURIWithFragment(uri.TrimFragment(), uri.Fragment())
	}
}

func (e *Encoder) encodeURIWithFragment(uri *ecore.URI, fragment string) error {
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

func (e *Encoder) relativizeURI(uri *ecore.URI) *ecore.URI {
	if e.baseURI != nil {
		return e.baseURI.Relativize(uri)
	}
	return uri
}

func (e *Encoder) newPackageData(ePackage ecore.EPackage) *binaryEncoderPackageData {
	return &binaryEncoderPackageData{
		id:        len(e.packageDataMap),
		classData: make([]*binaryEncoderClassData, ePackage.GetEClassifiers().Size()),
	}
}

func (e *Encoder) newClassID(ePackageData *binaryEncoderPackageData) int {
	for i, c := range ePackageData.classData {
		if c == nil {
			return i
		}
	}
	return -1
}

func (e *Encoder) newClassData(eClass ecore.EClass) (*binaryEncoderClassData, error) {
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
		eFeature := it.Next().(ecore.EStructuralFeature)
		eClassData.featureData = append(eClassData.featureData, e.newFeatureData(eFeature))
	}
	return eClassData, nil
}

const extendedMetaData = "http:///org/eclipse/emf/ecore/util/ExtendedMetaData"

func (e *Encoder) newFeatureData(eFeature ecore.EStructuralFeature) *binaryEncoderFeatureData {
	eFeatureData := &binaryEncoderFeatureData{
		name:        eFeature.GetName(),
		featureKind: getBinaryCodecFeatureKind(eFeature),
	}
	if eReference, _ := eFeature.(ecore.EReference); eReference != nil {
		eFeatureData.isTransient = eReference.IsTransient() || (eReference.IsContainer() && !eReference.IsResolveProxies())
	} else if eAttribute, _ := eFeature.(ecore.EAttribute); eAttribute != nil {
		eDataType := eAttribute.GetEAttributeType()
		eFeatureData.isTransient = eAttribute.IsTransient()
		eFeatureData.dataType = eDataType
		eFeatureData.factory = eDataType.GetEPackage().GetEFactoryInstance()
	}
	// if we have a xmlns prefix map or schema loction map, consider it as non transient to ensure
	// information is kept between binary encoder and xml encoder
	if e.isNamespaceAttributeEncoded {
		if eAnnotation := eFeature.GetEAnnotation(extendedMetaData); eAnnotation != nil {
			if name := eAnnotation.GetDetails().GetValue("name"); name == "xmlns:prefix" || name == "xsi:schemaLocation" {
				eFeatureData.isTransient = false
			}
		}
	}
	return eFeatureData
}
