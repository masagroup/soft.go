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
	"github.com/vmihailenco/msgpack/v5/msgpcode"
)

// the class and feature data are indexed by the ids of the encoder, which may not be the ones of
// the decoder model: a feature is found by its name
type binaryDecoderPackageData struct {
	ePackage   ecore.EPackage
	eFactory   ecore.EFactory
	eClassData map[int]*binaryDecoderClassData
}

type binaryDecoderClassData struct {
	eClass      ecore.EClass
	eFactory    ecore.EFactory
	featureData map[int]*binaryDecoderFeatureData
}

type binaryDecoderFeatureData struct {
	eFeature    ecore.EStructuralFeature
	eFactory    ecore.EFactory
	eDataType   ecore.EDataType
	featureID   int
	featureKind binaryFeatureKind
}

// maxPreallocation is the number of elements a list is allocated for before they are read when the
// length of the input is unknown: a longer list grows while its elements are read.
const maxPreallocation = 1024

// lenScanner is a reader that tells how many bytes it has left. The msgpack decoder reads directly
// from a reader that is also an io.ByteScanner (bytes.Reader, bytes.Buffer, strings.Reader) and
// wraps any other one in a bufio.Reader that reads ahead, so Len is the remaining input only for
// the former.
type lenScanner interface {
	io.Reader
	io.ByteScanner
	Len() int
}

