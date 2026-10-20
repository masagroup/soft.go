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
	"bytes"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/masagroup/soft.go/ecore"
	"github.com/vmihailenco/msgpack/v5"
)

type binaryDecoderPackageData struct {
	ePackage   ecore.EPackage
	eFactory   ecore.EFactory
	eClassData []*binaryDecoderClassData
}

type binaryDecoderClassData struct {
	eClass      ecore.EClass
	eFactory    ecore.EFactory
	featureData []*binaryDecoderFeatureData
}

type binaryDecoderFeatureData struct {
	eFeature    ecore.EStructuralFeature
	eFactory    ecore.EFactory
	eDataType   ecore.EDataType
	featureID   int
	featureKind binaryFeatureKind
}

type Decoder struct {
	resource         ecore.EResource
	r                io.Reader
	decoder          *msgpack.Decoder
	baseURI          *ecore.URI
	objects          []ecore.EObject
	uris             []*ecore.URI
	packageData      []*binaryDecoderPackageData
	enumLiterals     []string
	isResolveProxies bool
}

func NewDecoder(resource ecore.EResource, r io.Reader, options map[string]any) *Decoder {
	d := &Decoder{
		resource:     resource,
		r:            r,
		decoder:      msgpack.NewDecoder(r),
		objects:      []ecore.EObject{},
		uris:         []*ecore.URI{},
		packageData:  []*binaryDecoderPackageData{},
		enumLiterals: []string{},
	}
	if uri := resource.GetURI(); uri != nil {
		d.baseURI = uri
	}
	return d
}

func (d *Decoder) DecodeResource() {
	var err error
	defer func() {
		if err != nil {
			// add error to resource errors
			resourcePath := ""
			if d.resource.GetURI() != nil {
				resourcePath = d.resource.GetURI().String()
			}
			d.resource.GetErrors().Add(ecore.NewEDiagnosticImpl(err.Error(), resourcePath, 0, 0))
		}
	}()
	// signature
	if err = d.decodeSignature(); err != nil {
		return
	}
	// version
	if err = d.decodeVersion(); err != nil {
		return
	}
	// objects
	var size int
	if size, err = d.decodeInt(); err != nil {
		return
	}
	objects := make([]any, size)
	for i := 0; i < size; i++ {
		if objects[i], err = d.decodeObject(); err != nil {
			return
		}
	}

	// add objects to resource
	d.resource.GetContents().AddAll(ecore.NewImmutableEList(objects))
}

func (d *Decoder) DecodeObject() (ecore.EObject, error) {
	if err := d.decodeSignature(); err != nil {
		return nil, err
	}
	if err := d.decodeVersion(); err != nil {
		return nil, err
	}
	return d.decodeObject()
}

func (d *Decoder) decodeInt() (int, error) {
	return d.decoder.DecodeInt()
}

func (d *Decoder) decodeInt64() (int64, error) {
	return d.decoder.DecodeInt64()
}

func (d *Decoder) decodeInt32() (int32, error) {
	return d.decoder.DecodeInt32()

}

func (d *Decoder) decodeInt16() (int16, error) {
	return d.decoder.DecodeInt16()
}

func (d *Decoder) decodeByte() (byte, error) {
	return d.decoder.DecodeUint8()
}

func (d *Decoder) decodeBool() (bool, error) {
	return d.decoder.DecodeBool()
}

func (d *Decoder) decodeString() (string, error) {
	return d.decoder.DecodeString()
}

func (d *Decoder) decodeBytes() ([]byte, error) {
	return d.decoder.DecodeBytes()
}

func (d *Decoder) decodeDate() (*time.Time, error) {
	t, err := d.decoder.DecodeTime()
	if err != nil {
		return nil, err
	}
	t = t.UTC()
	return &t, nil
}

func (d *Decoder) decodeFloat64() (float64, error) {
	return d.decoder.DecodeFloat64()
}

func (d *Decoder) decodeFloat32() (float32, error) {
	return d.decoder.DecodeFloat32()
}

func (d *Decoder) decodeInterface() (any, error) {
	return d.decoder.DecodeInterface()
}

