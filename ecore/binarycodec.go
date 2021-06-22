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

const (
	object_container              = iota
	object_container_proxy        = iota
	object                        = iota
	object_proxy                  = iota
	object_list                   = iota
	object_list_proxy             = iota
	object_containment            = iota
	object_containment_proxy      = iota
	object_containment_list       = iota
	object_containment_list_proxy = iota
	data                          = iota
	data_list                     = iota
	enum                          = iota
	date                          = iota
	primitive                     = iota
)

func getBinaryCodecFeatureKind(eFeature EStructuralFeature) int {
	if eReference, _ := eFeature.(EReference); eReference != nil {
		if eReference.IsContainment() {
			if eReference.IsResolveProxies() {
				if eReference.IsMany() {
					return object_containment_list_proxy
				} else {
					return object_containment_proxy
				}
			} else {
				if eReference.IsMany() {
					return object_containment_list
				} else {
					return object_containment
				}
			}
		} else if eReference.IsContainer() {
			if eReference.IsResolveProxies() {
				return object_container_proxy
			} else {
				return object_container
			}
		} else if eReference.IsResolveProxies() {
			if eReference.IsMany() {
				return object_list_proxy
			} else {
				return object_proxy
			}
		} else {
			if eReference.IsMany() {
				return object_list
			} else {
				return object
			}
		}
	} else if eAttribute, _ := eFeature.(EAttribute); eAttribute != nil {
		if eAttribute.IsMany() {
			return data_list
		} else {
			eDataType := eAttribute.GetEAttributeType()
			if eEnum, _ := eDataType.(EEnum); eEnum != nil {
				return enum
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
				return primitive
			}
			if instanceTypeName == "*time.Time" {
				return date
			}
			return data
		}
	}
	return -1
}
