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
)

type BinaryCodec struct {
}

func (bc *BinaryCodec) NewEncoder(resource EResource, w io.Writer, options map[string]interface{}) EResourceEncoder {
	return NewBinaryEncoder(resource, w, options)
}
func (bc *BinaryCodec) NewDecoder(resource EResource, r io.Reader, options map[string]interface{}) EResourceDecoder {
	return NewBinaryDecoder(resource, r, options)
}

type binaryFeatureKind int

const (
	bfkObjectContainer binaryFeatureKind = iota
	bfkObjectContainerProxy
	bfkObject
	bfkObjectProxy
	bfkObjectList
	bfkObjectListProxy
	bfkObjectContainment
	bfkObjectContainmentProxy
	bfkObjectContainmentList
	bfkObjectContainmentListProxy
	bfkData
	bfkDataList
	bfkEnum
	bfkDate
	bfkPrimitive
)

func getBinaryCodecFeatureKind(eFeature EStructuralFeature) binaryFeatureKind {
	if eReference, _ := eFeature.(EReference); eReference != nil {
		if eReference.IsContainment() {
			if eReference.IsResolveProxies() {
				if eReference.IsMany() {
					return bfkObjectContainmentListProxy
				} else {
					return bfkObjectContainmentProxy
				}
			} else {
				if eReference.IsMany() {
					return bfkObjectContainmentList
				} else {
					return bfkObjectContainment
				}
			}
		} else if eReference.IsContainer() {
			if eReference.IsResolveProxies() {
				return bfkObjectContainerProxy
			} else {
				return bfkObjectContainer
			}
		} else if eReference.IsResolveProxies() {
			if eReference.IsMany() {
				return bfkObjectListProxy
			} else {
				return bfkObjectProxy
			}
		} else {
			if eReference.IsMany() {
				return bfkObjectList
			} else {
				return bfkObject
			}
		}
	} else if eAttribute, _ := eFeature.(EAttribute); eAttribute != nil {
		if eAttribute.IsMany() {
			return bfkDataList
		} else {
			eDataType := eAttribute.GetEAttributeType()
			if eEnum, _ := eDataType.(EEnum); eEnum != nil {
				return bfkEnum
			}
			instanceTypeName := eDataType.GetInstanceTypeName()
			if instanceTypeName == "float64" ||
				instanceTypeName == "float32" ||
				instanceTypeName == "int" ||
				instanceTypeName == "int64" ||
				instanceTypeName == "int32" ||
				instanceTypeName == "int16" ||
				instanceTypeName == "bool" ||
				instanceTypeName == "string" {
				return bfkPrimitive
			}
			if instanceTypeName == "*time.Time" {
				return bfkDate
			}
			return bfkData
		}
	}
	return -1
}