func (d *Decoder) decodeSignature() error {
	signature, err := d.decodeBytes()
	if err != nil {
		return err
	}
	if !bytes.Equal(signature, binarySignature) {
		return errors.New("invalid signature for a binary emf serialization")
	}
	return nil
}

func (d *Decoder) decodeVersion() error {
	version, err := d.decodeInt()
	if err != nil {
		return err
	}
	if version != binaryVersion {
		return errors.New("invalid version for binary emf serialization")
	}
	return nil
}

func (d *Decoder) decodeObject() (ecore.EObject, error) {
	id, err := d.decodeInt()
	if err != nil {
		return nil, err
	}
	switch id {
	case -1:
		return nil, nil
	default:
		if len(d.objects) <= int(id) {
			var eResult ecore.EObject
			eClassData, err := d.decodeClass()
			if err != nil {
				return nil, err
			}
			eObject := eClassData.eFactory.Create(eClassData.eClass).(ecore.EObjectInternal)
			eResult = eObject
			decodedInt, err := d.decodeInt()
			if err != nil {
				return nil, err
			}
			featureID := decodedInt - 1

			if featureID == -3 {
				// proxy object
				eProxyURI, err := d.decodeURI()
				if err != nil {
					return nil, err
				}
				eObject.ESetProxyURI(eProxyURI)
				if d.isResolveProxies {
					eResult = ecore.ResolveInResource(eObject, d.resource)
					d.objects = append(d.objects, eResult)
				} else {
					d.objects = append(d.objects, eObject)
				}
				decodedInt, err := d.decodeInt()
				if err != nil {
					return nil, err
				}
				featureID = decodedInt - 1
			} else {
				// standard object
				d.objects = append(d.objects, eObject)
			}

			if featureID == -2 {
				// object id attribute
				objectID, err := d.decodeInterface()
				if err != nil {
					return nil, err
				}
				if objectIDManager := d.resource.GetObjectIDManager(); objectIDManager != nil {
					if err := objectIDManager.SetID(eObject, objectID); err != nil {
						return nil, err
					}
				}
				decodedInt, err := d.decodeInt()
				if err != nil {
					return nil, err
				}
				featureID = decodedInt - 1
			}

			for featureID != -1 {
				eFeatureData := eClassData.featureData[featureID]
				if eFeatureData == nil {
					eFeatureData, err = d.newFeatureData(eClassData, featureID)
					if err != nil {
						return nil, err
					}
					eClassData.featureData[featureID] = eFeatureData
				}
				if err := d.decodeFeatureValue(eObject, eFeatureData); err != nil {
					return nil, err
				}

				decodedInt, err := d.decodeInt()
				if err != nil {
					return nil, err
				}
				featureID = decodedInt - 1
			}
			return eResult, nil

		} else {
			return d.objects[id], nil
		}
	}
}