type Decoder struct {
	resource     ecore.EResource
	r            io.Reader
	input        lenScanner // r when it tells its remaining length, else nil
	decoder      *msgpack.Decoder
	baseURI      *ecore.URI
	objects      []ecore.EObject
	uris         []*ecore.URI
	packageData  []*binaryDecoderPackageData
	enumLiterals []string
	depth        int
	maxDepth     int
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
		maxDepth:     DefaultMaxDepth,
	}
	d.input, _ = r.(lenScanner)
	if uri := resource.GetURI(); uri != nil {
		d.baseURI = uri
	}
	if maxDepth, isInt := options[OPTION_MAX_DEPTH].(int); isInt && maxDepth > 0 {
		d.maxDepth = maxDepth
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
	if size, err = d.decodeSize(); err != nil {
		return
	}
	objects := make([]any, 0, d.preallocation(size))
	for i := 0; i < size; i++ {
		var object ecore.EObject
		if object, err = d.decodeObject(); err != nil {
			return
		}
		objects = append(objects, object)
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

// decodeSize decodes the number of elements of a list. The size is read from the input and is not
// trusted: a negative one is refused, and so is one larger than the remaining input when the reader
// tells its length, since every element takes at least one byte.
func (d *Decoder) decodeSize() (int, error) {
	size, err := d.decodeInt()
	if err != nil {
		return 0, err
	}
	if size < 0 {
		return 0, fmt.Errorf("invalid size %d", size)
	}
	if d.input != nil && size > d.input.Len() {
		return 0, fmt.Errorf("invalid size %d: only %d bytes left", size, d.input.Len())
	}
	return size, nil
}

// preallocation returns the capacity to allocate for a list of size elements before reading them:
// size when decodeSize has checked it against the remaining input, at most maxPreallocation
// otherwise, the list growing while its elements are read.
func (d *Decoder) preallocation(size int) int {
	if d.input != nil {
		return size
	}
	return min(size, maxPreallocation)
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

// decodeInterface decodes an object id: a scalar (number, string, bytes) or an extension (uuid,
// ulid). An array or a map is refused before it is read: msgpack allocates them from their
// announced length and decodes their elements recursively, without bound.
func (d *Decoder) decodeInterface() (any, error) {
	c, err := d.decoder.PeekCode()
	if err != nil {
		return nil, err
	}
	if msgpcode.IsFixedArray(c) || msgpcode.IsFixedMap(c) ||
		c == msgpcode.Array16 || c == msgpcode.Array32 || c == msgpcode.Map16 || c == msgpcode.Map32 {
		return nil, fmt.Errorf("invalid object id: unexpected msgpack code %#x", c)
	}
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
	switch {
	case id == -1:
		return nil, nil
	case id < -1 || id > len(d.objects):
		// the objects are numbered in the order they are encoded
		return nil, fmt.Errorf("invalid object id %d", id)
	default:
		if len(d.objects) == id {
			// a new object is decoded inside the one that refers to it first
			if d.depth >= d.maxDepth {
				return nil, fmt.Errorf("invalid object %d: nested deeper than %d objects", id, d.maxDepth)
			}
			d.depth++
			defer func() { d.depth-- }()

			eClassData, err := d.decodeClass()
			if err != nil {
				return nil, err
			}
			eObject := eClassData.eFactory.Create(eClassData.eClass).(ecore.EObjectInternal)
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
				d.objects = append(d.objects, eObject)
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
				if featureID < 0 {
					return nil, fmt.Errorf("invalid feature id %d", featureID)
				}
				eFeatureData := eClassData.featureData[featureID]
				if eFeatureData == nil {
					eFeatureData, err = d.newFeatureData(eClassData)
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
			return eObject, nil

		} else {
			return d.objects[id], nil
		}
	}
}

func (d *Decoder) decodeObjects(list ecore.EList) error {
	size, err := d.decodeSize()
	if err != nil {
		return err
	}
	objects := make([]any, 0, d.preallocation(size))
	for i := 0; i < size; i++ {
		object, err := d.decodeObject()
		if err != nil {
			return err
		}
		objects = append(objects, object)
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
		value, err := createFromString(featureData.eFactory, featureData.eDataType, decoded)
		if err != nil {
			return err
		}
		eObject.ESetFromID(featureData.featureID, value)
	case bfkDataList:
		size, err := d.decodeSize()
		if err != nil {
			return err
		}
		values := make([]any, 0, d.preallocation(size))
		for range size {
			decoded, err := d.decodeString()
			if err != nil {
				return err
			}
			value, err := createFromString(featureData.eFactory, featureData.eDataType, decoded)
			if err != nil {
				return err
			}
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
		if id < 0 || id > len(d.enumLiterals) {
			return fmt.Errorf("invalid enum literal id %d", id)
		}
		if len(d.enumLiterals) == id {
			decoded, err := d.decodeString()
			if err != nil {
				return err
			}
			valueStr = decoded
			d.enumLiterals = append(d.enumLiterals, valueStr)
		} else {
			valueStr = d.enumLiterals[id]
		}
		value, err := createFromString(featureData.eFactory, featureData.eDataType, valueStr)
		if err != nil {
			return err
		}
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
	if id < 0 || id > len(d.packageData) {
		return nil, fmt.Errorf("invalid package id %d", id)
	}
	if len(d.packageData) == id {
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
	switch {
	case id == -1:
		return nil, nil
	case id < -1 || id > len(d.uris):
		return nil, fmt.Errorf("invalid uri id %d", id)
	default:
		var uri *ecore.URI
		if len(d.uris) == id {
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
		eClassData: map[int]*binaryDecoderClassData{},
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
		featureData: map[int]*binaryDecoderFeatureData{},
	}, nil
}

func (d *Decoder) newFeatureData(eClassData *binaryDecoderClassData) (*binaryDecoderFeatureData, error) {
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
		featureID:   eClassData.eClass.GetFeatureID(eFeature),
		featureKind: getBinaryCodecFeatureKind(eFeature),
	}
	if eAttribute, _ := eFeature.(ecore.EAttribute); eAttribute != nil {
		eFeatureData.eDataType = eAttribute.GetEAttributeType()
		eFeatureData.eFactory = eFeatureData.eDataType.GetEPackage().GetEFactoryInstance()
	}
	return eFeatureData, nil
}

// createFromString creates the value of eDataType from str: a factory panics on an invalid literal,
// which is returned as an error.
func createFromString(eFactory ecore.EFactory, eDataType ecore.EDataType, str string) (value any, err error) {
	defer func() {
		if r := recover(); r != nil {
			value, err = nil, fmt.Errorf("%v", r)
		}
	}()
	return eFactory.CreateFromString(eDataType, str), nil
}