func (d *Decoder) decodeObjects(list ecore.EList) error {
	size, err := d.decodeInt()
	if err != nil {
		return err
	}
	objects := make([]any, size)
	for i := 0; i < size; i++ {
		object, err := d.decodeObject()
		if err != nil {
			return err
		}
		objects[i] = object
	}

	// If the list is empty, we need to add all the objects,
	// otherwise, the reference is bidirectional and the list is at least partially populated.
	existingSize := list.Size()
	if existingSize == 0 {
		list.AddAll(ecore.NewImmutableEList(objects))
	} else {
		indices := make([]int, existingSize)
		duplicateCount := 0
		existingObjects := make([]any, existingSize)
		copy(existingObjects, list.ToArray())
	LOOP:
		for i := 0; i < size; i++ {
			o := objects[i]
			count := duplicateCount
			for j := range existingSize {
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
		list.AddAll(ecore.NewImmutableEList(objects))
		for i := range existingSize {
			newPosition := indices[i]
			oldPosition := size + i
			if newPosition != oldPosition {
				list.Move(oldPosition, newPosition)
			}
		}
	}
	return nil
}

func (d *Decoder) decodeFeatureValue(eObject ecore.EObjectInternal, featureData *binaryDecoderFeatureData) error {
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
		decoded, err := d.decodeObject()
		if err != nil {
			return err
		}
		eObject.ESetFromID(featureData.featureID, decoded)
	case bfkObjectList:
		fallthrough
	case bfkObjectListProxy:
		fallthrough
	case bfkObjectContainmentList:
		fallthrough
	case bfkObjectContainmentListProxy:
		l := eObject.EGetFromID(featureData.featureID, false).(ecore.EList)
		return d.decodeObjects(l)
	case bfkData:
		decoded, err := d.decodeString()
		if err != nil {
			return err
		}
		value := featureData.eFactory.CreateFromString(featureData.eDataType, decoded)
		eObject.ESetFromID(featureData.featureID, value)
	case bfkDataList:
		size, err := d.decodeInt()
		if err != nil {
			return err
		}
		values := []any{}
		for range size {
			decoded, err := d.decodeString()
			if err != nil {
				return err
			}
			value := featureData.eFactory.CreateFromString(featureData.eDataType, decoded)
			values = append(values, value)
		}
		l := eObject.EGetResolve(featureData.eFeature, false).(ecore.EList)
		l.AddAll(ecore.NewBasicEList(values))
	case bfkEnum:
		var valueStr string
		id, err := d.decodeInt()
		if err != nil {
			return err
		}
		if len(d.enumLiterals) <= id {
			decoded, err := d.decodeString()
			if err != nil {
				return err
			}
			valueStr = decoded
			d.enumLiterals = append(d.enumLiterals, valueStr)
		} else {
			valueStr = d.enumLiterals[id]
		}
		value := featureData.eFactory.CreateFromString(featureData.eDataType, valueStr)
		eObject.ESetFromID(featureData.featureID, value)
	case bfkDate:
		decoded, err := d.decodeDate()
		if err != nil {
			return err
		}
		eObject.ESetFromID(featureData.featureID, decoded)
	case bfkFloat64:
		decoded, err := d.decodeFloat64()
		if err != nil {
			return err
		}
		eObject.ESetFromID(featureData.featureID, decoded)
	case bfkFloat32:
		decoded, err := d.decodeFloat32()
		if err != nil {
			return err
		}
		eObject.ESetFromID(featureData.featureID, decoded)
	case bfkInt:
		decoded, err := d.decodeInt()
		if err != nil {
			return err
		}
		eObject.ESetFromID(featureData.featureID, decoded)
	case bfkInt64:
		decoded, err := d.decodeInt64()
		if err != nil {
			return err
		}
		eObject.ESetFromID(featureData.featureID, decoded)
	case bfkInt32:
		decoded, err := d.decodeInt32()
		if err != nil {
			return err
		}
		eObject.ESetFromID(featureData.featureID, decoded)
	case bfkInt16:
		decoded, err := d.decodeInt16()
		if err != nil {
			return err
		}
		eObject.ESetFromID(featureData.featureID, decoded)
	case bfkByte:
		decoded, err := d.decodeByte()
		if err != nil {
			return err
		}
		eObject.ESetFromID(featureData.featureID, decoded)
	case bfkBool:
		decoded, err := d.decodeBool()
		if err != nil {
			return err
		}
		eObject.ESetFromID(featureData.featureID, decoded)
	case bfkString:
		decoded, err := d.decodeString()
		if err != nil {
			return err
		}
		eObject.ESetFromID(featureData.featureID, decoded)
	case bfkByteArray:
		decoded, err := d.decodeBytes()
		if err != nil {
			return err
		}
		eObject.ESetFromID(featureData.featureID, decoded)
	}
	return nil
}

func (d *Decoder) decodeClass() (*binaryDecoderClassData, error) {
	ePackageData, err := d.decodePackage()
	if err != nil {
		return nil, err
	}
	id, err := d.decodeInt()
	if err != nil {
		return nil, err
	}
	eClassData := ePackageData.eClassData[id]
	if eClassData == nil {
		eClassData, err = d.newClassData(ePackageData)
		if err != nil {
			return nil, err
		}
		ePackageData.eClassData[id] = eClassData
	}
	return eClassData, nil
}

func (d *Decoder) decodePackage() (*binaryDecoderPackageData, error) {
	id, err := d.decodeInt()
	if err != nil {
		return nil, err
	}
	if len(d.packageData) <= id {
		// decode package parameters
		nsURI, err := d.decodeString()
		if err != nil {
			return nil, err
		}
		uri, err := d.decodeURI()
		if err != nil {
			return nil, err
		}

		// retrieve package
		packageRegistry := ecore.GetPackageRegistry()
		resourceSet := d.resource.GetResourceSet()
		if resourceSet != nil {
			packageRegistry = resourceSet.GetPackageRegistry()
		}
		ePackage := packageRegistry.GetPackage(nsURI)
		if ePackage == nil && resourceSet != nil {
			ePackage, _ = resourceSet.GetEObject(uri, true).(ecore.EPackage)
		}
		if ePackage == nil {
			return nil, fmt.Errorf("unable to find package '%s'", nsURI)
		}
		// create new package data
		ePackageData := d.newPackageData(ePackage)
		ePackageData.eFactory = ecore.GetRegistryFactory(packageRegistry, ePackage)
		d.packageData = append(d.packageData, ePackageData)
		return ePackageData, nil

	} else {
		return d.packageData[id], nil
	}
}

func (d *Decoder) decodeURI() (*ecore.URI, error) {
	id, err := d.decodeInt()
	if err != nil {
		return nil, err
	}
	switch id {
	case -1:
		return nil, nil
	default:
		var uri *ecore.URI
		if len(d.uris) <= int(id) {
			// build uri
			uriStr, err := d.decodeString()
			if err != nil {
				return nil, err
			}
			if uriStr == "" {
				uri = d.baseURI
			} else {
				uri = d.resolveURI(ecore.NewURI(uriStr))
			}
			// add it to the uri array
			d.uris = append(d.uris, uri)
		} else {
			uri = d.uris[id]
		}
		fragment, err := d.decodeString()
		if err != nil {
			return nil, err
		}
		return ecore.NewURIBuilder(uri).SetFragment(fragment).URI(), nil
	}
}

func (d *Decoder) resolveURI(uri *ecore.URI) *ecore.URI {
	if d.baseURI != nil {
		return d.baseURI.Resolve(uri)
	}
	return uri
}

func (p *binaryDecoderPackageData) getFactory() ecore.EFactory {
	if p.eFactory != nil {
		return p.eFactory
	}
	return p.ePackage.GetEFactoryInstance()
}

func (d *Decoder) newPackageData(ePackage ecore.EPackage) *binaryDecoderPackageData {
	return &binaryDecoderPackageData{
		ePackage:   ePackage,
		eClassData: make([]*binaryDecoderClassData, ePackage.GetEClassifiers().Size()),
	}
}

func (d *Decoder) newClassData(ePackageData *binaryDecoderPackageData) (*binaryDecoderClassData, error) {
	className, err := d.decodeString()
	if err != nil {
		return nil, err
	}
	ePackage := ePackageData.ePackage
	eClass, _ := ePackage.GetEClassifier(className).(ecore.EClass)
	if eClass == nil {
		return nil, fmt.Errorf("unable to find class '%v' in package '%v'", className, ePackage.GetNsURI())
	}
	return &binaryDecoderClassData{
		eClass:      eClass,
		eFactory:    ePackageData.getFactory(),
		featureData: make([]*binaryDecoderFeatureData, eClass.GetFeatureCount()),
	}, nil
}

func (d *Decoder) newFeatureData(eClassData *binaryDecoderClassData, featureID int) (*binaryDecoderFeatureData, error) {
	eFeatureName, err := d.decodeString()
	if err != nil {
		return nil, err
	}
	eFeature := eClassData.eClass.GetEStructuralFeatureFromName(eFeatureName)
	if eFeature == nil {
		return nil, fmt.Errorf("unable to find feature '%v' in '%v' EClass", eFeatureName, eClassData.eClass.GetName())
	}
	eFeatureData := &binaryDecoderFeatureData{
		eFeature:    eFeature,
		featureID:   featureID,
		featureKind: getBinaryCodecFeatureKind(eFeature),
	}
	if eAttribute, _ := eFeature.(ecore.EAttribute); eAttribute != nil {
		eFeatureData.eDataType = eAttribute.GetEAttributeType()
		eFeatureData.eFactory = eFeatureData.eDataType.GetEPackage().GetEFactoryInstance()
	}
	return eFeatureData, nil
}
